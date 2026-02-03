package ui

import (
	"log"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageNav struct {
	thm            *material.Theme
	backClickable  *widget.Clickable
	forthClickable *widget.Clickable
}

func NewPageNav(thm *material.Theme) *PageNav {
	return &PageNav{thm: thm, backClickable: new(widget.Clickable), forthClickable: new(widget.Clickable)}
}

func (pn PageNav) Layout(gtx C) D {
	backIcon, err := widget.NewIcon(icons.HardwareKeyboardArrowLeft)
	if err != nil {
		log.Fatal("Couldn't create new left icon")
	}
	backButton := material.IconButton(pn.thm, pn.backClickable, backIcon, "")
	backButton.Size = unit.Dp(25)
	backButton.Inset = layout.UniformInset(unit.Dp(5))

	forthIcon, err := widget.NewIcon(icons.HardwareKeyboardArrowRight)
	if err != nil {
		log.Fatal("Couldn't create new right icon")
	}
	forthButton := material.IconButton(pn.thm, pn.forthClickable, forthIcon, "")
	forthButton.Size = unit.Dp(25)
	forthButton.Inset = layout.UniformInset(unit.Dp(5))

	return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx, Rigid(backButton),
			Rigid(layout.Spacer{Width: unit.Dp(5)}), Rigid(forthButton))
	})
}

func (pn PageNav) BackClicked(gtx C) bool {
	return pn.backClickable.Clicked(gtx)
}

func (pn PageNav) ForthClicked(gtx C) bool {
	return pn.forthClickable.Clicked(gtx)
}
