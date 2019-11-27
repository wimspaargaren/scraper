package main

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"code.sajari.com/docconv"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/wimspaargaren/scraper/models"
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
	// cleanupCorruptPDFS()

	// downloadArticlePDFs()
}

func cleanupCorruptPDFS() {
	files, err := ioutil.ReadDir("final")
	if err != nil {
		log.Fatal(err)
	}
	counter := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".pdf") {
			go checkIfNeededToBeRemoved(f)
			counter++
		}
	}
	log.Infof("counter: %d", counter)
}

func checkIfNeededToBeRemoved(f os.FileInfo) {
	log.Infof("processing: %s", f.Name())
	_, err := docconv.ConvertPath("final/" + f.Name())
	if err != nil {
		firstPart := strings.Split(f.Name(), ".pdf")[0]
		articleID, err := uuid.FromString(firstPart)
		if err != nil {
			panic(err)
		}
		article, err := articleDB.Get(context.Background(), articleID)
		if err != nil {
			panic(err)
		}
		article.GotPdf = false

		err = articleDB.Db.Model(&article).Update("got_pdf", false).Error
		if err != nil {
			panic(err)
		}
		err = os.Remove("final/" + f.Name())
		if err != nil {
			panic(err)
		}
		log.Infof("deleted: %s", "final/"+f.Name())
	}
}
