package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Page is a component for rendering entire webpage
type Page struct {
	thm  *Theme
	list *widget.List
}

func NewPage(thm *Theme) *Page {
	list := new(widget.List)
	list.Axis = layout.Vertical
	list.Alignment = layout.Middle
	return &Page{thm: thm, list: list}
}

func (p *Page) Layout(gtx C, elements [][]Element) D {
	listUi := material.List(p.thm, p.list)

	pageMargin := layout.Inset{
		Left:  unit.Dp(5),
		Right: unit.Dp(5),
	}
	return pageMargin.Layout(gtx, func(gtx C) D {
		return listUi.Layout(gtx, len(elements), func(gtx C, idx int) D {
			line := elements[idx]
			if len(line) == 1 {
				return line[0].Layout(gtx)
			} else {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, elementsToFlexChildren(line)...)
			}
		})
	})
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
