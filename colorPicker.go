package rui

import (
	"strings"
)

// Constants for [ColorPicker] specific properties and events.
const (
	// ColorChangedEvent is the constant for "color-changed" property tag.
	//
	// Used by `ColorPicker`.
	// Event generated when color picker value has been changed.
	//
	// General listener format:
	//  func(picker rui.ColorPicker, newColor, oldColor rui.Color)
	//
	// where:
	//   - picker - Interface of a color picker which generated this event,
	//   - newColor - New color value,
	//   - oldColor - Old color value.
	//
	// Allowed listener formats:
	//  func(picker rui.ColorPicker, newColor rui.Color)
	//  func(newColor, oldColor rui.Color)
	//  func(newColor rui.Color)
	//  func(picker rui.ColorPicker)
	//  func()
	ColorChangedEvent PropertyName = "color-changed"

	// ColorPickerValue is the constant for "color-picker-value" property tag.
	//
	// Used by `ColorPicker`.
	// Define current color picker value.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	ColorPickerValue PropertyName = "color-picker-value"
)

// ColorPicker represent a ColorPicker view
type ColorPicker interface {
	View
}

type colorPickerData struct {
	viewData
}

// NewColorPicker create new ColorPicker object and return it
func NewColorPicker(session Session, params Params) ColorPicker {
	view := new(colorPickerData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newColorPicker(session Session) View {
	return new(colorPickerData)
}

func (picker *colorPickerData) init(session Session) {
	picker.viewData.init(session)
	picker.tag = "ColorPicker"
	picker.hasHtmlDisabled = true
	picker.properties[Padding] = Px(0)
	picker.normalize = normalizeColorPickerTag
	picker.set = picker.setFunc
	picker.changed = picker.propertyChanged
}

func normalizeColorPickerTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Value, ColorTag:
		return ColorPickerValue
	}

	return normalizeDataListTag(tag)
}

func (picker *colorPickerData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case ColorChangedEvent:
		return setTwoArgEventListener[ColorPicker, Color](picker, tag, value)

	case ColorPickerValue:
		oldColor := GetColorPickerValue(picker)
		result := setColorProperty(picker, ColorPickerValue, value)
		if result != nil {
			picker.setRaw("old-color", oldColor)
		}
		return result

	case DataList:
		return setDataList(picker, value, "")
	}

	return picker.viewData.setFunc(tag, value)
}

func (picker *colorPickerData) propertyChanged(tag PropertyName) {
	switch tag {
	case ColorPickerValue:
		color := GetColorPickerValue(picker)
		picker.Session().callFunc("setInputValue", picker.htmlID(), color.rgbString())

		if listeners := getTwoArgEventListeners[ColorPicker, Color](picker, nil, ColorChangedEvent); len(listeners) > 0 {
			oldColor := Color(0)
			if value := picker.getRaw("old-color"); value != nil {
				oldColor = value.(Color)
			}
			for _, listener := range listeners {
				listener.Run(picker, color, oldColor)
			}
		}

	default:
		picker.viewData.propertyChanged(tag)
	}

}

func (picker *colorPickerData) htmlTag() string {
	return "input"
}

func (picker *colorPickerData) htmlSubviews(self View, buffer *strings.Builder) {
	dataListHtmlSubviews(self, buffer, func(text string, session Session) string {
		text, _ = session.resolveConstants(text)
		return text
	})
}

func (picker *colorPickerData) htmlProperties(self View, buffer *strings.Builder) {
	picker.viewData.htmlProperties(self, buffer)

	buffer.WriteString(` type="color" value="`)
	buffer.WriteString(GetColorPickerValue(picker).rgbString())
	buffer.WriteByte('"')

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)
	if picker.getRaw(ClickEvent) == nil {
		buffer.WriteString(` onclick="stopEventPropagation(this, event)"`)
	}

	dataListHtmlProperties(picker, buffer)
}

func (picker *colorPickerData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if color, ok := StringToColor(text); ok {
				oldColor := GetColorPickerValue(picker)
				picker.properties[ColorPickerValue] = color
				if color != oldColor {
					for _, listener := range getTwoArgEventListeners[ColorPicker, Color](picker, nil, ColorChangedEvent) {
						listener.Run(picker, color, oldColor)
					}
					if listener, ok := picker.changeListener[ColorPickerValue]; ok {
						listener(picker, ColorPickerValue)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

// GetColorPickerValue returns the value of ColorPicker subview.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetColorPickerValue(view View, subviewID ...string) Color {
	if view = getSubview(view, subviewID); view != nil {
		if value, ok := colorProperty(view, ColorPickerValue, view.Session()); ok {
			return value
		}
		for _, tag := range []PropertyName{ColorPickerValue, Value, ColorTag} {
			if value := valueFromStyle(view, tag); value != nil {
				if result, ok := valueToColor(value, view.Session()); ok {
					return result
				}
			}
		}
	}
	return 0
}

// GetColorChangedListeners returns the ColorChangedListener list of an ColorPicker subview.
// If there are no listeners then the empty list is returned
//
// Result elements can be of the following types:
//   - func(rui.ColorPicker, rui.Color, rui.Color),
//   - func(rui.ColorPicker, rui.Color),
//   - func(rui.ColorPicker),
//   - func(rui.Color, rui.Color),
//   - func(rui.Color),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetColorChangedListeners(view View, subviewID ...string) []any {
	return getTwoArgEventRawListeners[ColorPicker, Color](view, subviewID, ColorChangedEvent)
}
