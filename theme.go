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

type ThemeManager struct {
	themes       map[string]Theme
	themesLoaded bool
	imagesURL    string
	defaultTheme Theme
	themeURL     string
	imagesCache  map[string]*Image
}

func NewThemeManager() *ThemeManager {
	manager := &ThemeManager{}
	manager.defaultTheme = NewTheme()
	manager.themeURL = "/ide/theme/default/theme.json"
	// manager.themes = make(map[string]interface{})

	return manager
}

type ThemeStyle struct {
	fillColor      string
	font           string
	fontSize       int
	textColor      string
	lineColor      string
	bgImage        interface{}
	fgImage        interface{}
	bgImageTips    interface{}
	checkedImage   interface{}
	uncheckedImage interface{}
}

func NewThemeStyle(font, fillColor, textColor, lineColor string) *ThemeStyle {
	style := &ThemeStyle{}

	if len(font) != 0 {
		style.font = font
	}

	if len(fillColor) != 0 {
		style.fillColor = fillColor
	}

	if len(textColor) != 0 {
		style.textColor = textColor
	}

	if len(lineColor) != 0 {
		style.lineColor = lineColor
	}

	return style
}

type Theme map[string]*ThemeStyle

func NewTheme() Theme {
	theme := make(Theme)
	theme[STATE_NORMAL] = NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000")
	theme[STATE_ACTIVE] = NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000")
	theme[STATE_OVER] = NewThemeStyle("13pt bold sans-serif ", "", "#000000", "#000000")
	theme[STATE_DISABLE] = NewThemeStyle("13pt bold sans-serif ", "", "Gray", "")
	theme[STATE_SELECTED] = NewThemeStyle("13pt bold sans-serif ", "", "Gray", "")

	return theme
}

type ThemeJson struct {
	imagesURL string
	widgets   map[string]Theme
	global    struct {
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

func (manager *ThemeManager) getIconImageURL() string {
	return manager.imagesURL
}

func (manager *ThemeManager) getImageURL() string {
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
	if len(manager.imagesURL) == 0 {
		return nil
	}

	url := manager.imagesURL + "#" + name

	return manager.createImage(url)
}

func (manager *ThemeManager) setTheme(theme *Theme) {
	//TODO
}

func (manager *ThemeManager) getDefaultFont(themeJson *ThemeJson) *ThemeFont {
	return &themeJson.global.font.macosx
}

func (manager *ThemeManager) applyDefaultFont(style *ThemeStyle, defaultFont *ThemeFont) {
	famlily := "sans"
	size := 10
	weight := "normal"

	style.fontSize = 10
	style.font = fmt.Sprintf("%s %dpx %s", weight, size, famlily)

	return
}

func (manager *ThemeManager) loadTheme(themeURL string, themeJson *ThemeJson) {
	dir := path.Dir(themeURL)
	imagesURL := dir + "/"
	if len(themeJson.imagesURL) > 0 {
		imagesURL += themeJson.imagesURL
	} else {
		imagesURL += "images.json"
	}

	manager.setImagesURL(imagesURL)
	font := manager.getDefaultFont(themeJson)
	widgetsTheme := themeJson.widgets

	for _, widgetTheme := range widgetsTheme {
		for _, style := range widgetTheme {
			if style.bgImage != nil {
				style.bgImage = manager.getImage(style.bgImage.(string))
			}
			if style.fgImage != nil {
				style.fgImage = manager.getImage(style.fgImage.(string))
			}
			if style.bgImageTips != nil {
				style.bgImageTips = manager.getImage(style.bgImageTips.(string))
			}
			if style.checkedImage != nil {
				style.checkedImage = manager.getImage(style.checkedImage.(string))
			}
			if style.uncheckedImage != nil {
				style.uncheckedImage = manager.getImage(style.uncheckedImage.(string))
			}
			manager.applyDefaultFont(style, font)
		}
	}

	manager.themes = widgetsTheme
	manager.themesLoaded = true

	return
}

func (manager *ThemeManager) getThemeURL() string {
	return manager.themeURL
}

func (manager *ThemeManager) loadThemeURL(url string) {
	if url != "" {
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
		if wm != nil {
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
	fmt.Printf(string(s))

	return
}

func (manager *ThemeManager) get(name string, noDefault bool) Theme {
	theme, ok := manager.themes[name]

	if !ok {
		if noDefault {
			manager.themes[name] = NewTheme()
			theme = manager.themes[name]
		} else {
			theme = manager.defaultTheme
		}
	}

	return theme
}

func (manager *ThemeManager) set(name, state, font, textColor, fillColor, lineColor, bgImage string) {
	if state == "" {
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
	theme, ok := manager.themes[name]

	if !ok {
		theme = NewTheme()
		manager.themes[name] = theme
	}

	if font != "" {
		theme[state].font = font
	}

	if textColor != "" {
		theme[state].textColor = textColor
	}

	if lineColor != "" {
		theme[state].lineColor = lineColor
	}

	if bgImage != "" {
		theme[state].bgImage = bgImage
	}

	return
}

var themeManager *ThemeManager

func GetThemeManagerInstance() *ThemeManager {
	if themeManager == nil {
		themeManager = NewThemeManager()
	}

	return themeManager
}
