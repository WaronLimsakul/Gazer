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
			input:    [2]*Style{nil, {color: &red}},
			expected: &Style{color: &red},
		},
		{
			name: "high priority wins on conflicts",
			input: [2]*Style{
				{color: &red, fontSize: &fontSize16},
				{color: &blue, fontSize: &fontSize12},
			},
			expected: &Style{color: &red, fontSize: &fontSize16},
		},
		{
			name: "low fills gaps in high",
			input: [2]*Style{
				{color: &red},
				{fontSize: &fontSize12, margin: &margin10},
			},
			expected: &Style{color: &red, fontSize: &fontSize12, margin: &margin10},
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
			input:    [2]map[string]*Style{nil, {"a": {color: &red}}},
			expected: map[string]*Style{"a": {color: &red}},
		},
		{
			name: "no key conflicts",
			input: [2]map[string]*Style{
				{"a": {color: &red}},
				{"b": {color: &blue}},
			},
			expected: map[string]*Style{
				"a": {color: &red},
				"b": {color: &blue},
			},
		},
		{
			name: "merge on same key",
			input: [2]map[string]*Style{
				{"a": {color: &red}, "b": {fontSize: &fontSize16}},
				{"a": {fontSize: &fontSize12}, "c": {color: &blue}},
			},
			expected: map[string]*Style{
				"a": {color: &red, fontSize: &fontSize12},
				"b": {fontSize: &fontSize16},
				"c": {color: &blue},
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
			input:    [2]*StyleSet{nil, {universal: &Style{color: &red}}},
			expected: &StyleSet{universal: &Style{color: &red}},
		},
		{
			name: "merge universal and id styles",
			input: [2]*StyleSet{
				{
					universal: &Style{color: &red},
					idStyles:  map[string]*Style{"a": {fontSize: &fontSize16}},
				},
				{
					universal: &Style{fontSize: &fontSize12},
					idStyles:  map[string]*Style{"b": {color: &blue}},
				},
			},
			expected: &StyleSet{
				universal: &Style{color: &red, fontSize: &fontSize12},
				idStyles: map[string]*Style{
					"a": {fontSize: &fontSize16},
					"b": {color: &blue},
				},
			},
		},
		{
			name: "merge all style types with conflicts",
			input: [2]*StyleSet{
				{
					classStyles: map[string]*Style{"c1": {color: &red}},
					tagStyles:   map[parser.Tag]*Style{parser.Div: {fontSize: &fontSize16}},
				},
				{
					classStyles: map[string]*Style{"c1": {fontSize: &fontSize12}},
					tagStyles:   map[parser.Tag]*Style{parser.Span: {color: &blue}},
				},
			},
			expected: &StyleSet{
				classStyles: map[string]*Style{
					"c1": {color: &red, fontSize: &fontSize12},
				},
				tagStyles: map[parser.Tag]*Style{
					parser.Div:  {fontSize: &fontSize16},
					parser.Span: {color: &blue},
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
	return ptrValEq(a.color, b.color) && ptrValEq(a.bgColor, b.bgColor) &&
		ptrValEq(a.margin, b.margin) && ptrValEq(a.padding, b.padding) &&
		ptrValEq(a.border, b.border) && ptrValEq(a.fontSize, b.fontSize)
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
