package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jrskg/go-restaurant/helpers"
	"github.com/jrskg/go-restaurant/utils"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ApiError(c, http.StatusUnauthorized, errors.New("unauthorized"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ApiError(c, http.StatusUnauthorized, errors.New("unauthorized"))
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims, err := helpers.ValidateAndDecodeToken(tokenStr)
		if err != nil {
			utils.ApiError(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("userId", claims.UserId)
		c.Set("name", claims.Name)
		c.Set("avatar", claims.Avatar)

		c.Next()
	}
}
