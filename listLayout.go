package rui

import (
	"strings"
)

// Constants which represent values of the "orientation" property of the [ListLayout]
const (
	// TopDownOrientation - subviews are arranged from top to bottom. Synonym of VerticalOrientation
	TopDownOrientation = 0

	// StartToEndOrientation - subviews are arranged from left to right. Synonym of HorizontalOrientation
	StartToEndOrientation = 1

	// BottomUpOrientation - subviews are arranged from bottom to top
	BottomUpOrientation = 2

	// EndToStartOrientation - subviews are arranged from right to left
	EndToStartOrientation = 3
)

// Constants which represent values of the "list-wrap" property of the [ListLayout]
const (
	// ListWrapOff - subviews are scrolled and "true" if a new row/column starts
	ListWrapOff = 0

	// ListWrapOn - the new row/column starts at bottom/right
	ListWrapOn = 1

	// ListWrapReverse - the new row/column starts at top/left
	ListWrapReverse = 2
)

// ListLayout represents a ListLayout view
type ListLayout interface {
	ViewsContainer
	// UpdateContent updates child Views if the "content" property value is set to ListAdapter,
	// otherwise does nothing
	UpdateContent()
}

type listLayoutData struct {
	viewsContainerData
	adapter ListAdapter
}

// NewListLayout create new ListLayout object and return it
func NewListLayout(session Session, params Params) ListLayout {
	view := new(listLayoutData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newListLayout(session Session) View {
	return NewListLayout(session, nil)
}

// Init initialize fields of ViewsAlignContainer by default values
func (listLayout *listLayoutData) init(session Session) {
	listLayout.viewsContainerData.init(session)
	listLayout.tag = "ListLayout"
	listLayout.systemClass = "ruiListLayout"
}

func (listLayout *listLayoutData) String() string {
	return getViewString(listLayout, nil)
}

func (listLayout *listLayoutData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case "wrap":
		tag = ListWrap

	case "row-gap":
		return ListRowGap

	case ColumnGap:
		return ListColumnGap
	}
	return tag
}

func (listLayout *listLayoutData) Get(tag string) any {
	return listLayout.get(listLayout.normalizeTag(tag))
}

func (listLayout *listLayoutData) get(tag string) any {
	if tag == Gap {
		if rowGap := GetListRowGap(listLayout); rowGap.Equal(GetListColumnGap(listLayout)) {
			return rowGap
		}
		return AutoSize()
	}

	return listLayout.viewsContainerData.get(tag)
}

func (listLayout *listLayoutData) Remove(tag string) {
	listLayout.remove(listLayout.normalizeTag(tag))
}

func (listLayout *listLayoutData) remove(tag string) {
	switch tag {
	case Gap:
		listLayout.remove(ListRowGap)
		listLayout.remove(ListColumnGap)
		return

	case Content:
		listLayout.adapter = nil
	}

	listLayout.viewsContainerData.remove(tag)
	if listLayout.created {
		switch tag {
		case Orientation, ListWrap, HorizontalAlign, VerticalAlign:
			updateCSSStyle(listLayout.htmlID(), listLayout.session)
		}
	}
}

func (listLayout *listLayoutData) Set(tag string, value any) bool {
	return listLayout.set(listLayout.normalizeTag(tag), value)
}

func (listLayout *listLayoutData) set(tag string, value any) bool {
	if value == nil {
		listLayout.remove(tag)
		return true
	}

	switch tag {
	case Gap:
		return listLayout.set(ListRowGap, value) && listLayout.set(ListColumnGap, value)

	case Content:
		if adapter, ok := value.(ListAdapter); ok {
			listLayout.adapter = adapter
			listLayout.UpdateContent()
			// TODO
			return true
		}
		listLayout.adapter = nil
	}

	if listLayout.viewsContainerData.set(tag, value) {
		if listLayout.created {
			switch tag {
			case Orientation, ListWrap, HorizontalAlign, VerticalAlign:
				updateCSSStyle(listLayout.htmlID(), listLayout.session)
			}
		}
		return true
	}
	return false
}

func (listLayout *listLayoutData) htmlSubviews(self View, buffer *strings.Builder) {
	if listLayout.views != nil {
		for _, view := range listLayout.views {
			view.addToCSSStyle(map[string]string{`flex`: `0 0 auto`})
			viewHTML(view, buffer)
		}
	}
}

func (listLayout *listLayoutData) UpdateContent() {
	if adapter := listLayout.adapter; adapter != nil {
		listLayout.views = []View{}

		session := listLayout.session
		htmlID := listLayout.htmlID()
		isDisabled := IsDisabled(listLayout)

		count := adapter.ListSize()
		for i := 0; i < count; i++ {
			if view := adapter.ListItem(i, session); view != nil {
				view.setParentID(htmlID)
				if isDisabled {
					view.Set(Disabled, true)
				}
				listLayout.views = append(listLayout.views, view)
			}
		}

		if listLayout.created {
			updateInnerHTML(htmlID, session)
		}

		listLayout.propertyChangedEvent(Content)
	}
}

// GetListVerticalAlign returns the vertical align of a ListLayout or ListView sibview:
// TopAlign (0), BottomAlign (1), CenterAlign (2), or StretchAlign (3)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetListVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, VerticalAlign, TopAlign, false)
}

// GetListHorizontalAlign returns the vertical align of a ListLayout or ListView subview:
// LeftAlign (0), RightAlign (1), CenterAlign (2), or StretchAlign (3)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetListHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, HorizontalAlign, LeftAlign, false)
}

// GetListOrientation returns the orientation of a ListLayout or ListView subview:
// TopDownOrientation (0), StartToEndOrientation (1), BottomUpOrientation (2), or EndToStartOrientation (3)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetListOrientation(view View, subviewID ...string) int {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if orientation, ok := valueToOrientation(view.Get(Orientation), view.Session()); ok {
			return orientation
		}

		if value := valueFromStyle(view, Orientation); value != nil {
			if orientation, ok := valueToOrientation(value, view.Session()); ok {
				return orientation
			}
		}
	}

	return 0
}

// GetListWrap returns the wrap type of a ListLayout or ListView subview:
// ListWrapOff (0), ListWrapOn (1), or ListWrapReverse (2)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetListWrap(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, ListWrap, ListWrapOff, false)
}

// GetListRowGap returns the gap between ListLayout or ListView rows.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetListRowGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ListRowGap, false)
}

// GetListColumnGap returns the gap between ListLayout or ListView columns.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetListColumnGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ListColumnGap, false)
}

// UpdateContent updates child Views of ListLayout/GridLayout subview if the "content" property value is set to ListAdapter/GridAdapter,
// otherwise does nothing.
// If the second argument (subviewID) is not specified or it is "" then the first argument (view) updates.
func UpdateContent(view View, subviewID ...string) {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		switch view := view.(type) {
		case GridLayout:
			view.UpdateGridContent()

		case ListLayout:
			view.UpdateContent()

		case ListView:
			view.ReloadListViewData()
		}
	}
}
