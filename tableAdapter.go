package rui

// TableAdapter describes the [TableView] content
type TableAdapter interface {
	// RowCount returns number of rows in the table
	RowCount() int

	// ColumnCount returns number of columns in the table
	ColumnCount() int

	// Cell returns the contents of a table cell. The function can return elements of the following types:
	//   - string
	//   - rune
	//   - float32, float64
	//   - integer values: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64
	//   - bool
	//   - rui.Color
	//   - rui.View
	//   - fmt.Stringer
	//   - rui.VerticalTableJoin, rui.HorizontalTableJoin
	Cell(row, column int) any
}

// TableColumnStyle describes the style of [TableView] columns.
//
// To set column styles, you must either implement the [TableColumnStyle] interface in the table adapter
// or assign its separate implementation to the "column-style" property.
type TableColumnStyle interface {
	// ColumnStyle returns a map of properties which describe the style of the column
	ColumnStyle(column int) Params
}

// TableRowStyle describes the style of [TableView] rows.
//
// To set row styles, you must either implement the [TableRowStyle] interface in the table adapter
// or assign its separate implementation to the "row-style" property.
type TableRowStyle interface {
	// RowStyle returns a map of properties which describe the style of the row
	RowStyle(row int) Params
}

// TableCellStyle describes the style of [TableView] cells.
//
// To set row cells, you must either implement the [TableCellStyle] interface in the table adapter
// or assign its separate implementation to the "cell-style" property.
type TableCellStyle interface {
	// CellStyle returns a map of properties which describe the style of the cell
	CellStyle(row, column int) Params
}

// TableAllowCellSelection determines whether [TableView] cell selection is allowed.
//
// It is only used if the "selection-mode" property is set to CellSelection (1).
//
// To set cell selection allowing, you must either implement the TableAllowCellSelection interface
// in the table adapter or assign its separate implementation to the "allow-selection" property.
type TableAllowCellSelection interface {
	// AllowCellSelection returns "true" if we allow the user to select particular cell at specific rows and columns
	AllowCellSelection(row, column int) bool
}

// TableAllowRowSelection determines whether [TableView] row selection is allowed.
//
// It is only used if the "selection-mode" property is set to RowSelection (2).
//
// To set row selection allowing, you must either implement the TableAllowRowSelection interface
// in the table adapter or assign its separate implementation to the "allow-selection" property.
type TableAllowRowSelection interface {
	// AllowRowSelection returns "true" if we allow the user to select particular row in the table
	AllowRowSelection(row int) bool
}

// SimpleTableAdapter is implementation of [TableAdapter] where the content
// defines as [][]any.
//
// When you assign [][]any value to the "content" property, it is converted to SimpleTableAdapter
type SimpleTableAdapter interface {
	TableAdapter
	TableCellStyle
}

type simpleTableAdapter struct {
	content     [][]any
	columnCount int
}

// TextTableAdapter is implementation of [TableAdapter] where the content
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
func NewSimpleTableAdapter(content [][]any) SimpleTableAdapter {
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

func (adapter *simpleTableAdapter) Cell(row, column int) any {
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

func (adapter *textTableAdapter) Cell(row, column int) any {
	if adapter.content != nil && row >= 0 && row < len(adapter.content) &&
		adapter.content[row] != nil && column >= 0 && column < len(adapter.content[row]) {
		return adapter.content[row][column]
	}
	return nil
}

type simpleTableLineStyle struct {
	params []Params
}

func (style *simpleTableLineStyle) ColumnStyle(column int) Params {
	if column < len(style.params) {
		params := style.params[column]
		if len(params) > 0 {
			return params
		}
	}
	return nil
}

func (style *simpleTableLineStyle) RowStyle(row int) Params {
	if row < len(style.params) {
		params := style.params[row]
		if len(params) > 0 {
			return params
		}
	}
	return nil
}
