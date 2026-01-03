package ui

import (
	"image/color"
	"log"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Tabs struct {
	Tabs     []*Tab
	addTab   *widget.Clickable
	thm      *Theme
	Selected int // idx of selected tab
}

type Tab struct {
	clickable    *widget.Clickable
	SearchEditor *widget.Editor
}

func NewTabs(thm *Theme) *Tabs {
	firstTab := newTab()
	tabs := []*Tab{firstTab}
	addTab := new(widget.Clickable)
	res := &Tabs{Tabs: tabs, addTab: addTab, thm: thm}
	res.Select(0)
	return res
}

func (t Tabs) Layout(gtx C) D {
	// TODO: use new theme system
	tabsBarBg := color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	tabsMargin := layout.Inset{
		Left:   unit.Dp(5),
		Right:  unit.Dp(5),
		Top:    unit.Dp(5),
		Bottom: unit.Dp(5),
	}
	flex := layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}

	flexChildren := make([]layout.FlexChild, len(t.Tabs)+1)
	for i, tab := range t.Tabs {
		flexChildren[i] = layout.Rigid(func(gtx C) D {
			isSelected := i == t.Selected
			return tab.Layout(t.thm, gtx, isSelected)
		})
	}

	plusIcon, err := widget.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatalf("Couldn't get icon: %v", err)
	}
	newTabButton := material.IconButton(t.thm, t.addTab, plusIcon, "")
	newTabButton.Size = unit.Dp(15)
	newTabButton.Inset = layout.Inset{
		Top:    10,
		Bottom: 10,
		Left:   10,
		Right:  10,
	}
	flexChildren[len(flexChildren)-1] = Rigid(newTabButton)
	return layout.Background{}.Layout(gtx,
		func(gtx C) D {
			// expand horizontal
			gtx.Constraints.Min.X = gtx.Constraints.Max.X

			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: tabsBarBg}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return D{Size: gtx.Constraints.Min}
		},
		func(gtx C) D {
			return tabsMargin.Layout(gtx, func(gtx C) D {
				return flex.Layout(gtx, flexChildren...)
			})
		},
	)
}

// AddTab adds a new tab to the Tabs (no select happen)
func (t *Tabs) AddTab() {
	t.Tabs = append(t.Tabs, newTab())
}

func (t *Tabs) Select(idx int) bool {
	if idx >= len(t.Tabs) {
		return false
	}

	t.Selected = idx
	return true
}

func (t Tabs) SelectedTab() *Tab {
	return t.Tabs[t.Selected]
}

func (t Tabs) AddTabClicked(gtx C) bool {
	return t.addTab.Clicked(gtx)
}

// TabClicked return index of the clicked tab if exist, otherwise returns -1
func (t Tabs) TabClicked(gtx C) int {
	for idx, tab := range t.Tabs {
		if tab.clickable.Clicked(gtx) {
			return idx
		}
	}
	return -1
}

func (t *Tab) Layout(thm *Theme, gtx C, isSelected bool) D {
	// TODO: change the name and icon if not hard
	tabMargin := layout.Inset{
		Right: unit.Dp(3),
	}
	tab := material.Button(thm, t.clickable, "New Tab")
	// tab.TextSize = thm.TextSize * 0.75
	tab.Inset.Left = unit.Dp(15)
	tab.Inset.Right = unit.Dp(15)
	tab.CornerRadius = unit.Dp(8)

	if isSelected {
		// TODO: button use ContrastBg by default, so I'm forced to only use Fg.
		// I think we should have GazerTheme type. That's a big style revolution.
		tab.Background = thm.Fg
	}

	return tabMargin.Layout(gtx, func(gtx C) D { return tab.Layout(gtx) })
}

func newTab() *Tab {
	clickable := new(widget.Clickable)
	searchEditor := setupSearchBarEditor()
	return &Tab{clickable: clickable, SearchEditor: searchEditor}
}
