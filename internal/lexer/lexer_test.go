package lexer

import (
	"strings"
	"testing"
)

type TestGetNextTokenCase struct {
	name     string
	input    string
	expected []Token
}

func TestGetNextToken(t *testing.T) {
	normal := `<!DOCTYPE html>
			<body>
			<p style="hello">hello, world</p>
			<img src="img_src.jpg" width="500" height="600"/>
			</body>`
	brackInInner := `<!DOCTYPE html>
		<body>
			<p>hello, <world</p>
		</body>`
	testCases := []TestGetNextTokenCase{
		{
			name:  "normal",
			input: normal,
			expected: []Token{
				{Type: DocType, Content: "html", Endpos: endPos(normal, "<!DOCTYPE html>")},
				{Type: Open, Content: "body", Endpos: endPos(normal, "<body>")},
				{Type: Open, Content: `p style="hello"`, Endpos: endPos(normal, `<p style="hello">`)},
				{Type: NoTag, Content: "hello, world", Endpos: endPos(normal, "hello, world")},
				{Type: Close, Content: "p", Endpos: endPos(normal, "</p>")},
				{Type: SClose,
					Content: `img src="img_src.jpg" width="500" height="600"`,
					Endpos:  endPos(normal, `<img src="img_src.jpg" width="500" height="600"/>`)},
				{Type: Close, Content: "body", Endpos: endPos(normal, "</body>")},
			},
		},
		{
			name:  "brack in inner content",
			input: brackInInner,
			expected: []Token{
				{Type: DocType, Content: "html", Endpos: endPos(brackInInner, "<!DOCTYPE html>")},
				{Type: Open, Content: "body", Endpos: endPos(brackInInner, "<body>")},
				{Type: Open, Content: "p", Endpos: endPos(brackInInner, "<p>")},
				{Type: NoTag, Content: "hello, ", Endpos: endPos(brackInInner, "hello, ")},
				{Type: NoTag, Content: "<world", Endpos: endPos(brackInInner, "<world")},
				{Type: Close, Content: "p", Endpos: endPos(brackInInner, "</p>")},
				{Type: Close, Content: "body", Endpos: endPos(brackInInner, "</body>")},
			},
		},
	}

	for _, test := range testCases {
		testTokenizeSeq(test.name, t, test.input, test.expected)
	}

}

func endPos(s string, target string) int {
	return strings.Index(s, target) + len(target)
}

func testTokenizeSeq(name string, t *testing.T, testcase string, expected []Token) {
	t.Run(name, func(t *testing.T) {
		for cur, reps := 0, 0; cur < len(testcase); reps++ {
			token := GetNextToken(testcase, cur)
			if reps < len(expected) {
				if token != expected[reps] {
					t.Errorf("#%d: Expect %+v | Got %+v", reps, expected[reps], token)
				}
			} else if token.Type != Void {
				t.Errorf("Extra token: %+v", token)
			}
			cur = token.Endpos
		}
	})
}
