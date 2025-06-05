package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
)

type mockTokenGenerator struct {
	mock.Mock
}

func (m *mockTokenGenerator) GenerateToken(claims *entity.JWTClaims) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}

type mockTokenValidator struct {
	mock.Mock
}

func (m *mockTokenValidator) ValidateToken(tokenString string) (*entity.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JWTClaims), args.Error(1)
}

type mockTimeProvider struct {
	mock.Mock
}

func (m *mockTimeProvider) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type mockConfig struct {
	AuthSecret      string
	TokenExpiration time.Duration
	Issuer          string
}

func (c *mockConfig) GetAuthSecret() string {
	return c.AuthSecret
}

func (c *mockConfig) GetTokenExpiration() time.Duration {
	return c.TokenExpiration
}

func (c *mockConfig) GetIssuer() string {
	return c.Issuer
}

func TestAuthService_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("successful authentication", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{
			AuthSecret:      "test-auth-secret",
			TokenExpiration: time.Hour,
			Issuer:          "test-issuer",
		}
		mockTokenGen := new(mockTokenGenerator)
		mockTimeProvider := new(mockTimeProvider)
		service := NewAuthService(config, mockTokenGen, nil, mockTimeProvider)

		req := entity.AuthRequest{Secret: "test-auth-secret"}
		now := time.Now()
		expectedToken := "generated-jwt-token"

		mockTimeProvider.On("Now").Return(now).Once()

		mockTokenGen.On("GenerateToken", mock.MatchedBy(func(claims *entity.JWTClaims) bool {
			return claims.Subject == "authenticated-user" &&
				claims.Issuer == "test-issuer"
		})).Return(expectedToken, nil).Once()

		result, err := service.Authenticate(req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedToken, result.Token)
		assert.Equal(t, now.Add(time.Hour), result.ExpiresAt)

		mockTokenGen.AssertExpectations(t)
		mockTimeProvider.AssertExpectations(t)
	})

	t.Run("invalid secret should return error", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{
			AuthSecret:      "test-auth-secret",
			TokenExpiration: time.Hour,
			Issuer:          "test-issuer",
		}
		service := NewAuthService(config, nil, nil, nil)

		req := entity.AuthRequest{Secret: "wrong-secret"}

		result, err := service.Authenticate(req)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid authentication secret")
	})

	t.Run("empty secret should return error", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{
			AuthSecret:      "test-auth-secret",
			TokenExpiration: time.Hour,
			Issuer:          "test-issuer",
		}
		service := NewAuthService(config, nil, nil, nil)

		req := entity.AuthRequest{Secret: ""}

		result, err := service.Authenticate(req)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid authentication secret")
	})

	t.Run("token generation failure should return error", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{
			AuthSecret:      "test-auth-secret",
			TokenExpiration: time.Hour,
			Issuer:          "test-issuer",
		}
		mockTokenGen := new(mockTokenGenerator)
		mockTimeProvider := new(mockTimeProvider)
		service := NewAuthService(config, mockTokenGen, nil, mockTimeProvider)

		req := entity.AuthRequest{Secret: "test-auth-secret"}
		now := time.Now()

		mockTimeProvider.On("Now").Return(now).Once()
		mockTokenGen.On("GenerateToken", mock.Anything).Return("", assert.AnError).Once()

		result, err := service.Authenticate(req)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to generate token")

		mockTokenGen.AssertExpectations(t)
		mockTimeProvider.AssertExpectations(t)
	})
}

func TestAuthService_GenerateToken(t *testing.T) {
	t.Parallel()

	t.Run("successful token generation", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{
			TokenExpiration: 2 * time.Hour,
			Issuer:          "test-issuer",
		}
		mockTokenGen := new(mockTokenGenerator)
		mockTimeProvider := new(mockTimeProvider)
		service := NewAuthService(config, mockTokenGen, nil, mockTimeProvider)

		subject := "test-user"
		now := time.Now()
		expectedToken := "generated-token"

		mockTimeProvider.On("Now").Return(now).Once()
		mockTokenGen.On("GenerateToken", mock.MatchedBy(func(claims *entity.JWTClaims) bool {
			return claims.Subject == subject &&
				claims.Issuer == "test-issuer"
		})).Return(expectedToken, nil).Once()

		result, err := service.GenerateToken(subject)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedToken, result.Token)
		assert.Equal(t, now.Add(2*time.Hour), result.ExpiresAt)

		mockTokenGen.AssertExpectations(t)
		mockTimeProvider.AssertExpectations(t)
	})

	t.Run("empty subject should return error", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{
			TokenExpiration: 2 * time.Hour,
			Issuer:          "test-issuer",
		}
		service := NewAuthService(config, nil, nil, nil)

		result, err := service.GenerateToken("")

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "subject cannot be empty")
	})

	t.Run("token generation failure should return error", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{
			TokenExpiration: 2 * time.Hour,
			Issuer:          "test-issuer",
		}
		mockTokenGen := new(mockTokenGenerator)
		mockTimeProvider := new(mockTimeProvider)
		service := NewAuthService(config, mockTokenGen, nil, mockTimeProvider)

		subject := "test-user"
		now := time.Now()

		mockTimeProvider.On("Now").Return(now).Once()
		mockTokenGen.On("GenerateToken", mock.Anything).Return("", assert.AnError).Once()

		result, err := service.GenerateToken(subject)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to generate token")

		mockTokenGen.AssertExpectations(t)
		mockTimeProvider.AssertExpectations(t)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	t.Parallel()

	t.Run("successful token validation", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{}
		mockTokenValidator := new(mockTokenValidator)
		service := NewAuthService(config, nil, mockTokenValidator, nil)

		tokenString := "valid-token"
		expectedClaims := &entity.JWTClaims{
			Subject: "test-user",
		}

		mockTokenValidator.On("ValidateToken", tokenString).Return(expectedClaims, nil).Once()

		result, err := service.ValidateToken(tokenString)

		require.NoError(t, err)
		assert.Equal(t, expectedClaims, result)

		mockTokenValidator.AssertExpectations(t)
	})

	t.Run("empty token should return error", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{}
		service := NewAuthService(config, nil, nil, nil)

		result, err := service.ValidateToken("")

		require.Error(t, err)
		require.Nil(t, result)
		assert.Contains(t, err.Error(), "token cannot be empty")
	})

	t.Run("invalid token should return error", func(t *testing.T) {
		t.Parallel()

		config := &mockConfig{}
		mockTokenValidator := new(mockTokenValidator)
		service := NewAuthService(config, nil, mockTokenValidator, nil)

		tokenString := "invalid-token"

		mockTokenValidator.On("ValidateToken", tokenString).Return(nil, assert.AnError).Once()

		result, err := service.ValidateToken(tokenString)

		require.Error(t, err)
		require.Nil(t, result)

		mockTokenValidator.AssertExpectations(t)
	})
}
