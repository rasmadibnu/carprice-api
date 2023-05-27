package routes

import (
	"api-cariprice/app/controller"
	"api-cariprice/app/repository"
	"api-cariprice/app/service"
	"api-cariprice/config"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func WebRouter(db config.Database) {
	// Repository Asset
	otrRepository := repository.NewOTRRepository(db)

	// Service Asset
	otrService := service.NewOTRServices(otrRepository)

	//Controller Asset
	searchController := controller.NewSearchController(otrService)
	otrController := controller.NewOTRController(otrService)

	// Route
	httpRouter := gin.Default()

	// Register routing
	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Testing  connection
	httpRouter.GET("/status-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"data": "Service ✅ API Up and Running"})
	})

	v1 := httpRouter.Group("/api/v1") // Grouping routes

	v1.GET("/search/:query", searchController.Index)

	v1.GET("/otr", otrController.Index)
	v1.POST("/otr", otrController.Store)
	v1.DELETE("/otr/:id", otrController.Delete)

	httpRouter.Run(":" + os.Getenv("APP_PORT")) // Run Routes with PORT
}
