package rui

// ResizeEvent is the constant for "resize-event" property tag.
// The "resize-event" is fired when the view changes its size.
// The main listener format:
//   func(View, Frame).
// The additional listener formats:
//   func(Frame), func(View), and func().
const ResizeEvent = "resize-event"

func (view *viewData) onResize(self View, x, y, width, height float64) {
	view.frame.Left = x
	view.frame.Top = y
	view.frame.Width = width
	view.frame.Height = height
	for _, listener := range GetResizeListeners(view, "") {
		listener(self, view.frame)
	}
}

func (view *viewData) onItemResize(self View, index string, x, y, width, height float64) {
}

func (view *viewData) setFrameListener(tag string, value any) bool {
	listeners, ok := valueToEventListeners[View, Frame](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		delete(view.properties, tag)
	} else {
		view.properties[tag] = listeners
	}
	view.propertyChangedEvent(tag)
	return true
}

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
// If the second argument (subviewID) is "" then the value of the first argument (view) is returned
func GetViewFrame(view View, subviewID string) Frame {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return Frame{}
	}
	return view.Frame()
}

// GetResizeListeners returns the list of "resize-event" listeners. If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then the listeners list of the first argument (view) is returned
func GetResizeListeners(view View, subviewID string) []func(View, Frame) {
	return getEventListeners[View, Frame](view, subviewID, ResizeEvent)
}
