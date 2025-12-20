package lexer

import (
	"strings"
)

type Type int

const (
	Void Type = iota // unknown state
	Open
	Close
	Inner   // anything that's not tags
	SClose  // self-closed
	DocType // special syntax for <!DOCTYPE ..>
	// TODO: Comment
)

type Token struct {
	Type    Type
	Content string
	Endpos  int // in the raw string (last idx + 1)
	// line?
}

// TODO: support comment <!--something-->
func GetNextToken(raw string, pos int) Token {
	var res Token
	idx := skipWhiteSpace(raw, pos)
	for ; idx < len(raw); idx++ {
		char := raw[idx]
		switch char {
		case '<':
			switch res.Type {
			case Inner:
				res.Type = Inner
				res.Endpos = idx
				return res
			case Void:
				if idx+1 >= len(raw) {
					res.Content += string(char)
					res.Endpos = idx + 1
					return res
				}

				if raw[idx+1] == '/' {
					res.Type = Close
					idx++ // skip '/'
				} else if dtLen := len("!DOCTYPE"); idx+dtLen < len(raw) &&
					raw[idx+1:idx+dtLen+1] == "!DOCTYPE" {
					res.Type = DocType
					idx += dtLen
				} else {
					res.Type = Open
				}
			case Open:
				res.Content = string('<') + res.Content
				res.Type = Inner // reinterpret itself to inner
				res.Endpos = idx
				return res
			default:
				res.Content += string('<')
			}
		case '>':
			if res.Type == Close || res.Type == Open || res.Type == DocType {
				if res.Type == DocType {
					res.Content = strings.TrimSpace(res.Content)
				}
				res.Endpos = idx + 1
				return res
			} else {
				res.Content += string('>')
			}
		case '/':
			if res.Type == Open {
				if idx+1 == len(raw) {
					res.Content += string(char)
					res.Endpos = idx + 1
					return res
				} else if raw[idx+1] == '>' {
					res.Type = SClose
					res.Endpos = idx + 2
					return res
				}
			} else {
				res.Content += string('>')
			}
		default:
			if res.Type == Void {
				res.Type = Inner
			}
			res.Content += string(char)
		}
	}

	res.Endpos = idx + 1
	return res
}

func skipWhiteSpace(s string, pos int) int {
	for idx, char := range s[pos:] {
		if char != ' ' && char != '\n' && char != '\t' {
			return idx + pos
		}
	}
	return len(s)
}
