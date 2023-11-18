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

type Advertisement struct {
	Title string
	Area  uint32
	Price uint32
	PPH   float64
	Href  string
}

func (adv Advertisement) Print() {
	fmt.Printf("Title: %v\nArea: %v\nPrice: %v\nPrice / ha: %v lei / ha\nLink: %v\n\n", adv.Title, adv.Area, adv.Price, adv.PPH, adv.Href)
}

func GetNumberOfPages(url string, c *colly.Collector) uint8 {
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

func ScrapePageIndex(url string, index uint8, c *colly.Collector, entries map[string]Advertisement, m *sync.RWMutex, wg *sync.WaitGroup) {
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

		if (areaErr == nil && priceErr == nil) && (area >= 5000) && (pricePerHa >= 15000 && pricePerHa <= 5*20000) {
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
	numberOfPages := GetNumberOfPages(url, c)
	wg := sync.WaitGroup{}

	fmt.Printf("Scraping %v ...\n", url)

	for i := 1; i <= int(numberOfPages); i++ {
		collector := colly.NewCollector(colly.AllowedDomains("www.olx.ro", "olx.ro"))
		go ScrapePageIndex(url, uint8(i), collector, entries, m, &wg)
	}

	wg.Wait()

	m.Lock()
	fmt.Printf("\nFinished scraping url %v\n", url)
	fmt.Printf("Time elapsed: %v\n", time.Since(t0))
	m.Unlock()

	wg1.Done()
}

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
