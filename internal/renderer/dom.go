package renderer

import (
	"log"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/parser"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

type DomRenderer struct {
	thm *material.Theme
	url string
	// All Texts' selectables elements based on its pointer.
	// NOTE: we can use pointer because it is one url per DomRenderer
	selectables    map[*parser.Node]*widget.Selectable
	linkClickables map[*parser.Node]*widget.Clickable
}

func newDomRenderer(thm *material.Theme, url string) *DomRenderer {
	return &DomRenderer{thm: thm, url: url,
		selectables:    make(map[*parser.Node]*widget.Selectable),
		linkClickables: make(map[*parser.Node]*widget.Clickable),
	}
}

// renderDOM takes a DOM root node and return [][]Element
// First layer (outer) is each horizontal line of rendering.
// Second layer (inner) is each element in that line from left to right.
func (dr *DomRenderer) render(root *parser.Node) [][]Element {
	res := make([][]Element, 0)

	// expect to be Root node
	if root == nil || root.Tag != parser.Root {
		return res
	}

	// expect root node to only have HTML tag
	// len(root.Children) != 1 ||
	if root.Children[0].Tag != parser.Html {
		return res
	}

	htmlNode := root.Children[0]
	for _, child := range htmlNode.Children {
		res = append(res, dr.renderNode(child)...)
	}

	return res
}

// renderNodes returns flex children needs for render a node and its children.
func (dr *DomRenderer) renderNode(node *parser.Node) [][]Element {
	res := make([][]Element, 0)
	switch node.Tag {
	case parser.Head:
		return res // TODO
	case parser.Body:
		for _, child := range node.Children {
			res = append(res, dr.renderNode(child)...)
		}
	case parser.Title:
		return res // TODO
	case parser.Meta:
		return res // TODO
	case parser.Div:
		for _, child := range node.Children {
			res = append(res, dr.renderNode(child)...)
		}
	case parser.Br:
		res = append(res, []Element{layout.Spacer{Height: unit.Dp(10)}})
	case parser.Hr:
		res = append(res, []Element{ui.HorizontalLine{Thm: dr.thm, Width: WINDOW_WIDTH, Height: unit.Dp(1)}})
	case parser.Img:
		img, err := ui.NewImg(node.Attrs["src"])
		if err != nil {
			log.Println("ui.NewImg: ", err)
			break
		}
		res = append(res, []Element{img})
	}

	if parser.TextElements[node.Tag] {
		res = append(res, dr.renderText(node)...)
	}

	return res
}

// renderText returns [][]Element needs for rendering a text node and its children.
func (dr *DomRenderer) renderText(node *parser.Node) [][]Element {
	// base case
	if node.Tag == parser.Text {
		selectable, ok := dr.selectables[node]
		if !ok {
			selectable = new(widget.Selectable)
			dr.selectables[node] = selectable
		}

		return [][]Element{{ui.Text(dr.thm, selectable, node.Inner)}}
	}

	// recursive case [phase 1]: aggregate the children elements
	res := make([][]Element, 0)
	// if inline-text and prev is also inline-text, put it in latest one don't append
	for i, child := range node.Children {
		if parser.InlineElements[child.Tag] && len(res) > 0 && i > 0 &&
			parser.InlineElements[node.Children[i-1].Tag] {
			childElems := dr.renderNode(child)
			if len(childElems) > 0 {
				res[len(res)-1] = append(res[len(res)-1], childElems[0]...)
			}
			if len(childElems) > 1 {
				res = append(res, childElems[1:]...)
			}

		} else {
			res = append(res, dr.renderNode(child)...)
		}

	}

	// recursive case [phase 2]: decorate the children that are Label

	// recursive case [phase 2.1]: take care of a special decorator
	// tag "A" decorator has different signature
	if node.Tag == parser.A {
		clickable, ok := dr.linkClickables[node]
		if !ok {
			clickable = new(widget.Clickable)
			dr.linkClickables[node] = clickable
		}

		for _, line := range res {
			for i, el := range line {
				if label, ok := el.(ui.Label); ok {
					line[i] = ui.A(clickable, label)
				}
			}
		}
		return res
	}

	// tag "ul" and "li" don't wanna apply to everyone in a row

	if node.Tag == parser.Ul {
		for _, line := range res {
			if len(line) > 0 {
				if label, ok := line[0].(ui.Label); ok {
					line[0] = ui.Ul(label)
				}
			}
		}
		return res
	}

	if node.Tag == parser.Ol {
		count := 1
		for _, line := range res {
			if len(line) > 0 {
				if label, ok := line[0].(ui.Label); ok {
					line[0] = ui.Ol(label, &count)
				}
			}
		}
		return res
	}

	// recursive case [phase 2.2]: normal text decorator
	dec := ui.P
	switch node.Tag {
	case parser.H1:
		dec = ui.H1
	case parser.H2:
		dec = ui.H2
	case parser.H3:
		dec = ui.H3
	case parser.H4:
		dec = ui.H4
	case parser.H5:
		dec = ui.H5
	case parser.P:
		dec = ui.P
	case parser.I:
		dec = ui.I
	case parser.B:
		dec = ui.B
	case parser.Li:
		dec = ui.Li
	}

	for _, line := range res {
		for i, el := range line {
			if label, ok := el.(ui.Label); ok {
				line[i] = dec(dr.thm, label)
			}
		}
	}

	return res
}

// update updates all elements ui in domrender
func (dr *DomRenderer) update(gtx C) {
	for _, clickable := range dr.linkClickables {
		clickable.Update(gtx)
		if clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}
	}
}

// linkClicked return whether the link in the page is clicked and
// if so, what does it linked to.
func (dr *DomRenderer) linkClicked(gtx C) (bool, string) {
	for node, clickable := range dr.linkClickables {
		if clickable.Clicked(gtx) {
			return true, node.Attrs["href"]
		}
	}
	return false, ""
}
