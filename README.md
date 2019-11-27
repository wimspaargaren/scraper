# Scraper tool

Simple scraping/processing tool for gathering information about articles from ACM, ScienceDirect, WebOfScience, IEEE, Springer and Google Scholar.

Please note, this repo does not provide a general interface for scraping neither does it provide exemplary GO code. This is repo is developed during a previous mapping study which can be used for another one by some small adjustment.

# Prerequisites

1. Go 1.13 (for mac simply run `brew install go`)
2. Postgresql 11.4 or higher (for mac simply run `brew install postgresql`)
3. A database :)

# Scrape articles

## ACM, Scholar, Springer

Build and run the root folder. Uncomment one of the functions:
1. `doACMScraper`
2. `doScholarScrape`
3. `doSpringerScrape`

The first parameter is the search url created by the search engine.

For example for scraping scholar with the term `software` go to [scholar](https://scholar.google.nl), type `software`, hit enter and copy `/scholar?hl=nl&as_sdt=0%2C5&q=scholar&btnG=`. The second and third option should be kept as `0`. The fourth option is saved to the database as reference for which query you used, so you can basically use anything.

## IEEE

IEEE is a bit simpler. Just use the provided CSV export and add it to the root of the project as: `ieee-xplore-export.csv` and run `ProcessIEEEExport`.

## Update DOIs

The scraper is not always able to correctly find the DOIs. Running `FindDOIs()` and `processDOILinks()` tries finding the DOIs, however in my experience this isn't always perfect as well.

## Downloading pdfs

cd into the `article-downloader` folder, build and run. It searches for articles which have `Userful` status in the database which means the `processed` column should equal `3`, so if you want to download everything execute `UPDATE articles SET processed = 3` in your db.

## Build and run

Execute `go build && ./scraper` for the scraper and `go build && ./article-downloader` for the downloader. The following env variables are available for the postgres connection:

1. PG_HOST
1. PG_DBNAME
1. PG_USERNAME
1. PG_PASSWORD
1. PG_PORT

Adjust by duplicating the scraper-example.env.

