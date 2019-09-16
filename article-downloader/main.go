package main

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/wimspaargaren/literature-scraper/models"
)

var (
	articleDB *models.ArticleDB
	db        *gorm.DB
	ctx       = context.Background()
)

func init() {
	db = models.InitDB(false)
	db.DB().SetMaxOpenConns(50)
	articleDB = models.NewArticleDB(db)
}

func main() {
	downloadArticlePDFs()
}
