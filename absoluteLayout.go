package rui

import "strings"

// AbsoluteLayout represents an interface for working with the AbsoluteLayout container where child views can be arbitrary.
//
// It extends the [ViewsContainer] interface and supports all its properties.
//
// Elements placed within a AbsoluteLayout container can have additional properties:
//   - [Left] ("left");
//   - [Top] ("top");
//   - [Right] ("right");
//   - [Bottom] ("bottom").
type AbsoluteLayout interface {
	ViewsContainer
}

type absoluteLayoutData struct {
	viewsContainerData
}

// NewAbsoluteLayout create a new [AbsoluteLayout] object and return its interface
func NewAbsoluteLayout(session Session, params Params) AbsoluteLayout {
	if session == nil {
		ErrorLog(`the Session argument cannot be nil`)
		return nil
	}

	view := new(absoluteLayoutData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newAbsoluteLayout(session Session) View {
	//return NewAbsoluteLayout(session, nil)
	return new(absoluteLayoutData)
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
			viewHTML(view, buffer, "")
		}
	}
}
