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

func TestParse(t *testing.T) {
	testCases := []string{
		`<!DOCTYPE html>
		<html>
			<h1 style="color:blue">This is a Heading</h1>
			<p>This is a paragraph.</p>
		</html>`,
	}

	// TODO: do something with the way to test this. It's kinda inconvenient
	root := newBaseNode()
	root.Children = append(root.Children, &Node{Tag: Html, Inner: "", Attrs: map[string]string{}, Children: []*Node{
		{Tag: H1, Inner: "This is a Heading", Attrs: map[string]string{"style": "color:blue"}, Children: []*Node{}, Parent: nil},
		{Tag: P, Inner: "This is a paragraph.", Attrs: map[string]string{}, Children: []*Node{}, Parent: nil},
	}, Parent: nil})

	expected := []*Node{root}

	for i, test := range testCases {
		actual, err := Parse(test)
		if err != nil {
			t.Errorf("#%d: unexpected error %v", i, err)
		} else {
			if !actual.equal(expected[i]) {
				t.Errorf("#%d: Expected %v \n Got %v", i, *expected[i], *actual)
			}
		}
	}
}
