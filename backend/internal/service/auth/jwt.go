package service

import (
	"errors"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	jwtSecret string
}

func NewJWTTokenGenerator(jwtSecret string) entity.TokenGenerator {
	return &JWTTokenGenerator{
		jwtSecret: jwtSecret,
	}
}

func (g *JWTTokenGenerator) GenerateToken(claims *entity.JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.jwtSecret))
	if err != nil {
		return "", errors.New("failed to sign token")
	}
	return tokenString, nil
}

type JWTTokenValidator struct {
	jwtSecret string
}

func NewJWTTokenValidator(jwtSecret string) entity.TokenValidator {
	return &JWTTokenValidator{
		jwtSecret: jwtSecret,
	}
}

func (v *JWTTokenValidator) ValidateToken(tokenString string) (*entity.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(v.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*entity.JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
