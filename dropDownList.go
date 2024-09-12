package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// DropDownEvent is the constant for "drop-down-event" property tag.
// The "drop-down-event" event occurs when a list item becomes selected.
// The main listener format: func(DropDownList, int), where the second argument is the item index.
const DropDownEvent = "drop-down-event"

// DropDownList represent a DropDownList view
type DropDownList interface {
	View
	getItems() []string
}

type dropDownListData struct {
	viewData
	items            []string
	disabledItems    []any
	itemSeparators   []any
	dropDownListener []func(DropDownList, int, int)
}

// NewDropDownList create new DropDownList object and return it
func NewDropDownList(session Session, params Params) DropDownList {
	view := new(dropDownListData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newDropDownList(session Session) View {
	return NewDropDownList(session, nil)
}

func (list *dropDownListData) init(session Session) {
	list.viewData.init(session)
	list.tag = "DropDownList"
	list.hasHtmlDisabled = true
	list.items = []string{}
	list.disabledItems = []any{}
	list.itemSeparators = []any{}
	list.dropDownListener = []func(DropDownList, int, int){}
}

func (list *dropDownListData) String() string {
	return getViewString(list, nil)
}

func (list *dropDownListData) Focusable() bool {
	return true
}

func (list *dropDownListData) Remove(tag string) {
	list.remove(strings.ToLower(tag))
}

func (list *dropDownListData) remove(tag string) {
	switch tag {
	case Items:
		if len(list.items) > 0 {
			list.items = []string{}
			if list.created {
				updateInnerHTML(list.htmlID(), list.session)
			}
			list.propertyChangedEvent(tag)
		}

	case DisabledItems:
		if len(list.disabledItems) > 0 {
			list.disabledItems = []any{}
			if list.created {
				updateInnerHTML(list.htmlID(), list.session)
			}
			list.propertyChangedEvent(tag)
		}

	case ItemSeparators, "separators":
		if len(list.itemSeparators) > 0 {
			list.itemSeparators = []any{}
			if list.created {
				updateInnerHTML(list.htmlID(), list.session)
			}
			list.propertyChangedEvent(ItemSeparators)
		}

	case DropDownEvent:
		if len(list.dropDownListener) > 0 {
			list.dropDownListener = []func(DropDownList, int, int){}
			list.propertyChangedEvent(tag)
		}

	case Current:
		oldCurrent := GetCurrent(list)
		delete(list.properties, Current)
		if oldCurrent != 0 {
			if list.created {
				list.session.callFunc("selectDropDownListItem", list.htmlID(), 0)
			}
			list.onSelectedItemChanged(0, oldCurrent)
		}

	default:
		list.viewData.remove(tag)
	}
}

func (list *dropDownListData) Set(tag string, value any) bool {
	return list.set(strings.ToLower(tag), value)
}

func (list *dropDownListData) set(tag string, value any) bool {
	if value == nil {
		list.remove(tag)
		return true
	}

	switch tag {
	case Items:
		return list.setItems(value)

	case DisabledItems:
		items, ok := list.parseIndicesArray(value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		}
		list.disabledItems = items
		if list.created {
			updateInnerHTML(list.htmlID(), list.session)
		}
		list.propertyChangedEvent(tag)
		return true

	case ItemSeparators, "separators":
		items, ok := list.parseIndicesArray(value)
		if !ok {
			notCompatibleType(ItemSeparators, value)
			return false
		}
		list.itemSeparators = items
		if list.created {
			updateInnerHTML(list.htmlID(), list.session)
		}
		list.propertyChangedEvent(ItemSeparators)
		return true

	case DropDownEvent:
		listeners, ok := valueToEventWithOldListeners[DropDownList, int](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(DropDownList, int, int){}
		}
		list.dropDownListener = listeners
		list.propertyChangedEvent(tag)
		return true

	case Current:
		oldCurrent := GetCurrent(list)
		if !list.setIntProperty(Current, value) {
			return false
		}

		if current := GetCurrent(list); oldCurrent != current {
			if list.created {
				list.session.callFunc("selectDropDownListItem", list.htmlID(), current)
			}
			list.onSelectedItemChanged(current, oldCurrent)
		}
		return true
	}

	return list.viewData.set(tag, value)
}

func (list *dropDownListData) setItems(value any) bool {
	items, ok := anyToStringArray(value)
	if !ok {
		notCompatibleType(Items, value)
		return false
	}

	list.items = items
	if list.created {
		updateInnerHTML(list.htmlID(), list.session)
	}

	list.propertyChangedEvent(Items)
	return true
}

func intArrayToStringArray[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64](array []T) []string {
	items := make([]string, len(array))
	for i, val := range array {
		items[i] = strconv.Itoa(int(val))
	}
	return items
}

func anyToStringArray(value any) ([]string, bool) {

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

func (list *dropDownListData) parseIndicesArray(value any) ([]any, bool) {
	switch value := value.(type) {
	case []int:
		items := make([]any, len(value))
		for i, n := range value {
			items[i] = n
		}
		return items, true

	case []any:
		items := make([]any, len(value))
		for i, val := range value {
			if val == nil {
				return nil, false
			}

			switch val := val.(type) {
			case string:
				if isConstantName(val) {
					items[i] = val
				} else {
					n, err := strconv.Atoi(val)
					if err != nil {
						return nil, false
					}
					items[i] = n
				}
			default:
				if n, ok := isInt(val); ok {
					items[i] = n
				} else {
					return nil, false
				}
			}

		}
		return items, true

	case string:
		values := strings.Split(value, ",")
		items := make([]any, len(values))
		for i, str := range values {
			str = strings.Trim(str, " ")
			if str == "" {
				return nil, false
			}
			if isConstantName(str) {
				items[i] = str
			} else {
				n, err := strconv.Atoi(str)
				if err != nil {
					return nil, false
				}
				items[i] = n
			}
		}
		return items, true

	case []DataValue:
		items := make([]any, 0, len(value))
		for _, val := range value {
			if !val.IsObject() {
				items = append(items, val.Value())
			}
		}
		return list.parseIndicesArray(items)
	}

	return nil, false
}

func (list *dropDownListData) Get(tag string) any {
	return list.get(strings.ToLower(tag))
}

func (list *dropDownListData) get(tag string) any {
	switch tag {
	case Items:
		return list.items

	case DisabledItems:
		return list.disabledItems

	case ItemSeparators:
		return list.itemSeparators

	case Current:
		result, _ := intProperty(list, Current, list.session, 0)
		return result

	case DropDownEvent:
		return list.dropDownListener
	}

	return list.viewData.get(tag)
}

func (list *dropDownListData) getItems() []string {
	return list.items
}

func (list *dropDownListData) htmlTag() string {
	return "select"
}

func (list *dropDownListData) htmlSubviews(self View, buffer *strings.Builder) {
	if list.items != nil {
		current := GetCurrent(list)
		notTranslate := GetNotTranslate(list)
		disabledItems := GetDropDownDisabledItems(list)
		separators := GetDropDownItemSeparators(list)
		for i, item := range list.items {
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

func (list *dropDownListData) onSelectedItemChanged(number, old int) {
	for _, listener := range list.dropDownListener {
		listener(list, number, old)
	}
	list.propertyChangedEvent(Current)
}

func (list *dropDownListData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "itemSelected":
		if text, ok := data.PropertyValue("number"); ok {
			if number, err := strconv.Atoi(text); err == nil {
				if GetCurrent(list) != number && number >= 0 && number < len(list.items) {
					old := GetCurrent(list)
					list.properties[Current] = number
					list.onSelectedItemChanged(number, old)
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
	return getEventWithOldListeners[DropDownList, int](view, subviewID, DropDownEvent)
}

// GetDropDownItems return the DropDownList items list.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDropDownItems(view View, subviewID ...string) []string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}
	if view != nil {
		if list, ok := view.(DropDownList); ok {
			return list.getItems()
		}
	}
	return []string{}
}

func getIndicesArray(view View, tag string) []int {
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
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	return getIndicesArray(view, DisabledItems)
}

// GetDropDownItemSeparators return an array of indices of DropDownList items after which a separator should be added.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDropDownItemSeparators(view View, subviewID ...string) []int {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	return getIndicesArray(view, ItemSeparators)
}
