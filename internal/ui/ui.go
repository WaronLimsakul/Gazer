package ui

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/engine"
)

const (
	WINDOW_WIDTH  = 1600
	WINDOW_HEIGHT = 900
)

// Draw takes gio's Window and Gazer's state
// and keep redrawing according to state
func Draw(window *app.Window, state *engine.State) {
	ops := op.Ops{}
	thm := material.NewTheme()
	srcInput := setupSrcInput()

	for {
		switch ev := window.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, ev)

			// Handle user search behavior
			for {
				editorEv, ok := srcInput.Update(gtx)
				if !ok {
					break
				}

				switch editorEv.(type) {
				// TODO: while loading, show something instead
				case widget.SubmitEvent:
					state.Url = srcInput.Text()
					state.Notifier <- engine.Search
				default:
					continue
				}

			}

			srcInputUi := material.Editor(thm, srcInput, "search")
			srcInputUi.TextSize = unit.Sp(20)

			inputMargin := layout.Inset{
				Top:    unit.Dp(25),
				Bottom: unit.Dp(25),
				Left:   unit.Dp(25),
				Right:  unit.Dp(25),
			}

			flexChildren := []layout.FlexChild{
				// search bar
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					border := widget.Border{
						Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
						CornerRadius: unit.Dp(2),
						Width:        unit.Dp(1),
					}
					return inputMargin.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return border.Layout(gtx, srcInputUi.Layout)
					})
				}),
			}

			siteMargin := layout.Inset{
				Left:  unit.Dp(25),
				Right: unit.Dp(25),
			}

			siteMargin.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// children from DOM rendering
				flexChildren = append(flexChildren, renderDOM(thm, state.Root)...)
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

// setupSrcInput create a new widget.Editor used as
// input behavior for search component
func setupSrcInput() *widget.Editor {
	srcInput := new(widget.Editor)
	srcInput.Alignment = text.Middle
	srcInput.SingleLine = true
	srcInput.Submit = true
	return srcInput
}

func newTheme() *material.Theme {
	thm := material.NewTheme()
	thm.Palette = material.Palette{
		Bg: color.NRGBA{R: 0xf5, G: 0xf5, B: 0xf5, A: 0xff}, // Light gray background
		Fg: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}, // Black text
		// Adjust ContrastBg, ContrastFg, etc.
	}
	return thm
}
