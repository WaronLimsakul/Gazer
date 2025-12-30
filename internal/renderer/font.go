package renderer

import (
	_ "embed"

	"gioui.org/font"
	"gioui.org/font/opentype"
)

//go:embed fonts/inter.ttf
var interBytes []byte

//go:embed fonts/inter_italic.ttf
var interItalicBytes []byte

func loadFont() ([]font.FontFace, error) {
	inter, err := opentype.Parse(interBytes)
	if err != nil {
		return nil, err
	}

	interItalic, err := opentype.Parse(interItalicBytes)
	if err != nil {
		return nil, err
	}

	return []font.FontFace{
		{Font: font.Font{Typeface: "Inter", Style: font.Regular}, Face: inter},
		{Font: font.Font{Typeface: "Inter", Style: font.Italic}, Face: interItalic},
	}, nil
}
