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

func H1(thm *Theme, label LabelStyle) LabelStyle {
	label.TextSize = thm.TextSize * 96.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func H2(thm *Theme, label LabelStyle) LabelStyle {
	label.TextSize = thm.TextSize * 60.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func H3(thm *Theme, label LabelStyle) LabelStyle {
	label.TextSize = thm.TextSize * 48.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func H4(thm *Theme, label LabelStyle) LabelStyle {
	label.TextSize = thm.TextSize * 34.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func H5(thm *Theme, label LabelStyle) LabelStyle {
	label.TextSize = thm.TextSize * 24.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func P(thm *Theme, label LabelStyle) LabelStyle {
	label.TextSize = thm.TextSize
	return label
}

func I(thm *Theme, label LabelStyle) LabelStyle {
	label.Font.Style = font.Italic
	return label
}

func B(thm *Theme, label LabelStyle) LabelStyle {
	label.Font.Weight = font.Bold
	return label
}

func A(thm *Theme, label LabelStyle) LabelStyle {
	label.Color = color.NRGBA{R: 0, G: 0, B: 238, A: 255}
	return label
}

func Text(thm *Theme, selectable *widget.Selectable, txt string) LabelStyle {
	text := material.Label(thm, thm.TextSize, txt)
	text.State = selectable
	return text
}
