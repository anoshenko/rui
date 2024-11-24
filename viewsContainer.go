package rui

import "strings"

// ParentView describe a view which can have a child views
type ParentView interface {
	// Views return a list of child views
	Views() []View
}

// ViewsContainer represent a mutable list-container of views
type ViewsContainer interface {
	View
	ParentView

	// Append appends a view to the end of the list of a view children
	Append(view View)

	// Insert inserts a view to the "index" position in the list of a view children
	Insert(view View, index int)

	// Remove removes a view from the list of a view children and return it
	RemoveView(index int) View

	// ViewIndex returns the index of view, -1 overwise
	ViewIndex(view View) int

	setContent(value any) bool
}

type viewsContainerData struct {
	viewData
	views []View
}

// Init initialize fields of ViewsContainer by default values
func (container *viewsContainerData) init(session Session) {
	container.viewData.init(session)
	container.tag = "ViewsContainer"
	container.views = []View{}
	container.get = container.getFunc
	container.set = container.setFunc
	container.remove = container.removeFunc
	container.changed = container.propertyChanged
}

func (container *viewsContainerData) setParentID(parentID string) {
	container.viewData.setParentID(parentID)
	htmlID := container.htmlID()
	for _, view := range container.views {
		view.setParentID(htmlID)
	}
}

// Views return a list of child views
func (container *viewsContainerData) Views() []View {
	if container.views == nil {
		container.views = []View{}
	} else if count := len(container.views); count > 0 {
		views := make([]View, count)
		copy(views, container.views)
		return views
	}
	return []View{}
}

// Append appends a view to the end of the list of a view children
func (container *viewsContainerData) Append(view View) {
	if view != nil {
		htmlID := container.htmlID()
		view.setParentID(htmlID)
		if len(container.views) == 0 {
			container.views = []View{view}
		} else {
			container.views = append(container.views, view)
		}

		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)

		viewHTML(view, buffer, "")
		container.Session().appendToInnerHTML(htmlID, buffer.String())

		if listener, ok := container.changeListener[Content]; ok {
			listener(container, Content)
		}
	}
}

// Insert inserts a view to the "index" position in the list of a view children
func (container *viewsContainerData) Insert(view View, index int) {
	if view != nil {
		if container.views == nil || index < 0 || index >= len(container.views) {
			container.Append(view)
			return
		}

		htmlID := container.htmlID()
		view.setParentID(htmlID)
		if index > 0 {
			container.views = append(container.views[:index], append([]View{view}, container.views[index:]...)...)
		} else {
			container.views = append([]View{view}, container.views...)
		}
		updateInnerHTML(htmlID, container.Session())
		if listener, ok := container.changeListener[Content]; ok {
			listener(container, Content)
		}
	}
}

// Remove removes view from list and return it
func (container *viewsContainerData) RemoveView(index int) View {
	if container.views == nil {
		container.views = []View{}
		return nil
	}

	count := len(container.views)
	if index < 0 || index >= count {
		return nil
	}

	view := container.views[index]
	if index == 0 {
		container.views = container.views[1:]
	} else if index == count-1 {
		container.views = container.views[:index]
	} else {
		container.views = append(container.views[:index], container.views[index+1:]...)
	}

	view.setParentID("")

	container.Session().callFunc("removeView", view.htmlID())
	if listener, ok := container.changeListener[Content]; ok {
		listener(container, Content)
	}
	return view
}

func (container *viewsContainerData) ViewIndex(view View) int {
	for index, v := range container.views {
		if v == view {
			return index
		}
	}
	return -1
}

func (container *viewsContainerData) cssStyle(self View, builder cssBuilder) {
	container.viewData.cssStyle(self, builder)
	builder.add(`overflow`, `auto`)
}

func (container *viewsContainerData) htmlSubviews(self View, buffer *strings.Builder) {
	if container.views != nil {
		for _, view := range container.views {
			viewHTML(view, buffer, "")
		}
	}
}

func viewFromTextValue(text string, session Session) View {
	if strings.Contains(text, "{") && strings.Contains(text, "}") {
		data := ParseDataText(text)
		if data != nil {
			if view := CreateViewFromObject(session, data); view != nil {
				return view
			}
		}
	}
	return NewTextView(session, Params{Text: text})
}

func (container *viewsContainerData) removeFunc(tag PropertyName) []PropertyName {
	switch tag {
	case Content:
		if len(container.views) > 0 {
			container.views = []View{}
			return []PropertyName{tag}
		}
		return []PropertyName{}

	case Disabled:
		if container.getRaw(Disabled) != nil {
			container.setRaw(Disabled, nil)
			for _, view := range container.views {
				view.Remove(Disabled)
			}
			return []PropertyName{tag}
		}
	}
	return container.viewData.removeFunc(tag)
}

func (container *viewsContainerData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Content:
		if container.setContent(value) {
			return []PropertyName{tag}
		}
		return nil

	case Disabled:
		oldDisabled := IsDisabled(container)
		result := container.viewData.setFunc(Disabled, value)
		if result != nil {
			disabled := IsDisabled(container)
			if oldDisabled != disabled {
				for _, view := range container.views {
					view.Set(Disabled, disabled)
				}
			}
		}
		return result
	}

	return container.viewData.setFunc(tag, value)
}

func (container *viewsContainerData) propertyChanged(tag PropertyName) {
	switch tag {
	case Content:
		updateInnerHTML(container.htmlID(), container.Session())

	default:
		container.viewData.propertyChanged(tag)
	}
}

func (container *viewsContainerData) setContent(value any) bool {
	session := container.Session()
	switch value := value.(type) {
	case View:
		container.views = []View{value}

	case []View:
		container.views = value

	case string:
		container.views = []View{viewFromTextValue(value, session)}

	case []string:
		views := []View{}
		for _, text := range value {
			views = append(views, viewFromTextValue(text, session))
		}
		container.views = views

	case []any:
		views := []View{}
		for _, v := range value {
			switch v := v.(type) {
			case View:
				views = append(views, v)

			case string:
				views = append(views, viewFromTextValue(v, session))

			default:
				notCompatibleType(Content, value)
				return false
			}
		}
		container.views = views

	case DataObject:
		if view := CreateViewFromObject(session, value); view != nil {
			container.views = []View{view}
		} else {
			return false
		}

	case []DataValue:
		views := []View{}
		for _, data := range value {
			if data.IsObject() {
				if view := CreateViewFromObject(session, data.Object()); view != nil {
					views = append(views, view)
				}
			} else {
				views = append(views, viewFromTextValue(data.Value(), session))
			}
		}
		container.views = views

	default:
		notCompatibleType(Content, value)
		return false
	}

	htmlID := container.htmlID()
	isDisabled := IsDisabled(container)
	for _, view := range container.views {
		view.setParentID(htmlID)
		if isDisabled {
			view.Set(Disabled, true)
		}
	}

	return true
}

func (container *viewsContainerData) getFunc(tag PropertyName) any {
	switch tag {
	case Content:
		return container.views

	default:
		return container.viewData.getFunc(tag)
	}
}

// AppendView appends a view to the end of the list of a view children
func AppendView(rootView View, containerID string, view View) bool {
	var container ViewsContainer = nil
	if containerID != "" {
		container = ViewsContainerByID(rootView, containerID)
	} else {
		if cont, ok := rootView.(ViewsContainer); ok {
			container = cont
		} else {
			ErrorLogF(`Unable to add a view to "%s"`, rootView.Tag())
		}
	}

	if container != nil {
		container.Append(view)
		return true
	}

	return false
}

// Insert inserts a view to the "index" position in the list of a view children
func InsertView(rootView View, containerID string, view View, index int) bool {
	var container ViewsContainer = nil
	if containerID != "" {
		container = ViewsContainerByID(rootView, containerID)
	} else {
		if cont, ok := rootView.(ViewsContainer); ok {
			container = cont
		} else {
			ErrorLogF(`Unable to add a view to "%s"`, rootView.Tag())
		}
	}

	if container != nil {
		container.Insert(view, index)
		return true
	}

	return false
}

// Remove removes a view from the list of a view children and return it
func RemoveView(rootView View, containerID string, index int) View {
	var container ViewsContainer = nil
	if containerID != "" {
		container = ViewsContainerByID(rootView, containerID)
	} else {
		if cont, ok := rootView.(ViewsContainer); ok {
			container = cont
		} else {
			ErrorLogF(`Unable to add a view to "%s"`, rootView.Tag())
		}
	}

	if container != nil {
		return container.RemoveView(index)
	}

	return nil
}
