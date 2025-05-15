package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
	"github.com/jrskg/go-restaurant/middlewares"
)

func OrderItemRoute(router *gin.Engine) {
	orderItemGroup := router.Group("/order-item")
	orderItemGroup.Use(middlewares.Authenticate())
	orderItemGroup.POST("/create", controllers.CreateOrderItem())
	orderItemGroup.PUT("/:orderItemId", controllers.UpdateOrderItem())
	orderItemGroup.GET("/order/:orderId", controllers.GetOrderItemsByOrder())
	orderItemGroup.GET("/:orderItemId", controllers.GetOrderItem())
	//todo: add more routes
}
