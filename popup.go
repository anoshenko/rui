package rui

import "strings"

const (
	// Title is the Popup string property
	Title = "title"
	// TitleStyle is the Popup string property
	TitleStyle = "title-style"
	// CloseButton is the Popup bool property
	CloseButton = "close-button"
	// OutsideClose is the Popup bool property
	OutsideClose = "outside-close"
	Buttons      = "buttons"
	ButtonsAlign = "buttons-align"
	DismissEvent = "dismiss-event"
)

type PopupButton struct {
	Title   string
	OnClick func(Popup)
}

// Popup interface
type Popup interface {
	//Properties
	View() View
	Session() Session
	Show()
	Dismiss()
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

func (popup *popupData) init(view View, params Params) {
	popup.view = view
	session := view.Session()

	popup.dismissListener = []func(Popup){}
	if value, ok := params[DismissEvent]; ok && value != nil {
		switch value := value.(type) {
		case func(Popup):
			popup.dismissListener = []func(Popup){value}

		case func():
			popup.dismissListener = []func(Popup){
				func(popup Popup) {
					value()
				},
			}

		case []func(Popup):
			for _, fn := range value {
				if fn != nil {
					popup.dismissListener = append(popup.dismissListener, fn)
				}
			}

		case []func():
			for _, fn := range value {
				if fn != nil {
					popup.dismissListener = append(popup.dismissListener, func(popup Popup) {
						fn()
					})
				}
			}

		case []interface{}:
			for _, val := range value {
				if val != nil {
					switch fn := val.(type) {
					case func(Popup):
						popup.dismissListener = append(popup.dismissListener, fn)

					case func():
						popup.dismissListener = append(popup.dismissListener, func(popup Popup) {
							fn()
						})
					}
				}
			}
		}
	}

	var title View = nil
	titleStyle := "ruiPopupTitle"
	closeButton, _ := boolProperty(params, CloseButton, session)
	outsideClose, _ := boolProperty(params, OutsideClose, session)
	vAlign, _ := enumProperty(params, VerticalAlign, session, CenterAlign)
	hAlign, _ := enumProperty(params, HorizontalAlign, session, CenterAlign)
	buttonsAlign, _ := enumProperty(params, ButtonsAlign, session, RightAlign)

	buttons := []PopupButton{}
	if value, ok := params[Buttons]; ok && value != nil {
		switch value := value.(type) {
		case PopupButton:
			buttons = []PopupButton{value}

		case []PopupButton:
			buttons = value
		}
	}

	popupView := NewGridLayout(view.Session(), Params{
		Style:               "ruiPopup",
		MaxWidth:            Percent(100),
		MaxHeight:           Percent(100),
		CellVerticalAlign:   StretchAlign,
		CellHorizontalAlign: StretchAlign,
		ClickEvent:          func(View) {},
	})

	for tag, value := range params {
		switch tag {
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

		case CloseButton, OutsideClose, VerticalAlign, HorizontalAlign, Buttons:
			// do nothing

		default:
			popupView.Set(tag, value)
		}
	}

	var cellHeight []SizeUnit
	viewRow := 0
	if title != nil || closeButton {
		viewRow = 1
		titleHeight, _ := sizeConstant(popup.Session(), "ruiPopupTitleHeight")
		titleView := NewGridLayout(session, Params{
			Row:               0,
			Style:             titleStyle,
			CellWidth:         []SizeUnit{Fr(1), titleHeight},
			CellVerticalAlign: CenterAlign,
			PaddingLeft:       Px(12),
		})
		if title != nil {
			titleView.Append(title)
		}
		if closeButton {
			titleView.Append(NewGridLayout(session, Params{
				Column:              1,
				Height:              titleHeight,
				Width:               titleHeight,
				CellHorizontalAlign: CenterAlign,
				CellVerticalAlign:   CenterAlign,
				TextSize:            Px(20),
				Content:             "✕",
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
		Content:             NewColumnLayout(session, Params{Content: popupView}),
		MaxWidth:            Percent(100),
		MaxHeight:           Percent(100),
	})

	if outsideClose {
		popup.layerView.Set(ClickEvent, func(View) {
			popup.Dismiss()
		})
	}
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

// NewPopup creates a new Popup
func NewPopup(view View, param Params) Popup {
	if view == nil {
		return nil
	}

	popup := new(popupData)
	popup.init(view, param)
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
			updateCSSProperty("ruiPopupLayer", "visibility", "hidden", session)
			session.runScript(`updateInnerHTML('ruiPopupLayer', '');`)
		} else {
			manager.popups = manager.popups[:count-1]
			manager.updatePopupLayerInnerHTML(session)
		}
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
			return
		}
	}
}
