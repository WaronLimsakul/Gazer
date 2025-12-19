package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func endPos(s string, target string) int {
	return strings.Index(s, target) + len(target)
}

func TestGetNextToken(t *testing.T) {
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

	for cur, reps := 0, 0; cur < len(t1); reps++ {
		token := GetNextToken(t1, cur)
		fmt.Println("reps:", reps, token)
		if reps < len(expectedTokens) {
			if token != expectedTokens[reps] {
				t.Errorf("#%d: Expect %+v | Got %+v", reps, expectedTokens[reps], token)
			}
		} else if token.Type != Void {
			t.Errorf("Extra token: %+v", token)
		}
		cur = token.Endpos
	}
}
