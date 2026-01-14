package css

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

// StyleSet is a (almost) ready-to-use style set of one CSS file (or more?)
type StyleSet struct {
	Universal   *Style
	IdStyles    map[string]*Style
	ClassStyles map[string]*Style
	TagStyles   map[parser.Tag]*Style
}

// Style is a property to style the rendering of any argument.
// The responsibility to intepret the struct is on caller.
// Change this =
// 1. modify how to parse the style name in Style.registerDecls
// 2. modify comparison in styleEq function (in style_test.go)
// 3. modify how to implement it in ui package
type Style struct {
	Color      *color.NRGBA
	BgColor    *color.NRGBA
	Margin     *layout.Inset
	Padding    *layout.Inset
	Border     *widget.Border
	FontSize   *unit.Sp // TODO: might have to change after supporting other type
	FontWeight *font.Weight
	FontStyle  *font.Style
}

// AddStyleSet adds 2 style sets with different importance (high/low priority)
func AddStyleSet(high, low *StyleSet) *StyleSet {
	if high == nil {
		return low
	} else if low == nil {
		return high
	}

	res := new(StyleSet)
	res.Universal = AddStylePtr(high.Universal, low.Universal)
	res.IdStyles = AddStyleMap(high.IdStyles, low.IdStyles)
	res.ClassStyles = AddStyleMap(high.ClassStyles, low.ClassStyles)
	res.TagStyles = AddStyleMap(high.TagStyles, low.TagStyles)

	return res
}

// AddStyleMap adds 2 style map (e.g. IdStyles) with different priority into one style map while
func AddStyleMap[K comparable](high, low map[K]*Style) map[K]*Style {
	var res map[K]*Style
	if high != nil {
		res = high
		for id, lowStyle := range low {
			highStyle, conflict := res[id]
			if conflict {
				res[id] = AddStylePtr(highStyle, lowStyle)
			} else {
				res[id] = lowStyle
			}
		}
	} else {
		res = low
	}
	return res
}

// AddStylePtr take 2 style pointers, and return a new style pointer that are
// the sum of both style, one with higher priority than another one
func AddStylePtr(sHigh *Style, sLow *Style) *Style {
	if sHigh == nil {
		return sLow
	} else if sLow == nil {
		return sHigh
	}

	res := new(Style)
	if sHigh.Color != nil {
		res.Color = sHigh.Color
	} else {
		res.Color = sLow.Color
	}

	if sHigh.BgColor != nil {
		res.BgColor = sHigh.BgColor
	} else {
		res.BgColor = sLow.BgColor
	}

	if sHigh.Margin != nil {
		res.Margin = sHigh.Margin
	} else {
		res.Margin = sLow.Margin
	}

	if sHigh.Border != nil {
		res.Border = sHigh.Border
	} else {
		res.Border = sLow.Border
	}

	if sHigh.FontSize != nil {
		res.FontSize = sHigh.FontSize
	} else {
		res.FontSize = sLow.FontSize
	}

	return res
}

// AddStyle takes 2 Style and merge 2 styles, one with high priority, one with low
func AddStyle(sHigh Style, sLow Style) Style {
	var res Style
	if sHigh.Color != nil {
		res.Color = sHigh.Color
	} else {
		res.Color = sLow.Color
	}

	if sHigh.BgColor != nil {
		res.BgColor = sHigh.BgColor
	} else {
		res.BgColor = sLow.BgColor
	}

	if sHigh.Margin != nil {
		res.Margin = sHigh.Margin
	} else {
		res.Margin = sLow.Margin
	}

	if sHigh.Border != nil {
		res.Border = sHigh.Border
	} else {
		res.Border = sLow.Border
	}

	if sHigh.FontSize != nil {
		res.FontSize = sHigh.FontSize
	} else {
		res.FontSize = sLow.FontSize
	}

	return res
}

// applyRule applies the css rule to the style set
func (s *StyleSet) applyRule(r rule) {
	for _, selector := range r.selectors {
		if selector == "*" {
			if s.Universal == nil {
				s.Universal = new(Style)
			}
			s.Universal.registerDecls(r.styles)
		} else if id, ok := strings.CutPrefix(selector, "#"); ok {
			style, ok := s.IdStyles[id]
			if !ok {
				style = new(Style)
				s.IdStyles[id] = style
			}
			style.registerDecls(r.styles)
		} else if class, ok := strings.CutPrefix(selector, "."); ok {
			// TODO: support tag.class syntax
			style, ok := s.ClassStyles[class]
			if !ok {
				style = new(Style)
				s.ClassStyles[class] = style
			}
			style.registerDecls(r.styles)
		} else {
			// tag name is case-insensitive
			tag, ok := parser.TagMap[strings.ToLower(selector)]
			if !ok {
				continue // tag not supported, skip
			}
			style, ok := s.TagStyles[tag]
			if !ok {
				style = new(Style)
				s.TagStyles[tag] = style
			}
			style.registerDecls(r.styles)
		}
	}
}

func newStyleSet() *StyleSet {
	return &StyleSet{
		Universal:   new(Style),
		IdStyles:    make(map[string]*Style),
		ClassStyles: make(map[string]*Style),
		TagStyles:   make(map[parser.Tag]*Style),
	}
}

// for debugging
func (s StyleSet) String() string {
	var builder strings.Builder

	builder.WriteString("{\n")
	builder.WriteString("\t" + "Universal: " + s.Universal.String() + "\n")

	builder.WriteString("\t" + "ids: " + "\n")
	for selector, style := range s.IdStyles {
		builder.WriteString("\t\t" + selector + ": " + style.String() + "\n")
	}
	builder.WriteString("\t" + "classes: " + "\n")
	for selector, style := range s.ClassStyles {
		builder.WriteString("\t\t" + selector + ": " + style.String() + "\n")
	}
	builder.WriteString("\t" + "tags: " + "\n")
	for tag, style := range s.TagStyles {
		builder.WriteString("\t\t" + tag.String() + ": " + style.String() + "\n")
	}

	return builder.String()
}

func (s Style) String() string {
	var builder strings.Builder

	builder.WriteString("{ ")
	if s.Color != nil {
		fmt.Fprintf(&builder, "Color: %v ", *s.Color)
	}
	if s.BgColor != nil {
		fmt.Fprintf(&builder, "BgColor: %v ", *s.BgColor)
	}
	if s.Margin != nil {
		fmt.Fprintf(&builder, "Margin: %v ", *s.Margin)
	}
	if s.Padding != nil {
		fmt.Fprintf(&builder, "Padding: %v ", *s.Padding)
	}
	if s.Border != nil {
		fmt.Fprintf(&builder, "Border: %v ", *s.Border)
	}
	if s.FontSize != nil {
		fmt.Fprintf(&builder, "FontSize: %v ", *s.FontSize)
	}
	if s.FontWeight != nil {
		fmt.Fprintf(&builder, "FontWeight: %v", *s.FontWeight)
	}
	builder.WriteString("}")
	return builder.String()
}

// registerDecls register CSS declarations (e.g. "color: red; fontSize: 10px") into the style struct
func (s *Style) registerDecls(decls map[string]string) {
	for prop, val := range decls {
		switch prop {
		case "color":
			c, err := s.parseColor(val)
			if err != nil {
				continue
			}
			s.Color = c
		case "background-color":
			c, err := s.parseColor(val)
			if err != nil {
				continue
			}
			s.BgColor = c
		case "margin":
			// TODO: "auto" value of margin is very interesting
			inset, err := s.parseInset(val)
			if err != nil {
				continue
			}
			s.Margin = inset
		case "margin-left":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Margin == nil {
				s.Margin = new(layout.Inset)
			}
			s.Margin.Left = length
		case "margin-right":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Margin == nil {
				s.Margin = new(layout.Inset)
			}
			s.Margin.Right = length
		case "margin-top":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Margin == nil {
				s.Margin = new(layout.Inset)
			}
			s.Margin.Top = length
		case "margin-bottom":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Margin == nil {
				s.Margin = new(layout.Inset)
			}
			s.Margin.Bottom = length
		case "border-width":
			width, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Border == nil {
				s.Border = new(widget.Border)
			}
			s.Border.Width = width
		case "border-radius":
			radius, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Border == nil {
				s.Border = new(widget.Border)
			}
			s.Border.CornerRadius = radius
		case "border-color":
			c, err := s.parseColor(val)
			if err != nil {
				continue
			}
			if s.Border == nil {
				s.Border = new(widget.Border)
			}
			s.Border.Color = *c
		case "padding":
			inset, err := s.parseInset(val)
			if err != nil {
				continue
			}
			s.Padding = inset
		case "padding-left":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Padding == nil {
				s.Padding = new(layout.Inset)
			}
			s.Padding.Left = length
		case "Padding-right":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Padding == nil {
				s.Padding = new(layout.Inset)
			}
			s.Padding.Right = length
		case "Padding-top":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Padding == nil {
				s.Padding = new(layout.Inset)
			}
			s.Padding.Top = length
		case "Padding-bottom":
			length, err := s.parseLength(val)
			if err != nil {
				continue
			}
			if s.Padding == nil {
				s.Padding = new(layout.Inset)
			}
			s.Padding.Left = length
		case "font-size":
			size, err := s.parseLength(val)
			if err != nil {
				continue
			}
			spSize := unit.Sp(size)
			s.FontSize = &spSize
		case "font-weight":
			weight, ok := fontWeights[val]
			if !ok {
				continue
			}
			s.FontWeight = &weight
		case "font-style":
			fstyle, ok := fontStyles[val]
			if !ok {
				continue
			}
			s.FontStyle = &fstyle
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

		r, errR := strconv.ParseInt(params[0], 10, 64)
		g, errG := strconv.ParseInt(params[1], 10, 64)
		b, errB := strconv.ParseInt(params[2], 10, 64)
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

		r, errR := strconv.ParseInt(params[0], 10, 64)
		g, errG := strconv.ParseInt(params[1], 10, 64)
		b, errB := strconv.ParseInt(params[2], 10, 64)
		a, errA := strconv.ParseInt(params[3], 10, 64)
		if errR != nil || errG != nil || errB != nil || errA != nil {
			return nil, fmt.Errorf(
				"Parse number error; r: %v, g: %v, b: %v, a: %v", errR, errG, errB, errA)
		}
		return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
	}

	// #RRGGBB (hex)
	if raw[0] == '#' && len(raw) == 7 {
		r, errR := strconv.ParseInt(raw[1:3], 16, 64)
		g, errG := strconv.ParseInt(raw[3:5], 16, 64)
		b, errB := strconv.ParseInt(raw[5:7], 16, 64)
		if errR != nil || errG != nil || errB != nil {
			return nil, fmt.Errorf("Parse number error; r: %v, g: %v, b: %v", errR, errG, errB)
		}
		return &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
	}

	// #RGB (hex)
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

// parseInset parses raw css string value that represent inset (e.g. margin, Padding)
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
