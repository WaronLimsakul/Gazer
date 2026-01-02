package parser

type Tag uint8

// Legal tag list.
// If change this part. Then change:
// 1. TagMap
// 2. Tag.String method
// 3. VoidElements if appropriate
// 4. TextElements if appropriate
// 5. ui.renderNode()
const (
	Root Tag = iota // Only for root node
	Html
	Head
	Body

	Title
	Meta

	Div

	H1
	H2
	H3
	H4
	H5
	P
	I
	B
	A

	Ul
	Ol
	Li

	Img

	Br

	Text // For no tag content or invalid tag

	// TODO: Img Ol Hr Span
)

var TagMap = map[string]Tag{
	"html":   Html,
	"head":   Head,
	"body":   Body,
	"title":  Title,
	"meta":   Meta,
	"div":    Div,
	"h1":     H1,
	"h2":     H2,
	"h3":     H3,
	"h4":     H4,
	"h5":     H5,
	"p":      P,
	"i":      I,
	"em":     I,
	"b":      B,
	"a":      A,
	"ul":     Ul,
	"ol":     Ol,
	"li":     Li,
	"strong": B,
	"br":     Br,
	"img":    Img,
}

func (t Tag) String() string {
	switch t {
	case Root:
		return "root"
	case Html:
		return "html"
	case Head:
		return "head"
	case Body:
		return "body"
	case Title:
		return "title"
	case Meta:
		return "meta"
	case Div:
		return "div"
	case H1:
		return "h1"
	case H2:
		return "h2"
	case H3:
		return "h3"
	case H4:
		return "h4"
	case H5:
		return "h5"
	case P:
		return "p"
	case I:
		return "i"
	case B:
		return "b"
	case A:
		return "a"
	case Ul:
		return "ul"
	case Ol:
		return "ol"
	case Li:
		return "li"
	case Br:
		return "br"
	case Text:
		return "text"
	case Img:
		return "img"
	default:
		return "unknown"
	}
}

// void elements = tag that always means self-close
var VoidElements = map[Tag]bool{
	Br:   true,
	Meta: true,
	Img:  true,
}

// text elements = tag that is rendered as a string
var TextElements = map[Tag]bool{
	H1:   true,
	H2:   true,
	H3:   true,
	H4:   true,
	H5:   true,
	P:    true,
	I:    true,
	B:    true,
	A:    true,
	Ul:   true,
	Ol:   true,
	Li:   true,
	Text: true,
}

// inline elements = element that will not break line when
// being child of another text element. E.g. <p>hello, <i>world</i></p> is one line
var InlineElements = map[Tag]bool{
	I:    true,
	B:    true,
	A:    true,
	Text: true,
	Img:  true,
}
