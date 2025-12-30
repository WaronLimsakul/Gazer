package ui

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
	"golang.org/x/exp/shiny/materialdesign/icons"
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
	thm := material.NewTheme()
	srcInput := setupSrcInput()
	searchClickable := new(widget.Clickable)

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
				// press "enter" search
				case widget.SubmitEvent:
					state.Url = srcInput.Text()
					state.Notifier <- engine.Search
				default:
					continue
				}

			}

			// click search
			if searchClickable.Clicked(gtx) {
				state.Url = srcInput.Text()
				state.Notifier <- engine.Search
			}

			flexChildren := []layout.FlexChild{newSearchBar(thm, srcInput, searchClickable)}

			siteMargin := layout.Inset{
				Left:  unit.Dp(25),
				Right: unit.Dp(25),
			}

			siteMargin.Layout(gtx, func(gtx C) D {
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
	srcInput.Alignment = text.Start
	srcInput.SingleLine = true
	srcInput.Submit = true
	return srcInput
}

func newSearchBar(
	thm *material.Theme,
	editor *widget.Editor,
	searchClickable *widget.Clickable,
) layout.FlexChild {
	srcInputUi := material.Editor(thm, editor, "search")
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
		Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
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
	searchButton := material.IconButton(thm, searchClickable, icon, "Search")
	searchButton.Size = unit.Dp(20)

	// search bar
	return layout.Rigid(func(gtx C) D {
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
	})
}
