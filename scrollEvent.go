package rui

import "fmt"

// ScrollEvent is the constant for "scroll-event" property tag.
// The "scroll-event" is fired when the content of the view is scrolled.
// The main listener format:
//   func(View, Frame).
// The additional listener formats:
//   func(Frame), func(View), and func().
const ScrollEvent = "scroll-event"

func (view *viewData) onScroll(self View, x, y, width, height float64) {
	view.scroll.Left = x
	view.scroll.Top = y
	view.scroll.Width = width
	view.scroll.Height = height
	for _, listener := range GetScrollListeners(view, "") {
		listener(self, view.scroll)
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
// If the second argument (subviewID) is "" then a value of the first argument (view) is returned
func GetViewScroll(view View, subviewID string) Frame {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return Frame{}
	}
	return view.Scroll()
}

// GetScrollListeners returns the list of "scroll-event" listeners. If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then the listeners list of the first argument (view) is returned
func GetScrollListeners(view View, subviewID string) []func(View, Frame) {
	return getEventListeners[View, Frame](view, subviewID, ResizeEvent)
}

// ScrollTo scrolls the view's content to the given position.
// If the second argument (subviewID) is "" then the first argument (view) is used
func ScrollViewTo(view View, subviewID string, x, y float64) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		view.Session().runScript(fmt.Sprintf(`scrollTo("%s", %g, %g)`, view.htmlID(), x, y))
	}
}

// ScrollViewToEnd scrolls the view's content to the start of view.
// If the second argument (subviewID) is "" then the first argument (view) is used
func ScrollViewToStart(view View, subviewID string) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		view.Session().runScript(`scrollToStart("` + view.htmlID() + `")`)
	}
}

// ScrollViewToEnd scrolls the view's content to the end of view.
// If the second argument (subviewID) is "" then the first argument (view) is used
func ScrollViewToEnd(view View, subviewID string) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		view.Session().runScript(`scrollToEnd("` + view.htmlID() + `")`)
	}
}
