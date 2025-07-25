package rui

import "reflect"

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
	canvasView.set = canvasView.setFunc
	canvasView.remove = canvasView.removeFunc
	canvasView.changed = canvasView.propertyChanged
}

func normalizeCanvasViewTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case "draw-func":
		tag = DrawFunction
	}
	return tag
}

func (canvasView *canvasViewData) removeFunc(tag PropertyName) []PropertyName {
	if tag == DrawFunction {
		if canvasView.getRaw(DrawFunction) != nil {
			canvasView.setRaw(DrawFunction, nil)
			//canvasView.Redraw()
			return []PropertyName{DrawFunction}
		}
		return []PropertyName{}
	}

	return canvasView.viewData.removeFunc(tag)
}

func (canvasView *canvasViewData) setFunc(tag PropertyName, value any) []PropertyName {
	if tag == DrawFunction {
		switch value := value.(type) {
		case func(Canvas):
			canvasView.setRaw(DrawFunction, value)

		case string:
			canvasView.setRaw(DrawFunction, value)

		default:
			notCompatibleType(tag, value)
			return nil
		}
		return []PropertyName{DrawFunction}
	}

	return canvasView.viewData.setFunc(tag, value)
}

func (canvasView *canvasViewData) propertyChanged(tag PropertyName) {
	if tag == DrawFunction {
		canvasView.Redraw()
	} else {
		canvasView.viewData.propertyChanged(tag)
	}
}

func (canvasView *canvasViewData) htmlTag() string {
	return "canvas"
}

func (canvasView *canvasViewData) Redraw() {
	canvas := newCanvas(canvasView)
	canvas.ClearRect(0, 0, canvasView.frame.Width, canvasView.frame.Height)
	if value := canvasView.getRaw(DrawFunction); value != nil {
		switch drawer := value.(type) {
		case func(Canvas):
			drawer(canvas)

		case string:
			bind := canvasView.binding()
			if bind == nil {
				ErrorLogF(`There is no a binding object for call "%s"`, drawer)
				break
			}

			val := reflect.ValueOf(bind)
			method := val.MethodByName(drawer)
			if !method.IsValid() {
				ErrorLogF(`The "%s" method is not valid`, drawer)
				break
			}

			methodType := method.Type()
			if methodType.NumIn() == 1 && equalType(methodType.In(0), reflect.TypeOf(canvas)) {
				method.Call([]reflect.Value{reflect.ValueOf(canvas)})
			} else {
				ErrorLogF(`Unsupported prototype of "%s" method`, drawer)
			}
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
