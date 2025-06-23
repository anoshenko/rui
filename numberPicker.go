package rui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Constants related to [NumberPicker] specific properties and events
const (
	// NumberChangedEvent is the constant for "number-changed" property tag.
	//
	// Used by NumberPicker.
	// Set listener(s) that track the change in the entered value.
	//
	// General listener format:
	//
	//  func(picker rui.NumberPicker, newValue float64, oldValue float64)
	//
	// where:
	//   - picker - Interface of a number picker which generated this event,
	//   - newValue - New value,
	//   - oldValue - Old Value.
	//
	// Allowed listener formats:
	//
	//  func(picker rui.NumberPicker, newValue float64)
	//  func(newValue float64, oldValue float64)
	//  func(newValue float64)
	//  func()
	NumberChangedEvent PropertyName = "number-changed"

	// NumberPickerType is the constant for "number-picker-type" property tag.
	//
	// Used by NumberPicker.
	// Sets the visual representation.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NumberEditor) or "editor" - Displayed as an editor.
	//   - 1 (NumberSlider) or "slider" - Displayed as a slider.
	NumberPickerType PropertyName = "number-picker-type"

	// NumberPickerMin is the constant for "number-picker-min" property tag.
	//
	// Used by NumberPicker.
	// Set the minimum value. The default value is 0.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	NumberPickerMin PropertyName = "number-picker-min"

	// NumberPickerMax is the constant for "number-picker-max" property tag.
	//
	// Used by NumberPicker.
	// Set the maximum value. The default value is 1.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	NumberPickerMax PropertyName = "number-picker-max"

	// NumberPickerStep is the constant for "number-picker-step" property tag.
	//
	// Used by NumberPicker.
	// Set the value change step.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	NumberPickerStep PropertyName = "number-picker-step"

	// NumberPickerValue is the constant for "number-picker-value" property tag.
	//
	// Used by NumberPicker.
	// Current value. The default value is 0.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	NumberPickerValue PropertyName = "number-picker-value"

	// NumberPickerValue is the constant for "number-picker-value" property tag.
	//
	// Used by NumberPicker.
	// Precision of displaying fractional part in editor. The default value is 0 (not used).
	//
	// Supported types: int, int8...int64, uint, uint8...uint64, string.
	//
	// Internal type is float, other types converted to it during assignment.
	NumberPickerPrecision PropertyName = "number-picker-precision"
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
	picker.set = picker.setFunc
	picker.changed = picker.propertyChanged
}

func (picker *numberPickerData) Focusable() bool {
	return true
}

func normalizeNumberPickerTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Type, Min, Max, Step, Value, "precision":
		return "number-picker-" + tag
	}

	return normalizeDataListTag(tag)
}

func (picker *numberPickerData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case NumberChangedEvent:
		return setTwoArgEventListener[NumberPicker, float64](picker, tag, value)

	case NumberPickerValue:
		picker.setRaw("old-number", GetNumberPickerValue(picker))
		min, max := GetNumberPickerMinMax(picker)

		return setFloatProperty(picker, NumberPickerValue, value, min, max)

	case DataList:
		return setDataList(picker, value, "")
	}

	return picker.viewData.setFunc(tag, value)
}

func (picker *numberPickerData) numberFormat() string {
	if precission := GetNumberPickerPrecision(picker); precission > 0 {
		return fmt.Sprintf("%%.%df", precission)
	}
	return "%g"
}

func (picker *numberPickerData) propertyChanged(tag PropertyName) {
	switch tag {
	case NumberPickerType:
		if GetNumberPickerType(picker) == NumberSlider {
			picker.Session().updateProperty(picker.htmlID(), "type", "range")
		} else {
			picker.Session().updateProperty(picker.htmlID(), "type", "number")
		}

	case NumberPickerMin:
		min, _ := GetNumberPickerMinMax(picker)
		picker.Session().updateProperty(picker.htmlID(), "min", fmt.Sprintf(picker.numberFormat(), min))

	case NumberPickerMax:
		_, max := GetNumberPickerMinMax(picker)
		picker.Session().updateProperty(picker.htmlID(), "max", fmt.Sprintf(picker.numberFormat(), max))

	case NumberPickerStep:
		if step := GetNumberPickerStep(picker); step > 0 {
			picker.Session().updateProperty(picker.htmlID(), "step", fmt.Sprintf(picker.numberFormat(), step))
		} else {
			picker.Session().updateProperty(picker.htmlID(), "step", "any")
		}

	case NumberPickerValue:
		value := GetNumberPickerValue(picker)
		format := picker.numberFormat()
		picker.Session().callFunc("setInputValue", picker.htmlID(), fmt.Sprintf(format, value))

		if listeners := getTwoArgEventListeners[NumberPicker, float64](picker, nil, NumberChangedEvent); len(listeners) > 0 {
			old := 0.0
			if val := picker.getRaw("old-number"); val != nil {
				if n, ok := val.(float64); ok {
					old = n
				}
			}
			if old != value {
				for _, listener := range listeners {
					listener.Run(picker, value, old)
				}
			}
		}

	default:
		picker.viewData.propertyChanged(tag)
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

	format := picker.numberFormat()
	min, max := GetNumberPickerMinMax(picker)
	if min != math.Inf(-1) {
		buffer.WriteString(` min="`)
		fmt.Fprintf(buffer, format, min)
		buffer.WriteByte('"')
	}

	if max != math.Inf(1) {
		buffer.WriteString(` max="`)
		fmt.Fprintf(buffer, format, max)
		buffer.WriteByte('"')
	}

	step := GetNumberPickerStep(picker)
	if step != 0 {
		buffer.WriteString(` step="`)
		fmt.Fprintf(buffer, format, step)
		buffer.WriteByte('"')
	} else {
		buffer.WriteString(` step="any"`)
	}

	buffer.WriteString(` value="`)
	fmt.Fprintf(buffer, format, GetNumberPickerValue(picker))
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
					for _, listener := range getTwoArgEventListeners[NumberPicker, float64](picker, nil, NumberChangedEvent) {
						listener.Run(picker, value, oldValue)
					}
					if listener, ok := picker.changeListener[NumberPickerValue]; ok {
						listener.Run(picker, NumberPickerValue)
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
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetNumberPickerType(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, NumberPickerType, NumberEditor, false)
}

// GetNumberPickerMinMax returns the min and max value of NumberPicker subview.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetNumberPickerMinMax(view View, subviewID ...string) (float64, float64) {
	view = getSubview(view, subviewID)
	pickerType := GetNumberPickerType(view)

	var defMin, defMax float64
	if pickerType == NumberSlider {
		defMin = 0
		defMax = 1
	} else {
		defMin = math.Inf(-1)
		defMax = math.Inf(1)
	}

	min := floatStyledProperty(view, nil, NumberPickerMin, defMin)
	max := floatStyledProperty(view, nil, NumberPickerMax, defMax)

	if min > max {
		return max, min
	}
	return min, max
}

// GetNumberPickerStep returns the value changing step of NumberPicker subview.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetNumberPickerStep(view View, subviewID ...string) float64 {
	view = getSubview(view, subviewID)
	_, max := GetNumberPickerMinMax(view)

	result := floatStyledProperty(view, nil, NumberPickerStep, 0)
	if result > max {
		return max
	}
	return result
}

// GetNumberPickerValue returns the value of NumberPicker subview.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetNumberPickerValue(view View, subviewID ...string) float64 {
	view = getSubview(view, subviewID)
	min, _ := GetNumberPickerMinMax(view)
	return floatStyledProperty(view, nil, NumberPickerValue, min)
}

// GetNumberChangedListeners returns the NumberChangedListener list of an NumberPicker subview.
// If there are no listeners then the empty list is returned
//
// Result elements can be of the following types:
//   - func(rui.NumberPicker, float64, float64),
//   - func(rui.NumberPicker, float64),
//   - func(rui.NumberPicker),
//   - func(float64, float64),
//   - func(float64),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetNumberChangedListeners(view View, subviewID ...string) []any {
	return getTwoArgEventRawListeners[NumberPicker, float64](view, subviewID, NumberChangedEvent)
}

// GetNumberPickerPrecision returns the precision of displaying fractional part in editor of NumberPicker subview.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetNumberPickerPrecision(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, NumberPickerPrecision, 0)
}
