package rui

func newTableCellView(session Session) *tableCellView {
	view := new(tableCellView)
	view.init(session)
	return view
}

func (cell *tableCellView) init(session Session) {
	cell.viewData.init(session)
	cell.normalize = func(tag PropertyName) PropertyName {
		if tag == VerticalAlign {
			return TableVerticalAlign
		}
		return tag
	}
}

func (cell *tableCellView) cssStyle(self View, builder cssBuilder) {
	session := cell.Session()
	cell.viewData.cssViewStyle(builder, session)

	if value, ok := enumProperty(cell, TableVerticalAlign, session, 0); ok {
		builder.add("vertical-align", enumProperties[TableVerticalAlign].values[value])
	}
}

// GetTableContent returns a TableAdapter which defines the TableView content.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableContent(view View, subviewID ...string) TableAdapter {
	if view = getSubview(view, subviewID); view != nil {
		if content := view.getRaw(Content); content != nil {
			if adapter, ok := content.(TableAdapter); ok {
				return adapter
			}
		}
	}

	return nil
}

// GetTableRowStyle returns a TableRowStyle which defines styles of TableView rows.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableRowStyle(view View, subviewID ...string) TableRowStyle {
	if view = getSubview(view, subviewID); view != nil {
		for _, tag := range []PropertyName{RowStyle, Content} {
			if value := view.getRaw(tag); value != nil {
				if style, ok := value.(TableRowStyle); ok {
					return style
				}
			}
		}
	}

	return nil
}

// GetTableColumnStyle returns a TableColumnStyle which defines styles of TableView columns.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableColumnStyle(view View, subviewID ...string) TableColumnStyle {
	if view = getSubview(view, subviewID); view != nil {
		for _, tag := range []PropertyName{ColumnStyle, Content} {
			if value := view.getRaw(tag); value != nil {
				if style, ok := value.(TableColumnStyle); ok {
					return style
				}
			}
		}
	}

	return nil
}

// GetTableCellStyle returns a TableCellStyle which defines styles of TableView cells.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableCellStyle(view View, subviewID ...string) TableCellStyle {
	if view = getSubview(view, subviewID); view != nil {
		for _, tag := range []PropertyName{CellStyle, Content} {
			if value := view.getRaw(tag); value != nil {
				if style, ok := value.(TableCellStyle); ok {
					return style
				}
			}
		}
		return nil
	}

	return nil
}

// GetTableSelectionMode returns the mode of the TableView elements selection.
// Valid values are NoneSelection (0), CellSelection (1), and RowSelection (2).
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableSelectionMode(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, SelectionMode, NoneSelection, false)
}

// GetTableVerticalAlign returns a vertical align in a TableView cell. Returns one of next values:
// TopAlign (0), BottomAlign (1), CenterAlign (2), and BaselineAlign (3)
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, TableVerticalAlign, TopAlign, false)
}

// GetTableHeadHeight returns the number of rows in the table header.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableHeadHeight(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, HeadHeight, 0)
}

// GetTableFootHeight returns the number of rows in the table footer.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableFootHeight(view View, subviewID ...string) int {
	return intStyledProperty(view, subviewID, FootHeight, 0)
}

// GetTableCurrent returns the row and column index of the TableView selected cell/row.
// If there is no selected cell/row or the selection mode is NoneSelection (0),
// then a value of the row and column index less than 0 is returned.
// If the selection mode is RowSelection (2) then the returned column index is less than 0.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableCurrent(view View, subviewID ...string) CellIndex {
	if view = getSubview(view, subviewID); view != nil {
		if selectionMode := GetTableSelectionMode(view); selectionMode != NoneSelection {
			return tableViewCurrent(view)
		}
	}
	return CellIndex{Row: -1, Column: -1}
}

// GetTableCellClickedListeners returns listeners of event which occurs when the user clicks on a table cell.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableCellClickedListeners(view View, subviewID ...string) []func(TableView, int, int) {
	return getTwoArgEventListeners[TableView, int](view, subviewID, TableCellClickedEvent)
}

// GetTableCellSelectedListeners returns listeners of event which occurs when a table cell becomes selected.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableCellSelectedListeners(view View, subviewID ...string) []func(TableView, int, int) {
	return getTwoArgEventListeners[TableView, int](view, subviewID, TableCellSelectedEvent)
}

// GetTableRowClickedListeners returns listeners of event which occurs when the user clicks on a table row.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableRowClickedListeners(view View, subviewID ...string) []func(TableView, int) {
	return getOneArgEventListeners[TableView, int](view, subviewID, TableRowClickedEvent)
}

// GetTableRowSelectedListeners returns listeners of event which occurs when a table row becomes selected.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTableRowSelectedListeners(view View, subviewID ...string) []func(TableView, int) {
	return getOneArgEventListeners[TableView, int](view, subviewID, TableRowSelectedEvent)
}

// ReloadTableViewData updates TableView
// If the second argument (subviewID) is not specified or it is "" then updates the first argument (TableView).
func ReloadTableViewData(view View, subviewID ...string) bool {
	if view = getSubview(view, subviewID); view != nil {
		if tableView, ok := view.(TableView); ok {
			tableView.ReloadTableData()
			return true
		}
	}
	return false
}

// ReloadTableViewCell updates the given table cell.
// If the last argument (subviewID) is not specified or it is "" then updates the cell of the first argument (TableView).
func ReloadTableViewCell(row, column int, view View, subviewID ...string) bool {
	if view = getSubview(view, subviewID); view != nil {
		if tableView, ok := view.(TableView); ok {
			tableView.ReloadCell(row, column)
			return true
		}
	}
	return false
}
