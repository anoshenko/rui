package rui

import (
	"fmt"
	"strconv"
	"strings"
)

func stringProperty(properties Properties, tag string, session Session) (string, bool) {
	if value := properties.getRaw(tag); value != nil {
		if text, ok := value.(string); ok {
			return session.resolveConstants(text)
		}
	}
	return "", false
}

func imageProperty(properties Properties, tag string, session Session) (string, bool) {
	if value := properties.getRaw(tag); value != nil {
		if text, ok := value.(string); ok {
			if text != "" && text[0] == '@' {
				if image, ok := session.ImageConstant(text[1:]); ok {
					return image, true
				} else {
					return "", false
				}
			}

			return text, true
		}
	}
	return "", false
}

func valueToSizeUnit(value any, session Session) (SizeUnit, bool) {
	if value != nil {
		switch value := value.(type) {
		case SizeUnit:
			return value, true

		case string:
			if text, ok := session.resolveConstants(value); ok {
				return StringToSizeUnit(text)
			}

		case float64:
			return Px(value), true

		case float32:
			return Px(float64(value)), true
		}

		if n, ok := isInt(value); ok {
			return Px(float64(n)), true
		}
	}

	return AutoSize(), false
}

func sizeProperty(properties Properties, tag string, session Session) (SizeUnit, bool) {
	return valueToSizeUnit(properties.getRaw(tag), session)
}

func angleProperty(properties Properties, tag string, session Session) (AngleUnit, bool) {
	if value := properties.getRaw(tag); value != nil {
		switch value := value.(type) {
		case AngleUnit:
			return value, true

		case string:
			if text, ok := session.resolveConstants(value); ok {
				return StringToAngleUnit(text)
			}
		}
	}

	return AngleUnit{Type: 0, Value: 0}, false
}

func valueToColor(value any, session Session) (Color, bool) {
	if value != nil {
		switch value := value.(type) {
		case Color:
			return value, true

		case string:
			if len(value) > 1 && value[0] == '@' {
				return session.Color(value[1:])
			}
			return StringToColor(value)
		}
	}

	return Color(0), false
}

func colorProperty(properties Properties, tag string, session Session) (Color, bool) {
	return valueToColor(properties.getRaw(tag), session)
}

func valueToEnum(value any, tag string, session Session, defaultValue int) (int, bool) {
	if value != nil {
		values := enumProperties[tag].values
		switch value := value.(type) {
		case int:
			if value >= 0 && value < len(values) {
				return value, true
			}

		case string:
			if text, ok := session.resolveConstants(value); ok {
				if tag == Orientation {
					switch strings.ToLower(text) {
					case "vertical":
						value = "up-down"

					case "horizontal":
						value = "left-to-right"
					}
				}
				if result, ok := enumStringToInt(text, values, true); ok {
					return result, true
				}
			}
		}
	}

	return defaultValue, false
}

func enumStringToInt(value string, enumValues []string, logError bool) (int, bool) {
	value = strings.Trim(value, " \t\n\r")

	for n, val := range enumValues {
		if val == value {
			return n, true
		}
	}

	if n, err := strconv.Atoi(value); err == nil {
		if n >= 0 && n < len(enumValues) {
			return n, true
		}

		if logError {
			ErrorLogF(`Out of bounds: value index = %d, valid values = [%v]`, n, enumValues)
		}
		return 0, false
	}

	value = strings.ToLower(value)
	for n, val := range enumValues {
		if val == value {
			return n, true
		}
	}

	if logError {
		ErrorLogF(`Unknown "%s" value. Valid values = [%v]`, value, enumValues)
	}
	return 0, false
}

func enumProperty(properties Properties, tag string, session Session, defaultValue int) (int, bool) {
	return valueToEnum(properties.getRaw(tag), tag, session, defaultValue)
}

func valueToBool(value any, session Session) (bool, bool) {
	if value != nil {
		switch value := value.(type) {
		case bool:
			return value, true

		case string:
			if text, ok := session.resolveConstants(value); ok {
				switch strings.ToLower(text) {
				case "true", "yes", "on", "1":
					return true, true

				case "false", "no", "off", "0":
					return false, true

				default:
					ErrorLog(`The error of converting of "` + text + `" to bool`)
				}
			}
		}
	}

	return false, false
}

func boolProperty(properties Properties, tag string, session Session) (bool, bool) {
	return valueToBool(properties.getRaw(tag), session)
}

func valueToInt(value any, session Session, defaultValue int) (int, bool) {
	if value != nil {
		switch value := value.(type) {
		case string:
			if text, ok := session.resolveConstants(value); ok {
				n, err := strconv.Atoi(strings.Trim(text, " \t"))
				if err == nil {
					return n, true
				}
				ErrorLog(err.Error())
			} else {
				n, err := strconv.Atoi(strings.Trim(value, " \t"))
				if err == nil {
					return n, true
				}
				ErrorLog(err.Error())
			}

		default:
			return isInt(value)
		}
	}

	return defaultValue, false
}

func intProperty(properties Properties, tag string, session Session, defaultValue int) (int, bool) {
	return valueToInt(properties.getRaw(tag), session, defaultValue)
}

func valueToFloat(value any, session Session, defaultValue float64) (float64, bool) {
	if value != nil {
		switch value := value.(type) {
		case float64:
			return value, true

		case string:
			if text, ok := session.resolveConstants(value); ok {
				f, err := strconv.ParseFloat(text, 64)
				if err == nil {
					return f, true
				}
				ErrorLog(err.Error())
			}
		}
	}

	return defaultValue, false
}

func floatProperty(properties Properties, tag string, session Session, defaultValue float64) (float64, bool) {
	return valueToFloat(properties.getRaw(tag), session, defaultValue)
}

func valueToFloatText(value any, session Session, defaultValue float64) (string, bool) {
	if value != nil {
		switch value := value.(type) {
		case float64:
			return fmt.Sprintf("%g", value), true

		case string:
			if text, ok := session.resolveConstants(value); ok {
				if _, err := strconv.ParseFloat(text, 64); err != nil {
					ErrorLog(err.Error())
					return fmt.Sprintf("%g", defaultValue), false
				}
				return text, true
			}
		}
	}

	return fmt.Sprintf("%g", defaultValue), false
}

func floatTextProperty(properties Properties, tag string, session Session, defaultValue float64) (string, bool) {
	return valueToFloatText(properties.getRaw(tag), session, defaultValue)
}

func valueToRange(value any, session Session) (Range, bool) {
	if value != nil {
		switch value := value.(type) {
		case Range:
			return value, true

		case int:
			return Range{First: value, Last: value}, true

		case string:
			if text, ok := session.resolveConstants(value); ok {
				var result Range
				if result.setValue(text) {
					return result, true
				}
			}
		}
	}
	return Range{}, false
}

func rangeProperty(properties Properties, tag string, session Session) (Range, bool) {
	return valueToRange(properties.getRaw(tag), session)
}
