package usecase

import (
	"log/slog"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
	"github.com/Schieck/packs-calculator/internal/domain/errs"
)

type AuthenticateUseCase struct {
	authService entity.AuthService
	logger      *slog.Logger
}

type ValidateTokenUseCase struct {
	authService entity.AuthService
	logger      *slog.Logger
}

func NewAuthenticateUseCase(authService entity.AuthService, logger *slog.Logger) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		authService: authService,
		logger:      logger,
	}
}

func NewValidateTokenUseCase(authService entity.AuthService, logger *slog.Logger) entity.ValidateTokenUseCase {
	return &ValidateTokenUseCase{
		authService: authService,
		logger:      logger,
	}
}

func (uc *AuthenticateUseCase) Execute(req entity.AuthRequest) (*entity.AuthResult, error) {
	uc.logger.Info("Executing authentication use case")

	if err := uc.validateAuthRequest(req); err != nil {
		uc.logger.Warn("Authentication request validation failed", "error", err)
		return nil, err
	}

	response, err := uc.authService.Authenticate(req)
	if err != nil {
		uc.logger.Warn("Authentication failed", "error", err)
		return nil, errs.ErrInvalidCredentials
	}

	uc.logger.Info("Authentication successful")
	return response, nil
}

func (uc *ValidateTokenUseCase) Execute(token string) (*entity.JWTClaims, error) {
	uc.logger.Debug("Executing token validation use case")

	if err := uc.validateToken(token); err != nil {
		uc.logger.Warn("Token validation failed", "error", err)
		return nil, err
	}

	claims, err := uc.authService.ValidateToken(token)
	if err != nil {
		uc.logger.Warn("Token validation failed", "error", err)
		return nil, errs.ErrInvalidToken
	}

	uc.logger.Debug("Token validation successful", "subject", claims.Subject)
	return claims, nil
}

func (uc *AuthenticateUseCase) validateAuthRequest(req entity.AuthRequest) error {
	if req.Secret == "" {
		return errs.ErrInvalidCredentials
	}
	return nil
}

func (uc *ValidateTokenUseCase) validateToken(token string) error {
	if token == "" {
		return errs.ErrInvalidToken
	}
	return nil
}
