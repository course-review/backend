package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
)

var language = "de"
var regexRuleCourseNumber = regexp.MustCompile("<b>(\\d{3}-\\d{4}-[A-Z0-9]{3})<\\/b>")
var regexRuleCourseName = regexp.MustCompile("\\w{2}\">(.*?)<\\/a><\\/b>")

func main() {
	vvzScraper("2025S")
}

func vvzScraper(semester string) {
	collector := colly.NewCollector(
		colly.AllowedDomains("www.vvz.ethz.ch"),
	)

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

	collector.Visit(vvzListUrl(1, semester, language))

	fmt.Println("Matched Items:")
	for _, item := range courseNames {
		fmt.Println("-", item)
	}
}

func vvzListUrl(page int, semester string, language string) string {
	return fmt.Sprintf("https://www.vvz.ethz.ch/Vorlesungsverzeichnis/sucheLehrangebot.view?seite=%d&semkez=%s&lang=%s", page, semester, language)
}
