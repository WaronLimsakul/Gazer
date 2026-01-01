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

type Notification uint8

const (
	Search Notification = iota
)

type State struct {
	Url string
	// a channel to notify the engine with the event
	Notifier chan Notification

	// DOM root
	Root *parser.Node
}

var client = &http.Client{Timeout: 3 * time.Second}

func Start(state *State, window *app.Window) {
	for noti := range state.Notifier {
		switch noti {
		case Search:
			root, err := search(state.Url)
			if err != nil {
				fmt.Println("search:", err)
				continue
			}
			state.Root = root
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

// search fetches the url and parse the DOM tree
// and the root of DOM tree and error if exists
func search(url string) (*parser.Node, error) {
	if len(url) == 0 {
		return nil, fmt.Errorf("Empty URL")
	}

	// handle prefix: we want https:// or http://
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	res, err := fetch(url)
	if err != nil {
		return nil, fmt.Errorf("fetch: %v", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %v", err)
	}

	log.Println("fetch:\n", string(resBody))

	root, err := parser.Parse(string(resBody))
	if err != nil {
		return nil, fmt.Errorf("parse: %v", err)

	}
	log.Println("parse:\n", *root)
	return root, nil
}

// fetch uses the url to fetch the response
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
