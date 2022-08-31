package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// EditTextChangedEvent is the constant for the "edit-text-changed" property tag.
	EditTextChangedEvent = "edit-text-changed"
	// EditViewType is the constant for the "edit-view-type" property tag.
	EditViewType = "edit-view-type"
	// EditViewPattern is the constant for the "edit-view-pattern" property tag.
	EditViewPattern = "edit-view-pattern"
	// Spellcheck is the constant for the "spellcheck" property tag.
	Spellcheck = "spellcheck"
)

const (
	// SingleLineText - single-line text type of EditView
	SingleLineText = 0
	// PasswordText - password type of EditView
	PasswordText = 1
	// EmailText - e-mail type of EditView. Allows to enter one email
	EmailText = 2
	// EmailsText - e-mail type of EditView. Allows to enter multiple emails separeted by comma
	EmailsText = 3
	// URLText - url type of EditView. Allows to enter one url
	URLText = 4
	// PhoneText - telephone type of EditView. Allows to enter one phone number
	PhoneText = 5
	// MultiLineText - multi-line text type of EditView
	MultiLineText = 6
)

// EditView - grid-container of View
type EditView interface {
	View
	AppendText(text string)
}

type editViewData struct {
	viewData
	textChangeListeners []func(EditView, string)
}

// NewEditView create new EditView object and return it
func NewEditView(session Session, params Params) EditView {
	view := new(editViewData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newEditView(session Session) View {
	return NewEditView(session, nil)
}

func (edit *editViewData) Init(session Session) {
	edit.viewData.Init(session)
	edit.textChangeListeners = []func(EditView, string){}
	edit.tag = "EditView"
}

func (edit *editViewData) String() string {
	return getViewString(edit)
}

func (edit *editViewData) Focusable() bool {
	return true
}

func (edit *editViewData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Type, "edit-type":
		return EditViewType

	case Pattern, "edit-pattern":
		return EditViewPattern

	case "maxlength", "maxlen":
		return MaxLength

	case "wrap":
		return EditWrap
	}

	return tag
}

func (edit *editViewData) Remove(tag string) {
	edit.remove(edit.normalizeTag(tag))
}

func (edit *editViewData) remove(tag string) {
	_, exists := edit.properties[tag]
	switch tag {
	case Hint:
		if exists {
			delete(edit.properties, Hint)
			if edit.created {
				removeProperty(edit.htmlID(), "placeholder", edit.session)
			}
			edit.propertyChangedEvent(tag)
		}

	case MaxLength:
		if exists {
			delete(edit.properties, MaxLength)
			if edit.created {
				removeProperty(edit.htmlID(), "maxlength", edit.session)
			}
			edit.propertyChangedEvent(tag)
		}

	case ReadOnly, Spellcheck:
		if exists {
			delete(edit.properties, tag)
			if edit.created {
				updateBoolProperty(edit.htmlID(), tag, false, edit.session)
			}
			edit.propertyChangedEvent(tag)
		}

	case EditTextChangedEvent:
		if len(edit.textChangeListeners) > 0 {
			edit.textChangeListeners = []func(EditView, string){}
			edit.propertyChangedEvent(tag)
		}

	case Text:
		if exists {
			oldText := GetText(edit)
			delete(edit.properties, tag)
			if oldText != "" {
				edit.textChanged("")
				if edit.created {
					edit.session.runScript(fmt.Sprintf(`setInputValue('%s', '%s')`, edit.htmlID(), ""))
				}
			}
		}

	case EditViewPattern:
		if exists {
			oldText := GetEditViewPattern(edit)
			delete(edit.properties, tag)
			if oldText != "" {
				if edit.created {
					removeProperty(edit.htmlID(), Pattern, edit.session)
				}
				edit.propertyChangedEvent(tag)
			}
		}

	case EditViewType:
		if exists {
			oldType := GetEditViewType(edit)
			delete(edit.properties, tag)
			if oldType != 0 {
				if edit.created {
					updateInnerHTML(edit.parentHTMLID(), edit.session)
				}
				edit.propertyChangedEvent(tag)
			}
		}

	case EditWrap:
		if exists {
			oldWrap := IsEditViewWrap(edit)
			delete(edit.properties, tag)
			if GetEditViewType(edit) == MultiLineText {
				if wrap := IsEditViewWrap(edit); wrap != oldWrap {
					if edit.created {
						if wrap {
							updateProperty(edit.htmlID(), "wrap", "soft", edit.session)
						} else {
							updateProperty(edit.htmlID(), "wrap", "off", edit.session)
						}
					}
					edit.propertyChangedEvent(tag)
				}
			}
		}

	default:
		edit.viewData.remove(tag)
		return
	}
}

func (edit *editViewData) Set(tag string, value any) bool {
	return edit.set(edit.normalizeTag(tag), value)
}

func (edit *editViewData) set(tag string, value any) bool {
	if value == nil {
		edit.remove(tag)
		return true
	}

	switch tag {
	case Text:
		oldText := GetText(edit)
		if text, ok := value.(string); ok {
			edit.properties[Text] = text
			if text = GetText(edit); oldText != text {
				edit.textChanged(text)
				if edit.created {
					if GetEditViewType(edit) == MultiLineText {
						updateInnerHTML(edit.htmlID(), edit.Session())
					} else {
						text = strings.ReplaceAll(text, `"`, `\"`)
						text = strings.ReplaceAll(text, `'`, `\'`)
						text = strings.ReplaceAll(text, "\n", `\n`)
						text = strings.ReplaceAll(text, "\r", `\r`)
						edit.session.runScript(fmt.Sprintf(`setInputValue('%s', '%s')`, edit.htmlID(), text))
					}
				}
			}
			return true
		}
		return false

	case Hint:
		oldText := GetHint(edit)
		if text, ok := value.(string); ok {
			edit.properties[Hint] = text
			if text = GetHint(edit); oldText != text {
				if edit.created {
					if text != "" {
						updateProperty(edit.htmlID(), "placeholder", text, edit.session)
					} else {
						removeProperty(edit.htmlID(), "placeholder", edit.session)
					}
				}
				edit.propertyChangedEvent(tag)
			}
			return true
		}
		return false

	case MaxLength:
		oldMaxLength := GetMaxLength(edit)
		if edit.setIntProperty(MaxLength, value) {
			if maxLength := GetMaxLength(edit); maxLength != oldMaxLength {
				if edit.created {
					if maxLength > 0 {
						updateProperty(edit.htmlID(), "maxlength", strconv.Itoa(maxLength), edit.session)
					} else {
						removeProperty(edit.htmlID(), "maxlength", edit.session)
					}
				}
				edit.propertyChangedEvent(tag)
			}
			return true
		}
		return false

	case ReadOnly:
		if edit.setBoolProperty(ReadOnly, value) {
			if edit.created {
				if IsReadOnly(edit) {
					updateProperty(edit.htmlID(), ReadOnly, "", edit.session)
				} else {
					removeProperty(edit.htmlID(), ReadOnly, edit.session)
				}
			}
			edit.propertyChangedEvent(tag)
			return true
		}
		return false

	case Spellcheck:
		if edit.setBoolProperty(Spellcheck, value) {
			if edit.created {
				updateBoolProperty(edit.htmlID(), Spellcheck, IsSpellcheck(edit), edit.session)
			}
			edit.propertyChangedEvent(tag)
			return true
		}
		return false

	case EditViewPattern:
		oldText := GetEditViewPattern(edit)
		if text, ok := value.(string); ok {
			edit.properties[EditViewPattern] = text
			if text = GetEditViewPattern(edit); oldText != text {
				if edit.created {
					if text != "" {
						updateProperty(edit.htmlID(), Pattern, text, edit.session)
					} else {
						removeProperty(edit.htmlID(), Pattern, edit.session)
					}
				}
				edit.propertyChangedEvent(tag)
			}
			return true
		}
		return false

	case EditViewType:
		oldType := GetEditViewType(edit)
		if edit.setEnumProperty(EditViewType, value, enumProperties[EditViewType].values) {
			if GetEditViewType(edit) != oldType {
				if edit.created {
					updateInnerHTML(edit.parentHTMLID(), edit.session)
				}
				edit.propertyChangedEvent(tag)
			}
			return true
		}
		return false

	case EditWrap:
		oldWrap := IsEditViewWrap(edit)
		if edit.setBoolProperty(EditWrap, value) {
			if GetEditViewType(edit) == MultiLineText {
				if wrap := IsEditViewWrap(edit); wrap != oldWrap {
					if edit.created {
						if wrap {
							updateProperty(edit.htmlID(), "wrap", "soft", edit.session)
						} else {
							updateProperty(edit.htmlID(), "wrap", "off", edit.session)
						}
					}
					edit.propertyChangedEvent(tag)
				}
			}
			return true
		}
		return false

	case EditTextChangedEvent:
		listeners, ok := valueToEventListeners[EditView, string](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(EditView, string){}
		}
		edit.textChangeListeners = listeners
		edit.propertyChangedEvent(tag)
		return true
	}

	return edit.viewData.set(tag, value)
}

func (edit *editViewData) Get(tag string) any {
	return edit.get(edit.normalizeTag(tag))
}

func (edit *editViewData) get(tag string) any {
	if tag == EditTextChangedEvent {
		return edit.textChangeListeners
	}
	return edit.viewData.get(tag)
}

func (edit *editViewData) AppendText(text string) {
	if GetEditViewType(edit) == MultiLineText {
		if value := edit.getRaw(Text); value != nil {
			if textValue, ok := value.(string); ok {
				textValue += text
				edit.properties[Text] = textValue

				text := strings.ReplaceAll(text, `"`, `\"`)
				text = strings.ReplaceAll(text, `'`, `\'`)
				text = strings.ReplaceAll(text, "\n", `\n`)
				text = strings.ReplaceAll(text, "\r", `\r`)

				edit.session.runScript(`appendToInnerHTML("` + edit.htmlID() + `", "` + text + `")`)

				edit.textChanged(textValue)
				return
			}
		}
		edit.set(Text, text)
	} else {
		edit.set(Text, GetText(edit)+text)
	}
}

func (edit *editViewData) textChanged(newText string) {
	for _, listener := range edit.textChangeListeners {
		listener(edit, newText)
	}
	edit.propertyChangedEvent(Text)
}

func (edit *editViewData) htmlTag() string {
	if GetEditViewType(edit) == MultiLineText {
		return "textarea"
	}
	return "input"
}

func (edit *editViewData) htmlProperties(self View, buffer *strings.Builder) {
	edit.viewData.htmlProperties(self, buffer)

	writeSpellcheck := func() {
		if spellcheck := IsSpellcheck(edit); spellcheck {
			buffer.WriteString(` spellcheck="true"`)
		} else {
			buffer.WriteString(` spellcheck="false"`)
		}
	}

	editType := GetEditViewType(edit)
	switch editType {
	case SingleLineText:
		buffer.WriteString(` type="text" inputmode="text"`)
		writeSpellcheck()

	case PasswordText:
		buffer.WriteString(` type="password" inputmode="text"`)

	case EmailText:
		buffer.WriteString(` type="email" inputmode="email"`)

	case EmailsText:
		buffer.WriteString(` type="email" inputmode="email" multiple`)

	case URLText:
		buffer.WriteString(` type="url" inputmode="url"`)

	case PhoneText:
		buffer.WriteString(` type="tel" inputmode="tel"`)

	case MultiLineText:
		if IsEditViewWrap(edit) {
			buffer.WriteString(` wrap="soft"`)
		} else {
			buffer.WriteString(` wrap="off"`)
		}
		writeSpellcheck()
	}

	if IsReadOnly(edit) {
		buffer.WriteString(` readonly`)
	}

	if maxLength := GetMaxLength(edit); maxLength > 0 {
		buffer.WriteString(` maxlength="`)
		buffer.WriteString(strconv.Itoa(maxLength))
		buffer.WriteByte('"')
	}

	convertText := func(text string) string {
		if strings.ContainsRune(text, '"') {
			text = strings.ReplaceAll(text, `"`, `&#34;`)
		}
		return textToJS(text)
	}

	if hint := GetHint(edit); hint != "" {
		buffer.WriteString(` placeholder="`)
		buffer.WriteString(convertText(hint))
		buffer.WriteByte('"')
	}

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)
	if pattern := GetEditViewPattern(edit); pattern != "" {
		buffer.WriteString(` pattern="`)
		buffer.WriteString(convertText(pattern))
		buffer.WriteByte('"')
	}

	if editType != MultiLineText {
		if text := GetText(edit); text != "" {
			buffer.WriteString(` value="`)
			buffer.WriteString(convertText(text))
			buffer.WriteByte('"')
		}
	}
}

func (edit *editViewData) htmlDisabledProperties(self View, buffer *strings.Builder) {
	if IsDisabled(self) {
		buffer.WriteString(` disabled`)
	}
	edit.viewData.htmlDisabledProperties(self, buffer)
}

func (edit *editViewData) htmlSubviews(self View, buffer *strings.Builder) {
	if GetEditViewType(edit) == MultiLineText {
		buffer.WriteString(textToJS(GetText(edit)))
	}
}

func (edit *editViewData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "textChanged":
		oldText := GetText(edit)
		if text, ok := data.PropertyValue("text"); ok {
			edit.properties[Text] = text
			if text := GetText(edit); text != oldText {
				edit.textChanged(text)
			}
		}
		return true
	}

	return edit.viewData.handleCommand(self, command, data)
}

// GetText returns a text of the EditView subview.
// If the second argument (subviewID) is not specified or it is "" then a text of the first argument (view) is returned.
func GetText(view View, subviewID ...string) string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if value := view.getRaw(Text); value != nil {
			if text, ok := value.(string); ok {
				return text
			}
		}
	}
	return ""
}

// GetHint returns a hint text of the subview.
// If the second argument (subviewID) is not specified or it is "" then a text of the first argument (view) is returned.
func GetHint(view View, subviewID ...string) string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if text, ok := stringProperty(view, Hint, view.Session()); ok {
			return text
		}
		if value := valueFromStyle(view, Hint); value != nil {
			if text, ok := value.(string); ok {
				if text, ok = view.Session().resolveConstants(text); ok {
					return text
				}
			}
		}
	}
	return ""
}

// GetMaxLength returns a maximal lenght of EditView. If a maximal lenght is not limited  then 0 is returned
// If the second argument (subviewID) is not specified or it is "" then a value of the first argument (view) is returned.
func GetMaxLength(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, MaxLength, 0)
}

// IsReadOnly returns the true if a EditView works in read only mode.
// If the second argument (subviewID) is not specified or it is "" then a value of the first argument (view) is returned.
func IsReadOnly(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, ReadOnly, false)
}

// IsSpellcheck returns a value of the Spellcheck property of EditView.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsSpellcheck(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, Spellcheck, false)
}

// GetTextChangedListeners returns the TextChangedListener list of an EditView or MultiLineEditView subview.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextChangedListeners(view View, subviewID ...string) []func(EditView, string) {
	return getEventListeners[EditView, string](view, subviewID, EditTextChangedEvent)
}

// GetEditViewType returns a value of the Type property of EditView.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetEditViewType(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, EditViewType, SingleLineText, false)
}

// GetEditViewPattern returns a value of the Pattern property of EditView.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetEditViewPattern(view View, subviewID ...string) string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if pattern, ok := stringProperty(view, EditViewPattern, view.Session()); ok {
			return pattern
		}
		if value := valueFromStyle(view, EditViewPattern); value != nil {
			if pattern, ok := value.(string); ok {
				if pattern, ok = view.Session().resolveConstants(pattern); ok {
					return pattern
				}
			}
		}
	}
	return ""
}

// IsEditViewWrap returns a value of the EditWrap property of MultiLineEditView.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsEditViewWrap(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, EditWrap, false)
}

// AppendEditText appends the text to the EditView content.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func AppendEditText(view View, subviewID string, text string) {
	if subviewID != "" {
		if edit := EditViewByID(view, subviewID); edit != nil {
			edit.AppendText(text)
			return
		}
	}

	if edit, ok := view.(EditView); ok {
		edit.AppendText(text)
	}
}

// GetCaretColor returns the color of the text input carret.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetCaretColor(view View, subviewID ...string) Color {
	return colorStyledProperty(view, subviewID, CaretColor, false)
}
