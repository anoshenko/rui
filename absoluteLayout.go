package rui

import "strings"

// AbsoluteLayout represent an AbsoluteLayout view where child views can be arbitrary positioned
type AbsoluteLayout interface {
	ViewsContainer
}

type absoluteLayoutData struct {
	viewsContainerData
}

// NewAbsoluteLayout create new AbsoluteLayout object and return it
func NewAbsoluteLayout(session Session, params Params) AbsoluteLayout {
	view := new(absoluteLayoutData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newAbsoluteLayout(session Session) View {
	return NewAbsoluteLayout(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (layout *absoluteLayoutData) init(session Session) {
	layout.viewsContainerData.init(session)
	layout.tag = "AbsoluteLayout"
	layout.systemClass = "ruiAbsoluteLayout"
}

func (layout *absoluteLayoutData) htmlSubviews(self View, buffer *strings.Builder) {
	if layout.views != nil {
		for _, view := range layout.views {
			view.addToCSSStyle(map[string]string{`position`: `absolute`})
			viewHTML(view, buffer)
		}
	}
}
