package rui

import (
	"strconv"
	"strings"
)

// Constants for [TabsLayout] specific view and events
const (
	// CurrentTabChangedEvent is the constant for "current-tab-changed" property tag.
	//
	// Used by `TabsLayout`.
	// Occur when the new tab becomes active.
	//
	// General listener format:
	// `func(tabsLayout rui.TabsLayout, newTab, oldTab int)`.
	//
	// where:
	// tabsLayout - Interface of a tabs layout which generated this event,
	// newTab - Index of a new active tab,
	// oldTab - Index of an old active tab.
	//
	// Allowed listener formats:
	// `func(tabsLayout rui.TabsLayout, newTab int)`,
	// `func(newTab, oldTab int)`,
	// `func(newTab int)`,
	// `func()`.
	CurrentTabChangedEvent PropertyName = "current-tab-changed"

	// Icon is the constant for "icon" property tag.
	//
	// Used by `TabsLayout`.
	// Defines the icon name that is displayed in the tab. The property is set for the child view of `TabsLayout`.
	//
	// Supported types: `string`.
	Icon = "icon"

	// TabCloseButton is the constant for "tab-close-button" property tag.
	//
	// Used by `TabsLayout`.
	// Controls whether to add close button to a tab(s). This property can be set separately for each child view or for tabs
	// layout itself. Property set for child view takes precedence over the value set for tabs layout. Default value is
	// `false`.
	//
	// Supported types: `bool`, `int`, `string`.
	//
	// Values:
	// `true` or `1` or "true", "yes", "on", "1" - Tab(s) has close button.
	// `false` or `0` or "false", "no", "off", "0" - No close button in tab(s).
	TabCloseButton PropertyName = "tab-close-button"

	// TabCloseEvent is the constant for "tab-close-event" property tag.
	//
	// Used by `TabsLayout`.
	// Occurs when the user clicks on the tab close button.
	//
	// General listener format:
	// `func(tabsLayout rui.TabsLayout, tab int)`.
	//
	// where:
	// tabsLayout - Interface of a tabs layout which generated this event,
	// tab - Index of the tab.
	//
	// Allowed listener formats:
	// `func(tab int)`,
	// `func(tabsLayout rui.TabsLayout)`,
	// `func()`.
	TabCloseEvent PropertyName = "tab-close-event"

	// Tabs is the constant for "tabs" property tag.
	//
	// Used by `TabsLayout`.
	// Sets where the tabs are located. Default value is "top".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`TopTabs`) or "top" - Tabs on the top.
	// `1`(`BottomTabs`) or "bottom" - Tabs on the bottom.
	// `2`(`LeftTabs`) or "left" - Tabs on the left. Each tab is rotated 90° counterclockwise.
	// `3`(`RightTabs`) or "right" - Tabs located on the right. Each tab is rotated 90° clockwise.
	// `4`(`LeftListTabs`) or "left-list" - Tabs on the left. The tabs are displayed as a list.
	// `5`(`RightListTabs`) or "right-list" - Tabs on the right. The tabs are displayed as a list.
	// `6`(`HiddenTabs`) or "hidden" - Tabs are hidden.
	Tabs PropertyName = "tabs"

	// TabBarStyle is the constant for "tab-bar-style" property tag.
	//
	// Used by `TabsLayout`.
	// Set the style for the display of the tab bar. The default value is "ruiTabBar".
	//
	// Supported types: `string`.
	TabBarStyle PropertyName = "tab-bar-style"

	// TabStyle is the constant for "tab-style" property tag.
	//
	// Used by `TabsLayout`.
	// Set the style for the display of the tab. The default value is "ruiTab" or "ruiVerticalTab".
	//
	// Supported types: `string`.
	TabStyle PropertyName = "tab-style"

	// CurrentTabStyle is the constant for "current-tab-style" property tag.
	//
	// Used by `TabsLayout`.
	// Set the style for the display of the current(selected) tab. The default value is "ruiCurrentTab" or
	// "ruiCurrentVerticalTab".
	//
	// Supported types: `string`.
	CurrentTabStyle PropertyName = "current-tab-style"

	inactiveTabStyle = "data-inactiveTabStyle"
	activeTabStyle   = "data-activeTabStyle"
)

// Constants that are the values of the "tabs" property of a [TabsLayout]
const (
	// TopTabs - tabs of TabsLayout are on the top
	TopTabs = 0
	// BottomTabs - tabs of TabsLayout are on the bottom
	BottomTabs = 1
	// LeftTabs - tabs of TabsLayout are on the left. Bookmarks are rotated counterclockwise 90 degrees.
	LeftTabs = 2
	// RightTabs - tabs of TabsLayout are on the right. Bookmarks are rotated clockwise 90 degrees.
	RightTabs = 3
	// LeftListTabs - tabs of TabsLayout are on the left
	LeftListTabs = 4
	// RightListTabs - tabs of TabsLayout are on the right
	RightListTabs = 5
	// HiddenTabs - tabs of TabsLayout are hidden
	HiddenTabs = 6
)

// TabsLayout represents a TabsLayout view
type TabsLayout interface {
	ViewsContainer
	ListAdapter
}

type tabsLayoutData struct {
	viewsContainerData
}

// NewTabsLayout create new TabsLayout object and return it
func NewTabsLayout(session Session, params Params) TabsLayout {
	view := new(tabsLayoutData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newTabsLayout(session Session) View {
	//return NewTabsLayout(session, nil)
	return new(tabsLayoutData)
}

// Init initialize fields of ViewsContainer by default values
func (tabsLayout *tabsLayoutData) init(session Session) {
	tabsLayout.viewsContainerData.init(session)
	tabsLayout.tag = "TabsLayout"
	tabsLayout.systemClass = "ruiTabsLayout"
	tabsLayout.set = tabsLayout.setFunc
	tabsLayout.changed = tabsLayout.propertyChanged
}

func tabsLayoutCurrent(view View, defaultValue int) int {
	result, _ := intProperty(view, Current, view.Session(), defaultValue)
	return result
}

func (tabsLayout *tabsLayoutData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case CurrentTabChangedEvent:
		return setTwoArgEventListener[TabsLayout, int](tabsLayout, tag, value)

	case TabCloseEvent:
		return setOneArgEventListener[TabsLayout, int](tabsLayout, tag, value)

	case Current:
		tabsLayout.setRaw("old-current", tabsLayoutCurrent(tabsLayout, -1))

		if current, ok := value.(int); ok && current < 0 {
			tabsLayout.setRaw(Current, nil)
			return []PropertyName{tag}
		}

		return setIntProperty(tabsLayout, Current, value)

	case TabStyle, CurrentTabStyle, TabBarStyle:
		if text, ok := value.(string); ok {
			return setStringPropertyValue(tabsLayout, tag, text)
		}
		notCompatibleType(tag, value)
		return nil
	}

	return tabsLayout.viewsContainerData.setFunc(tag, value)
}

func (tabsLayout *tabsLayoutData) propertyChanged(tag PropertyName) {
	switch tag {
	case Current:
		session := tabsLayout.Session()
		current := GetCurrent(tabsLayout)
		session.callFunc("activateTab", tabsLayout.htmlID(), current)

		if listeners := getTwoArgEventListeners[TabsLayout, int](tabsLayout, nil, CurrentTabChangedEvent); len(listeners) > 0 {
			oldCurrent, _ := intProperty(tabsLayout, "old-current", session, -1)
			for _, listener := range listeners {
				listener(tabsLayout, current, oldCurrent)
			}
		}

	case Tabs:
		htmlID := tabsLayout.htmlID()
		session := tabsLayout.Session()
		session.updateProperty(htmlID, inactiveTabStyle, tabsLayoutInactiveTabStyle(tabsLayout))
		session.updateProperty(htmlID, activeTabStyle, tabsLayoutActiveTabStyle(tabsLayout))
		updateCSSStyle(htmlID, session)
		updateInnerHTML(htmlID, session)

	case TabStyle, CurrentTabStyle, TabBarStyle:
		htmlID := tabsLayout.htmlID()
		session := tabsLayout.Session()
		session.updateProperty(htmlID, inactiveTabStyle, tabsLayoutInactiveTabStyle(tabsLayout))
		session.updateProperty(htmlID, activeTabStyle, tabsLayoutActiveTabStyle(tabsLayout))
		updateInnerHTML(htmlID, session)

	case TabCloseButton:
		updateInnerHTML(tabsLayout.htmlID(), tabsLayout.Session())

	default:
		tabsLayout.viewsContainerData.propertyChanged(tag)
	}
}

/*
	func (tabsLayout *tabsLayoutData) valueToTabListeners(value any) []func(TabsLayout, int, int) {
		if value == nil {
			return []func(TabsLayout, int, int){}
		}

		switch value := value.(type) {
		case func(TabsLayout, int, int):
			return []func(TabsLayout, int, int){value}

		case func(TabsLayout, int):
			fn := func(view TabsLayout, current, _ int) {
				value(view, current)
			}
			return []func(TabsLayout, int, int){fn}

		case func(TabsLayout):
			fn := func(view TabsLayout, _, _ int) {
				value(view)
			}
			return []func(TabsLayout, int, int){fn}

		case func(int, int):
			fn := func(_ TabsLayout, current, old int) {
				value(current, old)
			}
			return []func(TabsLayout, int, int){fn}

		case func(int):
			fn := func(_ TabsLayout, current, _ int) {
				value(current)
			}
			return []func(TabsLayout, int, int){fn}

		case func():
			fn := func(TabsLayout, int, int) {
				value()
			}
			return []func(TabsLayout, int, int){fn}

		case []func(TabsLayout, int, int):
			return value

		case []func(TabsLayout, int):
			listeners := make([]func(TabsLayout, int, int), len(value))
			for i, val := range value {
				if val == nil {
					return nil
				}
				listeners[i] = func(view TabsLayout, current, _ int) {
					val(view, current)
				}
			}
			return listeners

		case []func(TabsLayout):
			listeners := make([]func(TabsLayout, int, int), len(value))
			for i, val := range value {
				if val == nil {
					return nil
				}
				listeners[i] = func(view TabsLayout, _, _ int) {
					val(view)
				}
			}
			return listeners

		case []func(int, int):
			listeners := make([]func(TabsLayout, int, int), len(value))
			for i, val := range value {
				if val == nil {
					return nil
				}
				listeners[i] = func(_ TabsLayout, current, old int) {
					val(current, old)
				}
			}
			return listeners

		case []func(int):
			listeners := make([]func(TabsLayout, int, int), len(value))
			for i, val := range value {
				if val == nil {
					return nil
				}
				listeners[i] = func(_ TabsLayout, current, _ int) {
					val(current)
				}
			}
			return listeners

		case []func():
			listeners := make([]func(TabsLayout, int, int), len(value))
			for i, val := range value {
				if val == nil {
					return nil
				}
				listeners[i] = func(TabsLayout, int, int) {
					val()
				}
			}
			return listeners

		case []any:
			listeners := make([]func(TabsLayout, int, int), len(value))
			for i, val := range value {
				if val == nil {
					return nil
				}
				switch val := val.(type) {
				case func(TabsLayout, int, int):
					listeners[i] = val

				case func(TabsLayout, int):
					listeners[i] = func(view TabsLayout, current, _ int) {
						val(view, current)
					}

				case func(TabsLayout):
					listeners[i] = func(view TabsLayout, _, _ int) {
						val(view)
					}

				case func(int, int):
					listeners[i] = func(_ TabsLayout, current, old int) {
						val(current, old)
					}

				case func(int):
					listeners[i] = func(_ TabsLayout, current, _ int) {
						val(current)
					}

				case func():
					listeners[i] = func(TabsLayout, int, int) {
						val()
					}

				default:
					return nil
				}
			}
			return listeners
		}

		return nil
	}
*/

func (tabsLayout *tabsLayoutData) tabsLocation() int {
	tabs, _ := enumProperty(tabsLayout, Tabs, tabsLayout.session, 0)
	return tabs
}

func (tabsLayout *tabsLayoutData) tabBarStyle() string {
	if style, ok := stringProperty(tabsLayout, TabBarStyle, tabsLayout.session); ok {
		return style
	}
	if value := valueFromStyle(tabsLayout, TabBarStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = tabsLayout.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	return "ruiTabBar"
}

func tabsLayoutInactiveTabStyle(view View) string {
	session := view.Session()
	if style, ok := stringProperty(view, TabStyle, session); ok {
		return style
	}
	if value := valueFromStyle(view, TabStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = session.resolveConstants(style); ok {
				return style
			}
		}
	}

	tabs, _ := enumProperty(view, Tabs, session, 0)
	switch tabs {
	case LeftTabs, RightTabs:
		return "ruiVerticalTab"
	}
	return "ruiTab"
}

func tabsLayoutActiveTabStyle(view View) string {
	session := view.Session()
	if style, ok := stringProperty(view, CurrentTabStyle, session); ok {
		return style
	}
	if value := valueFromStyle(view, CurrentTabStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = session.resolveConstants(style); ok {
				return style
			}
		}
	}

	tabs, _ := enumProperty(view, Tabs, session, 0)
	switch tabs {
	case LeftTabs, RightTabs:
		return "ruiCurrentVerticalTab"
	}
	return "ruiCurrentTab"
}

func (tabsLayout *tabsLayoutData) ListSize() int {
	if tabsLayout.views == nil {
		tabsLayout.views = []View{}
	}
	return len(tabsLayout.views)
}

func (tabsLayout *tabsLayoutData) ListItem(index int, session Session) View {
	if tabsLayout.views == nil {
		tabsLayout.views = []View{}
	}

	if index < 0 || index >= len(tabsLayout.views) {
		return NewTextView(session, Params{
			Text: "Invalid index",
		})
	}

	views := []View{}
	page := tabsLayout.views[index]

	if icon, ok := imageProperty(page, Icon, session); ok && icon != "" {
		views = append(views, NewImageView(session, Params{
			Source: icon,
			Row:    0,
			Column: 0,
		}))
	}

	title, ok := stringProperty(page, Title, session)
	if !ok || title == "" {
		title = "No title"
	}
	if !GetNotTranslate(tabsLayout) {
		title, _ = tabsLayout.Session().GetString(title)
	}

	views = append(views, NewTextView(session, Params{
		Text:   title,
		Row:    0,
		Column: 1,
	}))

	closeButton, _ := boolProperty(tabsLayout, TabCloseButton, tabsLayout.session)
	if close, ok := boolProperty(page, TabCloseButton, tabsLayout.session); ok {
		closeButton = close
	}

	if closeButton {
		views = append(views, NewGridLayout(session, Params{
			Style:   "ruiTabCloseButton",
			Row:     0,
			Column:  2,
			Content: "✕",
			ClickEvent: func() {
				for _, listener := range getOneArgEventListeners[TabsLayout, int](tabsLayout, nil, TabCloseEvent) {
					listener(tabsLayout, index)
				}
			},
		}))
	}

	tabsHeight, _ := sizeConstant(session, "ruiTabHeight")

	return NewGridLayout(session, Params{
		Height:            tabsHeight,
		CellWidth:         []SizeUnit{AutoSize(), Fr(1), AutoSize()},
		CellVerticalAlign: CenterAlign,
		GridRowGap:        Px(8),
		Content:           views,
	})
}

func (tabsLayout *tabsLayoutData) IsListItemEnabled(index int) bool {
	return true
}

func (tabsLayout *tabsLayoutData) updateTitle(view View, tag PropertyName) {
	session := tabsLayout.session
	title, _ := stringProperty(view, Title, session)
	if !GetNotTranslate(tabsLayout) {
		title, _ = session.GetString(title)
	}
	session.updateInnerHTML(view.htmlID()+"-title", title)
}

func (tabsLayout *tabsLayoutData) updateIcon(view View, tag PropertyName) {
	session := tabsLayout.session
	icon, _ := stringProperty(view, Icon, session)
	session.updateProperty(view.htmlID()+"-icon", "src", icon)
}

func (tabsLayout *tabsLayoutData) updateTabCloseButton(view View, tag PropertyName) {
	updateInnerHTML(tabsLayout.htmlID(), tabsLayout.session)
}

// Append appends view to the end of view list
func (tabsLayout *tabsLayoutData) Append(view View) {
	if tabsLayout.views == nil {
		tabsLayout.views = []View{}
	}
	if view != nil {
		tabsLayout.viewsContainerData.Append(view)
		view.SetChangeListener(Title, tabsLayout.updateTitle)
		view.SetChangeListener(Icon, tabsLayout.updateIcon)
		view.SetChangeListener(TabCloseButton, tabsLayout.updateTabCloseButton)
		if len(tabsLayout.views) == 1 {
			tabsLayout.setRaw(Current, nil)
			tabsLayout.Set(Current, 0)
		}
	}
}

// Insert inserts view to the "index" position in view list
func (tabsLayout *tabsLayoutData) Insert(view View, index int) {
	if tabsLayout.views == nil {
		tabsLayout.views = []View{}
	}
	if view != nil {
		if current := GetCurrent(tabsLayout); current >= index {
			tabsLayout.setRaw(Current, current+1)
			defer tabsLayout.currentChanged(current+1, current)
		}
		tabsLayout.viewsContainerData.Insert(view, index)
		view.SetChangeListener(Title, tabsLayout.updateTitle)
		view.SetChangeListener(Icon, tabsLayout.updateIcon)
		view.SetChangeListener(TabCloseButton, tabsLayout.updateTabCloseButton)
	}
}

func (tabsLayout *tabsLayoutData) currentChanged(newCurrent, oldCurrent int) {
	for _, listener := range getTwoArgEventListeners[TabsLayout, int](tabsLayout, nil, CurrentTabChangedEvent) {
		listener(tabsLayout, newCurrent, oldCurrent)
	}
	if listener, ok := tabsLayout.changeListener[Current]; ok {
		listener(tabsLayout, Current)
	}
}

// Remove removes view from list and return it
func (tabsLayout *tabsLayoutData) RemoveView(index int) View {

	if index < 0 || index >= len(tabsLayout.views) {
		return nil
	}

	oldCurrent := GetCurrent(tabsLayout)
	newCurrent := oldCurrent
	if index < oldCurrent || (index == oldCurrent && oldCurrent > 0) {
		newCurrent--
	}

	if view := tabsLayout.viewsContainerData.RemoveView(index); view != nil {
		view.SetChangeListener(Title, nil)
		view.SetChangeListener(Icon, nil)
		view.SetChangeListener(TabCloseButton, nil)

		if newCurrent != oldCurrent {
			tabsLayout.setRaw(Current, newCurrent)
			tabsLayout.currentChanged(newCurrent, oldCurrent)
		}
	}
	return nil

	/*
		if tabsLayout.views == nil {
			tabsLayout.views = []View{}
			return nil
		}

		count := len(tabsLayout.views)
		if index < 0 || index >= count {
			return nil
		}

		view := tabsLayout.views[index]
		view.setParentID("")
		view.SetChangeListener(Title, nil)
		view.SetChangeListener(Icon, nil)
		view.SetChangeListener(TabCloseButton, nil)

		current := GetCurrent(tabsLayout)
		if index < current || (index == current && current > 0) {
			current--
		}

		if len(tabsLayout.views) == 1 {
			tabsLayout.views = []View{}
			current = -1
		} else if index == 0 {
			tabsLayout.views = tabsLayout.views[1:]
		} else if index == count-1 {
			tabsLayout.views = tabsLayout.views[:index]
		} else {
			tabsLayout.views = append(tabsLayout.views[:index], tabsLayout.views[index+1:]...)
		}

		updateInnerHTML(tabsLayout.parentHTMLID(), tabsLayout.session)
		tabsLayout.propertyChangedEvent(Content)

		delete(tabsLayout.view, Current)
		tabsLayout.Set(Current, current)
		return view
	*/
}

func (tabsLayout *tabsLayoutData) htmlProperties(self View, buffer *strings.Builder) {
	tabsLayout.viewsContainerData.htmlProperties(self, buffer)
	buffer.WriteString(` data-inactiveTabStyle="`)
	buffer.WriteString(tabsLayoutInactiveTabStyle(tabsLayout))
	buffer.WriteString(`" data-activeTabStyle="`)
	buffer.WriteString(tabsLayoutActiveTabStyle(tabsLayout))
	buffer.WriteString(`" data-current="`)
	buffer.WriteString(strconv.Itoa(GetCurrent(tabsLayout)))
	buffer.WriteRune('"')
}

func (tabsLayout *tabsLayoutData) cssStyle(self View, builder cssBuilder) {
	tabsLayout.viewsContainerData.cssStyle(self, builder)
	switch tabsLayout.tabsLocation() {
	case TopTabs:
		builder.add(`grid-template-rows`, `auto 1fr`)

	case BottomTabs:
		builder.add(`grid-template-rows`, `1fr auto`)

	case LeftTabs, LeftListTabs:
		builder.add(`grid-template-columns`, `auto 1fr`)

	case RightTabs, RightListTabs:
		builder.add(`grid-template-columns`, `1fr auto`)
	}
}

func (tabsLayout *tabsLayoutData) htmlSubviews(self View, buffer *strings.Builder) {
	if tabsLayout.views == nil {
		return
	}

	current := GetCurrent(tabsLayout)
	location := tabsLayout.tabsLocation()
	tabsLayoutID := tabsLayout.htmlID()

	if location != HiddenTabs {

		buffer.WriteString(`<div class="`)
		buffer.WriteString(tabsLayout.tabBarStyle())
		buffer.WriteString(`" style="display: flex;`)

		switch location {
		case LeftTabs, LeftListTabs, TopTabs:
			buffer.WriteString(` grid-row-start: 1; grid-row-end: 2; grid-column-start: 1; grid-column-end: 2;`)

		case RightTabs, RightListTabs:
			buffer.WriteString(` grid-row-start: 1; grid-row-end: 2; grid-column-start: 2; grid-column-end: 3;`)

		case BottomTabs:
			buffer.WriteString(` grid-row-start: 2; grid-row-end: 3; grid-column-start: 1; grid-column-end: 2;`)
		}

		buffer.WriteString(` flex-flow: `)
		switch location {
		case LeftTabs, LeftListTabs, RightTabs, RightListTabs:
			buffer.WriteString(`column nowrap; justify-content: flex-start; align-items: stretch;`)

		default:
			buffer.WriteString(`row nowrap; justify-content: flex-start; align-items: stretch;`)
		}

		buffer.WriteString(`">`)

		inactiveStyle := tabsLayoutInactiveTabStyle(tabsLayout)
		activeStyle := tabsLayoutActiveTabStyle(tabsLayout)

		notTranslate := GetNotTranslate(tabsLayout)
		closeButton, _ := boolProperty(tabsLayout, TabCloseButton, tabsLayout.session)

		var tabStyle, titleStyle string
		switch location {
		case LeftTabs, RightTabs:
			tabStyle = `display: grid; grid-template-rows: auto 1fr auto; align-items: center; justify-items: center; grid-row-gap: 8px;`

		case LeftListTabs, RightListTabs:
			tabStyle = `display: grid; grid-template-columns: auto 1fr auto; align-items: center; justify-items: start; grid-column-gap: 8px;`

		default:
			tabStyle = `display: grid; grid-template-columns: auto 1fr auto; align-items: center; justify-items: center; grid-column-gap: 8px;`
		}

		switch location {
		case LeftTabs:
			titleStyle = ` style="writing-mode: vertical-lr; transform: rotate(180deg); grid-row-start: 2; grid-row-end: 3; grid-column-start: 1; grid-column-end: 2;">`

		case RightTabs:
			titleStyle = ` style="writing-mode: vertical-lr; grid-row-start: 2; grid-row-end: 3; grid-column-start: 1; grid-column-end: 2;">`

		default:
			titleStyle = ` style="grid-row-start: 1; grid-row-end: 2; grid-column-start: 2; grid-column-end: 3;">`
		}

		for n, view := range tabsLayout.views {
			icon, _ := imageProperty(view, Icon, tabsLayout.session)
			title, _ := stringProperty(view, Title, tabsLayout.session)
			if !notTranslate {
				title, _ = tabsLayout.Session().GetString(title)
			}

			buffer.WriteString(`<div id="`)
			buffer.WriteString(tabsLayoutID)
			buffer.WriteByte('-')
			buffer.WriteString(strconv.Itoa(n))
			buffer.WriteString(`" class="`)
			if n == current {
				buffer.WriteString(activeStyle)
			} else {
				buffer.WriteString(inactiveStyle)
			}
			buffer.WriteString(`" tabindex="0" onclick="tabClickEvent(this, '`)
			buffer.WriteString(tabsLayoutID)
			buffer.WriteString(`', `)
			buffer.WriteString(strconv.Itoa(n))
			buffer.WriteString(`, event)" onkeydown="tabKeyClickEvent('`)
			buffer.WriteString(tabsLayoutID)
			buffer.WriteString(`', `)
			buffer.WriteString(strconv.Itoa(n))
			buffer.WriteString(`, event)" style="`)
			buffer.WriteString(tabStyle)
			buffer.WriteString(`" data-container="`)
			buffer.WriteString(tabsLayoutID)
			buffer.WriteString(`" data-view="`)
			buffer.WriteString(tabsLayoutID)
			buffer.WriteString(`-page`)
			buffer.WriteString(strconv.Itoa(n))
			buffer.WriteString(`">`)

			if icon != "" {
				buffer.WriteString(`<img id="`)
				buffer.WriteString(view.htmlID())
				switch location {
				case LeftTabs:
					buffer.WriteString(`-icon" style="grid-row-start: 3; grid-row-end: 4; grid-column-start: 1; grid-column-end: 2;" src="`)

				case RightTabs:
					buffer.WriteString(`-icon" style="grid-row-start: 1; grid-row-end: 2; grid-column-start: 1; grid-column-end: 2;" src="`)

				default:
					buffer.WriteString(`-icon" style="grid-row-start: 1; grid-row-end: 2; grid-column-start: 1; grid-column-end: 2;" src="`)
				}
				buffer.WriteString(icon)
				buffer.WriteString(`">`)
			}

			buffer.WriteString(`<div id="`)
			buffer.WriteString(view.htmlID())
			buffer.WriteString(`-title"`)
			buffer.WriteString(titleStyle)
			buffer.WriteString(title)
			buffer.WriteString(`</div>`)

			close, ok := boolProperty(view, TabCloseButton, tabsLayout.session)
			if !ok {
				close = closeButton
			}
			if close {
				buffer.WriteString(`<div class="ruiTabCloseButton" tabindex="0" onclick="tabCloseClickEvent(this, '`)
				buffer.WriteString(tabsLayoutID)
				buffer.WriteString(`', `)
				buffer.WriteString(strconv.Itoa(n))
				buffer.WriteString(`, event)" onkeydown="tabCloseKeyClickEvent('`)
				buffer.WriteString(tabsLayoutID)
				buffer.WriteString(`', `)
				buffer.WriteString(strconv.Itoa(n))
				buffer.WriteString(`, event)" style="display: grid; `)

				switch location {
				case LeftTabs:
					buffer.WriteString(`grid-row-start: 1; grid-row-end: 2; grid-column-start: 1; grid-column-end: 2;">`)

				case RightTabs:
					buffer.WriteString(`grid-row-start: 3; grid-row-end: 4; grid-column-start: 1; grid-column-end: 2;">`)

				default:
					buffer.WriteString(`grid-row-start: 1; grid-row-end: 2; grid-column-start: 3; grid-column-end: 4;">`)
				}

				buffer.WriteString(`✕</div>`)
			}

			buffer.WriteString(`</div>`)
		}

		buffer.WriteString(`</div>`)
	}

	for n, view := range tabsLayout.views {
		buffer.WriteString(`<div id="`)
		buffer.WriteString(tabsLayoutID)
		buffer.WriteString(`-page`)
		buffer.WriteString(strconv.Itoa(n))

		if current != n {
			buffer.WriteString(`" style="display: grid; align-items: stretch; justify-items: stretch; visibility: hidden; `)
		} else {
			buffer.WriteString(`" style="display: grid; align-items: stretch; justify-items: stretch; `)
		}

		switch location {
		case LeftTabs, LeftListTabs:
			buffer.WriteString(`grid-row-start: 1; grid-row-end: 2; grid-column-start: 2; grid-column-end: 3;">`)

		case TopTabs:
			buffer.WriteString(`grid-row-start: 2; grid-row-end: 3; grid-column-start: 1; grid-column-end: 2;">`)

		default:
			buffer.WriteString(`grid-row-start: 1; grid-row-end: 2; grid-column-start: 1; grid-column-end: 2;">`)
		}

		viewHTML(view, buffer)
		buffer.WriteString(`</div>`)
	}
}

func (tabsLayout *tabsLayoutData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case "tabClick":
		if numberText, ok := data.PropertyValue("number"); ok {
			if number, err := strconv.Atoi(numberText); err == nil {
				current := GetCurrent(tabsLayout)
				if current != number {
					tabsLayout.setRaw(Current, number)

					tabsLayout.currentChanged(number, current)
				}
			}
		}
		return true

	case "tabCloseClick":
		if numberText, ok := data.PropertyValue("number"); ok {
			if number, err := strconv.Atoi(numberText); err == nil {
				for _, listener := range getOneArgEventListeners[TabsLayout, int](tabsLayout, nil, TabCloseEvent) {
					listener(tabsLayout, number)
				}
			}
		}
		return true
	}
	return tabsLayout.viewsContainerData.handleCommand(self, command, data)
}
