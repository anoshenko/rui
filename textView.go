package rui

import (
	"fmt"
	"strings"
)

// TextView represents a TextView view
type TextView interface {
	View
}

type textViewData struct {
	viewData
}

// NewTextView create new TextView object and return it
func NewTextView(session Session, params Params) TextView {
	view := new(textViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newTextView(session Session) View {
	return NewTextView(session, nil)
}

// Init initialize fields of TextView by default values
func (textView *textViewData) init(session Session) {
	textView.viewData.init(session)
	textView.tag = "TextView"
}

func (textView *textViewData) String() string {
	return getViewString(textView, nil)
}

func (textView *textViewData) Get(tag string) any {
	return textView.get(strings.ToLower(tag))
}

func (textView *textViewData) Remove(tag string) {
	textView.remove(strings.ToLower(tag))
}

func (textView *textViewData) remove(tag string) {
	textView.viewData.remove(tag)
	if textView.created {
		switch tag {
		case Text:
			updateInnerHTML(textView.htmlID(), textView.session)

		case TextOverflow:
			textView.textOverflowUpdated()
		}
	}
}

func (textView *textViewData) Set(tag string, value any) bool {
	return textView.set(strings.ToLower(tag), value)
}

func (textView *textViewData) set(tag string, value any) bool {
	switch tag {
	case Text:
		switch value := value.(type) {
		case string:
			textView.properties[Text] = value

		case fmt.Stringer:
			textView.properties[Text] = value.String()

		case float32:
			textView.properties[Text] = fmt.Sprintf("%g", float64(value))

		case float64:
			textView.properties[Text] = fmt.Sprintf("%g", value)

		case []rune:
			textView.properties[Text] = string(value)

		case bool:
			if value {
				textView.properties[Text] = "true"
			} else {
				textView.properties[Text] = "false"
			}

		default:
			if n, ok := isInt(value); ok {
				textView.properties[Text] = fmt.Sprintf("%d", n)
			} else {
				notCompatibleType(tag, value)
				return false
			}
		}
		if textView.created {
			updateInnerHTML(textView.htmlID(), textView.session)
		}

	case TextOverflow:
		if !textView.viewData.set(tag, value) {
			return false
		}
		if textView.created {
			textView.textOverflowUpdated()
		}

	case NotTranslate:
		if !textView.viewData.set(tag, value) {
			return false
		}
		if textView.created {
			updateInnerHTML(textView.htmlID(), textView.Session())
		}

	default:
		return textView.viewData.set(tag, value)
	}

	textView.propertyChangedEvent(tag)
	return true
}

func (textView *textViewData) textOverflowUpdated() {
	session := textView.Session()
	if n, ok := enumProperty(textView, TextOverflow, session, 0); ok {
		values := enumProperties[TextOverflow].cssValues
		if n >= 0 && n < len(values) {
			session.updateCSSProperty(textView.htmlID(), TextOverflow, values[n])
			return
		}
	}
	session.updateCSSProperty(textView.htmlID(), TextOverflow, "")
}

func (textView *textViewData) htmlSubviews(self View, buffer *strings.Builder) {
	if value := textView.getRaw(Text); value != nil {
		if text, ok := value.(string); ok {
			if !GetNotTranslate(textView) {
				text, _ = textView.session.GetString(text)
			}
			buffer.WriteString(text)
		}
	}
}

// GetTextOverflow returns a value of the "text-overflow" property:
// TextOverflowClip (0) or TextOverflowEllipsis (1).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTextOverflow(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, TextOverflow, SingleLineText, false)
}
