package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
)

func MenuRoute(router *gin.Engine) {
	menuGroup := router.Group("/menu")
	menuGroup.POST("/create", controllers.CreateMenu())
	menuGroup.PUT("/:menuId", controllers.UpdateMenu())
	menuGroup.DELETE("/:menuId", controllers.DeleteMenu())
	menuGroup.GET("/:menuId", controllers.GetMenu())
	menuGroup.GET("/all", controllers.GetAllMenus())
}
