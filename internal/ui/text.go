package ui

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func H1(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h1 := material.H1(thm, txt)
	h1.State = selectable
	return h1
}

func H2(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h1 := material.H2(thm, txt)
	h1.State = selectable
	return h1
}

func H3(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h1 := material.H3(thm, txt)
	h1.State = selectable
	return h1
}

func H4(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h1 := material.H4(thm, txt)
	h1.State = selectable
	return h1
}

func H5(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h1 := material.H5(thm, txt)
	h1.State = selectable
	return h1
}

func Body1(thm *material.Theme, selectable *widget.Selectable, txt string) material.LabelStyle {
	h1 := material.Body1(thm, txt)
	h1.State = selectable
	return h1
}
