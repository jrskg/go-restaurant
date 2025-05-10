package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/database"
	"github.com/jrskg/go-restaurant/routes"
	"github.com/jrskg/go-restaurant/utils"
)

func main() {
	defer func() {
		database.DisconnectDB(database.DBClient)
	}()
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		utils.ApiSuccess(c, http.StatusOK, nil, "Server health is fine and running")
	})

	routes.UserRoute(router)
	routes.FoodRoute(router)
	routes.InvoiceRoute(router)
	routes.MenuRoute(router)
	routes.OrderItemRoute(router)
	routes.OrderRoute(router)
	routes.TableRoute(router)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
