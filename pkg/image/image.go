package image

import (
	"fmt"
	texturePacker "github.com/Luncher/gwk/pkg/texture_packer"
	"honnef.co/go/js/dom"
	"path/filepath"
	"strings"
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

type ImageSizeInfo struct {
	X       int
	Y       int
	W       int
	H       int
	Ox      int
	Oy      int
	Rw      int
	Rh      int
	Rotated bool
	Trimmed bool
}

type Image struct {
	src   string
	rect  *ImageSizeInfo
	image *dom.HTMLImageElement
}

var imagesCache = make(map[string]*Image)

func NewImage(url string) *Image {
	if image, exists := imagesCache[url]; exists {
		return image
	} else {
		image := &Image{src: url, rect: &ImageSizeInfo{}}
		image.SetImageSrc(url)
		imagesCache[url] = image
		return image
	}
}

func (image *Image) GetImage() *dom.HTMLImageElement {
	return image.image
}

func (image *Image) GetImageRect() *ImageSizeInfo {
	return image.rect
}

func (image *Image) isTexturePacker(url string) bool {
	return strings.Index(url, "#") > -1
}

func (image *Image) SetImageSrc(url string) {
	fmt.Printf("SetImageSrc: %s\n", url)
	fmt.Printf("SetImageSrc: %s\n", url)
	if image.isTexturePacker(url) {
		image.setupTexturePackerImage(url)
	} else {
		image.setupNormalImage(url)
	}

	return
}

func (image *Image) setupNormalImage(url string) {
	LoadImage(url, func(imageElement *dom.HTMLImageElement) {
		image.image = imageElement
		image.rect = GetImageRectDefault(imageElement)
	})

	return
}

func (image *Image) setupTexturePackerImage(url string) {
	sepIndex := strings.Index(url, "#")
	jsonPath := url[:sepIndex]
	texturePacker.LoadImagesURL(jsonPath, func(json *texturePacker.TexturePackerJSON) {
		imageName := url[sepIndex+1:]
		imagesName := json.Meta.Image
		imagesUrl := filepath.Dir(jsonPath) + "/" + imagesName

		imageJSON, exists := json.Frames[imageName]
		if !exists {
			return
		}
		image.rect = &ImageSizeInfo{
			X:       imageJSON.Frame.X,
			Y:       imageJSON.Frame.Y,
			W:       imageJSON.Frame.W,
			H:       imageJSON.Frame.H,
			Rotated: imageJSON.Rotated,
			Trimmed: imageJSON.Trimmed,
			Ox:      imageJSON.SpriteSourceSize.X,
			Oy:      imageJSON.SpriteSourceSize.Y,
			Rw:      imageJSON.SourceSize.W,
			Rh:      imageJSON.SourceSize.H,
		}
		LoadImage(imagesUrl, func(img *dom.HTMLImageElement) {
			image.image = img
		})
	})

	return
}

func (image *Image) Draw(context *dom.CanvasRenderingContext2D, display Display, x, y, dw, dh int) {
	imageVal := image.GetImage()
	rect := image.GetImageRect()

	DrawImage(context, imageVal, display, x, y, dw, dh, rect)

	return
}

func GetImageRectDefault(image *dom.HTMLImageElement) *ImageSizeInfo {
	return &ImageSizeInfo{W: image.Width, H: image.Height}
}

func LoadImage(url string, onDone func(*dom.HTMLImageElement)) {
	imageElement := dom.GetWindow().Document().CreateElement("img").(*dom.HTMLImageElement)
	imageElement.AddEventListener("error", false, func(event dom.Event) {
		fmt.Printf("loadImage error%v\n", event)
		onDone(imageElement)
	})

	imageElement.AddEventListener("load", false, func(event dom.Event) {
		fmt.Printf("loadImage %s done\n", url)
	})
	imageElement.Src = url

	return
}

func DrawImage(context *dom.CanvasRenderingContext2D, image *dom.HTMLImageElement, display Display, x, y, dw, dh int, srcRect *ImageSizeInfo) {
	if image == nil || image.Width == 0 {
		return
	}
	sr := srcRect
	if sr == nil {
		sr = GetImageRectDefault(image)
	}

	//TODO
}
