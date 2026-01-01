package ui

import (
	"log"

	"gioui.org/io/pointer"
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

type SearchBar struct {
	thm       *material.Theme
	editor    *widget.Editor
	clickable *widget.Clickable
}

func NewSearchBar(thm *material.Theme) *SearchBar {
	editor := setupSearchBarEditor()
	clickable := new(widget.Clickable)
	return &SearchBar{thm, editor, clickable}
}

func (s *SearchBar) Layout(gtx C) D {
	srcInputUi := material.Editor(s.thm, s.editor, "search")
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
		Color:        s.thm.Fg,
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
	searchButton := material.IconButton(s.thm, s.clickable, icon, "Search")
	searchButton.Size = unit.Dp(20)

	// search bar
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
}

// Searched return a bool whether user click or press "enter" to search
func (s *SearchBar) Searched(gtx C) bool {
	for {
		editorEv, ok := s.editor.Update(gtx)
		if !ok {
			break
		}
		switch editorEv.(type) {
		// press "enter" search
		case widget.SubmitEvent:
			return true
		default:
			continue
		}
	}

	// click search
	return s.clickable.Clicked(gtx)
}

// Text gets the text inside the search bar
func (s SearchBar) Text() string {
	return s.editor.Text()
}

// Update updates the ui when search bar is hovered
func (s SearchBar) Update(gtx C) {
	if s.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}
}

// SetupSearchEditor create a new widget.Editor used as
// input behavior for search component
func setupSearchBarEditor() *widget.Editor {
	srcInput := new(widget.Editor)
	srcInput.Alignment = text.Start
	srcInput.SingleLine = true
	srcInput.Submit = true
	return srcInput
}
