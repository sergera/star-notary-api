package slc

import (
	"reflect"
	"testing"
)

// TestMap checks if the Map function transforms slices correctly
func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		f        func(int) int
		expected []int
	}{
		{
			name:     "Identity Function",
			input:    []int{1, 2, 3},
			f:        func(i int) int { return i },
			expected: []int{1, 2, 3},
		},
		{
			name:     "Square Numbers",
			input:    []int{1, 2, 3},
			f:        func(i int) int { return i * i },
			expected: []int{1, 4, 9},
		},
		{
			name:     "Empty Slice",
			input:    []int{},
			f:        func(i int) int { return i * i },
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.f)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TestMap %s: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}
