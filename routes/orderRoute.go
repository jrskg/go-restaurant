package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
)

func OrderRoute(router *gin.Engine) {
	orderGroup := router.Group("/order")
	orderGroup.POST("/create", controllers.CreateOrder())
	orderGroup.PUT("/:orderId", controllers.UpdateOrder())
	orderGroup.DELETE("/:orderId", controllers.DeleteOrder())
	orderGroup.GET("/:orderId", controllers.GetOrder())
	orderGroup.GET("/all", controllers.GetAllOrders())
}
