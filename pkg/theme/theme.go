package theme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Luncher/gwk/pkg/image"
	"io"
	"net/http"
	"path"
	"strings"
)

type JSONThemeStyle struct {
	FillColor      string `json:"fillColor"`
	Font           string `json:"font"`
	FontSize       int    `json:"fontSize"`
	TextColor      string `json:"textColor"`
	LineColor      string `json:"lineColor"`
	BgImage        string `json:"bgImage"`
	FgImage        string `json:"fgImage"`
	BgImageTips    string `json:"bgImageTips"`
	CheckedImage   string `json:"checkedImage"`
	UncheckedImage string `json:"uncheckedImage"`
}

type ThemeStyle struct {
	FillColor      string
	Font           string
	FontSize       int
	TextColor      string
	LineColor      string
	BgImage        *image.Image
	FgImage        *image.Image
	BgImageTips    *image.Image
	CheckedImage   *image.Image
	UncheckedImage *image.Image
}

type ThemeWidget struct {
	StateNormal          *ThemeStyle `json:"state-normal"`
	StateActive          *ThemeStyle `json:"state-active"`
	StateOver            *ThemeStyle `json:"state-over"`
	StateDisable         *ThemeStyle `json:"state-disable"`
	StateDisableSelected *ThemeStyle `json:"state-disable-selected"`
	StateSelected        *ThemeStyle `json:"state-selected"`
	StateNormalCurrent   *ThemeStyle `json:"state-normal-current"`
}

type ThemeFont struct {
	Family string
	Size   int
	Weight string
}

type ThemeGlobalFont struct {
	font struct {
		windows *ThemeFont
		linux   *ThemeFont
		macosx  *ThemeFont
	}
}

// type Theme map[string]*ThemeWidget

type ThemeJson struct {
	Global    *ThemeGlobalFont        `json:"global"`
	Name      string                  `json:"name"`
	Version   string                  `json:"version"`
	ImagesURL string                  `json:"imagesURL"`
	Widgets   map[string]*ThemeWidget `json:"widgets"`
}

func (style *ThemeStyle) UnmarshalJSON(rawData []byte) error {
	reader := bytes.NewReader(rawData)
	decoder := json.NewDecoder(reader)
	var JsonStyle JSONThemeStyle
	if err := decoder.Decode(&JsonStyle); err != nil {
		return err
	}

	style.FillColor = JsonStyle.FillColor
	style.Font = JsonStyle.Font
	style.FontSize = JsonStyle.FontSize
	style.TextColor = JsonStyle.TextColor
	style.LineColor = JsonStyle.LineColor

	if len(JsonStyle.BgImage) > 0 {
		style.BgImage = image.NewImage(JsonStyle.BgImage)
	}

	if len(JsonStyle.FgImage) > 0 {
		style.FgImage = image.NewImage(JsonStyle.FgImage)
	}

	if len(JsonStyle.BgImageTips) > 0 {
		style.BgImageTips = image.NewImage(JsonStyle.BgImageTips)
	}

	if len(JsonStyle.CheckedImage) > 0 {
		style.CheckedImage = image.NewImage(JsonStyle.CheckedImage)
	}

	if len(JsonStyle.UncheckedImage) > 0 {
		style.UncheckedImage = image.NewImage(JsonStyle.UncheckedImage)
	}

	return nil
}

func NewThemeStyle(font, fillColor, textColor, lineColor string) *ThemeStyle {
	style := &ThemeStyle{}

	if len(font) != 0 {
		style.Font = font
	}

	if len(fillColor) != 0 {
		style.FillColor = fillColor
	}

	if len(textColor) != 0 {
		style.TextColor = textColor
	}

	if len(lineColor) != 0 {
		style.LineColor = lineColor
	}

	return style
}

func NewThemeWidget() *ThemeWidget {
	widgetTheme := &ThemeWidget{
		NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000"),
		NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000"),
		NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000"),
		NewThemeStyle("13pt bold sans-serif ", "", "Gray", ""),
		NewThemeStyle("13pt bold sans-serif ", "", "Gray", ""),
		NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000"),
		NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000"),
	}

	return widgetTheme
}

var themes map[string]*ThemeWidget
var themesLoaded bool
var imagesURL string
var defaultTheme *ThemeWidget
var themeURL string
var imagesCache map[string]*image.Image

func SetImagesURL(url string) {
	imagesURL = url
	return
}

func GetImagesURL() string {
	return imagesURL
}

func GetIconImageURL() string {
	return imagesURL
}

func createImage(url string) *image.Image {
	if image, ok := imagesCache[url]; ok {
		return image
	}
	return image.NewImage(url)
}

func GetImage(name string) *image.Image {
	if len(imagesURL) == 0 {
		return nil
	}
	return createImage(imagesURL + "#" + name)
}

func GetIconImage(name string) *image.Image {
	if ok := strings.HasSuffix(name, ".png"); ok {
		return GetImage(name)
	} else {
		return GetImage(name + ".png")
	}
}

func GetBgImage(name string) *image.Image {
	return GetImage(name)
}

// func SetTheme(theme *Theme) {
// 	//TODO
// }

func GetThemeURL() string {
	//TODO
	return ""
}

func getDefaultFont(themeJson *ThemeJson) *ThemeFont {
	if themeJson.Global != nil {
		return themeJson.Global.font.macosx
	}
	return nil
}

func applyDefaultFont(style *ThemeStyle, font *ThemeFont) {
	size := 10
	if font != nil && font.Size > 0 {
		size = font.Size
	}

	family := "sans"
	if font != nil && len(font.Family) > 0 {
		family = font.Family
	}

	weight := "normal"
	if font != nil && len(font.Weight) > 0 {
		weight = font.Weight
	}

	style.FontSize = size
	style.Font = fmt.Sprintf("%s %dpx %s", weight, size, family)

	return
}

func loadTheme(themeURL string, reader io.ReadCloser) error {
	dir := path.Dir(themeURL)
	url := dir + "/"

	decoder := json.NewDecoder(reader)
	var themeJson ThemeJson
	if err := decoder.Decode(&themeJson); err != nil {
		return err
	}
	defer reader.Close()

	if len(themeJson.ImagesURL) > 0 {
		url += themeJson.ImagesURL
	} else {
		url += "images.json"
	}
	SetImagesURL(url)
	font := getDefaultFont(&themeJson)

	for _, widgetTheme := range themeJson.Widgets {
		if widgetTheme.StateNormal != nil {
			applyDefaultFont(widgetTheme.StateNormal, font)
		}
		if widgetTheme.StateActive != nil {
			applyDefaultFont(widgetTheme.StateActive, font)
		}
		if widgetTheme.StateOver != nil {
			applyDefaultFont(widgetTheme.StateOver, font)
		}
		if widgetTheme.StateDisable != nil {
			applyDefaultFont(widgetTheme.StateDisable, font)
		}
		if widgetTheme.StateSelected != nil {
			applyDefaultFont(widgetTheme.StateSelected, font)
		}
		if widgetTheme.StateDisableSelected != nil {
			applyDefaultFont(widgetTheme.StateDisableSelected, font)
		}
		if widgetTheme.StateNormalCurrent != nil {
			applyDefaultFont(widgetTheme.StateNormalCurrent, font)
		}
	}
	themesLoaded = true
	themes = themeJson.Widgets

	return nil
}

func LoadThemeURL(url string) error {
	go func() {
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		} else {
			loadTheme(url, res.Body)
		}
	}()

	return nil
}

func Get(name string, noDefault bool) *ThemeWidget {
	theme := themes[name]

	if theme != nil {
		if noDefault {
			themes[name] = NewThemeWidget()
			theme = themes[name]
		} else {
			theme = defaultTheme
		}
	}

	return theme
}

func init() {
	imagesURL = ""
	themesLoaded = false
	themeURL = "/ide/theme/default/theme.json"
	defaultTheme = NewThemeWidget()

	return
}
