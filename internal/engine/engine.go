package engine

import (
	"io"
	"net/http"

	"gioui.org/app"
)

type Notification int

const (
	Search Notification = iota
)

type State struct {
	Url      string
	Notifier chan Notification
	Content  string
}

func NewState() *State {
	s := State{}
	s.Notifier = make(chan Notification)
	return &s
}

func Start(state *State, window *app.Window) {
	for noti := range state.Notifier {
		switch noti {
		case Search:
			url := state.Url
			if len(url) == 0 {
				// TODO
				continue
			}
			res, err := http.Get(url)
			if err != nil {
				// TODO
				continue
			}
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				// TODO
				continue
			}
			state.Content = string(resBody)
			window.Invalidate()
		default:
			continue
		}
	}
}
