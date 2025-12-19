package engine

import (
	"fmt"
	"io"
	"net/http"

	"gioui.org/app"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

type Notification int

const (
	Search Notification = iota
)

type State struct {
	Url      string
	Notifier chan Notification
	Content  string //  TODO: delete after finish html rendering

	// DOM root
	Root *parser.Node
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
				// TODO: make it formal log
				fmt.Println("Empty URL")
				continue
			}

			res, err := http.Get(url)
			if err != nil {
				// TODO: make it formal log
				fmt.Println("http.Get:", err)
				continue
			}

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				// TODO: make it formal log
				fmt.Println("io.ReadAll:", err)
				continue
			}

			state.Content = string(resBody)
			window.Invalidate()
		default:
			continue
		}
	}
}
