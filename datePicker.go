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
	DateChangedEvent PropertyName = "date-changed"

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
	DatePickerMin PropertyName = "date-picker-min"

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
	DatePickerMax PropertyName = "date-picker-max"

	// DatePickerStep is the constant for "date-picker-step" property tag.
	//
	// Used by `DatePicker`.
	// Date change step in days.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// >= `0` or >= "0" - Step value in days used to increment or decrement date.
	DatePickerStep PropertyName = "date-picker-step"

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
	DatePickerValue PropertyName = "date-picker-value"

	dateFormat = "2006-01-02"
)

// DatePicker represent a DatePicker view
type DatePicker interface {
	View
}

type datePickerData struct {
	viewData
}

// NewDatePicker create new DatePicker object and return it
func NewDatePicker(session Session, params Params) DatePicker {
	view := new(datePickerData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newDatePicker(session Session) View {
	return new(datePickerData) // NewDatePicker(session, nil)
}

func (picker *datePickerData) init(session Session) {
	picker.viewData.init(session)
	picker.tag = "DatePicker"
	picker.hasHtmlDisabled = true
	picker.normalize = normalizeDatePickerTag
	picker.set = datePickerSet
	picker.changed = datePickerPropertyChanged
}

func (picker *datePickerData) Focusable() bool {
	return true
}

func normalizeDatePickerTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Type, Min, Max, Step, Value:
		return "date-picker-" + tag
	}

	return normalizeDataListTag(tag)
}

func stringToDate(value string) (time.Time, bool) {
	format := "20060102"
	if strings.ContainsRune(value, '-') {
		if part := strings.Split(value, "-"); len(part) == 3 {
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
	} else if strings.ContainsRune(value, ' ') {
		if part := strings.Split(value, " "); len(part) == 3 {
			if part[0] != "" && part[0][0] > '9' {
				format = "January 02, 2006"
			} else {
				format = "02 January 2006"
			}
		}
	} else if strings.ContainsRune(value, '/') {
		if part := strings.Split(value, "/"); len(part) == 3 {
			if len(part[2]) == 2 {
				format = "01/02/06"
			} else {
				format = "01/02/2006"
			}
		}
	} else if len(value) == 6 {
		format = "010206"
	}

	if date, err := time.Parse(format, value); err == nil {
		return date, true
	}
	return time.Now(), false
}

func datePickerSet(view View, tag PropertyName, value any) []PropertyName {

	setDateValue := func(tag PropertyName) []PropertyName {
		switch value := value.(type) {
		case time.Time:
			view.setRaw(tag, value)
			return []PropertyName{tag}

		case string:
			if isConstantName(value) {
				view.setRaw(tag, value)
				return []PropertyName{tag}
			}

			if date, ok := stringToDate(value); ok {
				view.setRaw(tag, date)
				return []PropertyName{tag}
			}
		}

		notCompatibleType(tag, value)
		return nil
	}

	switch tag {
	case DatePickerMin, DatePickerMax:
		return setDateValue(tag)

	case DatePickerStep:
		return setIntProperty(view, DatePickerStep, value)

	case DatePickerValue:
		view.setRaw("old-date", GetDatePickerValue(view))
		return setDateValue(tag)

	case DateChangedEvent:
		return setEventWithOldListener[DatePicker, time.Time](view, tag, value)

	case DataList:
		return setDataList(view, value, dateFormat)
	}

	return viewSet(view, tag, value)
}

func datePickerPropertyChanged(view View, tag PropertyName) {

	session := view.Session()

	switch tag {

	case DatePickerMin:
		if date, ok := GetDatePickerMin(view); ok {
			session.updateProperty(view.htmlID(), "min", date.Format(dateFormat))
		} else {
			session.removeProperty(view.htmlID(), "min")
		}

	case DatePickerMax:
		if date, ok := GetDatePickerMax(view); ok {
			session.updateProperty(view.htmlID(), "max", date.Format(dateFormat))
		} else {
			session.removeProperty(view.htmlID(), "max")
		}

	case DatePickerStep:
		if step := GetDatePickerStep(view); step > 0 {
			session.updateProperty(view.htmlID(), "step", strconv.Itoa(step))
		} else {
			session.removeProperty(view.htmlID(), "step")
		}

	case DatePickerValue:
		date := GetDatePickerValue(view)
		session.callFunc("setInputValue", view.htmlID(), date.Format(dateFormat))

		if listeners := GetDateChangedListeners(view); len(listeners) > 0 {
			oldDate := time.Now()
			if value := view.getRaw("old-date"); value != nil {
				if date, ok := value.(time.Time); ok {
					oldDate = date
				}
			}
			for _, listener := range listeners {
				listener(view, date, oldDate)
			}
		}

	default:
		viewPropertyChanged(view, tag)
	}
}

func (picker *datePickerData) htmlTag() string {
	return "input"
}

func (picker *datePickerData) htmlSubviews(self View, buffer *strings.Builder) {
	dataListHtmlSubviews(self, buffer, func(text string, session Session) string {
		text, _ = session.resolveConstants(text)
		if date, ok := stringToDate(text); ok {
			return date.Format(dateFormat)
		}
		return text
	})
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

	dataListHtmlProperties(picker, buffer)
}

func (picker *datePickerData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "textChanged":
		if text, ok := data.PropertyValue("text"); ok {
			if value, err := time.Parse(dateFormat, text); err == nil {
				oldValue := GetDatePickerValue(picker)
				picker.properties[DatePickerValue] = value
				if value != oldValue {
					for _, listener := range GetDateChangedListeners(picker) {
						listener(picker, value, oldValue)
					}
					if listener, ok := picker.changeListener[DatePickerValue]; ok {
						listener(picker, DatePickerValue)
					}
				}
			}
		}
		return true
	}

	return picker.viewData.handleCommand(self, command, data)
}

func getDateProperty(view View, mainTag, shortTag PropertyName) (time.Time, bool) {
	valueToTime := func(value any) (time.Time, bool) {
		if value != nil {
			switch value := value.(type) {
			case time.Time:
				return value, true

			case string:
				if text, ok := view.Session().resolveConstants(value); ok {
					if result, ok := stringToDate(text); ok {
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
