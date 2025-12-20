package parser

import (
	"fmt"

	"github.com/WaronLimsakul/Gazer/internal/lexer"
)

type Tag int

const (
	Html Tag = iota
	Head
	Title
	H1
	H2
	H3
	H4
	H5
	P
	Br

	// TODO: A Img Ul Ol Li B (or Strong) I (or Em) Hr Div Span
)

var TagMap = map[string]Tag{
	"html": Html,
	"head": Head,
	"h1":   H1,
	"h2":   H2,
	"h3":   H3,
	"h4":   H4,
	"h5":   H5,
	"p":    P,
	"br":   Br,
}

type Node struct {
	Tag      string
	Inner    string
	Attrs    map[string]any
	Children []*Node
	Parent   *Node
}

// Parse parses raw html string and return root node of the DOM
// NOTE: if tag invalid, assume it's "p" tag
func Parse(src string) (*Node, error) {
	root := newRootNode()
	curNode := root

	var token lexer.Token
	// process token-by-token to create a DOM tree
	for idx := 0; idx < len(src); idx = token.Endpos {
		token = lexer.GetNextToken(src, idx)

		switch token.Type {
		// open-tag = create a child node
		case lexer.Open:
			child, err := newNode(token.Content)
			if err != nil {
				return nil, err
			}
			child.Parent = curNode
			curNode.Children = append(curNode.Children, child)
			curNode = child
		// close-tag = done with a child node, going back to parent
		case lexer.Close:
			if curNode.Parent != nil {
				curNode = curNode.Parent
			}
		case lexer.Inner:
			curNode.Inner += token.Content
		// same as having open tag, but not going to that child
		case lexer.SClose:
			child, err := newNode(token.Content)
			if err != nil {
				return nil, err
			}
			child.Parent = curNode
			curNode.Children = append(curNode.Children, child)
		// specailly just to check HTML
		case lexer.DocType:
			if token.Content != "html" {
				return nil, fmt.Errorf("Not html")
			}
		// if invalid token (Void), just ignore
		default:
			continue
		}

	}
	return root, nil
}

// newNode takes content of the open tag and return a new node
func newNode(content string) (*Node, error) {
	node := new(Node)
	return node, nil
}

func newRootNode() *Node {
	node := new(Node)
	node.Attrs = make(map[string]any)
	node.Children = make([]*Node, 0)
	return node
}
