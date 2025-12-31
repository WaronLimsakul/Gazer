package renderer

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/pointer"
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

// Draw takes gio's Window and Gazer's state
// and keep redrawing according to state
func Draw(window *app.Window, state *engine.State) {
	ops := op.Ops{}
	thm := newTheme()
	searchEditor := ui.SetupSearchEditor()
	searchClickable := new(widget.Clickable)
	domRenderer := newDomRenderer(thm, state.Url)

	for {
		switch ev := window.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, ev)

			// Handle user search behavior
			for {
				editorEv, ok := searchEditor.Update(gtx)
				if !ok {
					break
				}

				switch editorEv.(type) {
				// press "enter" search
				case widget.SubmitEvent:
					state.Url = searchEditor.Text()
					state.Notifier <- engine.Search
				default:
					continue
				}

			}

			// click search
			if searchClickable.Clicked(gtx) {
				state.Url = searchEditor.Text()
				state.Notifier <- engine.Search
			}

			// change the pointer cursor when hover
			if searchClickable.Hovered() {
				pointer.CursorPointer.Add(&ops)
			}

			flexChildren := []layout.FlexChild{
				ui.SearchBar(thm, searchEditor, searchClickable),
				ui.HorizontalLine(thm, unit.Dp(WINDOW_WIDTH)),
			}

			// from now, handle website rendering
			if domRenderer.url != state.Url {
				domRenderer = newDomRenderer(thm, state.Url)
			}

			siteMargin := layout.Inset{
				Left:  unit.Dp(25),
				Right: unit.Dp(25),
			}

			// TODO: make site margin apply only on site. It also apply to search bar
			siteMargin.Layout(gtx, func(gtx C) D {
				// children from DOM rendering
				flexChildren = append(flexChildren, domRenderer.render(state.Root)...)
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}.Layout(gtx, flexChildren...)
			})

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
