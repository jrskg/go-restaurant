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
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderItemPack struct {
	TableId    *string            `json:"tableId"`
	OrderItems []models.OrderItem `json:"orderItems"`
}

var orderItemCollection = database.OpenCollection(database.DBClient, constants.ORDER_ITEM_COLLECTION)

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderItemPack := OrderItemPack{}
		order := models.Order{}

		if err := c.BindJSON(&orderItemPack); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if orderItemPack.TableId == nil || len(orderItemPack.OrderItems) < 1 {
			utils.ApiError(c, http.StatusBadRequest, errors.New("tableId or orderItems is empty"))
			return
		}

		order.OrderDate = time.Now().UTC()
		order.TableId = *orderItemPack.TableId
		createdOrderId, err := OrderItemOrderCreator(order)
		if err != nil {
			slog.Error("Error while creating order", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		orderItemsToBeInserted := make([]models.OrderItem, 0)

		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.OrderId = createdOrderId
			num := utils.ToFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num
			if err := utils.Validate.Struct(orderItem); err != nil {
				utils.ApiError(c, http.StatusBadRequest, err)
				return
			}

			orderItem.ID = bson.NewObjectID()
			orderItem.OrderItemId = orderItem.ID.Hex()
			orderItem.CreatedAt = time.Now().UTC()
			orderItem.UpdatedAt = time.Now().UTC()

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		_, err = orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			slog.Error("Error while inserting order items", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusCreated, orderItemsToBeInserted, "Order items created successfully")
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItemId")
		if orderItemId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid order item id"))
			return
		}

		updateOrderItemDto := models.UpdateOrderItemDto{}
		if err := c.BindJSON(&updateOrderItemDto); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(updateOrderItemDto); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}
		if updateOrderItemDto.UnitPrice != nil {
			num := utils.ToFixed(*updateOrderItemDto.UnitPrice, 2)
			updateOrderItemDto.UnitPrice = &num
		}
		fieldsToUpdate := bson.M{
			"unitPrice": updateOrderItemDto.UnitPrice,
			"quantity":  updateOrderItemDto.Quantity,
			"foodId":    updateOrderItemDto.FoodId,
		}
		updateObj := bson.M{"updatedAt": time.Now().UTC()}
		for k, v := range fieldsToUpdate {
			if !utils.IsNil(v) {
				updateObj[k] = v
			}
		}
		filter := bson.M{"orderItemId": orderItemId}

		restaurant, err := orderItemCollection.UpdateOne(ctx, filter, bson.M{"$set": updateObj})
		if err != nil {
			slog.Error("Error while updating order items", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}
		if restaurant.ModifiedCount != 1 {
			utils.ApiError(c, http.StatusInternalServerError, errors.New("order item not found"))
			return
		}

		utils.ApiSuccess(c, http.StatusOK, nil, "Order item updated successfully")
	}
}

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := orderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			slog.Error("Error while fetching order items", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		allOrderItems := make([]models.OrderItem, 0)
		if err := result.All(ctx, &allOrderItems); err != nil {
			slog.Error("Error while fetching order items", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, allOrderItems, "Order items fetched successfully")
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("orderId")
		if orderId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid order id"))
			return
		}

		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, allOrderItems, "Order items fetched successfully")
	}
}

func ItemsByOrder(id string) (orderItems []bson.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{{Key: "orderId", Value: id}}},
	}
	lookupFoodStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "food"},
			{Key: "localField", Value: "foodId"},
			{Key: "foreignField", Value: "foodId"},
			{Key: "as", Value: "food"},
		}},
	}
	unwindFoodStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$food"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}
	lookupOrderStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "order"},
			{Key: "localField", Value: "orderId"},
			{Key: "foreignField", Value: "orderId"},
			{Key: "as", Value: "order"},
		}},
	}
	unwindOrderStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$order"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}
	lookupTableStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "table"},
			{Key: "localField", Value: "order.tableId"},
			{Key: "foreignField", Value: "tableId"},
			{Key: "as", Value: "table"},
		}},
	}
	unwindTableStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$table"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "amount", Value: "$food.price"},
			{Key: "foodName", Value: "$food.name"},
			{Key: "foodImage", Value: "$food.foodImage"},
			{Key: "totalCount", Value: 1},
			{Key: "tableNumber", Value: "$table.tableNumber"},
			{Key: "tableId", Value: "$table.tableId"},
			{Key: "orderId", Value: "$order.orderId"},
			{Key: "quantity", Value: 1},
		}},
	}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "orderId", Value: "$orderId"},
				{Key: "tableId", Value: "$tableId"},
				{Key: "tableNumber", Value: "$tableNumber"},
			}},
			{Key: "paymentDue", Value: bson.D{
				{Key: "$sum", Value: "$amount"},
			}},
			{Key: "totalCount", Value: bson.D{
				{Key: "$sum", Value: 1},
			}},
			{Key: "orderItems", Value: bson.D{
				{Key: "$push", Value: "$$ROOT"},
			}},
		}},
	}
	projectStage2 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "paymentDue", Value: 1},
			{Key: "totalCount", Value: 1},
			{Key: "tableNumber", Value: "$_id.tableNumber"},
			{Key: "orderItems", Value: 1},
		}},
	}

	cursor, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupFoodStage,
		unwindFoodStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &orderItems); err != nil {
		return nil, err
	}

	return orderItems, nil
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItemId")
		if orderItemId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid order item id"))
			return
		}

		orderItem := models.OrderItem{}
		err := orderItemCollection.FindOne(ctx, bson.M{"orderItemId": orderItemId}).Decode(&orderItem)
		if err != nil {
			slog.Error("Error while fetching order item", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusNotFound, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, orderItem, "Order item fetched successfully")
	}
}
