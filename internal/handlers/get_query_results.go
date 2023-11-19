package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tiberiu1204/olx_web_scraper/api"
	"github.com/tiberiu1204/olx_web_scraper/internal/scraper"
)

func GetQueryResults(w http.ResponseWriter, r *http.Request) {
	var body = api.ScrapeQueryBody{}
	var err error = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var queries = body.Queries
	var entries = map[string]scraper.Advertisement{}
	scraper.ScrapeMultipleOlxQueries(queries, entries)
	var filters api.ScrapeQueryFilters = body.Filters
	var advertisements = make([]scraper.Advertisement, 0, 150)

	for _, entry := range entries {
		if applyFilters(entry, filters) {
			advertisements = append(advertisements, entry)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(advertisements)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

func applyFilters(a scraper.Advertisement, f api.ScrapeQueryFilters) bool {
	if f.MaxArea == 0 {
		f.MaxArea = 1<<32 - 1
	}
	if f.MaxPPH == 0 {
		f.MaxPPH = math.Inf(1)
	}
	if f.MaxPrice == 0 {
		f.MaxPrice = 1<<32 - 1
	}
	return (a.Area >= f.MinArea) && (a.Area <= f.MaxArea) &&
		(a.Price >= f.MinPrice) && (a.Price <= f.MaxPrice) &&
		(a.PPH >= f.MinPPH) && (a.PPH <= f.MaxPPH)
}
