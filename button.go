package rui

import "strings"

// Button represent a Button view
type Button interface {
	ListLayout
}

type buttonData struct {
	listLayoutData
}

// NewButton create new Button object and return it
func NewButton(session Session, params Params) Button {
	button := new(buttonData)
	button.init(session)
	setInitParams(button, params)
	return button
}

func newButton(session Session) View {
	return new(buttonData)
}

func (button *buttonData) init(session Session) {
	button.listLayoutData.init(session)
	button.tag = "Button"
	button.systemClass = "ruiButton"
	button.setRaw(Style, "ruiEnabledButton")
	button.setRaw(StyleDisabled, "ruiDisabledButton")
	button.setRaw(Semantics, ButtonSemantics)
	button.setRaw(TabIndex, 0)
}

func (button *buttonData) Focusable() bool {
	return true
}

func (button *buttonData) htmlSubviews(self View, buffer *strings.Builder) {
	if button.views != nil {
		for _, view := range button.views {
			view.addToCSSStyle(map[string]string{`flex`: `0 0 auto`})
			viewHTML(view, buffer, "")
		}
	}
}

// GetButtonVerticalAlign returns the vertical align of a Button subview:
// TopAlign (0), BottomAlign (1), CenterAlign (2), or StretchAlign (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetButtonVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, VerticalAlign, CenterAlign, false)
}

// GetButtonHorizontalAlign returns the vertical align of a Button subview:
// LeftAlign (0), RightAlign (1), CenterAlign (2), or StretchAlign (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetButtonHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, HorizontalAlign, CenterAlign, false)
}

// GetButtonOrientation returns the orientation of a Button subview:
// TopDownOrientation (0), StartToEndOrientation (1), BottomUpOrientation (2), or EndToStartOrientation (3)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetButtonOrientation(view View, subviewID ...string) int {
	if view = getSubview(view, subviewID); view != nil {
		if orientation, ok := valueToOrientation(view.Get(Orientation), view.Session()); ok {
			return orientation
		}

		if value := valueFromStyle(view, Orientation); value != nil {
			if orientation, ok := valueToOrientation(value, view.Session()); ok {
				return orientation
			}
		}
	}

	return StartToEndOrientation
}
