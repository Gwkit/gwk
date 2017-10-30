package gwk

import (
	"honnef.co/go/js/dom"
	"github.com/Luncher/gwk/pkg/rt"
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
	app.canvas = rt.
}

func (app *Application) getView() interface{} {
	return app.view
}

func (app *Application) setMinHeight(minHeight int) {
	app.minHeight = minHeight

	return
}
