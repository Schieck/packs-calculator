package service

import (
	"math"
	"sync"
)

const (
	// maxInt chosen over math.MaxInt64 to reduce memory usage (4 bytes vs 8 bytes per element)
	maxInt = math.MaxInt32
	// poolInitialCapacity sized for typical order quantities to minimize reallocations
	poolInitialCapacity = 1024
)

type dpArrays struct {
	dp   []int
	last []int
}

// Object pooling prevents GC pressure in high-throughput scenarios
// where thousands of DP arrays would otherwise be allocated per second
var dpArrayPool = sync.Pool{
	New: func() interface{} {
		return &dpArrays{
			dp:   make([]int, poolInitialCapacity),
			last: make([]int, poolInitialCapacity),
		}
	},
}

// GetDPArraysFromPool retrieves DP arrays from the pool, ensuring proper size and initialization.
func GetDPArraysFromPool(requiredSize int) ([]int, []int) {
	arrays := dpArrayPool.Get().(*dpArrays)

	// Dynamic resizing balances memory efficiency with allocation cost
	if cap(arrays.dp) < requiredSize+1 {
		// Allocate new arrays when pool capacity is insufficient
		arrays.dp = make([]int, requiredSize+1)
		arrays.last = make([]int, requiredSize+1)
	} else {
		// Reuse existing capacity to avoid allocations
		arrays.dp = arrays.dp[:requiredSize+1]
		arrays.last = arrays.last[:requiredSize+1]
	}

	// maxInt initialization represents "unreachable" state in DP algorithm
	for i := 1; i <= requiredSize; i++ {
		arrays.dp[i] = maxInt
		arrays.last[i] = 0
	}
	arrays.dp[0] = 0 // Base case: 0 packs needed for 0 quantity

	return arrays.dp, arrays.last
}

// ReturnDPArraysToPool returns the DP arrays to the pool for reuse.
func ReturnDPArraysToPool(arrays *dpArrays) {
	// Size limit prevents memory bloat from exceptionally large orders
	// while keeping pools effective for typical workloads
	if cap(arrays.dp) <= poolInitialCapacity*10 {
		dpArrayPool.Put(arrays)
	}
	// Large arrays are discarded to prevent pool memory accumulation
}

// CreateDPArrays wraps the arrays for pool return - utility function
func CreateDPArrays(dp, last []int) *dpArrays {
	return &dpArrays{dp: dp, last: last}
}
