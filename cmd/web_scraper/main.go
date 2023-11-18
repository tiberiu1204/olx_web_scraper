package main

import (
	"fmt"

	"github.com/tiberiu1204/olx_web_scraper/internal/scraper"
)

func main() {
	entries := map[string]scraper.Advertisement{}
	scraper.ScrapeOlxQuery("teren agricol", entries)
	scraper.ScrapeOlxQuery("teren arabil", entries)
	scraper.ScrapeOlxQuery("teren cultivabil", entries)
	fmt.Println(len(entries))
}
