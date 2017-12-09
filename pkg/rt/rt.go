package rt

import (
	"fmt"
	"honnef.co/go/js/dom"
	"time"
)

type GwkRT struct {
	canvas            *dom.HTMLCanvasElement
	mainCanvasW       int
	mainCanvasH       int
	mainCanvasScale   *struct{ x, y float32 }
	mainCanvasPostion *struct{ x, y int }
}

var rt = &GwkRT{}

func GetRTInstance() *GwkRT {
	return rt
}

func (rt *GwkRT) init() {

}

func (rt *GwkRT) GetViewPort() (int, int) {
	height := dom.GetWindow().InnerHeight()
	width := dom.GetWindow().InnerWidth()

	return width, height
}

func (rt *GwkRT) GetMainCanvas(id string) *dom.HTMLCanvasElement {
	if rt.canvas != nil {
		return rt.canvas
	}

	if len(id) == 0 {
		id = "main_canvas"
	}

	document := dom.GetWindow().Document()
	fmt.Printf("GetMainCanvas ID:  %s\n", id)
	canvas := document.GetElementByID(id).(*dom.HTMLCanvasElement)
	if canvas == nil {
		canvas := document.CreateElement("canvas").(*dom.HTMLCanvasElement)
		canvas.SetID("main-canvas")
		canvas.Style().SetProperty("zIndex", "0", "")
		dom.GetWindow().Document().AppendChild(canvas)
	}

	rt.canvas = canvas

	return canvas
}

func (rt *GwkRT) moveMainCanvas(x, y int) {
	canvas := rt.canvas

	canvas.Style().SetProperty("position", "absolute", "")
	canvas.Style().SetProperty("top", fmt.Sprintf("%dpx", x), "")
	canvas.Style().SetProperty("left", fmt.Sprintf("%dpx", y), "")

	rt.mainCanvasPostion.x = x
	rt.mainCanvasPostion.y = y

	return
}

func (rt *GwkRT) resizeMainCanvas(w, h, styleW, styleH int) {
	canvas := rt.GetMainCanvas("")

	canvas.Style().SetProperty("width", fmt.Sprintf("%d", w), "")
	canvas.Style().SetProperty("height", fmt.Sprintf("%d", h), "")
	rt.mainCanvasW = w
	rt.mainCanvasH = h
	rt.mainCanvasScale.x = float32(w) / float32(styleW)
	rt.mainCanvasScale.y = float32(h) / float32(styleH)

	return
}

func (rw *GwkRT) getMainCanvasScale() (float32, float32) {
	return rt.mainCanvasScale.x, rt.mainCanvasScale.y
}

func (rt *GwkRT) getMainCanvasPosition() (int, int) {
	return rt.mainCanvasPostion.x, rt.mainCanvasPostion.y
}

func (rt *GwkRT) createImage(src string,
	onLoad func(*dom.HTMLImageElement),
	onError func(interface{})) dom.HTMLImageElement {
	image := dom.GetWindow().Document().CreateElement("image").(*dom.HTMLImageElement)
	image.AddEventListener("onload", false, func(event dom.Event) {
		onLoad(image)
	})

	image.AddEventListener("onerror", false, func(event dom.Event) {
		onError(image)
	})

	image.Src = src

	return *image
}

func (rt *GwkRT) requestAnimFrame(callback func(time.Duration)) {
	dom.GetWindow().RequestAnimationFrame(callback)
}
