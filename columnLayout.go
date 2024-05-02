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

	// ColumnFill is the constant for the "column-fill" property tag.
	// The "column-fill" int property controls how an ColumnLayout's contents are balanced when broken into columns.
	// Valid values are
	// * ColumnFillBalance (0) - Content is equally divided between columns (default value);
	// * ColumnFillAuto (1) - Columns are filled sequentially. Content takes up only the room it needs, possibly resulting in some columns remaining empty.
	ColumnFill = "column-fill"

	// ColumnSpanAll is the constant for the "column-span-all" property tag.
	// The "column-span-all" bool property makes it possible for a view to span across all columns when its value is set to true.
	ColumnSpanAll = "column-span-all"
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
	view.init(session)
	setInitParams(view, params)
	return view
}

func newColumnLayout(session Session) View {
	return NewColumnLayout(session, nil)
}

// Init initialize fields of ColumnLayout by default values
func (ColumnLayout *columnLayoutData) init(session Session) {
	ColumnLayout.viewsContainerData.init(session)
	ColumnLayout.tag = "ColumnLayout"
	//ColumnLayout.systemClass = "ruiColumnLayout"
}

func (columnLayout *columnLayoutData) String() string {
	return getViewString(columnLayout, nil)
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
			columnLayout.session.updateCSSProperty(columnLayout.htmlID(), tag, "")

		case ColumnSeparator:
			columnLayout.session.updateCSSProperty(columnLayout.htmlID(), "column-rule", "")
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
			session.updateCSSProperty(columnLayout.htmlID(), "column-rule", css)

		case ColumnCount:
			session := columnLayout.Session()
			if count, ok := intProperty(columnLayout, tag, session, 0); ok && count > 0 {
				session.updateCSSProperty(columnLayout.htmlID(), tag, strconv.Itoa(count))
			} else {
				session.updateCSSProperty(columnLayout.htmlID(), tag, "auto")
			}
		}
	}
	return true
}

// GetColumnCount returns int value which specifies number of columns into which the content of
// ColumnLayout is break. If the return value is 0 then the number of columns is calculated
// based on the "column-width" property.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetColumnCount(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, ColumnCount, 0)
}

// GetColumnWidth returns SizeUnit value which specifies the width of each column of ColumnLayout.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetColumnWidth(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ColumnWidth, false)
}

// GetColumnGap returns SizeUnit property which specifies the size of the gap (gutter) between columns of ColumnLayout.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetColumnGap(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, ColumnGap, false)
}

func getColumnSeparator(view View, subviewID []string) ViewBorder {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
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

// GetColumnSeparator returns ViewBorder struct which specifies the line drawn between
// columns in a multi-column ColumnLayout.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetColumnSeparator(view View, subviewID ...string) ViewBorder {
	return getColumnSeparator(view, subviewID)
}

// ColumnSeparatorStyle returns int value which specifies the style of the line drawn between
// columns in a multi-column layout.
// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetColumnSeparatorStyle(view View, subviewID ...string) int {
	border := getColumnSeparator(view, subviewID)
	return border.Style
}

// ColumnSeparatorWidth returns SizeUnit value which specifies the width of the line drawn between
// columns in a multi-column layout.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetColumnSeparatorWidth(view View, subviewID ...string) SizeUnit {
	border := getColumnSeparator(view, subviewID)
	return border.Width
}

// ColumnSeparatorColor returns Color value which specifies the color of the line drawn between
// columns in a multi-column layout.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetColumnSeparatorColor(view View, subviewID ...string) Color {
	border := getColumnSeparator(view, subviewID)
	return border.Color
}

// GetColumnFill returns a "column-fill" property value of the subview.
// Returns one of next values: ColumnFillBalance (0) or ColumnFillAuto (1)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetColumnFill(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, ColumnFill, ColumnFillBalance, true)
}

// IsColumnSpanAll returns a "column-span-all" property value of the subview.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsColumnSpanAll(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, ColumnSpanAll, false)
}
