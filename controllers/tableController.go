package controllers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/constants"
	"github.com/jrskg/go-restaurant/database"
	"github.com/jrskg/go-restaurant/models"
	"github.com/jrskg/go-restaurant/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var tableCollection = database.OpenCollection(database.DBClient, constants.TABLE_COLLECTION)

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		table := models.Table{}

		if err := c.BindJSON(&table); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(table); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		table.ID = bson.NewObjectID()
		table.TableId = table.ID.Hex()
		table.CreatedAt = time.Now()
		table.UpdatedAt = time.Now()

		_, err := tableCollection.InsertOne(ctx, table)
		if err != nil {
			slog.Error("Error while creating table", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusCreated, table, "Table created successfully")
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tableId := c.Param("tableId")
		if tableId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("tableId is empty"))
			return
		}

		updateTableDto := models.UpdateTableDto{}
		if err := c.BindJSON(&updateTableDto); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		updateObj := bson.M{"updatedAt": time.Now().UTC()}
		if updateTableDto.NumberOfGuests != nil {
			updateObj["numberOfGuests"] = updateTableDto.NumberOfGuests
		}
		if updateTableDto.TableNumber != nil {
			updateObj["tableNumber"] = updateTableDto.TableNumber
		}

		filter := bson.M{"tableId": tableId}
		result, err := tableCollection.UpdateOne(ctx, filter, bson.M{"$set": updateObj})

		if err != nil {
			slog.Error("Error while updating table", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if result.MatchedCount == 0 {
			utils.ApiError(c, http.StatusNotFound, errors.New("table not found"))
			return
		}

		utils.ApiSuccess(c, http.StatusOK, nil, "Table updated successfully")
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tableId := c.Param("tableId")
		if tableId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("tableId is empty"))
			return
		}
		table := models.Table{}

		err := tableCollection.FindOne(ctx, bson.M{"tableId": tableId}).Decode(&table)
		if err != nil {
			slog.Error("Error while fetching table", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusNotFound, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, table, "Table fetched successfully")
	}
}

func GetAllTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tables := make([]models.Table, 0)
		result, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			slog.Error("Error while fetching tables", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if err := result.All(ctx, &tables); err != nil {
			slog.Error("Error while fetching tables", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, tables, "Tables fetched successfully")
	}
}
