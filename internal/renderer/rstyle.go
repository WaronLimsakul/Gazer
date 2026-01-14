package renderer

import (
	"github.com/WaronLimsakul/Gazer/internal/css"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

// Contextual style passed around in renderNode
type RenderingStyle struct {
	base  css.Style
	label ui.LabelExtraStyle
	// NOTE: can add more extra style for other type of node
}

func newRenderingStyle() RenderingStyle {
	return RenderingStyle{base: css.Style{}, label: ui.NewLabelExtraStyle()}
}

// getLabelStyle creates ui.LabelStyle from rendering style information
func (r RenderingStyle) getLabelStyle() ui.LabelStyle {
	return ui.LabelStyle{Base: r.base, Extra: r.label}
}

// updateLabelStyle update rendering style according to ls
func (r *RenderingStyle) updateLabelStyle(ls ui.LabelStyle) {
	r.base = ls.Base
	r.label = ls.Extra
}
