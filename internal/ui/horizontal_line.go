package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func HorizontalLine(thm *material.Theme, width unit.Dp) layout.FlexChild {
	border := widget.Border{Color: thm.Fg, Width: unit.Dp(0.5)}
	margin := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5)}
	return layout.Rigid(func(gtx C) D {
		return margin.Layout(gtx, func(gtx C) D {
			return border.Layout(gtx, layout.Spacer{Width: width, Height: unit.Dp(0.1)}.Layout)
		})
	})
}
