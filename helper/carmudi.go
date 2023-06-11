package helper

import (
	"api-cariprice/app/entity"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func SearchCarmudi(query string, locations []entity.Location, filter string) ([]entity.Cars, []entity.Location, int) {
	var cars []entity.Cars
	var foundCarmudi int = 0

	var pagesToScrape []string

	queryLink := "https://www.carmudi.co.id/mobil-bekas-dijual/indonesia?page_size=50&keyword=" + strings.ReplaceAll(query, " ", "-")
	pageToScrape := queryLink + "&page_number=1"

	i := 1
	limit := 5

	pagesDiscovered := []string{pageToScrape}
	var foundCarm []int

	c := colly.NewCollector()

	c.OnHTML("div.masthead", func(e *colly.HTMLElement) {
		re := regexp.MustCompile("[0-9]+")
		regex := re.FindAllString(e.ChildText("h1.headline.delta"), -1)
		found, _ := strconv.Atoi(strings.Join(regex, ""))
		foundCarm = append(foundCarm, found)
	})

	c.OnHTML("ul.pagination li", func(e *colly.HTMLElement) {
		newPaginationLink := queryLink + "/" + e.ChildText("a") + "?" + filter
		if !Contains(pagesToScrape, newPaginationLink) && e.ChildText("a") != "1" {
			// if the page discovered should be scraped
			if !Contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	c.OnHTML("article", func(e *colly.HTMLElement) {
		title := e.Attr("data-title")
		if len(title) > 0 {
			var detail []string
			detail = append(detail, e.Attr("data-year"))
			specs := e.ChildTexts(".listing__specs .item")
			detail = append(detail, specs[:2]...)

			re := regexp.MustCompile("[0-9]+")
			regex := re.FindAllString(e.ChildText("div.listing__price"), -1)
			price, _ := strconv.Atoi(strings.Join(regex, ""))
			exists, index := ContainsCity(locations, specs[2])
			if exists {
				locations[index].Count++
			} else {
				location := entity.Location{
					Label: specs[2],
					Count: 1,
				}
				locations = append(locations, location)
			}
			car := entity.Cars{
				Title:       title,
				Description: e.ChildText("div.listing__excerpt"),
				Price:       price,
				Location:    specs[2],
				Brand:       e.Attr("data-make"),
				Model:       e.Attr("data-model"),
				Detail:      strings.Join(detail, ","),
				SourceName:  "carmudi.co.id",
				IsScraping:  true,
				SourceLink:  e.Attr("data-url"),
			}
			cars = append(cars, car)
		}
	})

	c.OnScraped(func(response *colly.Response) {
		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			// getting the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// incrementing the iteration counter
			i++

			// visiting a new page

			c.Visit(pageToScrape)
		}

	})

	// Set up the callback for errors
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping by visiting the URL
	c.Visit(pageToScrape)

	if len(cars) < 1 {
		return []entity.Cars{}, locations, foundCarmudi
	}

	foundCarmudi = foundCarm[0]

	return cars, locations, foundCarmudi
}
