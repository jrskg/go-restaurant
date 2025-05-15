package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Table struct {
	ID             bson.ObjectID `bson:"_id" json:"_id"`
	NumberOfGuests *int          `bson:"numberOfGuests" json:"numberOfGuests" validate:"required"`
	TableNumber    *int          `bson:"tableNumber" json:"tableNumber" validate:"required"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"updatedAt"`
	TableId        string        `bson:"tableId" json:"tableId"`
}

type UpdateTableDto struct {
	NumberOfGuests *int `json:"numberOfGuests"`
	TableNumber    *int `json:"tableNumber"`
}
