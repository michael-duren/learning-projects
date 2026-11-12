Here’s a LeetCode-style problem broken into parts. Each part builds on the last, and the follow-ups at the end connect it to system design.

Implement a B-Tree

Design a B-tree of integers with minimum degree t (t >= 2). The tree must keep these invariants after every operation:

Keys within a node are sorted in strictly increasing order.
Every node except the root has at least t - 1 keys and at most 2t - 1 keys. A non-empty root has at least 1 key.
An internal node with k keys has exactly k + 1 children. Every key in Children[i] is less than Keys[i], and every key in Children[i+1] is greater than Keys[i].
All leaves are at the same depth.

Keys are unique. Inserting a key that already exists does nothing.

go
func NewBTree(t int) *BTree
func (b *BTree) Search(key int) bool
func (b *BTree) Insert(key int)
func (b *BTree) Validate() error
func (b *BTree) InOrder() []int
func (b *BTree) Delete(key int)
Part 1: Search (easy)

Return whether key is in the tree. At each node, find the first key >= key. If it’s equal, you’re done. If the node is a leaf, it’s not there. Otherwise descend into the matching child.

Part 2: Insert (medium)

Insert key while keeping all invariants. Use single-pass, top-down insertion: while walking down, split any full child (2t - 1 keys) before entering it. That guarantees there’s always room to push a median key up into the parent, so you never have to walk back up the tree.

Splitting a full node: the middle key moves up into the parent. Keys left of it become one node, keys right of it become another. If the node is internal, its children are divided the same way. If the root is full, create a new empty root, make the old root its only child, then split. This is the only way the tree gets taller.

Example (t = 2, so max 3 keys per node):

Insert 10, 20, 30: [10 20 30]

Insert 40: root is full → split, 20 moves up
[20]
/ \
[10] [30 40]

Insert 50: [20]
/ \
[10] [30 40 50]

Insert 60: child [30 40 50] is full → split before entering, 40 moves up
[20 40]
/ | \
[10] [30] [50 60]
Part 3: Validate and InOrder (easy, but do it before Part 4)

Validate walks the tree and returns an error describing the first invariant violation it finds, or nil. InOrder returns all keys in sorted order.

Test them together: insert a few thousand random keys, then check that Validate() returns nil and that InOrder() equals the sorted input. Run this after every operation while debugging.

Part 4: Delete (hard)

Remove key if present, keeping all invariants. Use the top-down approach again: before descending into a child, make sure it has at least t keys (one more than the minimum), so removing a key from it can’t cause it to drop below the minimum.

Hints, roughly following CLRS’s cases:

Key is in a leaf: remove it.
Key is in an internal node: replace it with its predecessor (largest key in the left subtree) or successor (smallest key in the right subtree), whichever side has at least t keys, then delete that key from the subtree. If both children have only t - 1 keys, merge them along with the key, then recurse into the merged node.
Key isn’t in this node, and the child you’re about to enter has only t - 1 keys: borrow a key through the parent from a sibling that has at least t keys (a rotation), or if neither sibling can spare one, merge the child with a sibling and the separating parent key.
If the root ends up with 0 keys and one child, that child becomes the new root. This is the only way the tree gets shorter.

Constraints: 2 <= t <= 100, up to 10^5 operations. Search, Insert, and Delete should run in O(t · log_t n).

Follow-ups (the system design part)
Count disk reads. Add a counter that goes up every time you visit a node, and pretend each visit is one disk page read. Insert 1,000,000 keys with t = 2, then with t = 100, and compare the number of reads per search. This is why databases use B-trees: a node sized to fill a 4–16 KB page gives a very high t, so even a billion keys fit in a tree only 3–4 levels deep.
Range query. Implement Range(lo, hi int) []int. Then think about why it’s awkward here: you have to keep going up and down the tree.
Convert to a B+ tree. Store all keys in the leaves, keep only copies of keys in internal nodes as signposts, and link the leaves together in a list. Now Range is one search plus a straight walk along the leaves. This is the structure behind indexes in Postgres, MySQL (InnoDB), and SQLite.

If you want to check your work as you go, the visualizer at cs.usfca.edu/~galles/visualization/BTree.html lets you set the max degree and step through inserts and deletes one at a time.
