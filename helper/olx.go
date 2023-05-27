package helper

import (
	"api-cariprice/app/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type OLXResponse struct {
	Version  string `json:"version"`
	MetaData struct {
		TotalAds   int `json:"total_ads"`
		TotalPages int `json:"total_pages"`
	} `json:"metadata"`
	Data []Data `json:"data"`
}

type Parameters struct {
	Key   string `json:"key_name"`
	Value string `json:"formatted_value"`
}

type Data struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Price struct {
		Value struct {
			Raw     interface{} `json:"raw"`
			Display string      `json:"display"`
		} `json:"value"`
		Key string `json:"key"`
	} `json:"price"`
	Description string `json:"description"`
	Images      []struct {
		Url string `json:"url"`
	} `json:"images"`
	Location struct {
		Country  string `json:"COUNTRY_name"`
		Province string `json:"ADMIN_LEVEL_1_name"`
		City     string `json:"ADMIN_LEVEL_3_name"`
	} `json:"locations_resolved"`
	Parameters []Parameters `json:"parameters"`
	CreatedAt  string       `json:"display_date"`
}

func ContainsCity(s []entity.Location, e string) (bool, int) {
	for idx, v := range s {
		if v.Label == e {
			return true, idx
		}
	}
	return false, -1
}

func SearchOLX(query string, locations []entity.Location) ([]entity.Cars, []entity.Location, int) {
	link := "https://www.olx.co.id/api/relevance/v4/search?category=198&facet_limit=100&location=1000001&location_facet_limit=20&platform=web-desktop&query=" + url.QueryEscape(query) + "&relaxedFilters=true&size=40&spellcheck=true"
	response, err := http.Get(link + "&page=0")

	if err != nil {
		fmt.Print(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responseObject OLXResponse
	json.Unmarshal(responseData, &responseObject)

	var data []Data

	data = append(data, responseObject.Data...)

	for page := 1; page < 3; page++ {
		response, err := http.Get(link + "&page=" + fmt.Sprintf("%d", page))

		if err != nil {
			fmt.Print(err)
		}

		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Fatal(err)
		}

		var responseObject OLXResponse
		json.Unmarshal(responseData, &responseObject)
		data = append(data, responseObject.Data...)
	}

	var cars []entity.Cars
	for _, col := range data {
		brand := slices.IndexFunc(col.Parameters, func(c Parameters) bool { return c.Key == "Merek" })
		model := slices.IndexFunc(col.Parameters, func(c Parameters) bool { return c.Key == "Model" })

		var detail []string

		for idx := 0; idx < len(col.Parameters); idx++ {
			if idx == brand || idx == model {
			} else {
				detail = append(detail, col.Parameters[idx].Value)
			}
		}
		price := strings.Split(fmt.Sprintf("%f", col.Price.Value.Raw), ".")
		priceInt, _ := strconv.Atoi(price[0])
		location := col.Location.Province + ", " + col.Location.City
		ts, _ := time.Parse("2006-01-02T15:04:05-0700", col.CreatedAt)
		exists, index := ContainsCity(locations, location)
		if exists {
			locations[index].Count++
		} else {
			location := entity.Location{
				Label: location,
				Count: 1,
			}
			locations = append(locations, location)
		}

		car := entity.Cars{
			Thumbnail:       col.Images[0].Url,
			Title:           col.Title,
			Description:     col.Description,
			Price:           priceInt,
			Brand:           col.Parameters[brand].Value,
			Model:           col.Parameters[model].Value,
			Location:        location,
			Status:          "Bekas",
			SourceName:      "olx.co.id",
			SourceLink:      "https://www.olx.co.id/item/" + col.ID,
			Detail:          strings.Join(detail, ","),
			CreatedAtString: ts.Format("2006-01-02 15:04:05"),
			IsScraping:      true,
		}
		cars = append(cars, car)
	}

	return cars, locations, responseObject.MetaData.TotalAds
}
