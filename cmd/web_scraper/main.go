package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/tiberiu1204/olx_web_scraper/internal/scraper"
	"github.com/tiberiu1204/olx_web_scraper/internal/utils"
)

func main() {
	// var numberOfPages uint8 = getNumberOfPages(URL, c)
	c := colly.NewCollector(colly.AllowedDomains("www.olx.ro", "olx.ro"))
	numberOfPages := scraper.GetNumberOfPages(scraper.URL, c)
	var totalPrice float64 = 0
	entries := make(map[string]scraper.Advertisement)
	var wg = sync.WaitGroup{}
	var m = sync.RWMutex{}

	for i := 1; i <= int(numberOfPages); i++ {
		collector := colly.NewCollector(colly.AllowedDomains("www.olx.ro", "olx.ro"))
		go scraper.ScrapePageIndex(uint8(i), collector, entries, &m, &wg)
	}

	wg.Wait()

	for key := range entries {
		adv := entries[key]
		fmt.Printf("Title: %v\nArea: %v\nPrice: %v\nLink: %v\n\n", adv.Title, adv.Area, adv.Price, adv.Href)
		var pricePerHa float32 = float32(adv.Price) / float32(float32(adv.Area)/10000)
		totalPrice += float64(pricePerHa)
	}

	fmt.Printf("Found %v entries.\n", len(entries))
	fmt.Printf("Average price: %v lei / ha.", utils.ToFixed(totalPrice/float64(len(entries)), 2))
}
