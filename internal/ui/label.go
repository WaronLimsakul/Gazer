package ui

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

type LabelStyle = material.LabelStyle
type Theme = material.Theme

// label components are supposed to be used as decorator pattern
// Text() is a base component with LabelFunc type
type LabelFunc = func(*Theme, *widget.Selectable, string) LabelStyle

// Decorate any LabelStyle with LabelDecorator e.g. H1, H2, etc.
type LabelDecorator = func(*Theme, LabelStyle) LabelStyle

type Label struct {
	tags      map[parser.Tag]bool
	margin    layout.Inset
	clickable *widget.Clickable
	style     LabelStyle
}

func (l Label) Layout(gtx C) D {
	return l.margin.Layout(gtx, func(gtx C) D {
		// LabelStyle.Layout try to takes just what it need by default.
		// However, passed gtx might just give min = max = max
		gtx.Constraints.Min = image.Point{}

		if l.clickable != nil {
			return l.clickable.Layout(gtx, l.style.Layout)
		}
		return l.style.Layout(gtx)
	})
}

func Text(thm *Theme, selectable *widget.Selectable, txt string) Label {
	text := material.Label(thm, thm.TextSize, txt)
	text.State = selectable
	tags := make(map[parser.Tag]bool)
	tags[parser.Text] = true
	return Label{tags: tags, style: text}
}

func H1(thm *Theme, label Label) Label {
	label.tags[parser.H1] = true
	label.style.TextSize = thm.TextSize * 96.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func H2(thm *Theme, label Label) Label {
	label.tags[parser.H2] = true
	label.style.TextSize = thm.TextSize * 60.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func H3(thm *Theme, label Label) Label {
	label.tags[parser.H3] = true
	label.style.TextSize = thm.TextSize * 48.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func H4(thm *Theme, label Label) Label {
	label.tags[parser.H4] = true
	label.style.TextSize = thm.TextSize * 34.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func H5(thm *Theme, label Label) Label {
	label.tags[parser.H5] = true
	label.style.TextSize = thm.TextSize * 24.0 / 16.0
	label.style.Font.Weight = font.Bold
	return label
}

func P(thm *Theme, label Label) Label {
	label.tags[parser.P] = true
	label.style.TextSize = thm.TextSize
	return label
}

func I(thm *Theme, label Label) Label {
	label.tags[parser.I] = true
	label.style.Font.Style = font.Italic
	return label
}

func B(thm *Theme, label Label) Label {
	label.tags[parser.B] = true
	label.style.Font.Weight = font.Bold
	return label
}

func A(clickable *widget.Clickable, label Label) Label {
	label.tags[parser.A] = true
	label.style.Color = color.NRGBA{R: 0, G: 0, B: 238, A: 255}
	label.clickable = clickable
	return label
}

// we don't need thm, but just try to make it like the others
func Ul(thm *Theme, label Label) Label {
	// give Li label a bullet point of not yet
	if label.tags[parser.Li] && !label.tags[parser.Ul] {
		label.style.Text = "• " + label.style.Text
	}

	label.margin.Left += unit.Dp(10)
	label.tags[parser.Ul] = true
	return label
}

// we don't need thm, but just try to make it like the others
func Li(thm *Theme, label Label) Label {
	label.tags[parser.Li] = true
	return label
}
