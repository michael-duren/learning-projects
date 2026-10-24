package algorithms

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
