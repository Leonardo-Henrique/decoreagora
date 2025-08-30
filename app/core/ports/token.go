package ports

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenHandler interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(tokenString string) (*JWTCustomClaims, error)
}

type JWTCustomClaims struct {
	UserID int
	jwt.RegisteredClaims
}
