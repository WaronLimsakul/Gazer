package css

// TODO NOW

import (
	"testing"
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
					// TODO
				},
			},
		},
	}
}
