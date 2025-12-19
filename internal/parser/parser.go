package parser

import (
	"fmt"
	"strings"
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
}

// Parse parses raw html string and return root node of the DOM
func Parse(src string) (*Node, error) {
	src, isHtml := strings.CutPrefix(src, "<!DOCTYPE html>")
	if !isHtml {
		return nil, fmt.Errorf("parser: not HTML5")
	}

	// TODO
	return nil, nil
}

func newNode() *Node {
	node := Node{}
	node.Attrs = make(map[string]any)
	node.Children = make([]*Node, 0)
	return &node
}
