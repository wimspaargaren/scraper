package main

import (
	"context"

	"github.com/abadojack/whatlanggo"
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
	articles, err := articleDB.List(context.Background())
	if err != nil {
		panic(err)
	}
	for _, article := range articles {
		infoTitle := whatlanggo.Detect(article.Title)
		if article.Abstract != "" {
			infoAbstract := whatlanggo.Detect(article.Title)
			if infoAbstract.Lang.String() != infoTitle.Lang.String() {
				log.Warningf("warning: %s is %s, but %s is %s", article.Title, infoTitle.Lang.String(), article.Abstract, infoAbstract.Lang.String())
			}
			article.Lang = infoTitle.Lang.String()
			if infoAbstract.Confidence > infoTitle.Confidence {
				article.Lang = infoAbstract.Lang.String()
			}
			err := articleDB.Update(context.Background(), article)
			if err != nil {
				log.Fatalf("error updating article: %s", err)
			}
		} else {
			if infoTitle.Confidence < 0.2 {
				log.Warningf("unclear; %s, lang: %s", article.Title, infoTitle.Lang.String())
			}
			article.Lang = infoTitle.Lang.String()
			err := articleDB.Update(context.Background(), article)
			if err != nil {
				log.Errorf("error updating article: %s", err)
			}
		}
	}

}
