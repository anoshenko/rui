package rui

import "strings"

// TODO Expanded event

const (
	// Summary is the constant for the "summary" property tag.
	// The contents of the "summary" property are used as the label for the disclosure widget.
	Summary = "summary"
	// Expanded is the constant for the "expanded" property tag.
	// If the "expanded" boolean property is "true", then the content of view is visible.
	// If the value is "false" then the content is collapsed.
	Expanded = "expanded"
)

// DetailsView - collapsible container of View
type DetailsView interface {
	ViewsContainer
}

type detailsViewData struct {
	viewsContainerData
}

// NewDetailsView create new DetailsView object and return it
func NewDetailsView(session Session, params Params) DetailsView {
	view := new(detailsViewData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newDetailsView(session Session) View {
	return NewDetailsView(session, nil)
}

// Init initialize fields of DetailsView by default values
func (detailsView *detailsViewData) Init(session Session) {
	detailsView.viewsContainerData.Init(session)
	detailsView.tag = "DetailsView"
	//detailsView.systemClass = "ruiDetailsView"
}

func (detailsView *detailsViewData) Remove(tag string) {
	detailsView.remove(strings.ToLower(tag))
}

func (detailsView *detailsViewData) remove(tag string) {
	detailsView.viewsContainerData.remove(tag)
	if detailsView.created {
		switch tag {
		case Summary:
			updateInnerHTML(detailsView.htmlID(), detailsView.Session())

		case Expanded:
			removeProperty(detailsView.htmlID(), "open", detailsView.Session())
		}
	}
}

func (detailsView *detailsViewData) Set(tag string, value interface{}) bool {
	return detailsView.set(strings.ToLower(tag), value)
}

func (detailsView *detailsViewData) set(tag string, value interface{}) bool {
	if value == nil {
		detailsView.remove(tag)
		return true
	}

	switch tag {
	case Summary:
		switch value := value.(type) {
		case string:
			detailsView.properties[Summary] = value

		case View:
			detailsView.properties[Summary] = value

		case DataObject:
			if view := CreateViewFromObject(detailsView.Session(), value); view != nil {
				detailsView.properties[Summary] = view
			} else {
				return false
			}

		default:
			notCompatibleType(tag, value)
			return false
		}
		if detailsView.created {
			updateInnerHTML(detailsView.htmlID(), detailsView.Session())
		}

	case Expanded:
		if !detailsView.setBoolProperty(tag, value) {
			notCompatibleType(tag, value)
			return false
		}
		if detailsView.created {
			if IsDetailsExpanded(detailsView, "") {
				updateProperty(detailsView.htmlID(), "open", "", detailsView.Session())
			} else {
				removeProperty(detailsView.htmlID(), "open", detailsView.Session())
			}
		}

	default:
		return detailsView.viewsContainerData.Set(tag, value)
	}

	detailsView.propertyChangedEvent(tag)
	return true
}

func (detailsView *detailsViewData) Get(tag string) interface{} {
	return detailsView.get(strings.ToLower(tag))
}

func (detailsView *detailsViewData) get(tag string) interface{} {
	return detailsView.viewsContainerData.get(tag)
}

func (detailsView *detailsViewData) htmlTag() string {
	return "details"
}

func (detailsView *detailsViewData) htmlProperties(self View, buffer *strings.Builder) {
	detailsView.viewsContainerData.htmlProperties(self, buffer)
	if IsDetailsExpanded(detailsView, "") {
		buffer.WriteString(` open`)
	}
}

func (detailsView *detailsViewData) htmlSubviews(self View, buffer *strings.Builder) {
	if value, ok := detailsView.properties[Summary]; ok {
		switch value := value.(type) {
		case string:
			buffer.WriteString("<summary>")
			buffer.WriteString(value)
			buffer.WriteString("</summary>")

		case View:
			buffer.WriteString("<summary>")
			viewHTML(value, buffer)
			buffer.WriteString("</summary>")
		}
	}

	detailsView.viewsContainerData.htmlSubviews(self, buffer)
}

// GetDetailsSummary returns a value of the Summary property of DetailsView.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func GetDetailsSummary(view View, subviewID string) View {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.Get(Summary); value != nil {
			switch value := value.(type) {
			case string:
				return NewTextView(view.Session(), Params{Text: value})

			case View:
				return value
			}
		}
	}
	return nil
}

// IsDetailsExpanded returns a value of the Expanded property of DetailsView.
// If the second argument (subviewID) is "" then a value from the first argument (view) is returned.
func IsDetailsExpanded(view View, subviewID string) bool {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if result, ok := boolStyledProperty(view, Expanded); ok {
			return result
		}
	}
	return false
}
