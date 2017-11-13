package gwk

import (
	"github.com/Luncher/gwk/pkg/rt"
	"honnef.co/go/js/dom"
	"math"
)

var (
	TYPE_GENERAL       = "general"
	TYPE_WEBAPP        = "webapp"
	TYPE_PREVIEW       = "preview"
	TYPE_PC_VIEWER     = "pc_viewer"
	TYPE_PC_EDITOR     = "pc_editor"
	TYPE_MOBILE_EDITOR = "mobile_editor"
	TYPE_INLINE_EDITOR = "inline_editor"
)

type Application struct {
	win       *Window
	view      interface{}
	t         string
	canvasID  string
	minHeight int
	canvas    *dom.HTMLCanvasElement
	manager   *WindowManager
}

func NewApplication(canvasID, t string) *Application {
	app := &Application{}

	app.t = t
	app.canvasID = canvasID
	app.canvas = rt.GetRTInstance().GetMainCanvas(canvasID)
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

	h = int(math.Max(float64(h), float64(app.minHeight)))

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
