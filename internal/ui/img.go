package ui

import (
	"fmt"
	"image"
	"io"
	"net/url"
	"strings"
	"time"

	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"

	"gioui.org/op"
	"gioui.org/op/paint"
	"github.com/WaronLimsakul/Gazer/internal/engine"
)

var imgFormats = []string{".jpg", ".jpeg", ".png", ".gif"}

type Img struct {
	src    string
	format string
	img    image.Image
	isGif  bool
	gifImg *GifImg // nil if not gif format
}

// additional data Img needs to render gif
type GifImg struct {
	// img    *gif.GIF
	start  time.Time     // starting rendering time
	elapse time.Duration // elapse dur for each loop
	// cache for check which frame to rendering
	// usage: if gif (age % elapse) < frameTimeBounds[i],
	// then return composedFrames[i]
	frameTimeBounds []time.Duration
	// precomposed frames for rendering
	composedFrames []image.Image
}

// NewImg creates a new Img component from legal URL src
func NewImg(src string) (*Img, error) {
	parsedUrl, err := url.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %v", err)
	}

	// check if the format is supported
	var imgFormat string
	for _, format := range imgFormats {
		if strings.HasSuffix(parsedUrl.Path, format) {
			imgFormat = format
			break
		}
	}
	if len(imgFormat) == 0 {
		fmt.Println("not supported file format:", parsedUrl.Path)
		return nil, fmt.Errorf("Not supported file format")
	}

	// fetch the image content
	imgReader, err := engine.Fetch(*parsedUrl)
	if err != nil {
		fmt.Println("engine fetch error:", err)
		return nil, err
	}
	defer imgReader.Close()

	// decode the image
	var img image.Image
	var format string
	var gifImg *GifImg
	var isGif bool
	if imgFormat == ".gif" {
		isGif = true
		gifImg, err = newGifImg(imgReader)
		if err != nil {
			return nil, fmt.Errorf("newGifImg: %v", err)
		}
	} else {
		img, format, err = image.Decode(imgReader)
		if err != nil {
			return nil, fmt.Errorf("image.Decode: %v", err)
		}
	}

	return &Img{src: src, format: format, img: img, isGif: isGif, gifImg: gifImg}, nil
}

func (i Img) Layout(gtx C) D {
	var size image.Point
	var img image.Image
	if i.isGif {
		var nextFrameTime time.Time
		img, nextFrameTime = i.gifImg.getGifFrameAndNextFrameTime(time.Now())
		gtx.Execute(op.InvalidateCmd{At: nextFrameTime}) // no one tells me to use this...
	} else {
		img = i.img
	}
	imgOp := paint.NewImageOp(img)
	size = imgOp.Size()
	imgOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return D{Size: gtx.Constraints.Constrain(size)}

}

// newGifImg create a new *GifImg data from the reader r
func newGifImg(r io.Reader) (*GifImg, error) {
	img, err := gif.DecodeAll(r)
	if err != nil {
		return nil, fmt.Errorf("gif.DecodeAll: %v", err)
	}

	// calculate elapse time and frameTimeBounds
	var elapseAcc int // in 1/100 second unit
	frameTimeBounds := make([]time.Duration, len(img.Delay))
	for i, delay := range img.Delay {
		elapseAcc += delay
		frameTimeBounds[i] = time.Duration(elapseAcc) * 10 * time.Millisecond
	}
	elapse := time.Duration(elapseAcc) * 10 * time.Millisecond

	// build precomposed frames
	composedFrames := composeGifFrames(img)

	return &GifImg{
		start:           time.Now(),
		elapse:          elapse,
		frameTimeBounds: frameTimeBounds,
		composedFrames:  composedFrames,
	}, nil
}

// getGifFrameAndNextFrameTime get the gif frame for rendering at the
// exact time t, and also the next time next frame should come.
func (g GifImg) getGifFrameAndNextFrameTime(t time.Time) (image.Image, time.Time) {
	age := t.Sub(g.start) % g.elapse
	for i, bound := range g.frameTimeBounds {
		if age <= bound {
			return g.composedFrames[i], t.Add(bound - age)
		}
	}
	return g.composedFrames[len(g.composedFrames)-1], t.Add(g.elapse - age)
}

// composeGifFrames takes GIF data gg and create a
// complete set of gif frames ready to render one-by-one
func composeGifFrames(gg *gif.GIF) []image.Image {
	if gg == nil {
		return nil
	}

	frames := make([]image.Image, len(gg.Image))
	gifBounds := gg.Image[0].Bounds()
	canvas := image.NewRGBA(gifBounds)

	for i, frame := range gg.Image {
		// TODO: there are some other disposal methods, but hard to support rn
		if i > 0 && gg.Disposal[i-1] == gif.DisposalBackground {
			// clear canvas if specified
			draw.Draw(canvas, gifBounds, image.Transparent, image.Point{}, draw.Src)
		}
		// draw on top of canvas
		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		// save current frame
		curFrame := image.NewRGBA(gifBounds)
		draw.Draw(curFrame, gifBounds, canvas, image.Point{}, draw.Src)
		frames[i] = curFrame
	}
	return frames
}
