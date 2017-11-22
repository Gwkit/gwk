package image

import (
	"github.com/Luncher/gwk/pkg/structs"
)

import (
	"honnef.co/go/js/dom"
)

var (
	DISPLAY_9PATCH = 2
)

type Image struct {
	src   string
	rect  structs.Rect
	image *dom.HTMLImageElement
}

func NewImage(url string) *Image {
	return &Image{src: url}
}

func (image *Image) GetImage() *dom.HTMLImageElement {
	return image.image
}

func (image *Image) GetImageRect() structs.Rect {
	return image.rect
}

func (image *Image) Draw(context *dom.CanvasRenderingContext2D, display, x, y, dw, dh int, src structs.Rect) {
	//TODO
	return
}
