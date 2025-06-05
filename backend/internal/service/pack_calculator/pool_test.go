package service

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDPArraysFromPool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		upperBound int
		wantSize   int
	}{
		{
			name:       "small upper bound",
			upperBound: 100,
			wantSize:   101, // upperBound + 1
		},
		{
			name:       "medium upper bound",
			upperBound: 1000,
			wantSize:   1001,
		},
		{
			name:       "large upper bound",
			upperBound: 10000,
			wantSize:   10001,
		},
		{
			name:       "zero upper bound",
			upperBound: 0,
			wantSize:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dp, last := GetDPArraysFromPool(tt.upperBound)

			// Verify the returned arrays have correct size
			assert.Equal(t, tt.wantSize, len(dp))
			assert.Equal(t, tt.wantSize, len(last))

			// Verify arrays are properly initialized
			assert.Equal(t, 0, dp[0])
			assert.Equal(t, 0, last[0])

			for i := 1; i < len(dp); i++ {
				assert.Equal(t, maxInt, dp[i])
				assert.Equal(t, 0, last[i])
			}

			// Return to pool to avoid leaking
			ReturnDPArraysToPool(CreateDPArrays(dp, last))
		})
	}
}

func TestReturnDPArraysToPool(t *testing.T) {
	t.Parallel()

	t.Run("return valid arrays", func(t *testing.T) {
		t.Parallel()

		// Get arrays from pool
		dp, last := GetDPArraysFromPool(100)

		// Modify some values to verify they get reset
		dp[1] = 5
		dp[2] = 10
		last[1] = 25
		last[2] = 50

		// Return to pool
		ReturnDPArraysToPool(CreateDPArrays(dp, last))

		// Get arrays again and verify they are reset
		newDP, newLast := GetDPArraysFromPool(100)
		assert.Equal(t, 0, newDP[0])
		assert.Equal(t, maxInt, newDP[1])
		assert.Equal(t, maxInt, newDP[2])
		assert.Equal(t, 0, newLast[1])
		assert.Equal(t, 0, newLast[2])

		// Clean up
		ReturnDPArraysToPool(CreateDPArrays(newDP, newLast))
	})
}

func TestCreateDPArrays(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dp   []int
		last []int
	}{
		{
			name: "normal case",
			dp:   []int{0, maxInt, maxInt},
			last: []int{0, 0, 100},
		},
		{
			name: "empty arrays",
			dp:   []int{},
			last: []int{},
		},
		{
			name: "single element",
			dp:   []int{0},
			last: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dpArrays := CreateDPArrays(tt.dp, tt.last)

			assert.Equal(t, tt.dp, dpArrays.dp)
			assert.Equal(t, tt.last, dpArrays.last)
		})
	}
}

func TestDPArrayPool_Concurrency(t *testing.T) {
	// Note: Not using t.Parallel() here since this test specifically tests concurrency behavior
	const (
		numGoroutines = 100
		upperBound    = 1000
	)

	var wg sync.WaitGroup
	results := make(chan *dpArrays, numGoroutines)

	// Launch multiple goroutines that get arrays from pool
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dp, last := GetDPArraysFromPool(upperBound)

			// Verify arrays are properly initialized
			assert.Equal(t, 0, dp[0])
			assert.Equal(t, 0, last[0])
			for j := 1; j < len(dp); j++ {
				assert.Equal(t, maxInt, dp[j])
				assert.Equal(t, 0, last[j])
			}

			results <- CreateDPArrays(dp, last)
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Return all arrays to pool
	for dpArrays := range results {
		ReturnDPArraysToPool(dpArrays)
	}

	// Verify pool is working correctly after concurrent access
	dp, last := GetDPArraysFromPool(upperBound)
	assert.Equal(t, upperBound+1, len(dp))
	assert.Equal(t, upperBound+1, len(last))
	ReturnDPArraysToPool(CreateDPArrays(dp, last))
}

func TestDPArrays_Struct(t *testing.T) {
	t.Parallel()

	t.Run("verify struct fields", func(t *testing.T) {
		t.Parallel()

		dpArrays := &dpArrays{
			dp:   []int{1, 2, 3},
			last: []int{4, 5, 6},
		}

		assert.Equal(t, []int{1, 2, 3}, dpArrays.dp)
		assert.Equal(t, []int{4, 5, 6}, dpArrays.last)
	})
}

func TestPoolConstants(t *testing.T) {
	t.Parallel()

	t.Run("verify maxInt constant", func(t *testing.T) {
		t.Parallel()

		// maxInt should be a very large positive integer
		assert.Greater(t, maxInt, 1000000)
		assert.Positive(t, maxInt)
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
