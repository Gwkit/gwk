package gwk

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

type Widget struct {
	rect        *Rect
	pointerDown bool
	visible     bool
	state       string
}

func (w *Widget) onPointerDown(point *Point) {

}

func (w *Widget) onPointerMove(point *Point) {

}

func (w *Widget) onPointerUp(point *Point) {

}

func (w *Widget) onContextMenu(point *Point) {

}

func (w *Widget) onKeyDown(code string) {

}

func (w *Widget) onKeyUp(code string) {

}

func (w *Widget) show(visible bool) {

}

func (w *Widget) postRedraw() {

}

func (w *Widget) destroy() {

}

func (w *Widget) onDoubleClick(point *Point) {

}

func (w *Widget) onLongPress(point *Point) {

}

func (w *Widget) setState(state string) {

}

type Position struct {
	x int
	y int
}

func NewPosition(x, y int) *Position {
	return &Position{x, y}
}

type Rect struct {
	x, y, w, h int
}

func NewRect(x, y, w, h int) *Rect {
	return &Rect{x, y, w, h}
}

type Point struct {
	x, y int
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}
