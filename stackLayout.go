package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO PeekChangedEvent

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

// StackLayout - list-container of View
type StackLayout interface {
	ViewsContainer
	Peek() View
	MoveToFront(view View) bool
	MoveToFrontByID(viewID string) bool
	Push(view View, animation int, onPushFinished func())
	Pop(animation int, onPopFinished func(View)) bool
}

type stackLayoutData struct {
	viewsContainerData
	peek              uint
	pushView, popView View
	animationType     int
	onPushFinished    func()
	onPopFinished     func(View)
}

// NewStackLayout create new StackLayout object and return it
func NewStackLayout(session Session, params Params) StackLayout {
	view := new(stackLayoutData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newStackLayout(session Session) View {
	return NewStackLayout(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (layout *stackLayoutData) Init(session Session) {
	layout.viewsContainerData.Init(session)
	layout.tag = "StackLayout"
	layout.systemClass = "ruiStackLayout"
	layout.properties[TransitionEndEvent] = []func(View, string){layout.pushFinished, layout.popFinished}
}

func (layout *stackLayoutData) pushFinished(view View, tag string) {
	if tag == "ruiPush" {
		if layout.pushView != nil {
			layout.pushView = nil
			count := len(layout.views)
			if count > 0 {
				layout.peek = uint(count - 1)
			} else {
				layout.peek = 0
			}
			updateInnerHTML(layout.htmlID(), layout.session)
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

func (layout *stackLayoutData) Set(tag string, value interface{}) bool {
	return layout.set(strings.ToLower(tag), value)
}

func (layout *stackLayoutData) set(tag string, value interface{}) bool {
	if tag == TransitionEndEvent {
		listeners, ok := valueToAnimationListeners(value)
		if ok {
			listeners = append(listeners, layout.pushFinished)
			listeners = append(listeners, layout.popFinished)
			layout.properties[TransitionEndEvent] = listeners
			layout.propertyChangedEvent(TransitionEndEvent)
		}
		return ok
	}
	return layout.viewsContainerData.set(tag, value)
}

func (layout *stackLayoutData) Remove(tag string) {
	layout.remove(strings.ToLower(tag))
}

func (layout *stackLayoutData) remove(tag string) {
	if tag == TransitionEndEvent {
		layout.properties[TransitionEndEvent] = []func(View, string){layout.pushFinished, layout.popFinished}
		layout.propertyChangedEvent(TransitionEndEvent)
	} else {
		layout.viewsContainerData.remove(tag)
	}
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
					updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(peek), "visibility", "hidden", layout.Session())
				}

				layout.peek = uint(i)
				updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(i), "visibility", "visible", layout.Session())
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
					updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(peek), "visibility", "hidden", layout.Session())
				}

				layout.peek = uint(i)
				updateCSSProperty(layout.htmlID()+"page"+strconv.Itoa(i), "visibility", "visible", layout.Session())
			}
			return true
		}
	}

	ErrorLogF(`MoveToFront("%s") fail. Subview with "%s" not found."`, viewID, viewID)
	return false
}

func (layout *stackLayoutData) Append(view View) {
	if view != nil {
		layout.peek = uint(len(layout.views))
		layout.viewsContainerData.Append(view)
	} else {
		ErrorLog("StackLayout.Append(nil, ....) is forbidden")
	}
}

func (layout *stackLayoutData) Insert(view View, index uint) {
	if view != nil {
		count := uint(len(layout.views))
		if index < count {
			layout.peek = index
		} else {
			layout.peek = count
		}
		layout.viewsContainerData.Insert(view, index)
	} else {
		ErrorLog("StackLayout.Insert(nil, ....) is forbidden")
	}
}

func (layout *stackLayoutData) RemoveView(index uint) View {
	if layout.peek > 0 {
		layout.peek--
	}
	return layout.viewsContainerData.RemoveView(index)
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
	buffer.WriteString(`push" class="ruiStackPageLayout" ontransitionend="stackTransitionEndEvent(\'`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`\', \'ruiPush\', event)" style="`)

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

	appendToInnerHTML(htmlID, buffer.String(), session)
	updateCSSProperty(htmlID+"push", "transform", "translate(0px, 0px)", layout.session)

	layout.views = append(layout.views, view)
	view.setParentID(htmlID)
	layout.propertyChangedEvent(Content)
}

func (layout *stackLayoutData) Pop(animation int, onPopFinished func(View)) bool {
	count := uint(len(layout.views))
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
	buffer.WriteString(`pop" class="ruiStackPageLayout" ontransitionend="stackTransitionEndEvent(\'`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`\', \'ruiPop\', event)" style="transition: transform 1s ease;">`)
	viewHTML(layout.popView, buffer)
	buffer.WriteString(`</div>`)

	appendToInnerHTML(htmlID, buffer.String(), session)

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

	updateCSSProperty(htmlID+"pop", "transform", value, layout.session)
	layout.propertyChangedEvent(Content)
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
