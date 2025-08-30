package repositories

import (
	"fmt"
	"time"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct{}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) GenerateToken(userID int) (string, error) {
	claims := ports.JWTCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "decoreagora-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.C.JWT_SECRET)
}

func (j *JWT) ValidateToken(tokenString string) (*ports.JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ports.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.C.JWT_SECRET, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ports.JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
