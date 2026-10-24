package algorithms_test

import (
	"math/rand/v2"
	"slices"
	"testing"

	"algorithms"
)

func TestMergeSort(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected []int
	}{
		{name: "nil", arr: nil, expected: nil},
		{name: "empty", arr: []int{}, expected: []int{}},
		{name: "single", arr: []int{42}, expected: []int{42}},
		{name: "two sorted", arr: []int{1, 2}, expected: []int{1, 2}},
		{name: "two reversed", arr: []int{2, 1}, expected: []int{1, 2}},
		{name: "unsorted", arr: []int{1, 3, 100, 10}, expected: []int{1, 3, 10, 100}},
		{name: "already sorted", arr: []int{1, 2, 3, 4, 5}, expected: []int{1, 2, 3, 4, 5}},
		{name: "reversed", arr: []int{5, 4, 3, 2, 1}, expected: []int{1, 2, 3, 4, 5}},
		{name: "odd length", arr: []int{7, 2, 9, 4, 1}, expected: []int{1, 2, 4, 7, 9}},
		{name: "duplicates", arr: []int{3, 1, 3, 2, 1}, expected: []int{1, 1, 2, 3, 3}},
		{name: "all equal", arr: []int{4, 4, 4, 4}, expected: []int{4, 4, 4, 4}},
		{name: "negatives", arr: []int{-1, 5, -10, 0, 3}, expected: []int{-10, -1, 0, 3, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := slices.Clone(tt.arr)
			algorithms.MergeSort(arr)
			if !slices.Equal(arr, tt.expected) {
				t.Errorf("MergeSort(%v) = %v, want %v", tt.arr, arr, tt.expected)
			}
		})
	}
}

// TestMergeSortRandom compares against slices.Sort on random inputs.
func TestMergeSortRandom(t *testing.T) {
	for range 100 {
		n := rand.IntN(200)
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rand.IntN(1000) - 500
		}

		expected := slices.Clone(arr)
		slices.Sort(expected)

		got := slices.Clone(arr)
		algorithms.MergeSort(got)
		if !slices.Equal(got, expected) {
			t.Fatalf("MergeSort(%v) = %v, want %v", arr, got, expected)
		}
	}
}
