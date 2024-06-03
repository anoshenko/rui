package rui

// ListAdapter - the list data source
type ListAdapter interface {
	// ListSize returns the number of elements in the list
	ListSize() int

	// ListItem creates a View of a list item at the given index
	ListItem(index int, session Session) View
}

// ListItemEnabled implements the optional method of ListAdapter interface
type ListItemEnabled interface {
	// IsListItemEnabled returns the status (enabled/disabled) of a list item at the given index
	IsListItemEnabled(index int) bool
}

type textListAdapter struct {
	items  []string
	views  []View
	params Params
}

type viewListAdapter struct {
	items []View
}

// NewTextListAdapter create the new ListAdapter for a string list displaying. The second argument is parameters of a TextView item
func NewTextListAdapter(items []string, params Params) ListAdapter {
	if items == nil {
		return nil
	}
	adapter := new(textListAdapter)
	adapter.items = items
	if params != nil {
		adapter.params = params
	} else {
		adapter.params = Params{}
	}
	adapter.views = make([]View, len(items))
	return adapter
}

// NewTextListAdapter create the new ListAdapter for a view list displaying
func NewViewListAdapter(items []View) ListAdapter {
	if items != nil {
		adapter := new(viewListAdapter)
		adapter.items = items
		return adapter
	}
	return nil
}

func (adapter *textListAdapter) ListSize() int {
	return len(adapter.items)
}

func (adapter *textListAdapter) ListItem(index int, session Session) View {
	if index < 0 || index >= len(adapter.items) {
		return nil
	}

	if adapter.views[index] == nil {
		adapter.params[Text] = adapter.items[index]
		adapter.views[index] = NewTextView(session, adapter.params)
	}

	return adapter.views[index]
}

func (adapter *textListAdapter) IsListItemEnabled(index int) bool {
	return true
}

func (adapter *viewListAdapter) ListSize() int {
	return len(adapter.items)
}

func (adapter *viewListAdapter) ListItem(index int, session Session) View {
	if index >= 0 && index < len(adapter.items) {
		return adapter.items[index]
	}
	return nil
}

func (adapter *viewListAdapter) IsListItemEnabled(index int) bool {
	if index >= 0 && index < len(adapter.items) {
		return !IsDisabled(adapter.items[index])
	}
	return true
}
