package rui

import (
	"strings"
)

func getOrientation(style Properties, session Session) (int, bool) {
	if value := style.Get(Orientation); value != nil {
		switch value := value.(type) {
		case int:
			return value, true

		case string:
			text, ok := session.resolveConstants(value)
			if !ok {
				return 0, false
			}

			text = strings.ToLower(strings.Trim(text, " \t\n\r"))
			switch text {
			case "vertical":
				return TopDownOrientation, true

			case "horizontal":
				return StartToEndOrientation, true
			}

			if result, ok := enumStringToInt(text, enumProperties[Orientation].values, true); ok {
				return result, true
			}
		}
	}

	return 0, false
}

func (style *viewStyle) Get(tag string) interface{} {
	return style.get(strings.ToLower(tag))
}

func (style *viewStyle) get(tag string) interface{} {
	switch tag {
	case Border, CellBorder:
		return getBorder(&style.propertyList, tag)

	case BorderLeft, BorderRight, BorderTop, BorderBottom,
		BorderStyle, BorderLeftStyle, BorderRightStyle, BorderTopStyle, BorderBottomStyle,
		BorderColor, BorderLeftColor, BorderRightColor, BorderTopColor, BorderBottomColor,
		BorderWidth, BorderLeftWidth, BorderRightWidth, BorderTopWidth, BorderBottomWidth:
		if border := getBorder(style, Border); border != nil {
			return border.Get(tag)
		}
		return nil

	case CellBorderLeft, CellBorderRight, CellBorderTop, CellBorderBottom,
		CellBorderStyle, CellBorderLeftStyle, CellBorderRightStyle, CellBorderTopStyle, CellBorderBottomStyle,
		CellBorderColor, CellBorderLeftColor, CellBorderRightColor, CellBorderTopColor, CellBorderBottomColor,
		CellBorderWidth, CellBorderLeftWidth, CellBorderRightWidth, CellBorderTopWidth, CellBorderBottomWidth:
		if border := getBorder(style, CellBorder); border != nil {
			return border.Get(tag)
		}
		return nil

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		return getRadiusElement(style, tag)

	case ColumnSeparator:
		if val, ok := style.properties[ColumnSeparator]; ok {
			return val.(ColumnSeparatorProperty)
		}
		return nil

	case ColumnSeparatorStyle, ColumnSeparatorWidth, ColumnSeparatorColor:
		if val, ok := style.properties[ColumnSeparator]; ok {
			separator := val.(ColumnSeparatorProperty)
			return separator.Get(tag)
		}
		return nil
	}

	return style.propertyList.getRaw(tag)
}
