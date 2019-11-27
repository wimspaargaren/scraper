package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"

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
	ctx := context.Background()
	articles, err := articleDB.ListOnStatus(ctx, models.StatusUseful)
	if err != nil {
		panic(err)
	}
	res := ""
	for _, article := range articles {
		splitAuthors := strings.Split(article.Authors, ",")
		if len(splitAuthors) == 1 {
			splitAuthors = strings.Split(article.Authors, " ")
		}
		firstAuthor := splitAuthors[0]
		firstAuthor = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(firstAuthor, " ", ""), ".", ""), "-", "")
		keywords, err := article.GetKeywords()
		if err != nil {
			panic(err)
		}
		res += "@article{" + firstAuthor + strconv.Itoa(article.Year) + ",\n"
		res += "\t author = {" + strings.ReplaceAll(strings.ReplaceAll(article.Authors, "\n", ""), "  ", " ") + "},\n"
		res += "\t" + ` file = {:article-downloader/final/` + article.ID.String() + `.pdf:pdf},` + "\n"
		if len(keywords.List) != 0 {
			keywordString := ""
			for i, keyword := range keywords.List {
				keywordString += strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(keyword), "\n", ""), "   ", " "))
				if i != len(keywords.List)-1 {
					keywordString += ", "
				}
			}
			res += "\t" + ` keywords = {` + keywordString + `},` + "\n"
		}
		res += "\t" + ` year = {` + strconv.Itoa(article.Year) + `},` + "\n"
		res += "\t" + ` doi = {` + article.Doi + `},` + "\n"
		res += "\t" + ` title = {{` + strings.ReplaceAll(article.Title, "\n", "") + `}}` + "\n"
		res += "}\n"
	}

	log.Infof("res: %s", res)
	f, err := os.Create("result.bib")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(res)
	if err != nil {
		panic(err)
	}

	err = f.Sync()
	if err != nil {
		panic(err)
	}
}
