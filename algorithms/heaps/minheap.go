package heaps

import "container/heap"

type maxHeap []int

var _ (heap.Interface) = (*maxHeap)(nil)

func (h maxHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h maxHeap) Len() int {
	return len(h)
}

func (h *maxHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *maxHeap) Push(x any) {
	item := x.(int)
	(*h) = append(*h, item)
}

func (h *maxHeap) Pop() any {
	n := len(*h)
	rv := (*h)[n-1]
	(*h) = (*h)[:n-1]
	return rv
}

func newMaxHeap(arr []int) maxHeap {
	newH := maxHeap(arr)
	heap.Init(&newH)
	return newH
}

func lastStoneWeight(stones []int) int {
	h := newMaxHeap(stones)

	for h.Len() > 1 {
		x, y := heap.Pop(&h).(int), heap.Pop(&h).(int)
		if x == y {
			continue
		}
		var remainder int
		if x > y {
			remainder = x - y
		} else {
			remainder = y - x
		}
		heap.Push(&h, remainder)
	}
	if h.Len() == 1 {
		return h[0]
	}
	return 0
}
