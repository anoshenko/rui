package rui

import (
	"strconv"
	"strings"
	"time"
)

// Constants for [TimePicker] specific properties and events.
const (
	// TimeChangedEvent is the constant for "time-changed" property tag.
	//
	// Used by `TimePicker`.
	// Occur when current time of the time picker has been changed.
	//
	// General listener format:
	// `func(picker rui.TimePicker, newTime, oldTime time.Time)`.
	//
	// where:
	// picker - Interface of a time picker which generated this event,
	// newTime - New time value,
	// oldTime - Old time value.
	//
	// Allowed listener formats:
	// `func(picker rui.TimePicker, newTime time.Time)`,
	// `func(newTime, oldTime time.Time)`,
	// `func(newTime time.Time)`,
	// `func(picker rui.TimePicker)`,
	// `func()`.
	TimeChangedEvent PropertyName = "time-changed"

	// TimePickerMin is the constant for "time-picker-min" property tag.
	//
	// Used by `TimePicker`.
	// The minimum value of the time.
	//
	// Supported types: `time.Time`, `string`.
	//
	// Internal type is `time.Time`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - values of this type parsed and converted to `time.Time`. The following formats are supported:
	// "HH:MM:SS" - "08:15:00".
	// "HH:MM:SS PM" - "08:15:00 AM".
	// "HH:MM" - "08:15".
	// "HH:MM PM" - "08:15 AM".
	TimePickerMin PropertyName = "time-picker-min"

	// TimePickerMax is the constant for "time-picker-max" property tag.
	//
	// Used by `TimePicker`.
	// The maximum value of the time.
	//
	// Supported types: `time.Time`, `string`.
	//
	// Internal type is `time.Time`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - values of this type parsed and converted to `time.Time`. The following formats are supported:
	// "HH:MM:SS" - "08:15:00".
	// "HH:MM:SS PM" - "08:15:00 AM".
	// "HH:MM" - "08:15".
	// "HH:MM PM" - "08:15 AM".
	TimePickerMax PropertyName = "time-picker-max"

	// TimePickerStep is the constant for "time-picker-step" property tag.
	//
	// Used by `TimePicker`.
	// Time step in seconds.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// >= `0` or >= "0" - Step value in seconds used to increment or decrement time.
	TimePickerStep PropertyName = "time-picker-step"

	// TimePickerValue is the constant for "time-picker-value" property tag.
	//
	// Used by `TimePicker`.
	// Current value.
	//
	// Supported types: `time.Time`, `string`.
	//
	// Internal type is `time.Time`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - values of this type parsed and converted to `time.Time`. The following formats are supported:
	// "HH:MM:SS" - "08:15:00".
	// "HH:MM:SS PM" - "08:15:00 AM".
	// "HH:MM" - "08:15".
	// "HH:MM PM" - "08:15 AM".
	TimePickerValue PropertyName = "time-picker-value"

	timeFormat = "15:04:05"
)

// TimePicker represents a TimePicker view
type TimePicker interface {
	View
}

type timePickerData struct {
	viewData
}

// NewTimePicker create new TimePicker object and return it
func NewTimePicker(session Session, params Params) TimePicker {
	view := new(timePickerData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newTimePicker(session Session) View {
	return new(timePickerData)
}

func (picker *timePickerData) init(session Session) {
	picker.viewData.init(session)
	picker.tag = "TimePicker"
	picker.hasHtmlDisabled = true
	picker.normalize = normalizeTimePickerTag
	picker.set = picker.setFunc
	picker.changed = picker.propertyChanged
}

func (picker *timePickerData) Focusable() bool {
	return true
}

func normalizeTimePickerTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Type, Min, Max, Step, Value:
		return "time-picker-" + tag
	}

	return normalizeDataListTag(tag)
}

func stringToTime(value string) (time.Time, bool) {
	lowText := strings.ToUpper(value)
	pm := strings.HasSuffix(lowText, "PM") || strings.HasSuffix(lowText, "AM")

	var format string
	switch len(strings.Split(value, ":")) {
	case 2:
		if pm {
			format = "3:04 PM"
		} else {
			format = "15:04"
		}

	default:
		if pm {
			format = "03:04:05 PM"
		} else {
			format = "15:04:05"
		}
	}

	result, err := time.Parse(format, value)
	if err != nil {
		ErrorLog(err.Error())
		return time.Now(), false
	}
	return result, true
}

func (picker *timePickerData) setFunc(tag PropertyName, value any) []PropertyName {

	setTimeValue := func(tag PropertyName) []PropertyName {
		switch value := value.(type) {
		case time.Time:
			picker.setRaw(tag, value)
			return []PropertyName{tag}

		case string:
			if isConstantName(value) {
				picker.setRaw(tag, value)
				return []PropertyName{tag}
			}

			if time, ok := stringToTime(value); ok {
				picker.setRaw(tag, time)
				return []PropertyName{tag}
			}
		}

		notCompatibleType(tag, value)
		return nil
	}

	switch tag {
	case TimePickerMin:
		return setTimeValue(TimePickerMin)

	case TimePickerMax:
		return setTimeValue(TimePickerMax)

	case TimePickerStep:
		return setIntProperty(picker, TimePickerStep, value)

	case TimePickerValue:
		picker.setRaw("old-time", GetTimePickerValue(picker))
		return setTimeValue(tag)

	case TimeChangedEvent:
		return setTwoArgEventListener[TimePicker, time.Time](picker, tag, value)

	case DataList:
		return setDataList(picker, value, timeFormat)
	}

	return picker.viewData.setFunc(tag, value)
}

func (picker *timePickerData) propertyChanged(tag PropertyName) {

	session := picker.Session()

	switch tag {

	case TimePickerMin:
		if time, ok := GetTimePickerMin(picker); ok {
			session.updateProperty(picker.htmlID(), "min", time.Format(timeFormat))
		} else {
			session.removeProperty(picker.htmlID(), "min")
		}

	case TimePickerMax:
		if time, ok := GetTimePickerMax(picker); ok {
			session.updateProperty(picker.htmlID(), "max", time.Format(timeFormat))
		} else {
			session.removeProperty(picker.htmlID(), "max")
		}

	case TimePickerStep:
		if step := GetTimePickerStep(picker); step > 0 {
			session.updateProperty(picker.htmlID(), "step", strconv.Itoa(step))
		} else {
			session.removeProperty(picker.htmlID(), "step")
		}

	case TimePickerValue:
		value := GetTimePickerValue(picker)
		session.callFunc("setInputValue", picker.htmlID(), value.Format(timeFormat))

		if listeners := GetTimeChangedListeners(picker); len(listeners) > 0 {
			oldTime := time.Now()
			if val := picker.getRaw("old-time"); val != nil {
				if time, ok := val.(time.Time); ok {
					oldTime = time
				}
			}
			for _, listener := range listeners {
				listener(picker, value, oldTime)
			}
		}

	default:
		picker.viewData.propertyChanged(tag)
	}
}

func (picker *timePickerData) htmlTag() string {
	return "input"
}

func (picker *timePickerData) htmlSubviews(self View, buffer *strings.Builder) {
	dataListHtmlSubviews(self, buffer, func(text string, session Session) string {
		text, _ = session.resolveConstants(text)
		if time, ok := stringToTime(text); ok {
			return time.Format(timeFormat)
		}
		return text
	})
}

func (picker *timePickerData) htmlProperties(self View, buffer *strings.Builder) {
	picker.viewData.htmlProperties(self, buffer)

	buffer.WriteString(` type="time"`)

	if min, ok := getTimeProperty(picker, TimePickerMin, Min); ok {
		buffer.WriteString(` min="`)
		buffer.WriteString(min.Format(timeFormat))
		buffer.WriteByte('"')
	}

	if max, ok := getTimeProperty(picker, TimePickerMax, Max); ok {
		buffer.WriteString(` max="`)
		buffer.WriteString(max.Format(timeFormat))
		buffer.WriteByte('"')
	}

	if step, ok := intProperty(picker, TimePickerStep, picker.Session(), 0); ok && step > 0 {
		buffer.WriteString(` step="`)
		buffer.WriteString(strconv.Itoa(step))
		buffer.WriteByte('"')
	}

	buffer.WriteString(` value="`)
	buffer.WriteString(GetTimePickerValue(picker).Format(timeFormat))
	buffer.WriteByte('"')

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)
	if picker.getRaw(ClickEvent) == nil {
		buffer.WriteString(` onclick="stopEventPropagation(this, event)"`)
	}

	dataListHtmlProperties(picker, buffer)
}

func (picker *timePickerData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if value, err := time.Parse(timeFormat, text); err == nil {
				oldValue := GetTimePickerValue(picker)
				picker.properties[TimePickerValue] = value
				if value != oldValue {
					for _, listener := range GetTimeChangedListeners(picker) {
						listener(picker, value, oldValue)
					}
					if listener, ok := picker.changeListener[TimePickerValue]; ok {
						listener(picker, TimePickerValue)
					}

				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

func getTimeProperty(view View, mainTag, shortTag PropertyName) (time.Time, bool) {
	valueToTime := func(value any) (time.Time, bool) {
		if value != nil {
			switch value := value.(type) {
			case time.Time:
				return value, true

			case string:
				if text, ok := view.Session().resolveConstants(value); ok {
					if result, ok := stringToTime(text); ok {
						return result, true
					}
				}
			}
		}
		return time.Now(), false
	}

	if view != nil {
		if result, ok := valueToTime(view.getRaw(mainTag)); ok {
			return result, true
		}

		for _, tag := range []PropertyName{mainTag, shortTag} {
			if value := valueFromStyle(view, tag); value != nil {
				if result, ok := valueToTime(value); ok {
					return result, true
				}
			}
		}
	}

	return time.Now(), false
}

// GetTimePickerMin returns the min time of TimePicker subview and "true" as the second value if the min time is set,
// "false" as the second value otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTimePickerMin(view View, subviewID ...string) (time.Time, bool) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		return getTimeProperty(view, TimePickerMin, Min)
	}
	return time.Now(), false
}

// GetTimePickerMax returns the max time of TimePicker subview and "true" as the second value if the min time is set,
// "false" as the second value otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTimePickerMax(view View, subviewID ...string) (time.Time, bool) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		return getTimeProperty(view, TimePickerMax, Max)
	}
	return time.Now(), false
}

// GetTimePickerStep returns the time changing step in seconds of TimePicker subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTimePickerStep(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, TimePickerStep, 60)
}

// GetTimePickerValue returns the time of TimePicker subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTimePickerValue(view View, subviewID ...string) time.Time {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return time.Now()
	}
	time, _ := getTimeProperty(view, TimePickerValue, Value)
	return time
}

// GetTimeChangedListeners returns the TimeChangedListener list of an TimePicker subview.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTimeChangedListeners(view View, subviewID ...string) []func(TimePicker, time.Time, time.Time) {
	return getTwoArgEventListeners[TimePicker, time.Time](view, subviewID, TimeChangedEvent)
}
