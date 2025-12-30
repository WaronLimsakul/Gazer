package renderer

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

// renderDOM takes a DOM root node and return a slice of FlexChild
// NOTE: Cuz I want to plug this in Flex.layout()
func renderDOM(thm *material.Theme, root *parser.Node) []layout.FlexChild {
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
		res = append(res, renderNode(thm, child)...)
	}

	return res
}

// renderNodes returns flex children needs for render the node and its children.
func renderNode(thm *material.Theme, node *parser.Node) []layout.FlexChild {
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
		res = append(res, renderText(thm, material.Body1, node)...)
	case parser.Title:
		return res // TODO
	case parser.Meta:
		return res // TODO
	case parser.Div:
		// Text child from div should be rendered as body1
		res = append(res, renderText(thm, material.Body1, node)...)
	case parser.H1:
		res = append(res, renderText(thm, material.H1, node)...)
	case parser.H2:
		res = append(res, renderText(thm, material.H2, node)...)
	case parser.H3:
		res = append(res, renderText(thm, material.H3, node)...)
	case parser.H4:
		res = append(res, renderText(thm, material.H4, node)...)
	case parser.H5:
		res = append(res, renderText(thm, material.H5, node)...)
	case parser.P:
		res = append(res, renderText(thm, material.Body1, node)...)
	case parser.Br:
		res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}))
	}

	return res
}

type Label = func(*material.Theme, string) material.LabelStyle

// renderText returns flex children needs for rendering node
// with the direct child node that has Text tag being rendered as textFuc desire.
func renderText(thm *material.Theme, textFunc Label, node *parser.Node) []layout.FlexChild {
	res := make([]layout.FlexChild, 0)
	for _, child := range node.Children {
		if child.Tag == parser.Text {
			res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return textFunc(thm, child.Inner).Layout(gtx)
			}))
		} else {
			res = append(res, renderNode(thm, child)...)
		}
	}
	return res
}
