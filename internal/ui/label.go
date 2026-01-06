package ui

import (
	"image"
	"image/color"
	"strconv"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
	tags map[parser.Tag]bool
	// margin outside border (if exists)
	margin layout.Inset
	// margin inside border (if exists)
	textMargin layout.Inset
	border     widget.Border
	bgColor    color.NRGBA
	color      color.NRGBA // text color TODO: not sure if we needs to check zero value
	clickable  *widget.Clickable
	// for <li>: e.g. prefix "•"
	prefix string
	style  LabelStyle
}

func (l Label) Layout(gtx C) D {
	// handle ui interaction
	l.clickable.Update(gtx)
	if l.clickable.Hovered() {
		pointer.CursorNone.Add(gtx.Ops)
	}

	// layout
	return l.margin.Layout(gtx, func(gtx C) D {
		normalLabel := func(gtx C) D {
			return l.border.Layout(gtx, func(gtx C) D {
				var contentSize D
				var contentOp op.CallOp
				contentWidget := func(gtx C) D {
					return l.textMargin.Layout(gtx, func(gtx C) D {
						// LabelStyle.Layout try to takes just what it need by default.
						// However, passed gtx might just give min = max = max
						gtx.Constraints.Min = image.Point{}
						// TODO: not sure
						tmpStyle := l.style
						tmpStyle.Color = l.color
						return tmpStyle.Layout(gtx)
					})
				}
				macro := op.Record(gtx.Ops)
				if l.clickable != nil {
					contentSize = l.clickable.Layout(gtx, contentWidget)
				} else {
					contentSize = contentWidget(gtx)
				}
				contentOp = macro.Stop()
				rrect := clip.UniformRRect(image.Rectangle{Max: contentSize.Size}, gtx.Dp(l.border.CornerRadius))
				defer rrect.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, l.bgColor)
				contentOp.Add(gtx.Ops)
				return D{Size: contentSize.Size}
			})
		}
		if len(l.prefix) == 0 {
			return normalLabel(gtx)
		} else {
			prefixStyle := l.style
			prefixStyle.Text = l.prefix
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				Rigid(prefixStyle),
				layout.Rigid(normalLabel),
			)
		}
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
	label.style.TextSize = thm.TextSize * 2.25
	label.style.Font.Weight = font.Bold
	return label
}

func H2(thm *Theme, label Label) Label {
	label.tags[parser.H2] = true
	label.style.TextSize = thm.TextSize * 1.75
	label.style.Font.Weight = font.Bold
	return label
}

func H3(thm *Theme, label Label) Label {
	label.tags[parser.H3] = true
	label.style.TextSize = thm.TextSize * 1.375
	label.style.Font.Weight = font.Bold
	return label
}

func H4(thm *Theme, label Label) Label {
	label.tags[parser.H4] = true
	label.style.TextSize = thm.TextSize * 1.125
	label.style.Font.Weight = font.Bold
	return label
}

func H5(thm *Theme, label Label) Label {
	label.tags[parser.H5] = true
	label.style.TextSize = thm.TextSize
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
func Ul(label Label) Label {
	// give Li label a bullet point of not yet
	if label.tags[parser.Li] && !label.tags[parser.Ul] && !label.tags[parser.Ol] {
		label.prefix = "• "
	}

	label.margin.Left += unit.Dp(10)
	label.tags[parser.Ul] = true
	return label
}

func Ol(label Label, count *int) Label {
	if label.tags[parser.Li] && !label.tags[parser.Ol] && !label.tags[parser.Ul] {
		label.prefix = strconv.Itoa(*count) + ". "
		*count++
	}

	label.margin.Left += unit.Dp(10)
	label.tags[parser.Ol] = true
	return label
}

// we don't need thm, but just try to make it like the others
func Li(thm *Theme, label Label) Label {
	label.tags[parser.Li] = true
	return label
}

func Button(thm *Theme, clickable *widget.Clickable, label Label) Label {
	label.border = widget.Border{Color: thm.Fg, CornerRadius: unit.Dp(2), Width: unit.Dp(1)}

	// TODO: use the full theme set
	lightGray := color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	label.bgColor = lightGray

	label.margin.Left += unit.Dp(1)
	label.margin.Right += unit.Dp(1)
	label.margin.Top += unit.Dp(1)
	label.margin.Bottom += unit.Dp(1)

	label.textMargin.Top += unit.Dp(3)
	label.textMargin.Bottom += unit.Dp(3)
	label.textMargin.Left += unit.Dp(6)
	label.textMargin.Right += unit.Dp(6)

	label.clickable = clickable
	label.tags[parser.Button] = true
	return label
}
