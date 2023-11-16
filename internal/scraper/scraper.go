package scraper

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gocolly/colly"
	"github.com/tiberiu1204/olx_web_scraper/internal/utils"
)

const URL string = "https://www.olx.ro/oferte/q-teren-agricol/"

type Advertisement struct {
	Title string
	Area  uint32
	Price uint32
	Href  string
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

func ScrapePageIndex(index uint8, c *colly.Collector, entries map[string]Advertisement, m *sync.RWMutex, wg *sync.WaitGroup) {
	var pageUrl string = URL + "?page=" + strconv.Itoa(int(index))
	c.OnHTML("div.css-1sw7q4x", func(h *colly.HTMLElement) {
		selection := h.DOM
		title := selection.Find("h6").Text()
		area, areaErr := utils.GetNumberFromString(selection.Find("span.css-643j0o").Text())
		price, priceErr := utils.GetNumberFromString(selection.Find("p.css-10b0gli.er34gjf0").Text())
		href := h.ChildAttr("a", "href")
		var pricePerHa float32 = float32(price) / float32(area/10000)

		if len(href) > 0 && href[0] != 'h' {
			href = "https://www.olx.ro" + href
		}

		if (areaErr == nil && priceErr == nil) && (area >= 5000) && (pricePerHa >= 15000 && pricePerHa <= 5*20000) {
			adv := Advertisement{Title: title, Area: uint32(area), Price: uint32(price), Href: href}
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
