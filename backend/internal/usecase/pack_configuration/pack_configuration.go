package usecase

import (
	"fmt"
	"log/slog"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

type PackConfigurationService interface {
	GetAllConfigurations() ([]*entity.PackConfiguration, error)
	GetConfigurationByID(id int) (*entity.PackConfiguration, error)
	GetDefaultConfiguration() (*entity.PackConfiguration, error)
	CreateConfiguration(name string, packSizes []int) (*entity.PackConfiguration, error)
	UpdateConfiguration(id int, name string, packSizes []int, isDefault bool) (*entity.PackConfiguration, error)
	DeleteConfiguration(id int) error
	SetDefaultConfiguration(id int) error
}

type GetAllConfigurationsUseCase struct {
	service PackConfigurationService
	logger  *slog.Logger
}

func NewGetAllConfigurationsUseCase(service PackConfigurationService, logger *slog.Logger) *GetAllConfigurationsUseCase {
	return &GetAllConfigurationsUseCase{
		service: service,
		logger:  logger,
	}
}

func (uc *GetAllConfigurationsUseCase) Execute() ([]*entity.PackConfiguration, error) {
	uc.logger.Info("Executing get all pack configurations use case")

	configurations, err := uc.service.GetAllConfigurations()
	if err != nil {
		uc.logger.Error("Failed to get all pack configurations", "error", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved pack configurations", "count", len(configurations))
	return configurations, nil
}

type GetConfigurationByIDUseCase struct {
	service PackConfigurationService
	logger  *slog.Logger
}

func NewGetConfigurationByIDUseCase(service PackConfigurationService, logger *slog.Logger) *GetConfigurationByIDUseCase {
	return &GetConfigurationByIDUseCase{
		service: service,
		logger:  logger,
	}
}

func (uc *GetConfigurationByIDUseCase) Execute(id int) (*entity.PackConfiguration, error) {
	uc.logger.Info("Executing get pack configuration by ID use case", "id", id)

	if id <= 0 {
		uc.logger.Warn("Invalid pack configuration ID", "id", id)
		return nil, fmt.Errorf("invalid pack configuration ID: %d", id)
	}

	configuration, err := uc.service.GetConfigurationByID(id)
	if err != nil {
		uc.logger.Error("Failed to get pack configuration by ID", "id", id, "error", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved pack configuration", "id", id, "name", configuration.Name)
	return configuration, nil
}

type GetDefaultConfigurationUseCase struct {
	service PackConfigurationService
	logger  *slog.Logger
}

func NewGetDefaultConfigurationUseCase(service PackConfigurationService, logger *slog.Logger) *GetDefaultConfigurationUseCase {
	return &GetDefaultConfigurationUseCase{
		service: service,
		logger:  logger,
	}
}

func (uc *GetDefaultConfigurationUseCase) Execute() (*entity.PackConfiguration, error) {
	uc.logger.Info("Executing get default pack configuration use case")

	configuration, err := uc.service.GetDefaultConfiguration()
	if err != nil {
		uc.logger.Error("Failed to get default pack configuration", "error", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved default pack configuration", "id", configuration.ID, "name", configuration.Name)
	return configuration, nil
}

type CreateConfigurationUseCase struct {
	service PackConfigurationService
	logger  *slog.Logger
}

func NewCreateConfigurationUseCase(service PackConfigurationService, logger *slog.Logger) *CreateConfigurationUseCase {
	return &CreateConfigurationUseCase{
		service: service,
		logger:  logger,
	}
}

func (uc *CreateConfigurationUseCase) Execute(name string, packSizes []int) (*entity.PackConfiguration, error) {
	uc.logger.Info("Executing create pack configuration use case", "name", name, "pack_sizes", packSizes)

	if err := uc.validateInput(name, packSizes); err != nil {
		uc.logger.Warn("Create pack configuration input validation failed", "error", err)
		return nil, err
	}

	configuration, err := uc.service.CreateConfiguration(name, packSizes)
	if err != nil {
		uc.logger.Error("Failed to create pack configuration", "name", name, "error", err)
		return nil, err
	}

	uc.logger.Info("Successfully created pack configuration", "id", configuration.ID, "name", configuration.Name)
	return configuration, nil
}

func (uc *CreateConfigurationUseCase) validateInput(name string, packSizes []int) error {
	if name == "" {
		return fmt.Errorf("pack configuration name cannot be empty")
	}

	if len(packSizes) == 0 {
		return fmt.Errorf("pack configuration must have at least one pack size")
	}

	for _, size := range packSizes {
		if size <= 0 {
			return fmt.Errorf("pack sizes must be positive, got %d", size)
		}
	}

	return nil
}

type UpdateConfigurationUseCase struct {
	service PackConfigurationService
	logger  *slog.Logger
}

func NewUpdateConfigurationUseCase(service PackConfigurationService, logger *slog.Logger) *UpdateConfigurationUseCase {
	return &UpdateConfigurationUseCase{
		service: service,
		logger:  logger,
	}
}

func (uc *UpdateConfigurationUseCase) Execute(id int, name string, packSizes []int, isDefault bool) (*entity.PackConfiguration, error) {
	uc.logger.Info("Executing update pack configuration use case", "id", id, "name", name, "pack_sizes", packSizes, "is_default", isDefault)

	if err := uc.validateInput(id, name, packSizes); err != nil {
		uc.logger.Warn("Update pack configuration input validation failed", "error", err)
		return nil, err
	}

	configuration, err := uc.service.UpdateConfiguration(id, name, packSizes, isDefault)
	if err != nil {
		uc.logger.Error("Failed to update pack configuration", "id", id, "error", err)
		return nil, err
	}

	uc.logger.Info("Successfully updated pack configuration", "id", configuration.ID, "name", configuration.Name)
	return configuration, nil
}

func (uc *UpdateConfigurationUseCase) validateInput(id int, name string, packSizes []int) error {
	if id <= 0 {
		return fmt.Errorf("invalid pack configuration ID: %d", id)
	}

	if name == "" {
		return fmt.Errorf("pack configuration name cannot be empty")
	}

	if len(packSizes) == 0 {
		return fmt.Errorf("pack configuration must have at least one pack size")
	}

	for _, size := range packSizes {
		if size <= 0 {
			return fmt.Errorf("pack sizes must be positive, got %d", size)
		}
	}

	return nil
}

type DeleteConfigurationUseCase struct {
	service PackConfigurationService
	logger  *slog.Logger
}

func NewDeleteConfigurationUseCase(service PackConfigurationService, logger *slog.Logger) *DeleteConfigurationUseCase {
	return &DeleteConfigurationUseCase{
		service: service,
		logger:  logger,
	}
}

func (uc *DeleteConfigurationUseCase) Execute(id int) error {
	uc.logger.Info("Executing delete pack configuration use case", "id", id)

	if id <= 0 {
		uc.logger.Warn("Invalid pack configuration ID", "id", id)
		return fmt.Errorf("invalid pack configuration ID: %d", id)
	}

	err := uc.service.DeleteConfiguration(id)
	if err != nil {
		uc.logger.Error("Failed to delete pack configuration", "id", id, "error", err)
		return err
	}

	uc.logger.Info("Successfully deleted pack configuration", "id", id)
	return nil
}

type SetDefaultConfigurationUseCase struct {
	service PackConfigurationService
	logger  *slog.Logger
}

func NewSetDefaultConfigurationUseCase(service PackConfigurationService, logger *slog.Logger) *SetDefaultConfigurationUseCase {
	return &SetDefaultConfigurationUseCase{
		service: service,
		logger:  logger,
	}
}

func (uc *SetDefaultConfigurationUseCase) Execute(id int) error {
	uc.logger.Info("Executing set default pack configuration use case", "id", id)

	if id <= 0 {
		uc.logger.Warn("Invalid pack configuration ID", "id", id)
		return fmt.Errorf("invalid pack configuration ID: %d", id)
	}

	err := uc.service.SetDefaultConfiguration(id)
	if err != nil {
		uc.logger.Error("Failed to set default pack configuration", "id", id, "error", err)
		return err
	}

	uc.logger.Info("Successfully set default pack configuration", "id", id)
	return nil
}
