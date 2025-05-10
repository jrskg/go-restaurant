package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
)

func TableRoute(router *gin.Engine) {
	tableGroup := router.Group("/table")
	tableGroup.POST("/create", controllers.CreateTable())
	tableGroup.PUT("/:tableId", controllers.UpdateTable())
	tableGroup.DELETE("/:tableId", controllers.DeleteTable())
	tableGroup.GET("/:tableId", controllers.GetTable())
	tableGroup.GET("/all", controllers.GetAllTables())
}
