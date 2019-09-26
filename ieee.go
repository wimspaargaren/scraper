package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

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

		err = articleDB.Add(ctx, &models.Article{
			Year:         year,
			Description:  record[10],
			Title:        record[0],
			URL:          record[15],
			Platform:     models.PlatformIEEE,
			Query:        `("Document Title":decision* AND ( "Document Title":"software engineering" OR "Document Title":"software development"))`,
			ResultNumber: counter,
			Doi:          record[13],
			Metadata:     []byte("{}"),
		})
		if err != nil {
			log.Errorf("Error adding article: %s", err.Error())
		} else {
			counter++
		}
	}
	return nil
}
