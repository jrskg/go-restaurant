package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Note struct {
	ID        bson.ObjectID `bson:"_id" json:"_id"`
	Text      *string       `bson:"text" json:"text" validate:"required"`
	Title     *string       `bson:"title" json:"title" validate:"required"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
	NoteId    string        `bson:"noteId" json:"noteId"`
}
