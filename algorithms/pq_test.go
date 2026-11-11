package algorithms_test

import (
	"algorithms"
	"container/heap"
	"testing"
)

func TestPQ(t *testing.T) {
	t.Run("returns min", func(t *testing.T) {
		arr := []int{2000, 500, 3000000, 30, 2, 1000, 1}
		items := algorithms.PriorityQueue{}
		for i, v := range arr {
			items = append(items, &algorithms.Item{Index: i, Value: v, Priority: v})
		}
		heap.Init(&items)
		rv := items[0].Value
		if rv != 1 {
			t.Fatalf("expected: %d, got: %d", 1, rv)
		}
	})
}
