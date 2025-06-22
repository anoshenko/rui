package rui

import "strings"

// Constants for [DetailsView] specific properties and events
const (
	// Summary is the constant for "summary" property tag.
	//
	// Used by DetailsView.
	// The content of this property is used as the label for the disclosure widget.
	//
	// Supported types:
	//   - string - Summary as a text.
	//   - View - Summary as a view, in this case it can be quite complex if needed.
	Summary PropertyName = "summary"

	// Expanded is the constant for "expanded" property tag.
	//
	// Used by DetailsView.
	// Controls the content expanded state of the details view. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Content is visible.
	//   - false, 0, "false", "no", "off", or "0" - Content is collapsed (hidden).
	Expanded PropertyName = "expanded"

	// HideSummaryMarker is the constant for "hide-summary-marker" property tag.
	//
	// Used by DetailsView.
	// Allows you to hide the summary marker (▶︎). Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - The summary marker is hidden.
	//   - false, 0, "false", "no", "off", or "0" - The summary marker is displayed (default value).
	HideSummaryMarker PropertyName = "hide-summary-marker"
)

// DetailsView represent a DetailsView view, which is a collapsible container of views
type DetailsView interface {
	ViewsContainer
}

type detailsViewData struct {
	viewsContainerData
}

// NewDetailsView create new DetailsView object and return it
func NewDetailsView(session Session, params Params) DetailsView {
	view := new(detailsViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newDetailsView(session Session) View {
	return new(detailsViewData)
}

// Init initialize fields of DetailsView by default values
func (detailsView *detailsViewData) init(session Session) {
	detailsView.viewsContainerData.init(session)
	detailsView.tag = "DetailsView"
	detailsView.set = detailsView.setFunc
	detailsView.changed = detailsView.propertyChanged
	//detailsView.systemClass = "ruiDetailsView"
}

func (detailsView *detailsViewData) Views() []View {
	views := detailsView.viewsContainerData.Views()
	if summary := detailsView.Get(Summary); summary != nil {
		switch summary := summary.(type) {
		case View:
			return append([]View{summary}, views...)
		}
	}
	return views
}

func (detailsView *detailsViewData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Summary:
		switch value := value.(type) {
		case string:
			detailsView.setRaw(Summary, value)

		case View:
			detailsView.setRaw(Summary, value)
			value.setParentID(detailsView.htmlID())

		case DataObject:
			if view := CreateViewFromObject(detailsView.Session(), value); view != nil {
				detailsView.setRaw(Summary, view)
				view.setParentID(detailsView.htmlID())
			} else {
				return nil
			}

		default:
			notCompatibleType(tag, value)
			return nil
		}
		return []PropertyName{tag}
	}

	return detailsView.viewsContainerData.setFunc(tag, value)
}

func (detailsView *detailsViewData) propertyChanged(tag PropertyName) {
	switch tag {
	case Summary, HideSummaryMarker:
		updateInnerHTML(detailsView.htmlID(), detailsView.Session())

	case Expanded:
		if IsDetailsExpanded(detailsView) {
			detailsView.Session().updateProperty(detailsView.htmlID(), "open", "")
		} else {
			detailsView.Session().removeProperty(detailsView.htmlID(), "open")
		}

	case NotTranslate:
		updateInnerHTML(detailsView.htmlID(), detailsView.Session())

	default:
		detailsView.viewsContainerData.propertyChanged(tag)
	}
}

func (detailsView *detailsViewData) htmlTag() string {
	return "details"
}

func (detailsView *detailsViewData) htmlProperties(self View, buffer *strings.Builder) {
	detailsView.viewsContainerData.htmlProperties(self, buffer)
	buffer.WriteString(` ontoggle="detailsEvent(this)"`)
	if IsDetailsExpanded(detailsView) {
		buffer.WriteString(` open`)
	}
}

func (detailsView *detailsViewData) htmlSubviews(self View, buffer *strings.Builder) {
	summary := false
	hidden := IsSummaryMarkerHidden(detailsView)

	if value, ok := detailsView.properties[Summary]; ok {

		switch value := value.(type) {
		case string:
			if !GetNotTranslate(detailsView) {
				value, _ = detailsView.session.GetString(value)
			}
			if hidden {
				buffer.WriteString(`<summary class="hiddenMarker">`)
			} else {
				buffer.WriteString("<summary>")
			}
			buffer.WriteString(value)
			buffer.WriteString("</summary>")
			summary = true

		case View:
			if hidden {
				buffer.WriteString(`<summary class="hiddenMarker">`)
				viewHTML(value, buffer, "")
				buffer.WriteString("</summary>")
			} else if value.htmlTag() == "div" {
				viewHTML(value, buffer, "summary")
			} else {
				buffer.WriteString(`<summary><div style="display: inline-block;">`)
				viewHTML(value, buffer, "")
				buffer.WriteString("</div></summary>")
			}
			summary = true
		}
	}

	if !summary {
		if hidden {
			buffer.WriteString(`<summary class="hiddenMarker"></summary>`)
		} else {
			buffer.WriteString("<summary></summary>")
		}
	}

	detailsView.viewsContainerData.htmlSubviews(self, buffer)
}

func (detailsView *detailsViewData) handleCommand(self View, command PropertyName, data DataObject) bool {
	if command == "details-open" {
		if n, ok := dataIntProperty(data, "open"); ok {
			detailsView.properties[Expanded] = (n != 0)
			if listener, ok := detailsView.changeListener[Expanded]; ok {
				listener.Run(detailsView, Expanded)
			}
		}
		return true
	}
	return detailsView.viewsContainerData.handleCommand(self, command, data)
}

// GetDetailsSummary returns a value of the Summary property of DetailsView.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetDetailsSummary(view View, subviewID ...string) View {
	if view = getSubview(view, subviewID); view != nil {
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
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func IsDetailsExpanded(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, Expanded, false)
}

// IsDetailsExpanded returns a value of the HideSummaryMarker property of DetailsView.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func IsSummaryMarkerHidden(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, HideSummaryMarker, false)
}
