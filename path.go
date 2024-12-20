package rui

// Path is a path interface
type Path interface {
	// MoveTo begins a new sub-path at the point specified by the given (x, y) coordinates
	MoveTo(x, y float64)

	// LineTo adds a straight line to the current sub-path by connecting
	// the sub-path's last point to the specified (x, y) coordinates
	LineTo(x, y float64)

	// ArcTo adds a circular arc to the current sub-path, using the given control points and radius.
	// The arc is automatically connected to the path's latest point with a straight line, if necessary.
	//  - x0, y0 - coordinates of the first control point;
	//  - x1, y1 - coordinates of the second control point;
	//  - radius - the arc's radius. Must be non-negative.
	ArcTo(x0, y0, x1, y1, radius float64)

	// Arc adds a circular arc to the current sub-path.
	//   - x, y - coordinates of the arc's center;
	//   - radius - the arc's radius. Must be non-negative;
	//   - startAngle - the angle at which the arc starts, measured clockwise from the positive
	//                x-axis and expressed in radians.
	//   - endAngle - the angle at which the arc ends, measured clockwise from the positive
	//                x-axis and expressed in radians.
	//   - clockwise - if true, causes the arc to be drawn clockwise between the start and end angles,
	//               otherwise - counter-clockwise
	Arc(x, y, radius, startAngle, endAngle float64, clockwise bool)

	// BezierCurveTo adds a cubic Bézier curve to the current sub-path. The starting point is
	// the latest point in the current path.
	//   - cp0x, cp0y - coordinates of the first control point;
	//   - cp1x, cp1y - coordinates of the second control point;
	//   - x, y - coordinates of the end point.
	BezierCurveTo(cp0x, cp0y, cp1x, cp1y, x, y float64)

	// QuadraticCurveTo  adds a quadratic Bézier curve to the current sub-path.
	//   - cpx, cpy - coordinates of the control point;
	//   - x, y - coordinates of the end point.
	QuadraticCurveTo(cpx, cpy, x, y float64)

	// Ellipse adds an elliptical arc to the current sub-path
	//   - x, y - coordinates of the ellipse's center;
	//   - radiusX - the ellipse's major-axis radius. Must be non-negative;
	//   - radiusY - the ellipse's minor-axis radius. Must be non-negative;
	//   - rotation - the rotation of the ellipse, expressed in radians;
	//   - startAngle - the angle at which the ellipse starts, measured clockwise
	//                from the positive x-axis and expressed in radians;
	//   - endAngle - the angle at which the ellipse ends, measured clockwise
	//	            from the positive x-axis and expressed in radians.
	//   - clockwise - if true, draws the ellipse clockwise, otherwise draws counter-clockwise
	Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, clockwise bool)

	// Close adds a straight line from the current point to the start of the current sub-path.
	// If the shape has already been closed or has only one point, this function does nothing.
	Close()

	obj() any
}

type pathData struct {
	session Session
	varName any
}

// NewPath creates a new empty Path
func (canvas *canvasData) NewPath() Path {
	path := new(pathData)
	path.session = canvas.session
	path.varName = canvas.session.createPath("")
	return path
}

func (canvas *canvasData) NewPathFromSvg(data string) Path {
	path := new(pathData)
	path.session = canvas.session
	path.varName = canvas.session.createPath(data)
	return path
}

func (path *pathData) MoveTo(x, y float64) {
	path.session.callCanvasVarFunc(path.varName, "moveTo", x, y)
}

func (path *pathData) LineTo(x, y float64) {
	path.session.callCanvasVarFunc(path.varName, "lineTo", x, y)
}

func (path *pathData) ArcTo(x0, y0, x1, y1, radius float64) {
	if radius > 0 {
		path.session.callCanvasVarFunc(path.varName, "arcTo", x0, y0, x1, y1, radius)
	}
}

func (path *pathData) Arc(x, y, radius, startAngle, endAngle float64, clockwise bool) {
	if radius > 0 {
		if !clockwise {
			path.session.callCanvasVarFunc(path.varName, "arc", x, y, radius, startAngle, endAngle, true)
		} else {
			path.session.callCanvasVarFunc(path.varName, "arc", x, y, radius, startAngle, endAngle)
		}
	}
}

func (path *pathData) BezierCurveTo(cp0x, cp0y, cp1x, cp1y, x, y float64) {
	path.session.callCanvasVarFunc(path.varName, "bezierCurveTo", cp0x, cp0y, cp1x, cp1y, x, y)
}

func (path *pathData) QuadraticCurveTo(cpx, cpy, x, y float64) {
	path.session.callCanvasVarFunc(path.varName, "quadraticCurveTo", cpx, cpy, x, y)
}

func (path *pathData) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, clockwise bool) {
	if radiusX > 0 && radiusY > 0 {
		if !clockwise {
			path.session.callCanvasVarFunc(path.varName, "ellipse", x, y, radiusX, radiusY, rotation, startAngle, endAngle, true)
		} else {
			path.session.callCanvasVarFunc(path.varName, "ellipse", x, y, radiusX, radiusY, rotation, startAngle, endAngle)
		}
	}
}

func (path *pathData) Close() {
	path.session.callCanvasVarFunc(path.varName, "closePath")
}

func (path *pathData) obj() any {
	return path.varName
}
