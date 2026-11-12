package main

import "sort"

// Implement a B-Tree
//
// Design a B-tree of integers with minimum degree t (t >= 2). The tree must keep these invariants after every operation:
//
// Keys within a node are sorted in strictly increasing order.
// Every node except the root has at least t - 1 keys and at most 2t - 1 keys. A non-empty root has at least 1 key.
// An internal node with k keys has exactly k + 1 children. Every key in Children[i] is less than Keys[i], and every key in Children[i+1] is greater than Keys[i].
// All leaves are at the same depth.
//
// Keys are unique. Inserting a key that already exists does nothing.

type Node struct {
	Keys     []int
	Children []*Node
	Leaf     bool
}

type BTree struct {
	Root *Node
	T    int // minimum degree, t >= 2
}

func NewNode(t int, leaf bool) *Node {
	return &Node{
		Leaf:     leaf,
		Keys:     make([]int, 0, t),
		Children: make([]*Node, 0, t),
	}
}

func NewBTree(t int) *BTree {
	return &BTree{
		T:    t,
		Root: NewNode(t, true),
	}
}

func FromSlice(t int, keys []int) *BTree {
	b := NewBTree(t)
	for _, k := range keys {
		b.Insert(k)
	}
	return b
}

func (b *BTree) full(n *Node) bool {
	return len(n.Keys) == 2*b.T-1
}

func (b *BTree) Search(key int) bool {
	n := b.Root
	for {
		i := sort.SearchInts(n.Keys, key)
		if i < len(n.Keys) && n.Keys[i] == key {
			return true
		}
		if n.Leaf {
			return false
		}
		n = n.Children[i]
	}

}
func (b *BTree) Insert(key int) {
	if b.Search(key) {
		return
	}
}
func (b *BTree) Validate() error
func (b *BTree) InOrder() []int
func (b *BTree) Delete(key int)
