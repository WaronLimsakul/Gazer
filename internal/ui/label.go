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
	"github.com/WaronLimsakul/Gazer/internal/css"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

// type material.LabelStyle = material.material.LabelStyle
type Theme = material.Theme
type Style = css.Style

// Labels are supposed to be built using a decorator pattern.
// Start with empty LabelStyle, pass it around with all LabelStyleDecator and then
// finish the building with NewLabel()

// NewLabel() is a final function to build a label. Client should have built the LabelStyle
// as they like before calling this.
type LabelFunc = func(*Theme, *LabelStyle, *widget.Selectable, string) Label

// Decorate any material.LabelStyle with LabelStyleDecorator e.g. H1, H2, etc.
// NOTE: some decorator might have a little different signature
type LabelStyleDecorator = func(*Theme, LabelStyle) LabelStyle

type Label struct {
	tags map[parser.Tag]bool

	margin  layout.Inset // margin outside border (if exists)
	padding layout.Inset // margin inside border (if exists)
	border  widget.Border
	bgColor color.NRGBA
	color   color.NRGBA // text color

	// for <a> or <button>
	clickable *widget.Clickable

	// for <li>: e.g. prefix "•"
	prefix string

	style material.LabelStyle
}

// what we need when create/decorate a Label
type LabelStyle struct {
	base  css.Style
	extra LabelExtraStyle
}

// extra fields we need apart from css.Style
type LabelExtraStyle struct {
	tags      map[parser.Tag]bool
	clickable *widget.Clickable
	prefix    string
	count     *int // for <ol>
}

// TODO NOW: after finish acc recursion, do this
func (l Label) Layout(gtx C) D {
	// handle ui interaction
	l.clickable.Update(gtx)
	if l.clickable.Hovered() {
		pointer.CursorNone.Add(gtx.Ops)
	}

	// layout
	return l.margin.Layout(gtx, func(gtx C) D {
		nonPrefixLabel := func(gtx C) D {
			return l.border.Layout(gtx, func(gtx C) D {
				var contentSize D
				var contentOp op.CallOp
				contentWidget := func(gtx C) D {
					return l.padding.Layout(gtx, func(gtx C) D {
						// material.LabelStyle.Layout try to takes just what it need by default.
						// However, passed gtx might just give min = max = max
						gtx.Constraints.Min = image.Point{}
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
				rrect := clip.UniformRRect(
					image.Rectangle{Max: contentSize.Size}, gtx.Dp(l.border.CornerRadius))
				// NOTE: can do this or use layout.Background{}
				defer rrect.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, l.bgColor)
				contentOp.Add(gtx.Ops)
				return D{Size: contentSize.Size}
			})
		}

		if len(l.prefix) == 0 {
			return nonPrefixLabel(gtx)
		} else {
			prefixStyle := l.style
			prefixStyle.Text = l.prefix
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				Rigid(prefixStyle),
				layout.Rigid(nonPrefixLabel),
			)
		}
	})
}

// NewLabel create a new Label element from the theme, labeltyle, selectable and text string
func NewLabel(thm *Theme, lstyle *LabelStyle, selectable *widget.Selectable, txt string) Label {
	var text material.LabelStyle
	if lstyle.base.FontSize != nil {
		text = material.Label(thm, *lstyle.base.FontSize, txt)
	} else {
		text = material.Label(thm, thm.TextSize, txt)
	}

	text.State = selectable
	res := Label{tags: lstyle.extra.tags, prefix: lstyle.extra.prefix, style: text}
	if lstyle.base.Color != nil {
		res.color = *lstyle.base.Color
	}
	if lstyle.base.BgColor != nil {
		res.bgColor = *lstyle.base.BgColor
	}
	if lstyle.base.Border != nil {
		res.border = *lstyle.base.Border
	}
	if lstyle.base.Margin != nil {
		res.margin = *lstyle.base.Margin
	}
	if lstyle.base.Padding != nil {
		res.padding = *lstyle.base.Padding
	}

	return res
}

func H1(thm *Theme, style LabelStyle) LabelStyle {
	style.extra.tags[parser.H1] = true
	size := thm.TextSize * 2.25
	style.base.FontSize = &size
	bold := font.Bold
	style.base.FontWeight = &bold
	return style
}

func H2(thm *Theme, style LabelStyle) LabelStyle {
	style.extra.tags[parser.H2] = true
	size := thm.TextSize * 1.75
	style.base.FontSize = &size
	bold := font.Bold
	style.base.FontWeight = &bold
	return style
}

func H3(thm *Theme, style LabelStyle) LabelStyle {
	style.extra.tags[parser.H3] = true
	size := thm.TextSize * 1.375
	style.base.FontSize = &size
	bold := font.Bold
	style.base.FontWeight = &bold
	return style
}

func H4(thm *Theme, style LabelStyle) LabelStyle {
	style.extra.tags[parser.H4] = true
	size := thm.TextSize * 1.125
	style.base.FontSize = &size
	bold := font.Bold
	style.base.FontWeight = &bold
	return style
}

func H5(thm *Theme, style LabelStyle) LabelStyle {
	style.extra.tags[parser.H5] = true
	size := thm.TextSize
	style.base.FontSize = &size
	bold := font.Bold
	style.base.FontWeight = &bold
	return style
}

func P(thm *Theme, style LabelStyle) LabelStyle {
	style.extra.tags[parser.P] = true
	return style
}

func I(thm *Theme, style LabelStyle) LabelStyle {
	style.extra.tags[parser.I] = true
	italic := font.Italic
	style.base.FontStyle = &italic
	return style
}

func B(thm *Theme, style LabelStyle) LabelStyle {
	style.extra.tags[parser.B] = true
	bold := font.Bold
	style.base.FontWeight = &bold
	return style
}

func A(clickable *widget.Clickable, style LabelStyle) LabelStyle {
	style.extra.tags[parser.A] = true
	style.base.Color = &color.NRGBA{R: 0, G: 0, B: 238, A: 255}
	style.extra.clickable = clickable
	return style
}

// we don't need thm, but just try to make it like the others
func Ul(style LabelStyle) LabelStyle {
	style.extra.tags[parser.Ul] = true
	if style.base.Margin == nil {
		style.base.Margin = new(layout.Inset)
	}
	style.base.Margin.Left += unit.Dp(10)
	return style
}

func Ol(style LabelStyle, count *int) LabelStyle {
	style.extra.tags[parser.Ol] = true
	if style.base.Margin == nil {
		style.base.Margin = new(layout.Inset)
	}
	style.base.Margin.Left += unit.Dp(10)
	*style.extra.count = 1 // reset the counting to 1
	return style
}

// we don't need thm, but just try to make it like the others
func Li(thm *Theme, style LabelStyle) LabelStyle {
	// if we are a child of ul or ol, add prefix
	// TODO: find a way to only add prefix of the most inner one
	if style.extra.tags[parser.Ul] && style.extra.prefix == "" {
		style.extra.prefix = "• "
	}
	if style.extra.tags[parser.Ol] && style.extra.prefix == "" {
		style.extra.prefix = strconv.Itoa(*style.extra.count) + ". "
		*style.extra.count++
	}

	style.extra.tags[parser.Li] = true
	return style
}

func Button(thm *Theme, clickable *widget.Clickable, style LabelStyle) LabelStyle {
	// TODO: not sure if it should be all or none like this
	if style.base.Border == nil {
		style.base.Border = &widget.Border{Color: thm.Fg, CornerRadius: unit.Dp(2), Width: unit.Dp(1)}
	}

	// TODO: use the full theme set
	if style.base.BgColor == nil {
		lightGray := color.NRGBA{R: 240, G: 240, B: 240, A: 255}
		style.base.BgColor = &lightGray
	}

	// TODO: v8 just let the margin be 0, but I feel like it's a little weird
	if style.base.Margin == nil {
		buttonMargin := layout.UniformInset(unit.Dp(1))
		style.base.Margin = &buttonMargin
	}

	// TODO: v8 has separate each padding side so they can have all optional, I have to do all or none for now
	if style.base.Padding == nil {
		buttonPadding := layout.Inset{
			Top:    unit.Dp(3),
			Bottom: unit.Dp(3),
			Left:   unit.Dp(6),
			Right:  unit.Dp(6),
		}
		style.base.Padding = &buttonPadding
	}

	style.extra.clickable = clickable
	style.extra.tags[parser.Button] = true
	return style
}
