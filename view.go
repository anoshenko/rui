package rui

import (
	"fmt"
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

	// SetChangeListener set the function to track the change of the View property
	SetChangeListener(tag PropertyName, listener func(View, PropertyName))

	// HasFocus returns 'true' if the view has focus
	HasFocus() bool

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
	exscludeTags() []PropertyName
	htmlDisabledProperty() bool

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
	changeListener   map[PropertyName]func(View, PropertyName)
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
	view.changeListener = map[PropertyName]func(View, PropertyName){}
	view.addCSS = map[string]string{}
	//view.animation = map[string]AnimationEndListener{}
	view.singleTransition = map[PropertyName]AnimationProperty{}
	view.noResizeEvent = false
	view.created = false
	view.hasHtmlDisabled = false
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
		}

		for _, tag := range changedTags {
			if listener, ok := view.changeListener[tag]; ok {
				listener(view, tag)
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
		}

		for _, tag := range changedTags {
			if listener, ok := view.changeListener[tag]; ok {
				listener(view, tag)
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
	if tag == ID {
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
			result = setOneArgEventListener[View, string](view, tag, value)
			if result != nil {
				if listeners, ok := view.getRaw(tag).([]func(View, string)); ok {
					newListeners := make([]func(View, PropertyName), len(listeners))
					for i, listener := range listeners {
						newListeners[i] = func(view View, name PropertyName) {
							listener(view, string(name))
						}
					}
					view.setRaw(tag, newListeners)
					return result
				}
				view.setRaw(tag, nil)
				return nil
			}
		}
		return result

	case AnimationStartEvent, AnimationEndEvent, AnimationIterationEvent, AnimationCancelEvent:
		return setOneArgEventListener[View, string](view, tag, value)

	case ResizeEvent, ScrollEvent:
		return setOneArgEventListener[View, Frame](view, tag, value)
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
		if value, ok := intProperty(view, TabIndex, view.Session(), 0); ok {
			session.updateProperty(view.htmlID(), "tabindex", strconv.Itoa(value))
		} else if view.Focusable() {
			session.updateProperty(view.htmlID(), "tabindex", "0")
		} else {
			session.updateProperty(view.htmlID(), "tabindex", "-1")
		}

	case Style, StyleDisabled:
		session.updateProperty(view.htmlID(), "class", view.htmlClass(IsDisabled(view)))

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
		session.updateCSSProperty(htmlID, "perspective-origin", transformOriginCSS(x, y, AutoSize(), view.Session()))

	case BackfaceVisible:
		if GetBackfaceVisible(view) {
			session.updateCSSProperty(htmlID, string(BackfaceVisible), "visible")
		} else {
			session.updateCSSProperty(htmlID, string(BackfaceVisible), "hidden")
		}

	case TransformOriginX, TransformOriginY, TransformOriginZ:
		x, y, z := getTransformOrigin(view, session)
		session.updateCSSProperty(htmlID, "transform-origin", transformOriginCSS(x, y, z, view.Session()))

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
		AnimationStartEvent, AnimationEndEvent, AnimationIterationEvent, AnimationCancelEvent:

		updateEventListenerHtml(view, tag)

	case DataList:
		updateInnerHTML(view.htmlID(), view.Session())

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
		buffer.WriteString(fmt.Sprintf(` data-left="%g" data-top="%g" data-width="%g" data-height="%g"`,
			view.frame.Left, view.frame.Top, view.frame.Width, view.frame.Height))
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
	//mouseEventsHtml(view, buffer, hasTooltip)

	viewEventsHtml[PointerEvent](view, []PropertyName{PointerDown, PointerUp, PointerMove, PointerOut, PointerOver, PointerCancel}, buffer)
	//pointerEventsHtml(view, buffer)

	viewEventsHtml[TouchEvent](view, []PropertyName{TouchStart, TouchEnd, TouchMove, TouchCancel}, buffer)
	//touchEventsHtml(view, buffer)

	viewEventsHtml[string](view, []PropertyName{TransitionRunEvent, TransitionStartEvent, TransitionEndEvent, TransitionCancelEvent,
		AnimationStartEvent, AnimationEndEvent, AnimationIterationEvent, AnimationCancelEvent}, buffer)
	//transitionEventsHtml(view, buffer)
	//animationEventsHtml(view, buffer)

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

	case FocusEvent:
		view.hasFocus = true
		for _, listener := range getNoArgEventListeners[View](view, nil, command) {
			listener(self)
		}

	case LostFocusEvent:
		view.hasFocus = false
		for _, listener := range getNoArgEventListeners[View](view, nil, command) {
			listener(self)
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
	default:
		return false
	}
	return true

}

func (view *viewData) SetChangeListener(tag PropertyName, listener func(View, PropertyName)) {
	if listener == nil {
		delete(view.changeListener, tag)
	} else {
		view.changeListener[tag] = listener
	}
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

func (view *viewData) exscludeTags() []PropertyName {
	return nil
}
