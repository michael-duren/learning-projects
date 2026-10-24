package algorithms

type ListNode struct {
	key  int
	next *ListNode
}

type MyHashSet struct {
	len     int
	buckets []*ListNode
}

func Constructor() MyHashSet {
	return MyHashSet{
		len:     10000,
		buckets: make([]*ListNode, 0, 10000),
	}
}

func (this *MyHashSet) Add(key int) {
	bucket := this.len % key
	newnode := &ListNode{
		key:  key,
		next: nil,
	}
	if this.buckets[bucket] == nil {
		this.buckets[bucket] = newnode
	} else {
		curr := this.buckets[bucket]
		for curr.next != nil {
			curr = curr.next
		}
		curr.next = newnode
	}
}

func (this *MyHashSet) Remove(key int) {
	bucket := this.len % key
	curr := this.buckets[bucket]
	if curr == nil {
		return
	}
	var prev *ListNode
	for curr != nil {
		if curr.key == key {
			if prev == nil {
				this.buckets[bucket] = curr.next
			} else {
				prev.next = curr.next
			}
			return
		}
		prev = curr
		curr = curr.next
	}
}

func (this *MyHashSet) Contains(key int) bool {
	bucket := this.len % key
	curr := this.buckets[bucket]
	if curr == nil {
		return false
	}
	for curr != nil {
		if curr.key == key {
			return true
		}
		curr = curr.next
	}
	return false
}

/**
 * Your MyHashSet object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Add(key);
 * obj.Remove(key);
 * param_3 := obj.Contains(key);
 */
