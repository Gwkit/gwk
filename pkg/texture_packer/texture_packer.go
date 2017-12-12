package texture_packer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TextureImage struct {
	Frame            TextureImageFrame            `json:"frame"`
	Rotated          bool                         `json:"rotated"`
	Trimmed          bool                         `json:"trimmed"`
	SpriteSourceSize TextureImageSpriteSourceSize `json:"spriteSourceSize"`
	SourceSize       TextureImageSourceSize       `json:"sourceSize"`
	Pivot            TextureImagePivot            `json:"pivot"`
}

type TextureImageFrame struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type TextureImageSpriteSourceSize struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type TextureImageSourceSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type TextureImagePivot struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type TexturePackerMetaSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type TexturePackerMeta struct {
	Version string                `json:"version"`
	Image   string                `json:"image"`
	Format  string                `json:"format"`
	Size    TexturePackerMetaSize `json:"size"`
	Scale   string                `json:"scale"`
}

type TexturePackerJSON struct {
	Frames map[string]TextureImage `json:"frames"`
	Meta   TexturePackerMeta       `json:"meta"`
}

var texturePackerCache = make(map[string]*TexturePackerJSON)

func LoadImagesJSON(url string, reader io.ReadCloser) error {
	decoder := json.NewDecoder(reader)
	var imagesJson TexturePackerJSON
	if err := decoder.Decode(&imagesJson); err != nil {
		fmt.Println(err)
		return err
	}
	defer reader.Close()
	fmt.Printf("add image cache:%s\n", url)
	texturePackerCache[url] = &imagesJson

	return nil
}

func LoadImagesURL(url string, onDone func(*TexturePackerJSON)) error {
	if cache, exists := texturePackerCache[url]; exists {
		onDone(cache)
		return nil
	}
	go func() {
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		} else {
			LoadImagesJSON(url, res.Body)
			onDone(texturePackerCache[url])
		}
	}()

	return nil
}
