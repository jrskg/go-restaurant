package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ApiError(c *gin.Context, code int, err error) {
	c.JSON(
		code,
		gin.H{
			"success": false,
			"message": err.Error(),
		},
	)
}

func ApiSuccess(c *gin.Context, code int, data any, message string) {
	c.JSON(
		code,
		gin.H{
			"success": true,
			"message": message,
			"data":    data,
		},
	)
}

var Validate = validator.New()

func ToFixed(num float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Floor(num*factor) / factor
}

func ValidateAndParseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("invalid time format")
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format")
	}
	return parsedTime.UTC(), nil
}
