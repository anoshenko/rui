package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// ListItemClickedEvent is the constant for "list-item-clicked" property tag.
	// The "list-item-clicked" event occurs when the user clicks on an item in the list.
	// The main listener format: func(ListView, int), where the second argument is the item index.
	ListItemClickedEvent = "list-item-clicked"
	// ListItemSelectedEvent is the constant for "list-item-selected" property tag.
	// The "list-item-selected" event occurs when a list item becomes selected.
	// The main listener format: func(ListView, int), where the second argument is the item index.
	ListItemSelectedEvent = "list-item-selected"
	// ListItemCheckedEvent is the constant for "list-item-checked" property tag.
	// The "list-item-checked" event occurs when a list item checkbox becomes checked/unchecked.
	// The main listener format: func(ListView, []int), where the second argument is the array of checked item indexes.
	ListItemCheckedEvent = "list-item-checked"
	// ListItemStyle is the constant for "list-item-style" property tag.
	// The "list-item-style" string property defines the style of an unselected item
	ListItemStyle = "list-item-style"
	// CurrentStyle is the constant for "current-style" property tag.
	// The "current-style" string property defines the style of the selected item when the ListView is focused.
	CurrentStyle = "current-style"
	// CurrentInactiveStyle is the constant for "current-inactive-style" property tag.
	// The "current-inactive-style" string property defines the style of the selected item when the ListView is unfocused.
	CurrentInactiveStyle = "current-inactive-style"
)

const (
	// VerticalOrientation is the vertical ListView orientation
	VerticalOrientation = 0
	// HorizontalOrientation is the horizontal ListView orientation
	HorizontalOrientation = 1

	// NoneCheckbox is value of "checkbox" property: no checkbox
	NoneCheckbox = 0
	// SingleCheckbox is value of "checkbox" property: only one item can be checked
	SingleCheckbox = 1
	// MultipleCheckbox is value of "checkbox" property: several items can be checked
	MultipleCheckbox = 2
)

// ListView - the list view interface
type ListView interface {
	View
	ParanetView
	// ReloadListViewData updates ListView content
	ReloadListViewData()

	getCheckedItems() []int
	getItemFrames() []Frame
}

type listViewData struct {
	viewData
	adapter           ListAdapter
	clickedListeners  []func(ListView, int)
	selectedListeners []func(ListView, int)
	checkedListeners  []func(ListView, []int)
	items             []View
	itemFrame         []Frame
	checkedItem       []int
}

// NewListView creates the new list view
func NewListView(session Session, params Params) ListView {
	view := new(listViewData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newListView(session Session) View {
	return NewListView(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (listView *listViewData) Init(session Session) {
	listView.viewData.Init(session)
	listView.tag = "ListView"
	listView.systemClass = "ruiListView"
	listView.items = []View{}
	listView.itemFrame = []Frame{}
	listView.checkedItem = []int{}
	listView.clickedListeners = []func(ListView, int){}
	listView.selectedListeners = []func(ListView, int){}
	listView.checkedListeners = []func(ListView, []int){}
}

func (listView *listViewData) String() string {
	return getViewString(listView)
}

func (listView *listViewData) Views() []View {
	return listView.items
}

func (listView *listViewData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case HorizontalAlign:
		tag = ItemHorizontalAlign

	case VerticalAlign:
		tag = ItemVerticalAlign

	case "wrap":
		tag = ListWrap
	}
	return tag
}

func (listView *listViewData) Remove(tag string) {
	listView.remove(listView.normalizeTag(tag))
}

func (listView *listViewData) remove(tag string) {
	switch tag {
	case Checked:
		if len(listView.checkedItem) > 0 {
			listView.checkedItem = []int{}
			if listView.created {
				updateInnerHTML(listView.htmlID(), listView.session)
			}
			listView.propertyChangedEvent(tag)
		}

	case Items:
		if listView.adapter != nil {
			listView.adapter = nil
			if listView.created {
				updateInnerHTML(listView.htmlID(), listView.session)
			}
			listView.propertyChangedEvent(tag)
		}

	case Orientation, ListWrap:
		if _, ok := listView.properties[tag]; ok {
			delete(listView.properties, tag)
			if listView.created {
				updateCSSStyle(listView.htmlID(), listView.session)
			}
			listView.propertyChangedEvent(tag)
		}

	case Current:
		current := GetCurrent(listView, "")
		delete(listView.properties, tag)
		if listView.created {
			updateInnerHTML(listView.htmlID(), listView.session)
		}
		if current != -1 {
			for _, listener := range listView.selectedListeners {
				listener(listView, -1)
			}
			listView.propertyChangedEvent(tag)
		}

	case ItemWidth, ItemHeight, ItemHorizontalAlign, ItemVerticalAlign, ItemCheckbox,
		CheckboxHorizontalAlign, CheckboxVerticalAlign, ListItemStyle, CurrentStyle, CurrentInactiveStyle:
		if _, ok := listView.properties[tag]; ok {
			delete(listView.properties, tag)
			if listView.created {
				updateInnerHTML(listView.htmlID(), listView.session)
			}
			listView.propertyChangedEvent(tag)
		}

	case ListItemClickedEvent:
		if len(listView.clickedListeners) > 0 {
			listView.clickedListeners = []func(ListView, int){}
			listView.propertyChangedEvent(tag)
		}

	case ListItemSelectedEvent:
		if len(listView.selectedListeners) > 0 {
			listView.selectedListeners = []func(ListView, int){}
			listView.propertyChangedEvent(tag)
		}

	case ListItemCheckedEvent:
		if len(listView.checkedListeners) > 0 {
			listView.checkedListeners = []func(ListView, []int){}
			listView.propertyChangedEvent(tag)
		}

	default:
		listView.viewData.remove(tag)
	}
}

func (listView *listViewData) Set(tag string, value any) bool {
	return listView.set(listView.normalizeTag(tag), value)
}

func (listView *listViewData) set(tag string, value any) bool {
	if value == nil {
		listView.remove(tag)
		return true
	}

	switch tag {
	case ListItemClickedEvent:
		listeners, ok := valueToEventListeners[ListView, int](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(ListView, int){}
		}
		listView.clickedListeners = listeners
		listView.propertyChangedEvent(tag)
		return true

	case ListItemSelectedEvent:
		listeners, ok := valueToEventListeners[ListView, int](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(ListView, int){}
		}
		listView.selectedListeners = listeners
		listView.propertyChangedEvent(tag)
		return true

	case ListItemCheckedEvent:
		listeners, ok := valueToEventListeners[ListView, []int](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(ListView, []int){}
		}
		listView.checkedListeners = listeners
		listView.propertyChangedEvent(tag)
		return true

	case Checked:
		if !listView.setChecked(value) {
			return false
		}

	case Items:
		if !listView.setItems(value) {
			return false
		}

	case Current:
		oldCurrent := GetCurrent(listView, "")
		if !listView.setIntProperty(Current, value) {
			return false
		}
		current := GetCurrent(listView, "")
		if oldCurrent == current {
			return true
		}

		for _, listener := range listView.selectedListeners {
			listener(listView, current)
		}

	case Orientation, ListWrap, VerticalAlign, HorizontalAlign, Style, StyleDisabled, ItemWidth, ItemHeight:
		result := listView.viewData.set(tag, value)
		if result && listView.created {
			updateInnerHTML(listView.htmlID(), listView.session)
		}
		return result

	case ItemHorizontalAlign, ItemVerticalAlign, ItemCheckbox, CheckboxHorizontalAlign, CheckboxVerticalAlign:
		if !listView.setEnumProperty(tag, value, enumProperties[tag].values) {
			return false
		}

	case ListItemStyle, CurrentStyle, CurrentInactiveStyle:
		switch value := value.(type) {
		case string:
			listView.properties[tag] = value

		default:
			notCompatibleType(tag, value)
			return false
		}

	default:
		return listView.viewData.set(tag, value)
	}

	if listView.created {
		updateInnerHTML(listView.htmlID(), listView.session)
	}
	listView.propertyChangedEvent(tag)
	return true
}

func (listView *listViewData) Get(tag string) any {
	return listView.get(listView.normalizeTag(tag))
}

func (listView *listViewData) get(tag string) any {
	switch tag {
	case ListItemClickedEvent:
		return listView.clickedListeners

	case ListItemSelectedEvent:
		return listView.selectedListeners

	case ListItemCheckedEvent:
		return listView.checkedListeners

	case Checked:
		return listView.checkedItem

	case Items:
		return listView.adapter

	case ListItemStyle:
		return listView.listItemStyle()

	case CurrentStyle:
		return listView.currentStyle()

	case CurrentInactiveStyle:
		return listView.currentInactiveStyle()
	}
	return listView.viewData.get(tag)
}

func (listView *listViewData) setItems(value any) bool {
	switch value := value.(type) {
	case []string:
		listView.adapter = NewTextListAdapter(value, nil)

	case []DataValue:
		hasObject := false
		for _, val := range value {
			if val.IsObject() {
				hasObject = true
				break
			}
		}

		if hasObject {
			items := make([]View, len(value))
			for i, val := range value {
				if val.IsObject() {
					if view := CreateViewFromObject(listView.session, val.Object()); view != nil {
						items[i] = view
					} else {
						return false
					}
				} else {
					items[i] = NewTextView(listView.session, Params{Text: val.Value()})
				}
			}
			listView.adapter = NewViewListAdapter(items)
		} else {
			items := make([]string, len(value))
			for i, val := range value {
				items[i] = val.Value()
			}
			listView.adapter = NewTextListAdapter(items, nil)
		}

	case []any:
		items := make([]View, len(value))
		for i, val := range value {
			switch value := val.(type) {
			case View:
				items[i] = value

			case string:
				items[i] = NewTextView(listView.session, Params{Text: value})

			case fmt.Stringer:
				items[i] = NewTextView(listView.session, Params{Text: value.String()})

			case float32:
				items[i] = NewTextView(listView.session, Params{Text: fmt.Sprintf("%g", float64(value))})

			case float64:
				items[i] = NewTextView(listView.session, Params{Text: fmt.Sprintf("%g", value)})

			default:
				if n, ok := isInt(val); ok {
					items[i] = NewTextView(listView.session, Params{Text: strconv.Itoa(n)})
				} else {
					notCompatibleType(Items, value)
					return false
				}
			}
		}
		listView.adapter = NewViewListAdapter(items)

	case []View:
		listView.adapter = NewViewListAdapter(value)

	case ListAdapter:
		listView.adapter = value

	default:
		notCompatibleType(Items, value)
		return false
	}

	size := listView.adapter.ListSize()
	listView.items = make([]View, size)
	listView.itemFrame = make([]Frame, size)

	return true
}

func (listView *listViewData) setChecked(value any) bool {
	var checked []int
	if value == nil {
		checked = []int{}
	} else {
		switch value := value.(type) {
		case string:
			checked = []int{}
			for _, val := range strings.Split(value, ",") {
				n, err := strconv.Atoi(strings.Trim(val, " \t"))
				if err != nil {
					invalidPropertyValue(Checked, value)
					ErrorLog(err.Error())
					return false
				}
				checked = append(checked, n)
			}

		case int:
			checked = []int{value}

		case []int:
			checked = value

		default:
			return false
		}
	}

	switch GetListViewCheckbox(listView, "") {
	case SingleCheckbox:
		count := len(checked)
		if count > 1 {
			return false
		}

		if len(listView.checkedItem) > 0 &&
			(count == 0 || listView.checkedItem[0] != checked[0]) {
			listView.updateCheckboxItem(listView.checkedItem[0], false)
		}

		if count == 1 {
			listView.updateCheckboxItem(checked[0], true)
		}

	case MultipleCheckbox:
		inSlice := func(n int, slice []int) bool {
			for _, n2 := range slice {
				if n2 == n {
					return true
				}
			}
			return false
		}

		for _, n := range listView.checkedItem {
			if !inSlice(n, checked) {
				listView.updateCheckboxItem(n, false)
			}
		}

		for _, n := range checked {
			if !inSlice(n, listView.checkedItem) {
				listView.updateCheckboxItem(n, true)
			}
		}

	default:
		return false
	}

	listView.checkedItem = checked
	for _, listener := range listView.checkedListeners {
		listener(listView, listView.checkedItem)
	}
	return true
}

func (listView *listViewData) Focusable() bool {
	return true
}

func (listView *listViewData) ReloadListViewData() {
	itemCount := 0
	if listView.adapter != nil {
		itemCount = listView.adapter.ListSize()

		if itemCount != len(listView.items) {
			listView.items = make([]View, itemCount)
			listView.itemFrame = make([]Frame, itemCount)
		}

		for i := 0; i < itemCount; i++ {
			listView.items[i] = listView.adapter.ListItem(i, listView.Session())
		}
	} else if len(listView.items) > 0 {
		listView.items = []View{}
		listView.itemFrame = []Frame{}
	}

	updateInnerHTML(listView.htmlID(), listView.session)
}

func (listView *listViewData) getCheckedItems() []int {
	return listView.checkedItem
}

func (listView *listViewData) getItemFrames() []Frame {
	return listView.itemFrame
}

func (listView *listViewData) itemAlign(self View, buffer *strings.Builder) {
	values := enumProperties[ItemHorizontalAlign].cssValues
	if hAlign := GetListItemHorizontalAlign(listView, ""); hAlign >= 0 && hAlign < len(values) {
		buffer.WriteString(" justify-items: ")
		buffer.WriteString(values[hAlign])
		buffer.WriteRune(';')
	}

	values = enumProperties[ItemVerticalAlign].cssValues
	if vAlign := GetListItemVerticalAlign(listView, ""); vAlign >= 0 && vAlign < len(values) {
		buffer.WriteString(" align-items: ")
		buffer.WriteString(values[vAlign])
		buffer.WriteRune(';')
	}
}

func (listView *listViewData) itemSize(self View, buffer *strings.Builder) {
	if itemWidth := GetListItemWidth(listView, ""); itemWidth.Type != Auto {
		buffer.WriteString(` min-width: `)
		buffer.WriteString(itemWidth.cssString(""))
		buffer.WriteRune(';')
	}

	if itemHeight := GetListItemHeight(listView, ""); itemHeight.Type != Auto {
		buffer.WriteString(` min-height: `)
		buffer.WriteString(itemHeight.cssString(""))
		buffer.WriteRune(';')
	}
}

func (listView *listViewData) getDivs(self View, checkbox, hCheckboxAlign, vCheckboxAlign int) (string, string, string) {
	session := listView.Session()

	contentBuilder := allocStringBuilder()
	defer freeStringBuilder(contentBuilder)

	contentBuilder.WriteString(`<div style="display: grid;`)
	listView.itemAlign(self, contentBuilder)

	onDivBuilder := allocStringBuilder()
	defer freeStringBuilder(onDivBuilder)

	if hCheckboxAlign == CenterAlign {
		if vCheckboxAlign == BottomAlign {
			onDivBuilder.WriteString(`<div style="grid-row: 2 / 3; grid-column: 1 / 2; display: grid; justify-items: center;`)
			contentBuilder.WriteString(` grid-row: 1 / 2; grid-column: 1 / 2;">`)
		} else {
			vCheckboxAlign = TopAlign
			onDivBuilder.WriteString(`<div style="grid-row: 1 / 2; grid-column: 1 / 2; display: grid; justify-items: center;`)
			contentBuilder.WriteString(` grid-row: 2 / 3; grid-column: 1 / 2;">`)
		}
	} else {
		if hCheckboxAlign == RightAlign {
			onDivBuilder.WriteString(`<div style="grid-row: 1 / 2; grid-column: 2 / 3; display: grid;`)
			contentBuilder.WriteString(` grid-row: 1 / 2; grid-column: 1 / 2;">`)
		} else {
			onDivBuilder.WriteString(`<div style="grid-row: 1 / 2; grid-column: 1 / 2; display: grid;`)
			contentBuilder.WriteString(` grid-row: 1 / 2; grid-column: 2 / 3;">`)
		}
		switch vCheckboxAlign {
		case BottomAlign:
			onDivBuilder.WriteString(` align-items: end;`)

		case CenterAlign:
			onDivBuilder.WriteString(` align-items: center;`)

		default:
			onDivBuilder.WriteString(` align-items: start;`)
		}
	}

	onDivBuilder.WriteString(`">`)

	offDivBuilder := allocStringBuilder()
	defer freeStringBuilder(offDivBuilder)

	offDivBuilder.WriteString(onDivBuilder.String())

	if checkbox == SingleCheckbox {
		offDivBuilder.WriteString(session.radiobuttonOffImage())
		onDivBuilder.WriteString(session.radiobuttonOnImage())
	} else {
		offDivBuilder.WriteString(session.checkboxOffImage())
		onDivBuilder.WriteString(session.checkboxOnImage())
	}

	onDivBuilder.WriteString("</div>")
	offDivBuilder.WriteString("</div>")

	return onDivBuilder.String(), offDivBuilder.String(), contentBuilder.String()
}

func (listView *listViewData) checkboxItemDiv(self View, checkbox, hCheckboxAlign, vCheckboxAlign int) string {
	itemStyleBuilder := allocStringBuilder()
	defer freeStringBuilder(itemStyleBuilder)

	itemStyleBuilder.WriteString(`<div style="display: grid; justify-items: stretch; align-items: stretch;`)

	if hCheckboxAlign == CenterAlign {
		if vCheckboxAlign == BottomAlign {
			itemStyleBuilder.WriteString(` grid-template-columns: 1fr; grid-template-rows: 1fr auto;`)
		} else {
			vCheckboxAlign = TopAlign
			itemStyleBuilder.WriteString(` grid-template-columns: 1fr; grid-template-rows: auto 1fr;`)
		}
	} else {
		if hCheckboxAlign == RightAlign {
			itemStyleBuilder.WriteString(` grid-template-columns: 1fr auto; grid-template-rows: 1fr;`)
		} else {
			itemStyleBuilder.WriteString(` grid-template-columns: auto 1fr; grid-template-rows: 1fr;`)
		}
	}

	if gap, ok := sizeConstant(listView.session, "ruiCheckboxGap"); ok && gap.Type != Auto {
		itemStyleBuilder.WriteString(` grid-gap: `)
		itemStyleBuilder.WriteString(gap.cssString("auto"))
		itemStyleBuilder.WriteRune(';')
	}

	itemStyleBuilder.WriteString(`">`)
	return itemStyleBuilder.String()

}

func (listView *listViewData) getItemView(index int) View {
	if listView.adapter == nil || index < 0 || index >= listView.adapter.ListSize() {
		return nil
	}

	size := listView.adapter.ListSize()
	if size != len(listView.items) {
		listView.items = make([]View, size)
	}

	if listView.items[index] == nil {
		listView.items[index] = listView.adapter.ListItem(index, listView.Session())
	}

	return listView.items[index]
}

func (listView *listViewData) itemStyle(tag, defaultStyle string) string {
	if value := listView.getRaw(tag); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = listView.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	if value := valueFromStyle(listView, tag); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = listView.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	return defaultStyle
}

func (listView *listViewData) listItemStyle() string {
	return listView.itemStyle(ListItemStyle, "ruiListItem")
}

func (listView *listViewData) currentStyle() string {
	return listView.itemStyle(CurrentStyle, "ruiListItemFocused")
}

func (listView *listViewData) currentInactiveStyle() string {
	return listView.itemStyle(CurrentInactiveStyle, "ruiListItemSelected")
}

func (listView *listViewData) checkboxSubviews(self View, buffer *strings.Builder, checkbox int) {
	count := listView.adapter.ListSize()
	listViewID := listView.htmlID()

	hCheckboxAlign := GetListViewCheckboxHorizontalAlign(listView, "")
	vCheckboxAlign := GetListViewCheckboxVerticalAlign(listView, "")

	itemDiv := listView.checkboxItemDiv(self, checkbox, hCheckboxAlign, vCheckboxAlign)
	onDiv, offDiv, contentDiv := listView.getDivs(self, checkbox, hCheckboxAlign, vCheckboxAlign)

	current := GetCurrent(listView, "")
	checkedItems := GetListViewCheckedItems(listView, "")
	for i := 0; i < count; i++ {
		buffer.WriteString(`<div id="`)
		buffer.WriteString(listViewID)
		buffer.WriteRune('-')
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(`" class="ruiView `)
		buffer.WriteString(listView.listItemStyle())
		if i == current {
			buffer.WriteRune(' ')
			buffer.WriteString(listView.currentInactiveStyle())
		}
		buffer.WriteString(`" onclick="listItemClickEvent(this, event)" data-left="0" data-top="0" data-width="0" data-height="0" style="display: grid; justify-items: stretch; align-items: stretch;`)
		listView.itemSize(self, buffer)
		buffer.WriteString(`">`)
		buffer.WriteString(itemDiv)

		checked := false
		for _, index := range checkedItems {
			if index == i {
				buffer.WriteString(onDiv)
				checked = true
				break
			}
		}
		if !checked {
			buffer.WriteString(offDiv)
		}
		buffer.WriteString(contentDiv)

		if view := listView.getItemView(i); view != nil {
			//view.setNoResizeEvent()
			viewHTML(view, buffer)
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}

		buffer.WriteString(`</div></div></div>`)
	}
}

func (listView *listViewData) noneCheckboxSubviews(self View, buffer *strings.Builder) {
	count := listView.adapter.ListSize()
	listViewID := listView.htmlID()

	itemStyleBuilder := allocStringBuilder()
	defer freeStringBuilder(itemStyleBuilder)

	itemStyleBuilder.WriteString(`data-left="0" data-top="0" data-width="0" data-height="0" style="max-width: 100%; max-height: 100%; display: grid;`)

	listView.itemAlign(self, itemStyleBuilder)
	listView.itemSize(self, itemStyleBuilder)

	itemStyleBuilder.WriteString(`" onclick="listItemClickEvent(this, event)"`)
	itemStyle := itemStyleBuilder.String()

	current := GetCurrent(listView, "")
	for i := 0; i < count; i++ {
		buffer.WriteString(`<div id="`)
		buffer.WriteString(listViewID)
		buffer.WriteRune('-')
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(`" class="ruiView `)
		buffer.WriteString(listView.listItemStyle())
		if i == current {
			buffer.WriteRune(' ')
			buffer.WriteString(listView.currentInactiveStyle())
		}
		buffer.WriteString(`" `)
		buffer.WriteString(itemStyle)
		buffer.WriteString(`>`)

		if view := listView.getItemView(i); view != nil {
			//view.setNoResizeEvent()
			viewHTML(view, buffer)
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}

		buffer.WriteString(`</div>`)
	}
}

func (listView *listViewData) updateCheckboxItem(index int, checked bool) {

	checkbox := GetListViewCheckbox(listView, "")
	hCheckboxAlign := GetListViewCheckboxHorizontalAlign(listView, "")
	vCheckboxAlign := GetListViewCheckboxVerticalAlign(listView, "")
	onDiv, offDiv, contentDiv := listView.getDivs(listView, checkbox, hCheckboxAlign, vCheckboxAlign)

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`updateInnerHTML('`)
	buffer.WriteString(listView.htmlID())
	buffer.WriteRune('-')
	buffer.WriteString(strconv.Itoa(index))
	buffer.WriteString(`', '`)

	buffer.WriteString(listView.checkboxItemDiv(listView, checkbox, hCheckboxAlign, vCheckboxAlign))
	if checked {
		buffer.WriteString(onDiv)
	} else {
		buffer.WriteString(offDiv)
	}
	buffer.WriteString(contentDiv)

	session := listView.Session()
	if listView.adapter != nil {
		if view := listView.getItemView(index); view != nil {
			view.setNoResizeEvent()
			viewHTML(view, buffer)
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}
	}
	buffer.WriteString(`</div></div>');`)

	session.runScript(buffer.String())
}

func (listView *listViewData) htmlProperties(self View, buffer *strings.Builder) {
	listView.viewData.htmlProperties(self, buffer)
	buffer.WriteString(`onfocus="listViewFocusEvent(this, event)" onblur="listViewBlurEvent(this, event)"`)
	buffer.WriteString(` onkeydown="listViewKeyDownEvent(this, event)" data-focusitemstyle="`)
	buffer.WriteString(listView.currentStyle())
	buffer.WriteString(`" data-bluritemstyle="`)
	buffer.WriteString(listView.currentInactiveStyle())
	buffer.WriteString(`"`)
	current := GetCurrent(listView, "")
	if listView.adapter != nil && current >= 0 && current < listView.adapter.ListSize() {
		buffer.WriteString(` data-current="`)
		buffer.WriteString(listView.htmlID())
		buffer.WriteRune('-')
		buffer.WriteString(strconv.Itoa(current))
		buffer.WriteRune('"')
	}

	listView.viewData.htmlProperties(self, buffer)
}

/*
func (listView *listViewData) cssStyle(self View, builder cssBuilder) {
	listView.viewData.cssStyle(self, builder)

	if GetListWrap(listView, "") != WrapOff {
		switch GetListOrientation(listView, "") {
		case TopDownOrientation, BottomUpOrientation:
			builder.add(`max-height`, `100%`)
		default:
			builder.add(`max-width`, `100%`)
		}
	}
}
*/

func (listView *listViewData) htmlSubviews(self View, buffer *strings.Builder) {
	if listView.adapter == nil || listView.adapter.ListSize() == 0 {
		return
	}

	if !listView.session.ignoreViewUpdates() {
		listView.session.setIgnoreViewUpdates(true)
		defer listView.session.setIgnoreViewUpdates(false)
	}

	buffer.WriteString(`<div style="display: flex; align-content: stretch;`)

	wrap := GetListWrap(listView, "")
	orientation := GetListOrientation(listView, "")
	rows := (orientation == StartToEndOrientation || orientation == EndToStartOrientation)

	if rows {
		if wrap == ListWrapOff {
			buffer.WriteString(` min-width: 100%; height: 100%;`)
		} else {
			buffer.WriteString(` width: 100%; min-height: 100%;`)
		}
	} else {
		if wrap == ListWrapOff {
			buffer.WriteString(` width: 100%; min-height: 100%;`)
		} else {
			buffer.WriteString(` min-width: 100%; height: 100%;`)
		}
	}

	buffer.WriteString(` flex-flow: `)
	buffer.WriteString(enumProperties[Orientation].cssValues[orientation])

	switch wrap {
	case ListWrapOn:
		buffer.WriteString(` wrap;`)

	case ListWrapReverse:
		buffer.WriteString(` wrap-reverse;`)

	default:
		buffer.WriteString(`;`)
	}

	var hAlignTag, vAlignTag string
	if rows {
		hAlignTag = `justify-content`
		vAlignTag = `align-items`
	} else {
		hAlignTag = `align-items`
		vAlignTag = `justify-content`
	}

	value := ""
	switch enumStyledProperty(listView, "", HorizontalAlign, LeftAlign, false) {
	case LeftAlign:
		if (!rows && wrap == ListWrapReverse) || orientation == EndToStartOrientation {
			value = `flex-end`
		} else {
			value = `flex-start`
		}
	case RightAlign:
		if (!rows && wrap == ListWrapReverse) || orientation == EndToStartOrientation {
			value = `flex-start`
		} else {
			value = `flex-end`
		}
	case CenterAlign:
		value = `center`

	case StretchAlign:
		if rows {
			value = `space-between`
		} else {
			value = `stretch`
		}
	}

	if value != "" {
		buffer.WriteRune(' ')
		buffer.WriteString(hAlignTag)
		buffer.WriteString(`: `)
		buffer.WriteString(value)
		buffer.WriteRune(';')
	}

	value = ""
	switch enumStyledProperty(listView, "", VerticalAlign, TopAlign, false) {
	case TopAlign:
		if (rows && wrap == ListWrapReverse) || orientation == BottomUpOrientation {
			value = `flex-end`
		} else {
			value = `flex-start`
		}
	case BottomAlign:
		if (rows && wrap == ListWrapReverse) || orientation == BottomUpOrientation {
			value = `flex-start`
		} else {
			value = `flex-end`
		}
	case CenterAlign:
		value = `center`

	case StretchAlign:
		if rows {
			value = `stretch`
		} else {
			value = `space-between`
		}
	}

	if value != "" {
		buffer.WriteRune(' ')
		buffer.WriteString(vAlignTag)
		buffer.WriteString(`: `)
		buffer.WriteString(value)
		buffer.WriteRune(';')
	}

	buffer.WriteString(`">`)

	checkbox := GetListViewCheckbox(listView, "")
	if checkbox == NoneCheckbox {
		listView.noneCheckboxSubviews(self, buffer)
	} else {
		listView.checkboxSubviews(self, buffer, checkbox)
	}

	buffer.WriteString(`</div>`)
}

func (listView *listViewData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "itemSelected":
		if number, ok := dataIntProperty(data, `number`); ok {
			listView.properties[Current] = number
			for _, listener := range listView.selectedListeners {
				listener(listView, number)
			}
			listView.propertyChangedEvent(Current)
		}

	case "itemUnselected":
		if _, ok := listView.properties[Current]; ok {
			delete(listView.properties, Current)
			for _, listener := range listView.selectedListeners {
				listener(listView, -1)
			}
			listView.propertyChangedEvent(Current)
		}

	case "itemClick":
		listView.onItemClick()

	default:
		return listView.viewData.handleCommand(self, command, data)
	}

	return true
}

func (listView *listViewData) onItemClick() {
	current := GetCurrent(listView, "")
	if current >= 0 && !IsDisabled(listView, "") {
		checkbox := GetListViewCheckbox(listView, "")
	m:
		switch checkbox {
		case SingleCheckbox:
			if len(listView.checkedItem) == 0 {
				listView.checkedItem = []int{current}
				listView.updateCheckboxItem(current, true)
			} else if listView.checkedItem[0] != current {
				listView.updateCheckboxItem(listView.checkedItem[0], false)
				listView.checkedItem[0] = current
				listView.updateCheckboxItem(current, true)
			}

		case MultipleCheckbox:
			for i, index := range listView.checkedItem {
				if index == current {
					listView.updateCheckboxItem(index, false)
					count := len(listView.checkedItem)
					if count == 1 {
						listView.checkedItem = []int{}
					} else if i == 0 {
						listView.checkedItem = listView.checkedItem[1:]
					} else if i == count-1 {
						listView.checkedItem = listView.checkedItem[:i]
					} else {
						listView.checkedItem = append(listView.checkedItem[:i], listView.checkedItem[i+1:]...)
					}
					break m
				}
			}

			listView.updateCheckboxItem(current, true)
			listView.checkedItem = append(listView.checkedItem, current)
		}

		if checkbox != NoneCheckbox {
			for _, listener := range listView.checkedListeners {
				listener(listView, listView.checkedItem)
			}
			listView.propertyChangedEvent(Checked)
		}
		for _, listener := range listView.clickedListeners {
			listener(listView, current)
		}
	}
}

func (listView *listViewData) onItemResize(self View, index string, x, y, width, height float64) {
	n, err := strconv.Atoi(index)
	if err != nil {
		ErrorLog(err.Error())
	} else if n >= 0 && n < len(listView.itemFrame) {
		listView.itemFrame[n] = Frame{Left: x, Top: y, Width: width, Height: height}
	} else {
		ErrorLogF(`Invalid ListView item index: %d`, n)
	}
}

// GetVerticalAlign return the vertical align of a list: TopAlign (0), BottomAlign (1), CenterAlign (2), StretchAlign (3)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetVerticalAlign(view View, subviewID string) int {
	return enumStyledProperty(view, subviewID, VerticalAlign, TopAlign, false)
}

// GetHorizontalAlign return the vertical align of a list/checkbox: LeftAlign (0), RightAlign (1), CenterAlign (2), StretchAlign (3)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetHorizontalAlign(view View, subviewID string) int {
	return enumStyledProperty(view, subviewID, HorizontalAlign, LeftAlign, false)
}

// GetListItemClickedListeners returns a ListItemClickedListener of the ListView.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemClickedListeners(view View, subviewID string) []func(ListView, int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(ListItemClickedEvent); value != nil {
			if result, ok := value.([]func(ListView, int)); ok {
				return result
			}
		}
	}
	return []func(ListView, int){}
}

// GetListItemSelectedListeners returns a ListItemSelectedListener of the ListView.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemSelectedListeners(view View, subviewID string) []func(ListView, int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(ListItemSelectedEvent); value != nil {
			if result, ok := value.([]func(ListView, int)); ok {
				return result
			}
		}
	}
	return []func(ListView, int){}
}

// GetListItemCheckedListeners returns a ListItemCheckedListener of the ListView.
// If there are no listeners then the empty list is returned
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemCheckedListeners(view View, subviewID string) []func(ListView, []int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(ListItemCheckedEvent); value != nil {
			if result, ok := value.([]func(ListView, []int)); ok {
				return result
			}
		}
	}
	return []func(ListView, []int){}
}

// GetListItemWidth returns the width of a ListView item.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemWidth(view View, subviewID string) SizeUnit {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		result, _ := sizeProperty(view, ItemWidth, view.Session())
		return result
	}
	return AutoSize()
}

// GetListItemHeight returns the height of a ListView item.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemHeight(view View, subviewID string) SizeUnit {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		result, _ := sizeProperty(view, ItemHeight, view.Session())
		return result
	}
	return AutoSize()
}

// GetListViewCheckbox returns the ListView checkbox type: NoneCheckbox (0), SingleCheckbox (1), or MultipleCheckbox (2).
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCheckbox(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		result, _ := enumProperty(view, ItemCheckbox, view.Session(), 0)
		return result
	}
	return 0
}

// GetListViewCheckedItems returns the array of ListView checked items.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCheckedItems(view View, subviewID string) []int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if listView, ok := view.(ListView); ok {
			checkedItems := listView.getCheckedItems()
			switch GetListViewCheckbox(view, "") {
			case NoneCheckbox:
				return []int{}

			case SingleCheckbox:
				if len(checkedItems) > 1 {
					return []int{checkedItems[0]}
				}
			}

			return checkedItems
		}
	}
	return []int{}
}

// IsListViewCheckedItem returns true if the ListView item with index is checked, false otherwise.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func IsListViewCheckedItem(view View, subviewID string, index int) bool {
	for _, n := range GetListViewCheckedItems(view, subviewID) {
		if n == index {
			return true
		}
	}
	return false
}

// GetListViewCheckboxVerticalAlign returns the vertical align of the ListView checkbox:
// TopAlign (0), BottomAlign (1), CenterAlign (2)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCheckboxVerticalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if align, ok := enumProperty(view, CheckboxVerticalAlign, view.Session(), TopAlign); ok {
			return align
		}
	}
	return TopAlign
}

// GetListViewCheckboxHorizontalAlign returns the horizontal align of the ListView checkbox:
// LeftAlign (0), RightAlign (1), CenterAlign (2)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewCheckboxHorizontalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if align, ok := enumProperty(view, CheckboxHorizontalAlign, view.Session(), LeftAlign); ok {
			return align
		}
	}
	return LeftAlign
}

// GetListItemVerticalAlign returns the vertical align of the ListView item content:
// TopAlign (0), BottomAlign (1), CenterAlign (2)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemVerticalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if align, ok := enumProperty(view, ItemVerticalAlign, view.Session(), TopAlign); ok {
			return align
		}
	}
	return TopAlign
}

// ItemHorizontalAlign returns the horizontal align of the ListView item content:
// LeftAlign (0), RightAlign (1), CenterAlign (2), StretchAlign (3)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemHorizontalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if align, ok := enumProperty(view, ItemHorizontalAlign, view.Session(), LeftAlign); ok {
			return align
		}
	}
	return LeftAlign
}

// GetListItemFrame - returns the location and size of the ListView item in pixels.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListItemFrame(view View, subviewID string, index int) Frame {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if listView, ok := view.(ListView); ok {
			itemFrames := listView.getItemFrames()
			if index >= 0 && index < len(itemFrames) {
				return itemFrames[index]
			}
		}
	}
	return Frame{Left: 0, Top: 0, Width: 0, Height: 0}
}

// GetListViewAdapter - returns the ListView adapter.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListViewAdapter(view View, subviewID string) ListAdapter {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(Items); value != nil {
			if adapter, ok := value.(ListAdapter); ok {
				return adapter
			}
		}
	}
	return nil
}

// ReloadListViewData updates ListView content
// If the second argument (subviewID) is "" then content the first argument (view) is updated.
func ReloadListViewData(view View, subviewID string) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if listView, ok := view.(ListView); ok {
			listView.ReloadListViewData()
		}
	}
}
