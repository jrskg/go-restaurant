package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
	"github.com/jrskg/go-restaurant/middlewares"
)

func OrderRoute(router *gin.Engine) {
	orderGroup := router.Group("/order")
	orderGroup.Use(middlewares.Authenticate())
	orderGroup.POST("/create", controllers.CreateOrder())
	orderGroup.PUT("/:orderId", controllers.UpdateOrder())
	orderGroup.DELETE("/:orderId", controllers.DeleteOrder())
	orderGroup.GET("/:orderId", controllers.GetOrder())
	orderGroup.GET("/all", controllers.GetAllOrders())
}
