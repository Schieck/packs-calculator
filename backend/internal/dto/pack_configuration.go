package dto

import (
	"time"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

type CreatePackConfigurationRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=255" example:"Standard Packs"`
	PackSizes []int  `json:"pack_sizes" validate:"required,min=1,dive,min=1" swaggertype:"array,integer" example:"250,500,1000"`
}

type UpdatePackConfigurationRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=255" example:"Updated Standard Packs"`
	PackSizes []int  `json:"pack_sizes" validate:"required,min=1,dive,min=1" swaggertype:"array,integer" example:"250,500,1000,2000"`
	IsDefault bool   `json:"is_default" example:"false"`
}

type SetDefaultPackConfigurationRequest struct {
	IsDefault bool `json:"is_default" validate:"required" example:"true"`
}

type PackConfigurationResponse struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"Main Edge Case"`
	PackSizes []int     `json:"pack_sizes" swaggertype:"array,integer" example:"23,31,53"`
	IsDefault bool      `json:"is_default" example:"true"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type PackConfigurationListResponse struct {
	Configurations []*PackConfigurationResponse `json:"configurations"`
	Count          int                          `json:"count" example:"3"`
}

func ToPackConfigurationResponse(config *entity.PackConfiguration) *PackConfigurationResponse {
	return &PackConfigurationResponse{
		ID:        config.ID,
		Name:      config.Name,
		PackSizes: config.PackSizes,
		IsDefault: config.IsDefault,
		IsActive:  config.IsActive,
		CreatedAt: config.CreatedAt,
		UpdatedAt: config.UpdatedAt,
	}
}

func ToPackConfigurationListResponse(configs []*entity.PackConfiguration) *PackConfigurationListResponse {
	responses := make([]*PackConfigurationResponse, len(configs))
	for i, config := range configs {
		responses[i] = ToPackConfigurationResponse(config)
	}

	return &PackConfigurationListResponse{
		Configurations: responses,
		Count:          len(responses),
	}
}
