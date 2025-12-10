package parser

type Tag int

const (
	Html Tag = iota
	Head
	Title
	H1
	H2
	H3
	H4
	H5
	P
	Br

	// TODO:
	// A
	// Img
	// Ul
	// Ol
	// Li
	// B or Strong
	// I or Em
	// Hr
	// Div
	// Span
)

type Node struct {
	Tag      string
	Inner    string
	Attrs    map[string]any
	Children []*Node
}
