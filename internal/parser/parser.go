package parser

import (
	"fmt"
	"maps"
	"strings"
	"unicode"

	"github.com/WaronLimsakul/Gazer/internal/lexer"
)

type Tag int

const (
	Root Tag = iota // Only for root node
	Html
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

	Text // For no tag content or invalid tag

	// TODO: A Img Ul Ol Li B (or Strong) I (or Em) Hr Div Span
)

var TagMap = map[string]Tag{
	"html":  Html,
	"head":  Head,
	"body":  Body,
	"title": Title,
	"h1":    H1,
	"h2":    H2,
	"h3":    H3,
	"h4":    H4,
	"h5":    H5,
	"p":     P,
	"br":    Br,
}

type Node struct {
	Tag      Tag
	Inner    string // only for Text node content
	Attrs    map[string]string
	Children []*Node
	Parent   *Node
}

// Parse parses raw html string and return root node of the DOM
// NOTE: if tag invalid, assume it's Text node
// NOTE2: special tag <br>, <br/> or even </br> always means self-close <br/>
func Parse(src string) (*Node, error) {
	root := newBaseNode()
	curNode := root

	var token lexer.Token
	// process token-by-token to create a DOM tree
	for idx := 0; idx < len(src); idx = token.Endpos {
		token = lexer.GetNextToken(src, idx)

		// look at NOTE2
		if getTagFromContent(token.Content) == Br {
			token.Type = lexer.SClose
		}

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
		// no tag content = being a child Text node
		case lexer.NoTag:
			curNode.Children = append(curNode.Children, newText(token.Content, curNode))
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
// 1. get tag name.
// 2. assign attributes
//   - key=value works
//   - key="value and another value" also works
func newNode(content string) (*Node, error) {
	node := newBaseNode()
	content = strings.TrimSpace(content)

	node.Tag = getTagFromContent(content)

	// Assign attributes
	tagSepIdx := strings.Index(content, " ")
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
				key = ""
			} else {
				val += string(char)
			}
		case Valuing:
			if unicode.IsSpace(char) {
				(*attrs)[key] = val
				state = Keying
				key = ""
			} else {
				val += string(char)
			}
		}
	}
}

// getTagFromContent takes string content of the tag and appropriate Tag
// e.g. getTagFromContent("p style=color:white") will return P
func getTagFromContent(content string) Tag {
	var tagName string
	tagSepIdx := strings.Index(content, " ")
	if tagSepIdx != -1 {
		tagName = content[:tagSepIdx]
	} else {
		tagName = content
	}
	return getTag(tagName)
}

// getTag return Tag based on the name, if invalid, gives <p>
func getTag(tagName string) Tag {
	tag, ok := TagMap[tagName]
	if !ok {
		return Text
	} else {
		return tag
	}
}

func newBaseNode() *Node {
	node := new(Node)
	node.Attrs = make(map[string]string)
	node.Children = make([]*Node, 0)
	return node
}

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

func (n Node) String() string {
	return n.StringRec(0)
}

// recursively print node while informed layer
func (n Node) StringRec(layer int) string {
	tags := []string{"root", "html", "head", "body", "title", "h1", "h2", "h3", "h4", "h5", "p", "br"}
	res := "\n"
	if layer > 0 {
		for range layer {
			res += "\t"
		}
	}

	res += fmt.Sprintf("{%s | inner: %s | attrs: %v | parent: %p}", tags[n.Tag], n.Inner, n.Attrs, n.Parent)
	if len(n.Children) == 0 {
		return res
	}

	for _, child := range n.Children {
		res += child.StringRec(layer + 1)
	}
	return res
}

// newText create a new Text Node with target content and its parent
func newText(content string, parent *Node) *Node {
	text := newBaseNode()
	text.Tag = Text
	text.Inner = content
	text.Parent = parent
	return text
}
