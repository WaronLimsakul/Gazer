package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	urlPkg "net/url"
	"strings"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/engine"
	"golang.org/x/exp/shiny/materialdesign/icons"

	_ "github.com/mat/besticon/ico"
)

type Tabs struct {
	Tabs     []*Tab
	addTab   *widget.Clickable
	thm      *Theme
	Selected int // idx of selected tab
}

type Tab struct {
	clickable      *widget.Clickable
	closeClickable *widget.Clickable // for "close tab" button
	SearchEditor   *widget.Editor
	Title          string
	// map url to fetched favicon (cache)
	favIcons map[string]image.Image
}

func NewTabs(thm *Theme) *Tabs {
	firstTab := newTab()
	tabs := []*Tab{firstTab}
	addTab := new(widget.Clickable)
	res := &Tabs{Tabs: tabs, addTab: addTab, thm: thm}
	res.Select(0)
	return res
}

func (t Tabs) Layout(gtx C, stateTabs []*engine.Tab) D {
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
			var url string
			if i < len(stateTabs) {
				url = stateTabs[i].Url
			}
			return tab.Layout(t.thm, gtx, isSelected, url)
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

func (t *Tabs) DeleteTab(idx int) {
	t.Tabs = append(t.Tabs[:idx], t.Tabs[idx+1:]...)
	// deal with selected idx shifting
	if idx <= t.Selected {
		if t.Selected == 0 {
			t.Selected++
		} else {
			t.Selected--
		}
	}
}

// TabClosed returns index of the tab that got "close tab" clicked if exist, otherwise returns -1
func (t Tabs) TabClosed(gtx C) int {
	for idx, tab := range t.Tabs {
		if tab.closeClickable.Clicked(gtx) {
			return idx
		}
	}
	return -1
}

func (t *Tab) Layout(thm *Theme, gtx C, isSelected bool, url string) D {
	tabMargin := layout.Inset{
		Right: unit.Dp(3),
	}

	title := t.Title
	if title == "" {
		title = "New Tab"
	}

	// get favicon from the cache
	favicon, ok := t.favIcons[url]
	if !ok {
		fetched, _ := t.getFavIcon(url)
		favicon = fetched
		// even if getFavIcon fail, we still cache nil which means it fail
		t.favIcons[url] = fetched
	}

	// if we already try and failed, then use default favicon
	if favicon == nil {
		favicon = defaultFavIcon
	}

	return tabMargin.Layout(gtx, func(gtx C) D {
		return t.clickable.Layout(gtx, func(gtx C) D {
			// check the tab content size first
			macro := op.Record(gtx.Ops)
			tabContentDim := layout.Inset{
				Top: unit.Dp(8), Bottom: unit.Dp(8),
				Left: unit.Dp(15), Right: unit.Dp(15),
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if favicon == nil {
							return D{}
						}
						gtx.Constraints.Max = image.Point{X: gtx.Dp(16), Y: gtx.Dp(16)}
						img := widget.Image{
							Src: paint.NewImageOp(favicon),
							Fit: widget.Contain,
						}
						return img.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx C) D {
						label := material.Body1(thm, title)
						if isSelected {
							label.Color = thm.Bg
						}
						return label.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx C) D {
						closeIcon, err := widget.NewIcon(icons.NavigationCancel)
						if err != nil {
							panic("can't decode close icon")
						}
						closeButton := material.IconButton(thm, t.closeClickable, closeIcon, "")
						closeButton.Size = unit.Dp(20)
						closeButton.Inset = layout.Inset{}
						return closeButton.Layout(gtx)
					}),
				)
			})
			tabContentOp := macro.Stop()

			// NOTE: can do this or use layout.Background{}
			// draw background
			tabBgColor := thm.Bg
			if isSelected {
				tabBgColor = thm.Fg
			}
			tabShape := clip.UniformRRect(image.Rectangle{Max: tabContentDim.Size}, gtx.Dp(8))
			defer tabShape.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, tabBgColor)

			// draw the content on top
			tabContentOp.Add(gtx.Ops)

			return tabContentDim
		})
	})
}

// getFavIcon fetch favicon.ico from the raw string address then decode and return it in image.Image
func (t Tab) getFavIcon(raw string) (image.Image, error) {
	if raw == "" {
		return nil, fmt.Errorf("Empty raw string")
	}

	if !strings.HasPrefix(raw, "https://") && !strings.HasPrefix(raw, "http://") {
		raw = "https://" + raw
	}

	url, err := urlPkg.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %v", err)
	}

	url.RawQuery = ""
	url.RawFragment = ""
	url.Path = "/favicon.ico"

	reqUrl := url.String()
	res, err := http.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %v", err)
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, fmt.Errorf("image.Decode: %v", err)
	}

	return img, nil
}

func newTab() *Tab {
	clickable := new(widget.Clickable)
	closeClickable := new(widget.Clickable)
	searchEditor := setupSearchBarEditor()
	return &Tab{
		clickable:      clickable,
		closeClickable: closeClickable,
		SearchEditor:   searchEditor,
		favIcons:       make(map[string]image.Image)}
}
