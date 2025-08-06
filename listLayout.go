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
	//return NewListLayout(session, nil)
	return new(listLayoutData)
}

// Init initialize fields of ViewsAlignContainer by default values
func (listLayout *listLayoutData) init(session Session) {
	listLayout.viewsContainerData.init(session)
	listLayout.tag = "ListLayout"
	listLayout.systemClass = "ruiListLayout"
	listLayout.normalize = normalizeListLayoutTag
	listLayout.get = listLayout.getFunc
	listLayout.set = listLayout.setFunc
	listLayout.remove = listLayout.removeFunc
	listLayout.changed = listLayout.propertyChanged
}

func normalizeListLayoutTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
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

func (listLayout *listLayoutData) getFunc(tag PropertyName) any {
	switch tag {
	case Gap:
		if rowGap := GetListRowGap(listLayout); rowGap.Equal(GetListColumnGap(listLayout)) {
			return rowGap
		}
		return AutoSize()

	case Content:
		if listLayout.adapter != nil {
			return listLayout.adapter
		}
	}

	return listLayout.viewsContainerData.getFunc(tag)
}

func (listLayout *listLayoutData) removeFunc(tag PropertyName) []PropertyName {
	switch tag {
	case Gap:
		result := []PropertyName{}
		for _, tag := range []PropertyName{ListRowGap, ListColumnGap} {
			if listLayout.getRaw(tag) != nil {
				listLayout.setRaw(tag, nil)
				result = append(result, tag)
			}
		}
		return result

	case Content:
		result := listLayout.viewsContainerData.removeFunc(Content)
		if listLayout.adapter != nil {
			listLayout.adapter = nil
			return []PropertyName{Content}
		}
		return result
	}

	return listLayout.viewsContainerData.removeFunc(tag)
}

func (listLayout *listLayoutData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Gap:
		result := listLayout.setFunc(ListRowGap, value)
		if result != nil {
			if gap := listLayout.getRaw(ListRowGap); gap != nil {
				listLayout.setRaw(ListColumnGap, gap)
				result = append(result, ListColumnGap)
			}
		}
		return result

	case Content:
		if adapter, ok := value.(ListAdapter); ok {
			listLayout.adapter = adapter
			listLayout.createContent()
		} else if listLayout.setContent(value) {
			listLayout.adapter = nil
		} else {
			return nil
		}
		return []PropertyName{Content}
	}
	return listLayout.viewsContainerData.setFunc(tag, value)
}

func (listLayout *listLayoutData) propertyChanged(tag PropertyName) {
	switch tag {
	case Orientation, ListWrap, HorizontalAlign, VerticalAlign:
		updateCSSStyle(listLayout.htmlID(), listLayout.Session())

	default:
		listLayout.viewsContainerData.propertyChanged(tag)
	}
}

func (listLayout *listLayoutData) htmlSubviews(self View, buffer *strings.Builder) {
	if listLayout.views != nil {
		for _, view := range listLayout.views {
			view.addToCSSStyle(map[string]string{`flex`: `0 0 auto`})
			viewHTML(view, buffer, "")
		}
	}
}

func (listLayout *listLayoutData) createContent() bool {
	if adapter := listLayout.adapter; adapter != nil {
		listLayout.views = []View{}

		session := listLayout.session
		htmlID := listLayout.htmlID()
		isDisabled := IsDisabled(listLayout)

		for i := range adapter.ListSize() {
			if view := adapter.ListItem(i, session); view != nil {
				view.setParentID(htmlID)
				if isDisabled {
					view.Set(Disabled, true)
				}
				listLayout.views = append(listLayout.views, view)
			}
		}

		return true
	}
	return false
}

func (listLayout *listLayoutData) UpdateContent() {
	if listLayout.createContent() {
		if listLayout.created {
			updateInnerHTML(listLayout.htmlID(), listLayout.session)
		}
		listLayout.runChangeListener(Content)
	}
}

// GetListVerticalAlign returns the vertical align of a ListLayout or ListView subview:
// TopAlign (0), BottomAlign (1), CenterAlign (2), or StretchAlign (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, VerticalAlign, TopAlign, false)
}

// GetListHorizontalAlign returns the vertical align of a ListLayout or ListView subview:
// LeftAlign (0), RightAlign (1), CenterAlign (2), or StretchAlign (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, HorizontalAlign, LeftAlign, false)
}

// GetListOrientation returns the orientation of a ListLayout or ListView subview:
// TopDownOrientation (0), StartToEndOrientation (1), BottomUpOrientation (2), or EndToStartOrientation (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListOrientation(view View, subviewID ...string) int {
	if view = getSubview(view, subviewID); view != nil {
		if orientation, ok := valueToOrientation(view.Get(Orientation), view.Session()); ok {
			return orientation
		}

		if value := valueFromStyle(view, Orientation); value != nil {
			if orientation, ok := valueToOrientation(value, view.Session()); ok {
				return orientation
			}
		}
	}

	return TopDownOrientation
}

// GetListWrap returns the wrap type of a ListLayout or ListView subview:
// ListWrapOff (0), ListWrapOn (1), or ListWrapReverse (2)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListWrap(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, ListWrap, ListWrapOff, false)
}

// GetListRowGap returns the gap between ListLayout or ListView rows.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListRowGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ListRowGap, false)
}

// GetListColumnGap returns the gap between ListLayout or ListView columns.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetListColumnGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ListColumnGap, false)
}

// UpdateContent updates child Views of ListLayout/GridLayout subview if the "content" property value is set to ListAdapter/GridAdapter,
// otherwise does nothing.
// If the second argument (subviewID) is not specified or it is "" then the first argument (view) updates.
func UpdateContent(view View, subviewID ...string) {
	if view = getSubview(view, subviewID); view != nil {
		switch view := view.(type) {
		case GridLayout:
			view.UpdateGridContent()

		case ListLayout:
			view.UpdateContent()

		case ListView:
			view.ReloadListViewData()

		case TableView:
			view.ReloadTableData()
		}
	}
}
