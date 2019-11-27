package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/scraper/models"
)

var (
	articleDB *models.ArticleDB
	db        *gorm.DB
)

func init() {
	db = models.InitDB(false)
	db.DB().SetMaxOpenConns(50)
	articleDB = models.NewArticleDB(db)
}

func main() {
	data := [][]string{{"ID", "Title", "Journals/Conferences", "Author", "Interconnection AI & SE", "Topic", "Year", "Study method", "Context", "Link", "DOI", "Comment", "Keywords"}}
	articles, err := articleDB.ListOnStatus(context.Background(), models.StatusUseful)
	if err != nil {
		panic(err)
	}
	for _, article := range articles {
		input, err := ioutil.ReadFile(fmt.Sprintf("../article-downloader/pdfsfolder/%s.pdf", article.ID))
		if err == nil {
			err = ioutil.WriteFile(fmt.Sprintf("relevant-articles/%s", article.ID), input, 0644)
			if err != nil {
				log.Errorf("Error writing file: %s", err)
				return
			}
		} else {
			log.Error(err)
		}
	}

	for _, article := range articles {
		keywords, err := article.GetKeywords()
		if err != nil {
			panic(err)
		}
		for i, keyword := range keywords.List {
			keywords.List[i] = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(keyword), "\n", ""), "   ", " "))
		}
		err = article.AddKeywords(keywords)
		if err != nil {
			panic(err)
		}
		err = articleDB.Update(context.Background(), article)
		if err != nil {
			panic(err)
		}
		keywordString := ""
		for i, keyword := range keywords.List {
			keywordString += strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(keyword), "\n", ""), "   ", " "))
			if i != len(keywords.List)-1 {
				keywordString += " , "
			}
		}
		log.Infof("article ID: %s", article.ID)
		sheet := []string{
			article.ID.String(),
			article.Title,
			article.Journal,
			article.Authors,
			``,
			``,
			fmt.Sprintf(`%d`, article.Year),
			``,
			``,
			article.URL,
			article.Doi,
			article.Comment,
			keywordString,
		}
		data = append(data, sheet)
		log.Infof("%v", sheet)
	}

	log.Infof("Total: %d", len(articles))
	file, err := os.Create("result.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			panic(err)
		}
	}
}
