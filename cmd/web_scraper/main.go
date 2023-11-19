package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/tiberiu1204/olx_web_scraper/internal/handlers"
)

func main() {
	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()
	handlers.Handler(r)
	fmt.Println("Starting GO API service...")

	err := http.ListenAndServe("localhost:8000", r)
	if err != nil {
		log.Error(err)
	}
}

// func sampleQuery(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Hit sample-query endpoint")
// 	entriesArray := make([]scraper.Advertisement, 0, 1600)
// 	entries := map[string]scraper.Advertisement{}
// 	queries := []string{"teren agricol", "teren arabil", "teren cultivabil"}

// 	scraper.ScrapeMultipleOlxQueries(queries, entries)

// 	for _, entry := range entries {
// 		entriesArray = append(entriesArray, entry)
// 	}

// 	json.NewEncoder(w).Encode(entriesArray)
// }

// func handleRequests() {
// 	http.HandleFunc("/query", sampleQuery)
// 	http.ListenAndServe(":8081", nil)
// }
