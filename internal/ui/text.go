package ui

import (
	"gioui.org/font"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Text = func(*material.Theme, *widget.Selectable, string) material.LabelStyle

func H1(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h1 := material.H1(thm, txt)
	h1.State = selectable
	return h1
}

func H2(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h2 := material.H2(thm, txt)
	h2.State = selectable
	return h2
}

func H3(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h3 := material.H3(thm, txt)
	h3.State = selectable
	return h3
}

func H4(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h4 := material.H4(thm, txt)
	h4.State = selectable
	return h4
}

func H5(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h5 := material.H5(thm, txt)
	h5.State = selectable
	return h5
}

func P(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	p := material.Body1(thm, txt)
	p.State = selectable
	return p
}

func I(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	i := material.Body1(thm, txt)
	i.State = selectable
	i.Font.Style = font.Italic
	return i
}
