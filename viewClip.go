package rui

import (
	"fmt"
	"strings"
)

// ClipShape defines a View clipping area
type ClipShape interface {
	Properties
	fmt.Stringer
	stringWriter
	cssStyle(session Session) string
	valid(session Session) bool
}

type insetClip struct {
	dataProperty
}

type ellipseClip struct {
	dataProperty
}

type circleClip struct {
	dataProperty
}

type polygonClip struct {
	dataProperty
}

// InsetClip creates a rectangle View clipping area.
//   - top - offset from the top border of a View;
//   - right - offset from the right border of a View;
//   - bottom - offset from the bottom border of a View;
//   - left - offset from the left border of a View;
//   - radius - corner radius, pass nil if you don't need to round corners
func InsetClip(top, right, bottom, left SizeUnit, radius RadiusProperty) ClipShape {
	clip := new(insetClip)
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

// CircleClip creates a circle View clipping area.
func CircleClip(x, y, radius SizeUnit) ClipShape {
	clip := new(circleClip)
	clip.init()
	clip.setRaw(X, x)
	clip.setRaw(Y, y)
	clip.setRaw(Radius, radius)
	return clip
}

// EllipseClip creates a ellipse View clipping area.
func EllipseClip(x, y, rx, ry SizeUnit) ClipShape {
	clip := new(ellipseClip)
	clip.init()
	clip.setRaw(X, x)
	clip.setRaw(Y, y)
	clip.setRaw(RadiusX, rx)
	clip.setRaw(RadiusY, ry)
	return clip
}

// PolygonClip creates a polygon View clipping area.
//
// The elements of the function argument can be or text constants,
// or the text representation of SizeUnit, or elements of SizeUnit type.
func PolygonClip(points []any) ClipShape {
	clip := new(polygonClip)
	clip.init()
	if polygonClipSet(clip, Points, points) != nil {
		return clip
	}
	return nil
}

// PolygonPointsClip creates a polygon View clipping area.
func PolygonPointsClip(points []SizeUnit) ClipShape {
	clip := new(polygonClip)
	clip.init()
	if polygonClipSet(clip, Points, points) != nil {
		return clip
	}
	return nil
}

func (clip *insetClip) init() {
	clip.dataProperty.init()
	clip.set = insetClipSet
	clip.supportedProperties = []PropertyName{
		Top, Right, Bottom, Left, Radius,
		RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY,
	}
}

func insetClipSet(properties Properties, tag PropertyName, value any) []PropertyName {
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

func (clip *insetClip) String() string {
	return runStringWriter(clip)
}

func (clip *insetClip) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("inset { ")
	comma := false
	for _, tag := range []PropertyName{Top, Right, Bottom, Left, Radius} {
		if value, ok := clip.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(string(tag))
			buffer.WriteString(" = ")
			writePropertyValue(buffer, tag, value, indent)
			comma = true
		}
	}

	buffer.WriteString(" }")
}

func (clip *insetClip) cssStyle(session Session) string {

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

func (clip *insetClip) valid(session Session) bool {
	for _, tag := range []PropertyName{Top, Right, Bottom, Left, Radius, RadiusX, RadiusY} {
		if value, ok := sizeProperty(clip, tag, session); ok && value.Type != Auto && value.Value != 0 {
			return true
		}
	}
	return false
}

func (clip *circleClip) init() {
	clip.dataProperty.init()
	clip.set = circleClipSet
	clip.supportedProperties = []PropertyName{X, Y, Radius}
}

func circleClipSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case X, Y, Radius:
		return setSizeProperty(properties, tag, value)
	}

	ErrorLogF(`"%s" property is not supported by the circle clip shape`, tag)
	return nil
}

func (clip *circleClip) String() string {
	return runStringWriter(clip)
}

func (clip *circleClip) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("circle { ")
	comma := false
	for _, tag := range []PropertyName{Radius, X, Y} {
		if value, ok := clip.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(string(tag))
			buffer.WriteString(" = ")
			writePropertyValue(buffer, tag, value, indent)
			comma = true
		}
	}

	buffer.WriteString(" }")
}

func (clip *circleClip) cssStyle(session Session) string {

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

func (clip *circleClip) valid(session Session) bool {
	if value, ok := sizeProperty(clip, Radius, session); ok && value.Value == 0 {
		return false
	}
	return true
}

func (clip *ellipseClip) init() {
	clip.dataProperty.init()
	clip.set = ellipseClipSet
	clip.supportedProperties = []PropertyName{X, Y, Radius, RadiusX, RadiusY}
}

func ellipseClipSet(properties Properties, tag PropertyName, value any) []PropertyName {
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

func (clip *ellipseClip) String() string {
	return runStringWriter(clip)
}

func (clip *ellipseClip) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("ellipse { ")
	comma := false
	for _, tag := range []PropertyName{RadiusX, RadiusY, X, Y} {
		if value, ok := clip.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(string(tag))
			buffer.WriteString(" = ")
			writePropertyValue(buffer, tag, value, indent)
			comma = true
		}
	}

	buffer.WriteString(" }")
}

func (clip *ellipseClip) cssStyle(session Session) string {

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

func (clip *ellipseClip) valid(session Session) bool {
	rx, _ := sizeProperty(clip, RadiusX, session)
	ry, _ := sizeProperty(clip, RadiusY, session)
	return rx.Value != 0 && ry.Value != 0
}

func (clip *polygonClip) init() {
	clip.dataProperty.init()
	clip.set = polygonClipSet
	clip.supportedProperties = []PropertyName{Points}
}

func polygonClipSet(properties Properties, tag PropertyName, value any) []PropertyName {
	if Points == tag {
		switch value := value.(type) {
		case []any:
			points := make([]any, len(value))
			for i, val := range value {
				switch val := val.(type) {
				case string:
					if isConstantName(val) {
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
				if isConstantName(val) {
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

func (clip *polygonClip) String() string {
	return runStringWriter(clip)
}

func (clip *polygonClip) points() []any {
	if value := clip.getRaw(Points); value != nil {
		if points, ok := value.([]any); ok {
			return points
		}
	}
	return nil
}

func (clip *polygonClip) writeString(buffer *strings.Builder, indent string) {

	buffer.WriteString("inset { ")

	if points := clip.points(); points != nil {
		buffer.WriteString(string(Points))
		buffer.WriteString(` = "`)
		for i, value := range points {
			if i > 0 {
				buffer.WriteString(", ")
			}
			writePropertyValue(buffer, "", value, indent)
		}

		buffer.WriteString(`" `)
	}
	buffer.WriteRune('}')
}

func (clip *polygonClip) cssStyle(session Session) string {

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

func (clip *polygonClip) valid(session Session) bool {
	return len(clip.points()) > 0
}

func parseClipShape(obj DataObject) ClipShape {
	switch obj.Tag() {
	case "inset":
		clip := new(insetClip)
		clip.init()
		for _, tag := range []PropertyName{Top, Right, Bottom, Left, Radius, RadiusX, RadiusY} {
			if value, ok := obj.PropertyValue(string(tag)); ok {
				insetClipSet(clip, tag, value)
			}
		}
		return clip

	case "circle":
		clip := new(ellipseClip)
		clip.init()
		for _, tag := range []PropertyName{X, Y, Radius} {
			if value, ok := obj.PropertyValue(string(tag)); ok {
				circleClipSet(clip, tag, value)
			}
		}
		return clip

	case "ellipse":
		clip := new(ellipseClip)
		clip.init()
		for _, tag := range []PropertyName{X, Y, RadiusX, RadiusY} {
			if value, ok := obj.PropertyValue(string(tag)); ok {
				ellipseClipSet(clip, tag, value)
			}
		}
		return clip

	case "polygon":
		clip := new(ellipseClip)
		clip.init()
		if value, ok := obj.PropertyValue(string(Points)); ok {
			polygonClipSet(clip, Points, value)
		}
		return clip
	}

	return nil
}

func setClipShapeProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	switch value := value.(type) {
	case ClipShape:
		properties.setRaw(tag, value)
		return []PropertyName{tag}

	case string:
		if isConstantName(value) {
			properties.setRaw(tag, value)
			return []PropertyName{tag}
		}

		if obj := NewDataObject(value); obj == nil {
			if clip := parseClipShape(obj); clip != nil {
				properties.setRaw(tag, clip)
				return []PropertyName{tag}
			}
		}

	case DataObject:
		if clip := parseClipShape(value); clip != nil {
			properties.setRaw(tag, clip)
			return []PropertyName{tag}
		}

	case DataValue:
		if value.IsObject() {
			if clip := parseClipShape(value.Object()); clip != nil {
				properties.setRaw(tag, clip)
				return []PropertyName{tag}
			}
		}
	}

	notCompatibleType(tag, value)
	return nil
}

func getClipShape(prop Properties, tag PropertyName, session Session) ClipShape {
	if value := prop.getRaw(tag); value != nil {
		switch value := value.(type) {
		case ClipShape:
			return value

		case string:
			if text, ok := session.resolveConstants(value); ok {
				if obj := NewDataObject(text); obj == nil {
					return parseClipShape(obj)
				}
			}
		}
	}

	return nil
}

// GetClip returns a View clipping area.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetClip(view View, subviewID ...string) ClipShape {
	if view = getSubview(view, subviewID); view != nil {
		return getClipShape(view, Clip, view.Session())
	}

	return nil
}

// GetShapeOutside returns a shape around which adjacent inline content.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetShapeOutside(view View, subviewID ...string) ClipShape {
	if view = getSubview(view, subviewID); view != nil {
		return getClipShape(view, ShapeOutside, view.Session())
	}

	return nil
}
