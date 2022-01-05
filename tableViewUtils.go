package rui

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
