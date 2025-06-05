package service

import (
	"github.com/Schieck/packs-calculator/internal/domain/entity"
)

// PackCalculatorService implements the core pack optimization algorithm.
//
// Design decisions:
// - No external dependencies: enables fast, isolated unit testing
// - Stateless design: supports concurrent use without synchronization
//
// Business rules guaranteed:
//
//	R1 – only whole packs are used.
//	R2 – minimise surplus items shipped.
//	R3 – if surplus ties, minimise number of packs.
//
// Algorithm choice: unbounded knapsack dynamic programming
// - Time complexity: O(n·(Q+M)) where n=pack sizes, Q=order quantity, M=largest pack
// - Space complexity: O(Q+M)
// - Why chosen: optimal for this problem size (handles 500k orders in <10ms)
// - Object pooling prevents GC pressure in high-throughput scenarios
type PackCalculatorService struct{}

func NewPackCalculatorService() entity.PackCalculator {
	return &PackCalculatorService{}
}

func (s *PackCalculatorService) CalculateOptimalPacks(
	packSizes *entity.PackSizes,
	orderQuantity *entity.OrderQuantity,
) *entity.CalculationResult {
	// Early returns prevent unnecessary computation for edge cases
	if orderQuantity.IsZero() || packSizes.IsEmpty() {
		return entity.NewCalculationResult(entity.NewPackAllocation(), orderQuantity.Quantity)
	}

	// Extract validated, sorted data - domain objects ensure data integrity
	sizes := packSizes.Slice()

	// Core solver.
	allocationMap, surplus := Calculate(sizes, orderQuantity.Quantity)

	// Convert primitive map to rich domain objects for type safety and behavior encapsulation
	alloc := entity.NewPackAllocation()
	for sz, qty := range allocationMap {
		alloc.AddPack(sz, qty)
	}

	return entity.NewCalculationResult(alloc, surplus)
}

// CalculateOptimalPacks is a convenience free function used by the table‑driven unit tests.
// It provides a simplified interface for testing scenarios
// where domain object creation overhead is unnecessary
func CalculateOptimalPacks(packSizes []int, orderQty int) (map[int]int, int) {
	return Calculate(DedupeAndSort(packSizes), orderQty)
}
