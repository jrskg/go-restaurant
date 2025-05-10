package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Food struct {
	ID        bson.ObjectID `bson:"_id" json:"_id"`
	Name      *string       `bson:"name" json:"name" validate:"required,min=2,max=50"`
	Price     *float64      `bson:"price" json:"price" validate:"required"`
	FoodImage *string       `bson:"foodImage" json:"foodImage" validate:"required"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
	FoodId    string        `bson:"foodId" json:"foodId"`
	MenuId    *string       `bson:"menuId" json:"menuId" validate:"required"`
}

type UpdateFoodDto struct {
	Name      *string  `json:"name,omitempty" validate:"omitempty,required,min=2,max=50"`
	Price     *float64 `json:"price,omitempty" validate:"omitempty,required"`
	FoodImage *string  `json:"foodImage,omitempty" validate:"omitempty,required"`
	MenuId    *string  `json:"menuId,omitempty" validate:"omitempty,required"`
}
