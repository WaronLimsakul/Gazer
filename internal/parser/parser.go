package parser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/WaronLimsakul/Gazer/internal/lexer"
)

// Parse parses raw html string and return root node of the DOM
// NOTE: if tag invalid, assume it's Text node
// NOTE2: void elements (e.g. <br>) always means self-close
func Parse(src string) (*Node, error) {
	root := newNode()
	curNode := root

	var token lexer.Token
	// process token-by-token to create a DOM tree
	for idx := 0; idx < len(src); idx = token.Endpos {
		token = lexer.GetNextToken(src, idx)

		// handle void elements
		if VoidElements[getTagFromContent(token.Content)] {
			token.Type = lexer.SClose
		}

		switch token.Type {
		// open-tag = create a child node
		case lexer.Open:
			child, err := parseNode(token.Content)
			if err != nil {
				return nil, err
			}

			// if head is still not closed but there is body, auto-close head
			// NOTE: If any more rule like this, we have to centralize it.
			if curNode.Tag == Head && child.Tag == Body {
				curNode = curNode.Parent
			}

			child.Parent = curNode
			curNode.Children = append(curNode.Children, child)
			curNode = child
		// close-tag = done with a child node, going back to parent
		// NOTE: notice that the tag type doesn't matter. Close is close.
		case lexer.Close:
			if curNode.Parent != nil {
				curNode = curNode.Parent
			}
		// no tag content = being a child Text node
		case lexer.NoTag:
			curNode.Children = append(curNode.Children, newTextNode(token.Content, curNode))
		// same as having open tag, but not going to that child
		case lexer.SClose:
			child, err := parseNode(token.Content)
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
		// if invalid token (Void) or comment, just ignore
		default:
			continue
		}

	}
	return root, nil
}

// parseNode takes content of the open tag and return a new node
// 1. get tag name.
// 2. assign attributes
//   - key=value works
//   - key="value and another value" also works
func parseNode(content string) (*Node, error) {
	node := newNode()
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
type attrParsingState uint8

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
				// key is case-insensitive
				key += strings.ToLower(string(char))
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
	tag, ok := TagMap[strings.ToLower(tagName)]
	if !ok {
		return Div // assume unknown node to be a div
	} else {
		return tag
	}
}
