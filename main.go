package main

import (
	"gioui.org/app"
	"github.com/WaronLimsakul/Gazer/internal/engine"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

func main() {
	w := ui.NewWindow()
	state := engine.NewState()
	go ui.Draw(w, state)
	go engine.Start(state, w)

	app.Main()
}
