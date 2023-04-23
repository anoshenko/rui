package rui

import "strings"

const (
	// KeyDown is the constant for "key-down-event" property tag.
	// The "key-down-event" event is fired when a key is pressed.
	// The main listener format:
	//   func(View, KeyEvent).
	// The additional listener formats:
	//   func(KeyEvent), func(View), and func().
	KeyDownEvent = "key-down-event"

	// KeyPp is the constant for "key-up-event" property tag.
	// The "key-up-event" event is fired when a key is released.
	// The main listener format:
	//   func(View, KeyEvent).
	// The additional listener formats:
	//   func(KeyEvent), func(View), and func().
	KeyUpEvent = "key-up-event"
)

type KeyEvent struct {
	// TimeStamp is the time at which the event was created (in milliseconds).
	// This value is time since epoch—but in reality, browsers' definitions vary.
	TimeStamp uint64

	// Key is the key value of the key represented by the event. If the value has a printed representation,
	// this attribute's value is the same as the char property. Otherwise, it's one of the key value strings
	// specified in Key values. If the key can't be identified, its value is the string "Unidentified".
	Key string

	// Code holds a string that identifies the physical key being pressed. The value is not affected
	// by the current keyboard layout or modifier state, so a particular key will always return the same value.
	Code string

	// Repeat == true if a key has been depressed long enough to trigger key repetition, otherwise false.
	Repeat bool

	// CtrlKey == true if the control key was down when the event was fired. false otherwise.
	CtrlKey bool

	// ShiftKey == true if the shift key was down when the event was fired. false otherwise.
	ShiftKey bool

	// AltKey == true if the alt key was down when the event was fired. false otherwise.
	AltKey bool

	// MetaKey == true if the meta key was down when the event was fired. false otherwise.
	MetaKey bool
}

func (event *KeyEvent) init(data DataObject) {
	getBool := func(tag string) bool {
		if value, ok := data.PropertyValue(tag); ok && value == "1" {
			return true
		}
		return false
	}

	event.Key, _ = data.PropertyValue("key")
	event.Code, _ = data.PropertyValue("code")
	event.TimeStamp = getTimeStamp(data)
	event.Repeat = getBool("repeat")
	event.CtrlKey = getBool("ctrlKey")
	event.ShiftKey = getBool("shiftKey")
	event.AltKey = getBool("altKey")
	event.MetaKey = getBool("metaKey")
}

func valueToEventListeners[V View, E any](value any) ([]func(V, E), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(V, E):
		return []func(V, E){value}, true

	case func(E):
		fn := func(_ V, event E) {
			value(event)
		}
		return []func(V, E){fn}, true

	case func(V):
		fn := func(view V, _ E) {
			value(view)
		}
		return []func(V, E){fn}, true

	case func():
		fn := func(V, E) {
			value()
		}
		return []func(V, E){fn}, true

	case []func(V, E):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(_ V, event E) {
				v(event)
			}
		}
		return listeners, true

	case []func(V):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view V, _ E) {
				v(view)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(V, E) {
				v()
			}
		}
		return listeners, true

	case []any:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(V, E):
				listeners[i] = v

			case func(E):
				listeners[i] = func(_ V, event E) {
					v(event)
				}

			case func(V):
				listeners[i] = func(view V, _ E) {
					v(view)
				}

			case func():
				listeners[i] = func(V, E) {
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

func valueToEventWithOldListeners[V View, E any](value any) ([]func(V, E, E), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(V, E, E):
		return []func(V, E, E){value}, true

	case func(V, E):
		fn := func(v V, val, _ E) {
			value(v, val)
		}
		return []func(V, E, E){fn}, true

	case func(E, E):
		fn := func(_ V, val, old E) {
			value(val, old)
		}
		return []func(V, E, E){fn}, true

	case func(E):
		fn := func(_ V, val, _ E) {
			value(val)
		}
		return []func(V, E, E){fn}, true

	case func(V):
		fn := func(v V, _, _ E) {
			value(v)
		}
		return []func(V, E, E){fn}, true

	case func():
		fn := func(V, E, E) {
			value()
		}
		return []func(V, E, E){fn}, true

	case []func(V, E, E):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(V, E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(view V, val, _ E) {
				fn(view, val)
			}
		}
		return listeners, true

	case []func(E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(_ V, val, _ E) {
				fn(val)
			}
		}
		return listeners, true

	case []func(E, E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(_ V, val, old E) {
				fn(val, old)
			}
		}
		return listeners, true

	case []func(V):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(view V, _, _ E) {
				fn(view)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(V, E, E) {
				fn()
			}
		}
		return listeners, true

	case []any:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch fn := v.(type) {
			case func(V, E, E):
				listeners[i] = fn

			case func(V, E):
				listeners[i] = func(view V, val, _ E) {
					fn(view, val)
				}

			case func(E, E):
				listeners[i] = func(_ V, val, old E) {
					fn(val, old)
				}

			case func(E):
				listeners[i] = func(_ V, val, _ E) {
					fn(val)
				}

			case func(V):
				listeners[i] = func(view V, _, _ E) {
					fn(view)
				}

			case func():
				listeners[i] = func(V, E, E) {
					fn()
				}

			default:
				return nil, false
			}
		}
		return listeners, true
	}

	return nil, false
}

var keyEvents = map[string]struct{ jsEvent, jsFunc string }{
	KeyDownEvent: {jsEvent: "onkeydown", jsFunc: "keyDownEvent"},
	KeyUpEvent:   {jsEvent: "onkeyup", jsFunc: "keyUpEvent"},
}

func (view *viewData) setKeyListener(tag string, value any) bool {
	listeners, ok := valueToEventListeners[View, KeyEvent](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeKeyListener(tag)
	} else if js, ok := keyEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			view.session.updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)")
		}
	} else {
		return false
	}
	return true
}

func (view *viewData) removeKeyListener(tag string) {
	delete(view.properties, tag)
	if view.created {
		if js, ok := keyEvents[tag]; ok {
			view.session.removeProperty(view.htmlID(), js.jsEvent)
		}
	}
}

func getEventWithOldListeners[V View, E any](view View, subviewID []string, tag string) []func(V, E, E) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(V, E, E)); ok {
				return result
			}
		}
	}
	return []func(V, E, E){}
}
func getEventListeners[V View, E any](view View, subviewID []string, tag string) []func(V, E) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(V, E)); ok {
				return result
			}
		}
	}
	return []func(V, E){}
}

func keyEventsHtml(view View, buffer *strings.Builder) {
	for tag, js := range keyEvents {
		if listeners := getEventListeners[View, KeyEvent](view, nil, tag); len(listeners) > 0 {
			buffer.WriteString(js.jsEvent + `="` + js.jsFunc + `(this, event)" `)
		}
	}
}

func handleKeyEvents(view View, tag string, data DataObject) {
	listeners := getEventListeners[View, KeyEvent](view, nil, tag)
	if len(listeners) > 0 {
		var event KeyEvent
		event.init(data)

		for _, listener := range listeners {
			listener(view, event)
		}
	}
}

// GetKeyDownListeners returns the "key-down-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetKeyDownListeners(view View, subviewID ...string) []func(View, KeyEvent) {
	return getEventListeners[View, KeyEvent](view, subviewID, KeyDownEvent)
}

// GetKeyUpListeners returns the "key-up-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetKeyUpListeners(view View, subviewID ...string) []func(View, KeyEvent) {
	return getEventListeners[View, KeyEvent](view, subviewID, KeyUpEvent)
}
