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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

var orderCollection = database.OpenCollection(database.DBClient, constants.ORDER_COLLECTION)

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(order); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		var table models.Table
		err := tableCollection.FindOne(ctx, bson.M{"tableId": order.TableId}).Decode(&table)
		if err != nil {
			slog.Error("Error while fetching table", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusNotFound, err)
			return
		}

		if !order.OrderDate.After(time.Now()) {
			utils.ApiError(c, http.StatusBadRequest, errors.New("order date must be in future"))
			return
		}

		order.CreatedAt = time.Now().UTC()
		order.UpdatedAt = time.Now().UTC()
		order.ID = bson.NewObjectID()
		order.OrderID = order.ID.Hex()
		order.OrderDate = order.OrderDate.UTC()

		_, err = orderCollection.InsertOne(ctx, order)
		if err != nil {
			slog.Error("Error while creating order", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusCreated, order, "Order created successfully")
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderId := c.Param("orderId")
		if orderId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid order id"))
			return
		}

		var updateOrderDto models.UpdateOrderDto
		if err := c.BindJSON(&updateOrderDto); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(updateOrderDto); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if updateOrderDto.TableId == nil {
			utils.ApiError(c, http.StatusBadRequest, errors.New("tableId is empty"))
			return
		}

		count, err := tableCollection.CountDocuments(ctx, bson.M{"tableId": updateOrderDto.TableId})
		if err != nil || count < 1 {
			utils.ApiError(c, http.StatusBadRequest, errors.New("table not found"))
			return
		}

		filter := bson.M{"orderId": orderId}
		update := bson.M{"tableId": updateOrderDto.TableId, "updatedAt": time.Now().UTC()}

		result, err := orderCollection.UpdateOne(ctx, filter, bson.M{"$set": update})
		if err != nil {
			slog.Error("Error while updating order", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if result.ModifiedCount == 0 {
			utils.ApiError(c, http.StatusNotFound, errors.New("order not found"))
			return
		}

		utils.ApiSuccess(c, http.StatusOK, nil, "Order updated successfully")
	}
}

func DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderId := c.Param("orderId")
		if orderId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid order id"))
			return
		}

		wc := writeconcern.Majority()
		txnOptions := options.Transaction().SetWriteConcern(wc)

		session, err := database.DBClient.StartSession()
		if err != nil {
			slog.Error("Error while starting session", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}
		defer session.EndSession(context.Background())
		callback := func(ctx context.Context) (any, error) {
			_, err := orderCollection.DeleteOne(ctx, bson.M{"orderId": orderId})
			if err != nil {
				slog.Error("Error while deleting order", slog.String("error", err.Error()))
				return nil, err
			}

			_, err = orderItemCollection.DeleteMany(ctx, bson.M{"orderId": orderId})
			if err != nil {
				slog.Error("Error while deleting order items", slog.String("error", err.Error()))
				return nil, err
			}
			return nil, nil
		}

		_, err = session.WithTransaction(ctx, callback, txnOptions)
		if err != nil {
			slog.Error("Error while deleting order", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, nil, "Order deleted successfully")
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, candel := context.WithTimeout(context.Background(), 10*time.Second)
		defer candel()

		orderId := c.Param("orderId")
		if orderId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid order id"))
			return
		}

		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"orderId": orderId}).Decode(&order)
		if err != nil {
			slog.Error("Error while fetching order", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusNotFound, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, order, "Order fetched successfully")
	}
}

func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		allOrders := make([]models.Order, 0)
		result, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			slog.Error("Error while fetching orders", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if err := result.All(ctx, &allOrders); err != nil {
			slog.Error("Error while fetching orders", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, allOrders, "Orders fetched successfully")
	}
}

func OrderItemOrderCreator(order models.Order) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	order.CreatedAt = time.Now().UTC()
	order.UpdatedAt = time.Now().UTC()
	order.ID = bson.NewObjectID()
	order.OrderID = order.ID.Hex()

	_, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}
	return order.OrderID, nil
}
