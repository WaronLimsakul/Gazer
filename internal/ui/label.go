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

// object passed around in dom tree (top down)
// TODO NOW NOW: wait, but in the domrender.renderNode(), it's just css.Style, not ui.LabelStyle
type LabelStyle struct {
	Style
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

// NewLabel receives create a new Label element
func NewLabel(thm *Theme, lstyle *LabelStyle, selectable *widget.Selectable, txt string) Label {
	var text material.LabelStyle
	if lstyle.FontSize != nil {
		text = material.Label(thm, *lstyle.FontSize, txt)
	} else {
		text = material.Label(thm, thm.TextSize, txt)
	}

	text.State = selectable
	res := Label{tags: lstyle.tags, prefix: lstyle.prefix, style: text}
	if lstyle.Color != nil {
		res.color = *lstyle.Color
	}
	if lstyle.BgColor != nil {
		res.bgColor = *lstyle.BgColor
	}
	if lstyle.Border != nil {
		res.border = *lstyle.Border
	}
	if lstyle.Margin != nil {
		res.margin = *lstyle.Margin
	}
	if lstyle.Padding != nil {
		res.padding = *lstyle.Padding
	}

	return res
}

func H1(thm *Theme, style LabelStyle) LabelStyle {
	style.tags[parser.H1] = true
	size := thm.TextSize * 2.25
	style.FontSize = &size
	bold := font.Bold
	style.FontWeight = &bold
	return style
}

func H2(thm *Theme, style LabelStyle) LabelStyle {
	style.tags[parser.H2] = true
	size := thm.TextSize * 1.75
	style.FontSize = &size
	bold := font.Bold
	style.FontWeight = &bold
	return style
}

func H3(thm *Theme, style LabelStyle) LabelStyle {
	style.tags[parser.H3] = true
	size := thm.TextSize * 1.375
	style.FontSize = &size
	bold := font.Bold
	style.FontWeight = &bold
	return style
}

func H4(thm *Theme, style LabelStyle) LabelStyle {
	style.tags[parser.H4] = true
	size := thm.TextSize * 1.125
	style.FontSize = &size
	bold := font.Bold
	style.FontWeight = &bold
	return style
}

func H5(thm *Theme, style LabelStyle) LabelStyle {
	style.tags[parser.H5] = true
	size := thm.TextSize
	style.FontSize = &size
	bold := font.Bold
	style.FontWeight = &bold
	return style
}

func P(thm *Theme, style LabelStyle) LabelStyle {
	style.tags[parser.P] = true
	return style
}

func I(thm *Theme, style LabelStyle) LabelStyle {
	style.tags[parser.I] = true
	italic := font.Italic
	style.FontStyle = &italic
	return style
}

func B(thm *Theme, style LabelStyle) LabelStyle {
	style.tags[parser.B] = true
	bold := font.Bold
	style.FontWeight = &bold
	return style
}

func A(clickable *widget.Clickable, style LabelStyle) LabelStyle {
	style.tags[parser.A] = true
	style.Color = &color.NRGBA{R: 0, G: 0, B: 238, A: 255}
	style.clickable = clickable
	return style
}

// we don't need thm, but just try to make it like the others
func Ul(style LabelStyle) LabelStyle {
	style.tags[parser.Ul] = true
	if style.Margin == nil {
		style.Margin = new(layout.Inset)
	}
	style.Margin.Left += unit.Dp(10)
	return style
}

func Ol(style LabelStyle, count *int) LabelStyle {
	style.tags[parser.Ol] = true
	if style.Margin == nil {
		style.Margin = new(layout.Inset)
	}
	style.Margin.Left += unit.Dp(10)
	*style.count = 1 // reset the counting to 1
	return style
}

// we don't need thm, but just try to make it like the others
func Li(thm *Theme, style LabelStyle) LabelStyle {
	// if we are a child of ul or ol, add prefix
	// TODO: find a way to only add prefix of the most inner one
	if style.tags[parser.Ul] && style.prefix == "" {
		style.prefix = "• "
	}
	if style.tags[parser.Ol] && style.prefix == "" {
		style.prefix = strconv.Itoa(*style.count) + ". "
		*style.count++
	}

	style.tags[parser.Li] = true
	return style
}

func Button(thm *Theme, clickable *widget.Clickable, style LabelStyle) LabelStyle {
	// TODO: not sure if it should be all or none like this
	if style.Border == nil {
		style.Border = &widget.Border{Color: thm.Fg, CornerRadius: unit.Dp(2), Width: unit.Dp(1)}
	}

	// TODO: use the full theme set
	if style.BgColor == nil {
		lightGray := color.NRGBA{R: 240, G: 240, B: 240, A: 255}
		style.BgColor = &lightGray
	}

	// TODO: v8 just let the margin be 0, but I feel like it's a little weird
	if style.Margin == nil {
		buttonMargin := layout.UniformInset(unit.Dp(1))
		style.Margin = &buttonMargin
	}

	// TODO: v8 has separate each padding side so they can have all optional, I have to do all or none for now
	if style.Padding == nil {
		buttonPadding := layout.Inset{
			Top:    unit.Dp(3),
			Bottom: unit.Dp(3),
			Left:   unit.Dp(6),
			Right:  unit.Dp(6),
		}
		style.Padding = &buttonPadding
	}

	style.clickable = clickable
	style.tags[parser.Button] = true
	return style
}
