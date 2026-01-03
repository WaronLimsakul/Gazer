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

// Draw takes gio's Window and Gazer's state
// and keep redrawing according to state
func Draw(window *app.Window, state *engine.State) {
	ops := op.Ops{}
	thm := newTheme()
	if state.Tabs == nil {
		state.Tabs = ui.NewTabs(thm)
	}

	hLine := ui.HorizontalLine{Thm: thm, Width: WINDOW_WIDTH, Height: unit.Dp(1)}
	domRenderer := newDomRenderer(thm, "")
	page := ui.NewPage(thm)

	for {
		switch ev := window.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, ev)

			tab := state.Tabs.SelectedTab()

			searchBar := ui.NewSearchBar(thm, tab.SearchEditor)

			searchBar.RenderInteraction(gtx)
			if searchBar.Searched(gtx) {
				tab.Url = searchBar.Text()
				state.Notifier <- engine.Search
			}

			// update state if user search click a link
			jump, url := domRenderer.linkClicked(gtx)
			if jump {
				searchBar.SetText(url)
				tab.Url = url
				state.Notifier <- engine.Search
			}

			if state.Tabs.AddTabClicked(gtx) {
				state.Notifier <- engine.AddTab
			}

			tabClicked := state.Tabs.TabClicked(gtx)
			if tabClicked != -1 {
				// TODO: I want to send the thing to engine to deal with instead
				state.Tabs.Select(tabClicked)
				window.Invalidate()
			}

			appFlex := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}
			appFlexChildren := []layout.FlexChild{
				layout.Rigid(func(gtx C) D { return state.Tabs.Layout(gtx) }),
				ui.Rigid(searchBar), ui.Rigid(hLine)}
			// if loading the page, replace horizontal line with progress bar
			if tab.IsLoading {
				progress := <-state.LoadProgress
				appFlexChildren[1] = ui.Rigid(material.ProgressBar(thm, progress))
			}

			// from now, handle website rendering
			if domRenderer.url != tab.Url {
				domRenderer = newDomRenderer(thm, tab.Url)
			} else {
				domRenderer.update(gtx)
			}

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
