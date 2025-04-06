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

	"github.com/jackc/pgx/v5"
)

var webhookURL = os.Getenv("DISCORD_WEBHOOK_URL")

const (
	username     = "VVZ Scrape-Inator 6000"
	avatarURL    = "https://cdn.discordapp.com/attachments/819966095070330950/1076174116165537832/vvz.png"
	semester     = "2025S"
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

func sendScrapingStart() {
	sendDiscordMessage("Scraping new courses of "+semester, "", 1651554)
}

func sendScrapingEnd(newCourses int) {
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

func sendScrapingError(courseNr, url string) {
	description := fmt.Sprintf("[%s](%s) caused an issue", courseNr, url)
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

func routine(semester string, context context.Context, db *sql.Queries) {

	_ = os.MkdirAll(resultFolder, os.ModePerm)
	resultFilePath := resultFolder + semester + ".txt"
	_ = os.WriteFile(resultFilePath, []byte{}, 0644)

	sendScrapingStart()

	vvzUrl := fmt.Sprintf("https://www.vvz.ethz.ch/Vorlesungsverzeichnis/sucheLehrangebot.view?lang=%s&semkez=%s&seite=0&studiengangTyp=&deptId=&studiengangAbschnittId=&lerneinheitstitel=&lerneinheitscode=&unterbereichAbschnittId=&rufname=&kpRange=0,999&lehrsprache=&bereichAbschnittId=&ansicht=1&katalogdaten=&wahlinfo=", language, semester)

	courseUrl1 := fmt.Sprintf("https://www.vvz.ethz.ch/Vorlesungsverzeichnis/sucheLehrangebot.view?lang=%s&search=on&semkez=%s&lerneinheitscode=", language, semester)
	courseUrl2 := "&_strukturAus=on&search=Suchen"
	//search on and Suchen?

	newCourses := 0
	var text []byte

	cacheFile := semester + ".json"
	if data, err := os.ReadFile(cacheFile); err == nil {
		text = data
	} else {
		// print url
		fmt.Println("Fetching VVZ page:", vvzUrl)
		resp, err := http.Get(vvzUrl)
		if err != nil {
			fmt.Println("Failed to fetch VVZ page:", err)
			return
		}
		defer resp.Body.Close()
		text, _ = io.ReadAll(resp.Body)
		_ = os.WriteFile(cacheFile, text, 0644)
	}
	// log the response

	coursePattern := regexp.MustCompile(`<b>(\d{3}-\d{4}-[A-Z0-9]{3})<\/b>`)
	matches := coursePattern.FindAllStringSubmatch(string(text), -1)

	if len(matches) == 0 {
		sendScrapingEnd(404)
		return
	}

	fmt.Println("Amount of matches:", len(matches))

	for _, match := range matches {
		courseNr := match[1]

		// check db
		_, err := db.GetCourseName(context, courseNr)
		if err == nil {
			fmt.Println("Course already exists in DB:", courseNr)
			continue
		}

		courseURL := courseUrl1 + courseNr + courseUrl2
		fmt.Println("Fetching course URL:", courseURL)

		resp, err := http.Get(courseURL)
		if err != nil {
			sendScrapingError(courseNr, courseURL)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		titlePattern := regexp.MustCompile(language + `">(.*?)<\/a><\/b>`)
		titleMatch := titlePattern.FindStringSubmatch(string(body))

		if titleMatch == nil {
			sendScrapingError(courseNr, courseURL)
			fmt.Println("Error scraping", courseNr, "-", courseURL)
			continue
		}

		course := titleMatch[0]

		// add course to db
		_, err = db.AddCourse(context, sql.AddCourseParams{CourseNumber: courseNr, CourseName: course})
		if err != nil {
			fmt.Println("Error adding course to DB:", err)
			sendScrapingError(courseNr, courseURL)
			continue
		}

		line := fmt.Sprintf("vvzScrapeResults%s - %s\n", courseNr, course)
		f, _ := os.OpenFile(resultFilePath, os.O_APPEND|os.O_WRONLY, 0644)
		f.WriteString(line)
		f.Close()

		newCourses++

		// Sleep for 1 second before the next iteration
		time.Sleep(1 * time.Second)
	}

	sendScrapingEnd(newCourses)
}
