package ui

import (
	"image"
	"image/color"
	"image/draw"
)

var defaultFavIcon image.Image

func init() {
	// Create a simple 16x16 gray square with a lighter inner square
	// to represent a generic page/document icon
	size := 16
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Fill with gray background
	gray := color.RGBA{180, 180, 180, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{gray}, image.Point{}, draw.Src)

	// Draw a lighter inner rectangle to suggest a document
	lightGray := color.RGBA{220, 220, 220, 255}
	innerRect := image.Rect(3, 3, 13, 13)
	draw.Draw(img, innerRect, &image.Uniform{lightGray}, image.Point{}, draw.Src)

	defaultFavIcon = img
}
