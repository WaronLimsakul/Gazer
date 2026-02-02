package engine

import (
	"fmt"
	"io"
	"log"
	"net/http"
	urlPkg "net/url"
	"os"
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
}

type Tab struct {
	Url          string // processed URL
	Root         *parser.Node
	Styles       *css.StyleSet
	IsLoading    bool
	LoadProgress chan float32 // a channel reporting the loading progress
}

var client = &http.Client{Timeout: 3 * time.Second}

// checked content type in HTTPS header
var supportedContentType = map[string]bool{
	"text/html": true,
	"text/css":  true,
}

// Start starts the engine to watch for notification and serve the request>
func Start(state *State, window *app.Window) {
	serverNotifiers := make(map[*Tab]chan Notification) // map notification to channel to server

	for noti := range state.Notifier {
		// 2 operations that manager has to deal: open and close tab
		switch noti.Type {
		case AddTab:
			state.Tabs = append(state.Tabs, newTab())
			window.Invalidate()
		case CloseTab:
			delete(serverNotifiers, state.Tabs[noti.TabIdx]) // delete the serverNotifier
			state.Tabs = append(state.Tabs[:noti.TabIdx], state.Tabs[noti.TabIdx+1:]...)
			window.Invalidate()
		default:
			tab := state.Tabs[noti.TabIdx]
			serverNotifier, ok := serverNotifiers[tab]
			if !ok {
				serverNotifier = make(chan Notification)
				serverNotifiers[tab] = serverNotifier
				go serveTab(tab, serverNotifier, window)
			}
			serverNotifier <- noti
		}
	}
}

func serveTab(tab *Tab, notifier chan Notification, window *app.Window) {
	// 1 url = 1 root node
	cache := make(map[string]*parser.Node)
	for noti := range notifier {
		// TODO NOW: have go routine assigned for each tab and delegate
		// the work depends on the tabID
		switch noti.Type {
		case Search:
			preparedUrl, err := prepareUrl(noti.Url)
			if err != nil {
				fmt.Println("prepareUrl: ", err)
				continue
			}
			tab.Url = preparedUrl.String()

			cachedRoot, ok := cache[tab.Url]
			if ok {
				tab.Root = cachedRoot
				window.Invalidate()
				continue
			}

			tab.IsLoading = true
			go reportProgress(tab, window, tab.LoadProgress)

			root, err := getDom(*preparedUrl)
			styles := getStyles(root, preparedUrl)
			tab.IsLoading = false

			if err != nil {
				fmt.Println("search:", err)
				continue
			}

			cache[tab.Url] = root
			tab.Root = root
			tab.Styles = styles
			window.Invalidate()
		default:
			continue
		}
	}
}

func NewState() *State {
	s := State{}
	s.Notifier = make(chan Notification)
	s.Tabs = []*Tab{newTab()}
	return &s
}

func newTab() *Tab {
	return &Tab{LoadProgress: make(chan float32)}
}

// ResolveJumpTarget takes href string and the base url of the site
// to determine the jump target address
func ResolveJumpTarget(href, base string) (string, error) {
	url, err := urlPkg.ParseRequestURI(base)
	if err != nil {
		return "", err
	}
	target, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	return target.String(), nil
}

// getDom fetches the url and parse the DOM tree
// then return the root of DOM tree and error if exists
func getDom(url urlPkg.URL) (*parser.Node, error) {
	contentReader, err := Fetch(url)
	if err != nil {
		return nil, fmt.Errorf("Fetch: %v", err)
	}
	defer contentReader.Close()

	resBody, err := io.ReadAll(contentReader)
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
				contentReader, err := Fetch(*hrefUrl)
				if err != nil {
					log.Println("Fetch:", err)
					continue
				}

				defer contentReader.Close()
				content, err := io.ReadAll(contentReader)
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

// Fetch uses the url to fetch the content and return
// io.ReadCloser representing a content reader
func Fetch(url urlPkg.URL) (io.ReadCloser, error) {
	switch url.Scheme {
	case "file":
		file, err := os.Open(url.Path)
		if err != nil {
			return nil, fmt.Errorf("os.Open: %v", err)
		}
		return file, nil
	case "http":
		res, err := http.Get(url.String())
		if err != nil {
			return nil, fmt.Errorf("http.Get: %v", err)
		}
		return res.Body, nil
	case "https":
		// use serious client for https
		req, err := http.NewRequest(http.MethodGet, url.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("http.NewRequest: %v", err)
		}
		req.Header.Set("User-Agent", "Gazer")

		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("client.Do: %v", err)
		}

		contentTypes, ok := res.Header["Content-Type"]
		if !ok {
			return nil, fmt.Errorf("No content type provided")
		}

		contentType := contentTypes[0]
		contentType = strings.Split(contentType, ";")[0]
		if _, ok := supportedContentType[contentType]; !ok {
			return nil, fmt.Errorf("Unsupported content type: %v", contentType)

		}

		return res.Body, nil
	default:
		return nil, fmt.Errorf("Unsupported scheme: %v", url.Scheme)
	}
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

// prepareUrl takes a url string and return a new url.URL we can Fetch from
// supported scheme: HTTP, HTTPS, file system
func prepareUrl(rawUrl string) (*urlPkg.URL, error) {
	if len(rawUrl) == 0 {
		return nil, fmt.Errorf("Empty URL")
	}

	// handle prefix: we want https:// or http://
	if !strings.HasPrefix(rawUrl, "file://") &&
		!strings.HasPrefix(rawUrl, "https://") &&
		!strings.HasPrefix(rawUrl, "http://") {
		rawUrl = "https://" + rawUrl
	}

	// check normal valid url
	url, err := urlPkg.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %v", err)
	}

	// support local file
	if url.Scheme == "file" {
		return url, nil
	}

	// check valid HTTP request url
	url, err = urlPkg.ParseRequestURI(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("url.ParseRequestURI: %v", err)
	}

	return url, nil
}
