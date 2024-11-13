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
	FocusEvent PropertyName = "focus-event"

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
	LostFocusEvent PropertyName = "lost-focus-event"
)

func focusEventsHtml(view View, buffer *strings.Builder) {
	if view.Focusable() {
		for _, tag := range []PropertyName{FocusEvent, LostFocusEvent} {
			if js, ok := eventJsFunc[tag]; ok {
				buffer.WriteString(js.jsEvent)
				buffer.WriteString(`="`)
				buffer.WriteString(js.jsFunc)
				buffer.WriteString(`(this, event)" `)
			}
		}
	}
}

// GetFocusListeners returns a FocusListener list. If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetFocusListeners(view View, subviewID ...string) []func(View) {
	return getNoParamEventListeners[View](view, subviewID, FocusEvent)
}

// GetLostFocusListeners returns a LostFocusListener list. If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetLostFocusListeners(view View, subviewID ...string) []func(View) {
	return getNoParamEventListeners[View](view, subviewID, LostFocusEvent)
}
