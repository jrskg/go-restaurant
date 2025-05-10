package controllers

import (
	"context"
	"errors"
	"fmt"
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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var foodCollection = database.OpenCollection(database.DBClient, constants.FOOD_COLLECTION)
var menuCollection = database.OpenCollection(database.DBClient, constants.MENU_COLLECTION)

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(food); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}
		count, err := menuCollection.CountDocuments(ctx, bson.M{"menuId": food.MenuId})
		if err != nil || count < 1 {
			utils.ApiError(c, http.StatusBadRequest, errors.New("menu not found"))
			return
		}

		food.CreatedAt = time.Now().UTC()
		food.UpdatedAt = time.Now().UTC()
		food.ID = bson.NewObjectID()
		food.FoodId = food.ID.Hex()
		var num = utils.ToFixed(*food.Price, 2)
		food.Price = &num
		result, insertErr := foodCollection.InsertOne(ctx, food)

		if insertErr != nil || result.InsertedID == "" {
			utils.ApiError(c, http.StatusInternalServerError, insertErr)
			return
		}
		utils.ApiSuccess(c, http.StatusCreated, food, "Food created successfully")
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foodId := c.Param("foodId")
		if foodId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("foodId is empty"))
			return
		}

		var updateFoodDto models.UpdateFoodDto

		if err := c.BindJSON(&updateFoodDto); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(updateFoodDto); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if updateFoodDto.MenuId != nil {
			count, err := menuCollection.CountDocuments(ctx, bson.M{"menuId": updateFoodDto.MenuId})
			if err != nil || count < 1 {
				utils.ApiError(c, http.StatusBadRequest, errors.New("invalid menuId"))
				return
			}
		}

		filter := bson.M{"foodId": foodId}
		updateFields := bson.M{
			"name":      updateFoodDto.Name,
			"price":     updateFoodDto.Price,
			"foodImage": updateFoodDto.FoodImage,
			"menuId":    updateFoodDto.MenuId,
		}
		updateObj := bson.M{"updatedAt": time.Now().UTC()}

		for k, v := range updateFields {
			if v != nil {
				updateObj[k] = v
			}
		}

		result, err := foodCollection.UpdateOne(ctx, filter, bson.M{"$set": updateObj})
		if err != nil {
			slog.Error("Error while updating food", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if result.MatchedCount == 0 {
			utils.ApiError(c, http.StatusNotFound, errors.New("food not found"))
			return
		}

		utils.ApiSuccess(c, http.StatusOK, nil, "Food updated successfully")
	}
}

func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foodId := c.Param("foodId")
		if foodId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid food id"))
			return
		}

		result, err := foodCollection.DeleteOne(ctx, bson.M{"foodId": foodId})
		if err != nil {
			slog.Error("Error while deleting food", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if result.DeletedCount == 0 {
			utils.ApiError(c, http.StatusNotFound, errors.New("food not found"))
			return
		}

		utils.ApiSuccess(c, http.StatusOK, nil, "Food deleted successfully")
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		foodId := c.Param("foodId")
		var food models.Food
		err := foodCollection.FindOne(ctx, bson.M{"foodId": foodId}).Decode(&food)
		if err != nil {
			slog.Error("Error while fetching food", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusNotFound, err)
			return
		}
		utils.ApiSuccess(c, http.StatusOK, food, "Food fetched successfully")
	}
}

func GetAllFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		const limit int64 = 50
		cursor := c.Query("cursor")
		var cursorTime time.Time
		if cursor != "" {
			pt, err := utils.ValidateAndParseTime(cursor)
			if err != nil {
				utils.ApiError(c, http.StatusBadRequest, fmt.Errorf("invalid cursor format: %s", cursor))
				return
			}
			cursorTime = pt
		}

		var allFoods []models.Food = make([]models.Food, 0)

		filter := bson.M{}
		if cursor != "" {
			filter = bson.M{"createdAt": bson.M{"$gt": cursorTime}}
		}
		options := options.Find().SetLimit(limit + 1)
		result, err := foodCollection.Find(ctx, filter, options)
		if err != nil {
			slog.Error("Error while fetching food", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if err := result.All(ctx, &allFoods); err != nil {
			slog.Error("Error while saving to allFoods", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}
		hasMore := false
		nextCursor := time.Time{}
		if len(allFoods) > int(limit) {
			hasMore = true
			allFoods = allFoods[:len(allFoods)-1]
			nextCursor = allFoods[len(allFoods)-1].CreatedAt
		}
		utils.ApiSuccess(
			c,
			http.StatusOK,
			bson.M{
				"hasMore":    hasMore,
				"foods":      allFoods,
				"nextCursor": nextCursor.UTC(),
			},
			"Foods fetched successfully",
		)
	}
}
