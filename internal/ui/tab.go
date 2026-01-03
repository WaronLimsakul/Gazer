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
	"github.com/WaronLimsakul/Gazer/internal/parser"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

// TODO: consider further if we want "ui.Tabs" and also "engine.Tabs"
// cuz I feel like this ui.Tabs kinda do both logic and rendering
type Tabs struct {
	Tabs   []*Tab
	addTab *widget.Clickable
	thm    *Theme
}

type Tab struct {
	Url  string
	Root *parser.Node // DOM root

	IsSelected bool
	IsLoading  bool

	clickable *widget.Clickable

	SearchEditor *widget.Editor
}

func NewTabs(thm *Theme) *Tabs {
	firstTab := newTab("", nil)
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
			return tab.Layout(t.thm, gtx)
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
func (t *Tabs) AddTab(url string, root *parser.Node) {
	t.Tabs = append(t.Tabs, newTab(url, root))
}

func (t *Tabs) Select(idx int) bool {
	if idx >= len(t.Tabs) {
		return false
	}

	for _, tab := range t.Tabs {
		tab.IsSelected = false
	}

	t.Tabs[idx].IsSelected = true
	return true
}

func (t Tabs) SelectedTab() *Tab {
	for _, tab := range t.Tabs {
		if tab.IsSelected {
			return tab
		}
	}
	return nil
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

func (t *Tab) Layout(thm *Theme, gtx C) D {
	// TODO: change the name and icon if not hard
	tabMargin := layout.Inset{
		Right: unit.Dp(3),
	}
	tab := material.Button(thm, t.clickable, "New Tab")
	// tab.TextSize = thm.TextSize * 0.75
	tab.Inset.Left = unit.Dp(15)
	tab.Inset.Right = unit.Dp(15)
	tab.CornerRadius = unit.Dp(8)

	if t.IsSelected {
		// TODO: button use ContrastBg by default, so I'm forced to only use Fg.
		// I think we should have GazerTheme type. That's a big style revolution.
		tab.Background = thm.Fg
	}

	return tabMargin.Layout(gtx, func(gtx C) D { return tab.Layout(gtx) })
}

func newTab(url string, root *parser.Node) *Tab {
	clickable := new(widget.Clickable)
	searchEditor := setupSearchBarEditor()
	return &Tab{Url: url, Root: root, clickable: clickable, SearchEditor: searchEditor}
}
