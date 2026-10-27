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
