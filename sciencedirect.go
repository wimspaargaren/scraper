package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/scraper/models"
)

func ProcessScienceDirect() {
	file, err := os.Open("sciencedirect2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	totalCounter := 100

	fileContent := string(b)
	lines := strings.Split(fileContent, "\n")
	article := models.Article{
		Metadata:     []byte("{}"),
		Platform:     models.PlatformScienceDirect,
		ResultNumber: totalCounter,
		Keywords:     []byte("{}"),
		Query:        `(fintech OR "financial technology" OR banking) AND (ai OR "artificial intelligence" OR ml OR "machine learning" OR "deep learning")`,
	}
	counter := 0
	for _, line := range lines {
		log.Infof(fmt.Sprintf("D%sX", line))
		if line == "" {
			log.Infof("WTF")
		}
		if counter == 0 {
			article.Authors = line[0 : len(line)-2]
		}
		if counter == 1 {
			article.Title = line[0 : len(line)-2]
		}
		if counter == 2 {
			article.Journal = line[0 : len(line)-2]
		}
		if strings.HasPrefix(line, "https://doi.org/") {
			article.Doi = line[16 : len(line)-1]
		}
		if strings.HasPrefix(line, "Abstract: ") {
			article.Abstract = line[10 : len(line)-1]
		}
		if strings.HasPrefix(line, "Keywords: ") {
			keywords := strings.Split(line[9:len(line)-1], ";")
			article.AddKeywords(models.Keywords{
				List: keywords,
			})
		}
		if strings.HasPrefix(line, "(") {
			article.URL = line[1 : len(line)-2]
		}
		counter++
		if line == "" || line == " " {
			counter = 0
			log.Infof("Adding article")
			err := articleDB.Add(context.Background(), &article)
			if err != nil {
				panic(err)
			}
			totalCounter++
			article = models.Article{
				Metadata:     []byte("{}"),
				Platform:     models.PlatformScienceDirect,
				Keywords:     []byte("{}"),
				ResultNumber: totalCounter,
				Query:        `(fintech OR "financial technology" OR banking) AND (ai OR "artificial intelligence" OR ml OR "machine learning" OR "deep learning")`,
			}
		}

	}
	log.Infof("Counter: %d", counter)
}
