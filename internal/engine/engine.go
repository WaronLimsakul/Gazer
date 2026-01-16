package engine

import (
	"fmt"
	"io"
	"log"
	"net/http"
	urlPkg "net/url"
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
	CloseTab
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
	// cache for each tab's HTML parser, 1 tab = 1 cache
	caches := make(map[*Tab]map[string]*parser.Node)

	for noti := range state.Notifier {
		// TODO: might spin up go routine for each job
		switch noti.Type {
		case Search:
			tab := state.Tabs[noti.TabIdx]

			tabCache, ok := caches[tab]
			if !ok {
				tabCache = make(map[string]*parser.Node)
				caches[tab] = tabCache
			}

			preparedUrl, err := prepareUrl(noti.Url)
			if err != nil {
				fmt.Println("prepareUrl: ", err)
				continue
			}
			tab.Url = preparedUrl.String()

			cachedRoot, ok := tabCache[tab.Url]
			if ok {
				tab.Root = cachedRoot
				window.Invalidate()
				continue
			}

			tab.IsLoading = true
			go reportProgress(tab, window, state.LoadProgress)

			root, err := search(tab.Url)
			styles := getStyles(root, preparedUrl)
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
		case CloseTab:
			state.Tabs = append(state.Tabs[:noti.TabIdx], state.Tabs[noti.TabIdx+1:]...)
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

// getStyles get the CSS StyleSet from the DOM root (and might need the base url of the root).
// it returns nil if not found.
func getStyles(root *parser.Node, baseUrl *urlPkg.URL) *css.StyleSet {
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
				continue
			}
			internal = styles
		case parser.Link:
			if rel, ok := node.Attrs["rel"]; ok && rel == "stylesheet" {
				href, ok := node.Attrs["href"]
				if !ok {
					continue
				}
				hrefUrl, err := baseUrl.Parse(href) // OP function
				if err != nil {
					log.Println("baseUrl.Parse: ", err)
					continue
				}
				res, err := client.Get(hrefUrl.String())
				if err != nil {
					log.Println("client.Get: ", err)
					continue
				}
				defer res.Body.Close()
				content, err := io.ReadAll(res.Body)
				if err != nil {
					log.Println("io.ReadAll: ", err)
					continue
				}
				log.Printf("fetch CSS [%s]: %s", href, content)

				styles, err := css.Parse(string(content))
				if err != nil {
					fmt.Println("css.Parse", err)
					continue
				}
				log.Println("parse CSS: ", *styles)

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
	if root == nil {
		return nil
	} else if root.Tag != parser.Root {
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

// prepareUrl takes a url string and return a new url.URL we can fetch from
func prepareUrl(rawUrl string) (*urlPkg.URL, error) {
	if len(rawUrl) == 0 {
		return nil, fmt.Errorf("Empty URL")
	}

	// handle prefix: we want https:// or http://
	if !strings.HasPrefix(rawUrl, "https://") && !strings.HasPrefix(rawUrl, "http://") {
		rawUrl = "https://" + rawUrl
	}

	// check valid url (don't use Parse, it's for both absolute and relative)
	url, err := urlPkg.ParseRequestURI(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %v", err)
	}

	return url, nil
}
