package main

import (
	"fmt"

	"github.com/tiberiu1204/olx_web_scraper/internal/scraper"
)

func main() {
	entries := map[string]scraper.Advertisement{}
	queries := []string{"teren agricol", "teren arabil", "teren cultivabil"}
	scraper.ScrapeMultipleOlxQueries(queries, entries)
	fmt.Printf("\n%v", len(entries))
}
