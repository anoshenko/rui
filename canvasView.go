package rui

import "strings"

// DrawFunction is the constant for the "draw-function" property tag.
// The "draw-function" property sets the draw function of CanvasView.
// The function should have the following format: func(Canvas)
const DrawFunction = "draw-function"

// CanvasView interface of a custom draw view
type CanvasView interface {
	View
	Redraw()
}

type canvasViewData struct {
	viewData
	drawer func(Canvas)
}

// NewCanvasView creates the new custom draw view
func NewCanvasView(session Session, params Params) CanvasView {
	view := new(canvasViewData)
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newCanvasView(session Session) View {
	return NewCanvasView(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (canvasView *canvasViewData) Init(session Session) {
	canvasView.viewData.Init(session)
	canvasView.tag = "CanvasView"
}

func (canvasView *canvasViewData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case "draw-func":
		tag = DrawFunction
	}
	return tag
}

func (canvasView *canvasViewData) Remove(tag string) {
	canvasView.remove(canvasView.normalizeTag(tag))
}

func (canvasView *canvasViewData) remove(tag string) {
	if tag == DrawFunction {
		canvasView.drawer = nil
		canvasView.Redraw()
		canvasView.propertyChangedEvent(tag)
	} else {
		canvasView.viewData.remove(tag)
	}
}

func (canvasView *canvasViewData) Set(tag string, value interface{}) bool {
	return canvasView.set(canvasView.normalizeTag(tag), value)
}

func (canvasView *canvasViewData) set(tag string, value interface{}) bool {
	if tag == DrawFunction {
		if value == nil {
			canvasView.drawer = nil
		} else if fn, ok := value.(func(Canvas)); ok {
			canvasView.drawer = fn
		} else {
			notCompatibleType(tag, value)
			return false
		}
		canvasView.Redraw()
		canvasView.propertyChangedEvent(tag)
		return true
	}

	return canvasView.viewData.set(tag, value)
}

func (canvasView *canvasViewData) Get(tag string) interface{} {
	return canvasView.get(canvasView.normalizeTag(tag))
}

func (canvasView *canvasViewData) get(tag string) interface{} {
	if tag == DrawFunction {
		return canvasView.drawer
	}
	return canvasView.viewData.get(tag)
}

func (canvasView *canvasViewData) htmlTag() string {
	return "canvas"
}

func (canvasView *canvasViewData) Redraw() {
	if canvasView.drawer != nil {
		canvas := newCanvas(canvasView)
		canvas.ClearRect(0, 0, canvasView.frame.Width, canvasView.frame.Height)
		if canvasView.drawer != nil {
			canvasView.drawer(canvas)
		}
		canvasView.session.runScript(canvas.finishDraw())
	}
}

func (canvasView *canvasViewData) onResize(self View, x, y, width, height float64) {
	canvasView.viewData.onResize(self, x, y, width, height)
	canvasView.Redraw()
}

// RedrawCanvasView finds CanvasView with canvasViewID and redraws it
func RedrawCanvasView(rootView View, canvasViewID string) {
	if canvas := CanvasViewByID(rootView, canvasViewID); canvas != nil {
		canvas.Redraw()
	}
}
