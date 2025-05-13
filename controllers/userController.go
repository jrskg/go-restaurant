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
	"github.com/jrskg/go-restaurant/helpers"
	"github.com/jrskg/go-restaurant/models"
	"github.com/jrskg/go-restaurant/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var userCollection = database.OpenCollection(database.DBClient, constants.USER_COLLECTION)

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(user); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		result := userCollection.FindOne(ctx, bson.M{"email": user.Email})
		if err := result.Err(); err == nil {
			utils.ApiError(c, http.StatusBadRequest, errors.New("email already exist"))
			return
		} else if err != mongo.ErrNoDocuments {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		token, refreshToken, err := helpers.GetTokens(user.UserId, *user.Name, *user.Email, *user.Avatar)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		hashedPassword, err := utils.HashPassword(*user.Password)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}
		now := time.Now().UTC()
		user.ID = bson.NewObjectID()
		user.CreatedAt = now
		user.UpdatedAt = now
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserId = user.ID.Hex()
		user.Password = &hashedPassword

		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusCreated, user, "User signed up successfully")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var credentials models.LoginDto
		if err := c.BindJSON(&credentials); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(credentials); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&user)
		if err != nil {
			slog.Error("Error while fetching user", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusNotFound, err)
			return
		}

		if !utils.VerifyPassword(*credentials.Password, *user.Password) {
			utils.ApiError(c, http.StatusUnauthorized, errors.New("invalid credentials"))
			return
		}

		token, refreshToken, err := helpers.GetTokens(user.UserId, *user.Name, *user.Email, *user.Avatar)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		updates := bson.M{
			"$set": bson.M{
				"token":        &token,
				"refreshToken": &refreshToken,
				"updatedAt":    time.Now().UTC(),
			},
		}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = userCollection.FindOneAndUpdate(ctx, bson.M{"email": credentials.Email}, updates, opts).Decode(&user)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, user, "User logged in successfully")
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
