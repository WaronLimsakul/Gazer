package main

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

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Gazer"))
		w.Option(app.Size(unit.Dp(600), unit.Dp(800)))
		w.Option(app.MinSize(unit.Dp(600), unit.Dp(800)))
		w.Option(app.MaxSize(unit.Dp(600), unit.Dp(800)))

		ops := op.Ops{}
		thm := material.NewTheme()

		srcInput := widget.Editor{}
		srcInput.Alignment = text.Middle
		srcInput.SingleLine = true
		srcInput.Submit = true

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
					default:
						continue
					}

				}

				srcInputUi := material.Editor(thm, &srcInput, "search")
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
	}()

	app.Main()
}
