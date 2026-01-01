package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type HorizontalLine struct {
	Thm   *material.Theme
	Width unit.Dp
}

func (h HorizontalLine) Layout(gtx C) D {
	border := widget.Border{Color: h.Thm.Fg, Width: unit.Dp(0.5)}
	margin := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5)}
	return margin.Layout(gtx, func(gtx C) D {
		return border.Layout(gtx, layout.Spacer{Width: h.Width, Height: unit.Dp(0.1)}.Layout)
	})
}
