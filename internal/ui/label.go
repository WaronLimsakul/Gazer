package ui

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type LabelStyle = material.LabelStyle
type Theme = material.Theme

// label components are supposed to be used as decorator pattern
// Text() is a base component with LabelFunc type
type LabelFunc = func(*Theme, *widget.Selectable, string) LabelStyle

// Decorate any LabelStyle with LabelDecorator e.g. H1, H2, etc.
type LabelDecorator = func(*Theme, LabelStyle) LabelStyle

type Label struct {
	clickable *widget.Clickable
	style     LabelStyle
}

func (l Label) Layout(gtx C) D {
	if l.clickable != nil {
		return l.clickable.Layout(gtx, l.style.Layout)
	}
	return l.style.Layout(gtx)
}

func Text(thm *Theme, selectable *widget.Selectable, txt string) Label {
	text := material.Label(thm, thm.TextSize, txt)
	text.State = selectable
	return Label{clickable: nil, style: text}
}

func H1(thm *Theme, label Label) Label {
	label.style.TextSize = thm.TextSize * 96.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func H2(thm *Theme, label Label) Label {
	label.style.TextSize = thm.TextSize * 60.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func H3(thm *Theme, label Label) Label {
	label.style.TextSize = thm.TextSize * 48.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func H4(thm *Theme, label Label) Label {
	label.style.TextSize = thm.TextSize * 34.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func H5(thm *Theme, label Label) Label {
	label.style.TextSize = thm.TextSize * 24.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func P(thm *Theme, label Label) Label {
	label.style.TextSize = thm.TextSize
	return label
}

func I(thm *Theme, label Label) Label {
	label.style.Font.Style = font.Italic
	return label
}

func B(thm *Theme, label Label) Label {
	label.style.Font.Weight = font.Bold
	return label
}

func A(clickable *widget.Clickable, label Label) Label {
	label.style.Color = color.NRGBA{R: 0, G: 0, B: 238, A: 255}
	label.clickable = clickable
	return label
}
