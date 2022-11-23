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
func NewBackgroundConicGradient(params Params) BackgroundElement {
	result := new(backgroundConicGradient)
	result.properties = map[string]any{}
	for tag, value := range params {
		result.Set(tag, value)
	}
	return result
}

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
			if value != "" {
				if value[0] == '@' {
					if val, ok := session.Constant(value[1:]); ok {
						value = val
					} else {
						return
					}
				}

				if angle, ok := StringToAngleUnit(value); ok {
					buffer.WriteRune(' ')
					buffer.WriteString(angle.cssString())
				}
			}

		case AngleUnit:
			buffer.WriteRune(' ')
			buffer.WriteString(value.cssString())
		}
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

func (gradient *backgroundConicGradient) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case "x-center":
		tag = CenterX

	case "y-center":
		tag = CenterY
	}

	return tag
}

func (gradient *backgroundConicGradient) Set(tag string, value any) bool {
	tag = gradient.normalizeTag(tag)
	switch tag {
	case CenterX, CenterY, Repeating, From:
		return gradient.propertyList.Set(tag, value)

	case Gradient:
		return gradient.setGradient(value)
	}

	ErrorLogF(`"%s" property is not supported by BackgroundConicGradient`, tag)
	return false
}

func (gradient *backgroundConicGradient) stringToAngle(text string) (any, bool) {
	if text == "" {
		return nil, false
	} else if text[0] == '@' {
		return text, true
	}
	return StringToAngleUnit(text)
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

	if colorText[0] == '@' {
		result.Color = colorText
	} else if color, ok := StringToColor(colorText); ok {
		result.Color = color
	} else {
		return result, false
	}

	if pointText != "" {
		if angle, ok := gradient.stringToAngle(pointText); ok {
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

func (gradient *backgroundConicGradient) setGradient(value any) bool {
	if value == nil {
		delete(gradient.properties, Gradient)
		return true
	}

	switch value := value.(type) {
	case string:
		if value == "" {
			delete(gradient.properties, Gradient)
			return true
		}

		if strings.Contains(value, ",") || strings.Contains(value, " ") {
			if vector := gradient.parseGradientText(value); vector != nil {
				gradient.properties[Gradient] = vector
				return true
			}
			return false
		} else if value[0] == '@' {
			gradient.properties[Gradient] = value
			return true
		}

		ErrorLogF(`Invalid conic gradient: "%s"`, value)
		return false

	case []BackgroundGradientAngle:
		count := len(value)
		if count < 2 {
			ErrorLog("The gradient must contain at least 2 points")
			return false
		}

		for i, point := range value {
			if point.Color == nil {
				ErrorLogF("Invalid %d element of the conic gradient: Color is nil", i)
				return false
			}
		}
		gradient.properties[Gradient] = value
		return true
	}
	return false
}

func (gradient *backgroundConicGradient) Get(tag string) any {
	return gradient.backgroundElement.Get(gradient.normalizeTag(tag))
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
