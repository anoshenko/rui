package rui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// DataList is the constant for "data-list" property tag.
	//
	// Used by `ColorPicker`, `DatePicker`, `EditView`, `NumberPicker`, `TimePicker`.
	//
	// Usage in `ColorPicker`:
	// List of pre-defined colors.
	//
	// Supported types: `[]string`, `string`, `[]fmt.Stringer`, `[]Color`, `[]SizeUnit`, `[]AngleUnit`, `[]any` containing
	// elements of `string`, `fmt.Stringer`, `bool`, `rune`, `float32`, `float64`, `int`, `int8` … `int64`, `uint`, `uint8` …
	// `uint64`.
	//
	// Internal type is `[]string`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - contain single item.
	// `[]string` - an array of items.
	// `[]fmt.Stringer` - an array of objects convertible to a string.
	// `[]Color` - An array of color values which will be converted to a string array.
	// `[]SizeUnit` - an array of size unit values which will be converted to a string array.
	// `[]any` - this array must contain only types which were listed in Types section.
	//
	// Usage in `DatePicker`:
	// List of predefined dates. If we set this property, date picker may have a drop-down menu with a list of these values.
	// Some browsers may ignore this property, such as Safari for macOS. The value of this property must be an array of
	// strings in the format "YYYY-MM-DD".
	//
	// Supported types: `[]string`, `string`, `[]fmt.Stringer`, `[]Color`, `[]SizeUnit`, `[]AngleUnit`, `[]any` containing
	// elements of `string`, `fmt.Stringer`, `bool`, `rune`, `float32`, `float64`, `int`, `int8` … `int64`, `uint`, `uint8` …
	// `uint64`.
	//
	// Internal type is `[]string`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - contain single item.
	// `[]string` - an array of items.
	// `[]fmt.Stringer` - an array of objects convertible to a string.
	// `[]Color` - An array of color values which will be converted to a string array.
	// `[]SizeUnit` - an array of size unit values which will be converted to a string array.
	// `[]any` - this array must contain only types which were listed in Types section.
	//
	// Usage in `EditView`:
	// Array of recommended values.
	//
	// Supported types: `[]string`, `string`, `[]fmt.Stringer`, `[]Color`, `[]SizeUnit`, `[]AngleUnit`, `[]any` containing
	// elements of `string`, `fmt.Stringer`, `bool`, `rune`, `float32`, `float64`, `int`, `int8` … `int64`, `uint`, `uint8` …
	// `uint64`.
	//
	// Internal type is `[]string`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - contain single item.
	// `[]string` - an array of items.
	// `[]fmt.Stringer` - an array of objects convertible to a string.
	// `[]Color` - An array of color values which will be converted to a string array.
	// `[]SizeUnit` - an array of size unit values which will be converted to a string array.
	// `[]any` - this array must contain only types which were listed in Types section.
	//
	// Usage in `NumberPicker`:
	// Specify an array of recommended values.
	//
	// Supported types: `[]string`, `string`, `[]fmt.Stringer`, `[]Color`, `[]SizeUnit`, `[]AngleUnit`, `[]float`, `[]int`,
	// `[]bool`, `[]any` containing elements of `string`, `fmt.Stringer`, `bool`, `rune`, `float32`, `float64`, `int`, `int8`
	// … `int64`, `uint`, `uint8` … `uint64`.
	//
	// Internal type is `[]string`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - must contain integer or floating point number, converted to `[]string`.
	// `[]string` - an array of strings which must contain integer or floating point numbers, stored as is.
	// `[]fmt.Stringer` - object which implement this interface must contain integer or floating point numbers, converted to a `[]string`.
	// `[]Color` - an array of color values, converted to `[]string`.
	// `[]SizeUnit` - an array of size unit, converted to `[]string`.
	// `[]AngleUnit` - an array of angle unit, converted to `[]string`.
	// `[]float` - converted to `[]string`.
	// `[]int` - converted to `[]string`.
	// `[]bool` - converted to `[]string`.
	// `[]any` - an array which may contain types listed in Types section above, each value will be converted to a `string` and wrapped to array.
	//
	// Usage in `TimePicker`:
	// An array of recommended values. The value of this property must be an array of strings in the format "HH:MM:SS" or
	// "HH:MM".
	//
	// Supported types: `[]string`, `string`, `[]fmt.Stringer`, `[]Color`, `[]SizeUnit`, `[]AngleUnit`, `[]any` containing
	// elements of `string`, `fmt.Stringer`, `bool`, `rune`, `float32`, `float64`, `int`, `int8` … `int64`, `uint`, `uint8` …
	// `uint64`.
	//
	// Internal type is `[]string`, other types converted to it during assignment.
	//
	// Conversion rules:
	// `string` - contain single item.
	// `[]string` - an array of items.
	// `[]fmt.Stringer` - an array of objects convertible to a string.
	// `[]Color` - An array of color values which will be converted to a string array.
	// `[]SizeUnit` - an array of size unit values which will be converted to a string array.
	// `[]any` - this array must contain only types which were listed in Types section.
	DataList PropertyName = "data-list"
)

func dataListID(view View) string {
	return view.htmlID() + "-datalist"
}

func normalizeDataListTag(tag PropertyName) PropertyName {
	switch tag {
	case "datalist":
		return DataList
	}

	return tag
}

func setDataList(properties Properties, value any, dateTimeFormat string) []PropertyName {
	if items, ok := anyToStringArray(value, dateTimeFormat); ok {
		properties.setRaw(DataList, items)
		return []PropertyName{DataList}
	}

	notCompatibleType(DataList, value)
	return nil
}

func anyToStringArray(value any, dateTimeFormat string) ([]string, bool) {

	switch value := value.(type) {
	case string:
		return []string{value}, true

	case []string:
		return value, true

	case []DataValue:
		items := make([]string, 0, len(value))
		for _, val := range value {
			if !val.IsObject() {
				items = append(items, val.Value())
			}
		}
		return items, true

	case []fmt.Stringer:
		items := make([]string, len(value))
		for i, str := range value {
			items[i] = str.String()
		}
		return items, true

	case []Color:
		items := make([]string, len(value))
		for i, str := range value {
			items[i] = str.String()
		}
		return items, true

	case []SizeUnit:
		items := make([]string, len(value))
		for i, str := range value {
			items[i] = str.String()
		}
		return items, true

	case []AngleUnit:
		items := make([]string, len(value))
		for i, str := range value {
			items[i] = str.String()
		}
		return items, true

	case []float32:
		items := make([]string, len(value))
		for i, val := range value {
			items[i] = fmt.Sprintf("%g", float64(val))
		}
		return items, true

	case []float64:
		items := make([]string, len(value))
		for i, val := range value {
			items[i] = fmt.Sprintf("%g", val)
		}
		return items, true

	case []int:
		return intArrayToStringArray(value), true

	case []uint:
		return intArrayToStringArray(value), true

	case []int8:
		return intArrayToStringArray(value), true

	case []uint8:
		return intArrayToStringArray(value), true

	case []int16:
		return intArrayToStringArray(value), true

	case []uint16:
		return intArrayToStringArray(value), true

	case []int32:
		return intArrayToStringArray(value), true

	case []uint32:
		return intArrayToStringArray(value), true

	case []int64:
		return intArrayToStringArray(value), true

	case []uint64:
		return intArrayToStringArray(value), true

	case []bool:
		items := make([]string, len(value))
		for i, val := range value {
			if val {
				items[i] = "true"
			} else {
				items[i] = "false"
			}
		}
		return items, true

	case []time.Time:
		if dateTimeFormat == "" {
			dateTimeFormat = dateFormat + " " + timeFormat
		}

		items := make([]string, len(value))
		for i, val := range value {
			items[i] = val.Format(dateTimeFormat)
		}
		return items, true

	case []any:
		items := make([]string, 0, len(value))
		for _, v := range value {
			switch val := v.(type) {
			case string:
				items = append(items, val)

			case fmt.Stringer:
				items = append(items, val.String())

			case bool:
				if val {
					items = append(items, "true")
				} else {
					items = append(items, "false")
				}

			case float32:
				items = append(items, fmt.Sprintf("%g", float64(val)))

			case float64:
				items = append(items, fmt.Sprintf("%g", val))

			case rune:
				items = append(items, string(val))

			default:
				if n, ok := isInt(v); ok {
					items = append(items, strconv.Itoa(n))
				} else {
					return []string{}, false
				}
			}
		}

		return items, true
	}

	return []string{}, false
}

func getDataListProperty(properties Properties) []string {
	if value := properties.getRaw(DataList); value != nil {
		if items, ok := value.([]string); ok {
			return items
		}
	}
	return nil
}

func dataListHtmlSubviews(view View, buffer *strings.Builder, normalizeItem func(text string, session Session) string) {
	if items := getDataListProperty(view); len(items) > 0 {
		session := view.Session()
		buffer.WriteString(`<datalist id="`)
		buffer.WriteString(dataListID(view))
		buffer.WriteString(`">`)
		for _, text := range items {
			text = normalizeItem(text, session)

			if strings.ContainsRune(text, '"') {
				text = strings.ReplaceAll(text, `"`, `&#34;`)
			}
			if strings.ContainsRune(text, '\n') {
				text = strings.ReplaceAll(text, "\n", `\n`)
			}
			buffer.WriteString(`<option value="`)
			buffer.WriteString(text)
			buffer.WriteString(`"></option>`)
		}
		buffer.WriteString(`</datalist>`)
	}
}

func dataListHtmlProperties(view View, buffer *strings.Builder) {
	if len(getDataListProperty(view)) > 0 {
		buffer.WriteString(` list="`)
		buffer.WriteString(dataListID(view))
		buffer.WriteString(`"`)
	}
}

// GetDataList returns the data list of an editor.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDataList(view View, subviewID ...string) []string {
	if view = getSubview(view, subviewID); view != nil {
		return getDataListProperty(view)
	}

	return []string{}
}
