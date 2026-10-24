package algorithms

import "math"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func NewTreeNode(val int) *TreeNode {
	return &TreeNode{
		Val: val,
	}
}

// assumes array is sorted
func NewBstFromArr(arr []int) *TreeNode {
	if len(arr) == 0 {
		return nil
	}

	mid := len(arr) / 2
	root := NewTreeNode(arr[mid])
	root.Left = NewBstFromArr(arr[:mid])
	root.Right = NewBstFromArr(arr[mid+1:])

	return root
}

const Null = math.MinInt

// Example input: [3, 9, 20, null, null, 15, 7].
func NewTreeFromLevelOrder(vals []int) *TreeNode {
	if len(vals) == 0 || vals[0] == Null {
		return nil
	}
	root := NewTreeNode(vals[0])
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) && vals[i] != Null {
			node.Left = NewTreeNode(vals[i])
			queue = append(queue, node.Left)
		}
		i++
		if i < len(vals) && vals[i] != Null {
			node.Right = NewTreeNode(vals[i])
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}
