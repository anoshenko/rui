package rui

import (
	"strings"
)

const (
	// PointerDown is the constant for "pointer-down" property tag.
	// The "pointer-down" event is fired when a pointer becomes active. For mouse, it is fired when
	// the device transitions from no buttons depressed to at least one button depressed.
	// For touch, it is fired when physical contact is made with the digitizer.
	// For pen, it is fired when the stylus makes physical contact with the digitizer.
	// The main listener format: func(View, PointerEvent).
	// The additional listener formats: func(PointerEvent), func(View), and func().
	PointerDown = "pointer-down"

	// PointerUp is the constant for "pointer-up" property tag.
	// The "pointer-up" event is fired when a pointer is no longer active.
	// The main listener format: func(View, PointerEvent).
	// The additional listener formats: func(PointerEvent), func(View), and func().
	PointerUp = "pointer-up"

	// PointerMove is the constant for "pointer-move" property tag.
	// The "pointer-move" event is fired when a pointer changes coordinates.
	// The main listener format: func(View, PointerEvent).
	// The additional listener formats: func(PointerEvent), func(View), and func().
	PointerMove = "pointer-move"

	// PointerCancel is the constant for "pointer-cancel" property tag.
	// The "pointer-cancel" event is fired if the pointer will no longer be able to generate events
	// (for example the related device is deactivated).
	// The main listener format: func(View, PointerEvent).
	// The additional listener formats: func(PointerEvent), func(View), and func().
	PointerCancel = "pointer-cancel"

	// PointerOut is the constant for "pointer-out" property tag.
	// The "pointer-out" event is fired for several reasons including: pointing device is moved out
	// of the hit test boundaries of an element; firing the pointerup event for a device
	// that does not support hover (see "pointer-up"); after firing the pointercancel event (see "pointer-cancel");
	// when a pen stylus leaves the hover range detectable by the digitizer.
	// The main listener format: func(View, PointerEvent).
	// The additional listener formats: func(PointerEvent), func(View), and func().
	PointerOut = "pointer-out"

	// PointerOver is the constant for "pointer-over" property tag.
	// The "pointer-over" event is fired when a pointing device is moved into an view's hit test boundaries.
	// The main listener format: func(View, PointerEvent).
	// The additional listener formats: func(PointerEvent), func(View), and func().
	PointerOver = "pointer-over"
)

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

func valueToPointerListeners(value interface{}) ([]func(View, PointerEvent), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(View, PointerEvent):
		return []func(View, PointerEvent){value}, true

	case func(PointerEvent):
		fn := func(view View, event PointerEvent) {
			value(event)
		}
		return []func(View, PointerEvent){fn}, true

	case func(View):
		fn := func(view View, event PointerEvent) {
			value(view)
		}
		return []func(View, PointerEvent){fn}, true

	case func():
		fn := func(view View, event PointerEvent) {
			value()
		}
		return []func(View, PointerEvent){fn}, true

	case []func(View, PointerEvent):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(PointerEvent):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View, PointerEvent), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view View, event PointerEvent) {
				v(event)
			}
		}
		return listeners, true

	case []func(View):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View, PointerEvent), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view View, event PointerEvent) {
				v(view)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View, PointerEvent), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view View, event PointerEvent) {
				v()
			}
		}
		return listeners, true

	case []interface{}:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View, PointerEvent), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(View, PointerEvent):
				listeners[i] = v

			case func(PointerEvent):
				listeners[i] = func(view View, event PointerEvent) {
					v(event)
				}

			case func(View):
				listeners[i] = func(view View, event PointerEvent) {
					v(view)
				}

			case func():
				listeners[i] = func(view View, event PointerEvent) {
					v()
				}

			default:
				return nil, false
			}
		}
		return listeners, true
	}

	return nil, false
}

var pointerEvents = map[string]struct{ jsEvent, jsFunc string }{
	PointerDown:   {jsEvent: "onpointerdown", jsFunc: "pointerDownEvent"},
	PointerUp:     {jsEvent: "onpointerup", jsFunc: "pointerUpEvent"},
	PointerMove:   {jsEvent: "onpointermove", jsFunc: "pointerMoveEvent"},
	PointerCancel: {jsEvent: "onpointercancel", jsFunc: "pointerCancelEvent"},
	PointerOut:    {jsEvent: "onpointerout", jsFunc: "pointerOutEvent"},
	PointerOver:   {jsEvent: "onpointerover", jsFunc: "pointerOverEvent"},
}

func (view *viewData) setPointerListener(tag string, value interface{}) bool {
	listeners, ok := valueToPointerListeners(value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removePointerListener(tag)
	} else if js, ok := pointerEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)", view.Session())
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
			removeProperty(view.htmlID(), js.jsEvent, view.Session())
		}
	}
}

func getPointerListeners(view View, subviewID string, tag string) []func(View, PointerEvent) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(View, PointerEvent)); ok {
				return result
			}
		}
	}
	return []func(View, PointerEvent){}
}

func pointerEventsHtml(view View, buffer *strings.Builder) {
	for tag, js := range pointerEvents {
		if value := view.getRaw(tag); value != nil {
			if listeners, ok := value.([]func(View, PointerEvent)); ok && len(listeners) > 0 {
				buffer.WriteString(js.jsEvent + `="` + js.jsFunc + `(this, event)" `)
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
	listeners := getPointerListeners(view, "", tag)
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
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetPointerDownListeners(view View, subviewID string) []func(View, PointerEvent) {
	return getPointerListeners(view, subviewID, PointerDown)
}

// GetPointerUpListeners returns the "pointer-up" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetPointerUpListeners(view View, subviewID string) []func(View, PointerEvent) {
	return getPointerListeners(view, subviewID, PointerUp)
}

// GetPointerMoveListeners returns the "pointer-move" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetPointerMoveListeners(view View, subviewID string) []func(View, PointerEvent) {
	return getPointerListeners(view, subviewID, PointerMove)
}

// GetPointerCancelListeners returns the "pointer-cancel" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetPointerCancelListeners(view View, subviewID string) []func(View, PointerEvent) {
	return getPointerListeners(view, subviewID, PointerCancel)
}

// GetPointerOverListeners returns the "pointer-over" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetPointerOverListeners(view View, subviewID string) []func(View, PointerEvent) {
	return getPointerListeners(view, subviewID, PointerOver)
}

// GetPointerOutListeners returns the "pointer-out" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetPointerOutListeners(view View, subviewID string) []func(View, PointerEvent) {
	return getPointerListeners(view, subviewID, PointerOut)
}
