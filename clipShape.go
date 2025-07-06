package rui

import (
	"fmt"
	"strings"
)

type ClipShape string

const (
	InsetClip   ClipShape = "inset"
	CircleClip  ClipShape = "circle"
	EllipseClip ClipShape = "ellipse"
	PolygonClip ClipShape = "polygon"
)

// ClipShapeProperty defines a View clipping area
type ClipShapeProperty interface {
	Properties
	fmt.Stringer
	stringWriter

	// Shape returns the clip shape type
	Shape() ClipShape
	cssStyle(session Session) string
	valid(session Session) bool
}

type insetClipData struct {
	dataProperty
}

type ellipseClipData struct {
	dataProperty
}

type circleClipData struct {
	dataProperty
}

type polygonClipData struct {
	dataProperty
}

// NewClipShapeProperty creates ClipShapeProperty.
//
// The following properties can be used for shapes:
//
// InsetClip:
//   - "top" (Top) - offset (SizeUnit) from the top border of a View;
//   - "right" (Right) - offset (SizeUnit) from the right border of a View;
//   - "bottom" (Bottom) - offset (SizeUnit) from the bottom border of a View;
//   - "left" (Left) - offset (SizeUnit) from the left border of a View;
//   - "radius" (Radius) - corner radius (RadiusProperty).
//
// CircleClip:
//   - "x" (X) - x-axis position (SizeUnit) of the circle clip center;
//   - "y" (Y) - y-axis position (SizeUnit) of the circle clip center;
//   - "radius" (Radius) - radius (SizeUnit) of the circle clip center.
//
// EllipseClip:
//   - "x" (X) - x-axis position (SizeUnit) of the ellipse clip center;
//   - "y" (Y) - y-axis position (SizeUnit) of the ellipse clip center;
//   - "radius-x" (RadiusX) - x-axis radius (SizeUnit) of the ellipse clip center;
//   - "radius-y" (RadiusY) - y-axis radius (SizeUnit) of the ellipse clip center.
//
// PolygonClip:
//   - "points" (Points) - an array ([]SizeUnit) of corner points of the polygon in the following order: x1, y1, x2, y2, ….
//
// The function will return nil if no properties are specified, unsupported properties are specified, or at least one property has an invalid value.
func NewClipShapeProperty(shape ClipShape, params Params) ClipShapeProperty {
	if len(params) == 0 {
		ErrorLog("No ClipShapeProperty params")
		return nil
	}

	var result ClipShapeProperty

	switch shape {
	case InsetClip:
		clip := new(insetClipData)
		clip.init()
		result = clip

	case CircleClip:
		clip := new(circleClipData)
		clip.init()
		result = clip

	case EllipseClip:
		clip := new(ellipseClipData)
		clip.init()
		result = clip

	case PolygonClip:
		clip := new(polygonClipData)
		clip.init()
		result = clip

	default:
		ErrorLog("Unknown ClipShape: " + string(shape))
		return nil
	}

	for tag, value := range params {
		if !result.Set(tag, value) {
			return nil
		}
	}

	return result
}

// NewInsetClip creates a rectangle View clipping area.
//   - top - offset from the top border of a View;
//   - right - offset from the right border of a View;
//   - bottom - offset from the bottom border of a View;
//   - left - offset from the left border of a View;
//   - radius - corner radius, pass nil if you don't need to round corners
func NewInsetClip(top, right, bottom, left SizeUnit, radius RadiusProperty) ClipShapeProperty {
	clip := new(insetClipData)
	clip.init()
	clip.setRaw(Top, top)
	clip.setRaw(Right, right)
	clip.setRaw(Bottom, bottom)
	clip.setRaw(Left, left)
	if radius != nil {
		clip.setRaw(Radius, radius)
	}
	return clip
}

// NewCircleClip creates a circle View clipping area.
//   - x - x-axis position of the circle clip center;
//   - y - y-axis position of the circle clip center;
//   - radius - radius of the circle clip center.
func NewCircleClip(x, y, radius SizeUnit) ClipShapeProperty {
	clip := new(circleClipData)
	clip.init()
	clip.setRaw(X, x)
	clip.setRaw(Y, y)
	clip.setRaw(Radius, radius)
	return clip
}

// NewEllipseClip creates a ellipse View clipping area.
//   - x - x-axis position of the ellipse clip center;
//   - y - y-axis position of the ellipse clip center;
//   - rx - x-axis radius of the ellipse clip center;
//   - ry - y-axis radius of the ellipse clip center.
func NewEllipseClip(x, y, rx, ry SizeUnit) ClipShapeProperty {
	clip := new(ellipseClipData)
	clip.init()
	clip.setRaw(X, x)
	clip.setRaw(Y, y)
	clip.setRaw(RadiusX, rx)
	clip.setRaw(RadiusY, ry)
	return clip
}

// NewPolygonClip creates a polygon View clipping area.
//   - points - an array of corner points of the polygon in the following order: x1, y1, x2, y2, …
//
// The elements of the function argument can be or text constants,
// or the text representation of SizeUnit, or elements of SizeUnit type.
func NewPolygonClip(points []any) ClipShapeProperty {
	clip := new(polygonClipData)
	clip.init()
	if polygonClipDataSet(clip, Points, points) != nil {
		return clip
	}
	return nil
}

// NewPolygonPointsClip creates a polygon View clipping area.
//   - points - an array of corner points of the polygon in the following order: x1, y1, x2, y2, …
func NewPolygonPointsClip(points []SizeUnit) ClipShapeProperty {
	clip := new(polygonClipData)
	clip.init()
	if polygonClipDataSet(clip, Points, points) != nil {
		return clip
	}
	return nil
}

func (clip *insetClipData) init() {
	clip.dataProperty.init()
	clip.set = insetClipDataSet
	clip.supportedProperties = []PropertyName{
		Top, Right, Bottom, Left, Radius,
		RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY,
	}
}

func (clip *insetClipData) Shape() ClipShape {
	return InsetClip
}

func insetClipDataSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case Top, Right, Bottom, Left:
		return setSizeProperty(properties, tag, value)

	case Radius:
		return setRadiusProperty(properties, value)

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		if setRadiusPropertyElement(properties, tag, value) {
			return []PropertyName{tag, Radius}
		}
		return nil
	}

	ErrorLogF(`"%s" property is not supported by the inset clip shape`, tag)
	return nil
}

func (clip *insetClipData) String() string {
	return runStringWriter(clip)
}

func (clip *insetClipData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("inset { ")
	comma := false
	for _, tag := range []PropertyName{Top, Right, Bottom, Left, Radius} {
		if value, ok := clip.properties[tag]; ok {
			text := propertyValueToString(tag, value, indent)
			if text != "" {
				if comma {
					buffer.WriteString(", ")
				}
				buffer.WriteString(string(tag))
				buffer.WriteString(" = ")
				buffer.WriteString(text)
				comma = true
			}
		}
	}

	buffer.WriteString(" }")
}

func (clip *insetClipData) cssStyle(session Session) string {

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	leadText := "inset("
	for _, tag := range []PropertyName{Top, Right, Bottom, Left} {
		value, _ := sizeProperty(clip, tag, session)
		buffer.WriteString(leadText)
		buffer.WriteString(value.cssString("0px", session))
		leadText = " "
	}

	if radius := getRadiusProperty(clip); radius != nil {
		buffer.WriteString(" round ")
		buffer.WriteString(radius.BoxRadius(session).cssString(session))
	}

	buffer.WriteRune(')')
	return buffer.String()
}

func (clip *insetClipData) valid(session Session) bool {
	for _, tag := range []PropertyName{Top, Right, Bottom, Left, Radius, RadiusX, RadiusY} {
		if value, ok := sizeProperty(clip, tag, session); ok && value.Type != Auto && value.Value != 0 {
			return true
		}
	}
	return false
}

func (clip *circleClipData) init() {
	clip.dataProperty.init()
	clip.set = circleClipDataSet
	clip.supportedProperties = []PropertyName{X, Y, Radius}
}

func (clip *circleClipData) Shape() ClipShape {
	return CircleClip
}

func circleClipDataSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case X, Y, Radius:
		return setSizeProperty(properties, tag, value)
	}

	ErrorLogF(`"%s" property is not supported by the circle clip shape`, tag)
	return nil
}

func (clip *circleClipData) String() string {
	return runStringWriter(clip)
}

func (clip *circleClipData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("circle { ")
	comma := false
	for _, tag := range []PropertyName{Radius, X, Y} {
		if value, ok := clip.properties[tag]; ok {
			text := propertyValueToString(tag, value, indent)
			if text != "" {
				if comma {
					buffer.WriteString(", ")
				}
				buffer.WriteString(string(tag))
				buffer.WriteString(" = ")
				buffer.WriteString(text)
				comma = true
			}
		}
	}

	buffer.WriteString(" }")
}

func (clip *circleClipData) cssStyle(session Session) string {

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString("circle(")
	r, _ := sizeProperty(clip, Radius, session)
	buffer.WriteString(r.cssString("50%", session))

	buffer.WriteString(" at ")
	x, _ := sizeProperty(clip, X, session)
	buffer.WriteString(x.cssString("50%", session))
	buffer.WriteRune(' ')

	y, _ := sizeProperty(clip, Y, session)
	buffer.WriteString(y.cssString("50%", session))
	buffer.WriteRune(')')

	return buffer.String()
}

func (clip *circleClipData) valid(session Session) bool {
	if value, ok := sizeProperty(clip, Radius, session); ok && value.Value == 0 {
		return false
	}
	return true
}

func (clip *ellipseClipData) init() {
	clip.dataProperty.init()
	clip.set = ellipseClipDataSet
	clip.supportedProperties = []PropertyName{X, Y, Radius, RadiusX, RadiusY}
}

func (clip *ellipseClipData) Shape() ClipShape {
	return EllipseClip
}

func ellipseClipDataSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case X, Y, RadiusX, RadiusY:
		return setSizeProperty(properties, tag, value)

	case Radius:
		if result := setSizeProperty(properties, RadiusX, value); result != nil {
			properties.setRaw(RadiusY, properties.getRaw(RadiusX))
			return append(result, RadiusY)
		}
		return nil
	}

	ErrorLogF(`"%s" property is not supported by the ellipse clip shape`, tag)
	return nil
}

func (clip *ellipseClipData) String() string {
	return runStringWriter(clip)
}

func (clip *ellipseClipData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("ellipse { ")
	comma := false
	for _, tag := range []PropertyName{RadiusX, RadiusY, X, Y} {
		if value, ok := clip.properties[tag]; ok {
			text := propertyValueToString(tag, value, indent)
			if text != "" {
				if comma {
					buffer.WriteString(", ")
				}
				buffer.WriteString(string(tag))
				buffer.WriteString(" = ")
				buffer.WriteString(text)
				comma = true
			}
		}
	}

	buffer.WriteString(" }")
}

func (clip *ellipseClipData) cssStyle(session Session) string {

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	rx, _ := sizeProperty(clip, RadiusX, session)
	ry, _ := sizeProperty(clip, RadiusX, session)
	buffer.WriteString("ellipse(")
	buffer.WriteString(rx.cssString("50%", session))
	buffer.WriteRune(' ')
	buffer.WriteString(ry.cssString("50%", session))

	buffer.WriteString(" at ")
	x, _ := sizeProperty(clip, X, session)
	buffer.WriteString(x.cssString("50%", session))
	buffer.WriteRune(' ')

	y, _ := sizeProperty(clip, Y, session)
	buffer.WriteString(y.cssString("50%", session))
	buffer.WriteRune(')')

	return buffer.String()
}

func (clip *ellipseClipData) valid(session Session) bool {
	rx, _ := sizeProperty(clip, RadiusX, session)
	ry, _ := sizeProperty(clip, RadiusY, session)
	return rx.Value != 0 && ry.Value != 0
}

func (clip *polygonClipData) init() {
	clip.dataProperty.init()
	clip.set = polygonClipDataSet
	clip.supportedProperties = []PropertyName{Points}
}

func (clip *polygonClipData) Shape() ClipShape {
	return PolygonClip
}

func polygonClipDataSet(properties Properties, tag PropertyName, value any) []PropertyName {
	if Points == tag {
		switch value := value.(type) {
		case []any:
			points := make([]any, len(value))
			for i, val := range value {
				switch val := val.(type) {
				case string:
					if ok, _ := isConstantName(val); ok {
						points[i] = val
					} else if size, ok := StringToSizeUnit(val); ok {
						points[i] = size
					} else {
						notCompatibleType(tag, val)
						return nil
					}

				case SizeUnit:
					points[i] = val

				default:
					notCompatibleType(tag, val)
					points[i] = AutoSize()
					return nil
				}
			}
			properties.setRaw(Points, points)
			return []PropertyName{tag}

		case []SizeUnit:
			points := make([]any, len(value))
			for i, point := range value {
				points[i] = point
			}
			properties.setRaw(Points, points)
			return []PropertyName{tag}

		case string:
			values := strings.Split(value, ",")
			points := make([]any, len(values))
			for i, val := range values {
				val = strings.Trim(val, " \t\n\r")
				if ok, _ := isConstantName(val); ok {
					points[i] = val
				} else if size, ok := StringToSizeUnit(val); ok {
					points[i] = size
				} else {
					notCompatibleType(tag, val)
					return nil
				}
			}
			properties.setRaw(Points, points)
			return []PropertyName{tag}
		}
	}
	return nil
}

func (clip *polygonClipData) String() string {
	return runStringWriter(clip)
}

func (clip *polygonClipData) points() []any {
	if value := clip.getRaw(Points); value != nil {
		if points, ok := value.([]any); ok {
			return points
		}
	}
	return nil
}

func (clip *polygonClipData) writeString(buffer *strings.Builder, indent string) {

	buffer.WriteString("polygon { ")

	if points := clip.points(); points != nil {
		buffer.WriteString(string(Points))
		buffer.WriteString(` = "`)
		comma := false
		for _, value := range points {
			text := propertyValueToString("", value, indent)
			if text != "" {
				if comma {
					buffer.WriteString(", ")
				}
				buffer.WriteString(text)
				comma = true
			}
		}

		buffer.WriteString(`" `)
	}
	buffer.WriteRune('}')
}

func (clip *polygonClipData) cssStyle(session Session) string {

	points := clip.points()
	count := len(points)
	if count < 2 {
		return ""
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	writePoint := func(value any) {
		switch value := value.(type) {
		case string:
			if val, ok := session.resolveConstants(value); ok {
				if size, ok := StringToSizeUnit(val); ok {
					buffer.WriteString(size.cssString("0px", session))
					return
				}
			}

		case SizeUnit:
			buffer.WriteString(value.cssString("0px", session))
			return
		}

		buffer.WriteString("0px")
	}

	leadText := "polygon("
	for i := 1; i < count; i += 2 {
		buffer.WriteString(leadText)
		writePoint(points[i-1])
		buffer.WriteRune(' ')
		writePoint(points[i])
		leadText = ", "
	}

	buffer.WriteRune(')')
	return buffer.String()
}

func (clip *polygonClipData) valid(session Session) bool {
	return len(clip.points()) > 0
}

func parseClipShapeProperty(obj DataObject) ClipShapeProperty {
	switch obj.Tag() {
	case "inset":
		clip := new(insetClipData)
		clip.init()
		for _, tag := range []PropertyName{Top, Right, Bottom, Left, Radius, RadiusX, RadiusY} {
			if value, ok := obj.PropertyValue(string(tag)); ok {
				insetClipDataSet(clip, tag, value)
			}
		}
		return clip

	case "circle":
		clip := new(circleClipData)
		clip.init()
		for _, tag := range []PropertyName{X, Y, Radius} {
			if value, ok := obj.PropertyValue(string(tag)); ok {
				circleClipDataSet(clip, tag, value)
			}
		}
		return clip

	case "ellipse":
		clip := new(ellipseClipData)
		clip.init()
		for _, tag := range []PropertyName{X, Y, RadiusX, RadiusY} {
			if value, ok := obj.PropertyValue(string(tag)); ok {
				ellipseClipDataSet(clip, tag, value)
			}
		}
		return clip

	case "polygon":
		clip := new(polygonClipData)
		clip.init()
		if value, ok := obj.PropertyValue(string(Points)); ok {
			polygonClipDataSet(clip, Points, value)
		}
		return clip
	}

	return nil
}

func setClipShapePropertyProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	switch value := value.(type) {
	case ClipShapeProperty:
		properties.setRaw(tag, value)
		return []PropertyName{tag}

	case string:
		if ok, _ := isConstantName(value); ok {
			properties.setRaw(tag, value)
			return []PropertyName{tag}
		}

		if obj := NewDataObject(value); obj == nil {
			if clip := parseClipShapeProperty(obj); clip != nil {
				properties.setRaw(tag, clip)
				return []PropertyName{tag}
			}
		}

	case DataObject:
		if clip := parseClipShapeProperty(value); clip != nil {
			properties.setRaw(tag, clip)
			return []PropertyName{tag}
		}

	case DataValue:
		if value.IsObject() {
			if clip := parseClipShapeProperty(value.Object()); clip != nil {
				properties.setRaw(tag, clip)
				return []PropertyName{tag}
			}
		}
	}

	notCompatibleType(tag, value)
	return nil
}

func getClipShapeProperty(prop Properties, tag PropertyName, session Session) ClipShapeProperty {
	if value := prop.getRaw(tag); value != nil {
		switch value := value.(type) {
		case ClipShapeProperty:
			return value

		case string:
			if text, ok := session.resolveConstants(value); ok {
				if obj := NewDataObject(text); obj == nil {
					return parseClipShapeProperty(obj)
				}
			}
		}
	}

	return nil
}

// GetClip returns a View clipping area.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetClip(view View, subviewID ...string) ClipShapeProperty {
	if view = getSubview(view, subviewID); view != nil {
		return getClipShapeProperty(view, Clip, view.Session())
	}

	return nil
}

// GetShapeOutside returns a shape around which adjacent inline content.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetShapeOutside(view View, subviewID ...string) ClipShapeProperty {
	if view = getSubview(view, subviewID); view != nil {
		return getClipShapeProperty(view, ShapeOutside, view.Session())
	}

	return nil
}
