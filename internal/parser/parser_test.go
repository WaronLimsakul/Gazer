package parser

import (
	"maps"
	"testing"
)

func TestGetTag(t *testing.T) {
	testCases := map[string]Tag{
		"h1":   H1,
		"p":    P,
		"ahah": Text,
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

type testParseCase struct {
	name     string
	input    string
	expected *Node
}

func TestParse(t *testing.T) {
	testCases := []testParseCase{
		{
			name: "normal",
			input: `<!DOCTYPE html>
		<html>
			<h1 style="color:blue">This is a Heading</h1>
			<p>This is a paragraph.</p>
		</html>`,
			expected: newTestTree(Root, "", nil,
				newTestTree(Html, "", nil,
					newTestTree(H1, "", map[string]string{"style": "color:blue"},
						newText("This is a Heading", nil)),
					newTestTree(P, "", nil,
						newText("This is a paragraph.", nil)))),
		},
		{
			name: "nested",
			input: `<!DOCTYPE html>
		<html>
			<head>
				<title>My Page</title>
			</head>
			<body>
				<h1>Welcome</h1>
				<p>Hello world</p>
			</body>
		</html>`,
			expected: newTestTree(Root, "", nil,
				newTestTree(Html, "", nil,
					newTestTree(Head, "", nil,
						newTestTree(Title, "", nil,
							newText("My Page", nil))),
					newTestTree(Body, "", nil,
						newTestTree(H1, "", nil,
							newText("Welcome", nil)),
						newTestTree(P, "", nil,
							newText("Hello world", nil))))),
		},
		{
			name: "many attributes",
			input: `<!DOCTYPE html>
		<html>
			<h1 class="header" id="main" style="color:red">Title</h1>
		</html>`,
			expected: newTestTree(Root, "", nil,
				newTestTree(Html, "", nil,
					newTestTree(H1, "", map[string]string{
						"class": "header",
						"id":    "main",
						"style": "color:red"},
						newText("Title", nil)))),
		},
		{
			name: "self-closing br",
			input: `<!DOCTYPE html>
		<html>
			<p>Line one<br>Line two</p>
		</html>`,
			expected: newTestTree(Root, "", nil,
				newTestTree(Html, "", nil,
					newTestTree(P, "", nil,
						newText("Line one", nil),
						newTestTree(Br, "", nil),
						newText("Line two", nil)))),
		},
		{
			name: "multiple br tags",
			input: `<!DOCTYPE html>
		<html>
			<p>First<br>Second<br>Third</p>
		</html>`,
			expected: newTestTree(Root, "", nil,
				newTestTree(Html, "", nil,
					newTestTree(P, "", nil,
						newText("First", nil),
						newTestTree(Br, "", nil),
						newText("Second", nil),
						newTestTree(Br, "", nil),
						newText("Third", nil)))),
		},
		{
			name: "empty tags",
			input: `<!DOCTYPE html>
		<html>
			<h1></h1>
			<p></p>
		</html>`,
			expected: newTestTree(Root, "", nil,
				newTestTree(Html, "", nil,
					newTestTree(H1, "", nil),
					newTestTree(P, "", nil))),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Parse(test.input)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			} else {
				if !actual.equal(test.expected) {
					t.Errorf("Expected %v \n Got %v", *test.expected, *actual)
				}
			}
		})
	}
}

func newTestTree(tag Tag, inner string, attrs map[string]string, children ...*Node) *Node {
	if attrs == nil {
		attrs = make(map[string]string)
	}
	node := &Node{
		Tag:      tag,
		Inner:    inner,
		Attrs:    attrs,
		Children: children,
	}
	for _, child := range children {
		child.Parent = node
	}
	return node
}
