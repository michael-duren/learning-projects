package algorithms

type MinHeap struct {
	arr []int
}

func NewMinHeap() *MinHeap {
	return &MinHeap{arr: []int{}}
}

func (mh *MinHeap) Push(val int) {
	mh.arr = append(mh.arr, val)
	if len(mh.arr) == 1 {
		return
	}

	n := len(mh.arr)
	i := n - 1
	for i > 0 {
		p := (i - 1) / 2
		if mh.arr[i] >= mh.arr[p] {
			break
		}
		mh.arr[i], mh.arr[p] = mh.arr[p], mh.arr[i]
		i = p
	}
}

func (mh *MinHeap) Pop() int {
	if len(mh.arr) == 0 {
		return -1
	}

	tmp := mh.arr[0]
	mh.arr[0] = mh.arr[len(mh.arr)-1]
	mh.arr = mh.arr[:len(mh.arr)-1]

	i := 0
	n := len(mh.arr)
	for i*2+1 < n {
		smallest := i
		l, r := i*2+1, i*2+2
		if mh.arr[smallest] > mh.arr[l] {
			smallest = l
		}
		if r < n && mh.arr[smallest] > mh.arr[r] {
			smallest = r
		}
		if i == smallest {
			break
		}
		mh.arr[i], mh.arr[smallest] = mh.arr[smallest], mh.arr[i]

		i = smallest
	}

	return tmp
}

func (mh *MinHeap) Top() int {
	if len(mh.arr) == 0 {
		return -1
	}
	return mh.arr[0]
}

func (mh *MinHeap) Heapify(nums []int) {
	mh.arr = nums

	n := len(mh.arr)
	curr := (n - 1) / 2

	for curr >= 0 {
		i := curr
		for i*2+1 < n {
			smallest := i
			l, r := i*2+1, i*2+2
			if mh.arr[smallest] > mh.arr[l] {
				smallest = l
			}
			if r < n && mh.arr[smallest] > mh.arr[r] {
				smallest = r
			}
			mh.arr[i], mh.arr[smallest] = mh.arr[smallest], mh.arr[i]
			if smallest == i {
				break
			}
			i = smallest
		}
	}
}
