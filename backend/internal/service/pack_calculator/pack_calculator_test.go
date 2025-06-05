package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

func TestPackCalculatorService_NewPackCalculatorService(t *testing.T) {
	service := NewPackCalculatorService()
	assert.NotNil(t, service)
	assert.IsType(t, &PackCalculatorService{}, service)
}

func TestPackCalculatorService_CalculateOptimalPacks_EdgeCases(t *testing.T) {
	service := NewPackCalculatorService()

	t.Run("Zero order quantity", func(t *testing.T) {
		packSizes, err := entity.NewPackSizes([]int{100, 250, 500})
		require.NoError(t, err)

		orderQuantity, err := entity.NewOrderQuantity(0)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQuantity)

		assert.NotNil(t, result)
		assert.True(t, result.Allocation.IsEmpty())
		assert.Equal(t, 0, result.Surplus)
		assert.True(t, result.IsExactMatch())
	})

	t.Run("Empty pack sizes", func(t *testing.T) {
		packSizes, err := entity.NewPackSizes([]int{})
		require.NoError(t, err)

		orderQuantity, err := entity.NewOrderQuantity(100)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQuantity)

		assert.NotNil(t, result)
		assert.True(t, result.Allocation.IsEmpty())
		assert.Equal(t, 100, result.Surplus)
		assert.True(t, result.IsUnfulfillable())
	})

	t.Run("Both zero order and empty packs", func(t *testing.T) {
		packSizes, err := entity.NewPackSizes([]int{})
		require.NoError(t, err)

		orderQuantity, err := entity.NewOrderQuantity(0)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQuantity)

		assert.NotNil(t, result)
		assert.True(t, result.Allocation.IsEmpty())
		assert.Equal(t, 0, result.Surplus)
		assert.True(t, result.IsExactMatch())
	})
}

func TestPackCalculatorService_CalculateOptimalPacks_ValidScenarios(t *testing.T) {
	service := NewPackCalculatorService()

	t.Run("Exact match scenario", func(t *testing.T) {
		packSizes, err := entity.NewPackSizes([]int{250, 500})
		require.NoError(t, err)

		orderQuantity, err := entity.NewOrderQuantity(750)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQuantity)

		assert.NotNil(t, result)
		assert.False(t, result.Allocation.IsEmpty())
		assert.Equal(t, 0, result.Surplus)
		assert.True(t, result.IsExactMatch())
		assert.Equal(t, 750, result.Allocation.TotalItems())
	})

	t.Run("Minimal surplus scenario", func(t *testing.T) {
		packSizes, err := entity.NewPackSizes([]int{250, 500})
		require.NoError(t, err)

		orderQuantity, err := entity.NewOrderQuantity(251)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQuantity)

		assert.NotNil(t, result)
		assert.False(t, result.Allocation.IsEmpty())
		assert.Equal(t, 249, result.Surplus) // 500 - 251
		assert.True(t, result.HasSurplus())
		assert.False(t, result.IsExactMatch())
	})
}

func TestCalculateOptimalPacks_FreeFunction(t *testing.T) {
	t.Run("Basic functionality", func(t *testing.T) {
		allocation, surplus := CalculateOptimalPacks([]int{250, 500}, 750)

		expectedAllocation := map[int]int{250: 1, 500: 1}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("With duplicates and unsorted input", func(t *testing.T) {
		allocation, surplus := CalculateOptimalPacks([]int{500, 250, 250, 500}, 750)

		expectedAllocation := map[int]int{250: 1, 500: 1}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("With invalid pack sizes", func(t *testing.T) {
		allocation, surplus := CalculateOptimalPacks([]int{-100, 0, 250, 500}, 750)

		expectedAllocation := map[int]int{250: 1, 500: 1}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("Empty pack sizes", func(t *testing.T) {
		allocation, surplus := CalculateOptimalPacks([]int{}, 100)

		expectedAllocation := map[int]int{}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 100, surplus)
	})
}
