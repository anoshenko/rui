package rui

// TableAdapter describes the TableView content
type TableAdapter interface {
	// RowCount returns number of rows in the table
	RowCount() int

	// ColumnCount returns number of columns in the table
	ColumnCount() int

	// Cell returns the contents of a table cell. The function can return elements of the following types:
	// * string
	// * rune
	// * float32, float64
	// * integer values: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
	// * bool
	// * rui.Color
	// * rui.View
	// * fmt.Stringer
	// * rui.VerticalTableJoin, rui.HorizontalTableJoin
	Cell(row, column int) interface{}
}

// TableColumnStyle describes the style of TableView columns.
// To set column styles, you must either implement the TableColumnStyle interface in the table adapter
// or assign its separate implementation to the "column-style" property.
type TableColumnStyle interface {
	ColumnStyle(column int) Params
}

// TableRowStyle describes the style of TableView rows.
// To set row styles, you must either implement the TableRowStyle interface in the table adapter
// or assign its separate implementation to the "row-style" property.
type TableRowStyle interface {
	RowStyle(row int) Params
}

// TableCellStyle describes the style of TableView cells.
// To set row cells, you must either implement the TableCellStyle interface in the table adapter
// or assign its separate implementation to the "cell-style" property.
type TableCellStyle interface {
	CellStyle(row, column int) Params
}

// TableAllowCellSelection determines whether TableView cell selection is allowed.
// It is only used if the "selection-mode" property is set to CellSelection (1).
// To set cell selection allowing, you must either implement the TableAllowCellSelection interface
// in the table adapter or assign its separate implementation to the "allow-selection" property.
type TableAllowCellSelection interface {
	AllowCellSelection(row, column int) bool
}

// TableAllowRowSelection determines whether TableView row selection is allowed.
// It is only used if the "selection-mode" property is set to RowSelection (2).
// To set row selection allowing, you must either implement the TableAllowRowSelection interface
// in the table adapter or assign its separate implementation to the "allow-selection" property.
type TableAllowRowSelection interface {
	AllowRowSelection(row int) bool
}

// SimpleTableAdapter is implementation of TableAdapter where the content
// defines as [][]interface{}.
// When you assign [][]interface{} value to the "content" property, it is converted to SimpleTableAdapter
type SimpleTableAdapter interface {
	TableAdapter
	TableCellStyle
}

type simpleTableAdapter struct {
	content     [][]interface{}
	columnCount int
}

// TextTableAdapter is implementation of TableAdapter where the content
// defines as [][]string.
// When you assign [][]string value to the "content" property, it is converted to TextTableAdapter
type TextTableAdapter interface {
	TableAdapter
}

type textTableAdapter struct {
	content     [][]string
	columnCount int
}

// NewTextTableAdapter is an auxiliary structure. It used as cell content and
// specifies that the cell should be merged with the one above it
type VerticalTableJoin struct {
}

// HorizontalTableJoin is an auxiliary structure. It used as cell content and
// specifies that the cell should be merged with the one before it
type HorizontalTableJoin struct {
}

// NewSimpleTableAdapter creates the new SimpleTableAdapter
func NewSimpleTableAdapter(content [][]interface{}) SimpleTableAdapter {
	if content == nil {
		return nil
	}

	adapter := new(simpleTableAdapter)
	adapter.content = content
	adapter.columnCount = 0
	for _, row := range content {
		if row != nil {
			columnCount := len(row)
			if adapter.columnCount < columnCount {
				adapter.columnCount = columnCount
			}
		}
	}

	return adapter
}

func (adapter *simpleTableAdapter) RowCount() int {
	if adapter.content != nil {
		return len(adapter.content)
	}
	return 0
}

func (adapter *simpleTableAdapter) ColumnCount() int {
	return adapter.columnCount
}

func (adapter *simpleTableAdapter) Cell(row, column int) interface{} {
	if adapter.content != nil && row >= 0 && row < len(adapter.content) &&
		adapter.content[row] != nil && column >= 0 && column < len(adapter.content[row]) {
		return adapter.content[row][column]
	}
	return nil
}

func (adapter *simpleTableAdapter) CellStyle(row, column int) Params {
	if adapter.content == nil {
		return nil
	}

	getColumnSpan := func() int {
		count := 0
		for i := column + 1; i < adapter.columnCount; i++ {
			next := adapter.Cell(row, i)
			switch next.(type) {
			case HorizontalTableJoin:
				count++

			default:
				return count
			}
		}
		return count
	}

	getRowSpan := func() int {
		rowCount := len(adapter.content)
		count := 0
		for i := row + 1; i < rowCount; i++ {
			next := adapter.Cell(i, column)
			switch next.(type) {
			case VerticalTableJoin:
				count++

			default:
				return count
			}
		}
		return count
	}

	columnSpan := getColumnSpan()
	rowSpan := getRowSpan()

	var params Params = nil
	if rowSpan > 0 {
		params = Params{RowSpan: rowSpan + 1}
	}

	if columnSpan > 0 {
		if params == nil {
			params = Params{ColumnSpan: columnSpan + 1}
		} else {
			params[ColumnSpan] = columnSpan
		}
	}

	return params
}

// NewTextTableAdapter creates the new TextTableAdapter
func NewTextTableAdapter(content [][]string) TextTableAdapter {
	if content == nil {
		return nil
	}

	adapter := new(textTableAdapter)
	adapter.content = content
	adapter.columnCount = 0
	for _, row := range content {
		if row != nil {
			columnCount := len(row)
			if adapter.columnCount < columnCount {
				adapter.columnCount = columnCount
			}
		}
	}

	return adapter
}

func (adapter *textTableAdapter) RowCount() int {
	if adapter.content != nil {
		return len(adapter.content)
	}
	return 0
}

func (adapter *textTableAdapter) ColumnCount() int {
	return adapter.columnCount
}

func (adapter *textTableAdapter) Cell(row, column int) interface{} {
	if adapter.content != nil && row >= 0 && row < len(adapter.content) &&
		adapter.content[row] != nil && column >= 0 && column < len(adapter.content[row]) {
		return adapter.content[row][column]
	}
	return nil
}

type simpleTableRowStyle struct {
	params []Params
}

func (style *simpleTableRowStyle) RowStyle(row int) Params {
	if row < len(style.params) {
		params := style.params[row]
		if len(params) > 0 {
			return params
		}
	}
	return nil
}

func (table *tableViewData) setRowStyle(value interface{}) bool {
	newSimpleTableRowStyle := func(params []Params) TableRowStyle {
		if len(params) == 0 {
			return nil
		}
		result := new(simpleTableRowStyle)
		result.params = params
		return result
	}

	switch value := value.(type) {
	case TableRowStyle:
		table.properties[RowStyle] = value

	case []Params:
		if style := newSimpleTableRowStyle(value); style != nil {
			table.properties[RowStyle] = style
		} else {
			delete(table.properties, RowStyle)
		}

	case DataNode:
		if value.Type() == ArrayNode {
			params := make([]Params, value.ArraySize())
			for i, element := range value.ArrayElements() {
				params[i] = Params{}
				if element.IsObject() {
					obj := element.Object()
					for k := 0; k < obj.PropertyCount(); k++ {
						if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
							params[i][prop.Tag()] = prop.Text()
						}
					}
				} else {
					params[i][Style] = element.Value()
				}
			}
			if style := newSimpleTableRowStyle(params); style != nil {
				table.properties[RowStyle] = style
			} else {
				delete(table.properties, RowStyle)
			}
		} else {
			return false
		}

	default:
		return false
	}
	return true
}

func (table *tableViewData) getRowStyle() TableRowStyle {
	for _, tag := range []string{RowStyle, Content} {
		if value := table.getRaw(tag); value != nil {
			if style, ok := value.(TableRowStyle); ok {
				return style
			}
		}
	}
	return nil
}

type simpleTableColumnStyle struct {
	params []Params
}

func (style *simpleTableColumnStyle) ColumnStyle(row int) Params {
	if row < len(style.params) {
		params := style.params[row]
		if len(params) > 0 {
			return params
		}
	}
	return nil
}

func (table *tableViewData) setColumnStyle(value interface{}) bool {
	newSimpleTableColumnStyle := func(params []Params) TableColumnStyle {
		if len(params) == 0 {
			return nil
		}
		result := new(simpleTableColumnStyle)
		result.params = params
		return result
	}

	switch value := value.(type) {
	case TableColumnStyle:
		table.properties[ColumnStyle] = value

	case []Params:
		if style := newSimpleTableColumnStyle(value); style != nil {
			table.properties[ColumnStyle] = style
		} else {
			delete(table.properties, ColumnStyle)
		}

	case DataNode:
		if value.Type() == ArrayNode {
			params := make([]Params, value.ArraySize())
			for i, element := range value.ArrayElements() {
				params[i] = Params{}
				if element.IsObject() {
					obj := element.Object()
					for k := 0; k < obj.PropertyCount(); k++ {
						if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
							params[i][prop.Tag()] = prop.Text()
						}
					}
				} else {
					params[i][Style] = element.Value()
				}
			}
			if style := newSimpleTableColumnStyle(params); style != nil {
				table.properties[ColumnStyle] = style
			} else {
				delete(table.properties, ColumnStyle)
			}
		} else {
			return false
		}

	default:
		return false
	}
	return true
}

func (table *tableViewData) getColumnStyle() TableColumnStyle {
	for _, tag := range []string{ColumnStyle, Content} {
		if value := table.getRaw(tag); value != nil {
			if style, ok := value.(TableColumnStyle); ok {
				return style
			}
		}
	}
	return nil
}

func (table *tableViewData) getCellStyle() TableCellStyle {
	for _, tag := range []string{CellStyle, Content} {
		if value := table.getRaw(tag); value != nil {
			if style, ok := value.(TableCellStyle); ok {
				return style
			}
		}
	}
	return nil
}
