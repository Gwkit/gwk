package gwk

import (
	"github.com/Luncher/gwk/pkg/keyevent"
	"github.com/Luncher/gwk/pkg/structs"
	"honnef.co/go/js/dom"
	"math"
)

const (
	UNKNOW = iota
	SCROLL_TYPE_V
	SCROLL_TYPE_H
	SCROLL_TYPE_BOTH
	SCROLL_TYPE_AUTO
	SCROLL_TYPE_NONE
)

type ScrollType int

func (t ScrollType) String() string {
	switch t {
	case SCROLL_TYPE_V:
		return "vertical"
	case SCROLL_TYPE_H:
		return "horizon"
	case SCROLL_TYPE_BOTH:
		return "both"
	case SCROLL_TYPE_AUTO:
		return "auto"
	case SCROLL_TYPE_NONE:
		return "none"
	default:
		return "unknow"
	}
}

type ScrollView struct {
	*Widget
	xOffset       float64
	yOffset       float64
	scrollBarSize float64
	scrollType    ScrollType
	vScrollBar    *VScrollBar
	hScrollBar    *HScrollBar
	isScrollView  bool
	virtualSize   *structs.Rect
	workArea      *structs.Rect
}

func NewScrollView(parent *Widget, x, y, w, h float32) *ScrollView {
	scrollView := &ScrollView{
		Widget:        NewWidget(TYPE_SCROLL_VIEW, parent, x, y, w, h),
		isScrollView:  true,
		scrollBarSize: 8,
	}

	scrollView.vScrollBar =
		NewVScrollBar(scrollView.Widget, w-float32(scrollView.scrollBarSize), 0, float32(scrollView.scrollBarSize), h)
	scrollView.hScrollBar =
		NewHScrollBar(scrollView.Widget, 0, h-float32(scrollView.scrollBarSize), w, float32(scrollView.scrollBarSize))

	scrollView.vScrollBar.scrolledHandler = func(currentPosition, scrollRange float64) {
		if scrollView.virtualSize != nil {
			scrollView.yOffset = math.Min(math.Max(0, currentPosition), float64(scrollView.virtualSize.H-scrollView.rect.H))
		}
	}

	scrollView.hScrollBar.scrolledHandler = func(currentPosition, scrollRange float64) {
		if scrollView.virtualSize != nil {
			scrollView.xOffset = math.Min(math.Max(0, currentPosition), float64(scrollView.virtualSize.W-scrollView.rect.W))
		}
	}

	return scrollView
}

func (view *ScrollView) onScrolledH(offset float64) {

}

func (view *ScrollView) onScrolledV(offset float64) {

}

func (view *ScrollView) getScrollPositionH() float64 {
	return view.hScrollBar.getCurrentPosition()
}

func (view *ScrollView) getScrollPositionV() float64 {
	return view.vScrollBar.getCurrentPosition()
}

func (view *ScrollView) setScrollPositionH(position float64) {
	view.hScrollBar.setCurrentPosition(position)

	return
}

func (view *ScrollView) setScrollPositionV(position float64) {
	view.vScrollBar.setCurrentPosition(position)

	return
}

func (view *ScrollView) setScrollBarSize(scrollBarSize float64) {
	view.scrollBarSize = scrollBarSize

	return
}

func (view *ScrollView) setScrollType(t ScrollType) {
	view.scrollType = t

	return
}

func (view *ScrollView) onShow(visible bool) {
	view.needRelayout = visible

	return
}

func (view *ScrollView) onRelayout(workArea *structs.Rect, virtualSize *structs.Rect) {

}

func (view *ScrollView) setNeedRelayout(value bool) {
	view.needRelayout = value

	return
}

func (view *ScrollView) relayout(context *dom.CanvasRenderingContext2D, force bool) {
	if view.needRelayout || force {
		v := view.getScrollPositionV()
		view.updateScrollBar()
		view.needRelayout = false
		view.onRelayout(view.workArea, view.virtualSize)
		view.setScrollPositionV(v)
	}

	return
}

func (view *ScrollView) getVirtualSize() *structs.Rect {
	size := structs.NewRect(0, 0, view.rect.W, view.rect.H)

	for _, child := range view.children {
		if !child.visible {
			continue
		}
		right := child.rect.X + child.rect.W
		bottom := child.rect.Y + child.rect.H

		if right > size.W {
			size.W = right
		}
		if bottom > size.H {
			size.H = bottom
		}
	}

	size.W += int(view.scrollBarSize)
	size.H += int(view.scrollBarSize)

	size.W = int(math.Max(float64(view.rect.W), float64(size.W)))
	size.H = int(math.Max(float64(view.rect.H), float64(size.H)))

	view.virtualSize = size

	return size
}

func (view *ScrollView) updateScrollBar() bool {
	rect := view.rect
	size := view.getVirtualSize()
	vScrollBar := view.vScrollBar
	hScrollBar := view.hScrollBar
	scrollBarSize := view.scrollBarSize

	vScrollBar.resize(int(scrollBarSize), rect.H-int(scrollBarSize))
	vScrollBar.move(rect.W-int(scrollBarSize), 0)
	hScrollBar.resize(rect.W-int(scrollBarSize), int(scrollBarSize))
	hScrollBar.move(0, rect.H-int(scrollBarSize))

	switch view.scrollType {
	case SCROLL_TYPE_V:
		hScrollBar.show(false)
		vScrollBar.show(true)
		vScrollBar.setScrollRange(float64(size.H))
		vScrollBar.setCurrentPosition(0)
		vScrollBar.resize(int(scrollBarSize), rect.H)
	case SCROLL_TYPE_H:
		vScrollBar.show(false)
		hScrollBar.show(true)
		hScrollBar.setScrollRange(float64(size.W))
		hScrollBar.setCurrentPosition(0)
		hScrollBar.resize(rect.W, int(scrollBarSize))
	case SCROLL_TYPE_BOTH:
		vScrollBar.show(true)
		hScrollBar.show(true)
		vScrollBar.setScrollRange(float64(size.H))
		vScrollBar.setCurrentPosition(0)
		hScrollBar.setScrollRange(float64(size.W))
		hScrollBar.setCurrentPosition(0)
	case SCROLL_TYPE_NONE:
		vScrollBar.show(false)
		vScrollBar.show(false)
	default:
		if size.W > rect.W {
			hScrollBar.setScrollRange(float64(size.W))
			hScrollBar.setCurrentPosition(0)
		} else {
			hScrollBar.show(false)
		}
		if size.H > rect.H {
			vScrollBar.setScrollRange(float64(size.H))
			vScrollBar.setCurrentPosition(0)
		} else {
			vScrollBar.show(false)
		}
	}

	workArea := structs.NewRect(0, 0, rect.W, rect.H)
	if vScrollBar.visible {
		workArea.W = rect.W - int(scrollBarSize)
	}
	if hScrollBar.visible {
		workArea.H = rect.H - int(scrollBarSize)
	}

	view.workArea = workArea

	return true
}

func (view *ScrollView) paintChildren(context *dom.CanvasRenderingContext2D) {
	ww := view.rect.W
	hh := view.rect.H
	xOffset := view.getXOffset()
	yOffset := view.getYOffset()

	viewT := yOffset
	viewL := xOffset
	viewR := viewL + float64(ww)
	viewB := viewT + float64(hh)
	border := view.border

	context.Save()
	context.ClearRect(float64(border), float64(border), float64(ww-2*border), float64(hh-2*border))
	context.Translate(-xOffset, -yOffset)

	var focusChild *Widget
	paintChildrenFocusLater := view.paintChildrenFocusLater

	for _, child := range view.children {
		rect := child.rect
		iterFocused := false
		if child.state == STATE_OVER || child.state == STATE_ACTIVE {
			iterFocused = true
		}

		if !child.visible {
			continue
		}

		if rect.X > int(viewR) || rect.Y > int(viewB) {
			continue
		}

		if rect.X+rect.W < int(viewL) || int(viewT) > rect.Y+rect.H {
			continue
		}

		if paintChildrenFocusLater != nil && iterFocused {
			focusChild = child
		} else {
			child.draw(context)
		}
	}

	if paintChildrenFocusLater != nil && focusChild != nil {
		focusChild.draw(context)
	}

	context.Restore()
	view.hScrollBar.draw(context)
	view.vScrollBar.draw(context)

	return
}

func (view *ScrollView) onPointerDownScrollBar(point *structs.Point) bool {
	vScrollBar := view.vScrollBar
	hScrollBar := view.hScrollBar

	p := view.translatePoint(point)
	if vScrollBar.visible && isPointInRect(p, vScrollBar.rect) {
		vScrollBar.onPointerDown(point)
		return true
	}

	if hScrollBar.visible && isPointInRect(p, hScrollBar.rect) {
		hScrollBar.onPointerDown(point)
		return true
	}

	return false
}

func (view *ScrollView) onPointerDown(point *structs.Point) {
	if view.onPointerDownScrollBar(point) {
		return
	}

	p := structs.NewPoint(point.X+int(view.xOffset), point.Y+int(view.yOffset))
	view.Widget.onPointerDown(p)

	return
}

func (view *ScrollView) onPointerMove(point *structs.Point) {
	p := structs.NewPoint(point.X+int(view.xOffset), point.Y+int(view.yOffset))
	view.Widget.onPointerMove(p)

	return
}

func (view *ScrollView) onPointerUp(point *structs.Point) {
	p := structs.NewPoint(point.X+int(view.xOffset), point.Y+int(view.yOffset))
	view.Widget.onPointerUp(p)

	return
}

func (view *ScrollView) onLongPress(point *structs.Point) {
	p := structs.NewPoint(point.X+int(view.xOffset), point.Y+int(view.yOffset))
	view.Widget.onLongPress(p)

	return
}

func (view *ScrollView) onDoubleClick(point *structs.Point) {
	p := structs.NewPoint(point.X+int(view.xOffset), point.Y+int(view.yOffset))
	view.Widget.onDoubleClick(p)

	return
}

func (view *ScrollView) onContextMenu(point *structs.Point) {
	p := structs.NewPoint(point.X+int(view.xOffset), point.Y+int(view.yOffset))
	view.Widget.onContextMenu(p)

	return
}

func (view *ScrollView) getXOffset() float64 {
	return view.hScrollBar.getCurrentPosition()
}

func (view *ScrollView) getYOffset() float64 {
	return view.vScrollBar.getCurrentPosition()
}

func (view *ScrollView) getXScrollRange() float64 {
	return view.hScrollBar.getScrollRange()
}

func (view *ScrollView) getYScrollRange() float64 {
	return view.vScrollBar.getScrollRange()
}

func (view *ScrollView) setXOffset(xOffset float64) {
	view.hScrollBar.setCurrentPosition(xOffset)
}

func (view *ScrollView) setYOffset(yOffset float64) {
	view.vScrollBar.setCurrentPosition(yOffset)
}

func (view *ScrollView) onWheel(delta float64) {
	yOffset := view.getYOffset() + delta
	view.setYOffset(yOffset)
}

func (view *ScrollView) onKeyDown(code int) {
	delta := 10.0
	xOffset := 0.0
	yOffset := 0.0

	switch code {
	case keyevent.DOM_VK_LEFT:
		xOffset = view.getXOffset() - delta
	case keyevent.DOM_VK_RIGHT:
		xOffset = view.getXOffset() + delta
	case keyevent.DOM_VK_UP:
		yOffset = view.getYOffset() - delta
	case keyevent.DOM_VK_DOWN:
		yOffset = view.getYOffset() + delta
	case keyevent.DOM_VK_PAGE_UP:
		yOffset = view.getYOffset() - float64(view.getHeight()) + delta
	case keyevent.DOM_VK_PAGE_DOWN:
		xOffset = view.getXOffset() - float64(view.getWidth()) + delta
	case keyevent.DOM_VK_HOME:
		xOffset = 0
	case keyevent.DOM_VK_END:
		xOffset = view.getYScrollRange() - float64(view.getHeight())
	default:
		view.Widget.onKeyDown(code)
	}

	view.setXOffset(xOffset)
	view.setYOffset(yOffset)

	return
}
