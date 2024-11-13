package rui

import (
	"strconv"
	"strings"
)

// Constants for [EditView] specific properties and events
const (
	// EditTextChangedEvent is the constant for "edit-text-changed" property tag.
	//
	// Used by `EditView`.
	// Occur when edit view text has been changed.
	//
	// General listener format:
	// `func(editView rui.EditView, newText, oldText string)`.
	//
	// where:
	// editView - Interface of an edit view which generated this event,
	// newText - New edit view text,
	// oldText - Previous edit view text.
	//
	// Allowed listener formats:
	// `func(editView rui.EditView, newText string)`,
	// `func(newText, oldText string)`,
	// `func(newText string)`,
	// `func(editView rui.EditView)`,
	// `func()`.
	EditTextChangedEvent PropertyName = "edit-text-changed"

	// EditViewType is the constant for "edit-view-type" property tag.
	//
	// Used by `EditView`.
	// Type of the text input. Default value is "text".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`SingleLineText`) or "text" - One-line text editor.
	// `1`(`PasswordText`) or "password" - Password editor. The text is hidden by asterisks.
	// `2`(`EmailText`) or "email" - Single e-mail editor.
	// `3`(`EmailsText`) or "emails" - Multiple e-mail editor.
	// `4`(`URLText`) or "url" - Internet address input editor.
	// `5`(`PhoneText`) or "phone" - Phone number editor.
	// `6`(`MultiLineText`) or "multiline" - Multi-line text editor.
	EditViewType PropertyName = "edit-view-type"

	// EditViewPattern is the constant for "edit-view-pattern" property tag.
	//
	// Used by `EditView`.
	// Regular expression to limit editing of a text.
	//
	// Supported types: `string`.
	EditViewPattern PropertyName = "edit-view-pattern"

	// Spellcheck is the constant for "spellcheck" property tag.
	//
	// Used by `EditView`.
	// Enable or disable spell checker. Available in `SingleLineText` and `MultiLineText` types of edit view. Default value is
	// `false`.
	//
	// Supported types: `bool`, `int`, `string`.
	//
	// Values:
	// `true` or `1` or "true", "yes", "on", "1" - Enable spell checker for text.
	// `false` or `0` or "false", "no", "off", "0" - Disable spell checker for text.
	Spellcheck PropertyName = "spellcheck"
)

// Constants for the values of an [EditView] "edit-view-type" property
const (
	// SingleLineText - single-line text type of EditView
	SingleLineText = 0

	// PasswordText - password type of EditView
	PasswordText = 1

	// EmailText - e-mail type of EditView. Allows to enter one email
	EmailText = 2

	// EmailsText - e-mail type of EditView. Allows to enter multiple emails separated by comma
	EmailsText = 3

	// URLText - url type of EditView. Allows to enter one url
	URLText = 4

	// PhoneText - telephone type of EditView. Allows to enter one phone number
	PhoneText = 5

	// MultiLineText - multi-line text type of EditView
	MultiLineText = 6
)

// EditView represent an EditView view
type EditView interface {
	View

	// AppendText appends text to the current text of an EditView view
	AppendText(text string)
	textChanged(newText, oldText string)
}

type editViewData struct {
	viewData
}

// NewEditView create new EditView object and return it
func NewEditView(session Session, params Params) EditView {
	view := new(editViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newEditView(session Session) View {
	return new(editViewData) // NewEditView(session, nil)
}

func (edit *editViewData) init(session Session) {
	edit.viewData.init(session)
	edit.hasHtmlDisabled = true
	edit.tag = "EditView"
	edit.normalize = normalizeEditViewTag
	edit.set = editViewSet
	edit.changed = editViewPropertyChanged
}

func (edit *editViewData) Focusable() bool {
	return true
}

func normalizeEditViewTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
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

	return normalizeDataListTag(tag)
}

func editViewSet(view View, tag PropertyName, value any) []PropertyName {
	switch tag {
	case Text:
		if text, ok := value.(string); ok {
			old := ""
			if val := view.getRaw(Text); val != nil {
				if txt, ok := val.(string); ok {
					old = txt
				}
			}
			view.setRaw("old-text", old)
			view.setRaw(tag, text)
			return []PropertyName{tag}
		}

		notCompatibleType(tag, value)
		return nil

	case Hint:
		if text, ok := value.(string); ok {
			return setStringPropertyValue(view, tag, strings.Trim(text, " \t\n"))
		}
		notCompatibleType(tag, value)
		return nil

	case DataList:
		setDataList(view, value, "")

	case EditTextChangedEvent:
		return setEventWithOldListener[EditView, string](view, tag, value)
	}

	return viewSet(view, tag, value)
}

func editViewPropertyChanged(view View, tag PropertyName) {
	session := view.Session()

	switch tag {
	case Text:
		text := GetText(view)
		session.callFunc("setInputValue", view.htmlID(), text)

		if edit, ok := view.(EditView); ok {
			old := ""
			if val := view.getRaw("old-text"); val != nil {
				if txt, ok := val.(string); ok {
					old = txt
				}
			}
			edit.textChanged(text, old)
		}

	case Hint:
		if text := GetHint(view); text != "" {
			session.updateProperty(view.htmlID(), "placeholder", text)
		} else {
			session.removeProperty(view.htmlID(), "placeholder")
		}

	case MaxLength:
		if maxLength := GetMaxLength(view); maxLength > 0 {
			session.updateProperty(view.htmlID(), "maxlength", strconv.Itoa(maxLength))
		} else {
			session.removeProperty(view.htmlID(), "maxlength")
		}

	case ReadOnly:
		if IsReadOnly(view) {
			session.updateProperty(view.htmlID(), "readonly", "")
		} else {
			session.removeProperty(view.htmlID(), "readonly")
		}

	case Spellcheck:
		session.updateProperty(view.htmlID(), "spellcheck", IsSpellcheck(view))

	case EditViewPattern:
		if text := GetEditViewPattern(view); text != "" {
			session.updateProperty(view.htmlID(), "pattern", text)
		} else {
			session.removeProperty(view.htmlID(), "pattern")
		}

	case EditViewType:
		updateInnerHTML(view.parentHTMLID(), session)

	case EditWrap:
		if wrap := IsEditViewWrap(view); wrap {
			session.updateProperty(view.htmlID(), "wrap", "soft")
		} else {
			session.updateProperty(view.htmlID(), "wrap", "off")
		}

	case DataList:
		updateInnerHTML(view.htmlID(), session)

	default:
		viewPropertyChanged(view, tag)
	}
}

func (edit *editViewData) AppendText(text string) {
	if GetEditViewType(edit) == MultiLineText {
		if value := edit.getRaw(Text); value != nil {
			if textValue, ok := value.(string); ok {
				oldText := textValue
				textValue += text
				edit.properties[Text] = textValue
				edit.session.callFunc("appendToInnerHTML", edit.htmlID(), text)
				edit.session.callFunc("appendToInputValue", edit.htmlID(), text)
				edit.textChanged(textValue, oldText)
				return
			}
		}
		edit.setRaw(Text, text)
	} else {
		edit.setRaw(Text, GetText(edit)+text)
	}
}

func (edit *editViewData) textChanged(newText, oldText string) {
	for _, listener := range GetTextChangedListeners(edit) {
		listener(edit, newText, oldText)
	}
	if listener, ok := edit.changeListener[Text]; ok {
		listener(edit, Text)
	}
}

func (edit *editViewData) htmlTag() string {
	if GetEditViewType(edit) == MultiLineText {
		return "textarea"
	}
	return "input"
}

func (edit *editViewData) htmlSubviews(self View, buffer *strings.Builder) {
	if GetEditViewType(edit) == MultiLineText {
		if text := GetText(edit); text != "" {
			buffer.WriteString(text)
		}
	}
	dataListHtmlSubviews(self, buffer, func(text string, session Session) string {
		return text
	})
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
		if strings.ContainsRune(text, '\n') {
			text = strings.ReplaceAll(text, "\n", `\n`)
		}
		return text
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

	dataListHtmlProperties(edit, buffer)
}

func (edit *editViewData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "textChanged":
		oldText := GetText(edit)
		if text, ok := data.PropertyValue("text"); ok {
			edit.setRaw(Text, text)
			if text != oldText {
				edit.textChanged(text, oldText)
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

	session := view.Session()
	text := ""
	if view != nil {
		var ok bool
		text, ok = stringProperty(view, Hint, view.Session())
		if !ok {
			if value := valueFromStyle(view, Hint); value != nil {
				if text, ok = value.(string); ok {
					if text, ok = session.resolveConstants(text); !ok {
						text = ""
					}
				} else {
					text = ""
				}
			}
		}
	}

	if text != "" && !GetNotTranslate(view) {
		text, _ = session.GetString(text)
	}

	return text
}

// GetMaxLength returns a maximal length of EditView. If a maximal length is not limited  then 0 is returned
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
func GetTextChangedListeners(view View, subviewID ...string) []func(EditView, string, string) {
	return getEventWithOldListeners[EditView, string](view, subviewID, EditTextChangedEvent)
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

// GetCaretColor returns the color of the text input caret.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetCaretColor(view View, subviewID ...string) Color {
	return colorStyledProperty(view, subviewID, CaretColor, false)
}
