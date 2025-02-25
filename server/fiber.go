package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"coursereview/app/generated/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

//todo check if user is allowed to do the action

func main() {

	app := fiber.New()
	app.Use(logger.New())

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

	app.Get("/getUserData", func(c *fiber.Ctx) error {
		data, err := db.GetUserData(c.Context(), c.Query("user"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	app.Post("/insertReview", func(c *fiber.Ctx) error {
		data := new(sql.SetReviewParams)
		if err := c.BodyParser(data); err != nil {
			return err
		}

		//todo check here if mapping already exists

		review, err := db.SetReview(c.Context(), *data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	app.Post("/updateReview", func(c *fiber.Ctx) error {
		data := new(sql.UpdateReviewParams)
		if err := c.BodyParser(data); err != nil {
			return err
		}
		review, err := db.UpdateReview(c.Context(), *data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	app.Post("/deleteReview", func(c *fiber.Ctx) error {
		data, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}
		review, err := db.DeleteReview(c.Context(), int32(data))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(review)
	})

	app.Post("/insertRating", func(c *fiber.Ctx) error {
		data := new(sql.SetRatingParams)
		if err := c.BodyParser(data); err != nil {
			return err
		}
		rating, err := db.SetRating(c.Context(), *data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(rating)
	})

	app.Post("/updateRating", func(c *fiber.Ctx) error {
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

	app.Post("/updateSemester", func(c *fiber.Ctx) error {
		data := new(sql.UpdateSemesterParams)
		if err := c.BodyParser(data); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		semester, err := db.UpdateSemester(c.Context(), *data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(semester)
	})

	// // // // // // // //
	// mod / admin needed //
	// // // // // // // //

	app.Post("/setCurrentSemester", func(c *fiber.Ctx) error {
		db.RemoveCurrentSemester(c.Context())

		for _, semester := range c.Query("semester") {
			_, err := db.SetCurrentSemester(c.Context(), string(semester))
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}
		return c.JSON(fiber.Map{"success": "Semester set"})
	})

	app.Post("/setModerator", func(c *fiber.Ctx) error {
		val, err := db.SetModerator(c.Context(), c.String())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": val})
	})

	app.Get("/getUnverifiedReviews", func(c *fiber.Ctx) error {
		reviews, err := db.GetUnverifiedReviews(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(reviews)
	})

	app.Post("/verifyReview", func(c *fiber.Ctx) error {
		data, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}
		review, err := db.VerifyReview(c.Context(), int32(data))
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

	app.Post("/addCourse", func(c *fiber.Ctx) error {
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
