package rui

import "strings"

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
	DataList = "data-list"
)

type dataList struct {
	dataList     []string
	dataListHtml bool
}

func (list *dataList) dataListInit() {
	list.dataList = []string{}
}

func (list *dataList) dataListID(view View) string {
	return view.htmlID() + "-datalist"
}

func (list *dataList) normalizeDataListTag(tag string) string {
	switch tag {
	case "datalist":
		return DataList
	}

	return tag
}

func (list *dataList) setDataList(view View, value any, created bool) bool {
	items, ok := anyToStringArray(value)
	if !ok {
		notCompatibleType(DataList, value)
		return false
	}

	list.dataList = items
	if created {
		session := view.Session()
		dataListID := list.dataListID(view)
		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)

		if list.dataListHtml {
			list.dataListItemsHtml(buffer)
			session.updateInnerHTML(dataListID, buffer.String())
		} else {
			list.dataListHtmlCode(view, buffer)
			session.appendToInnerHTML(view.parentHTMLID(), buffer.String())
			list.dataListHtml = true
			session.updateProperty(view.htmlID(), "list", dataListID)
		}
	}

	return true
}

func (list *dataList) dataListHtmlSubviews(view View, buffer *strings.Builder) {
	if len(list.dataList) > 0 {
		list.dataListHtmlCode(view, buffer)
		list.dataListHtml = true
	} else {
		list.dataListHtml = false
	}
}

func (list *dataList) dataListHtmlCode(view View, buffer *strings.Builder) {
	buffer.WriteString(`<datalist id="`)
	buffer.WriteString(list.dataListID(view))
	buffer.WriteString(`">`)
	list.dataListItemsHtml(buffer)
	buffer.WriteString(`</datalist>`)
}

func (list *dataList) dataListItemsHtml(buffer *strings.Builder) {
	for _, text := range list.dataList {
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
}

func (list *dataList) dataListHtmlProperties(view View, buffer *strings.Builder) {
	if len(list.dataList) > 0 {
		buffer.WriteString(` list="`)
		buffer.WriteString(list.dataListID(view))
		buffer.WriteString(`"`)
	}
}

// GetDataList returns the data list of an editor.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDataList(view View, subviewID ...string) []string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if value := view.Get(DataList); value != nil {
			if list, ok := value.([]string); ok {
				return list
			}
		}
	}

	return []string{}
}
