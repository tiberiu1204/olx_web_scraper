package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gocolly/colly"
)

const URL string = "https://www.olx.ro/oferte/q-teren-agricol/"

type Advertisement struct {
	title string
	area  uint32
	price uint32
	href  string
}

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.olx.ro", "olx.ro"))
	// var numberOfPages uint8 = getNumberOfPages(URL, c)

	numberOfPages := getNumberOfPages(URL, c)
	numberOfEntries := 0
	var totalPrice float64 = 0

	for i := 1; i <= int(numberOfPages); i++ {
		advertisements := scrapePageIndex(uint8(i), c)
		for index := range advertisements {
			adv := advertisements[index]
			fmt.Printf("Title: %v\nArea: %v\nPrice: %v\nLink: %v\n\n", adv.title, adv.area, adv.price, adv.href)
			var pricePerHa float32 = float32(adv.price) / float32(float32(adv.area)/10000)
			totalPrice += float64(pricePerHa)
		}
		numberOfEntries += len(advertisements)
	}
	fmt.Printf("Found %v entries.\n", numberOfEntries)
	fmt.Printf("Average price: %v lei / ha.", toFixed(totalPrice/float64(numberOfEntries), 2))
}

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

func getNumberFromString(str string) (int32, error) {
	var number int32 = 0
	var err error
	ZERO, NINE := '0', '9'
	for index := range str {
		char := str[index]
		if ZERO <= rune(char) && rune(char) <= NINE {
			number = number*10 + (rune(char) - ZERO)
		}
	}
	if number == 0 {
		err = errors.New("")
	}
	return number, err
}

func scrapePageIndex(index uint8, c *colly.Collector) []Advertisement {
	var advertisements []Advertisement = make([]Advertisement, 0, 10)

	var pageUrl string = URL + "?page=" + strconv.Itoa(int(index))

	c.OnHTML("div.css-1sw7q4x", func(h *colly.HTMLElement) {
		selection := h.DOM
		title := selection.Find("h6").Text()
		area, areaErr := getNumberFromString(selection.Find("span.css-643j0o").Text())
		price, priceErr := getNumberFromString(selection.Find("p.css-10b0gli.er34gjf0").Text())
		href := h.ChildAttr("a", "href")
		var pricePerHa float32 = float32(price) / float32(area/10000)

		if (areaErr == nil && priceErr == nil) && (area >= 5000) && (pricePerHa >= 15000 && pricePerHa <= 5*15000) {
			advertisements = append(advertisements, Advertisement{title: title, area: uint32(area), price: uint32(price), href: href})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Erorr while scraping: %v\n", err.Error())
	})

	c.Visit(pageUrl)

	return advertisements
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
