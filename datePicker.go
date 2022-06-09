package rui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DateChangedEvent = "date-changed"
	DatePickerMin    = "date-picker-min"
	DatePickerMax    = "date-picker-max"
	DatePickerStep   = "date-picker-step"
	DatePickerValue  = "date-picker-value"
	dateFormat       = "2006-01-02"
)

// DatePicker - DatePicker view
type DatePicker interface {
	View
}

type datePickerData struct {
	viewData
	dateChangedListeners []func(DatePicker, time.Time)
}

// NewDatePicker create new DatePicker object and return it
func NewDatePicker(session Session, params Params) DatePicker {
	view := new(datePickerData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newDatePicker(session Session) View {
	return NewDatePicker(session, nil)
}

func (picker *datePickerData) Init(session Session) {
	picker.viewData.Init(session)
	picker.tag = "DatePicker"
	picker.dateChangedListeners = []func(DatePicker, time.Time){}
}

func (picker *datePickerData) String() string {
	return getViewString(picker)
}

func (picker *datePickerData) Focusable() bool {
	return true
}

func (picker *datePickerData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Type, Min, Max, Step, Value:
		return "date-picker-" + tag
	}

	return tag
}

func (picker *datePickerData) Remove(tag string) {
	picker.remove(picker.normalizeTag(tag))
}

func (picker *datePickerData) remove(tag string) {
	switch tag {
	case DateChangedEvent:
		if len(picker.dateChangedListeners) > 0 {
			picker.dateChangedListeners = []func(DatePicker, time.Time){}
			picker.propertyChangedEvent(tag)
		}
		return

	case DatePickerMin:
		delete(picker.properties, DatePickerMin)
		if picker.created {
			removeProperty(picker.htmlID(), Min, picker.session)
		}

	case DatePickerMax:
		delete(picker.properties, DatePickerMax)
		if picker.created {
			removeProperty(picker.htmlID(), Max, picker.session)
		}

	case DatePickerStep:
		delete(picker.properties, DatePickerStep)
		if picker.created {
			removeProperty(picker.htmlID(), Step, picker.session)
		}

	case DatePickerValue:
		if _, ok := picker.properties[DatePickerValue]; ok {
			delete(picker.properties, DatePickerValue)
			date := GetDatePickerValue(picker, "")
			if picker.created {
				picker.session.runScript(fmt.Sprintf(`setInputValue('%s', '%s')`, picker.htmlID(), date.Format(dateFormat)))
			}
			for _, listener := range picker.dateChangedListeners {
				listener(picker, date)
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

func (picker *datePickerData) Set(tag string, value interface{}) bool {
	return picker.set(picker.normalizeTag(tag), value)
}

func (picker *datePickerData) set(tag string, value interface{}) bool {
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
				format := "20060102"
				if strings.ContainsRune(text, '-') {
					if part := strings.Split(text, "-"); len(part) == 3 {
						if part[0] != "" && part[0][0] > '9' {
							if len(part[2]) == 2 {
								format = "Jan-02-06"
							} else {
								format = "Jan-02-2006"
							}
						} else if part[1] != "" && part[1][0] > '9' {
							format = "02-Jan-2006"
						} else {
							format = "2006-01-02"
						}
					}
				} else if strings.ContainsRune(text, ' ') {
					if part := strings.Split(text, " "); len(part) == 3 {
						if part[0] != "" && part[0][0] > '9' {
							format = "January 02, 2006"
						} else {
							format = "02 January 2006"
						}
					}
				} else if strings.ContainsRune(text, '/') {
					if part := strings.Split(text, "/"); len(part) == 3 {
						if len(part[2]) == 2 {
							format = "01/02/06"
						} else {
							format = "01/02/2006"
						}
					}
				} else if len(text) == 6 {
					format = "010206"
				}

				if date, err := time.Parse(format, text); err == nil {
					picker.properties[tag] = value
					return date, true
				}
			}
		}

		notCompatibleType(tag, value)
		return time.Now(), false
	}

	switch tag {
	case DatePickerMin:
		old, oldOK := getDateProperty(picker, DatePickerMin, Min)
		if date, ok := setTimeValue(DatePickerMin); ok {
			if !oldOK || date != old {
				if picker.created {
					updateProperty(picker.htmlID(), Min, date.Format(dateFormat), picker.session)
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case DatePickerMax:
		old, oldOK := getDateProperty(picker, DatePickerMax, Max)
		if date, ok := setTimeValue(DatePickerMax); ok {
			if !oldOK || date != old {
				if picker.created {
					updateProperty(picker.htmlID(), Max, date.Format(dateFormat), picker.session)
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case DatePickerStep:
		oldStep := GetDatePickerStep(picker, "")
		if picker.setIntProperty(DatePickerStep, value) {
			if step := GetDatePickerStep(picker, ""); oldStep != step {
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

	case DatePickerValue:
		oldDate := GetDatePickerValue(picker, "")
		if date, ok := setTimeValue(DatePickerValue); ok {
			if date != oldDate {
				if picker.created {
					picker.session.runScript(fmt.Sprintf(`setInputValue('%s', '%s')`, picker.htmlID(), date.Format(dateFormat)))
				}
				for _, listener := range picker.dateChangedListeners {
					listener(picker, date)
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case DateChangedEvent:
		switch value := value.(type) {
		case func(DatePicker, time.Time):
			picker.dateChangedListeners = []func(DatePicker, time.Time){value}

		case func(time.Time):
			fn := func(view DatePicker, date time.Time) {
				value(date)
			}
			picker.dateChangedListeners = []func(DatePicker, time.Time){fn}

		case []func(DatePicker, time.Time):
			picker.dateChangedListeners = value

		case []func(time.Time):
			listeners := make([]func(DatePicker, time.Time), len(value))
			for i, val := range value {
				if val == nil {
					notCompatibleType(tag, val)
					return false
				}

				listeners[i] = func(view DatePicker, date time.Time) {
					val(date)
				}
			}
			picker.dateChangedListeners = listeners

		case []interface{}:
			listeners := make([]func(DatePicker, time.Time), len(value))
			for i, val := range value {
				if val == nil {
					notCompatibleType(tag, val)
					return false
				}

				switch val := val.(type) {
				case func(DatePicker, time.Time):
					listeners[i] = val

				case func(time.Time):
					listeners[i] = func(view DatePicker, date time.Time) {
						val(date)
					}

				default:
					notCompatibleType(tag, val)
					return false
				}
			}
			picker.dateChangedListeners = listeners
		}
		picker.propertyChangedEvent(tag)
		return true

	default:
		return picker.viewData.set(tag, value)
	}
	return false
}

func (picker *datePickerData) Get(tag string) interface{} {
	return picker.get(picker.normalizeTag(tag))
}

func (picker *datePickerData) get(tag string) interface{} {
	switch tag {
	case DateChangedEvent:
		return picker.dateChangedListeners

	default:
		return picker.viewData.get(tag)
	}
}

func (picker *datePickerData) htmlTag() string {
	return "input"
}

func (picker *datePickerData) htmlProperties(self View, buffer *strings.Builder) {
	picker.viewData.htmlProperties(self, buffer)

	buffer.WriteString(` type="date"`)

	if min, ok := getDateProperty(picker, DatePickerMin, Min); ok {
		buffer.WriteString(` min="`)
		buffer.WriteString(min.Format(dateFormat))
		buffer.WriteByte('"')
	}

	if max, ok := getDateProperty(picker, DatePickerMax, Max); ok {
		buffer.WriteString(` max="`)
		buffer.WriteString(max.Format(dateFormat))
		buffer.WriteByte('"')
	}

	if step, ok := intProperty(picker, DatePickerStep, picker.Session(), 0); ok && step > 0 {
		buffer.WriteString(` step="`)
		buffer.WriteString(strconv.Itoa(step))
		buffer.WriteByte('"')
	}

	buffer.WriteString(` value="`)
	buffer.WriteString(GetDatePickerValue(picker, "").Format(dateFormat))
	buffer.WriteByte('"')

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)
	if picker.getRaw(ClickEvent) == nil {
		buffer.WriteString(` onclick="stopEventPropagation(this, event)"`)
	}
}

func (picker *datePickerData) htmlDisabledProperties(self View, buffer *strings.Builder) {
	if IsDisabled(self, "") {
		buffer.WriteString(` disabled`)
	}
	picker.viewData.htmlDisabledProperties(self, buffer)
}

func (picker *datePickerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if value, err := time.Parse(dateFormat, text); err == nil {
				oldValue := GetDatePickerValue(picker, "")
				picker.properties[DatePickerValue] = value
				if value != oldValue {
					for _, listener := range picker.dateChangedListeners {
						listener(picker, value)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

func getDateProperty(view View, mainTag, shortTag string) (time.Time, bool) {
	valueToTime := func(value interface{}) (time.Time, bool) {
		if value != nil {
			switch value := value.(type) {
			case time.Time:
				return value, true

			case string:
				if text, ok := view.Session().resolveConstants(value); ok {
					if result, err := time.Parse(dateFormat, text); err == nil {
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

// GetDatePickerMin returns the min date of DatePicker subview and "true" as the second value if the min date is set,
// "false" as the second value otherwise.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetDatePickerMin(view View, subviewID string) (time.Time, bool) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		return getDateProperty(view, DatePickerMin, Min)
	}
	return time.Now(), false
}

// GetDatePickerMax returns the max date of DatePicker subview and "true" as the second value if the min date is set,
// "false" as the second value otherwise.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetDatePickerMax(view View, subviewID string) (time.Time, bool) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		return getDateProperty(view, DatePickerMax, Max)
	}
	return time.Now(), false
}

// GetDatePickerStep returns the date changing step in days of DatePicker subview.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetDatePickerStep(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if result, _ := intStyledProperty(view, DatePickerStep, 0); result >= 0 {
			return result
		}
	}
	return 0
}

// GetDatePickerValue returns the date of DatePicker subview.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetDatePickerValue(view View, subviewID string) time.Time {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return time.Now()
	}
	date, _ := getDateProperty(view, DatePickerValue, Value)
	return date
}

// GetDateChangedListeners returns the DateChangedListener list of an DatePicker subview.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetDateChangedListeners(view View, subviewID string) []func(DatePicker, time.Time) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(DateChangedEvent); value != nil {
			if listeners, ok := value.([]func(DatePicker, time.Time)); ok {
				return listeners
			}
		}
	}
	return []func(DatePicker, time.Time){}
}
