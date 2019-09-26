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
	// FindDOIs()
	// processDOILinks()
	// ProcessWebOfScienceExport()
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

// https://link.springer.com/search?query=fintech+AND+%28AI+OR+%22artificial+OR+intelligence%22+OR+ML+OR+%22machine+OR+learning%22+OR+%22deep+OR+learning%22%29&date-facet-mode=between&showAll=true#
// https://link.springer.com/search?query=%22financial+AND+technology%22+AND+%28AI+OR+%22artificial+OR+intelligence%22+OR+ML+OR+%22machine+OR+learning%22+OR+%22deep+OR+learning%22%29&date-facet-mode=between&showAll=true#machine+OR+learning%22+OR+%22deep+OR+learning%22%29&date-facet-mode=between&showAll=true#
func doSpringerScrape() {
	linksToVisit := []string{
		`/search?query=fintech+AND+%28AI+OR+%22artificial+OR+intelligence%22+OR+ML+OR+%22machine+OR+learning%22+OR+%22deep+OR+learning%22%29&date-facet-mode=between&showAll=true#`,
		`/search?query=%22financial+AND+technology%22+AND+%28AI+OR+%22artificial+OR+intelligence%22+OR+ML+OR+%22machine+OR+learning%22+OR+%22deep+OR+learning%22%29&date-facet-mode=between&showAll=true#machine+OR+learning%22+OR+%22deep+OR+learning%22%29&date-facet-mode=between&showAll=true#`,
	}
	for _, search := range linksToVisit {
		log.Infof("Processing: %s", search)
		springerScraper(search, 0, 0, search[15:len(search)])

	}
}

// https://link.springer.com/search?query=fintech+AND+%28%22financial+OR+technology%22+OR+AI+OR+%22artificial+OR+intelligence%22+OR+ML+OR+%22
