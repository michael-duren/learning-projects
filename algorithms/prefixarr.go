package algorithms

import "fmt"

type NumMatrix struct {
	sumMat [][]int
}

func ConstructorMatrix(matrix [][]int) NumMatrix {
	rows, cols := len(matrix), len(matrix[0])
	sumMat := make([][]int, rows+1)
	for i := range sumMat {
		sumMat[i] = make([]int, cols+1)
	}
	for r := 0; r < rows; r++ {
		prefix := 0
		for c := 0; c < cols; c++ {
			prefix += matrix[r][c]
			above := sumMat[r][c+1]
			sumMat[r+1][c+1] = prefix + above
		}
	}

	return NumMatrix{sumMat: sumMat}
}

func main() {
	m := [][]int{
		{3, 2, 1},
		{2, 1, 2},
		{6, 1, 2},
	}
	for _, v := range m {
		for _, vv := range v {
			fmt.Printf(" %d ", vv)
		}
		fmt.Print("\n")
	}
	nm := ConstructorMatrix(m)
	for _, v := range nm.sumMat {
		for _, vv := range v {
			if vv == 0 {
				continue
			}
			fmt.Printf(" %d ", vv)
		}
		fmt.Print("\n")
	}
}
