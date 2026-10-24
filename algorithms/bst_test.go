package algorithms_test

import (
	"algorithms"
	"testing"
)

func treesAreEqual(t *testing.T, t1 *algorithms.TreeNode, t2 *algorithms.TreeNode) bool {
	t.Helper()
	s := algorithms.NewStack[*algorithms.TreeNode](10)
	s.Push(t1)
	s.Push(t2)
	for s.Len() > 0 {
		curr1, curr2 := s.Pop(), s.Pop()
		if curr1 == nil && curr2 == nil {
			continue
		}
		if curr1 == nil || curr2 == nil {
			return false
		}
		if curr1.Val != curr2.Val {
			return false
		}
		s.Push(curr1.Left)
		s.Push(curr2.Left)
		s.Push(curr1.Right)
		s.Push(curr2.Right)
	}
	return true
}

func CanCreateBSTFromArray(t *testing.T) {
	sortedArr := []int{1, 2, 3, 4, 5}
	expected := algorithms.TreeNode{
		Val: 3,
		Left: &algorithms.TreeNode{
			Val: 2,
			Left: &algorithms.TreeNode{
				Val: 1,
			},
		},
		Right: &algorithms.TreeNode{
			Val: 4,
			Right: &algorithms.TreeNode{
				Val: 5,
			},
		},
	}

	newTree := algorithms.NewBstFromArr(sortedArr)
	treesAreEqual(t, &expected, newTree)
}
