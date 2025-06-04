package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
)

func TestAuthService_Integration(t *testing.T) {
	jwtSecret := "test-jwt-secret"
	authSecret := "test-auth-secret"

	authService := NewAuthServiceWithDefaults(jwtSecret, authSecret)

	t.Run("complete authentication and validation flow", func(t *testing.T) {
		// Step 1: Authenticate
		authReq := entity.AuthRequest{Secret: authSecret}
		authResp, err := authService.Authenticate(authReq)

		assert.NoError(t, err)
		assert.NotNil(t, authResp)
		assert.NotEmpty(t, authResp.Token)
		assert.False(t, authResp.ExpiresAt.IsZero())

		// Step 2: Validate the generated token
		claims, err := authService.ValidateToken(authResp.Token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "authenticated-user", claims.Subject)
		assert.Equal(t, "packs-calculator", claims.Issuer)
	})

	t.Run("direct token generation and validation", func(t *testing.T) {
		// Step 1: Generate token directly
		subject := "direct-user"
		tokenResp, err := authService.GenerateToken(subject)

		assert.NoError(t, err)
		assert.NotNil(t, tokenResp)
		assert.NotEmpty(t, tokenResp.Token)

		// Step 2: Validate the generated token
		claims, err := authService.ValidateToken(tokenResp.Token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, subject, claims.Subject)
		assert.Equal(t, "packs-calculator", claims.Issuer)
	})

	t.Run("invalid authentication should not affect valid tokens", func(t *testing.T) {
		// Generate a valid token first
		validTokenResp, err := authService.GenerateToken("valid-user")
		assert.NoError(t, err)

		// Try invalid authentication
		invalidReq := entity.AuthRequest{Secret: "wrong-secret"}
		_, err = authService.Authenticate(invalidReq)
		assert.Error(t, err)

		// The valid token should still work
		claims, err := authService.ValidateToken(validTokenResp.Token)
		assert.NoError(t, err)
		assert.Equal(t, "valid-user", claims.Subject)
	})
}
