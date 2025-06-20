package rui

import (
	"fmt"
	"reflect"
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
	// for an animated Popup showing/hidding.
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

// Popup represents a Popup view
type Popup interface {
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
	dissmissAnimation(listener func(PropertyName)) bool
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
	name string
}

type popupData struct {
	layerView       GridLayout
	popupView       GridLayout
	contentView     View
	buttons         []PopupButton
	cancelable      bool
	dismissListener []popupListener
	showTransform   TransformProperty
	showOpacity     float64
	showDuration    float64
	showTiming      string
}

type popupManager struct {
	popups []Popup
}

type popupArrow struct {
	column, row      int
	location, align  int
	size, width, off SizeUnit
}

func (arrow *popupArrow) fixOff(popupView View) {
	if arrow.align == CenterAlign && arrow.off.Type == Auto {
		r := GetRadius(popupView)
		switch arrow.location {
		case TopArrow:
			switch arrow.align {
			case LeftAlign:
				arrow.off = r.TopLeftX

			case RightAlign:
				arrow.off = r.TopRightX
			}

		case BottomArrow:
			switch arrow.align {
			case LeftAlign:
				arrow.off = r.BottomLeftX

			case RightAlign:
				arrow.off = r.BottomRightX
			}

		case RightArrow:
			switch arrow.align {
			case TopAlign:
				arrow.off = r.TopRightY

			case BottomAlign:
				arrow.off = r.BottomRightY
			}

		case LeftArrow:
			switch arrow.align {
			case TopAlign:
				arrow.off = r.TopLeftY

			case BottomAlign:
				arrow.off = r.BottomLeftY
			}
		}
	}
}

func (arrow *popupArrow) createView(popupView View) View {
	session := popupView.Session()

	defaultSize := func(constTag string, defValue SizeUnit) SizeUnit {
		if value, ok := session.Constant(constTag); ok {
			if size, ok := StringToSizeUnit(value); ok && size.Type != Auto && size.Value != 0 {
				return size
			}
		}
		return defValue
	}

	if arrow.size.Type == Auto || arrow.size.Value == 0 {
		arrow.size = defaultSize("ruiArrowSize", Px(16))
	}

	if arrow.width.Type == Auto || arrow.width.Value == 0 {
		arrow.width = defaultSize("ruiArrowWidth", Px(16))
	}

	params := Params{BackgroundColor: GetBackgroundColor(popupView)}

	if shadow := GetShadowProperty(popupView); shadow != nil {
		params[Shadow] = shadow
	}

	if filter := GetBackdropFilter(popupView); filter != nil {
		params[BackdropFilter] = filter
	}

	switch arrow.location {
	case TopArrow:
		params[Row] = 0
		params[Column] = 1
		params[Clip] = NewPolygonClip([]any{"0%", "100%", "50%", "0%", "100%", "100%"})
		params[Width] = arrow.width
		params[Height] = arrow.size

	case RightArrow:
		params[Row] = 1
		params[Column] = 0
		params[Clip] = NewPolygonClip([]any{"0%", "0%", "100%", "50%", "0%", "100%"})
		params[Width] = arrow.size
		params[Height] = arrow.width

	case BottomArrow:
		params[Row] = 0
		params[Column] = 1
		params[Clip] = NewPolygonClip([]any{"0%", "0%", "50%", "100%", "100%", "0%"})
		params[Width] = arrow.width
		params[Height] = arrow.size

	case LeftArrow:
		params[Row] = 1
		params[Column] = 0
		params[Clip] = NewPolygonClip([]any{"100%", "0%", "0%", "50%", "100%", "100%"})
		params[Width] = arrow.size
		params[Height] = arrow.width
	}

	arrowView := NewView(session, params)

	params = Params{
		Row:     arrow.row,
		Column:  arrow.column,
		Content: arrowView,
	}

	arrow.fixOff(popupView)

	switch arrow.location {
	case TopArrow, BottomArrow:
		cellWidth := make([]SizeUnit, 3)
		switch arrow.align {
		case LeftAlign:
			cellWidth[0] = arrow.off
			cellWidth[2] = Fr(1)

		case RightAlign:
			cellWidth[0] = Fr(1)
			cellWidth[2] = arrow.off

		default:
			cellWidth[0] = Fr(1)
			cellWidth[2] = Fr(1)
			if arrow.off.Type != Auto && arrow.off.Value != 0 {
				arrowView.Set(MarginLeft, arrow.off)
			}
		}
		params[CellWidth] = cellWidth

	case RightArrow, LeftArrow:
		cellHeight := make([]SizeUnit, 3)
		switch arrow.align {
		case TopAlign:
			cellHeight[0] = arrow.off
			cellHeight[2] = Fr(1)

		case BottomAlign:
			cellHeight[0] = Fr(1)
			cellHeight[2] = arrow.off

		default:
			cellHeight[0] = Fr(1)
			cellHeight[2] = Fr(1)
			if arrow.off.Type != Auto && arrow.off.Value != 0 {
				arrowView.Set(MarginTop, arrow.off)
			}
		}
		params[CellHeight] = cellHeight
	}

	return NewGridLayout(session, params)
}

func (popup *popupData) layerCellWidth(arrowLocation int, popupParams Params, session Session) []SizeUnit {

	var columnCount int
	switch arrowLocation {
	case LeftArrow, RightArrow:
		columnCount = 4

	default:
		columnCount = 3
	}

	cellWidth := make([]SizeUnit, columnCount)
	switch hAlign, _ := enumProperty(popupParams, HorizontalAlign, session, CenterAlign); hAlign {
	case LeftAlign:
		cellWidth[columnCount-1] = Fr(1)

	case RightAlign:
		cellWidth[0] = Fr(1)

	default:
		cellWidth[0] = Fr(1)
		cellWidth[columnCount-1] = Fr(1)
	}
	return cellWidth
}

func (popup *popupData) layerCellHeight(arrowLocation int, popupParams Params, session Session) []SizeUnit {

	var rowCount int
	switch arrowLocation {
	case TopArrow, BottomArrow:
		rowCount = 4

	default:
		rowCount = 3
	}

	cellHeight := make([]SizeUnit, rowCount)
	switch vAlign, _ := enumProperty(popupParams, VerticalAlign, session, CenterAlign); vAlign {
	case LeftAlign:
		cellHeight[rowCount-1] = Fr(1)

	case RightAlign:
		cellHeight[0] = Fr(1)

	default:
		cellHeight[0] = Fr(1)
		cellHeight[rowCount-1] = Fr(1)
	}

	return cellHeight
}

func (popup *popupData) init(view View, popupParams Params) {
	popup.contentView = view
	popup.cancelable = false
	session := view.Session()

	popupRow := 1
	popupColumn := 1
	arrow := popupArrow{
		row:    1,
		column: 1,
		align:  CenterAlign,
	}

	switch arrow.location, _ = enumProperty(popupParams, Arrow, session, NoneArrow); arrow.location {
	case TopArrow:
		popupRow = 2

	case BottomArrow:
		arrow.row = 2

	case LeftArrow:
		popupColumn = 2

	case RightArrow:
		arrow.column = 2
	}

	layerParams := Params{
		Style:      "ruiPopupLayer",
		MaxWidth:   Percent(100),
		MaxHeight:  Percent(100),
		CellWidth:  popup.layerCellWidth(arrow.location, popupParams, session),
		CellHeight: popup.layerCellHeight(arrow.location, popupParams, session),
	}

	params := Params{
		Style:               "ruiPopup",
		ID:                  "ruiPopup",
		Row:                 popupRow,
		Column:              popupColumn,
		MaxWidth:            Percent(100),
		MaxHeight:           Percent(100),
		CellVerticalAlign:   StretchAlign,
		CellHorizontalAlign: StretchAlign,
		ClickEvent:          func(View) {},
		Shadow: NewShadowProperty(Params{
			SpreadRadius: Px(4),
			Blur:         Px(16),
			ColorTag:     "@ruiPopupShadow",
		}),
	}

	var closeButton View = nil
	var title View = nil
	outsideClose := false
	popup.buttons = []PopupButton{}
	titleStyle := "ruiPopupTitle"

	popup.showOpacity = 1.0
	popup.showDuration = 1.0
	popup.showTiming = "easy"

	for tag, value := range popupParams {
		if value != nil {
			switch tag {
			case VerticalAlign, HorizontalAlign, Arrow, Row, Column:
				// Do nothing

			case Margin:
				layerParams[Padding] = value

			case MarginLeft:
				layerParams[PaddingLeft] = value

			case MarginRight:
				layerParams[PaddingRight] = value

			case MarginTop:
				layerParams[PaddingTop] = value

			case MarginBottom:
				layerParams[PaddingBottom] = value

			case CloseButton:
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
				popup.cancelable = true

			case OutsideClose:
				outsideClose, _ = boolProperty(popupParams, OutsideClose, session)
				if outsideClose {
					popup.cancelable = true
				}

			case Buttons:
				switch value := value.(type) {
				case PopupButton:
					popup.buttons = []PopupButton{value}

				case []PopupButton:
					popup.buttons = value
				}

			case Title:
				switch value := value.(type) {
				case string:
					title = NewTextView(view.Session(), Params{Text: value})

				case View:
					title = value

				default:
					notCompatibleType(Title, value)
				}

			case TitleStyle:
				switch value := value.(type) {
				case string:
					titleStyle = value

				default:
					notCompatibleType(TitleStyle, value)
				}

			case DismissEvent:
				if listeners, ok := valueToPopupEventListeners(value); ok {
					if listeners != nil {
						popup.dismissListener = listeners
					}
				} else {
					notCompatibleType(tag, value)
				}

			case ArrowAlign:
				switch text := value.(type) {
				case string:
					switch text {
					case "top":
						value = "left"

					case "bottom":
						value = "right"
					}
				}
				arrow.align, _ = enumProperty(popupParams, ArrowAlign, session, CenterAlign)

			case ArrowSize:
				arrow.size, _ = sizeProperty(popupParams, ArrowSize, session)

			case ArrowOffset:
				arrow.off, _ = sizeProperty(popupParams, ArrowOffset, session)

			case ShowOpacity:
				if opacity, _ := floatProperty(popupParams, ShowOpacity, session, 1); opacity >= 0 && opacity < 1 {
					popup.showOpacity = opacity
				}

			case ShowTransform:
				if transform := valueToTransformProperty(value); transform != nil && !transform.empty() {
					popup.showTransform = transform
				}

			case ShowDuration:
				if duration, _ := floatProperty(popupParams, ShowDuration, session, 1); duration > 0 {
					popup.showDuration = duration
				}

			case ShowTiming:
				if text, ok := value.(string); ok {
					text, _ = session.resolveConstants(text)
					if isTimingFunctionValid(text) {
						popup.showTiming = text
					}
				}

			default:
				params[tag] = value
			}
		}
	}

	popup.popupView = NewGridLayout(view.Session(), params)

	var popupCellHeight []SizeUnit
	viewRow := 0
	if title != nil || closeButton != nil {
		titleContent := []View{}
		if title != nil {
			titleContent = append(titleContent, title)
		}
		if closeButton != nil {
			titleContent = append(titleContent, closeButton)
		}
		popup.popupView.Append(NewGridLayout(session, Params{
			Row:               0,
			Style:             titleStyle,
			CellWidth:         []any{Fr(1), AutoSize()},
			CellVerticalAlign: CenterAlign,
			PaddingLeft:       Px(12),
			Content:           titleContent,
		}))

		viewRow = 1
		popupCellHeight = []SizeUnit{AutoSize(), Fr(1)}
	} else {
		popupCellHeight = []SizeUnit{Fr(1)}
	}

	view.Set(Row, viewRow)
	popup.popupView.Append(view)

	if buttonCount := len(popup.buttons); buttonCount > 0 {
		buttonsAlign, _ := enumProperty(params, ButtonsAlign, session, RightAlign)
		popupCellHeight = append(popupCellHeight, AutoSize())
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

		for i, button := range popup.buttons {
			title := button.Title
			if title == "" && button.Type == CancelButton {
				title = "Cancel"
			}

			buttonView := NewButton(session, Params{
				Column:  i,
				Content: title,
			})

			if button.OnClick != nil {
				fn := button.OnClick
				buttonView.Set(ClickEvent, func() {
					fn(popup)
				})
			} else if button.Type == CancelButton {
				buttonView.Set(ClickEvent, popup.cancel)
			}

			if button.Type == DefaultButton {
				buttonView.Set(Style, "ruiDefaultButton")
			}

			buttonsPanel.Append(buttonView)
		}

		popup.popupView.Append(NewGridLayout(session, Params{
			Row:                 viewRow + 1,
			CellHorizontalAlign: buttonsAlign,
			Content:             buttonsPanel,
		}))
	}

	popup.popupView.Set(CellHeight, popupCellHeight)

	if arrow.location != NoneArrow {
		layerParams[Content] = []View{popup.popupView, arrow.createView(popup.popupView)}
	} else {
		layerParams[Content] = []View{popup.popupView}
	}

	popup.layerView = NewGridLayout(session, layerParams)

	if popup.showOpacity != 1 || popup.showTransform != nil {
		animation := NewAnimationProperty(Params{
			Duration:       popup.showDuration,
			TimingFunction: popup.showTiming,
		})
		if popup.showOpacity != 1 {
			popup.popupView.Set(Opacity, popup.showOpacity)
			popup.popupView.SetTransition(Opacity, animation)
		}
		if popup.showTransform != nil {
			popup.popupView.Set(Transform, popup.showTransform)
			popup.popupView.SetTransition(Transform, animation)
		}
	} else {
		session.updateCSSProperty("ruiPopupLayer", "transition", "")
	}

	if outsideClose {
		popup.layerView.Set(ClickEvent, popup.cancel)
	}
}

func (popup *popupData) showAnimation() {
	if popup.showOpacity != 1 || popup.showTransform != nil {
		htmlID := popup.popupView.htmlID()
		session := popup.Session()
		if popup.showOpacity != 1 {
			session.updateCSSProperty(htmlID, string(Opacity), "1")
		}
		if popup.showTransform != nil {
			session.updateCSSProperty(htmlID, string(Transform), "")
		}
	}
}

func (popup *popupData) dissmissAnimation(listener func(PropertyName)) bool {
	if popup.showOpacity != 1 || popup.showTransform != nil {
		session := popup.Session()
		popup.popupView.Set(TransitionEndEvent, listener)
		popup.popupView.Set(TransitionCancelEvent, listener)

		htmlID := popup.popupView.htmlID()
		if popup.showOpacity != 1 {
			session.updateCSSProperty(htmlID, string(Opacity), fmt.Sprintf("%.2f", popup.showOpacity))
		}
		if popup.showTransform != nil {
			session.updateCSSProperty(htmlID, string(Transform), popup.showTransform.transformCSS(session))
		}
		return true
	}
	return false
}

func (popup *popupData) View() View {
	return popup.contentView
}

func (popup *popupData) Session() Session {
	return popup.contentView.Session()
}

func (popup *popupData) cancel() {
	for _, button := range popup.buttons {
		if button.Type == CancelButton && button.OnClick != nil {
			button.OnClick(popup)
			return
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

func (popup *popupData) html(buffer *strings.Builder) {
	viewHTML(popup.layerView, buffer, "")
}

func (popup *popupData) viewByHTMLID(id string) View {
	return viewByHTMLID(id, popup.layerView)
}

func (popup *popupData) onDismiss() {
	popup.Session().callFunc("removeView", popup.layerView.htmlID())

	for _, listener := range popup.dismissListener {
		listener.Run(popup)
	}
}

func (popup *popupData) keyEvent(event KeyEvent) bool {
	if !event.AltKey && !event.CtrlKey && !event.ShiftKey && !event.MetaKey {
		switch event.Code {
		case EnterKey:
			for _, button := range popup.buttons {
				if button.Type == DefaultButton && button.OnClick != nil {
					button.OnClick(popup)
					return true
				}
			}

		case EscapeKey:
			if popup.cancelable {
				popup.Dismiss()
				return true
			}
		}
	}
	return false
}

// NewPopup creates a new Popup
func NewPopup(view View, param Params) Popup {
	if view == nil {
		return nil
	}

	popup := new(popupData)
	popup.init(view, param)
	return popup
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
		session.updateInnerHTML("ruiPopupLayer", "")
		return
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	for _, popup := range manager.popups {
		popup.html(buffer)
	}
	session.updateInnerHTML("ruiPopupLayer", buffer.String())
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
	session.updateCSSProperty("ruiPopupLayer", "visibility", "visible")
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
		if index == count-1 {
			if count == 1 {
				manager.popups = []Popup{}
				session.updateCSSProperty("ruiRoot", "pointer-events", "auto")
				session.updateCSSProperty("ruiPopupLayer", "visibility", "hidden")
			} else {
				manager.popups = manager.popups[:count-1]
			}
		} else if index == 0 {
			manager.popups = manager.popups[1:]
		} else {
			manager.popups = append(manager.popups[:index], manager.popups[index+1:]...)
		}
		popup.onDismiss()
	}

	if !popup.dissmissAnimation(listener) {
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

func newPopupListenerBinding(name string) popupListener {
	obj := new(popupListenerBinding)
	obj.name = name
	return obj
}

func (data *popupListenerBinding) Run(popup Popup) {
	bind := popup.View().binding()
	if bind == nil {
		ErrorLogF(`There is no a binding object for call "%s"`, data.name)
		return
	}

	val := reflect.ValueOf(bind)
	method := val.MethodByName(data.name)
	if !method.IsValid() {
		ErrorLogF(`The "%s" method is not valid`, data.name)
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
		return []popupListener{newPopupListenerBinding(value)}, true

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
					result = append(result, newPopupListenerBinding(v))

				default:
					return nil, false
				}
			}
		}
		return result, len(result) > 0
	}

	return nil, false
}
