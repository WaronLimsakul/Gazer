package main

import (
	"gioui.org/app"
	"github.com/WaronLimsakul/Gazer/internal/engine"
	"github.com/WaronLimsakul/Gazer/internal/renderer"
)

func main() {
	w := renderer.NewWindow()
	state := engine.NewState()
	go renderer.Draw(w, state)
	go engine.Start(state, w)

	app.Main()
}
