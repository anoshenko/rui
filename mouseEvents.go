package rui

import (
	"strconv"
	"strings"
)

// Constants related to [View] mouse events properties
const (
	// ClickEvent is the constant for "click-event" property tag.
	//
	// Used by View.
	// Occur when the user clicks on the view.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.MouseEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Mouse event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.MouseEvent)
	//  func()
	ClickEvent PropertyName = "click-event"

	// DoubleClickEvent is the constant for "double-click-event" property tag.
	//
	// Used by View.
	// Occur when the user double clicks on the view.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.MouseEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Mouse event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.MouseEvent)
	//  func()
	DoubleClickEvent PropertyName = "double-click-event"

	// MouseDown is the constant for "mouse-down" property tag.
	//
	// Used by View.
	// Is fired at a View when a pointing device button is pressed while the pointer is inside the view.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.MouseEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Mouse event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.MouseEvent)
	//  func()
	MouseDown PropertyName = "mouse-down"

	// MouseUp is the constant for "mouse-up" property tag.
	//
	// Used by View.
	// Is fired at a View when a button on a pointing device (such as a mouse or trackpad) is released while the pointer is
	// located inside it. "mouse-up" events are the counterpoint to "mouse-down" events.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.MouseEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Mouse event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.MouseEvent)
	//  func()
	MouseUp PropertyName = "mouse-up"

	// MouseMove is the constant for "mouse-move" property tag.
	//
	// Used by View.
	// Is fired at a view when a pointing device(usually a mouse) is moved while the cursor's hotspot is inside it.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.MouseEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Mouse event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.MouseEvent)
	//  func()
	MouseMove PropertyName = "mouse-move"

	// MouseOut is the constant for "mouse-out" property tag.
	//
	// Used by View.
	// Is fired at a View when a pointing device (usually a mouse) is used to move the cursor so that it is no longer
	// contained within the view or one of its children. "mouse-out" is also delivered to a view if the cursor enters a child
	// view, because the child view obscures the visible area of the view.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.MouseEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Mouse event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.MouseEvent)
	//  func()
	MouseOut PropertyName = "mouse-out"

	// MouseOver is the constant for "mouse-over" property tag.
	//
	// Used by View.
	// Is fired at a View when a pointing device (such as a mouse or trackpad) is used to move the cursor onto the view or one
	// of its child views.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.MouseEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Mouse event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.MouseEvent)
	//  func()
	MouseOver PropertyName = "mouse-over"

	// ContextMenuEvent is the constant for "context-menu-event" property tag.
	//
	// Used by View.
	// Occur when the user calls the context menu by the right mouse clicking.
	//
	// General listener format:
	//
	//  func(view rui.View, event rui.MouseEvent)
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - event - Mouse event.
	//
	// Allowed listener formats:
	//
	//  func(view rui.View)
	//  func(event rui.MouseEvent)
	//  func()
	ContextMenuEvent PropertyName = "context-menu-event"

	// PrimaryMouseButton is a number of the main pressed button, usually the left button or the un-initialized state
	PrimaryMouseButton = 0

	// AuxiliaryMouseButton is a number of the auxiliary pressed button, usually the wheel button
	// or the middle button (if present)
	AuxiliaryMouseButton = 1

	// SecondaryMouseButton is a number of the secondary pressed button, usually the right button
	SecondaryMouseButton = 2

	// MouseButton4 is a number of the fourth button, typically the Browser Back button
	MouseButton4 = 3

	// MouseButton5 is a number of the fifth button, typically the Browser Forward button
	MouseButton5 = 4

	// PrimaryMouseMask is the mask of the primary button (usually the left button)
	PrimaryMouseMask = 1

	// SecondaryMouseMask is the mask of the secondary button (usually the right button)
	SecondaryMouseMask = 2

	// AuxiliaryMouseMask  is the mask of the auxiliary button (usually the mouse wheel button or middle button)
	AuxiliaryMouseMask = 4

	// MouseMask4 is the mask of the 4th button (typically the "Browser Back" button)
	MouseMask4 = 8

	//MouseMask5 is the mask of the  5th button (typically the "Browser Forward" button)
	MouseMask5 = 16
)

// MouseEvent represent a mouse event
type MouseEvent struct {
	// TimeStamp is the time at which the event was created (in milliseconds).
	// This value is time since epoch—but in reality, browsers' definitions vary.
	TimeStamp uint64

	// Button indicates which button was pressed on the mouse to trigger the event:
	// PrimaryMouseButton (0), AuxiliaryMouseButton (1), SecondaryMouseButton (2),
	// MouseButton4 (3), and MouseButton5 (4)
	Button int

	// Buttons indicates which buttons are pressed on the mouse (or other input device)
	// when a mouse event is triggered. Each button that can be pressed is represented by a given mask:
	// PrimaryMouseMask (1), SecondaryMouseMask (2), AuxiliaryMouseMask (4), MouseMask4 (8), and MouseMask5 (16)
	Buttons int

	// X provides the horizontal coordinate within the view's viewport.
	X float64
	// Y provides the vertical coordinate within the view's viewport.
	Y float64

	// ClientX provides the horizontal coordinate within the application's viewport at which the event occurred.
	ClientX float64
	// ClientY provides the vertical coordinate within the application's viewport at which the event occurred.
	ClientY float64

	// ScreenX provides the horizontal coordinate (offset) of the mouse pointer in global (screen) coordinates.
	ScreenX float64
	// ScreenY provides the vertical coordinate (offset) of the mouse pointer in global (screen) coordinates.
	ScreenY float64

	// CtrlKey == true if the control key was down when the event was fired. false otherwise.
	CtrlKey bool
	// ShiftKey == true if the shift key was down when the event was fired. false otherwise.
	ShiftKey bool
	// AltKey == true if the alt key was down when the event was fired. false otherwise.
	AltKey bool
	// MetaKey == true if the meta key was down when the event was fired. false otherwise.
	MetaKey bool
}

func getTimeStamp(data DataObject) uint64 {
	if value, ok := data.PropertyValue("timeStamp"); ok {
		if index := strings.Index(value, "."); index > 0 {
			value = value[:index]
		}
		if n, err := strconv.ParseUint(value, 10, 64); err == nil {
			return n
		}
	}
	return 0
}

func (event *MouseEvent) init(data DataObject) {

	event.TimeStamp = getTimeStamp(data)
	event.Button, _ = dataIntProperty(data, "button")
	event.Buttons, _ = dataIntProperty(data, "buttons")
	event.X = dataFloatProperty(data, "x")
	event.Y = dataFloatProperty(data, "y")
	event.ClientX = dataFloatProperty(data, "clientX")
	event.ClientY = dataFloatProperty(data, "clientY")
	event.ScreenX = dataFloatProperty(data, "screenX")
	event.ScreenY = dataFloatProperty(data, "screenY")
	event.CtrlKey = dataBoolProperty(data, "ctrlKey")
	event.ShiftKey = dataBoolProperty(data, "shiftKey")
	event.AltKey = dataBoolProperty(data, "altKey")
	event.MetaKey = dataBoolProperty(data, "metaKey")
}

func handleMouseEvents(view View, tag PropertyName, data DataObject) {
	listeners := getOneArgEventListeners[View, MouseEvent](view, nil, tag)
	if len(listeners) > 0 {
		var event MouseEvent
		event.init(data)

		for _, listener := range listeners {
			listener.Run(view, event)
		}
	}
}

// GetClickListeners returns the "click-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(View, MouseEvent),
//   - func(View),
//   - func(MouseEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetClickListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, MouseEvent](view, subviewID, ClickEvent)
}

// GetDoubleClickListeners returns the "double-click-event" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(View, MouseEvent),
//   - func(View),
//   - func(MouseEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDoubleClickListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, MouseEvent](view, subviewID, DoubleClickEvent)
}

// GetContextMenuListeners returns the "context-menu" listener list.
// If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(View, MouseEvent),
//   - func(View),
//   - func(MouseEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetContextMenuListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, MouseEvent](view, subviewID, ContextMenuEvent)
}

// GetMouseDownListeners returns the "mouse-down" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(View, MouseEvent),
//   - func(View),
//   - func(MouseEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetMouseDownListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, MouseEvent](view, subviewID, MouseDown)
}

// GetMouseUpListeners returns the "mouse-up" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(View, MouseEvent),
//   - func(View),
//   - func(MouseEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetMouseUpListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, MouseEvent](view, subviewID, MouseUp)
}

// GetMouseMoveListeners returns the "mouse-move" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(View, MouseEvent),
//   - func(View),
//   - func(MouseEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetMouseMoveListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, MouseEvent](view, subviewID, MouseMove)
}

// GetMouseOverListeners returns the "mouse-over" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(View, MouseEvent),
//   - func(View),
//   - func(MouseEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetMouseOverListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, MouseEvent](view, subviewID, MouseOver)
}

// GetMouseOutListeners returns the "mouse-out" listener list. If there are no listeners then the empty list is returned.
//
// Result elements can be of the following types:
//   - func(View, MouseEvent),
//   - func(View),
//   - func(MouseEvent),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetMouseOutListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, MouseEvent](view, subviewID, MouseOut)
}
