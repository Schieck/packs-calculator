package usecase

import (
	"log/slog"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

type HealthService interface {
	CheckHealth() (*entity.HealthResponse, error)
}

type HealthUseCase struct {
	healthService HealthService
	logger        *slog.Logger
}

func NewHealthUseCase(healthService HealthService, logger *slog.Logger) *HealthUseCase {
	return &HealthUseCase{
		healthService: healthService,
		logger:        logger,
	}
}

func (uc *HealthUseCase) Execute() (*entity.HealthResponse, error) {
	uc.logger.Debug("Executing health check use case")

	response, err := uc.healthService.CheckHealth()
	if err != nil {
		uc.logger.Error("Health check failed", "error", err)
		return nil, err
	}

	uc.logger.Debug("Health check successful", "status", response.Status)
	return response, nil
}
