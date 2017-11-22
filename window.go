package gwk

import (
	"github.com/Luncher/gwk/pkg/rt"
	"github.com/Luncher/gwk/pkg/structs"
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
	downPosition structs.Position
	upPosition   structs.Position
	lastPosition structs.Position
}

func NewWindow(manager *WindowManager, x, y, w, h float32) *Window {
	window := &Window{
		Widget: *NewWidget(nil, x, y, w, h),
		t:      TYPE_WINDOW,
	}

	if manager != nil {
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
	width, height := rt.GetRTInstance().GetViewPort()
	var sw = math.Min(float64(window.manager.w), float64(width))
	var sh = math.Min(float64(window.manager.h), float64(height))

	var x = (int(sw) - window.rect.W) / 2
	var y = (int(sh) - window.rect.H) / 2

	window.rect.X = x
	window.rect.Y = y

	return window
}

func (window *Window) onPointerDown(point *structs.Point) {
	window.pointerDown = true
	window.downPosition.X = point.X
	window.downPosition.Y = point.Y
	window.lastPosition.X = point.X
	window.lastPosition.Y = point.Y

	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onPointerDown(point)
	} else {
		window.Widget.onPointerDown(point)
	}

	window.PostRedraw()

	return
}

func (window *Window) onPointerMove(point *structs.Point) {
	window.lastPosition.X = point.X
	window.lastPosition.Y = point.Y

	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onPointerMove(point)
	} else {
		window.Widget.onPointerMove(point)
	}

	window.PostRedraw()

	return
}

func (window *Window) onPointerUp(point *structs.Point) {
	window.upPosition.X = point.X
	window.upPosition.Y = point.Y

	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onPointerUp(point)
	} else {
		window.Widget.onPointerUp(point)
	}
	window.pointerDown = false

	window.PostRedraw()

	return
}

func (window *Window) isClicked() bool {
	dx := window.lastPosition.X - window.downPosition.X
	dy := window.lastPosition.Y - window.downPosition.Y

	return math.Abs(float64(dx)) < 5 && math.Abs(float64(dy)) < 5
}

func (window *Window) onContextMenu(point *structs.Point) {
	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onContextMenu(point)
	} else {
		window.Widget.onContextMenu(point)
	}

	return
}

func (window *Window) onKeyDown(code int) {
	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onKeyDown(code)
	} else {
		window.Widget.onKeyDown(code)
	}

	return
}

func (window *Window) onKeyUp(code int) {
	if widget, ok := window.grabWidget.(Widget); ok {
		widget.onKeyUp(code)
	} else {
		window.Widget.onKeyUp(code)
	}

	return
}

func (window *Window) beforePaint(ctx *dom.CanvasRenderingContext2D) {
	ctx.BeginPath()
	ctx.Rect(0, 0, float64(window.rect.W), float64(window.rect.H))
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
