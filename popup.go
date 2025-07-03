package rui

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
)

// Constants for [Popup] specific properties and events
const (
	// Title is the constant for "title" property tag.
	//
	// Used by Popup, TabsLayout.
	//
	// Usage in Popup:
	// Define the title.
	//
	// Supported types: string.
	//
	// Usage in TabsLayout:
	// Set the title of the tab. The property is set for the child view of TabsLayout.
	//
	// Supported types: string.
	Title PropertyName = "title"

	// TitleStyle is the constant for "title-style" property tag.
	//
	// Used by Popup.
	// Set popup title style. Default title style is "ruiPopupTitle".
	//
	// Supported types: string.
	TitleStyle PropertyName = "title-style"

	// CloseButton is the constant for "close-button" property tag.
	//
	// Used by Popup.
	// Controls whether a close button can be added to the popup. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", "1" - Close button will be added to a title bar of a window.
	//   - false, 0, "false", "no", "off", "0" - Popup without a close button.
	CloseButton PropertyName = "close-button"

	// OutsideClose is the constant for "outside-close" property tag.
	//
	// Used by Popup.
	// Controls whether popup can be closed by clicking outside of the window. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", "1" - Clicking outside the popup window will automatically call the Dismiss() method.
	//   - false, 0, "false", "no", "off", "0" - Clicking outside the popup window has no effect.
	OutsideClose PropertyName = "outside-close"

	// Buttons is the constant for "buttons" property tag.
	//
	// Used by Popup.
	// Buttons that will be placed at the bottom of the popup.
	//
	// Supported types: PopupButton, []PopupButton.
	//
	// Internal type is []PopupButton, other types converted to it during assignment.
	// See PopupButton description for more details.
	Buttons PropertyName = "buttons"

	// ButtonsAlign is the constant for "buttons-align" property tag.
	//
	// Used by Popup.
	// Set the horizontal alignment of popup buttons.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (LeftAlign) or "left" - Left alignment.
	//   - 1 (RightAlign) or "right" - Right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	//   - 3 (StretchAlign) or "stretch" - Width alignment.
	ButtonsAlign PropertyName = "buttons-align"

	// DismissEvent is the constant for "dismiss-event" property tag.
	//
	// Used by Popup.
	// Used to track the closing state of the Popup. It occurs after the Popup disappears from the screen.
	//
	// General listener format:
	//
	//  func(popup rui.Popup)
	//
	// where:
	// popup - Interface of a popup which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	DismissEvent PropertyName = "dismiss-event"

	// Arrow is the constant for "arrow" property tag.
	//
	// Used by Popup.
	// Add an arrow to popup. Default value is "none".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NoneArrow) or "none" - No arrow.
	//   - 1 (TopArrow) or "top" - Arrow at the top side of the pop-up window.
	//   - 2 (RightArrow) or "right" - Arrow on the right side of the pop-up window.
	//   - 3 (BottomArrow) or "bottom" - Arrow at the bottom of the pop-up window.
	//   - 4 (LeftArrow) or "left" - Arrow on the left side of the pop-up window.
	Arrow PropertyName = "arrow"

	// ArrowAlign is the constant for "arrow-align" property tag.
	//
	// Used by Popup.
	// Set the horizontal alignment of the popup arrow. Default value is "center".
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (TopAlign/LeftAlign) or "top" - Top/left alignment.
	//   - 1 (BottomAlign/RightAlign) or "bottom" - Bottom/right alignment.
	//   - 2 (CenterAlign) or "center" - Center alignment.
	ArrowAlign PropertyName = "arrow-align"

	// ArrowSize is the constant for "arrow-size" property tag.
	//
	// Used by Popup.
	// Set the size(length) of the popup arrow. Default value is 16px defined by @ruiArrowSize constant.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	ArrowSize PropertyName = "arrow-size"

	// ArrowWidth is the constant for "arrow-width" property tag.
	//
	// Used by Popup.
	// Set the width of the popup arrow. Default value is 16px defined by @ruiArrowWidth constant.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	ArrowWidth PropertyName = "arrow-width"

	// ShowTransform is the constant for "show-transform" property tag.
	//
	// Used by Popup.
	// Specify start translation, scale and rotation over x, y and z axes as well as a distortion
	// for an animated Popup showing/hiding.
	//
	// Supported types: TransformProperty, string.
	//
	// See TransformProperty description for more details.
	//
	// Conversion rules:
	//   - TransformProperty - stored as is, no conversion performed.
	//   - string - string representation of Transform interface. Example:
	//
	//	"_{ translate-x = 10px, scale-y = 1.1}"
	ShowTransform PropertyName = "show-transform"

	// ShowDuration is the constant for "show-duration" property tag.
	//
	// Used by Popup.
	// Sets the length of time in seconds that a Popup show/hide animation takes to complete.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	ShowDuration PropertyName = "show-duration"

	// ShowTiming is the constant for "show-timing" property tag.
	//
	// Used by Popup.
	// Set how a Popup show/hide animation progresses through the duration of each cycle.
	//
	// Supported types: string.
	//
	// Values:
	//   - "ease" (EaseTiming) - Speed increases towards the middle and slows down at the end.
	//   - "ease-in" (EaseInTiming) - Speed is slow at first, but increases in the end.
	//   - "ease-out" (EaseOutTiming) - Speed is fast at first, but decreases in the end.
	//   - "ease-in-out" (EaseInOutTiming) - Speed is slow at first, but quickly increases and at the end it decreases again.
	//   - "linear" (LinearTiming) - Constant speed.
	//   - "step(n)" (StepTiming(n int) function) - Timing function along stepCount stops along the transition, displaying each stop for equal lengths of time.
	//   - "cubic-bezier(x1, y1, x2, y2)" (CubicBezierTiming(x1, y1, x2, y2 float64) function) - Cubic-Bezier curve timing function. x1 and x2 must be in the range [0, 1].
	ShowTiming PropertyName = "show-timing"

	// ShowOpacity is the constant for "show-opacity" property tag.
	//
	// Used by Popup.
	// In [1..0] range sets the start opacity of Popup show animation (the finish animation opacity is 1).
	// Opacity is the degree to which content behind the view is hidden, and is the opposite of transparency.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	ShowOpacity PropertyName = "show-opacity"

	// ArrowOffset is the constant for "arrow-offset" property tag.
	//
	// Used by Popup.
	// Set the offset of the popup arrow.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	ArrowOffset PropertyName = "arrow-offset"

	// NoneArrow is value of the popup "arrow" property: no arrow
	NoneArrow = 0

	// TopArrow is value of the popup "arrow" property:
	// Arrow at the top side of the pop-up window
	TopArrow = 1

	// RightArrow is value of the popup "arrow" property:
	// Arrow on the right side of the pop-up window
	RightArrow = 2

	// BottomArrow is value of the popup "arrow" property:
	// Arrow at the bottom of the pop-up window
	BottomArrow = 3

	// LeftArrow is value of the popup "arrow" property:
	// Arrow on the left side of the pop-up window
	LeftArrow = 4
)

// Constants which are used as a values of [PopupButtonType] variables
const (
	// NormalButton is the constant of the popup button type: the normal button
	NormalButton PopupButtonType = 0

	// DefaultButton is the constant of the popup button type: button that fires when the "Enter" key is pressed
	DefaultButton PopupButtonType = 1

	// CancelButton is the constant of the popup button type: button that fires when the "Escape" key is pressed
	CancelButton PopupButtonType = 2
)

const (
	popupLayerID   = "ruiPopupLayer"
	popupArrowID   = "ruiPopupArrow"
	popupButtonsID = "ruiPopupButtons"
	popupID        = "ruiPopup"
	popupContentID = "ruiPopupContent"
	popupTitleID   = "ruiPopupTitle"
)

// PopupButtonType represent popup button type
type PopupButtonType int

// PopupButton describes a button that will be placed at the bottom of the window.
type PopupButton struct {
	// Title of the button
	Title string

	// Type of the button
	Type PopupButtonType

	// OnClick is the handler function that gets called when the button is pressed
	OnClick func(Popup)
}

type popupButton struct {
	title      string
	buttonType PopupButtonType
	onClick    popupListener
}

// Popup represents a Popup view
type Popup interface {
	Properties
	fmt.Stringer

	// View returns a content view of the popup
	View() View

	// Session returns current client session
	Session() Session

	// Show displays a popup
	Show()

	// Dismiss closes a popup
	Dismiss()

	onDismiss()
	html(buffer *strings.Builder)
	viewByHTMLID(id string) View
	keyEvent(event KeyEvent) bool
	showAnimation()
	dismissAnimation(listener func(PropertyName)) bool
}

type popupListener interface {
	Run(Popup)
	rawListener() any
}

type popupListener0 struct {
	fn func()
}

type popupListener1 struct {
	fn func(Popup)
}

type popupListenerBinding struct {
	name    string
	dismiss bool
}

type popupData struct {
	propertyList
	session          Session
	layerView        GridLayout
	popupView        GridLayout
	contentContainer ColumnLayout
	contentView      View
}

type popupManager struct {
	popups []Popup
}

func (popup *popupData) createArrowView(location int) View {

	session := popup.Session()

	getSize := func(tag PropertyName, constTag string) SizeUnit {
		size, _ := sizeProperty(popup, tag, session)
		if size.Type != Auto && size.Value > 0 {
			return size
		}

		if value, ok := session.Constant(constTag); ok {
			if size, ok := StringToSizeUnit(value); ok && size.Type != Auto && size.Value > 0 {
				return size
			}
		}
		return Px(16)
	}

	size := getSize(ArrowSize, "ruiArrowSize")
	width := getSize(ArrowWidth, "ruiArrowWidth")

	params := Params{BackgroundColor: GetBackgroundColor(popup.popupView)}

	if shadow := GetShadowProperty(popup.popupView); shadow != nil {
		params[Shadow] = shadow
	}

	if filter := GetBackdropFilter(popup.popupView); filter != nil {
		params[BackdropFilter] = filter
	}

	switch location {
	case TopArrow:
		params[Row] = 0
		params[Column] = 1
		params[Clip] = NewPolygonClip([]any{"0%", "100%", "50%", "0%", "100%", "100%"})
		params[Width] = width
		params[Height] = size

	case RightArrow:
		params[Row] = 1
		params[Column] = 0
		params[Clip] = NewPolygonClip([]any{"0%", "0%", "100%", "50%", "0%", "100%"})
		params[Width] = size
		params[Height] = width

	case BottomArrow:
		params[Row] = 0
		params[Column] = 1
		params[Clip] = NewPolygonClip([]any{"0%", "0%", "50%", "100%", "100%", "0%"})
		params[Width] = width
		params[Height] = size

	case LeftArrow:
		params[Row] = 1
		params[Column] = 0
		params[Clip] = NewPolygonClip([]any{"100%", "0%", "0%", "50%", "100%", "100%"})
		params[Width] = size
		params[Height] = width
	}

	arrowView := NewView(session, params)

	params = Params{
		ID:      popupArrowID,
		Content: arrowView,
	}

	switch location {
	case TopArrow:
		params[Row] = 1
		params[Column] = 2

	case BottomArrow:
		params[Row] = 3
		params[Column] = 2

	case LeftArrow:
		params[Row] = 2
		params[Column] = 1

	case RightArrow:
		params[Row] = 2
		params[Column] = 3
	}

	off, _ := sizeProperty(popup, ArrowOffset, session)
	align, _ := enumProperty(popup, ArrowAlign, session, CenterAlign)

	if align != CenterAlign && off.Type == Auto {
		r := GetRadius(popup.popupView)
		switch location {
		case TopArrow:
			switch align {
			case LeftAlign:
				off = r.TopLeftX

			case RightAlign:
				off = r.TopRightX
			}

		case BottomArrow:
			switch align {
			case LeftAlign:
				off = r.BottomLeftX

			case RightAlign:
				off = r.BottomRightX
			}

		case RightArrow:
			switch align {
			case TopAlign:
				off = r.TopRightY

			case BottomAlign:
				off = r.BottomRightY
			}

		case LeftArrow:
			switch align {
			case TopAlign:
				off = r.TopLeftY

			case BottomAlign:
				off = r.BottomLeftY
			}
		}
	}

	switch location {
	case TopArrow, BottomArrow:
		cellWidth := make([]SizeUnit, 3)
		switch align {
		case LeftAlign:
			cellWidth[0] = off
			cellWidth[2] = Fr(1)

		case RightAlign:
			cellWidth[0] = Fr(1)
			cellWidth[2] = off

		default:
			cellWidth[0] = Fr(1)
			cellWidth[2] = Fr(1)
			if off.Type != Auto && off.Value != 0 {
				arrowView.Set(MarginLeft, off)
			}
		}
		params[CellWidth] = cellWidth

	case RightArrow, LeftArrow:
		cellHeight := make([]SizeUnit, 3)
		switch align {
		case TopAlign:
			cellHeight[0] = off
			cellHeight[2] = Fr(1)

		case BottomAlign:
			cellHeight[0] = Fr(1)
			cellHeight[2] = off

		default:
			cellHeight[0] = Fr(1)
			cellHeight[2] = Fr(1)
			if off.Type != Auto && off.Value != 0 {
				arrowView.Set(MarginTop, off)
			}
		}
		params[CellHeight] = cellHeight
	}

	return NewGridLayout(session, params)
}

func (popup *popupData) layerCellWidth() []SizeUnit {
	cellWidth := make([]SizeUnit, 5)
	switch hAlign, _ := enumProperty(popup, HorizontalAlign, popup.session, CenterAlign); hAlign {
	case LeftAlign:
		cellWidth[4] = Fr(1)

	case RightAlign:
		cellWidth[0] = Fr(1)

	default:
		cellWidth[0] = Fr(1)
		cellWidth[4] = Fr(1)
	}
	return cellWidth
}

func (popup *popupData) layerCellHeight() []SizeUnit {
	cellHeight := make([]SizeUnit, 5)
	switch vAlign, _ := enumProperty(popup, VerticalAlign, popup.session, CenterAlign); vAlign {
	case LeftAlign:
		cellHeight[4] = Fr(1)

	case RightAlign:
		cellHeight[0] = Fr(1)

	default:
		cellHeight[0] = Fr(1)
		cellHeight[4] = Fr(1)
	}

	return cellHeight
}

func (popup *popupData) Get(tag PropertyName) any {

	switch tag = defaultNormalize(tag); tag {
	case Content:
		return popup.contentView

	case "layer-view":
		if popup.layerView == nil {
			popup.layerView = popup.createLayerView()
		}
		return popup.layerView
	}

	return popup.properties[tag]
}

func (popup *popupData) arrowType() int {
	result, _ := enumProperty(popup, Arrow, popup.session, NoneArrow)
	return result
}

func (popup *popupData) supported(tag PropertyName) bool {
	switch tag {
	case Row, Column, CellWidth, CellHeight, Gap, GridColumnGap, GridRowGap,
		CellVerticalAlign, CellHorizontalAlign,
		CellVerticalSelfAlign, CellHorizontalSelfAlign:
		return false
	}
	return true
}

func (popup *popupData) Set(tag PropertyName, value any) bool {

	if value == nil {
		popup.Remove(tag)
		return true
	}

	switch tag = defaultNormalize(tag); tag {
	case Buttons:
		return popup.setButtons(value)

	case Title:
		switch value := value.(type) {
		case string:
			popup.setRaw(Title, value)

		case View:
			popup.setRaw(Title, value)

		default:
			notCompatibleType(Title, value)
			return false
		}
		popup.propertyChanged(Title)
		return true

	case Content:
		switch value := value.(type) {
		case View:
			popup.contentView = value
			popup.setRaw(Content, popup.contentView)

		case DataObject:
			view := CreateViewFromObject(popup.session, value, nil)
			if view == nil {
				return false
			}
			popup.contentView = view
			popup.setRaw(Content, popup.contentView)

		case string:
			if len(value) > 0 && value[0] == '@' {
				view := CreateViewFromResources(popup.session, value[1:])
				if view != nil {
					popup.contentView = view
					break
				}
			}
			popup.contentView = NewTextView(popup.session, Params{Text: value})
			popup.setRaw(Content, value)

		default:
			notCompatibleType(Buttons, value)
			return false
		}

		if binding := popup.getRaw(Binding); binding != nil {
			popup.contentView.Set(Binding, binding)
		}

		popup.propertyChanged(Content)
		return true

	case Binding:
		if popup.contentView != nil {
			popup.contentView.Set(Binding, value)
		}
		popup.setRaw(Binding, value)
		popup.propertyChanged(Binding)
		return true

	case DismissEvent:
		if listeners, ok := valueToPopupEventListeners(value); ok {
			if listeners != nil {
				popup.setRaw(DismissEvent, listeners)
				popup.propertyChanged(DismissEvent)
				return true
			}
		}
		notCompatibleType(tag, value)
		return false

	case ShowTransform:
		return setTransformProperty(popup, tag, value)
	}

	if popup.supported(tag) {
		tags := viewStyleSet(popup, tag, value)
		if len(tags) > 0 {
			for _, tag := range tags {
				popup.propertyChanged(tag)
			}
			return true
		}
	} else {
		ErrorLogF(`"%s" property is not supported by the popup.`, string(tag))
	}

	return false
}

func (popup *popupData) Remove(tag PropertyName) {
	tag = defaultNormalize(tag)

	switch tag {
	case Content:
		popup.contentView = nil
	}

	if popup.supported(tag) {
		tags := viewStyleRemove(popup, tag)
		for _, tag := range tags {
			popup.propertyChanged(tag)
		}
	} else {
		ErrorLogF(`"%s" property is not supported by the popup.`, string(tag))
	}
}

func (popup *popupData) propertyChanged(tag PropertyName) {
	if popup.layerView == nil {
		return
	}

	switch tag {
	case Content:
		popup.contentContainer.Set(Content, popup.contentView)

	case Arrow, ArrowWidth, ArrowSize, ArrowAlign, ArrowOffset:
		popup.layerView.RemoveViewByID(popupArrowID)
		if location := popup.arrowType(); location != NoneArrow {
			popup.layerView.Append(popup.createArrowView(location))
		}

	case Buttons:
		popup.popupView.RemoveViewByID(popupButtonsID)
		if view := popup.createButtonsPanel(); view != nil {
			popup.popupView.Append(view)
		}

	case Margin:
		if margin, ok := getBounds(popup, Margin, popup.session); ok {
			popup.layerView.Set(Padding, margin)
		} else {
			popup.layerView.Remove(Padding)
		}

	case Title, CloseButton:
		popup.popupView.RemoveViewByID(popupTitleID)
		if view := popup.createTitleView(); view != nil {
			popup.popupView.Append(view)
		}

	case TitleStyle:
		if title := ViewByID(popup.popupView, popupTitleID); title != nil {
			titleStyle := "ruiPopupTitle"
			if style, ok := stringProperty(popup, TitleStyle, popup.session); ok {
				titleStyle = style
			}
			title.Set(Style, titleStyle)
		}

	case OutsideClose:
		outsideClose, _ := boolProperty(popup, OutsideClose, popup.session)
		if outsideClose {
			popup.layerView.Set(ClickEvent, popup.cancel)
		} else {
			popup.layerView.Set(ClickEvent, func() {})
		}

	case ShowDuration, ShowTiming:
		animation := popup.animationProperty()
		opacity, _ := floatProperty(popup, ShowOpacity, popup.session, 1)
		if opacity != 1 {
			popup.popupView.SetTransition(Opacity, animation)
		}
		transform := getTransformProperty(popup, ShowTransform)
		if transform != nil {
			popup.popupView.SetTransition(Transform, animation)
		}

	case ShowOpacity:
		opacity, _ := floatProperty(popup, ShowOpacity, popup.session, 1)
		if opacity != 1 {
			popup.popupView.SetTransition(Opacity, popup.animationProperty())
		}

	case ShowTransform:
		transform := getTransformProperty(popup, ShowTransform)
		if transform != nil {
			popup.popupView.SetTransition(Transform, popup.animationProperty())
		}

	default:
		if popup.supported(tag) {
			popup.popupView.Set(tag, popup.getRaw(tag))
		}
	}
}

func (popup *popupData) animationProperty() AnimationProperty {
	duration, _ := floatProperty(popup, ShowDuration, popup.session, 1)
	timing, ok := stringProperty(popup, ShowTiming, popup.session)
	if !ok {
		timing = EaseTiming
	}
	return NewAnimationProperty(Params{
		Duration:       duration,
		TimingFunction: timing,
	})
}

func (popup *popupData) View() View {
	return popup.contentView
}

func (popup *popupData) Session() Session {
	return popup.session
}

func (popup *popupData) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)
	writeViewStyle("Popup", popup, buffer, "", nil)
	return buffer.String()
}

func (popup *popupData) setButtons(value any) bool {
	popupButtonFromObject := func(obj DataObject) popupButton {
		var button popupButton

		button.title, _ = obj.PropertyValue(string(Title))

		if text, ok := obj.PropertyValue("type"); ok {
			text, _ = popup.session.resolveConstants(text)
			t, _ := enumStringToInt(text, []string{"normal", "default", "cancel"}, true)
			button.buttonType = PopupButtonType(t)
		}

		if fn, ok := obj.PropertyValue("click"); ok {
			button.onClick = newPopupListenerBinding(fn, true)
		} else if button.buttonType == CancelButton {
			button.onClick = newPopupListener0(popup.Dismiss)
		}

		return button
	}

	buttonConvert := func(button PopupButton) popupButton {
		result := popupButton{
			title:      button.Title,
			buttonType: button.Type,
		}
		if button.OnClick != nil {
			result.onClick = newPopupListener1(button.OnClick)
		}
		return result
	}

	switch value := value.(type) {
	case PopupButton:
		popup.setRaw(Buttons, []popupButton{buttonConvert(value)})

	case []PopupButton:
		buttons := make([]popupButton, 0, len(value))
		for _, button := range value {
			buttons = append(buttons, buttonConvert(button))
		}
		popup.setRaw(Buttons, buttons)

	case []popupButton:
		popup.setRaw(Buttons, value)

	case DataObject:
		popup.setRaw(Buttons, []popupButton{popupButtonFromObject(value)})

	case []DataValue:
		buttons := make([]popupButton, 0, len(value))
		for _, val := range value {
			if val.IsObject() {
				buttons = append(buttons, popupButtonFromObject(val.Object()))
			} else {
				notCompatibleType(Buttons, val)
			}
		}
		if len(buttons) > 0 {
			popup.setRaw(Buttons, buttons)
		}

	case []any:
		buttons := make([]popupButton, 0, len(value))
		for _, val := range value {
			switch val := val.(type) {
			case DataObject:
				buttons = append(buttons, popupButtonFromObject(val))

			case popupButton:
				buttons = append(buttons, val)

			case PopupButton:
				buttons = append(buttons, popupButton{
					title:      val.Title,
					buttonType: val.Type,
					onClick:    newPopupListener1(val.OnClick),
				})

			default:
				notCompatibleType(Buttons, val)
			}
		}
		if len(buttons) > 0 {
			popup.setRaw(Buttons, buttons)
		}

	default:
		notCompatibleType(Buttons, value)
		return false
	}

	popup.propertyChanged(Buttons)
	return true
}

func (popup *popupData) buttons() []popupButton {
	if value := popup.getRaw(Buttons); value != nil {
		if result, ok := value.([]popupButton); ok {
			return result
		}
	}
	return nil
}

func (popup *popupData) cancel() {
	if buttons := popup.buttons(); buttons != nil {
		for _, button := range buttons {
			if button.buttonType == CancelButton && button.onClick != nil {
				button.onClick.Run(popup)
				return
			}
		}
	}
	popup.Dismiss()
}

func (popup *popupData) Dismiss() {
	popup.Session().popupManager().dismissPopup(popup)
}

func (popup *popupData) Show() {
	popup.Session().popupManager().showPopup(popup)
}

func (popup *popupData) showAnimation() {
	opacity, _ := floatProperty(popup, ShowOpacity, popup.session, 1)
	transform := getTransformProperty(popup, ShowTransform)

	if opacity != 1 || transform != nil {
		htmlID := popup.popupView.htmlID()
		session := popup.Session()
		if opacity != 1 {
			session.updateCSSProperty(htmlID, string(Opacity), "1")
		}
		if transform != nil {
			session.updateCSSProperty(htmlID, string(Transform), "")
		}
	}
}

func (popup *popupData) dismissAnimation(listener func(PropertyName)) bool {
	opacity, _ := floatProperty(popup, ShowOpacity, popup.session, 1)
	transform := getTransformProperty(popup, ShowTransform)

	if opacity != 1 || transform != nil {
		session := popup.Session()
		popup.popupView.Set(TransitionEndEvent, listener)
		popup.popupView.Set(TransitionCancelEvent, listener)

		htmlID := popup.popupView.htmlID()
		if opacity != 1 {
			session.updateCSSProperty(htmlID, string(Opacity), fmt.Sprintf("%.2f", opacity))
		}
		if transform != nil {
			session.updateCSSProperty(htmlID, string(Transform), transform.transformCSS(session))
		}
		return true
	}
	return false
}

func (popup *popupData) html(buffer *strings.Builder) {
	if popup.layerView == nil {
		popup.layerView = popup.createLayerView()
	}
	viewHTML(popup.layerView, buffer, "")
}

func (popup *popupData) viewByHTMLID(id string) View {
	if popup.layerView != nil {
		return viewByHTMLID(id, popup.layerView)
	}
	return nil
}

func (popup *popupData) onDismiss() {
	if popup.layerView != nil {
		popup.Session().callFunc("removeView", popup.layerView.htmlID())

		if value := popup.getRaw(DismissEvent); value != nil {
			if listeners, ok := value.([]popupListener); ok {
				for _, listener := range listeners {
					listener.Run(popup)
				}
			}
		}
	}
}

func (popup *popupData) keyEvent(event KeyEvent) bool {
	if !event.AltKey && !event.CtrlKey && !event.ShiftKey && !event.MetaKey {
		switch event.Code {
		case EnterKey:
			for _, button := range popup.buttons() {
				if button.buttonType == DefaultButton && button.onClick != nil {
					button.onClick.Run(popup)
					return true
				}
			}

		case EscapeKey:
			cancelable := func() bool {
				if closeButton, _ := boolProperty(popup, CloseButton, popup.session); closeButton {
					return true
				}
				if outsideClose, _ := boolProperty(popup, OutsideClose, popup.session); outsideClose {
					return true
				}

				for _, button := range popup.buttons() {
					if button.buttonType == CancelButton {
						return true
					}
				}
				return false
			}

			if cancelable() {
				popup.cancel()
				return true
			}
		}
	}
	return false
}

func (popup *popupData) createButtonsPanel() GridLayout {
	buttons := popup.buttons()
	if buttonCount := len(buttons); buttonCount > 0 {
		session := popup.session
		buttonsAlign, _ := enumProperty(popup, ButtonsAlign, session, RightAlign)
		gap, _ := sizeConstant(session, "ruiPopupButtonGap")
		cellWidth := []SizeUnit{}
		for range buttonCount {
			cellWidth = append(cellWidth, Fr(1))
		}

		buttonsPanel := NewGridLayout(session, Params{
			CellWidth: cellWidth,
		})
		if gap.Type != Auto && gap.Value > 0 {
			buttonsPanel.Set(Gap, gap)
			buttonsPanel.Set(Margin, gap)
		}

		for i, button := range buttons {
			title := button.title
			if title == "" && button.buttonType == CancelButton {
				title = "Cancel"
			}

			buttonView := NewButton(session, Params{
				Column:  i,
				Content: title,
			})

			if button.onClick != nil {
				fn := button.onClick.Run
				buttonView.Set(ClickEvent, func() {
					fn(popup)
				})
			} else if button.buttonType == CancelButton {
				buttonView.Set(ClickEvent, popup.cancel)
			}

			if button.buttonType == DefaultButton {
				buttonView.Set(Style, "ruiDefaultButton")
			}

			buttonsPanel.Append(buttonView)
		}

		return NewGridLayout(session, Params{
			ID:                  popupButtonsID,
			Column:              0,
			Row:                 2,
			CellHorizontalAlign: buttonsAlign,
			Content:             buttonsPanel,
		})
	}
	return nil
}

func (popup *popupData) createTitleView() GridLayout {
	session := popup.Session()

	var closeButton View = nil
	if hasButton, _ := boolProperty(popup, CloseButton, popup.session); hasButton {
		closeButton = NewGridLayout(session, Params{
			Column:              1,
			Height:              "@ruiPopupTitleHeight",
			Width:               "@ruiPopupTitleHeight",
			CellHorizontalAlign: CenterAlign,
			CellVerticalAlign:   CenterAlign,
			TextSize:            Px(20),
			Content:             "✕",
			NotTranslate:        true,
			ClickEvent:          popup.cancel,
		})
	}

	var title View = nil
	if value := popup.getRaw(Title); value != nil {
		switch value := value.(type) {
		case string:
			if len(value) > 0 && value[0] == '@' {
				title = CreateViewFromResources(session, value[1:])
				if title != nil {
					break
				}
			}
			title = NewTextView(session, Params{Text: value})

		case View:
			title = value
		}
	}

	if title == nil && closeButton == nil {
		return nil
	}

	titleStyle := "ruiPopupTitle"
	if style, ok := stringProperty(popup, TitleStyle, session); ok {
		titleStyle = style
	}

	titleContent := []View{}
	if title != nil {
		titleContent = append(titleContent, title)
	}
	if closeButton != nil {
		titleContent = append(titleContent, closeButton)
	}

	return NewGridLayout(session, Params{
		ID:                popupTitleID,
		Row:               0,
		Column:            0,
		Style:             titleStyle,
		CellWidth:         []any{Fr(1), AutoSize()},
		CellVerticalAlign: CenterAlign,
		PaddingLeft:       Px(12),
		Content:           titleContent,
	})
}

func (popup *popupData) createContentContainer() ColumnLayout {
	params := Params{
		ID:     popupContentID,
		Column: 0,
		Row:    1,
	}

	if popup.contentView != nil {
		params[Content] = popup.contentView
	}

	popup.contentContainer = NewColumnLayout(popup.session, params)
	return popup.contentContainer
}

func (popup *popupData) createLayerView() GridLayout {

	session := popup.session

	params := Params{
		Style:               "ruiPopup",
		ID:                  popupID,
		Row:                 2,
		Column:              2,
		MaxWidth:            Percent(100),
		MaxHeight:           Percent(100),
		CellVerticalAlign:   StretchAlign,
		CellHorizontalAlign: StretchAlign,
		CellHeight:          []SizeUnit{AutoSize(), Fr(1), AutoSize()},
		ClickEvent:          func(View) {},
		Shadow: NewShadowProperty(Params{
			SpreadRadius: Px(4),
			Blur:         Px(16),
			ColorTag:     "@ruiPopupShadow",
		}),
	}

	popupProperties := []PropertyName{
		Content,
		Title,
		TitleStyle,
		CloseButton,
		OutsideClose,
		Buttons,
		ButtonsAlign,
		DismissEvent,
		Arrow,
		ArrowAlign,
		ArrowSize,
		ArrowWidth,
		ArrowOffset,
		ShowTransform,
		ShowDuration,
		ShowTiming,
		ShowOpacity,
		VerticalAlign,
		HorizontalAlign,
		Margin,
		Row,
		Column,
		CellWidth,
		CellHeight,
		CellVerticalAlign,
		CellHorizontalAlign,
	}

	for tag, value := range popup.properties {
		if !slices.Contains(popupProperties, tag) {
			params[tag] = value
		}
	}

	views := make([]View, 0, 3)
	if title := popup.createTitleView(); title != nil {
		views = append(views, title)
	}

	views = append(views, popup.createContentContainer())

	if buttons := popup.createButtonsPanel(); buttons != nil {
		views = append(views, buttons)
	}

	params[Content] = views

	popup.popupView = NewGridLayout(session, params)

	layerParams := Params{
		Style:      popupLayerID,
		MaxWidth:   Percent(100),
		MaxHeight:  Percent(100),
		CellWidth:  popup.layerCellWidth(),
		CellHeight: popup.layerCellHeight(),
	}

	if margin, ok := getBounds(popup, Margin, session); ok {
		layerParams[Padding] = margin
	}

	if location := popup.arrowType(); location != NoneArrow {
		layerParams[Content] = []View{popup.popupView, popup.createArrowView(location)}
	} else {
		layerParams[Content] = []View{popup.popupView}
	}

	if outsideClose, _ := boolProperty(popup, OutsideClose, session); outsideClose {
		layerParams[ClickEvent] = popup.cancel
	}

	popup.layerView = NewGridLayout(session, layerParams)

	opacity, _ := floatProperty(popup, ShowOpacity, session, 1)
	transform := getTransformProperty(popup, ShowTransform)
	if opacity != 1 || transform != nil {
		animation := popup.animationProperty()
		if opacity != 1 {
			popup.popupView.Set(Opacity, opacity)
			popup.popupView.SetTransition(Opacity, animation)
		}
		if transform != nil {
			popup.popupView.Set(Transform, transform)
			popup.popupView.SetTransition(Transform, animation)
		}
	}

	return popup.layerView
}

// NewPopup creates a new Popup
func NewPopup(view View, param Params) Popup {
	if view == nil {
		return nil
	}

	popup := new(popupData)
	popup.session = view.Session()
	popup.contentView = view
	popup.properties = map[PropertyName]any{}
	for tag, value := range param {
		popup.Set(tag, value)
	}
	return popup
}

// CreatePopupFromObject create new Popup and initialize it by content of object. Parameters:
//   - session - the session to which the view will be attached (should not be nil);
//   - text - text describing Popup;
//   - binding - object assigned to the Binding property (optional parameter).
//
// If the function fails, it returns nil and an error message is written to the log.
func CreatePopupFromObject(session Session, object DataObject, binding any) Popup {
	if session == nil {
		ErrorLog(`"session" argument is nil`)
		return nil
	}

	if object == nil {
		ErrorLog(`"object" argument is nil`)
		return nil
	}

	if strings.ToLower(object.Tag()) != "popup" {
		ErrorLog(`the object name must be "Popup"`)
		return nil
	}

	popup := new(popupData)
	popup.session = session
	popup.properties = map[PropertyName]any{}

	for key, value := range object.ToParams() {
		popup.Set(key, value)
	}

	if binding != nil {
		popup.Set(Binding, binding)
	}

	return popup
}

// CreatePopupFromText create new Popup and initialize it by content of text. Parameters:
//   - session - the session to which the view will be attached (should not be nil);
//   - text - text describing Popup;
//   - binding - object assigned to the Binding property (optional parameter).
//
// If the function fails, it returns nil and an error message is written to the log.
func CreatePopupFromText(session Session, text string, binding any) Popup {
	data, err := ParseDataText(text)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	return CreatePopupFromObject(session, data, binding)
}

// CreatePopupFromResources create new Popup and initialize it by the content of
// the resource file from "popups" directory. Parameters:
//   - session - the session to which the view will be attached (should not be nil);
//   - text - file name in the "popups" folder of the application resources (it is not necessary to specify the .rui extension, it is added automatically);
//   - binding - object assigned to the Binding property (optional parameter).
//
// If the function fails, it returns nil and an error message is written to the log.
func CreatePopupFromResources(session Session, name string, binding any) Popup {
	if strings.ToLower(filepath.Ext(name)) != ".rui" {
		name += ".rui"
	}

	createEmbed := func(fs *embed.FS, path string) Popup {
		if data, err := fs.ReadFile(path); err == nil {
			data, err := ParseDataText(string(data))
			if err == nil {
				return CreatePopupFromObject(session, data, binding)
			}
			ErrorLog(err.Error())
		}
		return nil
	}

	for _, fs := range resources.embedFS {
		rootDirs := resources.embedRootDirs(fs)
		for _, dir := range rootDirs {
			switch dir {
			case imageDir, themeDir, rawDir:
				// do nothing

			case viewDir:
				if result := createEmbed(fs, dir+"/"+name); result != nil {
					return result
				}

			case popupDir:
				if result := createEmbed(fs, dir+"/"+name); result != nil {
					return result
				}

			default:
				if result := createEmbed(fs, dir+"/"+popupDir+"/"+name); result != nil {
					return result
				}
				if result := createEmbed(fs, dir+"/"+viewDir+"/"+name); result != nil {
					return result
				}
			}
		}
	}

	if resources.path == "" {
		return nil
	}

	createFromFile := func(path string) Popup {
		if data, err := os.ReadFile(path); err == nil {
			data, err := ParseDataText(string(data))
			if err == nil {
				return CreatePopupFromObject(session, data, binding)
			}
			ErrorLog(err.Error())
		}
		return nil
	}

	if result := createFromFile(resources.path + popupDir + "/" + name); result != nil {
		return result
	}

	return createFromFile(resources.path + viewDir + "/" + name)
}

// ShowPopup creates a new Popup and shows it
func ShowPopup(view View, param Params) Popup {
	popup := NewPopup(view, param)
	if popup != nil {
		popup.Show()
	}
	return popup
}

func (manager *popupManager) updatePopupLayerInnerHTML(session Session) {
	if manager.popups == nil {
		manager.popups = []Popup{}
		session.updateInnerHTML(popupLayerID, "")
		return
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	for _, popup := range manager.popups {
		popup.html(buffer)
	}
	session.updateInnerHTML(popupLayerID, buffer.String())
}

func (manager *popupManager) showPopup(popup Popup) {
	if popup == nil {
		return
	}

	session := popup.Session()
	if len(manager.popups) == 0 {
		manager.popups = []Popup{popup}
	} else {
		manager.popups = append(manager.popups, popup)
	}

	session.callFunc("blurCurrent")
	manager.updatePopupLayerInnerHTML(session)
	session.updateCSSProperty("ruiTooltipLayer", "visibility", "hidden")
	session.updateCSSProperty("ruiTooltipLayer", "opacity", "0")
	session.updateCSSProperty(popupLayerID, "visibility", "visible")
	session.updateCSSProperty("ruiRoot", "pointer-events", "none")
	popup.showAnimation()
}

func (manager *popupManager) dismissPopup(popup Popup) {
	if manager.popups == nil {
		manager.popups = []Popup{}
		return
	}

	count := len(manager.popups)
	if count <= 0 || popup == nil {
		return
	}

	index := -1
	for n, p := range manager.popups {
		if p == popup {
			index = n
			break
		}
	}

	if index < 0 {
		return
	}

	session := popup.Session()
	listener := func(PropertyName) {
		switch index {
		case 0:
			if count == 1 {
				manager.popups = []Popup{}
				session.updateCSSProperty("ruiRoot", "pointer-events", "auto")
				session.updateCSSProperty(popupLayerID, "visibility", "hidden")
			} else {
				manager.popups = manager.popups[1:]
			}

		case count - 1:
			manager.popups = manager.popups[:count-1]

		default:
			manager.popups = append(manager.popups[:index], manager.popups[index+1:]...)
		}
		popup.onDismiss()
	}

	if !popup.dismissAnimation(listener) {
		listener("")
	}
}

func newPopupListener0(fn func()) popupListener {
	obj := new(popupListener0)
	obj.fn = fn
	return obj
}

func (data *popupListener0) Run(_ Popup) {
	data.fn()
}

func (data *popupListener0) rawListener() any {
	return data.fn
}

func newPopupListener1(fn func(Popup)) popupListener {
	obj := new(popupListener1)
	obj.fn = fn
	return obj
}

func (data *popupListener1) Run(popup Popup) {
	data.fn(popup)
}

func (data *popupListener1) rawListener() any {
	return data.fn
}

func newPopupListenerBinding(name string, dismiss bool) popupListener {
	obj := new(popupListenerBinding)
	obj.name = name
	obj.dismiss = dismiss
	return obj
}

func (data *popupListenerBinding) runDismiss(popup Popup) bool {
	if strings.ToLower(data.name) == "dismiss" {
		popup.Dismiss()
		return true
	}
	return false
}

func (data *popupListenerBinding) Run(popup Popup) {
	bind := popup.View().binding()
	if bind == nil {
		if !data.dismiss || !data.runDismiss(popup) {
			ErrorLogF(`There is no a binding object for call "%s"`, data.name)
		}
		return
	}

	val := reflect.ValueOf(bind)
	method := val.MethodByName(data.name)
	if !method.IsValid() {
		if !data.dismiss || !data.runDismiss(popup) {
			ErrorLogF(`The "%s" method is not valid`, data.name)
		}
		return
	}

	methodType := method.Type()
	var args []reflect.Value = nil
	switch methodType.NumIn() {
	case 0:
		args = []reflect.Value{}

	case 1:
		inType := methodType.In(0)
		if inType == reflect.TypeOf(popup) {
			args = []reflect.Value{reflect.ValueOf(popup)}
		}
	}

	if args != nil {
		method.Call(args)
	} else {
		ErrorLogF(`Unsupported prototype of "%s" method`, data.name)
	}
}

func (data *popupListenerBinding) rawListener() any {
	return data.name
}

func valueToPopupEventListeners(value any) ([]popupListener, bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case []popupListener:
		return value, true

	case popupListener:
		return []popupListener{value}, true

	case string:
		return []popupListener{newPopupListenerBinding(value, false)}, true

	case func(Popup):
		return []popupListener{newPopupListener1(value)}, true

	case func():
		return []popupListener{newPopupListener0(value)}, true

	case []func(Popup):
		result := make([]popupListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newPopupListener1(fn))
			}
		}
		return result, len(result) > 0

	case []func():
		result := make([]popupListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newPopupListener0(fn))
			}
		}
		return result, len(result) > 0

	case []any:
		result := make([]popupListener, 0, len(value))
		for _, v := range value {
			if v != nil {
				switch v := v.(type) {
				case func(Popup):
					result = append(result, newPopupListener1(v))

				case func():
					result = append(result, newPopupListener0(v))

				case string:
					result = append(result, newPopupListenerBinding(v, false))

				default:
					return nil, false
				}
			}
		}
		return result, len(result) > 0
	}

	return nil, false
}

func getPopupListenerBinding(listeners []popupListener) string {
	for _, listener := range listeners {
		raw := listener.rawListener()
		if text, ok := raw.(string); ok && text != "" {
			return text
		}
	}
	return ""
}
