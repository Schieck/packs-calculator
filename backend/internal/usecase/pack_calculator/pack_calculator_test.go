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
		shouldError        bool
		errorContains      string
	}{
		// Basic functionality tests
		{
			name:               "Exact match with one pack",
			packSizes:          []int{250, 500, 1000},
			orderQty:           500,
			expectedAllocation: map[int]int{500: 1}, // Exact match with 500 pack
			expectedSurplus:    0,                   // No surplus needed
			shouldError:        false,
		},
		{
			name:               "Best with minimum overage",
			packSizes:          []int{250, 500},
			orderQty:           251,
			expectedAllocation: map[int]int{500: 1}, // Use 500 pack for minimum surplus
			expectedSurplus:    249,                 // 500 - 251 = 249 surplus
			shouldError:        false,
		},
		{
			name:               "Zero order returns empty allocation",
			packSizes:          []int{100, 250},
			orderQty:           0,
			expectedAllocation: map[int]int{},
			expectedSurplus:    0,
			shouldError:        false,
		},
		{
			name:               "No packs available - unfulfillable order",
			packSizes:          []int{},
			orderQty:           100,
			expectedAllocation: map[int]int{},
			expectedSurplus:    100,
			shouldError:        false,
		},
		{
			name:               "Single item order with large packs",
			packSizes:          []int{100, 250, 500},
			orderQty:           1,
			expectedAllocation: map[int]int{100: 1}, // Use smallest pack
			expectedSurplus:    99,                  // 100 - 1 = 99 surplus
			shouldError:        false,
		},
		{
			name:               "Rule R2: Minimize surplus",
			packSizes:          []int{10, 15, 20},
			orderQty:           18,
			expectedAllocation: map[int]int{20: 1}, // 20 gives surplus 2, optimal choice
			expectedSurplus:    2,                  // 20 - 18 = 2
			shouldError:        false,
		},
		{
			name:               "Rule R3: Minimize packs when surplus equal",
			packSizes:          []int{10, 20},
			orderQty:           30,
			expectedAllocation: map[int]int{10: 1, 20: 1}, // Exact match with 1x10 and 1x20
			expectedSurplus:    0,                         // Perfect match
			shouldError:        false,
		},

		// Error cases
		{
			name:          "Negative order quantity should error",
			packSizes:     []int{100, 250},
			orderQty:      -10,
			shouldError:   true,
			errorContains: "order quantity cannot be negative",
		},
		{
			name:          "Invalid pack sizes should error",
			packSizes:     []int{100, -50, 250},
			orderQty:      100,
			shouldError:   true,
			errorContains: "pack sizes must be positive",
		},
		{
			name:          "Zero pack size should error",
			packSizes:     []int{100, 0, 250},
			orderQty:      100,
			shouldError:   true,
			errorContains: "pack sizes must be positive",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Execute use case
			result, err := useCase.Execute(test.packSizes, test.orderQty)

			if test.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.errorContains)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err, "Use case should not return error for test: %s", test.name)
			require.NotNil(t, result, "Result should not be nil")

			// Verify allocation
			actualAllocation := result.Allocation.GetAllocation()
			assert.Equal(t, test.expectedAllocation, actualAllocation,
				"Pack allocation mismatch for test: %s", test.name)

			// Verify surplus
			assert.Equal(t, test.expectedSurplus, result.Surplus,
				"Surplus mismatch for test: %s", test.name)

			// Additional business rule verifications
			if !result.Allocation.IsEmpty() {
				// Verify total items >= order quantity (Rule R1 - only whole packs)
				totalItems := result.Allocation.TotalItems()
				assert.GreaterOrEqual(t, totalItems, test.orderQty,
					"Total items should be >= order quantity (whole packs only)")

				// Verify surplus calculation
				expectedSurplusFromTotal := totalItems - test.orderQty
				assert.Equal(t, expectedSurplusFromTotal, result.Surplus,
					"Surplus should equal total items minus order quantity")
			}
		})
	}
}

func TestCalculatePacksUseCase_ValidateInput(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	calculator := packCalculatorService.NewPackCalculatorService()
	packProcessor := packCalculatorService.NewPackSizeProcessorService()
	useCase := NewCalculatePacksUseCase(calculator, packProcessor, logger)

	t.Run("Valid input should not error", func(t *testing.T) {
		err := useCase.validateInput([]int{100, 250, 500}, 100)
		assert.NoError(t, err)
	})

	t.Run("Negative order quantity should error", func(t *testing.T) {
		err := useCase.validateInput([]int{100, 250}, -1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order quantity cannot be negative")
	})

	t.Run("Zero pack size should error", func(t *testing.T) {
		err := useCase.validateInput([]int{100, 0, 250}, 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pack sizes must be positive")
	})

	t.Run("Negative pack size should error", func(t *testing.T) {
		err := useCase.validateInput([]int{100, -50}, 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pack sizes must be positive")
	})
}
