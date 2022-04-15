package rui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TimeChangedEvent = "time-changed"
	TimePickerMin    = "time-picker-min"
	TimePickerMax    = "time-picker-max"
	TimePickerStep   = "time-picker-step"
	TimePickerValue  = "time-picker-value"
	timeFormat       = "15:04"
)

// TimePicker - TimePicker view
type TimePicker interface {
	View
}

type timePickerData struct {
	viewData
	timeChangedListeners []func(TimePicker, time.Time)
}

// NewTimePicker create new TimePicker object and return it
func NewTimePicker(session Session, params Params) TimePicker {
	view := new(timePickerData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newTimePicker(session Session) View {
	return NewTimePicker(session, nil)
}

func (picker *timePickerData) Init(session Session) {
	picker.viewData.Init(session)
	picker.tag = "TimePicker"
	picker.timeChangedListeners = []func(TimePicker, time.Time){}
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
			picker.timeChangedListeners = []func(TimePicker, time.Time){}
			picker.propertyChangedEvent(tag)
		}
		return

	case TimePickerMin:
		delete(picker.properties, TimePickerMin)
		if picker.created {
			removeProperty(picker.htmlID(), Min, picker.session)
		}

	case TimePickerMax:
		delete(picker.properties, TimePickerMax)
		if picker.created {
			removeProperty(picker.htmlID(), Max, picker.session)
		}

	case TimePickerStep:
		delete(picker.properties, TimePickerMax)
		if picker.created {
			removeProperty(picker.htmlID(), Step, picker.session)
		}

	case TimePickerValue:
		if _, ok := picker.properties[TimePickerValue]; ok {
			delete(picker.properties, TimePickerValue)
			time := GetTimePickerValue(picker, "")
			if picker.created {
				picker.session.runScript(fmt.Sprintf(`setInputValue('%s', '%s')`, picker.htmlID(), time.Format(timeFormat)))
			}
			for _, listener := range picker.timeChangedListeners {
				listener(picker, time)
			}
		} else {
			return
		}

	default:
		picker.viewData.remove(tag)
		return
	}
	picker.propertyChangedEvent(tag)
}

func (picker *timePickerData) Set(tag string, value interface{}) bool {
	return picker.set(picker.normalizeTag(tag), value)
}

func (picker *timePickerData) set(tag string, value interface{}) bool {
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
				if time, err := time.Parse(timeFormat, text); err == nil {
					picker.properties[tag] = value
					return time, true
				}
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
					updateProperty(picker.htmlID(), Min, time.Format(timeFormat), picker.session)
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
					updateProperty(picker.htmlID(), Max, time.Format(timeFormat), picker.session)
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case TimePickerStep:
		oldStep := GetTimePickerStep(picker, "")
		if picker.setIntProperty(TimePickerStep, value) {
			if step := GetTimePickerStep(picker, ""); oldStep != step {
				if picker.created {
					if step > 0 {
						updateProperty(picker.htmlID(), Step, strconv.Itoa(step), picker.session)
					} else {
						removeProperty(picker.htmlID(), Step, picker.session)
					}
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case TimePickerValue:
		oldTime := GetTimePickerValue(picker, "")
		if time, ok := setTimeValue(TimePickerMax); ok {
			if time != oldTime {
				if picker.created {
					picker.session.runScript(fmt.Sprintf(`setInputValue('%s', '%s')`, picker.htmlID(), time.Format(timeFormat)))
				}
				for _, listener := range picker.timeChangedListeners {
					listener(picker, time)
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case TimeChangedEvent:
		switch value := value.(type) {
		case func(TimePicker, time.Time):
			picker.timeChangedListeners = []func(TimePicker, time.Time){value}

		case func(time.Time):
			fn := func(view TimePicker, time time.Time) {
				value(time)
			}
			picker.timeChangedListeners = []func(TimePicker, time.Time){fn}

		case []func(TimePicker, time.Time):
			picker.timeChangedListeners = value

		case []func(time.Time):
			listeners := make([]func(TimePicker, time.Time), len(value))
			for i, val := range value {
				if val == nil {
					notCompatibleType(tag, val)
					return false
				}

				listeners[i] = func(view TimePicker, time time.Time) {
					val(time)
				}
			}
			picker.timeChangedListeners = listeners

		case []interface{}:
			listeners := make([]func(TimePicker, time.Time), len(value))
			for i, val := range value {
				if val == nil {
					notCompatibleType(tag, val)
					return false
				}

				switch val := val.(type) {
				case func(TimePicker, time.Time):
					listeners[i] = val

				case func(time.Time):
					listeners[i] = func(view TimePicker, time time.Time) {
						val(time)
					}

				default:
					notCompatibleType(tag, val)
					return false
				}
			}
			picker.timeChangedListeners = listeners
		}
		picker.propertyChangedEvent(tag)
		return true

	default:
		return picker.viewData.set(tag, value)
	}
	return false
}

func (picker *timePickerData) Get(tag string) interface{} {
	return picker.get(picker.normalizeTag(tag))
}

func (picker *timePickerData) get(tag string) interface{} {
	switch tag {
	case TimeChangedEvent:
		return picker.timeChangedListeners

	default:
		return picker.viewData.get(tag)
	}
}

func (picker *timePickerData) htmlTag() string {
	return "input"
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
	buffer.WriteString(GetTimePickerValue(picker, "").Format(timeFormat))
	buffer.WriteByte('"')

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)
	if picker.getRaw(ClickEvent) == nil {
		buffer.WriteString(` onclick="stopEventPropagation(this, event)"`)
	}
}

func (picker *timePickerData) htmlDisabledProperties(self View, buffer *strings.Builder) {
	if IsDisabled(self, "") {
		buffer.WriteString(` disabled`)
	}
	picker.viewData.htmlDisabledProperties(self, buffer)
}

func (picker *timePickerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if value, err := time.Parse(timeFormat, text); err == nil {
				oldValue := GetTimePickerValue(picker, "")
				picker.properties[TimePickerValue] = value
				if value != oldValue {
					for _, listener := range picker.timeChangedListeners {
						listener(picker, value)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

func getTimeProperty(view View, mainTag, shortTag string) (time.Time, bool) {
	valueToTime := func(value interface{}) (time.Time, bool) {
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

		if value, ok := valueFromStyle(view, shortTag); ok {
			if result, ok := valueToTime(value); ok {
				return result, true
			}
		}
	}

	return time.Now(), false
}

// GetTimePickerMin returns the min time of TimePicker subview and "true" as the second value if the min time is set,
// "false" as the second value otherwise.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTimePickerMin(view View, subviewID string) (time.Time, bool) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		return getTimeProperty(view, TimePickerMin, Min)
	}
	return time.Now(), false
}

// GetTimePickerMax returns the max time of TimePicker subview and "true" as the second value if the min time is set,
// "false" as the second value otherwise.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTimePickerMax(view View, subviewID string) (time.Time, bool) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		return getTimeProperty(view, TimePickerMax, Max)
	}
	return time.Now(), false
}

// GetTimePickerStep returns the time changing step in seconds of TimePicker subview.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTimePickerStep(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return 60
	}

	result, ok := intStyledProperty(view, TimePickerStep, 60)
	if !ok {
		result, _ = intStyledProperty(view, Step, 60)
	}

	if result < 0 {
		return 60
	}
	return result
}

// GetTimePickerValue returns the time of TimePicker subview.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTimePickerValue(view View, subviewID string) time.Time {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return time.Now()
	}
	time, _ := getTimeProperty(view, TimePickerValue, Value)
	return time
}

// GetTimeChangedListeners returns the TimeChangedListener list of an TimePicker subview.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTimeChangedListeners(view View, subviewID string) []func(TimePicker, time.Time) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(TimeChangedEvent); value != nil {
			if listeners, ok := value.([]func(TimePicker, time.Time)); ok {
				return listeners
			}
		}
	}
	return []func(TimePicker, time.Time){}
}
