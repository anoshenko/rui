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
	if value == nil {
		delete(view.properties, tag)
		return true
	}

	switch value := value.(type) {
	case func(View, Frame):
		view.properties[tag] = []func(View, Frame){value}

	case []func(View, Frame):
		if len(value) > 0 {
			view.properties[tag] = value
		} else {
			delete(view.properties, tag)
			return true
		}

	case func(Frame):
		fn := func(_ View, frame Frame) {
			value(frame)
		}
		view.properties[tag] = []func(View, Frame){fn}

	case []func(Frame):
		count := len(value)
		if count == 0 {
			delete(view.properties, tag)
			return true
		}

		listeners := make([]func(View, Frame), count)
		for i, val := range value {
			if val == nil {
				notCompatibleType(tag, val)
				return false
			}

			listeners[i] = func(_ View, frame Frame) {
				val(frame)
			}
		}
		view.properties[tag] = listeners

	case func(View):
		fn := func(view View, _ Frame) {
			value(view)
		}
		view.properties[tag] = []func(View, Frame){fn}

	case []func(View):
		count := len(value)
		if count == 0 {
			delete(view.properties, tag)
			return true
		}

		listeners := make([]func(View, Frame), count)
		for i, val := range value {
			if val == nil {
				notCompatibleType(tag, val)
				return false
			}

			listeners[i] = func(view View, _ Frame) {
				val(view)
			}
		}
		view.properties[tag] = listeners

	case func():
		fn := func(View, Frame) {
			value()
		}
		view.properties[tag] = []func(View, Frame){fn}

	case []func():
		count := len(value)
		if count == 0 {
			delete(view.properties, tag)
			return true
		}

		listeners := make([]func(View, Frame), count)
		for i, val := range value {
			if val == nil {
				notCompatibleType(tag, val)
				return false
			}

			listeners[i] = func(View, Frame) {
				val()
			}
		}
		view.properties[tag] = listeners

	case []any:
		count := len(value)
		if count == 0 {
			delete(view.properties, tag)
			return true
		}

		listeners := make([]func(View, Frame), count)
		for i, val := range value {
			if val == nil {
				notCompatibleType(tag, val)
				return false
			}

			switch val := val.(type) {
			case func(View, Frame):
				listeners[i] = val

			case func(Frame):
				listeners[i] = func(_ View, frame Frame) {
					val(frame)
				}

			case func(View):
				listeners[i] = func(view View, _ Frame) {
					val(view)
				}

			case func():
				listeners[i] = func(View, Frame) {
					val()
				}

			default:
				notCompatibleType(tag, val)
				return false
			}
		}
		view.properties[tag] = listeners

	default:
		notCompatibleType(tag, value)
		return false
	}

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
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(ResizeEvent); value != nil {
			if result, ok := value.([]func(View, Frame)); ok {
				return result
			}
		}
	}
	return []func(View, Frame){}
}
