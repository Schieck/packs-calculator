package http

import (
	"errors"
	"log/slog"
	"net/http"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
	"github.com/Schieck/packs-calculator/internal/domain/errs"
	"github.com/Schieck/packs-calculator/internal/dto"
	authUseCase "github.com/Schieck/packs-calculator/internal/usecase/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authenticateUseCase *authUseCase.AuthenticateUseCase
	logger              *slog.Logger
	validator           *validator.Validate
}

func NewAuthHandler(authenticateUseCase *authUseCase.AuthenticateUseCase, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authenticateUseCase: authenticateUseCase,
		logger:              logger,
		validator:           validator.New(),
	}
}

// Authenticate handles authentication requests
// @Summary Authenticate
// @Description Authenticate using secret and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.AuthRequest true "Authentication secret"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 401 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Router /auth/token [post]
func (h AuthHandler) Authenticate(c *gin.Context) {
	var dtoReq dto.AuthRequest

	if err := c.ShouldBindJSON(&dtoReq); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	if err := h.validator.Struct(&dtoReq); err != nil {
		h.logger.Warn("Request validation failed", "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error:   "Validation failed",
			Details: errs.FormatValidationErrors(err),
		})
		return
	}

	// Convert DTO to domain entity
	entityReq := entity.AuthRequest{
		Secret: dtoReq.Secret,
	}

	// Execute use case with domain entity
	entityResponse, err := h.authenticateUseCase.Execute(entityReq)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
				Error: "Invalid credentials",
			})
			return
		}

		h.logger.Error("Authentication use case failed", "error", err)
		c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	// Convert domain entity to DTO for response
	dtoResponse := &dto.AuthResponse{
		Token:     entityResponse.Token,
		ExpiresAt: entityResponse.ExpiresAt.Unix(),
	}

	c.JSON(http.StatusOK, dtoResponse)
}
