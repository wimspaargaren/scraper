package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
	"github.com/wimspaargaren/literature-scraper/models"
)

//FindDOIs finding dois
func FindDOIs() error {
	for {
		articles, err := articleDB.ListNoDOI(ctx)
		if err != nil {
			return err
		}
		if len(articles) == 0 {
			log.Info("No articles to process")
		}
		for _, article := range articles {
			log.Infof("Processing article: %s", article.ID.String())
			err = processArticle(article)
			if err != nil {
				log.Errorf("Error processing article: %s, error:%s", article.ID.String(), err.Error())
			}
		}
		log.Infof("Sleep doi finder")
		time.Sleep(time.Second * 5)
	}
}

func processArticle(article *models.Article) error {
	fmt.Println("Searching:", article.Title)
	url := `https://search.crossref.org/?q=` + url.QueryEscape(article.Title)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error doing acm request")
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("Could not create document from read, error: %s", err.Error())
		return err
	}

	searchErr := errors.New("not found")
	found := false
	doc.Find(".item-data").Each(func(i int, s *goquery.Selection) {
		if !found {
			val, ok := s.Find("a").Attr("href")
			if ok {
				fmt.Println("SET TO:", val)
				if strings.Contains(val, "doi") {
					article.Doi = val
					err := articleDB.Update(ctx, article)
					if err != nil {
						log.Errorf("Aaaw, error: %s", err.Error())
					} else {
						searchErr = nil
						found = true
					}
				}
			}
		}

	})
	return searchErr
}

func processDOILinks() error {
	for {
		articles, err := articleDB.ListDOILinks(ctx)
		if err != nil {
			return err
		}
		for _, article := range articles {
			unescaped, err := url.PathUnescape(article.Doi)
			if err != nil {
				log.Error("error: %s", err.Error())
				continue
			}
			if strings.HasPrefix(unescaped, "https://doi.org/") {
				article.Doi = unescaped[16:len(unescaped)]
				articleDB.Update(ctx, article)
			} else if strings.HasPrefix(unescaped, "https://plu.mx/plum/a/?") {
				article.Doi = unescaped[23:len(unescaped)]
				articleDB.Update(ctx, article)
			} else if strings.HasPrefix(unescaped, "http://dx.doi.org/") {
				article.Doi = unescaped[18:len(unescaped)]
				articleDB.Update(ctx, article)
			} else if strings.HasPrefix(unescaped, "https://crossmark.crossref.org/dialog/?doi=") {
				article.Doi = unescaped[43:len(unescaped)]
				articleDB.Update(ctx, article)
			}
		}
		log.Infof("Sleep doi links")
		time.Sleep(time.Second * 5)
	}
}