package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/wimspaargaren/scraper/models"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func acmScraper(queryURL string, callDepth, number int, query string) {
	log.Info("Processing acm request")
	url := acmLink + queryURL
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error doing acm request")
		return
	}
	defer resp.Body.Close()
	processACMResponse(resp, callDepth, number, query)
}

func processACMResponse(response *http.Response, callDepth, number int, query string) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Errorf("Could not create document from read, error: %s", err.Error())
		return
	}
	doc.Find(".details").Each(func(i int, s *goquery.Selection) {
		mainElement := s.Find(".title").Find("a")
		link, ok := mainElement.Attr("href")
		if !ok {
			log.Error("Could not find article link")
		}
		link = acmLink + "/" + link
		title := mainElement.Text()

		description := s.Find(".abstract").Text()
		description = fixString(description)

		yeartext := s.Find(".publicationDate").Text()
		yearstring := yeartext[len(yeartext)-4 : len(yeartext)]
		year, err := strconv.Atoi(yearstring)
		if err != nil {
			log.Errorf("Error parsing year: %s, error: %s", yearstring, err.Error())
		}

		citedBy := 0
		citeString := s.Find(".citedCount").Text()
		temp := ""
		if len(citeString) > 16 {
			temp = citeString[16:len(citeString)]
		}
		citedBy, err = strconv.Atoi(temp)
		if err != nil {
			log.Errorf("Error getting cited by, error: %s", err.Error())
		}
		err = articleDB.Add(ctx, &models.Article{
			Year:         year,
			Abstract:  description,
			Title:        title,
			URL:          link,
			Platform:     models.PlatformACM,
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
	nextURL, nextURLErr := getNextACMURL(doc, callDepth)
	if nextURLErr == nil {
		acmScraper(nextURL, callDepth+1, number, query)
	} else {
		log.Error("Could not find next page link")
	}
}

func getNextACMURL(doc *goquery.Document, callDepth int) (string, error) {
	next := strconv.Itoa(callDepth + 2)
	resURL := ""
	err := errors.New("Next link not found")
	doc.Find(".pagelogic").Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == next {
			url, ok := s.Attr("href")
			if ok {
				resURL = "/" + url
				err = nil
				return
			}
		}
	})
	return resURL, err
}
