package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
)

func TestJWTTokenGenerator_GenerateToken(t *testing.T) {
	t.Parallel()

	generator := NewJWTTokenGenerator("test-secret")

	t.Run("valid claims should generate token", func(t *testing.T) {
		t.Parallel()

		now := time.Now()
		claims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "test-issuer",
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			},
		}

		token, err := generator.GenerateToken(claims)

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token can be parsed
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		})
		require.NoError(t, err)
		assert.True(t, parsedToken.Valid)
	})

	t.Run("nil claims should still work", func(t *testing.T) {
		t.Parallel()

		token, err := generator.GenerateToken(nil)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("empty secret should still work", func(t *testing.T) {
		t.Parallel()

		emptySecretGenerator := NewJWTTokenGenerator("")
		now := time.Now()
		claims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "test-issuer",
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			},
		}

		token, err := emptySecretGenerator.GenerateToken(claims)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("expired claims should still generate token", func(t *testing.T) {
		t.Parallel()

		pastTime := time.Now().Add(-2 * time.Hour)
		claims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "test-issuer",
				IssuedAt:  jwt.NewNumericDate(pastTime),
				ExpiresAt: jwt.NewNumericDate(pastTime.Add(time.Hour)), // Expired
			},
		}

		token, err := generator.GenerateToken(claims)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestJWTTokenValidator_ValidateToken(t *testing.T) {
	t.Parallel()

	secret := "test-secret"
	validator := NewJWTTokenValidator(secret)
	generator := NewJWTTokenGenerator(secret)

	t.Run("valid token should be validated", func(t *testing.T) {
		t.Parallel()

		now := time.Now()
		originalClaims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "test-issuer",
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			},
		}

		token, err := generator.GenerateToken(originalClaims)
		require.NoError(t, err)

		claims, err := validator.ValidateToken(token)

		require.NoError(t, err)
		require.NotNil(t, claims)
		assert.Equal(t, originalClaims.Subject, claims.Subject)
		assert.Equal(t, originalClaims.Issuer, claims.Issuer)
	})

	t.Run("invalid token should return error", func(t *testing.T) {
		t.Parallel()

		invalidToken := "invalid.token.here"

		claims, err := validator.ValidateToken(invalidToken)

		require.Error(t, err)
		require.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("empty token should return error", func(t *testing.T) {
		t.Parallel()

		claims, err := validator.ValidateToken("")

		require.Error(t, err)
		require.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("expired token should return error", func(t *testing.T) {
		t.Parallel()

		pastTime := time.Now().Add(-2 * time.Hour)
		expiredClaims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "test-issuer",
				IssuedAt:  jwt.NewNumericDate(pastTime),
				ExpiresAt: jwt.NewNumericDate(pastTime.Add(time.Hour)), // Expired
			},
		}

		token, err := generator.GenerateToken(expiredClaims)
		require.NoError(t, err)

		claims, err := validator.ValidateToken(token)

		require.Error(t, err)
		require.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("token signed with different secret should return error", func(t *testing.T) {
		t.Parallel()

		differentSecretGenerator := NewJWTTokenGenerator("different-secret")
		now := time.Now()
		claims := &entity.JWTClaims{
			Subject: "test-user",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "test-issuer",
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			},
		}

		token, err := differentSecretGenerator.GenerateToken(claims)
		require.NoError(t, err)

		validatedClaims, err := validator.ValidateToken(token)

		require.Error(t, err)
		require.Nil(t, validatedClaims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("malformed token should return error", func(t *testing.T) {
		t.Parallel()

		malformedToken := "not.a.jwt"

		claims, err := validator.ValidateToken(malformedToken)

		require.Error(t, err)
		require.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})
}

func TestDefaultAuthConfig(t *testing.T) {
	t.Parallel()

	config := NewDefaultAuthConfig("test-secret", "test-issuer")

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.GetAuthSecret())
	assert.Positive(t, config.GetTokenExpiration())
	assert.NotEmpty(t, config.GetIssuer())
}

func TestDefaultTimeProvider(t *testing.T) {
	t.Parallel()

	provider := &DefaultTimeProvider{}

	assert.NotNil(t, provider)

	now := provider.Now()
	assert.WithinDuration(t, time.Now(), now, time.Second)
}
