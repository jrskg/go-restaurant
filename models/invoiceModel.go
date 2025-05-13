package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Invoice struct {
	ID             bson.ObjectID `bson:"_id" json:"_id"`
	InvoiceId      string        `bson:"invoiceId" json:"invoiceId"`
	PaymentMethod  *string       `bson:"paymentMethod" json:"paymentMethod" validate:"eq=CASH|eq=CARD|ep="`
	PaymentStatus  *string       `bson:"paymentStatus" json:"paymentStatus" validate:"eq=PAID|eq=PENDING"`
	PaymentDueDate time.Time     `bson:"paymentDueDate" json:"paymentDueDate"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"updatedAt"`
	OrderId        string        `bson:"orderId" json:"orderId" validate:"required"`
}

type UpdateInvoiceDto struct {
	PaymentMethod *string `json:"paymentMethod" validate:"eq=CASH|eq=CARD|ep="`
	PaymentStatus *string `json:"paymentStatus" validate:"eq=PAID|eq=PENDING"`
}
