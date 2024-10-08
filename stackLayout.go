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
	peek              int
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
	return NewStackLayout(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (layout *stackLayoutData) init(session Session) {
	layout.viewsContainerData.init(session)
	layout.tag = "StackLayout"
	layout.systemClass = "ruiStackLayout"
	layout.properties[TransitionEndEvent] = []func(View, string){layout.pushFinished, layout.popFinished}
}

func (layout *stackLayoutData) String() string {
	return getViewString(layout, nil)
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
			layout.propertyChangedEvent(Current)
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

func (layout *stackLayoutData) Set(tag string, value any) bool {
	return layout.set(strings.ToLower(tag), value)
}

func (layout *stackLayoutData) set(tag string, value any) bool {
	if value == nil {
		layout.remove(tag)
		return true
	}

	switch tag {
	case TransitionEndEvent:
		listeners, ok := valueToEventListeners[View, string](value)
		if ok && listeners != nil {
			listeners = append(listeners, layout.pushFinished)
			listeners = append(listeners, layout.popFinished)
			layout.properties[TransitionEndEvent] = listeners
			layout.propertyChangedEvent(TransitionEndEvent)
		}
		return ok

	case Current:
		setCurrent := func(index int) {
			if index != layout.peek {
				if layout.peek < len(layout.views) {
					layout.Session().updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(layout.peek), "visibility", "hidden")
				}

				layout.peek = index
				layout.Session().updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(index), "visibility", "visible")
				layout.propertyChangedEvent(Current)
			}
		}
		switch value := value.(type) {
		case string:
			text, ok := layout.session.resolveConstants(value)
			if !ok {
				invalidPropertyValue(tag, value)
				return false
			}
			n, err := strconv.Atoi(strings.Trim(text, " \t"))
			if err != nil {
				invalidPropertyValue(tag, value)
				ErrorLog(err.Error())
				return false
			}
			setCurrent(n)

		default:
			n, ok := isInt(value)
			if !ok {
				notCompatibleType(tag, value)
				return false
			} else if n < 0 || n >= len(layout.views) {
				ErrorLogF(`The view index "%d" of "%s" property is out of range`, n, tag)
				return false
			}
			setCurrent(n)
		}
		return true
	}
	return layout.viewsContainerData.set(tag, value)
}

func (layout *stackLayoutData) Remove(tag string) {
	layout.remove(strings.ToLower(tag))
}

func (layout *stackLayoutData) remove(tag string) {
	switch tag {
	case TransitionEndEvent:
		layout.properties[TransitionEndEvent] = []func(View, string){layout.pushFinished, layout.popFinished}
		layout.propertyChangedEvent(TransitionEndEvent)

	case Current:
		layout.set(Current, 0)

	default:
		layout.viewsContainerData.remove(tag)
	}
}

func (layout *stackLayoutData) Get(tag string) any {
	return layout.get(strings.ToLower(tag))
}

func (layout *stackLayoutData) get(tag string) any {
	if tag == Current {
		return layout.peek
	}
	return layout.viewsContainerData.get(tag)
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
				layout.propertyChangedEvent(Current)
			}
			return true
		}
	}

	ErrorLog(`MoveToFront() fail. Subview not found."`)
	return false
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
				layout.propertyChangedEvent(Current)
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
		layout.propertyChangedEvent(Current)
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
		layout.propertyChangedEvent(Current)
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
	defer layout.propertyChangedEvent(Current)
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
	layout.propertyChangedEvent(Content)
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
