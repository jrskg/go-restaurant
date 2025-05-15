package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
	"github.com/jrskg/go-restaurant/middlewares"
)

func TableRoute(router *gin.Engine) {
	tableGroup := router.Group("/table")
	tableGroup.Use(middlewares.Authenticate())
	tableGroup.POST("/create", controllers.CreateTable())
	tableGroup.PUT("/:tableId", controllers.UpdateTable())
	tableGroup.GET("/:tableId", controllers.GetTable())
	tableGroup.GET("/all", controllers.GetAllTables())
}
