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
	"github.com/jrskg/go-restaurant/models"
	"github.com/jrskg/go-restaurant/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(menu); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		menu.CreatedAt = time.Now().UTC()
		menu.UpdatedAt = time.Now().UTC()
		menu.ID = bson.NewObjectID()
		menu.MenuId = menu.ID.Hex()

		_, err := menuCollection.InsertOne(ctx, menu)
		if err != nil {
			slog.Error("Error while creating menu", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusCreated, menu, "Menu created successfully")
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var updateDto models.MenuUpdateDto
		if err := c.BindJSON(&updateDto); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(updateDto); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		menuId := c.Param("menuId")
		if menuId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("menuId is empty"))
			return
		}

		if updateDto.StartDate != nil && updateDto.EndDate != nil {
			if !inTimeSpan(*updateDto.StartDate, *updateDto.EndDate, time.Now()) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("start date and end date must be in the future"))
				return
			}
		}

		filter := bson.M{"menuId": menuId}
		updateFields := bson.M{
			"name":      updateDto.Name,
			"category":  updateDto.Category,
			"startDate": updateDto.StartDate,
			"endDate":   updateDto.EndDate,
		}
		updateObj := bson.M{"updatedAt": time.Now().UTC()}
		for k, v := range updateFields {
			if v != nil {
				updateObj[k] = v
			}
		}
		result, err := menuCollection.UpdateOne(ctx, filter, bson.M{"$set": updateObj})
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}
		if result.MatchedCount < 1 {
			utils.ApiError(c, http.StatusNotFound, errors.New("menu not found"))
			return
		}
		utils.ApiSuccess(c, http.StatusOK, updateDto, "Menu updated successfully")
	}
}

func DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		menuId := c.Param("menuId")
		if menuId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid menu id"))
			return
		}

		result, err := menuCollection.DeleteOne(ctx, bson.M{"menuId": menuId})
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if result.DeletedCount < 1 {
			utils.ApiError(c, http.StatusNotFound, errors.New("menu not found"))
			return
		}

		utils.ApiSuccess(c, http.StatusOK, nil, "Menu deleted successfully")
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		menuId := c.Param("menuId")
		if menuId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("menuId is empty"))
			return
		}

		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menuId": menuId}).Decode(&menu)
		if err != nil {
			slog.Error("Error while fetching menu", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusNotFound, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, menu, "Menu fetched successfully")
	}
}

func GetAllMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var allMenus []models.Menu = make([]models.Menu, 0)

		result, err := menuCollection.Find(ctx, bson.M{})
		if err != nil {
			slog.Error("Error while fetching menu", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if err := result.All(ctx, &allMenus); err != nil {
			slog.Error("Error while fetching menu", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}
		fmt.Println(allMenus)
		utils.ApiSuccess(c, http.StatusOK, allMenus, "Menus fetched successfully")
	}
}
