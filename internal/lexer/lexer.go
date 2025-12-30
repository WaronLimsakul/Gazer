package lexer

import (
	"strings"
	"unicode"
)

type Type int

const (
	Void    Type = iota // initial state
	Open                // <openTag>
	Close               // </closeTag>
	NoTag               // content with no tag
	SClose              // <self-closed/>
	DocType             // only for <!DOCTYPE ..>
	Comment             // <!--something-->
)

// html token designed for parsing
type Token struct {
	Type    Type
	Content string
	Endpos  int // in the raw string (last idx + 1)
	// line?
}

// GetNextToken receive raw html string and starting position and
// return the next html token (e.g. "<hello>", "</word>", "foo") from the
// starting position.
// NOTE: support comment <!--something-->, <!foo>
// TODO NOW: can't just check <!...>, because what if <!-- <h1>Hello, world</h1> -->
func GetNextToken(raw string, pos int) Token {
	var res Token
	idx := skipWhiteSpace(raw, pos)
	for ; idx < len(raw); idx++ {
		char := raw[idx]
		switch char {
		case '<':
			switch res.Type {
			case NoTag:
				res.Type = NoTag
				res.Endpos = idx
				return res
			case Void:
				if idx+1 >= len(raw) {
					res.Content += string(char)
					res.Endpos = idx + 1
					return res
				}

				// handle close-tag
				if raw[idx+1] == '/' {
					res.Type = Close
					idx++ // skip '/'

					// for special <!doctype> tag
				} else if dtLen := len("!doctype"); idx+dtLen < len(raw) &&
					strings.ToLower(raw[idx+1:idx+dtLen+1]) == "!doctype" {
					res.Type = DocType
					idx += dtLen
				} else if raw[idx+1] == '!' {
					res.Type = Comment
					idx++ // skip '!'
				} else {
					res.Type = Open
				}

			// e.g. <p> 1 < 2 </p>
			// 				  ^
			// 				we are here
			case Open:
				res.Content = string('<') + res.Content
				res.Type = NoTag // reinterpret itself to "no tag" content
				res.Endpos = idx
				return res
			default:
				res.Content += string('<')
			}
		case '>':
			// Token that has '<' will ends with '>'
			if res.Type == Open || res.Type == Close || res.Type == Comment || res.Type == DocType {
				if res.Type == DocType {
					res.Content = strings.TrimSpace(res.Content)
				}
				if res.Type == Comment {
					res.Content = strings.Trim(res.Content, "--") // <!--hi--> = "hi"
					res.Content = strings.TrimSpace(res.Content)  // <! \n haha> = "haha"
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
				res.Type = NoTag
			}
			res.Content += string(char)
		}
	}

	res.Endpos = idx + 1
	return res
}

// skipWhiteSpace takes a string s and starting position pos and
// returns the position after pos that is not white space.
func skipWhiteSpace(s string, pos int) int {
	for idx, char := range s[pos:] {
		if !unicode.IsSpace(char) {
			return idx + pos
		}
	}
	return len(s)
}
