package engine

// represent page navigation history timeline
type navHistory struct {
	cur *navHistoryNode
}

type navHistoryNode struct {
	url  string
	prev *navHistoryNode
	next *navHistoryNode
}

func newNavHistory() *navHistory {
	// starts first node with empty search
	return &navHistory{cur: newNavHistoryNode("")}
}

// getUrl get the current url at present time
func (n navHistory) getUrl() string {
	return n.cur.url
}

// back goes back in history 1 step, it stays the same if no past
func (n *navHistory) back() {
	if n.cur.prev != nil {
		n.cur = n.cur.prev
	}
}

// forth goes forward in history, it stays the same if no future
func (n *navHistory) forth() {
	if n.cur.next != nil {
		n.cur = n.cur.next
	}
}

func (n *navHistory) nav(url string) {
	n.cur.next = newNavHistoryNode(url)
	n.cur.next.prev = n.cur
	n.cur = n.cur.next
}

func newNavHistoryNode(url string) *navHistoryNode {
	return &navHistoryNode{url: url}
}
