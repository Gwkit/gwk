package gwk

import (
	"fmt"
	"honnef.co/go/js/dom"
	"math"
	"time"
)

type BeforeDrawHandler func()

type WindowManager struct {
	window             *Window
	w, h               int
	canvas             *dom.HTMLCanvasElement
	app                *Application
	pointerDown        bool
	target             *Window
	drawCount          int
	requestCount       int
	startTime          time.Duration
	windows            []*Window
	grabWindows        []*Window
	eventLogging       bool
	pointerDownPoint   *Point
	lastPointerPoint   *Point
	enablePaint        bool
	beforeDrawHandlers []BeforeDrawHandler
	lastUpdateTime     time.Duration
	currentEvent       dom.Event
	xInputOffset       int
	yInputOffset       int
	xInputScale        float32
	yInputScale        float32
	maxFpsMode         bool
	shouldShowFPS      bool
	tipsWidget         *Widget
	needRedraw         bool
	ctx                *dom.CanvasRenderingContext2D
}

var manager = &WindowManager{}

func NewWindowManager(app *Application, canvas *dom.HTMLCanvasElement, eventElement *dom.HTMLElement) *WindowManager {
	SetEventsConsumer(manager, eventElement)
	return manager.init(app, canvas)
}

func GetWindowManagerInstance() *WindowManager {
	return manager
}

func (manager *WindowManager) init(app *Application, canvas *dom.HTMLCanvasElement) *WindowManager {
	manager.app = app
	manager.canvas = canvas
	manager.w = canvas.Width
	manager.h = canvas.Height

	return manager
}

func (manager *WindowManager) getApp() *Application {
	return manager.app
}

func (manager *WindowManager) onMultiTouch(action string, points []Point, event dom.Event) {
	for _, point := range points {
		manager.translatePoint(&point)
	}

	return
}

func (manager *WindowManager) preprocessEvent(t string, e dom.Event) bool {
	//TODO
	// manager.currentEvent = e.originalEvent ? e.originalEvent : e
	// manager.currentEvent = e.originalEvent ? e.originalEvent : e
	return true
}

func (manager *WindowManager) getCanvas() *dom.HTMLCanvasElement {
	return manager.canvas
}

func (manager *WindowManager) getWidget() int {
	return manager.w
}

func (manager *WindowManager) getHeight() int {
	return manager.h
}

func (manager *WindowManager) findTargetWin(point *Point) *Window {
	for _, window := range manager.grabWindows {
		if window.visible {
			return window
		}
	}

	for _, window := range manager.windows {
		if window.visible {
			if isPointInRect(point, window.rect) {
				return window
			}
		}
	}

	return nil
}

func (manager *WindowManager) resize(w, h int) {
	manager.w = w
	manager.h = h

	manager.postRedraw()

	return
}

func (manager *WindowManager) grab(window *Window) {
	manager.grabWindows = append(manager.grabWindows, window)

	return
}

func (manager *WindowManager) ungrab(window *Window) {
	for i, win := range manager.grabWindows {
		if window == win {
			manager.grabWindows = append(manager.grabWindows[:i], manager.grabWindows[i+1:]...)
			break
		}
	}

	return
}

func (manager *WindowManager) onDoubleClick(point *Point) {
	manager.translatePoint(point)
	manager.target = manager.findTargetWin(point)

	if manager.target != nil {
		manager.target.onDoubleClick(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d, y=%d", point.x, point.y)
	}

	return
}

func (manager *WindowManager) onLongPress(point *Point) {
	manager.target = manager.findTargetWin(point)

	if manager.target != nil {
		manager.target.onLongPress(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d, y=%d", point.x, point.y)
	}

	return
}

func (manager *WindowManager) setInputOffset(xInputOffset, yInputOffset int) {
	manager.xInputOffset = xInputOffset
	manager.yInputOffset = yInputOffset

	return
}

func (manager *WindowManager) setInputScale(xInputScale, yInputScale float32) {
	manager.xInputScale = xInputScale
	manager.yInputScale = yInputScale

	return
}

func (manager *WindowManager) getInputScale() (x, y float32) {
	return manager.xInputScale, manager.yInputScale
}

func (manager *WindowManager) translatePoint(point *Point) *Point {
	if manager.xInputOffset != 0 {
		point.x -= manager.xInputOffset
	}

	if manager.yInputOffset != 0 {
		point.y -= manager.yInputOffset
	}

	if manager.xInputScale != 0.0 {
		point.x = int(math.Ceil(float64(float32(point.x) * manager.xInputScale)))
	}

	if manager.yInputScale != 0.0 {
		point.y = int(math.Ceil(float64(float32(point.y) * manager.yInputScale)))
	}

	return point
}

func (manager *WindowManager) onPointerDown(point *Point) {
	manager.translatePoint(point)
	manager.target = manager.findTargetWin(point)

	for _, window := range manager.windows {
		if window.state == STATE_SELECTED && window != manager.target {
			window.setState(STATE_NORMAL)
		}
	}

	manager.pointerDown = true
	manager.pointerDownPoint.x = point.x
	manager.pointerDownPoint.y = point.y
	manager.lastPointerPoint.x = point.x
	manager.lastPointerPoint.y = point.y

	if manager.target != nil {
		manager.target.onPointerDown(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d y=%d", point.x, point.y)
	}

	return
}

func (manager *WindowManager) onPointerMove(point *Point) {
	manager.translatePoint(point)
	target := manager.findTargetWin(point)

	manager.lastPointerPoint.x = point.x
	manager.lastPointerPoint.y = point.y

	if manager.target != nil && manager.target != target {
		manager.target.onPointerMove(point)
	}
	manager.target = target
	if manager.target != nil {
		manager.target.onPointerMove(point)
	}

	return
}

func (manager *WindowManager) onPointerUp(point *Point) {
	manager.translatePoint(point)
	point = manager.lastPointerPoint
	manager.target = manager.findTargetWin(point)

	if manager.target != nil {
		manager.target.onPointerUp(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d y=%d", point.x, point.y)
	}
	manager.pointerDown = false

	return
}

func (manager *WindowManager) getLastPointerPoint() *Point {
	return manager.lastPointerPoint
}

func (manager *WindowManager) isPointerDown() bool {
	return manager.pointerDown
}

func (manager *WindowManager) isClicked() bool {
	dx := math.Abs(float64(manager.lastPointerPoint.x - manager.pointerDownPoint.x))
	dy := math.Abs(float64(manager.lastPointerPoint.y - manager.pointerDownPoint.y))

	return (dx < 10 && dy < 10)
}

func (manager *WindowManager) isCtrlDown() bool {
	return manager.currentEvent && manager.currentEvent.ctrlKey
}

func (manager *WindowManager) isAltDown() bool {
	return manager.currentEvent && manager.currentEvent.altKey
}

func (manager *WindowManager) onContextMenu(point *Point) {
	manager.target = manager.findTargetWin(point)

	if manager.target {
		manager.target.onContextMenu(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d y=%d", point.x, point.y)
	}

	return
}

func (manager *WindowManager) onKeyDown(code int) {
	if manager.target == nil {
		manager.target = manager.findTargetWin(&Point{x: 50, y: 50})
		fmt.Printf("onKeyDown findTargetWin=")
		fmt.Println(manager.target)
	}

	if manager.target != nil {
		manager.target.onKeyDown(code)
	}

	return
}

func (manager *WindowManager) onKeyUp(code int) {
	if manager.target != nil {
		manager.target.onKeyUp(code)
	}

	return
}

func (manager *WindowManager) onWheel(delta float64) bool {
	manager.postRedraw()

	if manager.target == nil {
		manager.target = manager.findTargetWin(&Point{x: 50, y: 50})
		fmt.Printf("onWheel findTargetWin=")
		fmt.Println(manager.target)
	}

	if manager.target != nil {
		return manager.target.onWheel(delta)
	}

	return false
}

func (manager *WindowManager) dispatchPointerMoveOut() *WindowManager {
	manager.onPointerMove(&Point{x: -1, y: -1})
	manager.target = nil

	return manager
}

func (manager *WindowManager) setTopWindowAsTarget() *WindowManager {
	manager.target = nil

	for _, window := range manager.windows {
		if window.visible {
			manager.target = window
			break
		}
	}

	return manager
}

func (manager *WindowManager) addWindow(win *Window) {
	manager.dispatchPointerMoveOut()
	manager.target = win
	manager.windows = append(manager.windows, win)
	manager.postRedraw()

	return
}

func (manager *WindowManager) removeWindow(win *Window) {
	manager.ungrab(win)

	if manager.target == win {
		manager.target = nil
	}

	for i, window := range manager.windows {
		if win == window {
			manager.windows = append(manager.windows[:i], manager.windows[i+1:]...)
			break
		}
	}
	manager.postRedraw()

	return
}

func (manager *WindowManager) getFrameRate() int {
	duration := time.Now() - manager.startTime
	fps := math.Floor(1000 * manager.drawCount / duration)

	if duration > 1000 {
		manager.drawCount = 0
		manager.startTime = time.Now()
	}

	return fps
}

func (manager *WindowManager) setMaxFPSMode(maxFpsMode bool) *WindowManager {
	manager.maxFpsMode = maxFpsMode

	return manager
}

func (manager *WindowManager) showFPS(shouldShowFPS bool) *WindowManager {
	manager.drawCount = 1
	manager.startTime = time.Now()
	manager.shouldShowFPS = shouldShowFPS

	return manager
}

func (manager *WindowManager) getPaintEnable() bool {
	return manager.enablePaint
}

func (manager *WindowManager) setPaintEnable(enablePaint bool) *WindowManager {
	manager.enablePaint = enablePaint
	fmt.Printf("setPaintEnable:%t", enablePaint)

	if manager.enablePaint {
		manager.postRedraw()
	}

	return manager
}

func (manager *WindowManager) onDrawFrame() {
	manager.drawCount++
	manager.requestCount = 0
	manager.draw()

	return
}

func (manager *WindowManager) postRedraw() {
	if !manager.enablePaint {
		return
	}

	manager.requestCount++
	if manager.requestCount < 2 {
		requestAnimationFrame(func() {
			manager.onDrawFrame()
		})
	}

	return
}

func (manager *WindowManager) setTipsWidget(widget *Widget) {
	manager.tipsWidget = widget

	return
}

func (manager *WindowManager) drawTips(context *dom.CanvasRenderingContext2D) {
	tipsWidget := manager.tipsWidget
	if !tipsWidget || !tipsWidget.parent {
		return
	}

	hideTipsCanvas()
	p := tipsWidget.getPositionInView()
	win := tipsWidget.getWindow()

	if win.canvas {
		context = win.canvas.getContext("2d")
		context.Save()
		context.Translate(p.x, p.y)
		context.BeginPath()
		tipsWidget.drawTips(context)
		context.Restore()
	} else {
		context.Save()
		context.Translate(p.x, p.y)
		context.BeginPath()
		tipsWidget.drawTips(context)
		context.Restore()
	}

	return
}

func (manager *WindowManager) beforeDrawWindows(context *dom.CanvasRenderingContext2D) {

}

func (manager *WindowManager) afterDrawWindows(context *dom.CanvasRenderingContext2D) {

}

func (manager *WindowManager) drawWindows(context *dom.CanvasRenderingContext2D) {
	manager.beforeDrawWindows(context)
	for _, window := range manager.windows {
		window.draw(context)
	}
	manager.drawTips(context)
	manager.afterDrawWindows(context)

	return
}

func (manager *WindowManager) checkNeedRedraw(timeStep int) bool {
	return true
}

func (manager *WindowManager) getCanvas2D() *dom.CanvasRenderingContext2D {
	if manager.ctx == nil {
		ctx := manager.canvas.GetContext2d("2d")
		rctx := reflect.ValueOf(ctx)
		if !rctx.FieldByName("BeginFrame") {
			ctx.BeginFrame = func() {}
		}
		if !rctx.FieldByName("EndFrame") {
			ctx.EndFrame = func() {}
		}

		if !rctx.FieldByName("clipRect") {
			ctx.ClipRect = func(x, y, w, h int) {
				ctx.BeginPath()
				ctx.Rect(x, y, w, h)
				ctx.Clip()
				ctx.BeginPath()
			}
		}
		manager.ctx = ctx
	}

	return manager.ctx
}

func (manager *WindowManager) doDraw(ctx *dom.CanvasRenderingContext2D) {
	now := time.Now()
	timeStep := now - (manager.lastUpdateTime || 0)

	if !manager.checkNeedRedraw(timeStep) {
		return
	}

	manager.needRedraw = 0
	ctx.Save()
	manager.drawWindows(ctx)
	ctx.Restore()

	if manager.shouldShowFPS {
		str := manager.getFrameRate()
		w, h := 100, 30
		ctx.BeginPath()
		ctx.Rect(0, 0, w, h)
		ctx.FillStyle = "Black"
		ctx.Fill()

		ctx.Save()
		ctx.TextAlign = "center"
		ctx.TextBaseline = "middle"
		ctx.Font = "20px Sans"
		ctx.FillStyle = "White"
		ctx.FillText(str, w>>1, h>>1)
		ctx.Restore()
	}

	if manager.maxFpsMode || manager.needRedraw > 0 {
		manager.postRedraw()
	}
	manager.lastUpdateTime = now

	return
}

func (manager *WindowManager) draw() {
	ctx := manager.getCanvas2D()

	ctx.BeginFrame()
	manager.doDraw(ctx)
	ctx.EndFrame()

	return
}
