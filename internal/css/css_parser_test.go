package css

import (
	"image/color"
	"testing"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/WaronLimsakul/Gazer/internal/parser"
)

type testParseCase struct {
	name     string
	input    string
	expected StyleSet
}

func TestParse(t *testing.T) {
	red, _ := colors["red"]
	size10dp := unit.Dp(10)

	cases := []testParseCase{
		{
			name: "normal",
			input: `
h1 {
	color: red;
	font-size: 10px;
}

#header {
	color: #123abc;
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
					"spacer": {
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
					"header": {
						color:   &color.NRGBA{R: 1*16 + 2, G: 3*16 + 10, B: 11*16 + 12, A: 255},
						bgColor: &color.NRGBA{200, 100, 10, 255},
						border: &widget.Border{
							Color: color.NRGBA{0, 0, 0, 255},
						},
					},
				},
				tagStyles: map[parser.Tag]*Style{
					parser.H1: {
						color:    &red,
						fontSize: &size10dp,
					},
				},
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := StyleParser{}.Parse(testCase.input)
			if err != nil {
				t.Errorf("Parse error: %v", err)
			} else if !styleSetEq(*actual, testCase.expected) {
				// TODO: use go-cmp package to diff string instead.
				t.Errorf("Expected: %v | Got: %v", testCase.expected, *actual)
			}
		})
	}

}
