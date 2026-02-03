package ui

import "gioui.org/layout"

// just a component that wrap the page nav and search bar
type TopBar struct {
	searchBar *SearchBar
	pageNav   *PageNav
}

func NewTopBar(searchBar *SearchBar, pageNav *PageNav) *TopBar {
	return &TopBar{searchBar: searchBar, pageNav: pageNav}
}

func (tb TopBar) Layout(gtx C) D {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx, Rigid(tb.pageNav), Rigid(tb.searchBar))
}
