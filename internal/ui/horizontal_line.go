package ui

import (
	"image"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type HorizontalLine struct {
	Thm    *material.Theme
	Width  unit.Dp
	Height unit.Dp
}

func (h HorizontalLine) Layout(gtx C) D {
	height := min(gtx.Dp(h.Height), gtx.Constraints.Max.Y)
	width := min(gtx.Dp(h.Width), gtx.Constraints.Max.X)

	line := clip.Rect{Max: image.Pt(width, height)}

	defer line.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: h.Thm.ContrastBg}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return D{Size: image.Point{width, height}}
}
