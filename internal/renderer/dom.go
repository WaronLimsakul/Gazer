package renderer

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/parser"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

type DomRenderer struct {
	thm *material.Theme
	url string
	// All Texts' selectables elements based on its pointer.
	// NOTE: we can use pointer because it is one url per DomRenderer
	selectables map[*parser.Node]*widget.Selectable
}

func newDomRenderer(thm *material.Theme, url string) *DomRenderer {
	return &DomRenderer{thm, url, make(map[*parser.Node]*widget.Selectable)}
}

// renderDOM takes a DOM root node and return a slice of FlexChild
// NOTE: Cuz I want to plug this in Flex.layout()
func (dr *DomRenderer) render(root *parser.Node) [][]Element {
	res := make([][]Element, 0)

	// expect to be Root node
	if root == nil || root.Tag != parser.Root {
		return res
	}

	// expect root node to only have HTML tag
	if len(root.Children) != 1 || root.Children[0].Tag != parser.Html {
		return res
	}

	htmlNode := root.Children[0]
	for _, child := range htmlNode.Children {
		res = append(res, dr.renderNode(child)...)
	}

	return res
}

// renderNodes returns flex children needs for render the node and its children.
func (dr *DomRenderer) renderNode(node *parser.Node) [][]Element {
	res := make([][]Element, 0)
	switch node.Tag {
	// Ignore root or html tag
	case parser.Root:
		break
	case parser.Html:
		break
	case parser.Head:
		return res // TODO
	case parser.Body:
		for _, child := range node.Children {
			res = append(res, dr.renderNode(child)...)
		}
	case parser.Title:
		return res // TODO
	case parser.Meta:
		return res // TODO
	case parser.Div:
		for _, child := range node.Children {
			res = append(res, dr.renderNode(child)...)
		}
	case parser.Br:
		res = append(res, []Element{layout.Spacer{Height: unit.Dp(10)}})
	}

	if parser.TextElements[node.Tag] {
		res = append(res, labelsToElements(dr.renderText(node))...)
	}

	return res
}

// renderText returns [][]LabelStyle needs for rendering node and its children.
// First layer (outer) is each horizontal line of rendering.
// Second layer (inner) is each element in that line from left to right.
func (dr *DomRenderer) renderText(node *parser.Node) [][]material.LabelStyle {
	// base case
	if node.Tag == parser.Text {
		selectable, ok := dr.selectables[node]
		if !ok {
			selectable = new(widget.Selectable)
			dr.selectables[node] = selectable
		}

		return [][]material.LabelStyle{{ui.Text(dr.thm, selectable, node.Inner)}}
	}

	res := make([][]material.LabelStyle, 0)
	// TODO: if inline-text and prev is also inline-text, put it in latest one don't append
	for i, child := range node.Children {
		if parser.InlineTextElements[child.Tag] && len(res) > 0 && i > 0 &&
			parser.InlineTextElements[node.Children[i-1].Tag] {
			childElems := dr.renderText(child)
			if len(childElems) > 0 {
				res[len(res)-1] = append(res[len(res)-1], childElems[0]...)
			}
			if len(childElems) > 1 {
				res = append(res, childElems[1:]...)
			}

		} else {
			res = append(res, dr.renderText(child)...)
		}

	}

	// recursive case: decorate
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
	case parser.A:
		dec = ui.A
	}

	for _, line := range res {
		for i := range line {
			line[i] = dec(dr.thm, line[i])
		}
	}
	return res
}

// labelsToElements wrap each LabelStyle with Rigid to get [][]FlexChild
func labelsToElements(labels [][]material.LabelStyle) [][]Element {
	res := make([][]Element, len(labels))
	for i, line := range labels {
		res[i] = make([]Element, len(line))
		for j, label := range line {
			res[i][j] = label
		}
	}
	return res
}
