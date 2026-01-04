package renderer

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/engine"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

const (
	WINDOW_WIDTH  = 1600
	WINDOW_HEIGHT = 900
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Element = ui.Element
type Noti = engine.Notification

// Draw takes gio's Window and Gazer's state
// and keep redrawing according to state
func Draw(window *app.Window, state *engine.State) {
	ops := op.Ops{}
	thm := newTheme()

	hLine := ui.HorizontalLine{Thm: thm, Width: WINDOW_WIDTH, Height: unit.Dp(1)}
	page := ui.NewPage(thm)     // page doesn't depend on the tab
	tabsView := ui.NewTabs(thm) // will have another "tabs" from state
	domRenderers := map[*ui.Tab]*DomRenderer{}

	for {
		switch ev := window.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, ev)

			tabView := tabsView.SelectedTab()

			tabs := state.Tabs
			tab := tabs[tabsView.Selected]

			searchBar := ui.NewSearchBar(thm, tabView.SearchEditor)

			// handle search bar event
			searchBar.RenderInteraction(gtx)
			if searchBar.Searched(gtx) {
				state.Tabs[tabsView.Selected].Url = searchBar.Text()
				state.Notifier <- Noti{Type: engine.Search, TabIdx: tabsView.Selected}
			}

			// get the cached dom renderer
			domRenderer, ok := domRenderers[tabView]
			if !ok {
				domRenderer = newDomRenderer(thm, tabView)
				domRenderers[tabView] = domRenderer
			}

			// handle link clicking event
			jump, url := domRenderer.linkClicked(gtx)
			if jump {
				searchBar.SetText(url)
				tab.Url = searchBar.Text()
				state.Notifier <- Noti{Type: engine.Search, TabIdx: tabsView.Selected}
			}

			// handle clicking add tab button
			if tabsView.AddTabClicked(gtx) {
				tabsView.AddTab()
				tabsView.Select(len(tabsView.Tabs) - 1)
				state.Notifier <- Noti{Type: engine.AddTab}
			}

			// handle clicking tab
			tabClicked := tabsView.TabClicked(gtx)
			if tabClicked != -1 {
				tabsView.Select(tabClicked)
				window.Invalidate()
			}

			// start render app
			appFlex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}
			appFlexChildren := []layout.FlexChild{
				layout.Rigid(func(gtx C) D { return tabsView.Layout(gtx, tabs) }),
				ui.Rigid(searchBar),
			}

			// if loading the page, replace horizontal line with progress bar
			if tab.IsLoading {
				progress := <-state.LoadProgress
				appFlexChildren = append(appFlexChildren, ui.Rigid(material.ProgressBar(thm, progress)))
			} else {
				appFlexChildren = append(appFlexChildren, ui.Rigid(hLine))
			}

			// handle page rendering
			domRenderer.update(gtx)          // update ui interaction
			domRenderer.handleHead(tab.Root) // set tab data
			pageElements := domRenderer.render(tab.Root)
			appFlexChildren = append(appFlexChildren, layout.Rigid(func(gtx C) D {
				return page.Layout(gtx, pageElements)
			}))

			appFlex.Layout(gtx, appFlexChildren...)

			ev.Frame(gtx.Ops)
		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

// NewWindow creates new Gazer window
func NewWindow() *app.Window {
	w := new(app.Window)
	w.Option(app.Title("Gazer"))
	w.Option(app.Size(unit.Dp(WINDOW_WIDTH), unit.Dp(WINDOW_HEIGHT)))
	w.Option(app.MinSize(unit.Dp(WINDOW_WIDTH), unit.Dp(WINDOW_HEIGHT)))
	w.Option(app.MaxSize(unit.Dp(WINDOW_WIDTH), unit.Dp(WINDOW_HEIGHT)))
	// w.Option(app.Decorated(false))
	return w
}

func newTheme() *material.Theme {
	thm := material.NewTheme()

	// use set font faces
	faces, err := loadFont()
	if err != nil {
		log.Fatalf("loadFont: %v", err)
	}
	thm.Shaper = text.NewShaper(text.WithCollection(faces))

	// Nordic Blue theme
	thm.Palette = material.Palette{
		Bg:         color.NRGBA{R: 236, G: 239, B: 244, A: 255},
		Fg:         color.NRGBA{R: 76, G: 86, B: 106, A: 255},
		ContrastBg: color.NRGBA{R: 94, G: 129, B: 172, A: 255},
		ContrastFg: color.NRGBA{R: 236, G: 239, B: 244, A: 255},
	}

	return thm
}
