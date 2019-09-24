// &as_ylo=2018&as_yhi=2018
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

const (
	base   string = "https://scholar.google.com/scholar?hl=en&as_sdt=1%2C5&as_vis=1&"
	suffix string = "&btnG=&"
)

type Finding struct {
	Year   int
	Amount string
}

type Query struct {
	Query    string
	Findings []Finding
}

func main() {

	queries := []string{
		`q=("fintech")+AND+("artificial+intelligence")`,
		`q=("fintech")+AND+("AI")`,
		`q=%28"fintech"%29+AND+%28"ML"%29`,
		`q=%28"fintech"%29+AND+%28"machine+learning"%29`,
		`q=%28"fintech"%29+AND+%28"deep+learning"%29`,
		`q=("financial+technology")+AND+("artificial+intelligence")`,
		`q=("financial+technology")+AND+("AI")`,
		`q=%28"financial+technology"%29+AND+%28"ML"%29`,
		`q=%28"financial+technology"%29+AND+%28"machine+learning"%29`,
		`q=%28"financial+technology"%29+AND+%28"deep+learning"%29`,
	}

	// queriesRes := []Query{}
	for _, query := range queries {
		queryRes := Query{
			Query:    query,
			Findings: []Finding{},
		}
		for year := 2015; year < 2021; year++ {

			hi := "as_ylo=%d&as_yhi=%d"
			test := fmt.Sprintf(hi, year, year)
			// log.Infof("URL: %s", fmt.Sprintf("%s%s%s%s", base, query, suffix, test))
			resp, err := http.Get(fmt.Sprintf("%s%s%s%s", base, query, suffix, test))
			if err != nil {
				log.Errorf("Error doing acm request")
				return
			}
			defer resp.Body.Close()
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Errorf("Could not create document from read, error: %s", err.Error())
				return
			}
			doc.Find(".gs_ab_mdw").Each(func(i int, s *goquery.Selection) {
				html, err := s.Html()
				if err != nil {
					log.Errorf("Error converting to html: %s", err)
					return
				}
				if strings.Contains(html, "results") {
					// log.Info(html)
					splitted := strings.Split(html, "results")
					// log.Infof("Res: %s", splitted[0])
					queryRes.Findings = append(queryRes.Findings, Finding{
						Year:   year,
						Amount: splitted[0],
					})
				}
			})
			rand := rand.Intn(1000)
			time.Sleep(time.Second*5 + time.Millisecond*time.Duration(rand))

		}
		sort.Slice(queryRes.Findings, func(i, j int) bool {
			return queryRes.Findings[i].Year > queryRes.Findings[j].Year
		})
		fmt.Println("Query: ", queryRes.Query)
		for _, finding := range queryRes.Findings {
			fmt.Println("Year", finding.Year, "amount", finding.Amount)
		}
	}
}
