package engine

import "github.com/WaronLimsakul/Gazer/internal/parser"

type Tab struct {
	Url       string
	Root      *parser.Node
	IsLoading bool
}
