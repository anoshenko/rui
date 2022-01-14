package rui

import "strings"

func (cell *tableCellView) Set(tag string, value interface{}) bool {
	return cell.set(strings.ToLower(tag), value)
}

func (cell *tableCellView) set(tag string, value interface{}) bool {
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

// GetSelectionMode returns the mode of the TableView elements selection.
// Valid values are NoneSelection (0), CellSelection (1), and RowSelection (2).
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetSelectionMode(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if result, ok := enumStyledProperty(view, SelectionMode, NoneSelection); ok {
			return result
		}
	}
	return NoneSelection
}

// GetSelectionMode returns the index of the TableView selected row.
// If there is no selected row, then a value less than 0 are returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetCurrentTableRow(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		if selectionMode := GetSelectionMode(view, ""); selectionMode != NoneSelection {
			if tableView, ok := view.(TableView); ok {
				return tableView.getCurrent().Row
			}
		}
	}
	return -1
}

// GetCurrentTableCell returns the row and column index of the TableView selected cell.
// If there is no selected cell, then a value of the row and column index less than 0 is returned.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetCurrentTableCell(view View, subviewID string) CellIndex {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		if selectionMode := GetSelectionMode(view, ""); selectionMode != NoneSelection {
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
