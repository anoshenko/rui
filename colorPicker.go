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
	// `func(picker rui.ColorPicker, newColor, oldColor rui.Color)`.
	//
	// where:
	// picker - Interface of a color picker which generated this event,
	// newColor - New color value,
	// oldColor - Old color value.
	//
	// Allowed listener formats:
	// `func(picker rui.ColorPicker, newColor rui.Color)`,
	// `func(newColor, oldColor rui.Color)`,
	// `func(newColor rui.Color)`,
	// `func(picker rui.ColorPicker)`,
	// `func()`.
	ColorChangedEvent = "color-changed"

	// ColorPickerValue is the constant for "color-picker-value" property tag.
	//
	// Used by `ColorPicker`.
	// Define current color picker value.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	ColorPickerValue = "color-picker-value"
)

// ColorPicker represent a ColorPicker view
type ColorPicker interface {
	View
}

type colorPickerData struct {
	viewData
	dataList
	colorChangedListeners []func(ColorPicker, Color, Color)
}

// NewColorPicker create new ColorPicker object and return it
func NewColorPicker(session Session, params Params) ColorPicker {
	view := new(colorPickerData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newColorPicker(session Session) View {
	return NewColorPicker(session, nil)
}

func (picker *colorPickerData) init(session Session) {
	picker.viewData.init(session)
	picker.tag = "ColorPicker"
	picker.hasHtmlDisabled = true
	picker.colorChangedListeners = []func(ColorPicker, Color, Color){}
	picker.properties[Padding] = Px(0)
	picker.dataListInit()
}

func (picker *colorPickerData) String() string {
	return getViewString(picker, nil)
}

func (picker *colorPickerData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Value, ColorTag:
		return ColorPickerValue
	}

	return picker.normalizeDataListTag(tag)
}

func (picker *colorPickerData) Remove(tag string) {
	picker.remove(picker.normalizeTag(tag))
}

func (picker *colorPickerData) remove(tag string) {
	switch tag {
	case ColorChangedEvent:
		if len(picker.colorChangedListeners) > 0 {
			picker.colorChangedListeners = []func(ColorPicker, Color, Color){}
			picker.propertyChangedEvent(tag)
		}

	case ColorPickerValue:
		oldColor := GetColorPickerValue(picker)
		delete(picker.properties, ColorPickerValue)
		picker.colorChanged(oldColor)

	case DataList:
		if len(picker.dataList.dataList) > 0 {
			picker.setDataList(picker, []string{}, true)
		}

	default:
		picker.viewData.remove(tag)
	}
}

func (picker *colorPickerData) Set(tag string, value any) bool {
	return picker.set(picker.normalizeTag(tag), value)
}

func (picker *colorPickerData) set(tag string, value any) bool {
	if value == nil {
		picker.remove(tag)
		return true
	}

	switch tag {
	case ColorChangedEvent:
		listeners, ok := valueToEventWithOldListeners[ColorPicker, Color](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(ColorPicker, Color, Color){}
		}
		picker.colorChangedListeners = listeners
		picker.propertyChangedEvent(tag)
		return true

	case ColorPickerValue:
		oldColor := GetColorPickerValue(picker)
		if picker.setColorProperty(ColorPickerValue, value) {
			picker.colorChanged(oldColor)
			return true
		}

	case DataList:
		return picker.setDataList(picker, value, picker.created)

	default:
		return picker.viewData.set(tag, value)
	}
	return false
}

func (picker *colorPickerData) colorChanged(oldColor Color) {
	if newColor := GetColorPickerValue(picker); oldColor != newColor {
		if picker.created {
			picker.session.callFunc("setInputValue", picker.htmlID(), newColor.rgbString())
		}
		for _, listener := range picker.colorChangedListeners {
			listener(picker, newColor, oldColor)
		}
		picker.propertyChangedEvent(ColorTag)
	}
}

func (picker *colorPickerData) Get(tag string) any {
	return picker.get(picker.normalizeTag(tag))
}

func (picker *colorPickerData) get(tag string) any {
	switch tag {
	case ColorChangedEvent:
		return picker.colorChangedListeners

	case DataList:
		return picker.dataList.dataList

	default:
		return picker.viewData.get(tag)
	}
}

func (picker *colorPickerData) htmlTag() string {
	return "input"
}

func (picker *colorPickerData) htmlSubviews(self View, buffer *strings.Builder) {
	picker.dataListHtmlSubviews(self, buffer)
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

	picker.dataListHtmlProperties(picker, buffer)
}

func (picker *colorPickerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			oldColor := GetColorPickerValue(picker)
			if color, ok := StringToColor(text); ok {
				picker.properties[ColorPickerValue] = color
				if color != oldColor {
					for _, listener := range picker.colorChangedListeners {
						listener(picker, color, oldColor)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

// GetColorPickerValue returns the value of ColorPicker subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetColorPickerValue(view View, subviewID ...string) Color {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if value, ok := colorProperty(view, ColorPickerValue, view.Session()); ok {
			return value
		}
		for _, tag := range []string{ColorPickerValue, Value, ColorTag} {
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
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetColorChangedListeners(view View, subviewID ...string) []func(ColorPicker, Color, Color) {
	return getEventWithOldListeners[ColorPicker, Color](view, subviewID, ColorChangedEvent)
}
