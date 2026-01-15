package renderer

import (
	"github.com/WaronLimsakul/Gazer/internal/css"
	"github.com/WaronLimsakul/Gazer/internal/parser"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

// Contextual style passed around in renderNode
type RenderingContext struct {
	ancestors []parser.Tag // ancestors history of the current, node, oldest at index 0
	base      css.Style
	label     ui.LabelExtraStyle
	// NOTE: can add more extra style for other type of node
}

func newRenderingContext() RenderingContext {
	return RenderingContext{ancestors: make([]parser.Tag, 0)}
}

// getLabelStyle creates ui.LabelStyle from rendering context information
func (r RenderingContext) getLabelStyle() ui.LabelStyle {
	return ui.LabelStyle{Base: r.base, Extra: r.label}
}

// updateLabelStyle update rendering context according to ls
func (r *RenderingContext) updateLabelStyle(ls ui.LabelStyle) {
	r.base = ls.Base
	r.label = ls.Extra
}
