package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Menu struct {
	ID        bson.ObjectID `bson:"_id" json:"_id"`
	Name      string        `bson:"name" json:"name" validate:"required"`
	Category  string        `bson:"category" json:"category" validate:"required"`
	StartDate *time.Time    `bson:"startDate" json:"startDate"`
	EndDate   *time.Time    `bson:"endDate" json:"endDate"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
	MenuId    string        `bson:"menuId" json:"menuId"`
}

type MenuUpdateDto struct {
	Name      *string    `json:"name,omitempty" validate:"omitempty,required,min=2,max=50"`
	Category  *string    `json:"category,omitempty" validate:"omitempty,required,min=2,max=50"`
	StartDate *time.Time `json:"startDate,omitempty" validate:"omitempty,required"`
	EndDate   *time.Time `json:"endDate,omitempty" validate:"omitempty,required"`
}
