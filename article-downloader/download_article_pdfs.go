package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"github.com/wimspaargaren/literature-scraper/models"
)

func downloadArticlePDFs() {
	ctx := context.Background()
	articles, err := articleDB.ListOnStatus(ctx, models.StatusUseful)
	if err != nil {
		panic(err)
	}
	for _, article := range articles {
		if article.GotPdf {
			continue
		}
		if strings.HasSuffix(article.URL, ".pdf") || article.Platform == models.PlatformIEEE {
			log.Infof("Article; %s", article.ID)
			err = DownloadFile("pdfsfolder/"+article.ID.String()+".pdf", article.URL)
			if err != nil {
				log.Errorf("Error downloading pdf: %s", err)
			} else {
				err = articleDB.UpdatePDFFound(ctx, article)
				if err != nil {
					log.Errorf("Error updating pdf found for article: %s, err: %s", article.ID, err)
				}
			}
		} else if article.Platform == models.PlatformSpringer {
			downloadSpringerLink(article)
		} else if article.Platform == models.PlatformACM {
			downloadACMLink(article)
		} else if article.Platform == models.PlatformGoogleScholar {
			if strings.HasPrefix(article.URL, "https://link.springer") {
				downloadSpringerLink(article)
			} else if strings.HasPrefix(article.URL, "https://dl.acm.org") {
				downloadACMLink(article)
			} else if strings.HasPrefix(article.URL, "https://ieeexplore.ieee.org") {
				downloadIEEE(article)
			} else if strings.HasPrefix(article.URL, "http://citeseerx.ist.psu.edu") {
				err = DownloadFile("pdfsfolder/"+article.ID.String()+".pdf", article.URL)
				if err != nil {
					log.Errorf("Error downloading pdf: %s", err)
				} else {
					err = articleDB.UpdatePDFFound(ctx, article)
					if err != nil {
						log.Errorf("Error updating pdf found for article: %s, err: %s", article.ID, err)
						return
					}
				}
			} else {
				err = open.Run(article.URL)
				if err != nil {
					fmt.Errorf("Error opening article url: %s, err: %s", article.URL, err)
				}
			}
		}
	}
}

func downloadSpringerLink(article *models.Article) {
	resp, err := http.Get(article.URL)
	if err != nil {
		log.Errorf("Error doing scholar request")
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("Could not create document from read, error: %s", err.Error())
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if article.GotPdf {
			return
		}
		link, ok := s.Attr("href")
		if ok && strings.HasSuffix(link, ".pdf") {
			if s.Text() == "Download" || (strings.Contains(s.Text(), "Download") && strings.Contains(s.Text(), "PDF") && !strings.Contains(s.Text(), "book")) {
				log.Infof("Hi: %s", link)
				if !strings.HasPrefix(link, "http") {
					link = "https://link.springer.com" + link
				}

				err = DownloadFile("pdfsfolder/"+article.ID.String()+".pdf", link)
				if err != nil {
					log.Errorf("Error downloading pdf: %s", err)
				} else {
					err = articleDB.UpdatePDFFound(ctx, article)
					if err != nil {
						log.Errorf("Error updating pdf found for article: %s, err: %s", article.ID, err)
						return
					}
				}
			}

		}
	})
}

func downloadACMLink(article *models.Article) {
	log.Infof("artic: %s", article.URL)
	resp, err := http.Get(article.URL)
	if err != nil {
		log.Errorf("Error doing scholar request")
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("Could not create document from read, error: %s", err.Error())
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if article.GotPdf {
			return
		}
		link, ok := s.Attr("href")
		if ok && s.Text() == "PDF" {
			if !strings.HasPrefix(link, "http") {
				link = "https://dl.acm.org/" + link
			}

			err = DownloadFile("pdfsfolder/"+article.ID.String()+".pdf", link)
			if err != nil {
				log.Errorf("Error downloading pdf: %s", err)
			} else {
				err = articleDB.UpdatePDFFound(ctx, article)
				if err != nil {
					log.Errorf("Error updating pdf found for article: %s, err: %s", article.ID, err)
					return
				}
			}

		}
	})
}

func downloadIEEE(article *models.Article) {
	splitted := strings.Split(article.URL, "/")
	id := splitted[len(splitted)-2]
	log.Infof("%s", id)
	downloadLINk := "https://ieeexplore.ieee.org/stamp/stamp.jsp?tp=&arnumber=" + id
	err := DownloadFile("pdfsfolder/"+article.ID.String()+".pdf", downloadLINk)
	if err != nil {
		log.Errorf("Error downloading pdf: %s", err)
	} else {
		err = articleDB.UpdatePDFFound(ctx, article)
		if err != nil {
			log.Errorf("Error updating pdf found for article: %s, err: %s", article.ID, err)
			return
		}
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	log.Infof("Downloading: %s", url)
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
