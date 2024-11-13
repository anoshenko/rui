package rui

import (
	"math"
	"strconv"
	"strings"
)

// Constants related to [NumberPicker] specific properties and events
const (
	// NumberChangedEvent is the constant for "number-changed" property tag.
	//
	// Used by `NumberPicker`.
	// Set listener(s) that track the change in the entered value.
	//
	// General listener format:
	// `func(picker rui.NumberPicker, newValue, oldValue float64)`.
	//
	// where:
	// picker - Interface of a number picker which generated this event,
	// newValue - New value,
	// oldValue - Old Value.
	//
	// Allowed listener formats:
	// `func(picker rui.NumberPicker, newValue float64)`,
	// `func(newValue, oldValue float64)`,
	// `func(newValue float64)`,
	// `func()`.
	NumberChangedEvent PropertyName = "number-changed"

	// NumberPickerType is the constant for "number-picker-type" property tag.
	//
	// Used by `NumberPicker`.
	// Sets the visual representation.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`NumberEditor`) or "editor" - Displayed as an editor.
	// `1`(`NumberSlider`) or "slider" - Displayed as a slider.
	NumberPickerType PropertyName = "number-picker-type"

	// NumberPickerMin is the constant for "number-picker-min" property tag.
	//
	// Used by `NumberPicker`.
	// Set the minimum value. The default value is 0.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	NumberPickerMin PropertyName = "number-picker-min"

	// NumberPickerMax is the constant for "number-picker-max" property tag.
	//
	// Used by `NumberPicker`.
	// Set the maximum value. The default value is 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	NumberPickerMax PropertyName = "number-picker-max"

	// NumberPickerStep is the constant for "number-picker-step" property tag.
	//
	// Used by `NumberPicker`.
	// Set the value change step.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	NumberPickerStep PropertyName = "number-picker-step"

	// NumberPickerValue is the constant for "number-picker-value" property tag.
	//
	// Used by `NumberPicker`.
	// Current value. The default value is 0.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	NumberPickerValue PropertyName = "number-picker-value"
)

// Constants which describe values of the "number-picker-type" property of a [NumberPicker]
const (
	// NumberEditor - type of NumberPicker. NumberPicker is presented by editor
	NumberEditor = 0

	// NumberSlider - type of NumberPicker. NumberPicker is presented by slider
	NumberSlider = 1
)

// NumberPicker represents a NumberPicker view
type NumberPicker interface {
	View
}

type numberPickerData struct {
	viewData
}

// NewNumberPicker create new NumberPicker object and return it
func NewNumberPicker(session Session, params Params) NumberPicker {
	view := new(numberPickerData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newNumberPicker(session Session) View {
	return new(numberPickerData)
}

func (picker *numberPickerData) init(session Session) {
	picker.viewData.init(session)
	picker.tag = "NumberPicker"
	picker.hasHtmlDisabled = true
	picker.normalize = normalizeNumberPickerTag
	picker.set = numberPickerSet
	picker.changed = numberPickerPropertyChanged
}

func (picker *numberPickerData) Focusable() bool {
	return true
}

func normalizeNumberPickerTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Type, Min, Max, Step, Value:
		return "number-picker-" + tag
	}

	return normalizeDataListTag(tag)
}

func numberPickerSet(view View, tag PropertyName, value any) []PropertyName {
	switch tag {
	case NumberChangedEvent:
		return setEventWithOldListener[NumberPicker, float64](view, tag, value)

	case NumberPickerValue:
		view.setRaw("old-number", GetNumberPickerValue(view))
		min, max := GetNumberPickerMinMax(view)

		return setFloatProperty(view, NumberPickerValue, value, min, max)

	case DataList:
		return setDataList(view, value, "")
	}

	return viewSet(view, tag, value)
}

func numberPickerPropertyChanged(view View, tag PropertyName) {
	switch tag {
	case NumberPickerType:
		if GetNumberPickerType(view) == NumberSlider {
			view.Session().updateProperty(view.htmlID(), "type", "range")
		} else {
			view.Session().updateProperty(view.htmlID(), "type", "number")
		}

	case NumberPickerMin:
		min, _ := GetNumberPickerMinMax(view)
		view.Session().updateProperty(view.htmlID(), "min", strconv.FormatFloat(min, 'f', -1, 32))

	case NumberPickerMax:
		_, max := GetNumberPickerMinMax(view)
		view.Session().updateProperty(view.htmlID(), "max", strconv.FormatFloat(max, 'f', -1, 32))

	case NumberPickerStep:
		if step := GetNumberPickerStep(view); step > 0 {
			view.Session().updateProperty(view.htmlID(), "step", strconv.FormatFloat(step, 'f', -1, 32))
		} else {
			view.Session().updateProperty(view.htmlID(), "step", "any")
		}

	case TimePickerValue:
		value := GetNumberPickerValue(view)
		view.Session().callFunc("setInputValue", view.htmlID(), value)

		if listeners := GetNumberChangedListeners(view); len(listeners) > 0 {
			old := 0.0
			if val := view.getRaw("old-number"); val != nil {
				if n, ok := val.(float64); ok {
					old = n
				}
			}
			for _, listener := range listeners {
				listener(view, value, old)
			}
		}

	default:
		viewPropertyChanged(view, tag)
	}
}

func (picker *numberPickerData) htmlTag() string {
	return "input"
}

func (picker *numberPickerData) htmlSubviews(self View, buffer *strings.Builder) {
	dataListHtmlSubviews(self, buffer, func(text string, session Session) string {
		text, _ = session.resolveConstants(text)
		return text
	})
}

func (picker *numberPickerData) htmlProperties(self View, buffer *strings.Builder) {
	picker.viewData.htmlProperties(self, buffer)

	if GetNumberPickerType(picker) == NumberSlider {
		buffer.WriteString(` type="range"`)
	} else {
		buffer.WriteString(` type="number"`)
	}

	min, max := GetNumberPickerMinMax(picker)
	if min != math.Inf(-1) {
		buffer.WriteString(` min="`)
		buffer.WriteString(strconv.FormatFloat(min, 'f', -1, 64))
		buffer.WriteByte('"')
	}

	if max != math.Inf(1) {
		buffer.WriteString(` max="`)
		buffer.WriteString(strconv.FormatFloat(max, 'f', -1, 64))
		buffer.WriteByte('"')
	}

	step := GetNumberPickerStep(picker)
	if step != 0 {
		buffer.WriteString(` step="`)
		buffer.WriteString(strconv.FormatFloat(step, 'f', -1, 64))
		buffer.WriteByte('"')
	} else {
		buffer.WriteString(` step="any"`)
	}

	buffer.WriteString(` value="`)
	buffer.WriteString(strconv.FormatFloat(GetNumberPickerValue(picker), 'f', -1, 64))
	buffer.WriteByte('"')

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)

	dataListHtmlProperties(picker, buffer)
}

func (picker *numberPickerData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if value, err := strconv.ParseFloat(text, 32); err == nil {
				oldValue := GetNumberPickerValue(picker)
				picker.properties[NumberPickerValue] = text
				if value != oldValue {
					for _, listener := range GetNumberChangedListeners(picker) {
						listener(picker, value, oldValue)
					}
					if listener, ok := picker.changeListener[NumberPickerValue]; ok {
						listener(picker, NumberPickerValue)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

// GetNumberPickerType returns the type of NumberPicker subview. Valid values:
// NumberEditor (0) - NumberPicker is presented by editor (default type);
// NumberSlider (1) - NumberPicker is presented by slider.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetNumberPickerType(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, NumberPickerType, NumberEditor, false)
}

// GetNumberPickerMinMax returns the min and max value of NumberPicker subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetNumberPickerMinMax(view View, subviewID ...string) (float64, float64) {
	var pickerType int
	if len(subviewID) > 0 && subviewID[0] != "" {
		pickerType = GetNumberPickerType(view, subviewID[0])
	} else {
		pickerType = GetNumberPickerType(view)
	}

	var defMin, defMax float64
	if pickerType == NumberSlider {
		defMin = 0
		defMax = 1
	} else {
		defMin = math.Inf(-1)
		defMax = math.Inf(1)
	}

	min := floatStyledProperty(view, subviewID, NumberPickerMin, defMin)
	max := floatStyledProperty(view, subviewID, NumberPickerMax, defMax)

	if min > max {
		return max, min
	}
	return min, max
}

// GetNumberPickerStep returns the value changing step of NumberPicker subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetNumberPickerStep(view View, subviewID ...string) float64 {
	var max float64
	if len(subviewID) > 0 && subviewID[0] != "" {
		_, max = GetNumberPickerMinMax(view, subviewID[0])
	} else {
		_, max = GetNumberPickerMinMax(view)
	}

	result := floatStyledProperty(view, subviewID, NumberPickerStep, 0)
	if result > max {
		return max
	}
	return result
}

// GetNumberPickerValue returns the value of NumberPicker subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetNumberPickerValue(view View, subviewID ...string) float64 {
	var min float64
	if len(subviewID) > 0 && subviewID[0] != "" {
		min, _ = GetNumberPickerMinMax(view, subviewID[0])
	} else {
		min, _ = GetNumberPickerMinMax(view)
	}

	result := floatStyledProperty(view, subviewID, NumberPickerValue, min)
	return result
}

// GetNumberChangedListeners returns the NumberChangedListener list of an NumberPicker subview.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetNumberChangedListeners(view View, subviewID ...string) []func(NumberPicker, float64, float64) {
	return getEventWithOldListeners[NumberPicker, float64](view, subviewID, NumberChangedEvent)
}
