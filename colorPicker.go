package rui

import (
	"fmt"
	"strings"
)

const (
	ColorChangedEvent = "color-changed"
	ColorPickerValue  = "color-picker-value"
)

// ColorPicker - ColorPicker view
type ColorPicker interface {
	View
}

type colorPickerData struct {
	viewData
	colorChangedListeners []func(ColorPicker, Color)
}

// NewColorPicker create new ColorPicker object and return it
func NewColorPicker(session Session, params Params) ColorPicker {
	view := new(colorPickerData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newColorPicker(session Session) View {
	return NewColorPicker(session, nil)
}

func (picker *colorPickerData) Init(session Session) {
	picker.viewData.Init(session)
	picker.tag = "ColorPicker"
	picker.colorChangedListeners = []func(ColorPicker, Color){}
	picker.properties[Padding] = Px(0)
}

func (picker *colorPickerData) String() string {
	return getViewString(picker)
}

func (picker *colorPickerData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Value, ColorTag:
		return ColorPickerValue
	}

	return tag
}

func (picker *colorPickerData) Remove(tag string) {
	picker.remove(picker.normalizeTag(tag))
}

func (picker *colorPickerData) remove(tag string) {
	switch tag {
	case ColorChangedEvent:
		if len(picker.colorChangedListeners) > 0 {
			picker.colorChangedListeners = []func(ColorPicker, Color){}
			picker.propertyChangedEvent(tag)
		}

	case ColorPickerValue:
		oldColor := GetColorPickerValue(picker, "")
		delete(picker.properties, ColorPickerValue)
		picker.colorChanged(oldColor)

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
		listeners, ok := valueToEventListeners[ColorPicker, Color](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(ColorPicker, Color){}
		}
		picker.colorChangedListeners = listeners
		picker.propertyChangedEvent(tag)
		return true

	case ColorPickerValue:
		oldColor := GetColorPickerValue(picker, "")
		if picker.setColorProperty(ColorPickerValue, value) {
			picker.colorChanged(oldColor)
			return true
		}

	default:
		return picker.viewData.set(tag, value)
	}
	return false
}

func (picker *colorPickerData) colorChanged(oldColor Color) {
	if newColor := GetColorPickerValue(picker, ""); oldColor != newColor {
		if picker.created {
			picker.session.runScript(fmt.Sprintf(`setInputValue('%s', '%s')`, picker.htmlID(), newColor.rgbString()))
		}
		for _, listener := range picker.colorChangedListeners {
			listener(picker, newColor)
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

	default:
		return picker.viewData.get(tag)
	}
}

func (picker *colorPickerData) htmlTag() string {
	return "input"
}

func (picker *colorPickerData) htmlProperties(self View, buffer *strings.Builder) {
	picker.viewData.htmlProperties(self, buffer)

	buffer.WriteString(` type="color" value="`)
	buffer.WriteString(GetColorPickerValue(picker, "").rgbString())
	buffer.WriteByte('"')

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)
	if picker.getRaw(ClickEvent) == nil {
		buffer.WriteString(` onclick="stopEventPropagation(this, event)"`)
	}
}

func (picker *colorPickerData) htmlDisabledProperties(self View, buffer *strings.Builder) {
	if IsDisabled(self, "") {
		buffer.WriteString(` disabled`)
	}
	picker.viewData.htmlDisabledProperties(self, buffer)
}

func (picker *colorPickerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			oldColor := GetColorPickerValue(picker, "")
			if color, ok := StringToColor(text); ok {
				picker.properties[ColorPickerValue] = color
				if color != oldColor {
					for _, listener := range picker.colorChangedListeners {
						listener(picker, color)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

// GetColorPickerValue returns the value of ColorPicker subview.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetColorPickerValue(view View, subviewID string) Color {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if result, ok := colorStyledProperty(view, ColorPickerValue); ok {
			return result
		}
		for _, tag := range []string{Value, ColorTag} {
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
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetColorChangedListeners(view View, subviewID string) []func(ColorPicker, Color) {
	return getEventListeners[ColorPicker, Color](view, subviewID, ColorChangedEvent)
}
