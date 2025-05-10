package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
)

func OrderItemRoute(router *gin.Engine) {
	orderItemGroup := router.Group("/order-item")
	orderItemGroup.POST("/create", controllers.CreateOrderItem())
	orderItemGroup.PUT("/:orderItemId", controllers.UpdateOrderItem())
	//todo: add more routes
}
