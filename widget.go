package gwk

import (
	"fmt"
	"github.com/Luncher/gwk/pkg/image"
	"github.com/Luncher/gwk/pkg/structs"
	"github.com/Luncher/gwk/pkg/theme"
	"github.com/Luncher/gwk/pkg/utils"
	"honnef.co/go/js/dom"
	"math"
	"strconv"
)

const (
	STATE_NORMAL           = "state-normal"
	STATE_ACTIVE           = "state-active"
	STATE_OVER             = "state-over"
	STATE_DISABLE          = "state-disable"
	STATE_DISABLE_SELECTED = "state-diable-selected"
	STATE_SELECTED         = "state-selected"
	STATE_NORMAL_CURRENT   = "state-normal-current"
)

const (
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

const (
	BORDER_STYLE_NONE   = 0
	BORDER_STYLE_LEFT   = 1
	BORDER_STYLE_RIGHT  = 2
	BORDER_STYLE_TOP    = 4
	BORDER_STYLE_BOTTOM = 8
	BORDER_STYLE_ALL    = 0xffff
)

type PointerEventHandler interface {
	onPointerDown(*structs.Point)
	onPointerUp(*structs.Point)
	onDoubleClick(*structs.Point)
	onPointerMove(*structs.Point)
	onContextMenu(*structs.Point)
	onLongPress(*structs.Point)
	onWheel(delta float64) bool
}

type KeyEventHandler interface {
	onKeyDown(code int)
	onKeyUp(code int)
}

type PaintEventHandler interface {
	ensureImages()
	draw(*dom.CanvasRenderingContext2D)
	relayout(*dom.CanvasRenderingContext2D, bool)
	beforePaint(*dom.CanvasRenderingContext2D)
	paintBackground(*dom.CanvasRenderingContext2D)
	paintSelf(*dom.CanvasRenderingContext2D)
	paintChildren(*dom.CanvasRenderingContext2D)
	drawInputTips(*dom.CanvasRenderingContext2D)
	afterPaint(*dom.CanvasRenderingContext2D)
}

type WidgetInterface interface {
	KeyEventHandler
	PaintEventHandler
	PointerEventHandler
	getX() int
	getY() int
	destroy()
	getWindow() *Window
	showAll(visible bool) *Widget
	setState(string, bool) *Widget
	getParent() *Widget
	setVisible(visible bool) *Widget
	// SetText(text string, notify bool) *Widget
	findTargetWidgetEx(point *structs.Point, recursive bool) *Widget
}

type CheckEnable func() bool
type SetChecked func(bool, bool) *Widget
type RemovedHandler func()
type WidgetVisit func(*Widget)
type StateChangedHandler func(state string)
type OnMovedHandler func()
type OnResizedHandler func()
type ContextMenuHandler func(*structs.Point)
type ClickedHandler func(*Widget, *structs.Point)
type DoubleClickedHandler func(*structs.Point)
type LongPressHandler func(*structs.Point)
type KeyDownHandler func(int)
type KeyUpHandler func(int)
type OnBeforePaintHandler func(*dom.CanvasRenderingContext2D)
type OnAfterPaintHandler func(*dom.CanvasRenderingContext2D)
type WheelHandler func(float64)
type OnChangedHandler func(interface{})

type Widget struct {
	I                    WidgetInterface
	id                   string
	t                    string
	name                 string
	rect                 *structs.Rect
	pointerDown          bool
	visible              bool
	state                string
	parent               *Widget
	text                 string
	tag                  string
	tips                 string
	enable               bool
	checkable            bool
	children             []*Widget
	point                structs.Point
	cursor               string
	imageDisplay         image.Display
	borderStyle          int
	border               int
	themeType            string
	selected             bool
	selectable           bool
	needRelayout         bool
	isScrollView         bool
	xOffset              int
	yOffset              int
	inputTips            string
	leftMargin           int
	lineWidth            int
	roundRadius          int
	editing              bool
	paintFocusLater      bool
	userData             interface{}
	setChecked           SetChecked
	checkEnable          CheckEnable
	removedHandler       RemovedHandler
	target               *Widget
	onMoved              OnMovedHandler
	stateChangedHandler  StateChangedHandler
	onSized              OnResizedHandler
	contextMenuHandler   ContextMenuHandler
	longPressHandler     LongPressHandler
	keyUpHandler         KeyUpHandler
	keyDownHandler       KeyDownHandler
	clickedHandler       ClickedHandler
	doubleClickedHandler DoubleClickedHandler
	wheelHandler         WheelHandler
	theme                *theme.ThemeWidget
	onBeforePaint        OnBeforePaintHandler
	onAfterPaint         OnAfterPaintHandler
	onChanged            OnChangedHandler
}

func NewWidget(t string, parent *Widget, x, y, w, h float32) *Widget {
	widget := &Widget{
		t:            t,
		parent:       parent,
		cursor:       "default",
		visible:      true,
		enable:       true,
		borderStyle:  BORDER_STYLE_ALL,
		imageDisplay: image.DISPLAY_9PATCH,
		rect:         &structs.Rect{X: int(x), Y: int(y), W: int(w), H: int(h)},
	}

	widget.setState(STATE_NORMAL, false)

	if widget.parent != nil {
		var border int
		if border = 0; parent.border > 0 {
			border = parent.border
		}

		pw := parent.rect.W - 2*border
		ph := parent.rect.H - 2*border

		if x > 0 && x < 1 {
			widget.rect.X = int(float32(pw)*x + float32(border))
		}

		if w > 0 && w <= 1 {
			widget.rect.W = int(float32(pw) * w)
		}

		if y > 0 && y < 1 {
			widget.rect.Y = int(float32(ph)*y + float32(border))
		}

		if h > 0 && h <= 1 {
			widget.rect.H = int(float32(ph) * w)
		}

		parent.appendChild(widget)
	}

	return widget
}

func (w *Widget) UseTheme(t string) *Widget {
	w.themeType = t

	return w
}

func (w *Widget) isSelected() bool {
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
	return GetWindowManagerInstance().isPointerDown()
}

//FIXME
func (w *Widget) isClicked() bool {
	return GetWindowManagerInstance().isClicked()
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

func (w *Widget) getLastPointerPoint() *structs.Point {
	return GetWindowManagerInstance().getLastPointerPoint()
}

func (w *Widget) getTopWindow() *Window {
	return w.getWindow()
}

func (w *Widget) getWindow() *Window {
	return w.parent.getWindow()
}

func (w *Widget) getParent() *Widget {
	return w.parent
}

func (w *Widget) getX() int {
	return w.rect.X
}

func (w *Widget) getY() int {
	return w.rect.Y
}

func (w *Widget) getWidth() int {
	return w.rect.W
}

func (w *Widget) getHeight() int {
	return w.rect.H
}

func (w *Widget) getPositionInView() *structs.Point {
	x := w.getX()
	y := w.getY()
	point := &structs.Point{}

	for iter := w.getParent(); iter != nil; iter = iter.getParent() {
		x += iter.getX()
		y += iter.getY()
		if iter.isScrollView {
			x -= iter.xOffset
			y -= iter.yOffset
		}
	}

	point.X = x
	point.Y = y

	return point
}

func (w *Widget) getAbsPosition() *structs.Point {
	x := w.rect.X
	y := w.rect.Y

	for parent := w.parent; parent != nil; parent = parent.parent {
		x += parent.getX()
		y += parent.getY()
	}

	return &structs.Point{X: x, Y: y}
}

func (w *Widget) getPositionInWindow() *structs.Point {
	point := &structs.Point{}

	if w.parent != nil {
		for iter := w; iter != nil; iter = iter.parent {
			if iter.parent == nil {
				break
			}

			point.X += iter.rect.X
			point.Y += iter.rect.Y
		}
	}

	return point
}

func (w *Widget) translatePoint(point *structs.Point) *structs.Point {
	p := w.getAbsPosition()

	return &structs.Point{X: point.X - p.X, Y: point.Y - p.Y}
}

func (w *Widget) postRedrawAll() {
	GetWindowManagerInstance().postRedraw()

	return
}

func (w *Widget) PostRedraw() {
	GetWindowManagerInstance().postRedraw()

	return
}

func (w *Widget) redraw(rect *structs.Rect) {
	// p := w.getAbsPosition()

	// if rect == nil {
	// 	rect = &structs.Rect{X: 0, Y: 0, W: w.rect.W, H: w.rect.H}
	// }
	// rect.X = p.X + rect.X
	// rect.Y = p.Y + rect.Y

	// GetWindowManagerInstance().redraw(rect)

	return
}

func (w *Widget) isPointIn(point *structs.Point) bool {
	return isPointInRect(point, w.rect)
}

func (w *Widget) findTargetWidgetEx(point *structs.Point, recursive bool) *Widget {
	if !w.visible || !w.isPointIn(point) {
		return nil
	}

	if recursive && len(w.children) > 0 {
		n := len(w.children) - 1
		p := w.point
		p.X = point.X - w.rect.X
		p.Y = point.Y - w.rect.Y

		for i := n; i > 0; i-- {
			iter := w.children[i]
			ret := iter.findTargetWidgetEx(&p, false)
			if ret != nil {
				return ret
			}
		}
	}

	return w
}

func (w *Widget) findTargetWidget(point *structs.Point) *Widget {
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

		if t := parent.target; t == w {
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

func (w *Widget) forEachChild(onVisit WidgetVisit) {
	for _, child := range w.children {
		onVisit(child)
	}

	return
}

func (w *Widget) setTextOf(name, text string, notify bool) *Widget {
	child := w.lookup(name, true)

	if child != nil {
		child.SetText(text, notify)
	} else {
		fmt.Printf("not found %s", name)
	}

	return child
}

func (w *Widget) setVisibleOf(name string, value bool) *Widget {
	child := w.lookup(name, true)

	if child != nil {
		child.setVisible(value)
	} else {
		fmt.Printf("not found %s", name)
	}

	return child
}

func (w *Widget) SetText(text string, notify bool) *Widget {
	w.text = text
	w.setNeedRelayout(true)

	return w
}

func (w *Widget) GetText() string {
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
	h := widget.rect.H
	w := widget.rect.W
	y := widget.rect.H >> 1
	x := widget.leftMargin
	text := widget.GetText()
	inputTips := widget.getInputTips()

	if len(text) > 0 || len(inputTips) == 0 || widget.t != TYPE_EDIT || widget.editing {
		return
	}

	style := widget.getStyle("")
	context.Save()
	context.Font = style.Font
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
		style := w.getStyle("")
		x := w.rect.W >> 1
		y := w.rect.H >> 1
		font := style.Font
		textColor := style.TextColor

		if len(font) > 0 && len(textColor) > 0 {
			context.TextAlign = "center"
			context.TextBaseline = "middle"
			context.Font = font
			context.FillStyle = textColor
			context.FillText(tips, float64(x), float64(y), -1)
		}
	}

	return
}

func (w *Widget) setID(id string) *Widget {
	w.id = id

	return w
}

func (w *Widget) getID() string {
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
		w.target.setState(state, recursive)
	}

	return w
}

func (w *Widget) move(x, y int) *Widget {
	w.rect.X = x
	w.rect.Y = y
	if w.onMoved != nil {
		w.onMoved()
	}

	return w
}

func (w *Widget) MoveToCenter(moveX, moveY bool) *Widget {
	parent := w.parent
	pw := parent.rect.W
	ph := parent.rect.H

	if moveX {
		w.rect.X = (pw - w.rect.W) >> 1
	}

	if moveY {
		w.rect.Y = (ph - w.rect.H) >> 1
	}

	return w
}

func (w *Widget) moveToBottom(border int) *Widget {
	ph := w.parent.rect.H
	w.rect.Y = ph - w.rect.H - border

	return w
}

func (w *Widget) moveDelta(dx, dy int) *Widget {
	w.rect.X = w.rect.X + dx
	w.rect.Y = w.rect.Y + dy
	if w.onMoved != nil {
		w.onMoved()
	}

	return w
}

func (widget *Widget) resize(w, h int) *Widget {
	widget.rect.W = w
	widget.rect.H = h
	if widget.onSized != nil {
		widget.onSized()
	}
	widget.setNeedRelayout(true)

	return widget
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

func (w *Widget) onClicked(point *structs.Point) bool {
	if w.clickedHandler != nil {
		w.clickedHandler(w, point)
	}

	w.PostRedraw()

	return w.clickedHandler != nil
}

func (w *Widget) lookup(id string, recursive bool) *Widget {
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

func (w *Widget) onRelayout(context *dom.CanvasRenderingContext2D, force bool) {

}

func (w *Widget) relayout(context *dom.CanvasRenderingContext2D, force bool) {
	if !w.needRelayout || !force || len(w.children) == 0 {
		return
	}

	w.onRelayout(context, force)
	w.needRelayout = false

	return
}

func (w *Widget) setLineWidth(lineWidth int) *Widget {
	w.lineWidth = lineWidth

	return w
}

func (w *Widget) getLineWidth(style *theme.ThemeStyle) int {
	if w.lineWidth > 0 {
		return w.lineWidth
	}
	return 0
}

func (w *Widget) setRoundRadius(roundRadius int) *Widget {
	w.roundRadius = roundRadius

	return w
}

func (w *Widget) ensureTheme() *Widget {
	if len(w.themeType) > 0 {
		w.theme = theme.Get(w.themeType, false)
	} else {
		w.theme = theme.Get(w.t, false)
	}

	return w
}

func (w *Widget) getStyle(_state string) *theme.ThemeStyle {
	var style *theme.ThemeStyle
	w.ensureTheme()
	var state string
	if state = _state; len(state) == 0 {
		state = w.state
	}

	if !w.enable {
		if w.selectable && w.isSelected() {
			style = w.theme.StateSelected
		} else {
			style = w.theme.StateDisable
		}
	} else {
		if w.selectable && w.selected {
			style = w.theme.StateSelected
		} else if state == STATE_OVER {
			style = w.theme.StateOver
		} else if state == STATE_ACTIVE {
			style = w.theme.StateActive
		} else {
			style = w.theme.StateNormal
		}
	}

	if style != nil {
		style = w.theme.StateNormal
	}

	return style
}

func (w *Widget) setImageDisplay(imageDisplay image.Display) *Widget {
	w.imageDisplay = imageDisplay

	return w
}

func (w *Widget) setBorderStyle(borderStyle int) *Widget {
	w.borderStyle = borderStyle

	return w
}

func (w *Widget) paintBackground(context *dom.CanvasRenderingContext2D) {
	style := w.getStyle("")
	if style != nil {
		if style.BgImage != nil {
			w.paintBackgroundImage(context, style)
		} else {
			w.paintBackgroundColor(context, style)
		}
	}
}

func (w *Widget) paintBackgroundImage(context *dom.CanvasRenderingContext2D, style *theme.ThemeStyle) {
	dst := w.rect
	bgImage := style.BgImage
	imageDisplay := w.imageDisplay

	if bgImage.GetImage() != nil {
		var topOut, leftOut, rightOut, bottomOut int
		x := -leftOut
		y := topOut
		w := dst.W + rightOut + leftOut
		h := dst.H + bottomOut + topOut

		bgImage.Draw(context, imageDisplay, x, y, w, h)
	}

	return
}

func (widget *Widget) paintLeftBorder(context *dom.CanvasRenderingContext2D, w, h int) {
	context.BeginPath()
	context.MoveTo(0, 0)
	context.LineTo(0, float64(h))
	context.Stroke()
}

func (widget *Widget) paintRightBorder(context *dom.CanvasRenderingContext2D, w, h int) {
	context.BeginPath()
	context.MoveTo(float64(w), 0)
	context.LineTo(float64(w), float64(h))
	context.Stroke()
}

func (widget *Widget) paintTopBorder(context *dom.CanvasRenderingContext2D, w, h int) {
	context.BeginPath()
	context.MoveTo(0, 0)
	context.LineTo(float64(w), 0)
	context.Stroke()
}

func (widget *Widget) paintBottomBorder(context *dom.CanvasRenderingContext2D, w, h int) {
	context.BeginPath()
	context.MoveTo(0, float64(h))
	context.LineTo(float64(w), float64(h))
	context.Stroke()
}

func (w *Widget) paintBackgroundColor(context *dom.CanvasRenderingContext2D, style *theme.ThemeStyle) {
	dst := w.rect
	context.BeginPath()
	if w.roundRadius != 0 {
		roundRadius := math.Min(float64((dst.H>>1)-1), float64(w.roundRadius))
		utils.DrawRoundRect(context, float64(dst.W), float64(dst.H), roundRadius, 0)
	} else {
		context.Rect(0, 0, float64(dst.W), float64(dst.H))
	}

	if style.FillColor != "" {
		context.FillStyle = style.FillColor
		context.Fill()
	}

	lineWidth := w.getLineWidth(style)
	if lineWidth > 0 || style.LineColor != "" || w.borderStyle == BORDER_STYLE_NONE {
		context.BeginPath()
		return
	}

	width := w.getWidth()
	height := w.getHeight()
	context.LineWidth = lineWidth
	context.StrokeStyle = style.LineColor
	if w.borderStyle == BORDER_STYLE_ALL {
		context.Stroke()
		context.BeginPath()
		return
	}

	if w.borderStyle&BORDER_STYLE_LEFT != 0 {
		w.paintLeftBorder(context, width, height)
	}

	if w.borderStyle&BORDER_STYLE_RIGHT != 0 {
		w.paintRightBorder(context, width, height)
	}

	if w.borderStyle&BORDER_STYLE_TOP != 0 {
		w.paintTopBorder(context, width, height)
	}

	if w.borderStyle&BORDER_STYLE_BOTTOM != 0 {
		w.paintBottomBorder(context, width, height)
	}
	context.BeginPath()

	return
}

func (w *Widget) paintSelf(context *dom.CanvasRenderingContext2D) {
	return
}

func (w *Widget) beforePaint(context *dom.CanvasRenderingContext2D) {
	if w.onBeforePaint != nil {
		w.onBeforePaint(context)
	}
	return
}

func (w *Widget) afterPaint(context *dom.CanvasRenderingContext2D) {
	if w.onAfterPaint != nil {
		w.onAfterPaint(context)
	}
	return
}

func (w *Widget) setPaintFocusLater(paintFocusLater bool) *Widget {
	w.paintFocusLater = paintFocusLater

	return w
}

func (w *Widget) paintChildren(context *dom.CanvasRenderingContext2D) {
	fmt.Printf("%s paintChildren\n", w.t)
	if w.paintFocusLater {
		w.paintChildrenFocusLater(context)
	} else {
		w.paintChildrenDefault(context)
	}

	return
}

func (w *Widget) paintChildrenDefault(context *dom.CanvasRenderingContext2D) {
	for _, child := range w.children {
		child.I.draw(context)
	}

	return
}

func (w *Widget) paintChildrenFocusLater(context *dom.CanvasRenderingContext2D) {
	var focusChild *Widget
	for _, child := range w.children {
		if child.state == STATE_OVER || child.state == STATE_ACTIVE {
			focusChild = child
		} else {
			child.I.draw(context)
		}
	}

	if focusChild != nil {
		focusChild.I.draw(context)
	}

	return
}

func (w *Widget) ensureImages() {
	return
}

func (w *Widget) draw(context *dom.CanvasRenderingContext2D) {
	if !w.visible {
		return
	}

	if w.checkEnable != nil {
		w.setEnable(w.checkEnable())
	}

	w.I.ensureImages()

	context.Save()
	w.I.relayout(context, false)

	context.Translate(float64(w.rect.X), float64(w.rect.Y))
	w.I.beforePaint(context)
	w.I.paintBackground(context)
	w.I.paintSelf(context)
	w.I.paintChildren(context)
	w.I.drawInputTips(context)
	w.I.afterPaint(context)
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

func (w *Widget) onShow(visible bool) bool {
	return true
}

func (w *Widget) show(visible bool) *Widget {
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

	if w.parent != nil {
		w.PostRedraw()
	}

	return w
}

func (w *Widget) selectAllChildren(selected bool) *Widget {
	for _, child := range w.children {
		if child.checkable {
			child.setChecked(selected, false)
		}
	}

	return w
}

func (w *Widget) closeWindow(retInfo interface{}) *Widget {
	//TODO
	// w.getWindow().close(retInfo)

	return w
}

func (w *Widget) findTarget(point *structs.Point) *Widget {
	p := w.getAbsPosition()
	w.point.X = point.X - p.X
	w.point.Y = point.Y - p.Y

	for i := len(w.children) - 1; i >= 0; i-- {
		child := w.children[i]
		if !child.visible {
			continue
		}

		if isPointInRect(&w.point, child.rect) {
			return child
		}
	}

	return nil
}

/////////////////////////////////////////////////////
func (w *Widget) onPointerDown(point *structs.Point) {
	if !w.enable {
		return
	}

	target := w.findTarget(point)
	if w.target != nil && w.target != target {
		w.target.setState(STATE_NORMAL, false)
	}

	if target != nil {
		target.setState(STATE_ACTIVE, false)
		target.onPointerDown(point)
	} else {
		w.changeCursor()
	}

	w.target = target
	w.PostRedraw()

	return
}

func (w *Widget) onPointerMove(point *structs.Point) {
	if !w.enable {
		return
	}

	var target *Widget
	if w.isPointerDown() {
		target = w.target
	} else {
		target = w.findTarget(point)
	}

	if w.target != nil && target != w.target {
		w.target.setState(STATE_NORMAL, true)
	}

	if target != nil {
		if w.isPointerDown() {
			target.setState(STATE_ACTIVE, false)
		} else {
			target.setState(STATE_OVER, false)
		}
	} else {
		w.changeCursor()
	}

	w.target = target
	w.PostRedraw()

	return
}

func (w *Widget) onPointerUp(point *structs.Point) {
	if !w.enable {
		return
	}

	target := w.findTarget(point)
	if target != nil && w.target != target {
		w.target.setState(STATE_NORMAL, false)
		w.target.onPointerUp(point)
	}

	if target != nil {
		target.setState(STATE_OVER, false)
		target.onPointerUp(point)
	} else {
		w.changeCursor()
	}

	if w.isClicked() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		w.setState(STATE_ACTIVE, false)
		w.onClicked(point)
	}

	w.target = target
	w.PostRedraw()

	return
}

func (w *Widget) onKeyDown(code int) {
	if w.target != nil {
		w.target.onKeyDown(code)
	}

	if w.keyDownHandler != nil {
		w.keyDownHandler(code)
	}

	fmt.Printf("onKeyDown Widget:%s, code=%d ", w.t, code)

	return
}

func (w *Widget) onKeyUp(code int) {
	if w.target != nil {
		w.target.onKeyUp(code)
	}

	if w.keyUpHandler != nil {
		w.keyUpHandler(code)
	}

	fmt.Printf("onKeyUp Widget:%s, code=%d ", w.t, code)

	return
}

func (w *Widget) onWheel(delta float64) bool {
	if w.target != nil {
		return w.target.onWheel(delta)
	}

	if w.wheelHandler != nil {
		w.wheelHandler(delta)
	}

	return false
}

func (w *Widget) onDoubleClick(point *structs.Point) {
	target := w.findTarget(point)
	if target != nil {
		target.onDoubleClick(point)
		w.target = target
	}

	if w.state != STATE_DISABLE && w.doubleClickedHandler != nil {
		w.doubleClickedHandler(point)
	}

	return
}

func (w *Widget) onContextMenu(point *structs.Point) {
	target := w.findTarget(point)

	if target != nil {
		target.onContextMenu(point)
		w.target = target
	}

	if w.state != STATE_DISABLE && w.contextMenuHandler != nil {
		w.contextMenuHandler(point)
	}

	return
}

func (w *Widget) onLongPress(point *structs.Point) {
	target := w.findTarget(point)

	if target != nil {
		target.onLongPress(point)
		w.target = target
	}

	if w.state != STATE_DISABLE && w.longPressHandler != nil {
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

	if len(canvasPool) != 0 {
		canvas = canvasPool[len(canvasPool)-1]
		canvasPool = canvasPool[:len(canvasPool)-1]
	} else {
		document := dom.GetWindow().Document()
		canvas = document.CreateElement("canvas").(dom.HTMLCanvasElement)
	}

	resizeCanvas(canvas, w, h)
	canvas.Style().SetProperty("position", "absolute", "")
	canvas.Style().SetProperty("opacity", fmt.Sprintf("%d", 1), "")
	canvas.Style().SetProperty("left", fmt.Sprintf("%dpx", x), "")
	canvas.Style().SetProperty("top", fmt.Sprintf("%dpx", y), "")
	canvas.Style().SetProperty("width", fmt.Sprintf("%dpx", w), "")
	canvas.Style().SetProperty("height", fmt.Sprintf("%dpx", h), "")
	canvas.Style().SetProperty("zIndex", fmt.Sprintf("%d", zIndex), "")

	return canvas
}

func putCanvas(canvas dom.HTMLCanvasElement) {
	canvas.Style().SetProperty("zIndex", strconv.Itoa(-1), "")
	canvas.Style().SetProperty("opacity", strconv.Itoa(0), "")
	canvasPool = append(canvasPool, canvas)
}

var tipsCanvas *dom.HTMLCanvasElement

func getTipsCanvas(x, y, w, h, zIndex int) dom.HTMLCanvasElement {
	if tipsCanvas == nil {
		canvas := getCanvas(x, y, w, h, zIndex)
		tipsCanvas = &canvas
		document := dom.GetWindow().Document()
		body := document.GetElementsByTagName("body")[0]
		body.AppendChild(tipsCanvas)
	}

	canvas := *tipsCanvas

	canvas.Width = w
	canvas.Height = h
	canvas.Style().SetProperty("position", "absolute", "")
	canvas.Style().SetProperty("opacity", strconv.Itoa(1), "")
	canvas.Style().SetProperty("left", fmt.Sprintf("%dpx", x), "")
	canvas.Style().SetProperty("top", fmt.Sprintf("%dpx", y), "")
	canvas.Style().SetProperty("width", fmt.Sprintf("%dpx", w), "")
	canvas.Style().SetProperty("height", fmt.Sprintf("%dpx", h), "")
	canvas.Style().SetProperty("zIndex", strconv.Itoa(zIndex), "")

	return canvas
}

func hideTipsCanvas() {
	if tipsCanvas != nil {
		tipsCanvas.Style().SetProperty("zIndex", strconv.Itoa(-1), "")
	}
}
