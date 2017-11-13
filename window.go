package gwk

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
	"math"
)

type CloseHandler func()
type OnClosed func(info interface{})

type Window struct {
	Widget
	t            string
	grabWidget   interface{}
	onClosed     OnClosed
	closeHandler CloseHandler
	manager      *WindowManager
	downPosition Position
	upPosition   Position
	lastPosition Position
}

func NewWindow(manager *WindowManager, x, y, w, h float32) *Window {
	window := &Window{
		Widget: NewWidget(nil, x, y, w, h),
		t:      TYPE_WINDOW,
	}

	if manager {
		window.manager = manager
	} else {
		window.manager = GetWindowManagerInstance()
	}

	return window
}

func (window *Window) grab(widget *Widget) *Window {
	window.grabWidget = widget
	window.manager.grab(window)

	return window
}

func (window *Window) ungrab() *Window {
	window.grabWidget = nil
	return window
}

func (window *Window) moveToCenter() *Window {
	width, height := GetRTInstance().GetViewPort()
	var sw = math.Min(float64(window.manager.w), float64(width))
	var sh = math.Min(float64(window.manager.h), float64(height))

	var x = (int(sw) - window.rect.w) / 2
	var y = (int(sh) - window.rect.h) / 2

	window.rect.x = x
	window.rect.y = y

	return window
}

func (window *Window) onPointerDown(point *Point) {
	window.pointerDown = true
	window.downPosition.x = point.x
	window.downPosition.y = point.y
	window.lastPosition.x = point.x
	window.lastPosition.y = point.y

	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onPointerDown(point)
	} else {
		window.Widget.onPointerDown(point)
	}

	window.postRedraw()

	return
}

func (window *Window) onPointerMove(point *Point) {
	window.lastPosition.x = point.x
	window.lastPosition.y = point.y

	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onPointerMove(point)
	} else {
		window.Widget.onPointerMove(point)
	}

	window.postRedraw()

	return
}

func (window *Window) onPointerUp(point *Point) {
	window.upPosition.x = point.x
	window.upPosition.y = point.y

	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onPointerUp(point)
	} else {
		window.Widget.onPointerUp(point)
	}
	window.pointerDown = false

	window.postRedraw()

	return
}

func (window *Window) isClicked() bool {
	dx := window.lastPosition.x - window.downPosition.x
	dy := window.lastPosition.y - window.downPosition.y

	return math.Abs(float64(dx)) < 5 && math.Abs(float64(dy)) < 5
}

func (window *Window) onContextMenu(point *Point) {
	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onContextMenu(point)
	} else {
		window.Widget.onContextMenu(point)
	}

	return
}

func (window *Window) onKeyDown(code string) {
	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onKeyDown(code)
	} else {
		window.Widget.onKeyDown(code)
	}

	return
}

func (window *Window) onKeyUp(code string) {
	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onKeyUp(code)
	} else {
		window.Widget.onKeyUp(code)
	}

	return
}

func (window *Window) beforePaint(ctx *dom.CanvasRenderingContext2D) {
	ctx.BeginPath()
	ctx.Rect(0, 0, float64(window.rect.w), float64(window.rect.h))
	ctx.Clip()
	ctx.BeginPath()
}

func (window *Window) show(visible bool) {
	window.Widget.show(visible)
}

func (window *Window) close(retInfo interface{}) {
	if window.onClosed != nil {
		window.onClosed(retInfo)
	}

	if window.closeHandler != nil {
		window.closeHandler()
	}

	window.manager.ungrab(window)
	window.manager.removeWindow(window)
	window.destroy()

	return
}

func (window *Window) getCanvas2D() *dom.CanvasRenderingContext2D {
	return GetWindowManagerInstance().getCanvas2D()
}

func (window *Window) getCanvas() *dom.HTMLCanvasElement {
	return GetWindowManagerInstance().getCanvas()
}
