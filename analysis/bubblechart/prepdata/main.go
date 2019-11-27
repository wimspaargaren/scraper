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
)

func main() {
	countMap := make(map[string]int)
	csvReader, err := GetCSVReader("input.csv", ',', '#', 2)
	if err != nil {
		log.Errorf("Error reading csv, error: %s", err.Error())
		panic(err)
	}
	counter := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if counter == 0 {
			counter++
			continue
		}
		counter++
		if len(record) != 2 {
			log.Errorf("not enough records: %d, for row: %d", len(record), counter)
			panic(fmt.Errorf("not enough records"))
		}
		key := fmt.Sprintf(`%s,%s`, translate(record[0]), translate(record[1]))
		if strings.Contains(key, "NA") {
			continue
		}
		log.Infof(key)
		countMap[key]++
	}

	file, err := os.Create("../input.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"amount", "fintech_topic", "ai_type", "stat"})
	if err != nil {
		panic(err)
	}
	for key, val := range countMap {
		log.Infof("%d,%s", val, key)
		countMap := strconv.Itoa(val)
		err := writer.Write(strings.Split(fmt.Sprintf(`%d,%s,%s`, val, key, countMap), ","))
		if err != nil {
			panic(err)
		}
	}

}

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

func translate(input string) string {
	switch input {
	case "Wealth management":
		return "Private banking"
	case "User trust":
		return "Overall Bank Topics"
	case "Personal loans":
		return "Loans and other credit products"
	case "Money Laundering":
		return "Customer Due Diligence"
	case "Monetary policy":
		return "Market Risk"
	case "Intrusion Detection":
		return "Malware"
	case "Checking and savings accounts":
		return "Loans and other credit products"
	case "Stock brokerage (discount and full-service)":
		return "Stock brokerage"
	default:
		return input
	}
}
