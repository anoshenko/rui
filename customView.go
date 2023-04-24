package rui

import "strings"

// CustomView defines a custom view interface
type CustomView interface {
	ViewsContainer
	CreateSuperView(session Session) View
	SuperView() View
	setSuperView(view View)
	setTag(tag string)
}

// CustomViewData defines a data of a basic custom view
type CustomViewData struct {
	tag       string
	superView View
}

// InitCustomView initializes fields of CustomView by default values
func InitCustomView(customView CustomView, tag string, session Session, params Params) bool {
	customView.setTag(tag)
	if view := customView.CreateSuperView(session); view != nil {
		customView.setSuperView(view)
		setInitParams(customView, params)
	} else {
		ErrorLog(`nil SuperView of "` + tag + `" view`)
		return false
	}
	return true
}

// SuperView returns a super view
func (customView *CustomViewData) SuperView() View {
	return customView.superView
}

func (customView *CustomViewData) setSuperView(view View) {
	customView.superView = view
}

func (customView *CustomViewData) setTag(tag string) {
	customView.tag = tag
}

// Get returns a value of the property with name defined by the argument.
// The type of return value depends on the property. If the property is not set then nil is returned.
func (customView *CustomViewData) Get(tag string) any {
	return customView.superView.Get(tag)
}

func (customView *CustomViewData) getRaw(tag string) any {
	return customView.superView.getRaw(tag)
}

func (customView *CustomViewData) setRaw(tag string, value any) {
	customView.superView.setRaw(tag, value)
}

// Set sets the value (second argument) of the property with name defined by the first argument.
// Return "true" if the value has been set, in the opposite case "false" are returned and
// a description of the error is written to the log
func (customView *CustomViewData) Set(tag string, value any) bool {
	return customView.superView.Set(tag, value)
}

func (customView *CustomViewData) SetAnimated(tag string, value any, animation Animation) bool {
	return customView.superView.SetAnimated(tag, value, animation)
}

func (customView *CustomViewData) SetChangeListener(tag string, listener func(View, string)) {
	customView.superView.SetChangeListener(tag, listener)
}

// Remove removes the property with name defined by the argument
func (customView *CustomViewData) Remove(tag string) {
	customView.superView.Remove(tag)
}

// AllTags returns an array of the set properties
func (customView *CustomViewData) AllTags() []string {
	return customView.superView.AllTags()
}

// Clear removes all properties
func (customView *CustomViewData) Clear() {
	customView.superView.Clear()
}

func (customView *CustomViewData) cssViewStyle(buffer cssBuilder, session Session) {
	customView.superView.cssViewStyle(buffer, session)
}

// Session returns a current Session interface
func (customView *CustomViewData) Session() Session {
	return customView.superView.Session()
}

// Parent returns a parent view
func (customView *CustomViewData) Parent() View {
	return customView.superView.Parent()
}

func (customView *CustomViewData) parentHTMLID() string {
	return customView.superView.parentHTMLID()
}

func (customView *CustomViewData) setParentID(parentID string) {
	customView.superView.setParentID(parentID)
}

// Tag returns a tag of View interface
func (customView *CustomViewData) Tag() string {
	if customView.tag != "" {
		return customView.tag
	}
	return customView.superView.Tag()
}

// ID returns a id of the view
func (customView *CustomViewData) ID() string {
	return customView.superView.ID()
}

// Focusable returns true if the view receives the focus
func (customView *CustomViewData) Focusable() bool {
	return customView.superView.Focusable()
}

/*
// SetTransitionEndListener sets the new listener of the transition end event
func (customView *CustomViewData) SetTransitionEndListener(property string, listener TransitionEndListener) {
	customView.superView.SetTransitionEndListener(property, listener)
}

// SetTransitionEndFunc sets the new listener function of the transition end event
func (customView *CustomViewData) SetTransitionEndFunc(property string, listenerFunc func(View, string)) {
	customView.superView.SetTransitionEndFunc(property, listenerFunc)
}
*/

// Frame returns a location and size of the view in pixels
func (customView *CustomViewData) Frame() Frame {
	return customView.superView.Frame()
}

func (customView *CustomViewData) Scroll() Frame {
	return customView.superView.Scroll()
}

func (customView *CustomViewData) HasFocus() bool {
	return customView.superView.HasFocus()
}

func (customView *CustomViewData) onResize(self View, x, y, width, height float64) {
	customView.superView.onResize(customView.superView, x, y, width, height)
}

func (customView *CustomViewData) onItemResize(self View, index string, x, y, width, height float64) {
	customView.superView.onItemResize(customView.superView, index, x, y, width, height)
}

func (customView *CustomViewData) handleCommand(self View, command string, data DataObject) bool {
	return customView.superView.handleCommand(customView.superView, command, data)
}

func (customView *CustomViewData) htmlClass(disabled bool) string {
	return customView.superView.htmlClass(disabled)
}

func (customView *CustomViewData) htmlTag() string {
	return customView.superView.htmlTag()
}

func (customView *CustomViewData) closeHTMLTag() bool {
	return customView.superView.closeHTMLTag()
}

func (customView *CustomViewData) htmlID() string {
	return customView.superView.htmlID()
}

func (customView *CustomViewData) htmlSubviews(self View, buffer *strings.Builder) {
	customView.superView.htmlSubviews(customView.superView, buffer)
}

func (customView *CustomViewData) htmlProperties(self View, buffer *strings.Builder) {
	customView.superView.htmlProperties(customView.superView, buffer)
}

func (customView *CustomViewData) htmlDisabledProperties(self View, buffer *strings.Builder) {
	customView.superView.htmlDisabledProperties(customView.superView, buffer)
}

func (customView *CustomViewData) cssStyle(self View, builder cssBuilder) {
	customView.superView.cssStyle(customView.superView, builder)
}

func (customView *CustomViewData) addToCSSStyle(addCSS map[string]string) {
	customView.superView.addToCSSStyle(addCSS)
}

func (customView *CustomViewData) setNoResizeEvent() {
	customView.superView.setNoResizeEvent()
}

func (customView *CustomViewData) isNoResizeEvent() bool {
	return customView.superView.isNoResizeEvent()
}

// Views return a list of child views
func (customView *CustomViewData) Views() []View {
	if customView.superView != nil {
		if container, ok := customView.superView.(ViewsContainer); ok {
			return container.Views()
		}
	}
	return []View{}
}

// Append appends a view to the end of the list of a view children
func (customView *CustomViewData) Append(view View) {
	if customView.superView != nil {
		if container, ok := customView.superView.(ViewsContainer); ok {
			container.Append(view)
		}
	}
}

// Insert inserts a view to the "index" position in the list of a view children
func (customView *CustomViewData) Insert(view View, index int) {
	if customView.superView != nil {
		if container, ok := customView.superView.(ViewsContainer); ok {
			container.Insert(view, index)
		}
	}
}

// Remove removes a view from the list of a view children and return it
func (customView *CustomViewData) RemoveView(index int) View {
	if customView.superView != nil {
		if container, ok := customView.superView.(ViewsContainer); ok {
			return container.RemoveView(index)
		}
	}
	return nil
}

// Remove removes a view from the list of a view children and return it
func (customView *CustomViewData) ViewIndex(view View) int {
	if customView.superView != nil {
		if container, ok := customView.superView.(ViewsContainer); ok {
			return container.ViewIndex(view)
		}
	}
	return -1
}

func (customView *CustomViewData) String() string {
	if customView.superView != nil {
		return getViewString(customView)
	}
	return customView.tag + " { }"
}

func (customView *CustomViewData) setScroll(x, y, width, height float64) {
	if customView.superView != nil {
		customView.superView.setScroll(x, y, width, height)
	}
}

func (customView *CustomViewData) Transition(tag string) Animation {
	if customView.superView != nil {
		return customView.superView.Transition(tag)
	}
	return nil
}

func (customView *CustomViewData) Transitions() map[string]Animation {
	if customView.superView != nil {
		return customView.superView.Transitions()
	}
	return map[string]Animation{}
}

func (customView *CustomViewData) SetTransition(tag string, animation Animation) {
	if customView.superView != nil {
		customView.superView.SetTransition(tag, animation)
	}
}
