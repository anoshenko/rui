package rui

// Button represent a Button view
type Button interface {
	CustomView
}

type buttonData struct {
	CustomViewData
}

// NewButton create new Button object and return it
func NewButton(session Session, params Params) Button {
	button := new(buttonData)
	InitCustomView(button, "Button", session, params)
	return button
}

func newButton(session Session) View {
	return NewButton(session, nil)
}

func (button *buttonData) CreateSuperView(session Session) View {
	return NewListLayout(session, Params{
		Semantics:       ButtonSemantics,
		Style:           "ruiButton",
		StyleDisabled:   "ruiDisabledButton",
		HorizontalAlign: CenterAlign,
		VerticalAlign:   CenterAlign,
		Orientation:     StartToEndOrientation,
		TabIndex:        0,
	})
}

func (button *buttonData) Focusable() bool {
	return true
}
