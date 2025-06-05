package usecase

import (
	"fmt"
	"log/slog"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

type CalculatePacksUseCase struct {
	calculator    entity.PackCalculator
	packProcessor entity.PackSizeProcessor
	logger        *slog.Logger
}

func NewCalculatePacksUseCase(calculator entity.PackCalculator, packProcessor entity.PackSizeProcessor, logger *slog.Logger) *CalculatePacksUseCase {
	return &CalculatePacksUseCase{
		calculator:    calculator,
		packProcessor: packProcessor,
		logger:        logger,
	}
}

func (uc *CalculatePacksUseCase) Execute(packSizes []int, orderQuantity int) (*entity.CalculationResult, error) {
	uc.logger.Info("Executing pack calculation use case",
		"order_quantity", orderQuantity,
		"pack_sizes", packSizes)

	if err := uc.validateInput(packSizes, orderQuantity); err != nil {
		uc.logger.Warn("Pack calculation input validation failed", "error", err)
		return nil, err
	}

	packSizesEntity, err := uc.packProcessor.ProcessPackSizes(packSizes)
	if err != nil {
		uc.logger.Error("Failed to process pack sizes", "error", err)
		return nil, err
	}

	orderQuantityEntity, err := entity.NewOrderQuantity(orderQuantity)
	if err != nil {
		uc.logger.Error("Failed to create order quantity entity", "error", err)
		return nil, err
	}

	result := uc.calculator.CalculateOptimalPacks(packSizesEntity, orderQuantityEntity)
	if result == nil {
		uc.logger.Error("Calculator returned nil result")
		return nil, fmt.Errorf("calculation failed")
	}

	uc.logger.Info("Pack calculation completed successfully",
		"total_packs", result.Allocation.TotalPacks(),
		"total_items", result.Allocation.TotalItems(),
		"surplus", result.Surplus)

	return result, nil
}

func (uc *CalculatePacksUseCase) validateInput(packSizes []int, orderQuantity int) error {
	if orderQuantity < 0 {
		return fmt.Errorf("order quantity cannot be negative")
	}

	for _, size := range packSizes {
		if size <= 0 {
			return fmt.Errorf("pack sizes must be positive")
		}
	}

	return nil
}
