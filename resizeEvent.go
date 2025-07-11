package rui

// ResizeEvent is the constant for "resize-event" property tag.
//
// Used by View.
// Is fired when the view changes its size.
//
// General listener format:
//
//	func(view rui.View, frame rui.Frame)
//
// where:
//   - view - Interface of a view which generated this event,
//   - frame - New offset and size of the view's visible area.
//
// Allowed listener formats:
//
//	func(frame rui.Frame)
//	func(view rui.View)
//	func()
const ResizeEvent PropertyName = "resize-event"

func (view *viewData) onResize(self View, x, y, width, height float64) {
	view.frame.Left = x
	view.frame.Top = y
	view.frame.Width = width
	view.frame.Height = height
	for _, listener := range getOneArgEventListeners[View, Frame](view, nil, ResizeEvent) {
		listener.Run(self, view.frame)
	}
}

func (view *viewData) onItemResize(self View, index string, x, y, width, height float64) {
}

/*
func setFrameListener(properties Properties, tag PropertyName, value any) bool {
	if listeners, ok := valueToOneArgEventListeners[View, Frame](value); ok {
		if len(listeners) == 0 {
			properties.setRaw(tag, nil)
		} else {
			properties.setRaw(tag, listeners)
		}
		return true
	}
	notCompatibleType(tag, value)
	return false
}
*/

func (view *viewData) setNoResizeEvent() {
	view.noResizeEvent = true
}

func (view *viewData) isNoResizeEvent() bool {
	return view.noResizeEvent
}

func (container *viewsContainerData) isNoResizeEvent() bool {
	if container.noResizeEvent {
		return true
	}

	if parent := container.Parent(); parent != nil {
		return parent.isNoResizeEvent()
	}

	return false
}

func (view *viewData) Frame() Frame {
	return view.frame
}

// GetViewFrame returns the size and location of view's viewport.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetViewFrame(view View, subviewID ...string) Frame {
	view = getSubview(view, subviewID)
	if view == nil {
		return Frame{}
	}
	return view.Frame()
}

// GetResizeListeners returns the list of "resize-event" listeners. If there are no listeners then the empty list is returned
//
// Result elements can be of the following types:
//   - func(rui.View, rui.Frame),
//   - func(rui.View),
//   - func(rui.Frame),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetResizeListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, Frame](view, subviewID, ResizeEvent)
}
