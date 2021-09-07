package rui

import "strings"

// AbsoluteLayout - list-container of View
type AbsoluteLayout interface {
	ViewsContainer
}

type absoluteLayoutData struct {
	viewsContainerData
}

// NewAbsoluteLayout create new AbsoluteLayout object and return it
func NewAbsoluteLayout(session Session, params Params) AbsoluteLayout {
	view := new(absoluteLayoutData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newAbsoluteLayout(session Session) View {
	return NewAbsoluteLayout(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (layout *absoluteLayoutData) Init(session Session) {
	layout.viewsContainerData.Init(session)
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
