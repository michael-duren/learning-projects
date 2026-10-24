package algorithms

// merge combines the sorted halves arr[p:q] and arr[q:r] into arr[p:r].
func merge(arr []int, p, q, r int) {
	lArr := arr[p:q]
	rArr := arr[q:r]
	lPtr, rPtr := 0, 0
	newArr := make([]int, 0, r-p)

	for range r - p {
		if lPtr < len(lArr) && rPtr < len(rArr) {
			if lArr[lPtr] > rArr[rPtr] {
				newArr = append(newArr, rArr[rPtr])
				rPtr++
			} else {
				newArr = append(newArr, lArr[lPtr])
				lPtr++
			}
		} else if lPtr < len(lArr) {
			newArr = append(newArr, lArr[lPtr])
			lPtr++
		} else {
			newArr = append(newArr, rArr[rPtr])
			rPtr++
		}
	}
	copy(arr[p:r], newArr)
}

// mergeSort sorts the half-open range arr[p:r].
func mergeSort(arr []int, p, r int) {
	// base case: 0 or 1 elements
	if r-p <= 1 {
		return
	}
	q := p + (r-p)/2

	// 1st half
	mergeSort(arr, p, q)
	// 2nd half
	mergeSort(arr, q, r)

	merge(arr, p, q, r)
}

func MergeSort(arr []int) {
	mergeSort(arr, 0, len(arr))
}
