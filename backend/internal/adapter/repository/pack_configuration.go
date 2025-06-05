package repository

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
	"github.com/lib/pq"
)

type PackConfigurationRepository struct {
	db *sql.DB
}

func NewPackConfigurationRepository(db *sql.DB) *PackConfigurationRepository {
	return &PackConfigurationRepository{
		db: db,
	}
}

// When using[]int with Postgres we need to convert to pq.Int64Array
func intSliceToInt64Array(intSlice []int) (pq.Int64Array, error) {
	int64Array := make(pq.Int64Array, len(intSlice))
	for i, val := range intSlice {
		if val < 0 || val > math.MaxInt32 {
			return nil, fmt.Errorf("pack size %d is out of valid range", val)
		}
		int64Array[i] = int64(val)
	}
	return int64Array, nil
}

func int64ArrayToIntSlice(int64Array pq.Int64Array) ([]int, error) {
	intSlice := make([]int, len(int64Array))
	for i, val := range int64Array {
		if val < 0 || val > math.MaxInt32 {
			return nil, fmt.Errorf("database pack size %d is out of valid range", val)
		}
		intSlice[i] = int(val)
	}
	return intSlice, nil
}

func (r *PackConfigurationRepository) scanPackConfiguration(scanner interface {
	Scan(dest ...interface{}) error
}) (*entity.PackConfiguration, error) {
	config := &entity.PackConfiguration{}
	var packSizes pq.Int64Array

	err := scanner.Scan(
		&config.ID,
		&config.Name,
		&packSizes,
		&config.IsDefault,
		&config.IsActive,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	config.PackSizes, err = int64ArrayToIntSlice(packSizes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pack sizes: %w", err)
	}

	return config, nil
}

func (r *PackConfigurationRepository) GetAll() ([]*entity.PackConfiguration, error) {
	query := `
		SELECT id, name, pack_sizes, is_default, is_active, created_at, updated_at 
		FROM pack_configurations 
		WHERE is_active = true
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pack configurations: %w", err)
	}
	defer rows.Close()

	var configs []*entity.PackConfiguration
	for rows.Next() {
		config, err := r.scanPackConfiguration(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pack configuration: %w", err)
		}
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return configs, nil
}

func (r *PackConfigurationRepository) GetByID(id int) (*entity.PackConfiguration, error) {
	query := `
		SELECT id, name, pack_sizes, is_default, is_active, created_at, updated_at 
		FROM pack_configurations 
		WHERE id = $1 AND is_active = true
	`

	config, err := r.scanPackConfiguration(r.db.QueryRow(query, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pack configuration with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get pack configuration: %w", err)
	}

	return config, nil
}

func (r *PackConfigurationRepository) GetDefault() (*entity.PackConfiguration, error) {
	query := `
		SELECT id, name, pack_sizes, is_default, is_active, created_at, updated_at 
		FROM pack_configurations 
		WHERE is_default = true AND is_active = true
		LIMIT 1
	`

	config, err := r.scanPackConfiguration(r.db.QueryRow(query))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no default pack configuration found")
		}
		return nil, fmt.Errorf("failed to get default pack configuration: %w", err)
	}

	return config, nil
}

func (r *PackConfigurationRepository) Create(config *entity.PackConfiguration) (*entity.PackConfiguration, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO pack_configurations (name, pack_sizes, is_default, is_active) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
	`

	packSizes, err := intSliceToInt64Array(config.PackSizes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pack sizes: %w", err)
	}

	err = r.db.QueryRow(query, config.Name, packSizes, config.IsDefault, config.IsActive).Scan(
		&config.ID,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pack configuration: %w", err)
	}

	return config, nil
}

func (r *PackConfigurationRepository) Update(config *entity.PackConfiguration) (*entity.PackConfiguration, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	config.UpdatedAt = time.Now()

	query := `
		UPDATE pack_configurations 
		SET name = $1, pack_sizes = $2, is_default = $3, is_active = $4, updated_at = $5
		WHERE id = $6 AND is_active = true
		RETURNING created_at
	`

	packSizes, err := intSliceToInt64Array(config.PackSizes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pack sizes: %w", err)
	}

	err = r.db.QueryRow(query,
		config.Name,
		packSizes,
		config.IsDefault,
		config.IsActive,
		config.UpdatedAt,
		config.ID,
	).Scan(&config.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pack configuration with id %d not found or inactive", config.ID)
		}
		return nil, fmt.Errorf("failed to update pack configuration: %w", err)
	}

	return config, nil
}

func (r *PackConfigurationRepository) Delete(id int) error {
	query := `UPDATE pack_configurations SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND is_active = true`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pack configuration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pack configuration with id %d not found or already inactive", id)
	}

	return nil
}

func (r *PackConfigurationRepository) SetDefault(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, unset all existing defaults
	_, err = tx.Exec(`UPDATE pack_configurations SET is_default = false, updated_at = CURRENT_TIMESTAMP WHERE is_default = true`)
	if err != nil {
		return fmt.Errorf("failed to unset existing defaults: %w", err)
	}

	// Then set the new default
	result, err := tx.Exec(`UPDATE pack_configurations SET is_default = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND is_active = true`, id)
	if err != nil {
		return fmt.Errorf("failed to set new default: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pack configuration with id %d not found or inactive", id)
	}

	return tx.Commit()
}
