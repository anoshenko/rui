package rui

import "strings"

// Constants related to view's background gradient description
const (
	// EllipseGradient is value of the Shape property of a radial gradient background:
	// the shape is an axis-aligned ellipse
	EllipseGradient = 0

	// CircleGradient is value of the Shape property of a radial gradient background:
	// the gradient's shape is a circle with constant radius
	CircleGradient = 1

	// ClosestSideGradient is value of the Radius property of a radial gradient background:
	// The gradient's ending shape meets the side of the box closest to its center (for circles)
	// or meets both the vertical and horizontal sides closest to the center (for ellipses).
	ClosestSideGradient = 0

	// ClosestCornerGradient is value of the Radius property of a radial gradient background:
	// The gradient's ending shape is sized so that it exactly meets the closest corner
	// of the box from its center.
	ClosestCornerGradient = 1

	// FarthestSideGradient is value of the Radius property of a radial gradient background:
	// Similar to closest-side, except the ending shape is sized to meet the side of the box
	// farthest from its center (or vertical and horizontal sides).
	FarthestSideGradient = 2

	// FarthestCornerGradient is value of the Radius property of a radial gradient background:
	// The default value, the gradient's ending shape is sized so that it exactly meets
	// the farthest corner of the box from its center.
	FarthestCornerGradient = 3
)

type backgroundRadialGradient struct {
	backgroundGradient
}

// NewBackgroundRadialGradient creates the new background radial gradient
func NewBackgroundRadialGradient(params Params) BackgroundElement {
	result := new(backgroundRadialGradient)
	result.init()
	for tag, value := range params {
		result.Set(tag, value)
	}
	return result
}

func (gradient *backgroundRadialGradient) init() {
	gradient.backgroundElement.init()
	gradient.normalize = normalizeRadialGradientTag
	gradient.set = backgroundRadialGradientSet
	gradient.supportedProperties = []PropertyName{
		RadialGradientRadius, RadialGradientShape, CenterX, CenterY, Gradient, Repeating,
	}
}

func (gradient *backgroundRadialGradient) Tag() string {
	return "radial-gradient"
}

func (image *backgroundRadialGradient) Clone() BackgroundElement {
	result := NewBackgroundRadialGradient(nil)
	for tag, value := range image.properties {
		result.setRaw(tag, value)
	}
	return result
}

func normalizeRadialGradientTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case Radius:
		tag = RadialGradientRadius

	case Shape:
		tag = RadialGradientShape

	case "x-center":
		tag = CenterX

	case "y-center":
		tag = CenterY
	}

	return tag
}

func backgroundRadialGradientSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case RadialGradientRadius:
		switch value := value.(type) {
		case []SizeUnit:
			switch len(value) {
			case 0:
				properties.setRaw(RadialGradientRadius, nil)

			case 1:
				if value[0].Type == Auto {
					properties.setRaw(RadialGradientRadius, nil)
				} else {
					properties.setRaw(RadialGradientRadius, value[0])
				}

			default:
				properties.setRaw(RadialGradientRadius, value)
			}
			return []PropertyName{tag}

		case []any:
			switch len(value) {
			case 0:
				properties.setRaw(RadialGradientRadius, nil)
				return []PropertyName{tag}

			case 1:
				return backgroundRadialGradientSet(properties, RadialGradientRadius, value[0])

			default:
				properties.setRaw(RadialGradientRadius, value)
				return []PropertyName{tag}
			}

		case string:
			if setSimpleProperty(properties, RadialGradientRadius, value) {
				return []PropertyName{tag}
			}
			if size, err := stringToSizeUnit(value); err == nil {
				if size.Type == Auto {
					properties.setRaw(RadialGradientRadius, nil)
				} else {
					properties.setRaw(RadialGradientRadius, size)
				}
				return []PropertyName{tag}
			}
			return setEnumProperty(properties, RadialGradientRadius, value, enumProperties[RadialGradientRadius].values)

		case SizeUnit:
			if value.Type == Auto {
				properties.setRaw(RadialGradientRadius, nil)
			} else {
				properties.setRaw(RadialGradientRadius, value)
			}
			return []PropertyName{tag}

		case int:
			return setEnumProperty(properties, RadialGradientRadius, value, enumProperties[RadialGradientRadius].values)
		}

		ErrorLogF(`Invalid value of "%s" property: %v`, tag, value)
		return nil

	case RadialGradientShape, CenterX, CenterY:
		return propertiesSet(properties, tag, value)
	}

	return backgroundGradientSet(properties, tag, value)
}

func (gradient *backgroundRadialGradient) cssStyle(session Session) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if repeating, _ := boolProperty(gradient, Repeating, session); repeating {
		buffer.WriteString(`repeating-radial-gradient(`)
	} else {
		buffer.WriteString(`radial-gradient(`)
	}

	var shapeText string
	if shape, ok := enumProperty(gradient, RadialGradientShape, session, EllipseGradient); ok && shape == CircleGradient {
		shapeText = `circle `
	} else {
		shapeText = `ellipse `
	}

	if value, ok := gradient.properties[RadialGradientRadius]; ok {
		switch value := value.(type) {
		case string:
			if text, ok := session.resolveConstants(value); ok {
				values := enumProperties[RadialGradientRadius]
				if n, ok := enumStringToInt(text, values.values, false); ok {
					buffer.WriteString(shapeText)
					shapeText = ""
					buffer.WriteString(values.cssValues[n])
					buffer.WriteString(" ")
				} else {
					if r, ok := StringToSizeUnit(text); ok && r.Type != Auto {
						buffer.WriteString("ellipse ")
						shapeText = ""
						buffer.WriteString(r.cssString("", session))
						buffer.WriteString(" ")
						buffer.WriteString(r.cssString("", session))
						buffer.WriteString(" ")
					} else {
						ErrorLog(`Invalid radial gradient radius: ` + text)
					}
				}
			} else {
				ErrorLog(`Invalid radial gradient radius: ` + value)
			}

		case int:
			values := enumProperties[RadialGradientRadius].cssValues
			if value >= 0 && value < len(values) {
				buffer.WriteString(shapeText)
				shapeText = ""
				buffer.WriteString(values[value])
				buffer.WriteString(" ")
			} else {
				ErrorLogF(`Invalid radial gradient radius: %d`, value)
			}

		case SizeUnit:
			if value.Type != Auto {
				buffer.WriteString("ellipse ")
				shapeText = ""
				buffer.WriteString(value.cssString("", session))
				buffer.WriteString(" ")
				buffer.WriteString(value.cssString("", session))
				buffer.WriteString(" ")
			}

		case []SizeUnit:
			count := len(value)
			if count > 2 {
				count = 2
			}
			buffer.WriteString("ellipse ")
			shapeText = ""
			for i := 0; i < count; i++ {
				buffer.WriteString(value[i].cssString("50%", session))
				buffer.WriteString(" ")
			}

		case []any:
			count := len(value)
			if count > 2 {
				count = 2
			}
			buffer.WriteString("ellipse ")
			shapeText = ""
			for i := 0; i < count; i++ {
				if value[i] != nil {
					switch value := value[i].(type) {
					case SizeUnit:
						buffer.WriteString(value.cssString("50%", session))
						buffer.WriteString(" ")

					case string:
						if text, ok := session.resolveConstants(value); ok {
							if size, err := stringToSizeUnit(text); err == nil {
								buffer.WriteString(size.cssString("50%", session))
								buffer.WriteString(" ")
							} else {
								buffer.WriteString("50% ")
							}
						} else {
							buffer.WriteString("50% ")
						}
					}
				} else {
					buffer.WriteString("50% ")
				}
			}
		}
	}

	x, _ := sizeProperty(gradient, CenterX, session)
	y, _ := sizeProperty(gradient, CenterX, session)
	if x.Type != Auto || y.Type != Auto {
		if shapeText != "" {
			buffer.WriteString(shapeText)
		}
		buffer.WriteString("at ")
		buffer.WriteString(x.cssString("50%", session))
		buffer.WriteString(" ")
		buffer.WriteString(y.cssString("50%", session))
	} else if shapeText != "" {
		buffer.WriteString(shapeText)
	}

	buffer.WriteString(", ")
	if !gradient.writeGradient(session, buffer) {
		return ""
	}

	buffer.WriteString(") ")

	return buffer.String()
}

func (gradient *backgroundRadialGradient) writeString(buffer *strings.Builder, indent string) {
	gradient.writeToBuffer(buffer, indent, gradient.Tag(), []PropertyName{
		Gradient,
		CenterX,
		CenterY,
		Repeating,
		RadialGradientShape,
		RadialGradientRadius,
	})
}

func (gradient *backgroundRadialGradient) String() string {
	return runStringWriter(gradient)
}
