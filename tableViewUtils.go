package rui

import "strings"

func (cell *tableCellView) Set(tag string, value any) bool {
	return cell.set(strings.ToLower(tag), value)
}

func (cell *tableCellView) set(tag string, value any) bool {
	switch tag {
	case VerticalAlign:
		tag = TableVerticalAlign
	}
	return cell.viewData.set(tag, value)
}

func (cell *tableCellView) cssStyle(self View, builder cssBuilder) {
	session := cell.Session()
	cell.viewData.cssViewStyle(builder, session)

	if value, ok := enumProperty(cell, TableVerticalAlign, session, 0); ok {
		builder.add("vertical-align", enumProperties[TableVerticalAlign].values[value])
	}
}

// GetTableContent returns a TableAdapter which defines the TableView content.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableContent(view View, subviewID string) TableAdapter {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		if tableView, ok := view.(TableView); ok {
			return tableView.content()
		}
	}

	return nil
}

// GetTableRowStyle returns a TableRowStyle which defines styles of TableView rows.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableRowStyle(view View, subviewID string) TableRowStyle {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		if tableView, ok := view.(TableView); ok {
			return tableView.getRowStyle()
		}
	}

	return nil
}

// GetTableColumnStyle returns a TableColumnStyle which defines styles of TableView columns.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableColumnStyle(view View, subviewID string) TableColumnStyle {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		if tableView, ok := view.(TableView); ok {
			return tableView.getColumnStyle()
		}
	}

	return nil
}

// GetTableCellStyle returns a TableCellStyle which defines styles of TableView cells.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableCellStyle(view View, subviewID string) TableCellStyle {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		if tableView, ok := view.(TableView); ok {
			return tableView.getCellStyle()
		}
	}

	return nil
}

// GetTableSelectionMode returns the mode of the TableView elements selection.
// Valid values are NoneSelection (0), CellSelection (1), and RowSelection (2).
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableSelectionMode(view View, subviewID string) int {
	return enumStyledProperty(view, subviewID, SelectionMode, NoneSelection, false)
}

// GetTableVerticalAlign returns a vertical align in a TavleView cell. Returns one of next values:
// TopAlign (0), BottomAlign (1), CenterAlign (2), and BaselineAlign (3)
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableVerticalAlign(view View, subviewID string) int {
	return enumStyledProperty(view, subviewID, TableVerticalAlign, TopAlign, false)
}

// GetTableHeadHeight returns the number of rows in the table header.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableHeadHeight(view View, subviewID string) int {
	return intStyledProperty(view, subviewID, HeadHeight, 0)
}

// GetTableFootHeight returns the number of rows in the table footer.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableFootHeight(view View, subviewID string) int {
	return intStyledProperty(view, subviewID, FootHeight, 0)
}

// GetTableCurrent returns the row and column index of the TableView selected cell/row.
// If there is no selected cell/row or the selection mode is NoneSelection (0),
// then a value of the row and column index less than 0 is returned.
// If the selection mode is RowSelection (2) then the returned column index is less than 0.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableCurrent(view View, subviewID string) CellIndex {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		if selectionMode := GetTableSelectionMode(view, ""); selectionMode != NoneSelection {
			if tableView, ok := view.(TableView); ok {
				return tableView.getCurrent()
			}
		}
	}
	return CellIndex{Row: -1, Column: -1}
}

// GetTableCellClickedListeners returns listeners of event which occurs when the user clicks on a table cell.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableCellClickedListeners(view View, subviewID string) []func(TableView, int, int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(TableCellClickedEvent); value != nil {
			if result, ok := value.([]func(TableView, int, int)); ok {
				return result
			}
		}
	}
	return []func(TableView, int, int){}
}

// GetTableCellSelectedListeners returns listeners of event which occurs when a table cell becomes selected.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableCellSelectedListeners(view View, subviewID string) []func(TableView, int, int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(TableCellSelectedEvent); value != nil {
			if result, ok := value.([]func(TableView, int, int)); ok {
				return result
			}
		}
	}
	return []func(TableView, int, int){}
}

// GetTableRowClickedListeners returns listeners of event which occurs when the user clicks on a table row.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableRowClickedListeners(view View, subviewID string) []func(TableView, int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(TableRowClickedEvent); value != nil {
			if result, ok := value.([]func(TableView, int)); ok {
				return result
			}
		}
	}
	return []func(TableView, int){}
}

// GetTableRowSelectedListeners returns listeners of event which occurs when a table row becomes selected.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetTableRowSelectedListeners(view View, subviewID string) []func(TableView, int) {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(TableRowSelectedEvent); value != nil {
			if result, ok := value.([]func(TableView, int)); ok {
				return result
			}
		}
	}
	return []func(TableView, int){}
}

// ReloadTableViewData updates TableView
func ReloadTableViewData(view View, subviewID string) bool {
	var tableView TableView
	if subviewID != "" {
		if tableView = TableViewByID(view, subviewID); tableView == nil {
			return false
		}
	} else {
		var ok bool
		if tableView, ok = view.(TableView); !ok {
			return false
		}
	}

	tableView.ReloadTableData()
	return true
}
