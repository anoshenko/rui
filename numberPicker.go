package rui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	NumberChangedEvent = "number-changed"
	NumberPickerType   = "number-picker-type"
	NumberPickerMin    = "number-picker-min"
	NumberPickerMax    = "number-picker-max"
	NumberPickerStep   = "number-picker-step"
	NumberPickerValue  = "number-picker-value"
)

const (
	// NumberEditor - type of NumberPicker. NumberPicker is presented by editor
	NumberEditor = 0
	// NumberSlider - type of NumberPicker. NumberPicker is presented by slider
	NumberSlider = 1
)

// NumberPicker - NumberPicker view
type NumberPicker interface {
	View
}

type numberPickerData struct {
	viewData
	numberChangedListeners []func(NumberPicker, float64)
}

// NewNumberPicker create new NumberPicker object and return it
func NewNumberPicker(session Session, params Params) NumberPicker {
	view := new(numberPickerData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newNumberPicker(session Session) View {
	return NewNumberPicker(session, nil)
}

func (picker *numberPickerData) Init(session Session) {
	picker.viewData.Init(session)
	picker.tag = "NumberPicker"
	picker.numberChangedListeners = []func(NumberPicker, float64){}
}

func (picker *numberPickerData) Focusable() bool {
	return true
}

func (picker *numberPickerData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Type, Min, Max, Step, Value:
		return "number-picker-" + tag
	}

	return tag
}

func (picker *numberPickerData) Remove(tag string) {
	picker.remove(picker.normalizeTag(tag))
}

func (picker *numberPickerData) remove(tag string) {
	switch tag {
	case NumberChangedEvent:
		if len(picker.numberChangedListeners) > 0 {
			picker.numberChangedListeners = []func(NumberPicker, float64){}
			picker.propertyChangedEvent(tag)
		}

	default:
		picker.viewData.remove(tag)
		picker.propertyChanged(tag)
	}
}

func (picker *numberPickerData) Set(tag string, value interface{}) bool {
	return picker.set(picker.normalizeTag(tag), value)
}

func (picker *numberPickerData) set(tag string, value interface{}) bool {
	if value == nil {
		picker.remove(tag)
		return true
	}

	switch tag {
	case NumberChangedEvent:
		switch value := value.(type) {
		case func(NumberPicker, float64):
			picker.numberChangedListeners = []func(NumberPicker, float64){value}

		case func(float64):
			fn := func(view NumberPicker, newValue float64) {
				value(newValue)
			}
			picker.numberChangedListeners = []func(NumberPicker, float64){fn}

		case []func(NumberPicker, float64):
			picker.numberChangedListeners = value

		case []func(float64):
			listeners := make([]func(NumberPicker, float64), len(value))
			for i, val := range value {
				if val == nil {
					notCompatibleType(tag, val)
					return false
				}

				listeners[i] = func(view NumberPicker, newValue float64) {
					val(newValue)
				}
			}
			picker.numberChangedListeners = listeners

		case []interface{}:
			listeners := make([]func(NumberPicker, float64), len(value))
			for i, val := range value {
				if val == nil {
					notCompatibleType(tag, val)
					return false
				}

				switch val := val.(type) {
				case func(NumberPicker, float64):
					listeners[i] = val

				default:
					notCompatibleType(tag, val)
					return false
				}
			}
			picker.numberChangedListeners = listeners
		}
		picker.propertyChangedEvent(tag)
		return true

	case NumberPickerValue:
		oldValue := GetNumberPickerValue(picker, "")
		min, max := GetNumberPickerMinMax(picker, "")
		if picker.setFloatProperty(NumberPickerValue, value, min, max) {
			if newValue := GetNumberPickerValue(picker, ""); oldValue != newValue {
				if picker.created {
					picker.session.runScript(fmt.Sprintf(`setInputValue('%s', '%f')`, picker.htmlID(), newValue))
				}
				for _, listener := range picker.numberChangedListeners {
					listener(picker, newValue)
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	default:
		if picker.viewData.set(tag, value) {
			picker.propertyChanged(tag)
			return true
		}
	}
	return false
}

func (picker *numberPickerData) propertyChanged(tag string) {
	if picker.created {
		switch tag {
		case NumberPickerType:
			if GetNumberPickerType(picker, "") == NumberSlider {
				updateProperty(picker.htmlID(), "type", "range", picker.session)
			} else {
				updateProperty(picker.htmlID(), "type", "number", picker.session)
			}

		case NumberPickerMin:
			min, _ := GetNumberPickerMinMax(picker, "")
			updateProperty(picker.htmlID(), Min, strconv.FormatFloat(min, 'f', -1, 32), picker.session)

		case NumberPickerMax:
			_, max := GetNumberPickerMinMax(picker, "")
			updateProperty(picker.htmlID(), Max, strconv.FormatFloat(max, 'f', -1, 32), picker.session)

		case NumberPickerStep:
			if step := GetNumberPickerStep(picker, ""); step > 0 {
				updateProperty(picker.htmlID(), Step, strconv.FormatFloat(step, 'f', -1, 32), picker.session)
			} else {
				updateProperty(picker.htmlID(), Step, "any", picker.session)
			}

		case NumberPickerValue:
			value := GetNumberPickerValue(picker, "")
			picker.session.runScript(fmt.Sprintf(`setInputValue('%s', '%f')`, picker.htmlID(), value))
			for _, listener := range picker.numberChangedListeners {
				listener(picker, value)
			}
		}
	}
}

func (picker *numberPickerData) Get(tag string) interface{} {
	return picker.get(picker.normalizeTag(tag))
}

func (picker *numberPickerData) get(tag string) interface{} {
	switch tag {
	case NumberChangedEvent:
		return picker.numberChangedListeners

	default:
		return picker.viewData.get(tag)
	}
}

func (picker *numberPickerData) htmlTag() string {
	return "input"
}

func (picker *numberPickerData) htmlProperties(self View, buffer *strings.Builder) {
	picker.viewData.htmlProperties(self, buffer)

	if GetNumberPickerType(picker, "") == NumberSlider {
		buffer.WriteString(` type="range"`)
	} else {
		buffer.WriteString(` type="number"`)
	}

	min, max := GetNumberPickerMinMax(picker, "")
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

	step := GetNumberPickerStep(picker, "")
	if step != 0 {
		buffer.WriteString(` step="`)
		buffer.WriteString(strconv.FormatFloat(step, 'f', -1, 64))
		buffer.WriteByte('"')
	} else {
		buffer.WriteString(` step="any"`)
	}

	buffer.WriteString(` value="`)
	buffer.WriteString(strconv.FormatFloat(GetNumberPickerValue(picker, ""), 'f', -1, 64))
	buffer.WriteByte('"')

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)
}

func (picker *numberPickerData) htmlDisabledProperties(self View, buffer *strings.Builder) {
	if IsDisabled(self, "") {
		buffer.WriteString(` disabled`)
	}
	picker.viewData.htmlDisabledProperties(self, buffer)
}

func (picker *numberPickerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if value, err := strconv.ParseFloat(text, 32); err == nil {
				oldValue := GetNumberPickerValue(picker, "")
				picker.properties[NumberPickerValue] = value
				if value != oldValue {
					for _, listener := range picker.numberChangedListeners {
						listener(picker, value)
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
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetNumberPickerType(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return 0
	}

	t, _ := enumStyledProperty(view, NumberPickerType, NumberEditor)
	return t
}

// GetNumberPickerMinMax returns the min and max value of NumberPicker subview.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetNumberPickerMinMax(view View, subviewID string) (float64, float64) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		t, _ := enumStyledProperty(view, NumberPickerType, NumberEditor)
		var defMin, defMax float64

		if t == NumberSlider {
			defMin = 0
			defMax = 1
		} else {
			defMin = math.Inf(-1)
			defMax = math.Inf(1)
		}
		min, ok := floatStyledProperty(view, NumberPickerMin, defMin)
		if !ok {
			min, _ = floatStyledProperty(view, Min, defMin)
		}

		max, ok := floatStyledProperty(view, NumberPickerMax, defMax)
		if !ok {
			max, _ = floatStyledProperty(view, Max, defMax)
		}

		if min > max {
			return max, min
		}
		return min, max
	}
	return 0, 1
}

// GetNumberPickerStep returns the value changing step of NumberPicker subview.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetNumberPickerStep(view View, subviewID string) float64 {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return 0
	}

	result, ok := floatStyledProperty(view, NumberPickerStep, 0)
	if !ok {
		result, _ = floatStyledProperty(view, Step, 0)
	}

	_, max := GetNumberPickerMinMax(view, "")
	if result > max {
		return max
	}
	return result
}

// GetNumberPickerValue returns the value of NumberPicker subview.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetNumberPickerValue(view View, subviewID string) float64 {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return 0
	}

	min, _ := GetNumberPickerMinMax(view, "")
	result, ok := floatStyledProperty(view, NumberPickerValue, min)
	if !ok {
		result, _ = floatStyledProperty(view, Value, min)
	}
	return result
}

// GetNumberChangedListeners returns the NumberChangedListener list of an NumberPicker subview.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetNumberChangedListeners(view View, subviewID string) []func(NumberPicker, float64) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(NumberChangedEvent); value != nil {
			if listeners, ok := value.([]func(NumberPicker, float64)); ok {
				return listeners
			}
		}
	}
	return []func(NumberPicker, float64){}
}
