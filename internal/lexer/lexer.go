package lexer

import (
	"strings"
	"unicode"
)

type Type uint8

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
// NOTE: support comment <!--...-->
func GetNextToken(raw string, pos int) Token {
	var res Token
	idx := SkipWhiteSpace(raw, pos)
	for ; idx < len(raw); idx++ {
		char := raw[idx]
		switch res.Type {
		case Void:
			if char == '<' {
				if idx+1 < len(raw) && raw[idx+1] == '/' {
					res.Type = Close
					idx++ // skip /
				} else if dtLen := len("<!doctype"); idx+dtLen <= len(raw) &&
					strings.ToLower(raw[idx:idx+dtLen]) == "<!doctype" {
					res.Type = DocType
					idx += dtLen - 1 // skip <!doctype
				} else if cLen := len("<!--"); idx+cLen <= len(raw) &&
					raw[idx:idx+cLen] == "<!--" {
					res.Type = Comment
					idx += cLen - 1 // skip <!--
				} else {
					res.Type = Open
				}
			} else {
				res.Type = NoTag
				res.Content += string(char)
			}

		case Open:
			if char == '>' {
				res.Endpos = idx + 1
				return res
			}
			if idx+1 < len(raw) && raw[idx:idx+2] == "/>" {
				res.Type = SClose
				res.Endpos = idx + 2
				return res
			}

			// Reinterpret ourselves back to notag.
			// e.g. <p> 1 < 2 </p>
			// 				  ^
			// 				we are here
			if char == '<' {
				res.Type = NoTag
				res.Content = "<" + res.Content
				res.Endpos = idx
				return res
			}
			res.Content += string(char)

		case Close:
			if char == '>' {
				res.Endpos = idx + 1
				return res
			}
			res.Content += string(char)

		case DocType:
			if char == '>' {
				res.Content = strings.TrimSpace(res.Content)
				res.Endpos = idx + 1
				return res
			}
			res.Content += string(char)

		case Comment:
			if cLen := len("-->"); idx+cLen <= len(raw) && raw[idx:idx+cLen] == "-->" {
				res.Content = strings.TrimSpace(res.Content)
				res.Endpos = idx + cLen
				return res
			}
			res.Content += string(char)

		case NoTag:
			if char == '<' {
				res.Endpos = idx
				return res
			}
			res.Content += string(char)
		}
	}

	res.Endpos = idx + 1
	return res
}

// SkipWhiteSpace takes a string s and starting position pos and
// returns the position after pos that is not white space.
func SkipWhiteSpace(s string, pos int) int {
	for idx, char := range s[pos:] {
		if !unicode.IsSpace(char) {
			return idx + pos
		}
	}
	return len(s)
}
