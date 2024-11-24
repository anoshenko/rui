package rui

import (
	"fmt"
	"strings"
)

// Constants which represent [StackLayout] animation type during pushing or popping views
const (
	// PushTransform is the constant for "push-transform" property tag.
	//
	// Used by `StackLayout`.
	// Specify start translation, scale and rotation over x, y and z axes as well as a distortion
	// for an animated pushing of a child view.
	//
	// Supported types: `TransformProperty`, `string`.
	//
	// See `TransformProperty` description for more details.
	//
	// Conversion rules:
	// `TransformProperty` - stored as is, no conversion performed.
	// `string` - string representation of `Transform` interface. Example: "_{translate-x = 10px, scale-y = 1.1}".
	PushTransform = "push-transform"

	// PushDuration is the constant for "push-duration" property tag.
	//
	// Used by `StackLayout`.
	// Sets the length of time in seconds that an push/pop animation takes to complete.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	PushDuration = "push-duration"

	// PushTiming is the constant for "push-timing" property tag.
	//
	// Used by `StackLayout`.
	// Set how an push/pop animation progresses through the duration of each cycle.
	//
	// Supported types: `string`.
	//
	// Values:
	// "ease"(`EaseTiming`) - Speed increases towards the middle and slows down at the end.
	// "ease-in"(`EaseInTiming`) - Speed is slow at first, but increases in the end.
	// "ease-out"(`EaseOutTiming`) - Speed is fast at first, but decreases in the end.
	// "ease-in-out"(`EaseInOutTiming`) - Speed is slow at first, but quickly increases and at the end it decreases again.
	// "linear"(`LinearTiming`) - Constant speed.
	PushTiming = "push-timing"

	// MoveToFrontAnimation is the constant for "move-to-front-animation" property tag.
	//
	// Used by `StackLayout`.
	// Specifies whether animation is used when calling the MoveToFront/MoveToFrontByID method of StackLayout interface.
	//
	// Supported types: `bool`, `int`, `string`.
	//
	// Values:
	// `true` or `1` or "true", "yes", "on", "1" - animation is used (default value).
	// `false` or `0` or "false", "no", "off", "0" - animation is not used.
	MoveToFrontAnimation = "move-to-front-animation"
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
	Push(view View, onPushFinished func())

	// Pop removes the current View from the container using animation.
	// The second argument `onPopFinished`` is the function to be called when the animation ends. It may be nil.
	// The function will return false if the StackLayout is empty and true if the current item has been removed.
	Pop(onPopFinished func(View)) bool
}

type pushFinished struct {
	peekID   string
	listener func()
}

type popFinished struct {
	view     View
	listener func(View)
}

type stackLayoutData struct {
	viewsContainerData
	onPushFinished map[string]pushFinished
	onPopFinished  map[string]popFinished
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
	layout.onPushFinished = map[string]pushFinished{}
	layout.onPopFinished = map[string]popFinished{}
	layout.set = layout.setFunc
	layout.remove = layout.removeFunc

	layout.setRaw(TransitionEndEvent, []func(View, PropertyName){layout.transitionFinished})
	if session.TextDirection() == RightToLeftDirection {
		layout.setRaw(PushTransform, NewTransformProperty(Params{TranslateX: Percent(-100)}))
	} else {
		layout.setRaw(PushTransform, NewTransformProperty(Params{TranslateX: Percent(100)}))
	}
}

func (layout *stackLayoutData) transitionFinished(view View, tag PropertyName) {
	if tags := strings.Split(string(tag), "-"); len(tags) >= 2 {
		session := layout.Session()
		viewID := tags[1]

		switch tags[0] {
		case "push":
			if finished, ok := layout.onPushFinished[viewID]; ok {
				if finished.peekID != "" {
					pageID := finished.peekID + "page"
					session.startUpdateScript(pageID)
					session.updateCSSProperty(pageID, "visibility", "hidden")
					session.updateCSSProperty(pageID, "transition", "")
					session.updateCSSProperty(pageID, "transform", "")
					session.removeProperty(pageID, "ontransitionend")
					session.removeProperty(pageID, "ontransitioncancel")
					session.finishUpdateScript(pageID)
				}

				pageID := viewID + "page"
				session.startUpdateScript(pageID)
				session.updateCSSProperty(pageID, "z-index", "auto")
				session.updateCSSProperty(pageID, "transition", "")
				session.removeProperty(pageID, "ontransitionend")
				session.removeProperty(pageID, "ontransitioncancel")
				session.finishUpdateScript(pageID)

				if finished.listener != nil {
					finished.listener()
				}
				delete(layout.onPushFinished, viewID)
				layout.contentChanged()
			}

		case "pop":
			if finished, ok := layout.onPopFinished[viewID]; ok {
				session.callFunc("removeView", viewID+"page")
				if finished.listener != nil {
					finished.listener(finished.view)
				}
				delete(layout.onPopFinished, viewID)

				if count := len(layout.views); count > 0 {
					peekID := layout.views[count-1].htmlID() + "page"
					session.startUpdateScript(peekID)
					session.removeProperty(peekID, "ontransitionend")
					session.removeProperty(peekID, "ontransitioncancel")
					session.finishUpdateScript(peekID)
				}
			}

		case "move":
			if count := len(layout.views); count > 1 {
				pageID := layout.views[count-2].htmlID() + "page"
				session.startUpdateScript(pageID)
				session.updateCSSProperty(pageID, "visibility", "hidden")
				session.updateCSSProperty(pageID, "transition", "")
				session.updateCSSProperty(pageID, "transform", "")
				session.finishUpdateScript(pageID)
			}

			pageID := viewID + "page"
			session.startUpdateScript(pageID)
			session.updateCSSProperty(pageID, "z-index", "auto")
			session.updateCSSProperty(pageID, "transition", "")
			session.removeProperty(pageID, "ontransitionend")
			session.removeProperty(pageID, "ontransitioncancel")
			session.finishUpdateScript(pageID)
			layout.contentChanged()
		}
	}
}

func (layout *stackLayoutData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case TransitionEndEvent:
		// TODO
		listeners, ok := valueToOneArgEventListeners[View, PropertyName](value)
		if ok && listeners != nil {
			listeners = append(listeners, layout.transitionFinished)
			layout.setRaw(TransitionEndEvent, listeners)
			return []PropertyName{tag}
		}
		return nil

	case PushTiming:
		if text, ok := value.(string); ok {
			layout.setRaw(tag, text)
			return []PropertyName{tag}
		}
	}
	return layout.viewsContainerData.setFunc(tag, value)
}

func (layout *stackLayoutData) removeFunc(tag PropertyName) []PropertyName {
	switch tag {
	case TransitionEndEvent:
		layout.setRaw(TransitionEndEvent, []func(View, PropertyName){layout.transitionFinished})
		return []PropertyName{tag}
	}
	return layout.viewsContainerData.removeFunc(tag)
}

func (layout *stackLayoutData) Peek() View {
	if count := len(layout.views); count > 0 {
		return layout.views[count-1]
	}
	return nil
}

func (layout *stackLayoutData) MoveToFront(view View) bool {
	if view == nil {
		ErrorLog(`MoveToFront(nil) forbidden`)
		return false
	}

	htmlID := view.htmlID()
	switch count := len(layout.views); count {
	case 0:
		// do nothing

	case 1:
		if layout.views[0].htmlID() == htmlID {
			return true
		}

	default:
		for i, view := range layout.views {
			if view.htmlID() == htmlID {
				layout.moveToFrontByIndex(i)
				return true
			}
		}
	}

	ErrorLog(`MoveToFront() fail. Subview not found.`)
	return false
}

func (layout *stackLayoutData) MoveToFrontByID(viewID string) bool {
	switch count := len(layout.views); count {
	case 0:
		// do nothing

	case 1:
		if layout.views[0].ID() == viewID {
			return true
		}

	default:
		for i, view := range layout.views {
			if view.ID() == viewID {
				layout.moveToFrontByIndex(i)
				return true
			}
		}
	}

	ErrorLogF(`MoveToFront("%s") fail. Subview with "%s" not found.`, viewID, viewID)
	return false
}

func (layout *stackLayoutData) moveToFrontByIndex(index int) {

	count := len(layout.views)
	if index == count-1 {
		return
	}

	view := layout.views[index]
	peekID := layout.views[count-1].htmlID()
	if index == 0 {
		layout.views = append(layout.views[1:], view)
	} else {
		layout.views = append(append(layout.views[:index], layout.views[index+1:]...), view)
	}

	session := layout.Session()
	pageID := view.htmlID() + "page"
	peekPageID := peekID + "page"

	animated := IsMoveToFrontAnimation(layout)

	var transform TransformProperty = nil
	if animated {
		transform = GetPushTransform(layout)
	}

	if transform == nil {
		session.updateCSSProperty(peekPageID, "visibility", "hidden")
		session.updateCSSProperty(pageID, "visibility", "visible")
		layout.contentChanged()
		return
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`stackTransitionEndEvent('`)
	buffer.WriteString(layout.htmlID())
	buffer.WriteString(`', 'move-`)
	buffer.WriteString(view.htmlID())
	buffer.WriteString(`', event)`)

	listener := buffer.String()

	transformCSS := transformMirror(transform, session).transformCSS(session)
	transitionCSS := layout.pushTransitionCSS()

	session.updateCSSProperty(peekPageID, "transition", transitionCSS)

	session.startUpdateScript(pageID)
	session.updateProperty(pageID, "ontransitionend", listener)
	session.updateProperty(pageID, "ontransitioncancel", listener)
	session.updateCSSProperty(pageID, "transform", transformCSS)
	session.updateCSSProperty(pageID, "z-index", "100")
	session.updateCSSProperty(pageID, "visibility", "visible")
	session.finishUpdateScript(pageID)

	session.updateCSSProperty(pageID, "transition", transitionCSS)
	session.updateCSSProperty(pageID, "transform", "")

	session.updateCSSProperty(peekPageID, "transform", transformCSS)
}

func (layout *stackLayoutData) contentChanged() {
	if listener, ok := layout.changeListener[Content]; ok {
		listener(layout, Content)
	}
}

func (layout *stackLayoutData) RemovePeek() View {
	return layout.RemoveView(len(layout.views) - 1)
}

func (layout *stackLayoutData) pushTransitionCSS() string {
	return fmt.Sprintf("transform %.2fs %s", GetPushDuration(layout), GetPushTiming(layout))
}

func transformMirror(transform TransformProperty, session Session) TransformProperty {
	result := NewTransformProperty(nil)

	for _, tag := range []PropertyName{Perspective, RotateX, RotateY, RotateZ, ScaleX, ScaleY, ScaleZ, TranslateZ} {
		if value := transform.getRaw(tag); value != nil {
			result.Set(tag, value)
		}
	}

	for _, tag := range []PropertyName{Rotate, SkewX, SkewY} {
		if angle, ok := angleProperty(transform, tag, session); ok {
			angle.Value = -angle.Value
			result.Set(tag, angle)
		}
	}

	for _, tag := range []PropertyName{TranslateX, TranslateY} {
		if size, ok := sizeProperty(transform, tag, session); ok {
			size.Value = -size.Value
			result.Set(tag, size)
		}
	}

	return result
}

// Append appends a view to the end of the list of a view children
func (layout *stackLayoutData) Append(view View) {
	if view == nil {
		ErrorLog("StackLayout.Append(nil) is forbidden")
		return
	}

	stackID := layout.htmlID()
	view.setParentID(stackID)

	count := len(layout.views)
	if count == 0 {
		layout.views = []View{view}
	} else {
		layout.views = append(layout.views, view)
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`<div id="`)
	buffer.WriteString(view.htmlID())
	buffer.WriteString(`page" class="ruiStackPageLayout">`)
	viewHTML(view, buffer, "")
	buffer.WriteString(`</div>`)

	session := layout.Session()
	if count > 0 {
		session.updateCSSProperty(layout.views[count-1].htmlID()+"page", "visibility", "hidden")
	}
	session.appendToInnerHTML(stackID, buffer.String())

	if listener, ok := layout.changeListener[Content]; ok {
		listener(layout, Content)
	}
}

// Insert inserts a view to the "index" position in the list of a view children
func (layout *stackLayoutData) Insert(view View, index int) {
	if view == nil {
		ErrorLog("StackLayout.Insert(nil, ...) is forbidden")
		return
	}

	if layout.views == nil || index < 0 || index >= len(layout.views) {
		layout.Append(view)
		return
	}

	stackID := layout.htmlID()
	view.setParentID(stackID)
	if index > 0 {
		layout.views = append(layout.views[:index], append([]View{view}, layout.views[index:]...)...)
	} else {
		layout.views = append([]View{view}, layout.views...)
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`<div id="`)
	buffer.WriteString(view.htmlID())
	buffer.WriteString(`page" class="ruiStackPageLayout" style="visibility: hidden;">`)
	viewHTML(view, buffer, "")
	buffer.WriteString(`</div>`)

	session := layout.Session()
	session.appendToInnerHTML(stackID, buffer.String())

	if listener, ok := layout.changeListener[Content]; ok {
		listener(layout, Content)
	}
}

// Remove removes view from list and return it
func (layout *stackLayoutData) RemoveView(index int) View {
	if layout.views == nil {
		layout.views = []View{}
		return nil
	}

	count := len(layout.views)
	if index < 0 || index >= count {
		return nil
	}

	session := layout.Session()
	view := layout.views[index]
	view.setParentID("")

	if count == 1 {
		layout.views = []View{}
	} else if index == 0 {
		layout.views = layout.views[1:]
	} else if index == count-1 {
		layout.views = layout.views[:index]
		session.updateCSSProperty(layout.views[count-2].htmlID()+"page", "visibility", "visible")
	} else {
		layout.views = append(layout.views[:index], layout.views[index+1:]...)
	}

	layout.Session().callFunc("removeView", view.htmlID()+"page")

	if listener, ok := layout.changeListener[Content]; ok {
		listener(layout, Content)
	}
	return view
}

func (layout *stackLayoutData) Push(view View, onPushFinished func()) {
	if view == nil {
		ErrorLog("StackLayout.Push(nil, ....) is forbidden")
		return
	}

	transform := GetPushTransform(layout)
	if transform == nil {
		layout.Append(view)
		if onPushFinished != nil {
			onPushFinished()
		}
		return
	}

	prevPeek := ""
	finished := pushFinished{
		listener: onPushFinished,
	}
	if count := len(layout.views); count > 0 {
		finished.peekID = layout.views[count-1].htmlID()
		prevPeek = finished.peekID + "page"
	}

	htmlID := view.htmlID()
	layout.onPushFinished[htmlID] = finished

	view.setParentID(layout.htmlID())
	layout.views = append(layout.views, view)

	session := layout.Session()

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`<div id="`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`page" class="ruiStackPageLayout" ontransitionend="stackTransitionEndEvent('`)
	buffer.WriteString(layout.htmlID())
	buffer.WriteString(`', 'push-`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`', event)" style="z-index: 100; transform: `)
	buffer.WriteString(transform.transformCSS(layout.session))
	buffer.WriteRune(';')

	transitionCSS := layout.pushTransitionCSS()
	buffer.WriteString(" transition: ")
	buffer.WriteString(transitionCSS)
	buffer.WriteString(`;">`)

	viewHTML(view, buffer, "")
	buffer.WriteString(`</div>`)

	session.appendToInnerHTML(layout.htmlID(), buffer.String())

	if prevPeek != "" {
		mirror := transformMirror(transform, session)
		layout.session.updateCSSProperty(prevPeek, "transition", transitionCSS)
		layout.session.updateCSSProperty(prevPeek, "transform", mirror.transformCSS(session))
	}

	layout.session.updateCSSProperty(htmlID+"page", "transform", "")
}

func (layout *stackLayoutData) Pop(onPopFinished func(View)) bool {
	count := len(layout.views)
	if count == 0 {
		ErrorLog("StackLayout is empty")
		return false
	}

	transform := GetPushTransform(layout)
	if transform == nil {
		if view := layout.RemovePeek(); view != nil {
			if onPopFinished != nil {
				onPopFinished(view)
			}
			return true
		}
		return false
	}

	peek := count - 1
	view := layout.views[peek]
	view.setParentID("")

	layout.views = layout.views[:peek]
	layout.contentChanged()

	layout.onPopFinished[view.htmlID()] = popFinished{
		view:     view,
		listener: onPopFinished,
	}

	htmlID := view.htmlID()
	session := layout.Session()

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`stackTransitionEndEvent('`)
	buffer.WriteString(layout.htmlID())
	buffer.WriteString(`', 'pop-`)
	buffer.WriteString(htmlID)
	buffer.WriteString(`', event)`)

	listener := buffer.String()
	pageID := htmlID + "page"

	transitionCSS := layout.pushTransitionCSS()

	session.startUpdateScript(pageID)
	session.updateProperty(pageID, "ontransitionend", listener)
	session.updateProperty(pageID, "ontransitioncancel", listener)
	session.updateCSSProperty(pageID, "z-index", "100")
	session.updateCSSProperty(pageID, "transition", transitionCSS)
	session.finishUpdateScript(pageID)

	peek--
	if peek >= 0 {
		peekID := layout.views[peek].htmlID() + "page"
		session.updateCSSProperty(peekID, "transition", "")
		session.startUpdateScript(peekID)
		session.updateCSSProperty(peekID, "transform", transformMirror(transform, session).transformCSS(session))
		session.updateCSSProperty(peekID, "visibility", "visible")
		session.finishUpdateScript(peekID)
		session.updateCSSProperty(peekID, "transition", transitionCSS)
		session.updateCSSProperty(peekID, "transform", "")
	}
	session.updateCSSProperty(pageID, "transform", transform.transformCSS(session))

	return true
}

func (layout *stackLayoutData) htmlSubviews(self View, buffer *strings.Builder) {
	if count := len(layout.views); count > 0 {
		peek := count - 1
		for i, view := range layout.views {
			buffer.WriteString(`<div id="`)
			buffer.WriteString(view.htmlID())
			buffer.WriteString(`page`)
			buffer.WriteString(`" class="ruiStackPageLayout"`)
			if i != peek {
				buffer.WriteString(` style="visibility: hidden;"`)
			}
			buffer.WriteString(`>`)
			viewHTML(view, buffer, "")
			buffer.WriteString(`</div>`)
		}
	}
}

// IsMoveToFrontAnimation returns "true" if an animation is used when calling the MoveToFront/MoveToFrontByID method of StackLayout interface.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsMoveToFrontAnimation(view View, subviewID ...string) bool {
	if view = getSubview(view, subviewID); view != nil {
		if value, ok := boolProperty(view, MoveToFrontAnimation, view.Session()); ok {
			return value
		}
		if value := valueFromStyle(view, MoveToFrontAnimation); value != nil {
			if b, ok := valueToBool(value, view.Session()); ok {
				return b
			}
		}
	}

	return true
}

// GetPushDuration returns the length of time in seconds that an push/pop StackLayout animation takes to complete.
// If the second argument (subviewID) is not specified or it is "" then a width of the first argument (view) is returned
func GetPushDuration(view View, subviewID ...string) float64 {
	return floatStyledProperty(view, subviewID, PushDuration, 1)
}

// GetPushTiming returns the function which sets how an push/pop animation progresses.
// If the second argument (subviewID) is not specified or it is "" then a width of the first argument (view) is returned
func GetPushTiming(view View, subviewID ...string) string {
	result := stringStyledProperty(view, subviewID, PushTiming, false)
	if isTimingFunctionValid(result) {
		return result
	}

	return "easy"
}

// GetTransform returns the start transform (translation, scale and rotation over x, y and z axes as well as a distortion)
// for an animated pushing of a child view.
// The default value is nil (no transform).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetPushTransform(view View, subviewID ...string) TransformProperty {
	return transformStyledProperty(view, subviewID, PushTransform)
}
