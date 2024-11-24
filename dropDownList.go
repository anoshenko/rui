package rui

import (
	"strconv"
	"strings"
)

// DropDownEvent is the constant for "drop-down-event" property tag.
//
// Used by `DropDownList`.
// Occur when a list item becomes selected.
//
// General listener format:
// `func(list rui.DropDownList, index int)`.
//
// where:
// list - Interface of a drop down list which generated this event,
// index - Index of a newly selected item.
//
// Allowed listener formats:
const DropDownEvent PropertyName = "drop-down-event"

// DropDownList represent a DropDownList view
type DropDownList interface {
	View
}

type dropDownListData struct {
	viewData
}

// NewDropDownList create new DropDownList object and return it
func NewDropDownList(session Session, params Params) DropDownList {
	view := new(dropDownListData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newDropDownList(session Session) View {
	return new(dropDownListData)
}

func (list *dropDownListData) init(session Session) {
	list.viewData.init(session)
	list.tag = "DropDownList"
	list.hasHtmlDisabled = true
	list.normalize = normalizeDropDownListTag
	list.set = list.setFunc
	list.changed = list.propertyChanged
}

func (list *dropDownListData) Focusable() bool {
	return true
}

func normalizeDropDownListTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	if tag == "separators" {
		return ItemSeparators
	}
	return tag
}

func (list *dropDownListData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Items:
		if items, ok := anyToStringArray(value, ""); ok {
			return setArrayPropertyValue(list, tag, items)
		}
		notCompatibleType(Items, value)
		return nil

	case DisabledItems, ItemSeparators:
		if items, ok := parseIndicesArray(value); ok {
			return setArrayPropertyValue(list, tag, items)
		}
		notCompatibleType(tag, value)
		return nil

	case DropDownEvent:
		return setTwoArgEventListener[DropDownList, int](list, tag, value)

	case Current:
		list.setRaw("old-current", GetCurrent(list))
		return setIntProperty(list, Current, value)
	}

	return list.viewData.setFunc(tag, value)
}

func (list *dropDownListData) propertyChanged(tag PropertyName) {
	switch tag {
	case Items, DisabledItems, ItemSeparators:
		updateInnerHTML(list.htmlID(), list.Session())

	case Current:
		current := GetCurrent(list)
		list.Session().callFunc("selectDropDownListItem", list.htmlID(), current)

		oldCurrent, _ := intProperty(list, "old-current", list.Session(), -1)
		for _, listener := range GetDropDownListeners(list) {
			listener(list, current, oldCurrent)
		}

	default:
		list.viewData.propertyChanged(tag)
	}
}

func intArrayToStringArray[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64](array []T) []string {
	items := make([]string, len(array))
	for i, val := range array {
		items[i] = strconv.Itoa(int(val))
	}
	return items
}

func parseIndicesArray(value any) ([]any, bool) {
	switch value := value.(type) {
	case int:
		return []any{value}, true

	case []int:
		items := make([]any, len(value))
		for i, n := range value {
			items[i] = n
		}
		return items, true

	case []any:
		items := make([]any, 0, len(value))
		for _, val := range value {
			if val != nil {
				switch val := val.(type) {
				case string:
					if isConstantName(val) {
						items = append(items, val)
					} else if n, err := strconv.Atoi(val); err == nil {
						items = append(items, n)
					} else {
						return nil, false
					}

				default:
					if n, ok := isInt(val); ok {
						items = append(items, n)
					} else {
						return nil, false
					}
				}
			}
		}
		return items, true

	case []string:
		items := make([]any, 0, len(value))
		for _, str := range value {
			if str = strings.Trim(str, " \t"); str != "" {
				if isConstantName(str) {
					items = append(items, str)
				} else if n, err := strconv.Atoi(str); err == nil {
					items = append(items, n)
				} else {
					return nil, false
				}
			}
		}
		return items, true

	case string:
		return parseIndicesArray(strings.Split(value, ","))

	case []DataValue:
		items := make([]string, 0, len(value))
		for _, val := range value {
			if !val.IsObject() {
				items = append(items, val.Value())
			}
		}
		return parseIndicesArray(items)
	}

	return nil, false
}

func (list *dropDownListData) htmlTag() string {
	return "select"
}

func (list *dropDownListData) htmlSubviews(self View, buffer *strings.Builder) {
	if items := GetDropDownItems(list); len(items) > 0 {
		current := GetCurrent(list)
		notTranslate := GetNotTranslate(list)
		disabledItems := GetDropDownDisabledItems(list)
		separators := GetDropDownItemSeparators(list)
		for i, item := range items {
			disabled := false
			for _, index := range disabledItems {
				if i == index {
					disabled = true
					break
				}
			}

			if disabled {
				buffer.WriteString("<option disabled>")
			} else if i == current {
				buffer.WriteString("<option selected>")
			} else {
				buffer.WriteString("<option>")
			}
			if !notTranslate {
				item, _ = list.session.GetString(item)
			}

			buffer.WriteString(item)
			buffer.WriteString("</option>")
			for _, index := range separators {
				if i == index {
					buffer.WriteString("<hr>")
					break
				}
			}
		}
	}
}

func (list *dropDownListData) htmlProperties(self View, buffer *strings.Builder) {
	list.viewData.htmlProperties(self, buffer)
	buffer.WriteString(` size="1" onchange="dropDownListEvent(this, event)"`)
}

func (list *dropDownListData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "itemSelected":
		if text, ok := data.PropertyValue("number"); ok {
			if number, err := strconv.Atoi(text); err == nil {
				items := GetDropDownItems(list)
				if GetCurrent(list) != number && number >= 0 && number < len(items) {
					old := GetCurrent(list)
					list.properties[Current] = number
					for _, listener := range GetDropDownListeners(list) {
						listener(list, number, old)
					}
					if listener, ok := list.changeListener[Current]; ok {
						listener(list, Current)
					}
				}
			} else {
				ErrorLog(err.Error())
			}
		}

	default:
		return list.viewData.handleCommand(self, command, data)
	}
	return true
}

// GetDropDownListeners returns the "drop-down-event" listener list. If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDropDownListeners(view View, subviewID ...string) []func(DropDownList, int, int) {
	return getTwoArgEventListeners[DropDownList, int](view, subviewID, DropDownEvent)
}

// GetDropDownItems return the DropDownList items list.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDropDownItems(view View, subviewID ...string) []string {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.Get(Items); value != nil {
			if items, ok := value.([]string); ok {
				return items
			}
		}
	}
	return []string{}
}

func getIndicesArray(view View, tag PropertyName) []int {
	if view != nil {
		if value := view.Get(tag); value != nil {
			if values, ok := value.([]any); ok {
				count := len(values)
				if count > 0 {
					result := make([]int, 0, count)
					for _, value := range values {
						switch value := value.(type) {
						case int:
							result = append(result, value)

						case string:
							if value != "" && value[0] == '@' {
								if val, ok := view.Session().Constant(value[1:]); ok {
									if n, err := strconv.Atoi(val); err == nil {
										result = append(result, n)
									}
								}
							}
						}
					}
					return result
				}
			}
		}
	}
	return []int{}
}

// GetDropDownDisabledItems return an array of disabled(non selectable) items indices of DropDownList.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDropDownDisabledItems(view View, subviewID ...string) []int {
	view = getSubview(view, subviewID)
	return getIndicesArray(view, DisabledItems)
}

// GetDropDownItemSeparators return an array of indices of DropDownList items after which a separator should be added.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDropDownItemSeparators(view View, subviewID ...string) []int {
	view = getSubview(view, subviewID)
	return getIndicesArray(view, ItemSeparators)
}
