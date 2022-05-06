package rui

import (
	"fmt"
	"strings"
)

const (
	// Radius is the SizeUnit view property that determines the corners rounding radius
	// of an element's outer border edge.
	Radius = "radius"
	// RadiusX is the SizeUnit view property that determines the x-axis corners elliptic rounding
	// radius of an element's outer border edge.
	RadiusX = "radius-x"
	// RadiusY is the SizeUnit view property that determines the y-axis corners elliptic rounding
	// radius of an element's outer border edge.
	RadiusY = "radius-y"
	// RadiusTopLeft is the SizeUnit view property that determines the top-left corner rounding radius
	// of an element's outer border edge.
	RadiusTopLeft = "radius-top-left"
	// RadiusTopLeftX is the SizeUnit view property that determines the x-axis top-left corner elliptic
	// rounding radius of an element's outer border edge.
	RadiusTopLeftX = "radius-top-left-x"
	// RadiusTopLeftY is the SizeUnit view property that determines the y-axis top-left corner elliptic
	// rounding radius of an element's outer border edge.
	RadiusTopLeftY = "radius-top-left-y"
	// RadiusTopRight is the SizeUnit view property that determines the top-right corner rounding radius
	// of an element's outer border edge.
	RadiusTopRight = "radius-top-right"
	// RadiusTopRightX is the SizeUnit view property that determines the x-axis top-right corner elliptic
	// rounding radius of an element's outer border edge.
	RadiusTopRightX = "radius-top-right-x"
	// RadiusTopRightY is the SizeUnit view property that determines the y-axis top-right corner elliptic
	// rounding radius of an element's outer border edge.
	RadiusTopRightY = "radius-top-right-y"
	// RadiusBottomLeft is the SizeUnit view property that determines the bottom-left corner rounding radius
	// of an element's outer border edge.
	RadiusBottomLeft = "radius-bottom-left"
	// RadiusBottomLeftX is the SizeUnit view property that determines the x-axis bottom-left corner elliptic
	// rounding radius of an element's outer border edge.
	RadiusBottomLeftX = "radius-bottom-left-x"
	// RadiusBottomLeftY is the SizeUnit view property that determines the y-axis bottom-left corner elliptic
	// rounding radius of an element's outer border edge.
	RadiusBottomLeftY = "radius-bottom-left-y"
	// RadiusBottomRight is the SizeUnit view property that determines the bottom-right corner rounding radius
	// of an element's outer border edge.
	RadiusBottomRight = "radius-bottom-right"
	// RadiusBottomRightX is the SizeUnit view property that determines the x-axis bottom-right corner elliptic
	// rounding radius of an element's outer border edge.
	RadiusBottomRightX = "radius-bottom-right-x"
	// RadiusBottomRightY is the SizeUnit view property that determines the y-axis bottom-right corner elliptic
	// rounding radius of an element's outer border edge.
	RadiusBottomRightY = "radius-bottom-right-y"
	// X is the SizeUnit property of the ShadowProperty that determines the x-axis corners elliptic rounding
	// radius of an element's outer border edge.
	X = "x"
	// Y is the SizeUnit property of the ShadowProperty that determines the y-axis corners elliptic rounding
	// radius of an element's outer border edge.
	Y = "y"
	// TopLeft is the SizeUnit property of the ShadowProperty that determines the top-left corner rounding radius
	// of an element's outer border edge.
	TopLeft = "top-left"
	// TopLeftX is the SizeUnit property of the ShadowProperty that determines the x-axis top-left corner elliptic
	// rounding radius of an element's outer border edge.
	TopLeftX = "top-left-x"
	// TopLeftY is the SizeUnit property of the ShadowProperty that determines the y-axis top-left corner elliptic
	// rounding radius of an element's outer border edge.
	TopLeftY = "top-left-y"
	// TopRight is the SizeUnit property of the ShadowProperty that determines the top-right corner rounding radius
	// of an element's outer border edge.
	TopRight = "top-right"
	// TopRightX is the SizeUnit property of the ShadowProperty that determines the x-axis top-right corner elliptic
	// rounding radius of an element's outer border edge.
	TopRightX = "top-right-x"
	// TopRightY is the SizeUnit property of the ShadowProperty that determines the y-axis top-right corner elliptic
	// rounding radius of an element's outer border edge.
	TopRightY = "top-right-y"
	// BottomLeft is the SizeUnit property of the ShadowProperty that determines the bottom-left corner rounding radius
	// of an element's outer border edge.
	BottomLeft = "bottom-left"
	// BottomLeftX is the SizeUnit property of the ShadowProperty that determines the x-axis bottom-left corner elliptic
	// rounding radius of an element's outer border edge.
	BottomLeftX = "bottom-left-x"
	// BottomLeftY is the SizeUnit property of the ShadowProperty that determines the y-axis bottom-left corner elliptic
	// rounding radius of an element's outer border edge.
	BottomLeftY = "bottom-left-y"
	// BottomRight is the SizeUnit property of the ShadowProperty that determines the bottom-right corner rounding radius
	// of an element's outer border edge.
	BottomRight = "bottom-right"
	// BottomRightX is the SizeUnit property of the ShadowProperty that determines the x-axis bottom-right corner elliptic
	// rounding radius of an element's outer border edge.
	BottomRightX = "bottom-right-x"
	// BottomRightY is the SizeUnit property of the ShadowProperty that determines the y-axis bottom-right corner elliptic
	// rounding radius of an element's outer border edge.
	BottomRightY = "bottom-right-y"
)

type RadiusProperty interface {
	Properties
	ruiStringer
	fmt.Stringer
	BoxRadius(session Session) BoxRadius
}

type radiusPropertyData struct {
	propertyList
}

// NewRadiusProperty creates the new RadiusProperty
func NewRadiusProperty(params Params) RadiusProperty {
	result := new(radiusPropertyData)
	result.properties = map[string]interface{}{}
	if params != nil {
		for _, tag := range []string{X, Y, TopLeft, TopRight, BottomLeft, BottomRight, TopLeftX, TopLeftY,
			TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY} {
			if value, ok := params[tag]; ok {
				result.Set(tag, value)
			}
		}
	}
	return result
}

func (radius *radiusPropertyData) normalizeTag(tag string) string {
	return strings.TrimPrefix(strings.ToLower(tag), "radius-")
}

func (radius *radiusPropertyData) ruiString(writer ruiWriter) {
	writer.startObject("_")

	for _, tag := range []string{X, Y, TopLeft, TopLeftX, TopLeftY, TopRight, TopRightX, TopRightY,
		BottomLeft, BottomLeftX, BottomLeftY, BottomRight, BottomRightX, BottomRightY} {
		if value, ok := radius.properties[tag]; ok {
			writer.writeProperty(Style, value)
		}
	}

	writer.endObject()
}

func (radius *radiusPropertyData) String() string {
	writer := newRUIWriter()
	radius.ruiString(writer)
	return writer.finish()
}

func (radius *radiusPropertyData) delete(tags []string) {
	for _, tag := range tags {
		delete(radius.properties, tag)
	}
}

func (radius *radiusPropertyData) deleteUnusedTags() {
	for _, tag := range []string{X, Y} {
		if _, ok := radius.properties[tag]; ok {
			unused := true
			for _, t := range []string{TopLeft, TopRight, BottomLeft, BottomRight} {
				if _, ok := radius.properties[t+"-"+tag]; !ok {
					if _, ok := radius.properties[t]; !ok {
						unused = false
						break
					}
				}
			}
			if unused {
				delete(radius.properties, tag)
			}
		}
	}

	equalValue := func(value1, value2 interface{}) bool {
		switch value1 := value1.(type) {
		case string:
			switch value2 := value2.(type) {
			case string:
				return value1 == value2
			}

		case SizeUnit:
			switch value2 := value2.(type) {
			case SizeUnit:
				return value1.Equal(value2)
			}
		}
		return false
	}

	for _, tag := range []string{TopLeft, TopRight, BottomLeft, BottomRight} {
		tagX := tag + "-x"
		tagY := tag + "-y"
		valueX, okX := radius.properties[tagX]
		valueY, okY := radius.properties[tagY]

		if value, ok := radius.properties[tag]; ok {
			if okX && okY {
				delete(radius.properties, tag)
			} else if okX && !okY {
				if equalValue(value, valueX) {
					delete(radius.properties, tagX)
				} else {
					radius.properties[tagY] = value
					delete(radius.properties, tag)
				}
			} else if !okX && okY {
				if equalValue(value, valueY) {
					delete(radius.properties, tagY)
				} else {
					radius.properties[tagX] = value
					delete(radius.properties, tag)
				}
			}
		} else if okX && okY && equalValue(valueX, valueY) {
			radius.properties[tag] = valueX
			delete(radius.properties, tagX)
			delete(radius.properties, tagY)
		}
	}
}

func (radius *radiusPropertyData) Remove(tag string) {
	tag = radius.normalizeTag(tag)

	switch tag {
	case X, Y:
		if _, ok := radius.properties[tag]; ok {
			radius.Set(tag, AutoSize())
			delete(radius.properties, tag)
		}

	case TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY:
		delete(radius.properties, tag)

	case TopLeft, TopRight, BottomLeft, BottomRight:
		radius.delete([]string{tag, tag + "-x", tag + "-y"})

	default:
		ErrorLogF(`"%s" property is not compatible with the RadiusProperty`, tag)
	}

}

func (radius *radiusPropertyData) Set(tag string, value interface{}) bool {
	if value == nil {
		radius.Remove(tag)
		return true
	}

	tag = radius.normalizeTag(tag)
	switch tag {
	case X:
		if radius.setSizeProperty(tag, value) {
			radius.delete([]string{TopLeftX, TopRightX, BottomLeftX, BottomRightX})
			for _, t := range []string{TopLeft, TopRight, BottomLeft, BottomRight} {
				if val, ok := radius.properties[t]; ok {
					if _, ok := radius.properties[t+"-y"]; !ok {
						radius.properties[t+"-y"] = val
					}
					delete(radius.properties, t)
				}
			}
			return true
		}

	case Y:
		if radius.setSizeProperty(tag, value) {
			radius.delete([]string{TopLeftY, TopRightY, BottomLeftY, BottomRightY})
			for _, t := range []string{TopLeft, TopRight, BottomLeft, BottomRight} {
				if val, ok := radius.properties[t]; ok {
					if _, ok := radius.properties[t+"-x"]; !ok {
						radius.properties[t+"-x"] = val
					}
					delete(radius.properties, t)
				}
			}
			return true
		}

	case TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY:
		if radius.setSizeProperty(tag, value) {
			radius.deleteUnusedTags()
			return true
		}

	case TopLeft, TopRight, BottomLeft, BottomRight:
		switch value := value.(type) {
		case SizeUnit:
			radius.properties[tag] = value
			radius.delete([]string{tag + "-x", tag + "-y"})
			radius.deleteUnusedTags()
			return true

		case string:
			if strings.Contains(value, "/") {
				if values := strings.Split(value, "/"); len(values) == 2 {
					xOK := radius.Set(tag+"-x", value[0])
					yOK := radius.Set(tag+"-y", value[1])
					return xOK && yOK
				} else {
					notCompatibleType(tag, value)
				}
			} else {
				if radius.setSizeProperty(tag, value) {
					radius.delete([]string{tag + "-x", tag + "-y"})
					radius.deleteUnusedTags()
					return true
				}
			}
		}

	default:
		ErrorLogF(`"%s" property is not compatible with the RadiusProperty`, tag)
	}

	return false
}

func (radius *radiusPropertyData) Get(tag string) interface{} {
	tag = radius.normalizeTag(tag)
	if value, ok := radius.properties[tag]; ok {
		return value
	}

	switch tag {
	case TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY:
		tagLen := len(tag)
		if value, ok := radius.properties[tag[:tagLen-2]]; ok {
			return value
		}
		if value, ok := radius.properties[tag[tagLen-1:]]; ok {
			return value
		}
	}

	switch tag {
	case TopLeftX, TopRightX, BottomLeftX, BottomRightX:
		if value, ok := radius.properties[X]; ok {
			return value
		}
	case TopLeftY, TopRightY, BottomLeftY, BottomRightY:
		if value, ok := radius.properties[Y]; ok {
			return value
		}
	}

	return nil
}

func (radius *radiusPropertyData) BoxRadius(session Session) BoxRadius {
	x, _ := sizeProperty(radius, X, session)
	y, _ := sizeProperty(radius, Y, session)

	getRadius := func(tag string) (SizeUnit, SizeUnit) {
		rx := x
		ry := y
		if r, ok := sizeProperty(radius, tag, session); ok {
			rx = r
			ry = r
		}
		if r, ok := sizeProperty(radius, tag+"-x", session); ok {
			rx = r
		}
		if r, ok := sizeProperty(radius, tag+"-y", session); ok {
			ry = r
		}

		return rx, ry
	}

	var result BoxRadius

	result.TopLeftX, result.TopLeftY = getRadius(TopLeft)
	result.TopRightX, result.TopRightY = getRadius(TopRight)
	result.BottomLeftX, result.BottomLeftY = getRadius(BottomLeft)
	result.BottomRightX, result.BottomRightY = getRadius(BottomRight)

	return result
}

// BoxRadius defines radii of rounds the corners of an element's outer border edge
type BoxRadius struct {
	TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY SizeUnit
}

// AllAnglesIsEqual returns 'true' if all angles is equal, 'false' otherwise
func (radius BoxRadius) AllAnglesIsEqual() bool {
	return radius.TopLeftX.Equal(radius.TopRightX) &&
		radius.TopLeftY.Equal(radius.TopRightY) &&
		radius.TopLeftX.Equal(radius.BottomLeftX) &&
		radius.TopLeftY.Equal(radius.BottomLeftY) &&
		radius.TopLeftX.Equal(radius.BottomRightX) &&
		radius.TopLeftY.Equal(radius.BottomRightY)
}

// String returns a string representation of a BoxRadius struct
func (radius BoxRadius) String() string {

	if radius.AllAnglesIsEqual() {
		if radius.TopLeftX.Equal(radius.TopLeftY) {
			return radius.TopLeftX.String()
		} else {
			return fmt.Sprintf("_{ x = %s, y = %s }", radius.TopLeftX.String(), radius.TopLeftY.String())
		}
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString("_{ ")

	if radius.TopLeftX.Equal(radius.TopLeftY) {
		buffer.WriteString("top-left = ")
		buffer.WriteString(radius.TopLeftX.String())
	} else {
		buffer.WriteString("top-left-x = ")
		buffer.WriteString(radius.TopLeftX.String())
		buffer.WriteString("top-left-y = ")
		buffer.WriteString(radius.TopLeftY.String())
	}

	if radius.TopRightX.Equal(radius.TopRightY) {
		buffer.WriteString(", top-right = ")
		buffer.WriteString(radius.TopRightX.String())
	} else {
		buffer.WriteString(", top-right-x = ")
		buffer.WriteString(radius.TopRightX.String())
		buffer.WriteString(", top-right-y = ")
		buffer.WriteString(radius.TopRightY.String())
	}

	if radius.BottomLeftX.Equal(radius.BottomLeftY) {
		buffer.WriteString(", bottom-left = ")
		buffer.WriteString(radius.BottomLeftX.String())
	} else {
		buffer.WriteString(", bottom-left-x = ")
		buffer.WriteString(radius.BottomLeftX.String())
		buffer.WriteString(", bottom-left-y = ")
		buffer.WriteString(radius.BottomLeftY.String())
	}

	if radius.BottomRightX.Equal(radius.BottomRightY) {
		buffer.WriteString(", bottom-right = ")
		buffer.WriteString(radius.BottomRightX.String())
	} else {
		buffer.WriteString(", bottom-right-x = ")
		buffer.WriteString(radius.BottomRightX.String())
		buffer.WriteString(", bottom-right-y = ")
		buffer.WriteString(radius.BottomRightY.String())
	}

	buffer.WriteString(" }")
	return buffer.String()
}

func (radius BoxRadius) cssValue(builder cssBuilder) {

	if (radius.TopLeftX.Type == Auto || radius.TopLeftX.Value == 0) &&
		(radius.TopLeftY.Type == Auto || radius.TopLeftY.Value == 0) &&
		(radius.TopRightX.Type == Auto || radius.TopRightX.Value == 0) &&
		(radius.TopRightY.Type == Auto || radius.TopRightY.Value == 0) &&
		(radius.BottomRightX.Type == Auto || radius.BottomRightX.Value == 0) &&
		(radius.BottomRightY.Type == Auto || radius.BottomRightY.Value == 0) &&
		(radius.BottomLeftX.Type == Auto || radius.BottomLeftX.Value == 0) &&
		(radius.BottomLeftY.Type == Auto || radius.BottomLeftY.Value == 0) {
		return
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(radius.TopLeftX.cssString("0"))

	if radius.AllAnglesIsEqual() {

		if !radius.TopLeftX.Equal(radius.TopLeftY) {
			buffer.WriteString(" / ")
			buffer.WriteString(radius.TopLeftY.cssString("0"))
		}

	} else {

		buffer.WriteRune(' ')
		buffer.WriteString(radius.TopRightX.cssString("0"))
		buffer.WriteRune(' ')
		buffer.WriteString(radius.BottomRightX.cssString("0"))
		buffer.WriteRune(' ')
		buffer.WriteString(radius.BottomLeftX.cssString("0"))

		if !radius.TopLeftX.Equal(radius.TopLeftY) ||
			!radius.TopRightX.Equal(radius.TopRightY) ||
			!radius.BottomLeftX.Equal(radius.BottomLeftY) ||
			!radius.BottomRightX.Equal(radius.BottomRightY) {

			buffer.WriteString(" / ")
			buffer.WriteString(radius.TopLeftY.cssString("0"))
			buffer.WriteRune(' ')
			buffer.WriteString(radius.TopRightY.cssString("0"))
			buffer.WriteRune(' ')
			buffer.WriteString(radius.BottomRightY.cssString("0"))
			buffer.WriteRune(' ')
			buffer.WriteString(radius.BottomLeftY.cssString("0"))
		}
	}

	builder.add("border-radius", buffer.String())
}

func (radius BoxRadius) cssString() string {
	var builder cssValueBuilder
	radius.cssValue(&builder)
	return builder.finish()
}

func getRadiusProperty(style Properties) RadiusProperty {
	if value := style.Get(Radius); value != nil {
		switch value := value.(type) {
		case RadiusProperty:
			return value

		case BoxRadius:
			result := NewRadiusProperty(nil)
			if value.AllAnglesIsEqual() {
				result.Set(X, value.TopLeftX)
				result.Set(Y, value.TopLeftY)
			} else {
				if value.TopLeftX.Equal(value.TopLeftY) {
					result.Set(TopLeft, value.TopLeftX)
				} else {
					result.Set(TopLeftX, value.TopLeftX)
					result.Set(TopLeftY, value.TopLeftY)
				}
				if value.TopRightX.Equal(value.TopRightY) {
					result.Set(TopRight, value.TopRightX)
				} else {
					result.Set(TopRightX, value.TopRightX)
					result.Set(TopRightY, value.TopRightY)
				}
				if value.BottomLeftX.Equal(value.BottomLeftY) {
					result.Set(BottomLeft, value.BottomLeftX)
				} else {
					result.Set(BottomLeftX, value.BottomLeftX)
					result.Set(BottomLeftY, value.BottomLeftY)
				}
				if value.BottomRightX.Equal(value.BottomRightY) {
					result.Set(BottomRight, value.BottomRightX)
				} else {
					result.Set(BottomRightX, value.BottomRightX)
					result.Set(BottomRightY, value.BottomRightY)
				}
			}
			return result

		case SizeUnit:
			return NewRadiusProperty(Params{
				X: value,
				Y: value,
			})

		case string:
			return NewRadiusProperty(Params{
				X: value,
				Y: value,
			})
		}
	}

	return NewRadiusProperty(nil)
}

func (properties *propertyList) setRadius(value interface{}) bool {

	if value == nil {
		delete(properties.properties, Radius)
		return true
	}

	switch value := value.(type) {
	case RadiusProperty:
		properties.properties[Radius] = value
		return true

	case SizeUnit:
		properties.properties[Radius] = value
		return true

	case BoxRadius:
		radius := NewRadiusProperty(nil)
		if value.AllAnglesIsEqual() {
			radius.Set(X, value.TopLeftX)
			radius.Set(Y, value.TopLeftY)
		} else {
			if value.TopLeftX.Equal(value.TopLeftY) {
				radius.Set(TopLeft, value.TopLeftX)
			} else {
				radius.Set(TopLeftX, value.TopLeftX)
				radius.Set(TopLeftY, value.TopLeftY)
			}
			if value.TopRightX.Equal(value.TopRightY) {
				radius.Set(TopRight, value.TopRightX)
			} else {
				radius.Set(TopRightX, value.TopRightX)
				radius.Set(TopRightY, value.TopRightY)
			}
			if value.BottomLeftX.Equal(value.BottomLeftY) {
				radius.Set(BottomLeft, value.BottomLeftX)
			} else {
				radius.Set(BottomLeftX, value.BottomLeftX)
				radius.Set(BottomLeftY, value.BottomLeftY)
			}
			if value.BottomRightX.Equal(value.BottomRightY) {
				radius.Set(BottomRight, value.BottomRightX)
			} else {
				radius.Set(BottomRightX, value.BottomRightX)
				radius.Set(BottomRightY, value.BottomRightY)
			}
		}
		properties.properties[Radius] = radius
		return true

	case string:
		if strings.Contains(value, "/") {
			values := strings.Split(value, "/")
			if len(values) == 2 {
				okX := properties.setRadiusElement(RadiusX, values[0])
				okY := properties.setRadiusElement(RadiusY, values[1])
				return okX && okY
			} else {
				notCompatibleType(Radius, value)
			}
		} else {
			return properties.setSizeProperty(Radius, value)
		}

	case DataObject:
		radius := NewRadiusProperty(nil)
		for _, tag := range []string{X, Y, TopLeft, TopRight, BottomLeft, BottomRight, TopLeftX, TopLeftY,
			TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY} {
			if value, ok := value.PropertyValue(tag); ok {
				radius.Set(tag, value)
			}
		}
		properties.properties[Radius] = radius
		return true

	default:
		notCompatibleType(Radius, value)
	}

	return false
}

func (properties *propertyList) removeRadiusElement(tag string) {
	if value, ok := properties.properties[Radius]; ok && value != nil {
		radius := getRadiusProperty(properties)
		radius.Remove(tag)
		if len(radius.AllTags()) == 0 {
			delete(properties.properties, Radius)
		} else {
			properties.properties[Radius] = radius
		}
	}
}

func (properties *propertyList) setRadiusElement(tag string, value interface{}) bool {
	if value == nil {
		properties.removeRadiusElement(tag)
		return true
	}

	radius := getRadiusProperty(properties)
	if radius.Set(tag, value) {
		properties.properties[Radius] = radius
		return true
	}

	return false
}

func getRadiusElement(style Properties, tag string) interface{} {
	value := style.Get(Radius)
	if value != nil {
		switch value := value.(type) {
		case string:
			return value

		case SizeUnit:
			return value

		case RadiusProperty:
			return value.Get(tag)

		case BoxRadius:
			switch tag {
			case RadiusX:
				if value.TopLeftX.Equal(value.TopRightX) &&
					value.TopLeftX.Equal(value.BottomLeftX) &&
					value.TopLeftX.Equal(value.BottomRightX) {
					return value.TopLeftX
				}

			case RadiusY:
				if value.TopLeftY.Equal(value.TopRightY) &&
					value.TopLeftY.Equal(value.BottomLeftY) &&
					value.TopLeftY.Equal(value.BottomRightY) {
					return value.TopLeftY
				}

			case RadiusTopLeft:
				if value.TopLeftX.Equal(value.TopLeftY) {
					return value.TopLeftY
				}

			case RadiusTopRight:
				if value.TopRightX.Equal(value.TopRightY) {
					return value.TopRightY
				}

			case RadiusBottomLeft:
				if value.BottomLeftX.Equal(value.BottomLeftY) {
					return value.BottomLeftY
				}

			case RadiusBottomRight:
				if value.BottomRightX.Equal(value.BottomRightY) {
					return value.BottomRightY
				}

			case RadiusTopLeftX:
				return value.TopLeftX

			case RadiusTopLeftY:
				return value.TopLeftY

			case RadiusTopRightX:
				return value.TopRightX

			case RadiusTopRightY:
				return value.TopRightY

			case RadiusBottomLeftX:
				return value.BottomLeftX

			case RadiusBottomLeftY:
				return value.BottomLeftY

			case RadiusBottomRightX:
				return value.BottomRightX

			case RadiusBottomRightY:
				return value.BottomRightY
			}
		}
	}

	return nil
}

func getRadius(properties Properties, session Session) BoxRadius {
	if value := properties.Get(Radius); value != nil {
		switch value := value.(type) {
		case BoxRadius:
			return value

		case RadiusProperty:
			return value.BoxRadius(session)

		case SizeUnit:
			return BoxRadius{TopLeftX: value, TopLeftY: value, TopRightX: value, TopRightY: value,
				BottomLeftX: value, BottomLeftY: value, BottomRightX: value, BottomRightY: value}

		case string:
			if text, ok := session.resolveConstants(value); ok {
				if size, ok := StringToSizeUnit(text); ok {
					return BoxRadius{TopLeftX: size, TopLeftY: size, TopRightX: size, TopRightY: size,
						BottomLeftX: size, BottomLeftY: size, BottomRightX: size, BottomRightY: size}
				}
			}
		}
	}

	return BoxRadius{}
}
