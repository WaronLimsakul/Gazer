package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

// renderDOM takes a DOM root node and return a slice of FlexChild
// NOTE: Cuz I want to plug this in Flex.layout()
func renderDOM(gtx *layout.Context, thm *material.Theme, root *parser.Node) []layout.FlexChild {
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
		res = append(res, renderNode(gtx, thm, child)...)
	}

	return res
}

// renderNodes return flex children needs for render the node and its children
// NOTE: do we need the gtx here?
func renderNode(gtx *layout.Context, thm *material.Theme, node *parser.Node) []layout.FlexChild {
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
			res = append(res, renderNode(gtx, thm, child)...)
		}
	case parser.Title:
		return res // TODO
	case parser.H1:
		for _, child := range node.Children {
			if child.Tag == parser.Text {
				res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H1(thm, child.Inner).Layout(gtx)
				}))
			} else {
				res = append(res, renderNode(gtx, thm, child)...)
			}
		}
	case parser.H2:
		for _, child := range node.Children {
			if child.Tag == parser.Text {
				res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H2(thm, child.Inner).Layout(gtx)
				}))
			} else {
				res = append(res, renderNode(gtx, thm, child)...)
			}
		}
	case parser.H3:
		for _, child := range node.Children {
			if child.Tag == parser.Text {
				res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H3(thm, child.Inner).Layout(gtx)
				}))
			} else {
				res = append(res, renderNode(gtx, thm, child)...)
			}
		}
	case parser.H4:
		for _, child := range node.Children {
			if child.Tag == parser.Text {
				res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H4(thm, child.Inner).Layout(gtx)
				}))
			} else {
				res = append(res, renderNode(gtx, thm, child)...)
			}
		}
	case parser.H5:
		for _, child := range node.Children {
			if child.Tag == parser.Text {
				res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H5(thm, child.Inner).Layout(gtx)
				}))
			} else {
				res = append(res, renderNode(gtx, thm, child)...)
			}
		}
	case parser.P:
		for _, child := range node.Children {
			if child.Tag == parser.Text {
				res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Body1(thm, child.Inner).Layout(gtx)
				}))
			} else {
				res = append(res, renderNode(gtx, thm, child)...)
			}
		}
	case parser.Br:
		res = append(res, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
		}))
	}

	return res
}
