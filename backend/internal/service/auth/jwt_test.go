package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
)

func TestJWTTokenGenerator_GenerateToken(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	generator := NewJWTTokenGenerator(jwtSecret)

	t.Run("successful token generation", func(t *testing.T) {
		now := time.Now()
		claims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "test-app",
			},
		}

		token, err := generator.GenerateToken(claims)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Contains(t, token, ".") // JWT should have dots separating parts
	})

	t.Run("generated token should be valid JWT", func(t *testing.T) {
		now := time.Now()
		claims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "test-app",
			},
		}

		tokenString, err := generator.GenerateToken(claims)
		assert.NoError(t, err)

		// Verify the token can be parsed
		token, err := jwt.ParseWithClaims(tokenString, &entity.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		assert.NoError(t, err)
		assert.True(t, token.Valid)

		parsedClaims, ok := token.Claims.(*entity.JWTClaims)
		assert.True(t, ok)
		assert.Equal(t, "test-user", parsedClaims.Subject)
		assert.Equal(t, "test-app", parsedClaims.Issuer)
	})
}

func TestJWTTokenValidator_ValidateToken(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	validator := NewJWTTokenValidator(jwtSecret)
	generator := NewJWTTokenGenerator(jwtSecret)

	t.Run("successful token validation", func(t *testing.T) {
		// First generate a valid token
		now := time.Now()
		originalClaims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "test-app",
			},
		}

		tokenString, err := generator.GenerateToken(originalClaims)
		assert.NoError(t, err)

		// Now validate it
		claims, err := validator.ValidateToken(tokenString)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "test-user", claims.Subject)
		assert.Equal(t, "test-app", claims.Issuer)
	})

	t.Run("invalid token should return error", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		claims, err := validator.ValidateToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("token with wrong secret should return error", func(t *testing.T) {
		wrongSecretGenerator := NewJWTTokenGenerator("wrong-secret")
		now := time.Now()
		claims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				Issuer:    "test-app",
			},
		}

		tokenString, err := wrongSecretGenerator.GenerateToken(claims)
		assert.NoError(t, err)

		// Try to validate with correct secret validator
		result, err := validator.ValidateToken(tokenString)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("expired token should return error", func(t *testing.T) {
		// Create an expired token
		pastTime := time.Now().Add(-time.Hour)
		expiredClaims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(pastTime), // Expired
				IssuedAt:  jwt.NewNumericDate(pastTime.Add(-time.Hour)),
				NotBefore: jwt.NewNumericDate(pastTime.Add(-time.Hour)),
				Issuer:    "test-app",
			},
		}

		tokenString, err := generator.GenerateToken(expiredClaims)
		assert.NoError(t, err)

		claims, err := validator.ValidateToken(tokenString)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})
}

func TestDefaultAuthConfig(t *testing.T) {
	authSecret := "test-auth-secret"
	issuer := "test-issuer"

	config := NewDefaultAuthConfig(authSecret, issuer)

	t.Run("should return correct values", func(t *testing.T) {
		assert.Equal(t, authSecret, config.GetAuthSecret())
		assert.Equal(t, issuer, config.GetIssuer())
		assert.Equal(t, 24*time.Hour, config.GetTokenExpiration())
	})
}

func TestDefaultTimeProvider(t *testing.T) {
	provider := &DefaultTimeProvider{}

	t.Run("should return current time", func(t *testing.T) {
		before := time.Now()
		result := provider.Now()
		after := time.Now()

		// The result should be between before and after (allowing for small execution time)
		assert.True(t, result.After(before.Add(-time.Second)))
		assert.True(t, result.Before(after.Add(time.Second)))
	})
}
