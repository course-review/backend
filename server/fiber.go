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
	dbURL := "postgres://postgres:mysecretpassword@localhost:5432/postgres"
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

	app.Get("/getStats", func(c *fiber.Ctx) error {
		stats, err := db.GetStats(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(stats)
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
		user, err := DecodeJWT(c.Query("token"))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		if user.Exp < time.Now().Unix() {
			return c.Status(401).JSON(fiber.Map{"error": "Token expired"})
		}
		c.Locals("unique_id", user.UniqueID)
		log.Println(user.UniqueID)
		return c.Next()
	})

	auth.Get("/getUserData", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		log.Println(uniqueId)
		data, err := db.GetUserData(c.Context(), uniqueId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	auth.Post("/insertReview", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		//todo check here if mapping already exists

		id, err := db.GetCourseEvaluationMap(c.Context(), sql.GetCourseEvaluationMapParams{})
		if err != nil {
			id, err = db.SetCourseEvaluationMap(c.Context(), sql.SetCourseEvaluationMapParams{UserID: uniqueId, CourseNumber: c.Query("courseNumber"), Semester: pgtype.Text{String: c.Query("semester"), Valid: true}})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			_, err = db.UpdateSemester(c.Context(), sql.UpdateSemesterParams{EvaluationID: id, Semester: pgtype.Text{String: c.Query("semester"), Valid: true}})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		reviewText := strings.TrimSpace(c.Query("review"))

		if reviewText != "" {
			_, err := db.SetReview(c.Context(), sql.SetReviewParams{EvaluationID: id, Review: c.Query("review")})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		if c.Query("rating") != "" {
			_, err := db.SetRating(c.Context(), sql.SetRatingParams{
				EvaluationID: id,
				Recommended:  pgtype.Int4{Int32: int32(c.QueryInt("Recommended"))},
				Engaging:     pgtype.Int4{Int32: int32(c.QueryInt("Engaging"))},
				Difficulty:   pgtype.Int4{Int32: int32(c.QueryInt("Difficulty"))},
				Effort:       pgtype.Int4{Int32: int32(c.QueryInt("Effort"))},
				Resources:    pgtype.Int4{Int32: int32(c.QueryInt("Resources"))},
			})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		return c.JSON(fiber.Map{"success": "Review added"})
	})

	auth.Post("/updateReview", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		id := int32(c.QueryInt("id"))
		_, err := db.CheckUserWithId(c.Context(), sql.CheckUserWithIdParams{EvaluationID: id, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		review, err := db.UpdateReview(c.Context(), sql.UpdateReviewParams{Review: c.Query("review"), EvaluationID: id, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	auth.Post("/deleteReview", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)
		id := int32(c.QueryInt("id"))
		_, err := db.CheckUserWithId(c.Context(), sql.CheckUserWithIdParams{EvaluationID: id, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		review, err := db.DeleteReview(c.Context(), id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	// auth.Post("/insertRating", func(c *fiber.Ctx) error {
	// 	//todo check if mapId already exists

	// 	//if not create it

	// 	data := new(sql.SetRatingParams)
	// 	if err := c.BodyParser(data); err != nil {
	// 		return err
	// 	}
	// 	rating, err := db.SetRating(c.Context(), *data)
	// 	if err != nil {
	// 		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// 	return c.JSON(rating)
	// })

	auth.Post("/updateRating", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)

		_, err := db.CheckUserWithId(c.Context(), sql.CheckUserWithIdParams{EvaluationID: int32(c.QueryInt("id")), UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		data := new(sql.UpdateRatingParams)
		if err := c.BodyParser(data); err != nil {
			return err
		}
		rating, err := db.UpdateRating(c.Context(), *data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(rating)
	})

	auth.Post("/updateSemester", func(c *fiber.Ctx) error {
		uniqueId, _ := c.Locals("unique_id").(string)

		semester, err := db.UpdateSemester(c.Context(), sql.UpdateSemesterParams{EvaluationID: int32(c.QueryInt("id")), Semester: pgtype.Text{String: c.Query("semester"), Valid: true}, UserID: uniqueId})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(semester)
	})

	// // // // // // // //
	// mod / admin needed //
	// // // // // // // //
	moderator := app.Group("/moderator")

	moderator.Use("/", func(c *fiber.Ctx) error {
		user, err := DecodeJWT(c.Query("token"))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		mod, err := db.GetUser(c.Context(), user.UniqueID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if mod.Moderator.Bool == false {
			return c.Status(401).JSON(fiber.Map{"error": "Not authorized"})
		}
		return c.Next()
	})

	admin := app.Group("/admin")
	admin.Use("/", func(c *fiber.Ctx) error {
		user, err := DecodeJWT(c.Query("token"))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		mod, err := db.GetUser(c.Context(), user.UniqueID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if mod.Admin.Bool == false {
			return c.Status(401).JSON(fiber.Map{"error": "Not authorized"})
		}
		return c.Next()
	})

	moderator.Post("/setCurrentSemester", func(c *fiber.Ctx) error {
		db.RemoveCurrentSemester(c.Context())
		for _, semester := range strings.Split(c.Query("list"), ",") {
			_, err := db.SetCurrentSemester(c.Context(), string(semester))
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}
		return c.JSON(fiber.Map{"success": "Semester set"})
	})

	admin.Post("/setModerator", func(c *fiber.Ctx) error {
		val, err := db.SetModerator(c.Context(), c.Query("user"))
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
		review, err := db.VerifyReview(c.Context(), int32(c.QueryInt("id")))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	moderator.Post("/rejectReview", func(c *fiber.Ctx) error {
		review, err := db.RejectReview(c.Context(), int32(c.QueryInt("id")))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	app.Post("/setUser", func(c *fiber.Ctx) error {
		data := new(sql.SetUserParams)
		if err := c.BodyParser(data); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		user, err := db.SetUser(c.Context(), *data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(user)
	})

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
