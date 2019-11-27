package main

import (
	"context"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"github.com/wimspaargaren/scraper/models"
	"golang.org/x/net/publicsuffix"
)

func downloadArticlePDFs() {
	ctx := context.Background()
	articles, err := articleDB.ListOnStatus(ctx, models.StatusUseful)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())

	for _, article := range articles {
		if article.GotPdf {
			continue
		}
		if article.URL == "" {
			continue
		}
		if article.Platform == models.PlatformWebOfScience {
			// log.Warningf("Skipping ACM and Web of Science for now")
			continue
		}
		rand1 := rand.Intn(10)
		rand2 := rand.Intn(3000)
		log.Infof("Sleeping for: %s", time.Duration(rand1)*time.Second+time.Second*15+time.Millisecond*time.Duration(rand2))
		time.Sleep(time.Duration(rand1)*time.Second + time.Second*15 + time.Millisecond*time.Duration(rand2))

		log.Infof("Processing article: %s, from: %s", article.ID, article.Platform)
		if strings.HasSuffix(article.URL, ".pdf") {
			downloadPDFLink(article)
		} else if article.Platform == models.PlatformSpringer {
			downloadSpringerLink(article)
		} else if article.Platform == models.PlatformIEEE {
			downloadIEEE(article)
		} else if article.Platform == models.PlatformACM {
			if strings.HasPrefix(article.URL, "https://dlnext.acm.org/") {
				downloadACMNext(article)
			} else {
				downloadACMLink(article)
			}
		} else if article.Platform == models.PlatformGoogleScholar {
			if strings.HasPrefix(article.URL, "https://link.springer") {
				downloadSpringerLink(article)
			} else if strings.HasPrefix(article.URL, "https://dl.acm.org") {
				downloadACMLink(article)
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
					log.Errorf("Error opening article url: %s, err: %s", article.URL, err)
				}
			}
		}
		if !article.GotPdf {
			log.Warningf("Didnt find any article for: %s", article.ID)
		}
	}
}

func downloadACMNext(article *models.Article) {
	// https://dlnext.acm.org/doi/pdf/10.1145/3274287?download=true
	// url := "https://dlnext.acm.org/doi/abs/10.1145/3274287"
	url := strings.Replace(article.URL, "abs", "pdf", -1)
	url += "?download=true"
	err := DownloadFile("pdfsfolder/"+article.ID.String()+".pdf", url)
	if err != nil {
		log.Errorf("Error downloading pdf: %s", err)
	} else {
		err = articleDB.UpdatePDFFound(ctx, article)
		if err != nil {
			log.Errorf("Error updating pdf found for article: %s, err: %s", article.ID, err)
		}
	}
}

func downloadPDFLink(article *models.Article) {
	err := DownloadFile("pdfsfolder/"+article.ID.String()+".pdf", article.URL)
	if err != nil {
		log.Errorf("Error downloading pdf: %s", err)
	} else {
		err = articleDB.UpdatePDFFound(ctx, article)
		if err != nil {
			log.Errorf("Error updating pdf found for article: %s, err: %s", article.ID, err)
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
	err := DownloadIEEEFile("pdfsfolder/"+article.ID.String()+".pdf", article.URL)
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
	client, err := createClientWithCookie()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
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

// write as it downloads and not load the whole file into memory.
func DownloadIEEEFile(filepath string, url string) error {
	log.Infof("Downloading IEEE: %s", url)
	client, err := createClientWithCookie()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	splitted := strings.Split(bodyString, `"`)
	urlPDF := ""
	for _, part := range splitted {
		if strings.Contains(part, ".pdf") {

			urlPDF = part
			if strings.HasPrefix(urlPDF, "/iel") {
				urlPDF = "https://ieeexplore.ieee.org" + urlPDF
			}
			break
		}
	}
	time.Sleep(time.Second)
	log.Info(urlPDF)
	req, err = http.NewRequest("GET", urlPDF, nil)
	if err != nil {
		return err
	}

	resp, err = client.Do(req)
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

func createClientWithCookie() (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Jar: jar,
	}, nil
}

func writeToFile(fileName string, data string) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err = f.Write([]byte(data))
	if err != nil {
		panic(err)
	}
}
