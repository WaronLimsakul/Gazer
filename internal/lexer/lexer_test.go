package lexer

import (
	"strings"
	"testing"
)

func endPos(s string, target string) int {
	return strings.Index(s, target) + len(target)
}

func testTokenizeSeq(t *testing.T, testcase string, expected []Token) {
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
}

func TestGetNextTokenBasic(t *testing.T) {
	t1 := `
	<!DOCTYPE html>
		<body>
			<p>hello, world</p>
		</body>
	`
	expectedTokens := []Token{
		{Type: DocType, Content: "html", Endpos: endPos(t1, "<!DOCTYPE html>")},
		{Type: Open, Content: "body", Endpos: endPos(t1, "<body>")},
		{Type: Open, Content: "p", Endpos: endPos(t1, "<p>")},
		{Type: Inner, Content: "hello, world", Endpos: endPos(t1, "hello, world")},
		{Type: Close, Content: "p", Endpos: endPos(t1, "</p>")},
		{Type: Close, Content: "body", Endpos: endPos(t1, "</body>")},
	}

	testTokenizeSeq(t, t1, expectedTokens)

}

func TestGetNextTokenBracketInner(t *testing.T) {
	tc := `
	<!DOCTYPE html>
		<body>
			<p>hello, <world</p>
		</body>
	`
	expectedTokens := []Token{
		{Type: DocType, Content: "html", Endpos: endPos(tc, "<!DOCTYPE html>")},
		{Type: Open, Content: "body", Endpos: endPos(tc, "<body>")},
		{Type: Open, Content: "p", Endpos: endPos(tc, "<p>")},
		{Type: Inner, Content: "hello, ", Endpos: endPos(tc, "hello, ")},
		{Type: Inner, Content: "<world", Endpos: endPos(tc, "<world")},
		{Type: Close, Content: "p", Endpos: endPos(tc, "</p>")},
		{Type: Close, Content: "body", Endpos: endPos(tc, "</body>")},
	}

	testTokenizeSeq(t, tc, expectedTokens)
}
