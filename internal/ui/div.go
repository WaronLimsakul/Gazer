package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	// "gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/css"
)

// Children of all container-type component
type ContainerChildren [][]Element

// A basic container component, used for rendering <div>
type Div struct {
	thm *material.Theme

	margin  layout.Inset
	padding layout.Inset
	border  widget.Border
	bgColor color.NRGBA

	children ContainerChildren
}

// NewDiv creates new Div from a theme, css style and children it supposed to have
// NOTE: users must gather all ther children elements before create a new Div.
func NewDiv(thm *material.Theme, style css.Style, children [][]Element) Div {
	res := Div{children: children, thm: thm}
	if style.Margin != nil {
		res.margin = *style.Margin
	}
	if style.Padding != nil {
		res.padding = *style.Padding
	}
	if style.Border != nil {
		res.border = *style.Border
	}
	if style.BgColor != nil {
		res.bgColor = *style.BgColor
	}
	return res
}

func (d Div) Layout(gtx C) D {
	bg := func(gtx C) D {
		// foreground dimension is passed through gtx.Constraints.Min
		rrect := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(d.border.CornerRadius))
		defer rrect.Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, d.bgColor)

		return D{Size: gtx.Constraints.Min}
	}

	return d.margin.Layout(gtx, func(gtx C) D {
		return d.border.Layout(gtx, func(gtx C) D {
			return layout.Background{}.Layout(gtx, bg, func(gtx C) D {
				return d.padding.Layout(gtx, func(gtx C) D {
					return d.children.Layout(gtx)
				})
			})
		})
	})
}

// LayoutChildren lays the children of the Div out without style using flex for both axis
func (c ContainerChildren) Layout(gtx C) D {
	main := layout.Flex{Axis: layout.Vertical}
	mainChildren := make([]layout.FlexChild, len(c))
	rows := make([]layout.Flex, len(c))
	for i, row := range c {
		rowChildren := make([]layout.FlexChild, len(row))
		for j, el := range row {
			rowChildren[j] = Rigid(el)
		}
		mainChildren[i] = layout.Rigid(func(gtx C) D {
			return rows[i].Layout(gtx, rowChildren...)
		})
	}
	return main.Layout(gtx, mainChildren...)
}
