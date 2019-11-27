package main

import (
	"context"
	"net/http"
	"strconv"

	uuid "github.com/satori/go.uuid"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/scraper/models"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
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

func UpdateArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	amountString := ps.ByName("amount")
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		log.Errorf("Err: %s", err)
		return
	}
	articleID := uuid.FromStringOrNil(ps.ByName("id"))
	article := models.Article{
		ID:    articleID,
		Cited: amount,
	}
	err = articleDB.Update(context.Background(), &article)
	if err != nil {
		log.Errorf("Unable to update article id: %s, err: %s", articleID, err)
	}
	w.WriteHeader(200)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Errorf("cant response properly: %s", err)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/:id/:amount", UpdateArticle)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	})
	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))
}
