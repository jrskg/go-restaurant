package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Order struct {
	ID        bson.ObjectID `bson:"_id" json:"_id"`
	OrderDate time.Time     `bson:"orderDate" json:"orderDate" validate:"required"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
	OrderID   string        `bson:"orderId" json:"orderId"`
	TableId   string        `bson:"tableId" json:"tableId" validate:"required"`
}

type UpdateOrderDto struct {
	TableId *string `json:"tableId,omitempty" validate:"omitempty,required"`
}
