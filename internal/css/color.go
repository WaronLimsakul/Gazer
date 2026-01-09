package css

import "image/color"

var colors = map[string]color.NRGBA{
	"black":       {R: 0, G: 0, B: 0, A: 255},
	"white":       {R: 255, G: 255, B: 255, A: 255},
	"red":         {R: 255, G: 0, B: 0, A: 255},
	"green":       {R: 0, G: 128, B: 0, A: 255},
	"blue":        {R: 0, G: 0, B: 255, A: 255},
	"yellow":      {R: 255, G: 255, B: 0, A: 255},
	"cyan":        {R: 0, G: 255, B: 255, A: 255},
	"magenta":     {R: 255, G: 0, B: 255, A: 255},
	"gray":        {R: 128, G: 128, B: 128, A: 255},
	"grey":        {R: 128, G: 128, B: 128, A: 255},
	"orange":      {R: 255, G: 165, B: 0, A: 255},
	"purple":      {R: 128, G: 0, B: 128, A: 255},
	"pink":        {R: 255, G: 192, B: 203, A: 255},
	"brown":       {R: 165, G: 42, B: 42, A: 255},
	"navy":        {R: 0, G: 0, B: 128, A: 255},
	"lime":        {R: 0, G: 255, B: 0, A: 255},
	"silver":      {R: 192, G: 192, B: 192, A: 255},
	"transparent": {R: 0, G: 0, B: 0, A: 0},
}
