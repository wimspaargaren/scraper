package models

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq" // Postgres driver
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

func InitDB(logMode bool) *gorm.DB {
	host := os.Getenv("PG_HOST")
	databaseName := os.Getenv("PG_DBNAME")
	user := os.Getenv("PG_USERNAME")
	pass := os.Getenv("PG_PASSWORD")
	port := os.Getenv("PG_PORT")
	if host == "" {
		host = "localhost"
	}
	defaultDBName := "scraper"
	if databaseName == "" {
		databaseName = defaultDBName
	}
	if user == "" {
		user = "postgres"
	}
	if pass == "" {
		pass = ""
	}
	pgPort, err := strconv.Atoi(port)
	if err != nil {
		pgPort = 5432
	}

	// Connect to DB
	dbHost := flag.String("host", host, "Defaults to localhost")
	dbName := flag.String("db", databaseName, "Defaults to "+defaultDBName)
	dbUser := flag.String("user", user, "Defaults to postgres")
	dbPass := flag.String("pass", pass, "Defaults to empty")
	dbPort := flag.Int("port", pgPort, "Defaults to 5432")
	flag.Parse()

	log.Infof("DB configuration:")
	log.Infof("DB HOST: %s", *dbHost)
	log.Infof("DB Name: %s", *dbName)
	log.Infof("DB User: %s", *dbUser)
	log.Infof("DB Port: %d", *dbPort)

	url := fmt.Sprintf("dbname=%s user=%s password=%s port=%d host=%s sslmode=disable ", *dbName, *dbUser, *dbPass, *dbPort, *dbHost)

	db, err := gorm.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	db.LogMode(logMode)

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	models := []interface{}{
		&Article{},
	}
	db.AutoMigrate(models...)

	return db
}
