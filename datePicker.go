package rui

import (
	"strconv"
	"strings"
	"time"
)

// Constants for [DatePicker] specific properties and events.
const (
	// DateChangedEvent is the constant for "date-changed" property tag.
	//
	// Used by `DatePicker`.
	// Occur when date picker value has been changed.
	//
	// General listener format:
	// `func(picker rui.DatePicker, newDate, oldDate time.Time)`.
	//
	// where:
	// picker - Interface of a date picker which generated this event,
	// newDate - New date value,
	// oldDate - Old date value.
	//
	// Allowed listener formats:
	// `func(picker rui.DatePicker, newDate time.Time)`,
	// `func(newDate, oldDate time.Time)`,
	// `func(newDate time.Time)`,
	// `func(picker rui.DatePicker)`,
	// `func()`.
	DateChangedEvent = "date-changed"

	// DatePickerMin is the constant for "date-picker-min" property tag.
	//
	// Used by `DatePicker`.
	// Minimum date value.
	//
	// Supported types: `time.Time`, `string`.
	//
	// Internal type is `time.Time`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - values of this type parsed and converted to `time.Time`. The following formats are supported:
	// "YYYYMMDD" - "20240102".
	// "Mon-DD-YYYY" - "Jan-02-24".
	// "Mon-DD-YY" - "Jan-02-2024".
	// "DD-Mon-YYYY" - "02-Jan-2024".
	// "YYYY-MM-DD" - "2024-01-02".
	// "Month DD, YYYY" - "January 02, 2024".
	// "DD Month YYYY" - "02 January 2024".
	// "MM/DD/YYYY" - "01/02/2024".
	// "MM/DD/YY" - "01/02/24".
	// "MMDDYY" - "010224".
	DatePickerMin = "date-picker-min"

	// DatePickerMax is the constant for "date-picker-max" property tag.
	//
	// Used by `DatePicker`.
	// Maximum date value.
	//
	// Supported types: `time.Time`, `string`.
	//
	// Internal type is `time.Time`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - values of this type parsed and converted to `time.Time`. The following formats are supported:
	// "YYYYMMDD" - "20240102".
	// "Mon-DD-YYYY" - "Jan-02-24".
	// "Mon-DD-YY" - "Jan-02-2024".
	// "DD-Mon-YYYY" - "02-Jan-2024".
	// "YYYY-MM-DD" - "2024-01-02".
	// "Month DD, YYYY" - "January 02, 2024".
	// "DD Month YYYY" - "02 January 2024".
	// "MM/DD/YYYY" - "01/02/2024".
	// "MM/DD/YY" - "01/02/24".
	// "MMDDYY" - "010224".
	DatePickerMax = "date-picker-max"

	// DatePickerStep is the constant for "date-picker-step" property tag.
	//
	// Used by `DatePicker`.
	// Date change step in days.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// >= `0` or >= "0" - Step value in days used to increment or decrement date.
	DatePickerStep = "date-picker-step"

	// DatePickerValue is the constant for "date-picker-value" property tag.
	//
	// Used by `DatePicker`.
	// Current value.
	//
	// Supported types: `time.Time`, `string`.
	//
	// Internal type is `time.Time`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - values of this type parsed and converted to `time.Time`. The following formats are supported:
	// "YYYYMMDD" - "20240102".
	// "Mon-DD-YYYY" - "Jan-02-24".
	// "Mon-DD-YY" - "Jan-02-2024".
	// "DD-Mon-YYYY" - "02-Jan-2024".
	// "YYYY-MM-DD" - "2024-01-02".
	// "Month DD, YYYY" - "January 02, 2024".
	// "DD Month YYYY" - "02 January 2024".
	// "MM/DD/YYYY" - "01/02/2024".
	// "MM/DD/YY" - "01/02/24".
	// "MMDDYY" - "010224".
	DatePickerValue = "date-picker-value"

	dateFormat = "2006-01-02"
)

// DatePicker represent a DatePicker view
type DatePicker interface {
	View
}

type datePickerData struct {
	viewData
	dataList
	dateChangedListeners []func(DatePicker, time.Time, time.Time)
}

// NewDatePicker create new DatePicker object and return it
func NewDatePicker(session Session, params Params) DatePicker {
	view := new(datePickerData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newDatePicker(session Session) View {
	return NewDatePicker(session, nil)
}

func (picker *datePickerData) init(session Session) {
	picker.viewData.init(session)
	picker.tag = "DatePicker"
	picker.hasHtmlDisabled = true
	picker.dateChangedListeners = []func(DatePicker, time.Time, time.Time){}
	picker.dataListInit()
}

func (picker *datePickerData) String() string {
	return getViewString(picker, nil)
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
			picker.dateChangedListeners = []func(DatePicker, time.Time, time.Time){}
			picker.propertyChangedEvent(tag)
		}
		return

	case DatePickerMin:
		delete(picker.properties, DatePickerMin)
		if picker.created {
			picker.session.removeProperty(picker.htmlID(), Min)
		}

	case DatePickerMax:
		delete(picker.properties, DatePickerMax)
		if picker.created {
			picker.session.removeProperty(picker.htmlID(), Max)
		}

	case DatePickerStep:
		delete(picker.properties, DatePickerStep)
		if picker.created {
			picker.session.removeProperty(picker.htmlID(), Step)
		}

	case DatePickerValue:
		if _, ok := picker.properties[DatePickerValue]; ok {
			oldDate := GetDatePickerValue(picker)
			delete(picker.properties, DatePickerValue)
			date := GetDatePickerValue(picker)
			if picker.created {
				picker.session.callFunc("setInputValue", picker.htmlID(), date.Format(dateFormat))
			}
			for _, listener := range picker.dateChangedListeners {
				listener(picker, date, oldDate)
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

func (picker *datePickerData) Set(tag string, value any) bool {
	return picker.set(picker.normalizeTag(tag), value)
}

func (picker *datePickerData) set(tag string, value any) bool {
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
					picker.session.updateProperty(picker.htmlID(), Min, date.Format(dateFormat))
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
					picker.session.updateProperty(picker.htmlID(), Max, date.Format(dateFormat))
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case DatePickerStep:
		oldStep := GetDatePickerStep(picker)
		if picker.setIntProperty(DatePickerStep, value) {
			if step := GetDatePickerStep(picker); oldStep != step {
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

	case DatePickerValue:
		oldDate := GetDatePickerValue(picker)
		if date, ok := setTimeValue(DatePickerValue); ok {
			if date != oldDate {
				if picker.created {
					picker.session.callFunc("setInputValue", picker.htmlID(), date.Format(dateFormat))
				}
				for _, listener := range picker.dateChangedListeners {
					listener(picker, date, oldDate)
				}
				picker.propertyChangedEvent(tag)
			}
			return true
		}

	case DateChangedEvent:
		listeners, ok := valueToEventWithOldListeners[DatePicker, time.Time](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(DatePicker, time.Time, time.Time){}
		}
		picker.dateChangedListeners = listeners
		picker.propertyChangedEvent(tag)
		return true

	case DataList:
		return picker.setDataList(picker, value, picker.created)

	default:
		return picker.viewData.set(tag, value)
	}
	return false
}

func (picker *datePickerData) Get(tag string) any {
	return picker.get(picker.normalizeTag(tag))
}

func (picker *datePickerData) get(tag string) any {
	switch tag {
	case DateChangedEvent:
		return picker.dateChangedListeners

	case DataList:
		return picker.dataList.dataList

	default:
		return picker.viewData.get(tag)
	}
}

func (picker *datePickerData) htmlTag() string {
	return "input"
}

func (picker *datePickerData) htmlSubviews(self View, buffer *strings.Builder) {
	picker.dataListHtmlSubviews(self, buffer)
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
	buffer.WriteString(GetDatePickerValue(picker).Format(dateFormat))
	buffer.WriteByte('"')

	buffer.WriteString(` oninput="editViewInputEvent(this)"`)
	if picker.getRaw(ClickEvent) == nil {
		buffer.WriteString(` onclick="stopEventPropagation(this, event)"`)
	}

	picker.dataListHtmlProperties(picker, buffer)
}

func (picker *datePickerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if value, err := time.Parse(dateFormat, text); err == nil {
				oldValue := GetDatePickerValue(picker)
				picker.properties[DatePickerValue] = value
				if value != oldValue {
					for _, listener := range picker.dateChangedListeners {
						listener(picker, value, oldValue)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

func getDateProperty(view View, mainTag, shortTag string) (time.Time, bool) {
	valueToTime := func(value any) (time.Time, bool) {
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
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDatePickerMin(view View, subviewID ...string) (time.Time, bool) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		return getDateProperty(view, DatePickerMin, Min)
	}
	return time.Now(), false
}

// GetDatePickerMax returns the max date of DatePicker subview and "true" as the second value if the min date is set,
// "false" as the second value otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDatePickerMax(view View, subviewID ...string) (time.Time, bool) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		return getDateProperty(view, DatePickerMax, Max)
	}
	return time.Now(), false
}

// GetDatePickerStep returns the date changing step in days of DatePicker subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDatePickerStep(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, DatePickerStep, 0)
}

// GetDatePickerValue returns the date of DatePicker subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDatePickerValue(view View, subviewID ...string) time.Time {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view == nil {
		return time.Now()
	}
	date, _ := getDateProperty(view, DatePickerValue, Value)
	return date
}

// GetDateChangedListeners returns the DateChangedListener list of an DatePicker subview.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDateChangedListeners(view View, subviewID ...string) []func(DatePicker, time.Time, time.Time) {
	return getEventWithOldListeners[DatePicker, time.Time](view, subviewID, DateChangedEvent)
}
