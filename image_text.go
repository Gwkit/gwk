package gwk

import (
	"github.com/Luncher/gwk/pkg/image"
	"honnef.co/go/js/dom"
)

type ImageText struct {
	*Widget
	spacer        int
	vertical      bool
	textAlign     string
	textOverImage bool
	image         *image.Image
	fgImageDiplay image.Display
}

func NewImageText(parent *Widget, x, y, w, h float32) *ImageText {
	imageText := &ImageText{
		Widget:        NewWidget(TYPE_IMAGE_TEXT, parent, x, y, w, h),
		fgImageDiplay: image.DISPLAY_AUTO_SIZE_DOWN,
		spacer:        10,
	}
	imageText.border = 2
	imageText.Widget.I = imageText

	return imageText
}

func (imageText *ImageText) getImage() *image.Image {
	return imageText.image
}

func (imageText *ImageText) setImage(image *image.Image) {
	imageText.image = image
}

func (imageText *ImageText) setBorder(border int) {
	imageText.border = border
}

func (imageText *ImageText) setSpacer(spacer int) {
	imageText.spacer = spacer
}

func (imageText *ImageText) setTextOverImage(overImage bool) {
	imageText.textOverImage = overImage
}

func (imageText *ImageText) setVertical(vertical bool) {
	imageText.vertical = vertical
}

func (imageText *ImageText) setFgImageDisplay(display image.Display) {
	imageText.fgImageDiplay = display
}

func (imageText *ImageText) PaintSelf(context *dom.CanvasRenderingContext2D) {
	var x, y, w, h int
	rect := imageText.rect
	border := imageText.border
	text := imageText.GetText()
	image := imageText.getImage()
	style := imageText.getStyle("")
	fontSize := style.FontSize
	if style.FontSize == 0 {
		fontSize = 12
	}

	context.Font = style.Font
	context.FillStyle = style.TextColor
	if len(text) > 0 && image != nil {
		if imageText.textOverImage {
			x = rect.W >> 1
			y = rect.H >> 1
			w = rect.W
			h = rect.H
			image.Draw(context, imageText.fgImageDiplay, 0, 0, w, h)
			context.TextAlign = "center"
			context.TextBaseline = "middle"
		} else {
			if imageText.vertical {
				x = border
				y = border
				w = rect.W - 2*border
				h = rect.H - 2*border - fontSize - 4
			} else {
				x = border
				y = border
				w = rect.H
				h = rect.H - 2*border
			}
			image.Draw(context, imageText.fgImageDiplay, x, y, w, h)
			if imageText.vertical {
				context.TextAlign = "center"
				context.TextBaseline = "bottom"
				x = rect.W >> 1
				y = rect.H - border
			} else {
				context.TextAlign = "left"
				context.TextBaseline = "middle"
				x = rect.H + imageText.spacer
				y = rect.H >> 1
			}
		}
		context.FillText(text, float64(x), float64(y), -1)
	} else if len(text) > 0 {
		if imageText.textAlign == "left" {
			x = border
			y = rect.H >> 1
			context.TextAlign = "left"
			context.TextBaseline = "middle"
			context.FillText(text, float64(x), float64(y), -1)
		} else {
			x = rect.W >> 1
			y = rect.H >> 1
			context.TextAlign = "center"
			context.TextBaseline = "middle"
			context.FillText(text, float64(x), float64(y), -1)
		}
	} else if image != nil {
		x = border
		y = border
		w = rect.W - 2*border
		h = rect.H - 2*border
		image.Draw(context, imageText.fgImageDiplay, x, y, w, h)
	}

	return
}
