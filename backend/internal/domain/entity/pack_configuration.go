package entity

import (
	"fmt"
	"time"
)

type PackConfiguration struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	PackSizes []int     `db:"pack_sizes" json:"pack_sizes"`
	IsDefault bool      `db:"is_default" json:"is_default"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func NewPackConfiguration(name string, packSizes []int) (*PackConfiguration, error) {
	if name == "" {
		return nil, fmt.Errorf("pack configuration name cannot be empty")
	}

	if len(packSizes) == 0 {
		return nil, fmt.Errorf("pack configuration must have at least one pack size")
	}

	for _, size := range packSizes {
		if size <= 0 {
			return nil, fmt.Errorf("pack size must be positive, got %d", size)
		}
	}

	return &PackConfiguration{
		Name:      name,
		PackSizes: packSizes,
		IsDefault: false,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (pc *PackConfiguration) Validate() error {
	if pc.Name == "" {
		return fmt.Errorf("pack configuration name cannot be empty")
	}

	if len(pc.PackSizes) == 0 {
		return fmt.Errorf("pack configuration must have at least one pack size")
	}

	for _, size := range pc.PackSizes {
		if size <= 0 {
			return fmt.Errorf("pack size must be positive, got %d", size)
		}
	}

	return nil
}

func (pc *PackConfiguration) GetRawPackSizes() []int {
	return pc.PackSizes
}

type PackConfigurationRepository interface {
	GetAll() ([]*PackConfiguration, error)
	GetByID(id int) (*PackConfiguration, error)
	GetDefault() (*PackConfiguration, error)
	Create(config *PackConfiguration) (*PackConfiguration, error)
	Update(config *PackConfiguration) (*PackConfiguration, error)
	Delete(id int) error
	SetDefault(id int) error
}
