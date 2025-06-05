package service

import (
	"sort"
)

// DedupeAndSort prepares pack sizes for the DP algorithm which requires:
// - No duplicates: prevents redundant calculations in the DP loops
// - Positive values only: negative/zero pack sizes have no business meaning
// - Ascending order: enables algorithm optimizations and predictable results
func DedupeAndSort(in []int) []int {
	if len(in) == 0 {
		return nil
	}

	uniq := make(map[int]struct{}, len(in))
	out := make([]int, 0, len(in))

	// Map-based deduplication is O(n) vs O(n²) for slice-based approaches
	for _, v := range in {
		if v > 0 {
			uniq[v] = struct{}{}
		}
	}

	for v := range uniq {
		out = append(out, v)
	}

	sort.Ints(out)
	return out
}

func IsValidPackSize(size int) bool {
	return size > 0
}

func FilterValidPackSizes(sizes []int) []int {
	result := make([]int, 0, len(sizes))
	for _, size := range sizes {
		if IsValidPackSize(size) {
			result = append(result, size)
		}
	}
	return result
}
