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
	return new(textViewData)
}

// Init initialize fields of TextView by default values
func (textView *textViewData) init(session Session) {
	textView.viewData.init(session)
	textView.tag = "TextView"
	textView.set = textView.setFunc
	textView.changed = textView.propertyChanged
}

func (textView *textViewData) propertyChanged(tag PropertyName) {
	switch tag {
	case Text:
		updateInnerHTML(textView.htmlID(), textView.Session())

	case TextOverflow:
		session := textView.Session()
		if n, ok := enumProperty(textView, TextOverflow, session, 0); ok {
			values := enumProperties[TextOverflow].cssValues
			if n >= 0 && n < len(values) {
				session.updateCSSProperty(textView.htmlID(), string(TextOverflow), values[n])
				return
			}
		}
		session.updateCSSProperty(textView.htmlID(), string(TextOverflow), "")

	case NotTranslate:
		updateInnerHTML(textView.htmlID(), textView.Session())

	default:
		textView.viewData.propertyChanged(tag)
	}
}

func (textView *textViewData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Text:
		switch value := value.(type) {
		case string:
			textView.setRaw(Text, value)

		case fmt.Stringer:
			textView.setRaw(Text, value.String())

		case float32:
			textView.setRaw(Text, fmt.Sprintf("%g", float64(value)))

		case float64:
			textView.setRaw(Text, fmt.Sprintf("%g", value))

		case []rune:
			textView.setRaw(Text, string(value))

		case bool:
			if value {
				textView.setRaw(Text, "true")
			} else {
				textView.setRaw(Text, "false")
			}

		default:
			if n, ok := isInt(value); ok {
				textView.setRaw(Text, fmt.Sprintf("%d", n))
			} else {
				notCompatibleType(tag, value)
				return nil
			}
		}
		return []PropertyName{Text}
	}

	return textView.viewData.setFunc(tag, value)
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
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetTextOverflow(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, TextOverflow, SingleLineText, false)
}
