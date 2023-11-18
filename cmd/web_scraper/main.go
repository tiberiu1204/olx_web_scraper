package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/tiberiu1204/olx_web_scraper/internal/scraper"
)

func main() {
	handleRequests()
}

func sampleQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit sample-query endpoint")
	entriesArray := make([]scraper.Advertisement, 0, 1600)
	entries := map[string]scraper.Advertisement{}
	queries := []string{"teren agricol", "teren arabil", "teren cultivabil"}

	scraper.ScrapeMultipleOlxQueries(queries, entries)

	for _, entry := range entries {
		entriesArray = append(entriesArray, entry)
	}

	json.NewEncoder(w).Encode(entriesArray)
}

func handleRequests() {
	http.HandleFunc("/query", sampleQuery)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
