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
	textView.set = textViewSet
	textView.changed = textViewPropertyChanged
}

func textViewPropertyChanged(view View, tag PropertyName) {
	switch tag {
	case Text:
		updateInnerHTML(view.htmlID(), view.Session())

	case TextOverflow:
		session := view.Session()
		if n, ok := enumProperty(view, TextOverflow, session, 0); ok {
			values := enumProperties[TextOverflow].cssValues
			if n >= 0 && n < len(values) {
				session.updateCSSProperty(view.htmlID(), string(TextOverflow), values[n])
				return
			}
		}
		session.updateCSSProperty(view.htmlID(), string(TextOverflow), "")

	case NotTranslate:
		updateInnerHTML(view.htmlID(), view.Session())

	default:
		viewPropertyChanged(view, tag)
	}
}

func textViewSet(view View, tag PropertyName, value any) []PropertyName {
	switch tag {
	case Text:
		switch value := value.(type) {
		case string:
			view.setRaw(Text, value)

		case fmt.Stringer:
			view.setRaw(Text, value.String())

		case float32:
			view.setRaw(Text, fmt.Sprintf("%g", float64(value)))

		case float64:
			view.setRaw(Text, fmt.Sprintf("%g", value))

		case []rune:
			view.setRaw(Text, string(value))

		case bool:
			if value {
				view.setRaw(Text, "true")
			} else {
				view.setRaw(Text, "false")
			}

		default:
			if n, ok := isInt(value); ok {
				view.setRaw(Text, fmt.Sprintf("%d", n))
			} else {
				notCompatibleType(tag, value)
				return nil
			}
		}
		return []PropertyName{Text}
	}

	return viewSet(view, tag, value)
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
