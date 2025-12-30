package engine

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

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

	// DOM root
	Root *parser.Node
}

var client = &http.Client{Timeout: 3 * time.Second}

func Start(state *State, window *app.Window) {
	for noti := range state.Notifier {
		switch noti {
		case Search:
			url := state.Url
			if len(url) == 0 {
				log.Println("Empty URL")
				continue
			}

			// handle prefix: we want https:// or http://
			if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
				url = "https://" + url
			}

			res, err := fetch(url)
			if err != nil {
				log.Println("fetch:", err)
				continue
			}

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Println("io.ReadAll:", err)
				continue
			}

			log.Println("fetch:\n", string(resBody))

			root, err := parser.Parse(string(resBody))
			if err != nil {
				log.Println(err)
				continue
			}
			state.Root = root
			log.Println("parse:\n", *root)

			window.Invalidate()
		default:
			continue
		}
	}
}

func NewState() *State {
	s := State{}
	s.Notifier = make(chan Notification)
	return &s
}

func fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}
	req.Header.Set("User-Agent", "Gazer")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}
	return res, nil
}
