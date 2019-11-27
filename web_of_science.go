package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/scraper/models"
)

// ProcessWebOfScienceExport save web of science export to the database
func ProcessWebOfScienceExport(fileName string, counter int) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("error reading file: %s", err)
	}
	reader := bytes.NewReader(dat)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Errorf("Could not create document from read, error: %s", err.Error())
		return
	}

	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		mapTest := make(map[string]string)

		s.Find("tr").Each(func(j int, s *goquery.Selection) {
			tableColumns := s.Children()
			tableColumns.Get(0)
			first := ""
			s.Find("td").Each(func(x int, s *goquery.Selection) {
				if x == 0 {
					first = s.Text()
				} else if x == 1 {
					mapTest[first] = s.Text()
				}
			})
		})
		if len(mapTest) > 6 {
			article := models.Article{
				Metadata:     []byte("{}"),
				Keywords:     []byte("{}"),
				Platform:     models.PlatformWebOfScience,
				ResultNumber: counter,
				Query:        `("fintech" OR "banking" OR "financial technology") AND ("AI" OR "artificial intelligence" OR "ML" OR "machine learning"OR "deep learning")`,
			}
			for key, val := range mapTest {
				switch key {
				case "AB ":
					article.Abstract = val
				case "PY ":
					year, err := strconv.Atoi(val)
					if err != nil {
						log.Errorf("cant convert year string to int: %s", val)
					}
					article.Year = year
				case "TI ":
					log.Infof("Title: %s", val)
					article.Title = val
				case "TC ":
					log.Infof("Cited: %s", val)
					cited, err := strconv.Atoi(val)
					if err != nil {
						log.Errorf("cant convert cited string to int: %s", val)
					}
					article.Cited = cited
				case "DI ":
					log.Infof("DOI: %s", val)
					article.Doi = val
				case "LA ":
					article.Lang = val
				case "DE ":
					log.Infof("KEywords: %s", val)
					article.AddKeywords(models.Keywords{
						List: strings.Split(val, ";"),
					})
				case "AF ":
					article.Authors = strings.Replace(val, "<br>", "", -1)
				case "SO ":
					article.Journal = val
				}
			}
			err := articleDB.Add(context.Background(), &article)
			if err != nil {
				log.Errorf("error adding article: %s", err)
			}
			counter++
			log.Infof("length: %d", len(mapTest))

		}
	})
	log.Infof("Articles added: %d", counter)
}
