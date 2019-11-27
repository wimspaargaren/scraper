package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/scraper/models"
)

//GetCSVReader from filepath
func GetCSVReader(filePath string, seperator, comment rune, fieldsPerRecord int) (*csv.Reader, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	//CSV settings
	r := csv.NewReader(bufio.NewReader(csvFile))
	r.Comma = seperator
	r.Comment = comment
	r.FieldsPerRecord = fieldsPerRecord
	return r, nil
}

//ProcessIEEEExport processIEEE export csv
func ProcessIEEEExport() error {
	csvReader, err := GetCSVReader("ieee-xplore-export.csv", ',', '#', 30)
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
		year, err := strconv.Atoi(record[5])
		if err != nil {
			log.Errorf("Error parsing year, error: %s on row: %d", err.Error(), counter)
			return err
		}

		article := &models.Article{
			Year:         year,
			Abstract:  record[10],
			Title:        record[0],
			Authors:      record[1],
			URL:          record[15],
			Platform:     models.PlatformIEEE,
			Query:        `(("All Metadata":"financial technology" OR "All Metadata":"banking" OR "All Metadata":"fintech") AND ( "All Metadata":"AI" OR "All Metadata":"artificial intelligence" OR "All Metadata":"ML" OR "All Metadata":"machine learning" OR "All Metadata":"deep learning"))`,
			ResultNumber: counter,
			Doi:          record[13],
			Metadata:     []byte("{}"),
		}
		article.AddKeywords(models.Keywords{
			List: strings.Split(record[16], ";"),
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
