package main_test

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	Keys     []int
	Children []*Node // len(Children) == len(Keys)+1 for internal nodes, 0 for leaves
	Leaf     bool
}

type BTree struct {
	Root *Node
	T    int // minimum degree: nodes hold T-1 .. 2T-1 keys (root may hold fewer)
}

func NewBTree(t int) *BTree {
	return &BTree{Root: &Node{Leaf: true}, T: t}
}

// full reports whether a node has the maximum number of keys (2T-1).
func (b *BTree) full(n *Node) bool {
	return len(n.Keys) == 2*b.T-1
}

// ---------------------------------------------------------------------------
// Insert
// ---------------------------------------------------------------------------

// Insert adds key to the tree. Duplicates are ignored.
//
// The strategy is "split on the way down": before stepping into any node,
// make sure it isn't full. That way, if we ever need to push a key up into
// a parent, the parent is guaranteed to have room, and we never have to
// walk back up the tree.
func (b *BTree) Insert(key int) {
	if b.Search(key) {
		return
	}

	// Special case: the root is full. Nothing sits above the root to receive
	// its middle key, so create a new empty root, hang the old root under it,
	// and split. This is the ONLY place the tree grows taller.
	if b.full(b.Root) {
		oldRoot := b.Root
		newRoot := &Node{Leaf: false, Children: []*Node{oldRoot}}
		b.splitChild(newRoot, 0)
		b.Root = newRoot
	}

	b.insertNonFull(b.Root, key)
}

// insertNonFull inserts key into the subtree rooted at n.
// Precondition: n is NOT full, so it can always accept one more key.
func (b *BTree) insertNonFull(n *Node, key int) {
	// i is the first position whose key is >= key. In a leaf, that's where
	// the new key goes. In an internal node, Children[i] is the subtree
	// whose range contains key.
	i := sort.SearchInts(n.Keys, key)

	if n.Leaf {
		// Leaves are where keys actually get added. Since n isn't full,
		// there's room: just slot it in at position i.
		n.Keys = insertAt(n.Keys, i, key)
		return
	}

	// Internal node: we're about to step down into Children[i].
	// If that child is full, split it FIRST, while we still have its
	// parent (n) in hand. n isn't full, so it has room for the median.
	if b.full(n.Children[i]) {
		b.splitChild(n, i)

		// The split moved a median key up into n.Keys[i], and the child is
		// now two nodes: Children[i] (smaller keys) and Children[i+1]
		// (larger keys). Pick the half that key belongs in.
		if key > n.Keys[i] {
			i++
		}
	}

	b.insertNonFull(n.Children[i], key)
}

// splitChild splits the full child parent.Children[i] into two nodes and
// moves its middle key up into parent.
// Precondition: parent is not full, and parent.Children[i] is full.
//
// With T = 2, a full node has 3 keys:
//
//	parent:   [ ... ]                 parent:  [ ...  20  ... ]
//	              |          ==>                    /      \
//	child:   [10 20 30]                         [10]      [30]
//	                                            left      right
//
// In general a full node has 2T-1 keys: T-1 go left, 1 (the median) goes
// up, and T-1 go right. Each half ends up with exactly the minimum.
func (b *BTree) splitChild(parent *Node, i int) {
	t := b.T
	left := parent.Children[i]
	median := left.Keys[t-1]

	// Build the right half from the upper keys (and upper children).
	// We COPY into fresh slices on purpose: if right shared left's backing
	// array, a later append to left could silently overwrite right's data.
	right := &Node{
		Leaf: left.Leaf,
		Keys: append([]int(nil), left.Keys[t:]...),
	}
	if !left.Leaf {
		// A full internal node has 2T children: the first T stay left,
		// the last T go right.
		right.Children = append([]*Node(nil), left.Children[t:]...)
	}

	// Shrink the left half in place to the lower T-1 keys (and T children).
	left.Keys = left.Keys[:t-1]
	if !left.Leaf {
		left.Children = left.Children[:t]
	}

	// Hook everything into the parent: median at Keys[i], new right node
	// just after the left one at Children[i+1].
	parent.Keys = insertAt(parent.Keys, i, median)
	parent.Children = insertAt(parent.Children, i+1, right)
}

// insertAt returns s with v inserted at index i, shifting later items right.
func insertAt[E any](s []E, i int, v E) []E {
	var zero E
	s = append(s, zero)  // grow by one
	copy(s[i+1:], s[i:]) // shift s[i:] right by one
	s[i] = v
	return s
}

// ---------------------------------------------------------------------------
// Search and helpers
// ---------------------------------------------------------------------------

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

// String prints the tree one level per line, e.g.
//
//	[20 40]
//	[10] [30] [50 60]
func (b *BTree) String() string {
	var sb strings.Builder
	level := []*Node{b.Root}
	for len(level) > 0 {
		var next []*Node
		for j, n := range level {
			if j > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprint(&sb, n.Keys)
			next = append(next, n.Children...)
		}
		sb.WriteString("\n")
		level = next
	}
	return sb.String()
}

func main() {
	b := NewBTree(2) // max 3 keys per node, so splits happen often

	for _, k := range []int{10, 20, 30, 40, 50, 60, 70, 80, 90} {
		b.Insert(k)
		fmt.Printf("after inserting %d:\n%s\n", k, b)
	}
}
