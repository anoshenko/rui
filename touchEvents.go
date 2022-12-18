package rui

import (
	"strconv"
	"strings"
)

const (
	// TouchStart is the constant for "touch-start" property tag.
	// The "touch-start" event is fired when one or more touch points are placed on the touch surface.
	// The main listener format: func(View, TouchEvent).
	// The additional listener formats: func(TouchEvent), func(View), and func().
	TouchStart = "touch-start"

	// TouchEnd is the constant for "touch-end" property tag.
	// The "touch-end" event fires when one or more touch points are removed from the touch surface.
	// The main listener format: func(View, TouchEvent).
	// The additional listener formats: func(TouchEvent), func(View), and func().
	TouchEnd = "touch-end"

	// TouchMove is the constant for "touch-move" property tag.
	// The "touch-move" event is fired when one or more touch points are moved along the touch surface.
	// The main listener format: func(View, TouchEvent).
	// The additional listener formats: func(TouchEvent), func(View), and func().
	TouchMove = "touch-move"

	// TouchCancel is the constant for "touch-cancel" property tag.
	// The "touch-cancel" event is fired when one or more touch points have been disrupted
	// in an implementation-specific manner (for example, too many touch points are created).
	// The main listener format: func(View, TouchEvent).
	// The additional listener formats: func(TouchEvent), func(View), and func().
	TouchCancel = "touch-cancel"
)

// Touch contains parameters of a single touch of a touch event
type Touch struct {
	// Identifier is a unique identifier for this Touch object. A given touch point (say, by a finger)
	// will have the same identifier for the duration of its movement around the surface.
	// This lets you ensure that you're tracking the same touch all the time.
	Identifier int

	// X provides the horizontal coordinate within the view's viewport.
	X float64
	// Y provides the vertical coordinate within the view's viewport.
	Y float64

	// ClientX provides the horizontal coordinate within the application's viewport at which the event occurred.
	ClientX float64
	// ClientY provides the vertical coordinate within the application's viewport at which the event occurred.
	ClientY float64

	// ScreenX provides the horizontal coordinate (offset) of the touch pointer in global (screen) coordinates.
	ScreenX float64
	// ScreenY provides the vertical coordinate (offset) of the touch pointer in global (screen) coordinates.
	ScreenY float64

	// RadiusX is the X radius of the ellipse that most closely circumscribes the area of contact with the screen.
	// The value is in pixels of the same scale as screenX.
	RadiusX float64
	// RadiusY is the Y radius of the ellipse that most closely circumscribes the area of contact with the screen.
	// The value is in pixels of the same scale as screenX.
	RadiusY float64

	// RotationAngle is the angle (in degrees) that the ellipse described by radiusX and radiusY must be rotated,
	// clockwise, to most accurately cover the area of contact between the user and the surface.
	RotationAngle float64

	// Force is the amount of pressure being applied to the surface by the user, as a float
	// between 0.0 (no pressure) and 1.0 (maximum pressure).
	Force float64
}

// TouchEvent contains parameters of a touch event
type TouchEvent struct {
	// TimeStamp is the time at which the event was created (in milliseconds).
	// This value is time since epoch—but in reality, browsers' definitions vary.
	TimeStamp uint64

	// Touches is the array of all the Touch objects representing all current points
	// of contact with the surface, regardless of target or changed status.
	Touches []Touch

	// CtrlKey == true if the control key was down when the event was fired. false otherwise.
	CtrlKey bool
	// ShiftKey == true if the shift key was down when the event was fired. false otherwise.
	ShiftKey bool
	// AltKey == true if the alt key was down when the event was fired. false otherwise.
	AltKey bool
	// MetaKey == true if the meta key was down when the event was fired. false otherwise.
	MetaKey bool
}

var touchEvents = map[string]struct{ jsEvent, jsFunc string }{
	TouchStart:  {jsEvent: "ontouchstart", jsFunc: "touchStartEvent"},
	TouchEnd:    {jsEvent: "ontouchend", jsFunc: "touchEndEvent"},
	TouchMove:   {jsEvent: "ontouchmove", jsFunc: "touchMoveEvent"},
	TouchCancel: {jsEvent: "ontouchcancel", jsFunc: "touchCancelEvent"},
}

func (view *viewData) setTouchListener(tag string, value any) bool {
	listeners, ok := valueToEventListeners[View, TouchEvent](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeTouchListener(tag)
	} else if js, ok := touchEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			view.session.updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)")
		}
	} else {
		return false
	}
	return true
}

func (view *viewData) removeTouchListener(tag string) {
	delete(view.properties, tag)
	if view.created {
		if js, ok := touchEvents[tag]; ok {
			view.session.removeProperty(view.htmlID(), js.jsEvent)
		}
	}
}

func touchEventsHtml(view View, buffer *strings.Builder) {
	for tag, js := range touchEvents {
		if value := view.getRaw(tag); value != nil {
			if listeners, ok := value.([]func(View, TouchEvent)); ok && len(listeners) > 0 {
				buffer.WriteString(js.jsEvent + `="` + js.jsFunc + `(this, event)" `)
			}
		}
	}
}

func (event *TouchEvent) init(data DataObject) {

	event.Touches = []Touch{}
	event.TimeStamp = getTimeStamp(data)
	if node := data.PropertyByTag("touches"); node != nil && node.Type() == ArrayNode {
		for i := 0; i < node.ArraySize(); i++ {
			if element := node.ArrayElement(i); element != nil && element.IsObject() {
				if obj := element.Object(); obj != nil {
					var touch Touch
					if value, ok := obj.PropertyValue("identifier"); ok {
						touch.Identifier, _ = strconv.Atoi(value)
					}
					touch.X = dataFloatProperty(obj, "x")
					touch.Y = dataFloatProperty(obj, "y")
					touch.ClientX = dataFloatProperty(obj, "clientX")
					touch.ClientY = dataFloatProperty(obj, "clientY")
					touch.ScreenX = dataFloatProperty(obj, "screenX")
					touch.ScreenY = dataFloatProperty(obj, "screenY")
					touch.RadiusX = dataFloatProperty(obj, "radiusX")
					touch.RadiusY = dataFloatProperty(obj, "radiusY")
					touch.RotationAngle = dataFloatProperty(obj, "rotationAngle")
					touch.Force = dataFloatProperty(obj, "force")
					event.Touches = append(event.Touches, touch)
				}
			}
		}
	}
	event.CtrlKey = dataBoolProperty(data, "ctrlKey")
	event.ShiftKey = dataBoolProperty(data, "shiftKey")
	event.AltKey = dataBoolProperty(data, "altKey")
	event.MetaKey = dataBoolProperty(data, "metaKey")
}

func handleTouchEvents(view View, tag string, data DataObject) {
	listeners := getEventListeners[View, TouchEvent](view, nil, tag)
	if len(listeners) == 0 {
		return
	}

	var event TouchEvent
	event.init(data)

	for _, listener := range listeners {
		listener(view, event)
	}
}

// GetTouchStartListeners returns the "touch-start" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTouchStartListeners(view View, subviewID ...string) []func(View, TouchEvent) {
	return getEventListeners[View, TouchEvent](view, subviewID, TouchStart)
}

// GetTouchEndListeners returns the "touch-end" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTouchEndListeners(view View, subviewID ...string) []func(View, TouchEvent) {
	return getEventListeners[View, TouchEvent](view, subviewID, TouchEnd)
}

// GetTouchMoveListeners returns the "touch-move" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTouchMoveListeners(view View, subviewID ...string) []func(View, TouchEvent) {
	return getEventListeners[View, TouchEvent](view, subviewID, TouchMove)
}

// GetTouchCancelListeners returns the "touch-cancel" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTouchCancelListeners(view View, subviewID ...string) []func(View, TouchEvent) {
	return getEventListeners[View, TouchEvent](view, subviewID, TouchCancel)
}
