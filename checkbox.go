package rui

import (
	"strings"
)

// CheckboxChangedEvent is the constant for "checkbox-event" property tag.
//
// Used by `Checkbox`.
// Event occurs when the checkbox becomes checked/unchecked.
//
// General listener format:
// `func(checkbox rui.Checkbox, checked bool)`.
//
// where:
// checkbox - Interface of a checkbox which generated this event,
// checked - Checkbox state.
//
// Allowed listener formats:
// `func(checkbox rui.Checkbox)`,
// `func(checked bool)`,
// `func()`.
const CheckboxChangedEvent PropertyName = "checkbox-event"

// Checkbox represent a Checkbox view
type Checkbox interface {
	ViewsContainer
}

type checkboxData struct {
	viewsContainerData
}

// NewCheckbox create new Checkbox object and return it
func NewCheckbox(session Session, params Params) Checkbox {
	view := new(checkboxData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newCheckbox(session Session) View {
	return new(checkboxData)
}

func (button *checkboxData) init(session Session) {
	button.viewsContainerData.init(session)
	button.tag = "Checkbox"
	button.systemClass = "ruiGridLayout ruiCheckbox"
	button.set = button.setFunc
	button.remove = button.removeFunc
	button.changed = checkboxPropertyChanged

	button.setRaw(ClickEvent, checkboxClickListener)
	button.setRaw(KeyDownEvent, checkboxKeyListener)
}

func (button *checkboxData) Focusable() bool {
	return true
}

func checkboxPropertyChanged(view View, tag PropertyName) {
	switch tag {

	case Checked:
		session := view.Session()
		checked := IsCheckboxChecked(view)
		if listeners := GetCheckboxChangedListeners(view); len(listeners) > 0 {
			if checkbox, ok := view.(Checkbox); ok {
				for _, listener := range listeners {
					listener(checkbox, checked)
				}
			}
		}

		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)

		checkboxHtml(view, buffer, checked)
		session.updateInnerHTML(view.htmlID()+"checkbox", buffer.String())

	case CheckboxHorizontalAlign, CheckboxVerticalAlign:
		htmlID := view.htmlID()
		session := view.Session()
		updateCSSStyle(htmlID, session)
		updateInnerHTML(htmlID, session)

	case VerticalAlign:
		view.Session().updateCSSProperty(view.htmlID()+"content", "align-items", checkboxVerticalAlignCSS(view))

	case HorizontalAlign:
		view.Session().updateCSSProperty(view.htmlID()+"content", "justify-items", checkboxHorizontalAlignCSS(view))

	case AccentColor:
		updateInnerHTML(view.htmlID(), view.Session())

	default:
		viewsContainerPropertyChanged(view, tag)
	}
}

func (button *checkboxData) setFunc(view View, tag PropertyName, value any) []PropertyName {
	switch tag {
	case ClickEvent:
		if button.viewsContainerData.setFunc(view, ClickEvent, value) != nil {
			if value := view.getRaw(ClickEvent); value != nil {
				if listeners, ok := value.([]func(View, MouseEvent)); ok {
					listeners = append(listeners, checkboxClickListener)
					view.setRaw(ClickEvent, listeners)
					return []PropertyName{ClickEvent}
				}
			}

			return button.viewsContainerData.setFunc(view, ClickEvent, checkboxClickListener)
		}
		return nil

	case KeyDownEvent:
		if button.viewsContainerData.setFunc(view, KeyDownEvent, value) != nil {
			if value := view.getRaw(KeyDownEvent); value != nil {
				if listeners, ok := value.([]func(View, KeyEvent)); ok {
					listeners = append(listeners, checkboxKeyListener)
					view.setRaw(KeyDownEvent, listeners)
					return []PropertyName{KeyDownEvent}
				}
			}

			return button.viewsContainerData.setFunc(view, KeyDownEvent, checkboxKeyListener)
		}
		return nil

	case CheckboxChangedEvent:
		return setViewEventListener[Checkbox, bool](view, tag, value)

	case Checked:
		return setBoolProperty(view, Checked, value)

	case CellVerticalAlign, CellHorizontalAlign, CellWidth, CellHeight:
		ErrorLogF(`"%s" property is not compatible with the BoundsProperty`, string(tag))
		return nil
	}

	return button.viewsContainerData.setFunc(view, tag, value)
}

func (button *checkboxData) removeFunc(view View, tag PropertyName) []PropertyName {
	switch tag {
	case ClickEvent:
		button.setRaw(ClickEvent, checkboxClickListener)
		return []PropertyName{ClickEvent}

	case KeyDownEvent:
		button.setRaw(KeyDownEvent, checkboxKeyListener)
		return []PropertyName{ClickEvent}
	}

	return button.viewsContainerData.removeFunc(view, tag)
}

func (button *checkboxData) checked() bool {
	checked, _ := boolProperty(button, Checked, button.Session())
	return checked
}

/*
func (button *checkboxData) changedCheckboxState(state bool) {
	for _, listener := range GetCheckboxChangedListeners(button) {
		listener(button, state)
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	button.htmlCheckbox(buffer, state)
	button.Session().updateInnerHTML(button.htmlID()+"checkbox", buffer.String())
}
*/

func checkboxClickListener(view View, _ MouseEvent) {
	view.Set(Checked, !IsCheckboxChecked(view))
	BlurView(view)
}

func checkboxKeyListener(view View, event KeyEvent) {
	switch event.Code {
	case "Enter", "Space":
		view.Set(Checked, !IsCheckboxChecked(view))
	}
}

func (button *checkboxData) cssStyle(self View, builder cssBuilder) {
	session := button.Session()
	vAlign := GetCheckboxVerticalAlign(button)
	hAlign := GetCheckboxHorizontalAlign(button)
	switch hAlign {
	case CenterAlign:
		if vAlign == BottomAlign {
			builder.add("grid-template-rows", "1fr auto")
		} else {
			builder.add("grid-template-rows", "auto 1fr")
		}

	case RightAlign:
		builder.add("grid-template-columns", "1fr auto")

	default:
		builder.add("grid-template-columns", "auto 1fr")
	}

	if gap, ok := sizeConstant(session, "ruiCheckboxGap"); ok && gap.Type != Auto && gap.Value > 0 {
		builder.add("gap", gap.cssString("0", session))
	}

	builder.add("align-items", "stretch")
	builder.add("justify-items", "stretch")

	button.viewsContainerData.cssStyle(self, builder)
}

func checkboxHtml(button View, buffer *strings.Builder, checked bool) (int, int) {
	//func (button *checkboxData) htmlCheckbox(buffer *strings.Builder, checked bool) (int, int) {
	vAlign := GetCheckboxVerticalAlign(button)
	hAlign := GetCheckboxHorizontalAlign(button)

	buffer.WriteString(`<div id="`)
	buffer.WriteString(button.htmlID())
	buffer.WriteString(`checkbox" style="display: grid;`)
	if hAlign == CenterAlign {
		buffer.WriteString(" justify-items: center; grid-column-start: 1; grid-column-end: 2;")
		if vAlign == BottomAlign {
			buffer.WriteString(" grid-row-start: 2; grid-row-end: 3;")
		} else {
			buffer.WriteString(" grid-row-start: 1; grid-row-end: 2;")
		}
	} else {
		if hAlign == RightAlign {
			buffer.WriteString(" grid-column-start: 2; grid-column-end: 3;")
		} else {
			buffer.WriteString(" grid-column-start: 1; grid-column-end: 2;")
		}
		buffer.WriteString(" grid-row-start: 1; grid-row-end: 2;")
		switch vAlign {
		case BottomAlign:
			buffer.WriteString(" align-items: end;")

		case CenterAlign:
			buffer.WriteString(" align-items: center;")

		default:
			buffer.WriteString(" align-items: start;")
		}
	}

	buffer.WriteString(`">`)

	accentColor := Color(0)
	if color := GetAccentColor(button, ""); color != 0 {
		accentColor = color
	}

	if checked {
		buffer.WriteString(button.Session().checkboxOnImage(accentColor))
	} else {
		buffer.WriteString(button.Session().checkboxOffImage(accentColor))
	}
	buffer.WriteString(`</div>`)

	return vAlign, hAlign
}

func (button *checkboxData) htmlSubviews(self View, buffer *strings.Builder) {

	vCheckboxAlign, hCheckboxAlign := checkboxHtml(button, buffer, IsCheckboxChecked(button))

	buffer.WriteString(`<div id="`)
	buffer.WriteString(button.htmlID())
	buffer.WriteString(`content" style="display: grid;`)
	if hCheckboxAlign == LeftAlign {
		buffer.WriteString(" grid-column-start: 2; grid-column-end: 3;")
	} else {
		buffer.WriteString(" grid-column-start: 1; grid-column-end: 2;")
	}

	if hCheckboxAlign == CenterAlign && vCheckboxAlign != BottomAlign {
		buffer.WriteString(" grid-row-start: 2; grid-row-end: 3;")
	} else {
		buffer.WriteString(" grid-row-start: 1; grid-row-end: 2;")
	}

	buffer.WriteString(" align-items: ")
	buffer.WriteString(checkboxVerticalAlignCSS(button))
	buffer.WriteRune(';')

	buffer.WriteString(" justify-items: ")
	buffer.WriteString(checkboxHorizontalAlignCSS(button))
	buffer.WriteRune(';')

	buffer.WriteString(`">`)
	button.viewsContainerData.htmlSubviews(self, buffer)
	buffer.WriteString(`</div>`)
}

func checkboxHorizontalAlignCSS(view View) string {
	align := GetHorizontalAlign(view)
	values := enumProperties[CellHorizontalAlign].cssValues
	if align >= 0 && align < len(values) {
		return values[align]
	}
	return values[0]
}

func checkboxVerticalAlignCSS(view View) string {
	align := GetVerticalAlign(view)
	values := enumProperties[CellVerticalAlign].cssValues
	if align >= 0 && align < len(values) {
		return values[align]
	}
	return values[0]
}

// IsCheckboxChecked returns true if the Checkbox is checked, false otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsCheckboxChecked(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, Checked, false)
}

// GetCheckboxVerticalAlign return the vertical align of a Checkbox subview: TopAlign (0), BottomAlign (1), CenterAlign (2)
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetCheckboxVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CheckboxVerticalAlign, LeftAlign, false)
}

// GetCheckboxHorizontalAlign return the vertical align of a Checkbox subview: LeftAlign (0), RightAlign (1), CenterAlign (2)
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetCheckboxHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CheckboxHorizontalAlign, TopAlign, false)
}

// GetCheckboxChangedListeners returns the CheckboxChangedListener list of an Checkbox subview.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetCheckboxChangedListeners(view View, subviewID ...string) []func(Checkbox, bool) {
	return getEventListeners[Checkbox, bool](view, subviewID, CheckboxChangedEvent)
}
