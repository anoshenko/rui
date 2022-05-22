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
	propertyList
}

type ellipseClip struct {
	propertyList
}

type circleClip struct {
	propertyList
}

type polygonClip struct {
	points []interface{}
}

// InsetClip creates a rectangle View clipping area.
// top - offset from the top border of a View;
// right - offset from the right border of a View;
// bottom - offset from the bottom border of a View;
// left - offset from the left border of a View;
// radius - corner radius, pass nil if you don't need to round corners
func InsetClip(top, right, bottom, left SizeUnit, radius RadiusProperty) ClipShape {
	clip := new(insetClip)
	clip.init()
	clip.Set(Top, top)
	clip.Set(Right, right)
	clip.Set(Bottom, bottom)
	clip.Set(Left, left)
	if radius != nil {
		clip.Set(Radius, radius)
	}
	return clip
}

// CircleClip creates a circle View clipping area.
func CircleClip(x, y, radius SizeUnit) ClipShape {
	clip := new(circleClip)
	clip.init()
	clip.Set(X, x)
	clip.Set(Y, y)
	clip.Set(Radius, radius)
	return clip
}

// EllipseClip creates a ellipse View clipping area.
func EllipseClip(x, y, rx, ry SizeUnit) ClipShape {
	clip := new(ellipseClip)
	clip.init()
	clip.Set(X, x)
	clip.Set(Y, y)
	clip.Set(RadiusX, rx)
	clip.Set(RadiusY, ry)
	return clip
}

// PolygonClip creates a polygon View clipping area.
// The elements of the function argument can be or text constants,
// or the text representation of SizeUnit, or elements of SizeUnit type.
func PolygonClip(points []interface{}) ClipShape {
	clip := new(polygonClip)
	clip.points = []interface{}{}
	if clip.Set(Points, points) {
		return clip
	}
	return nil
}

// PolygonPointsClip creates a polygon View clipping area.
func PolygonPointsClip(points []SizeUnit) ClipShape {
	clip := new(polygonClip)
	clip.points = []interface{}{}
	if clip.Set(Points, points) {
		return clip
	}
	return nil
}

func (clip *insetClip) Set(tag string, value interface{}) bool {
	switch strings.ToLower(tag) {
	case Top, Right, Bottom, Left:
		if value == nil {
			clip.Remove(tag)
			return true
		}
		return clip.setSizeProperty(tag, value)

	case Radius:
		return clip.setRadius(value)

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		return clip.setRadiusElement(tag, value)
	}

	ErrorLogF(`"%s" property is not supported by the inset clip shape`, tag)
	return false
}

func (clip *insetClip) String() string {
	return runStringWriter(clip)
}

func (clip *insetClip) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("inset { ")
	comma := false
	for _, tag := range []string{Top, Right, Bottom, Left, Radius} {
		if value, ok := clip.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(tag)
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
	for _, tag := range []string{Top, Right, Bottom, Left} {
		value, _ := sizeProperty(clip, tag, session)
		buffer.WriteString(leadText)
		buffer.WriteString(value.cssString("0px"))
		leadText = " "
	}

	if radius := getRadiusProperty(clip); radius != nil {
		buffer.WriteString(" round ")
		buffer.WriteString(radius.BoxRadius(session).cssString())
	}

	buffer.WriteRune(')')
	return buffer.String()
}

func (clip *insetClip) valid(session Session) bool {
	for _, tag := range []string{Top, Right, Bottom, Left, Radius, RadiusX, RadiusY} {
		if value, ok := sizeProperty(clip, tag, session); ok && value.Type != Auto && value.Value != 0 {
			return true
		}
	}
	return false
}

func (clip *circleClip) Set(tag string, value interface{}) bool {
	if value == nil {
		clip.Remove(tag)
	}

	switch strings.ToLower(tag) {
	case X, Y, Radius:
		return clip.setSizeProperty(tag, value)
	}

	ErrorLogF(`"%s" property is not supported by the circle clip shape`, tag)
	return false
}

func (clip *circleClip) String() string {
	return runStringWriter(clip)
}

func (clip *circleClip) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("circle { ")
	comma := false
	for _, tag := range []string{Radius, X, Y} {
		if value, ok := clip.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(tag)
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
	buffer.WriteString(r.cssString("50%"))

	buffer.WriteString(" at ")
	x, _ := sizeProperty(clip, X, session)
	buffer.WriteString(x.cssString("50%"))
	buffer.WriteRune(' ')

	y, _ := sizeProperty(clip, Y, session)
	buffer.WriteString(y.cssString("50%"))
	buffer.WriteRune(')')

	return buffer.String()
}

func (clip *circleClip) valid(session Session) bool {
	if value, ok := sizeProperty(clip, Radius, session); ok && value.Value == 0 {
		return false
	}
	return true
}

func (clip *ellipseClip) Set(tag string, value interface{}) bool {
	if value == nil {
		clip.Remove(tag)
	}

	switch strings.ToLower(tag) {
	case X, Y, RadiusX, RadiusY:
		return clip.setSizeProperty(tag, value)

	case Radius:
		return clip.setSizeProperty(RadiusX, value) &&
			clip.setSizeProperty(RadiusY, value)
	}

	ErrorLogF(`"%s" property is not supported by the ellipse clip shape`, tag)
	return false
}

func (clip *ellipseClip) String() string {
	return runStringWriter(clip)
}

func (clip *ellipseClip) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("ellipse { ")
	comma := false
	for _, tag := range []string{RadiusX, RadiusY, X, Y} {
		if value, ok := clip.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(tag)
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
	buffer.WriteString(rx.cssString("50%"))
	buffer.WriteRune(' ')
	buffer.WriteString(ry.cssString("50%"))

	buffer.WriteString(" at ")
	x, _ := sizeProperty(clip, X, session)
	buffer.WriteString(x.cssString("50%"))
	buffer.WriteRune(' ')

	y, _ := sizeProperty(clip, Y, session)
	buffer.WriteString(y.cssString("50%"))
	buffer.WriteRune(')')

	return buffer.String()
}

func (clip *ellipseClip) valid(session Session) bool {
	rx, _ := sizeProperty(clip, RadiusX, session)
	ry, _ := sizeProperty(clip, RadiusY, session)
	return rx.Value != 0 && ry.Value != 0
}

func (clip *polygonClip) Get(tag string) interface{} {
	if Points == strings.ToLower(tag) {
		return clip.points
	}
	return nil
}

func (clip *polygonClip) getRaw(tag string) interface{} {
	return clip.Get(tag)
}

func (clip *polygonClip) Set(tag string, value interface{}) bool {
	if Points == strings.ToLower(tag) {
		switch value := value.(type) {
		case []interface{}:
			result := true
			clip.points = make([]interface{}, len(value))
			for i, val := range value {
				switch val := val.(type) {
				case string:
					if isConstantName(val) {
						clip.points[i] = val
					} else if size, ok := StringToSizeUnit(val); ok {
						clip.points[i] = size
					} else {
						notCompatibleType(tag, val)
						result = false
					}

				case SizeUnit:
					clip.points[i] = val

				default:
					notCompatibleType(tag, val)
					clip.points[i] = AutoSize()
					result = false
				}
			}
			return result

		case []SizeUnit:
			clip.points = make([]interface{}, len(value))
			for i, point := range value {
				clip.points[i] = point
			}
			return true

		case string:
			result := true
			values := strings.Split(value, ",")
			clip.points = make([]interface{}, len(values))
			for i, val := range values {
				val = strings.Trim(val, " \t\n\r")
				if isConstantName(val) {
					clip.points[i] = val
				} else if size, ok := StringToSizeUnit(val); ok {
					clip.points[i] = size
				} else {
					notCompatibleType(tag, val)
					result = false
				}
			}
			return result
		}
	}
	return false
}

func (clip *polygonClip) setRaw(tag string, value interface{}) {
	clip.Set(tag, value)
}

func (clip *polygonClip) Remove(tag string) {
	if Points == strings.ToLower(tag) {
		clip.points = []interface{}{}
	}
}

func (clip *polygonClip) Clear() {
	clip.points = []interface{}{}
}

func (clip *polygonClip) AllTags() []string {
	return []string{Points}
}

func (clip *polygonClip) String() string {
	return runStringWriter(clip)
}

func (clip *polygonClip) writeString(buffer *strings.Builder, indent string) {

	buffer.WriteString("inset { ")

	if clip.points != nil {
		buffer.WriteString(Points)
		buffer.WriteString(` = "`)
		for i, value := range clip.points {
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

	count := len(clip.points)
	if count < 2 {
		return ""
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	writePoint := func(value interface{}) {
		switch value := value.(type) {
		case string:
			if val, ok := session.resolveConstants(value); ok {
				if size, ok := StringToSizeUnit(val); ok {
					buffer.WriteString(size.cssString("0px"))
					return
				}
			}

		case SizeUnit:
			buffer.WriteString(value.cssString("0px"))
			return
		}

		buffer.WriteString("0px")
	}

	leadText := "polygon("
	for i := 1; i < count; i += 2 {
		buffer.WriteString(leadText)
		writePoint(clip.points[i-1])
		buffer.WriteRune(' ')
		writePoint(clip.points[i])
		leadText = ", "
	}

	buffer.WriteRune(')')
	return buffer.String()
}

func (clip *polygonClip) valid(session Session) bool {
	if clip.points == nil || len(clip.points) == 0 {
		return false
	}

	return true
}

func parseClipShape(obj DataObject) ClipShape {
	switch obj.Tag() {
	case "inset":
		clip := new(insetClip)
		for _, tag := range []string{Top, Right, Bottom, Left, Radius, RadiusX, RadiusY} {
			if value, ok := obj.PropertyValue(tag); ok {
				clip.Set(tag, value)
			}
		}
		return clip

	case "circle":
		clip := new(ellipseClip)
		for _, tag := range []string{X, Y, Radius} {
			if value, ok := obj.PropertyValue(tag); ok {
				clip.Set(tag, value)
			}
		}
		return clip

	case "ellipse":
		clip := new(ellipseClip)
		for _, tag := range []string{X, Y, RadiusX, RadiusY} {
			if value, ok := obj.PropertyValue(tag); ok {
				clip.Set(tag, value)
			}
		}
		return clip

	case "polygon":
		clip := new(ellipseClip)
		if value, ok := obj.PropertyValue(Points); ok {
			clip.Set(Points, value)
		}
		return clip
	}

	return nil
}

func (style *viewStyle) setClipShape(tag string, value interface{}) bool {
	switch value := value.(type) {
	case ClipShape:
		style.properties[tag] = value
		return true

	case string:
		if isConstantName(value) {
			style.properties[tag] = value
			return true
		}

		if obj := NewDataObject(value); obj == nil {
			if clip := parseClipShape(obj); clip != nil {
				style.properties[tag] = clip
				return true
			}
		}

	case DataObject:
		if clip := parseClipShape(value); clip != nil {
			style.properties[tag] = clip
			return true
		}

	case DataValue:
		if value.IsObject() {
			if clip := parseClipShape(value.Object()); clip != nil {
				style.properties[tag] = clip
				return true
			}
		}
	}

	notCompatibleType(tag, value)
	return false
}

func getClipShape(prop Properties, tag string, session Session) ClipShape {
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
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetClip(view View, subviewID string) ClipShape {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		return getClipShape(view, Clip, view.Session())
	}

	return nil
}

// GetShapeOutside returns a shape around which adjacent inline content.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetShapeOutside(view View, subviewID string) ClipShape {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		return getClipShape(view, ShapeOutside, view.Session())
	}

	return nil
}
