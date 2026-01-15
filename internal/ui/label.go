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
	margin  layout.Inset // margin outside border (if exists)
	padding layout.Inset // margin inside border (if exists)
	border  widget.Border
	bgColor color.NRGBA
	// color   color.NRGBA // text color

	// for <a> or <button>
	clickable *widget.Clickable

	// for <li>: e.g. Prefix "•"
	prefix string

	style material.LabelStyle
}

// what we need when create/decorate a Label
type LabelStyle struct {
	Base  css.Style
	Extra LabelExtraStyle
}

// Extra fields we need apart from css.Style
type LabelExtraStyle struct {
	Clickable *widget.Clickable
	Prefix    string
	Count     *int // for <ol>
}

func (l Label) Layout(gtx C) D {
	// handle ui interaction
	if l.clickable != nil {
		l.clickable.Update(gtx)
		if l.clickable.Hovered() {
			pointer.CursorNone.Add(gtx.Ops)
		}
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
						return l.style.Layout(gtx)
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
// TODO: make the impl looks better
func NewLabel(thm *Theme, lstyle LabelStyle, selectable *widget.Selectable, txt string) Label {
	var text material.LabelStyle
	if lstyle.Base.FontSize != nil {
		text = material.Label(thm, *lstyle.Base.FontSize, txt)
	} else {
		text = material.Label(thm, thm.TextSize, txt)
	}

	text.State = selectable
	if lstyle.Base.FontStyle != nil {
		text.Font.Style = *lstyle.Base.FontStyle
	}
	if lstyle.Base.FontWeight != nil {
		text.Font.Weight = *lstyle.Base.FontWeight
	}
	if lstyle.Base.Color != nil {
		text.Color = *lstyle.Base.Color
	}

	res := Label{
		prefix:    lstyle.Extra.Prefix,
		clickable: lstyle.Extra.Clickable,
		style:     text,
	}

	if lstyle.Base.BgColor != nil {
		res.bgColor = *lstyle.Base.BgColor
	}
	if lstyle.Base.Border != nil {
		res.border = *lstyle.Base.Border
	}
	if lstyle.Base.Margin != nil {
		res.margin = *lstyle.Base.Margin
	}
	if lstyle.Base.Padding != nil {
		res.padding = *lstyle.Base.Padding
	}

	return res
}

// TODO: FIXME: Bold and italic not rendered
func H1(thm *Theme, style LabelStyle) LabelStyle {
	size := thm.TextSize * 2.25
	style.Base.FontSize = &size
	bold := font.Bold
	style.Base.FontWeight = &bold
	return style
}

func H2(thm *Theme, style LabelStyle) LabelStyle {
	size := thm.TextSize * 1.75
	style.Base.FontSize = &size
	bold := font.Bold
	style.Base.FontWeight = &bold
	return style
}

func H3(thm *Theme, style LabelStyle) LabelStyle {
	size := thm.TextSize * 1.375
	style.Base.FontSize = &size
	bold := font.Bold
	style.Base.FontWeight = &bold
	return style
}

func H4(thm *Theme, style LabelStyle) LabelStyle {
	size := thm.TextSize * 1.125
	style.Base.FontSize = &size
	bold := font.Bold
	style.Base.FontWeight = &bold
	return style
}

func H5(thm *Theme, style LabelStyle) LabelStyle {
	size := thm.TextSize
	style.Base.FontSize = &size
	bold := font.Bold
	style.Base.FontWeight = &bold
	return style
}

func P(thm *Theme, style LabelStyle) LabelStyle {
	return style
}

func I(thm *Theme, style LabelStyle) LabelStyle {
	italic := font.Italic
	style.Base.FontStyle = &italic
	return style
}

func B(thm *Theme, style LabelStyle) LabelStyle {
	bold := font.Bold
	style.Base.FontWeight = &bold
	return style
}

func A(Clickable *widget.Clickable, style LabelStyle) LabelStyle {
	style.Base.Color = &color.NRGBA{R: 0, G: 0, B: 238, A: 255}
	style.Extra.Clickable = Clickable
	return style
}

func Ul(style LabelStyle) LabelStyle {
	if style.Base.Margin == nil {
		style.Base.Margin = new(layout.Inset)
	}
	style.Base.Margin.Left += unit.Dp(10)
	return style
}

func Ol(style LabelStyle) LabelStyle {
	if style.Base.Margin == nil {
		style.Base.Margin = new(layout.Inset)
	}
	style.Base.Margin.Left += unit.Dp(10)
	if style.Extra.Count == nil {
		style.Extra.Count = new(int)
	}
	*style.Extra.Count = 1 // reset the counting to 1
	return style
}

// we don't need thm, but just try to make it like the others
// TODO NOW: FIXME: Find other way to check if to add prefix, can't use this map anymore
func Li(thm *Theme, style LabelStyle, ancestors []parser.Tag) LabelStyle {
	if style.Extra.Prefix == "" {
		for i := len(ancestors) - 1; i >= 0; i-- {
			anc := ancestors[i]
			if anc == parser.Ol {
				style.Extra.Prefix = strconv.Itoa(*style.Extra.Count) + ". "
				*style.Extra.Count++
				break
			}

			if anc == parser.Ul {
				style.Extra.Prefix = "• "
				break
			}
		}
	}

	return style
}

func Button(thm *Theme, Clickable *widget.Clickable, style LabelStyle) LabelStyle {
	// TODO: not sure if it should be all or none like this
	if style.Base.Border == nil {
		style.Base.Border = &widget.Border{Color: thm.Fg, CornerRadius: unit.Dp(2), Width: unit.Dp(1)}
	}

	// TODO: use the full theme set
	if style.Base.BgColor == nil {
		lightGray := color.NRGBA{R: 240, G: 240, B: 240, A: 255}
		style.Base.BgColor = &lightGray
	}

	// TODO: v8 just let the margin be 0, but I feel like it's a little weird
	if style.Base.Margin == nil {
		buttonMargin := layout.UniformInset(unit.Dp(1))
		style.Base.Margin = &buttonMargin
	}

	// TODO: v8 has separate each padding side so they can have all optional, I have to do all or none for now
	if style.Base.Padding == nil {
		buttonPadding := layout.Inset{
			Top:    unit.Dp(3),
			Bottom: unit.Dp(3),
			Left:   unit.Dp(6),
			Right:  unit.Dp(6),
		}
		style.Base.Padding = &buttonPadding
	}

	style.Extra.Clickable = Clickable
	return style
}
