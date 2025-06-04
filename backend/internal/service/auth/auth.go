package service

import (
	"errors"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	config         entity.AuthConfig
	tokenGenerator entity.TokenGenerator
	tokenValidator entity.TokenValidator
	timeProvider   entity.TimeProvider
}

func NewAuthService(config entity.AuthConfig, tokenGenerator entity.TokenGenerator, tokenValidator entity.TokenValidator, timeProvider entity.TimeProvider) entity.AuthService {
	return &authService{
		config:         config,
		tokenGenerator: tokenGenerator,
		tokenValidator: tokenValidator,
		timeProvider:   timeProvider,
	}
}

func NewAuthServiceWithDefaults(jwtSecret, authSecret string) entity.AuthService {
	config := NewDefaultAuthConfig(authSecret, "packs-calculator")
	tokenGenerator := NewJWTTokenGenerator(jwtSecret)
	tokenValidator := NewJWTTokenValidator(jwtSecret)
	timeProvider := &DefaultTimeProvider{}

	return NewAuthService(config, tokenGenerator, tokenValidator, timeProvider)
}

func (s *authService) Authenticate(req entity.AuthRequest) (*entity.AuthResult, error) {
	if req.Secret == "" || req.Secret != s.config.GetAuthSecret() {
		return nil, errors.New("invalid authentication secret")
	}

	return s.GenerateToken("authenticated-user")
}

func (s *authService) GenerateToken(subject string) (*entity.AuthResult, error) {
	if subject == "" {
		return nil, errors.New("subject cannot be empty")
	}

	now := s.timeProvider.Now()
	expirationTime := now.Add(s.config.GetTokenExpiration())

	claims := &entity.JWTClaims{
		Subject: subject,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.GetIssuer(),
		},
	}

	tokenString, err := s.tokenGenerator.GenerateToken(claims)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &entity.AuthResult{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*entity.JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token cannot be empty")
	}

	return s.tokenValidator.ValidateToken(tokenString)
}
