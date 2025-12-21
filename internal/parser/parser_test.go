package parser

import (
	"maps"
	"testing"
)

func TestGetTag(t *testing.T) {
	testCases := map[string]Tag{
		"h1":   H1,
		"p":    P,
		"ahah": P,
		"br":   Br,
		"body": Body,
		"head": Head,
	}

	for test, expected := range testCases {
		if getTag(test) != expected {
			t.Errorf("Expected %v | Got %v", expected, getTag(test))
		}
	}
}

func TestAssignAttrs(t *testing.T) {
	testCases := []string{
		`style="color:red" height=100 width="200px"`,
		`  id=main   class="container fluid"  disabled=true `,
		`data-x="" title='hello world' tabindex=0`,
		`width=100 width=200 height=50`,
	}
	expected := []map[string]string{
		{"style": "color:red", "height": "100", "width": "200px"},
		{"id": "main", "class": "container fluid", "disabled": "true"},
		{"data-x": "", "title": "hello world", "tabindex": "0"},
		{"width": "200", "height": "50"},
	}

	for i, test := range testCases {
		dummy := make(map[string]string)
		assignAttrs(&dummy, test)
		if maps.Equal(dummy, expected[i]) {
			t.Errorf("#%d: Expected %v | Got %v", i, expected[i], dummy)
		}
	}
}
