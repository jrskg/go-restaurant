package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
	UserId string
	Name   string
	Email  string
	Avatar string
	jwt.RegisteredClaims
}

var secretKey = os.Getenv("JWT_SECRET")

func GetTokens(userId, name, email, avatar string) (token, refreshToken string, err error) {
	tokenClaims := &SignedDetails{
		UserId: userId,
		Name:   name,
		Email:  email,
		Avatar: avatar,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	refreshTokenClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * 24 * time.Hour)),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func ValidateAndDecodeToken(signedString string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(signedString, &SignedDetails{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(*SignedDetails); ok && token.Valid {
		exp := claims.ExpiresAt.Time.Unix()
		if exp < time.Now().Unix() {
			return nil, errors.New("token expired")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}
