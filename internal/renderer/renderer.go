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
	"gioui.org/widget"
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

type Element = interface {
	Layout(gtx C) D
}

// Draw takes gio's Window and Gazer's state
// and keep redrawing according to state
func Draw(window *app.Window, state *engine.State) {
	ops := op.Ops{}
	thm := newTheme()
	searchBar := ui.NewSearchBar(thm)
	domRenderer := newDomRenderer(thm, state.Url)
	pageWidget := new(widget.List)
	pageWidget.Axis = layout.Vertical

	for {
		switch ev := window.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, ev)

			searchBar.Update(gtx)
			if searchBar.Searched(gtx) {
				state.Url = searchBar.Text()
				state.Notifier <- engine.Search
			}

			appFlexChildren := []layout.FlexChild{
				// TODO: write a component that hold its state
				rigid(searchBar),
				ui.HorizontalLine(thm, unit.Dp(WINDOW_WIDTH)),
			}

			// from now, handle website rendering
			if domRenderer.url != state.Url {
				domRenderer = newDomRenderer(thm, state.Url)
			}

			pageMargin := layout.Inset{
				Left:  unit.Dp(25),
				Right: unit.Dp(25),
			}

			// TODO: make it look better than this
			pageElements := domRenderer.render(state.Root)
			page := material.List(thm, pageWidget)
			appFlexChildren = append(appFlexChildren, layout.Rigid(func(gtx C) D {
				return pageMargin.Layout(gtx, func(gtx C) D {
					return page.Layout(gtx, len(pageElements), func(gtx C, idx int) D {
						line := pageElements[idx]
						if len(line) == 1 {
							return line[0].Layout(gtx)
						} else {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, elementsToFlexChildren(line)...)
						}
					})
				})
			}))

			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx, appFlexChildren...)

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

// elementsToFlexChildren wrap each element in elements with layout.Rigid and return
func elementsToFlexChildren(elements []Element) []layout.FlexChild {
	res := make([]layout.FlexChild, len(elements))
	for i, elem := range elements {
		res[i] = layout.Rigid(func(gtx C) D {
			return elem.Layout(gtx)
		})
	}
	return res
}
