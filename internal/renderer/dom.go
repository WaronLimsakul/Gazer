package renderer

import (
	"gioui.org/font"
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
		// Text child from body should be rendered as body1
		res = append(res, dr.renderText(ui.Body1, node)...)
	case parser.Title:
		return res // TODO
	case parser.Meta:
		return res // TODO
	case parser.Div:
		// Text child from div should be rendered as body1
		res = append(res, dr.renderText(ui.Body1, node)...)
	case parser.H1:
		res = append(res, dr.renderText(ui.H1, node)...)
	case parser.H2:
		res = append(res, dr.renderText(ui.H2, node)...)
	case parser.H3:
		res = append(res, dr.renderText(ui.H3, node)...)
	case parser.H4:
		res = append(res, dr.renderText(ui.H4, node)...)
	case parser.H5:
		res = append(res, dr.renderText(ui.H5, node)...)
	case parser.P:
		res = append(res, dr.renderText(ui.Body1, node)...)
	case parser.I:
		res = append(res, dr.renderText(
			func(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
				label := material.Body1(dr.thm, txt)
				label.Font.Style = font.Italic
				label.State = selectable
				return label
			}, node)...)
	case parser.Br:
		res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}))
	}

	return res
}

type Label = func(*material.Theme, *widget.Selectable, string) material.LabelStyle

// renderText returns flex children needs for rendering node
// with the direct child node that has Text tag being rendered as textFuc desire.
func (dr *DomRenderer) renderText(textFunc Label, node *parser.Node) []layout.FlexChild {
	res := make([]layout.FlexChild, 0)
	for _, child := range node.Children {
		if child.Tag == parser.Text {
			// get text's selectable before layout
			selectable, ok := dr.selectables[node]
			if !ok {
				selectable = new(widget.Selectable)
				dr.selectables[node] = selectable
			}
			res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return textFunc(dr.thm, selectable, child.Inner).Layout(gtx)
			}))
		} else {
			res = append(res, dr.renderNode(child)...)
		}
	}
	return res
}
