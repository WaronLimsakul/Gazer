package main

import (
	"gioui.org/app"
	"github.com/WaronLimsakul/Gazer/internal/ui"
)

func main() {
	var url string
	go ui.Draw(ui.UserInput{Url: &url})
	app.Main()
}
