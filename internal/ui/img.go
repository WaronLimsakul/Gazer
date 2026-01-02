package ui

import (
	"fmt"
	"image"
	"net/http"
	"net/url"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"gioui.org/op/paint"
	// TODO: _ "image/gif"
)

var imgFormats = []string{".jpg", ".jpeg", ".png"}

type Img struct {
	src    string
	format string
	img    image.Image
}

func NewImg(src string) (*Img, error) {
	parsedUrl, err := url.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %v", err)
	}
	if parsedUrl.Scheme != "https" {
		return nil, fmt.Errorf("Not https")
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

	res, err := http.Get(src)
	if err != nil {
		return nil, err
	}

	img, format, err := image.Decode(res.Body)
	if err != nil {
		return nil, fmt.Errorf("image.Decode: %v", err)
	}

	return &Img{src: src, format: format, img: img}, nil
}

func (i Img) Layout(gtx C) D {
	// not sure
	imgOp := paint.NewImageOp(i.img)
	size := imgOp.Size()
	imgOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return D{Size: gtx.Constraints.Constrain(size)}
}
