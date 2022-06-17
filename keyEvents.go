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

func valueToKeyListeners(value interface{}) ([]func(View, KeyEvent), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(View, KeyEvent):
		return []func(View, KeyEvent){value}, true

	case func(KeyEvent):
		fn := func(view View, event KeyEvent) {
			value(event)
		}
		return []func(View, KeyEvent){fn}, true

	case func(View):
		fn := func(view View, event KeyEvent) {
			value(view)
		}
		return []func(View, KeyEvent){fn}, true

	case func():
		fn := func(view View, event KeyEvent) {
			value()
		}
		return []func(View, KeyEvent){fn}, true

	case []func(View, KeyEvent):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(KeyEvent):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View, KeyEvent), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view View, event KeyEvent) {
				v(event)
			}
		}
		return listeners, true

	case []func(View):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View, KeyEvent), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view View, event KeyEvent) {
				v(view)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View, KeyEvent), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view View, event KeyEvent) {
				v()
			}
		}
		return listeners, true

	case []interface{}:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View, KeyEvent), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(View, KeyEvent):
				listeners[i] = v

			case func(KeyEvent):
				listeners[i] = func(view View, event KeyEvent) {
					v(event)
				}

			case func(View):
				listeners[i] = func(view View, event KeyEvent) {
					v(view)
				}

			case func():
				listeners[i] = func(view View, event KeyEvent) {
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

var keyEvents = map[string]struct{ jsEvent, jsFunc string }{
	KeyDownEvent: {jsEvent: "onkeydown", jsFunc: "keyDownEvent"},
	KeyUpEvent:   {jsEvent: "onkeyup", jsFunc: "keyUpEvent"},
}

func (view *viewData) setKeyListener(tag string, value interface{}) bool {
	listeners, ok := valueToKeyListeners(value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeKeyListener(tag)
	} else if js, ok := keyEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)", view.Session())
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
			removeProperty(view.htmlID(), js.jsEvent, view.Session())
		}
	}
}

func getKeyListeners(view View, subviewID string, tag string) []func(View, KeyEvent) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(View, KeyEvent)); ok {
				return result
			}
		}
	}
	return []func(View, KeyEvent){}
}

func keyEventsHtml(view View, buffer *strings.Builder) {
	for tag, js := range keyEvents {
		if listeners := getKeyListeners(view, "", tag); len(listeners) > 0 {
			buffer.WriteString(js.jsEvent + `="` + js.jsFunc + `(this, event)" `)
		}
	}
}

func handleKeyEvents(view View, tag string, data DataObject) {
	listeners := getKeyListeners(view, "", tag)
	if len(listeners) == 0 {
		return
	}

	getBool := func(tag string) bool {
		if value, ok := data.PropertyValue(tag); ok && value == "1" {
			return true
		}
		return false
	}

	key, _ := data.PropertyValue("key")
	code, _ := data.PropertyValue("code")
	event := KeyEvent{
		TimeStamp: getTimeStamp(data),
		Key:       key,
		Code:      code,
		Repeat:    getBool("repeat"),
		CtrlKey:   getBool("ctrlKey"),
		ShiftKey:  getBool("shiftKey"),
		AltKey:    getBool("altKey"),
		MetaKey:   getBool("metaKey"),
	}

	for _, listener := range listeners {
		listener(view, event)
	}
}

// GetKeyDownListeners returns the "key-down-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetKeyDownListeners(view View, subviewID string) []func(View, KeyEvent) {
	return getKeyListeners(view, subviewID, KeyDownEvent)
}

// GetKeyUpListeners returns the "key-up-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetKeyUpListeners(view View, subviewID string) []func(View, KeyEvent) {
	return getKeyListeners(view, subviewID, KeyUpEvent)
}
