package algorithms_test

import (
	"algorithms"
	"testing"
)

// Null marks a missing child in level order input, e.g. []int{2, Null, 3}.
const Null = algorithms.Null

func treesAreEqual(t1 *algorithms.TreeNode, t2 *algorithms.TreeNode) bool {
	if t1 == nil || t2 == nil {
		return t1 == t2
	}
	return t1.Val == t2.Val &&
		treesAreEqual(t1.Left, t2.Left) &&
		treesAreEqual(t1.Right, t2.Right)
}

// inorder walks the tree without going through the code under test.
func inorder(root *algorithms.TreeNode) []int {
	if root == nil {
		return nil
	}
	out := inorder(root.Left)
	out = append(out, root.Val)
	return append(out, inorder(root.Right)...)
}

func assertTree(t *testing.T, got, want *algorithms.TreeNode) {
	t.Helper()
	if !treesAreEqual(got, want) {
		t.Fatalf("trees differ\n got inorder: %v\nwant inorder: %v", inorder(got), inorder(want))
	}
}

func TestNewBstFromArr(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int // level order
	}{
		{"empty", nil, nil},
		{"single", []int{1}, []int{1}},
		{"two", []int{1, 2}, []int{2, 1}},
		{"odd", []int{1, 2, 3, 4, 5}, []int{3, 2, 5, 1, Null, 4}},
		{"even", []int{1, 2, 3, 4, 5, 6}, []int{4, 2, 6, 1, 3, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := algorithms.NewBstFromArr(tt.arr)
			assertTree(t, got, algorithms.NewTreeFromLevelOrder(tt.want))
		})
	}
}

// func TestInsert(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		tree []int // level order
// 		val  int
// 		want []int // level order
// 	}{
// 		{"into empty", nil, 5, []int{5}},
// 		{"left of root", []int{5}, 3, []int{5, 3}},
// 		{"right of root", []int{5}, 7, []int{5, Null, 7}},
// 		{"deep left", []int{5, 3, 7}, 1, []int{5, 3, 7, 1}},
// 		{"deep right", []int{5, 3, 7}, 9, []int{5, 3, 7, Null, Null, Null, 9}},
// 		{"between", []int{5, 3, 7}, 4, []int{5, 3, 7, Null, 4}},
// 		{"duplicate ignored", []int{5, 3, 7}, 3, []int{5, 3, 7}},
// 		{"negative", []int{0}, -1, []int{0, -1}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			root := algorithms.NewTreeFromLevelOrder(tt.tree)
// 			got := algorithms.Insert(root, tt.val)
// 			assertTree(t, got, algorithms.NewTreeFromLevelOrder(tt.want))
// 		})
// 	}
// }

// func TestInsertKeepsOrder(t *testing.T) {
// 	var root *algorithms.TreeNode
// 	for _, v := range []int{50, 30, 70, 20, 40, 60, 80, 35, 45, 65} {
// 		root = algorithms.Insert(root, v)
// 	}
// 	want := []int{20, 30, 35, 40, 45, 50, 60, 65, 70, 80}
// 	if got := inorder(root); !slices.Equal(got, want) {
// 		t.Fatalf("inorder = %v, want %v", got, want)
// 	}
// }

// func TestSearch(t *testing.T) {
// 	root := algorithms.NewTreeFromLevelOrder([]int{8, 3, 10, 1, 6, Null, 14, Null, Null, 4, 7})
// 	tests := []struct {
// 		name  string
// 		val   int
// 		found bool
// 	}{
// 		{"root", 8, true},
// 		{"leaf", 7, true},
// 		{"inner", 6, true},
// 		{"rightmost", 14, true},
// 		{"missing smaller than all", 0, false},
// 		{"missing larger than all", 15, false},
// 		{"missing in gap", 5, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := algorithms.Search(root, tt.val)
// 			if !tt.found {
// 				if got != nil {
// 					t.Fatalf("Search(%d) = %d, want nil", tt.val, got.Val)
// 				}
// 				return
// 			}
// 			if got == nil || got.Val != tt.val {
// 				t.Fatalf("Search(%d) = %v, want node with %d", tt.val, got, tt.val)
// 			}
// 		})
// 	}
// }

// func TestSearchEmpty(t *testing.T) {
// 	if got := algorithms.Search(nil, 1); got != nil {
// 		t.Fatalf("Search(nil, 1) = %v, want nil", got)
// 	}
// }
//
// func TestMinMax(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		tree     []int // level order
// 		min, max int
// 	}{
// 		{"single", []int{5}, 5, 5},
// 		{"balanced", []int{8, 3, 10, 1, 6, Null, 14}, 1, 14},
// 		{"left only", []int{3, 2, Null, 1}, 1, 3},
// 		{"right only", []int{1, Null, 2, Null, 3}, 1, 3},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			root := algorithms.NewTreeFromLevelOrder(tt.tree)
// 			if got := algorithms.Min(root); got == nil || got.Val != tt.min {
// 				t.Errorf("Min = %v, want %d", got, tt.min)
// 			}
// 			if got := algorithms.Max(root); got == nil || got.Val != tt.max {
// 				t.Errorf("Max = %v, want %d", got, tt.max)
// 			}
// 		})
// 	}
// }
//
// func TestMinMaxEmpty(t *testing.T) {
// 	if got := algorithms.Min(nil); got != nil {
// 		t.Errorf("Min(nil) = %v, want nil", got)
// 	}
// 	if got := algorithms.Max(nil); got != nil {
// 		t.Errorf("Max(nil) = %v, want nil", got)
// 	}
// }
//
// func TestDelete(t *testing.T) {
// 	// a node with two children is replaced by its in-order successor
// 	tests := []struct {
// 		name string
// 		tree []int // level order
// 		val  int
// 		want []int // level order
// 	}{
// 		{"empty", nil, 1, nil},
// 		{"only node", []int{5}, 5, nil},
// 		{"missing", []int{5, 3, 7}, 4, []int{5, 3, 7}},
// 		{"leaf", []int{5, 3, 7}, 3, []int{5, Null, 7}},
// 		{"left child only", []int{5, 3, 7, 1}, 3, []int{5, 1, 7}},
// 		{"right child only", []int{5, 3, 7, Null, 4}, 3, []int{5, 4, 7}},
// 		{"two children", []int{5, 3, 8, Null, Null, 7, 9}, 8, []int{5, 3, 9, Null, Null, 7}},
// 		{"root with two children", []int{5, 3, 7}, 5, []int{7, 3}},
// 		{
// 			"successor has right child",
// 			[]int{5, 3, 10, Null, Null, 7, 12, Null, 8},
// 			5,
// 			[]int{7, 3, 10, Null, Null, 8, 12},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			root := algorithms.NewTreeFromLevelOrder(tt.tree)
// 			got := algorithms.Delete(root, tt.val)
// 			assertTree(t, got, algorithms.NewTreeFromLevelOrder(tt.want))
// 		})
// 	}
// }
//
// func TestDeleteAll(t *testing.T) {
// 	vals := []int{4, 2, 6, 1, 3, 5, 7}
// 	root := algorithms.NewBstFromArr([]int{1, 2, 3, 4, 5, 6, 7})
// 	for i, v := range vals {
// 		root = algorithms.Delete(root, v)
// 		want := slices.Sorted(slices.Values(vals[i+1:]))
// 		if got := inorder(root); !slices.Equal(got, want) {
// 			t.Fatalf("after Delete(%d): inorder = %v, want %v", v, got, want)
// 		}
// 	}
// 	if root != nil {
// 		t.Fatalf("root = %v after deleting everything, want nil", root)
// 	}
// }
//
// func TestInOrder(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		tree []int // level order
// 		want []int
// 	}{
// 		{"empty", nil, nil},
// 		{"single", []int{1}, []int{1}},
// 		{"balanced", []int{4, 2, 6, 1, 3, 5, 7}, []int{1, 2, 3, 4, 5, 6, 7}},
// 		{"skewed", []int{1, Null, 2, Null, 3}, []int{1, 2, 3}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := algorithms.InOrder(algorithms.NewTreeFromLevelOrder(tt.tree))
// 			if len(got) == 0 && len(tt.want) == 0 {
// 				return
// 			}
// 			if !slices.Equal(got, tt.want) {
// 				t.Fatalf("InOrder = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestHeight(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		tree []int // level order
// 		want int
// 	}{
// 		{"empty", nil, 0},
// 		{"single", []int{1}, 1},
// 		{"balanced", []int{4, 2, 6, 1, 3, 5, 7}, 3},
// 		{"skewed", []int{1, Null, 2, Null, 3, Null, 4}, 4},
// 		{"uneven", []int{4, 2, 6, 1}, 3},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := algorithms.Height(algorithms.NewTreeFromLevelOrder(tt.tree))
// 			if got != tt.want {
// 				t.Fatalf("Height = %d, want %d", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestIsValidBST(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		tree []int // level order
// 		want bool
// 	}{
// 		{"empty", nil, true},
// 		{"single", []int{1}, true},
// 		{"valid", []int{4, 2, 6, 1, 3, 5, 7}, true},
// 		{"left child too big", []int{4, 5, 6}, false},
// 		{"right child too small", []int{4, 2, 3}, false},
// 		// 5 is right of 2 but also left of 4, so it must be < 4
// 		{"grandchild breaks ancestor bound", []int{4, 2, 6, 1, 5}, false},
// 		{"duplicate", []int{2, 2}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := algorithms.IsValidBST(algorithms.NewTreeFromLevelOrder(tt.tree))
// 			if got != tt.want {
// 				t.Fatalf("IsValidBST = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
