package rui

import (
	"strings"
)

const (
	// Title is the constant for the "title" property tag.
	// The "title" property is defined the Popup/Tabs title
	Title = "title"

	// TitleStyle is the constant for the "title-style" property tag.
	// The "title-style" string property is used to set the title style of the Popup.
	TitleStyle = "title-style"

	// CloseButton is the constant for the "close-button" property tag.
	// The "close-button" bool property allow to add the close button to the Popup.
	// Setting this property to "true" adds a window close button to the title bar (the default value is "false").
	CloseButton = "close-button"

	// OutsideClose is the constant for the "outside-close" property tag.
	// The "outside-close" is a bool property. If it is set to "true",
	// then clicking outside the popup window automatically calls the Dismiss() method.
	OutsideClose = "outside-close"

	// Buttons is the constant for the "buttons" property tag.
	// Using the "buttons" property you can add buttons that will be placed at the bottom of the Popup.
	// The "buttons" property can be assigned the following data types: PopupButton and []PopupButton
	Buttons = "buttons"

	// ButtonsAlign is the constant for the "buttons-align" property tag.
	// The "buttons-align" int property is used for set the horizontal alignment of Popup buttons.
	// Valid values: LeftAlign (0), RightAlign (1), CenterAlign (2), and StretchAlign (3)
	ButtonsAlign = "buttons-align"

	// DismissEvent is the constant for the "dismiss-event" property tag.
	// The "dismiss-event" event is used to track the closing of the Popup.
	// It occurs after the Popup disappears from the screen.
	// The main listener for this event has the following format: func(Popup)
	DismissEvent = "dismiss-event"

	// Arrow is the constant for the "arrow" property tag.
	// Using the "popup-arrow" int property you can add ...
	Arrow = "arrow"

	// ArrowAlign is the constant for the "arrow-align" property tag.
	// The "arrow-align" int property is used for set the horizontal alignment of the Popup arrow.
	// Valid values: LeftAlign (0), RightAlign (1), TopAlign (0), BottomAlign (1), CenterAlign (2)
	ArrowAlign = "arrow-align"

	// ArrowSize is the constant for the "arrow-size" property tag.
	// The "arrow-size" SizeUnit property is used for set the size (length) of the Popup arrow.
	ArrowSize = "arrow-size"

	// ArrowWidth is the constant for the "arrow-width" property tag.
	// The "arrow-width" SizeUnit property is used for set the width of the Popup arrow.
	ArrowWidth = "arrow-width"

	// ArrowOffset is the constant for the "arrow-offset" property tag.
	// The "arrow-offset" SizeUnit property is used for set the offset of the Popup arrow.
	ArrowOffset = "arrow-offset"

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

	// NormalButton is the constant of the popup button type: the normal button
	NormalButton PopupButtonType = 0
	// DefaultButton is the constant of the popup button type: button that fires when the "Enter" key is pressed
	DefaultButton PopupButtonType = 1
	// CancelButton is the constant of the popup button type: button that fires when the "Escape" key is pressed
	CancelButton PopupButtonType = 2
)

type PopupButtonType int

// PopupButton describes a button that will be placed at the bottom of the window.
type PopupButton struct {
	Title   string
	Type    PopupButtonType
	OnClick func(Popup)
}

// Popup interface
type Popup interface {
	View() View
	Session() Session
	Show()
	Dismiss()
	onDismiss()
	html(buffer *strings.Builder)
	viewByHTMLID(id string) View
	keyEvent(event KeyEvent) bool
}

type popupData struct {
	layerView       View
	view            View
	buttons         []PopupButton
	cancelable      bool
	dismissListener []func(Popup)
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

	if shadow := GetViewShadows(popupView); shadow != nil {
		params[Shadow] = shadow
	}

	if filter := GetBackdropFilter(popupView); filter != nil {
		params[BackdropFilter] = filter
	}

	switch arrow.location {
	case TopArrow:
		params[Row] = 0
		params[Column] = 1
		params[Clip] = PolygonClip([]any{"0%", "100%", "50%", "0%", "100%", "100%"})
		params[Width] = arrow.width
		params[Height] = arrow.size

	case RightArrow:
		params[Row] = 1
		params[Column] = 0
		params[Clip] = PolygonClip([]any{"0%", "0%", "100%", "50%", "0%", "100%"})
		params[Width] = arrow.size
		params[Height] = arrow.width

	case BottomArrow:
		params[Row] = 0
		params[Column] = 1
		params[Clip] = PolygonClip([]any{"0%", "0%", "50%", "100%", "100%", "0%"})
		params[Width] = arrow.width
		params[Height] = arrow.size

	case LeftArrow:
		params[Row] = 1
		params[Column] = 0
		params[Clip] = PolygonClip([]any{"100%", "0%", "0%", "50%", "100%", "100%"})
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

func (popup *popupData) init(view View, popupParams Params) {
	popup.view = view
	popup.cancelable = false
	session := view.Session()

	columnCount := 3
	rowCount := 3
	popupRow := 1
	popupColumn := 1
	arrow := popupArrow{
		row:    1,
		column: 1,
		align:  CenterAlign,
	}

	switch arrow.location, _ = enumProperty(popupParams, Arrow, session, NoneArrow); arrow.location {
	case TopArrow:
		rowCount = 4
		popupRow = 2

	case BottomArrow:
		rowCount = 4
		arrow.row = 2

	case LeftArrow:
		columnCount = 4
		popupColumn = 2

	case RightArrow:
		columnCount = 4
		arrow.column = 2

	default:
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

	layerParams := Params{
		Style:      "ruiPopupLayer",
		MaxWidth:   Percent(100),
		MaxHeight:  Percent(100),
		CellWidth:  cellWidth,
		CellHeight: cellHeight,
	}

	params := Params{
		Style:               "ruiPopup",
		Row:                 popupRow,
		Column:              popupColumn,
		MaxWidth:            Percent(100),
		MaxHeight:           Percent(100),
		CellVerticalAlign:   StretchAlign,
		CellHorizontalAlign: StretchAlign,
		ClickEvent:          func(View) {},
		Shadow: NewShadowWithParams(Params{
			SpreadRadius: Px(4),
			Blur:         Px(16),
			ColorTag:     "@ruiPopupShadow",
		}),
	}

	var closeButton View = nil
	outsideClose := false
	buttons := []PopupButton{}
	titleStyle := "ruiPopupTitle"
	var title View = nil

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
					buttons = []PopupButton{value}

				case []PopupButton:
					buttons = value
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
				if listeners, ok := valueToNoParamListeners[Popup](value); ok {
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

			default:
				params[tag] = value
			}
		}
	}

	popupView := NewGridLayout(view.Session(), params)

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
		popupView.Append(NewGridLayout(session, Params{
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
	popupView.Append(view)

	popup.buttons = buttons
	if buttonCount := len(buttons); buttonCount > 0 {
		buttonsAlign, _ := enumProperty(params, ButtonsAlign, session, RightAlign)
		popupCellHeight = append(popupCellHeight, AutoSize())
		gap, _ := sizeConstant(session, "ruiPopupButtonGap")
		cellWidth := []SizeUnit{}
		for i := 0; i < buttonCount; i++ {
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

		popupView.Append(NewGridLayout(session, Params{
			Row:                 viewRow + 1,
			CellHorizontalAlign: buttonsAlign,
			Content:             buttonsPanel,
		}))
	}

	popupView.Set(CellHeight, popupCellHeight)

	if arrow.location != NoneArrow {
		layerParams[Content] = []View{popupView, arrow.createView(popupView)}
	} else {
		layerParams[Content] = []View{popupView}
	}

	popup.layerView = NewGridLayout(session, layerParams)
	if outsideClose {
		popup.layerView.Set(ClickEvent, popup.cancel)
	}
}

func (popup popupData) View() View {
	return popup.view
}

func (popup *popupData) Session() Session {
	return popup.view.Session()
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
	for _, listener := range popup.dismissListener {
		listener(popup)
	}
}

func (popup *popupData) Show() {
	popup.Session().popupManager().showPopup(popup)
}

func (popup *popupData) html(buffer *strings.Builder) {

	viewHTML(popup.layerView, buffer)
}

func (popup *popupData) viewByHTMLID(id string) View {
	return viewByHTMLID(id, popup.layerView)
}

func (popup *popupData) onDismiss() {
	for _, listener := range popup.dismissListener {
		listener(popup)
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

	session := popup.Session()
	if manager.popups[count-1] == popup {
		if count == 1 {
			manager.popups = []Popup{}
			session.updateCSSProperty("ruiRoot", "pointer-events", "auto")
			session.updateCSSProperty("ruiPopupLayer", "visibility", "hidden")
			session.updateInnerHTML("ruiPopupLayer", "")
		} else {
			manager.popups = manager.popups[:count-1]
			manager.updatePopupLayerInnerHTML(session)
		}
		popup.onDismiss()
		return
	}

	for n, p := range manager.popups {
		if p == popup {
			if n == 0 {
				manager.popups = manager.popups[1:]
			} else {
				manager.popups = append(manager.popups[:n], manager.popups[n+1:]...)
			}
			manager.updatePopupLayerInnerHTML(session)
			popup.onDismiss()
			return
		}
	}
}
