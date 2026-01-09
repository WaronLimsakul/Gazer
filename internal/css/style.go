package css

import (
	"fmt"
	"image/color"
	"maps"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

// StyleSet is a (almost) ready-to-use style set of one CSS file (or more?)
type StyleSet struct {
	universal   *Style
	idStyles    map[string]*Style
	classStyles map[string]*Style
	tagStyles   map[parser.Tag]*Style
}

// Style is a property to style the rendering of any argument.
// The responsibility to intepret the struct is on caller.
// Change this =
// 1. modify how to parse the style name in Style.apply
// 2. modify the ui component that support style to have the suppored fields
// 3. modify the implementation that ui component
// 4. modify comparison in styleEq function
type Style struct {
	color    *color.NRGBA
	bgColor  *color.NRGBA
	margin   *layout.Inset
	padding  *layout.Inset
	border   *widget.Border
	fontSize *unit.Dp // TODO: might have to change after supporting other type
}

// AddStyle adds 2 style struct, one with higher priority than another one
func AddStyle(sHigh Style, sLow Style) Style {
	var res Style
	if sHigh.color != nil {
		res.color = sHigh.color
	} else {
		res.color = sLow.color
	}

	if sHigh.bgColor != nil {
		res.bgColor = sHigh.bgColor
	} else {
		res.bgColor = sLow.bgColor
	}

	if sHigh.margin != nil {
		res.margin = sHigh.margin
	} else {
		res.margin = sLow.margin
	}

	if sHigh.border != nil {
		res.border = sHigh.border
	} else {
		res.border = sLow.border
	}

	if sHigh.fontSize != nil {
		res.fontSize = sHigh.fontSize
	} else {
		res.fontSize = sLow.fontSize
	}

	return res
}

// applyRule applies the css rule to the style set
func (s *StyleSet) applyRule(r rule) {
	for _, selector := range r.selectors {
		if selector == "*" {
			if s.universal == nil {
				s.universal = new(Style)
			}
			s.universal.add(r.styles)
		} else if id, ok := strings.CutPrefix(selector, "#"); ok {
			style, ok := s.idStyles[id]
			if !ok {
				style = new(Style)
				s.idStyles[id] = style
			}
			style.add(r.styles)
		} else if class, ok := strings.CutPrefix(selector, "."); ok {
			// TODO: support tag.clas syntax
			style, ok := s.classStyles[class]
			if !ok {
				style = new(Style)
				s.classStyles[class] = style
			}
			style.add(r.styles)
		} else {
			// tag name is case-insensitive
			tag, ok := parser.TagMap[strings.ToLower(selector)]
			if !ok {
				continue // tag not supported, skip
			}
			style, ok := s.tagStyles[tag]
			if !ok {
				style = new(Style)
				s.tagStyles[tag] = style
			}
			style.add(r.styles)
		}
	}
}

func newStyleSet() *StyleSet {
	return &StyleSet{
		universal:   nil,
		idStyles:    make(map[string]*Style),
		classStyles: make(map[string]*Style),
		tagStyles:   make(map[parser.Tag]*Style),
	}
}

// add adds the CSS declaration into the style struct
func (s *Style) add(decl map[string]string) {
	for prop, val := range decl {
		switch prop {
		case "color":
			c, err := s.parseColor(val)
			if err != nil {
				continue
			}
			s.color = c
		case "background-color":
			c, err := s.parseColor(val)
			if err != nil {
				continue
			}
			s.bgColor = c
		case "margin":
			// TODO: "auto" value of margin is very interesting
			inset, err := s.parseInset(val)
			if err != nil {
				continue
			}
			s.margin = inset
		case "margin-left":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.margin == nil {
				s.margin = new(layout.Inset)
			}
			s.margin.Left = length
		case "margin-right":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.margin == nil {
				s.margin = new(layout.Inset)
			}
			s.margin.Right = length
		case "margin-top":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.margin == nil {
				s.margin = new(layout.Inset)
			}
			s.margin.Top = length
		case "margin-bottom":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.margin == nil {
				s.margin = new(layout.Inset)
			}
			s.margin.Bottom = length
		case "border-width":
			width, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.border == nil {
				s.border = new(widget.Border)
			}
			s.border.Width = width
		case "border-radius":
			radius, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.border == nil {
				s.border = new(widget.Border)
			}
			s.border.CornerRadius = radius
		case "border-color":
			c, err := s.parseColor(val)
			if err != nil {
				continue
			}
			if s.border == nil {
				s.border = new(widget.Border)
			}
			s.border.Color = *c
		case "padding":
			inset, err := s.parseInset(val)
			if err != nil {
				continue
			}
			if s.padding == nil {
				s.padding = new(layout.Inset)
			}
			s.padding = inset
		case "padding-left":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.padding == nil {
				s.padding = new(layout.Inset)
			}
			s.padding.Left = length
		case "padding-right":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.padding == nil {
				s.padding = new(layout.Inset)
			}
			s.padding.Right = length
		case "padding-top":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.padding == nil {
				s.padding = new(layout.Inset)
			}
			s.padding.Top = length
		case "padding-bottom":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.padding == nil {
				s.padding = new(layout.Inset)
			}
			s.padding.Left = length
		}

	}
}

// parseLength parses a length string value into Dp unit
// e.g. parseLength("10px") -> unit.Dp(10)
func (s Style) parseLength(raw string) (unit.Dp, error) {
	// TODO: support em, rem, %
	pxlen, ok := strings.CutSuffix(raw, "px")
	if !ok {
		return unit.Dp(0), fmt.Errorf("unsupported unit: %s", raw)
	}

	res, err := strconv.ParseFloat(pxlen, 32)
	if err != nil {
		return unit.Dp(0), fmt.Errorf("strconv.ParseFloat: %v", err)
	}

	return unit.Dp(float32(res)), nil
}

func (s Style) parseColor(raw string) (*color.NRGBA, error) {
	// keyword color
	c, ok := colors[raw]
	if ok {
		return &c, nil
	}

	// rgb
	if strings.HasPrefix(raw, "rgb(") && raw[len(raw)-1] == ')' {
		content := raw[4 : len(raw)-1]
		params := strings.Split(content, ",")
		if len(params) != 3 {
			return nil, fmt.Errorf("Invalid input %v", params)
		}
		for i, param := range params {
			params[i] = strings.TrimSpace(param)
		}

		r, errR := strconv.ParseInt(params[0], 16, 64)
		g, errG := strconv.ParseInt(params[1], 16, 64)
		b, errB := strconv.ParseInt(params[2], 16, 64)
		if errR != nil || errG != nil || errB != nil {
			return nil, fmt.Errorf("Parse number error; r: %v, g: %v, b: %v", errR, errG, errB)
		}

		return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
	}

	// rgba
	if strings.HasPrefix(raw, "rgba(") && raw[len(raw)-1] == ')' {
		content := raw[4 : len(raw)-1]
		params := strings.Split(content, ",")
		if len(params) != 4 {
			return &color.NRGBA{}, fmt.Errorf("Invalid input %v", params)
		}
		for i, param := range params {
			params[i] = strings.TrimSpace(param)
		}

		r, errR := strconv.ParseInt(params[0], 16, 64)
		g, errG := strconv.ParseInt(params[1], 16, 64)
		b, errB := strconv.ParseInt(params[2], 16, 64)
		a, errA := strconv.ParseInt(params[3], 16, 64)
		if errR != nil || errG != nil || errB != nil || errA != nil {
			return nil, fmt.Errorf(
				"Parse number error; r: %v, g: %v, b: %v, a: %v", errR, errG, errB, errA)
		}
		return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
	}

	// #RRGGBB
	if raw[0] == '#' && len(raw) == 7 {
		r, errR := strconv.ParseInt(raw[1:3], 16, 64)
		g, errG := strconv.ParseInt(raw[3:5], 16, 64)
		b, errB := strconv.ParseInt(raw[5:7], 16, 64)
		if errR != nil || errG != nil || errB != nil {
			return nil, fmt.Errorf("Parse number error; r: %v, g: %v, b: %v", errR, errG, errB)
		}
		return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
	}

	// #RGB
	if raw[0] == '#' && len(raw) == 4 {
		r, errR := strconv.ParseInt(string(raw[1])+string(raw[1]), 16, 64)
		g, errG := strconv.ParseInt(string(raw[2])+string(raw[2]), 16, 64)
		b, errB := strconv.ParseInt(string(raw[3])+string(raw[3]), 16, 64)
		if errR != nil || errG != nil || errB != nil {
			return nil, fmt.Errorf("Parse number error; r: %v, g: %v, b: %v", errR, errG, errB)
		}
		return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
	}

	return nil, fmt.Errorf("Invalid format: %v", raw)
}

// parseInset parses raw css string value that represent inset (e.g. margin, padding)
func (s Style) parseInset(raw string) (*layout.Inset, error) {
	vals := strings.Fields(raw)
	length := len(vals)
	if length > 4 || length < 1 {
		return nil, fmt.Errorf("Invalid format: %v", vals)
	}
	// I have OCD
	for i, v := range vals {
		vals[i] = strings.TrimSpace(v)
	}
	switch length {
	case 4: // top right bottom left
		t, errT := s.parseLength(vals[0])
		r, errR := s.parseLength(vals[1])
		b, errB := s.parseLength(vals[2])
		l, errL := s.parseLength(vals[3])
		if errT != nil || errR != nil || errB != nil || errL != nil {
			return nil, fmt.Errorf(
				"Parse length error; t: %v | r: %v | b: %v | l: %v",
				errT, errR, errB, errL,
			)
		}
		return &layout.Inset{Top: t, Right: r, Bottom: b, Left: l}, nil

	case 3: // top left-right bottom
		t, errT := s.parseLength(vals[0])
		lr, errLR := s.parseLength(vals[1])
		b, errB := s.parseLength(vals[2])
		if errT != nil || errLR != nil || errB != nil {
			return nil, fmt.Errorf(
				"Parse length error; t: %v | lr: %v | b: %v",
				errT, errLR, errB,
			)
		}
		return &layout.Inset{Top: t, Right: lr, Bottom: b, Left: lr}, nil
	case 2: // top-bottom left-right
		tb, errTB := s.parseLength(vals[0])
		lr, errLR := s.parseLength(vals[1])
		if errTB != nil || errLR != nil {
			return nil, fmt.Errorf(
				"Parse length error; tb: %v | lr: %v",
				errTB, errLR,
			)
		}
		return &layout.Inset{Top: tb, Bottom: tb, Left: lr, Right: lr}, nil
	case 1: // all
		m, err := s.parseLength(vals[0])
		if err != nil {
			return nil, fmt.Errorf("Parse length error: %v", err)
		}
		res := layout.UniformInset(m)
		return &res, nil
	}

	return nil, fmt.Errorf("Invalid format: %v", vals)
}

func styleSetEq(a, b StyleSet) bool {
	return styleEqual(a.universal, b.universal) &&
		maps.EqualFunc(a.idStyles, b.idStyles, styleEqual) &&
		maps.EqualFunc(a.classStyles, b.classStyles, styleEqual) &&
		maps.EqualFunc(a.tagStyles, b.tagStyles, styleEqual)
}

func ptrValEq[T comparable](a *T, b *T) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	} else {
		return *a == *b
	}
}

func styleEqual(a, b *Style) bool {
	return ptrValEq(a.color, b.color) && ptrValEq(a.bgColor, b.bgColor) &&
		ptrValEq(a.margin, b.margin) && ptrValEq(a.border, b.border) &&
		ptrValEq(a.fontSize, b.fontSize)
}
