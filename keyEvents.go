package rui

import "strings"

// Constants which represent [View] specific keyboard events properties
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

// ControlKeyMask represent ORed state of keyboard's control keys like [AltKey], [CtrlKey], [ShiftKey] and [MetaKey]
type ControlKeyMask int

// KeyCode is a string representation the a physical key being pressed.
// The value is not affected by the current keyboard layout or modifier state,
// so a particular key will always have the same value.
type KeyCode string

// Constants for specific keyboard keys.
const (
	// AltKey is the mask of the "alt" key
	AltKey ControlKeyMask = 1
	// CtrlKey is the mask of the "ctrl" key
	CtrlKey ControlKeyMask = 2
	// ShiftKey is the mask of the "shift" key
	ShiftKey ControlKeyMask = 4
	// MetaKey is the mask of the "meta" key
	MetaKey ControlKeyMask = 8

	// KeyA reresent "A" key on the keyboard
	KeyA KeyCode = "KeyA"

	// KeyB reresent "B" key on the keyboard
	KeyB KeyCode = "KeyB"

	// KeyC reresent "C" key on the keyboard
	KeyC KeyCode = "KeyC"

	// KeyD reresent "D" key on the keyboard
	KeyD KeyCode = "KeyD"

	// KeyE reresent "E" key on the keyboard
	KeyE KeyCode = "KeyE"

	// KeyF reresent "F" key on the keyboard
	KeyF KeyCode = "KeyF"

	// KeyG reresent "G" key on the keyboard
	KeyG KeyCode = "KeyG"

	// KeyH reresent "H" key on the keyboard
	KeyH KeyCode = "KeyH"

	// KeyI reresent "I" key on the keyboard
	KeyI KeyCode = "KeyI"

	// KeyJ reresent "J" key on the keyboard
	KeyJ KeyCode = "KeyJ"

	// KeyK reresent "K" key on the keyboard
	KeyK KeyCode = "KeyK"

	// KeyL reresent "L" key on the keyboard
	KeyL KeyCode = "KeyL"

	// KeyM reresent "M" key on the keyboard
	KeyM KeyCode = "KeyM"

	// KeyN reresent "N" key on the keyboard
	KeyN KeyCode = "KeyN"

	// KeyO reresent "O" key on the keyboard
	KeyO KeyCode = "KeyO"

	// KeyP reresent "P" key on the keyboard
	KeyP KeyCode = "KeyP"

	// KeyQ reresent "Q" key on the keyboard
	KeyQ KeyCode = "KeyQ"

	// KeyR reresent "R" key on the keyboard
	KeyR KeyCode = "KeyR"

	// KeyS reresent "S" key on the keyboard
	KeyS KeyCode = "KeyS"

	// KeyT reresent "T" key on the keyboard
	KeyT KeyCode = "KeyT"

	// KeyU reresent "U" key on the keyboard
	KeyU KeyCode = "KeyU"

	// KeyV reresent "V" key on the keyboard
	KeyV KeyCode = "KeyV"

	// KeyW reresent "W" key on the keyboard
	KeyW KeyCode = "KeyW"

	// KeyX reresent "X" key on the keyboard
	KeyX KeyCode = "KeyX"

	// KeyY reresent "Y" key on the keyboard
	KeyY KeyCode = "KeyY"

	// KeyZ reresent "Z" key on the keyboard
	KeyZ KeyCode = "KeyZ"

	// Digit0Key reresent "Digit0" key on the keyboard
	Digit0Key KeyCode = "Digit0"

	// Digit1Key reresent "Digit1" key on the keyboard
	Digit1Key KeyCode = "Digit1"

	// Digit2Key reresent "Digit2" key on the keyboard
	Digit2Key KeyCode = "Digit2"

	// Digit3Key reresent "Digit3" key on the keyboard
	Digit3Key KeyCode = "Digit3"

	// Digit4Key reresent "Digit4" key on the keyboard
	Digit4Key KeyCode = "Digit4"

	// Digit5Key reresent "Digit5" key on the keyboard
	Digit5Key KeyCode = "Digit5"

	// Digit6Key reresent "Digit6" key on the keyboard
	Digit6Key KeyCode = "Digit6"

	// Digit7Key reresent "Digit7" key on the keyboard
	Digit7Key KeyCode = "Digit7"

	// Digit8Key reresent "Digit8" key on the keyboard
	Digit8Key KeyCode = "Digit8"

	// Digit9Key reresent "Digit9" key on the keyboard
	Digit9Key KeyCode = "Digit9"

	// SpaceKey reresent "Space" key on the keyboard
	SpaceKey KeyCode = "Space"

	// MinusKey reresent "Minus" key on the keyboard
	MinusKey KeyCode = "Minus"

	// EqualKey reresent "Equal" key on the keyboard
	EqualKey KeyCode = "Equal"

	// IntlBackslashKey reresent "IntlBackslash" key on the keyboard
	IntlBackslashKey KeyCode = "IntlBackslash"

	// BracketLeftKey reresent "BracketLeft" key on the keyboard
	BracketLeftKey KeyCode = "BracketLeft"

	// BracketRightKey reresent "BracketRight" key on the keyboard
	BracketRightKey KeyCode = "BracketRight"

	// SemicolonKey reresent "Semicolon" key on the keyboard
	SemicolonKey KeyCode = "Semicolon"

	// CommaKey reresent "Comma" key on the keyboard
	CommaKey KeyCode = "Comma"

	// PeriodKey reresent "Period" key on the keyboard
	PeriodKey KeyCode = "Period"

	// QuoteKey reresent "Quote" key on the keyboard
	QuoteKey KeyCode = "Quote"

	// BackquoteKey reresent "Backquote" key on the keyboard
	BackquoteKey KeyCode = "Backquote"

	// SlashKey reresent "Slash" key on the keyboard
	SlashKey KeyCode = "Slash"

	// EscapeKey reresent "Escape" key on the keyboard
	EscapeKey KeyCode = "Escape"

	// EnterKey reresent "Enter" key on the keyboard
	EnterKey KeyCode = "Enter"

	// TabKey reresent "Tab" key on the keyboard
	TabKey KeyCode = "Tab"

	// CapsLockKey reresent "CapsLock" key on the keyboard
	CapsLockKey KeyCode = "CapsLock"

	// DeleteKey reresent "Delete" key on the keyboard
	DeleteKey KeyCode = "Delete"

	// InsertKey reresent "Insert" key on the keyboard
	InsertKey KeyCode = "Insert"

	// HelpKey reresent "Help" key on the keyboard
	HelpKey KeyCode = "Help"

	// BackspaceKey reresent "Backspace" key on the keyboard
	BackspaceKey KeyCode = "Backspace"

	// PrintScreenKey reresent "PrintScreen" key on the keyboard
	PrintScreenKey KeyCode = "PrintScreen"

	// ScrollLockKey reresent "ScrollLock" key on the keyboard
	ScrollLockKey KeyCode = "ScrollLock"

	// PauseKey reresent "Pause" key on the keyboard
	PauseKey KeyCode = "Pause"

	// ContextMenuKey reresent "ContextMenu" key on the keyboard
	ContextMenuKey KeyCode = "ContextMenu"

	// ArrowLeftKey reresent "ArrowLeft" key on the keyboard
	ArrowLeftKey KeyCode = "ArrowLeft"

	// ArrowRightKey reresent "ArrowRight" key on the keyboard
	ArrowRightKey KeyCode = "ArrowRight"

	// ArrowUpKey reresent "ArrowUp" key on the keyboard
	ArrowUpKey KeyCode = "ArrowUp"

	// ArrowDownKey reresent "ArrowDown" key on the keyboard
	ArrowDownKey KeyCode = "ArrowDown"

	// HomeKey reresent "Home" key on the keyboard
	HomeKey KeyCode = "Home"

	// EndKey reresent "End" key on the keyboard
	EndKey KeyCode = "End"

	// PageUpKey reresent "PageUp" key on the keyboard
	PageUpKey KeyCode = "PageUp"

	// PageDownKey reresent "PageDown" key on the keyboard
	PageDownKey KeyCode = "PageDown"

	// F1Key reresent "F1" key on the keyboard
	F1Key KeyCode = "F1"

	// F2Key reresent "F2" key on the keyboard
	F2Key KeyCode = "F2"

	// F3Key reresent "F3" key on the keyboard
	F3Key KeyCode = "F3"

	// F4Key reresent "F4" key on the keyboard
	F4Key KeyCode = "F4"

	// F5Key reresent "F5" key on the keyboard
	F5Key KeyCode = "F5"

	// F6Key reresent "F6" key on the keyboard
	F6Key KeyCode = "F6"

	// F7Key reresent "F7" key on the keyboard
	F7Key KeyCode = "F7"

	// F8Key reresent "F8" key on the keyboard
	F8Key KeyCode = "F8"

	// F9Key reresent "F9" key on the keyboard
	F9Key KeyCode = "F9"

	// F10Key reresent "F10" key on the keyboard
	F10Key KeyCode = "F10"

	// F11Key reresent "F11" key on the keyboard
	F11Key KeyCode = "F11"

	// F12Key reresent "F12" key on the keyboard
	F12Key KeyCode = "F12"

	// F13Key reresent "F13" key on the keyboard
	F13Key KeyCode = "F13"

	// NumLockKey reresent "NumLock" key on the keyboard
	NumLockKey KeyCode = "NumLock"

	// NumpadKey0 reresent "Numpad0" key on the keyboard
	NumpadKey0 KeyCode = "Numpad0"

	// NumpadKey1 reresent "Numpad1" key on the keyboard
	NumpadKey1 KeyCode = "Numpad1"

	// NumpadKey2 reresent "Numpad2" key on the keyboard
	NumpadKey2 KeyCode = "Numpad2"

	// NumpadKey3 reresent "Numpad3" key on the keyboard
	NumpadKey3 KeyCode = "Numpad3"

	// NumpadKey4 reresent "Numpad4" key on the keyboard
	NumpadKey4 KeyCode = "Numpad4"

	// NumpadKey5 reresent "Numpad5" key on the keyboard
	NumpadKey5 KeyCode = "Numpad5"

	// NumpadKey6 reresent "Numpad6" key on the keyboard
	NumpadKey6 KeyCode = "Numpad6"

	// NumpadKey7 reresent "Numpad7" key on the keyboard
	NumpadKey7 KeyCode = "Numpad7"

	// NumpadKey8 reresent "Numpad8" key on the keyboard
	NumpadKey8 KeyCode = "Numpad8"

	// NumpadKey9 reresent "Numpad9" key on the keyboard
	NumpadKey9 KeyCode = "Numpad9"

	// NumpadDecimalKey reresent "NumpadDecimal" key on the keyboard
	NumpadDecimalKey KeyCode = "NumpadDecimal"

	// NumpadEnterKey reresent "NumpadEnter" key on the keyboard
	NumpadEnterKey KeyCode = "NumpadEnter"

	// NumpadAddKey reresent "NumpadAdd" key on the keyboard
	NumpadAddKey KeyCode = "NumpadAdd"

	// NumpadSubtractKey reresent "NumpadSubtract" key on the keyboard
	NumpadSubtractKey KeyCode = "NumpadSubtract"

	// NumpadMultiplyKey reresent "NumpadMultiply" key on the keyboard
	NumpadMultiplyKey KeyCode = "NumpadMultiply"

	// NumpadDivideKey reresent "NumpadDivide" key on the keyboard
	NumpadDivideKey KeyCode = "NumpadDivide"

	// ShiftLeftKey reresent "ShiftLeft" key on the keyboard
	ShiftLeftKey KeyCode = "ShiftLeft"

	// ShiftRightKey reresent "ShiftRight" key on the keyboard
	ShiftRightKey KeyCode = "ShiftRight"

	// ControlLeftKey reresent "ControlLeft" key on the keyboard
	ControlLeftKey KeyCode = "ControlLeft"

	// ControlRightKey reresent "ControlRight" key on the keyboard
	ControlRightKey KeyCode = "ControlRight"

	// AltLeftKey reresent "AltLeft" key on the keyboard
	AltLeftKey KeyCode = "AltLeft"

	// AltRightKey reresent "AltRight" key on the keyboard
	AltRightKey KeyCode = "AltRight"

	// MetaLeftKey reresent "MetaLeft" key on the keyboard
	MetaLeftKey KeyCode = "MetaLeft"

	// MetaRightKey reresent "MetaRight" key on the keyboard
	MetaRightKey KeyCode = "MetaRight"
)

// KeyEvent represent a keyboard event
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
	Code KeyCode

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
	code, _ := data.PropertyValue("code")
	event.Code = KeyCode(code)
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

func (view *viewData) setKeyListener(tag string, value any) bool {
	listeners, ok := valueToEventListeners[View, KeyEvent](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeKeyListener(tag)
	} else {
		switch tag {
		case KeyDownEvent:
			view.properties[tag] = listeners
			if view.created {
				view.session.updateProperty(view.htmlID(), "onkeydown", "keyDownEvent(this, event)")
			}

		case KeyUpEvent:
			view.properties[tag] = listeners
			if view.created {
				view.session.updateProperty(view.htmlID(), "onkeyup", "keyUpEvent(this, event)")
			}

		default:
			return false
		}
	}

	return true
}

func (view *viewData) removeKeyListener(tag string) {
	delete(view.properties, tag)
	if view.created {
		switch tag {
		case KeyDownEvent:
			if !view.Focusable() {
				view.session.removeProperty(view.htmlID(), "onkeydown")
			}

		case KeyUpEvent:
			view.session.removeProperty(view.htmlID(), "onkeyup")
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
	if len(getEventListeners[View, KeyEvent](view, nil, KeyDownEvent)) > 0 {
		buffer.WriteString(`onkeydown="keyDownEvent(this, event)" `)
	} else if view.Focusable() {
		if len(getEventListeners[View, MouseEvent](view, nil, ClickEvent)) > 0 {
			buffer.WriteString(`onkeydown="keyDownEvent(this, event)" `)
		}
	}

	if listeners := getEventListeners[View, KeyEvent](view, nil, KeyUpEvent); len(listeners) > 0 {
		buffer.WriteString(`onkeyup="keyUpEvent(this, event)" `)
	}
}

func handleKeyEvents(view View, tag string, data DataObject) {
	var event KeyEvent
	event.init(data)
	listeners := getEventListeners[View, KeyEvent](view, nil, tag)

	if len(listeners) > 0 {
		for _, listener := range listeners {
			listener(view, event)
		}
		return
	}

	if tag == KeyDownEvent && view.Focusable() && (event.Key == " " || event.Key == "Enter") && !IsDisabled(view) {
		if listeners := getEventListeners[View, MouseEvent](view, nil, ClickEvent); len(listeners) > 0 {
			clickEvent := MouseEvent{
				TimeStamp: event.TimeStamp,
				Button:    PrimaryMouseButton,
				Buttons:   PrimaryMouseMask,
				CtrlKey:   event.CtrlKey,
				AltKey:    event.AltKey,
				ShiftKey:  event.ShiftKey,
				MetaKey:   event.MetaKey,
				ClientX:   view.Frame().Width / 2,
				ClientY:   view.Frame().Height / 2,
				X:         view.Frame().Width / 2,
				Y:         view.Frame().Height / 2,
				ScreenX:   view.Frame().Left + view.Frame().Width/2,
				ScreenY:   view.Frame().Top + view.Frame().Height/2,
			}
			for _, listener := range listeners {
				listener(view, clickEvent)
			}
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
