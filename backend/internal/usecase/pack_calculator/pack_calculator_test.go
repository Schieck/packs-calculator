package usecase

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	packCalculatorService "github.com/Schieck/packs-calculator/internal/service/pack_calculator"
)

func TestCalculatePacksUseCase_Execute(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	calculator := packCalculatorService.NewPackCalculatorService()
	packProcessor := packCalculatorService.NewPackSizeProcessorService()
	useCase := NewCalculatePacksUseCase(calculator, packProcessor, logger)

	tests := []struct {
		name               string
		packSizes          []int
		orderQty           int
		expectedAllocation map[int]int
		expectedSurplus    int
		errorContains      string
	}{
		{
			name:               "exact match scenario",
			packSizes:          []int{250, 500, 1000, 2000, 5000},
			orderQty:           250,
			expectedAllocation: map[int]int{250: 1},
			expectedSurplus:    0,
		},
		{
			name:               "combination pack scenario",
			packSizes:          []int{250, 500, 1000, 2000, 5000},
			orderQty:           750,
			expectedAllocation: map[int]int{250: 1, 500: 1},
			expectedSurplus:    0,
		},
		{
			name:               "surplus scenario",
			packSizes:          []int{250, 500, 1000, 2000, 5000},
			orderQty:           251,
			expectedAllocation: map[int]int{500: 1},
			expectedSurplus:    249, // 500 - 251
		},
		{
			name:               "large order scenario",
			packSizes:          []int{250, 500, 1000, 2000, 5000},
			orderQty:           12001,
			expectedAllocation: map[int]int{250: 1, 5000: 2, 2000: 1},
			expectedSurplus:    249, // 12250 - 12001
		},
		{
			name:               "zero order quantity",
			packSizes:          []int{250, 500},
			orderQty:           0,
			expectedAllocation: map[int]int{},
			expectedSurplus:    0,
		},
		{
			name:               "large pack for small order",
			packSizes:          []int{500},
			orderQty:           100,
			expectedAllocation: map[int]int{500: 1},
			expectedSurplus:    400, // 500 - 100 = 400
		},
		{
			name:          "negative order quantity",
			packSizes:     []int{250, 500},
			orderQty:      -1,
			errorContains: "order quantity cannot be negative",
		},
		{
			name:          "negative pack size",
			packSizes:     []int{-250, 500},
			orderQty:      100,
			errorContains: "pack sizes must be positive",
		},
		{
			name:          "zero pack size",
			packSizes:     []int{0, 500},
			orderQty:      100,
			errorContains: "pack sizes must be positive",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := useCase.Execute(test.packSizes, test.orderQty)

			if test.errorContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.errorContains)
				require.Nil(t, result)
			} else {
				require.NoError(t, err, "Use case should not return error for test: %s", test.name)
				require.NotNil(t, result, "Result should not be nil")

				// Convert the allocation result to a comparable map
				actualAllocation := result.Allocation.GetAllocation()
				assert.Equal(t, test.expectedAllocation, actualAllocation,
					"Allocation mismatch for test: %s", test.name)

				assert.Equal(t, test.expectedSurplus, result.Surplus,
					"Surplus mismatch for test: %s", test.name)

				// Validate business rule: total items >= order quantity
				totalItems := result.Allocation.TotalItems()
				assert.GreaterOrEqual(t, totalItems, test.orderQty,
					"Total items should be >= order quantity for test: %s", test.name)

				// Cross-validate surplus calculation
				expectedSurplusFromTotal := totalItems - test.orderQty
				assert.Equal(t, expectedSurplusFromTotal, result.Surplus,
					"Surplus calculation mismatch for test: %s", test.name)
			}
		})
	}
}

func TestCalculatePacksUseCase_ValidateInput(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	calculator := packCalculatorService.NewPackCalculatorService()
	packProcessor := packCalculatorService.NewPackSizeProcessorService()
	useCase := NewCalculatePacksUseCase(calculator, packProcessor, logger)

	t.Run("valid input should pass", func(t *testing.T) {
		t.Parallel()
		err := useCase.validateInput([]int{250, 500}, 100)
		require.NoError(t, err)
	})

	t.Run("negative order quantity should fail", func(t *testing.T) {
		t.Parallel()
		err := useCase.validateInput([]int{250, 500}, -1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "order quantity cannot be negative")
	})

	t.Run("negative pack size should fail", func(t *testing.T) {
		t.Parallel()
		err := useCase.validateInput([]int{-250, 500}, 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pack sizes must be positive")
	})

	t.Run("zero pack size should fail", func(t *testing.T) {
		t.Parallel()
		err := useCase.validateInput([]int{0, 500}, 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pack sizes must be positive")
	})
}
