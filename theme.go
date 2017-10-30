package gwk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
)

type ThemeFont struct {
	family string
	size   int
	weight string
}

type ThemeStyle struct {
	font           string
	fillColor      string
	textColor      string
	lineColor      string
	fontSize       int
	bgImage        *Image
	fgImage        *Image
	bgImageTips    *Image
	checkedImage   *Image
	uncheckedImage *Image
}

func NewThemeStyle(font, fillColor, textColor, lineColor, bgImage string) *ThemeStyle {
	style := &ThemeStyle{}

	if font {
		style.font = font
	}

	if bgImage {
		style.bgImage = bgImage
	}

	if fillColor {
		style.fillColor = fillColor
	}

	if textColor {
		style.textColor = textColor
	}

	if lineColor {
		style.lineColor = lineColor
	}

	return style
}

type Theme struct {
	normal   *ThemeStyle
	active   *ThemeStyle
	over     *ThemeStyle
	disable  *ThemeStyle
	selected *ThemeStyle
}

func NewTheme() *Theme {
	theme := &Theme{
		NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000", ""),
		NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000", ""),
		NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000", ""),
		NewThemeStyle("13pt bold sans-serif ", "", "Gray", "", ""),
		NewThemeStyle("13pt bold sans-serif ", "", "Gray", "", ""),
	}

	return theme
}

type ThemeManager struct {
	themes       map[string]interface{}
	themesLoaded bool
	imagesURL    string
	defaultTheme *Theme
	themeURL     string
	imagesCache  map[string]Image
}

func NewThemeManager() *ThemeManager {
	manager := &ThemeManager{}
	manager.defaultTheme = NewTheme()
	manager.themeURL = "/ide/theme/default/theme.json"
	manager.themes = make(map[string]interface{})

	return manager
}

type WidgetStates struct {
	stateNormal struct {
		fillColor      string
		bgImage        string
		font           string
		textColor      string
		lineColor      string
		bgImage        interface{}
		fgImage        interface{}
		bgImageTips    interface{}
		checkedImage   interface{}
		uncheckedImage interface{}
	} `json:state-normal`
	stateOver struct {
		fillColor      string
		bgImage        string
		font           string
		textColor      string
		lineColor      string
		bgImage        interface{}
		fgImage        interface{}
		bgImageTips    interface{}
		checkedImage   interface{}
		uncheckedImage interface{}
	} `json:state-over`
	stateActive struct {
		fillColor      string
		bgImage        string
		font           string
		textColor      string
		lineColor      string
		bgImage        interface{}
		fgImage        interface{}
		bgImageTips    interface{}
		checkedImage   interface{}
		uncheckedImage interface{}
	} `json:state-active`
	stateDisable struct {
		fillColor      string
		bgImage        string
		font           string
		textColor      string
		lineColor      string
		bgImage        interface{}
		fgImage        interface{}
		bgImageTips    interface{}
		checkedImage   interface{}
		uncheckedImage interface{}
	} `json:state-disable`
}

type ThemeJson struct {
	imagesURL string
	widgets   struct {
		window WidgetStates
	}
	global struct {
		font struct {
			windows ThemeFont
			linux   ThemeFont
			macosx  ThemeFont
		}
	}
}

func (manager *ThemeManager) setImagesURL(imagesURL string) {
	manager.imagesURL = imagesURL

	return
}

func (manager *ThemeManager) getIconImageURL() {
	return manager.imagesURL
}

func (manager *ThemeManager) getImageURL() {
	return manager.imagesURL
}

func (manager *ThemeManager) createImage(url string) *Image {
	if image, ok := manager.imagesCache[url]; ok {
		return image
	}

	return NewImage(url)
}

func (manager *ThemeManager) getIconImage(name string) *Image {
	if ok := strings.HasSuffix(name, ".png"); ok {
		return manager.getImage(name)
	} else {
		return manager.getImage(name + ".png")
	}
}

func (manager *ThemeManager) getBgImage(name string) *Image {
	return manager.getImage(name)
}

func (manager *ThemeManager) getImage(name string) *Image {
	if _, ok := manager["imagesURL"]; !ok {
		return nil
	}

	url := manager["imagesURL"] + "#" + name
	return manager.createImage(url)
}

func (manager *ThemeManager) setTheme(theme) {
	//TODO
}

func (manager *ThemeManager) getDefaultFont(themeJson *ThemeJson) *ThemeFont {
	return themeJson.global.font.macosx
}

func (manager *ThemeManager) applyDefaultFont(style *ThemeStyle) {
	famlily := "sans"
	size := 10
	weight := "normal"

	style.fontSize = 10
	style.font = fmt.Sprintf("%d %dpx %s", weight, size, famlily)

	return
}

func (manager *ThemeManager) loadTheme(themeURL string, themeJson *ThemeJson) {
	dir := path.Dir(themeURL)
	imagesURL := dir + "/"
	if themeJson.imagesURL {
		imagesURL += themeJson.imagesURL
	} else {
		imagesURL += "images.json"
	}

	manager.setImagesURL(imagesURL)
	font := manager.getDefaultFont(themeJson)
	widgetsTheme := themeJson.widgets

	for _, widgetTheme := range widgetsTheme {
		for _, style := range widgetTheme {
			if style.bgImage {
				style.bgImage = manager.getImage(style.bgImage)
			}
			if style.fgImage {
				style.fgImage = manager.getImage(style.fgImage)
			}
			if style.bgImageTips {
				style.bgImageTips = manager.getImage(style.bgImageTips)
			}
			if style.checkedImage {
				style.checkedImage = manager.getImage(style.checkedImage)
			}
			if style.uncheckedImage {
				style.uncheckedImage = manager.getImage(style.uncheckedImage)
			}
			manager.applyDefaultFont(style, font)
		}
	}

	manager.themes = widgetsTheme
	manager.themesLoaded = true

	return
}

func (manager *ThemeManager) getThemeURL() string {
	return manager["themeURL"]
}

func (manager *ThemeManager) loadThemeURL(url string) {
	if !url {
		url := manager.getThemeURL()
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		themeJSON := &ThemeJson{}
		json.NewDecoder(resp.Body).Decode(themeJSON)
		manager.loadTheme(url, themeJSON)
		wm := GetWindowManagerInstance()
		if wm {
			wm.postRedraw()
		}
	}

	return
}

func (manager *ThemeManager) exist(name string) bool {
	_, exists := manager.themes[name]

	return exists
}

func (manager *ThemeManager) dump() {
	s, _ := json.Marshal(manager.themes)
	fmt.Printf(s)

	return
}

func (manager *ThemeManager) get(name string, noDefault bool) *Theme {
	theme := manager.themes[name]

	if !theme {
		if noDefault {
			manager.themes[name] = NewThemeManager()
			theme = manager.themes[name]
		} else {
			theme = manager.defaultTheme
		}
	}

	return theme
}

func (manager *ThemeManager) set(name, state, font, textColor, fillColor, lineColor, bgImage string) {
	if !state {
		manager.setOneState(name, "normal", font, textColor, fillColor, lineColor, bgImage)
		manager.setOneState(name, "active", font, textColor, fillColor, lineColor, bgImage)
		manager.setOneState(name, "over", font, textColor, fillColor, lineColor, bgImage)
		manager.setOneState(name, "disable", font, textColor, fillColor, lineColor, bgImage)
		manager.setOneState(name, "selected", font, textColor, fillColor, lineColor, bgImage)
	} else {
		manager.setOneState(name, state, font, textColor, fillColor, lineColor, bgImage)
	}

	return
}

func (manager *ThemeManager) setOneState(name, state, font, textColor, fillColor, lineColor, bgImage string) {
	theme := manager.themes[name]

	if !theme {
		theme := NewTheme()
		manager.themes[name] = theme
	}

	if font {
		theme[state].font = font
	}

	if textColor {
		theme[state].textColor = textColor
	}

	if lineColor {
		theme[state].lineColor = lineColor
	}

	if bgImage {
		theme[state].bgImage = bgImage
	}

	return
}
