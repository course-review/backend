package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"coursereview/app/generated/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectDB() (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

type TokenProperties struct {
	Student  bool   `json:"student"`
	Exp      int64  `json:"exp"`
	UniqueID string `json:"unique_id"`
}

type Ratings struct { //can I parse this directly to pgtype.Int4?
	Recommended int32 `json:"recommended"`
	Engaging    int32 `json:"engaging"`
	Difficulty  int32 `json:"difficulty"`
	Effort      int32 `json:"effort"`
	Resources   int32 `json:"resources"`
}

// DecodeJWT decodes a JWT payload without verifying the signature
func DecodeJWT(token string) (*TokenProperties, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid JWT token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims TokenProperties
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	if claims.Exp < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

//todo check if user is allowed to do the action

func main() {
	RunMigration()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	pool, err := connectDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	db := sql.New(pool)

	// Testing endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Playing with": "Duckies"})
	})

	app.Get("/all", func(c *fiber.Ctx) error {
		data, err := db.GetAllTheData(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	app.Get("/stats", func(c *fiber.Ctx) error {
		stats, err := db.GetStats(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(stats)
	})

	app.Get("/latestReviews", func(c *fiber.Ctx) error {
		reviews, err := db.GetReviewedCourses(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(reviews)
	})

	app.Get("/getReviews", func(c *fiber.Ctx) error {
		reviews, err := db.GetReviews(c.Context(), c.Query("course"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(reviews)
	})

	app.Get("/getRatings", func(c *fiber.Ctx) error {
		ratings, err := db.GetRatings(c.Context(), c.Query("course"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(ratings)
	})

	app.Get("/courses", func(c *fiber.Ctx) error {
		data, err := db.GetCourses(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	app.Get("/searchCourses", func(c *fiber.Ctx) error {
		return c.Status(500).JSON(fiber.Map{"error": "Not implemented"})
	})

	app.Get("/currentSemesters", func(c *fiber.Ctx) error {
		semester, err := db.GetCurrentSemester(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(semester)
	})

	app.Get("/courseName", func(c *fiber.Ctx) error {
		data, err := db.GetCourseName(c.Context(), c.Query("course"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	// // // // // // // // //
	// authentication needed //
	// // // // // // // // //
	// todo: could make an id behind auth group that checks if evalid belongs to user instead of duplicated code
	auth := app.Group("/auth")

	auth.Use("/", func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			type Token struct {
				Token string `json:"token"`
			}
			var data Token
			if err := c.BodyParser(&data); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
			}
			token = data.Token
		}
		user, err := DecodeJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		if user.Exp < time.Now().Unix() {
			return c.Status(401).JSON(fiber.Map{"error": "Token expired"})
		}
		c.Locals("unique_id", user.UniqueID)
		return c.Next()
	})

	auth.Get("/getUserData", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		data, err := db.GetUserData(c.Context(), uniqueId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	auth.Post("/insertReview", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		type payload struct {
			CourseNumber string `json:"courseNumber"`
			Semester     string `json:"semester"`
			Review       string `json:"review"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		id, err := db.GetCourseEvaluationMap(c.Context(), sql.GetCourseEvaluationMapParams{UserID: uniqueId, CourseNumber: data.CourseNumber})
		if err != nil {
			//check here if there is data to set
			id, err = db.SetCourseEvaluationMap(c.Context(), sql.SetCourseEvaluationMapParams{UserID: uniqueId, CourseNumber: data.CourseNumber, Semester: pgtype.Text{String: data.Semester, Valid: true}})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			_, err = db.UpdateSemester(c.Context(), sql.UpdateSemesterParams{EvaluationID: id, Semester: pgtype.Text{String: data.Semester, Valid: true}})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		err = reviewChange(c, db, id, data.Review)
		if err != nil {
			return err
		}
		return ratingChange(c, db, id)
	})

	auth.Post("/updateReview", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		type payload struct {
			Id     int32  `json:"id"`
			Review string `json:"review"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		_, err := db.CheckUserWithId(c.Context(), sql.CheckUserWithIdParams{EvaluationID: data.Id, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return reviewChange(c, db, data.Id, data.Review)
	})

	auth.Post("/deleteRating", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		type payload struct {
			Id int32 `json:"id"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		_, err := db.CheckUserWithId(c.Context(), sql.CheckUserWithIdParams{EvaluationID: data.Id, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		review, err := db.DeleteRating(c.Context(), data.Id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = db.CheckRatingAndReview(c.Context(), data.Id)
		if err != nil {
			_, err = db.DeleteCourseEvaluationMap(c.Context(), data.Id)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		return c.JSON(review)
	})

	auth.Post("/deleteReview", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)

		type payload struct {
			Id int32 `json:"id"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		_, err := db.CheckUserWithId(c.Context(), sql.CheckUserWithIdParams{EvaluationID: data.Id, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		review, err := db.DeleteReview(c.Context(), data.Id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = db.CheckRatingAndReview(c.Context(), data.Id)

		if err != nil {
			_, err = db.DeleteCourseEvaluationMap(c.Context(), data.Id)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		return c.JSON(review)
	})

	auth.Post("/updateRating", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)

		type payload struct {
			Id int32 `json:"id"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		_, err := db.CheckUserWithId(c.Context(), sql.CheckUserWithIdParams{EvaluationID: data.Id, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return ratingChange(c, db, data.Id)
	})

	auth.Post("/updateSemester", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		type payload struct { //can I parse this directly to pgtype.Text?
			Id       int32  `json:"id"`
			Semester string `json:"semester"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		semester, err := db.UpdateSemester(c.Context(), sql.UpdateSemesterParams{EvaluationID: data.Id, Semester: pgtype.Text{String: data.Semester, Valid: true}, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(semester)
	})

	// // // // // // // //
	// mod / admin needed //
	// // // // // // // //
	moderator := auth.Group("/moderator")

	moderator.Use("/", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		user, err := db.GetUser(c.Context(), uniqueId)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
		}
		if !user.Moderator.Bool && !user.Admin.Bool {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden"})
		}
		return c.Next()
	})

	admin := auth.Group("/admin")
	admin.Use("/", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		user, err := db.GetUser(c.Context(), uniqueId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if !user.Admin.Bool {
			return c.Status(401).JSON(fiber.Map{"error": "Not authorized"})
		}
		return c.Next()
	})

	moderator.Post("/setCurrentSemester", func(c *fiber.Ctx) error {
		type payload struct {
			List []string `json:"list"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		db.RemoveCurrentSemester(c.Context())
		for _, semester := range data.List {
			_, err := db.SetCurrentSemester(c.Context(), string(semester))
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}
		return c.JSON(fiber.Map{"success": "Semester set"})
	})

	admin.Post("/setModerator", func(c *fiber.Ctx) error {
		type payload struct {
			User string `json:"user"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		val, err := db.SetModerator(c.Context(), data.User)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": val})
	})

	moderator.Get("/getUnverifiedReviews", func(c *fiber.Ctx) error {
		reviews, err := db.GetUnverifiedReviews(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(reviews)
	})

	moderator.Post("/verifyReview", func(c *fiber.Ctx) error {
		type payload struct {
			Id int32 `json:"id"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		review, err := db.VerifyReview(c.Context(), data.Id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	moderator.Post("/rejectReview", func(c *fiber.Ctx) error {
		type payload struct {
			Id               int32  `json:"id"`
			RequestedChanges string `json:"requested_changes"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		review, err := db.RejectReview(c.Context(), sql.RejectReviewParams{EvaluationID: data.Id, RequestedChanges: pgtype.Text{String: data.RequestedChanges, Valid: true}})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	app.Post("/setUser", func(c *fiber.Ctx) error {
		type payload struct {
			User string `json:"user"`
		}
		var data payload
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		user, err := db.SetUser(c.Context(), data.User)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	})

	//todo: not used yet
	admin.Post("/addCourse", func(c *fiber.Ctx) error {
		data := new(sql.SetCourseParams)
		if err := c.BodyParser(data); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		course, err := db.SetCourse(c.Context(), *data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(course)
	})

	log.Fatal(app.Listen(":3000"))
}

func reviewChange(c *fiber.Ctx, db *sql.Queries, evalId int32, review string) error {
	review = strings.TrimSpace(review)
	if review == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Review cannot be empty"})
	}

	//check if eval id exists in review table
	//if it doesn't, insert
	//if it does, update

	_, err := db.GetReviewWithId(c.Context(), evalId)
	if err != nil {
		_, err = db.SetReview(c.Context(), sql.SetReviewParams{EvaluationID: evalId, Review: review})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		_, err = db.UpdateReview(c.Context(), sql.UpdateReviewParams{EvaluationID: evalId, Review: review})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.JSON(fiber.Map{"success": "Set review"})
}

func ratingChange(c *fiber.Ctx, db *sql.Queries, evalId int32) error {
	//check if eval id exists in rating table
	//if it doesn't, insert
	//if it does, update
	var newRating Ratings
	if err := c.BodyParser(&newRating); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	ratings := sql.SetRatingParams{
		EvaluationID: evalId,
		Recommended:  pgtype.Int4{Int32: newRating.Recommended, Valid: newRating.Recommended != 0},
		Engaging:     pgtype.Int4{Int32: newRating.Engaging, Valid: newRating.Engaging != 0},
		Difficulty:   pgtype.Int4{Int32: newRating.Difficulty, Valid: newRating.Difficulty != 0},
		Effort:       pgtype.Int4{Int32: newRating.Effort, Valid: newRating.Effort != 0},
		Resources:    pgtype.Int4{Int32: newRating.Resources, Valid: newRating.Resources != 0},
	}

	_, err := db.GetRatingWithId(c.Context(), evalId)
	if err != nil {
		if newRating.Recommended+newRating.Engaging+newRating.Difficulty+newRating.Effort+newRating.Resources == 0 {
			return c.Status(200).JSON(fiber.Map{"error": "Ratings not set"})
		}
		_, err = db.SetRating(c.Context(), ratings)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": "Set rating"})
	} else {
		if newRating.Recommended+newRating.Engaging+newRating.Difficulty+newRating.Effort+newRating.Resources == 0 {
			return c.Status(500).JSON(fiber.Map{"error": "Ratings cannot be empty"})
		}
		updateRatings := sql.UpdateRatingParams(ratings)
		_, err := db.UpdateRating(c.Context(), updateRatings)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": "Updated rating"})
	}
}
