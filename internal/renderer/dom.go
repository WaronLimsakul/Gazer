package renderer

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/css"
	"github.com/WaronLimsakul/Gazer/internal/parser"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

// Main renderering of a website. One of these per tab.
type DomRenderer struct {
	thm *material.Theme
	tab *ui.Tab
	// Cache the matrix of elements with root node pointer.
	// Can cache it because engine also cache by pointer
	// (same url + same tab = same root ptr).
	cache map[*Node]*[][]Element
	// All Texts' selectables elements based on its pointer.
	// These pointers will not be cleaned because the map still refer to it.
	selectables      map[*Node]*widget.Selectable
	linkClickables   map[*Node]*widget.Clickable
	buttonClickables map[*Node]*widget.Clickable
	inputEditors     map[*Node]*widget.Editor
}

type Node = parser.Node
type StyleSet = css.StyleSet
type Style = css.Style

func newDomRenderer(thm *material.Theme, tab *ui.Tab) *DomRenderer {
	return &DomRenderer{thm: thm, tab: tab, cache: make(map[*Node]*[][]Element),
		selectables:      make(map[*Node]*widget.Selectable),
		linkClickables:   make(map[*Node]*widget.Clickable),
		buttonClickables: make(map[*Node]*widget.Clickable),
		inputEditors:     make(map[*Node]*widget.Editor),
	}
}

// renderDOM takes a DOM root node and return [][]Element
// First layer (outer) is each horizontal line of rendering.
// Second layer (inner) is each element in that line from left to right.
func (dr *DomRenderer) render(root *Node, styles *StyleSet) [][]Element {
	res := make([][]Element, 0)
	// expect to be Root node
	if root == nil || root.Tag != parser.Root {
		return res
	}

	if cachedRes, ok := dr.cache[root]; ok {
		return *cachedRes
	}

	// expect root node to only have HTML tag
	// len(root.Children) != 1 ||
	if root.Children[0].Tag != parser.Html {
		return res
	}

	htmlNode := root.Children[0]
	for _, child := range htmlNode.Children {
		res = append(res, dr.renderNode(child, styles)...)
	}

	dr.cache[root] = &res
	return res
}

// renderNodes returns flex children needs for render a node and its children.
func (dr *DomRenderer) renderNode(node *Node, styles *StyleSet) [][]Element {
	res := make([][]Element, 0)
	switch node.Tag {
	case parser.Body:
		res = dr.gatherElements(node)
	case parser.Div:
		res = dr.gatherElements(node)
	case parser.Span:
		res = dr.gatherElements(node)
	case parser.Section:
		res = dr.gatherElements(node)
	case parser.Br:
		res = append(res, []Element{layout.Spacer{Height: unit.Dp(10)}})
	case parser.Hr:
		res = append(res, []Element{ui.HorizontalLine{Thm: dr.thm, Width: WINDOW_WIDTH, Height: unit.Dp(1)}})
	case parser.Img:
		img, err := dr.renderImg(node)
		if err != nil {
			break
		}
		res = append(res, []Element{img})
	case parser.Input:
		res = append(res, []Element{dr.renderInput(node)})
	}

	if parser.TextElements[node.Tag] {
		res = append(res, dr.renderText(node, styles)...)
	}

	return res
}

// renderText returns [][]Element needs for rendering a text node and its children.
// requires: node must be of the text type (check by using parser.TextElements)
func (dr *DomRenderer) renderText(node *Node, styles *StyleSet) [][]Element {
	// TODO NOW: solve conflict between inline styles and global style we have: use acc rec

	// base case
	if node.Tag == parser.Text {
		selectable, ok := dr.selectables[node]
		if !ok {
			selectable = new(widget.Selectable)
			dr.selectables[node] = selectable
		}
		// Text type node is not a real html tag, so it will never be affected by style (base case)
		return [][]Element{{ui.Text(dr.thm, selectable, node.Inner)}}
	}

	// recursive case [phase 1]: aggregate the children elements
	res := dr.gatherElements(node, styles)

	// recursive case [phase 2]: decorate the children that are Label
	// recursive case [phase 2.1]: take care of a special decorator
	var inlineStyle *Style
	styleStr, ok := node.Attrs["style"]
	if ok {
		res := css.ParseStyle(styleStr)
		inlineStyle = &res
	}

	// recursive case [phase 2.2]: take care of a special decorator
	// tag "A" decorator has different signature
	if node.Tag == parser.A {
		clickable, ok := dr.linkClickables[node]
		if !ok {
			clickable = new(widget.Clickable)
			dr.linkClickables[node] = clickable
		}

		for _, line := range res {
			for i, el := range line {
				if label, ok := el.(ui.Label); ok {
					line[i] = ui.A(clickable, label)
				}
			}
		}
		return res
	}

	if node.Tag == parser.Button {
		// TODO: v8 just wrap all text around and treat it like one big button
		clickable, ok := dr.buttonClickables[node]
		if !ok {
			clickable = new(widget.Clickable)
			dr.buttonClickables[node] = clickable
		}

		for _, line := range res {
			for i, el := range line {
				if label, ok := el.(ui.Label); ok {
					line[i] = ui.Button(dr.thm, clickable, label)
				}
			}
		}
		return res
	}

	// tag "ul" and "li" don't wanna apply to everyone in a row

	if node.Tag == parser.Ul {
		for _, line := range res {
			if len(line) > 0 {
				if label, ok := line[0].(ui.Label); ok {
					line[0] = ui.Ul(label)
				}
			}
		}
		return res
	}

	if node.Tag == parser.Ol {
		count := 1
		for _, line := range res {
			if len(line) > 0 {
				if label, ok := line[0].(ui.Label); ok {
					line[0] = ui.Ol(label, &count)
				}
			}
		}
		return res
	}

	// recursive case [phase 2.3]: normal text decorator
	dec := ui.P
	switch node.Tag {
	case parser.H1:
		dec = ui.H1
	case parser.H2:
		dec = ui.H2
	case parser.H3:
		dec = ui.H3
	case parser.H4:
		dec = ui.H4
	case parser.H5:
		dec = ui.H5
	case parser.P:
		dec = ui.P
	case parser.I:
		dec = ui.I
	case parser.B:
		dec = ui.B
	case parser.Li:
		dec = ui.Li
	}

	for _, line := range res {
		for i, el := range line {
			if label, ok := el.(ui.Label); ok {
				line[i] = dec(dr.thm, label)
			}
		}
	}

	return res
}

// renderImg receive Img tag node and return Img ui element.
// Img is void element, don't have to gather more
func (dr *DomRenderer) renderImg(node *Node) (Element, error) {
	if node == nil {
		return layout.Spacer{}, fmt.Errorf("nil node")
	}
	if node.Tag != parser.Img {
		return layout.Spacer{}, fmt.Errorf("invalid tag: %v", node.Tag.String())
	}
	img, err := ui.NewImg(node.Attrs["src"])
	if err != nil {
		return layout.Spacer{}, fmt.Errorf("ui.NewImg: %v", err)
	}
	return img, nil
}

// renderInput receive Input tag node and return Input ui element.
// Input is void element, don't have to gather more.
// requires: node must not be nil and have input tag
func (dr *DomRenderer) renderInput(node *Node) Element {
	editor, ok := dr.inputEditors[node]
	if !ok {
		editor = new(widget.Editor)
		dr.inputEditors[node] = editor
	}

	inputTypeStr := strings.ToLower(node.Attrs["type"])
	inputType, ok := ui.InputTypes[inputTypeStr]
	if !ok {
		inputType = ui.TextInput
	}

	hint := node.Attrs["placeholder"]
	return ui.NewInput(dr.thm, inputType, editor, hint)
}

// handleHead set the tabview data by processing <head> node in the DOM tree (except css-related)
func (dr *DomRenderer) handleHead(root *Node) {
	head := dr.findHead(root)
	if head == nil {
		dr.tab.Title = ""
		return
	}

	var titleSet bool
	for _, node := range head.Children {
		// handle title tag
		switch node.Tag {
		case parser.Title:
			for _, titleChild := range node.Children {
				if titleChild.Tag == parser.Text {
					titleSet = true
					dr.tab.Title = titleChild.Inner
				}
			}
			// TODO: support some other links
		}
	}

	if !titleSet {
		dr.tab.Title = ""
	}
}

func (dr DomRenderer) findHead(node *Node) *Node {
	if node == nil {
		return nil
	}
	if node.Tag == parser.Head {
		return node
	}
	for _, child := range node.Children {
		found := dr.findHead(child)
		if found != nil {
			return found
		}
	}
	return nil
}

// gaterElements recieves a node and gather all elements of the node's children
// according the tag rule (inline, block)
func (dr DomRenderer) gatherElements(node *Node, styles *StyleSet) [][]Element {
	res := make([][]Element, 0)
	// if inline-text and prev is also inline-text, put it in latest one don't append
	for i, child := range node.Children {
		if parser.InlineElements[child.Tag] && len(res) > 0 && i > 0 &&
			parser.InlineElements[node.Children[i-1].Tag] {
			childElems := dr.renderNode(child, styles)
			if len(childElems) > 0 {
				res[len(res)-1] = append(res[len(res)-1], childElems[0]...)
			}
			if len(childElems) > 1 {
				res = append(res, childElems[1:]...)
			}

		} else {
			res = append(res, dr.renderNode(child, styles)...)
		}

	}
	return res
}

// linkClicked return whether the link in the page is clicked and
// if so, what does it linked to.
func (dr *DomRenderer) linkClicked(gtx C) (bool, string) {
	for node, clickable := range dr.linkClickables {
		if clickable.Clicked(gtx) {
			return true, node.Attrs["href"]
		}
	}
	return false, ""
}
