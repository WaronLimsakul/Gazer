package css

import "strings"

// rule represent CSS rule: think of one block in CSS file
// e.g. p, h1 { color: green; margin: 10px; }
type rule struct {
	selectors []string          // e.g. p, h1, #class
	styles    map[string]string // map property->value of the style
}

// Parse parses raw CSS content string into a StyleSet
// requires: the syntax must be correct.
func Parse(raw string) (*StyleSet, error) {
	lexer := newLexer(raw)
	res := newStyleSet()

	state := Selector
	curRule := newRule()
	var tmpProp string
mainLoop: // first time in my life using this. Haha
	for {
		token := lexer.getNextToken()
		switch token.Type {
		case Void:
			continue // skip invalid token
		case Comment:
			continue // skip comment token
		case End:
			// TODO: not sure if this will break, have to test
			if state == Value {
				res.applyRule(curRule)
			}
			break mainLoop
		case Selector:
			if state == Value {
				res.applyRule(curRule)
				curRule.reset() // after finish processing one rule, reset the current rule
			}
			state = Selector
			content := strings.TrimSpace(token.Content)
			selectors := strings.Split(content, ",")
			for i, selector := range selectors {
				selectors[i] = strings.TrimSpace(selector)
			}
			curRule.selectors = selectors

		case Property:
			state = Property
			// property = case-insensitive
			content := strings.TrimSpace(strings.ToLower(token.Content))
			tmpProp = content
		case Value:
			state = Value
			content := strings.TrimSpace(token.Content)
			curRule.styles[tmpProp] = content
		}
	}

	return res, nil
}

// ParseStyle recieve a raw string of HTML inline "style" attribute and return a css.Style
func ParseStyle(raw string) (s Style) {
	raw = strings.TrimSpace(raw)
	decls := strings.Split(raw, ";")
	declsMap := make(map[string]string)
	for _, decl := range decls {
		kv := strings.Split(decl, ":")
		if len(kv) != 2 {
			continue
		}
		declsMap[strings.TrimSpace(strings.ToLower(kv[0]))] = strings.TrimSpace(kv[1])
	}
	s.registerDecls(declsMap)
	return s
}

func newRule() rule {
	return rule{make([]string, 0), make(map[string]string)}
}

func (r *rule) reset() {
	r.selectors = make([]string, 0)
	clear(r.styles)
}
