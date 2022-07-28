package rui

import (
	"strconv"
	"strings"
)

const (
	// ColumnCount is the constant for the "column-count" property tag.
	// The "column-count" int property specifies number of columns into which the content is break
	// Values less than zero are not valid. if the "column-count" property value is 0 then
	// the number of columns is calculated based on the "column-width" property
	ColumnCount = "column-count"
	// ColumnWidth is the constant for the "column-width" property tag.
	// The "column-width" SizeUnit property specifies the width of each column.
	ColumnWidth = "column-width"
	// ColumnGap is the constant for the "column-gap" property tag.
	// The "column-width" SizeUnit property sets the size of the gap (gutter) between columns.
	ColumnGap = "column-gap"
	// ColumnSeparator is the constant for the "column-separator" property tag.
	// The "column-separator" property specifies the line drawn between columns in a multi-column layout.
	ColumnSeparator = "column-separator"
	// ColumnSeparatorStyle is the constant for the "column-separator-style" property tag.
	// The "column-separator-style" int property sets the style of the line drawn between
	// columns in a multi-column layout.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	ColumnSeparatorStyle = "column-separator-style"
	// ColumnSeparatorWidth is the constant for the "column-separator-width" property tag.
	// The "column-separator-width" SizeUnit property sets the width of the line drawn between
	// columns in a multi-column layout.
	ColumnSeparatorWidth = "column-separator-width"
	// ColumnSeparatorColor is the constant for the "column-separator-color" property tag.
	// The "column-separator-color" Color property sets the color of the line drawn between
	// columns in a multi-column layout.
	ColumnSeparatorColor = "column-separator-color"
)

// ColumnLayout - grid-container of View
type ColumnLayout interface {
	ViewsContainer
}

type columnLayoutData struct {
	viewsContainerData
}

// NewColumnLayout create new ColumnLayout object and return it
func NewColumnLayout(session Session, params Params) ColumnLayout {
	view := new(columnLayoutData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newColumnLayout(session Session) View {
	return NewColumnLayout(session, nil)
}

// Init initialize fields of ColumnLayout by default values
func (ColumnLayout *columnLayoutData) Init(session Session) {
	ColumnLayout.viewsContainerData.Init(session)
	ColumnLayout.tag = "ColumnLayout"
	//ColumnLayout.systemClass = "ruiColumnLayout"
}

func (columnLayout *columnLayoutData) String() string {
	return getViewString(columnLayout)
}

func (columnLayout *columnLayoutData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Gap:
		return ColumnGap
	}
	return tag
}

func (columnLayout *columnLayoutData) Get(tag string) any {
	return columnLayout.get(columnLayout.normalizeTag(tag))
}

func (columnLayout *columnLayoutData) Remove(tag string) {
	columnLayout.remove(columnLayout.normalizeTag(tag))
}

func (columnLayout *columnLayoutData) remove(tag string) {
	columnLayout.viewsContainerData.remove(tag)
	if columnLayout.created {
		switch tag {
		case ColumnCount, ColumnWidth, ColumnGap:
			updateCSSProperty(columnLayout.htmlID(), tag, "", columnLayout.Session())

		case ColumnSeparator:
			updateCSSProperty(columnLayout.htmlID(), "column-rule", "", columnLayout.Session())
		}
	}
}

func (columnLayout *columnLayoutData) Set(tag string, value any) bool {
	return columnLayout.set(columnLayout.normalizeTag(tag), value)
}

func (columnLayout *columnLayoutData) set(tag string, value any) bool {
	if value == nil {
		columnLayout.remove(tag)
		return true
	}

	if !columnLayout.viewsContainerData.set(tag, value) {
		return false
	}

	if columnLayout.created {
		switch tag {
		case ColumnSeparator:
			css := ""
			session := columnLayout.Session()
			if val, ok := columnLayout.properties[ColumnSeparator]; ok {
				separator := val.(ColumnSeparatorProperty)
				css = separator.cssValue(columnLayout.Session())
			}
			updateCSSProperty(columnLayout.htmlID(), "column-rule", css, session)

		case ColumnCount:
			session := columnLayout.Session()
			if count, ok := intProperty(columnLayout, tag, session, 0); ok && count > 0 {
				updateCSSProperty(columnLayout.htmlID(), tag, strconv.Itoa(count), session)
			} else {
				updateCSSProperty(columnLayout.htmlID(), tag, "auto", session)
			}
		}
	}
	return true
}

// GetColumnCount returns int value which specifies number of columns into which the content of
// ColumnLayout is break. If the return value is 0 then the number of columns is calculated
// based on the "column-width" property.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetColumnCount(view View, subviewID string) int {
	return intStyledProperty(view, subviewID, ColumnCount, 0)
}

// GetColumnWidth returns SizeUnit value which specifies the width of each column of ColumnLayout.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetColumnWidth(view View, subviewID string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ColumnWidth, false)
}

// GetColumnGap returns SizeUnit property which specifies the size of the gap (gutter) between columns of ColumnLayout.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetColumnGap(view View, subviewID string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ColumnGap, false)
}

// GetColumnSeparator returns ViewBorder struct which specifies the line drawn between
// columns in a multi-column ColumnLayout.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetColumnSeparator(view View, subviewID string) ViewBorder {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		value := view.Get(ColumnSeparator)
		if value == nil {
			value = valueFromStyle(view, ColumnSeparator)
		}

		if value != nil {
			if separator, ok := value.(ColumnSeparatorProperty); ok {
				return separator.ViewBorder(view.Session())
			}
		}
	}

	return ViewBorder{}
}

// ColumnSeparatorStyle returns int value which specifies the style of the line drawn between
// columns in a multi-column layout.
// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetColumnSeparatorStyle(view View, subviewID string) int {
	border := GetColumnSeparator(view, subviewID)
	return border.Style
}

// ColumnSeparatorWidth returns SizeUnit value which specifies the width of the line drawn between
// columns in a multi-column layout.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetColumnSeparatorWidth(view View, subviewID string) SizeUnit {
	border := GetColumnSeparator(view, subviewID)
	return border.Width
}

// ColumnSeparatorColor returns Color value which specifies the color of the line drawn between
// columns in a multi-column layout.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetColumnSeparatorColor(view View, subviewID string) Color {
	border := GetColumnSeparator(view, subviewID)
	return border.Color
}
