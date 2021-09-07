package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// HiddenTabs - tabs of TabsLayout are hidden
	HiddenTabs = 0
	// TopTabs - tabs of TabsLayout are on the top
	TopTabs = 1
	// BottomTabs - tabs of TabsLayout are on the bottom
	BottomTabs = 2
	// LeftTabs - tabs of TabsLayout are on the left
	LeftTabs = 3
	// RightTabs - tabs of TabsLayout are on the right
	RightTabs = 4
	// LeftListTabs - tabs of TabsLayout are on the left
	LeftListTabs = 5
	// RightListTabs - tabs of TabsLayout are on the right
	RightListTabs = 6
)

// TabsLayoutCurrentChangedListener - listener of the current tab changing
type TabsLayoutCurrentChangedListener interface {
	OnTabsLayoutCurrentChanged(tabsLayout TabsLayout, newCurrent int, newCurrentView View, oldCurrent int, oldCurrentView View)
}

type tabsLayoutCurrentChangedListenerFunc struct {
	listenerFunc func(tabsLayout TabsLayout, newCurrent int, newCurrentView View, oldCurrent int, oldCurrentView View)
}

func (listener *tabsLayoutCurrentChangedListenerFunc) OnTabsLayoutCurrentChanged(tabsLayout TabsLayout,
	newCurrent int, newCurrentView View, oldCurrent int, oldCurrentView View) {
	if listener.listenerFunc != nil {
		listener.listenerFunc(tabsLayout, newCurrent, newCurrentView, oldCurrent, oldCurrentView)
	}
}

// TabsLayout - multi-tab container of View
type TabsLayout interface {
	ViewsContainer
	/*
		// Current return the index of active tab
		currentItem() int
		// SetCurrent set the index of active tab
		SetCurrent(current int)
		// TabsLocation return the location of tabs. It returns one of the following values: HiddenTabs (0),
		//   TopTabs (1), BottomTabs (2), LeftTabs (3), RightTabs (4), LeftListTabs (5), RightListTabs (6)
		tabsLocation() int
		// TabsLocation set the location of tabs. Valid values: HiddenTabs (0), TopTabs (1),
		//   BottomTabs (2), LeftTabs (3), RightTabs (4), LeftListTabs (5), RightListTabs (6)
		SetTabsLocation(location int)
		// TabStyle() return styles of tab in the passive and the active state
		TabStyle() (string, string)
		SetTabStyle(tabStyle string, activeTabStyle string)
	*/
	// SetCurrentTabChangedListener add the listener of the current tab changing
	SetCurrentTabChangedListener(listener TabsLayoutCurrentChangedListener)
	// SetCurrentTabChangedListener add the listener function of the current tab changing
	SetCurrentTabChangedListenerFunc(listenerFunc func(tabsLayout TabsLayout,
		newCurrent int, newCurrentView View, oldCurrent int, oldCurrentView View))
}

type tabsLayoutData struct {
	viewsContainerData
	//currentTab, tabsLocation int
	//tabStyle, activeTabStyle string
	tabListener TabsLayoutCurrentChangedListener
}

// NewTabsLayout create new TabsLayout object and return it
func NewTabsLayout(session Session) TabsLayout {
	view := new(tabsLayoutData)
	view.Init(session)
	return view
}

func newTabsLayout(session Session) View {
	return NewTabsLayout(session)
}

// Init initialize fields of ViewsContainer by default values
func (tabsLayout *tabsLayoutData) Init(session Session) {
	tabsLayout.viewsContainerData.Init(session)
	tabsLayout.tag = "TabsLayout"
	tabsLayout.systemClass = "ruiTabsLayout"
	tabsLayout.tabListener = nil
}

func (tabsLayout *tabsLayoutData) currentItem() int {
	result, _ := intProperty(tabsLayout, Current, tabsLayout.session, 0)
	return result
}

func (tabsLayout *tabsLayoutData) Set(tag string, value interface{}) bool {
	switch tag {
	case Current:
		oldCurrent := tabsLayout.currentItem()
		if !tabsLayout.setIntProperty(Current, value) {
			return false
		}

		if !tabsLayout.session.ignoreViewUpdates() {
			current := tabsLayout.currentItem()
			if oldCurrent != current {
				tabsLayout.session.runScript(fmt.Sprintf("activateTab(%v, %d);", tabsLayout.htmlID(), current))
				if tabsLayout.tabListener != nil {
					oldView := tabsLayout.views[oldCurrent]
					view := tabsLayout.views[current]
					tabsLayout.tabListener.OnTabsLayoutCurrentChanged(tabsLayout, current, view, oldCurrent, oldView)
				}
			}
		}

	case Tabs:
		if !tabsLayout.setEnumProperty(Tabs, value, enumProperties[Tabs].values) {
			return false
		}
		if !tabsLayout.session.ignoreViewUpdates() {
			htmlID := tabsLayout.htmlID()
			updateCSSStyle(htmlID, tabsLayout.session)
			updateInnerHTML(htmlID, tabsLayout.session)
		}

	case TabStyle, CurrentTabStyle:
		if value == nil {
			delete(tabsLayout.properties, tag)
		} else if text, ok := value.(string); ok {
			if text == "" {
				delete(tabsLayout.properties, tag)
			} else {
				tabsLayout.properties[tag] = text
			}
		} else {
			notCompatibleType(tag, value)
			return false
		}

		if !tabsLayout.session.ignoreViewUpdates() {
			htmlID := tabsLayout.htmlID()
			updateProperty(htmlID, "data-tabStyle", tabsLayout.inactiveTabStyle(), tabsLayout.session)
			updateProperty(htmlID, "data-activeTabStyle", tabsLayout.activeTabStyle(), tabsLayout.session)
			updateInnerHTML(htmlID, tabsLayout.session)
		}

	default:
		return tabsLayout.viewsContainerData.Set(tag, value)
	}

	return true
}

func (tabsLayout *tabsLayoutData) tabsLocation() int {
	tabs, _ := enumProperty(tabsLayout, Tabs, tabsLayout.session, 0)
	return tabs
}

func (tabsLayout *tabsLayoutData) inactiveTabStyle() string {
	if style, ok := stringProperty(tabsLayout, TabStyle, tabsLayout.session); ok {
		return style
	}
	switch tabsLayout.tabsLocation() {
	case LeftTabs, RightTabs:
		return "ruiInactiveVerticalTab"
	}
	return "ruiInactiveTab"
}

func (tabsLayout *tabsLayoutData) activeTabStyle() string {
	if style, ok := stringProperty(tabsLayout, CurrentTabStyle, tabsLayout.session); ok {
		return style
	}
	switch tabsLayout.tabsLocation() {
	case LeftTabs, RightTabs:
		return "ruiActiveVerticalTab"
	}
	return "ruiActiveTab"
}

func (tabsLayout *tabsLayoutData) TabStyle() (string, string) {
	return tabsLayout.inactiveTabStyle(), tabsLayout.activeTabStyle()
}

func (tabsLayout *tabsLayoutData) SetCurrentTabChangedListener(listener TabsLayoutCurrentChangedListener) {
	tabsLayout.tabListener = listener
}

/*
// SetCurrentTabChangedListener add the listener of the current tab changing
func (tabsLayout *tabsLayoutData) SetCurrentTabChangedListener(listener TabsLayoutCurrentChangedListener) {
	tabsLayout.tabListener = listener
}

// SetCurrentTabChangedListener add the listener function of the current tab changing
func (tabsLayout *tabsLayoutData) SetCurrentTabChangedListenerFunc(listenerFunc func(tabsLayout TabsLayout,
	newCurrent int, newCurrentView View, oldCurrent int, oldCurrentView View)) {
	}
*/

func (tabsLayout *tabsLayoutData) SetCurrentTabChangedListenerFunc(listenerFunc func(tabsLayout TabsLayout,
	newCurrent int, newCurrentView View, oldCurrent int, oldCurrentView View)) {
	listener := new(tabsLayoutCurrentChangedListenerFunc)
	listener.listenerFunc = listenerFunc
	tabsLayout.SetCurrentTabChangedListener(listener)
}

// Append appends view to the end of view list
func (tabsLayout *tabsLayoutData) Append(view View) {
	if tabsLayout.views == nil {
		tabsLayout.views = []View{}
	}
	tabsLayout.viewsContainerData.Append(view)
	if len(tabsLayout.views) == 1 {
		tabsLayout.setIntProperty(Current, 0)
		if tabsLayout.tabListener != nil {
			tabsLayout.tabListener.OnTabsLayoutCurrentChanged(tabsLayout, 0, tabsLayout.views[0], -1, nil)
		}
	}
	updateInnerHTML(tabsLayout.htmlID(), tabsLayout.session)
}

// Insert inserts view to the "index" position in view list
func (tabsLayout *tabsLayoutData) Insert(view View, index uint) {
	if tabsLayout.views == nil {
		tabsLayout.views = []View{}
	}
	tabsLayout.viewsContainerData.Insert(view, index)
	current := tabsLayout.currentItem()
	if current >= int(index) {
		tabsLayout.Set(Current, current+1)
	}
}

// Remove removes view from list and return it
func (tabsLayout *tabsLayoutData) RemoveView(index uint) View {
	if tabsLayout.views == nil {
		tabsLayout.views = []View{}
		return nil
	}
	i := int(index)
	count := len(tabsLayout.views)
	if i >= count {
		return nil
	}

	if count == 1 {
		view := tabsLayout.views[0]
		tabsLayout.views = []View{}
		if tabsLayout.tabListener != nil {
			tabsLayout.tabListener.OnTabsLayoutCurrentChanged(tabsLayout, 0, nil, 0, view)
		}
		return view
	}

	current := tabsLayout.currentItem()
	removeCurrent := (i == current)
	if i < current || (removeCurrent && i == count-1) {
		tabsLayout.properties[Current] = current - 1
		if tabsLayout.tabListener != nil {
			currentView := tabsLayout.views[current-1]
			oldCurrentView := tabsLayout.views[current]
			tabsLayout.tabListener.OnTabsLayoutCurrentChanged(tabsLayout, current-1, currentView, current, oldCurrentView)
		}
	}

	return tabsLayout.viewsContainerData.RemoveView(index)
}

func (tabsLayout *tabsLayoutData) htmlProperties(self View, buffer *strings.Builder) {
	tabsLayout.viewsContainerData.htmlProperties(self, buffer)
	buffer.WriteString(` data-inactiveTabStyle="`)
	buffer.WriteString(tabsLayout.inactiveTabStyle())
	buffer.WriteString(`" data-activeTabStyle="`)
	buffer.WriteString(tabsLayout.activeTabStyle())
	buffer.WriteString(`" data-current="`)
	buffer.WriteString(tabsLayout.htmlID())
	buffer.WriteRune('-')
	buffer.WriteString(strconv.Itoa(tabsLayout.currentItem()))
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

	//viewCount := len(tabsLayout.views)
	current := tabsLayout.currentItem()
	location := tabsLayout.tabsLocation()
	tabsLayoutID := tabsLayout.htmlID()

	if location != HiddenTabs {
		tabsHeight, _ := sizeConstant(tabsLayout.session, "ruiTabHeight")
		tabsSpace, _ := sizeConstant(tabsLayout.session, "ruiTabSpace")
		rowLayout := false
		buffer.WriteString(`<div style="display: flex;`)

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
			if tabsHeight.Type != Auto {
				buffer.WriteString(` height: `)
				buffer.WriteString(tabsHeight.cssString(""))
				buffer.WriteByte(';')
			}
			rowLayout = true
		}

		var tabsPadding Bounds
		if value, ok := tabsLayout.session.Constant("ruiTabPadding"); ok {
			if tabsPadding.parse(value, tabsLayout.session) {
				if !tabsPadding.allFieldsAuto() {
					buffer.WriteByte(' ')
					buffer.WriteString(Padding)
					buffer.WriteString(`: `)
					tabsPadding.writeCSSString(buffer, "0")
					buffer.WriteByte(';')
				}
			}
		}

		if tabsBackground, ok := tabsLayout.session.Color("tabsBackgroundColor"); ok {
			buffer.WriteString(` background-color: `)
			buffer.WriteString(tabsBackground.cssString())
			buffer.WriteByte(';')
		}

		buffer.WriteString(`">`)

		inactiveStyle := tabsLayout.inactiveTabStyle()
		activeStyle := tabsLayout.activeTabStyle()

		notTranslate := GetNotTranslate(tabsLayout, "")
		last := len(tabsLayout.views) - 1
		for n, view := range tabsLayout.views {
			title, _ := stringProperty(view, "title", tabsLayout.session)
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
			buffer.WriteString(`" tabindex="0" onclick="tabClickEvent(\'`)
			buffer.WriteString(tabsLayoutID)
			buffer.WriteString(`\', `)
			buffer.WriteString(strconv.Itoa(n))
			buffer.WriteString(`, event)`)
			buffer.WriteString(`" onclick="tabKeyClickEvent(\'`)
			buffer.WriteString(tabsLayoutID)
			buffer.WriteString(`\', `)
			buffer.WriteString(strconv.Itoa(n))
			buffer.WriteString(`, event)" style="display: flex; flex-flow: row nowrap; justify-content: center; align-items: center; `)

			if n != last && tabsSpace.Type != Auto && tabsSpace.Value > 0 {
				if rowLayout {
					buffer.WriteString(` margin-right: `)
					buffer.WriteString(tabsSpace.cssString(""))
				} else {
					buffer.WriteString(` margin-bottom: `)
					buffer.WriteString(tabsSpace.cssString(""))
				}
				buffer.WriteByte(';')
			}

			switch location {
			case LeftListTabs, RightListTabs:
				if tabsHeight.Type != Auto {
					buffer.WriteString(` height: `)
					buffer.WriteString(tabsHeight.cssString(""))
					buffer.WriteByte(';')
				}
			}

			buffer.WriteString(`" data-container="`)
			buffer.WriteString(tabsLayoutID)
			buffer.WriteString(`" data-view="`)
			//buffer.WriteString(view.htmlID())
			buffer.WriteString(tabsLayoutID)
			buffer.WriteString(`-page`)
			buffer.WriteString(strconv.Itoa(n))
			buffer.WriteString(`"><div`)

			switch location {
			case LeftTabs:
				buffer.WriteString(` style="writing-mode: vertical-lr; transform: rotate(180deg)">`)

			case RightTabs:
				buffer.WriteString(` style="writing-mode: vertical-lr;">`)

			default:
				buffer.WriteByte('>')
			}
			buffer.WriteString(title)
			buffer.WriteString(`</div></div>`)
		}

		buffer.WriteString(`</div>`)
	}

	for n, view := range tabsLayout.views {
		buffer.WriteString(`<div id="`)
		buffer.WriteString(tabsLayoutID)
		buffer.WriteString(`-page`)
		buffer.WriteString(strconv.Itoa(n))

		switch location {
		case LeftTabs, LeftListTabs:
			buffer.WriteString(`" style="position: relative; grid-row-start: 1; grid-row-end: 2; grid-column-start: 2; grid-column-end: 3;`)

		case TopTabs:
			buffer.WriteString(`" style="position: relative; grid-row-start: 2; grid-row-end: 3; grid-column-start: 1; grid-column-end: 2;`)

		default:
			buffer.WriteString(`" style="position: relative; grid-row-start: 1; grid-row-end: 2; grid-column-start: 1; grid-column-end: 2;`)
		}

		if current != n {
			buffer.WriteString(` display: none;`)
		}
		buffer.WriteString(`">`)

		view.addToCSSStyle(map[string]string{`position`: `absolute`, `left`: `0`, `right`: `0`, `top`: `0`, `bottom`: `0`})
		viewHTML(tabsLayout.views[n], buffer)
		buffer.WriteString(`</div>`)
	}
}

func (tabsLayout *tabsLayoutData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "tabClick":
		if numberText, ok := data.PropertyValue("number"); ok {
			if number, err := strconv.Atoi(numberText); err == nil {
				current := tabsLayout.currentItem()
				if current != number {
					tabsLayout.properties[Current] = number
					if tabsLayout.tabListener != nil {
						oldView := tabsLayout.views[current]
						view := tabsLayout.views[number]
						tabsLayout.tabListener.OnTabsLayoutCurrentChanged(tabsLayout, number, view, current, oldView)
					}
				}
			}
		}
		return true
	}
	return tabsLayout.viewsContainerData.handleCommand(self, command, data)
}
