package renderer

import (
	"fmt"
	"strings"
	"unicode"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/WaronLimsakul/Gazer/internal/css"
	"github.com/WaronLimsakul/Gazer/internal/parser"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

type Node = parser.Node
type StyleSet = css.StyleSet

// Main renderering of a website. One of these per tab.
type DomRenderer struct {
	thm *material.Theme
	tab *ui.Tab
	// Cache the matrix of elements with root node pointer.
	// Can cache it because engine also cache by pointer
	// (same url + same tab = same root ptr).
	cache map[*Node]*[][]Element
	// All Texts' selectables elements based on its pointer.
	// These pointers will not be cleaned because the map still refer to it.
	selectables      map[*Node]*widget.Selectable
	linkClickables   map[*Node]*widget.Clickable
	buttonClickables map[*Node]*widget.Clickable
	inputEditors     map[*Node]*widget.Editor
}

func newDomRenderer(thm *material.Theme, tab *ui.Tab) *DomRenderer {
	return &DomRenderer{thm: thm, tab: tab, cache: make(map[*Node]*[][]Element),
		selectables:      make(map[*Node]*widget.Selectable),
		linkClickables:   make(map[*Node]*widget.Clickable),
		buttonClickables: make(map[*Node]*widget.Clickable),
		inputEditors:     make(map[*Node]*widget.Editor),
	}
}

// renderDOM takes a DOM root node and return [][]Element
// First layer (outer) is each horizontal line of rendering.
// Second layer (inner) is each element in that line from left to right.
// TODO: doc
func (dr *DomRenderer) render(root *Node, styles *StyleSet) [][]Element {
	res := make([][]Element, 0)
	// expect to be Root node
	if root == nil || root.Tag != parser.Root {
		return res
	}

	if cachedRes, ok := dr.cache[root]; ok {
		return *cachedRes
	}

	// expect root node to only have HTML tag
	// len(root.Children) != 1 ||
	if root.Children[0].Tag != parser.Html {
		return res
	}

	htmlNode := root.Children[0]
	for _, child := range htmlNode.Children {
		res = append(res, dr.renderNode(child, styles, newRenderingContext())...)
	}

	dr.cache[root] = &res
	return res
}

// renderNode returns flex children needs for render a node and its children.
// TODO: doc
func (dr *DomRenderer) renderNode(node *Node, styles *StyleSet, rctx RenderingContext) [][]Element {
	rctx.ancestors = append(rctx.ancestors, node.Tag) // show the kid who is there pop

	res := make([][]Element, 0)
	switch node.Tag {
	case parser.Body:
		res = dr.gatherElements(node, styles, rctx)
	case parser.Br:
		res = append(res, []Element{layout.Spacer{Height: unit.Dp(10)}})
	case parser.Hr:
		res = append(res, []Element{ui.HorizontalLine{Thm: dr.thm, Width: WINDOW_WIDTH, Height: unit.Dp(1)}})
	case parser.Img:
		img, err := dr.renderImg(node)
		if err != nil {
			break
		}
		res = append(res, []Element{img})
	case parser.Input:
		res = append(res, []Element{dr.renderInput(node)})
	}

	if parser.ContainerElements[node.Tag] {
		res = append(res, []Element{dr.renderContainer(node, styles, rctx)})
	}
	if parser.TextElements[node.Tag] {
		res = append(res, dr.renderText(node, styles, rctx)...)
	}

	// pop from ancestors stack, done with this node
	rctx.ancestors = rctx.ancestors[:len(rctx.ancestors)-1]
	return res
}

// TODO: support other than Div
func (dr *DomRenderer) renderContainer(node *Node, styles *StyleSet, rctx RenderingContext) Element {
	// TODO: this entire style handling process should be centralized
	var localStyle css.Style
	localStylePtr := getNodeStyleFromStyleSet(styles, node)
	if localStylePtr != nil {
		localStyle = *localStylePtr
	}

	var inlineStyle css.Style
	inlineStyleStr, ok := node.Attrs["style"]
	if ok {
		inlineStyle = css.ParseStyle(inlineStyleStr)
	}

	curStyle := css.AddStyle(inlineStyle, css.AddStyle(localStyle, rctx.base))

	childrenRctx := rctx
	childrenRctx.base = containerInheritStyle(curStyle)
	children := dr.gatherElements(node, styles, childrenRctx)

	return ui.NewDiv(dr.thm, curStyle, children)
}

// renderText returns [][]Element needs for rendering a text node and its children.
// requires: node must be of the text type (check by using parser.TextElements)
// TODO: doc
// TODO: should we pass rctx by pointer or value?
func (dr *DomRenderer) renderText(node *Node, styles *StyleSet, rctx RenderingContext) [][]Element {
	// base case
	if node.Tag == parser.Text {
		selectable, ok := dr.selectables[node]
		if !ok {
			selectable = new(widget.Selectable)
			dr.selectables[node] = selectable
		}

		return [][]Element{{ui.NewLabel(dr.thm, rctx.getLabelStyle(), selectable, node.Inner)}}
	}

	// recursive case: decorate the label style
	// phase 1: update local style with tag-specific style
	switch node.Tag {
	case parser.A:
		clickable, ok := dr.linkClickables[node]
		if !ok {
			clickable = new(widget.Clickable)
			dr.linkClickables[node] = clickable
		}
		rctx.updateLabelStyle(ui.A(clickable, rctx.getLabelStyle()))
	case parser.Button:
		// TODO: v8 just wrap all text around and treat it like one big button
		clickable, ok := dr.buttonClickables[node]
		if !ok {
			clickable = new(widget.Clickable)
			dr.buttonClickables[node] = clickable
		}
		rctx.updateLabelStyle(ui.Button(dr.thm, clickable, rctx.getLabelStyle()))
	case parser.Ul:
		rctx.updateLabelStyle(ui.Ul(rctx.getLabelStyle()))
	case parser.Ol:
		rctx.updateLabelStyle(ui.Ol(rctx.getLabelStyle()))
	case parser.H1:
		rctx.updateLabelStyle(ui.H1(dr.thm, rctx.getLabelStyle()))
	case parser.H2:
		rctx.updateLabelStyle(ui.H2(dr.thm, rctx.getLabelStyle()))
	case parser.H3:
		rctx.updateLabelStyle(ui.H3(dr.thm, rctx.getLabelStyle()))
	case parser.H4:
		rctx.updateLabelStyle(ui.H4(dr.thm, rctx.getLabelStyle()))
	case parser.H5:
		rctx.updateLabelStyle(ui.H5(dr.thm, rctx.getLabelStyle()))
	case parser.P:
		rctx.updateLabelStyle(ui.P(dr.thm, rctx.getLabelStyle()))
	case parser.I:
		rctx.updateLabelStyle(ui.I(dr.thm, rctx.getLabelStyle()))
	case parser.B:
		rctx.updateLabelStyle(ui.B(dr.thm, rctx.getLabelStyle()))
	case parser.Li:
		rctx.updateLabelStyle(ui.Li(dr.thm, rctx.getLabelStyle(), rctx.ancestors))
	}

	// phase 2: update local style with CSS
	var externalStyle css.Style
	externalStylePtr := getNodeStyleFromStyleSet(styles, node)
	if externalStylePtr != nil {
		externalStyle = *externalStylePtr
	}

	var inlineStyle css.Style
	styleStr, ok := node.Attrs["style"]
	if ok {
		inlineStyle = css.ParseStyle(styleStr)
	}

	// priority: inline > external > inherit
	rctx.base = css.AddStyle(inlineStyle, css.AddStyle(externalStyle, rctx.base))

	// after modify the rctx, pass it up and gather elements
	return dr.gatherElements(node, styles, rctx)
}

// renderImg receive Img tag node and return Img ui element.
// Img is void element, don't have to gather more
func (dr *DomRenderer) renderImg(node *Node) (Element, error) {
	if node == nil {
		return layout.Spacer{}, fmt.Errorf("nil node")
	}
	if node.Tag != parser.Img {
		return layout.Spacer{}, fmt.Errorf("invalid tag: %v", node.Tag.String())
	}
	img, err := ui.NewImg(node.Attrs["src"])
	if err != nil {
		return layout.Spacer{}, fmt.Errorf("ui.NewImg: %v", err)
	}
	return img, nil
}

// renderInput receive Input tag node and return Input ui element.
// Input is void element, don't have to gather more.
// requires: node must not be nil and have input tag
func (dr *DomRenderer) renderInput(node *Node) Element {
	editor, ok := dr.inputEditors[node]
	if !ok {
		editor = new(widget.Editor)
		dr.inputEditors[node] = editor
	}

	inputTypeStr := strings.ToLower(node.Attrs["type"])
	inputType, ok := ui.InputTypes[inputTypeStr]
	if !ok {
		inputType = ui.TextInput
	}

	hint := node.Attrs["placeholder"]
	return ui.NewInput(dr.thm, inputType, editor, hint)
}

// handleHead set the tabview data by processing <head> node in the DOM tree (except css-related)
func (dr *DomRenderer) handleHead(root *Node) {
	head := dr.findHead(root)
	if head == nil {
		dr.tab.Title = ""
		return
	}

	var titleSet bool
	for _, node := range head.Children {
		// handle title tag
		switch node.Tag {
		case parser.Title:
			for _, titleChild := range node.Children {
				if titleChild.Tag == parser.Text {
					titleSet = true
					dr.tab.Title = titleChild.Inner
				}
			}
			// TODO: support some other links
		}
	}

	if !titleSet {
		dr.tab.Title = ""
	}
}

func (dr DomRenderer) findHead(node *Node) *Node {
	if node == nil {
		return nil
	}
	if node.Tag == parser.Head {
		return node
	}
	for _, child := range node.Children {
		found := dr.findHead(child)
		if found != nil {
			return found
		}
	}
	return nil
}

// gaterElements recieves a node and gather all elements of the node's children
// according the tag rule (inline, block)
func (dr DomRenderer) gatherElements(node *Node, styles *StyleSet, rctx RenderingContext) [][]Element {
	res := make([][]Element, 0)
	// if inline-text and prev is also inline-text, put it in latest one don't append
	for i, child := range node.Children {
		if parser.InlineElements[child.Tag] && len(res) > 0 && i > 0 &&
			parser.InlineElements[node.Children[i-1].Tag] {
			childElems := dr.renderNode(child, styles, rctx)
			if len(childElems) > 0 {
				res[len(res)-1] = append(res[len(res)-1], childElems[0]...)
			}
			if len(childElems) > 1 {
				res = append(res, childElems[1:]...)
			}

		} else {
			res = append(res, dr.renderNode(child, styles, rctx)...)
		}

	}
	return res
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

// getNodeStyleFromStyleSet return a css.Style of the node according to
// the CSS that styleset represent.
func getNodeStyleFromStyleSet(ss *StyleSet, node *parser.Node) *css.Style {
	if node == nil || ss == nil {
		return nil
	}

	var res *css.Style
	tagStyle, ok := ss.TagStyles[node.Tag]
	if ok {
		res = css.AddStylePtr(tagStyle, res)
	}

	classesStr, ok := node.Attrs["class"]
	if ok {
		classesStr := strings.TrimSpace(classesStr)
		classes := strings.FieldsFunc(classesStr, unicode.IsSpace)
		for _, class := range classes {
			classStyle, ok := ss.ClassStyles[class]
			if ok {
				res = css.AddStylePtr(classStyle, res)
			}
		}
	}

	id, ok := node.Attrs["id"]
	if ok {
		idStyle, ok := ss.IdStyles[id]
		if ok {
			res = css.AddStylePtr(idStyle, res)
		}
	}

	return res
}

// containerInheritStyle returns a container style with only inheritable fields
func containerInheritStyle(style css.Style) css.Style {
	var res css.Style
	res.Color = style.Color
	res.FontSize = style.FontSize
	res.FontStyle = style.FontStyle
	res.FontWeight = style.FontWeight
	return res
}
