package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// MiterJoin - Connected segments are joined by extending their outside edges
	// to connect at a single point, with the effect of filling an additional
	// lozenge-shaped area. This setting is affected by the miterLimit property
	MiterJoin = 0
	// RoundJoin - rounds off the corners of a shape by filling an additional sector
	// of disc centered at the common endpoint of connected segments.
	// The radius for these rounded corners is equal to the line width.
	RoundJoin = 1
	// BevelJoin - Fills an additional triangular area between the common endpoint
	// of connected segments, and the separate outside rectangular corners of each segment.
	BevelJoin = 2

	// ButtCap - the ends of lines are squared off at the endpoints. Default value.
	ButtCap = 0
	// RoundCap - the ends of lines are rounded.
	RoundCap = 1
	// SquareCap - the ends of lines are squared off by adding a box with an equal width
	// and half the height of the line's thickness.
	SquareCap = 2

	// AlphabeticBaseline - the text baseline is the normal alphabetic baseline. Default value.
	AlphabeticBaseline = 0
	// TopBaseline - the text baseline is the top of the em square.
	TopBaseline = 1
	// MiddleBaseline - the text baseline is the middle of the em square.
	MiddleBaseline = 2
	// BottomBaseline - the text baseline is the bottom of the bounding box.
	// This differs from the ideographic baseline in that the ideographic baseline doesn't consider descenders.
	BottomBaseline = 3
	// HangingBaseline - the text baseline is the hanging baseline. (Used by Tibetan and other Indic scripts.)
	HangingBaseline = 4
	// IdeographicBaseline - the text baseline is the ideographic baseline; this is
	// the bottom of the body of the characters, if the main body of characters protrudes
	// beneath the alphabetic baseline. (Used by Chinese, Japanese, and Korean scripts.)
	IdeographicBaseline = 5

	// StartAlign - the text is aligned at the normal start of the line (left-aligned
	// for left-to-right locales, right-aligned for right-to-left locales).
	StartAlign = 3
	// EndAlign - the text is aligned at the normal end of the line (right-aligned
	// for left-to-right locales, left-aligned for right-to-left locales).
	EndAlign = 4
)

// GradientPoint defined by an offset and a color, to a linear or radial gradient
type GradientPoint struct {
	// Offset - a number between 0 and 1, inclusive, representing the position of the color stop
	Offset float64
	// Color - the color of the stop
	Color Color
}

// FontParams defined optionally font properties
type FontParams struct {
	// Italic - if true then a font is italic
	Italic bool
	// SmallCaps - if true then a font uses small-caps glyphs
	SmallCaps bool
	// Weight - a font weight. Valid values: 0...9, there
	//   0 - a weight does not specify;
	//   1 - a minimal weight;
	//   4 - a normal weight;
	//   7 - a bold weight;
	//   9 - a maximal weight.
	Weight int
	// LineHeight - the height (relative to the font size of the element itself) of a line box.
	LineHeight SizeUnit
}

// TextMetrics is the result of the Canvas.TextMetrics function
type TextMetrics struct {
	// Width is the calculated width of a segment of inline text in pixels
	Width float64
	// Ascent is the distance from the horizontal baseline to the top of the bounding rectangle used to render the text, in pixels.
	Ascent float64
	// Descent is the distance from the horizontal baseline to the bottom of the bounding rectangle used to render the text, in pixels.
	Descent float64
	// Left is the distance to the left side of the bounding rectangle of the given text, in  pixels;
	// positive numbers indicating a distance going left from the given alignment point.
	Left float64
	// Right is the distance to the right side of the bounding rectangle of the given text, CSS pixels.
	Right float64
}

// Canvas is a drawing interface
type Canvas interface {
	// View return the view for the drawing
	View() CanvasView
	// Width returns the width in pixels of the canvas area
	Width() float64
	// Height returns the height in pixels of the canvas area
	Height() float64

	// Save saves the entire state of the canvas by pushing the current state onto a stack.
	Save()
	// Restore restores the most recently saved canvas state by popping the top entry
	// in the drawing state stack. If there is no saved state, this method does nothing.
	Restore()

	// ClipPath turns the rectangle into the current clipping region. It replaces any previous clipping region.
	ClipRect(x, y, width, height float64)
	// ClipPath turns the path into the current clipping region. It replaces any previous clipping region.
	ClipPath(path Path)

	// SetScale adds a scaling transformation to the canvas units horizontally and/or vertically.
	//   x - scaling factor in the horizontal direction. A negative value flips pixels across
	//       the vertical axis. A value of 1 results in no horizontal scaling;
	//   y - scaling factor in the vertical direction. A negative value flips pixels across
	//       the horizontal axis. A value of 1 results in no vertical scaling.
	SetScale(x, y float64)

	// SetTranslation adds a translation transformation to the current matrix.
	//   x - distance to move in the horizontal direction. Positive values are to the right, and negative to the left;
	//   y - distance to move in the vertical direction. Positive values are down, and negative are up.
	SetTranslation(x, y float64)

	// SetRotation adds a rotation to the transformation matrix.
	//   angle - the rotation angle, clockwise in radians
	SetRotation(angle float64)

	// SetTransformation multiplies the current transformation with the matrix described by the arguments
	// of this method. This lets you scale, rotate, translate (move), and skew the context.
	// The transformation matrix is described by:
	// ⎡ xScale xSkew  dx ⎤
	// ⎢ ySkew  yScale dy ⎥
	// ⎣   0      0     1 ⎦
	//   xScale, yScale - horizontal and vertical scaling. A value of 1 results in no scaling;
	//   xSkew, ySkew - horizontal and vertical skewing;
	//   dx, dy - horizontal and vertical translation (moving).
	SetTransformation(xScale, yScale, xSkew, ySkew, dx, dy float64)

	// ResetTransformation resets the current transform to the identity matrix
	ResetTransformation()

	// SetSolidColorFillStyle sets the color to use inside shapes
	SetSolidColorFillStyle(color Color)

	// SetSolidColorStrokeStyle sets color to use for the strokes (outlines) around shapes
	SetSolidColorStrokeStyle(color Color)

	// SetLinearGradientFillStyle sets a gradient along the line connecting two given coordinates to use inside shapes
	//   x0, y0 - coordinates of the start point;
	//   x1, y1 - coordinates of the end point;
	//   startColor, endColor - the start and end color
	//   stopPoints - the array of stop points
	SetLinearGradientFillStyle(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint)

	// SetLinearGradientStrokeStyle sets a gradient along the line connecting two given coordinates to use for the strokes (outlines) around shapes
	//   x0, y0 - coordinates of the start point;
	//   x1, y1 - coordinates of the end point;
	//   color0, color1 - the start and end color
	//   stopPoints - the array of stop points
	SetLinearGradientStrokeStyle(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint)

	// SetRadialGradientFillStyle sets a radial gradient using the size and coordinates of two circles
	// to use inside shapes
	//   x0, y0 - coordinates of the center of the start circle;
	//   r0 - the radius of the start circle;
	//   x1, y1 - coordinates the center of the end circle;
	//   r1 - the radius of the end circle;
	//   color0, color1 - the start and end color
	//   stopPoints - the array of stop points
	SetRadialGradientFillStyle(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint)

	// SetRadialGradientStrokeStyle sets a radial gradient using the size and coordinates of two circles
	// to use for the strokes (outlines) around shapes
	//   x0, y0 - coordinates of the center of the start circle;
	//   r0 - the radius of the start circle;
	//   x1, y1 - coordinates the center of the end circle;
	//   r1 - the radius of the end circle;
	//   color0, color1 - the start and end color
	//   stopPoints - the array of stop points
	SetRadialGradientStrokeStyle(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint)

	// SetImageFillStyle set the image as the filling pattern.
	//   repeate - indicating how to repeat the pattern's image. Possible values are:
	//     NoRepeat (0) - neither direction,
	//     RepeatXY (1) - both directions,
	//     RepeatX (2) - horizontal only,
	//     RepeatY (3) - vertical only.
	SetImageFillStyle(image Image, repeat int)

	// SetLineWidth the line width, in coordinate space units. Zero, negative, Infinity, and NaN values are ignored.
	SetLineWidth(width float64)

	// SetLineJoin sets the shape used to join two line segments where they meet.
	// Valid values: MiterJoin (0), RoundJoin (1), BevelJoin (2). All other values are ignored.
	SetLineJoin(join int)

	// SetLineJoin sets the shape used to draw the end points of lines.
	// Valid values: ButtCap (0), RoundCap (1), SquareCap (2). All other values are ignored.
	SetLineCap(cap int)

	// SetLineDash sets the line dash pattern used when stroking lines.
	// dash - an array of values that specify alternating lengths of lines and gaps which describe the pattern.
	// offset - the line dash offset
	SetLineDash(dash []float64, offset float64)

	// SetFont sets the current text style to use when drawing text
	SetFont(name string, size SizeUnit)
	// SetFontWithParams sets the current text style to use when drawing text
	SetFontWithParams(name string, size SizeUnit, params FontParams)

	// TextWidth calculates metrics of the text drawn by a given font
	TextMetrics(text string, fontName string, fontSize SizeUnit, fontParams FontParams) TextMetrics

	// SetTextBaseline sets the current text baseline used when drawing text. Valid values:
	// AlphabeticBaseline (0), TopBaseline (1), MiddleBaseline (2), BottomBaseline (3),
	// HangingBaseline (4), and IdeographicBaseline (5). All other values are ignored.
	SetTextBaseline(baseline int)

	// SetTextAlign sets the current text alignment used when drawing text. Valid values:
	// LeftAlign (0), RightAlign (1), CenterAlign (2), StartAlign (3), and EndAlign(4). All other values are ignored.
	SetTextAlign(align int)

	// SetShadow sets shadow parameters:
	//   offsetX, offsetY - the distance that shadows will be offset horizontally and vertically;
	//   blur - the amount of blur applied to shadows. Must be non-negative;
	//   color - the color of shadows.
	SetShadow(offsetX, offsetY, blur float64, color Color)
	// ResetShadow sets shadow parameters to default values (invisible shadow)
	ResetShadow()

	// ClearRect erases the pixels in a rectangular area by setting them to transparent black
	ClearRect(x, y, width, height float64)
	// FillRect draws a rectangle that is filled according to the current FillStyle.
	FillRect(x, y, width, height float64)
	// StrokeRect draws a rectangle that is stroked (outlined) according to the current strokeStyle
	// and other context settings
	StrokeRect(x, y, width, height float64)
	// FillAndStrokeRect draws a rectangle that is filled according to the current FillStyle and
	// is stroked (outlined) according to the current strokeStyle and other context settings
	FillAndStrokeRect(x, y, width, height float64)

	// FillRoundedRect draws a rounded rectangle that is filled according to the current FillStyle.
	FillRoundedRect(x, y, width, height, r float64)
	// StrokeRoundedRect draws a rounded rectangle that is stroked (outlined) according
	// to the current strokeStyle and other context settings
	StrokeRoundedRect(x, y, width, height, r float64)
	// FillAndStrokeRoundedRect draws a rounded rectangle that is filled according to the current FillStyle
	// and is stroked (outlined) according to the current strokeStyle and other context settings
	FillAndStrokeRoundedRect(x, y, width, height, r float64)

	// FillEllipse draws a ellipse that is filled according to the current FillStyle.
	//   x, y - coordinates of the ellipse's center;
	//   radiusX - the ellipse's major-axis radius. Must be non-negative;
	//   radiusY - the ellipse's minor-axis radius. Must be non-negative;
	//   rotation - the rotation of the ellipse, expressed in radians.
	FillEllipse(x, y, radiusX, radiusY, rotation float64)
	// StrokeRoundedRect draws a ellipse that is stroked (outlined) according
	// to the current strokeStyle and other context settings
	StrokeEllipse(x, y, radiusX, radiusY, rotation float64)
	// FillAndStrokeEllipse draws a ellipse that is filled according to the current FillStyle
	// and is stroked (outlined) according to the current strokeStyle and other context settings
	FillAndStrokeEllipse(x, y, radiusX, radiusY, rotation float64)

	// FillPath draws a path that is filled according to the current FillStyle.
	FillPath(path Path)
	// StrokePath draws a path that is stroked (outlined) according to the current strokeStyle
	// and other context settings
	StrokePath(path Path)
	// FillAndStrokeRect draws a path that is filled according to the current FillStyle and
	// is stroked (outlined) according to the current strokeStyle and other context settings
	FillAndStrokePath(path Path)

	// DrawLine draws a line according to the current strokeStyle and other context settings
	DrawLine(x0, y0, x1, y1 float64)

	// FillText draws a text string at the specified coordinates, filling the string's characters
	// with the current FillStyle
	FillText(x, y float64, text string)
	// StrokeText strokes — that is, draws the outlines of — the characters of a text string
	// at the specified coordinates
	StrokeText(x, y float64, text string)

	// DrawImage draws the image at the (x, y) position
	DrawImage(x, y float64, image Image)
	// DrawImageInRect draws the image in the rectangle (x, y, width, height), scaling in height and width if necessary
	DrawImageInRect(x, y, width, height float64, image Image)
	// DrawImageFragment draws the frament (described by srcX, srcY, srcWidth, srcHeight) of image
	// in the rectangle (dstX, dstY, dstWidth, dstHeight), scaling in height and width if necessary
	DrawImageFragment(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight float64, image Image)

	finishDraw() string
}

type canvasData struct {
	view   CanvasView
	script strings.Builder
}

func newCanvas(view CanvasView) Canvas {
	canvas := new(canvasData)
	canvas.view = view
	canvas.script.Grow(4096)
	canvas.script.WriteString(`const canvas = document.getElementById('`)
	canvas.script.WriteString(view.htmlID())
	canvas.script.WriteString(`');
const ctx = canvas.getContext('2d');
const dpr = window.devicePixelRatio || 1;
var gradient;
var path;
var img;
ctx.canvas.width = dpr * canvas.clientWidth;
ctx.canvas.height = dpr * canvas.clientHeight;
ctx.scale(dpr, dpr);`)
	/*
	   	canvas.script.WriteString(strconv.FormatFloat(view.canvasWidth(), 'g', -1, 64))
	   	canvas.script.WriteString(`;
	   ctx.canvas.height = dpr * `)
	   	canvas.script.WriteString(strconv.FormatFloat(view.canvasHeight(), 'g', -1, 64))
	   	canvas.script.WriteString(";\nctx.scale(dpr, dpr);")
	*/
	return canvas
}

func (canvas *canvasData) finishDraw() string {
	canvas.script.WriteString("\n")
	return canvas.script.String()
}

func (canvas *canvasData) View() CanvasView {
	return canvas.view
}

func (canvas *canvasData) Width() float64 {
	if canvas.view != nil {
		return canvas.view.Frame().Width
	}
	return 0
}

func (canvas *canvasData) Height() float64 {
	if canvas.view != nil {
		return canvas.view.Frame().Height
	}
	return 0
}

func (canvas *canvasData) Save() {
	canvas.script.WriteString("\nctx.save();")
}

func (canvas *canvasData) Restore() {
	canvas.script.WriteString("\nctx.restore();")
}

func (canvas *canvasData) ClipRect(x, y, width, height float64) {
	canvas.script.WriteString("\nctx.beginPath();\nctx.rect(")
	canvas.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(width, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(height, 'g', -1, 64))
	canvas.script.WriteString(");\nctx.clip();")
}

func (canvas *canvasData) ClipPath(path Path) {
	canvas.script.WriteString(path.scriptText())
	canvas.script.WriteString("\nctx.clip();")
}

func (canvas *canvasData) SetScale(x, y float64) {
	canvas.script.WriteString("\nctx.scale(")
	canvas.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	canvas.script.WriteString(");")
}

func (canvas *canvasData) SetTranslation(x, y float64) {
	canvas.script.WriteString("\nctx.translate(")
	canvas.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	canvas.script.WriteString(");")
}

func (canvas *canvasData) SetRotation(angle float64) {
	canvas.script.WriteString("\nctx.rotate(")
	canvas.script.WriteString(strconv.FormatFloat(angle, 'g', -1, 64))
	canvas.script.WriteString(");")
}

func (canvas *canvasData) SetTransformation(xScale, yScale, xSkew, ySkew, dx, dy float64) {
	canvas.script.WriteString("\nctx.transform(")
	canvas.script.WriteString(strconv.FormatFloat(xScale, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(ySkew, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(xSkew, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(yScale, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(dx, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(dy, 'g', -1, 64))
	canvas.script.WriteString(");")
}

func (canvas *canvasData) ResetTransformation() {
	canvas.script.WriteString("\nctx.resetTransform();\nctx.scale(dpr, dpr);")
}

func (canvas *canvasData) SetSolidColorFillStyle(color Color) {
	canvas.script.WriteString("\nctx.fillStyle = \"")
	canvas.script.WriteString(color.cssString())
	canvas.script.WriteString(`";`)
}

func (canvas *canvasData) SetSolidColorStrokeStyle(color Color) {
	canvas.script.WriteString("\nctx.strokeStyle = \"")
	canvas.script.WriteString(color.cssString())
	canvas.script.WriteString(`";`)
}

func (canvas *canvasData) setLinearGradient(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint) {
	canvas.script.WriteString("\ngradient = ctx.createLinearGradient(")
	canvas.script.WriteString(strconv.FormatFloat(x0, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y0, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(x1, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y1, 'g', -1, 64))
	canvas.script.WriteString(");\ngradient.addColorStop(0, '")
	canvas.script.WriteString(color0.cssString())
	canvas.script.WriteString("');")

	for _, point := range stopPoints {
		if point.Offset >= 0 && point.Offset <= 1 {
			canvas.script.WriteString("\ngradient.addColorStop(")
			canvas.script.WriteString(strconv.FormatFloat(point.Offset, 'g', -1, 64))
			canvas.script.WriteString(", '")
			canvas.script.WriteString(point.Color.cssString())
			canvas.script.WriteString("');")
		}
	}

	canvas.script.WriteString("\ngradient.addColorStop(1, '")
	canvas.script.WriteString(color1.cssString())
	canvas.script.WriteString("');")
}

func (canvas *canvasData) SetLinearGradientFillStyle(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint) {
	canvas.setLinearGradient(x0, y0, color0, x1, y1, color1, stopPoints)
	canvas.script.WriteString("\nctx.fillStyle = gradient;")
}

func (canvas *canvasData) SetLinearGradientStrokeStyle(x0, y0 float64, color0 Color, x1, y1 float64, color1 Color, stopPoints []GradientPoint) {
	canvas.setLinearGradient(x0, y0, color0, x1, y1, color1, stopPoints)
	canvas.script.WriteString("\nctx.strokeStyle = gradient;")
}

func (canvas *canvasData) setRadialGradient(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint) {
	canvas.script.WriteString("\ngradient = ctx.createRadialGradient(")
	canvas.script.WriteString(strconv.FormatFloat(x0, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y0, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(r0, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(x1, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y1, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(r1, 'g', -1, 64))
	canvas.script.WriteString(");\ngradient.addColorStop(0, '")
	canvas.script.WriteString(color0.cssString())
	canvas.script.WriteString("');")

	for _, point := range stopPoints {
		if point.Offset >= 0 && point.Offset <= 1 {
			canvas.script.WriteString("\ngradient.addColorStop(")
			canvas.script.WriteString(strconv.FormatFloat(point.Offset, 'g', -1, 64))
			canvas.script.WriteString(", '")
			canvas.script.WriteString(point.Color.cssString())
			canvas.script.WriteString("');")
		}
	}

	canvas.script.WriteString("\ngradient.addColorStop(1, '")
	canvas.script.WriteString(color1.cssString())
	canvas.script.WriteString("');")
}

func (canvas *canvasData) SetRadialGradientFillStyle(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint) {
	canvas.setRadialGradient(x0, y0, r0, color0, x1, y1, r1, color1, stopPoints)
	canvas.script.WriteString("\nctx.fillStyle = gradient;")
}

func (canvas *canvasData) SetRadialGradientStrokeStyle(x0, y0, r0 float64, color0 Color, x1, y1, r1 float64, color1 Color, stopPoints []GradientPoint) {
	canvas.setRadialGradient(x0, y0, r0, color0, x1, y1, r1, color1, stopPoints)
	canvas.script.WriteString("\nctx.strokeStyle = gradient;")
}

func (canvas *canvasData) SetImageFillStyle(image Image, repeat int) {
	if image == nil || image.LoadingStatus() != ImageReady {
		return
	}

	var repeatText string
	switch repeat {
	case NoRepeat:
		repeatText = "no-repeat"

	case RepeatXY:
		repeatText = "repeat"

	case RepeatX:
		repeatText = "repeat-x"

	case RepeatY:
		repeatText = "repeat-y"

	default:
		return
	}

	canvas.script.WriteString("\nimg = images.get('")
	canvas.script.WriteString(image.URL())
	canvas.script.WriteString("');\nif (img) {\nctx.fillStyle = ctx.createPattern(img,'")
	canvas.script.WriteString(repeatText)
	canvas.script.WriteString("');\n}")
}

func (canvas *canvasData) SetLineWidth(width float64) {
	if width > 0 {
		canvas.script.WriteString("\nctx.lineWidth = '")
		canvas.script.WriteString(strconv.FormatFloat(width, 'g', -1, 64))
		canvas.script.WriteString("';")
	}
}

func (canvas *canvasData) SetLineJoin(join int) {
	switch join {
	case MiterJoin:
		canvas.script.WriteString("\nctx.lineJoin = 'miter';")

	case RoundJoin:
		canvas.script.WriteString("\nctx.lineJoin = 'round';")

	case BevelJoin:
		canvas.script.WriteString("\nctx.lineJoin = 'bevel';")
	}
}

func (canvas *canvasData) SetLineCap(cap int) {
	switch cap {
	case ButtCap:
		canvas.script.WriteString("\nctx.lineCap = 'butt';")

	case RoundCap:
		canvas.script.WriteString("\nctx.lineCap = 'round';")

	case SquareCap:
		canvas.script.WriteString("\nctx.lineCap = 'square';")
	}
}

func (canvas *canvasData) SetLineDash(dash []float64, offset float64) {
	canvas.script.WriteString("\nctx.setLineDash([")
	for i, d := range dash {
		if i > 0 {
			canvas.script.WriteString(",")
		}
		canvas.script.WriteString(strconv.FormatFloat(d, 'g', -1, 64))
	}

	canvas.script.WriteString("]);")
	if offset >= 0 {
		canvas.script.WriteString("\nctx.lineDashOffset = '")
		canvas.script.WriteString(strconv.FormatFloat(offset, 'g', -1, 64))
		canvas.script.WriteString("';")
	}
}

func (canvas *canvasData) writeFont(name string, script *strings.Builder) {
	names := strings.Split(name, ",")
	lead := " "
	for _, font := range names {
		font = strings.Trim(font, " \n\"'")
		script.WriteString(lead)
		lead = ","
		if strings.Contains(font, " ") {
			script.WriteRune('"')
			script.WriteString(font)
			script.WriteRune('"')
		} else {
			script.WriteString(font)
		}

	}
	script.WriteString("';")
}

func (canvas *canvasData) SetFont(name string, size SizeUnit) {
	canvas.script.WriteString("\nctx.font = '")
	canvas.script.WriteString(size.cssString("1rem", canvas.View().Session()))
	canvas.writeFont(name, &canvas.script)
}

func (canvas *canvasData) fontWithParams(name string, size SizeUnit, params FontParams) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if params.Italic {
		buffer.WriteString("italic ")
	}
	if params.SmallCaps {
		buffer.WriteString("small-caps ")
	}
	if params.Weight > 0 && params.Weight <= 9 {
		switch params.Weight {
		case 4:
			buffer.WriteString("normal ")
		case 7:
			buffer.WriteString("bold ")
		default:
			buffer.WriteString(strconv.Itoa(params.Weight * 100))
			buffer.WriteRune(' ')
		}
	}

	buffer.WriteString(size.cssString("1rem", canvas.View().Session()))
	switch params.LineHeight.Type {
	case Auto:

	case SizeInPercent:
		if params.LineHeight.Value != 100 {
			buffer.WriteString("/")
			buffer.WriteString(strconv.FormatFloat(params.LineHeight.Value/100, 'g', -1, 64))
		}

	case SizeInFraction:
		if params.LineHeight.Value != 1 {
			buffer.WriteString("/")
			buffer.WriteString(strconv.FormatFloat(params.LineHeight.Value, 'g', -1, 64))
		}

	default:
		buffer.WriteString("/")
		buffer.WriteString(params.LineHeight.cssString("", canvas.View().Session()))
	}

	names := strings.Split(name, ",")
	lead := " "
	for _, font := range names {
		font = strings.Trim(font, " \n\"'")
		buffer.WriteString(lead)
		lead = ","
		if strings.Contains(font, " ") {
			buffer.WriteRune('"')
			buffer.WriteString(font)
			buffer.WriteRune('"')
		} else {
			buffer.WriteString(font)
		}
	}

	return buffer.String()
}

func (canvas *canvasData) setFontWithParams(name string, size SizeUnit, params FontParams, script *strings.Builder) {
	script.WriteString("\nctx.font = '")
	if params.Italic {
		script.WriteString("italic ")
	}
	if params.SmallCaps {
		script.WriteString("small-caps ")
	}
	if params.Weight > 0 && params.Weight <= 9 {
		switch params.Weight {
		case 4:
			script.WriteString("normal ")
		case 7:
			script.WriteString("bold ")
		default:
			script.WriteString(strconv.Itoa(params.Weight * 100))
			script.WriteRune(' ')
		}
	}

	script.WriteString(size.cssString("1rem", canvas.View().Session()))
	switch params.LineHeight.Type {
	case Auto:

	case SizeInPercent:
		if params.LineHeight.Value != 100 {
			script.WriteString("/")
			script.WriteString(strconv.FormatFloat(params.LineHeight.Value/100, 'g', -1, 64))
		}

	case SizeInFraction:
		if params.LineHeight.Value != 1 {
			script.WriteString("/")
			script.WriteString(strconv.FormatFloat(params.LineHeight.Value, 'g', -1, 64))
		}

	default:
		script.WriteString("/")
		script.WriteString(params.LineHeight.cssString("", canvas.View().Session()))
	}

	canvas.writeFont(name, script)
}

func (canvas *canvasData) SetFontWithParams(name string, size SizeUnit, params FontParams) {
	canvas.setFontWithParams(name, size, params, &canvas.script)
}

func (canvas *canvasData) TextMetrics(text string, fontName string, fontSize SizeUnit, fontParams FontParams) TextMetrics {
	view := canvas.View()
	return view.Session().canvasTextMetrics(view.htmlID(), canvas.fontWithParams(fontName, fontSize, fontParams), text)
}

func (canvas *canvasData) SetTextBaseline(baseline int) {
	switch baseline {
	case AlphabeticBaseline:
		canvas.script.WriteString("\nctx.textBaseline = 'alphabetic';")
	case TopBaseline:
		canvas.script.WriteString("\nctx.textBaseline = 'top';")
	case MiddleBaseline:
		canvas.script.WriteString("\nctx.textBaseline = 'middle';")
	case BottomBaseline:
		canvas.script.WriteString("\nctx.textBaseline = 'bottom';")
	case HangingBaseline:
		canvas.script.WriteString("\nctx.textBaseline = 'hanging';")
	case IdeographicBaseline:
		canvas.script.WriteString("\nctx.textBaseline = 'ideographic';")
	}
}

func (canvas *canvasData) SetTextAlign(align int) {
	switch align {
	case LeftAlign:
		canvas.script.WriteString("\nctx.textAlign = 'left';")
	case RightAlign:
		canvas.script.WriteString("\nctx.textAlign = 'right';")
	case CenterAlign:
		canvas.script.WriteString("\nctx.textAlign = 'center';")
	case StartAlign:
		canvas.script.WriteString("\nctx.textAlign = 'start';")
	case EndAlign:
		canvas.script.WriteString("\nctx.textAlign = 'end';")
	}
}

func (canvas *canvasData) SetShadow(offsetX, offsetY, blur float64, color Color) {
	if color.Alpha() > 0 && blur >= 0 {
		canvas.script.WriteString("\nctx.shadowColor = '")
		canvas.script.WriteString(color.cssString())
		canvas.script.WriteString("';\nctx.shadowOffsetX = ")
		canvas.script.WriteString(strconv.FormatFloat(offsetX, 'g', -1, 64))
		canvas.script.WriteString(";\nctx.shadowOffsetY = ")
		canvas.script.WriteString(strconv.FormatFloat(offsetY, 'g', -1, 64))
		canvas.script.WriteString(";\nctx.shadowBlur = ")
		canvas.script.WriteString(strconv.FormatFloat(blur, 'g', -1, 64))
		canvas.script.WriteString(";")
	}
}

func (canvas *canvasData) ResetShadow() {
	canvas.script.WriteString("\nctx.shadowColor = 'rgba(0,0,0,0)';\nctx.shadowOffsetX = 0;\nctx.shadowOffsetY = 0;\nctx.shadowBlur = 0;")
}

func (canvas *canvasData) writeRectArgs(x, y, width, height float64) {
	canvas.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(width, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(height, 'g', -1, 64))
}

func (canvas *canvasData) ClearRect(x, y, width, height float64) {
	canvas.script.WriteString("\nctx.clearRect(")
	canvas.writeRectArgs(x, y, width, height)
	canvas.script.WriteString(");")
}

func (canvas *canvasData) FillRect(x, y, width, height float64) {
	canvas.script.WriteString("\nctx.fillRect(")
	canvas.writeRectArgs(x, y, width, height)
	canvas.script.WriteString(");")
}

func (canvas *canvasData) StrokeRect(x, y, width, height float64) {
	canvas.script.WriteString("\nctx.strokeRect(")
	canvas.writeRectArgs(x, y, width, height)
	canvas.script.WriteString(");")
}

func (canvas *canvasData) FillAndStrokeRect(x, y, width, height float64) {
	canvas.FillRect(x, y, width, height)
	canvas.StrokeRect(x, y, width, height)
}

func (canvas *canvasData) writeRoundedRect(x, y, width, height, r float64) {
	left := strconv.FormatFloat(x, 'g', -1, 64)
	top := strconv.FormatFloat(y, 'g', -1, 64)
	right := strconv.FormatFloat(x+width, 'g', -1, 64)
	bottom := strconv.FormatFloat(y+height, 'g', -1, 64)
	leftR := strconv.FormatFloat(x+r, 'g', -1, 64)
	topR := strconv.FormatFloat(y+r, 'g', -1, 64)
	rightR := strconv.FormatFloat(x+width-r, 'g', -1, 64)
	bottomR := strconv.FormatFloat(y+height-r, 'g', -1, 64)
	radius := strconv.FormatFloat(r, 'g', -1, 64)
	canvas.script.WriteString("\nctx.beginPath();\nctx.moveTo(")
	canvas.script.WriteString(left)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(topR)
	canvas.script.WriteString(");\nctx.arc(")
	canvas.script.WriteString(leftR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(topR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(radius)
	canvas.script.WriteString(",Math.PI,Math.PI*3/2);\nctx.lineTo(")
	canvas.script.WriteString(rightR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(top)
	canvas.script.WriteString(");\nctx.arc(")
	canvas.script.WriteString(rightR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(topR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(radius)
	canvas.script.WriteString(",Math.PI*3/2,Math.PI*2);\nctx.lineTo(")
	canvas.script.WriteString(right)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(bottomR)
	canvas.script.WriteString(");\nctx.arc(")
	canvas.script.WriteString(rightR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(bottomR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(radius)
	canvas.script.WriteString(",0,Math.PI/2);\nctx.lineTo(")
	canvas.script.WriteString(leftR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(bottom)
	canvas.script.WriteString(");\nctx.arc(")
	canvas.script.WriteString(leftR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(bottomR)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(radius)
	canvas.script.WriteString(",Math.PI/2,Math.PI);\nctx.closePath();")
}

func (canvas *canvasData) FillRoundedRect(x, y, width, height, r float64) {
	canvas.writeRoundedRect(x, y, width, height, r)
	canvas.script.WriteString("\nctx.fill();")
}

func (canvas *canvasData) StrokeRoundedRect(x, y, width, height, r float64) {
	canvas.writeRoundedRect(x, y, width, height, r)
	canvas.script.WriteString("\nctx.stroke();")
}

func (canvas *canvasData) FillAndStrokeRoundedRect(x, y, width, height, r float64) {
	canvas.writeRoundedRect(x, y, width, height, r)
	canvas.script.WriteString("\nctx.fill();\nctx.stroke();")
}

func (canvas *canvasData) writeEllipse(x, y, radiusX, radiusY, rotation float64) {
	yText := strconv.FormatFloat(y, 'g', -1, 64)
	canvas.script.WriteString("\nctx.beginPath();\nctx.moveTo(")
	canvas.script.WriteString(strconv.FormatFloat(x+radiusX, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(yText)
	canvas.script.WriteString(");\nctx.ellipse(")
	canvas.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(yText)
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(radiusX, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(radiusY, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(rotation, 'g', -1, 64))
	canvas.script.WriteString(",0,Math.PI*2);")
}

func (canvas *canvasData) FillEllipse(x, y, radiusX, radiusY, rotation float64) {
	if radiusX >= 0 && radiusY >= 0 {
		canvas.writeEllipse(x, y, radiusX, radiusY, rotation)
		canvas.script.WriteString("\nctx.fill();")
	}
}

func (canvas *canvasData) StrokeEllipse(x, y, radiusX, radiusY, rotation float64) {
	if radiusX >= 0 && radiusY >= 0 {
		canvas.writeEllipse(x, y, radiusX, radiusY, rotation)
		canvas.script.WriteString("\nctx.stroke();")
	}
}

func (canvas *canvasData) FillAndStrokeEllipse(x, y, radiusX, radiusY, rotation float64) {
	if radiusX >= 0 && radiusY >= 0 {
		canvas.writeEllipse(x, y, radiusX, radiusY, rotation)
		canvas.script.WriteString("\nctx.fill();\nctx.stroke();")
	}
}

func (canvas *canvasData) writePointArgs(x, y float64) {
	canvas.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
}

func (canvas *canvasData) writeStringArgs(text string, script *strings.Builder) {
	//rText := []rune(text)
	for _, ch := range text {
		switch ch {
		case '\t':
			script.WriteString(`\t`)
		case '\n':
			script.WriteString(`\n`)
		case '\r':
			script.WriteString(`\r`)
		case '\\':
			script.WriteString(`\\`)
		case '"':
			script.WriteString(`\"`)
		case '\'':
			script.WriteString(`\'`)
		default:
			if ch < ' ' {
				script.WriteString(fmt.Sprintf("\\x%02X", int(ch)))
			} else {
				script.WriteRune(ch)
			}
		}
	}
}

func (canvas *canvasData) FillText(x, y float64, text string) {
	canvas.script.WriteString("\nctx.fillText('")
	canvas.writeStringArgs(text, &canvas.script)
	canvas.script.WriteString(`',`)
	canvas.writePointArgs(x, y)
	canvas.script.WriteString(");")
}

func (canvas *canvasData) StrokeText(x, y float64, text string) {
	canvas.script.WriteString("\nctx.strokeText('")
	canvas.writeStringArgs(text, &canvas.script)
	canvas.script.WriteString(`',`)
	canvas.writePointArgs(x, y)
	canvas.script.WriteString(");")
}

func (canvas *canvasData) FillPath(path Path) {
	canvas.script.WriteString(path.scriptText())
	canvas.script.WriteString("\nctx.fill();")
}

func (canvas *canvasData) StrokePath(path Path) {
	canvas.script.WriteString(path.scriptText())
	canvas.script.WriteString("\nctx.stroke();")
}

func (canvas *canvasData) FillAndStrokePath(path Path) {
	canvas.script.WriteString(path.scriptText())
	canvas.script.WriteString("\nctx.fill();\nctx.stroke();")
}

func (canvas *canvasData) DrawLine(x0, y0, x1, y1 float64) {
	canvas.script.WriteString("\nctx.beginPath();\nctx.moveTo(")
	canvas.script.WriteString(strconv.FormatFloat(x0, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y0, 'g', -1, 64))
	canvas.script.WriteString(");\nctx.lineTo(")
	canvas.script.WriteString(strconv.FormatFloat(x1, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y1, 'g', -1, 64))
	canvas.script.WriteString(");\nctx.stroke();")
}

func (canvas *canvasData) DrawImage(x, y float64, image Image) {
	if image == nil || image.LoadingStatus() != ImageReady {
		return
	}
	canvas.script.WriteString("\nimg = images.get('")
	canvas.script.WriteString(image.URL())
	canvas.script.WriteString("');\nif (img) {\nctx.drawImage(img,")
	canvas.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	canvas.script.WriteString(");\n}")
}

func (canvas *canvasData) DrawImageInRect(x, y, width, height float64, image Image) {
	if image == nil || image.LoadingStatus() != ImageReady {
		return
	}
	canvas.script.WriteString("\nimg = images.get('")
	canvas.script.WriteString(image.URL())
	canvas.script.WriteString("');\nif (img) {\nctx.drawImage(img,")
	canvas.script.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(y, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(width, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(height, 'g', -1, 64))
	canvas.script.WriteString(");\n}")
}

func (canvas *canvasData) DrawImageFragment(srcX, srcY, srcWidth, srcHeight, dstX, dstY, dstWidth, dstHeight float64, image Image) {
	if image == nil || image.LoadingStatus() != ImageReady {
		return
	}
	canvas.script.WriteString("\nimg = images.get('")
	canvas.script.WriteString(image.URL())
	canvas.script.WriteString("');\nif (img) {\nctx.drawImage(img,")
	canvas.script.WriteString(strconv.FormatFloat(srcX, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(srcY, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(srcWidth, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(srcHeight, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(dstX, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(dstY, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(dstWidth, 'g', -1, 64))
	canvas.script.WriteRune(',')
	canvas.script.WriteString(strconv.FormatFloat(dstHeight, 'g', -1, 64))
	canvas.script.WriteString(");\n}")
}
