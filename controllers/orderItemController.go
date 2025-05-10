package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/constants"
	"github.com/jrskg/go-restaurant/database"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var orderItemCollection = database.OpenCollection(database.DBClient, constants.ORDER_ITEM_COLLECTION)

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ItemsByOrder(id string) (OrderItems []bson.M, err error) {
	var ok []bson.M = nil
	return ok, nil
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
