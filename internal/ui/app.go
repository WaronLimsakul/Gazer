package ui

type Element interface {
	Layout(gtx C) D
}

type App struct {
}
