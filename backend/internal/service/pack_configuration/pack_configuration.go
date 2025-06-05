package service

import (
	"fmt"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

type PackConfigurationService struct {
	repository entity.PackConfigurationRepository
}

func NewPackConfigurationService(repository entity.PackConfigurationRepository) *PackConfigurationService {
	return &PackConfigurationService{
		repository: repository,
	}
}

func (s *PackConfigurationService) GetAllConfigurations() ([]*entity.PackConfiguration, error) {
	configurations, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get pack configurations: %w", err)
	}
	return configurations, nil
}

func (s *PackConfigurationService) GetConfigurationByID(id int) (*entity.PackConfiguration, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid pack configuration ID: %d", id)
	}

	configuration, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get pack configuration by ID: %w", err)
	}
	return configuration, nil
}

func (s *PackConfigurationService) GetDefaultConfiguration() (*entity.PackConfiguration, error) {
	configuration, err := s.repository.GetDefault()
	if err != nil {
		return nil, fmt.Errorf("failed to get default pack configuration: %w", err)
	}
	return configuration, nil
}

func (s *PackConfigurationService) CreateConfiguration(name string, packSizes []int) (*entity.PackConfiguration, error) {
	configuration, err := entity.NewPackConfiguration(name, packSizes)
	if err != nil {
		return nil, fmt.Errorf("failed to create pack configuration entity: %w", err)
	}

	createdConfig, err := s.repository.Create(configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to save pack configuration: %w", err)
	}

	return createdConfig, nil
}

func (s *PackConfigurationService) UpdateConfiguration(id int, name string, packSizes []int, isDefault bool) (*entity.PackConfiguration, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid pack configuration ID: %d", id)
	}

	existingConfig, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing configuration: %w", err)
	}

	updatedConfig := &entity.PackConfiguration{
		ID:        existingConfig.ID,
		Name:      name,
		PackSizes: packSizes,
		IsDefault: isDefault,
		IsActive:  existingConfig.IsActive,
		CreatedAt: existingConfig.CreatedAt,
	}

	if err := updatedConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pack configuration: %w", err)
	}

	savedConfig, err := s.repository.Update(updatedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to update pack configuration: %w", err)
	}

	return savedConfig, nil
}

func (s *PackConfigurationService) DeleteConfiguration(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid pack configuration ID: %d", id)
	}

	// Check if this is the default configuration
	defaultConfig, err := s.repository.GetDefault()
	if err == nil && defaultConfig.ID == id {
		return fmt.Errorf("cannot delete the default pack configuration")
	}

	err = s.repository.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete pack configuration: %w", err)
	}

	return nil
}

func (s *PackConfigurationService) SetDefaultConfiguration(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid pack configuration ID: %d", id)
	}

	// Verify the configuration exists and is active
	_, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("configuration not found or inactive: %w", err)
	}

	err = s.repository.SetDefault(id)
	if err != nil {
		return fmt.Errorf("failed to set default configuration: %w", err)
	}

	return nil
}
