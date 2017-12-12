package image

import (
	"github.com/Luncher/gwk/pkg/structs"
)

import (
	"honnef.co/go/js/dom"
)

type Display int

const (
	UNKNOW Display = iota
	DISPLAY_9PATCH
	DISPLAY_AUTO_SIZE_DOWN
)

func (s Display) String() string {
	switch s {
	case DISPLAY_9PATCH:
		return "9patch"
	case DISPLAY_AUTO_SIZE_DOWN:
		return "auto_size_down"
	default:
		return "unknow"
	}
}

type Image struct {
	src   string
	rect  *structs.Rect
	image *dom.HTMLImageElement
}

func NewImage(url string) *Image {
	return &Image{src: url}
}

func (image *Image) GetImage() *dom.HTMLImageElement {
	return image.image
}

func (image *Image) GetImageRect() *structs.Rect {
	return image.rect
}

func (image *Image) Draw(context *dom.CanvasRenderingContext2D, display Display, x, y, dw, dh int) {
	imageVal := image.GetImage()
	rect := image.GetImageRect()

	DrawImage(context, imageVal, display, x, y, dw, dh, rect)

	return
}

func DrawImage(context *dom.CanvasRenderingContext2D, image *dom.HTMLImageElement, display Display, x, y, dw, dh int, srcRect *structs.Rect) {
	if image == nil || image.Width == 0 {
		return
	}
	sr := srcRect
	if sr == nil {
		sr = GetImageRectDefault(image)
	}

	//TODO
}

func GetImageRectDefault(image *dom.HTMLImageElement) *structs.Rect {
	return structs.NewRect(0, 0, image.Width, image.Height)
}
