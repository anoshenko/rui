package rui

import "strings"

// Constants related to view's background gradient description
const (

	// ToTopGradient is value of the Direction property of a linear gradient. The value is equivalent to the 0deg angle
	ToTopGradient = 0
	// ToRightTopGradient is value of the Direction property of a linear gradient.
	ToRightTopGradient = 1
	// ToRightGradient is value of the Direction property of a linear gradient. The value is equivalent to the 90deg angle
	ToRightGradient = 2
	// ToRightBottomGradient is value of the Direction property of a linear gradient.
	ToRightBottomGradient = 3
	// ToBottomGradient is value of the Direction property of a linear gradient. The value is equivalent to the 180deg angle
	ToBottomGradient = 4
	// ToLeftBottomGradient is value of the Direction property of a linear gradient.
	ToLeftBottomGradient = 5
	// ToLeftGradient is value of the Direction property of a linear gradient. The value is equivalent to the 270deg angle
	ToLeftGradient = 6
	// ToLeftTopGradient is value of the Direction property of a linear gradient.
	ToLeftTopGradient = 7

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

// BackgroundGradientPoint define point on gradient straight line
type BackgroundGradientPoint struct {
	// Color - the color of the point. Must not be nil.
	// Can take a value of Color type or string (color constant or textual description of the color)
	Color any
	// Pos - the distance from the start of the gradient straight line. Optional (may be nil).
	// Can take a value of SizeUnit type or string (angle constant or textual description of the SizeUnit)
	Pos any
}

type backgroundGradient struct {
	backgroundElement
}

type backgroundLinearGradient struct {
	backgroundGradient
}

type backgroundRadialGradient struct {
	backgroundGradient
}

// NewBackgroundLinearGradient creates the new background linear gradient
func NewBackgroundLinearGradient(params Params) BackgroundElement {
	result := new(backgroundLinearGradient)
	result.init()
	for tag, value := range params {
		result.Set(tag, value)
	}
	return result
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

func parseGradientText(value string) []BackgroundGradientPoint {
	elements := strings.Split(value, ",")
	count := len(elements)
	if count < 2 {
		ErrorLog("The gradient must contain at least 2 points")
		return nil
	}

	points := make([]BackgroundGradientPoint, count)
	for i, element := range elements {
		if !points[i].setValue(element) {
			ErrorLogF(`Invalid %d element of the conic gradient: "%s"`, i, element)
			return nil
		}
	}
	return points
}

func backgroundGradientSet(properties Properties, tag PropertyName, value any) []PropertyName {

	switch tag {
	case Repeating:
		return setBoolProperty(properties, tag, value)

	case Gradient:
		switch value := value.(type) {
		case string:
			if value != "" {
				if strings.Contains(value, " ") || strings.Contains(value, ",") {
					if points := parseGradientText(value); len(points) >= 2 {
						properties.setRaw(Gradient, points)
						return []PropertyName{tag}
					}
				} else if value[0] == '@' {
					properties.setRaw(Gradient, value)
					return []PropertyName{tag}
				}
			}

		case []BackgroundGradientPoint:
			if len(value) >= 2 {
				properties.setRaw(Gradient, value)
				return []PropertyName{tag}
			}

		case []Color:
			count := len(value)
			if count >= 2 {
				points := make([]BackgroundGradientPoint, count)
				for i, color := range value {
					points[i].Color = color
				}
				properties.setRaw(Gradient, points)
				return []PropertyName{tag}
			}

		case []GradientPoint:
			count := len(value)
			if count >= 2 {
				points := make([]BackgroundGradientPoint, count)
				for i, point := range value {
					points[i].Color = point.Color
					points[i].Pos = Percent(point.Offset * 100)
				}
				properties.setRaw(Gradient, points)
				return []PropertyName{tag}
			}
		}

		ErrorLogF("Invalid gradient %v", value)
		return nil
	}

	ErrorLogF("Property %s is not supported by a background gradient", tag)
	return nil
}

func (point *BackgroundGradientPoint) setValue(text string) bool {
	text = strings.Trim(text, " ")

	colorText := text
	pointText := ""

	if index := strings.Index(text, " "); index > 0 {
		colorText = text[:index]
		pointText = strings.Trim(text[index+1:], " ")
	}

	if colorText == "" {
		return false
	}

	if colorText[0] == '@' {
		point.Color = colorText
	} else if color, ok := StringToColor(colorText); ok {
		point.Color = color
	} else {
		return false
	}

	if pointText == "" {
		point.Pos = nil
	} else if pointText[0] == '@' {
		point.Pos = pointText
	} else if pos, ok := StringToSizeUnit(pointText); ok {
		point.Pos = pos
	} else {
		return false
	}

	return true
}

func (point *BackgroundGradientPoint) color(session Session) (Color, bool) {
	if point.Color != nil {
		switch color := point.Color.(type) {
		case string:
			if color != "" {
				if color[0] == '@' {
					if clr, ok := session.Color(color[1:]); ok {
						return clr, true
					}
				} else {
					if clr, ok := StringToColor(color); ok {
						return clr, true
					}
				}
			}

		case Color:
			return color, true

		default:
			if n, ok := isInt(point.Color); ok {
				return Color(n), true
			}
		}
	}
	return 0, false
}

// String convert internal representation of [BackgroundGradientPoint] into a string.
func (point *BackgroundGradientPoint) String() string {
	result := "black"
	if point.Color != nil {
		switch color := point.Color.(type) {
		case string:
			result = color

		case Color:
			result = color.String()
		}
	}

	if point.Pos != nil {
		switch value := point.Pos.(type) {
		case string:
			result += " " + value

		case SizeUnit:
			if value.Type != Auto {
				result += " " + value.String()
			}
		}
	}

	return result
}

func (gradient *backgroundGradient) writeGradient(session Session, buffer *strings.Builder) bool {

	value, ok := gradient.properties[Gradient]
	if !ok {
		return false
	}

	var points []BackgroundGradientPoint = nil

	switch value := value.(type) {
	case string:
		if value != "" && value[0] == '@' {
			if text, ok := session.Constant(value[1:]); ok {
				points = parseGradientText(text)
			}
		}

	case []BackgroundGradientPoint:
		points = value
	}

	if len(points) > 0 {
		for i, point := range points {
			if i > 0 {
				buffer.WriteString(`, `)
			}

			if color, ok := point.color(session); ok {
				buffer.WriteString(color.cssString())
			} else {
				return false
			}

			if point.Pos != nil {
				switch value := point.Pos.(type) {
				case string:
					if value != "" {
						if value, ok := session.resolveConstants(value); ok {
							if pos, ok := StringToSizeUnit(value); ok && pos.Type != Auto {
								buffer.WriteRune(' ')
								buffer.WriteString(pos.cssString("", session))
							}
						}
					}

				case SizeUnit:
					if value.Type != Auto {
						buffer.WriteRune(' ')
						buffer.WriteString(value.cssString("", session))
					}
				}
			}
		}
		return true
	}

	return false
}

func (gradient *backgroundLinearGradient) init() {
	gradient.backgroundElement.init()
	gradient.set = backgroundLinearGradientSet
	gradient.supportedProperties = []PropertyName{Direction, Repeating, Gradient}

}

func (gradient *backgroundLinearGradient) Tag() string {
	return "linear-gradient"
}

func (image *backgroundLinearGradient) Clone() BackgroundElement {
	result := NewBackgroundLinearGradient(nil)
	for tag, value := range image.properties {
		result.setRaw(tag, value)
	}
	return result
}

func backgroundLinearGradientSet(properties Properties, tag PropertyName, value any) []PropertyName {
	if tag == Direction {
		switch value := value.(type) {
		case AngleUnit:
			properties.setRaw(Direction, value)
			return []PropertyName{tag}

		case string:
			if setSimpleProperty(properties, tag, value) {
				return []PropertyName{tag}
			}
			if angle, ok := StringToAngleUnit(value); ok {
				properties.setRaw(Direction, angle)
				return []PropertyName{tag}
			}
		}
		return setEnumProperty(properties, tag, value, enumProperties[Direction].values)
	}

	return backgroundGradientSet(properties, tag, value)
}

func (gradient *backgroundLinearGradient) cssStyle(session Session) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if repeating, _ := boolProperty(gradient, Repeating, session); repeating {
		buffer.WriteString(`repeating-linear-gradient(`)
	} else {
		buffer.WriteString(`linear-gradient(`)
	}

	if value, ok := gradient.properties[Direction]; ok {
		switch value := value.(type) {
		case string:
			if text, ok := session.resolveConstants(value); ok {
				direction := enumProperties[Direction]
				if n, ok := enumStringToInt(text, direction.values, false); ok {
					buffer.WriteString(direction.cssValues[n])
					buffer.WriteString(", ")
				} else {
					if angle, ok := StringToAngleUnit(text); ok {
						buffer.WriteString(angle.cssString())
						buffer.WriteString(", ")
					} else {
						ErrorLog(`Invalid linear gradient direction: ` + text)
					}
				}
			} else {
				ErrorLog(`Invalid linear gradient direction: ` + value)
			}

		case int:
			values := enumProperties[Direction].cssValues
			if value >= 0 && value < len(values) {
				buffer.WriteString(values[value])
				buffer.WriteString(", ")
			} else {
				ErrorLogF(`Invalid linear gradient direction: %d`, value)
			}

		case AngleUnit:
			buffer.WriteString(value.cssString())
			buffer.WriteString(", ")
		}
	}

	if !gradient.writeGradient(session, buffer) {
		return ""
	}

	buffer.WriteString(") ")
	return buffer.String()
}

func (gradient *backgroundLinearGradient) writeString(buffer *strings.Builder, indent string) {
	gradient.writeToBuffer(buffer, indent, gradient.Tag(), []PropertyName{
		Gradient,
		Repeating,
		Direction,
	})
}

func (gradient *backgroundLinearGradient) String() string {
	return runStringWriter(gradient)
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
