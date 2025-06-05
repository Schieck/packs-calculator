package service

// Business constraints: prefer exact quantities (no surplus) and minimal pack count
const (
	minimalSurplus = 0
	minimalPacks   = 1
)

func Calculate(packSizes []int, orderQty int) (map[int]int, int) {
	if orderQty <= 0 {
		return map[int]int{}, 0
	}
	if len(packSizes) == 0 {
		return map[int]int{}, orderQty
	}

	maxPack := packSizes[len(packSizes)-1]
	upper := orderQty + maxPack // ensures at least one feasible quantity

	dp, last := GetDPArraysFromPool(upper)
	defer ReturnDPArraysToPool(CreateDPArrays(dp, last))

	// Unbounded knapsack approach: allows reusing pack sizes multiple times
	for _, p := range packSizes {
		for q := p; q <= upper; q++ {
			if dp[q-p] != maxInt && dp[q-p]+1 < dp[q] {
				dp[q] = dp[q-p] + 1
				last[q] = p
			}
		}
	}

	bestQty := findOptimalQuantity(dp, orderQty, upper)
	if bestQty == -1 {
		return map[int]int{}, orderQty
	}

	alloc := reconstructAllocation(bestQty, last)
	return alloc, bestQty - orderQty
}

func findOptimalQuantity(dp []int, orderQty, upper int) int {
	bestQty, bestSurplus, bestPacks := -1, maxInt, maxInt

	for q := orderQty; q <= upper; q++ {
		if dp[q] == maxInt {
			continue
		}

		surplus := q - orderQty
		if surplus < bestSurplus || (surplus == bestSurplus && dp[q] < bestPacks) {
			bestQty, bestSurplus, bestPacks = q, surplus, dp[q]
			// Early termination when optimal business constraints are met
			if bestSurplus == minimalSurplus && bestPacks == minimalPacks {
				break
			}
		}
	}

	return bestQty
}

func reconstructAllocation(bestQty int, last []int) map[int]int {
	// Pre-allocate with estimated capacity based on typical pack distributions
	alloc := make(map[int]int, 8)

	for q := bestQty; q > 0; {
		p := last[q]
		alloc[p]++
		q -= p
	}

	return alloc
}

// This function is kept for backwards compatibility and direct testing.
func InitializeDPArrays(upper int) ([]int, []int) {
	dp := make([]int, upper+1)
	last := make([]int, upper+1)

	// Initialize DP array with maxInt (more efficient than a loop)
	for i := 1; i <= upper; i++ {
		dp[i] = maxInt
	}
	dp[0] = 0

	return dp, last
}
