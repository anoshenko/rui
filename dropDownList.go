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

// DropDownList - the interface of a drop-down list view
type DropDownList interface {
	View
	getItems() []string
}

type dropDownListData struct {
	viewData
	items            []string
	disabledItems    []any
	dropDownListener []func(DropDownList, int)
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
	list.items = []string{}
	list.disabledItems = []any{}
	list.dropDownListener = []func(DropDownList, int){}
}

func (list *dropDownListData) String() string {
	return getViewString(list)
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

	case DropDownEvent:
		if len(list.dropDownListener) > 0 {
			list.dropDownListener = []func(DropDownList, int){}
			list.propertyChangedEvent(tag)
		}

	case Current:
		oldCurrent := GetCurrent(list)
		delete(list.properties, Current)
		if oldCurrent != 0 {
			if list.created {
				list.session.runScript(fmt.Sprintf(`selectDropDownListItem('%s', %d)`, list.htmlID(), 0))
			}
			list.onSelectedItemChanged(0)
		}

	default:
		list.viewData.remove(tag)
		return
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
		return list.setDisabledItems(value)

	case DropDownEvent:
		listeners, ok := valueToEventListeners[DropDownList, int](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(DropDownList, int){}
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
				list.session.runScript(fmt.Sprintf(`selectDropDownListItem('%s', %d)`, list.htmlID(), current))
			}
			list.onSelectedItemChanged(current)
		}
		return true
	}

	return list.viewData.set(tag, value)
}

func (list *dropDownListData) setItems(value any) bool {
	switch value := value.(type) {
	case string:
		list.items = []string{value}

	case []string:
		list.items = value

	case []DataValue:
		list.items = make([]string, 0, len(value))
		for _, val := range value {
			if !val.IsObject() {
				list.items = append(list.items, val.Value())
			}
		}

	case []fmt.Stringer:
		list.items = make([]string, len(value))
		for i, str := range value {
			list.items[i] = str.String()
		}

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
					notCompatibleType(Items, value)
					return false
				}
			}
		}

		list.items = items

	default:
		notCompatibleType(Items, value)
		return false
	}

	if list.created {
		updateInnerHTML(list.htmlID(), list.session)
	}

	list.propertyChangedEvent(Items)
	return true
}

func (list *dropDownListData) setDisabledItems(value any) bool {
	switch value := value.(type) {
	case []int:
		list.disabledItems = make([]any, len(value))
		for i, n := range value {
			list.disabledItems[i] = n
		}

	case []any:
		disabledItems := make([]any, len(value))
		for i, val := range value {
			if val == nil {
				notCompatibleType(DisabledItems, value)
				return false
			}

			switch val := val.(type) {
			case string:
				if isConstantName(val) {
					disabledItems[i] = val
				} else {
					n, err := strconv.Atoi(val)
					if err != nil {
						notCompatibleType(DisabledItems, value)
						return false
					}
					disabledItems[i] = n
				}
			default:
				if n, ok := isInt(val); ok {
					disabledItems[i] = n
				} else {
					notCompatibleType(DisabledItems, value)
					return false
				}
			}

		}
		list.disabledItems = disabledItems

	case string:
		values := strings.Split(value, ",")
		disabledItems := make([]any, len(values))
		for i, str := range values {
			str = strings.Trim(str, " ")
			if str == "" {
				notCompatibleType(DisabledItems, value)
				return false
			}
			if isConstantName(str) {
				disabledItems[i] = str
			} else {
				n, err := strconv.Atoi(str)
				if err != nil {
					notCompatibleType(DisabledItems, value)
					return false
				}
				disabledItems[i] = n
			}
		}
		list.disabledItems = disabledItems

	case []DataValue:
		disabledItems := make([]string, 0, len(value))
		for _, val := range value {
			if !val.IsObject() {
				disabledItems = append(disabledItems, val.Value())
			}
		}
		return list.setDisabledItems(disabledItems)

	default:
		notCompatibleType(DisabledItems, value)
		return false
	}

	if list.created {
		updateInnerHTML(list.htmlID(), list.session)
	}

	list.propertyChangedEvent(Items)
	return true

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
		}
	}
}

func (list *dropDownListData) htmlProperties(self View, buffer *strings.Builder) {
	list.viewData.htmlProperties(self, buffer)
	buffer.WriteString(` size="1" onchange="dropDownListEvent(this, event)"`)
}

func (list *dropDownListData) htmlDisabledProperties(self View, buffer *strings.Builder) {
	list.viewData.htmlDisabledProperties(self, buffer)
	if IsDisabled(list) {
		buffer.WriteString(`disabled`)
	}
}

func (list *dropDownListData) onSelectedItemChanged(number int) {
	for _, listener := range list.dropDownListener {
		listener(list, number)
	}
	list.propertyChangedEvent(Current)
}

func (list *dropDownListData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "itemSelected":
		if text, ok := data.PropertyValue("number"); ok {
			if number, err := strconv.Atoi(text); err == nil {
				if GetCurrent(list) != number && number >= 0 && number < len(list.items) {
					list.properties[Current] = number
					list.onSelectedItemChanged(number)
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
func GetDropDownListeners(view View, subviewID ...string) []func(DropDownList, int) {
	return getEventListeners[DropDownList, int](view, subviewID, DropDownEvent)
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

// GetDropDownDisabledItems return the list of DropDownList disabled item indexes.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetDropDownDisabledItems(view View, subviewID ...string) []int {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if value := view.Get(DisabledItems); value != nil {
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
