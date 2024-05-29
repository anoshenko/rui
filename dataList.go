package rui

import "strings"

const (
	// DataList is the constant for the "data-list" property tag.
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
