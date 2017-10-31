package gwk

import (
	"github.com/Luncher/gwk/pkg/rt"
	"honnef.co/go/js/dom"
	"math"
)

var (
	TYPE_GENERAL       = 0
	TYPE_WEBAPP        = 1
	TYPE_PREVIEW       = 2
	TYPE_PC_VIEWER     = 3
	TYPE_PC_EDITOR     = 4
	TYPE_MOBILE_EDITOR = 5
	TYPE_INLINE_EDITOR = 6
)

type Application struct {
	win       *Window
	view      interface{}
	t         string
	minHeight int
	canvas    dom.HTMLCanvasElement
	manager   *WindowManager
}

func NewApplication(t string) *Application {
	app := &Application{}

	app.t = t
	app.canvas = rt.GetRTInstance().GetMainCanvas()
	app.adjustCanvasSize()
	app.manager = NewWindowManager(app, app.canvas, app.canvas)

	return app
}

func (app *Application) adjustCanvasSize() {
	var w, h int

	canvas := app.canvas
	width, height := rt.GetRTInstance().GetViewPort()

	switch app.t {
	case TYPE_GENERAL:
		w = width - 20
		h = height
	case TYPE_WEBAPP:
		w = width
		h = height
	case TYPE_PREVIEW:
		w = width
		h = height
		app.setMinHeight(1500)
	default:
		if app.minHeight == 0 {
			app.setMinHeight(800)
		}
		w = width - 20
		h = height
	}

	h = math.Max(h, app.minHeight)

	app.resizeCanvasTo(w, h)

	return
}

func (app *Application) resizeCanvasTo(w, h int) {
	canvas := app.canvas

	canvas.Width = w
	canvas.Height = h
	canvas.Style().SetProperty("top", "0px", "")
	canvas.Style().SetProperty("left", "0px", "")
	canvas.Style().SetProperty("position", "absolute", "")

	return
}

func (app *Application) getView() interface{} {
	return app.view
}

func (app *Application) setMinHeight(minHeight int) {
	app.minHeight = minHeight

	return
}

func (app *Application) isDevApp() bool {
	return false
}
