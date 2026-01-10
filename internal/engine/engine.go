package engine

import (
	"fmt"
	"io"
	"log"
	"net/http"
	urlPackage "net/url"
	"strings"
	"time"

	"gioui.org/app"
	"github.com/WaronLimsakul/Gazer/internal/css"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

type NotificationType uint8

const (
	Search NotificationType = iota
	AddTab
	ChangeTab
)

type Notification struct {
	Type   NotificationType
	TabIdx int
	Url    string
}

type State struct {
	Tabs []*Tab
	// a channel for client to notify the engine with the event
	Notifier chan Notification

	// a channel for client to check the loading progress
	LoadProgress chan float32
}

type Tab struct {
	Url       string
	Root      *parser.Node
	Styles    *css.StyleSet
	IsLoading bool
}

var client = &http.Client{Timeout: 3 * time.Second}

// Start starts the engine to watch for notification and serve the request>
func Start(state *State, window *app.Window) {
	cache := make(map[*Tab]map[string]*parser.Node)

	for noti := range state.Notifier {
		// TODO: might spin up go routine for each job
		switch noti.Type {
		case Search:
			tab := state.Tabs[noti.TabIdx]

			tabCache, ok := cache[tab]
			if !ok {
				tabCache = make(map[string]*parser.Node)
				cache[tab] = tabCache
			}

			tab.Url = noti.Url
			cachedRoot, ok := tabCache[tab.Url]
			if ok {
				tab.Root = cachedRoot
				window.Invalidate()
				continue
			}

			tab.IsLoading = true
			go reportProgress(tab, window, state.LoadProgress)
			root, err := search(tab.Url)
			styles := getStyles(root)
			tab.IsLoading = false
			if err != nil {
				fmt.Println("search:", err)
				continue
			}

			tabCache[tab.Url] = root
			tab.Root = root
			tab.Styles = styles
			window.Invalidate()
		case AddTab:
			state.Tabs = append(state.Tabs, &Tab{})
			window.Invalidate()
		default:
			continue
		}
	}
}

func NewState() *State {
	s := State{}
	s.Notifier = make(chan Notification)
	s.LoadProgress = make(chan float32)
	s.Tabs = []*Tab{{}}
	return &s
}

// search fetches the url and parse the DOM tree
// then return the root of DOM tree and error if exists
func search(url string) (*parser.Node, error) {
	if len(url) == 0 {
		return nil, fmt.Errorf("Empty URL")
	}

	// handle prefix: we want https:// or http://
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	// check valid url
	_, err := urlPackage.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %v", err)
	}

	res, err := fetch(url)
	if err != nil {
		return nil, fmt.Errorf("fetch: %v", err)
	}
	defer res.Body.Close()

	if !strings.HasPrefix(res.Header["Content-Type"][0], "text/html") {
		return nil, fmt.Errorf("Not HTML")
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

// getStyles get the CSS StyleSet from the DOM root. it returns nil if not found.
func getStyles(root *parser.Node) *css.StyleSet {
	head := findHead(root)
	if head == nil {
		return nil
	}

	var internal *css.StyleSet // styleset from <style></style>
	var external *css.StyleSet // styleset from <link ref="stylesheet" href="..">
	for _, node := range head.Children {
		// handle title tag
		switch node.Tag {
		case parser.Style:
			var contentBuilder strings.Builder
			for _, txt := range node.Children {
				if txt.Tag == parser.Text {
					contentBuilder.WriteString(txt.Inner)
				}
			}
			styles, err := css.Parse(contentBuilder.String())
			if err != nil {
				break
			}
			internal = styles
		case parser.Link:
			if rel, ok := node.Attrs["rel"]; ok && rel == "stylesheet" {
				href, ok := node.Attrs["href"]
				if !ok {
					break
				}
				_, err := urlPackage.Parse(href)
				if err != nil {
					break
				}
				res, err := client.Get(href)
				if err != nil {
					break
				}
				defer res.Body.Close()
				content, err := io.ReadAll(res.Body)
				if err != nil {
					break
				}
				styles, err := css.Parse(string(content))
				external = styles
			}
		}
	}

	if internal != nil && external != nil {
		return css.AddStyleSet(internal, external)
	} else if internal != nil {
		return internal
	} else {
		return external
	}
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

// reportProgress keep reporting synthetic progress to the channel in the state
// and also keep invalidating the window intending to rerender the progress ui.
func reportProgress(t *Tab, w *app.Window, progressChan chan float32) {
	var progress float32 = 0.0
	const rate = 0.1

	for t.IsLoading {
		time.Sleep(25 * time.Millisecond)
		progress += (1 - progress) * rate
		progressChan <- progress
		w.Invalidate()
	}

	progressChan <- 1.0
	w.Invalidate()
}

func findHead(root *parser.Node) *parser.Node {
	if root.Tag != parser.Root {
		return nil
	} else if len(root.Children) == 0 {
		return nil
	} else if root.Children[0].Tag != parser.Html {
		return nil
	}

	html := root.Children[0]
	for _, child := range html.Children {
		if child.Tag == parser.Head {
			return child
		}
	}
	return nil
}
