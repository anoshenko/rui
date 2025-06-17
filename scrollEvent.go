package rui

// ScrollEvent is the constant for "scroll-event" property tag.
//
// Used by View.
// Is fired when the content of the view is scrolled.
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
const ScrollEvent PropertyName = "scroll-event"

func (view *viewData) onScroll(self View, x, y, width, height float64) {
	view.scroll.Left = x
	view.scroll.Top = y
	view.scroll.Width = width
	view.scroll.Height = height
	for _, listener := range getOneArgEventListeners[View, Frame](view, nil, ScrollEvent) {
		listener.Run(self, view.scroll)
	}
}

func (view *viewData) Scroll() Frame {
	return view.scroll
}

func (view *viewData) setScroll(x, y, width, height float64) {
	view.scroll.Left = x
	view.scroll.Top = y
	view.scroll.Width = width
	view.scroll.Height = height
}

// GetViewScroll returns ...
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetViewScroll(view View, subviewID ...string) Frame {
	view = getSubview(view, subviewID)
	if view == nil {
		return Frame{}
	}
	return view.Scroll()
}

// GetScrollListeners returns the list of "scroll-event" listeners. If there are no listeners then the empty list is returned
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
func GetScrollListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[View, Frame](view, subviewID, ScrollEvent)
}

// ScrollTo scrolls the view's content to the given position.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func ScrollViewTo(view View, subviewID string, x, y float64) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		view.Session().callFunc("scrollTo", view.htmlID(), x, y)
	}
}

// ScrollViewToEnd scrolls the view's content to the start of view.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func ScrollViewToStart(view View, subviewID ...string) {
	if view = getSubview(view, subviewID); view != nil {
		view.Session().callFunc("scrollToStart", view.htmlID())
	}
}

// ScrollViewToEnd scrolls the view's content to the end of view.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func ScrollViewToEnd(view View, subviewID ...string) {
	if view = getSubview(view, subviewID); view != nil {
		view.Session().callFunc("scrollToEnd", view.htmlID())
	}
}
