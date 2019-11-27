package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/wimspaargaren/scraper/models"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func acmNextScraper(queryURL string, callDepth, number int, query string) {
	log.Info("Processing acm next request")
	files, err := ioutil.ReadDir("acmnextpages/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		log.Infof(f.Name())
		dat, err := ioutil.ReadFile("acmnextpages/" + f.Name())
		if err != nil {
			log.Errorf("error reading file: %s", err)
		}
		reader := bytes.NewReader(dat)
		processACMENextResponse(reader, callDepth, number, query)
	}
}

func processACMENextResponse(reader io.Reader, callDepth, number int, query string) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Errorf("Could not create document from read, error: %s", err.Error())
		return
	}
	counter := 0
	doc.Find(".issue-item--search").Each(func(i int, s *goquery.Selection) {
		mainElement := s.Find(".issue-item__title").Find("a")
		link, ok := mainElement.Attr("href")
		if !ok {
			log.Error("Could not find article link")
		}
		doi := ""
		if strings.HasPrefix(link, "/doi/abs/") {
			doi = link[9:len(link)]
		}
		link = acmeNext + link
		title := mainElement.Text()

		description := s.Find(".issue-item__abstract").Text()
		description = fixString(description)
		yeartext := s.Find(".issue-item__detail").Text()
		year := getYear(yeartext)

		citedBy := 0
		citeString := s.Find(".citation").Text()
		temp := ""
		if len(citeString) > 16 {
			temp = citeString[16:len(citeString)]
		}
		citedBy, err = strconv.Atoi(temp)
		if err != nil {
			log.Errorf("Error getting cited by, error: %s", err.Error())
		}
		if description == "" {
			counter++
		}
		err = articleDB.Add(ctx, &models.Article{
			Year:         year,
			Abstract:     description,
			Title:        title,
			URL:          link,
			Platform:     models.PlatformACM,
			Query:        query,
			ResultNumber: number,
			Doi:          doi,
			Cited:        citedBy,
			Metadata:     []byte("{}"),
			Keywords:     []byte("{}"),
		})

		if err != nil {
			log.Errorf("Error adding article: %s", err.Error())
		} else {
			number++
		}
	})
}

func getYear(text string) int {
	log.Infof("Text: %s", text)
	for i := 1980; i < 2021; i++ {
		if strings.Contains(text, strconv.Itoa(i)) {
			log.Infof("FOUND: %d", i)
			return i
		}
	}
	return 0
}
