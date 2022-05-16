package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// TableVerticalAlign is the constant for the "table-vertical-align" property tag.
	// The "table-vertical-align" int property sets the vertical alignment of the content inside a table cell.
	// Valid values are TopAlign (0), BottomAlign (1), CenterAlign (2), and BaselineAlign (3, 4)
	TableVerticalAlign = "table-vertical-align"

	// HeadHeight is the constant for the "head-height" property tag.
	// The "head-height" int property sets the number of rows in the table header.
	// The default value is 0 (no header)
	HeadHeight = "head-height"

	// HeadStyle is the constant for the "head-style" property tag.
	// The "head-style" string property sets the header style name
	HeadStyle = "head-style"

	// FootHeight is the constant for the "foot-height" property tag.
	// The "foot-height" int property sets the number of rows in the table footer.
	// The default value is 0 (no footer)
	FootHeight = "foot-height"

	// FootStyle is the constant for the "foot-style" property tag.
	// The "foot-style" string property sets the footer style name
	FootStyle = "foot-style"

	// RowSpan is the constant for the "row-span" property tag.
	// The "row-span" int property sets the number of table row to span.
	// Used only when specifying cell parameters in the implementation of TableCellStyle
	RowSpan = "row-span"

	// ColumnSpan is the constant for the "column-span" property tag.
	// The "column-span" int property sets the number of table column to span.
	// Used only when specifying cell parameters in the implementation of TableCellStyle
	ColumnSpan = "column-span"

	// RowStyle is the constant for the "row-style" property tag.
	// The "row-style" property sets the adapter which specifies styles of each table row.
	// This property can be assigned or by an implementation of TableRowStyle interface, or by an array of Params.
	RowStyle = "row-style"

	// ColumnStyle is the constant for the "column-style" property tag.
	// The "column-style" property sets the adapter which specifies styles of each table column.
	// This property can be assigned or by an implementation of TableColumnStyle interface, or by an array of Params.
	ColumnStyle = "column-style"

	// CellStyle is the constant for the "cell-style" property tag.
	// The "cell-style" property sets the adapter which specifies styles of each table cell.
	// This property can be assigned only by an implementation of TableCellStyle interface.
	CellStyle = "cell-style"

	// CellPadding is the constant for the "cell-padding" property tag.
	// The "cell-padding" Bounds property sets the padding area on all four sides of a table call at once.
	// An element's padding area is the space between its content and its border.
	CellPadding = "cell-padding"

	// CellPaddingLeft is the constant for the "cell-padding-left" property tag.
	// The "cell-padding-left" SizeUnit property sets the width of the padding area to the left of a cell content.
	// An element's padding area is the space between its content and its border.
	CellPaddingLeft = "cell-padding-left"

	// CellPaddingRight is the constant for the "cell-padding-right" property tag.
	// The "cell-padding-right" SizeUnit property sets the width of the padding area to the left of a cell content.
	// An element's padding area is the space between its content and its border.
	CellPaddingRight = "cell-padding-right"

	// CellPaddingTop is the constant for the "cell-padding-top" property tag.
	// The "cell-padding-top" SizeUnit property sets the height of the padding area to the top of a cell content.
	// An element's padding area is the space between its content and its border.
	CellPaddingTop = "cell-padding-top"

	// CellPaddingBottom is the constant for the "cell-padding-bottom" property tag.
	// The "cell-padding-bottom" SizeUnit property sets the height of the padding area to the bottom of a cell content.
	CellPaddingBottom = "cell-padding-bottom"

	// CellBorder is the constant for the "cell-border" property tag.
	// The "cell-border" property sets a table cell's border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorder = "cell-border"

	// CellBorderLeft is the constant for the "cell-border-left" property tag.
	// The "cell-border-left" property sets a view's left border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorderLeft = "cell-border-left"

	// CellBorderRight is the constant for the "cell-border-right" property tag.
	// The "cell-border-right" property sets a view's right border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorderRight = "cell-border-right"

	// CellBorderTop is the constant for the "cell-border-top" property tag.
	// The "cell-border-top" property sets a view's top border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorderTop = "cell-border-top"

	// CellBorderBottom is the constant for the "cell-border-bottom" property tag.
	// The "cell-border-bottom" property sets a view's bottom border. It sets the values of a border width, style, and color.
	// This property can be assigned a value of BorderProperty type, or ViewBorder type, or BorderProperty text representation.
	CellBorderBottom = "cell-border-bottom"

	// CellBorderStyle is the constant for the "cell-border-style" property tag.
	// The "cell-border-style" int property sets the line style for all four sides of a table cell's border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderStyle = "cell-border-style"

	// CellBorderLeftStyle is the constant for the "cell-border-left-style" property tag.
	// The "cell-border-left-style" int property sets the line style of a table cell's left border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderLeftStyle = "cell-border-left-style"

	// CellBorderRightStyle is the constant for the "cell-border-right-style" property tag.
	// The "cell-border-right-style" int property sets the line style of a table cell's right border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderRightStyle = "cell-border-right-style"

	// CellBorderTopStyle is the constant for the "cell-border-top-style" property tag.
	// The "cell-border-top-style" int property sets the line style of a table cell's top border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderTopStyle = "cell-border-top-style"

	// CellBorderBottomStyle is the constant for the "cell-border-bottom-style" property tag.
	// The "cell-border-bottom-style" int property sets the line style of a table cell's bottom border.
	// Valid values are NoneLine (0), SolidLine (1), DashedLine (2), DottedLine (3), and DoubleLine (4).
	CellBorderBottomStyle = "cell-border-bottom-style"

	// CellBorderWidth is the constant for the "cell-border-width" property tag.
	// The "cell-border-width" property sets the line width for all four sides of a table cell's border.
	CellBorderWidth = "cell-border-width"

	// CellBorderLeftWidth is the constant for the "cell-border-left-width" property tag.
	// The "cell-border-left-width" SizeUnit property sets the line width of a table cell's left border.
	CellBorderLeftWidth = "cell-border-left-width"

	// CellBorderRightWidth is the constant for the "cell-border-right-width" property tag.
	// The "cell-border-right-width" SizeUnit property sets the line width of a table cell's right border.
	CellBorderRightWidth = "cell-border-right-width"

	// CellBorderTopWidth is the constant for the "cell-border-top-width" property tag.
	// The "cell-border-top-width" SizeUnit property sets the line width of a table cell's top border.
	CellBorderTopWidth = "cell-border-top-width"

	// CellBorderBottomWidth is the constant for the "cell-border-bottom-width" property tag.
	// The "cell-border-bottom-width" SizeUnit property sets the line width of a table cell's bottom border.
	CellBorderBottomWidth = "cell-border-bottom-width"

	// CellBorderColor is the constant for the "cell-border-color" property tag.
	// The "cell-border-color" property sets the line color for all four sides of a table cell's border.
	CellBorderColor = "cell-border-color"

	// CellBorderLeftColor is the constant for the "cell-border-left-color" property tag.
	// The "cell-border-left-color" property sets the line color of a table cell's left border.
	CellBorderLeftColor = "cell-border-left-color"

	// CellBorderRightColor is the constant for the "cell-border-right-color" property tag.
	// The "cell-border-right-color" property sets the line color of a table cell's right border.
	CellBorderRightColor = "cell-border-right-color"

	// CellBorderTopColor is the constant for the "cell-border-top-color" property tag.
	// The "cell-border-top-color" property sets the line color of a table cell's top border.
	CellBorderTopColor = "cell-border-top-color"

	// CellBorderBottomColor is the constant for the "cell-border-bottom-color" property tag.
	// The "cell-border-bottom-color" property sets the line color of a table cell's bottom border.
	CellBorderBottomColor = "cell-border-bottom-color"

	// SelectionMode is the constant for the "selection-mode" property tag.
	// The "selection-mode" int property sets the mode of the table elements selection.
	// Valid values are NoneSelection (0), CellSelection (1), and RowSelection (2)
	SelectionMode = "selection-mode"

	// TableCellClickedEvent is the constant for "table-cell-clicked" property tag.
	// The "table-cell-clicked" event occurs when the user clicks on a table cell.
	// The main listener format: func(TableView, int, int), where the second argument is the row number,
	// and third argument is the column number.
	TableCellClickedEvent = "table-cell-clicked"

	// TableCellSelectedEvent is the constant for "table-cell-selected" property tag.
	// The "table-cell-selected" event occurs when a table cell becomes selected.
	// The main listener format: func(TableView, int, int), where the second argument is the row number,
	// and third argument is the column number.
	TableCellSelectedEvent = "table-cell-selected"

	// TableRowClickedEvent is the constant for "table-row-clicked" property tag.
	// The "table-row-clicked" event occurs when the user clicks on a table row.
	// The main listener format: func(TableView, int), where the second argument is the row number.
	TableRowClickedEvent = "table-row-clicked"

	// TableRowSelectedEvent is the constant for "table-row-selected" property tag.
	// The "table-row-selected" event occurs when a table row becomes selected.
	// The main listener format: func(TableView, int), where the second argument is the row number.
	TableRowSelectedEvent = "table-row-selected"

	// AllowSelection is the constant for the "allow-selection" property tag.
	// The "allow-selection" property sets the adapter which specifies styles of each table row.
	// This property can be assigned or by an implementation of TableAllowCellSelection
	// or TableAllowRowSelection interface.
	AllowSelection = "allow-selection"

	// NoneSelection is the value of "selection-mode" property: the selection is forbidden.
	NoneSelection = 0
	// CellSelection is the value of "selection-mode" property: the selection of a single cell only is enabled.
	CellSelection = 1
	// RowSelection is the value of "selection-mode" property: the selection of a table row only is enabled.
	RowSelection = 2
)

// CellIndex defines coordinates of the TableView cell
type CellIndex struct {
	Row, Column int
}

// TableView - text View
type TableView interface {
	View
	ParanetView
	ReloadTableData()
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
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newTableView(session Session) View {
	return NewTableView(session, nil)
}

// Init initialize fields of TableView by default values
func (table *tableViewData) Init(session Session) {
	table.viewData.Init(session)
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
	return GetTableSelectionMode(table, "") != NoneSelection
}

func (table *tableViewData) Get(tag string) interface{} {
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

func (table *tableViewData) Set(tag string, value interface{}) bool {
	return table.set(table.normalizeTag(tag), value)
}

func (table *tableViewData) set(tag string, value interface{}) bool {
	if value == nil {
		table.remove(tag)
		return true
	}

	switch tag {
	case Content:
		switch val := value.(type) {
		case TableAdapter:
			table.properties[Content] = value

		case [][]interface{}:
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
		listeners := table.valueToRowListeners(value)
		if listeners == nil {
			notCompatibleType(tag, value)
			return false
		}
		table.rowClickedListener = listeners

	case TableRowSelectedEvent:
		listeners := table.valueToRowListeners(value)
		if listeners == nil {
			notCompatibleType(tag, value)
			return false
		}
		table.rowSelectedListener = []func(TableView, int){}

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

		case DataNode:
			switch value.Type() {
			case ObjectNode:
				obj := value.Object()
				params := Params{}
				for k := 0; k < obj.PropertyCount(); k++ {
					if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
						params[prop.Tag()] = prop.Text()
					}
				}
				if len(params) > 0 {
					table.properties[tag] = params
				} else {
					delete(table.properties, tag)
				}

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
		case int:
			table.current.Row = value
			table.current.Column = -1

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
			notCompatibleType(tag, value)
			return false
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
			TableRowSelectedEvent, AllowSelection:
			table.ReloadTableData()

		case Gap:
			htmlID := table.htmlID()
			session := table.Session()
			gap, ok := sizeProperty(table, Gap, session)
			if !ok || gap.Type == Auto || gap.Value <= 0 {
				updateCSSProperty(htmlID, "border-spacing", "0", session)
				updateCSSProperty(htmlID, "border-collapse", "collapse", session)
			} else {
				updateCSSProperty(htmlID, "border-spacing", gap.cssString("0"), session)
				updateCSSProperty(htmlID, "border-collapse", "separate", session)
			}

		case SelectionMode:
			htmlID := table.htmlID()
			session := table.Session()

			switch GetTableSelectionMode(table, "") {
			case CellSelection:
				updateProperty(htmlID, "tabindex", "0", session)
				updateProperty(htmlID, "onfocus", "tableViewFocusEvent(this, event)", session)
				updateProperty(htmlID, "onblur", "tableViewBlurEvent(this, event)", session)
				updateProperty(htmlID, "data-selection", "cell", session)
				updateProperty(htmlID, "data-focusitemstyle", table.currentStyle(), session)
				updateProperty(htmlID, "data-bluritemstyle", table.currentInactiveStyle(), session)

				if table.current.Row >= 0 && table.current.Column >= 0 {
					updateProperty(htmlID, "data-current", table.cellID(table.current.Row, table.current.Column), session)
				} else {
					removeProperty(htmlID, "data-current", session)
				}
				updateProperty(htmlID, "onkeydown", "tableViewCellKeyDownEvent(this, event)", session)

			case RowSelection:
				updateProperty(htmlID, "tabindex", "0", session)
				updateProperty(htmlID, "onfocus", "tableViewFocusEvent(this, event)", session)
				updateProperty(htmlID, "onblur", "tableViewBlurEvent(this, event)", session)
				updateProperty(htmlID, "data-selection", "cell", session)
				updateProperty(htmlID, "data-focusitemstyle", table.currentStyle(), session)
				updateProperty(htmlID, "data-bluritemstyle", table.currentInactiveStyle(), session)

				if table.current.Row >= 0 {
					updateProperty(htmlID, "data-current", table.rowID(table.current.Row), session)
				} else {
					removeProperty(htmlID, "data-current", session)
				}
				updateProperty(htmlID, "onkeydown", "tableViewRowKeyDownEvent(this, event)", session)

			default: // NoneSelection
				for _, prop := range []string{"tabindex", "data-current", "onfocus", "onblur", "onkeydown", "data-selection"} {
					removeProperty(htmlID, prop, session)
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
	return "ruiCurrentTableCell"
}

func (table *tableViewData) valueToCellListeners(value interface{}) []func(TableView, int, int) {
	if value == nil {
		return []func(TableView, int, int){}
	}

	switch value := value.(type) {
	case func(TableView, int, int):
		return []func(TableView, int, int){value}

	case func(int, int):
		fn := func(view TableView, row, column int) {
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
			listeners[i] = func(view TableView, row, column int) {
				val(row, column)
			}
		}
		return listeners

	case []interface{}:
		listeners := make([]func(TableView, int, int), len(value))
		for i, val := range value {
			if val == nil {
				return nil
			}
			switch val := val.(type) {
			case func(TableView, int, int):
				listeners[i] = val

			case func(int, int):
				listeners[i] = func(view TableView, row, column int) {
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

func (table *tableViewData) valueToRowListeners(value interface{}) []func(TableView, int) {
	if value == nil {
		return []func(TableView, int){}
	}

	switch value := value.(type) {
	case func(TableView, int):
		return []func(TableView, int){value}

	case func(int):
		fn := func(view TableView, index int) {
			value(index)
		}
		return []func(TableView, int){fn}

	case []func(TableView, int):
		return value

	case []func(int):
		listeners := make([]func(TableView, int), len(value))
		for i, val := range value {
			if val == nil {
				return nil
			}
			listeners[i] = func(view TableView, index int) {
				val(index)
			}
		}
		return listeners

	case []interface{}:
		listeners := make([]func(TableView, int), len(value))
		for i, val := range value {
			if val == nil {
				return nil
			}
			switch val := val.(type) {
			case func(TableView, int):
				listeners[i] = val

			case func(int):
				listeners[i] = func(view TableView, index int) {
					val(index)
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

	if selectionMode := GetTableSelectionMode(table, ""); selectionMode != NoneSelection {
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

	var cssBuilder viewCSSBuilder
	cssBuilder.buffer = allocStringBuilder()
	defer freeStringBuilder(cssBuilder.buffer)

	var view tableCellView
	view.Init(session)

	ignorCells := []struct{ row, column int }{}
	selectionMode := GetTableSelectionMode(table, "")

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
	vAlignValue := GetTableVerticalAlign(table, "")
	if vAlignValue < 0 || vAlignValue >= len(vAlignCss) {
		vAlignValue = 0
	}

	vAlign := vAlignCss[vAlignValue]

	tableCSS := func(startRow, endRow int, cellTag string, cellBorder BorderProperty, cellPadding BoundsProperty) {
		var namedColors []NamedColor = nil

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
				for _, cell := range ignorCells {
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
							ignorCells = append(ignorCells, struct {
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
								ignorCells = append(ignorCells, struct {
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

					switch value := adapter.Cell(row, column).(type) {
					case string:
						buffer.WriteString(textToJS(value))

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
						buffer.WriteString(textToJS(value.String()))

					case rune:
						buffer.WriteString(textToJS(string(value)))

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

	headHeight := GetTableHeadHeight(table, "")
	footHeight := GetTableFootHeight(table, "")
	cellBorder := table.getCellBorder()
	cellPadding := table.boundsProperty(CellPadding)
	if cellPadding == nil {
		if style, ok := stringProperty(table, Style, table.Session()); ok {
			if style, ok := table.Session().resolveConstants(style); ok {
				cellPadding = table.cellPaddingFromStyle(style)
			}
		}
	}

	headFootStart := func(htmlTag, styleTag string) (BorderProperty, BoundsProperty) {
		buffer.WriteRune('<')
		buffer.WriteString(htmlTag)
		if value := table.getRaw(styleTag); value != nil {
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
	session := table.Session()
	var result BoundsProperty = nil

	if node := session.stylePropertyNode(style, CellPadding); node != nil && node.Type() == ObjectNode {
		for _, tag := range []string{Left, Right, Top, Bottom} {
			if node := node.Object().PropertyWithTag(tag); node != nil && node.Type() == TextNode {
				if result == nil {
					result = NewBoundsProperty(nil)
				}
				result.Set(tag, node.Text())
			}
		}
	}

	for _, tag := range []string{CellPaddingLeft, CellPaddingRight, CellPaddingTop, CellPaddingBottom} {
		if value, ok := session.styleProperty(style, CellPadding); ok {
			if result == nil {
				result = NewBoundsProperty(nil)
			}
			result.Set(tag, value)
		}
	}

	return result
}

func (table *tableViewData) cellBorderFromStyle(style string) BorderProperty {

	border := new(borderProperty)
	border.properties = map[string]interface{}{}

	session := table.Session()
	if node := session.stylePropertyNode(style, CellBorder); node != nil && node.Type() == ObjectNode {
		border.setBorderObject(node.Object())
	}

	for _, tag := range []string{
		CellBorderLeft,
		CellBorderRight,
		CellBorderTop,
		CellBorderBottom,
		CellBorderStyle,
		CellBorderLeftStyle,
		CellBorderRightStyle,
		CellBorderTopStyle,
		CellBorderBottomStyle,
		CellBorderWidth,
		CellBorderLeftWidth,
		CellBorderRightWidth,
		CellBorderTopWidth,
		CellBorderBottomWidth,
		CellBorderColor,
		CellBorderLeftColor,
		CellBorderRightColor,
		CellBorderTopColor,
		CellBorderBottomColor,
	} {
		if value, ok := session.styleProperty(style, tag); ok {
			border.Set(tag, value)
		}
	}

	if len(border.properties) == 0 {
		return nil
	}
	return border
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
	table.viewData.cssViewStyle(builder, table.Session())

	gap, ok := sizeProperty(table, Gap, table.Session())
	if !ok || gap.Type == Auto || gap.Value <= 0 {
		builder.add("border-spacing", "0")
		builder.add("border-collapse", "collapse")
	} else {
		builder.add("border-spacing", gap.cssString("0"))
		builder.add("border-collapse", "separate")
	}
}

func (table *tableViewData) ReloadTableData() {
	if content := table.content(); content != nil {
		updateProperty(table.htmlID(), "data-rows", strconv.Itoa(content.RowCount()), table.Session())
		updateProperty(table.htmlID(), "data-columns", strconv.Itoa(content.ColumnCount()), table.Session())
	}
	updateInnerHTML(table.htmlID(), table.Session())
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
