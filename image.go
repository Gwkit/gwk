package gwk

import (
	"honnef.co/go/js/dom"
)

var (
	DISPLAY_9PATCH = 2
)

type Image struct {
	src   string
	rect  Rect
	image *dom.HTMLImageElement
}

func NewImage(url string) *Image {
	return &Image{src: url}
}

func (image *Image) getImage() *dom.HTMLImageElement {
	return image.image
}

func (image *Image) getImageRect() Rect {
	return image.rect
}

func (image *Image) draw(context *dom.CanvasRenderingContext2D, display, x, y, dw, dh int, src Rect) {
	//TODO
	return
}
