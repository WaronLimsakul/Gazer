package ui

import (
	"gioui.org/font"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// label components are supposed to be used as decorator pattern
// Text() is a base component with LabelFunc type
type LabelFunc = func(*material.Theme, *widget.Selectable, string) material.LabelStyle

// Decorate any LabelStyle with LabelDecorator e.g. H1, H2, etc.
type LabelDecorator = func(*material.Theme, material.LabelStyle) material.LabelStyle

func H1(thm *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.TextSize = thm.TextSize * 96.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func H2(thm *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.TextSize = thm.TextSize * 60.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func H3(thm *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.TextSize = thm.TextSize * 48.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func H4(thm *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.TextSize = thm.TextSize * 34.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func H5(thm *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.TextSize = thm.TextSize * 24.0 / 16.0
	label.Font.Weight = font.Bold
	return label
}

func P(thm *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.TextSize = thm.TextSize
	return label
}

func I(thm *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.Font.Style = font.Italic
	return label
}

func B(thm *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.Font.Weight = font.Bold
	return label
}

func Text(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	text := material.Label(thm, thm.TextSize, txt)
	text.State = selectable
	return text
}
