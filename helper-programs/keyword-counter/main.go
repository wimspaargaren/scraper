package main

import (
	"context"
	"fmt"

	"github.com/bradfitz/slice"
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	"github.com/wimspaargaren/scraper/models"
)

var (
	articleDB  *models.ArticleDB
	db         *gorm.DB
	keywordMap map[string]int
)

func init() {
	db = models.InitDB(false)
	db.DB().SetMaxOpenConns(50)
	articleDB = models.NewArticleDB(db)
	keywordMap = make(map[string]int)
}

type Keyword struct {
	Word   string
	Amount int
}

func main() {
	articles, err := articleDB.ListOnStatus(context.Background(), models.StatusUseful)
	if err != nil {
		panic(err)
	}

	for _, article := range articles {
		keywords, err := article.GetKeywords()
		if err != nil {
			panic(err)
		}
		for _, keyword := range keywords.List {
			keywordMap[keyword]++
		}
	}

	keywordList := []Keyword{}
	for key, val := range keywordMap {
		keywordList = append(keywordList, Keyword{Word: key, Amount: val})
	}
	slice.Sort(keywordList[:], func(i, j int) bool {
		return keywordList[i].Amount > keywordList[j].Amount
	})

	for _, keyword := range keywordList {
		if keyword.Word == "" {
			continue
		}
		log.Infof(fmt.Sprintf("%s: %d", keyword.Word, keyword.Amount))
	}
}
