package rui

// ShowMessage displays the popup with text message
func ShowMessage(title, text string, session Session) {
	textView := NewTextView(session, Params{
		Text:  text,
		Style: "ruiMessageText",
	})
	params := Params{
		CloseButton:  true,
		OutsideClose: true,
	}
	if title != "" {
		params[Title] = title
	}
	NewPopup(textView, params).Show()
}

func ShowQuestion(title, text string, session Session, onYes func(), onNo func()) {
	textView := NewTextView(session, Params{
		Text:  text,
		Style: "ruiMessageText",
	})
	params := Params{
		CloseButton:  false,
		OutsideClose: false,
		Buttons: []PopupButton{
			{
				Title: "No",
				OnClick: func(popup Popup) {
					popup.Dismiss()
					if onNo != nil {
						onNo()
					}
				},
			},
			{
				Title: "Yes",
				OnClick: func(popup Popup) {
					popup.Dismiss()
					if onYes != nil {
						onYes()
					}
				},
			},
		},
	}
	if title != "" {
		params[Title] = title
	}
	NewPopup(textView, params).Show()
}

func ShowCancellableQuestion(title, text string, session Session, onYes func(), onNo func(), onCancel func()) {
	textView := NewTextView(session, Params{
		Text:  text,
		Style: "ruiMessageText",
	})

	params := Params{
		CloseButton:  false,
		OutsideClose: false,
		Buttons: []PopupButton{
			{
				Title: "Cancel",
				OnClick: func(popup Popup) {
					popup.Dismiss()
					if onCancel != nil {
						onCancel()
					}
				},
			},
			{
				Title: "No",
				OnClick: func(popup Popup) {
					popup.Dismiss()
					if onNo != nil {
						onNo()
					}
				},
			},
			{
				Title: "Yes",
				OnClick: func(popup Popup) {
					popup.Dismiss()
					if onYes != nil {
						onYes()
					}
				},
			},
		},
	}
	if title != "" {
		params[Title] = title
	}
	NewPopup(textView, params).Show()
}

type popupMenuData struct {
	items   []string
	session Session
	popup   Popup
	result  func(int)
}

func (popup *popupMenuData) itemClick(list ListView, n int) {
	if popup.popup != nil {
		popup.popup.Dismiss()
		popup.popup = nil
	}
	if popup.result != nil {
		popup.result(n)
	}
}

func (popup *popupMenuData) ListSize() int {
	return len(popup.items)
}

func (popup *popupMenuData) ListItem(index int, session Session) View {
	return NewTextView(popup.session, Params{
		Text:  popup.items[index],
		Style: "ruiPopupMenuItem",
	})
}

func (popup *popupMenuData) IsListItemEnabled(index int) bool {
	return true
}

const PopupMenuResult = "popup-menu-result"

// ShowMenu displays the popup with text message
func ShowMenu(session Session, params Params) Popup {
	value, ok := params[Items]
	if !ok || value == nil {
		ErrorLog("Unable to show empty menu")
		return nil
	}

	var adapter ListAdapter
	data := new(popupMenuData)
	data.session = session

	switch value := value.(type) {
	case []string:
		data.items = value
		adapter = data

	case ListAdapter:
		adapter = value

	default:
		notCompatibleType(Items, value)
		return nil
	}

	value, ok = params[PopupMenuResult]
	if ok && value != nil {
		if result, ok := value.(func(int)); ok {
			data.result = result
		}
	}

	listView := NewListView(session, Params{
		Items:                adapter,
		Orientation:          TopDownOrientation,
		ListItemClickedEvent: data.itemClick,
	})

	popupParams := Params{}
	for tag, value := range params {
		switch tag {
		case Items, PopupMenuResult:
			// do nothing

		default:
			popupParams[tag] = value
		}
	}

	data.popup = NewPopup(listView, popupParams)
	data.popup.Show()
	FocusView(listView)
	return data.popup
}
