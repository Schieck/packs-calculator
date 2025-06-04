package entity

import (
	"time"
)

type AuthRequest struct {
	Secret string
}

type AuthResult struct {
	Token     string
	ExpiresAt time.Time
}

type AuthService interface {
	Authenticate(req AuthRequest) (*AuthResult, error)
	ValidateToken(token string) (*JWTClaims, error)
	GenerateToken(subject string) (*AuthResult, error)
}

type AuthenticateUseCase interface {
	Execute(req AuthRequest) (*AuthResult, error)
}

type AuthConfig interface {
	GetAuthSecret() string
	GetTokenExpiration() time.Duration
	GetIssuer() string
}
