package service

import (
	"time"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
	"github.com/Schieck/packs-calculator/pkg/db"
)

type HealthService struct {
	database *db.DB
	version  string
}

func NewHealthService(database *db.DB, version string) *HealthService {
	return &HealthService{
		database: database,
		version:  version,
	}
}

func (s *HealthService) CheckHealth() (*entity.HealthResponse, error) {
	err := s.database.Ping()
	if err != nil {
		return nil, err
	}

	return &entity.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   s.version,
	}, nil
}
