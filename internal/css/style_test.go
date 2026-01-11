package css

import (
	"maps"
	"testing"

	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

func TestAddStyle(t *testing.T) {
	red := colors["red"]
	blue := colors["blue"]
	fontSize12 := unit.Dp(12)
	fontSize16 := unit.Dp(16)
	margin10 := layout.UniformInset(unit.Dp(10))

	cases := []struct {
		name     string
		input    [2]*Style
		expected *Style
	}{
		{
			name:     "nil handling",
			input:    [2]*Style{nil, {Color: &red}},
			expected: &Style{Color: &red},
		},
		{
			name: "high priority wins on conflicts",
			input: [2]*Style{
				{Color: &red, FontSize: &fontSize16},
				{Color: &blue, FontSize: &fontSize12},
			},
			expected: &Style{Color: &red, FontSize: &fontSize16},
		},
		{
			name: "low fills gaps in high",
			input: [2]*Style{
				{Color: &red},
				{FontSize: &fontSize12, Margin: &margin10},
			},
			expected: &Style{Color: &red, FontSize: &fontSize12, Margin: &margin10},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := AddStyle(tc.input[0], tc.input[1])
			if !styleEq(result, tc.expected) {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestMergeStyleMap(t *testing.T) {
	red := colors["red"]
	blue := colors["blue"]
	fontSize12 := unit.Dp(12)
	fontSize16 := unit.Dp(16)

	cases := []struct {
		name     string
		input    [2]map[string]*Style
		expected map[string]*Style
	}{
		{
			name:     "nil maps",
			input:    [2]map[string]*Style{nil, {"a": {Color: &red}}},
			expected: map[string]*Style{"a": {Color: &red}},
		},
		{
			name: "no key conflicts",
			input: [2]map[string]*Style{
				{"a": {Color: &red}},
				{"b": {Color: &blue}},
			},
			expected: map[string]*Style{
				"a": {Color: &red},
				"b": {Color: &blue},
			},
		},
		{
			name: "merge on same key",
			input: [2]map[string]*Style{
				{"a": {Color: &red}, "b": {FontSize: &fontSize16}},
				{"a": {FontSize: &fontSize12}, "c": {Color: &blue}},
			},
			expected: map[string]*Style{
				"a": {Color: &red, FontSize: &fontSize12},
				"b": {FontSize: &fontSize16},
				"c": {Color: &blue},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := AddStyleMap(tc.input[0], tc.input[1])
			if !maps.EqualFunc(result, tc.expected, styleEq) {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestAddStyleSet(t *testing.T) {
	red := colors["red"]
	blue := colors["blue"]
	fontSize12 := unit.Dp(12)
	fontSize16 := unit.Dp(16)

	cases := []struct {
		name     string
		input    [2]*StyleSet
		expected *StyleSet
	}{
		{
			name:     "nil stylesets",
			input:    [2]*StyleSet{nil, {universal: &Style{Color: &red}}},
			expected: &StyleSet{universal: &Style{Color: &red}},
		},
		{
			name: "merge universal and id styles",
			input: [2]*StyleSet{
				{
					universal: &Style{Color: &red},
					idStyles:  map[string]*Style{"a": {FontSize: &fontSize16}},
				},
				{
					universal: &Style{FontSize: &fontSize12},
					idStyles:  map[string]*Style{"b": {Color: &blue}},
				},
			},
			expected: &StyleSet{
				universal: &Style{Color: &red, FontSize: &fontSize12},
				idStyles: map[string]*Style{
					"a": {FontSize: &fontSize16},
					"b": {Color: &blue},
				},
			},
		},
		{
			name: "merge all style types with conflicts",
			input: [2]*StyleSet{
				{
					classStyles: map[string]*Style{"c1": {Color: &red}},
					tagStyles:   map[parser.Tag]*Style{parser.Div: {FontSize: &fontSize16}},
				},
				{
					classStyles: map[string]*Style{"c1": {FontSize: &fontSize12}},
					tagStyles:   map[parser.Tag]*Style{parser.Span: {Color: &blue}},
				},
			},
			expected: &StyleSet{
				classStyles: map[string]*Style{
					"c1": {Color: &red, FontSize: &fontSize12},
				},
				tagStyles: map[parser.Tag]*Style{
					parser.Div:  {FontSize: &fontSize16},
					parser.Span: {Color: &blue},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := AddStyleSet(tc.input[0], tc.input[1])
			if result == nil && tc.expected == nil {
				return
			}
			if result == nil || tc.expected == nil {
				t.Errorf("nil mismatch: got %v, want %v", result, tc.expected)
				return
			}
			if !styleSetEq(*result, *tc.expected) {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func styleSetEq(a, b StyleSet) bool {
	return styleEq(a.universal, b.universal) &&
		maps.EqualFunc(a.idStyles, b.idStyles, styleEq) &&
		maps.EqualFunc(a.classStyles, b.classStyles, styleEq) &&
		maps.EqualFunc(a.tagStyles, b.tagStyles, styleEq)
}

func styleEq(a, b *Style) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}
	return ptrValEq(a.Color, b.Color) && ptrValEq(a.BgColor, b.BgColor) &&
		ptrValEq(a.Margin, b.Margin) && ptrValEq(a.Padding, b.Padding) &&
		ptrValEq(a.Border, b.Border) && ptrValEq(a.FontSize, b.FontSize) &&
		ptrValEq(a.FontWeight, b.FontWeight)
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
