package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/controllers"
	"github.com/jrskg/go-restaurant/middlewares"
)

func InvoiceRoute(router *gin.Engine) {
	invoiceGroup := router.Group("/invoice")
	invoiceGroup.Use(middlewares.Authenticate())
	invoiceGroup.POST("/create", controllers.CreateInvoice())
	invoiceGroup.PUT("/:invoiceId", controllers.UpdateInvoice())
	// invoiceGroup.DELETE("/:invoiceId", controllers.DeleteInvoice())
	invoiceGroup.GET("/:invoiceId", controllers.GetInvoice())
	invoiceGroup.GET("/all", controllers.GetAllInvoices())
}
