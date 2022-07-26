package rui

import (
	"fmt"
	"strings"
)

// CheckboxChangedEvent is the constant for "checkbox-event" property tag.
// The "checkbox-event" event occurs when the checkbox becomes checked/unchecked.
// The main listener format: func(Checkbox, bool), where the second argument is the checkbox state.
const CheckboxChangedEvent = "checkbox-event"

// Checkbox - checkbox view
type Checkbox interface {
	ViewsContainer
}

type checkboxData struct {
	viewsContainerData
	checkedListeners []func(Checkbox, bool)
}

// NewCheckbox create new Checkbox object and return it
func NewCheckbox(session Session, params Params) Checkbox {
	view := new(checkboxData)
	view.Init(session)
	setInitParams(view, Params{
		ClickEvent:   checkboxClickListener,
		KeyDownEvent: checkboxKeyListener,
	})
	setInitParams(view, params)
	return view
}

func newCheckbox(session Session) View {
	return NewCheckbox(session, nil)
}

func (button *checkboxData) Init(session Session) {
	button.viewsContainerData.Init(session)
	button.tag = "Checkbox"
	button.systemClass = "ruiGridLayout ruiCheckbox"
	button.checkedListeners = []func(Checkbox, bool){}
}

func (button *checkboxData) String() string {
	return getViewString(button)
}

func (button *checkboxData) Focusable() bool {
	return true
}

func (button *checkboxData) Get(tag string) any {
	switch strings.ToLower(tag) {
	case CheckboxChangedEvent:
		return button.checkedListeners
	}

	return button.viewsContainerData.Get(tag)
}

func (button *checkboxData) Set(tag string, value any) bool {
	return button.set(tag, value)
}

func (button *checkboxData) set(tag string, value any) bool {
	switch tag {
	case CheckboxChangedEvent:
		if !button.setChangedListener(value) {
			notCompatibleType(tag, value)
			return false
		}

	case Checked:
		oldChecked := button.checked()
		if !button.setBoolProperty(Checked, value) {
			return false
		}
		if button.created {
			checked := button.checked()
			if checked != oldChecked {
				button.changedCheckboxState(checked)
			}
		}

	case CheckboxHorizontalAlign, CheckboxVerticalAlign:
		if !button.setEnumProperty(tag, value, enumProperties[tag].values) {
			return false
		}
		if button.created {
			htmlID := button.htmlID()
			updateCSSStyle(htmlID, button.session)
			updateInnerHTML(htmlID, button.session)
		}

	case VerticalAlign:
		if !button.setEnumProperty(tag, value, enumProperties[tag].values) {
			return false
		}
		if button.created {
			updateCSSProperty(button.htmlID()+"content", "align-items", button.cssVerticalAlign(), button.session)
		}

	case HorizontalAlign:
		if !button.setEnumProperty(tag, value, enumProperties[tag].values) {
			return false
		}
		if button.created {
			updateCSSProperty(button.htmlID()+"content", "justify-items", button.cssHorizontalAlign(), button.session)
		}

	case CellVerticalAlign, CellHorizontalAlign, CellWidth, CellHeight:
		return false

	default:
		return button.viewsContainerData.set(tag, value)
	}

	button.propertyChangedEvent(tag)
	return true
}

func (button *checkboxData) Remove(tag string) {
	button.remove(strings.ToLower(tag))
}

func (button *checkboxData) remove(tag string) {
	switch tag {
	case ClickEvent:
		if !button.viewsContainerData.set(ClickEvent, checkboxClickListener) {
			delete(button.properties, tag)
		}

	case KeyDownEvent:
		if !button.viewsContainerData.set(KeyDownEvent, checkboxKeyListener) {
			delete(button.properties, tag)
		}

	case CheckboxChangedEvent:
		if len(button.checkedListeners) > 0 {
			button.checkedListeners = []func(Checkbox, bool){}
		}

	case Checked:
		oldChecked := button.checked()
		delete(button.properties, tag)
		if button.created && oldChecked {
			button.changedCheckboxState(false)
		}

	case CheckboxHorizontalAlign, CheckboxVerticalAlign:
		delete(button.properties, tag)
		if button.created {
			htmlID := button.htmlID()
			updateCSSStyle(htmlID, button.session)
			updateInnerHTML(htmlID, button.session)
		}

	case VerticalAlign:
		delete(button.properties, tag)
		if button.created {
			updateCSSProperty(button.htmlID()+"content", "align-items", button.cssVerticalAlign(), button.session)
		}

	case HorizontalAlign:
		delete(button.properties, tag)
		if button.created {
			updateCSSProperty(button.htmlID()+"content", "justify-items", button.cssHorizontalAlign(), button.session)
		}

	default:
		button.viewsContainerData.remove(tag)
		return
	}
	button.propertyChangedEvent(tag)
}

func (button *checkboxData) checked() bool {
	checked, _ := boolProperty(button, Checked, button.Session())
	return checked
}

func (button *checkboxData) changedCheckboxState(state bool) {
	for _, listener := range button.checkedListeners {
		listener(button, state)
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	button.htmlCheckbox(buffer, state)
	button.Session().runScript(fmt.Sprintf(`updateInnerHTML('%v', '%v');`, button.htmlID()+"checkbox", buffer.String()))
}

func checkboxClickListener(view View) {
	view.Set(Checked, !IsCheckboxChecked(view, ""))
	BlurView(view)
}

func checkboxKeyListener(view View, event KeyEvent) {
	switch event.Code {
	case "Enter", "Space":
		view.Set(Checked, !IsCheckboxChecked(view, ""))
	}
}

func (button *checkboxData) setChangedListener(value any) bool {
	if value == nil {
		if len(button.checkedListeners) > 0 {
			button.checkedListeners = []func(Checkbox, bool){}
		}
		return true
	}

	switch value := value.(type) {
	case func(Checkbox, bool):
		button.checkedListeners = []func(Checkbox, bool){value}

	case func(bool):
		fn := func(_ Checkbox, checked bool) {
			value(checked)
		}
		button.checkedListeners = []func(Checkbox, bool){fn}

	case []func(Checkbox, bool):
		button.checkedListeners = value

	case []func(bool):
		listeners := make([]func(Checkbox, bool), len(value))
		for i, val := range value {
			if val == nil {
				return false
			}

			listeners[i] = func(_ Checkbox, checked bool) {
				val(checked)
			}
		}
		button.checkedListeners = listeners

	case []any:
		listeners := make([]func(Checkbox, bool), len(value))
		for i, val := range value {
			if val == nil {
				return false
			}

			switch val := val.(type) {
			case func(Checkbox, bool):
				listeners[i] = val

			case func(bool):
				listeners[i] = func(_ Checkbox, date bool) {
					val(date)
				}

			default:
				return false
			}
		}
		button.checkedListeners = listeners
	}
	return true
}

func (button *checkboxData) cssStyle(self View, builder cssBuilder) {
	session := button.Session()
	vAlign, _ := enumStyledProperty(button, CheckboxVerticalAlign, LeftAlign)
	hAlign, _ := enumStyledProperty(button, CheckboxHorizontalAlign, TopAlign)
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
		builder.add("gap", gap.cssString("0"))
	}

	builder.add("align-items", "stretch")
	builder.add("justify-items", "stretch")

	button.viewsContainerData.cssStyle(self, builder)
}

func (button *checkboxData) htmlCheckbox(buffer *strings.Builder, checked bool) (int, int) {
	vAlign, _ := enumStyledProperty(button, CheckboxVerticalAlign, LeftAlign)
	hAlign, _ := enumStyledProperty(button, CheckboxHorizontalAlign, TopAlign)

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
	if checked {
		buffer.WriteString(button.Session().checkboxOnImage())
	} else {
		buffer.WriteString(button.Session().checkboxOffImage())
	}
	buffer.WriteString(`</div>`)

	return vAlign, hAlign
}

func (button *checkboxData) htmlSubviews(self View, buffer *strings.Builder) {

	vCheckboxAlign, hCheckboxAlign := button.htmlCheckbox(buffer, IsCheckboxChecked(button, ""))

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
	buffer.WriteString(button.cssVerticalAlign())
	buffer.WriteRune(';')

	buffer.WriteString(" justify-items: ")
	buffer.WriteString(button.cssHorizontalAlign())
	buffer.WriteRune(';')

	buffer.WriteString(`">`)
	button.viewsContainerData.htmlSubviews(self, buffer)
	buffer.WriteString(`</div>`)
}

func (button *checkboxData) cssHorizontalAlign() string {
	align, _ := enumStyledProperty(button, HorizontalAlign, TopAlign)
	values := enumProperties[CellHorizontalAlign].cssValues
	if align >= 0 && align < len(values) {
		return values[align]
	}
	return values[0]
}

func (button *checkboxData) cssVerticalAlign() string {
	align, _ := enumStyledProperty(button, VerticalAlign, TopAlign)
	values := enumProperties[CellVerticalAlign].cssValues
	if align >= 0 && align < len(values) {
		return values[align]
	}
	return values[0]
}

// IsCheckboxChecked returns true if the Checkbox is checked, false otherwise.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func IsCheckboxChecked(view View, subviewID string) bool {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if checked := view.Get(Checked); checked != nil {
			if b, ok := checked.(bool); ok {
				return b
			}
		}
	}
	return false
}
