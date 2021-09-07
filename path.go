package rui

import (
	"strconv"
	"strings"
)

// Path is a path interface
type Path interface {
	// Reset erases the Path
	Reset()

	// MoveTo begins a new sub-path at the point specified by the given (x, y) coordinates
	MoveTo(x, y float64)

	// LineTo adds a straight line to the current sub-path by connecting
	// the sub-path's last point to the specified (x, y) coordinates
	LineTo(x, y float64)

	// ArcTo adds a circular arc to the current sub-path, using the given control points and radius.
	// The arc is automatically connected to the path's latest point with a straight line, if necessary.
	//   x0, y0 - coordinates of the first control point;
	//   x1, y1 - coordinates of the second control point;
	//   radius - the arc's radius. Must be non-negative.
	ArcTo(x0, y0, x1, y1, radius float64)

	// Arc adds a circular arc to the current sub-path.
	//   x, y - coordinates of the arc's center;
	//   radius - the arc's radius. Must be non-negative;
	//   startAngle - the angle at which the arc starts, measured clockwise from the positive
	//                x-axis and expressed in radians.
	//   endAngle - the angle at which the arc ends, measured clockwise from the positive
	//                x-axis and expressed in radians.
	//   clockwise - if true, causes the arc to be drawn clockwise between the start and end angles,
	//               otherwise - counter-clockwise
	Arc(x, y, radius, startAngle, endAngle float64, clockwise bool)

	// BezierCurveTo adds a cubic Bézier curve to the current sub-path. The starting point is
	// the latest point in the current path.
	//   cp0x, cp0y - coordinates of the first control point;
	//   cp1x, cp1y - coordinates of the second control point;
	//   x, y - coordinates of the end point.
	BezierCurveTo(cp0x, cp0y, cp1x, cp1y, x, y float64)

	// QuadraticCurveTo  adds a quadratic Bézier curve to the current sub-path.
	//   cpx, cpy - coordinates of the control point;
	//   x, y - coordinates of the end point.
	QuadraticCurveTo(cpx, cpy, x, y float64)

	// Ellipse adds an elliptical arc to the current sub-path
	//   x, y - coordinates of the ellipse's center;
	//   radiusX - the ellipse's major-axis radius. Must be non-negative;
	//   radiusY - the ellipse's minor-axis radius. Must be non-negative;
	//   rotation - the rotation of the ellipse, expressed in radians;
	//   startAngle - the angle at which the ellipse starts, measured clockwise
	//                from the positive x-axis and expressed in radians;
	//   endAngle - the angle at which the ellipse ends, measured clockwise
	//	            from the positive x-axis and expressed in radians.
	//   clockwise - if true, draws the ellipse clockwise, otherwise draws counter-clockwise
	Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, clockwise bool)

	// Close adds a straight line from the current point to the start of the current sub-path.
	// If the shape has already been closed or has only one point, this function does nothing.
	Close()

	scriptText() string
}

type pathData struct {
	script strings.Builder
}

// NewPath creates a new empty Path
func NewPath() Path {
	path := new(pathData)
	path.script.Grow(4096)
	path.script.WriteString("\nctx.beginPath();")
	return path
}

func (path *pathData) Reset() {
	path.script.Reset()
	path.script.WriteString("\nctx.beginPath();")
}

func (path *pathData) MoveTo(x, y float64) {
	path.script.WriteString("\nctx.moveTo(")
	path.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	path.script.WriteString(");")
}

func (path *pathData) LineTo(x, y float64) {
	path.script.WriteString("\nctx.lineTo(")
	path.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	path.script.WriteString(");")
}

func (path *pathData) ArcTo(x0, y0, x1, y1, radius float64) {
	if radius > 0 {
		path.script.WriteString("\nctx.arcTo(")
		path.script.WriteString(strconv.FormatFloat(x0, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(y0, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(x1, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(y1, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(radius, 'g', -1, 64))
		path.script.WriteString(");")
	}
}

func (path *pathData) Arc(x, y, radius, startAngle, endAngle float64, clockwise bool) {
	if radius > 0 {
		path.script.WriteString("\nctx.arc(")
		path.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(radius, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(startAngle, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(endAngle, 'g', -1, 64))
		if !clockwise {
			path.script.WriteString(",true);")
		} else {
			path.script.WriteString(");")
		}
	}
}

func (path *pathData) BezierCurveTo(cp0x, cp0y, cp1x, cp1y, x, y float64) {
	path.script.WriteString("\nctx.bezierCurveTo(")
	path.script.WriteString(strconv.FormatFloat(cp0x, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(cp0y, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(cp1x, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(cp1y, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	path.script.WriteString(");")
}

func (path *pathData) QuadraticCurveTo(cpx, cpy, x, y float64) {
	path.script.WriteString("\nctx.quadraticCurveTo(")
	path.script.WriteString(strconv.FormatFloat(cpx, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(cpy, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	path.script.WriteRune(',')
	path.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	path.script.WriteString(");")
}

func (path *pathData) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, clockwise bool) {
	if radiusX > 0 && radiusY > 0 {
		path.script.WriteString("\nctx.ellipse(")
		path.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(radiusX, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(radiusY, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(rotation, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(startAngle, 'g', -1, 64))
		path.script.WriteRune(',')
		path.script.WriteString(strconv.FormatFloat(endAngle, 'g', -1, 64))
		if !clockwise {
			path.script.WriteString(",true);")
		} else {
			path.script.WriteString(");")
		}
	}
}

func (path *pathData) Close() {
	path.script.WriteString("\nctx.close();")
}

func (path *pathData) scriptText() string {
	return path.script.String()
}
