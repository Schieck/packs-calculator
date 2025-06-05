package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedupeAndSort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "basic deduplication and sorting",
			input:    []int{500, 250, 500, 250, 1000},
			expected: []int{250, 500, 1000},
		},
		{
			name:     "already sorted with no duplicates",
			input:    []int{100, 200, 300},
			expected: []int{100, 200, 300},
		},
		{
			name:     "reverse sorted with duplicates",
			input:    []int{300, 200, 100, 200},
			expected: []int{100, 200, 300},
		},
		{
			name:     "single element",
			input:    []int{100},
			expected: []int{100},
		},
		{
			name:     "all same elements",
			input:    []int{100, 100, 100},
			expected: []int{100},
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: nil,
		},
		{
			name:     "unsorted with multiple duplicates",
			input:    []int{53, 31, 23, 53, 31, 23, 100},
			expected: []int{23, 31, 53, 100},
		},
		{
			name:     "with zeros (should be filtered out)",
			input:    []int{250, 0, 500, 0},
			expected: []int{250, 500},
		},
		{
			name:     "with negative numbers (should be filtered out)",
			input:    []int{250, -100, 500, -200},
			expected: []int{250, 500},
		},
		{
			name:     "mixed valid and invalid",
			input:    []int{-50, 0, 250, 500, -100, 0, 250},
			expected: []int{250, 500},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := DedupeAndSort(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidPackSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"positive number", 100, true},
		{"zero", 0, false},
		{"negative number", -50, false},
		{"large positive", 99999, true},
		{"one", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := IsValidPackSize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterValidPackSizes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "mixed valid and invalid",
			input:    []int{-100, 0, 250, 500, -50},
			expected: []int{250, 500},
		},
		{
			name:     "all valid",
			input:    []int{100, 250, 500},
			expected: []int{100, 250, 500},
		},
		{
			name:     "all invalid",
			input:    []int{-100, 0, -50},
			expected: []int{},
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single valid",
			input:    []int{100},
			expected: []int{100},
		},
		{
			name:     "single invalid",
			input:    []int{-100},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := FilterValidPackSizes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests
func BenchmarkDedupeAndSort(b *testing.B) {
	input := []int{500, 250, 1000, 250, 2000, 500, 100, 750, 1500, 250}

	b.Run("DedupeAndSort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DedupeAndSort(input)
		}
	})
}

func BenchmarkFilterValidPackSizes(b *testing.B) {
	input := []int{-100, 500, 0, 250, -50, 1000, 250, -200, 2000, 500}

	b.Run("FilterValidPackSizes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FilterValidPackSizes(input)
		}
	})
}
