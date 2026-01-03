package ui

import "gioui.org/layout"

// rigid recieves an element and wraps it with Rigid
func Rigid(e Element) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return e.Layout(gtx)
	})
}

// rigidMargin recieves an Inset and Element, wrap element with inset and wrap
// inset with rigid then return.
func RigidMargin(m layout.Inset, e Element) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return m.Layout(gtx, func(gtx C) D {
			return e.Layout(gtx)
		})
	})
}
