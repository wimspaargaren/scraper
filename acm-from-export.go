package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/scraper/models"
)

//ProcessACMExport process ACM export csv
func ProcessACMExport(fileName string) error {
	csvReader, err := GetCSVReader(fileName, ',', '#', 27)
	if err != nil {
		log.Errorf("Error reading csv, error: %s", err.Error())
		return err
	}
	counter := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if len(record) < 6 {
			log.Errorf("not enough records: %d, for row: %d", len(record), counter)
			return fmt.Errorf("not enough records")
		}
		year, err := strconv.Atoi(record[18])
		if err != nil {
			log.Errorf("Error parsing year, error: %s on row: %d", err.Error(), counter)
			return err
		}

		article := &models.Article{
			Year:         year,
			Abstract:     record[16],
			Title:        record[6],
			Authors:      record[2],
			URL:          fmt.Sprintf("https://dl.acm.org/citation.cfm?id=%s", record[1]),
			Platform:     models.PlatformACM,
			Query:        ` ("fintech" OR "banking" OR "financial technology") AND ("AI" OR "artificial intelligence" OR "ML" OR "machine learning"OR "deep learning")`,
			ResultNumber: counter,
			Doi:          record[11],
			Metadata:     []byte("{}"),
			Journal:      record[12],
		}
		article.AddKeywords(models.Keywords{
			List: strings.Split(record[10], ","),
		})
		err = articleDB.Add(ctx, article)
		if err != nil {
			log.Errorf("Error adding article: %s", err.Error())
		} else {
			counter++
		}
	}
	return nil
}
