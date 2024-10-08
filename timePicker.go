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
	TimeChangedEvent = "time-changed"

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
	TimePickerMin = "time-picker-min"

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
	TimePickerMax = "time-picker-max"

	// TimePickerStep is the constant for "time-picker-step" property tag.
	//
	// Used by `TimePicker`.
	// Time step in seconds.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// >= `0` or >= "0" - Step value in seconds used to increment or decrement time.
	TimePickerStep = "time-picker-step"

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
	TimePickerValue = "time-picker-value"

	timeFormat = "15:04:05"
)

// TimePicker represents a TimePicker view
type TimePicker interface {
	View
}

type timePickerData struct {
	viewData
	dataList
	timeChangedListeners []func(TimePicker, time.Time, time.Time)
}

// NewTimePicker create new TimePicker object and return it
func NewTimePicker(session Session, params Params) TimePicker {
	view := new(timePickerData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newTimePicker(session Session) View {
	return NewTimePicker(session, nil)
}

func (picker *timePickerData) init(session Session) {
	picker.viewData.init(session)
	picker.tag = "TimePicker"
	picker.hasHtmlDisabled = true
	picker.timeChangedListeners = []func(TimePicker, time.Time, time.Time){}
	picker.dataListInit()
}

func (picker *timePickerData) String() string {
	return getViewString(picker, nil)
}

func (picker *timePickerData) Focusable() bool {
	return true
}

func (picker *timePickerData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Type, Min, Max, Step, Value:
		return "time-picker-" + tag
	}

	return tag
}

func (picker *timePickerData) Remove(tag string) {
	picker.remove(picker.normalizeTag(tag))
}

func (picker *timePickerData) remove(tag string) {
	switch tag {
	case TimeChangedEvent:
		if len(picker.timeChangedListeners) > 0 {
			picker.timeChangedListeners = []func(TimePicker, time.Time, time.Time){}
			picker.propertyChangedEvent(tag)
		}
		return

	case TimePickerMin:
		delete(picker.properties, TimePickerMin)
		if picker.created {
			picker.session.removeProperty(picker.htmlID(), Min)
		}

	case TimePickerMax:
		delete(picker.properties, TimePickerMax)
		if picker.created {
			picker.session.removeProperty(picker.htmlID(), Max)
		}

	case TimePickerStep:
		delete(picker.properties, TimePickerStep)
		if picker.created {
			picker.session.removeProperty(picker.htmlID(), Step)
		}

	case TimePickerValue:
		if _, ok := picker.properties[TimePickerValue]; ok {
			oldTime := GetTimePickerValue(picker)
			delete(picker.properties, TimePickerValue)
			time := GetTimePickerValue(picker)
			if picker.created {
				picker.session.callFunc("setInputValue", picker.htmlID(), time.Format(timeFormat))
			}
			for _, listener := range picker.timeChangedListeners {
				listener(picker, time, oldTime)
			}
		} else {
			return
		}

	case DataList:
		if len(picker.dataList.dataList) > 0 {
			picker.setDataList(picker, []string{}, true)
		}

	default:
		picker.viewData.remove(tag)
		return
	}
	picker.propertyChangedEvent(tag)
}

func (picker *timePickerData) Set(tag string, value any) bool {
	return picker.set(picker.normalizeTag(tag), value)
}

func (picker *timePickerData) set(tag string, value any) bool {
	if value == nil {
		picker.remove(tag)
		return true
	}

	setTimeValue := func(tag string) (time.Time, bool) {
		switch value := value.(type) {
		case time.Time:
			picker.properties[tag] = value
			return value, true

		case string:
			if text, ok := picker.Session().resolveConstants(value); ok {
				lowText := strings.ToLower(text)
				pm := strings.HasSuffix(lowText, "pm") || strings.HasSuffix(lowText, "am")

				var format string
				switch len(strings.Split(text, ":")) {
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

				if time, err := time.Parse(format, text); err == nil {
					picker.properties[tag] = value
					return time, true
				} else {
					ErrorLog(err.Error())
				}
				return time.Now(), false
			}
		}

		notCompatibleType(tag, value)
		return time.Now(), false
	}

	switch tag {
	case TimePickerMin:
		old, oldOK := getTimeProperty(picker, TimePickerMin, Min)
		if time, ok := setTimeValue(TimePickerMin); ok {
			if !oldOK || time != old {
				if picker.created {
					picker.session.updateProperty(picker.htmlID(), Min, time.Format(timeFormat))
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case TimePickerMax:
		old, oldOK := getTimeProperty(picker, TimePickerMax, Max)
		if time, ok := setTimeValue(TimePickerMax); ok {
			if !oldOK || time != old {
				if picker.created {
					picker.session.updateProperty(picker.htmlID(), Max, time.Format(timeFormat))
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case TimePickerStep:
		oldStep := GetTimePickerStep(picker)
		if picker.setIntProperty(TimePickerStep, value) {
			if step := GetTimePickerStep(picker); oldStep != step {
				if picker.created {
					if step > 0 {
						picker.session.updateProperty(picker.htmlID(), Step, strconv.Itoa(step))
					} else {
						picker.session.removeProperty(picker.htmlID(), Step)
					}
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case TimePickerValue:
		oldTime := GetTimePickerValue(picker)
		if time, ok := setTimeValue(TimePickerValue); ok {
			if time != oldTime {
				if picker.created {
					picker.session.callFunc("setInputValue", picker.htmlID(), time.Format(timeFormat))
				}
				for _, listener := range picker.timeChangedListeners {
					listener(picker, time, oldTime)
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case TimeChangedEvent:
		listeners, ok := valueToEventWithOldListeners[TimePicker, time.Time](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(TimePicker, time.Time, time.Time){}
		}
		picker.timeChangedListeners = listeners
		picker.propertyChangedEvent(tag)
		return true

	case DataList:
		return picker.setDataList(picker, value, picker.created)

	default:
		return picker.viewData.set(tag, value)
	}
	return false
}

func (picker *timePickerData) Get(tag string) any {
	return picker.get(picker.normalizeTag(tag))
}

func (picker *timePickerData) get(tag string) any {
	switch tag {
	case TimeChangedEvent:
		return picker.timeChangedListeners

	case DataList:
		return picker.dataList.dataList

	default:
		return picker.viewData.get(tag)
	}
}

func (picker *timePickerData) htmlTag() string {
	return "input"
}

func (picker *timePickerData) htmlSubviews(self View, buffer *strings.Builder) {
	picker.dataListHtmlSubviews(self, buffer)
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

	picker.dataListHtmlProperties(picker, buffer)
}

func (picker *timePickerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if value, err := time.Parse(timeFormat, text); err == nil {
				oldValue := GetTimePickerValue(picker)
				picker.properties[TimePickerValue] = value
				if value != oldValue {
					for _, listener := range picker.timeChangedListeners {
						listener(picker, value, oldValue)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

func getTimeProperty(view View, mainTag, shortTag string) (time.Time, bool) {
	valueToTime := func(value any) (time.Time, bool) {
		if value != nil {
			switch value := value.(type) {
			case time.Time:
				return value, true

			case string:
				if text, ok := view.Session().resolveConstants(value); ok {
					if result, err := time.Parse(timeFormat, text); err == nil {
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

		if value := valueFromStyle(view, shortTag); value != nil {
			if result, ok := valueToTime(value); ok {
				return result, true
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
	return getEventWithOldListeners[TimePicker, time.Time](view, subviewID, TimeChangedEvent)
}
