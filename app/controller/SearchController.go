package controller

import (
	"net/http"
	"otr-api/app/entity"
	"otr-api/app/service"
	"otr-api/helper"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchController struct {
	service service.OTRServices
}

func NewSearchController(s service.OTRServices) SearchController {
	return SearchController{
		service: s,
	}
}

// @Summary Get Car list
// @Description Get carlist from mobil.cari.co and olx.co.id
// @Author Rasmad Ibnu
// @Success 200 {object} []entity.Car
// @Failure 404 {object} nil
// @method [GET]
// @Router /search
func (controller SearchController) Index(ctx *gin.Context) {
	q := ctx.Param("query")
	param := ctx.Request.URL.Query()

	m := []string{}
	for k, v := range param {
		m = append(m, k+"="+v[0])
	}

	otr, foundDatabase, _ := controller.service.Search(q)
	cars, location, foundCariMobil := helper.SearchCariMobil(q, strings.Join(m, ""))
	ety, locationOLX, foundOLX := helper.SearchOLX(q, location)
	ap := append(otr, ety...)
	ap = append(ap, cars...)

	var metadata entity.FilterOption
	metadata.TotalData = len(ap)
	metadata.Found = foundDatabase + foundCariMobil + foundOLX
	metadata.Location = locationOLX

	resp := helper.SuccessJSON(ctx, "Scraping successfully", http.StatusOK, metadata, ap)

	ctx.JSON(http.StatusOK, resp)

	return
}
