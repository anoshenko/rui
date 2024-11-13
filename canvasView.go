package rui

// DrawFunction is the constant for "draw-function" property tag.
//
// Used by `CanvasView`.
// Property sets the draw function of `CanvasView`.
//
// Supported types: `func(Canvas)`.
const DrawFunction PropertyName = "draw-function"

// CanvasView interface of a custom draw view
type CanvasView interface {
	View

	// Redraw forces CanvasView to redraw its content
	Redraw()
}

type canvasViewData struct {
	viewData
}

// NewCanvasView creates the new custom draw view
func NewCanvasView(session Session, params Params) CanvasView {
	view := new(canvasViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newCanvasView(session Session) View {
	return new(canvasViewData)
}

// Init initialize fields of ViewsContainer by default values
func (canvasView *canvasViewData) init(session Session) {
	canvasView.viewData.init(session)
	canvasView.tag = "CanvasView"
	canvasView.normalize = normalizeCanvasViewTag
	canvasView.set = canvasViewSet
	canvasView.remove = canvasViewRemove

}

func normalizeCanvasViewTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case "draw-func":
		tag = DrawFunction
	}
	return tag
}

func canvasViewRemove(view View, tag PropertyName) []PropertyName {
	if tag == DrawFunction {
		if view.getRaw(DrawFunction) != nil {
			view.setRaw(DrawFunction, nil)
			if canvasView, ok := view.(CanvasView); ok {
				canvasView.Redraw()
			}
			return []PropertyName{DrawFunction}
		}
		return []PropertyName{}
	}

	return viewRemove(view, tag)
}

func canvasViewSet(view View, tag PropertyName, value any) []PropertyName {
	if tag == DrawFunction {
		if fn, ok := value.(func(Canvas)); ok {
			view.setRaw(DrawFunction, fn)
		} else {
			notCompatibleType(tag, value)
			return nil
		}
		if canvasView, ok := view.(CanvasView); ok {
			canvasView.Redraw()
		}
		return []PropertyName{DrawFunction}
	}

	return viewSet(view, tag, value)
}

func (canvasView *canvasViewData) htmlTag() string {
	return "canvas"
}

func (canvasView *canvasViewData) Redraw() {
	canvas := newCanvas(canvasView)
	canvas.ClearRect(0, 0, canvasView.frame.Width, canvasView.frame.Height)
	if value := canvasView.getRaw(DrawFunction); value != nil {
		if drawer, ok := value.(func(Canvas)); ok {
			drawer(canvas)
		}
	}
	canvas.finishDraw()
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
