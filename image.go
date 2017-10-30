package gwk

var (
	DISPLAY_9PATCH = 2
)

type Image struct {
	url string
}

func NewImage(url string) *Image {
	return &Image{url}
}
