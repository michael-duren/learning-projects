package tries

type node struct {
	children [26]*node
	word     bool
}

type WordDictionary struct {
	root node
}

func Constructor() WordDictionary {
	return WordDictionary{
		root: node{
			children: [26]*node{},
		},
	}
}

func (this *WordDictionary) AddWord(word string) {
	curr := &this.root
	for _, c := range word {
		idx := c - 'a'
		if curr.children[idx] == nil {
			nn := node{
				children: [26]*node{},
			}
			curr.children[idx] = &nn
			curr = curr.children[idx]
		} else {
			curr = curr.children[idx]
		}
	}
	curr.word = true
}

// bob
// b..
func (this *WordDictionary) Search(word string) bool {
	q := make([]node, 0, 26)
	q[0] = this.root

	for _, c := range word {
		n := len(q)
		for range n {
			curr := q[0]
			q = q[1:]
			if c == '.' {
				for _, child := range curr.children {
					if child == nil {
						continue
					}
					q = append(q, *child)
				}
				if len(q) == 0 {
					return false
				}
				continue
			}
			idx := c - 'a'
			if curr.children[idx] != nil {
				q = append(q, *curr.children[idx])
			}
		}
		if len(q) < 1 {
			return false
		}
	}
	for _, child := range q {
		if child.word {
			return true
		}
	}
	return false
}
