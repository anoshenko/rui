package rui

import "strings"

// Constants which represent [View] specific keyboard events properties
const (
	// KeyDownEvent is the constant for "key-down-event" property tag.
	//
	// Used by View.
	// Is fired when a key is pressed.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.KeyEvent).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Key event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.KeyEvent)
	//  func()
	KeyDownEvent PropertyName = "key-down-event"

	// KeyUpEvent is the constant for "key-up-event" property tag.
	//
	// Used by View.
	// Is fired when a key is released.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.KeyEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Key event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.KeyEvent)
	//  func()
	KeyUpEvent PropertyName = "key-up-event"
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

	// KeyA represent "A" key on the keyboard
	KeyA KeyCode = "KeyA"

	// KeyB represent "B" key on the keyboard
	KeyB KeyCode = "KeyB"

	// KeyC represent "C" key on the keyboard
	KeyC KeyCode = "KeyC"

	// KeyD represent "D" key on the keyboard
	KeyD KeyCode = "KeyD"

	// KeyE represent "E" key on the keyboard
	KeyE KeyCode = "KeyE"

	// KeyF represent "F" key on the keyboard
	KeyF KeyCode = "KeyF"

	// KeyG represent "G" key on the keyboard
	KeyG KeyCode = "KeyG"

	// KeyH represent "H" key on the keyboard
	KeyH KeyCode = "KeyH"

	// KeyI represent "I" key on the keyboard
	KeyI KeyCode = "KeyI"

	// KeyJ represent "J" key on the keyboard
	KeyJ KeyCode = "KeyJ"

	// KeyK represent "K" key on the keyboard
	KeyK KeyCode = "KeyK"

	// KeyL represent "L" key on the keyboard
	KeyL KeyCode = "KeyL"

	// KeyM represent "M" key on the keyboard
	KeyM KeyCode = "KeyM"

	// KeyN represent "N" key on the keyboard
	KeyN KeyCode = "KeyN"

	// KeyO represent "O" key on the keyboard
	KeyO KeyCode = "KeyO"

	// KeyP represent "P" key on the keyboard
	KeyP KeyCode = "KeyP"

	// KeyQ represent "Q" key on the keyboard
	KeyQ KeyCode = "KeyQ"

	// KeyR represent "R" key on the keyboard
	KeyR KeyCode = "KeyR"

	// KeyS represent "S" key on the keyboard
	KeyS KeyCode = "KeyS"

	// KeyT represent "T" key on the keyboard
	KeyT KeyCode = "KeyT"

	// KeyU represent "U" key on the keyboard
	KeyU KeyCode = "KeyU"

	// KeyV represent "V" key on the keyboard
	KeyV KeyCode = "KeyV"

	// KeyW represent "W" key on the keyboard
	KeyW KeyCode = "KeyW"

	// KeyX represent "X" key on the keyboard
	KeyX KeyCode = "KeyX"

	// KeyY represent "Y" key on the keyboard
	KeyY KeyCode = "KeyY"

	// KeyZ represent "Z" key on the keyboard
	KeyZ KeyCode = "KeyZ"

	// Digit0Key represent "Digit0" key on the keyboard
	Digit0Key KeyCode = "Digit0"

	// Digit1Key represent "Digit1" key on the keyboard
	Digit1Key KeyCode = "Digit1"

	// Digit2Key represent "Digit2" key on the keyboard
	Digit2Key KeyCode = "Digit2"

	// Digit3Key represent "Digit3" key on the keyboard
	Digit3Key KeyCode = "Digit3"

	// Digit4Key represent "Digit4" key on the keyboard
	Digit4Key KeyCode = "Digit4"

	// Digit5Key represent "Digit5" key on the keyboard
	Digit5Key KeyCode = "Digit5"

	// Digit6Key represent "Digit6" key on the keyboard
	Digit6Key KeyCode = "Digit6"

	// Digit7Key represent "Digit7" key on the keyboard
	Digit7Key KeyCode = "Digit7"

	// Digit8Key represent "Digit8" key on the keyboard
	Digit8Key KeyCode = "Digit8"

	// Digit9Key represent "Digit9" key on the keyboard
	Digit9Key KeyCode = "Digit9"

	// SpaceKey represent "Space" key on the keyboard
	SpaceKey KeyCode = "Space"

	// MinusKey represent "Minus" key on the keyboard
	MinusKey KeyCode = "Minus"

	// EqualKey represent "Equal" key on the keyboard
	EqualKey KeyCode = "Equal"

	// IntlBackslashKey represent "IntlBackslash" key on the keyboard
	IntlBackslashKey KeyCode = "IntlBackslash"

	// BracketLeftKey represent "BracketLeft" key on the keyboard
	BracketLeftKey KeyCode = "BracketLeft"

	// BracketRightKey represent "BracketRight" key on the keyboard
	BracketRightKey KeyCode = "BracketRight"

	// SemicolonKey represent "Semicolon" key on the keyboard
	SemicolonKey KeyCode = "Semicolon"

	// CommaKey represent "Comma" key on the keyboard
	CommaKey KeyCode = "Comma"

	// PeriodKey represent "Period" key on the keyboard
	PeriodKey KeyCode = "Period"

	// QuoteKey represent "Quote" key on the keyboard
	QuoteKey KeyCode = "Quote"

	// BackquoteKey represent "Backquote" key on the keyboard
	BackquoteKey KeyCode = "Backquote"

	// SlashKey represent "Slash" key on the keyboard
	SlashKey KeyCode = "Slash"

	// EscapeKey represent "Escape" key on the keyboard
	EscapeKey KeyCode = "Escape"

	// EnterKey represent "Enter" key on the keyboard
	EnterKey KeyCode = "Enter"

	// TabKey represent "Tab" key on the keyboard
	TabKey KeyCode = "Tab"

	// CapsLockKey represent "CapsLock" key on the keyboard
	CapsLockKey KeyCode = "CapsLock"

	// DeleteKey represent "Delete" key on the keyboard
	DeleteKey KeyCode = "Delete"

	// InsertKey represent "Insert" key on the keyboard
	InsertKey KeyCode = "Insert"

	// HelpKey represent "Help" key on the keyboard
	HelpKey KeyCode = "Help"

	// BackspaceKey represent "Backspace" key on the keyboard
	BackspaceKey KeyCode = "Backspace"

	// PrintScreenKey represent "PrintScreen" key on the keyboard
	PrintScreenKey KeyCode = "PrintScreen"

	// ScrollLockKey represent "ScrollLock" key on the keyboard
	ScrollLockKey KeyCode = "ScrollLock"

	// PauseKey represent "Pause" key on the keyboard
	PauseKey KeyCode = "Pause"

	// ContextMenuKey represent "ContextMenu" key on the keyboard
	ContextMenuKey KeyCode = "ContextMenu"

	// ArrowLeftKey represent "ArrowLeft" key on the keyboard
	ArrowLeftKey KeyCode = "ArrowLeft"

	// ArrowRightKey represent "ArrowRight" key on the keyboard
	ArrowRightKey KeyCode = "ArrowRight"

	// ArrowUpKey represent "ArrowUp" key on the keyboard
	ArrowUpKey KeyCode = "ArrowUp"

	// ArrowDownKey represent "ArrowDown" key on the keyboard
	ArrowDownKey KeyCode = "ArrowDown"

	// HomeKey represent "Home" key on the keyboard
	HomeKey KeyCode = "Home"

	// EndKey represent "End" key on the keyboard
	EndKey KeyCode = "End"

	// PageUpKey represent "PageUp" key on the keyboard
	PageUpKey KeyCode = "PageUp"

	// PageDownKey represent "PageDown" key on the keyboard
	PageDownKey KeyCode = "PageDown"

	// F1Key represent "F1" key on the keyboard
	F1Key KeyCode = "F1"

	// F2Key represent "F2" key on the keyboard
	F2Key KeyCode = "F2"

	// F3Key represent "F3" key on the keyboard
	F3Key KeyCode = "F3"

	// F4Key represent "F4" key on the keyboard
	F4Key KeyCode = "F4"

	// F5Key represent "F5" key on the keyboard
	F5Key KeyCode = "F5"

	// F6Key represent "F6" key on the keyboard
	F6Key KeyCode = "F6"

	// F7Key represent "F7" key on the keyboard
	F7Key KeyCode = "F7"

	// F8Key represent "F8" key on the keyboard
	F8Key KeyCode = "F8"

	// F9Key represent "F9" key on the keyboard
	F9Key KeyCode = "F9"

	// F10Key represent "F10" key on the keyboard
	F10Key KeyCode = "F10"

	// F11Key represent "F11" key on the keyboard
	F11Key KeyCode = "F11"

	// F12Key represent "F12" key on the keyboard
	F12Key KeyCode = "F12"

	// F13Key represent "F13" key on the keyboard
	F13Key KeyCode = "F13"

	// NumLockKey represent "NumLock" key on the keyboard
	NumLockKey KeyCode = "NumLock"

	// NumpadKey0 represent "Numpad0" key on the keyboard
	NumpadKey0 KeyCode = "Numpad0"

	// NumpadKey1 represent "Numpad1" key on the keyboard
	NumpadKey1 KeyCode = "Numpad1"

	// NumpadKey2 represent "Numpad2" key on the keyboard
	NumpadKey2 KeyCode = "Numpad2"

	// NumpadKey3 represent "Numpad3" key on the keyboard
	NumpadKey3 KeyCode = "Numpad3"

	// NumpadKey4 represent "Numpad4" key on the keyboard
	NumpadKey4 KeyCode = "Numpad4"

	// NumpadKey5 represent "Numpad5" key on the keyboard
	NumpadKey5 KeyCode = "Numpad5"

	// NumpadKey6 represent "Numpad6" key on the keyboard
	NumpadKey6 KeyCode = "Numpad6"

	// NumpadKey7 represent "Numpad7" key on the keyboard
	NumpadKey7 KeyCode = "Numpad7"

	// NumpadKey8 represent "Numpad8" key on the keyboard
	NumpadKey8 KeyCode = "Numpad8"

	// NumpadKey9 represent "Numpad9" key on the keyboard
	NumpadKey9 KeyCode = "Numpad9"

	// NumpadDecimalKey represent "NumpadDecimal" key on the keyboard
	NumpadDecimalKey KeyCode = "NumpadDecimal"

	// NumpadEnterKey represent "NumpadEnter" key on the keyboard
	NumpadEnterKey KeyCode = "NumpadEnter"

	// NumpadAddKey represent "NumpadAdd" key on the keyboard
	NumpadAddKey KeyCode = "NumpadAdd"

	// NumpadSubtractKey represent "NumpadSubtract" key on the keyboard
	NumpadSubtractKey KeyCode = "NumpadSubtract"

	// NumpadMultiplyKey represent "NumpadMultiply" key on the keyboard
	NumpadMultiplyKey KeyCode = "NumpadMultiply"

	// NumpadDivideKey represent "NumpadDivide" key on the keyboard
	NumpadDivideKey KeyCode = "NumpadDivide"

	// ShiftLeftKey represent "ShiftLeft" key on the keyboard
	ShiftLeftKey KeyCode = "ShiftLeft"

	// ShiftRightKey represent "ShiftRight" key on the keyboard
	ShiftRightKey KeyCode = "ShiftRight"

	// ControlLeftKey represent "ControlLeft" key on the keyboard
	ControlLeftKey KeyCode = "ControlLeft"

	// ControlRightKey represent "ControlRight" key on the keyboard
	ControlRightKey KeyCode = "ControlRight"

	// AltLeftKey represent "AltLeft" key on the keyboard
	AltLeftKey KeyCode = "AltLeft"

	// AltRightKey represent "AltRight" key on the keyboard
	AltRightKey KeyCode = "AltRight"

	// MetaLeftKey represent "MetaLeft" key on the keyboard
	MetaLeftKey KeyCode = "MetaLeft"

	// MetaRightKey represent "MetaRight" key on the keyboard
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

func keyEventsHtml(view View, buffer *strings.Builder) {
	if len(getOneArgEventListeners[View, KeyEvent](view, nil, KeyDownEvent)) > 0 ||
		(view.Focusable() && len(getOneArgEventListeners[View, MouseEvent](view, nil, ClickEvent)) > 0) {

		buffer.WriteString(`onkeydown="keyDownEvent(this, event)" `)
	}

	if len(getOneArgEventListeners[View, KeyEvent](view, nil, KeyUpEvent)) > 0 {
		buffer.WriteString(`onkeyup="keyUpEvent(this, event)" `)
	}
}

func handleKeyEvents(view View, tag PropertyName, data DataObject) {
	var event KeyEvent
	event.init(data)
	listeners := getOneArgEventListeners[View, KeyEvent](view, nil, tag)

	if len(listeners) > 0 {
		for _, listener := range listeners {
			listener.Run(view, event)
		}
		return
	}

	if tag == KeyDownEvent && view.Focusable() && (event.Key == " " || event.Key == "Enter") &&
		!IsDisabled(view) && GetSemantics(view) != ButtonSemantics {

		switch view.Tag() {
		case "EditView", "ListView", "TableView", "TabsLayout", "TimePicker", "DatePicker", "AudioPlayer", "VideoPlayer":
			return
		}

		if listeners := getOneArgEventListeners[View, MouseEvent](view, nil, ClickEvent); len(listeners) > 0 {
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
				listener.Run(view, clickEvent)
			}
		}
	}
}

// GetKeyDownListeners returns the "key-down-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(rui.View, rui.KeyEvent),
//   - func(rui.View),
//   - func(rui.KeyEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetKeyDownListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, KeyEvent](view, subviewID, KeyDownEvent)
}

// GetKeyUpListeners returns the "key-up-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(rui.View, rui.KeyEvent),
//   - func(rui.View),
//   - func(rui.KeyEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetKeyUpListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, KeyEvent](view, subviewID, KeyUpEvent)
}
