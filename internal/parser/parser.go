package parser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/WaronLimsakul/Gazer/internal/lexer"
)

type Tag int

const (
	Html Tag = iota
	Head
	Body
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
	"body": Body,
	"h1":   H1,
	"h2":   H2,
	"h3":   H3,
	"h4":   H4,
	"h5":   H5,
	"p":    P,
	"br":   Br,
}

type Node struct {
	Tag      Tag
	Inner    string
	Attrs    map[string]string
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
// 1. get tag name
// 2. assign attributes
//   - key=value works
//   - key="value and another value" also works
func newNode(content string) (*Node, error) {
	node := new(Node)
	content = strings.TrimSpace(content)

	// Get tag
	var tagName string
	tagSepIdx := strings.Index(content, " ")
	if tagSepIdx != -1 {
		tagName = content[:tagSepIdx]
	} else {
		tagName = content
	}
	node.Tag = getTag(tagName)

	// Assign attributes
	if tagSepIdx == -1 || tagSepIdx+1 >= len(content) {
		return node, nil
	}
	content = content[tagSepIdx+1:]
	assignAttrs(&node.Attrs, content)
	return node, nil
}

// enum for attribute processing states
type attrParsingState int

const (
	Keying    attrParsingState = iota // processing key part
	Observing                         // observe wheter it will be key=value or key="value" format
	QValuing                          // processing value in key="value" format
	Valuing                           // processing value in key=value format
)

// assignAttrs takes a map and raw string in the attribute part of HTML tag
// then assign all of them to the map
func assignAttrs(attrs *map[string]string, s string) {
	s = strings.TrimSpace(s)
	var key, val string
	state := Keying
	for _, char := range s {
		switch state {
		case Keying:
			if unicode.IsSpace(char) {
				continue
			}

			if char == '=' {
				state = Observing
			} else {
				key += string(char)
			}
		case Observing:
			if unicode.IsSpace(char) {
				continue
			}

			if char == '"' {
				state = QValuing
			} else {
				state = Valuing
			}
			val = ""
		case QValuing:
			if char == '"' {
				(*attrs)[key] = val
				state = Keying
			} else {
				val += string(char)
			}
		case Valuing:
			if unicode.IsSpace(char) {
				(*attrs)[key] = val
				state = Keying
			} else {
				val += string(char)
			}
		}
	}
}

// getTag return Tag based on the name, if invalid, gives <p>
func getTag(tagName string) Tag {
	tag, ok := TagMap[tagName]
	if !ok {
		return P
	} else {
		return tag
	}
}

func newRootNode() *Node {
	node := new(Node)
	node.Attrs = make(map[string]string)
	node.Children = make([]*Node, 0)
	return node
}
