package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	"github.com/wimspaargaren/scraper/models"
)

var (
	articleDB *models.ArticleDB
	db        *gorm.DB
	counter   = make(map[string]int)
)

func init() {
	db = models.InitDB(false)
	db.DB().SetMaxOpenConns(50)
	articleDB = models.NewArticleDB(db)
}

func main() {
	missingDois := 0
	articlesToBeRemoved := 0
	articles, err := articleDB.List(context.Background())
	if err != nil {
		panic(err)
	}
	for _, article := range articles {
		if article.Doi != "" {
			counter[article.Doi]++
		} else {
			missingDois++
		}
	}
	for key, val := range counter {
		if val > 1 {
			log.Infof("Duplicate doi: %s", key)
			articlesToBeRemoved += val - 1
			articles, err := articleDB.ListOnDoi(context.Background(), key)
			if err != nil{
				log.Errorf("unable to list articles on doi: %s",err)
			}
			temp := val
			for _, article := range articles {
				log.Infof("URL: %s",article.URL)
				if temp > 1 {
					err := articleDB.Delete(context.Background(), article.ID)
					if err != nil{
						log.Errorf("unable to delete article")
					}
					temp--
				}

			}
		}
	}
	log.Infof("Number of articles to be removed: %d", articlesToBeRemoved)
	log.Infof("Number of missing dois: %d", missingDois)
}
