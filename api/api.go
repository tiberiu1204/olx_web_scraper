package api

import (
	"encoding/json"
	"net/http"

	"github.com/tiberiu1204/olx_web_scraper/internal/scraper"
)

// This structure holds the parameters for a scrape query

type ScrapeQueryFilters struct {
	MinPrice uint32  // Minimum price
	MaxPrice uint32  // Maximum price
	MinArea  uint32  // Minimum area
	MaxArea  uint32  // Maximum area
	MinPPH   float64 // Minimum price per hectar
	MaxPPH   float64 // Maximum price per hectar
}

// This structure holds the query paramenters

type ScrapeQueryParams struct {
	Username []string // The username of the caller

}

type ScrapeQueryBody struct {
	Queries []string           // The array of search queries that should be scraped
	Filters ScrapeQueryFilters // The filters applied to the search query
}

// This structure holds the query response

type ScrapeQueriesResponse struct {
	Code    int                     // Success Code, 200
	Entries []scraper.Advertisement // Scraped advertisements
}

// This structure hold request error information

type Error struct {
	Code    int    // Error code
	Message string // Error message
}

func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An Unexpexted Error Occurred.", http.StatusInternalServerError)
	}
)
