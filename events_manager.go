package gwk

import (
	"github.com/Luncher/gwk/pkg/keyevent"
	"honnef.co/go/js/dom"
	"strings"
)

type EventsManager struct {
	point             *Point
	lastPoint         *Point
	lastPointEvent    dom.Event
	pointerDown       bool
	pointerDownPoint  *Point
	eventsConsumer    *WindowManager
	longPressDuration int
	pointerDeviceType string
	points            []Point
}

func NewEventsManager() *EventsManager {
	em := &EventsManager{}
	em.longPressDuration = 600

	return em
}

func (em *EventsManager) setEventsConsumer(eventsConsumer *WindowManager, element dom.HTMLCanvasElement) {
	if em.eventsConsumer == nil {
		em.eventsConsumer = eventsConsumer
		em.addEventListeners(element)
	}

	return
}

func (em *EventsManager) cancelDefaultAction(event dom.Event) bool {
	event.PreventDefault()

	return false
}

func (em *EventsManager) targetIsEditor(event dom.Event) bool {
	input := event.Target().(*dom.HTMLInputElement)
	name := strings.ToLower(input.TagName())

	if name != "body" && name != "canvas" {
		return true
	}

	return false
}

func (em *EventsManager) shouldIgnoreKey(event dom.KeyboardEvent) bool {
	code := event.KeyCode

	if code == keyevent.DOM_VK_F5 ||
		code == keyevent.DOM_VK_F12 ||
		code == keyevent.DOM_VK_F11 {
		return true
	}

	if em.targetIsEditor(event) {
		return true
	}

	return false
}

func (em *EventsManager) onKeyDownGlobal(event dom.KeyboardEvent) bool {
	code := event.KeyCode

	if em.shouldIgnoreKey(event) {
		return true
	} else {
		em.onKeyDown(code, event)
		return em.cancelDefaultAction(event)
	}
}

func (em *EventsManager) onKeyUpGlobal(event dom.KeyboardEvent) bool {
	code := event.KeyCode

	if em.shouldIgnoreKey(event) {
		return true
	} else {
		em.onKeyUp(code, event)
		return em.cancelDefaultAction(event)
	}
}

func (em *EventsManager) onWheelGlobal(event dom.WheelEvent) bool {
	if event.Target().GetAttribute("localName") != "canvas" {
		return em.cancelDefaultAction(event)
	}

	delta := event.DeltaY
	if delta != 0 {
		em.onWheel(delta, event)
		return em.cancelDefaultAction(event)
	}

	return true
}

func (em *EventsManager) addEventListeners(element dom.HTMLElement) {
	em.pointerDeviceType = "mouse"
	element.AddEventListener("dblclick", false, func(event dom.Event) {
		em.onDoubleClickGlobal(event.(dom.MouseEvent))
	})

	element.AddEventListener("mousewheel", false, func(event dom.Event) {
		em.onWheelGlobal(event.(dom.WheelEvent))
	})

	element.AddEventListener("DOMMouseScroll", false, func(event dom.Event) {
		em.onWheelGlobal(event.(dom.WheelEvent))
	})

	element.AddEventListener("mousedown", false, func(event dom.Event) {
		em.onMouseDownGlobal(event.(dom.MouseEvent))
	})

	element.AddEventListener("mousemove", false, func(event dom.Event) {
		em.onMouseMoveGlobal(event.(dom.MouseEvent))
	})

	element.AddEventListener("mouseup", false, func(event dom.Event) {
		em.onMouseUpGlobal(event.(dom.MouseEvent))
	})

	element.AddEventListener("keyup", false, func(event dom.Event) {
		em.onKeyUpGlobal(event.(dom.KeyboardEvent))
	})

	element.AddEventListener("keydown", false, func(event dom.Event) {
		em.onKeyDownGlobal(event.(dom.KeyboardEvent))
	})

	return
}

func (em *EventsManager) getAbsPoint(event dom.Event, i int) *Point {
	p := &em.points[i]

	if event != nil {
		me := event.(dom.MouseEvent)
		p.x = me.ClientX
		p.y = me.ClientY

		em.lastPoint.x = p.x
		em.lastPoint.y = p.y
		em.lastPointEvent = me
	} else {
		p = em.lastPoint
	}

	return p
}

func (em *EventsManager) isRightMouseEvent(event dom.MouseEvent) bool {
	return event.Button > 2 && event.Button != 4
}

func (em *EventsManager) onDoubleClickGlobal(event dom.MouseEvent) bool {
	if em.targetIsEditor(event) {
		return true
	}

	if !em.isRightMouseEvent(event) {
		em.onDoubleClick(em.getAbsPoint(event, 0), event)
	}

	return em.cancelDefaultAction(event)
}

func (em *EventsManager) onMouseDownGlobal(event dom.MouseEvent) bool {
	if em.targetIsEditor(event) {
		return true
	}

	if !em.isRightMouseEvent(event) {
		em.onPointerDown(em.getAbsPoint(event, 0), event)
	}

	return em.cancelDefaultAction(event)
}

func (em *EventsManager) onMouseMoveGlobal(event dom.MouseEvent) bool {
	if em.targetIsEditor(event) || !em.pointerDown {
		return true
	}

	em.onPointerMove(em.getAbsPoint(event, 0), event)

	return em.cancelDefaultAction(event)
}

func (em *EventsManager) onMouseUpGlobal(event dom.MouseEvent) bool {
	if em.targetIsEditor(event) || !em.pointerDown {
		return true
	}

	if em.isRightMouseEvent(event) {
		em.onContextMenu(em.getAbsPoint(event, 0), event)
	} else {
		em.onPointerUp(em.getAbsPoint(event, 0), event)
	}

	return em.cancelDefaultAction(event)
}

func (em *EventsManager) onKeyDown(code int, event dom.Event) {
	if em.eventsConsumer.preprocessEvent("", event) {
		em.eventsConsumer.onKeyDown(code)
	}
}

func (em *EventsManager) onKeyUp(code int, event dom.Event) {
	if em.eventsConsumer.preprocessEvent("", event) {
		em.eventsConsumer.onKeyUp(code)
	}
}

func (em *EventsManager) onWheel(delta float64, event dom.Event) {
	if em.eventsConsumer.preprocessEvent("", event) {
		em.eventsConsumer.onWheel(delta)
	}
}

func (em *EventsManager) onLongPress() {
	em.onContextMenu(em.lastPoint, em.lastPointEvent)
}

func (em *EventsManager) onPointerDown(point *Point, event dom.Event) {
	em.pointerDown = true

	if em.eventsConsumer.preprocessEvent("", event) {
		em.eventsConsumer.onPointerDown(point)
	}
}

func (em *EventsManager) onPointerMove(point *Point, event dom.Event) {
	if em.eventsConsumer.preprocessEvent("", event) {
		em.eventsConsumer.onPointerMove(point)
	}
}

func (em *EventsManager) onPointerUp(point *Point, event dom.Event) {
	if em.eventsConsumer.preprocessEvent("", event) {
		em.eventsConsumer.onPointerUp(point)
	}
	em.pointerDown = false
}

func (em *EventsManager) onDoubleClick(point *Point, event dom.Event) {
	if em.eventsConsumer.preprocessEvent("", event) {
		em.eventsConsumer.onDoubleClick(point)
	}
}

func (em *EventsManager) onContextMenu(point *Point, event dom.Event) {
	if em.eventsConsumer.preprocessEvent("", event) {
		em.eventsConsumer.onContextMenu(point)
	}
}

func (em *EventsManager) getPointerDeviceType() string {
	return em.pointerDeviceType
}

func (em *EventsManager) getInputScale() (x, y float32) {
	return em.eventsConsumer.getInputScale()
}

var emInstance *EventsManager = nil

func GetEventsManagerInstance() *EventsManager {
	if emInstance == nil {
		emInstance = NewEventsManager()
	}

	return emInstance
}
