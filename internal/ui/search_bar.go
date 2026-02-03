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

// match pair clickable with editor
var searchButtonClickables = map[*widget.Editor]*widget.Clickable{}

func NewSearchBar(thm *material.Theme, editor *widget.Editor) *SearchBar {
	clickable, ok := searchButtonClickables[editor]
	if !ok {
		clickable = new(widget.Clickable)
		searchButtonClickables[editor] = clickable
	}
	return &SearchBar{thm: thm, editor: editor, clickable: clickable}
}

// TODO: make it longer than this
func (s *SearchBar) Layout(gtx C) D {
	// handle ui interaction
	s.clickable.Update(gtx)
	if s.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	// editor material
	srcInputUi := material.Editor(s.thm, s.editor, "search")
	srcInputUi.TextSize = unit.Sp(20)

	// search bar spacing
	margin := layout.Inset{
		Top:    unit.Dp(5),
		Bottom: unit.Dp(5),
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
			Rigid(layout.Spacer{Width: unit.Dp(5)}),
			Rigid(searchButton),
		)
	})
}

// Searched updates the editor and clickable in search bar and
// returns a bool whether user click or press "enter" to search
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

// SetText sets the text inside search bar
func (s SearchBar) SetText(txt string) {
	s.editor.SetText(txt)
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
