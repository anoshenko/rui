package rui

import "strings"

// Constants which represent [View] specific focus events properties
const (
	// FocusEvent is the constant for "focus-event" property tag.
	//
	// Used by `View`.
	// Occur when the view takes input focus.
	//
	// General listener format:
	// `func(View)`.
	//
	// where:
	// view - Interface of a view which generated this event.
	//
	// Allowed listener formats:
	// `func()`.
	FocusEvent = "focus-event"

	// LostFocusEvent is the constant for "lost-focus-event" property tag.
	//
	// Used by `View`.
	// Occur when the View lost input focus.
	//
	// General listener format:
	// `func(view rui.View)`.
	//
	// where:
	// view - Interface of a view which generated this event.
	//
	// Allowed listener formats:
	// `func()`.
	LostFocusEvent = "lost-focus-event"
)

func valueToNoParamListeners[V any](value any) ([]func(V), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(V):
		return []func(V){value}, true

	case func():
		fn := func(V) {
			value()
		}
		return []func(V){fn}, true

	case []func(V):
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
		listeners := make([]func(V), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(V) {
				v()
			}
		}
		return listeners, true

	case []any:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(V):
				listeners[i] = v

			case func():
				listeners[i] = func(V) {
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

func (view *viewData) setFocusListener(tag string, value any) bool {
	listeners, ok := valueToNoParamListeners[View](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeFocusListener(tag)
	} else if js, ok := focusEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			view.session.updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)")
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
			view.session.removeProperty(view.htmlID(), js.jsEvent)
		}
	}
}

func getFocusListeners(view View, subviewID []string, tag string) []func(View) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
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
			buffer.WriteString(js.jsEvent)
			buffer.WriteString(`="`)
			buffer.WriteString(js.jsFunc)
			buffer.WriteString(`(this, event)" `)
		}
	}
}

// GetFocusListeners returns a FocusListener list. If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetFocusListeners(view View, subviewID ...string) []func(View) {
	return getFocusListeners(view, subviewID, FocusEvent)
}

// GetLostFocusListeners returns a LostFocusListener list. If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetLostFocusListeners(view View, subviewID ...string) []func(View) {
	return getFocusListeners(view, subviewID, LostFocusEvent)
}
