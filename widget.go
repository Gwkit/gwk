package gwk

import (
	"fmt"
	"github.com/Luncher/gwk/pkg/utils"
	"honnef.co/go/js/dom"
	"math"
	"strings"
)

var (
	STATE_NORMAL           = "state-normal"
	STATE_ACTIVE           = "state-active"
	STATE_OVER             = "state-over"
	STATE_DISABLE          = "state-disable"
	STATE_DISABLE_SELECTED = "state-diable-selected"
	STATE_SELECTED         = "state-selected"
	STATE_NORMAL_CURRENT   = "state-normal-current"
)

var (
	TYPE_NONE                = 0
	TYPE_USER                = 13
	TYPE_FRAME               = "frame"
	TYPE_FRAMES              = "frames"
	TYPE_TOOLBAR             = "toolbar"
	TYPE_TITLEBAR            = "titlebar"
	TYPE_MINIMIZE_BUTTON     = "button.minimize"
	TYPE_FLOAT_MENU_BAR      = "float-menubar"
	TYPE_POPUP               = "popup"
	TYPE_DIALOG              = "dialog"
	TYPE_DRAGGALE_DIALOG     = "draggable-dialog"
	TYPE_WINDOW              = "window"
	TYPE_VBOX                = "vbox"
	TYPE_HBOX                = "hbox"
	TYPE_MENU                = "menu"
	TYPE_MENU_BAR            = "menu-bar"
	TYPE_MENU_BUTTON         = "menu.button"
	TYPE_GRID_ITEM           = "grid-item"
	TYPE_MENU_ITEM           = "menu.item"
	TYPE_MENU_BAR_ITEM       = "menubar.item"
	TYPE_CONTEXT_MENU_ITEM   = "contextmenu.item"
	TYPE_CONTEXT_MENU_BAR    = "contextmenu-bar"
	TYPE_VSCROLL_BAR         = "vscroll-bar"
	TYPE_HSCROLL_BAR         = "hscroll-bar"
	TYPE_SCROLL_VIEW         = "scroll-bar"
	TYPE_GRID_VIEW           = "grid-view"
	TYPE_LIST_VIEW           = "list-view"
	TYPE_LIST_ITEM           = "list-item"
	TYPE_LIST_ITEM_RADIO     = "list-item-radio"
	TYPE_IMAGE_VIEW          = "image-view"
	TYPE_TREE_VIEW           = "tree-view"
	TYPE_TREE_ITEM           = "tree-item"
	TYPE_ACCORDION           = "accordion"
	TYPE_ACCORDION_ITEM      = "accordion-item"
	TYPE_ACCORDION_TITLE     = "accordion-title"
	TYPE_PROPERTY_TITLE      = "property-title"
	TYPE_PROPERTY_SHEET      = "property-sheet"
	TYPE_PROPERTY_SHEETS     = "property-sheets"
	TYPE_VIEW_BASE           = "view-base"
	TYPE_COMPONENT_MENU_ITEM = "menuitem.component"
	TYPE_WINDOW_MENU_ITEM    = "menuitem.window"
	TYPE_MESSAGE_BOX         = "messagebox"
	TYPE_IMAGE_TEXT          = "icon-text"
	TYPE_BUTTON              = "button"
	TYPE_KEY_VALUE           = "key-value"
	TYPE_LABEL               = "label"
	TYPE_LINK                = "link"
	TYPE_EDIT                = "edit"
	TYPE_TEXT_AREA           = "text-area"
	TYPE_COMBOBOX            = "combobox"
	TYPE_SLIDER              = "slider"
	TYPE_PROGRESSBAR         = "progressbar"
	TYPE_RADIO_BUTTON        = "radio-button"
	TYPE_CHECK_BUTTON        = "check-button"
	TYPE_COLOR_BUTTON        = "color-button"
	TYPE_TAB_BUTTON          = "tab-button"
	TYPE_TAB_CONTROL         = "tab-control"
	TYPE_TAB_BUTTON_GROUP    = "tab-button-group"
	TYPE_TIPS                = "tips"
	TYPE_HLAYOUT             = "h-layout"
	TYPE_VLAYOUT             = "v-layout"
	TYPE_BUTTON_GROUP        = "button-group"
	TYPE_COMBOBOX_POPUP      = "combobox-popup"
	TYPE_COMBOBOX_POPUP_ITEM = "combobox-popup-item"
	TYPE_COLOR_EDIT          = "color-edit"
	TYPE_RANGE_EDIT          = "range-edit"
	TYPE_FILENAME_EDIT       = "filename-edit"
	TYPE_FILENAMES_EDIT      = "filenames-edit"
	TYPE_CANVAS_IMAGE        = "canvas-image"
	TYPE_ICON_BUTTON         = "icon-button"
)

var (
	BORDER_STYLE_NONE   = 0
	BORDER_STYLE_LEFT   = 1
	BORDER_STYLE_RIGHT  = 2
	BORDER_STYLE_TOP    = 4
	BORDER_STYLE_BOTTOM = 8
	BORDER_STYLE_ALL    = 0xffff
)

type CheckEnable func()
type RemovedHandler func()
type WidgeVisit func(*Widget)
type StateChangedHandler func(state string)
type OnMovedHandler func()
type OnResizedHandler func()
type ContextMenuHandler func(*Point)
type ClickedHandler func(*Widget, *Point)
type DoubleClickedHandler func(*Point)
type LongPressHandler func(*Point)
type KeyDownHandler func()
type KeyUpHandler func()
type OnBeforePaintHandler func(*dom.CanvasRenderingContext2D)
type OnAfterPaintHandler func(*dom.CanvasRenderingContext2D)
type WheelHandler func(float64)

type Widget struct {
	id                   int
	t                    string
	name                 string
	rect                 *Rect
	pointerDown          bool
	visible              bool
	state                string
	parent               *Widget
	text                 string
	tag                  string
	tips                 string
	enable               bool
	checkEnable          CheckEnable
	removedHandler       RemovedHandler
	children             []*Widget
	point                *Point
	cursor               string
	imageDisplay         int
	borderStyle          int
	border               int
	themeType            string
	selected             bool
	selectable           bool
	needRelayout         bool
	isScrollView         bool
	xOffset              int
	yOffset              int
	target               *Widget
	inputTips            string
	leftMargin           int
	editing              bool
	userData             interface{}
	onMoved              OnMovedHandler
	stateChangedHandler  StateChangedHandler
	onSized              OnResizedHandler
	contextMenuHandler   ContextMenuHandler
	keyUpHandler         KeyUpHandler
	keyDownHandler       KeyDownHandler
	clickedHandler       ClickedHandler
	doubleClickedHandler DoubleClickedHandler
	wheelHandler         WheelHandler
	lineWidth            int
	roundRadius          int
	theme                map[string]*ThemeStyle
	onBeforePaint        OnBeforePaintHandler
	onAfterPaint         OnAfterPaintHandler
	paintFocusLater      bool
}

func NewWidget(parent *Widget, x, y, w, h float32) *Widget {
	widget := &Widget{
		parent:       parent,
		cursor:       "default",
		borderStyle:  BORDER_STYLE_ALL,
		imageDisplay: DISPLAY_9PATCH,
		rect:         &Rect{x: int(x), y: int(y), w: int(w), h: int(h)},
	}

	widget.setState(STATE_NORMAL)

	if widget.parent != nil {
		var border int
		if border = 0; parent.border > 0 {
			border = parent.border
		}

		pw := parent.rect.w - 2*border
		ph := parent.rect.h - 2*border

		if x > 0 && x < 1 {
			widget.rect.x = int(float32(pw)*x + float32(border))
		}

		if w > 0 && w <= 1 {
			widget.rect.w = int(float32(pw) * w)
		}

		if y > 0 && y < 1 {
			widget.rect.y = int(float32(ph)*y + float32(border))
		}

		if h > 0 && h <= 1 {
			widget.rect.h = int(float32(ph) * w)
		}

		parent.appendChild(widget)
	}

	return widget
}

func (w *Widget) useTheme(t string) *Widget {
	w.themeType = t

	return w
}

func (w *Widget) isSelected(value bool) bool {
	return w.selected
}

func (w *Widget) setSelected(value bool) *Widget {
	w.selected = value

	return w
}

func (w *Widget) setSelectable(selectable bool) bool {
	w.selectable = selectable

	return true
}

func (w *Widget) setNeedRelayout(value bool) {
	w.needRelayout = value

	return
}

func (w *Widget) onAppendChild(child *Widget) {

}

func (w *Widget) appendChild(child *Widget) {
	child.parent = w
	w.children = append(w.children, child)
	w.onAppendChild(child)
	w.setNeedRelayout(true)

	return
}

func (w *Widget) getFrameRate() int {
	return GetWindowManagerInstance().getFrameRate()
}

func (w *Widget) showFPS(maxFpsMode bool) {
	GetWindowManagerInstance().showFPS(maxFpsMode)
}

func (w *Widget) isPointerDown() bool {
	GetWindowManagerInstance().isPointerDown()
}

func (w *Widget) isClicked() bool {
	win := w.getWindow()
	if win != nil {
		return win.isClicked()
	} else {
		return GetWindowManagerInstance().isClicked()
	}
}

func (w *Widget) isAltDown() bool {
	return GetWindowManagerInstance().isAltDown()
}

func (w *Widget) isCtrlDown() bool {
	return GetWindowManagerInstance().isCtrlDown()
}

func (w *Widget) getApp() *Application {
	return GetWindowManagerInstance().getApp()
}

func (w *Widget) getCanvas2D() *dom.CanvasRenderingContext2D {
	return GetWindowManagerInstance().getCanvas2D()
}

func (w *Widget) getCanvas() *dom.HTMLCanvasElement {
	return GetWindowManagerInstance().getCanvas()
}

func (w *Widget) getLastPointerPoint() *Point {
	return GetWindowManagerInstance().getLastPointerPoint()
}

func (w *Widget) getTopWindow() *Widget {
	return w.getWindow()
}

func (w *Widget) getWindow() *Widget {
	if w.parent == nil {
		return w
	}

	p := w.parent
	for ; p != nil; p = p.parent {
	}

	return p
}

func (w *Widget) getParent() *Widget {
	return w.parent
}

func (w *Widget) getX() int {
	return w.rect.x
}

func (w *Widget) getY() int {
	return w.rect.y
}

func (w *Widget) getWidth() int {
	return w.rect.w
}

func (w *Widget) getHeight() int {
	return w.rect.h
}

func (w *Widget) getPositionInView() *Point {
	x := w.getX()
	y := w.getY()
	point := &Point{}

	for iter := w.getParent(); iter != nil; iter = iter.getParent() {
		x += iter.getX()
		y += iter.getY()
		if iter.isScrollView {
			x -= iter.xOffset
			y -= iter.yOffset
		}
	}

	point.x = x
	point.y = y

	return point
}

func (w *Widget) getAbsPosition() *Point {
	x := w.rect.x
	y := w.rect.y

	for parent := w.parent; parent != nil; parent = parent.parent {
		x += parent.getX()
		y += parent.getY()
	}

	return &Point{x, y}
}

func (w *Widget) getPositionInWindow() *Point {
	point := &Point{}

	if w.parent != nil {
		for iter := w; iter != nil; iter = iter.parent {
			if iter.parent == nil {
				break
			}

			point.x += iter.rect.x
			point.y += iter.rect.y
		}
	}

	return point
}

func (w *Widget) translatePoint(point *Point) *Point {
	p := w.getAbsPosition()

	return &Point{x: point.x - p.x, y: point.y - p.y}
}

func (w *Widget) postRedrawAll() {
	GetWindowManagerInstance().postRedraw()

	return
}

func (w *Widget) postRedraw() {
	GetWindowManagerInstance().postRedraw()

	return
}

func (w *Widget) redraw(rect *Rect) {
	p := w.getAbsPosition()

	if rect == nil {
		rect = &Rect{x: 0, y: 0, w: w.rect.w, h: w.rect.h}
	}

	//TODO
	// GetWindowManagerInstance().redraw(rect)

	return
}

func (w *Widget) isPointIn(point *Point) bool {
	return isPointInRect(point, w.rect)
}

func (w *Widget) findTargetWidgetEx(point *Point, recursive bool) *Widget {
	if !w.visible || !w.isPointIn(point) {
		return nil
	}

	if recursive && len(w.children) > 0 {
		n := len(w.children) - 1
		p := w.point
		p.x = point.x - w.rect.x
		p.y = point.y - w.rect.y

		for i := n; i > 0; i-- {
			iter := w.children[i]
			ret := iter.findTargetWidgetEx(p, false)
			if ret != nil {
				return ret
			}
		}
	}

	return w
}

func (w *Widget) findTargetWidget(point *Point) *Widget {
	return w.findTargetWidgetEx(point, true)
}

func (w *Widget) setRemovedHandler(removeHandler RemovedHandler) *Widget {
	w.removedHandler = removeHandler

	return w
}

func (w *Widget) onRemoved() {
	if w.removedHandler != nil {
		w.removedHandler()
	}

	return
}

func (w *Widget) removeChild(child *Widget) *Widget {
	child.remove()

	return w
}

func (w *Widget) remove() *Widget {
	parent := w.parent
	if parent != nil {
		for i, child := range parent.children {
			if child == w {
				parent.children = append(parent.children[0:i], parent.children[i+1:]...)
				break
			}
		}

		if parent.target == w {
			parent.target = nil
		}

		w.parent = nil
		w.onRemoved()
		parent.setNeedRelayout(true)
	}

	return w
}

func (w *Widget) cleanUp() {

}

func (w *Widget) destroy() {
	if len(w.children) > 0 {
		w.destroyChildren()
	}

	w.remove()
	w.cleanUp()

	return
}

func (w *Widget) destroyChildren() {
	for _, child := range w.children {
		child.destroy()
	}
	w.target = nil
	w.children = w.children[:0]
	w.setNeedRelayout(true)

	return
}

func (w *Widget) forEachChild(onVisit WidgeVisit) {
	for _, child := range w.children {
		onVisit(child)
	}

	return
}

func (w *Widget) setTextOf(name, text string, notify bool) *Widget {
	child := w.lookup(name, true)

	if child != nil {
		child.setText(text, notify)
	} else {
		fmt.Printf("not found %s", name)
	}

	return child
}

func (w *Widget) setVisibleOf(name, value string) *Widget {
	child := w.lookup(name, true)

	if child {
		child.setVisible(value)
	} else {
		fmt.Printf("not found %s", name)
	}

	return child
}

func (w *Widget) setText(text string) *Widget {
	w.text = text
	w.setNeedRelayout(true)

	return w
}

func (w *Widget) getText() string {
	return w.text
}

func (w *Widget) setTips(tips string) *Widget {
	w.tips = tips

	return w
}

func (w *Widget) getTips() string {
	return w.tips
}

func (w *Widget) setInputTips(tips string) *Widget {
	w.inputTips = tips

	return w
}

func (w *Widget) getInputTips() string {
	return w.inputTips
}

func (widget *Widget) drawInputTips(context *dom.CanvasRenderingContext2D) {
	h := widget.rect.h
	w := widget.rect.w
	y := widget.rect.h >> 1
	x := widget.leftMargin
	text := widget.getText()
	inputTips := widget.getInputTips()

	if len(text) > 0 || len(inputTips) == 0 || widget.t != TYPE_EDIT || widget.editing {
		return
	}

	style := widget.getStyle()
	context.Save()
	context.Font = style.font
	context.FillStyle = "#E0E0E0"

	context.BeginPath()
	context.Rect(0, 0, float64(w-x), float64(h))
	context.Clip()

	context.TextAlign = "left"
	context.TextBaseline = "middle"
	context.FillText(inputTips, float64(x), float64(y), -1)

	context.Restore()

	return
}

func (w *Widget) drawTips(context *dom.CanvasRenderingContext2D) {
	tips := w.getTips()
	if len(tips) > 0 {
		style := w.getStyle()
		x := w.rect.w >> 1
		y := w.rect.h >> 1
		font := style.tipsFont
		if len(font) == 0 {
			font = style.font
		}
		textColor := style.tipsTextColor
		if len(textColor) == 0 {
			textColor = style.textColor
		}

		if len(font) && len(textColor) {
			context.TextAlign = "center"
			context.TextBaseline = "middle"
			context.Font = font
			context.FillStyle = textColor
			context.FillText(tips, float64(x), float64(y), -1)
		}
	}

	return
}

func (w *Widget) setID(id int) *Widget {
	w.id = id

	return w
}

func (w *Widget) getID() int {
	return w.id
}

func (w *Widget) setName(name string) *Widget {
	w.name = name

	return w
}

func (w *Widget) getName() string {
	return w.name
}

func (w *Widget) setTag(tag string) *Widget {
	w.tag = tag

	return w
}

func (w *Widget) getTag() string {
	return w.tag
}

func (w *Widget) setUserDate(data interface{}) *Widget {
	w.userData = data

	return w
}

func (w *Widget) getUserData() interface{} {
	return w.userData
}

func (w *Widget) setEnable(enable bool) *Widget {
	w.enable = enable

	return w
}

func (w *Widget) changeCursor() *Widget {
	canvas := w.getCanvas()
	if canvas.Style().GetPropertyValue("cursor") != w.cursor {
		canvas.Style().SetProperty("cursor", w.cursor, "")
	}

	return w
}

func (w *Widget) onStateChanged(state string) *Widget {
	if w.stateChangedHandler != nil {
		w.stateChangedHandler(state)
	}

	if state == STATE_OVER || state == STATE_ACTIVE {
		w.changeCursor()
	}

	return w
}

func (w *Widget) setState(state string, recursive bool) *Widget {
	if w.state != state {
		w.state = state
		w.onStateChanged(state)
		if state == STATE_OVER {
			GetWindowManagerInstance().setTipsWidget(w)
		}
	}

	if recursive && w.target != nil {
		me.target.setState(state, recursive)
	}

	return w
}

func (w *Widget) move(x, y int) *Widget {
	w.rect.x = x
	w.rect.y = y
	if me.onMove != nil {
		w.onMoved()
	}

	return w
}

func (w *Widget) moveToCenter(moveX, moveY int) *Widget {
	pw := w.parent.rect.w
	ph := w.parent.rect.h

	if moveX {
		w.rect.x = (pw - w.rect.w) >> 1
	}

	if moveY {
		w.rect.y = (ph - w.rect.h) >> 1
	}

	return w
}

func (w *Widget) moveToBottom(border int) *Widget {
	ph := w.parent.rect.h
	w.rect.y = ph - w.rect.h - border

	return w
}

func (w *Widget) moveDelta(dx, dy int) *Widget {
	w.rect.x = w.rect.x + dx
	w.rect.y = w.rect.y + dy
	if w.onMoved {
		w.onMoved()
	}

	return w
}

func (widget *Widget) resize(w, h int) *Widget {
	widget.rect.w = w
	widget.rect.h = h
	if w.onSized {
		w.onSized()
	}
	w.setNeedRelayout(true)

	return w
}

func (w *Widget) setStateChangedHandler(stateChangedHandler StateChangedHandler) *Widget {
	w.stateChangedHandler = stateChangedHandler

	return w
}

func (w *Widget) setContextMenuHandler(contextMenuHandler ContextMenuHandler) *Widget {
	w.contextMenuHandler = contextMenuHandler

	return w
}

func (w *Widget) setClickedHandler(clickedHandler ClickedHandler) *Widget {
	w.clickedHandler = clickedHandler

	return w
}

func (w *Widget) setKeyDownHandler(keyDownHandler KeyDownHandler) *Widget {
	w.keyDownHandler = keyDownHandler

	return w
}

func (w *Widget) setKeyUpHandler(keyUpHandler KeyUpHandler) *Widget {
	w.keyUpHandler = keyUpHandler

	return w
}

func (w *Widget) onClicked(point *Point) bool {
	if w.clickedHandler {
		w.clickedHandler(w, point)
	}

	w.postRedraw()

	return w.clickedHandler != nil
}

func (w *Widget) lookup(id int, recursive bool) *Widget {
	for _, child := range w.children {
		if child.id == id {
			return child
		}
	}

	if recursive {
		for _, child := range w.children {
			ret := child.lookup(id, recursive)
			if ret != nil {
				return ret
			}
		}
	}

	return nil
}

func (w *Widget) onRelayout(context *dom.HTMLCanvasElement, force bool) {

}

func (w *Widget) relayout(canvas *dom.CanvasRenderingContext2D, force bool) *Widget {
	if !w.needRelayout || !force || !len(w.children) {
		return w
	}

	w.onRelayout(canvas, force)
	w.needRelayout = false

	return w
}

func (w *Widget) setLineWidth(lineWidth int) *Widget {
	w.lineWidth = lineWidth

	return w
}

func (w *Widget) getLineWidth(style dom.CSSStyleDeclaration) int {
	if w.lineWidth {
		return w.lineWidth
	} else {
		return style.GetPropertyValue("lineWidth")
	}
}

func (w *Widget) setRoundRadius(roundRadius int) *Widget {
	w.roundRadius = roundRadius

	return w
}

func (w *Widget) ensureTheme() *Widget {
	if len(w.themeType) {
		w.theme = GetThemeManagerInstance().get(w.themeType)
	} else {
		w.theme = GetThemeManagerInstance().get(w.t)
	}

	return w
}

func (w *Widget) getStyle(_state string) *ThemeStyle {
	var style *ThemeStyle
	w.ensureTheme()
	var state string
	if state = _state; len(state) == 0 {
		state = w.state
	}

	if !w.enable {
		if w.selectable && w.isSelected() {
			style = w.theme[STATE_DISABLE_SELECTED]
		} else {
			style = w.theme[STATE_DISABLE_DISABLE]
		}
	} else {
		if w.selectable && w.selected {
			style = w.theme[STATE_SELECTED]
		} else if state == STATE_OVER {
			style = w.theme[STATE_OVER]
		} else if state == STATE_ACTIVE {
			style = w.theme[STATE_ACTIVE]
		} else {
			style = w.theme[STATE_NORMAL]
		}
	}

	if !style {
		style = w.theme[STATE_NORMAL]
	}

	return style
}

func (w *Widget) setImageDisplay(imageDisplay int) *Widget {
	w.imageDisplay = imageDisplay

	return w
}

func (w *Widget) setBorderStyle(borderStyle int) {
	w.borderStyle = borderStyle

	return w
}

func (w *Widget) paintBackground(canvas dom.HTMLCanvasElement) {
	style := w.getStyle()
	if style {
		if style.bgImage {
			w.paintBackgroundImage(canvas, style)
		} else {
			w.paintBackgroundColor(canvas, style)
		}
	}
}

func (w *Widget) paintBackgroundImage(canvas dom.HTMLCanvasElement, style *ThemeStyle) {
	dst := w.rect
	image := style.bgImage.getImage()
	src := style.bgImage.getImageRect()

	var imageDisplay int
	imageDisplay = w.imageDisplay
	if style.imageDisplay {
		imageDisplay = style.imageDisplay
	}

	if image {
		var topOut, leftOut, rightOut, bottomOut int
		if style.topOut {
			topOut = style.topOut
		}
		if style.leftOut {
			leftOut = style.leftOut
		}
		if style.rightOut {
			rightOut = style.rightOut
		}
		if style.bottomOut {
			bottomOut = style.bottomOut
		}

		x := -leftOut
		y := topOut
		w := dst.w + rightOut + leftOut
		h := dst.h + bottomOut + topOut

		style.bgImage.draw(canvas, imageDisplay, x, y, w, h, src)
	}

	return
}

func (widget *Widget) paintLeftBorder(context *dom.CanvasRenderingContext2D, w, h int) {
	context.BeginPath()
	context.MoveTo(0, 0)
	context.LineTo(0, h)
	context.Stroke()
}

func (widget *Widget) paintRightBorder(context *dom.CanvasRenderingContext2D, w, h int) {
	context.BeginPath()
	context.MoveTo(w, 0)
	context.LineTo(w, h)
	context.Stroke()
}

func (widget *Widget) paintTopBorder(context *dom.CanvasRenderingContext2D, w, h int) {
	context.BeginPath()
	context.MoveTo(0, 0)
	context.LineTo(w, 0)
	context.Stroke()
}

func (widget *Widget) paintBottomBorder(context *dom.CanvasRenderingContext2D, w, h int) {
	context.BeginPath()
	context.MoveTo(0, h)
	context.LineTo(w, h)
	context.Stroke()
}

func (w *Widget) paintBackgroundColor(context dom.CanvasRenderingContext2D, style ThemeStyle) {
	dst := w.rect
	context.BeginPath()
	if w.roundRadius || style.roundRadius {
		roundRadius := math.Min((dst.h>>1)-1, style.roundRadius)
		utils.DrawRoundRect(context, dst.w, dst.h, roundRadius)
	} else {
		context.Rect(0, 0, dst.w, dst.h)
	}

	if style.fillColor {
		context.FillStyle = style.fillColor
		context.Fill()
	}

	lineWidth := w.getLineWidth(style)
	if !lineWidth || !style.lineColor || w.borderStyle == BORDER_STYLE_NONE {
		return
	}

	width := w.getWidth()
	height := w.getHeight()
	context.LineWidth = lineWidth
	context.StrokeStyle = style.lineColor
	if w.borderStyle == BORDER_STYLE_ALL {
		context.Stroke()
		context.BeginPath()
		return
	}

	if w.borderStyle & BORDER_STYLE_LEFT {
		w.paintLeftBorder(context, width, height)
	}

	if w.borderStyle & BORDER_STYLE_RIGHT {
		w.paintRightBorder(context, width, height)
	}

	if w.borderStyle & BORDER_STYLE_TOP {
		w.paintTopBorder(context, width, height)
	}

	if w.borderStyle & BORDER_STYLE_BOTTOM {
		w.paintBottomBorder(context, width, height)
	}
	context.BeginPath()

	return
}

func (w *Widget) paintSelf(context dom.CanvasRenderingContext2D) *Widget {
	return w
}

func (w *Widget) beforePaint(context dom.CanvasRenderingContext2D) *Widget {
	if w.onBeforePaint {
		w.onBeforePaint(context)
	}
	return w
}

func (w *Widget) afterPaint(context dom.CanvasRenderingContext2D) *Widget {
	if w.onAfterPaint {
		w.onAfterPaint(context)
	}
	return w
}

func (w *Widget) setPaintFocusLater(paintFocusLater bool) *Widget {
	w.paintFocusLater = paintFocusLater

	return w
}

func (w *Widget) paintChildren(context dom.CanvasRenderingContext2D) *Widget {
	if w.paintFocusLater {
		w.paintChildrenFocusLater(context)
	} else {
		w.paintChildrenDefault(context)
	}

	return w
}

func (w *Widget) paintChildrenDefault(context dom.CanvasRenderingContext2D) *Widget {
	for _, child := range w.children {
		child.draw(context)
	}

	return
}

func (w *Widget) paintChildrenFocusLater(context dom.CanvasRenderingContext2D) {
	var focusChild *Widget
	for _, child := range w.children {
		if child.state == STATE_OVER || child.state == STATE_ACTIVE {
			focusChild = child
		} else {
			child.draw(context)
		}
	}

	if focusChild {
		focusChild.draw(context)
	}

	return
}

func (w *Widget) ensureImages() {
	return
}

func (w *Widget) draw(context dom.CanvasRenderingContext2D) {
	if !w.visible {
		return
	}

	if w.checkEnable {
		w.setEnable(w.checkEnable())
	}

	w.ensureImages()

	context.save()
	w.relayout(context, false)

	context.Translate(w.rect.x, w.rect.y)
	w.beforePaint(context)
	w.paintBackground(context)
	w.paintSelf(context)
	w.paintChildren(context)
	w.drawInputTips(context)
	w.afterPaint(context)
	context.ClosePath()
	context.Restore()

	return
}

func (w *Widget) setVisible(visible bool) *Widget {
	w.visible = visible

	return w
}

func (w *Widget) isVisible() bool {
	return w.visible
}

func (w *Widget) onShow(visible bool) {
	return true
}

func (w *Widget) show(visible bool) {
	if visible != w.visible {
		w.visible = visible
		w.onShow(visible)
	}

	return w
}

func (w *Widget) showAll(visible bool) *Widget {
	w.show(visible)
	for _, child := range w.children {
		child.showAll(visible)
	}

	if !w.parent {
		w.postRedraw()
	}

	return w
}

func (w *Widget) selectAllChildren(selected bool) *Widget {
	for _, child := range w.children {
		if child.checkable {
			child.setChecked(selected)
		}
	}

	return w
}

func (w *Widget) closeWindow(retInfo interface{}) *Widget {
	w.getWindow().close(retInfo)

	return w
}

func (w *Widget) findTarget(point *Point) *Widget {
	p := w.getAbsPosition()
	w.point.x = point.x - p.x
	w.point.y = point.y - p.y

	for i := len(w.children) - 1; i >= 0; i-- {
		child := w.children[i]
		if !child.visible {
			continue
		}

		if isPointInRect(w.point, child.rect) {
			return child
		}
	}

	return nil
}

/////////////////////////////////////////////////////
func (w *Widget) onPointerDown(point *Point) bool {
	if !w.enable {
		return false
	}

	target := w.findTarget(point)
	if w.target && w.target != target {
		w.target.setState(STATE_NORMAL)
	}

	if target {
		target.setState(STATE_ACTIVE)
		target.onPointerDown(point)
	} else {
		w.changeCursor()
	}

	w.target = target
	w.postRedraw()

	return true
}

func (w *Widget) onPointerMove(point *Point) bool {
	if !w.enable {
		return false
	}

	var target *Widget
	if w.isPointerDown() {
		target = w.target
	} else {
		target = w.findTarget(point)
	}

	if w.target && target != w.target {
		w.target.setState(STATE_NORMAL, true)
	}

	if target {
		if w.isPointerDown() {
			target.setState(STATE_ACTIVE)
		} else {
			target.setState(STATE_OVER)
		}
	} else {
		w.changeCursor()
	}

	w.target = target
	w.postRedraw()

	return true
}

func (w *Widget) onPoingterUp(point *Point) bool {
	if !w.enable {
		return false
	}

	target := w.findTarget(point)
	if target && w.target != target {
		w.target.setState(STATE_NORMAL)
		w.target.onPoingterUp(point)
	}

	if target {
		target.setState(STATE_OVER)
		target.onPoingterUp(point)
	} else {
		w.changeCursor()
	}

	if w.isClicked() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		w.setState(STATE_ACTIVE)
		w.onClicked(point)
	}

	w.target = target
	w.postRedraw()

	return true
}

func (w *Widget) onKeyDown(code int) {
	if w.target {
		w.target.onKeyDown(code)
	}

	if w.keyDownHandler {
		w.keyDownHandler(code)
	}

	fmt.Printf("onKeyDown Widget:%s, code=%d ", w.t, code)

	return
}

func (w *Widget) onKeyUp(code int) {
	if w.target {
		w.target.onKeyUp(code)
	}

	if w.keyUpHandler {
		w.keyUpHandler(code)
	}

	fmt.Printf("onKeyUp Widget:%s, code=%d ", w.t, code)

	return
}

func (w *Widget) onWheel(delta float64) bool {
	if w.target {
		return w.target.onWheel(delta)
	}

	if w.wheelHandler {
		w.wheelHandler(delta)
	}

	return false
}

func (w *Widget) onDoubleClick(point *Point) {
	var target *Widget

	if win, ok := w.(*Window); ok && win.grabWidget {
		target = win.grabWidget
	} else {
		target = w.findTarget(point)
	}

	if target {
		target.onDoubleClick(point)
		w.target = target
	}

	if w.state != STATE_DISABLE && w.doubleClickedHandler {
		w.doubleClickedHandler(point)
	}

	return
}

func (w *Widget) onContextMenu(point *Point) {
	target := w.findTarget(point)

	if target {
		target.onContextMenu(point)
		w.target = target
	}

	if w.state != STATE_DISABLE && w.contextMenuHandler {
		w.contextMenuHandler(point)
	}

	return
}

func (w *Widget) onLongPress(point *Point) {
	target := w.findTarget(point)

	if target {
		target.onLongPress(point)
		w.target = target
	}

	if w.state != STATE_DISABLE && w.longPressHandler {
		w.longPressHandler(point)
	}

	return
}

func (w *Widget) setCursor(cursor string) *Widget {
	w.cursor = cursor

	return w
}

var canvasPool []dom.HTMLCanvasElement

func resizeCanvas(canvas dom.HTMLCanvasElement, w, h int) {
	canvas.Width = w
	canvas.Height = h
}

func getCanvas(x, y, w, h, zIndex int) dom.HTMLCanvasElement {
	var canvas dom.HTMLCanvasElement

	if len(canvasPool) {
		canvas = canvasPool[canvas.Length-1]
		canvasPool = append(canvasPool[:canvas.Length-1])
	} else {
		canvas = dom.Document.CreateElement("canvas")
	}

	resizeCanvas(canvas, w, h)
	canvas.Style().SetProperty("position", "absolute")
	canvas.Style().SetProperty("opacity", 1)
	canvas.Style().SetProperty("left", fmt.Sprintf("%dpx", x))
	canvas.Style().SetProperty("top", fmt.Sprintf("%dpx", y))
	canvas.Style().SetProperty("width", fmt.Sprintf("%dpx", w))
	canvas.Style().SetProperty("height", fmt.Sprintf("%dpx", h))
	canvas.Style().SetProperty("zIndex", zIndex)

	return
}

func putCanvas(canvas dom.HTMLCanvasElement) {
	canvas.Style().SetProperty("zIndex", -1)
	canvas.Style().SetProperty("opacity", 0)
	canvasPool = append(canvasPool, canvas)
}

var tipsCanvas dom.HTMLCanvasElement

func getTipsCanvas(x, y, w, h, zIndex int) dom.HTMLCanvasElement {
	if !tipsCanvas {
		tipsCanvas = getCanvas(x, y, w, h, zIndex)
		body := dom.Document.GetElementsByTagName("body")[0]
		body.AppendChild(tipsCanvas)
	}

	canvas := tipsCanvas

	canvas.Width = w
	canvas.Height = h
	canvas.Style().SetProperty("position", "absolute")
	canvas.Style().SetProperty("opacity", 1)
	canvas.Style().SetProperty("left", fmt.Sprintf("%dpx", x))
	canvas.Style().SetProperty("top", fmt.Sprintf("%dpx", y))
	canvas.Style().SetProperty("width", fmt.Sprintf("%dpx", w))
	canvas.Style().SetProperty("height", fmt.Sprintf("%dpx", h))
	canvas.Style().SetProperty("zIndex", zIndex)

	return canvas
}

func hideTipsCanvas() {
	if tipsCanvas {
		tipsCanvas.Style().SetProperty("zIndex", -1)
	}
}
