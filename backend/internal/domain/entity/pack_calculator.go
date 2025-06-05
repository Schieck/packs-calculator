package entity

import (
	"fmt"
)

type PackSizes struct {
	sizes []int
	index map[int]struct{}
}

func NewPackSizes(processedSizes []int) *PackSizes {
	if len(processedSizes) == 0 {
		return &PackSizes{sizes: []int{}, index: make(map[int]struct{})}
	}

	index := make(map[int]struct{}, len(processedSizes))
	for _, size := range processedSizes {
		index[size] = struct{}{}
	}

	return &PackSizes{sizes: processedSizes, index: index}
}

func (ps *PackSizes) Slice() []int { return ps.sizes }

func (ps *PackSizes) IsEmpty() bool { return len(ps.sizes) == 0 }

func (ps *PackSizes) Contains(size int) bool {
	_, ok := ps.index[size]
	return ok
}

type OrderQuantity struct {
	Quantity int
}

func NewOrderQuantity(quantity int) (*OrderQuantity, error) {
	if quantity < 0 {
		return nil, fmt.Errorf("order quantity cannot be negative, got %d", quantity)
	}
	return &OrderQuantity{Quantity: quantity}, nil
}

func (oq *OrderQuantity) IsZero() bool {
	return oq.Quantity == 0
}

type PackAllocation struct {
	allocation map[int]int
}

func NewPackAllocation() *PackAllocation {
	return &PackAllocation{
		allocation: make(map[int]int),
	}
}

func (pa *PackAllocation) AddPack(size int, count int) {
	if count > 0 {
		pa.allocation[size] = count
	}
}

func (pa *PackAllocation) GetAllocation() map[int]int {
	result := make(map[int]int)
	for size, count := range pa.allocation {
		result[size] = count
	}
	return result
}

func (pa *PackAllocation) TotalPacks() int {
	total := 0
	for _, count := range pa.allocation {
		total += count
	}
	return total
}

func (pa *PackAllocation) TotalItems() int {
	total := 0
	for size, count := range pa.allocation {
		total += size * count
	}
	return total
}

func (pa *PackAllocation) IsEmpty() bool {
	return len(pa.allocation) == 0
}

type CalculationResult struct {
	Allocation *PackAllocation
	Surplus    int
}

func NewCalculationResult(allocation *PackAllocation, surplus int) *CalculationResult {
	return &CalculationResult{
		Allocation: allocation,
		Surplus:    surplus,
	}
}

func (cr *CalculationResult) IsExactMatch() bool {
	return cr.Surplus == 0
}

func (cr *CalculationResult) HasSurplus() bool {
	return cr.Surplus > 0
}

func (cr *CalculationResult) IsUnfulfillable() bool {
	return cr.Surplus > 0 && cr.Allocation.IsEmpty()
}

type PackCalculator interface {
	CalculateOptimalPacks(packSizes *PackSizes, orderQuantity *OrderQuantity) *CalculationResult
}

// PackSizeProcessor handles business logic for processing raw pack size input
type PackSizeProcessor interface {
	ProcessPackSizes(rawSizes []int) (*PackSizes, error)
}
