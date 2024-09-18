package rui

import (
	"strings"
)

// Constants for [View] specific pointer events properties
const (
	// PointerDown is the constant for "pointer-down" property tag.
	//
	// Used by `View`.
	// Fired when a pointer becomes active. For mouse, it is fired when the device transitions from no buttons depressed to at 
	// least one button depressed. For touch, it is fired when physical contact is made with the digitizer. For pen, it is 
	// fired when the stylus makes physical contact with the digitizer.
	//
	// General listener format:
	// `func(view rui.View, event rui.PointerEvent)`.
	//
	// where:
	// view - Interface of a view which generated this event,
	// event - Pointer event.
	//
	// Allowed listener formats:
	// `func(event rui.PointerEvent)`,
	// `func(view rui.View)`,
	// `func()`.
	PointerDown = "pointer-down"

	// PointerUp is the constant for "pointer-up" property tag.
	//
	// Used by `View`.
	// Is fired when a pointer is no longer active.
	//
	// General listener format:
	// `func(view rui.View, event rui.PointerEvent)`.
	//
	// where:
	// view - Interface of a view which generated this event,
	// event - Pointer event.
	//
	// Allowed listener formats:
	// `func(event rui.PointerEvent)`,
	// `func(view rui.View)`,
	// `func()`.
	PointerUp = "pointer-up"

	// PointerMove is the constant for "pointer-move" property tag.
	//
	// Used by `View`.
	// Is fired when a pointer changes coordinates.
	//
	// General listener format:
	// `func(view rui.View, event rui.PointerEvent)`.
	//
	// where:
	// view - Interface of a view which generated this event,
	// event - Pointer event.
	//
	// Allowed listener formats:
	// `func(event rui.PointerEvent)`,
	// `func(view rui.View)`,
	// `func()`.
	PointerMove = "pointer-move"

	// PointerCancel is the constant for "pointer-cancel" property tag.
	//
	// Used by `View`.
	// Is fired if the pointer will no longer be able to generate events (for example the related device is deactivated).
	//
	// General listener format:
	// `func(view rui.View, event rui.PointerEvent)`.
	//
	// where:
	// view - Interface of a view which generated this event,
	// event - Pointer event.
	//
	// Allowed listener formats:
	// `func(event rui.PointerEvent)`,
	// `func(view rui.View)`,
	// `func()`.
	PointerCancel = "pointer-cancel"

	// PointerOut is the constant for "pointer-out" property tag.
	//
	// Used by `View`.
	// Is fired for several reasons including: pointing device is moved out of the hit test boundaries of an element; firing 
	// the "pointer-up" event for a device that does not support hover (see "pointer-up"); after firing the "pointer-cancel" 
	// event (see "pointer-cancel"); when a pen stylus leaves the hover range detectable by the digitizer.
	//
	// General listener format:
	// `func(view rui.View, event rui.PointerEvent)`.
	//
	// where:
	// view - Interface of a view which generated this event,
	// event - Pointer event.
	//
	// Allowed listener formats:
	// `func(event rui.PointerEvent)`,
	// `func(view rui.View)`,
	// `func()`.
	PointerOut = "pointer-out"

	// PointerOver is the constant for "pointer-over" property tag.
	//
	// Used by `View`.
	// Is fired when a pointing device is moved into an view's hit test boundaries.
	//
	// General listener format:
	// `func(view rui.View, event rui.PointerEvent)`.
	//
	// where:
	// view - Interface of a view which generated this event,
	// event - Pointer event.
	//
	// Allowed listener formats:
	// `func(event rui.PointerEvent)`,
	// `func(view rui.View)`,
	// `func()`.
	PointerOver = "pointer-over"
)

// PointerEvent represent a stylus events. Also inherit [MouseEvent] attributes
type PointerEvent struct {
	MouseEvent

	// PointerID is a unique identifier for the pointer causing the event.
	PointerID int

	// Width is the width (magnitude on the X axis), in pixels, of the contact geometry of the pointer.
	Width float64
	// Height is the height (magnitude on the Y axis), in pixels, of the contact geometry of the pointer.
	Height float64

	// Pressure is the normalized pressure of the pointer input in the range 0 to 1, where 0 and 1 represent
	// the minimum and maximum pressure the hardware is capable of detecting, respectively.
	Pressure float64

	// TangentialPressure is the normalized tangential pressure of the pointer input (also known
	// as barrel pressure or cylinder stress) in the range -1 to 1, where 0 is the neutral position of the control.
	TangentialPressure float64

	// TiltX is the plane angle (in degrees, in the range of -90 to 90) between the Y–Z plane
	// and the plane containing both the pointer (e.g. pen stylus) axis and the Y axis.
	TiltX float64

	// TiltY is the plane angle (in degrees, in the range of -90 to 90) between the X–Z plane
	// and the plane containing both the pointer (e.g. pen stylus) axis and the X axis.
	TiltY float64

	// Twist is the clockwise rotation of the pointer (e.g. pen stylus) around its major axis in degrees,
	// with a value in the range 0 to 359.
	Twist float64

	// PointerType indicates the device type that caused the event ("mouse", "pen", "touch", etc.)
	PointerType string

	// IsPrimary indicates if the pointer represents the primary pointer of this pointer type.
	IsPrimary bool
}

var pointerEvents = map[string]struct{ jsEvent, jsFunc string }{
	PointerDown:   {jsEvent: "onpointerdown", jsFunc: "pointerDownEvent"},
	PointerUp:     {jsEvent: "onpointerup", jsFunc: "pointerUpEvent"},
	PointerMove:   {jsEvent: "onpointermove", jsFunc: "pointerMoveEvent"},
	PointerCancel: {jsEvent: "onpointercancel", jsFunc: "pointerCancelEvent"},
	PointerOut:    {jsEvent: "onpointerout", jsFunc: "pointerOutEvent"},
	PointerOver:   {jsEvent: "onpointerover", jsFunc: "pointerOverEvent"},
}

func (view *viewData) setPointerListener(tag string, value any) bool {
	listeners, ok := valueToEventListeners[View, PointerEvent](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removePointerListener(tag)
	} else if js, ok := pointerEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			view.session.updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)")
		}
	} else {
		return false
	}
	return true
}

func (view *viewData) removePointerListener(tag string) {
	delete(view.properties, tag)
	if view.created {
		if js, ok := pointerEvents[tag]; ok {
			view.session.removeProperty(view.htmlID(), js.jsEvent)
		}
	}
}

func pointerEventsHtml(view View, buffer *strings.Builder) {
	for tag, js := range pointerEvents {
		if value := view.getRaw(tag); value != nil {
			if listeners, ok := value.([]func(View, PointerEvent)); ok && len(listeners) > 0 {
				buffer.WriteString(js.jsEvent)
				buffer.WriteString(`="`)
				buffer.WriteString(js.jsFunc)
				buffer.WriteString(`(this, event)" `)
			}
		}
	}
}

func (event *PointerEvent) init(data DataObject) {
	event.MouseEvent.init(data)

	event.PointerID, _ = dataIntProperty(data, "pointerId")
	event.Width = dataFloatProperty(data, "width")
	event.Height = dataFloatProperty(data, "height")
	event.Pressure = dataFloatProperty(data, "pressure")
	event.TangentialPressure = dataFloatProperty(data, "tangentialPressure")
	event.TiltX = dataFloatProperty(data, "tiltX")
	event.TiltY = dataFloatProperty(data, "tiltY")
	event.Twist = dataFloatProperty(data, "twist")
	value, _ := data.PropertyValue("pointerType")
	event.PointerType = value
	event.IsPrimary = dataBoolProperty(data, "isPrimary")
}

func handlePointerEvents(view View, tag string, data DataObject) {
	listeners := getEventListeners[View, PointerEvent](view, nil, tag)
	if len(listeners) == 0 {
		return
	}

	var event PointerEvent
	event.init(data)

	for _, listener := range listeners {
		listener(view, event)
	}
}

// GetPointerDownListeners returns the "pointer-down" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPointerDownListeners(view View, subviewID ...string) []func(View, PointerEvent) {
	return getEventListeners[View, PointerEvent](view, subviewID, PointerDown)
}

// GetPointerUpListeners returns the "pointer-up" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPointerUpListeners(view View, subviewID ...string) []func(View, PointerEvent) {
	return getEventListeners[View, PointerEvent](view, subviewID, PointerUp)
}

// GetPointerMoveListeners returns the "pointer-move" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPointerMoveListeners(view View, subviewID ...string) []func(View, PointerEvent) {
	return getEventListeners[View, PointerEvent](view, subviewID, PointerMove)
}

// GetPointerCancelListeners returns the "pointer-cancel" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPointerCancelListeners(view View, subviewID ...string) []func(View, PointerEvent) {
	return getEventListeners[View, PointerEvent](view, subviewID, PointerCancel)
}

// GetPointerOverListeners returns the "pointer-over" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPointerOverListeners(view View, subviewID ...string) []func(View, PointerEvent) {
	return getEventListeners[View, PointerEvent](view, subviewID, PointerOver)
}

// GetPointerOutListeners returns the "pointer-out" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPointerOutListeners(view View, subviewID ...string) []func(View, PointerEvent) {
	return getEventListeners[View, PointerEvent](view, subviewID, PointerOut)
}
