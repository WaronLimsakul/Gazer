package css

import (
	"gioui.org/font"
)

// a little tricky map to convert css font-weight to gio UI's font weight
var fontWeights = map[string]font.Weight{
	"normal":  font.Normal,
	"bold":    font.Bold,
	"bolder":  font.SemiBold,
	"ligther": font.Light,
	"100":     font.Thin,
	"200":     font.ExtraLight,
	"300":     font.Light,
	"400":     font.Normal,
	"500":     font.Medium,
	"600":     font.SemiBold,
	"700":     font.Bold,
	"800":     font.ExtraBold,
	"900":     font.Black,
}

var fontStyles = map[string]font.Style{
	"normal":  font.Regular,
	"italic":  font.Italic,
	"oblique": font.Italic,
}
