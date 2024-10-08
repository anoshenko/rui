package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// Constants for [TableView] specific properties and events
const (
	// TableVerticalAlign is the constant for "table-vertical-align" property tag.
	//
	// Used by `TableView`.
	// Set the vertical alignment of the content inside a table cell.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`TopAlign`) or "top" - Top alignment.
	// `1`(`BottomAlign`) or "bottom" - Bottom alignment.
	// `2`(`CenterAlign`) or "center" - Center alignment.
	// `3`(`StretchAlign`) or "stretch" - Work as baseline alignment, see below.
	// `4`(`BaselineAlign`) or "baseline" - Baseline alignment.
	TableVerticalAlign = "table-vertical-align"

	// HeadHeight is the constant for "head-height" property tag.
	//
	// Used by `TableView`.
	// Sets the number of rows in the table header. The default value is `0` (no header).
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0` or "0" - No header.
	// > `0` or > "0" - Number of rows act as a header.
	HeadHeight = "head-height"

	// HeadStyle is the constant for "head-style" property tag.
	//
	// Used by `TableView`.
	// Set the header style name or description of style properties.
	//
	// Supported types: `string`, `Params`.
	//
	// Internal type is either `string` or `Params`.
	//
	// Conversion rules:
	// `string` - must contain style name defined in resources.
	// `Params` - must contain style properties.
	HeadStyle = "head-style"

	// FootHeight is the constant for "foot-height" property tag.
	//
	// Used by `TableView`.
	// Sets the number of rows in the table footer. The default value is `0` (no footer).
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0` or "0" - No footer.
	// > `0` or > "0" - Number of rows act as a footer.
	FootHeight = "foot-height"

	// FootStyle is the constant for "foot-style" property tag.
	//
	// Used by `TableView`.
	// Set the footer style name or description of style properties.
	//
	// Supported types: `string`, `Params`.
	//
	// Internal type is either `string` or `Params`.
	//
	// Conversion rules:
	// `string` - must contain style name defined in resources.
	// `Params` - must contain style properties.
	FootStyle = "foot-style"

	// RowSpan is the constant for "row-span" property tag.
	//
	// Used by `TableView`.
	// Set the number of table row to span. Used only when specifying cell parameters in the implementation of 
	// `TableCellStyle`.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0` or "0" - No merging will be applied.
	// > `0` or > "0" - Number of rows including current one to be merged together.
	RowSpan = "row-span"

	// ColumnSpan is the constant for "column-span" property tag.
	//
	// Used by `TableView`.
	// Sets the number of table column cells to be merged together. Used only when specifying cell parameters in the 
	// implementation of `TableCellStyle`.
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0` or "0" - No merging will be applied.
	// > `0` or > "0" - Number of columns including current one to be merged together.
	ColumnSpan = "column-span"

	// RowStyle is the constant for "row-style" property tag.
	//
	// Used by `TableView`.
	// Set the adapter which specifies styles of each table row.
	//
	// Supported types: `TableRowStyle`, `[]Params`.
	//
	// Internal type is `TableRowStyle`, other types converted to it during assignment.
	// See `TableRowStyle` description for more details.
	RowStyle = "row-style"

	// ColumnStyle is the constant for "column-style" property tag.
	//
	// Used by `TableView`.
	// Set the adapter which specifies styles of each table column.
	//
	// Supported types: `TableColumnStyle`, `[]Params`.
	//
	// Internal type is `TableColumnStyle`, other types converted to it during assignment.
	// See `TableColumnStyle` description for more details.
	ColumnStyle = "column-style"

	// CellStyle is the constant for "cell-style" property tag.
	//
	// Used by `TableView`.
	// Set the adapter which specifies styles of each table cell. This property can be assigned only by an implementation of 
	// `TableCellStyle` interface.
	//
	// Supported types: `TableCellStyle`.
	CellStyle = "cell-style"

	// CellPadding is the constant for "cell-padding" property tag.
	//
	// Used by `TableView`.
	// Sets the padding area on all four sides of a table cell at once. An element's padding area is the space between its 
	// content and its border.
	//
	// Supported types: `BoundsProperty`, `Bounds`, `SizeUnit`, `float32`, `float64`, `int`.
	//
	// Internal type is `BoundsProperty`, other types converted to it during assignment.
	// See `BoundsProperty`, `Bounds` and `SizeUnit` description for more details.
	CellPadding = "cell-padding"

	// CellPaddingLeft is the constant for "cell-padding-left" property tag.
	//
	// Used by `TableView`.
	// Set the width of the padding area to the left of a cell content. An element's padding area is the space between its 
	// content and its border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellPaddingLeft = "cell-padding-left"

	// CellPaddingRight is the constant for "cell-padding-right" property tag.
	//
	// Used by `TableView`.
	// Set the width of the padding area to the left of a cell content. An element's padding area is the space between its 
	// content and its border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellPaddingRight = "cell-padding-right"

	// CellPaddingTop is the constant for "cell-padding-top" property tag.
	//
	// Used by `TableView`.
	// Set the height of the padding area to the top of a cell content. An element's padding area is the space between its 
	// content and its border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellPaddingTop = "cell-padding-top"

	// CellPaddingBottom is the constant for "cell-padding-bottom" property tag.
	//
	// Used by `TableView`.
	// Set the height of the padding area to the bottom of a cell content.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellPaddingBottom = "cell-padding-bottom"

	// CellBorder is the constant for "cell-border" property tag.
	//
	// Used by `TableView`.
	// Set a table cell's border. It sets the values of a border width, style, and color. Can also be used when setting 
	// parameters in properties "row-style", "column-style", "foot-style" and "head-style".
	//
	// Supported types: `BorderProperty`, `ViewBorder`, `ViewBorders`.
	//
	// Internal type is `BorderProperty`, other types converted to it during assignment.
	// See `BorderProperty`, `ViewBorder` and `ViewBorders` description for more details.
	CellBorder = "cell-border"

	// CellBorderLeft is the constant for "cell-border-left" property tag.
	//
	// Used by `TableView`.
	// Set a view's left border. It sets the values of a border width, style, and color. This property can be assigned a value 
	// of `BorderProperty`, `ViewBorder` types or `BorderProperty` text representation.
	//
	// Supported types: `ViewBorder`, `BorderProperty`, `string`.
	//
	// Internal type is `BorderProperty`, other types converted to it during assignment.
	// See `ViewBorder` and `BorderProperty` description for more details.
	CellBorderLeft = "cell-border-left"

	// CellBorderRight is the constant for "cell-border-right" property tag.
	//
	// Used by `TableView`.
	// Set a view's right border. It sets the values of a border width, style, and color. This property can be assigned a 
	// value of `BorderProperty`, `ViewBorder` types or `BorderProperty` text representation.
	//
	// Supported types: `ViewBorder`, `BorderProperty`, `string`.
	//
	// Internal type is `BorderProperty`, other types converted to it during assignment.
	// See `ViewBorder` and `BorderProperty` description for more details.
	CellBorderRight = "cell-border-right"

	// CellBorderTop is the constant for "cell-border-top" property tag.
	//
	// Used by `TableView`.
	// Set a view's top border. It sets the values of a border width, style, and color. This property can be assigned a value 
	// of `BorderProperty`, `ViewBorder` types or `BorderProperty` text representation.
	//
	// Supported types: `ViewBorder`, `BorderProperty`, `string`.
	//
	// Internal type is `BorderProperty`, other types converted to it during assignment.
	// See `ViewBorder` and `BorderProperty` description for more details.
	CellBorderTop = "cell-border-top"

	// CellBorderBottom is the constant for "cell-border-bottom" property tag.
	//
	// Used by `TableView`.
	// Set a view's bottom border. It sets the values of a border width, style, and color.
	//
	// Supported types: `ViewBorder`, `BorderProperty`, `string`.
	//
	// Internal type is `BorderProperty`, other types converted to it during assignment.
	// See `ViewBorder` and `BorderProperty` description for more details.
	CellBorderBottom = "cell-border-bottom"

	// CellBorderStyle is the constant for "cell-border-style" property tag.
	//
	// Used by `TableView`.
	// Set the line style for all four sides of a table cell's border. Default value is "none".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`NoneLine`) or "none" - The border will not be drawn.
	// `1`(`SolidLine`) or "solid" - Solid line as a border.
	// `2`(`DashedLine`) or "dashed" - Dashed line as a border.
	// `3`(`DottedLine`) or "dotted" - Dotted line as a border.
	// `4`(`DoubleLine`) or "double" - Double line as a border.
	CellBorderStyle = "cell-border-style"

	// CellBorderLeftStyle is the constant for "cell-border-left-style" property tag.
	//
	// Used by `TableView`.
	// Set the line style of a table cell's left border. Default value is "none".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`NoneLine`) or "none" - The border will not be drawn.
	// `1`(`SolidLine`) or "solid" - Solid line as a border.
	// `2`(`DashedLine`) or "dashed" - Dashed line as a border.
	// `3`(`DottedLine`) or "dotted" - Dotted line as a border.
	// `4`(`DoubleLine`) or "double" - Double line as a border.
	CellBorderLeftStyle = "cell-border-left-style"

	// CellBorderRightStyle is the constant for "cell-border-right-style" property tag.
	//
	// Used by `TableView`.
	// Set the line style of a table cell's right border. Default value is "none".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`NoneLine`) or "none" - The border will not be drawn.
	// `1`(`SolidLine`) or "solid" - Solid line as a border.
	// `2`(`DashedLine`) or "dashed" - Dashed line as a border.
	// `3`(`DottedLine`) or "dotted" - Dotted line as a border.
	// `4`(`DoubleLine`) or "double" - Double line as a border.
	CellBorderRightStyle = "cell-border-right-style"

	// CellBorderTopStyle is the constant for "cell-border-top-style" property tag.
	//
	// Used by `TableView`.
	// Set the line style of a table cell's top border. Default value is "none".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`NoneLine`) or "none" - The border will not be drawn.
	// `1`(`SolidLine`) or "solid" - Solid line as a border.
	// `2`(`DashedLine`) or "dashed" - Dashed line as a border.
	// `3`(`DottedLine`) or "dotted" - Dotted line as a border.
	// `4`(`DoubleLine`) or "double" - Double line as a border.
	CellBorderTopStyle = "cell-border-top-style"

	// CellBorderBottomStyle is the constant for "cell-border-bottom-style" property tag.
	//
	// Used by `TableView`.
	// Sets the line style of a table cell's bottom border. Default value is "none".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`NoneLine`) or "none" - The border will not be drawn.
	// `1`(`SolidLine`) or "solid" - Solid line as a border.
	// `2`(`DashedLine`) or "dashed" - Dashed line as a border.
	// `3`(`DottedLine`) or "dotted" - Dotted line as a border.
	// `4`(`DoubleLine`) or "double" - Double line as a border.
	CellBorderBottomStyle = "cell-border-bottom-style"

	// CellBorderWidth is the constant for "cell-border-width" property tag.
	//
	// Used by `TableView`.
	// Set the line width for all four sides of a table cell's border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellBorderWidth = "cell-border-width"

	// CellBorderLeftWidth is the constant for "cell-border-left-width" property tag.
	//
	// Used by `TableView`.
	// Set the line width of a table cell's left border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellBorderLeftWidth = "cell-border-left-width"

	// CellBorderRightWidth is the constant for "cell-border-right-width" property tag.
	//
	// Used by `TableView`.
	// Set the line width of a table cell's right border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellBorderRightWidth = "cell-border-right-width"

	// CellBorderTopWidth is the constant for "cell-border-top-width" property tag.
	//
	// Used by `TableView`.
	// Set the line width of a table cell's top border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellBorderTopWidth = "cell-border-top-width"

	// CellBorderBottomWidth is the constant for "cell-border-bottom-width" property tag.
	//
	// Used by `TableView`.
	// Set the line width of a table cell's bottom border.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	CellBorderBottomWidth = "cell-border-bottom-width"

	// CellBorderColor is the constant for "cell-border-color" property tag.
	//
	// Used by `TableView`.
	// Set the line color for all four sides of a table cell's border.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	CellBorderColor = "cell-border-color"

	// CellBorderLeftColor is the constant for "cell-border-left-color" property tag.
	//
	// Used by `TableView`.
	// Set the line color of a table cell's left border.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	CellBorderLeftColor = "cell-border-left-color"

	// CellBorderRightColor is the constant for "cell-border-right-color" property tag.
	//
	// Used by `TableView`.
	// Set the line color of a table cell's right border.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	CellBorderRightColor = "cell-border-right-color"

	// CellBorderTopColor is the constant for "cell-border-top-color" property tag.
	//
	// Used by `TableView`.
	// Set the line color of a table cell's top border.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	CellBorderTopColor = "cell-border-top-color"

	// CellBorderBottomColor is the constant for "cell-border-bottom-color" property tag.
	//
	// Used by `TableView`.
	// Set the line color of a table cell's bottom border.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	CellBorderBottomColor = "cell-border-bottom-color"

	// SelectionMode is the constant for "selection-mode" property tag.
	//
	// Used by `TableView`.
	// Sets the mode of the table elements selection. Default value is "none".
	//
	// Supported types: `int`, `string`.
	//
	// Values:
	// `0`(`NoneSelection`) or "none" - Table elements are not selectable. The table cannot receive input focus.
	// `1`(`CellSelection`) or "cell" - One table cell can be selected(highlighted). The cell is selected interactively using the mouse or keyboard(using the cursor keys).
	// `2`(`RowSelection`) or "row" - The entire table row can be selected (highlighted). The row is selected interactively using the mouse or keyboard (using the cursor keys).
	SelectionMode = "selection-mode"

	// TableCellClickedEvent is the constant for "table-cell-clicked" property tag.
	//
	// Used by `TableView`.
	// Occur when the user clicks on a table cell.
	//
	// General listener format:
	// `func(table rui.TableView, row, col int)`.
	//
	// where:
	// table - Interface of a table view which generated this event,
	// row - Row of the clicked cell,
	// col - Column of the clicked cell.
	//
	// Allowed listener formats:
	// `func(row, col int)`.
	TableCellClickedEvent = "table-cell-clicked"

	// TableCellSelectedEvent is the constant for "table-cell-selected" property tag.
	//
	// Used by `TableView`.
	// Occur when a table cell becomes selected.
	//
	// General listener format:
	// `func(table rui.TableView, row, col int)`.
	//
	// where:
	// table - Interface of a table view which generated this event,
	// row - Row of the selected cell,
	// col - Column of the selected cell.
	//
	// Allowed listener formats:
	// `func(row, col int)`.
	TableCellSelectedEvent = "table-cell-selected"

	// TableRowClickedEvent is the constant for "table-row-clicked" property tag.
	//
	// Used by `TableView`.
	// Occur when the user clicks on a table row.
	//
	// General listener format:
	// `func(table rui.TableView, row int)`.
	//
	// where:
	// table - Interface of a table view which generated this event,
	// row - Clicked row.
	//
	// Allowed listener formats:
	// `func(row int)`.
	TableRowClickedEvent = "table-row-clicked"

	// TableRowSelectedEvent is the constant for "table-row-selected" property tag.
	//
	// Used by `TableView`.
	// Occur when a table row becomes selected.
	//
	// General listener format:
	// `func(table rui.TableView, row int)`.
	//
	// where:
	// table - Interface of a table view which generated this event,
	// row - Selected row.
	//
	// Allowed listener formats:
	// `func(row int)`.
	TableRowSelectedEvent = "table-row-selected"

	// AllowSelection is the constant for "allow-selection" property tag.
	//
	// Used by `TableView`.
	// Set the adapter which specifies whether cell/row selection is allowed. This property can be assigned by an 
	// implementation of `TableAllowCellSelection` or `TableAllowRowSelection` interface.
	//
	// Supported types: `TableAllowCellSelection`, `TableAllowRowSelection`.
	//
	// Internal type is either `TableAllowCellSelection`, `TableAllowRowSelection`, see their description for more details.
	AllowSelection = "allow-selection"
)

// Constants which represent values of "selection-mode" property of a [TableView]
const (
	// NoneSelection the selection is forbidden.
	NoneSelection = 0
	// CellSelection the selection of a single cell only is enabled.
	CellSelection = 1
	// RowSelection the selection of a table row only is enabled.
	RowSelection = 2
)

// CellIndex defines coordinates of the [TableView] cell
type CellIndex struct {
	Row, Column int
}

// TableView represents a TableView view
type TableView interface {
	View
	ParentView

	// ReloadTableData forces the table view to reload all data and redraw the entire table
	ReloadTableData()

	// ReloadCell forces the table view to reload the data for a specific cell and redraw it
	ReloadCell(row, column int)

	// CellFrame returns the frame of a specific cell, describing its position and size within the table view
	CellFrame(row, column int) Frame

	content() TableAdapter
	getCurrent() CellIndex
	getRowStyle() TableRowStyle
	getColumnStyle() TableColumnStyle
	getCellStyle() TableCellStyle
}

type tableViewData struct {
	viewData
	cellViews                                 []View
	cellFrame                                 []Frame
	cellSelectedListener, cellClickedListener []func(TableView, int, int)
	rowSelectedListener, rowClickedListener   []func(TableView, int)
	current                                   CellIndex
}

type tableCellView struct {
	viewData
}

// NewTableView create new TableView object and return it
func NewTableView(session Session, params Params) TableView {
	view := new(tableViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newTableView(session Session) View {
	return NewTableView(session, nil)
}

// Init initialize fields of TableView by default values
func (table *tableViewData) init(session Session) {
	table.viewData.init(session)
	table.tag = "TableView"
	table.cellViews = []View{}
	table.cellFrame = []Frame{}
	table.cellSelectedListener = []func(TableView, int, int){}
	table.cellClickedListener = []func(TableView, int, int){}
	table.rowSelectedListener = []func(TableView, int){}
	table.rowClickedListener = []func(TableView, int){}
	table.current.Row = -1
	table.current.Column = -1
}

func (table *tableViewData) String() string {
	return getViewString(table, nil)
}

func (table *tableViewData) normalizeTag(tag string) string {
	switch tag = strings.ToLower(tag); tag {
	case "top-cell-padding":
		tag = CellPaddingTop

	case "right-cell-padding":
		tag = CellPaddingRight

	case "bottom-cell-padding":
		tag = CellPaddingBottom

	case "left-cell-padding":
		tag = CellPaddingLeft
	}
	return tag
}

func (table *tableViewData) Focusable() bool {
	return GetTableSelectionMode(table) != NoneSelection
}

func (table *tableViewData) Get(tag string) any {
	return table.get(table.normalizeTag(tag))
}

func (table *tableViewData) Remove(tag string) {
	table.remove(table.normalizeTag(tag))
}

func (table *tableViewData) remove(tag string) {
	switch tag {
	case CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft:
		table.removeBoundsSide(CellPadding, tag)
		table.propertyChanged(tag)

	case SelectionMode, TableVerticalAlign, Gap, CellBorder, CellPadding, RowStyle,
		ColumnStyle, CellStyle, HeadHeight, HeadStyle, FootHeight, FootStyle, AllowSelection:
		if _, ok := table.properties[tag]; ok {
			delete(table.properties, tag)
			table.propertyChanged(tag)
		}

	case TableCellClickedEvent:
		table.cellClickedListener = []func(TableView, int, int){}
		table.propertyChanged(tag)

	case TableCellSelectedEvent:
		table.cellSelectedListener = []func(TableView, int, int){}
		table.propertyChanged(tag)

	case TableRowClickedEvent:
		table.rowClickedListener = []func(TableView, int){}
		table.propertyChanged(tag)

	case TableRowSelectedEvent:
		table.rowSelectedListener = []func(TableView, int){}
		table.propertyChanged(tag)

	case Current:
		table.current.Row = -1
		table.current.Column = -1
		table.propertyChanged(tag)

	default:
		table.viewData.remove(tag)
	}
}

func (table *tableViewData) Set(tag string, value any) bool {
	return table.set(table.normalizeTag(tag), value)
}

func (table *tableViewData) set(tag string, value any) bool {
	if value == nil {
		table.remove(tag)
		return true
	}

	switch tag {
	case Content:
		switch val := value.(type) {
		case TableAdapter:
			table.properties[Content] = value

		case [][]any:
			table.properties[Content] = NewSimpleTableAdapter(val)

		case [][]string:
			table.properties[Content] = NewTextTableAdapter(val)

		default:
			notCompatibleType(tag, value)
			return false
		}

	case TableCellClickedEvent:
		listeners := table.valueToCellListeners(value)
		if listeners == nil {
			notCompatibleType(tag, value)
			return false
		}
		table.cellClickedListener = listeners

	case TableCellSelectedEvent:
		listeners := table.valueToCellListeners(value)
		if listeners == nil {
			notCompatibleType(tag, value)
			return false
		}
		table.cellSelectedListener = listeners

	case TableRowClickedEvent:
		listeners, ok := valueToEventListeners[TableView, int](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(TableView, int){}
		}
		table.rowClickedListener = listeners

	case TableRowSelectedEvent:
		listeners, ok := valueToEventListeners[TableView, int](value)
		if !ok {
			notCompatibleType(tag, value)
			return false
		} else if listeners == nil {
			listeners = []func(TableView, int){}
		}
		table.rowSelectedListener = listeners

	case CellStyle:
		if style, ok := value.(TableCellStyle); ok {
			table.properties[tag] = style
		} else {
			notCompatibleType(tag, value)
			return false
		}

	case RowStyle:
		if !table.setRowStyle(value) {
			notCompatibleType(tag, value)
			return false
		}

	case ColumnStyle:
		if !table.setColumnStyle(value) {
			notCompatibleType(tag, value)
			return false
		}

	case HeadHeight, FootHeight:
		switch value := value.(type) {
		case string:
			if isConstantName(value) {
				table.properties[tag] = value
			} else if n, err := strconv.Atoi(value); err == nil {
				table.properties[tag] = n
			} else {
				ErrorLog(err.Error())
				notCompatibleType(tag, value)
				return false
			}

		default:
			if n, ok := isInt(value); ok {
				table.properties[tag] = n
			}
		}

	case HeadStyle, FootStyle:
		switch value := value.(type) {
		case string:
			table.properties[tag] = value

		case Params:
			if len(value) > 0 {
				table.properties[tag] = value
			} else {
				delete(table.properties, tag)
			}

		case DataObject:
			params := Params{}
			for k := 0; k < value.PropertyCount(); k++ {
				if prop := value.Property(k); prop != nil && prop.Type() == TextNode {
					params[prop.Tag()] = prop.Text()
				}
			}
			if len(params) > 0 {
				table.properties[tag] = params
			} else {
				delete(table.properties, tag)
			}

		case DataNode:
			switch value.Type() {
			case ObjectNode:
				return table.set(tag, value.Object())

			case TextNode:
				table.properties[tag] = value.Text()

			default:
				notCompatibleType(tag, value)
				return false
			}

		default:
			notCompatibleType(tag, value)
			return false
		}

	case CellPadding:
		if !table.setBounds(tag, value) {
			return false
		}

	case CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft:
		if !table.setBoundsSide(CellPadding, tag, value) {
			return false
		}

	case Gap:
		if !table.setSizeProperty(Gap, value) {
			return false
		}

	case SelectionMode, TableVerticalAlign, CellBorder, CellBorderStyle, CellBorderColor, CellBorderWidth,
		CellBorderLeft, CellBorderLeftStyle, CellBorderLeftColor, CellBorderLeftWidth,
		CellBorderRight, CellBorderRightStyle, CellBorderRightColor, CellBorderRightWidth,
		CellBorderTop, CellBorderTopStyle, CellBorderTopColor, CellBorderTopWidth,
		CellBorderBottom, CellBorderBottomStyle, CellBorderBottomColor, CellBorderBottomWidth:
		if !table.viewData.set(tag, value) {
			return false
		}

	case AllowSelection:
		switch value.(type) {
		case TableAllowCellSelection:
			table.properties[tag] = value

		case TableAllowRowSelection:
			table.properties[tag] = value

		default:
			notCompatibleType(tag, value)
			return false
		}

	case Current:
		switch value := value.(type) {
		case CellIndex:
			table.current = value

		case DataObject:
			if row, ok := dataIntProperty(value, "row"); ok {
				table.current.Row = row
			}
			if column, ok := dataIntProperty(value, "column"); ok {
				table.current.Column = column
			}

		case string:
			if strings.Contains(value, ",") {
				if values := strings.Split(value, ","); len(values) == 2 {
					var n = []int{0, 0}
					for i := 0; i < 2; i++ {
						var err error
						if n[i], err = strconv.Atoi(values[i]); err != nil {
							ErrorLog(err.Error())
							return false
						}
					}
					table.current.Row = n[0]
					table.current.Column = n[1]
				} else {
					notCompatibleType(tag, value)
					return false
				}
			} else {
				n, err := strconv.Atoi(value)
				if err != nil {
					ErrorLog(err.Error())
					return false
				}
				table.current.Row = n
				table.current.Column = -1
			}

		default:
			if n, ok := isInt(value); ok {
				table.current.Row = n
				table.current.Column = -1
			} else {
				notCompatibleType(tag, value)
				return false
			}
		}

	default:
		return table.viewData.set(tag, value)
	}

	table.propertyChanged(tag)
	return true
}

func (table *tableViewData) propertyChanged(tag string) {
	if table.created {
		switch tag {
		case Content, TableVerticalAlign, RowStyle, ColumnStyle, CellStyle, CellPadding,
			CellBorder, HeadHeight, HeadStyle, FootHeight, FootStyle,
			CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft,
			TableCellClickedEvent, TableCellSelectedEvent, TableRowClickedEvent,
			TableRowSelectedEvent, AllowSelection, AccentColor:
			table.ReloadTableData()

		case Current:
			switch GetTableSelectionMode(table) {
			case CellSelection:
				table.session.callFunc("setTableCellCursorByID", table.htmlID(), table.current.Row, table.current.Column)

			case RowSelection:
				table.session.callFunc("setTableRowCursorByID", table.htmlID(), table.current.Row)
			}

		case Gap:
			htmlID := table.htmlID()
			session := table.Session()
			gap, ok := sizeProperty(table, Gap, session)
			if !ok || gap.Type == Auto || gap.Value <= 0 {
				session.updateCSSProperty(htmlID, "border-spacing", "0")
				session.updateCSSProperty(htmlID, "border-collapse", "collapse")
			} else {
				session.updateCSSProperty(htmlID, "border-spacing", gap.cssString("0", session))
				session.updateCSSProperty(htmlID, "border-collapse", "separate")
			}

		case SelectionMode:
			htmlID := table.htmlID()
			session := table.Session()

			switch GetTableSelectionMode(table) {
			case CellSelection:
				tabIndex, _ := intProperty(table, TabIndex, session, 0)
				session.updateProperty(htmlID, "tabindex", tabIndex)
				session.updateProperty(htmlID, "onfocus", "tableViewFocusEvent(this, event)")
				session.updateProperty(htmlID, "onblur", "tableViewBlurEvent(this, event)")
				session.updateProperty(htmlID, "data-selection", "cell")
				session.updateProperty(htmlID, "data-focusitemstyle", table.currentStyle())
				session.updateProperty(htmlID, "data-bluritemstyle", table.currentInactiveStyle())

				if table.current.Row >= 0 && table.current.Column >= 0 {
					session.updateProperty(htmlID, "data-current", table.cellID(table.current.Row, table.current.Column))
				} else {
					session.removeProperty(htmlID, "data-current")
				}
				session.updateProperty(htmlID, "onkeydown", "tableViewCellKeyDownEvent(this, event)")

			case RowSelection:
				tabIndex, _ := intProperty(table, TabIndex, session, 0)
				session.updateProperty(htmlID, "tabindex", tabIndex)
				session.updateProperty(htmlID, "onfocus", "tableViewFocusEvent(this, event)")
				session.updateProperty(htmlID, "onblur", "tableViewBlurEvent(this, event)")
				session.updateProperty(htmlID, "data-selection", "row")
				session.updateProperty(htmlID, "data-focusitemstyle", table.currentStyle())
				session.updateProperty(htmlID, "data-bluritemstyle", table.currentInactiveStyle())

				if table.current.Row >= 0 {
					session.updateProperty(htmlID, "data-current", table.rowID(table.current.Row))
				} else {
					session.removeProperty(htmlID, "data-current")
				}
				session.updateProperty(htmlID, "onkeydown", "tableViewRowKeyDownEvent(this, event)")

			default: // NoneSelection
				if tabIndex, ok := intProperty(table, TabIndex, session, -1); !ok || tabIndex < 0 {
					session.removeProperty(htmlID, "tabindex")
				}

				for _, prop := range []string{"data-current", "onfocus", "onblur", "onkeydown", "data-selection"} {
					session.removeProperty(htmlID, prop)
				}
			}
			updateInnerHTML(htmlID, session)
		}
	}
	table.propertyChangedEvent(tag)
}

func (table *tableViewData) currentStyle() string {
	if value := table.getRaw(CurrentStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = table.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	if value := valueFromStyle(table, CurrentStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = table.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	return "ruiCurrentTableCellFocused"
}

func (table *tableViewData) currentInactiveStyle() string {
	if value := table.getRaw(CurrentInactiveStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = table.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	if value := valueFromStyle(table, CurrentInactiveStyle); value != nil {
		if style, ok := value.(string); ok {
			if style, ok = table.session.resolveConstants(style); ok {
				return style
			}
		}
	}
	return "ruiCurrentTableCell"
}

func (table *tableViewData) valueToCellListeners(value any) []func(TableView, int, int) {
	if value == nil {
		return []func(TableView, int, int){}
	}

	switch value := value.(type) {
	case func(TableView, int, int):
		return []func(TableView, int, int){value}

	case func(int, int):
		fn := func(_ TableView, row, column int) {
			value(row, column)
		}
		return []func(TableView, int, int){fn}

	case []func(TableView, int, int):
		return value

	case []func(int, int):
		listeners := make([]func(TableView, int, int), len(value))
		for i, val := range value {
			if val == nil {
				return nil
			}
			listeners[i] = func(_ TableView, row, column int) {
				val(row, column)
			}
		}
		return listeners

	case []any:
		listeners := make([]func(TableView, int, int), len(value))
		for i, val := range value {
			if val == nil {
				return nil
			}
			switch val := val.(type) {
			case func(TableView, int, int):
				listeners[i] = val

			case func(int, int):
				listeners[i] = func(_ TableView, row, column int) {
					val(row, column)
				}

			default:
				return nil
			}
		}
		return listeners
	}

	return nil
}

func (table *tableViewData) htmlTag() string {
	return "table"
}

func (table *tableViewData) rowID(index int) string {
	return fmt.Sprintf("%s-%d", table.htmlID(), index)
}

func (table *tableViewData) cellID(row, column int) string {
	return fmt.Sprintf("%s-%d-%d", table.htmlID(), row, column)
}

func (table *tableViewData) htmlProperties(self View, buffer *strings.Builder) {

	if content := table.content(); content != nil {
		buffer.WriteString(` data-rows="`)
		buffer.WriteString(strconv.Itoa(content.RowCount()))
		buffer.WriteString(`" data-columns="`)
		buffer.WriteString(strconv.Itoa(content.ColumnCount()))
		buffer.WriteRune('"')
	}

	if selectionMode := GetTableSelectionMode(table); selectionMode != NoneSelection {
		buffer.WriteString(` onfocus="tableViewFocusEvent(this, event)" onblur="tableViewBlurEvent(this, event)" data-focusitemstyle="`)
		buffer.WriteString(table.currentStyle())
		buffer.WriteString(`" data-bluritemstyle="`)
		buffer.WriteString(table.currentInactiveStyle())
		buffer.WriteRune('"')

		switch selectionMode {
		case RowSelection:
			buffer.WriteString(` data-selection="row" onkeydown="tableViewRowKeyDownEvent(this, event)"`)
			if table.current.Row >= 0 {
				buffer.WriteString(` data-current="`)
				buffer.WriteString(table.rowID(table.current.Row))
				buffer.WriteRune('"')
			}

		case CellSelection:
			buffer.WriteString(` data-selection="cell" onkeydown="tableViewCellKeyDownEvent(this, event)"`)
			if table.current.Row >= 0 && table.current.Column >= 0 {
				buffer.WriteString(` data-current="`)
				buffer.WriteString(table.cellID(table.current.Row, table.current.Column))
				buffer.WriteRune('"')
			}
		}
	}

	table.viewData.htmlProperties(self, buffer)
}

func (table *tableViewData) content() TableAdapter {
	if content := table.getRaw(Content); content != nil {
		if adapter, ok := content.(TableAdapter); ok {
			return adapter
		}
	}

	return nil
}

func (table *tableViewData) htmlSubviews(self View, buffer *strings.Builder) {
	table.cellViews = []View{}
	table.cellFrame = []Frame{}

	adapter := table.content()
	if adapter == nil {
		return
	}

	rowCount := adapter.RowCount()
	columnCount := adapter.ColumnCount()
	if rowCount == 0 || columnCount == 0 {
		return
	}

	table.cellFrame = make([]Frame, rowCount*columnCount)

	rowStyle := table.getRowStyle()
	cellStyle := table.getCellStyle()

	session := table.Session()

	if !session.ignoreViewUpdates() {
		session.setIgnoreViewUpdates(true)
		defer session.setIgnoreViewUpdates(false)
	}

	cssBuilder := viewCSSBuilder{buffer: allocStringBuilder()}
	defer freeStringBuilder(cssBuilder.buffer)

	var view tableCellView
	view.init(session)

	ignoreCells := []struct{ row, column int }{}
	selectionMode := GetTableSelectionMode(table)

	var allowCellSelection TableAllowCellSelection = nil
	if allow, ok := adapter.(TableAllowCellSelection); ok {
		allowCellSelection = allow
	}
	if value := table.getRaw(AllowSelection); value != nil {
		if style, ok := value.(TableAllowCellSelection); ok {
			allowCellSelection = style
		}
	}

	var allowRowSelection TableAllowRowSelection = nil
	if allow, ok := adapter.(TableAllowRowSelection); ok {
		allowRowSelection = allow
	}
	if value := table.getRaw(AllowSelection); value != nil {
		if style, ok := value.(TableAllowRowSelection); ok {
			allowRowSelection = style
		}
	}

	vAlignCss := enumProperties[TableVerticalAlign].cssValues
	vAlignValue := GetTableVerticalAlign(table)
	if vAlignValue < 0 || vAlignValue >= len(vAlignCss) {
		vAlignValue = 0
	}

	vAlign := vAlignCss[vAlignValue]

	tableCSS := func(startRow, endRow int, cellTag string, cellBorder BorderProperty, cellPadding BoundsProperty) {
		//var namedColors []NamedColor = nil

		for row := startRow; row < endRow; row++ {
			cssBuilder.buffer.Reset()
			if rowStyle != nil {
				if styles := rowStyle.RowStyle(row); styles != nil {
					view.Clear()
					for tag, value := range styles {
						view.Set(tag, value)
					}
					view.cssStyle(&view, &cssBuilder)
				}
			}

			buffer.WriteString(`<tr id="`)
			buffer.WriteString(table.rowID(row))
			buffer.WriteRune('"')

			if selectionMode == RowSelection {
				if row == table.current.Row {
					buffer.WriteString(` class="`)
					if table.HasFocus() {
						buffer.WriteString(table.currentStyle())
					} else {
						buffer.WriteString(table.currentInactiveStyle())
					}
					buffer.WriteRune('"')
				}

				buffer.WriteString(` onclick="tableRowClickEvent(this, event)"`)

				if allowRowSelection != nil && !allowRowSelection.AllowRowSelection(row) {
					buffer.WriteString(` data-disabled="1"`)
				}
			}

			if cssBuilder.buffer.Len() > 0 {
				buffer.WriteString(` style="`)
				buffer.WriteString(cssBuilder.buffer.String())
				buffer.WriteString(`"`)
			}
			buffer.WriteString(">")

			for column := 0; column < columnCount; column++ {
				ignore := false
				for _, cell := range ignoreCells {
					if cell.row == row && cell.column == column {
						ignore = true
						break
					}
				}

				if !ignore {
					rowSpan := 0
					columnSpan := 0

					cssBuilder.buffer.Reset()
					view.Clear()

					if cellBorder != nil {
						view.set(Border, cellBorder)
					}

					if cellPadding != nil {
						view.set(Padding, cellPadding)
					}

					if cellStyle != nil {
						if styles := cellStyle.CellStyle(row, column); styles != nil {
							for tag, value := range styles {
								valueToInt := func() int {
									switch value := value.(type) {
									case int:
										return value

									case string:
										if value, ok := session.resolveConstants(value); ok {
											if n, err := strconv.Atoi(value); err == nil {
												return n
											}
										}
									}
									return 0
								}

								switch tag = strings.ToLower(tag); tag {
								case RowSpan:
									rowSpan = valueToInt()

								case ColumnSpan:
									columnSpan = valueToInt()

								default:
									view.set(tag, value)
								}
							}
						}
					}

					if len(view.properties) > 0 {
						view.cssStyle(&view, &cssBuilder)
					}

					buffer.WriteRune('<')
					buffer.WriteString(cellTag)
					buffer.WriteString(` id="`)
					buffer.WriteString(table.cellID(row, column))
					buffer.WriteString(`" class="ruiView`)

					if selectionMode == CellSelection && row == table.current.Row && column == table.current.Column {
						buffer.WriteRune(' ')
						if table.HasFocus() {
							buffer.WriteString(table.currentStyle())
						} else {
							buffer.WriteString(table.currentInactiveStyle())
						}
					}
					buffer.WriteRune('"')

					if selectionMode == CellSelection {
						buffer.WriteString(` onclick="tableCellClickEvent(this, event)"`)
						if allowCellSelection != nil && !allowCellSelection.AllowCellSelection(row, column) {
							buffer.WriteString(` data-disabled="1"`)
						}
					}

					if columnSpan > 1 {
						buffer.WriteString(` colspan="`)
						buffer.WriteString(strconv.Itoa(columnSpan))
						buffer.WriteRune('"')
						for c := column + 1; c < column+columnSpan; c++ {
							ignoreCells = append(ignoreCells, struct {
								row    int
								column int
							}{row: row, column: c})
						}
					}

					if rowSpan > 1 {
						buffer.WriteString(` rowspan="`)
						buffer.WriteString(strconv.Itoa(rowSpan))
						buffer.WriteRune('"')
						if columnSpan < 1 {
							columnSpan = 1
						}
						for r := row + 1; r < row+rowSpan; r++ {
							for c := column; c < column+columnSpan; c++ {
								ignoreCells = append(ignoreCells, struct {
									row    int
									column int
								}{row: r, column: c})
							}
						}
					}

					if cssBuilder.buffer.Len() > 0 {
						buffer.WriteString(` style="`)
						buffer.WriteString(cssBuilder.buffer.String())
						buffer.WriteRune('"')
					}
					buffer.WriteRune('>')

					table.writeCellHtml(adapter, row, column, buffer)
					/*
						switch value := adapter.Cell(row, column).(type) {
						case string:
							buffer.WriteString(value)

						case View:
							viewHTML(value, buffer)
							table.cellViews = append(table.cellViews, value)

						case Color:
							buffer.WriteString(`<div style="display: inline; height: 1em; background-color: `)
							buffer.WriteString(value.cssString())
							buffer.WriteString(`">&nbsp;&nbsp;&nbsp;&nbsp;</div> `)
							buffer.WriteString(value.String())
							if namedColors == nil {
								namedColors = NamedColors()
							}
							for _, namedColor := range namedColors {
								if namedColor.Color == value {
									buffer.WriteString(" (")
									buffer.WriteString(namedColor.Name)
									buffer.WriteRune(')')
									break
								}
							}

						case fmt.Stringer:
							buffer.WriteString(value.String())

						case rune:
							buffer.WriteString(string(value))

						case float32:
							buffer.WriteString(fmt.Sprintf("%g", float64(value)))

						case float64:
							buffer.WriteString(fmt.Sprintf("%g", value))

						case bool:
							if value {
								buffer.WriteString(session.checkboxOnImage())
							} else {
								buffer.WriteString(session.checkboxOffImage())
							}

						default:
							if n, ok := isInt(value); ok {
								buffer.WriteString(fmt.Sprintf("%d", n))
							} else {
								buffer.WriteString("<Unsupported value>")
							}
						}
					*/

					buffer.WriteString(`</`)
					buffer.WriteString(cellTag)
					buffer.WriteRune('>')
				}
			}

			buffer.WriteString("</tr>")
		}
	}

	if columnStyle := table.getColumnStyle(); columnStyle != nil {
		buffer.WriteString("<colgroup>")
		for column := 0; column < columnCount; column++ {
			cssBuilder.buffer.Reset()
			if styles := columnStyle.ColumnStyle(column); styles != nil {
				view.Clear()
				for tag, value := range styles {
					view.Set(tag, value)
				}
				view.cssStyle(&view, &cssBuilder)
			}

			if cssBuilder.buffer.Len() > 0 {
				buffer.WriteString(`<col style="`)
				buffer.WriteString(cssBuilder.buffer.String())
				buffer.WriteString(`">`)
			} else {
				buffer.WriteString("<col>")
			}
		}
		buffer.WriteString("</colgroup>")
	}

	headHeight := GetTableHeadHeight(table)
	footHeight := GetTableFootHeight(table)
	cellBorder := table.getCellBorder()
	cellPadding := table.boundsProperty(CellPadding)
	if cellPadding == nil || len(cellPadding.AllTags()) == 0 {
		cellPadding = nil
		if style, ok := stringProperty(table, Style, table.Session()); ok {
			if style, ok := table.Session().resolveConstants(style); ok {
				cellPadding = table.cellPaddingFromStyle(style)
			}
		}
	}

	headFootStart := func(htmlTag, styleTag string) (BorderProperty, BoundsProperty) {
		buffer.WriteRune('<')
		buffer.WriteString(htmlTag)
		value := table.getRaw(styleTag)
		if value == nil {
			value = valueFromStyle(table, styleTag)
		}
		if value != nil {
			switch value := value.(type) {
			case string:
				if style, ok := session.resolveConstants(value); ok {
					buffer.WriteString(` class="`)
					buffer.WriteString(style)
					buffer.WriteString(`" style="vertical-align: `)
					buffer.WriteString(vAlign)
					buffer.WriteString(`;">`)

					return table.cellBorderFromStyle(style), table.cellPaddingFromStyle(style)
				}

			case Params:
				cssBuilder.buffer.Reset()
				view.Clear()
				view.Set(TableVerticalAlign, vAlignValue)
				for tag, val := range value {
					view.Set(tag, val)
				}

				var border BorderProperty = nil
				if value := view.Get(CellBorder); value != nil {
					border = value.(BorderProperty)
				}
				var padding BoundsProperty = nil
				if value := view.Get(CellPadding); value != nil {
					switch value := value.(type) {
					case SizeUnit:
						padding = NewBoundsProperty(Params{
							Top:    value,
							Right:  value,
							Bottom: value,
							Left:   value,
						})

					case BoundsProperty:
						padding = value
					}
				}

				view.cssStyle(&view, &cssBuilder)
				if cssBuilder.buffer.Len() > 0 {
					buffer.WriteString(` style="`)
					buffer.WriteString(cssBuilder.buffer.String())
					buffer.WriteString(`"`)
				}
				buffer.WriteRune('>')
				return border, padding
			}
		}

		buffer.WriteString(` style="vertical-align: `)
		buffer.WriteString(vAlign)
		buffer.WriteString(`;">`)
		return nil, nil
	}

	if headHeight > 0 {
		headCellBorder := cellBorder
		headCellPadding := cellPadding

		if headHeight > rowCount {
			headHeight = rowCount
		}

		border, padding := headFootStart("thead", HeadStyle)
		if border != nil {
			headCellBorder = border
		}
		if padding != nil {
			headCellPadding = padding
		}
		tableCSS(0, headHeight, "th", headCellBorder, headCellPadding)
		buffer.WriteString("</thead>")
	}

	if footHeight > rowCount-headHeight {
		footHeight = rowCount - headHeight
	}

	if rowCount > footHeight+headHeight {
		buffer.WriteString(`<tbody  style="vertical-align: `)
		buffer.WriteString(vAlign)
		buffer.WriteString(`;">`)
		tableCSS(headHeight, rowCount-footHeight, "td", cellBorder, cellPadding)
		buffer.WriteString("</tbody>")
	}

	if footHeight > 0 {
		footCellBorder := cellBorder
		footCellPadding := cellPadding

		border, padding := headFootStart("tfoot", FootStyle)
		if border != nil {
			footCellBorder = border
		}
		if padding != nil {
			footCellPadding = padding
		}
		tableCSS(rowCount-footHeight, rowCount, "td", footCellBorder, footCellPadding)
		buffer.WriteString("</tfoot>")
	}
}

func (table *tableViewData) cellPaddingFromStyle(style string) BoundsProperty {
	if value := table.Session().styleProperty(style, CellPadding); value != nil {
		switch value := value.(type) {
		case SizeUnit:
			return NewBoundsProperty(Params{
				Top:    value,
				Right:  value,
				Bottom: value,
				Left:   value,
			})

		case BoundsProperty:
			return value

		case string:
			if value, ok := table.Session().resolveConstants(value); ok {
				if strings.Contains(value, ",") {
					values := split4Values(value)
					switch len(values) {
					case 1:
						value = values[0]

					case 4:
						result := NewBoundsProperty(nil)
						n := 0
						for i, tag := range []string{Top, Right, Bottom, Left} {
							if size, ok := StringToSizeUnit(values[i]); ok && size.Type != Auto {
								result.Set(tag, size)
								n++
							}
						}
						if n > 0 {
							return result
						}
						return nil

					default:
						return nil
					}
				}

				if size, ok := StringToSizeUnit(value); ok && size.Type != Auto {
					return NewBoundsProperty(Params{
						Top:    size,
						Right:  size,
						Bottom: size,
						Left:   size,
					})
				}
			}
		}
	}

	return nil
}

func (table *tableViewData) writeCellHtml(adapter TableAdapter, row, column int, buffer *strings.Builder) {
	switch value := adapter.Cell(row, column).(type) {
	case string:
		buffer.WriteString(value)

	case View:
		viewHTML(value, buffer)
		table.cellViews = append(table.cellViews, value)

	case Color:
		buffer.WriteString(`<div style="display: inline; height: 1em; background-color: `)
		buffer.WriteString(value.cssString())
		buffer.WriteString(`">&nbsp;&nbsp;&nbsp;&nbsp;</div> `)
		buffer.WriteString(value.String())

		namedColors := NamedColors()
		for _, namedColor := range namedColors {
			if namedColor.Color == value {
				buffer.WriteString(" (")
				buffer.WriteString(namedColor.Name)
				buffer.WriteRune(')')
				break
			}
		}

	case fmt.Stringer:
		buffer.WriteString(value.String())

	case rune:
		buffer.WriteString(string(value))

	case float32:
		buffer.WriteString(fmt.Sprintf("%g", float64(value)))

	case float64:
		buffer.WriteString(fmt.Sprintf("%g", value))

	case bool:
		accentColor := Color(0)
		if color := GetAccentColor(table, ""); color != 0 {
			accentColor = color
		}
		if value {
			buffer.WriteString(table.Session().checkboxOnImage(accentColor))
		} else {
			buffer.WriteString(table.Session().checkboxOffImage(accentColor))
		}

	default:
		if n, ok := isInt(value); ok {
			buffer.WriteString(fmt.Sprintf("%d", n))
		} else {
			buffer.WriteString("<Unsupported value>")
		}
	}
}

func (table *tableViewData) cellBorderFromStyle(style string) BorderProperty {
	if value := table.Session().styleProperty(style, CellBorder); value != nil {
		if border, ok := value.(BorderProperty); ok {
			return border
		}
	}
	return nil
}

func (table *tableViewData) getCellBorder() BorderProperty {
	if value := table.getRaw(CellBorder); value != nil {
		if border, ok := value.(BorderProperty); ok {
			return border
		}
	}

	if style, ok := stringProperty(table, Style, table.Session()); ok {
		if style, ok := table.Session().resolveConstants(style); ok {
			return table.cellBorderFromStyle(style)
		}
	}

	return nil
}

func (table *tableViewData) getCurrent() CellIndex {
	return table.current
}

func (table *tableViewData) cssStyle(self View, builder cssBuilder) {
	session := table.Session()
	table.viewData.cssViewStyle(builder, session)

	gap, ok := sizeProperty(table, Gap, session)
	if !ok || gap.Type == Auto || gap.Value <= 0 {
		builder.add("border-spacing", "0")
		builder.add("border-collapse", "collapse")
	} else {
		builder.add("border-spacing", gap.cssString("0", session))
		builder.add("border-collapse", "separate")
	}
}

func (table *tableViewData) ReloadTableData() {
	session := table.Session()
	htmlID := table.htmlID()
	if content := table.content(); content != nil {
		session.updateProperty(htmlID, "data-rows", strconv.Itoa(content.RowCount()))
		session.updateProperty(htmlID, "data-columns", strconv.Itoa(content.ColumnCount()))
	}
	updateInnerHTML(htmlID, session)
}

func (table *tableViewData) onItemResize(self View, index string, x, y, width, height float64) {
	if n := strings.IndexRune(index, '-'); n > 0 {
		if row, err := strconv.Atoi(index[:n]); err == nil {
			if column, err := strconv.Atoi(index[n+1:]); err == nil {
				if content := table.content(); content != nil {
					i := row*content.ColumnCount() + column
					if i < len(table.cellFrame) {
						table.cellFrame[i].Left = x
						table.cellFrame[i].Top = y
						table.cellFrame[i].Width = width
						table.cellFrame[i].Height = height
					}
				}
			} else {
				ErrorLog(err.Error())
			}
		} else {
			ErrorLog(err.Error())
		}
	} else {
		ErrorLogF(`Invalid cell index: %s`, index)
	}
}

func (table *tableViewData) CellFrame(row, column int) Frame {
	if content := table.content(); content != nil {
		i := row*content.ColumnCount() + column
		if i < len(table.cellFrame) {
			return table.cellFrame[i]
		}
	}
	return Frame{}
}

func (table *tableViewData) ReloadCell(row, column int) {
	adapter := table.content()
	if adapter == nil {
		return
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	table.writeCellHtml(adapter, row, column, buffer)
	table.session.updateInnerHTML(table.cellID(row, column), buffer.String())
}

func (table *tableViewData) Views() []View {
	return table.cellViews
}

func (table *tableViewData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "currentRow":
		if row, ok := dataIntProperty(data, "row"); ok && row != table.current.Row {
			table.current.Row = row
			for _, listener := range table.rowSelectedListener {
				listener(table, row)
			}
		}

	case "currentCell":
		if row, ok := dataIntProperty(data, "row"); ok {
			if column, ok := dataIntProperty(data, "column"); ok {
				if row != table.current.Row || column != table.current.Column {
					table.current.Row = row
					table.current.Column = column
					for _, listener := range table.cellSelectedListener {
						listener(table, row, column)
					}
				}
			}
		}

	case "rowClick":
		if row, ok := dataIntProperty(data, "row"); ok {
			for _, listener := range table.rowClickedListener {
				listener(table, row)
			}
		}

	case "cellClick":
		if row, ok := dataIntProperty(data, "row"); ok {
			if column, ok := dataIntProperty(data, "column"); ok {
				for _, listener := range table.cellClickedListener {
					listener(table, row, column)
				}
			}
		}

	default:
		return table.viewData.handleCommand(self, command, data)
	}

	return true
}
