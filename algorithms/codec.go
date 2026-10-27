package algorithms

import (
	"strconv"
	"strings"
)

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type Codec struct {
}

func NewCode() Codec { return Codec{} }

const NULL = "N"

// Serializes a tree to a single string.
func (this *Codec) Serialize(root *TreeNode) string {
	if root == nil {
		return ""
	}
	var strs []string
	q := []*TreeNode{root}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]
		if curr == nil {
			strs = append(strs, NULL)
			continue
		}
		strs = append(strs, strconv.Itoa(curr.Val))
		q = append(q, curr.Left, curr.Right)
	}
	return strings.Join(strs, ",")
}

// Deserializes your encoded data to tree.
func (this *Codec) Deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}

	tokens := strings.Split(data, ",")
	newNode := func(tok string) *TreeNode {
		if tok == NULL {
			return nil
		}
		val, _ := strconv.Atoi(tok)
		return &TreeNode{Val: val}
	}

	root := newNode(tokens[0])
	q := []*TreeNode{root}

	i := 1
	// 1 2 3
	// q: 0
	for len(q) > 0 && i < len(tokens) {
		curr := q[0]
		q = q[1:]

		curr.Left = newNode(tokens[i])
		i++
		if curr.Left != nil {
			q = append(q, curr.Left)
		}

		if i < len(tokens) {
			curr.Right = newNode(tokens[i])
			i++
			if curr.Right != nil {
				q = append(q, curr.Right)
			}
		}
	}

	return root
}
