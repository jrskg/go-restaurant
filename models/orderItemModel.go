package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderItem struct {
	ID          bson.ObjectID `bson:"_id" json:"_id"`
	Quantity    *string       `bson:"quantity" json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice   *float64      `bson:"unitPrice" json:"unitPrice" validate:"required"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
	OrderItemId string        `bson:"orderItemId" json:"orderItemId"`
	OrderId     string        `bson:"orderId" json:"orderId" validate:"required"`
	FoodId      string        `bson:"foodId" json:"foodId" validate:"required"`
}
