package main

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/scraper/models"
)

const (
	scholarLink  = "https://scholar.google.nl"
	acmLink      = "https://dl.acm.org"
	springerLink = "https://link.springer.com"
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
	log.Infof("Started scraping")
	start := time.Now()
	// Enable everything you need
	// go FindDOIs()
	// processDOILinks()
	// ProcessIEEEExport()
	// doSpringerScrape()
	// doScholarScrape()
	// doACMScraper()
	log.Infof("Done scraping in %s", time.Now().Sub(start))
}

// /results.cfm?query=acmdlTitle%3A%28"decision+*"+AND+%28"software+development"+"software process modeling"+"software organization"+"requirements engineering"+"software maintenance"+"information systems"+"software architecture"+"software design"+"software project management"+"software+engineering"%29%29&Go.x=40&Go.y=9
func doACMScraper() {
	acmScraper(`/results.cfm?query=acmdlTitle%3A%28"decision+*"+AND+%28"software+development"+"software%20process%20modeling"+"software%20organization"+"requirements%20engineering"+"software%20maintenance"+"information%20systems"+"software%20architecture"+"software%20design"+"software%20project%20management"+"software+engineering"%29%29&Go.x=40&Go.y=9`,
		0,
		0,
		`acmdlTitle%3A%28"decision+*"+AND+%28"software+development"+"software process modeling"+"software organization"+"requirements engineering"+"software maintenance"+"information systems"+"software architecture"+"software design"+"software project management"+"software+engineering"%29%`)
}

//Missing: information systems
func doScholarScrape() {
	// https: //scholar.google.nl/scholar?q=allintitle:+decision+%22software+development%22+OR+%22software%20organization%22+OR+%22requirements%20engineering%22+OR+%22software%20maintenance%22+OR+%22software%20architecture%22+OR+%22software%20design%22+OR+%22software%20project%20management%22+OR+%22software+engineering%22&hl=en&as_sdt=0,5
	scholarScraper("/scholar?q=allintitle:+decision+%22software+development%22+OR+%22software%20organization%22+OR+%22requirements%20engineering%22+OR+%22software%20maintenance%22+OR+%22software%20architecture%22+OR+%22software%20design%22+OR+%22software%20project%20management%22+OR+%22software+engineering%22&hl=en&as_sdt=0,5",
		0,
		0,
		`allintitle: decision "software development" OR "software organization" OR "requirements engineering" OR "software maintenance" OR "software architecture" OR "software design" OR "software project management" OR "software engineering"`)

	// https://scholar.google.nl/scholar?start=550&q=allintitle:+decision+%22software+development%22+OR+%22software+organization%22+OR+%22requirements+engineering%22+OR+%22software+maintenance%22+OR+%22software+architecture%22+OR+%22software+design%22+OR+%22software+project+management%22+OR+%22software+engineering%22&hl=en&as_sdt=0,5
	// scholarScraper("/scholar?start=550&q=allintitle:+decision+%22software+development%22+OR+%22software+organization%22+OR+%22requirements+engineering%22+OR+%22software+maintenance%22+OR+%22software+architecture%22+OR+%22software+design%22+OR+%22software+project+management%22+OR+%22software+engineering%22&hl=en&as_sdt=0,5",
	// 	55,
	// 	550,
	// 	`allintitle: decision "software development" OR "software organization" OR "requirements engineering" OR "software maintenance" OR "software architecture" OR "software design" OR "software project management" OR "software engineering"`)
	//https://scholar.google.nl/scholar?hl=en&as_sdt=0%2C5&as_vis=1&q=allintitle%3A+%22software+process+modeling%22+%22decision+*%22&btnG=
	scholarScraper("/scholar?hl=en&as_sdt=0%2C5&as_vis=1&q=allintitle%3A+%22software+process+modeling%22+%22decision+*%22&btnG=",
		0,
		0,
		`allintitle: "software process modeling" "decision *"`)
}

func doSpringerScrape() {
	linksToVisit := []string{
		`/search/page/1?dc.title=%22decision%22+%22software+engineering%22&date-facet-mode=between&showAll=true&query=%22software+development%22`,
		`/search/page/1?dc.title=decision+making&date-facet-mode=between&showAll=true&query=%22software+project+management%22`,
		`/search/page/1?dc.title=decision+making&date-facet-mode=between&showAll=true&query=%22software+design%22`,
		`/search/page/1?dc.title=decision+making&date-facet-mode=between&showAll=true&query=%22software+architecture%22`,
		`/search/page/1?dc.title=decision+making&date-facet-mode=between&showAll=true&query=%22software+maintenance%22`,
		`/search/page/1?dc.title=decision+making&date-facet-mode=between&showAll=true&query=%22requirements+engineering%22`,
		`/search/page/1?dc.title=decision+making&date-facet-mode=between&showAll=true&query=%22software+organization%22`,
		`/search/page/1?dc.title=decision+making&date-facet-mode=between&showAll=true&query=%22software+process+modeling%22`,
	}
	for _, search := range linksToVisit {
		log.Infof("Processing: %s", search)
		springerScraper(search, 0, 0, search[15:len(search)])

	}
}
