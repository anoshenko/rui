package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// Constants which represent [StackLayout] animation type during pushing or popping views
const (
	// DefaultAnimation - default animation of StackLayout push
	DefaultAnimation = 0
	// StartToEndAnimation - start to end animation of StackLayout push
	StartToEndAnimation = 1
	// EndToStartAnimation - end to start animation of StackLayout push
	EndToStartAnimation = 2
	// TopDownAnimation - top down animation of StackLayout push
	TopDownAnimation = 3
	// BottomUpAnimation - bottom up animation of StackLayout push
	BottomUpAnimation = 4
)

// StackLayout represents a StackLayout view
type StackLayout interface {
	ViewsContainer

	// Peek returns the current (visible) View. If StackLayout is empty then it returns nil.
	Peek() View

	// RemovePeek removes the current View and returns it. If StackLayout is empty then it doesn't do anything and returns nil.
	RemovePeek() View

	// MoveToFront makes the given View current. Returns true if successful, false otherwise.
	MoveToFront(view View) bool

	// MoveToFrontByID makes the View current by viewID. Returns true if successful, false otherwise.
	MoveToFrontByID(viewID string) bool

	// Push adds a new View to the container and makes it current.
	// It is similar to Append, but the addition is done using an animation effect.
	// The animation type is specified by the second argument and can take the following values:
	// * DefaultAnimation (0) - Default animation. For the Push function it is EndToStartAnimation, for Pop - StartToEndAnimation;
	// * StartToEndAnimation (1) - Animation from beginning to end. The beginning and the end are determined by the direction of the text output;
	// * EndToStartAnimation (2) - End-to-Beginning animation;
	// * TopDownAnimation (3) - Top-down animation;
	// * BottomUpAnimation (4) - Bottom up animation.
	// The third argument `onPushFinished` is the function to be called when the animation ends. It may be nil.
	Push(view View, animation int, onPushFinished func())

	// Pop removes the current View from the container using animation.
	// The second argument `onPopFinished`` is the function to be called when the animation ends. It may be nil.
	// The function will return false if the StackLayout is empty and true if the current item has been removed.
	Pop(animation int, onPopFinished func(View)) bool
}

type stackLayoutData struct {
	viewsContainerData
	peek, prevPeek    int
	pushView, popView View
	animationType     int
	onPushFinished    func()
	onPopFinished     func(View)
}

// NewStackLayout create new StackLayout object and return it
func NewStackLayout(session Session, params Params) StackLayout {
	view := new(stackLayoutData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newStackLayout(session Session) View {
	//return NewStackLayout(session, nil)
	return new(stackLayoutData)
}

// Init initialize fields of ViewsContainer by default values
func (layout *stackLayoutData) init(session Session) {
	layout.viewsContainerData.init(session)
	layout.tag = "StackLayout"
	layout.systemClass = "ruiStackLayout"
	layout.properties[TransitionEndEvent] = []func(View, string){layout.pushFinished, layout.popFinished}
	layout.getFunc = layout.get
	layout.set = layout.setFunc
	layout.remove = layout.removeFunc
}

func (layout *stackLayoutData) pushFinished(view View, tag string) {
	if tag == "ruiPush" {
		if layout.pushView != nil {
			layout.pushView = nil
			count := len(layout.views)
			if count > 0 {
				layout.peek = count - 1
			} else {
				layout.peek = 0
			}
			updateInnerHTML(layout.htmlID(), layout.session)
			layout.currentChanged()
		}

		if layout.onPushFinished != nil {
			onPushFinished := layout.onPushFinished
			layout.onPushFinished = nil
			onPushFinished()
		}
	}
}

func (layout *stackLayoutData) popFinished(view View, tag string) {
	if tag == "ruiPop" {
		popView := layout.popView
		layout.popView = nil
		updateInnerHTML(layout.htmlID(), layout.session)
		if layout.onPopFinished != nil {
			onPopFinished := layout.onPopFinished
			layout.onPopFinished = nil
			onPopFinished(popView)
		}
	}
}

func (layout *stackLayoutData) setFunc(view View, tag PropertyName, value any) []PropertyName {
	switch tag {
	case TransitionEndEvent:
		listeners, ok := valueToEventListeners[View, string](value)
		if ok && listeners != nil {
			listeners = append(listeners, layout.pushFinished)
			listeners = append(listeners, layout.popFinished)
			view.setRaw(TransitionEndEvent, listeners)
			return []PropertyName{tag}
		}
		return nil

	case Current:
		newCurrent := 0
		switch value := value.(type) {
		case string:
			text, ok := layout.session.resolveConstants(value)
			if !ok {
				invalidPropertyValue(tag, value)
				return nil
			}
			n, err := strconv.Atoi(strings.Trim(text, " \t"))
			if err != nil {
				invalidPropertyValue(tag, value)
				ErrorLog(err.Error())
				return nil
			}
			newCurrent = n

		default:
			n, ok := isInt(value)
			if !ok {
				notCompatibleType(tag, value)
				return nil
			} else if n < 0 || n >= len(layout.views) {
				ErrorLogF(`The view index "%d" of "%s" property is out of range`, n, tag)
				return nil
			}
			newCurrent = n
		}

		layout.prevPeek = layout.peek
		if newCurrent == layout.peek {
			return []PropertyName{}
		}

		layout.peek = newCurrent
		return []PropertyName{tag}
	}
	return layout.viewsContainerData.setFunc(view, tag, value)
}

func (layout *stackLayoutData) propertyChanged(view View, tag PropertyName) {
	switch tag {
	case Current:
		if layout.prevPeek != layout.peek {
			if layout.prevPeek < len(layout.views) {
				layout.Session().updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(layout.prevPeek), "visibility", "hidden")
			}
			layout.Session().updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(layout.prevPeek), "visibility", "visible")
			layout.prevPeek = layout.peek
		}
	default:
		viewsContainerPropertyChanged(view, tag)
	}
}

func (layout *stackLayoutData) removeFunc(view View, tag PropertyName) []PropertyName {
	switch tag {
	case TransitionEndEvent:
		view.setRaw(TransitionEndEvent, []func(View, string){layout.pushFinished, layout.popFinished})
		return []PropertyName{tag}

	case Current:
		view.setRaw(Current, 0)
		return []PropertyName{tag}
	}
	return layout.viewsContainerData.removeFunc(view, tag)
}

func (layout *stackLayoutData) get(view View, tag PropertyName) any {
	if tag == Current {
		return layout.peek
	}
	return layout.viewsContainerData.get(view, tag)
}

func (layout *stackLayoutData) Peek() View {
	if int(layout.peek) < len(layout.views) {
		return layout.views[layout.peek]
	}
	return nil
}

func (layout *stackLayoutData) MoveToFront(view View) bool {
	peek := int(layout.peek)
	htmlID := view.htmlID()
	for i, view2 := range layout.views {
		if view2.htmlID() == htmlID {
			if i != peek {
				if peek < len(layout.views) {
					layout.Session().updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(peek), "visibility", "hidden")
				}

				layout.peek = i
				layout.Session().updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(i), "visibility", "visible")
				layout.currentChanged()
			}
			return true
		}
	}

	ErrorLog(`MoveToFront() fail. Subview not found."`)
	return false
}

func (layout *stackLayoutData) currentChanged() {
	if listener, ok := layout.changeListener[Current]; ok {
		listener(layout, Current)
	}
}

func (layout *stackLayoutData) MoveToFrontByID(viewID string) bool {
	peek := int(layout.peek)
	for i, view := range layout.views {
		if view.ID() == viewID {
			if i != peek {
				if peek < len(layout.views) {
					layout.Session().updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(peek), "visibility", "hidden")
				}

				layout.peek = i
				layout.Session().updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(i), "visibility", "visible")
				layout.currentChanged()
			}
			return true
		}
	}

	ErrorLogF(`MoveToFront("%s") fail. Subview with "%s" not found."`, viewID, viewID)
	return false
}

func (layout *stackLayoutData) Append(view View) {
	if view != nil {
		layout.peek = len(layout.views)
		layout.viewsContainerData.Append(view)
		layout.currentChanged()
	} else {
		ErrorLog("StackLayout.Append(nil, ....) is forbidden")
	}
}

func (layout *stackLayoutData) Insert(view View, index int) {
	if view != nil {
		count := len(layout.views)
		if index < count {
			layout.peek = int(index)
		} else {
			layout.peek = count
		}
		layout.viewsContainerData.Insert(view, index)
		layout.currentChanged()
	} else {
		ErrorLog("StackLayout.Insert(nil, ....) is forbidden")
	}
}

func (layout *stackLayoutData) RemoveView(index int) View {
	if index < 0 || index >= len(layout.views) {
		return nil
	}

	if layout.peek > 0 {
		layout.peek--
	}
	defer layout.currentChanged()
	return layout.viewsContainerData.RemoveView(index)
}

func (layout *stackLayoutData) RemovePeek() View {
	return layout.RemoveView(len(layout.views) - 1)
}

func (layout *stackLayoutData) Push(view View, animation int, onPushFinished func()) {
	if view == nil {
		ErrorLog("StackLayout.Push(nil, ....) is forbidden")
		return
	}

	layout.pushView = view
	layout.animationType = animation
	//layout.animation["ruiPush"] = Animation{FinishListener: layout}
	layout.onPushFinished = onPushFinished

	htmlID := layout.htmlID()
	session := layout.Session()

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`<div id="`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`push" class="ruiStackPageLayout" ontransitionend="stackTransitionEndEvent('`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`', 'ruiPush', event)" style="`)

	switch layout.animationType {
	case StartToEndAnimation:
		buffer.WriteString(fmt.Sprintf("transform: translate(-%gpx, 0px); transition: transform ", layout.frame.Width))

	case TopDownAnimation:
		buffer.WriteString(fmt.Sprintf("transform: translate(0px, -%gpx); transition: transform ", layout.frame.Height))

	case BottomUpAnimation:
		buffer.WriteString(fmt.Sprintf("transform: translate(0px, %gpx); transition: transform ", layout.frame.Height))

	default:
		buffer.WriteString(fmt.Sprintf("transform: translate(%gpx, 0px); transition: transform ", layout.frame.Width))
	}

	buffer.WriteString(`1s ease;">`)

	viewHTML(layout.pushView, buffer)
	buffer.WriteString(`</div>`)

	session.appendToInnerHTML(htmlID, buffer.String())
	layout.session.updateCSSProperty(htmlID+"push", "transform", "translate(0px, 0px)")

	layout.views = append(layout.views, view)
	view.setParentID(htmlID)

	if listener, ok := layout.changeListener[Content]; ok {
		listener(layout, Content)
	}
}

func (layout *stackLayoutData) Pop(animation int, onPopFinished func(View)) bool {
	count := len(layout.views)
	if count == 0 || layout.peek >= count {
		ErrorLog("StackLayout is empty")
		return false
	}

	layout.popView = layout.views[layout.peek]
	layout.RemoveView(layout.peek)

	layout.animationType = animation
	//layout.animation["ruiPop"] = Animation{FinishListener: layout}
	layout.onPopFinished = onPopFinished

	htmlID := layout.htmlID()
	session := layout.Session()

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`<div id="`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`pop" class="ruiStackPageLayout" ontransitionend="stackTransitionEndEvent('`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`', 'ruiPop', event)" ontransitioncancel="stackTransitionEndEvent('`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`', 'ruiPop', event)" style="transition: transform 1s ease;">`)
	viewHTML(layout.popView, buffer)
	buffer.WriteString(`</div>`)

	session.appendToInnerHTML(htmlID, buffer.String())

	var value string
	switch layout.animationType {
	case TopDownAnimation:
		value = fmt.Sprintf("translate(0px, -%gpx)", layout.frame.Height)

	case BottomUpAnimation:
		value = fmt.Sprintf("translate(0px, %gpx)", layout.frame.Height)

	case StartToEndAnimation:
		value = fmt.Sprintf("translate(-%gpx, 0px)", layout.frame.Width)

	default:
		value = fmt.Sprintf("translate(%gpx, 0px)", layout.frame.Width)
	}

	layout.session.updateCSSProperty(htmlID+"pop", "transform", value)
	return true
}

func (layout *stackLayoutData) htmlSubviews(self View, buffer *strings.Builder) {
	count := len(layout.views)
	if count > 0 {
		htmlID := layout.htmlID()
		peek := int(layout.peek)
		if peek >= count {
			peek = count - 1
		}

		for i, view := range layout.views {
			buffer.WriteString(`<div id="`)
			buffer.WriteString(htmlID)
			buffer.WriteString(`page`)
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString(`" class="ruiStackPageLayout"`)
			if i != peek {
				buffer.WriteString(` style="visibility: hidden;"`)
			}
			buffer.WriteString(`>`)
			viewHTML(view, buffer)
			buffer.WriteString(`</div>`)
		}
	}
}
