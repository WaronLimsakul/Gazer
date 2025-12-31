package ui

import (
	"log"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// SetupSearchEditor create a new widget.Editor used as
// input behavior for search component
func SetupSearchEditor() *widget.Editor {
	srcInput := new(widget.Editor)
	srcInput.Alignment = text.Start
	srcInput.SingleLine = true
	srcInput.Submit = true
	return srcInput
}

func SearchBar(
	thm *material.Theme,
	editor *widget.Editor,
	searchClickable *widget.Clickable,
) layout.FlexChild {
	srcInputUi := material.Editor(thm, editor, "search")
	srcInputUi.TextSize = unit.Sp(20)

	// search bar spacing
	margin := layout.Inset{
		Top:    unit.Dp(25),
		Bottom: unit.Dp(25),
		// press to get the search bar width
		Left:  unit.Dp(400),
		Right: unit.Dp(400),
	}
	border := widget.Border{
		Color:        thm.Fg,
		CornerRadius: unit.Dp(2),
		Width:        unit.Dp(1),
	}
	insideBorderMargin := layout.Inset{
		Top:    unit.Dp(8),
		Bottom: unit.Dp(8),
		Left:   unit.Dp(10),
		Right:  unit.Dp(10),
	}

	// search button
	icon, err := widget.NewIcon(icons.ActionSearch)
	if err != nil {
		log.Fatal("Couldn't create new search icon")
	}
	searchButton := material.IconButton(thm, searchClickable, icon, "Search")
	searchButton.Size = unit.Dp(20)

	// search bar
	return layout.Rigid(func(gtx C) D {
		return margin.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Spacing: 2}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return border.Layout(gtx, func(gtx C) D {
						return insideBorderMargin.Layout(gtx, srcInputUi.Layout)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Spacer{Width: unit.Dp(5)}.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return searchButton.Layout(gtx)
				}),
			)
		})
	})
}
