package helper

import (
	"api-cariprice/app/entity"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/golang-module/carbon/v2"
	"golang.org/x/exp/slices"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SearchCariMobil(query string, locations []entity.Location, filter string) ([]entity.Cars, []entity.Location, int) {
	var cars []entity.Cars
	var foundCariMobil int = 0

	var pagesToScrape []string
	months := []string{
		"Jan", "Feb", "Mar", "Apr",
		"May", "Jun", "Jul", "Aug",
		"Sep", "Oct", "Nov", "Dec",
	}

	queryLink := "https://mobil.cari.co/mobil/" + strings.ReplaceAll(query, " ", "-")
	pageToScrape := queryLink + "/1?" + filter

	i := 1
	limit := 5

	pagesDiscovered := []string{pageToScrape}

	c := colly.NewCollector()

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

	c.OnHTML("div.three-fifth", func(e *colly.HTMLElement) {
		total, _ := strconv.Atoi(strings.Replace(e.ChildText("h2 span.red"), ".", "", -1))
		foundCariMobil = total
	})

	c.OnHTML("div.item", func(c *colly.HTMLElement) {
		price, _ := strconv.Atoi(c.ChildAttr(".itemRight .price meta[itemprop=price]", "content"))

		if len(c.DOM.Find(".itemTitle").Nodes) > 0 {
			var createdAt string
			date := strings.Split(strings.Split(c.ChildText(".itemDetail .itemDate"), " - ")[0], " ")
			if len(date) > 2 {
				if date[2] != "lalu" {
					month := slices.IndexFunc(months, func(c string) bool { return c == date[1] })
					date := fmt.Sprintf(date[2]+"-"+fmt.Sprintf("%02d", month+1)+"-"+"%02s", date[0])
					ts, _ := time.Parse("2006-01-02", date)
					createdAt = ts.Format("2006-01-02 15:04:05")

				} else {
					hour, _ := strconv.Atoi(date[0])
					carbon := carbon.Parse("now").SubHours(hour).ToDateTimeString()
					createdAt = carbon
				}
			} else {
				if date[1] != "" {
					createdAt = carbon.Parse("yesterday").ToDateString() + " " + date[1] + ":00"

				} else {
					createdAt = carbon.Parse("yesterday").ToDateString() + " 00:00:00"

				}
			}

			var locationLabel string
			locationEl := c.ChildText(".itemRight .itemCity")
			if len(locationEl) > 0 {
				locationLabel = locationEl
			} else {
				locationLabel = "Unknown"
			}
			exists, index := ContainsCity(locations, locationLabel)

			if exists {
				locations[index].Count++
			} else {
				location := entity.Location{
					Label: locationLabel,
					Count: 1,
				}
				locations = append(locations, location)
			}

			car := entity.Cars{
				Thumbnail:       c.ChildAttr("div.picture [itemprop=image]", "content"),
				Title:           c.ChildText(".itemTitle"),
				Price:           price,
				Description:     c.ChildText(".itemDetail p"),
				Brand:           c.ChildText(".itemDesc span[itemprop=brand]"),
				Model:           c.ChildText(".itemDesc span[itemprop=model]"),
				Location:        locationLabel,
				Detail:          strings.Join(c.ChildTexts(".itemDesc span:not([itemprop=brand],[itemprop=model],[itemprop=itemCondition])"), ","),
				Status:          c.ChildText(".itemDesc span[itemprop=itemCondition]"),
				SourceName:      strings.Split(c.ChildText(".itemDetail .itemDate"), " - ")[1],
				SourceLink:      c.ChildAttr(".itemTitle a", "href"),
				CreatedAtString: createdAt,
				IsScraping:      true,
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
		// log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping by visiting the URL
	c.Visit(pageToScrape)

	if len(cars) < 1 {
		return []entity.Cars{}, locations, foundCariMobil
	}

	return cars, locations, foundCariMobil
}
