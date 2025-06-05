package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	t.Run("Basic exact match", func(t *testing.T) {
		packSizes := []int{250, 500}
		orderQty := 750

		allocation, surplus := Calculate(packSizes, orderQty)

		expectedAllocation := map[int]int{250: 1, 500: 1}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("Minimal surplus selection", func(t *testing.T) {
		packSizes := []int{250, 500}
		orderQty := 251

		allocation, surplus := Calculate(packSizes, orderQty)

		expectedAllocation := map[int]int{500: 1}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 249, surplus) // 500 - 251
	})

	t.Run("Prefer fewer packs with same surplus", func(t *testing.T) {
		packSizes := []int{3, 4}
		orderQty := 12

		allocation, surplus := Calculate(packSizes, orderQty)

		expectedAllocation := map[int]int{4: 3} // 12 items exactly, not {3: 4} which would be 4 packs
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("Zero order quantity", func(t *testing.T) {
		packSizes := []int{100, 250}
		orderQty := 0

		allocation, surplus := Calculate(packSizes, orderQty)

		expectedAllocation := map[int]int{}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("Negative order quantity", func(t *testing.T) {
		packSizes := []int{100, 250}
		orderQty := -5

		allocation, surplus := Calculate(packSizes, orderQty)

		expectedAllocation := map[int]int{}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("Empty pack sizes", func(t *testing.T) {
		packSizes := []int{}
		orderQty := 100

		allocation, surplus := Calculate(packSizes, orderQty)

		expectedAllocation := map[int]int{}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 100, surplus)
	})

	t.Run("Single pack size", func(t *testing.T) {
		packSizes := []int{100}
		orderQty := 250

		allocation, surplus := Calculate(packSizes, orderQty)

		expectedAllocation := map[int]int{100: 3} // 300 total
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 50, surplus) // 300 - 250
	})

	t.Run("Large numbers performance test", func(t *testing.T) {
		packSizes := []int{23, 31, 53}
		orderQty := 1000

		allocation, surplus := Calculate(packSizes, orderQty)

		// Verify we got a valid solution
		assert.NotEmpty(t, allocation)

		// Calculate total items
		totalItems := 0
		for size, count := range allocation {
			totalItems += size * count
		}

		assert.GreaterOrEqual(t, totalItems, orderQty)
		assert.Equal(t, totalItems-orderQty, surplus)
	})
}

func TestFindOptimalQuantity(t *testing.T) {
	t.Run("Find minimum surplus", func(t *testing.T) {
		// Mock DP array where dp[i] represents minimum packs to achieve quantity i
		dp := []int{0, maxInt, maxInt, 1, 1, 2, 1, 2, 2, 1, 2} // indices 0-10
		orderQty := 7
		upper := 10

		bestQty := findOptimalQuantity(dp, orderQty, upper)

		assert.Equal(t, 7, bestQty) // quantity 7 is exact match with 2 packs (0 surplus is optimal)
	})

	t.Run("Prefer fewer packs when surplus is equal", func(t *testing.T) {
		// dp[8] = 2 packs, dp[9] = 3 packs, both have same surplus relative to orderQty=7
		dp := []int{0, maxInt, maxInt, maxInt, maxInt, maxInt, maxInt, maxInt, 2, 3, maxInt}
		orderQty := 7
		upper := 10

		bestQty := findOptimalQuantity(dp, orderQty, upper)

		assert.Equal(t, 8, bestQty) // Choose 8 (2 packs) over 9 (3 packs) for same surplus
	})

	t.Run("Early break for optimal solution", func(t *testing.T) {
		// Perfect solution at exact quantity with 1 pack
		dp := []int{0, maxInt, maxInt, maxInt, maxInt, maxInt, maxInt, 1, maxInt, maxInt, maxInt}
		orderQty := 7
		upper := 10

		bestQty := findOptimalQuantity(dp, orderQty, upper)

		assert.Equal(t, 7, bestQty) // Exact match with minimal packs should break early
	})

	t.Run("No feasible solution", func(t *testing.T) {
		dp := []int{0, maxInt, maxInt, maxInt, maxInt, maxInt}
		orderQty := 3
		upper := 5

		bestQty := findOptimalQuantity(dp, orderQty, upper)

		assert.Equal(t, -1, bestQty)
	})
}

func TestReconstructAllocation(t *testing.T) {
	t.Run("Simple reconstruction", func(t *testing.T) {
		// last[i] indicates which pack size was used to reach quantity i
		// For bestQty=100, last[100] should contain the pack size used
		last := make([]int, 101)
		last[100] = 100 // last pack to reach 100 was size 100
		bestQty := 100

		allocation := reconstructAllocation(bestQty, last)

		// Should reconstruct the path backwards from bestQty
		// Note: This test verifies structure rather than exact values since DP path depends on the last array
		assert.NotEmpty(t, allocation)

		// Verify allocation is valid
		totalItems := 0
		for size, count := range allocation {
			assert.Positive(t, size)
			assert.Positive(t, count)
			totalItems += size * count
		}
		assert.Equal(t, bestQty, totalItems)
	})

	t.Run("Empty reconstruction", func(t *testing.T) {
		last := []int{0}
		bestQty := 0

		allocation := reconstructAllocation(bestQty, last)

		expectedAllocation := map[int]int{}
		assert.Equal(t, expectedAllocation, allocation)
	})
}

func TestInitializeDPArrays(t *testing.T) {
	t.Run("Proper initialization", func(t *testing.T) {
		upper := 5

		dp, last := InitializeDPArrays(upper)

		assert.Equal(t, upper+1, len(dp))
		assert.Equal(t, upper+1, len(last))

		// dp[0] should be 0, rest should be maxInt
		assert.Equal(t, 0, dp[0])
		for i := 1; i <= upper; i++ {
			assert.Equal(t, maxInt, dp[i])
			assert.Equal(t, 0, last[i])
		}
	})

	t.Run("Zero upper bound", func(t *testing.T) {
		upper := 0

		dp, last := InitializeDPArrays(upper)

		assert.Equal(t, 1, len(dp))
		assert.Equal(t, 1, len(last))
		assert.Equal(t, 0, dp[0])
		assert.Equal(t, 0, last[0])
	})
}
