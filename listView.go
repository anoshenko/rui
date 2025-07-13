package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// Constants which represent [ListView] specific properties and events
const (
	// ListItemClickedEvent is the constant for "list-item-clicked" property tag.
	//
	// Used by ListView.
	// Occur when the user clicks on an item in the list.
	//
	// General listener format:
	//
	//  func(list rui.ListView, item int)
	//
	// where:
	//   - list - Interface of a list which generated this event,
	//   - item - An index of an item clicked.
	//
	// Allowed listener formats:
	//
	//  func(item int)
	//  func(list rui.ListView)
	//  func()
	ListItemClickedEvent PropertyName = "list-item-clicked"

	// ListItemSelectedEvent is the constant for "list-item-selected" property tag.
	//
	// Used by ListView.
	// Occur when a list item becomes selected.
	//
	// General listener format:
	//
	//  func(list rui.ListView, item int)
	//
	// where:
	//   - list - Interface of a list which generated this event,
	//   - item - An index of an item selected.
	//
	// Allowed listener formats:
	//
	//  func(item int)
	//  func(list rui.ListView)
	//  func()
	ListItemSelectedEvent PropertyName = "list-item-selected"

	// ListItemCheckedEvent is the constant for "list-item-checked" property tag.
	//
	// Used by ListView.
	// Occur when a list item checkbox becomes checked or unchecked.
	//
	// General listener format:
	//
	//  func(list rui.ListView, checkedItems []int).
	//
	// where:
	//   - list - Interface of a list which generated this event,
	//   - checkedItems - Array of indices of marked elements.
	//
	// Allowed listener formats:
	//
	//  func(checkedItems []int)
	//  func(list rui.ListView)
	//  func()
	ListItemCheckedEvent PropertyName = "list-item-checked"

	// ListItemStyle is the constant for "list-item-style" property tag.
	//
	// Used by ListView.
	// Defines the style of an unselected item.
	//
	// Supported types: string.
	ListItemStyle PropertyName = "list-item-style"

	// CurrentStyle is the constant for "current-style" property tag.
	//
	// Used by ListView.
	// Defines the style of the selected item when the ListView is focused.
	//
	// Supported types: string.
	CurrentStyle PropertyName = "current-style"

	// CurrentInactiveStyle is the constant for "current-inactive-style" property tag.
	//
	// Used by ListView.
	// Defines the style of the selected item when the ListView is unfocused.
	//
	// Supported types: string.
	CurrentInactiveStyle PropertyName = "current-inactive-style"
)

// Constants which represent values of the "orientation" property of the [ListView]. These are aliases for values used in
// [ListLayout] "orientation" property like TopDownOrientation and StartToEndOrientation
const (
	// VerticalOrientation is the vertical ListView orientation
	VerticalOrientation = 0

	// HorizontalOrientation is the horizontal ListView orientation
	HorizontalOrientation = 1
)

// Constants which represent values of a "checkbox" property of [ListView]
const (
	// NoneCheckbox is value of "checkbox" property: no checkbox
	NoneCheckbox = 0

	// SingleCheckbox is value of "checkbox" property: only one item can be checked
	SingleCheckbox = 1

	// MultipleCheckbox is value of "checkbox" property: several items can be checked
	MultipleCheckbox = 2
)

// ListView represents a ListView view
type ListView interface {
	View
	ParentView
	// ReloadListViewData updates ListView content
	ReloadListViewData()

	getItemFrames() []Frame
}

type listViewData struct {
	viewData
	items     []View
	itemFrame []Frame
}

// NewListView creates the new list view
func NewListView(session Session, params Params) ListView {
	view := new(listViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newListView(session Session) View {
	return new(listViewData) // NewListView(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (listView *listViewData) init(session Session) {
	listView.viewData.init(session)
	listView.tag = "ListView"
	listView.systemClass = "ruiListView"
	listView.items = []View{}
	listView.itemFrame = []Frame{}
	listView.normalize = normalizeListViewTag
	listView.get = listView.getFunc
	listView.set = listView.setFunc
	listView.remove = listView.removeFunc
	listView.changed = listView.propertyChanged
}

func (listView *listViewData) Views() []View {
	return listView.items
}

func normalizeListViewTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Content:
		tag = Items

	case HorizontalAlign:
		tag = ItemHorizontalAlign

	case VerticalAlign:
		tag = ItemVerticalAlign

	case "wrap":
		tag = ListWrap

	case "row-gap":
		return ListRowGap

	case ColumnGap:
		return ListColumnGap
	}
	return tag
}

func (listView *listViewData) removeFunc(tag PropertyName) []PropertyName {
	switch tag {
	case Gap:
		result := listView.removeFunc(ListRowGap)
		if result != nil {
			if result2 := listView.removeFunc(ListColumnGap); result2 != nil {
				result = append(result, result2...)
			}
		}
		return result
	}
	return listView.viewData.removeFunc(tag)
}

func (listView *listViewData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Gap:
		result := listView.setFunc(ListRowGap, value)
		if result != nil {
			if result2 := listView.setFunc(ListColumnGap, value); result2 != nil {
				result = append(result, result2...)
			}
		}
		return result

	case ListItemClickedEvent, ListItemSelectedEvent:
		return setOneArgEventListener[ListView, int](listView, tag, value)

	case ListItemCheckedEvent:
		return setOneArgEventListener[ListView, []int](listView, tag, value)

	case Checked:
		var checked []int
		switch value := value.(type) {
		case string:
			elements := strings.Split(value, ",")
			checked = make([]int, 0, len(elements))
			for _, val := range elements {
				if val = strings.Trim(val, " \t"); val != "" {
					n, err := strconv.Atoi(val)
					if err != nil {
						invalidPropertyValue(Checked, value)
						ErrorLog(err.Error())
						return nil
					}
					checked = append(checked, n)
				}
			}

		case int:
			checked = []int{value}

		case []int:
			checked = value

		default:
			notCompatibleType(tag, value)
			return nil
		}

		return setArrayPropertyValue(listView, Checked, checked)

	case Items:
		return listView.setItems(value)

	case ListItemStyle, CurrentStyle, CurrentInactiveStyle:
		if text, ok := value.(string); ok {
			return setStringPropertyValue(listView, tag, text)
		}
		notCompatibleType(tag, value)
		return nil

	case Current:
		return setIntProperty(listView, Current, value)
	}

	return listView.viewData.setFunc(tag, value)

}

func (listView *listViewData) propertyChanged(tag PropertyName) {
	switch tag {

	case Current:
		updateInnerHTML(listView.htmlID(), listView.Session())
		if listeners := getOneArgEventListeners[ListView, int](listView, nil, ListItemSelectedEvent); len(listeners) > 0 {
			current := GetCurrent(listView)
			for _, listener := range listeners {
				listener.Run(listView, current)
			}
		}

	case Checked:
		updateInnerHTML(listView.htmlID(), listView.Session())
		if listeners := getOneArgEventListeners[ListView, []int](listView, nil, ListItemCheckedEvent); len(listeners) > 0 {
			checked := GetListViewCheckedItems(listView)
			for _, listener := range listeners {
				listener.Run(listView, checked)
			}
		}

	case Items, Orientation, ListWrap, ListRowGap, ListColumnGap, VerticalAlign, HorizontalAlign, Style, StyleDisabled, ItemWidth, ItemHeight,
		ItemHorizontalAlign, ItemVerticalAlign, ItemCheckbox, CheckboxHorizontalAlign, CheckboxVerticalAlign, ListItemStyle, AccentColor:
		updateInnerHTML(listView.htmlID(), listView.Session())

	case CurrentStyle:
		listView.Session().updateProperty(listView.htmlID(), "data-focusitemstyle", listViewCurrentStyle(listView))
		updateInnerHTML(listView.htmlID(), listView.Session())

	case CurrentInactiveStyle:
		listView.Session().updateProperty(listView.htmlID(), "data-bluritemstyle", listViewCurrentInactiveStyle(listView))
		updateInnerHTML(listView.htmlID(), listView.Session())

	default:
		listView.viewData.propertyChanged(tag)
	}
}

func (listView *listViewData) getFunc(tag PropertyName) any {
	switch tag {
	case Gap:
		if rowGap := GetListRowGap(listView); rowGap.Equal(GetListColumnGap(listView)) {
			return rowGap
		}
		return AutoSize()

	case ListItemClickedEvent, ListItemSelectedEvent:
		if listeners := getOneArgEventRawListeners[ListView, int](listView, nil, tag); len(listeners) > 0 {
			return listeners
		}
		return nil

	case ListItemCheckedEvent:
		if listeners := getOneArgEventRawListeners[ListView, []int](listView, nil, tag); len(listeners) > 0 {
			return listeners
		}
		return nil
	}
	return listView.viewData.getFunc(tag)
}

func (listView *listViewData) setItems(value any) []PropertyName {
	var adapter ListAdapter = nil

	session := listView.session

	switch value := value.(type) {
	case []string:
		adapter = NewTextListAdapter(value, nil)

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
					if view := CreateViewFromObject(session, val.Object(), nil); view != nil {
						items[i] = view
					} else {
						return nil
					}
				} else {
					items[i] = NewTextView(session, Params{Text: val.Value()})
				}
			}
			adapter = NewViewListAdapter(items)
		} else {
			items := make([]string, len(value))
			for i, val := range value {
				items[i] = val.Value()
			}
			adapter = NewTextListAdapter(items, nil)
		}

	case []any:
		items := make([]View, len(value))
		for i, val := range value {
			switch value := val.(type) {
			case View:
				items[i] = value

			case string:
				items[i] = NewTextView(session, Params{Text: value})

			case fmt.Stringer:
				items[i] = NewTextView(session, Params{Text: value.String()})

			case float32:
				items[i] = NewTextView(session, Params{Text: fmt.Sprintf("%g", float64(value))})

			case float64:
				items[i] = NewTextView(session, Params{Text: fmt.Sprintf("%g", value)})

			default:
				if n, ok := isInt(val); ok {
					items[i] = NewTextView(session, Params{Text: strconv.Itoa(n)})
				} else {
					notCompatibleType(Items, value)
					return nil
				}
			}
		}
		adapter = NewViewListAdapter(items)

	case []View:
		adapter = NewViewListAdapter(value)

	case ListAdapter:
		adapter = value

	default:
		notCompatibleType(Items, value)
		return nil
	}

	listView.setRaw(Items, adapter)
	return []PropertyName{Items}
}

func (listView *listViewData) Focusable() bool {
	return true
}

func (listView *listViewData) getAdapter() ListAdapter {
	if value := listView.getRaw(Items); value != nil {
		if adapter, ok := value.(ListAdapter); ok {
			return adapter
		}
	}
	if obj := listView.binding(); obj != nil {
		if adapter, ok := obj.(ListAdapter); ok {
			return adapter
		}
	}
	return nil
}

func (listView *listViewData) ReloadListViewData() {
	itemCount := 0
	if adapter := listView.getAdapter(); adapter != nil {
		itemCount = adapter.ListSize()

		if itemCount != len(listView.items) {
			listView.items = make([]View, itemCount)
			listView.itemFrame = make([]Frame, itemCount)
		}

		for i := range itemCount {
			listView.items[i] = adapter.ListItem(i, listView.Session())
		}
	} else if len(listView.items) > 0 {
		listView.items = []View{}
		listView.itemFrame = []Frame{}
	}

	updateInnerHTML(listView.htmlID(), listView.session)
}

func (listView *listViewData) getItemFrames() []Frame {
	return listView.itemFrame
}

func (listView *listViewData) itemAlign(buffer *strings.Builder) {
	values := enumProperties[ItemHorizontalAlign].cssValues
	if hAlign := GetListItemHorizontalAlign(listView); hAlign >= 0 && hAlign < len(values) {
		buffer.WriteString(" justify-items: ")
		buffer.WriteString(values[hAlign])
		buffer.WriteRune(';')
	}

	values = enumProperties[ItemVerticalAlign].cssValues
	if vAlign := GetListItemVerticalAlign(listView); vAlign >= 0 && vAlign < len(values) {
		buffer.WriteString(" align-items: ")
		buffer.WriteString(values[vAlign])
		buffer.WriteRune(';')
	}
}

func (listView *listViewData) itemSize(buffer *strings.Builder) {
	if itemWidth := GetListItemWidth(listView); itemWidth.Type != Auto {
		buffer.WriteString(` min-width: `)
		buffer.WriteString(itemWidth.cssString("", listView.Session()))
		buffer.WriteRune(';')
	}

	if itemHeight := GetListItemHeight(listView); itemHeight.Type != Auto {
		buffer.WriteString(` min-height: `)
		buffer.WriteString(itemHeight.cssString("", listView.Session()))
		buffer.WriteRune(';')
	}
}

func (listView *listViewData) getDivs(checkbox, hCheckboxAlign, vCheckboxAlign int) (string, string, string) {
	session := listView.Session()

	contentBuilder := allocStringBuilder()
	defer freeStringBuilder(contentBuilder)

	contentBuilder.WriteString(`<div style="display: grid;`)
	listView.itemAlign(contentBuilder)

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

	accentColor := Color(0)
	if color := GetAccentColor(listView, ""); color != 0 {
		accentColor = color
	}

	if checkbox == SingleCheckbox {
		offDivBuilder.WriteString(session.radiobuttonOffImage())
		onDivBuilder.WriteString(session.radiobuttonOnImage(accentColor))
	} else {
		offDivBuilder.WriteString(session.checkboxOffImage(accentColor))
		onDivBuilder.WriteString(session.checkboxOnImage(accentColor))
	}

	onDivBuilder.WriteString("</div>")
	offDivBuilder.WriteString("</div>")

	return onDivBuilder.String(), offDivBuilder.String(), contentBuilder.String()
}

func (listView *listViewData) checkboxItemDiv(hCheckboxAlign, vCheckboxAlign int) string {
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
		itemStyleBuilder.WriteString(gap.cssString("auto", listView.Session()))
		itemStyleBuilder.WriteRune(';')
	}

	itemStyleBuilder.WriteString(`">`)
	return itemStyleBuilder.String()

}

func (listView *listViewData) getItemView(adapter ListAdapter, index int) View {
	if adapter == nil || index < 0 || index >= adapter.ListSize() {
		return nil
	}

	if size := adapter.ListSize(); size != len(listView.items) {
		listView.items = make([]View, size)
		listView.itemFrame = make([]Frame, size)
	}

	if listView.items[index] == nil {
		listView.items[index] = adapter.ListItem(index, listView.Session())
	}

	return listView.items[index]
}

func listViewItemStyle(view View, tag PropertyName, defaultStyle string) string {
	session := view.Session()
	if value := view.getRaw(tag); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = session.resolveConstants(style); ok {
				return style
			}
		}
	}
	if value := valueFromStyle(view, tag); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = session.resolveConstants(style); ok {
				return style
			}
		}
	}
	return defaultStyle
}

func (listView *listViewData) listItemStyle() string {
	return listViewItemStyle(listView, ListItemStyle, "ruiListItem")
}

func listViewCurrentStyle(view View) string {
	return listViewItemStyle(view, CurrentStyle, "ruiListItemFocused")
}

func listViewCurrentInactiveStyle(view View) string {
	return listViewItemStyle(view, CurrentInactiveStyle, "ruiListItemSelected")
}

func (listView *listViewData) checkboxSubviews(adapter ListAdapter, buffer *strings.Builder, checkbox int) {
	count := adapter.ListSize()
	listViewID := listView.htmlID()

	hCheckboxAlign := GetListViewCheckboxHorizontalAlign(listView)
	vCheckboxAlign := GetListViewCheckboxVerticalAlign(listView)

	itemDiv := listView.checkboxItemDiv(hCheckboxAlign, vCheckboxAlign)
	onDiv, offDiv, contentDiv := listView.getDivs(checkbox, hCheckboxAlign, vCheckboxAlign)

	current := GetCurrent(listView)
	checkedItems := GetListViewCheckedItems(listView)
	for i := range count {
		buffer.WriteString(`<div id="`)
		buffer.WriteString(listViewID)
		buffer.WriteRune('-')
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(`" class="ruiView `)
		buffer.WriteString(listView.listItemStyle())
		if i == current {
			buffer.WriteRune(' ')
			buffer.WriteString(listViewCurrentInactiveStyle(listView))
		}
		buffer.WriteString(`" onclick="listItemClickEvent(this, event)" data-left="0" data-top="0" data-width="0" data-height="0" style="display: grid; justify-items: stretch; align-items: stretch;`)
		listView.itemSize(buffer)
		if ext, ok := adapter.(ListItemEnabled); ok {
			if !ext.IsListItemEnabled(i) {
				buffer.WriteString(`" data-disabled="1`)
			}
		}
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

		if view := listView.getItemView(adapter, i); view != nil {
			//view.setNoResizeEvent()
			viewHTML(view, buffer, "")
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}

		buffer.WriteString(`</div></div></div>`)
	}
}

func (listView *listViewData) noneCheckboxSubviews(adapter ListAdapter, buffer *strings.Builder) {
	count := adapter.ListSize()
	listViewID := listView.htmlID()

	itemStyleBuilder := allocStringBuilder()
	defer freeStringBuilder(itemStyleBuilder)

	itemStyleBuilder.WriteString(`data-left="0" data-top="0" data-width="0" data-height="0" style="max-width: 100%; max-height: 100%; display: grid;`)

	listView.itemAlign(itemStyleBuilder)
	listView.itemSize(itemStyleBuilder)

	itemStyleBuilder.WriteString(`" onclick="listItemClickEvent(this, event)"`)
	itemStyle := itemStyleBuilder.String()

	current := GetCurrent(listView)
	for i := range count {
		buffer.WriteString(`<div id="`)
		buffer.WriteString(listViewID)
		buffer.WriteRune('-')
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(`" class="ruiView `)
		buffer.WriteString(listView.listItemStyle())
		if i == current {
			buffer.WriteRune(' ')
			buffer.WriteString(listViewCurrentInactiveStyle(listView))
		}
		buffer.WriteString(`" `)
		buffer.WriteString(itemStyle)
		if ext, ok := adapter.(ListItemEnabled); ok {
			if !ext.IsListItemEnabled(i) {
				buffer.WriteString(` data-disabled="1"`)
			}
		}
		buffer.WriteString(`>`)

		if view := listView.getItemView(adapter, i); view != nil {
			//view.setNoResizeEvent()
			viewHTML(view, buffer, "")
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}

		buffer.WriteString(`</div>`)
	}
}

func (listView *listViewData) updateCheckboxItem(index int, checked bool) {

	checkbox := GetListViewCheckbox(listView)
	hCheckboxAlign := GetListViewCheckboxHorizontalAlign(listView)
	vCheckboxAlign := GetListViewCheckboxVerticalAlign(listView)
	onDiv, offDiv, contentDiv := listView.getDivs(checkbox, hCheckboxAlign, vCheckboxAlign)

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(listView.checkboxItemDiv(hCheckboxAlign, vCheckboxAlign))
	if checked {
		buffer.WriteString(onDiv)
	} else {
		buffer.WriteString(offDiv)
	}
	buffer.WriteString(contentDiv)

	session := listView.Session()
	if adapter := listView.getAdapter(); adapter != nil {
		if view := listView.getItemView(adapter, index); view != nil {
			view.setNoResizeEvent()
			viewHTML(view, buffer, "")
		} else {
			buffer.WriteString("ERROR: invalid item view")
		}
	}
	buffer.WriteString(`</div></div>`)
	session.updateInnerHTML(listView.htmlID()+"-"+strconv.Itoa(index), buffer.String())
}

func (listView *listViewData) htmlProperties(self View, buffer *strings.Builder) {
	listView.viewData.htmlProperties(self, buffer)
	buffer.WriteString(`onfocus="listViewFocusEvent(this, event)" onblur="listViewBlurEvent(this, event)"`)
	buffer.WriteString(` onkeydown="listViewKeyDownEvent(this, event)" data-focusitemstyle="`)
	buffer.WriteString(listViewCurrentStyle(listView))
	buffer.WriteString(`" data-bluritemstyle="`)
	buffer.WriteString(listViewCurrentInactiveStyle(listView))
	buffer.WriteString(`"`)

	if adapter := listView.getAdapter(); adapter != nil {
		if current := GetCurrent(listView); current >= 0 && current < adapter.ListSize() {
			buffer.WriteString(` data-current="`)
			buffer.WriteString(listView.htmlID())
			buffer.WriteRune('-')
			buffer.WriteString(strconv.Itoa(current))
			buffer.WriteRune('"')
		}
	}
	listView.viewData.htmlProperties(self, buffer)
}

/*
func (listView *listViewData) cssStyle(self View, builder cssBuilder) {
	listView.viewData.cssStyle(self, builder)

	if GetListWrap(listView) != WrapOff {
		switch GetListOrientation(listView) {
		case TopDownOrientation, BottomUpOrientation:
			builder.add(`max-height`, `100%`)
		default:
			builder.add(`max-width`, `100%`)
		}
	}
}
*/

func (listView *listViewData) htmlSubviews(self View, buffer *strings.Builder) {
	adapter := listView.getAdapter()
	if adapter == nil {
		return
	}

	if listSize := adapter.ListSize(); listSize == 0 {
		return
	}

	if !listView.session.ignoreViewUpdates() {
		listView.session.setIgnoreViewUpdates(true)
		defer listView.session.setIgnoreViewUpdates(false)
	}

	buffer.WriteString(`<div style="display: flex; align-content: stretch;`)

	if gap := GetListRowGap(listView); gap.Type != Auto {
		buffer.WriteString(` row-gap: `)
		buffer.WriteString(gap.cssString("0", listView.Session()))
		buffer.WriteRune(';')
	}

	if gap := GetListColumnGap(listView); gap.Type != Auto {
		buffer.WriteString(` column-gap: `)
		buffer.WriteString(gap.cssString("0", listView.Session()))
		buffer.WriteRune(';')
	}

	wrap := GetListWrap(listView)
	orientation := GetListOrientation(listView)
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
	switch GetListHorizontalAlign(listView) {
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
	switch GetListVerticalAlign(listView) {
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

	checkbox := GetListViewCheckbox(listView)
	if checkbox == NoneCheckbox {
		listView.noneCheckboxSubviews(adapter, buffer)
	} else {
		listView.checkboxSubviews(adapter, buffer, checkbox)
	}

	buffer.WriteString(`</div>`)
}

func (listView *listViewData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "itemSelected":
		if number, ok := dataIntProperty(data, `number`); ok {
			listView.handleCurrent(number)
		}

	case "itemUnselected":
		if _, ok := listView.properties[Current]; ok {
			listView.handleCurrent(-1)
		}

	case "itemClick":
		if number, ok := dataIntProperty(data, `number`); ok {
			listView.onItemClick(number)
		}

	default:
		return listView.viewData.handleCommand(self, command, data)
	}

	return true
}

func (listView *listViewData) handleCurrent(number int) {
	listView.properties[Current] = number
	for _, listener := range getOneArgEventListeners[ListView, int](listView, nil, ListItemSelectedEvent) {
		listener.Run(listView, number)
	}
	listView.runChangeListener(Current)
}

func (listView *listViewData) onItemClick(number int) {

	if IsDisabled(listView) {
		return
	}

	if current := GetCurrent(listView); current != number {
		listView.Set(Current, number)
	}

	if checkbox := GetListViewCheckbox(listView); checkbox != NoneCheckbox {
		checkedItem := GetListViewCheckedItems(listView)

		switch checkbox {
		case SingleCheckbox:
			if len(checkedItem) == 0 {
				checkedItem = []int{number}
				listView.updateCheckboxItem(number, true)
			} else if checkedItem[0] != number {
				listView.updateCheckboxItem(checkedItem[0], false)
				checkedItem = []int{number}
				listView.updateCheckboxItem(number, true)
			} else {
				checkedItem = []int{}
				listView.updateCheckboxItem(number, false)
			}

		case MultipleCheckbox:
			uncheck := false
			for i, index := range checkedItem {
				if index == number {
					uncheck = true
					listView.updateCheckboxItem(index, false)
					count := len(checkedItem)
					if count == 1 {
						checkedItem = []int{}
					} else if i == 0 {
						checkedItem = checkedItem[1:]
					} else if i == count-1 {
						checkedItem = checkedItem[:i]
					} else {
						checkedItem = append(checkedItem[:i], checkedItem[i+1:]...)
					}
					break
				}
			}

			if !uncheck {
				listView.updateCheckboxItem(number, true)
				checkedItem = append(checkedItem, number)
			}
		}

		setArrayPropertyValue(listView, Checked, checkedItem)
		listView.runChangeListener(Checked)

		for _, listener := range getOneArgEventListeners[ListView, []int](listView, nil, ListItemCheckedEvent) {
			listener.Run(listView, checkedItem)
		}
	}

	for _, listener := range getOneArgEventListeners[ListView, int](listView, nil, ListItemClickedEvent) {
		listener.Run(listView, number)
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
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, VerticalAlign, TopAlign, false)
}

// GetHorizontalAlign return the vertical align of a list/checkbox: LeftAlign (0), RightAlign (1), CenterAlign (2), StretchAlign (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, HorizontalAlign, LeftAlign, false)
}

// GetListItemClickedListeners returns a ListItemClickedListener of the ListView.
// If there are no listeners then the empty list is returned
//
// Result elements can be of the following types:
//   - func(rui.ListView, int),
//   - func(rui.ListView),
//   - func(int),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListItemClickedListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[ListView, int](view, subviewID, ListItemClickedEvent)
}

// GetListItemSelectedListeners returns a ListItemSelectedListener of the ListView.
// If there are no listeners then the empty list is returned
//
// Result elements can be of the following types:
//   - func(rui.ListView, int),
//   - func(rui.ListView),
//   - func(int),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListItemSelectedListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[ListView, int](view, subviewID, ListItemSelectedEvent)
}

// GetListItemCheckedListeners returns a ListItemCheckedListener of the ListView.
// If there are no listeners then the empty list is returned
//
// Result elements can be of the following types:
//   - func(rui.ListView, []int),
//   - func(rui.ListView),
//   - func([]int),
//   - func(),
//   - string.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListItemCheckedListeners(view View, subviewID ...string) []any {
	return getOneArgEventRawListeners[ListView, []int](view, subviewID, ListItemCheckedEvent)
}

// GetListItemWidth returns the width of a ListView item.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListItemWidth(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ItemWidth, false)
}

// GetListItemHeight returns the height of a ListView item.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListItemHeight(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ItemHeight, false)
}

// GetListViewCheckbox returns the ListView checkbox type: NoneCheckbox (0), SingleCheckbox (1), or MultipleCheckbox (2).
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListViewCheckbox(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, ItemCheckbox, 0, false)
}

// GetListViewCheckedItems returns the array of ListView checked items.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListViewCheckedItems(view View, subviewID ...string) []int {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.getRaw(Checked); value != nil {
			if checkedItems, ok := value.([]int); ok {
				switch GetListViewCheckbox(view) {
				case MultipleCheckbox:
					return checkedItems

				case SingleCheckbox:
					if len(checkedItems) > 0 {
						return []int{checkedItems[0]}
					}
				}
			}
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
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListViewCheckboxVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CheckboxVerticalAlign, TopAlign, false)
}

// GetListViewCheckboxHorizontalAlign returns the horizontal align of the ListView checkbox:
// LeftAlign (0), RightAlign (1), CenterAlign (2)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListViewCheckboxHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CheckboxHorizontalAlign, LeftAlign, false)
}

// GetListItemVerticalAlign returns the vertical align of the ListView item content:
// TopAlign (0), BottomAlign (1), CenterAlign (2)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListItemVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, ItemVerticalAlign, TopAlign, false)
}

// ItemHorizontalAlign returns the horizontal align of the ListView item content:
// LeftAlign (0), RightAlign (1), CenterAlign (2), StretchAlign (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListItemHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, ItemHorizontalAlign, LeftAlign, false)
}

// GetListItemFrame - returns the location and size of the ListView item in pixels.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
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
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListViewAdapter(view View, subviewID ...string) ListAdapter {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.Get(Items); value != nil {
			if adapter, ok := value.(ListAdapter); ok {
				return adapter
			}
		}
	}
	return nil
}

// ReloadListViewData updates ListView content
// If the second argument (subviewID) is not specified or it is "" then content the first argument (view) is updated.
func ReloadListViewData(view View, subviewID ...string) {
	if view = getSubview(view, subviewID); view != nil {
		if listView, ok := view.(ListView); ok {
			listView.ReloadListViewData()
		}
	}
}
