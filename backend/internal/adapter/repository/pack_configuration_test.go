package repository

import (
	"math"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntSliceToInt64Array(t *testing.T) {
	tests := []struct {
		name        string
		input       []int
		expected    pq.Int64Array
		expectError bool
	}{
		{
			name:        "valid positive integers",
			input:       []int{23, 31, 53},
			expected:    pq.Int64Array{23, 31, 53},
			expectError: false,
		},
		{
			name:        "empty slice",
			input:       []int{},
			expected:    pq.Int64Array{},
			expectError: false,
		},
		{
			name:        "single element",
			input:       []int{100},
			expected:    pq.Int64Array{100},
			expectError: false,
		},
		{
			name:        "negative number should fail",
			input:       []int{-5},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "max int32 should work",
			input:       []int{math.MaxInt32},
			expected:    pq.Int64Array{math.MaxInt32},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := intSliceToInt64Array(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestInt64ArrayToIntSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       pq.Int64Array
		expected    []int
		expectError bool
	}{
		{
			name:        "valid positive integers",
			input:       pq.Int64Array{23, 31, 53},
			expected:    []int{23, 31, 53},
			expectError: false,
		},
		{
			name:        "empty array",
			input:       pq.Int64Array{},
			expected:    []int{},
			expectError: false,
		},
		{
			name:        "single element",
			input:       pq.Int64Array{100},
			expected:    []int{100},
			expectError: false,
		},
		{
			name:        "negative number should fail",
			input:       pq.Int64Array{-5},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "max int32 should work",
			input:       pq.Int64Array{math.MaxInt32},
			expected:    []int{math.MaxInt32},
			expectError: false,
		},
		{
			name:        "value exceeding int32 should fail",
			input:       pq.Int64Array{math.MaxInt32 + 1},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := int64ArrayToIntSlice(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	original := []int{23, 31, 53, 100, 250, 500, 1000}

	// Convert to int64 array
	int64Array, err := intSliceToInt64Array(original)
	require.NoError(t, err)

	// Convert back to int slice
	result, err := int64ArrayToIntSlice(int64Array)
	require.NoError(t, err)

	// Should be identical to original
	assert.Equal(t, original, result)
}
