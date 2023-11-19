package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/tiberiu1204/olx_web_scraper/api"
	"github.com/tiberiu1204/olx_web_scraper/internal/scraper"
	"github.com/tiberiu1204/olx_web_scraper/internal/utils"
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
	var averagePPH float64 = 0
	var medianPPH float64 = 0

	for _, entry := range entries {
		if applyFilters(entry, filters) {
			advertisements = append(advertisements, entry)
			averagePPH += entry.PPH
		}
	}
	sort.Slice(advertisements, func(i, j int) bool {
		return advertisements[i].PPH < advertisements[j].PPH
	})

	var n = len(advertisements)
	averagePPH /= float64(n)
	if n&1 == 0 {
		medianPPH = (advertisements[n/2-1].PPH + advertisements[n/2].PPH) / 2
	} else {
		medianPPH = advertisements[n/2+1].PPH
	}

	var response api.ScrapeQueriesResponse = api.ScrapeQueriesResponse{
		Code:         http.StatusOK,
		Entries:      advertisements,
		MinPPHEntry:  advertisements[0],
		MaxPPHEntry:  advertisements[len(advertisements)-1],
		AveragePPH:   utils.ToFixed(averagePPH, 2),
		MedianPPH:    utils.ToFixed(medianPPH, 2),
		EntriesCount: uint32(len(advertisements)),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
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
