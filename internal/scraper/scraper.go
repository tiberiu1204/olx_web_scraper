package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/tiberiu1204/olx_web_scraper/internal/utils"
)

// This structure holds all the relevant information that will be scraped from each entry

type Advertisement struct {
	Title string  `json:"title"`            // The title of the advert
	Area  uint32  `json:"area"`             // The area of the land advertised, in m^2
	Price uint32  `json:"price"`            // The price of the advertised land, in RON
	PPH   float64 `json:"price_per_hectar"` // The RON per ha of the advertised land
	Href  string  `json:"href"`             // The link the advert points to
}

// This function prints an Advertisement struct object

func (adv Advertisement) Print() {
	fmt.Printf("Title: %v\nArea: %v\nPrice: %v\nPrice / ha: %v lei / ha\nLink: %v\n\n", adv.Title, adv.Area, adv.Price, adv.PPH, adv.Href)
}

// This function takes in a query url and a colly.Collector ponter
// and returns the number of pages the query contains

func getNumberOfPages(url string, c *colly.Collector) uint8 {
	var numberOfPages uint8 = 0
	c.OnHTML("li.pagination-item a.css-1mi714g", func(h *colly.HTMLElement) {
		number, err := strconv.Atoi(h.Text)
		if err == nil {
			numberOfPages = max(numberOfPages, uint8(number))
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Erorr while scraping: %v\n", err.Error())
	})

	c.Visit(url)
	return numberOfPages
}

// This function takes in a query url, the index of the page, a colly.Collector pointer, a map representing scraped entries,
// a RWMutex pointer and a WaitGroup pointer. The function will store scraped Advertisements in the entries map,
// where the keys represent the links asociated with the Advertisement struct value

func scrapePageIndex(url string, index uint8, c *colly.Collector, entries map[string]Advertisement, m *sync.RWMutex, wg *sync.WaitGroup) {
	var pageUrl string = url + "?page=" + strconv.Itoa(int(index))
	c.OnHTML("div.css-1sw7q4x", func(h *colly.HTMLElement) {
		selection := h.DOM
		title := selection.Find("h6").Text()
		area, areaErr := utils.GetNumberFromString(selection.Find("span.css-643j0o").Text())
		price, priceErr := utils.GetNumberFromString(selection.Find("p.css-10b0gli.er34gjf0").Text())
		href := h.ChildAttr("a", "href")
		pricePerHa := utils.PricePerHa(uint32(price), uint32(area))

		if len(href) > 0 && href[0] != 'h' {
			href = "https://www.olx.ro" + href
		}

		if areaErr == nil && priceErr == nil {
			adv := Advertisement{Title: title, Area: uint32(area), Price: uint32(price), PPH: pricePerHa, Href: href}
			m.Lock()
			entries[href] = adv
			m.Unlock()
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Erorr while scraping: %v\n", err.Error())
	})
	wg.Add(1)
	c.Visit(pageUrl)
	wg.Done()
}

// This function takes a query where each word is separated by a space ' '
// and a map where scraped entries related to land from the website olx.ro should be saved
// where the keys represent a link and the values represent an Advertisement struct

func ScrapeOlxQuery(query string, entries map[string]Advertisement, m *sync.RWMutex, wg1 *sync.WaitGroup) {
	t0 := time.Now()
	url := "https://www.olx.ro/oferte/q-" + strings.Join(strings.Split(query, " "), "-")
	c := colly.NewCollector(colly.AllowedDomains("www.olx.ro", "olx.ro"))
	numberOfPages := getNumberOfPages(url, c)
	wg := sync.WaitGroup{}

	fmt.Printf("Scraping %v ...\n", url)

	for i := 1; i <= int(numberOfPages); i++ {
		collector := colly.NewCollector(colly.AllowedDomains("www.olx.ro", "olx.ro"))
		go scrapePageIndex(url, uint8(i), collector, entries, m, &wg)
	}

	wg.Wait()

	m.Lock()
	fmt.Printf("\nFinished scraping url %v\n", url)
	fmt.Printf("Time elapsed: %v\n", time.Since(t0))
	m.Unlock()

	wg1.Done()
}

// This function takes an array of quries, where each query is a string of words separated by a space ' ',
// a representing the scraped entries, where the key is a string representing a link and the value is
// its corresponding Advertisement struct

func ScrapeMultipleOlxQueries(quries []string, entries map[string]Advertisement) {
	t0 := time.Now()
	m := sync.RWMutex{}
	wg := sync.WaitGroup{}

	for _, query := range quries {
		wg.Add(1)
		go ScrapeOlxQuery(query, entries, &m, &wg)
	}

	wg.Wait()

	fmt.Printf("\nFinished scraping all %v queries\n", len(quries))
	fmt.Printf("Found %v entries in %v", len(entries), time.Since(t0))
}
