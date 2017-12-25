package gwk

import (
	"fmt"
	"github.com/Luncher/gwk/pkg/image"
	"github.com/Luncher/gwk/pkg/structs"
	"honnef.co/go/js/dom"
	"math"
)

type ScrollBar struct {
	*Widget
	dragging             bool
	scrollRange          float64
	currentPosition      float64
	draggerRect          *structs.Rect
	pointerDownPoint     *structs.Point
	lastPointerPoint     *structs.Point
	currentPositionSaved float64
	scrolledHandler      onScrolledHandler
}

type onScrolledHandler func(currentPosition, scrollRange float64)

func NewScrollBar(t string, parent *Widget, x, y, w, h float32) *ScrollBar {
	bar := &ScrollBar{
		Widget:      NewWidget(t, parent, x, y, w, h),
		scrollRange: 100,
	}
	bar.I = bar

	return bar
}

func (bar *ScrollBar) onScrolled(currentPosition, scrollRange float64) {
	fmt.Printf("Scroll: %f\b", currentPosition)
	if bar.scrolledHandler != nil {
		bar.scrolledHandler(currentPosition, scrollRange)
	}

	return
}

func (bar *ScrollBar) onPointerDown(point *structs.Point) {
	bar.getWindow().grab(bar.Widget)
	bar.pointerDownPoint.X = point.X
	bar.pointerDownPoint.Y = point.Y
	p := bar.translatePoint(point)

	if bar.draggerRect != nil && isPointInRect(p, bar.draggerRect) {
		bar.dragging = true
		bar.setState(STATE_ACTIVE, false)
		bar.currentPositionSaved = bar.currentPosition
	} else {
		bar.setState(STATE_NORMAL, false)
	}

	return
}

func (bar *ScrollBar) onPointerMove(point *structs.Point) {
	if bar.dragging {
		ww := float64(bar.rect.W)
		hh := float64(bar.rect.H)
		if ww > hh {
			dx := float64(point.X - bar.pointerDownPoint.X)
			bar.setCurrentPosition(bar.currentPositionSaved + (dx/ww)*bar.scrollRange)
		} else {
			dy := float64(point.Y - bar.pointerDownPoint.Y)
			bar.setCurrentPosition(bar.currentPositionSaved + (dy/hh)*bar.scrollRange)
		}
	}

	return
}

func (bar *ScrollBar) onPointerUp(point *structs.Point) {
	ww := bar.rect.W
	hh := bar.rect.H
	if bar.draggerRect == nil {
		return
	}

	if !bar.dragging {
		r := bar.draggerRect
		p := bar.translatePoint(point)
		if ww > hh {
			if p.X < r.X {
				bar.addToCurrentPosition(float64(-ww))
			} else if p.X > (r.X + r.W) {
				bar.addToCurrentPosition(float64(ww))
			}
		} else {
			if p.Y < r.Y {
				bar.addToCurrentPosition(float64(-hh))
			} else if p.Y > (r.Y + r.H) {
				bar.addToCurrentPosition(float64(hh))
			}
		}
	}
	bar.dragging = false
	bar.getWindow().ungrab()
	bar.setState(STATE_NORMAL, false)

	return
}

func (bar *ScrollBar) setScrollRange(val float64) {
	bar.scrollRange = val
	bar.updateDraggerSize()

	return
}

func (bar *ScrollBar) getScrollRange() float64 {
	return bar.scrollRange
}

func (bar *ScrollBar) addToCurrentPosition(delta float64) {
	currentPosition := bar.currentPosition + delta
	bar.setCurrentPosition(currentPosition)

	return
}

func (bar *ScrollBar) getCurrentPosition() float64 {
	return bar.currentPosition
}

func (bar *ScrollBar) setCurrentPosition(currentPosition float64) {
	size := math.Max(float64(bar.getWidth()), float64(bar.getHeight()))
	bar.currentPosition = math.Max(math.Min(bar.scrollRange-size, currentPosition), 0)

	bar.PostRedraw()
	bar.updateDraggerSize()
	bar.onScrolled(bar.currentPosition, bar.scrollRange)

	return
}

func (bar *ScrollBar) onSized(w, h int) {
	bar.updateDraggerSize()

	return
}

func (bar *ScrollBar) updateDraggerSize() {
	var xx, yy, ww, hh, percent float64
	rect := bar.rect

	if rect.W > rect.H {
		bar.scrollRange = math.Max(float64(rect.W), bar.scrollRange)
		percent = bar.currentPosition / bar.scrollRange

		yy = 2
		hh = float64(rect.H) - yy - yy
		ww = math.Max(20, math.Floor(float64(rect.W)*float64(rect.W)/bar.scrollRange))
		xx = math.Floor(math.Min(math.Max(0, percent*float64(rect.W)), float64(rect.W)-ww))
	} else {
		bar.scrollRange = math.Max(float64(rect.H), bar.scrollRange)
		percent = bar.currentPosition / bar.scrollRange

		xx = 2
		ww = float64(rect.W) - xx - xx
		hh = math.Max(20, math.Floor(float64(rect.H)*float64(rect.H)/bar.scrollRange))
		yy = math.Floor(math.Min(math.Max(0, percent*float64(rect.H)), float64(rect.H)-hh))
	}

	bar.draggerRect = structs.NewRect(int(xx), int(yy), int(ww), int(hh))
	bar.currentPosition = bar.scrollRange * percent

	return
}

func (bar *ScrollBar) paintSelf(context *dom.CanvasRenderingContext2D) {
	if bar.draggerRect != nil {
		r := bar.draggerRect
		style := bar.getStyle("")
		if style.FgImage != nil {
			style.FgImage.Draw(context, image.DISPLAY_9PATCH, r.X, r.Y, r.W, r.H)
		} else {
			context.FillStyle = style.DragColor
			context.FillRect(float64(r.X), float64(r.Y), float64(r.W), float64(r.H))
		}
	}

	return
}

type VScrollBar struct {
	*ScrollBar
}

func NewVScrollBar(parent *Widget, x, y, w, h float32) *VScrollBar {
	return &VScrollBar{
		ScrollBar: NewScrollBar(TYPE_VSCROLL_BAR, parent, x, y, w, h),
	}
}

type HScrollBar struct {
	*ScrollBar
}

func NewHScrollBar(parent *Widget, x, y, w, h float32) *HScrollBar {
	return &HScrollBar{
		ScrollBar: NewScrollBar(TYPE_HSCROLL_BAR, parent, x, y, w, h),
	}
}
