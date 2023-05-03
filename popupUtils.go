package rui

// ShowMessage displays the popup with the title given in the "title" argument and the message text given in the "text" argument.
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

// ShowQuestion displays a message with the given title and text and two buttons "Yes" and "No".
// When the "Yes" button is clicked, the message is closed and the onYes function is called (if it is not nil).
// When the "No" button is pressed, the message is closed and the onNo function is called (if it is not nil).
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

// ShowCancellableQuestion displays a message with the given title and text and three buttons "Yes", "No" and "Cancel".
// When the "Yes", "No" or "Cancel" button is pressed, the message is closed and the onYes, onNo or onCancel function
// (if it is not nil) is called, respectively.
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
	items    []string
	disabled []int
	session  Session
	popup    Popup
	result   func(int)
}

func (popup *popupMenuData) itemClick(list ListView, n int) {
	if popup.IsListItemEnabled(n) {
		if popup.popup != nil {
			popup.popup.Dismiss()
			popup.popup = nil
		}
		if popup.result != nil {
			popup.result(n)
		}
	}
}

func (popup *popupMenuData) ListSize() int {
	return len(popup.items)
}

func (popup *popupMenuData) ListItem(index int, session Session) View {
	view := NewTextView(popup.session, Params{
		Text:  popup.items[index],
		Style: "ruiPopupMenuItem",
	})
	if !popup.IsListItemEnabled(index) {
		view.Set(TextColor, "@ruiDisabledTextColor")
	}
	return view
}

func (popup *popupMenuData) IsListItemEnabled(index int) bool {
	if popup.disabled != nil {
		for _, n := range popup.disabled {
			if index == n {
				return false
			}
		}
	}
	return true
}

// PopupMenuResult is the constant for the "popup-menu-result" property tag.
// The "popup-menu-result" property sets the function (format: func(int)) to be called when
// a menu item of popup menu is selected.
const PopupMenuResult = "popup-menu-result"

// ShowMenu displays the menu. Menu items are set using the Items property.
// The "popup-menu-result" property sets the function (format: func(int)) to be called when a menu item is selected.
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

	if value, ok := params[PopupMenuResult]; ok && value != nil {
		if result, ok := value.(func(int)); ok {
			data.result = result
		}
	}

	if value, ok := params[DisabledItems]; ok && value != nil {
		if value, ok := value.([]int); ok {
			data.disabled = value
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
		case Items, PopupMenuResult, DisabledItems:
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
