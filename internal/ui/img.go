package ui

import (
	"fmt"
	"image"
	"net/url"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"gioui.org/op/paint"
	"github.com/WaronLimsakul/Gazer/internal/engine"
	// TODO: _ "image/gif"
)

var imgFormats = []string{".jpg", ".jpeg", ".png", ".gif"}

type Img struct {
	src    string
	format string
	img    image.Image
}

// NewImg creates a new Img component from legal URL src
func NewImg(src string) (*Img, error) {
	parsedUrl, err := url.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %v", err)
	}

	var supported bool
	for _, format := range imgFormats {
		if strings.HasSuffix(parsedUrl.Path, format) {
			supported = true
			break
		}
	}
	if !supported {
		return nil, fmt.Errorf("Not supported file format")
	}

	imgReader, err := engine.Fetch(*parsedUrl)
	if err != nil {
		return nil, err
	}
	defer imgReader.Close()

	// TODO NOW: If it's gif, you might need to use gif.DecodeAll() to get all
	// the frames, save it somewhere, then figure out the frame based on time.
	img, format, err := image.Decode(imgReader)
	if err != nil {
		return nil, fmt.Errorf("image.Decode: %v", err)
	}

	return &Img{src: src, format: format, img: img}, nil
}

func (i Img) Layout(gtx C) D {
	// TODO: rescale a bit to fit (like can't be bigger than screen)
	imgOp := paint.NewImageOp(i.img)
	size := imgOp.Size()
	imgOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return D{Size: gtx.Constraints.Constrain(size)}
}
