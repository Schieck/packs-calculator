package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

func createPackSizesIntegration(rawSizes []int) (*entity.PackSizes, error) {
	processor := NewPackSizeProcessorService()
	return processor.ProcessPackSizes(rawSizes)
}

type testCase struct {
	name               string
	packSizes          []int
	orderQty           int
	expectedAllocation map[int]int
	expectedSurplus    int
	shouldError        bool
	errorMessage       string
}

func TestPackCalculator_CalculateOptimalPacks(t *testing.T) {
	calculator := NewPackCalculatorService()

	tests := []testCase{
		// Basic functionality tests
		{
			name:               "Exact match with one pack",
			packSizes:          []int{250, 500, 1000},
			orderQty:           500,
			expectedAllocation: map[int]int{500: 1}, // Exact match with 500 pack
			expectedSurplus:    0,                   // No surplus needed
		},
		{
			name:               "Best with minimum overage",
			packSizes:          []int{250, 500},
			orderQty:           251,
			expectedAllocation: map[int]int{500: 1}, // Use 500 pack for minimum surplus
			expectedSurplus:    249,                 // 500 - 251 = 249 surplus
		},
		{
			name:               "Min packs preferred if same overage",
			packSizes:          []int{250},
			orderQty:           500,
			expectedAllocation: map[int]int{250: 2}, // Exact match with 2x250
			expectedSurplus:    0,                   // No surplus needed
		},
		{
			name:               "Large pack combo",
			packSizes:          []int{250, 500, 1000, 2000, 5000},
			orderQty:           12001,
			expectedAllocation: map[int]int{5000: 2, 2000: 1, 250: 1}, // 12250 total, minimal surplus
			expectedSurplus:    249,                                   // 12250 - 12001 = 249
		},
		{
			name:               "Performance edge case - large quantity",
			packSizes:          []int{23, 31, 53},
			orderQty:           500000,
			expectedAllocation: map[int]int{53: 9429, 31: 7, 23: 2}, // Exact match: 499737+217+46=500000
			expectedSurplus:    0,                                   // Perfect exact match
		},

		// Edge cases
		{
			name:               "Zero order returns empty allocation",
			packSizes:          []int{100, 250},
			orderQty:           0,
			expectedAllocation: map[int]int{},
			expectedSurplus:    0,
		},
		{
			name:               "No packs available - unfulfillable order",
			packSizes:          []int{},
			orderQty:           100,
			expectedAllocation: map[int]int{},
			expectedSurplus:    100,
		},
		{
			name:               "Single item order with large packs",
			packSizes:          []int{100, 250, 500},
			orderQty:           1,
			expectedAllocation: map[int]int{100: 1}, // Use smallest pack
			expectedSurplus:    99,                  // 100 - 1 = 99 surplus
		},

		// Business rule verification tests
		{
			name:               "Rule R2: Minimize surplus - prefer larger pack over multiple smaller",
			packSizes:          []int{10, 15, 20},
			orderQty:           18,
			expectedAllocation: map[int]int{20: 1}, // 20 gives surplus 2, better than 2x10 (same surplus but more packs)
			expectedSurplus:    2,                  // 20 - 18 = 2
		},
		{
			name:               "Rule R3: Minimize packs when surplus is equal",
			packSizes:          []int{10, 20},
			orderQty:           30,
			expectedAllocation: map[int]int{10: 1, 20: 1}, // Exact match with 1x10 and 1x20
			expectedSurplus:    0,                         // Perfect match
		},

		// Additional edge cases for comprehensive coverage
		{
			name:               "Very large pack sizes relative to order",
			packSizes:          []int{1000, 5000},
			orderQty:           50,
			expectedAllocation: map[int]int{1000: 1}, // Use smallest available pack
			expectedSurplus:    950,                  // 1000 - 50 = 950
		},
		{
			name:               "All pack sizes identical",
			packSizes:          []int{100, 100, 100},
			orderQty:           250,
			expectedAllocation: map[int]int{100: 3}, // Need 3 packs to exceed order
			expectedSurplus:    50,                  // 300 - 250 = 50
		},
		{
			name:               "Canonical example – 501 order (spec)",
			packSizes:          []int{250, 500, 1000},
			orderQty:           501,
			expectedAllocation: map[int]int{500: 1, 250: 1}, // 751 total, surplus 249
			expectedSurplus:    249,
		},
		{
			name:               "Order not divisible by GCD (must overshoot)",
			packSizes:          []int{4, 6},
			orderQty:           7,
			expectedAllocation: map[int]int{4: 2}, // 8 total
			expectedSurplus:    1,
		},
		{
			name:               "Duplicate & unsorted pack sizes",
			packSizes:          []int{1000, 250, 250, 500},
			orderQty:           750,
			expectedAllocation: map[int]int{500: 1, 250: 1}, // 0 surplus
			expectedSurplus:    0,
		},
		{
			name:               "Tie on surplus, choose fewer packs",
			packSizes:          []int{3, 4},
			orderQty:           12,
			expectedAllocation: map[int]int{4: 3}, // {3:4} is same surplus but 4 packs
			expectedSurplus:    0,
		},
		{
			name:               "Size-1 pack present – huge order stress",
			packSizes:          []int{1, 1000},
			orderQty:           12345,
			expectedAllocation: map[int]int{1000: 12, 1: 345}, // exact fit
			expectedSurplus:    0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			packSizes, err := createPackSizesIntegration(test.packSizes)
			require.NoError(t, err, "Failed to create pack sizes")

			orderQuantity, err := entity.NewOrderQuantity(test.orderQty)
			require.NoError(t, err, "Failed to create order quantity")

			result := calculator.CalculateOptimalPacks(packSizes, orderQuantity)
			require.NotNil(t, result, "Result should not be nil")

			actualAllocation := result.Allocation.GetAllocation()
			assert.Equal(t, test.expectedAllocation, actualAllocation,
				"Pack allocation mismatch for test: %s", test.name)

			assert.Equal(t, test.expectedSurplus, result.Surplus,
				"Surplus mismatch for test: %s", test.name)

			if !result.Allocation.IsEmpty() {
				totalItems := result.Allocation.TotalItems()
				assert.GreaterOrEqual(t, totalItems, test.orderQty,
					"Total items should be >= order quantity (whole packs only)")

				expectedSurplusFromTotal := totalItems - test.orderQty
				assert.Equal(t, expectedSurplusFromTotal, result.Surplus,
					"Surplus should equal total items minus order quantity")
			}
		})
	}
}

func TestPackSizes_Creation(t *testing.T) {
	processor := NewPackSizeProcessorService()

	t.Run("Valid pack sizes", func(t *testing.T) {
		sizes := []int{100, 250, 500}
		packSizes, err := processor.ProcessPackSizes(sizes)

		assert.NoError(t, err)
		assert.False(t, packSizes.IsEmpty())
		assert.True(t, packSizes.Contains(250))
		assert.False(t, packSizes.Contains(300))
	})

	t.Run("Empty pack sizes allowed", func(t *testing.T) {
		sizes := []int{}
		packSizes, err := processor.ProcessPackSizes(sizes)

		assert.NoError(t, err)
		assert.True(t, packSizes.IsEmpty())
	})

	t.Run("Invalid pack sizes should error", func(t *testing.T) {
		sizes := []int{100, -50, 250}
		packSizes, err := processor.ProcessPackSizes(sizes)

		assert.Error(t, err)
		assert.Nil(t, packSizes)
	})
}

func TestOrderQuantity_Creation(t *testing.T) {
	t.Run("Valid order quantities", func(t *testing.T) {
		quantities := []int{0, 1, 100, 500000}

		for _, qty := range quantities {
			orderQty, err := entity.NewOrderQuantity(qty)
			assert.NoError(t, err)
			assert.Equal(t, qty, orderQty.Quantity)
			assert.Equal(t, qty == 0, orderQty.IsZero())
		}
	})

	t.Run("Negative quantities should error", func(t *testing.T) {
		orderQty, err := entity.NewOrderQuantity(-10)

		assert.Error(t, err)
		assert.Nil(t, orderQty)
	})
}

func TestPackAllocation_Methods(t *testing.T) {
	allocation := entity.NewPackAllocation()

	// Test empty allocation
	assert.True(t, allocation.IsEmpty())
	assert.Equal(t, 0, allocation.TotalPacks())
	assert.Equal(t, 0, allocation.TotalItems())

	// Add some packs
	allocation.AddPack(250, 2)
	allocation.AddPack(500, 1)

	// Test populated allocation
	assert.False(t, allocation.IsEmpty())
	assert.Equal(t, 3, allocation.TotalPacks())    // 2 + 1
	assert.Equal(t, 1000, allocation.TotalItems()) // 250*2 + 500*1

	expectedMap := map[int]int{250: 2, 500: 1}
	assert.Equal(t, expectedMap, allocation.GetAllocation())
}

func TestCalculationResult_Methods(t *testing.T) {
	allocation := entity.NewPackAllocation()
	allocation.AddPack(500, 1)

	t.Run("Exact match result", func(t *testing.T) {
		result := entity.NewCalculationResult(allocation, 0)

		assert.True(t, result.IsExactMatch())
		assert.False(t, result.HasSurplus())
		assert.False(t, result.IsUnfulfillable())
	})

	t.Run("Surplus result", func(t *testing.T) {
		result := entity.NewCalculationResult(allocation, 50)

		assert.False(t, result.IsExactMatch())
		assert.True(t, result.HasSurplus())
		assert.False(t, result.IsUnfulfillable())
	})

	t.Run("Unfulfillable result", func(t *testing.T) {
		emptyAllocation := entity.NewPackAllocation()
		result := entity.NewCalculationResult(emptyAllocation, 100)

		assert.False(t, result.IsExactMatch())
		assert.True(t, result.HasSurplus())
		assert.True(t, result.IsUnfulfillable())
	})
}
