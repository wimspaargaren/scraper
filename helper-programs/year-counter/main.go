package main

import (
	"context"
	"sort"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	"github.com/wimspaargaren/scraper/models"
)

type Finding struct {
	Year   int
	Amount int
}

var (
	articleDB   *models.ArticleDB
	db          *gorm.DB
	counter     = make(map[int]int)
	counterList = []Finding{}
)

func init() {
	db = models.InitDB(false)
	db.DB().SetMaxOpenConns(50)
	articleDB = models.NewArticleDB(db)
}

func main() {
	articles, err := articleDB.ListOnStatus(context.Background(), models.StatusUseful)
	if err != nil {
		panic(err)
	}
	for _, article := range articles {
		counter[article.Year]++
	}

	for key, val := range counter {
		counterList = append(counterList, Finding{
			Year:   key,
			Amount: val,
		})
	}
	sort.Slice(counterList, func(i, j int) bool {
		return counterList[i].Year > counterList[j].Year
	})

	for _, finding := range counterList {
		log.Infof("Year: %d, counter: %d", finding.Year, finding.Amount)
	}
}
