package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

func createPackSizes(rawSizes []int) (*entity.PackSizes, error) {
	processor := NewPackSizeProcessorService()
	return processor.ProcessPackSizes(rawSizes)
}

func TestPackCalculatorService_NewPackCalculatorService(t *testing.T) {
	t.Parallel()

	service := NewPackCalculatorService()
	assert.NotNil(t, service)
	assert.IsType(t, &PackCalculatorService{}, service)
}

func TestPackCalculatorService_CalculateOptimalPacks_EdgeCases(t *testing.T) {
	t.Parallel()

	service := NewPackCalculatorService()

	t.Run("zero order quantity", func(t *testing.T) {
		t.Parallel()
		packSizes, err := createPackSizes([]int{250, 500})
		require.NoError(t, err)

		orderQty, err := entity.NewOrderQuantity(0)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQty)

		require.NotNil(t, result)
		assert.True(t, result.Allocation.IsEmpty())
		assert.Equal(t, 0, result.Surplus)
		assert.True(t, result.IsExactMatch())
	})

	t.Run("large pack for small order", func(t *testing.T) {
		t.Parallel()
		packSizes, err := createPackSizes([]int{500})
		require.NoError(t, err)

		orderQty, err := entity.NewOrderQuantity(100)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQty)

		require.NotNil(t, result)
		assert.False(t, result.Allocation.IsEmpty())
		assert.Equal(t, 400, result.Surplus) // 500 - 100 = 400
		assert.True(t, result.HasSurplus())
	})

	t.Run("empty pack sizes", func(t *testing.T) {
		t.Parallel()
		packSizes, err := createPackSizes([]int{})
		require.NoError(t, err)

		orderQty, err := entity.NewOrderQuantity(100)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQty)

		require.NotNil(t, result)
		assert.True(t, result.Allocation.IsEmpty())
		assert.Equal(t, 100, result.Surplus) // Cannot fulfill, so surplus equals original order
		assert.True(t, result.IsUnfulfillable())
	})
}

func TestPackCalculatorService_CalculateOptimalPacks_ValidScenarios(t *testing.T) {
	t.Parallel()

	service := NewPackCalculatorService()

	t.Run("exact match scenario", func(t *testing.T) {
		t.Parallel()
		packSizes, err := createPackSizes([]int{250, 500})
		require.NoError(t, err)

		orderQty, err := entity.NewOrderQuantity(750)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQty)

		require.NotNil(t, result)
		assert.False(t, result.Allocation.IsEmpty())
		assert.Equal(t, 0, result.Surplus)
		assert.True(t, result.IsExactMatch())
		assert.Equal(t, 750, result.Allocation.TotalItems())
	})

	t.Run("surplus scenario", func(t *testing.T) {
		t.Parallel()
		packSizes, err := createPackSizes([]int{250, 500})
		require.NoError(t, err)

		orderQty, err := entity.NewOrderQuantity(251)
		require.NoError(t, err)

		result := service.CalculateOptimalPacks(packSizes, orderQty)

		require.NotNil(t, result)
		assert.False(t, result.Allocation.IsEmpty())
		assert.Equal(t, 249, result.Surplus) // 500 - 251
		assert.True(t, result.HasSurplus())
		assert.False(t, result.IsExactMatch())
	})
}

func TestCalculateOptimalPacks_FreeFunction(t *testing.T) {
	t.Parallel()

	t.Run("basic calculation", func(t *testing.T) {
		t.Parallel()
		packSizes := []int{250, 500}
		orderQty := 750

		allocation, surplus := CalculateOptimalPacks(packSizes, orderQty)

		expectedAllocation := map[int]int{250: 1, 500: 1}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("exact match large", func(t *testing.T) {
		t.Parallel()
		packSizes := []int{100, 200, 500}
		orderQty := 1000

		allocation, surplus := CalculateOptimalPacks(packSizes, orderQty)

		expectedAllocation := map[int]int{500: 2}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 0, surplus)
	})

	t.Run("surplus calculation", func(t *testing.T) {
		t.Parallel()
		packSizes := []int{250, 500}
		orderQty := 251

		allocation, surplus := CalculateOptimalPacks(packSizes, orderQty)

		expectedAllocation := map[int]int{500: 1}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 249, surplus) // 500 - 251 = 249
	})

	t.Run("large pack scenario", func(t *testing.T) {
		t.Parallel()
		packSizes := []int{500}
		orderQty := 100

		allocation, surplus := CalculateOptimalPacks(packSizes, orderQty)

		expectedAllocation := map[int]int{500: 1}
		assert.Equal(t, expectedAllocation, allocation)
		assert.Equal(t, 400, surplus) // 500 - 100 = 400
	})
}
