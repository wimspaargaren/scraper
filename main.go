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
	acmeNext     = "https://dlnext.acm.org"
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
	// doACMNextScrape()
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

// /results.cfm?query=acmdlTitle%3A%28"decision+*"+AND+%28"software+development"+"software process modeling"+"software organization"+"requirements engineering"+"software maintenance"+"information systems"+"software architecture"+"software design"+"software project management"+"software+engineering"%29%29&Go.x=40&Go.y=9
func doACMNextScrape() {
	acmNextScraper(`/action/doSearch?AllField=%28"fintech"+OR+"financial+technology"%29+AND+%28"AI"+OR+"artificial+intelligence"+OR+"ML"+OR+"machine+learning"OR+"deep+learning"%29&expand=all&startPage=`,
		0,
		0,
		`("fintech" OR "financial technology") AND ("AI" OR "artificial intelligence" OR "ML" OR "machine learning"OR "deep learning")`)
}

//Missing: information systems
func doScholarScrape() {
	scholarScraper(`/scholar?start=980&q=(("fintech"+OR+"financial+technology")+AND+("AI"+OR+"artificial+intelligence"+OR+"ML"+OR+"machine+learning"+OR+"deep+learning"))&hl=en&as_sdt=0,5`,
		0,
		0,
		`("fintech" OR "financial technology") AND ("AI" OR "artificial intelligence" OR "ML" OR "machine learning"OR "deep learning"))`)

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
