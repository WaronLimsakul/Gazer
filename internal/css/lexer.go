package css

import (
	"strings"

	"github.com/WaronLimsakul/Gazer/internal/lexer"
)

type tokenType uint8

const (
	Selector tokenType = iota
	Property
	Value
	End  // end of the string
	Void // for invalid token, parser can just skip if found one
)

type Token struct {
	Type    tokenType
	Content string
}

type Lexer struct {
	raw   string
	pos   int
	state tokenType // lexer state is just the current token type
}

// newLexer create a new CSS Lexer with the cursor pointing
// at the beginning of the raw CSS string
func newLexer(raw string) *Lexer {
	return &Lexer{raw: raw, pos: 0, state: Selector}
}

// getNextToken returns the next CSS Token from the current position
// of the Lexer's cursor. If the sequence is ended, return Token type End.
func (sl *Lexer) getNextToken() Token {
	sl.pos = lexer.SkipWhiteSpace(sl.raw, sl.pos)
	if sl.pos >= len(sl.raw) {
		return Token{Type: End}
	}

	content := ""
	for i := sl.pos; i < len(sl.raw); i++ {
		ch := sl.raw[i]
		switch sl.state {
		case Selector:
			if ch == '{' {
				sl.state = Property
				sl.pos = i + 1
				return Token{Type: Selector, Content: strings.TrimSpace(content)}
			}
		case Property:
			switch ch {
			case ':':
				sl.state = Value
				sl.pos = i + 1
				// TODO: not sure if case-insensitive
				return Token{Type: Property, Content: strings.TrimSpace(content)}
			case '}':
				sl.state = Selector
				sl.pos = i + 1
				return Token{Type: Void}
			}
		case Value:
			switch ch {
			case ';':
				sl.state = Property
				sl.pos = i + 1
				return Token{Type: Value, Content: strings.TrimSpace(content)}
			case '}':
				sl.state = Selector
				sl.pos = i + 1
				return Token{Type: Value, Content: strings.TrimSpace(content)}
			}
		}
		content += string(ch)
	}
	return Token{Type: End}
}

// Think it's more natural for CSS to let the lexer hold the state and keep getting next token
