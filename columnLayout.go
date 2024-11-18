package rui

import (
	"strconv"
)

// Constants for [ColumnLayout] specific properties and events
const (
	// ColumnCount is the constant for "column-count" property tag.
	//
	// Used by `ColumnLayout`.
	// Specifies number of columns into which the content is break. Values less than zero are not valid. If this property
	// value is 0 then the number of columns is calculated based on the "column-width" property.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0` or "0" - Use "column-width" to control how many columns will be created.
	// >= `0` or >= "0" - Тhe number of columns into which the content is divided.
	ColumnCount PropertyName = "column-count"

	// ColumnWidth is the constant for "column-width" property tag.
	//
	// Used by `ColumnLayout`.
	// Specifies the width of each column.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	ColumnWidth PropertyName = "column-width"

	// ColumnGap is the constant for "column-gap" property tag.
	//
	// Used by `ColumnLayout`.
	// Set the size of the gap (gutter) between columns.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	ColumnGap PropertyName = "column-gap"

	// ColumnSeparator is the constant for "column-separator" property tag.
	//
	// Used by `ColumnLayout`.
	// Specifies the line drawn between columns in a multi-column layout.
	//
	// Supported types: `ColumnSeparatorProperty`, `ViewBorder`.
	//
	// Internal type is `ColumnSeparatorProperty`, other types converted to it during assignment.
	// See `ColumnSeparatorProperty` and `ViewBorder` description for more details.
	ColumnSeparator PropertyName = "column-separator"

	// ColumnSeparatorStyle is the constant for "column-separator-style" property tag.
	//
	// Used by `ColumnLayout`.
	// Controls the style of the line drawn between columns in a multi-column layout.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`NoneLine`) or "none" - The separator will not be drawn.
	// `1`(`SolidLine`) or "solid" - Solid line as a separator.
	// `2`(`DashedLine`) or "dashed" - Dashed line as a separator.
	// `3`(`DottedLine`) or "dotted" - Dotted line as a separator.
	// `4`(`DoubleLine`) or "double" - Double line as a separator.
	ColumnSeparatorStyle PropertyName = "column-separator-style"

	// ColumnSeparatorWidth is the constant for "column-separator-width" property tag.
	//
	// Used by `ColumnLayout`.
	// Set the width of the line drawn between columns in a multi-column layout.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	ColumnSeparatorWidth PropertyName = "column-separator-width"

	// ColumnSeparatorColor is the constant for "column-separator-color" property tag.
	//
	// Used by `ColumnLayout`.
	// Set the color of the line drawn between columns in a multi-column layout.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	ColumnSeparatorColor PropertyName = "column-separator-color"

	// ColumnFill is the constant for "column-fill" property tag.
	//
	// Used by `ColumnLayout`.
	// Controls how a `ColumnLayout`'s content is balanced when broken into columns. Default value is "balance".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`ColumnFillBalance`) or "balance" - Content is equally divided between columns.
	// `1`(`ColumnFillAuto`) or "auto" - Columns are filled sequentially. Content takes up only the room it needs, possibly resulting in some columns remaining empty.
	ColumnFill PropertyName = "column-fill"

	// ColumnSpanAll is the constant for "column-span-all" property tag.
	//
	// Used by `ColumnLayout`.
	// Property used in views placed inside the column layout container. Makes it possible for a view to span across all
	// columns. Default value is `false`.
	//
	// Supported types: `bool`, `int`, `string`.
	//
	// Values:
	// `true` or `1` or "true", "yes", "on", "1" - View will span across all columns.
	// `false` or `0` or "false", "no", "off", "0" - View will be a part of a column.
	ColumnSpanAll PropertyName = "column-span-all"
)

// ColumnLayout represent a ColumnLayout view
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
	return new(columnLayoutData)
}

// Init initialize fields of ColumnLayout by default values
func (columnLayout *columnLayoutData) init(session Session) {
	columnLayout.viewsContainerData.init(session)
	columnLayout.tag = "ColumnLayout"
	columnLayout.normalize = normalizeColumnLayoutTag
	columnLayout.changed = columnLayout.propertyChanged
	//columnLayout.systemClass = "ruiColumnLayout"
}

func normalizeColumnLayoutTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Gap:
		return ColumnGap
	}
	return tag
}

func (columnLayout *columnLayoutData) propertyChanged(tag PropertyName) {
	switch tag {
	case ColumnSeparator:
		css := ""
		session := columnLayout.Session()
		if value := columnLayout.getRaw(ColumnSeparator); value != nil {
			separator := value.(ColumnSeparatorProperty)
			css = separator.cssValue(session)
		}
		session.updateCSSProperty(columnLayout.htmlID(), "column-rule", css)

	case ColumnCount:
		session := columnLayout.Session()
		if count := GetColumnCount(columnLayout); count > 0 {
			session.updateCSSProperty(columnLayout.htmlID(), string(ColumnCount), strconv.Itoa(count))
		} else {
			session.updateCSSProperty(columnLayout.htmlID(), string(ColumnCount), "auto")
		}

	default:
		columnLayout.viewsContainerData.propertyChanged(tag)
	}
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
