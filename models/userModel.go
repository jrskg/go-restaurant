package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `bson:"_id" json:"_id"`
	Name         *string       `bson:"name" json:"name" validate:"required,min=2,max=50"`
	Email        *string       `bson:"email" json:"email" validate:"email,required"`
	Password     *string       `bson:"password" json:"password" validate:"required,min=6"`
	Avatar       *string       `bson:"avatar" json:"avatar"`
	Token        *string       `bson:"token" json:"token"`
	RefreshToken *string       `bson:"refreshToken" json:"refreshToken"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
	UserId       string        `bson:"userId" json:"userId"`
}

type LoginDto struct {
	Email    *string `bson:"email" json:"email" validate:"email,required"`
	Password *string `bson:"password" json:"password" validate:"required"`
}
