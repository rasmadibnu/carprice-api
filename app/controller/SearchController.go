package controller

import (
	"api-cariprice/app/entity"
	"api-cariprice/app/service"
	"api-cariprice/helper"
	"net/http"
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

	otr, foundDatabase, locationDB, _ := controller.service.Search(q)
	cars, locationCarmudi, foundCariMobil := helper.SearchCariMobil(q, locationDB, strings.Join(m, ""))
	olx, locationOLX, foundOLX := helper.SearchOLX(q, locationCarmudi)
	carmudi, locationCarmudi, foundCarmudi := helper.SearchCarmudi(q, locationOLX, "")
	ap := append(otr, olx...)
	ap2 := append(ap, cars...)
	ap3 := append(ap2, carmudi...)

	var metadata entity.FilterOption
	metadata.TotalData = len(ap3)
	metadata.Found = foundDatabase + foundCariMobil + foundOLX + foundCarmudi
	metadata.Location = locationCarmudi

	resp := helper.SuccessJSON(ctx, "Scraping successfully", http.StatusOK, metadata, ap3)

	ctx.JSON(http.StatusOK, resp)

	return
}
