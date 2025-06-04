package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Subject string `json:"sub"`
	jwt.RegisteredClaims
}

type ValidateTokenUseCase interface {
	Execute(token string) (*JWTClaims, error)
}

type TokenGenerator interface {
	GenerateToken(claims *JWTClaims) (string, error)
}

type TokenValidator interface {
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type TimeProvider interface {
	Now() time.Time
}
