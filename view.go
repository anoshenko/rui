package rui

import (
	"encoding/base64"
	"fmt"
	"maps"
	"strconv"
	"strings"
)

// Frame - the location and size of a rectangle area
type Frame struct {
	// Left - the left border
	Left float64
	// Top - the top border
	Top float64
	// Width - the width of a rectangle area
	Width float64
	// Height - the height of a rectangle area
	Height float64
}

// Right returns the right border
func (frame Frame) Right() float64 {
	return frame.Left + frame.Width
}

// Bottom returns the bottom border
func (frame Frame) Bottom() float64 {
	return frame.Top + frame.Height
}

const changeListeners PropertyName = "change-listeners"

// View represents a base view interface
type View interface {
	ViewStyle
	fmt.Stringer

	// Session returns the current Session interface
	Session() Session

	// Parent returns the parent view
	Parent() View

	// Tag returns the tag of View interface
	Tag() string

	// ID returns the id of the view
	ID() string

	// Focusable returns true if the view receives the focus
	Focusable() bool

	// Frame returns the location and size of the view in pixels
	Frame() Frame

	// Scroll returns the location size of the scrollable view in pixels
	Scroll() Frame

	// SetParams sets properties with name "tag" of the "rootView" subview. Result:
	//   - true - all properties were set successful,
	//   - false - error (incompatible type or invalid format of a string value, see AppLog).
	SetParams(params Params) bool

	// SetAnimated sets the value (second argument) of the property with name defined by the first argument.
	// Return "true" if the value has been set, in the opposite case "false" are returned and
	// a description of the error is written to the log
	SetAnimated(tag PropertyName, value any, animation AnimationProperty) bool

	// SetChangeListener set the function (the second argument) to track the change of the View property (the first argument).
	//
	// Allowed listener function formats:
	//
	//  func(view rui.View, property rui.PropertyName)
	//  func(view rui.View)
	//  func(property rui.PropertyName)
	//  func()
	//  string
	//
	// If the second argument is given as a string, it specifies the name of the binding function.
	SetChangeListener(tag PropertyName, listener any) bool

	// HasFocus returns 'true' if the view has focus
	HasFocus() bool

	// LoadFile loads the content of the dropped file by drag-and-drop mechanism for all views except FilePicker.
	// The selected file is loaded for FilePicker view.
	//
	// This function is asynchronous. The "result" function will be called after loading the data.
	LoadFile(file FileInfo, result func(FileInfo, []byte))

	init(session Session)
	handleCommand(self View, command PropertyName, data DataObject) bool
	htmlClass(disabled bool) string
	htmlTag() string
	closeHTMLTag() bool
	htmlID() string
	parentHTMLID() string
	setParentID(parentID string)
	htmlSubviews(self View, buffer *strings.Builder)
	htmlProperties(self View, buffer *strings.Builder)
	cssStyle(self View, builder cssBuilder)
	addToCSSStyle(addCSS map[string]string)
	excludeTags() []PropertyName
	htmlDisabledProperty() bool
	binding() any

	onResize(self View, x, y, width, height float64)
	onItemResize(self View, index string, x, y, width, height float64)
	setNoResizeEvent()
	isNoResizeEvent() bool
	setScroll(x, y, width, height float64)
}

// viewData - base implementation of View interface
type viewData struct {
	viewStyle
	session          Session
	tag              string
	viewID           string
	_htmlID          string
	parentID         string
	systemClass      string
	changeListener   map[PropertyName]oneArgListener[View, PropertyName]
	singleTransition map[PropertyName]AnimationProperty
	addCSS           map[string]string
	frame            Frame
	scroll           Frame
	noResizeEvent    bool
	created          bool
	hasFocus         bool
	hasHtmlDisabled  bool
	get              func(tag PropertyName) any
	set              func(tag PropertyName, value any) []PropertyName
	remove           func(tag PropertyName) []PropertyName
	changed          func(tag PropertyName)
	fileLoader       map[string]func(FileInfo, []byte)
}

func newView(session Session) View {
	return new(viewData)
}

// NewView create new View object and return it
func NewView(session Session, params Params) View {
	view := new(viewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func setInitParams(view View, params Params) {
	if params != nil {
		session := view.Session()
		if !session.ignoreViewUpdates() {
			session.setIgnoreViewUpdates(true)
			defer session.setIgnoreViewUpdates(false)
		}
		for _, tag := range params.AllTags() {
			if value, ok := params[tag]; ok {
				view.Set(tag, value)
			}
		}
	}
}

func (view *viewData) init(session Session) {
	view.viewStyle.init()
	view.get = view.getFunc
	view.set = view.setFunc
	view.remove = view.removeFunc
	view.normalize = normalizeViewTag
	view.changed = view.propertyChanged
	view.tag = "View"
	view.session = session
	view.changeListener = map[PropertyName]oneArgListener[View, PropertyName]{}
	view.addCSS = map[string]string{}
	//view.animation = map[string]AnimationEndListener{}
	view.singleTransition = map[PropertyName]AnimationProperty{}
	view.noResizeEvent = false
	view.created = false
	view.hasHtmlDisabled = false
	view.fileLoader = map[string]func(FileInfo, []byte){}
}

func (view *viewData) Session() Session {
	return view.session
}

func (view *viewData) Parent() View {
	return view.session.viewByHTMLID(view.parentID)
}

func (view *viewData) parentHTMLID() string {
	return view.parentID
}

func (view *viewData) setParentID(parentID string) {
	view.parentID = parentID
}

func (view *viewData) Tag() string {
	return view.tag
}

func (view *viewData) ID() string {
	return view.viewID
}

func (view *viewData) binding() any {
	if result := view.getRaw(Binding); result != nil {
		return result
	}

	if parent := view.Parent(); parent != nil {
		return parent.binding()
	}

	return nil
}

func (view *viewData) ViewByID(id string) View {
	if id == view.ID() {
		if v := view.session.viewByHTMLID(view.htmlID()); v != nil {
			return v
		}
		return view
	}
	return nil
}

func (view *viewData) Focusable() bool {
	if focus, ok := boolProperty(view, Focusable, view.session); ok {
		return focus
	}

	if style, ok := stringProperty(view, Style, view.session); ok {
		if style, ok := view.session.resolveConstants(style); ok {
			if value := view.session.styleProperty(style, Focusable); ok {
				if focus, ok := valueToBool(value, view.Session()); ok {
					return focus
				}
			}
		}
	}

	return false
}

func (view *viewData) Remove(tag PropertyName) {
	changedTags := view.removeFunc(view.normalize(tag))

	if view.created && len(changedTags) > 0 {
		for _, tag := range changedTags {
			view.changed(tag)
			if listener, ok := view.changeListener[tag]; ok {
				listener.Run(view, tag)
			}
		}
	}
}

func (view *viewData) Get(tag PropertyName) any {
	switch tag {
	case ID:
		return view.ID()
	}
	return view.get(view.normalize(tag))

}

func (view *viewData) Set(tag PropertyName, value any) bool {
	if value == nil {
		view.Remove(tag)
		return true
	}

	tag = view.normalize(tag)
	changedTags := view.set(tag, value)

	if view.created && len(changedTags) > 0 {
		for _, tag := range changedTags {
			view.changed(tag)
			if listener, ok := view.changeListener[tag]; ok {
				listener.Run(view, tag)
			}
		}
	}

	return changedTags != nil
}

func normalizeViewTag(tag PropertyName) PropertyName {
	tag = normalizeViewStyleTag(tag)
	switch tag {
	case "tab-index":
		return TabIndex
	}
	return tag
}

func (view *viewData) getFunc(tag PropertyName) any {
	switch tag {
	case ID:
		if id := view.ID(); id != "" {
			return id
		} else {
			return nil
		}
	}
	return viewStyleGet(view, tag)
}

func (view *viewData) removeFunc(tag PropertyName) []PropertyName {
	var changedTags []PropertyName = nil

	switch tag {
	case ID:
		if view.viewID != "" {
			view.viewID = ""
			changedTags = []PropertyName{ID}
		} else {
			changedTags = []PropertyName{}
		}

	case Binding:
		if view.getRaw(Binding) != nil {
			view.setRaw(Binding, nil)
			changedTags = []PropertyName{Binding}
		} else {
			changedTags = []PropertyName{}
		}

	case Animation:
		if val := view.getRaw(Animation); val != nil {
			if animations, ok := val.([]AnimationProperty); ok {
				for _, animation := range animations {
					animation.unused(view.session)
				}
			}

			view.setRaw(Animation, nil)
			changedTags = []PropertyName{Animation}
		}

	default:
		changedTags = viewStyleRemove(view, tag)
	}

	return changedTags
}

func (view *viewData) setFunc(tag PropertyName, value any) []PropertyName {

	switch tag {

	case ID:
		if text, ok := value.(string); ok {
			view.viewID = text
			view.setRaw(ID, text)
			return []PropertyName{ID}
		}
		notCompatibleType(ID, value)
		return nil

	case Binding:
		view.setRaw(Binding, value)
		return []PropertyName{Binding}

	case changeListeners:
		switch value := value.(type) {
		case DataObject:
			for i := range value.PropertyCount() {
				node := value.Property(i)
				if node.Type() == TextNode {
					if text := node.Text(); text != "" {
						view.changeListener[PropertyName(node.Tag())] = newOneArgListenerBinding[View, PropertyName](text)
					}
				}
			}
			if len(view.changeListener) > 0 {
				view.setRaw(changeListeners, view.changeListener)
			}
			return []PropertyName{tag}

		case DataNode:
			if value.Type() == ObjectNode {
				return view.setFunc(tag, value.Object())
			}
		}
		return []PropertyName{}

	case Animation:
		oldAnimations := []AnimationProperty{}
		if val := view.getRaw(Animation); val != nil {
			if animation, ok := val.([]AnimationProperty); ok {
				oldAnimations = animation
			}
		}

		if !setAnimationProperty(view, tag, value) {
			return nil
		}

		for _, animation := range oldAnimations {
			animation.unused(view.session)
		}
		return []PropertyName{Animation}

	case TabIndex, "tab-index":
		return setIntProperty(view, TabIndex, value)

	case UserData:
		view.setRaw(tag, value)
		return []PropertyName{UserData}

	case Style, StyleDisabled:
		if text, ok := value.(string); ok {
			view.setRaw(tag, text)
			return []PropertyName{tag}
		}
		notCompatibleType(ID, value)
		return nil

	case FocusEvent, LostFocusEvent:
		return setNoArgEventListener[View](view, tag, value)

	case KeyDownEvent, KeyUpEvent:
		return setOneArgEventListener[View, KeyEvent](view, tag, value)

	case ClickEvent, DoubleClickEvent, MouseDown, MouseUp, MouseMove, MouseOut, MouseOver, ContextMenuEvent:
		return setOneArgEventListener[View, MouseEvent](view, tag, value)

	case PointerDown, PointerUp, PointerMove, PointerOut, PointerOver, PointerCancel:
		return setOneArgEventListener[View, PointerEvent](view, tag, value)

	case TouchStart, TouchEnd, TouchMove, TouchCancel:
		return setOneArgEventListener[View, TouchEvent](view, tag, value)

	case TransitionRunEvent, TransitionStartEvent, TransitionEndEvent, TransitionCancelEvent:
		result := setOneArgEventListener[View, PropertyName](view, tag, value)
		if result == nil {
			if listeners, ok := valueToOneArgEventListeners[View, string](view); ok && len(listeners) > 0 {
				newListeners := make([]oneArgListener[View, PropertyName], len(listeners))
				for i, listener := range listeners {
					newListeners[i] = newOneArgListenerVE(func(view View, name PropertyName) {
						listener.Run(view, string(name))
					})
				}
				view.setRaw(tag, newListeners)
				result = []PropertyName{tag}
			}
		}
		return result

	case AnimationStartEvent, AnimationEndEvent, AnimationIterationEvent, AnimationCancelEvent:
		return setOneArgEventListener[View, string](view, tag, value)

	case ResizeEvent, ScrollEvent:
		return setOneArgEventListener[View, Frame](view, tag, value)

	case DragData:
		switch value := value.(type) {
		case map[string]string:
			if len(value) == 0 {
				view.setRaw(DragData, nil)
			} else {
				view.setRaw(DragData, maps.Clone(value))
			}

		case string:
			if value == "" {
				view.setRaw(DragData, nil)
			} else {
				data := map[string]string{}
				for _, line := range strings.Split(value, ";") {
					index := strings.IndexRune(line, ':')
					if index < 0 {
						invalidPropertyValue(DragData, value)
						return nil
					}
					mime := line[:index]
					val := line[index+1:]
					if len(mime) > 0 || len(val) > 0 {
						data[mime] = val
					}
				}

				if len(data) == 0 {
					view.setRaw(DragData, nil)
				} else {
					view.setRaw(DragData, data)
				}
			}

		case DataObject:
			data := map[string]string{}
			count := value.PropertyCount()
			for i := range count {
				node := value.Property(i)
				if node.Type() == TextNode {
					data[node.Tag()] = node.Text()
				} else {
					invalidPropertyValue(DragData, value)
					return nil
				}
			}
			if len(data) == 0 {
				view.setRaw(DragData, nil)
			} else {
				view.setRaw(DragData, data)
			}

		case DataNode:
			switch value.Type() {
			case TextNode:
				return view.setFunc(DragData, value.Text())

			case ObjectNode:
				return view.setFunc(DragData, value.Object())
			}
			invalidPropertyValue(DragData, value)
			return nil

		case DataValue:
			if value.IsObject() {
				return view.setFunc(DragData, value.Object())
			}
			return view.setFunc(DragData, value.Value())

		default:
			notCompatibleType(DragData, value)
		}

		return []PropertyName{DragData}

	case DragStartEvent, DragEndEvent, DragEnterEvent, DragLeaveEvent, DragOverEvent, DropEvent:
		return setOneArgEventListener[View, DragAndDropEvent](view, tag, value)

	case DropEffect:
		return view.setDropEffect(value)

	case DropEffectAllowed:
		return view.setDropEffectAllowed(value)
	}

	return viewStyleSet(view, tag, value)
}

func (view *viewData) SetParams(params Params) bool {
	if params == nil {
		errorLog("Argument of function SetParams is nil")
		return false
	}

	session := view.Session()
	session.startUpdateScript(view.htmlID())
	result := true
	for _, tag := range params.AllTags() {
		if value, ok := params[tag]; ok {
			result = view.Set(tag, value) && result
		}
	}
	session.finishUpdateScript(view.htmlID())
	return result
}

func (view *viewData) propertyChanged(tag PropertyName) {

	htmlID := view.htmlID()
	session := view.Session()

	switch tag {
	case TabIndex:
		if value, ok := intProperty(view, TabIndex, session, 0); ok {
			session.updateProperty(htmlID, "tabindex", strconv.Itoa(value))
		} else if view.Focusable() {
			session.updateProperty(htmlID, "tabindex", "0")
		} else {
			session.updateProperty(htmlID, "tabindex", "-1")
		}

	case Style, StyleDisabled:
		session.updateProperty(htmlID, "class", view.htmlClass(IsDisabled(view)))

	case Disabled:
		tabIndex := GetTabIndex(view, htmlID)
		enabledClass := view.htmlClass(false)
		disabledClass := view.htmlClass(true)
		session.startUpdateScript(htmlID)
		if IsDisabled(view) {
			session.updateProperty(htmlID, "data-disabled", "1")
			if view.htmlDisabledProperty() {
				session.updateProperty(htmlID, "disabled", true)
			}
			if tabIndex >= 0 {
				session.updateProperty(htmlID, "tabindex", -1)
			}
			if enabledClass != disabledClass {
				session.updateProperty(htmlID, "class", disabledClass)
			}
		} else {
			session.updateProperty(htmlID, "data-disabled", "0")
			if view.htmlDisabledProperty() {
				session.removeProperty(htmlID, "disabled")
			}
			if tabIndex >= 0 {
				session.updateProperty(htmlID, "tabindex", tabIndex)
			}
			if enabledClass != disabledClass {
				session.updateProperty(htmlID, "class", enabledClass)
			}
		}
		session.finishUpdateScript(htmlID)
		updateInnerHTML(htmlID, session)

	case Visibility:
		switch GetVisibility(view) {
		case Invisible:
			session.updateCSSProperty(htmlID, string(Visibility), "hidden")
			session.updateCSSProperty(htmlID, "display", "")
			session.callFunc("hideTooltip")

		case Gone:
			session.updateCSSProperty(htmlID, string(Visibility), "hidden")
			session.updateCSSProperty(htmlID, "display", "none")
			session.callFunc("hideTooltip")

		default:
			session.updateCSSProperty(htmlID, string(Visibility), "visible")
			session.updateCSSProperty(htmlID, "display", "")
		}

	case Background:
		session.updateCSSProperty(htmlID, string(Background), backgroundCSS(view, session))

	case Mask:
		session.updateCSSProperty(htmlID, "mask", maskCSS(view, session))

	case Border, BorderLeft, BorderRight, BorderTop, BorderBottom:
		cssWidth := ""
		cssColor := ""
		cssStyle := "none"

		if border := getBorderProperty(view, Border); border != nil {
			cssWidth = border.cssWidthValue(session)
			cssColor = border.cssColorValue(session)
			cssStyle = border.cssStyleValue(session)
		}

		session.updateCSSProperty(htmlID, string(BorderWidth), cssWidth)
		session.updateCSSProperty(htmlID, string(BorderColor), cssColor)
		session.updateCSSProperty(htmlID, string(BorderStyle), cssStyle)

	case BorderStyle, BorderLeftStyle, BorderRightStyle, BorderTopStyle, BorderBottomStyle:
		if border := getBorderProperty(view, Border); border != nil {
			session.updateCSSProperty(htmlID, string(BorderStyle), border.cssStyleValue(session))
		}

	case BorderColor, BorderLeftColor, BorderRightColor, BorderTopColor, BorderBottomColor:
		if border := getBorderProperty(view, Border); border != nil {
			session.updateCSSProperty(htmlID, string(BorderColor), border.cssColorValue(session))
		}

	case BorderWidth, BorderLeftWidth, BorderRightWidth, BorderTopWidth, BorderBottomWidth:
		if border := getBorderProperty(view, Border); border != nil {
			session.updateCSSProperty(htmlID, string(BorderWidth), border.cssWidthValue(session))
		}

	case Outline, OutlineColor, OutlineStyle, OutlineWidth:
		session.updateCSSProperty(htmlID, string(Outline), GetOutline(view).cssString(session))

	case Shadow:
		session.updateCSSProperty(htmlID, "box-shadow", shadowCSS(view, Shadow, session))

	case TextShadow:
		session.updateCSSProperty(htmlID, "text-shadow", shadowCSS(view, TextShadow, session))

	case Radius, RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		radius := GetRadius(view)
		session.updateCSSProperty(htmlID, "border-radius", radius.cssString(session))

	case Margin, MarginTop, MarginRight, MarginBottom, MarginLeft,
		"top-margin", "right-margin", "bottom-margin", "left-margin":
		margin := GetMargin(view)
		session.updateCSSProperty(htmlID, string(Margin), margin.cssString(session))

	case Padding, PaddingTop, PaddingRight, PaddingBottom, PaddingLeft,
		"top-padding", "right-padding", "bottom-padding", "left-padding":
		padding := GetPadding(view)
		session.updateCSSProperty(htmlID, string(Padding), padding.cssString(session))

	case AvoidBreak:
		if avoid, ok := boolProperty(view, AvoidBreak, session); ok {
			if avoid {
				session.updateCSSProperty(htmlID, "break-inside", "avoid")
			} else {
				session.updateCSSProperty(htmlID, "break-inside", "auto")
			}
		}

	case Clip:
		if clip := getClipShapeProperty(view, Clip, session); clip != nil && clip.valid(session) {
			session.updateCSSProperty(htmlID, `clip-path`, clip.cssStyle(session))
		} else {
			session.updateCSSProperty(htmlID, `clip-path`, "none")
		}

	case ShapeOutside:
		if clip := getClipShapeProperty(view, ShapeOutside, session); clip != nil && clip.valid(session) {
			session.updateCSSProperty(htmlID, string(ShapeOutside), clip.cssStyle(session))
		} else {
			session.updateCSSProperty(htmlID, string(ShapeOutside), "none")
		}

	case Filter:
		text := ""
		if value := view.getRaw(Filter); value != nil {
			if filter, ok := value.(FilterProperty); ok {
				text = filter.cssStyle(session)
			}
		}
		session.updateCSSProperty(htmlID, string(Filter), text)

	case BackdropFilter:
		text := ""
		if value := view.getRaw(BackdropFilter); value != nil {
			if filter, ok := value.(FilterProperty); ok {
				text = filter.cssStyle(session)
			}
		}
		if session.startUpdateScript(htmlID) {
			defer session.finishUpdateScript(htmlID)
		}
		session.updateCSSProperty(htmlID, "-webkit-backdrop-filter", text)
		session.updateCSSProperty(htmlID, string(BackdropFilter), text)

	case FontName:
		if font, ok := stringProperty(view, FontName, session); ok {
			session.updateCSSProperty(htmlID, "font-family", font)
		} else {
			session.updateCSSProperty(htmlID, "font-family", "")
		}

	case Italic:
		if state, ok := boolProperty(view, tag, session); ok {
			if state {
				session.updateCSSProperty(htmlID, "font-style", "italic")
			} else {
				session.updateCSSProperty(htmlID, "font-style", "normal")
			}
		} else {
			session.updateCSSProperty(htmlID, "font-style", "")
		}

	case SmallCaps:
		if state, ok := boolProperty(view, tag, session); ok {
			if state {
				session.updateCSSProperty(htmlID, "font-variant", "small-caps")
			} else {
				session.updateCSSProperty(htmlID, "font-variant", "normal")
			}
		} else {
			session.updateCSSProperty(htmlID, "font-variant", "")
		}

	case Strikethrough, Overline, Underline:
		session.updateCSSProperty(htmlID, "text-decoration", textDecorationCSS(view, session))
		/*
			for _, tag2 := range []PropertyName{TextLineColor, TextLineStyle, TextLineThickness} {
				view.propertyChanged(tag2)
			}
		*/

	case Transition:
		session.updateCSSProperty(htmlID, "transition", transitionCSS(view, session))

	case Animation:
		session.updateCSSProperty(htmlID, "animation", animationCSS(view, session))

	case AnimationPaused:
		paused, ok := boolProperty(view, AnimationPaused, session)
		if !ok {
			session.updateCSSProperty(htmlID, `animation-play-state`, ``)
		} else if paused {
			session.updateCSSProperty(htmlID, `animation-play-state`, `paused`)
		} else {
			session.updateCSSProperty(htmlID, `animation-play-state`, `running`)
		}

	case ZIndex, Order, TabSize:
		if i, ok := intProperty(view, tag, session, 0); ok {
			session.updateCSSProperty(htmlID, string(tag), strconv.Itoa(i))
		} else {
			session.updateCSSProperty(htmlID, string(tag), "")
		}

	case Row, Column:
		if parentID := view.parentHTMLID(); parentID != "" {
			updateInnerHTML(parentID, session)
		}

	case UserSelect:
		if session.startUpdateScript(htmlID) {
			defer session.finishUpdateScript(htmlID)
		}
		if userSelect, ok := boolProperty(view, UserSelect, session); ok {
			if userSelect {
				session.updateCSSProperty(htmlID, "-webkit-user-select", "auto")
				session.updateCSSProperty(htmlID, "user-select", "auto")
			} else {
				session.updateCSSProperty(htmlID, "-webkit-user-select", "none")
				session.updateCSSProperty(htmlID, "user-select", "none")
			}
		} else {
			session.updateCSSProperty(htmlID, "-webkit-user-select", "")
			session.updateCSSProperty(htmlID, "user-select", "")
		}

	case ColumnSpanAll:
		if spanAll, ok := boolProperty(view, ColumnSpanAll, session); ok && spanAll {
			session.updateCSSProperty(htmlID, `column-span`, `all`)
		} else {
			session.updateCSSProperty(htmlID, `column-span`, `none`)
		}

	case Tooltip:
		if tooltip := GetTooltip(view); tooltip == "" {
			session.removeProperty(htmlID, "data-tooltip")
		} else {
			session.updateProperty(htmlID, "data-tooltip", tooltip)
			session.updateProperty(htmlID, "onmouseenter", "mouseEnterEvent(this, event)")
			session.updateProperty(htmlID, "onmouseleave", "mouseLeaveEvent(this, event)")
		}

	case PerspectiveOriginX, PerspectiveOriginY:
		x, y := GetPerspectiveOrigin(view)
		session.updateCSSProperty(htmlID, "perspective-origin", transformOriginCSS(x, y, AutoSize(), session))

	case BackfaceVisible:
		if GetBackfaceVisible(view) {
			session.updateCSSProperty(htmlID, string(BackfaceVisible), "visible")
		} else {
			session.updateCSSProperty(htmlID, string(BackfaceVisible), "hidden")
		}

	case TransformOriginX, TransformOriginY, TransformOriginZ:
		x, y, z := getTransformOrigin(view, session)
		session.updateCSSProperty(htmlID, "transform-origin", transformOriginCSS(x, y, z, session))

	case Transform:
		css := ""
		if transform := getTransformProperty(view, Transform); transform != nil {
			css = transform.transformCSS(session)
		}
		session.updateCSSProperty(htmlID, "transform", css)

	case Perspective, SkewX, SkewY, TranslateX, TranslateY, TranslateZ,
		ScaleX, ScaleY, ScaleZ, Rotate, RotateX, RotateY, RotateZ:
		// do nothing

	case FocusEvent, LostFocusEvent, ResizeEvent, ScrollEvent, KeyDownEvent, KeyUpEvent,
		ClickEvent, DoubleClickEvent, MouseDown, MouseUp, MouseMove, MouseOut, MouseOver, ContextMenuEvent,
		PointerDown, PointerUp, PointerMove, PointerOut, PointerOver, PointerCancel,
		TouchStart, TouchEnd, TouchMove, TouchCancel,
		TransitionRunEvent, TransitionStartEvent, TransitionEndEvent, TransitionCancelEvent,
		AnimationStartEvent, AnimationEndEvent, AnimationIterationEvent, AnimationCancelEvent,
		DragEndEvent, DragEnterEvent, DragLeaveEvent:

		updateEventListenerHtml(view, tag)

	case DragStartEvent:
		if view.getRaw(DragStartEvent) != nil || view.getRaw(DragData) != nil {
			session.updateProperty(htmlID, "ondragstart", "dragStartEvent(this, event)")
		} else {
			session.removeProperty(htmlID, "ondragstart")
		}

	case DropEvent:
		if view.getRaw(DropEvent) != nil {
			session.updateProperty(htmlID, "ondrop", "dropEvent(this, event)")
			session.updateProperty(htmlID, "ondragover", "dragOverEvent(this, event)")
			if view.getRaw(DragOverEvent) != nil {
				session.updateProperty(htmlID, "data-drag-over", "1")
			} else {
				session.removeProperty(htmlID, "data-drag-over")
			}
		} else {
			session.removeProperty(htmlID, "ondrop")
			session.removeProperty(htmlID, "ondragover")
		}

	case DragOverEvent:
		if view.getRaw(DragOverEvent) != nil {
			session.updateProperty(htmlID, "data-drag-over", "1")
		} else {
			session.removeProperty(htmlID, "data-drag-over")
		}

	case DragData:
		if data := base64DragData(view); data != "" {
			session.updateProperty(htmlID, "draggable", "true")
			session.updateProperty(htmlID, "data-drag", data)
			session.updateProperty(htmlID, "ondragstart", "dragStartEvent(this, event)")
		} else {
			session.removeProperty(htmlID, "draggable")
			session.removeProperty(htmlID, "data-drag")
			if view.getRaw(DragStartEvent) == nil {
				session.removeProperty(htmlID, "ondragstart")
			}
		}

	case DragImage:
		if img, ok := stringProperty(view, DragImage, session); ok && img != "" {
			img = strings.Trim(img, " \t")
			if img[0] == '@' {
				img, ok = session.ImageConstant(img[1:])
				if !ok {
					session.removeProperty(htmlID, "data-drag-image")
					return
				}
			}
			session.updateProperty(htmlID, "data-drag-image", img)
		} else {
			session.removeProperty(htmlID, "data-drag-image")
		}

	case DragImageXOffset:
		if f := GetDragImageXOffset(view); f != 0 {
			session.updateProperty(htmlID, "data-drag-image-x", f)
		} else {
			session.removeProperty(htmlID, "data-drag-image-x")
		}

	case DragImageYOffset:
		if f := GetDragImageXOffset(view); f != 0 {
			session.updateProperty(htmlID, "data-drag-image-y", f)
		} else {
			session.removeProperty(htmlID, "data-drag-image-y")
		}

	case DropEffect:
		effect := GetDropEffect(view)
		switch effect {
		case DropEffectCopy:
			session.updateProperty(htmlID, "data-drop-effect", "copy")
		case DropEffectMove:
			session.updateProperty(htmlID, "data-drop-effect", "move")
		case DropEffectLink:
			session.updateProperty(htmlID, "data-drop-effect", "link")
		default:
			session.removeProperty(htmlID, "data-drop-effect")
		}

	case DropEffectAllowed:
		effect := GetDropEffectAllowed(view)
		if effect >= DropEffectCopy && effect >= DropEffectAll {
			values := []string{"undefined", "copy", "move", "copyMove", "link", "copyLink", "linkMove", "all"}
			session.updateProperty(htmlID, "data-drop-effect-allowed", values[effect])
		} else {
			session.removeProperty(htmlID, "data-drop-effect-allowed")
		}

	case DataList:
		updateInnerHTML(htmlID, session)

	case Opacity:
		if f, ok := floatTextProperty(view, Opacity, session, 0); ok {
			session.updateCSSProperty(htmlID, string(Opacity), f)
		} else {
			session.updateCSSProperty(htmlID, string(Opacity), "")
		}

	default:
		if cssTag, ok := sizeProperties[tag]; ok {
			if size, ok := sizeProperty(view, tag, session); ok {
				session.updateCSSProperty(htmlID, cssTag, size.cssString("", session))
			} else {
				session.updateCSSProperty(htmlID, cssTag, "")
			}
			return
		}

		colorTags := map[PropertyName]string{
			BackgroundColor: string(BackgroundColor),
			TextColor:       "color",
			TextLineColor:   "text-decoration-color",
			CaretColor:      string(CaretColor),
			AccentColor:     string(AccentColor),
		}
		if cssTag, ok := colorTags[tag]; ok {
			if color, ok := colorProperty(view, tag, session); ok {
				session.updateCSSProperty(htmlID, cssTag, color.cssString())
			} else {
				session.updateCSSProperty(htmlID, cssTag, "")
			}
			return
		}

		if valuesData, ok := enumProperties[tag]; ok && valuesData.cssTag != "" {
			if n, ok := enumProperty(view, tag, session, 0); ok {
				session.updateCSSProperty(htmlID, valuesData.cssTag, valuesData.cssValues[n])
			} else {
				session.updateCSSProperty(htmlID, valuesData.cssTag, "")
			}
			return
		}
	}
}

func (view *viewData) htmlTag() string {
	if semantics := GetSemantics(view); semantics > DefaultSemantics {
		values := enumProperties[Semantics].cssValues
		if semantics < len(values) {
			return values[semantics]
		}
	}
	return "div"
}

func (view *viewData) closeHTMLTag() bool {
	return true
}

func (view *viewData) htmlID() string {
	if view._htmlID == "" {
		view._htmlID = view.session.nextViewID()
	}
	return view._htmlID
}

func (view *viewData) htmlSubviews(self View, buffer *strings.Builder) {
}

func (view *viewData) addToCSSStyle(addCSS map[string]string) {
	view.addCSS = addCSS
}

func (view *viewData) cssStyle(self View, builder cssBuilder) {
	view.viewStyle.cssViewStyle(builder, view.session)
	if view.addCSS != nil {
		for tag, value := range view.addCSS {
			builder.add(tag, value)
		}
	}
}

func (view *viewData) htmlDisabledProperty() bool {
	return view.hasHtmlDisabled
}

func (view *viewData) htmlProperties(self View, buffer *strings.Builder) {
	view.created = true

	if IsDisabled(self) {
		buffer.WriteString(` data-disabled="1"`)
		if view.hasHtmlDisabled {
			buffer.WriteString(`  disabled`)
		}
	} else {
		buffer.WriteString(` data-disabled="0"`)
	}

	if view.frame.Left != 0 || view.frame.Top != 0 || view.frame.Width != 0 || view.frame.Height != 0 {
		fmt.Fprintf(buffer, ` data-left="%g" data-top="%g" data-width="%g" data-height="%g"`,
			view.frame.Left, view.frame.Top, view.frame.Width, view.frame.Height)
	}
}

func viewHTML(view View, buffer *strings.Builder, htmlTag string) {
	if htmlTag == "" {
		htmlTag = view.htmlTag()
	}
	//viewHTMLTag := view.htmlTag()
	buffer.WriteRune('<')
	buffer.WriteString(htmlTag)
	buffer.WriteString(` id="`)
	buffer.WriteString(view.htmlID())
	buffer.WriteRune('"')

	disabled := IsDisabled(view)

	if cls := view.htmlClass(disabled); cls != "" {
		buffer.WriteString(` class="`)
		buffer.WriteString(cls)
		buffer.WriteRune('"')
	}

	cssBuilder := viewCSSBuilder{buffer: allocStringBuilder()}
	view.cssStyle(view, &cssBuilder)

	if style := cssBuilder.finish(); style != "" {
		buffer.WriteString(` style="`)
		buffer.WriteString(style)
		buffer.WriteRune('"')
	}

	buffer.WriteRune(' ')
	view.htmlProperties(view, buffer)

	if view.isNoResizeEvent() {
		buffer.WriteString(` data-noresize="1" `)
	} else {
		buffer.WriteRune(' ')
	}

	if !disabled {
		if tabIndex := GetTabIndex(view); tabIndex >= 0 {
			buffer.WriteString(`tabindex="`)
			buffer.WriteString(strconv.Itoa(tabIndex))
			buffer.WriteString(`" `)
		}
	}

	if tooltip := GetTooltip(view); tooltip != "" {
		buffer.WriteString(`data-tooltip=" `)
		buffer.WriteString(tooltip)
		buffer.WriteString(`" onmouseenter="mouseEnterEvent(this, event)" onmouseleave="mouseLeaveEvent(this, event)" `)
	}

	buffer.WriteString(`onscroll="scrollEvent(this, event)" `)

	focusEventsHtml(view, buffer)
	keyEventsHtml(view, buffer)

	viewEventsHtml[MouseEvent](view, []PropertyName{ClickEvent, DoubleClickEvent, MouseDown, MouseUp, MouseMove, MouseOut, MouseOver, ContextMenuEvent}, buffer)
	viewEventsHtml[PointerEvent](view, []PropertyName{PointerDown, PointerUp, PointerMove, PointerOut, PointerOver, PointerCancel}, buffer)
	viewEventsHtml[TouchEvent](view, []PropertyName{TouchStart, TouchEnd, TouchMove, TouchCancel}, buffer)
	viewEventsHtml[string](view, []PropertyName{TransitionRunEvent, TransitionStartEvent, TransitionEndEvent, TransitionCancelEvent,
		AnimationStartEvent, AnimationEndEvent, AnimationIterationEvent, AnimationCancelEvent}, buffer)

	dragAndDropHtml(view, buffer)

	buffer.WriteRune('>')
	view.htmlSubviews(view, buffer)
	if view.closeHTMLTag() {
		buffer.WriteString(`</`)
		buffer.WriteString(htmlTag)
		buffer.WriteRune('>')
	}
}

func (view *viewData) htmlClass(disabled bool) string {
	cls := "ruiView"
	disabledStyle := false
	if disabled {
		if value, ok := stringProperty(view, StyleDisabled, view.Session()); ok && value != "" {
			cls += " " + value
			disabledStyle = true
		}
	}
	if !disabledStyle {
		if value, ok := stringProperty(view, Style, view.Session()); ok {
			cls += " " + value
		}
	}

	if view.systemClass != "" {
		cls = view.systemClass + " " + cls
	}

	return cls
}

func (view *viewData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {

	case KeyDownEvent, KeyUpEvent:
		if !IsDisabled(self) {
			handleKeyEvents(self, command, data)
		}

	case ClickEvent, DoubleClickEvent, MouseDown, MouseUp, MouseMove, MouseOut, MouseOver, ContextMenuEvent:
		handleMouseEvents(self, command, data)

	case PointerDown, PointerUp, PointerMove, PointerOut, PointerOver, PointerCancel:
		handlePointerEvents(self, command, data)

	case TouchStart, TouchEnd, TouchMove, TouchCancel:
		handleTouchEvents(self, command, data)

	case DragStartEvent:
		handleDragAndDropEvents(self, command, data)

	case DragEndEvent, DragEnterEvent, DragLeaveEvent, DragOverEvent, DropEvent:
		handleDragAndDropEvents(self, command, data)

	case FocusEvent:
		view.hasFocus = true
		for _, listener := range getNoArgEventListeners[View](view, nil, command) {
			listener.Run(self)
		}

	case LostFocusEvent:
		view.hasFocus = false
		for _, listener := range getNoArgEventListeners[View](view, nil, command) {
			listener.Run(self)
		}

	case TransitionRunEvent, TransitionStartEvent, TransitionEndEvent, TransitionCancelEvent:
		view.handleTransitionEvents(command, data)

	case AnimationStartEvent, AnimationEndEvent, AnimationIterationEvent, AnimationCancelEvent:
		view.handleAnimationEvents(command, data)

	case "scroll":
		view.onScroll(view, dataFloatProperty(data, "x"), dataFloatProperty(data, "y"), dataFloatProperty(data, "width"), dataFloatProperty(data, "height"))

	case "widthChanged":
		if value, ok := data.PropertyValue("width"); ok {
			if width, ok := StringToSizeUnit(value); ok {
				self.setRaw(Width, width)
			}
		}

	case "heightChanged":
		if value, ok := data.PropertyValue("height"); ok {
			if height, ok := StringToSizeUnit(value); ok {
				self.setRaw(Height, height)
			}
		}

		/*
			case "resize":
				floatProperty := func(tag string) float64 {
					if value, ok := data.PropertyValue(tag); ok {
						if result, err := strconv.ParseFloat(value, 64); err == nil {
							return result
						}
					}
					return 0
				}

				self.onResize(self, floatProperty("x"), floatProperty("y"), floatProperty("width"), floatProperty("height"))
				return true
		*/

	case "fileLoaded":
		file := dataToFileInfo(data)
		key := file.key()

		if listener := view.fileLoader[key]; listener != nil {
			delete(view.fileLoader, key)

			if base64Data, ok := data.PropertyValue("data"); ok {
				if index := strings.LastIndex(base64Data, ","); index >= 0 {
					base64Data = base64Data[index+1:]
				}
				decode, err := base64.StdEncoding.DecodeString(base64Data)
				if err == nil {
					listener(file, decode)
				} else {
					ErrorLog(err.Error())
				}
			}
		}
		return true

	case "fileLoadingError":
		file := dataToFileInfo(data)
		key := file.key()

		if error, ok := data.PropertyValue("error"); ok {
			ErrorLogF(`Load "%s" file error: %s`, file.Name, error)
		}

		if listener := view.fileLoader[key]; listener != nil {
			delete(view.fileLoader, key)
			listener(file, nil)
		}
		return true

	default:
		return false
	}
	return true

}

func (view *viewData) SetChangeListener(tag PropertyName, listener any) bool {
	if listener == nil {
		delete(view.changeListener, tag)
	} else {
		switch listener := listener.(type) {
		case func():
			view.changeListener[tag] = newOneArgListener0[View, PropertyName](listener)

		case func(View):
			view.changeListener[tag] = newOneArgListenerV[View, PropertyName](listener)

		case func(PropertyName):
			view.changeListener[tag] = newOneArgListenerE[View](listener)

		case func(View, PropertyName):
			view.changeListener[tag] = newOneArgListenerVE(listener)

		case string:
			view.changeListener[tag] = newOneArgListenerBinding[View, PropertyName](listener)

		default:
			return false
		}

		view.setRaw(changeListeners, view.changeListener)
	}
	return true
}

func (view *viewData) HasFocus() bool {
	return view.hasFocus
}

func (view *viewData) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)
	writeViewStyle(view.tag, view, buffer, "", nil)
	return buffer.String()
}

func (view *viewData) excludeTags() []PropertyName {
	return nil
}
