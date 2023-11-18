package api

import (
	"github.com/tiberiu1204/olx_web_scraper/internal/scraper"
)

type ScrapeQueryParams struct {
	Queries []string // The array of search queries that should be scraped
}

type ScrapeQueriesResponse struct {
	Code    int                     // Success Code, 200
	Entries []scraper.Advertisement // Scraped advertisements
}

type Error struct {
	Code    int    // Error code
	Message string // Error message
}
