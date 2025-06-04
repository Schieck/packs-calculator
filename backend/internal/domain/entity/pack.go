package entity

import (
	"fmt"
	"sort"
)

// PackSize represents a fixed pack size that can be shipped
type PackSize struct {
	Size int
}

// NewPackSize creates a new PackSize with validation
func NewPackSize(size int) (*PackSize, error) {
	if size <= 0 {
		return nil, fmt.Errorf("pack size must be positive, got %d", size)
	}
	return &PackSize{Size: size}, nil
}

// PackSizes represents a collection of available pack sizes
type PackSizes struct {
	sizes []PackSize
}

// NewPackSizes creates a new PackSizes collection
func NewPackSizes(sizes []int) (*PackSizes, error) {
	if len(sizes) == 0 {
		return &PackSizes{sizes: []PackSize{}}, nil
	}

	packSizes := make([]PackSize, 0, len(sizes))
	for _, size := range sizes {
		if size <= 0 {
			return nil, fmt.Errorf("all pack sizes must be positive, got %d", size)
		}
		packSizes = append(packSizes, PackSize{Size: size})
	}

	// Sort sizes for optimal calculation
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i].Size < packSizes[j].Size
	})

	return &PackSizes{sizes: packSizes}, nil
}

// ToSlice returns the pack sizes as a slice of integers
func (ps *PackSizes) ToSlice() []int {
	result := make([]int, len(ps.sizes))
	for i, size := range ps.sizes {
		result[i] = size.Size
	}
	return result
}

// IsEmpty returns true if no pack sizes are available
func (ps *PackSizes) IsEmpty() bool {
	return len(ps.sizes) == 0
}

// Contains checks if a specific pack size exists
func (ps *PackSizes) Contains(size int) bool {
	for _, packSize := range ps.sizes {
		if packSize.Size == size {
			return true
		}
	}
	return false
}

// OrderQuantity represents the quantity of items ordered
type OrderQuantity struct {
	Quantity int
}

// NewOrderQuantity creates a new OrderQuantity with validation
func NewOrderQuantity(quantity int) (*OrderQuantity, error) {
	if quantity < 0 {
		return nil, fmt.Errorf("order quantity cannot be negative, got %d", quantity)
	}
	return &OrderQuantity{Quantity: quantity}, nil
}

// IsZero returns true if the order quantity is zero
func (oq *OrderQuantity) IsZero() bool {
	return oq.Quantity == 0
}

// PackAllocation represents how many packs of each size to ship
type PackAllocation struct {
	allocation map[int]int
}

// NewPackAllocation creates a new PackAllocation
func NewPackAllocation() *PackAllocation {
	return &PackAllocation{
		allocation: make(map[int]int),
	}
}

// AddPack adds a pack of the specified size to the allocation
func (pa *PackAllocation) AddPack(size int, count int) {
	if count > 0 {
		pa.allocation[size] = count
	}
}

// GetAllocation returns the pack allocation map
func (pa *PackAllocation) GetAllocation() map[int]int {
	// Return a copy to prevent external modification
	result := make(map[int]int)
	for size, count := range pa.allocation {
		result[size] = count
	}
	return result
}

// TotalPacks returns the total number of packs in this allocation
func (pa *PackAllocation) TotalPacks() int {
	total := 0
	for _, count := range pa.allocation {
		total += count
	}
	return total
}

// TotalItems returns the total number of items in this allocation
func (pa *PackAllocation) TotalItems() int {
	total := 0
	for size, count := range pa.allocation {
		total += size * count
	}
	return total
}

// IsEmpty returns true if no packs are allocated
func (pa *PackAllocation) IsEmpty() bool {
	return len(pa.allocation) == 0
}

// CalculationResult represents the result of a pack calculation
type CalculationResult struct {
	Allocation *PackAllocation
	Surplus    int
}

// NewCalculationResult creates a new CalculationResult
func NewCalculationResult(allocation *PackAllocation, surplus int) *CalculationResult {
	return &CalculationResult{
		Allocation: allocation,
		Surplus:    surplus,
	}
}

// IsExactMatch returns true if there's no surplus
func (cr *CalculationResult) IsExactMatch() bool {
	return cr.Surplus == 0
}

// HasSurplus returns true if there are surplus items
func (cr *CalculationResult) HasSurplus() bool {
	return cr.Surplus > 0
}

// IsUnfulfillable returns true if the order cannot be fulfilled
func (cr *CalculationResult) IsUnfulfillable() bool {
	return cr.Surplus > 0 && cr.Allocation.IsEmpty()
}

// PackCalculator defines the interface for pack calculation
type PackCalculator interface {
	CalculateOptimalPacks(packSizes *PackSizes, orderQuantity *OrderQuantity) *CalculationResult
}
