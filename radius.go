package rui

import (
	"fmt"
	"strings"
)

// Constants for [RadiusProperty] specific properties
const (
	// Radius is the constant for "radius" property tag.
	//
	// Used by View, BackgroundElement, ClipShapeProperty.
	//
	// Usage in View:
	// Specifies the corners rounding radius of an element's outer border edge.
	//
	// Supported types: RadiusProperty, SizeUnit, SizeFunc, BoxRadius, string, float, int.
	//
	// Internal type is either RadiusProperty or SizeUnit, other types converted to them during assignment.
	// See RadiusProperty, SizeUnit, SizeFunc and BoxRadius description for more details.
	//
	// Conversion rules:
	//   - RadiusProperty - stored as is, no conversion performed.
	//   - SizeUnit - stored as is and set all corners to have the same value.
	//   - BoxRadius - a new RadiusProperty will be created and all corresponding elliptical radius values will be set.
	//   - string - if one value will be provided then it will be set as a radius for all corners. If two values will be provided divided by (/) then x and y radius will be set for all corners. Examples: "1em", "1em/0.5em", "2/4". Values which doesn't have size prefix will use size in pixels by default.
	//   - float - values of this type will set radius for all corners in pixels.
	//   - int - values of this type will set radius for all corners in pixels.
	//
	// Usage in BackgroundElement:
	// Same as "radial-gradient-radius".
	//
	// Usage in ClipShapeProperty:
	// Specifies the radius of the corners or the radius of the cropping area.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	Radius PropertyName = "radius"

	// RadiusX is the constant for "radius-x" property tag.
	//
	// Used by View, ClipShapeProperty.
	//
	// Usage in View:
	// Specifies the x-axis corners elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	//
	// Usage in ClipShapeProperty:
	// Specifies the x-axis corners elliptic rounding radius of the elliptic clip shape.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusX PropertyName = "radius-x"

	// RadiusY is the constant for "radius-y" property tag.
	//
	// Used by View, ClipShapeProperty.
	//
	// Usage in View:
	// Specifies the y-axis corners elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	//
	// Usage in ClipShapeProperty:
	// Specifies the y-axis corners elliptic rounding radius of of the elliptic clip shape.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusY PropertyName = "radius-y"

	// RadiusTopLeft is the constant for "radius-top-left" property tag.
	//
	// Used by View.
	// Specifies the top-left corner rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusTopLeft PropertyName = "radius-top-left"

	// RadiusTopLeftX is the constant for "radius-top-left-x" property tag.
	//
	// Used by View.
	// Specifies the x-axis top-left corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusTopLeftX PropertyName = "radius-top-left-x"

	// RadiusTopLeftY is the constant for "radius-top-left-y" property tag.
	//
	// Used by View.
	// Specifies the y-axis top-left corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusTopLeftY PropertyName = "radius-top-left-y"

	// RadiusTopRight is the constant for "radius-top-right" property tag.
	//
	// Used by View.
	// Specifies the top-right corner rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusTopRight PropertyName = "radius-top-right"

	// RadiusTopRightX is the constant for "radius-top-right-x" property tag.
	//
	// Used by View.
	// Specifies the x-axis top-right corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusTopRightX PropertyName = "radius-top-right-x"

	// RadiusTopRightY is the constant for "radius-top-right-y" property tag.
	//
	// Used by View.
	// Specifies the y-axis top-right corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusTopRightY PropertyName = "radius-top-right-y"

	// RadiusBottomLeft is the constant for "radius-bottom-left" property tag.
	//
	// Used by View.
	// Specifies the bottom-left corner rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusBottomLeft PropertyName = "radius-bottom-left"

	// RadiusBottomLeftX is the constant for "radius-bottom-left-x" property tag.
	//
	// Used by View.
	// Specifies the x-axis bottom-left corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusBottomLeftX PropertyName = "radius-bottom-left-x"

	// RadiusBottomLeftY is the constant for "radius-bottom-left-y" property tag.
	//
	// Used by View.
	// Specifies the y-axis bottom-left corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusBottomLeftY PropertyName = "radius-bottom-left-y"

	// RadiusBottomRight is the constant for "radius-bottom-right" property tag.
	//
	// Used by View.
	// Specifies the bottom-right corner rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusBottomRight PropertyName = "radius-bottom-right"

	// RadiusBottomRightX is the constant for "radius-bottom-right-x" property tag.
	//
	// Used by View.
	// Specifies the x-axis bottom-right corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusBottomRightX PropertyName = "radius-bottom-right-x"

	// RadiusBottomRightY is the constant for "radius-bottom-right-y" property tag.
	//
	// Used by View.
	// Specifies the y-axis bottom-right corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	RadiusBottomRightY PropertyName = "radius-bottom-right-y"

	// X is the constant for "x" property tag.
	//
	// Used by ClipShapeProperty, RadiusProperty.
	//
	// Usage in ClipShapeProperty:
	// Specifies x-axis position of the clip shape center.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	//
	// Usage in RadiusProperty:
	// Determines the x-axis elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	X PropertyName = "x"

	// Y is the constant for "y" property tag.
	//
	// Used by ClipShapeProperty, RadiusProperty.
	//
	// Usage in ClipShapeProperty:
	// Specifies y-axis position of the clip shape center.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	//
	// Usage in RadiusProperty:
	// Determines the y-axis elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	Y PropertyName = "y"

	// TopLeft is the constant for "top-left" property tag.
	//
	// Used by RadiusProperty.
	// Determines the top-left corner rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TopLeft PropertyName = "top-left"

	// TopLeftX is the constant for "top-left-x" property tag.
	//
	// Used by RadiusProperty.
	// Determines the x-axis top-left corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TopLeftX PropertyName = "top-left-x"

	// TopLeftY is the constant for "top-left-y" property tag.
	//
	// Used by RadiusProperty.
	// Determines the y-axis top-left corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TopLeftY PropertyName = "top-left-y"

	// TopRight is the constant for "top-right" property tag.
	//
	// Used by RadiusProperty.
	// Determines the top-right corner rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TopRight PropertyName = "top-right"

	// TopRightX is the constant for "top-right-x" property tag.
	//
	// Used by RadiusProperty.
	// Determines the x-axis top-right corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TopRightX PropertyName = "top-right-x"

	// TopRightY is the constant for "top-right-y" property tag.
	//
	// Used by RadiusProperty.
	// Determines the y-axis top-right corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TopRightY PropertyName = "top-right-y"

	// BottomLeft is the constant for "bottom-left" property tag.
	//
	// Used by RadiusProperty.
	// Determines the bottom-left corner rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	BottomLeft PropertyName = "bottom-left"

	// BottomLeftX is the constant for "bottom-left-x" property tag.
	//
	// Used by RadiusProperty.
	// Determines the x-axis bottom-left corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	BottomLeftX PropertyName = "bottom-left-x"

	// BottomLeftY is the constant for "bottom-left-y" property tag.
	//
	// Used by RadiusProperty.
	// Determines the y-axis bottom-left corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	BottomLeftY PropertyName = "bottom-left-y"

	// BottomRight is the constant for "bottom-right" property tag.
	//
	// Used by RadiusProperty.
	// Determines the bottom-right corner rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	BottomRight PropertyName = "bottom-right"

	// BottomRightX is the constant for "bottom-right-x" property tag.
	//
	// Used by RadiusProperty.
	// Determines the x-axis bottom-right corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	BottomRightX PropertyName = "bottom-right-x"

	// BottomRightY is the constant for "bottom-right-y" property tag.
	//
	// Used by RadiusProperty.
	// Determines the y-axis bottom-right corner elliptic rounding radius of an element's outer border edge.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	BottomRightY PropertyName = "bottom-right-y"
)

// RadiusProperty is a description of the [View] (shape) elliptical corner radius.
type RadiusProperty interface {
	Properties
	stringWriter
	fmt.Stringer

	// BoxRadius returns x and y radius of the corners of the element
	BoxRadius(session Session) BoxRadius
}

type radiusPropertyData struct {
	dataProperty
}

// NewRadiusProperty creates the new RadiusProperty
//
// The following properties can be used:
//   - "x" (X) - Determines the x-axis elliptic rounding radius of an element's outer border edge.
//   - "y" (Y) - Determines the y-axis corner elliptic rounding radius of an element's outer border edge.
//   - "top-left" (TopLeft) - Determines the top-left corner rounding radius of an element's outer border edge.
//   - "top-left-x" (TopLeftX) - Determines the x-axis top-left corner elliptic rounding radius of an element's outer border edge.
//   - "top-left-y" (TopLeftY) - Determines the y-axis top-left corner elliptic rounding radius of an element's outer border edge.
//   - "top-right" (TopRight) -  Determines the top-right corner rounding radius of an element's outer border edge.
//   - "top-right-x" (TopRightX) - Determines the x-axis top-right corner elliptic rounding radius of an element's outer border edge.
//   - "top-right-y" (TopRightY) - Determines the y-axis top-right corner elliptic rounding radius of an element's outer border edge.
//   - "bottom-left" (BottomLeft) - Determines the bottom-left corner rounding radius of an element's outer border edge.
//   - "bottom-left-x" (BottomLeftX) - Determines the x-axis bottom-left corner elliptic rounding radius of an element's outer border edge.
//   - "bottom-left-y" (BottomLeftY) -  Determines the y-axis bottom-left corner elliptic rounding radius of an element's outer border edge.
//   - "bottom-right" (BottomRight) - Determines the bottom-right corner rounding radius of an element's outer border edge.
//   - "bottom-right-x" (BottomRightX) - Determines the x-axis bottom-right corner elliptic rounding radius of an element's outer border edge.
//   - "bottom-right-y" (BottomRightY) - Determines the y-axis bottom-right corner elliptic rounding radius of an element's outer border edge.
func NewRadiusProperty(params Params) RadiusProperty {
	result := new(radiusPropertyData)
	result.init()

	if params != nil {
		for _, tag := range result.supportedProperties {
			if value, ok := params[tag]; ok {
				radiusPropertySet(result, tag, value)
			}
		}
	}
	return result
}

// NewRadiusProperty creates the new RadiusProperty which having the same elliptical radii for all angles.
//
// Arguments determines the x- and y-axis elliptic rounding radius. if an argument is specified as int or float64, the value is considered to be in pixels
func NewEllipticRadius[xType SizeUnit | int | float64, yType SizeUnit | int | float64](x xType, y yType) RadiusProperty {
	return NewRadiusProperty(Params{
		X: x,
		Y: y,
	})
}

// NewRadius creates the new RadiusProperty.
//
// The arguments specify the radii in a clockwise direction: "top-right", "bottom-right", "bottom-left", and "top-left".
//
// if an argument is specified as int or float64, the value is considered to be in pixels
func NewRadii[topRightType SizeUnit | int | float64, bottomRightType SizeUnit | int | float64, bottomLeftType SizeUnit | int | float64, topLeftType SizeUnit | int | float64](
	topRight topRightType, bottomRight bottomRightType, bottomLeft bottomLeftType, topLeft topLeftType) RadiusProperty {
	return NewRadiusProperty(Params{
		TopRight:    topRight,
		BottomRight: bottomRight,
		BottomLeft:  bottomLeft,
		TopLeft:     topLeft,
	})
}

func radiusPropertyNormalize(tag PropertyName) PropertyName {
	name := strings.TrimPrefix(strings.ToLower(string(tag)), "radius-")
	return PropertyName(name)
}

func (radius *radiusPropertyData) init() {
	radius.dataProperty.init()
	radius.normalize = radiusPropertyNormalize
	radius.get = radiusPropertyGet
	radius.remove = radiusPropertyRemove
	radius.set = radiusPropertySet
	radius.supportedProperties = []PropertyName{
		X, Y, TopLeft, TopRight, BottomLeft, BottomRight, TopLeftX, TopLeftY,
		TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY,
	}
}

func (radius *radiusPropertyData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("_{ ")
	comma := false
	for _, tag := range radius.supportedProperties {
		if value, ok := radius.properties[tag]; ok {
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

func (radius *radiusPropertyData) String() string {
	return runStringWriter(radius)
}

func radiusPropertyRemove(properties Properties, tag PropertyName) []PropertyName {
	result := []PropertyName{}
	removeTag := func(tag PropertyName) {
		if properties.getRaw(tag) != nil {
			properties.setRaw(tag, nil)
			result = append(result, tag)
		}
	}

	switch tag {
	case X, Y:
		if properties.getRaw(tag) == nil {
			for _, prefix := range []PropertyName{TopLeft, TopRight, BottomLeft, BottomRight} {
				removeTag(prefix + "-" + tag)
			}
		} else {
			removeTag(tag)
		}

	case TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY:
		removeTag(tag)

	case TopLeft, TopRight, BottomLeft, BottomRight:
		for _, tag := range []PropertyName{tag, tag + "-x", tag + "-y"} {
			removeTag(tag)
		}

	default:
		ErrorLogF(`"%s" property is not compatible with the RadiusProperty`, tag)
	}

	return result
}

func deleteRadiusUnusedTags(radius Properties, result []PropertyName) {

	for _, tag := range []PropertyName{X, Y} {
		if radius.getRaw(tag) != nil {
			unused := true
			for _, t := range []PropertyName{TopLeft, TopRight, BottomLeft, BottomRight} {
				if radius.getRaw(t+"-"+tag) == nil && radius.getRaw(t) == nil {
					unused = false
					break
				}
			}
			if unused {
				radius.setRaw(tag, nil)
				result = append(result, tag)
			}
		}
	}

	equalValue := func(value1, value2 any) bool {
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

	for _, tag := range []PropertyName{TopLeft, TopRight, BottomLeft, BottomRight} {
		tagX := tag + "-x"
		tagY := tag + "-y"
		valueX := radius.getRaw(tagX)
		valueY := radius.getRaw(tagY)

		if value := radius.getRaw(tag); value != nil {
			if valueX != nil && valueY != nil {
				radius.setRaw(tag, nil)
				result = append(result, tag)
			} else if valueX != nil && valueY == nil {
				if equalValue(value, valueX) {
					radius.setRaw(tagX, nil)
					result = append(result, tagX)
				} else {
					radius.setRaw(tagY, value)
					result = append(result, tagY)
					radius.setRaw(tag, nil)
					result = append(result, tag)
				}
			} else if valueX == nil && valueY != nil {
				if equalValue(value, valueY) {
					radius.setRaw(tagY, nil)
					result = append(result, tagY)
				} else {
					radius.setRaw(tagX, value)
					result = append(result, tagX)
					radius.setRaw(tag, nil)
					result = append(result, tag)
				}
			}
		} else if valueX != nil && valueY != nil && equalValue(valueX, valueY) {
			radius.setRaw(tag, valueX)
			result = append(result, tag)
			radius.setRaw(tagX, nil)
			result = append(result, tagX)
			radius.setRaw(tagY, nil)
			result = append(result, tagY)
		}
	}
}

func radiusPropertySet(radius Properties, tag PropertyName, value any) []PropertyName {
	var result []PropertyName = nil

	deleteTags := func(tags []PropertyName) {
		for _, tag := range tags {
			if radius.getRaw(tag) != nil {
				radius.setRaw(tag, nil)
				result = append(result, tag)
			}
		}
	}

	switch tag {
	case X:
		if result = setSizeProperty(radius, tag, value); result != nil {
			deleteTags([]PropertyName{TopLeftX, TopRightX, BottomLeftX, BottomRightX})
			for _, t := range []PropertyName{TopLeft, TopRight, BottomLeft, BottomRight} {
				if val := radius.getRaw(t); val != nil {
					t2 := t + "-y"
					if radius.getRaw(t2) != nil {
						radius.setRaw(t2, val)
						result = append(result, t2)
					}
					radius.setRaw(t, nil)
					result = append(result, t)
				}
			}
		}

	case Y:
		if result = setSizeProperty(radius, tag, value); result != nil {
			deleteTags([]PropertyName{TopLeftY, TopRightY, BottomLeftY, BottomRightY})
			for _, t := range []PropertyName{TopLeft, TopRight, BottomLeft, BottomRight} {
				if val := radius.getRaw(t); val != nil {
					t2 := t + "-x"
					if radius.getRaw(t2) != nil {
						radius.setRaw(t2, val)
						result = append(result, t2)
					}
					radius.setRaw(t, nil)
					result = append(result, t)
				}
			}
		}

	case TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY:
		if result = setSizeProperty(radius, tag, value); result != nil {
			deleteRadiusUnusedTags(radius, result)
		}

	case TopLeft, TopRight, BottomLeft, BottomRight:
		switch value := value.(type) {
		case SizeUnit:
			radius.setRaw(tag, value)
			result = []PropertyName{tag}
			deleteTags([]PropertyName{tag + "-x", tag + "-y"})
			deleteRadiusUnusedTags(radius, result)

		case string:
			if strings.Contains(value, "/") {
				if values := strings.Split(value, "/"); len(values) == 2 {
					if result = radiusPropertySet(radius, tag+"-x", values[0]); result != nil {
						if resultY := radiusPropertySet(radius, tag+"-y", values[1]); resultY != nil {
							result = append(result, resultY...)
						}

					}
				} else {
					notCompatibleType(tag, value)
				}
			} else {
				if result = setSizeProperty(radius, tag, value); result != nil {
					deleteTags([]PropertyName{tag + "-x", tag + "-y"})
					deleteRadiusUnusedTags(radius, result)
				}
			}
		}

	default:
		ErrorLogF(`"%s" property is not compatible with the RadiusProperty`, tag)
	}

	return result
}

func radiusPropertyGet(properties Properties, tag PropertyName) any {
	if value := properties.getRaw(tag); value != nil {
		return value
	}

	switch tag {
	case TopLeftX, TopLeftY, TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY:
		tagLen := len(tag)
		if value := properties.getRaw(tag[:tagLen-2]); value != nil {
			return value
		}
		if value := properties.getRaw(tag[tagLen-1:]); value != nil {
			return value
		}
	}

	switch tag {
	case TopLeftX, TopRightX, BottomLeftX, BottomRightX:
		if value := properties.getRaw(X); value != nil {
			return value
		}
	case TopLeftY, TopRightY, BottomLeftY, BottomRightY:
		if value := properties.getRaw(Y); value != nil {
			return value
		}
	}

	return nil
}

func (radius *radiusPropertyData) BoxRadius(session Session) BoxRadius {
	x, _ := sizeProperty(radius, X, session)
	y, _ := sizeProperty(radius, Y, session)

	getRadius := func(tag PropertyName) (SizeUnit, SizeUnit) {
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

func (radius BoxRadius) cssValue(builder cssBuilder, session Session) {

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

	buffer.WriteString(radius.TopLeftX.cssString("0", session))

	if radius.AllAnglesIsEqual() {

		if !radius.TopLeftX.Equal(radius.TopLeftY) {
			buffer.WriteString(" / ")
			buffer.WriteString(radius.TopLeftY.cssString("0", session))
		}

	} else {

		buffer.WriteRune(' ')
		buffer.WriteString(radius.TopRightX.cssString("0", session))
		buffer.WriteRune(' ')
		buffer.WriteString(radius.BottomRightX.cssString("0", session))
		buffer.WriteRune(' ')
		buffer.WriteString(radius.BottomLeftX.cssString("0", session))

		if !radius.TopLeftX.Equal(radius.TopLeftY) ||
			!radius.TopRightX.Equal(radius.TopRightY) ||
			!radius.BottomLeftX.Equal(radius.BottomLeftY) ||
			!radius.BottomRightX.Equal(radius.BottomRightY) {

			buffer.WriteString(" / ")
			buffer.WriteString(radius.TopLeftY.cssString("0", session))
			buffer.WriteRune(' ')
			buffer.WriteString(radius.TopRightY.cssString("0", session))
			buffer.WriteRune(' ')
			buffer.WriteString(radius.BottomRightY.cssString("0", session))
			buffer.WriteRune(' ')
			buffer.WriteString(radius.BottomLeftY.cssString("0", session))
		}
	}

	builder.add("border-radius", buffer.String())
}

func (radius BoxRadius) cssString(session Session) string {
	var builder cssValueBuilder
	radius.cssValue(&builder, session)
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

func setRadiusProperty(properties Properties, value any) []PropertyName {

	if value == nil {
		return propertiesRemove(properties, Radius)
	}

	switch value := value.(type) {
	case RadiusProperty:
		properties.setRaw(Radius, value)

	case SizeUnit:
		properties.setRaw(Radius, value)

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
		properties.setRaw(Radius, radius)

	case string:
		if strings.Contains(value, "/") {
			values := strings.Split(value, "/")
			if len(values) == 2 {
				if setRadiusPropertyElement(properties, RadiusX, values[0]) {
					result := []PropertyName{Radius, RadiusX}
					if setRadiusPropertyElement(properties, RadiusY, values[1]) {
						result = append(result, RadiusY)
					}
					return result
				}
			}
			notCompatibleType(Radius, value)
			return nil

		} else {
			return setSizeProperty(properties, Radius, value)
		}

	case DataObject:
		radius := NewRadiusProperty(nil)
		for _, tag := range []PropertyName{X, Y, TopLeft, TopRight, BottomLeft, BottomRight, TopLeftX, TopLeftY,
			TopRightX, TopRightY, BottomLeftX, BottomLeftY, BottomRightX, BottomRightY} {
			if value, ok := value.PropertyValue(string(tag)); ok {
				radius.Set(tag, value)
			}
		}
		properties.setRaw(Radius, radius)

	case float32:
		properties.setRaw(Radius, Px(float64(value)))

	case float64:
		properties.setRaw(Radius, Px(value))

	default:
		if n, ok := isInt(value); ok {
			properties.setRaw(Radius, Px(float64(n)))
		} else {
			notCompatibleType(Radius, value)
			return nil
		}
	}

	return []PropertyName{Radius}
}

func removeRadiusPropertyElement(properties Properties, tag PropertyName) bool {
	if value := properties.getRaw(Radius); value != nil {
		radius := getRadiusProperty(properties)
		radius.Remove(tag)
		if radius.empty() {
			properties.setRaw(Radius, nil)
		} else {
			properties.setRaw(Radius, radius)
		}
		return true
	}
	return false
}

func setRadiusPropertyElement(properties Properties, tag PropertyName, value any) bool {
	if value == nil {
		removeRadiusPropertyElement(properties, tag)
		return true
	}

	radius := getRadiusProperty(properties)
	if radius.Set(tag, value) {
		properties.setRaw(Radius, radius)
		return true
	}

	return false
}

func getRadiusElement(style Properties, tag PropertyName) any {
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
