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
)

type InvoiceViewFromat struct {
	InvoiceId      string    `json:"invoiceId"`
	PaymentMethod  string    `json:"paymentMethod"`
	OrderId        string    `json:"orderId"`
	PaymentStatus  *string   `json:"paymentStatus"`
	PaymentDue     any       `json:"paymentDue"`
	TableNumber    any       `json:"tableNumber"`
	PaymentDueDate time.Time `json:"paymentDueDate"`
	OrderDetails   any       `json:"orderDetails"`
}

var invoiceCollection = database.OpenCollection(database.DBClient, constants.INVOICE_COLLECTION)

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		invoice := models.Invoice{}

		if err := c.BindJSON(&invoice); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(invoice); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		orderResult := orderCollection.FindOne(ctx, bson.M{"orderId": invoice.OrderId})
		if err := orderResult.Err(); err != nil {
			slog.Error("Error while fetching order", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusNotFound, err)
			return
		}

		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		invoice.CreatedAt = time.Now().UTC()
		invoice.UpdatedAt = time.Now().UTC()
		invoice.PaymentDueDate = time.Now().Add(time.Hour * 24).UTC()
		invoice.ID = bson.NewObjectID()
		invoice.InvoiceId = invoice.ID.Hex()

		_, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			slog.Error("Error while creating invoice", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusCreated, invoice, "Invoice created successfully")
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		invoiceId := c.Param("invoiceId")
		if invoiceId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid invoice id"))
			return
		}

		updateInvoiceDto := models.UpdateInvoiceDto{}

		if err := c.BindJSON(&updateInvoiceDto); err != nil {
			if errors.Is(err, io.EOF) {
				utils.ApiError(c, http.StatusBadRequest, errors.New("request body is empty"))
				return
			}
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(updateInvoiceDto); err != nil {
			utils.ApiError(c, http.StatusBadRequest, err)
			return
		}

		status := "PENDING"
		if updateInvoiceDto.PaymentStatus == nil {
			updateInvoiceDto.PaymentStatus = &status
		}

		options := options.FindOneAndUpdate().SetReturnDocument(options.After)
		filter := bson.M{"invoiceId": invoiceId}
		update := bson.M{"$set": bson.M{
			"paymentMethod": updateInvoiceDto.PaymentMethod,
			"paymentStatus": updateInvoiceDto.PaymentStatus,
			"updatedAt":     time.Now().UTC(),
		}}

		updateInvoice := models.Invoice{}
		err := invoiceCollection.FindOneAndUpdate(ctx, filter, update, options).Decode(&updateInvoice)
		if err != nil {
			slog.Error("Error while updating invoice", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, updateInvoice, "Invoice updated successfully")
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		invoiceId := c.Param("invoiceId")
		if invoiceId == "" {
			utils.ApiError(c, http.StatusBadRequest, errors.New("invalid invoice id"))
			return
		}

		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoiceId": invoiceId}).Decode(&invoice)
		if err != nil {
			slog.Error("Error while fetching invoice", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		var invoiceView InvoiceViewFromat
		allOrderItems, err := ItemsByOrder(invoice.OrderId)
		if err != nil {
			slog.Error("Error while fetching order items", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}
		invoiceView.OrderId = invoice.OrderId
		invoiceView.PaymentDueDate = invoice.PaymentDueDate

		invoiceView.PaymentMethod = "null"
		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		}

		invoiceView.InvoiceId = invoice.InvoiceId
		invoiceView.PaymentStatus = invoice.PaymentStatus
		invoiceView.PaymentDue = allOrderItems[0]["paymentDue"]
		invoiceView.TableNumber = allOrderItems[0]["tableNumber"]
		invoiceView.OrderDetails = allOrderItems[0]["orderItems"]

		utils.ApiSuccess(c, http.StatusOK, invoiceView, "Invoice fetched successfully")
	}
}

func GetAllInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		allInvoices := make([]models.Invoice, 0)
		result, err := invoiceCollection.Find(ctx, bson.M{})

		if err != nil {
			slog.Error("Error while fetching invoice", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		if err = result.All(ctx, &allInvoices); err != nil {
			slog.Error("Error while fetching invoice", slog.String("error", err.Error()))
			utils.ApiError(c, http.StatusInternalServerError, err)
			return
		}

		utils.ApiSuccess(c, http.StatusOK, allInvoices, "Invoice fetched successfully")
	}
}
