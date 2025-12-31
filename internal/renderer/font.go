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

//go:embed fonts/inter_bold.ttf
var interBoldBytes []byte

//go:embed fonts/inter_bold_italic.ttf
var interBoldItalicBytes []byte

func loadFont() ([]font.FontFace, error) {
	inter, err := opentype.Parse(interBytes)
	if err != nil {
		return nil, err
	}

	interItalic, err := opentype.Parse(interItalicBytes)
	if err != nil {
		return nil, err
	}

	interBold, err := opentype.Parse(interBoldBytes)
	if err != nil {
		return nil, err
	}

	interBoldItalic, err := opentype.Parse(interBoldItalicBytes)
	if err != nil {
		return nil, err
	}

	return []font.FontFace{
		{Font: font.Font{Typeface: "Inter", Style: font.Regular, Weight: font.Normal}, Face: inter},
		{Font: font.Font{Typeface: "Inter", Style: font.Italic, Weight: font.Normal}, Face: interItalic},
		{Font: font.Font{Typeface: "Inter", Style: font.Regular, Weight: font.Bold}, Face: interBold},
		{Font: font.Font{Typeface: "Inter", Style: font.Italic, Weight: font.Bold}, Face: interBoldItalic},
	}, nil
}
