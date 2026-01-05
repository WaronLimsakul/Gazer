package ui

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type InputType uint8

// All supported input type.
// Change this = change the InputTypes below
const (
	TextInput InputType = iota
	PasswordInput
	NumberInput
	// TODO: Email, Checkbox
)

// for rendering from DOM node
var InputTypes = map[string]InputType{
	"text":     TextInput,
	"password": PasswordInput,
	"number":   NumberInput,
}

type Input struct {
	thm       *Theme
	inputType InputType
	hint      string
	editor    *widget.Editor
	// size?
	// border?
	// margin?
	// min-width?
}

func NewInput(thm *Theme, inputType InputType, editor *widget.Editor, hint string) Input {
	editor.SingleLine = true
	switch inputType {
	case TextInput:
		editor.InputHint = key.HintText
	case PasswordInput:
		editor.Mask = '●'
		editor.InputHint = key.HintPassword
	case NumberInput:
		editor.InputHint = key.HintNumeric
	}

	return Input{thm: thm, inputType: inputType, editor: editor, hint: hint}
}

func (i Input) Layout(gtx C) D {
	border := widget.Border{Color: i.thm.Fg, CornerRadius: unit.Dp(1), Width: unit.Dp(1)}
	contentMargin := layout.UniformInset(unit.Dp(4))
	input := material.Editor(i.thm, i.editor, i.hint)
	minWidth := unit.Dp(100)
	return border.Layout(gtx, func(gtx C) D {
		return contentMargin.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Dp(minWidth)
			return input.Layout(gtx)
		})
	})
}
