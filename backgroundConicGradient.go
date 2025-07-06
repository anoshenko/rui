package rui

import (
	"strings"
)

type backgroundConicGradient struct {
	backgroundElement
}

// BackgroundGradientAngle defined an element of the conic gradient
type BackgroundGradientAngle struct {
	// Color - the color of the key angle. Must not be nil.
	// Can take a value of Color type or string (color constant or textual description of the color)
	Color any
	// Angle - the key angle. Optional (may be nil).
	// Can take a value of AngleUnit type or string (angle constant or textual description of the angle)
	Angle any
}

// NewBackgroundConicGradient creates the new background conic gradient
//
// The following properties can be used:
//   - "gradient" [Gradient] - Describes gradient stop points. This is a mandatory property while describing background gradients.
//   - "center-x" [CenterX] - center X point of the gradient.
//   - "center-y" [CenterY] - center Y point of the gradient.
//   - "from" [From] - start angle position of the gradient.
//   - "repeating" [Repeating] - Defines whether stop points needs to be repeated after the last one.
func NewBackgroundConicGradient(params Params) BackgroundElement {
	result := new(backgroundConicGradient)
	result.init()
	for tag, value := range params {
		result.Set(tag, value)
	}
	return result
}

// String convert internal representation of [BackgroundGradientAngle] into a string.
func (point *BackgroundGradientAngle) String() string {
	result := "black"
	if point.Color != nil {
		switch color := point.Color.(type) {
		case string:
			result = color

		case Color:
			result = color.String()
		}
	}

	if point.Angle != nil {
		switch value := point.Angle.(type) {
		case string:
			result += " " + value

		case AngleUnit:
			result += " " + value.String()
		}
	}

	return result
}

func (point *BackgroundGradientAngle) color(session Session) (Color, bool) {
	if point.Color != nil {
		switch color := point.Color.(type) {
		case string:
			if ok, constName := isConstantName(color); ok {
				if clr, ok := session.Color(constName); ok {
					return clr, true
				}
			} else if clr, ok := StringToColor(color); ok {
				return clr, true
			}

		case Color:
			return color, true

		default:
			if n, ok := isInt(color); ok {
				return Color(n), true
			}
		}
	}
	return 0, false
}

func (point *BackgroundGradientAngle) isValid(session Session) bool {
	_, ok := point.color(session)
	return ok
}

func (point *BackgroundGradientAngle) cssString(session Session, buffer *strings.Builder) {
	if color, ok := point.color(session); ok {
		buffer.WriteString(color.cssString())
	} else {
		return
	}

	if point.Angle != nil {
		switch value := point.Angle.(type) {
		case string:
			if ok, constName := isConstantName(value); ok {
				if value, ok = session.Constant(constName); !ok {
					return
				}
			}

			if angle, ok := StringToAngleUnit(value); ok {
				buffer.WriteRune(' ')
				buffer.WriteString(angle.cssString())
			}

		case AngleUnit:
			buffer.WriteRune(' ')
			buffer.WriteString(value.cssString())
		}
	}
}

func (gradient *backgroundConicGradient) init() {
	gradient.backgroundElement.init()
	gradient.normalize = normalizeConicGradientTag
	gradient.set = backgroundConicGradientSet
	gradient.supportedProperties = []PropertyName{
		CenterX, CenterY, Repeating, From, Gradient,
	}
}

func (gradient *backgroundConicGradient) Tag() string {
	return "conic-gradient"
}

func (image *backgroundConicGradient) Clone() BackgroundElement {
	result := NewBackgroundConicGradient(nil)
	for tag, value := range image.properties {
		result.setRaw(tag, value)
	}
	return result
}

func normalizeConicGradientTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case "x-center":
		tag = CenterX

	case "y-center":
		tag = CenterY
	}

	return tag
}

func backgroundConicGradientSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case Gradient:
		switch value := value.(type) {
		case string:
			if value == "" {
				return propertiesRemove(properties, tag)
			}

			if strings.Contains(value, ",") || strings.Contains(value, " ") {
				if vector := parseGradientText(value); vector != nil {
					properties.setRaw(Gradient, vector)
					return []PropertyName{tag}
				}
			} else if ok, _ := isConstantName(value); ok {
				properties.setRaw(Gradient, value)
				return []PropertyName{tag}
			}

			ErrorLogF(`Invalid conic gradient: "%s"`, value)

		case []BackgroundGradientAngle:
			count := len(value)
			if count < 2 {
				ErrorLog("The gradient must contain at least 2 points")
				return nil
			}

			for i, point := range value {
				if point.Color == nil {
					ErrorLogF("Invalid %d element of the conic gradient: Color is nil", i)
					return nil
				}
			}
			properties.setRaw(Gradient, value)
			return []PropertyName{tag}

		default:
			notCompatibleType(tag, value)
		}
		return nil
	}

	return propertiesSet(properties, tag, value)
}

func (gradient *backgroundConicGradient) stringToGradientPoint(text string) (BackgroundGradientAngle, bool) {
	var result BackgroundGradientAngle
	colorText := ""
	pointText := ""

	if index := strings.Index(text, " "); index > 0 {
		colorText = text[:index]
		pointText = strings.Trim(text[index+1:], " ")
	} else {
		colorText = text
	}

	if colorText == "" {
		return result, false
	}

	if ok, _ := isConstantName(colorText); ok {
		result.Color = colorText
	} else if color, ok := StringToColor(colorText); ok {
		result.Color = color
	} else {
		return result, false
	}

	if pointText != "" {
		if ok, _ := isConstantName(pointText); ok {
			result.Angle = pointText
		} else if angle, ok := StringToAngleUnit(text); ok {
			result.Angle = angle
		} else {
			return result, false
		}
	}

	return result, true
}

func (gradient *backgroundConicGradient) parseGradientText(value string) []BackgroundGradientAngle {
	elements := strings.Split(value, ",")
	count := len(elements)
	if count < 2 {
		ErrorLog("The gradient must contain at least 2 points")
		return nil
	}

	vector := make([]BackgroundGradientAngle, count)
	for i, element := range elements {
		var ok bool
		if vector[i], ok = gradient.stringToGradientPoint(strings.Trim(element, " ")); !ok {
			ErrorLogF(`Invalid %d element of the conic gradient: "%s"`, i, element)
			return nil
		}
	}
	return vector
}
func (gradient *backgroundConicGradient) cssStyle(session Session) string {

	points := []BackgroundGradientAngle{}
	if value, ok := gradient.properties[Gradient]; ok {
		switch value := value.(type) {
		case string:
			if text, ok := session.resolveConstants(value); ok && text != "" {
				if points = gradient.parseGradientText(text); points == nil {
					return ""
				}
			} else {
				ErrorLog(`Invalid conic gradient: ` + value)
				return ""
			}

		case []BackgroundGradientAngle:
			points = value
		}
	} else {
		return ""
	}

	if len(points) < 2 {
		ErrorLog("The gradient must contain at least 2 points")
		return ""
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if repeating, _ := boolProperty(gradient, Repeating, session); repeating {
		buffer.WriteString(`repeating-conic-gradient(`)
	} else {
		buffer.WriteString(`conic-gradient(`)
	}

	comma := false
	if angle, ok := angleProperty(gradient, From, session); ok {
		buffer.WriteString("from ")
		buffer.WriteString(angle.cssString())
		comma = true
	}

	x, _ := sizeProperty(gradient, CenterX, session)
	y, _ := sizeProperty(gradient, CenterX, session)
	if x.Type != Auto || y.Type != Auto {
		if comma {
			buffer.WriteRune(' ')
		}
		buffer.WriteString("at ")
		buffer.WriteString(x.cssString("50%", session))
		buffer.WriteString(" ")
		buffer.WriteString(y.cssString("50%", session))
		comma = true
	}

	for _, point := range points {
		if point.isValid(session) {
			if comma {
				buffer.WriteString(`, `)
			}
			point.cssString(session, buffer)
			comma = true
		}
	}

	buffer.WriteString(") ")

	return buffer.String()
}

func (gradient *backgroundConicGradient) writeString(buffer *strings.Builder, indent string) {
	gradient.writeToBuffer(buffer, indent, gradient.Tag(), []PropertyName{
		Gradient,
		CenterX,
		CenterY,
		Repeating,
	})
}

func (gradient *backgroundConicGradient) String() string {
	return runStringWriter(gradient)
}
