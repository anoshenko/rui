package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const DropDownEvent = "drop-down-event"

// DropDownList - the interface of a drop-down list view
type DropDownList interface {
	View
	getItems() []string
}

type dropDownListData struct {
	viewData
	items            []string
	dropDownListener []func(DropDownList, int)
}

// NewDropDownList create new DropDownList object and return it
func NewDropDownList(session Session, params Params) DropDownList {
	view := new(dropDownListData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newDropDownList(session Session) View {
	return NewDropDownList(session, nil)
}

func (list *dropDownListData) Init(session Session) {
	list.viewData.Init(session)
	list.tag = "DropDownList"
	list.items = []string{}
	list.dropDownListener = []func(DropDownList, int){}
}

func (list *dropDownListData) Remove(tag string) {
	list.remove(strings.ToLower(tag))
}

func (list *dropDownListData) remove(tag string) {
	switch tag {
	case Items:
		if len(list.items) > 0 {
			list.items = []string{}
			updateInnerHTML(list.htmlID(), list.session)
		}

	case Current:
		list.set(Current, 0)

	case DropDownEvent:
		if len(list.dropDownListener) > 0 {
			list.dropDownListener = []func(DropDownList, int){}
		}

	default:
		list.viewData.remove(tag)
	}
}

func (list *dropDownListData) Set(tag string, value interface{}) bool {
	return list.set(strings.ToLower(tag), value)
}

func (list *dropDownListData) set(tag string, value interface{}) bool {
	switch tag {
	case Items:
		return list.setItems(value)

	case Current:
		oldCurrent := GetDropDownCurrent(list, "")
		if !list.setIntProperty(Current, value) {
			return false
		}

		if !list.session.ignoreViewUpdates() {
			current := GetDropDownCurrent(list, "")
			if oldCurrent != current {
				list.session.runScript(fmt.Sprintf(`selectDropDownListItem('%s', %d)`, list.htmlID(), current))
				list.onSelectedItemChanged(current)
			}
		}
		return true

	case DropDownEvent:
		return list.setDropDownListener(value)
	}

	return list.viewData.set(tag, value)
}

func (list *dropDownListData) setItems(value interface{}) bool {
	switch value := value.(type) {
	case string:
		list.items = []string{value}

	case []string:
		list.items = value

	case []DataValue:
		list.items = []string{}
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

	case []interface{}:
		items := []string{}
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

	if !list.session.ignoreViewUpdates() {
		updateInnerHTML(list.htmlID(), list.session)
	}
	return true
}

func (list *dropDownListData) setDropDownListener(value interface{}) bool {
	switch value := value.(type) {
	case func(DropDownList, int):
		list.dropDownListener = []func(DropDownList, int){value}
		return true

	case func(int):
		list.dropDownListener = []func(DropDownList, int){func(list DropDownList, index int) {
			value(index)
		}}
		return true

	case []func(DropDownList, int):
		list.dropDownListener = value
		return true

	case []func(int):
		listeners := make([]func(DropDownList, int), len(value))
		for i, val := range value {
			if val == nil {
				notCompatibleType(DropDownEvent, value)
				return false
			}
			listeners[i] = func(list DropDownList, index int) {
				val(index)
			}
		}
		list.dropDownListener = listeners
		return true

	case []interface{}:
		listeners := make([]func(DropDownList, int), len(value))
		for i, val := range value {
			if val == nil {
				notCompatibleType(DropDownEvent, value)
				return false
			}
			switch val := val.(type) {
			case func(DropDownList, int):
				listeners[i] = val

			case func(int):
				listeners[i] = func(list DropDownList, index int) {
					val(index)
				}

			default:
				notCompatibleType(DropDownEvent, value)
				return false
			}
			list.dropDownListener = listeners
		}
		return true
	}

	notCompatibleType(DropDownEvent, value)
	return false
}

func (list *dropDownListData) Get(tag string) interface{} {
	return list.get(strings.ToLower(tag))
}

func (list *dropDownListData) get(tag string) interface{} {
	switch tag {
	case Items:
		return list.items

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
		current := GetDropDownCurrent(list, "")
		notTranslate := GetNotTranslate(list, "")
		for i, item := range list.items {
			if i == current {
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
}

func (list *dropDownListData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "itemSelected":
		if text, ok := data.PropertyValue("number"); ok {
			if number, err := strconv.Atoi(text); err == nil {
				if GetDropDownCurrent(list, "") != number && number >= 0 && number < len(list.items) {
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

func GetDropDownListeners(view View) []func(DropDownList, int) {
	if value := view.Get(DropDownEvent); value != nil {
		if listeners, ok := value.([]func(DropDownList, int)); ok {
			return listeners
		}
	}
	return []func(DropDownList, int){}
}

// func GetDropDownItems return the view items list
func GetDropDownItems(view View, subviewID string) []string {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if list, ok := view.(DropDownList); ok {
			return list.getItems()
		}
	}
	return []string{}
}

// func GetDropDownCurrentItem return the number of the selected item
func GetDropDownCurrent(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		result, _ := intProperty(view, Current, view.Session(), 0)
		return result
	}
	return 0
}
