package css

// TODO NOW

import (
	"image/color"
	"testing"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

/*
	type Style struct {
		color    *color.NRGBA
		bgColor  *color.NRGBA
		margin   *layout.Inset
		padding  *layout.Inset
		border   *widget.Border
		fontSize *unit.Dp // TODO: might have to change after supporting other type
	}
*/

type testParseCase struct {
	name     string
	input    string
	expected StyleSet
}

func TestParse(t *testing.T) {
	cases := []testParseCase{
		{
			name: "normal",
			input: `
h1 {
	color: red;
	font-size: 10px;
}

#header {
	background-color: rgb(200, 100, 10);
	border-color: #000;
}

*, .spacer {
	margin: 10px 15px;
	padding: 10px 27px 18px 30px;
}`,
			expected: StyleSet{
				universal: &Style{
					margin: &layout.Inset{
						Top:    unit.Dp(10),
						Bottom: unit.Dp(10),
						Left:   unit.Dp(15),
						Right:  unit.Dp(15),
					},
					padding: &layout.Inset{
						Top:    unit.Dp(10),
						Right:  unit.Dp(27),
						Bottom: unit.Dp(18),
						Left:   unit.Dp(30),
					},
				},
				classStyles: map[string]*Style{
					"spacer": &Style{
						margin: &layout.Inset{
							Top:    unit.Dp(10),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(15),
							Right:  unit.Dp(15),
						},
						padding: &layout.Inset{
							Top:    unit.Dp(10),
							Right:  unit.Dp(27),
							Bottom: unit.Dp(18),
							Left:   unit.Dp(30),
						},
					},
				},
				idStyles: map[string]*Style{
					"header": &Style{
						bgColor: &color.NRGBA{200, 100, 10, 255},
						border: &widget.Border{
							Color: color.NRGBA{0, 0, 0, 255},
						},
					},
				},
			},
		},
	}
}
