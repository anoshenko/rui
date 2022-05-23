package rui

import (
	"strings"
)

const (
	// TopDownOrientation - subviews are arranged from top to bottom. Synonym of VerticalOrientation
	TopDownOrientation = 0
	// StartToEndOrientation - subviews are arranged from left to right. Synonym of HorizontalOrientation
	StartToEndOrientation = 1
	// BottomUpOrientation - subviews are arranged from bottom to top
	BottomUpOrientation = 2
	// EndToStartOrientation - subviews are arranged from right to left
	EndToStartOrientation = 3
	// WrapOff - subviews are scrolled and "true" if a new row/column starts
	WrapOff = 0
	// WrapOn - the new row/column starts at bottom/right
	WrapOn = 1
	// WrapReverse - the new row/column starts at top/left
	WrapReverse = 2
)

// ListLayout - list-container of View
type ListLayout interface {
	ViewsContainer
}

type listLayoutData struct {
	viewsContainerData
}

// NewListLayout create new ListLayout object and return it
func NewListLayout(session Session, params Params) ListLayout {
	view := new(listLayoutData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newListLayout(session Session) View {
	return NewListLayout(session, nil)
}

// Init initialize fields of ViewsAlignContainer by default values
func (listLayout *listLayoutData) Init(session Session) {
	listLayout.viewsContainerData.Init(session)
	listLayout.tag = "ListLayout"
	listLayout.systemClass = "ruiListLayout"
}

func (listLayout *listLayoutData) String() string {
	return getViewString(listLayout)
}

func (listLayout *listLayoutData) Remove(tag string) {
	listLayout.remove(strings.ToLower(tag))
}

func (listLayout *listLayoutData) remove(tag string) {
	listLayout.viewsContainerData.remove(tag)
	if listLayout.created {
		switch tag {
		case Orientation, Wrap, HorizontalAlign, VerticalAlign:
			updateCSSStyle(listLayout.htmlID(), listLayout.session)
		}
	}
}

func (listLayout *listLayoutData) Set(tag string, value interface{}) bool {
	return listLayout.set(strings.ToLower(tag), value)
}

func (listLayout *listLayoutData) set(tag string, value interface{}) bool {
	if value == nil {
		listLayout.remove(tag)
		return true
	}

	if listLayout.viewsContainerData.set(tag, value) {
		if listLayout.created {
			switch tag {
			case Orientation, Wrap, HorizontalAlign, VerticalAlign:
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

// GetListVerticalAlign returns the vertical align of a ListLayout or ListView sibview:
// TopAlign (0), BottomAlign (1), CenterAlign (2), or StretchAlign (3)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListVerticalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return LeftAlign
	}
	result, _ := enumProperty(view, VerticalAlign, view.Session(), 0)
	return result
}

// GetListHorizontalAlign returns the vertical align of a ListLayout or ListView subview:
// LeftAlign (0), RightAlign (1), CenterAlign (2), or StretchAlign (3)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListHorizontalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view == nil {
		return TopAlign
	}
	result, _ := enumProperty(view, HorizontalAlign, view.Session(), 0)
	return result
}

// GetListOrientation returns the orientation of a ListLayout or ListView subview:
// TopDownOrientation (0), StartToEndOrientation (1), BottomUpOrientation (2), or EndToStartOrientation (3)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListOrientation(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
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
// WrapOff (0), WrapOn (1), or WrapReverse (2)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetListWrap(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if result, ok := enumStyledProperty(view, Wrap, 0); ok {
			return result
		}
	}
	return WrapOff
}
