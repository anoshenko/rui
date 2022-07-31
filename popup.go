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
	// The "arrow-align" int property is used for set the horizontal alignment of Popup arrow.
	// Valid values: LeftAlign (0), RightAlign (1), TopAlign (0), BottomAlign (1), CenterAlign (2)
	ArrowAlign = "arrow-align"

	// ArrowSize is the constant for the "arrow-size" property tag.
	// The "arrow-size" SizeUnit property is used for set the size of Popup arrow.
	ArrowSize = "arrow-size"

	// ArrowOffset is the constant for the "arrow-offset" property tag.
	// The "arrow-offset" SizeUnit property is used for set the offset of Popup arrow.
	ArrowOffset = "arrow-offset"

	NoneArrow   = 0
	TopArrow    = 1
	RightArrow  = 2
	BottomArrow = 3
	LeftArrow   = 4
)

// PopupButton describes a button that will be placed at the bottom of the window.
type PopupButton struct {
	Title   string
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
}

type popupData struct {
	layerView       View
	view            View
	dismissListener []func(Popup)
}

type popupManager struct {
	popups []Popup
}

func (popup *popupData) init(view View, popupParams Params) {
	popup.view = view
	session := view.Session()

	params := Params{
		Style:               "ruiPopup",
		MaxWidth:            Percent(100),
		MaxHeight:           Percent(100),
		CellVerticalAlign:   StretchAlign,
		CellHorizontalAlign: StretchAlign,
		ClickEvent:          func(View) {},
	}

	closeButton := false
	outsideClose := false
	vAlign := CenterAlign
	hAlign := CenterAlign
	arrow := 0
	arrowAlign := CenterAlign
	arrowSize := AutoSize()
	arrowOff := AutoSize()
	buttons := []PopupButton{}
	titleStyle := "ruiPopupTitle"
	var title View = nil

	for tag, value := range popupParams {
		if value != nil {
			switch tag {
			case CloseButton:
				closeButton, _ = boolProperty(popupParams, CloseButton, session)

			case OutsideClose:
				outsideClose, _ = boolProperty(popupParams, OutsideClose, session)

			case VerticalAlign:
				vAlign, _ = enumProperty(popupParams, VerticalAlign, session, CenterAlign)

			case HorizontalAlign:
				hAlign, _ = enumProperty(popupParams, HorizontalAlign, session, CenterAlign)

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

			case Arrow:
				arrow, _ = enumProperty(popupParams, Arrow, session, NoneArrow)

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
				arrowAlign, _ = enumProperty(popupParams, ArrowAlign, session, CenterAlign)

			case ArrowSize:
				arrowSize, _ = sizeProperty(popupParams, ArrowSize, session)

			case ArrowOffset:
				arrowOff, _ = sizeProperty(popupParams, ArrowOffset, session)

			default:
				params[tag] = value
			}
		}
	}

	popupView := NewGridLayout(view.Session(), params)

	var cellHeight []SizeUnit
	viewRow := 0
	if title != nil || closeButton {
		viewRow = 1
		titleView := NewGridLayout(session, Params{
			Row:               0,
			Style:             titleStyle,
			CellWidth:         []any{Fr(1), "@ruiPopupTitleHeight"},
			CellVerticalAlign: CenterAlign,
			PaddingLeft:       Px(12),
		})
		if title != nil {
			titleView.Append(title)
		}
		if closeButton {
			titleView.Append(NewGridLayout(session, Params{
				Column:              1,
				Height:              "@ruiPopupTitleHeight",
				Width:               "@ruiPopupTitleHeight",
				CellHorizontalAlign: CenterAlign,
				CellVerticalAlign:   CenterAlign,
				TextSize:            Px(20),
				Content:             "✕",
				NotTranslate:        true,
				ClickEvent: func(View) {
					popup.Dismiss()
				},
			}))
		}

		popupView.Append(titleView)
		cellHeight = []SizeUnit{AutoSize(), Fr(1)}
	} else {
		cellHeight = []SizeUnit{Fr(1)}
	}

	view.Set(Row, viewRow)
	popupView.Append(view)

	if buttonCount := len(buttons); buttonCount > 0 {
		buttonsAlign, _ := enumProperty(params, ButtonsAlign, session, RightAlign)
		cellHeight = append(cellHeight, AutoSize())
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

		createButton := func(n int, button PopupButton) Button {
			return NewButton(session, Params{
				Column:  n,
				Content: button.Title,
				ClickEvent: func() {
					if button.OnClick != nil {
						button.OnClick(popup)
					} else {
						popup.Dismiss()
					}
				},
			})
		}
		for i, button := range buttons {
			buttonsPanel.Append(createButton(i, button))
		}

		popupView.Append(NewGridLayout(session, Params{
			Row:                 viewRow + 1,
			CellHorizontalAlign: buttonsAlign,
			Content:             buttonsPanel,
		}))
	}
	popupView.Set(CellHeight, cellHeight)

	popup.layerView = NewGridLayout(session, Params{
		Style:               "ruiPopupLayer",
		CellVerticalAlign:   vAlign,
		CellHorizontalAlign: hAlign,
		MaxWidth:            Percent(100),
		MaxHeight:           Percent(100),
		Filter: NewViewFilter(Params{
			DropShadow: NewShadowWithParams(Params{
				SpreadRadius: Px(4),
				Blur:         Px(16),
				ColorTag:     "@ruiPopupShadow",
			}),
		}),
		Content: NewColumnLayout(session, Params{
			Content: popup.createArrow(arrow, arrowAlign, arrowSize, arrowOff, popupView),
		}),
	})

	if outsideClose {
		popup.layerView.Set(ClickEvent, func(View) {
			popup.Dismiss()
		})
	}
}

func (popup popupData) createArrow(arrow int, align int, size SizeUnit, off SizeUnit, popupView View) View {
	if arrow == NoneArrow {
		return popupView
	}

	session := popupView.Session()

	if size.Type == Auto {
		size = Px(16)
		if value, ok := session.Constant("ruiArrowSize"); ok {
			size, _ = StringToSizeUnit(value)
		}
	}

	color := GetBackgroundColor(popupView, "")
	border := NewBorder(Params{
		Style:    SolidLine,
		ColorTag: color,
		Width:    size,
	})
	border2 := NewBorder(Params{
		Style:    SolidLine,
		ColorTag: 1,
		Width:    SizeUnit{Type: size.Type, Value: size.Value / 2},
	})

	popupParams := Params{
		MaxWidth:  Percent(100),
		MaxHeight: Percent(100),
		Content:   popupView,
	}
	arrowParams := Params{
		Width:  size,
		Height: size,
	}
	containerParams := Params{
		MaxWidth:  Percent(100),
		MaxHeight: Percent(100),
	}

	switch arrow {
	case TopArrow:
		arrowParams[BorderBottom] = border
		arrowParams[BorderLeft] = border2
		arrowParams[BorderRight] = border2
		popupParams[Row] = 1
		containerParams[CellHeight] = []SizeUnit{AutoSize(), Fr(1)}

	case RightArrow:
		arrowParams[Column] = 1
		arrowParams[BorderLeft] = border
		arrowParams[BorderTop] = border2
		arrowParams[BorderBottom] = border2
		containerParams[CellWidth] = []SizeUnit{Fr(1), AutoSize()}

	case BottomArrow:
		arrowParams[Row] = 1
		arrowParams[BorderTop] = border
		arrowParams[BorderLeft] = border2
		arrowParams[BorderRight] = border2
		containerParams[CellHeight] = []SizeUnit{Fr(1), AutoSize()}

	case LeftArrow:
		arrowParams[BorderRight] = border
		arrowParams[BorderTop] = border2
		arrowParams[BorderBottom] = border2
		popupParams[Column] = 1
		containerParams[CellWidth] = []SizeUnit{AutoSize(), Fr(1)}

	default:
		return popupView
	}

	switch align {
	case LeftAlign:
		switch arrow {
		case TopArrow:
			if off.Type == Auto {
				off = GetRadius(popupView, "").TopLeftX
			}
			if off.Type != Auto {
				arrowParams[MarginLeft] = off
			}

		case RightArrow:
			if off.Type == Auto {
				off = GetRadius(popupView, "").TopRightY
			}
			if off.Type != Auto {
				arrowParams[MarginTop] = off
			}

		case BottomArrow:
			if off.Type == Auto {
				off = GetRadius(popupView, "").BottomLeftX
			}
			if off.Type != Auto {
				arrowParams[MarginLeft] = off
			}

		case LeftArrow:
			if off.Type == Auto {
				off = GetRadius(popupView, "").TopLeftY
			}
			if off.Type != Auto {
				arrowParams[MarginTop] = off
			}
		}

	case RightAlign:
		switch arrow {
		case TopArrow:
			if off.Type == Auto {
				off = GetRadius(popupView, "").TopRightX
			}
			if off.Type != Auto {
				arrowParams[MarginRight] = off
			}

		case RightArrow:
			if off.Type == Auto {
				off = GetRadius(popupView, "").BottomRightY
			}
			if off.Type != Auto {
				arrowParams[MarginBottom] = off
			}

		case BottomArrow:
			if off.Type == Auto {
				off = GetRadius(popupView, "").BottomRightX
			}
			if off.Type != Auto {
				arrowParams[MarginRight] = off
			}

		case LeftArrow:
			if off.Type == Auto {
				off = GetRadius(popupView, "").BottomLeftY
			}
			if off.Type != Auto {
				arrowParams[MarginBottom] = off
			}
		}

	default:
		align = CenterAlign
		if off.Type != Auto {
			switch arrow {
			case TopArrow, BottomArrow:
				if off.Value < 0 {
					off.Value = -off.Value
					arrowParams[MarginRight] = off
				} else {
					arrowParams[MarginLeft] = off
				}

			case RightArrow, LeftArrow:
				if off.Value < 0 {
					off.Value = -off.Value
					arrowParams[MarginBottom] = off
				} else {
					arrowParams[MarginTop] = off
				}
			}
		}
	}

	if margin := popupView.Get(Margin); margin != nil {
		popupView.Remove(Margin)
		containerParams[Padding] = margin
	}

	/*
		for key, value := range popupParams {
			if key != Content {
				popupView.Set(key, value)
			}
		}
	*/

	containerParams[CellVerticalAlign] = align
	containerParams[CellHorizontalAlign] = align
	containerParams[Content] = []View{
		NewView(session, arrowParams),
		//popupView,
		NewColumnLayout(session, popupParams),
	}

	return NewGridLayout(session, containerParams)
}

func (popup popupData) View() View {
	return popup.view
}

func (popup *popupData) Session() Session {
	return popup.view.Session()
}

func (popup *popupData) Dismiss() {
	popup.Session().popupManager().dismissPopup(popup)
	for _, listener := range popup.dismissListener {
		listener(popup)
	}
	// TODO
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
		session.runScript(`updateInnerHTML('ruiPopupLayer', '');`)
		return
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`updateInnerHTML('ruiPopupLayer', '`)
	for _, p := range manager.popups {
		p.html(buffer)
	}
	buffer.WriteString(`');`)
	session.runScript(buffer.String())
}

func (manager *popupManager) showPopup(popup Popup) {
	if popup == nil {
		return
	}

	session := popup.Session()
	if manager.popups == nil || len(manager.popups) == 0 {
		manager.popups = []Popup{popup}
	} else {
		manager.popups = append(manager.popups, popup)
	}

	session.runScript(`if (document.activeElement != document.body) document.activeElement.blur();`)
	manager.updatePopupLayerInnerHTML(session)
	updateCSSProperty("ruiPopupLayer", "visibility", "visible", session)
	updateCSSProperty("ruiRoot", "pointer-events", "none", session)
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
			updateCSSProperty("ruiRoot", "pointer-events", "auto", session)
			updateCSSProperty("ruiPopupLayer", "visibility", "hidden", session)
			session.runScript(`updateInnerHTML('ruiPopupLayer', '');`)
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
