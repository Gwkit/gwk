package gwk

import (
	"fmt"
	"github.com/Luncher/gwk/pkg/rt"
	"github.com/Luncher/gwk/pkg/structs"
	"honnef.co/go/js/dom"
	"math"
	"time"
)

type WindowCloseHandler func()

type WindowHandler interface {
	show(visible bool) *Widget
	close(interface{}) *Widget
	onClose(interface{})
	onShow(visible bool)
}

type Window struct {
	*Widget
	grabWidget   *Widget
	closeHandler WindowCloseHandler
	manager      *WindowManager
	downPosition structs.Position
	upPosition   structs.Position
	lastPosition structs.Position
}

func NewWindow(manager *WindowManager, x, y, w, h float32) *Window {
	window := &Window{
		Widget: NewWidget(TYPE_WINDOW, nil, x, y, w, h),
	}
	window.I = window

	if manager != nil {
		window.manager = manager
	} else {
		window.manager = GetWindowManagerInstance()
	}

	time.AfterFunc(time.Millisecond*10, func() {
		manager.addWindow(window)
	})

	return window
}

func (w *Window) getWindow() *Window {
	return w
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

	if window.grabWidget != nil {
		window.grabWidget.onPointerDown(point)
	} else {
		window.Widget.onPointerDown(point)
	}

	fmt.Printf("window PostRedraw\n")
	window.PostRedraw()

	return
}

func (window *Window) onPointerMove(point *structs.Point) {
	window.lastPosition.X = point.X
	window.lastPosition.Y = point.Y

	fmt.Printf("window onPointerMove \n")
	if window.grabWidget != nil {
		window.grabWidget.onPointerMove(point)
	} else {
		window.Widget.onPointerMove(point)
	}

	window.PostRedraw()

	return
}

func (window *Window) onPointerUp(point *structs.Point) {
	window.upPosition.X = point.X
	window.upPosition.Y = point.Y

	if window.grabWidget != nil {
		window.grabWidget.onPointerUp(point)
	} else {
		window.Widget.onPointerUp(point)
	}
	window.pointerDown = false

	window.PostRedraw()

	return
}

func (window *Window) onDoubleClick(point *structs.Point) {
	if window.grabWidget != nil {
		window.grabWidget.onDoubleClick(point)
		window.target = window.grabWidget
		if window.state != STATE_DISABLE && window.doubleClickedHandler != nil {
			window.doubleClickedHandler(point)
		}
	}
}

func (window *Window) isClicked() bool {
	dx := window.lastPosition.X - window.downPosition.X
	dy := window.lastPosition.Y - window.downPosition.Y

	return math.Abs(float64(dx)) < 5 && math.Abs(float64(dy)) < 5
}

func (window *Window) onContextMenu(point *structs.Point) {
	if window.grabWidget != nil {
		window.grabWidget.onContextMenu(point)
	} else {
		window.Widget.onContextMenu(point)
	}

	return
}

func (window *Window) onKeyDown(code int) {
	if window.grabWidget != nil {
		window.grabWidget.onKeyDown(code)
	} else {
		window.Widget.onKeyDown(code)
	}

	return
}

func (window *Window) onKeyUp(code int) {
	if window.grabWidget != nil {
		window.grabWidget.onKeyUp(code)
	} else {
		window.Widget.onKeyUp(code)
	}

	return
}

func (window *Window) beforePaint(ctx *dom.CanvasRenderingContext2D) {
	fmt.Printf("window beforePaint\n")
	ctx.BeginPath()
	ctx.ClearRect(0, 0, float64(window.rect.W), float64(window.rect.H))

	return
}

func (window *Window) show(visible bool) {
	window.Widget.show(visible)
}

func (window *Window) close(retInfo interface{}) {
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
