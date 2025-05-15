package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
	"github.com/jrskg/go-restaurant/middlewares"
)

func FoodRoute(router *gin.Engine) {
	foodGroup := router.Group("/food")
	foodGroup.Use(middlewares.Authenticate())
	foodGroup.POST("/create", controllers.CreateFood())
	foodGroup.PUT("/:foodId", controllers.UpdateFood())
	foodGroup.DELETE("/:foodId", controllers.DeleteFood())
	foodGroup.GET("/:foodId", controllers.GetFood())
	foodGroup.GET("/all", controllers.GetAllFoods())
}
