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
func (dr *DomRenderer) render(root *parser.Node) []layout.FlexChild {
	res := make([]layout.FlexChild, 0)

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
func (dr *DomRenderer) renderNode(node *parser.Node) []layout.FlexChild {
	res := make([]layout.FlexChild, 0)
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
		res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}))
	}

	if parser.TextElements[node.Tag] {
		res = append(res, labelsToFlexChildren(dr.renderText(node))...)
	}

	return res
}

// renderText returns flex children needs for rendering node
// with the direct child node that has Text tag being rendered as textFuc desire.
// TODO NOW: have to make the rendering inherited
// e.g. <i><h1>hello</h1></i>, <h1><i>hello</i></h1> or  has to be big and italic
func (dr *DomRenderer) renderText(node *parser.Node) []material.LabelStyle {
	// base case
	if node.Tag == parser.Text {
		selectable, ok := dr.selectables[node]
		if !ok {
			selectable = new(widget.Selectable)
			dr.selectables[node] = selectable
		}
		return []material.LabelStyle{ui.Text(dr.thm, selectable, node.Inner)}
	}

	res := make([]material.LabelStyle, 0)
	for _, child := range node.Children {
		res = append(res, dr.renderText(child)...)
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

	for i := range res {
		res[i] = dec(dr.thm, res[i])
	}
	return res
}

func labelsToFlexChildren(labels []material.LabelStyle) []layout.FlexChild {
	res := make([]layout.FlexChild, len(labels))
	for i, label := range labels {
		res[i] = layout.Rigid(func(gtx C) D {
			return label.Layout(gtx)
		})
	}
	return res
}
