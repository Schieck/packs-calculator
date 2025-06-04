package http

import (
	"log/slog"
	"net/http"

	"github.com/Schieck/packs-calculator/internal/domain/errs"
	"github.com/Schieck/packs-calculator/internal/dto"
	healthUseCase "github.com/Schieck/packs-calculator/internal/usecase/health"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthUseCase *healthUseCase.HealthUseCase
	logger        *slog.Logger
}

func NewHealthHandler(healthUseCase *healthUseCase.HealthUseCase, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		healthUseCase: healthUseCase,
		logger:        logger,
	}
}

// Health handles health check requests
// @Summary Health Check
// @Description Check the health status of the API and database connection
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} errs.ErrorResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	entityResponse, err := h.healthUseCase.Execute()
	if err != nil {
		h.logger.Error("Health check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, errs.ErrorResponse{
			Error: "Database connection failed",
		})
		return
	}

	dtoResponse := &dto.HealthResponse{
		Status:    entityResponse.Status,
		Timestamp: entityResponse.Timestamp,
		Version:   entityResponse.Version,
	}

	c.JSON(http.StatusOK, dtoResponse)
}
