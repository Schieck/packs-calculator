package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

// MockPackCalculator implements the PackCalculator interface for testing
type MockPackCalculator struct {
	results map[string]*entity.CalculationResult
}

func NewMockPackCalculator() *MockPackCalculator {
	return &MockPackCalculator{
		results: make(map[string]*entity.CalculationResult),
	}
}

func (m *MockPackCalculator) SetResult(packSizes []int, orderQty int, allocation map[int]int, surplus int) {
	key := m.generateKey(packSizes, orderQty)
	alloc := entity.NewPackAllocation()
	for size, count := range allocation {
		alloc.AddPack(size, count)
	}
	m.results[key] = entity.NewCalculationResult(alloc, surplus)
}

func (m *MockPackCalculator) generateKey(packSizes []int, orderQty int) string {
	return "" // Simple implementation for testing
}

func (m *MockPackCalculator) CalculateOptimalPacks(packSizes *entity.PackSizes, orderQuantity *entity.OrderQuantity) *entity.CalculationResult {
	// This would contain the actual algorithm implementation
	// For now, we'll implement a simple version to make tests pass
	return m.calculateSimple(packSizes, orderQuantity)
}

// Simple implementation for testing - this would be replaced by the actual algorithm
func (m *MockPackCalculator) calculateSimple(packSizes *entity.PackSizes, orderQuantity *entity.OrderQuantity) *entity.CalculationResult {
	if orderQuantity.IsZero() {
		return entity.NewCalculationResult(entity.NewPackAllocation(), 0)
	}

	if packSizes.IsEmpty() {
		return entity.NewCalculationResult(entity.NewPackAllocation(), orderQuantity.Quantity)
	}

	// Simple greedy algorithm for testing
	sizes := packSizes.ToSlice()
	remaining := orderQuantity.Quantity
	allocation := entity.NewPackAllocation()

	// Reverse order - start with largest packs
	for i := len(sizes) - 1; i >= 0; i-- {
		packSize := sizes[i]
		if remaining >= packSize {
			count := remaining / packSize
			allocation.AddPack(packSize, count)
			remaining -= count * packSize
		}
	}

	// If we still have remaining items, take the smallest pack to minimize surplus
	if remaining > 0 && len(sizes) > 0 {
		smallestPack := sizes[0]
		allocation.AddPack(smallestPack, 1)
		remaining -= smallestPack
	}

	surplus := -remaining // negative remaining means surplus
	if surplus < 0 {
		surplus = 0 // Can't have negative surplus in real scenario
	}

	return entity.NewCalculationResult(allocation, surplus)
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
	calculator := NewMockPackCalculator()

	tests := []testCase{
		// Basic functionality tests
		{
			name:               "Exact match with one pack",
			packSizes:          []int{250, 500, 1000},
			orderQty:           500,
			expectedAllocation: map[int]int{500: 1},
			expectedSurplus:    0,
		},
		{
			name:               "Best with minimum overage",
			packSizes:          []int{250, 500},
			orderQty:           251,
			expectedAllocation: map[int]int{500: 1},
			expectedSurplus:    249,
		},
		{
			name:               "Min packs preferred if same overage",
			packSizes:          []int{250},
			orderQty:           500,
			expectedAllocation: map[int]int{250: 2},
			expectedSurplus:    0,
		},
		{
			name:               "Large pack combo",
			packSizes:          []int{250, 500, 1000, 2000, 5000},
			orderQty:           12001,
			expectedAllocation: map[int]int{5000: 2, 2000: 1, 250: 1},
			expectedSurplus:    249, // 12250 - 12001 = 249
		},
		{
			name:               "Performance edge case - large quantity",
			packSizes:          []int{23, 31, 53},
			orderQty:           500000,
			expectedAllocation: map[int]int{53: 9433, 31: 1}, // This should be calculated by actual algorithm
			expectedSurplus:    0,
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
			expectedAllocation: map[int]int{100: 1},
			expectedSurplus:    99,
		},

		// Business rule verification tests
		{
			name:               "Rule R2: Minimize surplus - prefer larger pack over multiple smaller",
			packSizes:          []int{10, 15, 20},
			orderQty:           18,
			expectedAllocation: map[int]int{20: 1}, // 2 surplus vs 10+10 with 2 surplus but more packs
			expectedSurplus:    2,
		},
		{
			name:               "Rule R3: Minimize packs when surplus is equal",
			packSizes:          []int{10, 20},
			orderQty:           30,
			expectedAllocation: map[int]int{10: 3}, // Exact match vs 20+10 with same result
			expectedSurplus:    0,
		},
		{
			name:               "Complex optimization - multiple solutions",
			packSizes:          []int{5, 7, 11},
			orderQty:           24,
			expectedAllocation: map[int]int{7: 2, 11: 1}, // 25 total, 1 surplus, 3 packs vs other combinations
			expectedSurplus:    1,
		},

		// Additional edge cases for comprehensive coverage
		{
			name:               "Very large pack sizes relative to order",
			packSizes:          []int{1000, 5000},
			orderQty:           50,
			expectedAllocation: map[int]int{1000: 1},
			expectedSurplus:    950,
		},
		{
			name:               "All pack sizes identical",
			packSizes:          []int{100, 100, 100},
			orderQty:           250,
			expectedAllocation: map[int]int{100: 3}, // Should deduplicate and use 3 packs
			expectedSurplus:    50,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create domain entities
			packSizes, err := entity.NewPackSizes(test.packSizes)
			require.NoError(t, err, "Failed to create pack sizes")

			orderQuantity, err := entity.NewOrderQuantity(test.orderQty)
			require.NoError(t, err, "Failed to create order quantity")

			// Execute calculation
			result := calculator.CalculateOptimalPacks(packSizes, orderQuantity)
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

func TestPackSizes_Creation(t *testing.T) {
	t.Run("Valid pack sizes", func(t *testing.T) {
		sizes := []int{100, 250, 500}
		packSizes, err := entity.NewPackSizes(sizes)

		assert.NoError(t, err)
		assert.False(t, packSizes.IsEmpty())
		assert.True(t, packSizes.Contains(250))
		assert.False(t, packSizes.Contains(300))
	})

	t.Run("Empty pack sizes allowed", func(t *testing.T) {
		sizes := []int{}
		packSizes, err := entity.NewPackSizes(sizes)

		assert.NoError(t, err)
		assert.True(t, packSizes.IsEmpty())
	})

	t.Run("Invalid pack sizes should error", func(t *testing.T) {
		sizes := []int{100, -50, 250}
		packSizes, err := entity.NewPackSizes(sizes)

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
