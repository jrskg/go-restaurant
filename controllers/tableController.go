package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/constants"
	"github.com/jrskg/go-restaurant/database"
)

var tableCollection = database.OpenCollection(database.DBClient, constants.TABLE_COLLECTION)

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetAllTables() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
