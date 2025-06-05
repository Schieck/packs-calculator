package service

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDPArraysFromPool(t *testing.T) {
	t.Run("Basic functionality", func(t *testing.T) {
		requiredSize := 100

		dp, last := GetDPArraysFromPool(requiredSize)

		// Verify size
		assert.Equal(t, requiredSize+1, len(dp))
		assert.Equal(t, requiredSize+1, len(last))

		// Verify initialization
		assert.Equal(t, 0, dp[0])
		for i := 1; i <= requiredSize; i++ {
			assert.Equal(t, maxInt, dp[i])
			assert.Equal(t, 0, last[i])
		}
	})

	t.Run("Large size allocation", func(t *testing.T) {
		requiredSize := poolInitialCapacity * 2

		dp, last := GetDPArraysFromPool(requiredSize)

		assert.Equal(t, requiredSize+1, len(dp))
		assert.Equal(t, requiredSize+1, len(last))
		assert.Equal(t, 0, dp[0])
		assert.Equal(t, maxInt, dp[1])
	})

	t.Run("Small size reuses pool", func(t *testing.T) {
		requiredSize := poolInitialCapacity / 2

		dp, last := GetDPArraysFromPool(requiredSize)

		// Should reuse existing pool arrays with proper resizing
		assert.GreaterOrEqual(t, cap(dp), requiredSize+1)
		assert.GreaterOrEqual(t, cap(last), requiredSize+1)
		assert.Equal(t, requiredSize+1, len(dp))
		assert.Equal(t, requiredSize+1, len(last))
	})

	t.Run("Zero size", func(t *testing.T) {
		requiredSize := 0

		dp, last := GetDPArraysFromPool(requiredSize)

		assert.Equal(t, 1, len(dp))
		assert.Equal(t, 1, len(last))
		assert.Equal(t, 0, dp[0])
		assert.Equal(t, 0, last[0])
	})
}

func TestReturnDPArraysToPool(t *testing.T) {
	t.Run("Return normal-sized arrays", func(t *testing.T) {
		arrays := &dpArrays{
			dp:   make([]int, poolInitialCapacity),
			last: make([]int, poolInitialCapacity),
		}

		// This should not panic and should return arrays to pool
		ReturnDPArraysToPool(arrays)
		// We can't easily verify the array was returned, but we can verify no panic
	})

	t.Run("Don't return oversized arrays", func(t *testing.T) {
		oversized := poolInitialCapacity * 20
		arrays := &dpArrays{
			dp:   make([]int, oversized),
			last: make([]int, oversized),
		}

		// This should not panic and should NOT return arrays to pool (let GC handle)
		ReturnDPArraysToPool(arrays)
		// Arrays should be discarded, not returned to pool
	})
}

func TestCreateDPArrays(t *testing.T) {
	t.Run("Wrapper creation", func(t *testing.T) {
		dp := []int{0, maxInt, maxInt}
		last := []int{0, 0, 100}

		arrays := CreateDPArrays(dp, last)

		assert.NotNil(t, arrays)
		assert.Equal(t, dp, arrays.dp)
		assert.Equal(t, last, arrays.last)
	})
}

func TestDPArrayPool_Concurrency(t *testing.T) {
	t.Run("Concurrent access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		iterations := 100

		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					requiredSize := 100 + j // Vary the size

					dp, last := GetDPArraysFromPool(requiredSize)

					// Verify arrays are properly initialized
					assert.Equal(t, requiredSize+1, len(dp))
					assert.Equal(t, requiredSize+1, len(last))
					assert.Equal(t, 0, dp[0])
					if requiredSize > 0 {
						assert.Equal(t, maxInt, dp[1])
					}

					// Return to pool
					ReturnDPArraysToPool(CreateDPArrays(dp, last))
				}
			}()
		}

		wg.Wait()
	})
}

func TestDPArrays_Struct(t *testing.T) {
	t.Run("Struct creation and access", func(t *testing.T) {
		dp := []int{0, 1, 2}
		last := []int{0, 100, 200}

		arrays := dpArrays{dp: dp, last: last}

		assert.Equal(t, dp, arrays.dp)
		assert.Equal(t, last, arrays.last)
	})
}

func TestPoolConstants(t *testing.T) {
	t.Run("Constants are defined correctly", func(t *testing.T) {
		assert.Equal(t, 1024, poolInitialCapacity)
		assert.Greater(t, maxInt, 1000000) // Should be a very large number
	})
}

// Benchmark tests to verify pooling performance
func BenchmarkGetDPArraysFromPool(b *testing.B) {
	b.Run("With pooling", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dp, last := GetDPArraysFromPool(500)
			ReturnDPArraysToPool(CreateDPArrays(dp, last))
		}
	})

	b.Run("Without pooling", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			InitializeDPArrays(500) // Direct allocation without pooling
		}
	})
}
