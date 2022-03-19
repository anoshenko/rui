package rui

import "strings"

const (
	// FocusEvent is the constant for "focus-event" property tag.
	// The "focus-event" event occurs when the View takes input focus.
	// The main listener format:
	//   func(View).
	// The additional listener format:
	//   func().
	FocusEvent = "focus-event"

	// LostFocusEvent is the constant for "lost-focus-event" property tag.
	// The "lost-focus-event" event occurs when the View lost input focus.
	// The main listener format:
	//   func(View).
	// The additional listener format:
	//   func().
	LostFocusEvent = "lost-focus-event"
)

func valueToFocusListeners(value interface{}) ([]func(View), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(View):
		return []func(View){value}, true

	case func():
		fn := func(View) {
			value()
		}
		return []func(View){fn}, true

	case []func(View):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(View) {
				v()
			}
		}
		return listeners, true

	case []interface{}:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(View), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(View):
				listeners[i] = v

			case func():
				listeners[i] = func(View) {
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

var focusEvents = map[string]struct{ jsEvent, jsFunc string }{
	FocusEvent:     {jsEvent: "onfocus", jsFunc: "focusEvent"},
	LostFocusEvent: {jsEvent: "onblur", jsFunc: "blurEvent"},
}

func (view *viewData) setFocusListener(tag string, value interface{}) bool {
	listeners, ok := valueToFocusListeners(value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeFocusListener(tag)
	} else if js, ok := focusEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)", view.Session())
		}
	} else {
		return false
	}
	return true
}

func (view *viewData) removeFocusListener(tag string) {
	delete(view.properties, tag)
	if view.created {
		if js, ok := focusEvents[tag]; ok {
			updateProperty(view.htmlID(), js.jsEvent, "", view.Session())
		}
	}
}

func getFocusListeners(view View, subviewID string, tag string) []func(View) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(View)); ok {
				return result
			}
		}
	}
	return []func(View){}
}

func focusEventsHtml(view View, buffer *strings.Builder) {
	if view.Focusable() {
		for _, js := range focusEvents {
			buffer.WriteString(js.jsEvent + `="` + js.jsFunc + `(this, event)" `)
		}
	}
}

// GetFocusListeners returns a FocusListener list. If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetFocusListeners(view View, subviewID string) []func(View) {
	return getFocusListeners(view, subviewID, FocusEvent)
}

// GetLostFocusListeners returns a LostFocusListener list. If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetLostFocusListeners(view View, subviewID string) []func(View) {
	return getFocusListeners(view, subviewID, LostFocusEvent)
}
