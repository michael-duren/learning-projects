package algorithms

type PrefixTree struct {
	children [26]*PrefixTree
	word     bool
}

func Constructor() PrefixTree {
	return PrefixTree{
		children: [26]*PrefixTree{},
	}
}

func (this *PrefixTree) Insert(word string) {
	curr := this
	for _, v := range word {
		idx := v - 'a'
		if curr.children[idx] == nil {
			curr.children[idx] = new(Constructor())
		}
		curr = curr.children[idx]
	}
	curr.word = true
}

func (this *PrefixTree) Search(word string) bool {
	curr := this
	for _, v := range word {
		idx := v - 'a'
		if curr.children[idx] == nil {
			return false
		}
		curr = curr.children[idx]
	}
	return curr.word
}

func (this *PrefixTree) StartsWith(prefix string) bool {
	curr := this
	for _, v := range prefix {
		idx := v - 'a'
		if curr.children[idx] == nil {
			return false
		}
		curr = curr.children[idx]
	}
	return true
}

type TrieNode struct {
	children [26]*TrieNode
	word     bool
	orig     string
}

func NewTrieNode() TrieNode {
	return TrieNode{
		children: [26]*TrieNode{},
	}
}

func (t *TrieNode) Insert(word string) {
	curr := t
	for _, v := range word {
		idx := v - 'a'
		if curr.children[idx] == nil {
			curr.children[idx].children = [26]*TrieNode{}
		}
		curr = curr.children[idx]
	}
	curr.word = true
	curr.orig = word
}

func findWords(board [][]byte, words []string) []string {
	t := NewTrieNode()
	for _, w := range words {
		t.Insert(w)
	}

	rv := make([]string, 0, len(words))
	var dfs func(row, col int, t *TrieNode)
	dfs = func(row, col int, t *TrieNode) {
		if row > len(board) || row < 0 || col > len(board[0]) || col < 0 {
			return
		}

		curr := t.children[board[row][col]-'a']
		// no words down this path
		if curr == nil {
			return
		}

		if curr.word {
			rv = append(rv, curr.orig)
		}

		directions := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
		for _, d := range directions {
			dfs(row+d[0], col+d[1], curr)
		}
	}

	for i := range len(board) {
		for j := range len(board[0]) {
			dfs(i, j, &t)
		}
	}
	return rv
}
