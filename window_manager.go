package gwk

import (
	"fmt"
	"github.com/Luncher/gwk/pkg/event"
	"github.com/Luncher/gwk/pkg/structs"
	"honnef.co/go/js/dom"
	"math"
	"strconv"
	"time"
)

type BeforeDrawHandler func()

type WindowManager struct {
	window             *Window
	w, h               int
	canvas             dom.HTMLCanvasElement
	app                *Application
	pointerDown        bool
	target             *Window
	drawCount          int
	requestCount       int
	startTime          time.Time
	windows            []*Window
	grabWindows        []*Window
	eventLogging       bool
	pointerDownPoint   structs.Point
	lastPointerPoint   structs.Point
	enablePaint        bool
	beforeDrawHandlers []BeforeDrawHandler
	lastUpdateTime     time.Time
	currentEvent       dom.Event
	xInputOffset       int
	yInputOffset       int
	xInputScale        float32
	yInputScale        float32
	maxFpsMode         bool
	shouldShowFPS      bool
	tipsWidget         *Widget
	needRedraw         int
	ctx                *dom.CanvasRenderingContext2D
}

var manager = &WindowManager{}

func NewWindowManager(app *Application, canvas dom.HTMLCanvasElement, eventElement dom.HTMLCanvasElement) *WindowManager {
	event.SetEventsConsumer(event.EventConsumer(manager), eventElement)
	return manager.init(app, canvas)
}

func GetWindowManagerInstance() *WindowManager {
	return manager
}

func (manager *WindowManager) init(app *Application, canvas dom.HTMLCanvasElement) *WindowManager {
	manager.app = app
	manager.canvas = canvas
	manager.w = canvas.Width
	manager.h = canvas.Height
	manager.enablePaint = true

	return manager
}

func (manager *WindowManager) getApp() *Application {
	return manager.app
}

func (manager *WindowManager) onMultiTouch(action string, points []structs.Point, event dom.Event) {
	for _, point := range points {
		manager.translatePoint(&point)
	}

	return
}

func (manager *WindowManager) PreprocessEvent(t string, e dom.Event) bool {
	//TODO
	// manager.currentEvent = e.originalEvent ? e.originalEvent : e
	// manager.currentEvent = e.originalEvent ? e.originalEvent : e
	return true
}

func (manager *WindowManager) GetInputScale() (x, y float32) {
	return manager.xInputScale, manager.yInputScale
}

func (manager *WindowManager) getCanvas() *dom.HTMLCanvasElement {
	return &manager.canvas
}

func (manager *WindowManager) getWidget() int {
	return manager.w
}

func (manager *WindowManager) getHeight() int {
	return manager.h
}

func (manager *WindowManager) findTargetWin(point *structs.Point) *Window {
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

func (manager *WindowManager) OnDoubleClick(point *structs.Point) {
	manager.translatePoint(point)
	manager.target = manager.findTargetWin(point)

	if manager.target != nil {
		manager.target.onDoubleClick(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d, y=%d", point.X, point.Y)
	}

	return
}

func (manager *WindowManager) onLongPress(point *structs.Point) {
	manager.target = manager.findTargetWin(point)

	if manager.target != nil {
		manager.target.onLongPress(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d, y=%d", point.X, point.Y)
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

func (manager *WindowManager) translatePoint(point *structs.Point) *structs.Point {
	if manager.xInputOffset != 0 {
		point.X -= manager.xInputOffset
	}

	if manager.yInputOffset != 0 {
		point.Y -= manager.yInputOffset
	}

	if manager.xInputScale != 0.0 {
		point.X = int(math.Ceil(float64(float32(point.X) * manager.xInputScale)))
	}

	if manager.yInputScale != 0.0 {
		point.Y = int(math.Ceil(float64(float32(point.Y) * manager.yInputScale)))
	}

	return point
}

func (manager *WindowManager) OnPointerDown(point *structs.Point) {
	manager.translatePoint(point)
	manager.target = manager.findTargetWin(point)

	for _, window := range manager.windows {
		if window.state == STATE_SELECTED && window != manager.target {
			window.setState(STATE_NORMAL, false)
		}
	}

	manager.pointerDown = true
	manager.pointerDownPoint.X = point.X
	manager.pointerDownPoint.Y = point.Y
	manager.lastPointerPoint.X = point.X
	manager.lastPointerPoint.Y = point.Y

	if manager.target != nil {
		manager.target.onPointerDown(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d y=%d\n", point.X, point.Y)
	}

	return
}

func (manager *WindowManager) OnPointerMove(point *structs.Point) {
	manager.translatePoint(point)
	target := manager.findTargetWin(point)

	manager.lastPointerPoint.X = point.X
	manager.lastPointerPoint.Y = point.Y

	if manager.target != nil && manager.target != target {
		manager.target.onPointerMove(point)
	}
	manager.target = target
	if manager.target != nil {
		manager.target.onPointerMove(point)
	}

	return
}

func (manager *WindowManager) OnPointerUp(point *structs.Point) {
	manager.translatePoint(point)
	point = &manager.lastPointerPoint
	manager.target = manager.findTargetWin(point)

	if manager.target != nil {
		manager.target.onPointerUp(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d y=%d\n", point.X, point.Y)
	}
	manager.pointerDown = false

	return
}

func (manager *WindowManager) getLastPointerPoint() *structs.Point {
	return &manager.lastPointerPoint
}

func (manager *WindowManager) isPointerDown() bool {
	return manager.pointerDown
}

func (manager *WindowManager) isClicked() bool {
	dx := math.Abs(float64(manager.lastPointerPoint.X - manager.pointerDownPoint.X))
	dy := math.Abs(float64(manager.lastPointerPoint.Y - manager.pointerDownPoint.Y))

	return (dx < 10 && dy < 10)
}

func (manager *WindowManager) isCtrlDown() bool {
	if manager.currentEvent != nil {
		if keyEvent, ok := manager.currentEvent.(dom.KeyboardEvent); ok {
			return keyEvent.CtrlKey
		}
	}
	return false
}

func (manager *WindowManager) isAltDown() bool {
	if manager.currentEvent != nil {
		if keyEvent, ok := manager.currentEvent.(dom.KeyboardEvent); ok {
			return keyEvent.AltKey
		}
	}
	return false
}

func (manager *WindowManager) OnContextMenu(point *structs.Point) {
	manager.target = manager.findTargetWin(point)

	if manager.target != nil {
		manager.target.onContextMenu(point)
	} else {
		fmt.Printf("Window Manager: no target for x=%d y=%d\n", point.X, point.Y)
	}

	return
}

func (manager *WindowManager) OnKeyDown(code int) {
	if manager.target == nil {
		manager.target = manager.findTargetWin(&structs.Point{X: 50, Y: 50})
	}

	if manager.target != nil {
		manager.target.onKeyDown(code)
	}

	return
}

func (manager *WindowManager) OnKeyUp(code int) {
	if manager.target != nil {
		manager.target.onKeyUp(code)
	}

	return
}

func (manager *WindowManager) OnWheel(delta float64) {
	manager.postRedraw()

	if manager.target == nil {
		manager.target = manager.findTargetWin(&structs.Point{X: 50, Y: 50})
	}

	if manager.target != nil {
		manager.target.onWheel(delta)
	}

	return
}

func (manager *WindowManager) dispatchPointerMoveOut() *WindowManager {
	manager.OnPointerMove(&structs.Point{X: -1, Y: -1})
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
	duration := time.Now().Sub(manager.startTime).Seconds()
	fps := math.Floor(float64(1000*manager.drawCount) / duration)

	if duration > 1000 {
		manager.drawCount = 0
		manager.startTime = time.Now()
	}

	return int(fps)
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

	fmt.Printf("postRedraw \n")
	manager.requestCount++
	if manager.requestCount < 2 {
		dom.GetWindow().RequestAnimationFrame(func(d time.Duration) {
			fmt.Printf("RequestAnimationFrame %v\n", d)
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
	if tipsWidget == nil || tipsWidget.parent == nil {
		return
	}

	hideTipsCanvas()
	p := tipsWidget.getPositionInView()

	context.Save()
	context.Translate(float64(p.X), float64(p.Y))
	context.BeginPath()
	tipsWidget.drawTips(context)
	context.Restore()

	return
}

func (manager *WindowManager) beforeDrawWindows(context *dom.CanvasRenderingContext2D) {

}

func (manager *WindowManager) afterDrawWindows(context *dom.CanvasRenderingContext2D) {

}

func (manager *WindowManager) drawWindows(context *dom.CanvasRenderingContext2D) {
	fmt.Printf("drawWindows \n")

	manager.beforeDrawWindows(context)
	for _, window := range manager.windows {
		window.draw(context)
	}
	manager.drawTips(context)
	manager.afterDrawWindows(context)

	return
}

func (manager *WindowManager) checkNeedRedraw(timeStep float64) bool {
	return true
}

func (manager *WindowManager) getCanvas2D() *dom.CanvasRenderingContext2D {
	if manager.ctx == nil {
		manager.ctx = manager.canvas.GetContext2d()
	}

	return manager.ctx
}

func (manager *WindowManager) doDraw(ctx *dom.CanvasRenderingContext2D) {
	now := time.Now()
	timeStep := now.Sub(manager.lastUpdateTime)
	fmt.Printf("doDraw \n")
	if !manager.checkNeedRedraw(timeStep.Seconds()) {
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
		ctx.Rect(0, 0, float64(w), float64(h))
		ctx.FillStyle = "Black"
		ctx.Fill()

		ctx.Save()
		ctx.TextAlign = "center"
		ctx.TextBaseline = "middle"
		ctx.Font = "20px Sans"
		ctx.FillStyle = "White"
		ctx.FillText(strconv.Itoa(str), float64(w>>1), float64(h>>1), -1)
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

	// ctx.BeginFrame()
	manager.doDraw(ctx)
	// ctx.EndFrame()

	return
}
