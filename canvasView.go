package rui

import "strings"

// DrawFunction is the constant for "draw-function" property tag.
//
// Used by `CanvasView`.
// Property sets the draw function of `CanvasView`.
//
// Supported types: `func(Canvas)`.
const DrawFunction = "draw-function"

// CanvasView interface of a custom draw view
type CanvasView interface {
	View

	// Redraw forces CanvasView to redraw its content
	Redraw()
}

type canvasViewData struct {
	viewData
	drawer func(Canvas)
}

// NewCanvasView creates the new custom draw view
func NewCanvasView(session Session, params Params) CanvasView {
	view := new(canvasViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newCanvasView(session Session) View {
	return NewCanvasView(session, nil)
}

// Init initialize fields of ViewsContainer by default values
func (canvasView *canvasViewData) init(session Session) {
	canvasView.viewData.init(session)
	canvasView.tag = "CanvasView"
}

func (canvasView *canvasViewData) String() string {
	return getViewString(canvasView, nil)
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

func (canvasView *canvasViewData) Set(tag string, value any) bool {
	return canvasView.set(canvasView.normalizeTag(tag), value)
}

func (canvasView *canvasViewData) set(tag string, value any) bool {
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

func (canvasView *canvasViewData) Get(tag string) any {
	return canvasView.get(canvasView.normalizeTag(tag))
}

func (canvasView *canvasViewData) get(tag string) any {
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
		canvas.finishDraw()
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
