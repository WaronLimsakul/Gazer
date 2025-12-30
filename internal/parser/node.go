package parser

import (
	"fmt"
	"maps"
)

type Node struct {
	Tag      Tag
	Inner    string // only for Text node content
	Attrs    map[string]string
	Children []*Node
	Parent   *Node
}

func (n Node) String() string {
	return n.recursiveString(0)
}

// recursively print node while informed layer
func (n Node) recursiveString(layer int) string {
	res := "\n"
	if layer > 0 {
		for range layer {
			res += "\t"
		}
	}

	res += fmt.Sprintf("{%s | inner: %s | attrs: %v | parent: %p}", n.Tag, n.Inner, n.Attrs, n.Parent)
	if len(n.Children) == 0 {
		return res
	}

	for _, child := range n.Children {
		res += child.recursiveString(layer + 1)
	}
	return res
}

// isChildOfTag traces the tree and and return boolean whether the
// node is a child of a node with target tag.
// func (n Node) isChildOfTag(tag Tag) bool {
// 	if n.Parent == nil {
// 		return false
// 	} else if n.Parent.Tag == tag {
// 		return true
// 	}
// 	return n.Parent.isChildOfTag(tag)
// }

func (n Node) equal(other *Node) bool {
	// check simple fields
	if n.Tag != other.Tag ||
		n.Inner != other.Inner ||
		len(n.Children) != len(other.Children) ||
		!maps.Equal(n.Attrs, other.Attrs) {
		return false
	}

	// recursively check children
	if len(n.Children) == 0 && len(other.Children) == 0 {
		return true
	} else {
		for i, child := range n.Children {
			if !child.equal(other.Children[i]) {
				return false
			}
		}
	}
	return true
}

// newNode returns a new basic, ready-to-use node
func newNode() *Node {
	node := new(Node)
	node.Attrs = make(map[string]string)
	node.Children = make([]*Node, 0)
	return node
}

// newTextNode create a new Text Node with target content and its parent
func newTextNode(content string, parent *Node) *Node {
	text := newNode()
	text.Tag = Text
	text.Inner = content
	text.Parent = parent
	return text
}
