package main

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/literature-scraper/models"

	"github.com/PuerkitoBio/goquery"
)

func springerScraper(queryURL string, callDepth, number int, query string) {
	log.Info("Processing springer request")
	url := springerLink + queryURL
	log.Infof("Url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error doing springer request")
		return
	}
	defer resp.Body.Close()
	processSpringerResponse(resp, callDepth, number, query)
}

func processSpringerResponse(response *http.Response, callDepth, number int, query string) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Errorf("Could not create document from read, error: %s", err.Error())
		return
	}
	doc.Find("#results-list").Find("li").Each(func(i int, s *goquery.Selection) {
		mainElement := s.Find(".title")
		link, ok := mainElement.Attr("href")
		if !ok {
			log.Error("Could not find article link")
		}
		link = springerLink + link
		title := mainElement.Text()

		description := s.Find(".snippet").Text()
		description = fixString(description)
		year := 0
		yearString, ok := s.Find(".year").Attr("title")
		if !ok {
			log.Errorf("Could not find year")
		} else {
			year, err = strconv.Atoi(yearString)
			if err != nil {
				splitted := strings.Split(yearString, " ")
				yearStringCorrection := splitted[len(splitted)-1]
				year, err = strconv.Atoi(yearStringCorrection)
				if err != nil {
					log.Errorf("error parsing year string, error: %s", err.Error())
				}
			}
		}
		if title != "" {
			err = articleDB.Add(ctx, &models.Article{
				Year:         year,
				Description:  description,
				Title:        title,
				URL:          link,
				Platform:     models.PlatformSpringer,
				Query:        query,
				ResultNumber: number,
				Metadata:     []byte("{}"),
			})
			if err != nil {
				log.Errorf("Error adding article: %s", err.Error())
			} else {
				number++
			}
		}
	})
	nextURL, nextURLErr := getNextSpringerURL(doc, callDepth)
	if nextURLErr == nil {
		rand := rand.Intn(3000)
		log.Infof("next url: %s, sleeping for: %d milliseconds", nextURL, rand)
		time.Sleep(time.Duration(rand) * time.Millisecond)
		springerScraper(nextURL, callDepth+1, number, query)
	} else {
		log.Error("Could not find next page link")
	}
}

func getNextSpringerURL(doc *goquery.Document, callDepth int) (string, error) {
	link, ok := doc.Find(".functions-bar").Find(".next").Attr("href")
	if !ok {
		return "", errors.New("Error finding next URLß")
	}
	return link, nil
}
