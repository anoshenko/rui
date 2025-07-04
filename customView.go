package rui

import (
	"iter"
	"strings"
)

// CustomView defines a custom view interface
type CustomView interface {
	ViewsContainer

	// CreateSuperView must be implemented to create a base view from which custom control has been built
	CreateSuperView(session Session) View

	// SuperView must be implemented to return a base view from which custom control has been built
	SuperView() View

	setSuperView(view View)
	setTag(tag string)
}

// CustomViewData defines a data of a basic custom view
type CustomViewData struct {
	tag           string
	superView     View
	defaultParams Params
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

func (customView *CustomViewData) init(session Session) {
}

// SuperView returns a super view
func (customView *CustomViewData) SuperView() View {
	return customView.superView
}

func (customView *CustomViewData) setSuperView(view View) {
	customView.superView = view
	customView.defaultParams = Params{}
	for tag, value := range view.All() {
		if value != nil {
			customView.defaultParams[tag] = value
		}
	}
}

func (customView *CustomViewData) setTag(tag string) {
	customView.tag = tag
}

// Get returns a value of the property with name defined by the argument.
// The type of return value depends on the property. If the property is not set then nil is returned.
func (customView *CustomViewData) Get(tag PropertyName) any {
	return customView.superView.Get(tag)
}

func (customView *CustomViewData) getRaw(tag PropertyName) any {
	return customView.superView.getRaw(tag)
}

func (customView *CustomViewData) setRaw(tag PropertyName, value any) {
	customView.superView.setRaw(tag, value)
}

func (customView *CustomViewData) setContent(value any) bool {
	if container, ok := customView.superView.(ViewsContainer); ok {
		return container.setContent(value)
	}
	return false
}

// Set sets the value (second argument) of the property with name defined by the first argument.
// Return "true" if the value has been set, in the opposite case "false" are returned and
// a description of the error is written to the log
func (customView *CustomViewData) Set(tag PropertyName, value any) bool {
	return customView.superView.Set(tag, value)
}

// SetAnimated sets the value (second argument) of the property with name defined by the first argument.
// Return "true" if the value has been set, in the opposite case "false" are returned and
// a description of the error is written to the log
func (customView *CustomViewData) SetAnimated(tag PropertyName, value any, animation AnimationProperty) bool {
	return customView.superView.SetAnimated(tag, value, animation)
}

func (customView *CustomViewData) SetParams(params Params) bool {
	return customView.superView.SetParams(params)
}

// SetChangeListener set the function to track the change of the View property
func (customView *CustomViewData) SetChangeListener(tag PropertyName, listener any) bool {
	return customView.superView.SetChangeListener(tag, listener)
}

// Remove removes the property with name defined by the argument
func (customView *CustomViewData) Remove(tag PropertyName) {
	customView.superView.Remove(tag)
}

func (customView *CustomViewData) AllTags() []PropertyName {
	return customView.superView.AllTags()
}

// AllTags returns an array of the set properties
func (customView *CustomViewData) All() iter.Seq2[PropertyName, any] {
	return customView.superView.All()
}

func (customView *CustomViewData) IsEmpty() bool {
	return customView.superView.IsEmpty()
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

// Scroll returns a location and size of a scrollable view in pixels
func (customView *CustomViewData) Scroll() Frame {
	return customView.superView.Scroll()
}

// HasFocus returns "true" if the view has focus
func (customView *CustomViewData) HasFocus() bool {
	return customView.superView.HasFocus()
}

func (customView *CustomViewData) onResize(self View, x, y, width, height float64) {
	customView.superView.onResize(customView.superView, x, y, width, height)
}

func (customView *CustomViewData) onItemResize(self View, index string, x, y, width, height float64) {
	customView.superView.onItemResize(customView.superView, index, x, y, width, height)
}

func (customView *CustomViewData) handleCommand(self View, command PropertyName, data DataObject) bool {
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

func (customView *CustomViewData) htmlDisabledProperty() bool {
	return customView.superView.htmlDisabledProperty()
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

func (customView *CustomViewData) RemoveViewByID(id string) View {
	if customView.superView != nil {
		if container, ok := customView.superView.(ViewsContainer); ok {
			return container.RemoveViewByID(id)
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

func (customView *CustomViewData) excludeTags() []PropertyName {
	if customView.superView != nil {
		exclude := []PropertyName{}
		for tag, value := range customView.defaultParams {
			if value == customView.superView.getRaw(tag) {
				exclude = append(exclude, tag)
			}
		}
		return exclude
	}
	return nil
}

// String convert internal representation of a [CustomViewData] into a string.
func (customView *CustomViewData) String() string {
	if customView.superView != nil {
		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)
		writeViewStyle(customView.tag, customView, buffer, "", customView.excludeTags())
		return buffer.String()
	}
	return customView.tag + " { }"
}

func (customView *CustomViewData) setScroll(x, y, width, height float64) {
	if customView.superView != nil {
		customView.superView.setScroll(x, y, width, height)
	}
}

// Transition returns the transition animation of the property(tag). Returns nil is there is no transition animation.
func (customView *CustomViewData) Transition(tag PropertyName) AnimationProperty {
	if customView.superView != nil {
		return customView.superView.Transition(tag)
	}
	return nil
}

// Transitions returns a map of transition animations. The result is always non-nil.
func (customView *CustomViewData) Transitions() map[PropertyName]AnimationProperty {
	if customView.superView != nil {
		return customView.superView.Transitions()
	}
	return map[PropertyName]AnimationProperty{}
}

// SetTransition sets the transition animation for the property if "animation" argument is not nil, and
// removes the transition animation of the property if "animation" argument  is nil.
// The "tag" argument is the property name.
func (customView *CustomViewData) SetTransition(tag PropertyName, animation AnimationProperty) {
	if customView.superView != nil {
		customView.superView.SetTransition(tag, animation)
	}
}

func (customView *CustomViewData) LoadFile(file FileInfo, result func(FileInfo, []byte)) {
	if customView.superView != nil {
		customView.superView.LoadFile(file, result)
	}
}

func (customView *CustomViewData) binding() any {
	if customView.superView != nil {
		return customView.superView.binding()
	}
	return nil
}
