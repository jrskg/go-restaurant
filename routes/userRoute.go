package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
)

func UserRoute(router *gin.Engine) {
	userGroup := router.Group("/user")
	userGroup.POST("/signup", controllers.Signup())
	userGroup.POST("/login", controllers.Login())
	userGroup.GET("/logout", controllers.Logout())

	userGroup.GET("/:userId", controllers.GetUser())
	userGroup.GET("/all", controllers.GetAllUsers())
}
