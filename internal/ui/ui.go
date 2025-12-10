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
)

type UserInput struct {
	Url *string
}

const WINDOW_WIDTH = 600
const WINDOW_HEIGHT = 800

func Draw(input UserInput) {
	w := setupWindow()
	ops := op.Ops{}
	thm := material.NewTheme()
	srcInput := setupSrcInput()

	for {
		switch ev := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, ev)

			for {
				editorEv, ok := srcInput.Update(gtx)
				if !ok {
					break
				}

				switch editorEv.(type) {
				case widget.SubmitEvent:
					*input.Url = srcInput.Text()
				default:
					continue
				}

			}

			srcInputUi := material.Editor(thm, srcInput, "search")
			srcInputUi.TextSize = unit.Sp(20)

			margin := layout.Inset{
				Top:    unit.Dp(25),
				Bottom: unit.Dp(25),
				Left:   unit.Dp(25),
				Right:  unit.Dp(25),
			}

			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					border := widget.Border{
						Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
						CornerRadius: unit.Dp(2),
						Width:        unit.Dp(1),
					}
					return margin.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return border.Layout(gtx, srcInputUi.Layout)
					})
				}),
			)

			ev.Frame(gtx.Ops)
		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

func setupWindow() *app.Window {
	w := new(app.Window)
	w.Option(app.Title("Gazer"))
	w.Option(app.Size(unit.Dp(WINDOW_WIDTH), unit.Dp(WINDOW_HEIGHT)))
	w.Option(app.MinSize(unit.Dp(WINDOW_WIDTH), unit.Dp(WINDOW_HEIGHT)))
	w.Option(app.MaxSize(unit.Dp(WINDOW_WIDTH), unit.Dp(WINDOW_HEIGHT)))
	return w
}

func setupSrcInput() *widget.Editor {
	srcInput := new(widget.Editor)
	srcInput.Alignment = text.Middle
	srcInput.SingleLine = true
	srcInput.Submit = true
	return srcInput
}
