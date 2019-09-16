package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wimspaargaren/literature-scraper/models"

	"github.com/PuerkitoBio/goquery"
)

func scholarScraper(queryURL string, callDepth, number int, query string) {
	log.Info("Processing scolar request")
	url := scholarLink + queryURL
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error doing scholar request")
		return
	}
	processScholarResponse(resp, callDepth, number, query)
	defer resp.Body.Close()
}

func processScholarResponse(response *http.Response, callDepth, number int, query string) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Errorf("Could not create document from read, error: %s", err.Error())
		return
	}
	doc.Find(".gs_ri").Each(func(i int, s *goquery.Selection) {
		mainElement := s.Find("h3").Find("a")
		link, ok := mainElement.Attr("href")
		if !ok {
			log.Warningf("Could not find article link for number: %d", number)
		}
		title := mainElement.Text()
		if title == "" {
			title = s.Find("h3").Text()
		}
		title = fixString(title)

		description := s.Find(".gs_rs").Text()
		description = fixString(description)

		citedBy := 0
		s.Find(".gs_fl").Children().Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "Cited by") {
				citeString := s.Text()
				temp := citeString[9:len(citeString)]
				citedBy, err = strconv.Atoi(temp)
				if err != nil {
					log.Errorf("Error getting cited by, error: %s", err.Error())
				}
			}
		})

		year := getScholaryear(s.Find(".gs_a").Text())
		log.Infof("Adding article number: %d", number)
		err := articleDB.Add(ctx, &models.Article{
			Year:         year,
			Description:  description,
			Title:        title,
			URL:          link,
			Platform:     models.PlatformGoogleScholar,
			Query:        query,
			ResultNumber: number,
			Cited:        citedBy,
			Metadata:     []byte("{}"),
		})
		if err != nil {
			log.Errorf("Error adding article: %s", err.Error())
		} else {
			number++
		}
	})

	nextURL, ok := doc.Find(".gs_ico_nav_next").Parent().Attr("href")
	if ok {
		rand := rand.Intn(3000)
		time.Sleep(time.Duration(rand) * time.Millisecond)
		scholarScraper(nextURL, callDepth+1, number, query)
	} else {
		log.Errorf("Could not find next page link at calldepth: %d, next result to add should be: %d", callDepth, number)
	}

}

func getScholaryear(yearstring string) int {
	splitted := strings.Split(yearstring, " ")
	for _, possibleYear := range splitted {
		yearParsed, err := strconv.Atoi(strings.TrimSpace(possibleYear))
		if err == nil {
			return yearParsed

		}
	}

	log.Error("Could not find article year")
	return 0
}
