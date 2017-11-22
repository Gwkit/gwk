package event

import (
	"github.com/Luncher/gwk/pkg/keyevent"
	"github.com/Luncher/gwk/pkg/structs"
	"honnef.co/go/js/dom"
	"strings"
)

type EventConsumer interface {
	OnKeyDown(code int)
	OnKeyUp(code int)
	OnWheel(delta float64)
	OnPointerDown(point *structs.Point)
	OnPointerMove(point *structs.Point)
	OnPointerUp(point *structs.Point)
	OnDoubleClick(point *structs.Point)
	OnContextMenu(point *structs.Point)
	PreprocessEvent(string, dom.Event) bool
	GetInputScale() (x, y float32)
}

type EventsManager struct {
	point             *structs.Point
	lastPoint         *structs.Point
	lastPointEvent    dom.Event
	pointerDown       bool
	pointerDownPoint  *structs.Point
	eventsConsumer    EventConsumer
	longPressDuration int
	pointerDeviceType string
	points            []structs.Point
}

var eventManager *EventsManager

func newEventsManager() *EventsManager {
	manager := &EventsManager{}
	manager.longPressDuration = 600

	return manager
}

func (manager *EventsManager) setEventsConsumer(consumer EventConsumer, element dom.HTMLCanvasElement) {
	if manager.eventsConsumer == nil {
		manager.eventsConsumer = consumer
		manager.addEventListeners(element)
	}

	return
}

func (manager *EventsManager) addEventListeners(element dom.Element) {
	manager.pointerDeviceType = "mouse"
	element.AddEventListener("dbclick", false, func(event dom.Event) {
		manager.onDoubleClickGlobal(event.(dom.MouseEvent))
	})

	element.AddEventListener("mousewheel", false, func(event dom.Event) {
		manager.onWheelGlobal(event.(dom.WheelEvent))
	})

	element.AddEventListener("DOMMouseScroll", false, func(event dom.Event) {
		manager.onWheelGlobal(event.(dom.WheelEvent))
	})

	element.AddEventListener("mousedown", false, func(event dom.Event) {
		manager.onMouseDownGlobal(event.(dom.MouseEvent))
	})

	element.AddEventListener("mousemove", false, func(event dom.Event) {
		manager.onMouseMoveGlobal(event.(dom.MouseEvent))
	})

	element.AddEventListener("mouseup", false, func(event dom.Event) {
		manager.onMouseUpGlobal(event.(dom.MouseEvent))
	})

	element.AddEventListener("keyup", false, func(event dom.Event) {
		manager.onKeyUpGlobal(event.(dom.KeyboardEvent))
	})

	element.AddEventListener("keydown", false, func(event dom.Event) {
		manager.onKeyDownGlobal(event.(dom.KeyboardEvent))
	})

	return
}

func (manager *EventsManager) targetIsEditor(event dom.Event) bool {
	input := event.Target().(*dom.HTMLInputElement)
	name := strings.ToLower(input.TagName())

	if name != "body" && name != "canvas" {
		return true
	}

	return false
}

func (manager *EventsManager) shouldIgnoreKey(event dom.KeyboardEvent) bool {
	code := event.KeyCode

	if code == keyevent.DOM_VK_F5 ||
		code == keyevent.DOM_VK_F12 ||
		code == keyevent.DOM_VK_F11 {
		return true
	}

	if manager.targetIsEditor(event) {
		return true
	}

	return false
}

func (manager *EventsManager) isRightMouseEvent(event dom.MouseEvent) bool {
	return event.Button > 2 && event.Button != 4
}

func (manager *EventsManager) getAbsPoint(event dom.Event, i int) *structs.Point {
	p := &manager.points[i]

	if event != nil {
		me := event.(dom.MouseEvent)
		p.X = me.ClientX
		p.X = me.ClientY

		manager.lastPoint.X = p.X
		manager.lastPoint.Y = p.Y
		manager.lastPointEvent = me
	} else {
		p = manager.lastPoint
	}

	return p
}

func (manager *EventsManager) cancelDefaultAction(event dom.Event) bool {
	event.PreventDefault()

	return false
}

func (manager *EventsManager) onDoubleClick(point *structs.Point, event dom.MouseEvent) {
	if manager.eventsConsumer.PreprocessEvent("", event) {
		manager.eventsConsumer.OnDoubleClick(point)
	}
}

func (manager *EventsManager) onDoubleClickGlobal(event dom.MouseEvent) bool {
	if manager.targetIsEditor(event) {
		return true
	}

	if !manager.isRightMouseEvent(event) {
		manager.onDoubleClick(manager.getAbsPoint(event, 0), event)
	}

	return manager.cancelDefaultAction(event)
}

func (manager *EventsManager) onPointerDown(point *structs.Point, event dom.Event) {
	manager.pointerDown = true

	if manager.eventsConsumer.PreprocessEvent("", event) {
		manager.eventsConsumer.OnPointerDown(point)
	}

	return
}

func (manager *EventsManager) onMouseDownGlobal(event dom.MouseEvent) bool {
	if manager.targetIsEditor(event) {
		return true
	}

	if !manager.isRightMouseEvent(event) {
		manager.onPointerDown(manager.getAbsPoint(event, 0), event)
	}

	return manager.cancelDefaultAction(event)
}

func (manager *EventsManager) onPointerMove(point *structs.Point, event dom.Event) {
	if manager.eventsConsumer.PreprocessEvent("", event) {
		manager.eventsConsumer.OnPointerMove(point)
	}
}

func (manager *EventsManager) onMouseMoveGlobal(event dom.MouseEvent) bool {
	if manager.targetIsEditor(event) || !manager.pointerDown {
		return true
	}

	manager.onPointerMove(manager.getAbsPoint(event, 0), event)

	return manager.cancelDefaultAction(event)
}

func (manager *EventsManager) onContextMenu(point *structs.Point, event dom.Event) {
	if manager.eventsConsumer.PreprocessEvent("", event) {
		manager.eventsConsumer.OnContextMenu(point)
	}
}

func (manager *EventsManager) onPointerUp(point *structs.Point, event dom.Event) {
	if manager.eventsConsumer.PreprocessEvent("", event) {
		manager.eventsConsumer.OnPointerUp(point)
	}
	manager.pointerDown = false
}

func (manager *EventsManager) onMouseUpGlobal(event dom.MouseEvent) bool {
	if manager.targetIsEditor(event) || !manager.pointerDown {
		return true
	}

	if manager.isRightMouseEvent(event) {
		manager.onContextMenu(manager.getAbsPoint(event, 0), event)
	} else {
		manager.onPointerUp(manager.getAbsPoint(event, 0), event)
	}

	return manager.cancelDefaultAction(event)
}

func (manager *EventsManager) onWheel(delta float64, event dom.Event) {
	if manager.eventsConsumer.PreprocessEvent("", event) {
		manager.eventsConsumer.OnWheel(delta)
	}
}

func (manager *EventsManager) onWheelGlobal(event dom.WheelEvent) bool {
	if event.Target().GetAttribute("localName") != "canvas" {
		return manager.cancelDefaultAction(event)
	}

	delta := event.DeltaY
	if delta != 0 {
		manager.onWheel(delta, event)
		return manager.cancelDefaultAction(event)
	}

	return true
}

func (em *EventsManager) onKeyUp(code int, event dom.Event) {
	if em.eventsConsumer.PreprocessEvent("", event) {
		em.eventsConsumer.OnKeyUp(code)
	}
}

func (manager *EventsManager) onKeyUpGlobal(event dom.KeyboardEvent) bool {
	code := event.KeyCode

	if manager.shouldIgnoreKey(event) {
		return true
	} else {
		manager.onKeyUp(code, event)
		return manager.cancelDefaultAction(event)
	}
}

func (manager *EventsManager) onKeyDown(code int, event dom.Event) {
	if manager.eventsConsumer.PreprocessEvent("", event) {
		manager.eventsConsumer.OnKeyDown(code)
	}
}

func (manager *EventsManager) onKeyDownGlobal(event dom.KeyboardEvent) bool {
	code := event.KeyCode

	if manager.shouldIgnoreKey(event) {
		return true
	} else {
		manager.onKeyDown(code, event)
		return manager.cancelDefaultAction(event)
	}
}

func (manager *EventsManager) getPointerDeviceType() string {
	return manager.pointerDeviceType
}

func (manager *EventsManager) getInputScale() (x, y float32) {
	return manager.eventsConsumer.GetInputScale()
}

func SetEventsConsumer(consumer EventConsumer, element dom.HTMLCanvasElement) {
	eventManager.setEventsConsumer(consumer, element)
}

func init() {
	eventManager = newEventsManager()
}
