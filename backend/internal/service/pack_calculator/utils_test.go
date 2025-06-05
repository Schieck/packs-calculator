package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedupeAndSort(t *testing.T) {
	t.Run("Basic functionality", func(t *testing.T) {
		input := []int{500, 250, 250, 1000, 500}
		expected := []int{250, 500, 1000}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Already sorted and unique", func(t *testing.T) {
		input := []int{100, 250, 500}
		expected := []int{100, 250, 500}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Reverse sorted with duplicates", func(t *testing.T) {
		input := []int{1000, 500, 500, 250, 250, 100}
		expected := []int{100, 250, 500, 1000}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Filter out non-positive values", func(t *testing.T) {
		input := []int{-100, 0, 250, -50, 500}
		expected := []int{250, 500}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("All non-positive values", func(t *testing.T) {
		input := []int{-100, 0, -50}
		expected := []int{}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Empty input", func(t *testing.T) {
		input := []int{}

		result := DedupeAndSort(input)

		assert.Nil(t, result)
	})

	t.Run("Single positive value", func(t *testing.T) {
		input := []int{250}
		expected := []int{250}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Single non-positive value", func(t *testing.T) {
		input := []int{-250}
		expected := []int{}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("All same positive values", func(t *testing.T) {
		input := []int{250, 250, 250, 250}
		expected := []int{250}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Mixed with zeros", func(t *testing.T) {
		input := []int{0, 100, 0, 200, 0}
		expected := []int{100, 200}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Large numbers", func(t *testing.T) {
		input := []int{5000, 2000, 10000, 2000, 1}
		expected := []int{1, 2000, 5000, 10000}

		result := DedupeAndSort(input)

		assert.Equal(t, expected, result)
	})
}

func TestIsValidPackSize(t *testing.T) {
	t.Run("Positive values are valid", func(t *testing.T) {
		testCases := []int{1, 100, 250, 500, 1000, 5000}

		for _, size := range testCases {
			assert.True(t, IsValidPackSize(size), "Size %d should be valid", size)
		}
	})

	t.Run("Non-positive values are invalid", func(t *testing.T) {
		testCases := []int{0, -1, -100, -250}

		for _, size := range testCases {
			assert.False(t, IsValidPackSize(size), "Size %d should be invalid", size)
		}
	})
}

func TestFilterValidPackSizes(t *testing.T) {
	t.Run("Filter mixed valid and invalid", func(t *testing.T) {
		input := []int{-100, 0, 250, -50, 500, 1000}
		expected := []int{250, 500, 1000}

		result := FilterValidPackSizes(input)

		assert.Equal(t, expected, result)
	})

	t.Run("All valid sizes", func(t *testing.T) {
		input := []int{100, 250, 500}
		expected := []int{100, 250, 500}

		result := FilterValidPackSizes(input)

		assert.Equal(t, expected, result)
	})

	t.Run("All invalid sizes", func(t *testing.T) {
		input := []int{-100, 0, -250}
		expected := []int{}

		result := FilterValidPackSizes(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Empty input", func(t *testing.T) {
		input := []int{}
		expected := []int{}

		result := FilterValidPackSizes(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Single valid size", func(t *testing.T) {
		input := []int{250}
		expected := []int{250}

		result := FilterValidPackSizes(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Single invalid size", func(t *testing.T) {
		input := []int{-250}
		expected := []int{}

		result := FilterValidPackSizes(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Preserve order of valid sizes", func(t *testing.T) {
		input := []int{1000, -100, 250, 0, 500}
		expected := []int{1000, 250, 500}

		result := FilterValidPackSizes(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Preserve duplicates", func(t *testing.T) {
		input := []int{250, 250, -100, 500, 500}
		expected := []int{250, 250, 500, 500}

		result := FilterValidPackSizes(input)

		assert.Equal(t, expected, result)
	})
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
