package usecase

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
	"github.com/Schieck/packs-calculator/internal/domain/errs"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Authenticate(req entity.AuthRequest) (*entity.AuthResult, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthResult), args.Error(1)
}

func (m *mockAuthService) ValidateToken(token string) (*entity.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JWTClaims), args.Error(1)
}

func (m *mockAuthService) GenerateToken(subject string) (*entity.AuthResult, error) {
	args := m.Called(subject)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuthResult), args.Error(1)
}

func TestAuthenticateUseCase_Execute(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	mockService := new(mockAuthService)
	useCase := NewAuthenticateUseCase(mockService, logger)

	t.Run("successful authentication", func(t *testing.T) {
		req := entity.AuthRequest{Secret: "valid-secret"}
		expiresAt := time.Now().Add(time.Hour)

		expectedResponse := &entity.AuthResult{
			Token:     "jwt-token",
			ExpiresAt: expiresAt,
		}

		mockService.On("Authenticate", req).Return(expectedResponse, nil).Once()

		result, err := useCase.Execute(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockService.AssertExpectations(t)
	})

	t.Run("empty secret should return error", func(t *testing.T) {
		req := entity.AuthRequest{Secret: ""}

		result, err := useCase.Execute(req)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidCredentials, err)
		assert.Nil(t, result)
	})

	t.Run("service error should return invalid credentials", func(t *testing.T) {
		req := entity.AuthRequest{Secret: "invalid-secret"}

		mockService.On("Authenticate", req).Return(nil, errs.ErrInvalidCredentials).Once()

		result, err := useCase.Execute(req)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidCredentials, err)
		assert.Nil(t, result)
		mockService.AssertExpectations(t)
	})
}

func TestValidateTokenUseCase_Execute(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	mockService := new(mockAuthService)
	useCase := NewValidateTokenUseCase(mockService, logger)

	t.Run("successful token validation", func(t *testing.T) {
		token := "valid-jwt-token"
		expectedClaims := &entity.JWTClaims{
			Subject: "test-user",
		}

		mockService.On("ValidateToken", token).Return(expectedClaims, nil).Once()

		result, err := useCase.Execute(token)

		assert.NoError(t, err)
		assert.Equal(t, expectedClaims, result)
		mockService.AssertExpectations(t)
	})

	t.Run("empty token should return error", func(t *testing.T) {
		token := ""

		result, err := useCase.Execute(token)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidToken, err)
		assert.Nil(t, result)
	})

	t.Run("service error should return invalid token", func(t *testing.T) {
		token := "invalid-token"

		mockService.On("ValidateToken", token).Return(nil, errs.ErrInvalidToken).Once()

		result, err := useCase.Execute(token)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidToken, err)
		assert.Nil(t, result)
		mockService.AssertExpectations(t)
	})
}
