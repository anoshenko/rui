package rui

import "strings"

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
	// Pos - the distance from the start of the gradient straight line
	Pos SizeUnit
	// Color - the color of the point
	Color Color
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
	result.properties = map[string]interface{}{}
	for tag, value := range params {
		result.Set(tag, value)
	}
	return result
}

// NewBackgroundRadialGradient creates the new background radial gradient
func NewBackgroundRadialGradient(params Params) BackgroundElement {
	result := new(backgroundRadialGradient)
	result.properties = map[string]interface{}{}
	for tag, value := range params {
		result.Set(tag, value)
	}
	return result
}

func (gradient *backgroundGradient) Set(tag string, value interface{}) bool {

	switch tag = strings.ToLower(tag); tag {
	case Repeating:
		return gradient.setBoolProperty(tag, value)

	case Gradient:
		switch value := value.(type) {
		case string:
			if value != "" {
				elements := strings.Split(value, `,`)
				if count := len(elements); count > 1 {
					points := make([]interface{}, count)
					for i, element := range elements {
						if strings.Contains(element, "@") {
							points[i] = element
						} else {
							var point BackgroundGradientPoint
							if point.setValue(element) {
								points[i] = point
							} else {
								ErrorLogF("Invalid gradient element #%d: %s", i, element)
								return false
							}
						}
					}
					gradient.properties[Gradient] = points
					return true
				}

				text := strings.Trim(value, " \n\r\t")
				if text[0] == '@' {
					gradient.properties[Gradient] = text
					return true
				}
			}

		case []BackgroundGradientPoint:
			if len(value) >= 2 {
				gradient.properties[Gradient] = value
				return true
			}

		case []Color:
			count := len(value)
			if count >= 2 {
				points := make([]BackgroundGradientPoint, count)
				for i, color := range value {
					points[i].Color = color
					points[i].Pos = AutoSize()
				}
				gradient.properties[Gradient] = points
				return true
			}

		case []GradientPoint:
			count := len(value)
			if count >= 2 {
				points := make([]BackgroundGradientPoint, count)
				for i, point := range value {
					points[i].Color = point.Color
					points[i].Pos = Percent(point.Offset * 100)
				}
				gradient.properties[Gradient] = points
				return true
			}

		case []interface{}:
			if count := len(value); count > 1 {
				points := make([]interface{}, count)
				for i, element := range value {
					switch element := element.(type) {
					case string:
						if strings.Contains(element, "@") {
							points[i] = element
						} else {
							var point BackgroundGradientPoint
							if !point.setValue(element) {
								ErrorLogF("Invalid gradient element #%d: %s", i, element)
								return false
							}
							points[i] = point
						}

					case BackgroundGradientPoint:
						points[i] = element

					case GradientPoint:
						points[i] = BackgroundGradientPoint{Color: element.Color, Pos: Percent(element.Offset * 100)}

					case Color:
						points[i] = BackgroundGradientPoint{Color: element, Pos: AutoSize()}

					default:
						ErrorLogF("Invalid gradient element #%d: %v", i, element)
						return false
					}
				}
				gradient.properties[Gradient] = points
				return true
			}

		default:
			ErrorLogF("Invalid gradient %v", value)
			return false
		}
	}

	return gradient.backgroundElement.Set(tag, value)
}

func (point *BackgroundGradientPoint) setValue(value string) bool {
	var ok bool

	switch elements := strings.Split(value, `:`); len(elements) {
	case 2:
		if point.Color, ok = StringToColor(elements[0]); !ok {
			return false
		}
		if point.Pos, ok = StringToSizeUnit(elements[1]); !ok {
			return false
		}

	case 1:
		if point.Color, ok = StringToColor(elements[0]); !ok {
			return false
		}
		point.Pos = AutoSize()

	default:
		return false
	}

	return false
}

func (gradient *backgroundGradient) writeGradient(session Session, buffer *strings.Builder) bool {

	value, ok := gradient.properties[Gradient]
	if !ok {
		return false
	}

	points := []BackgroundGradientPoint{}

	switch value := value.(type) {
	case string:
		if text, ok := session.resolveConstants(value); ok && text != "" {
			elements := strings.Split(text, `,`)
			points := make([]BackgroundGradientPoint, len(elements))
			for i, element := range elements {
				if !points[i].setValue(element) {
					ErrorLogF(`Invalid gradient point #%d: "%s"`, i, element)
					return false
				}
			}
		} else {
			ErrorLog(`Invalid gradient: ` + value)
			return false
		}

	case []BackgroundGradientPoint:
		points = value

	case []interface{}:
		points = make([]BackgroundGradientPoint, len(value))
		for i, element := range value {
			switch element := element.(type) {
			case string:
				if text, ok := session.resolveConstants(element); ok && text != "" {
					if !points[i].setValue(text) {
						ErrorLogF(`Invalid gradient point #%d: "%s"`, i, text)
						return false
					}
				} else {
					ErrorLogF(`Invalid gradient point #%d: "%s"`, i, text)
					return false
				}

			case BackgroundGradientPoint:
				points[i] = element
			}
		}
	}

	if len(points) > 0 {
		for i, point := range points {
			if i > 0 {
				buffer.WriteString(`, `)
			}

			buffer.WriteString(point.Color.cssString())
			if point.Pos.Type != Auto {
				buffer.WriteRune(' ')
				buffer.WriteString(point.Pos.cssString(""))
			}
		}
		return true
	}

	return false
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

func (gradient *backgroundLinearGradient) Set(tag string, value interface{}) bool {
	if strings.ToLower(tag) == Direction {
		switch value := value.(type) {
		case AngleUnit:
			gradient.properties[Direction] = value
			return true

		case string:
			var angle AngleUnit
			if ok, _ := angle.setValue(value); ok {
				gradient.properties[Direction] = angle
				return true
			}
		}
		return gradient.setEnumProperty(tag, value, enumProperties[Direction].values)
	}

	return gradient.backgroundGradient.Set(tag, value)
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

func (gradient *backgroundRadialGradient) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
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

func (gradient *backgroundRadialGradient) Set(tag string, value interface{}) bool {
	tag = gradient.normalizeTag(tag)
	switch tag {
	case RadialGradientRadius:
		switch value := value.(type) {
		case string, SizeUnit:
			return gradient.propertyList.Set(RadialGradientRadius, value)

		case int:
			n := value
			if n >= 0 && n < len(enumProperties[RadialGradientRadius].values) {
				return gradient.propertyList.Set(RadialGradientRadius, value)
			}
		}
		ErrorLogF(`Invalid value of "%s" property: %v`, tag, value)

	case RadialGradientShape:
		return gradient.propertyList.Set(RadialGradientShape, value)

	case CenterX, CenterY:
		return gradient.propertyList.Set(tag, value)
	}

	return gradient.backgroundGradient.Set(tag, value)
}

func (gradient *backgroundRadialGradient) Get(tag string) interface{} {
	return gradient.backgroundGradient.Get(gradient.normalizeTag(tag))
}

func (gradient *backgroundRadialGradient) cssStyle(session Session) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if repeating, _ := boolProperty(gradient, Repeating, session); repeating {
		buffer.WriteString(`repeating-radial-gradient(`)
	} else {
		buffer.WriteString(`radial-gradient(`)
	}

	if shape, ok := enumProperty(gradient, RadialGradientShape, session, EllipseGradient); ok && shape == CircleGradient {
		buffer.WriteString(`circle `)
	} else {
		buffer.WriteString(`ellipse `)
	}

	if value, ok := gradient.properties[RadialGradientRadius]; ok {
		switch value := value.(type) {
		case string:
			if text, ok := session.resolveConstants(value); ok {
				values := enumProperties[RadialGradientRadius]
				if n, ok := enumStringToInt(text, values.values, false); ok {
					buffer.WriteString(values.cssValues[n])
					buffer.WriteString(" ")
				} else {
					if r, ok := StringToSizeUnit(text); ok && r.Type != Auto {
						buffer.WriteString(r.cssString(""))
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
				buffer.WriteString(values[value])
				buffer.WriteString(" ")
			} else {
				ErrorLogF(`Invalid radial gradient radius: %d`, value)
			}

		case SizeUnit:
			if value.Type != Auto {
				buffer.WriteString(value.cssString(""))
				buffer.WriteString(" ")
			}
		}
	}

	x, _ := sizeProperty(gradient, CenterX, session)
	y, _ := sizeProperty(gradient, CenterX, session)
	if x.Type != Auto || y.Type != Auto {
		buffer.WriteString("at ")
		buffer.WriteString(x.cssString("50%"))
		buffer.WriteString(" ")
		buffer.WriteString(y.cssString("50%"))
	}

	buffer.WriteString(", ")
	if !gradient.writeGradient(session, buffer) {
		return ""
	}

	buffer.WriteString(") ")

	return buffer.String()
}
