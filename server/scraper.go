package main

import (
	"bytes"
	"context"
	"coursereview/app/generated/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5"
)

var webhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
var regexRuleCourseNumber = regexp.MustCompile("<b>(\\d{3}-\\d{4}-[A-Z0-9]{3})<\\/b>")
var regexRuleCourseName = regexp.MustCompile("\\w{2}\">(.*?)<\\/a><\\/b>")

const (
	username     = "VVZ Scrape-Inator 6000"
	avatarURL    = "https://cdn.discordapp.com/attachments/819966095070330950/1076174116165537832/vvz.png"
	language     = "de"
	resultFolder = "vvzScrapeResults/"
)

func sendDiscordMessage(title, description string, color int) {
	payload := map[string]interface{}{
		"username":   username,
		"avatar_url": avatarURL,
		"embeds": []map[string]interface{}{
			{
				"title":       title,
				"description": description,
				"color":       color,
			},
		},
	}
	data, _ := json.Marshal(payload)
	http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
}

func sendScrapingStart(semester string) {
	sendDiscordMessage("Scraping new courses of "+semester, "", 1651554)
}

func sendScrapingEnd(newCourses int, semester string) {
	title := fmt.Sprintf("Finished scraping %d new courses of %s", newCourses, semester)
	sendDiscordMessage(title, "", 1651554)

	if newCourses > 0 {
		filePath := resultFolder + semester + ".txt"
		fileData, err := os.ReadFile(filePath)
		if err == nil {
			body := &bytes.Buffer{}
			writer := io.MultiWriter(body)

			req, err := http.NewRequest("POST", webhookURL, body)
			if err != nil {
				return
			}
			writer.Write(fileData)

			req.Header.Set("Content-Type", "multipart/form-data")
			http.DefaultClient.Do(req)
		}
	}
}

func sendScrapingError(err string) {
	description := fmt.Sprintf("%s caused an issue", err)
	sendDiscordMessage("uh-oh", description, 6428441)
}

func vvzScraper(semester string, mainContext context.Context) {
	scrapeContext, cancel := context.WithCancel(mainContext)
	defer cancel()

	// connect to db
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("DB_URL environment variable not set")
		return
	}
	pool, err := pgx.Connect(scrapeContext, dbURL)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	db := sql.New(pool)
	routine(semester, scrapeContext, db)
	pool.Close(scrapeContext)
	cancel()
}

func vvzListUrl(semester string, language string) string {
	return fmt.Sprintf("https://www.vvz.ethz.ch/Vorlesungsverzeichnis/sucheLehrangebot.view?seite=0&semkez=%s&lang=%s", semester, language)
}

func routine(semester string, context context.Context, db *sql.Queries) {
	sendScrapingStart(semester)

	collector := colly.NewCollector(
		colly.AllowedDomains("www.vvz.ethz.ch"),
	)
	collector.SetRequestTimeout(120 * time.Second)

	var courseNumbers []string
	var courseNames []string

	collector.OnHTML("tr", func(e *colly.HTMLElement) {
		tableRow, _ := e.DOM.Html()
		matchesCourseNumber := regexRuleCourseNumber.FindAllStringSubmatch(tableRow, -1)
		matchesCourseName := regexRuleCourseName.FindAllStringSubmatch(tableRow, -1)

		for _, match := range matchesCourseNumber {
			if len(match) > 1 {
				courseNumbers = append(courseNumbers, match[1])
			}
		}

		for _, match := range matchesCourseName {
			if len(match) > 1 {
				courseNames = append(courseNames, match[1])
			}
		}
	})

	err := collector.Visit(vvzListUrl(semester, language))
	if err != nil {
		fmt.Println("Error visiting URL:", err)
		sendScrapingError("Error visiting URL: " + err.Error())
		return
	}

	if len(courseNumbers) != len(courseNames) {
		sendScrapingError("Course numbers and names are not equal in length")
		return
	}
	newCourses := 0
	for i, item := range courseNumbers {
		// check db
		_, err := db.GetCourseName(context, item)
		if err == nil {
			continue
		}
		// add course to db
		_, err = db.AddCourse(context, sql.AddCourseParams{CourseNumber: item, CourseName: courseNames[i]})
		if err != nil {
			fmt.Println("Error adding course to DB:", err)
			sendScrapingError(item)
			continue
		}
		newCourses++
	}

	sendScrapingEnd(newCourses, semester)
}
