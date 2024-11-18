package rui

import (
	"strings"
)

func setTransitionProperty(properties Properties, value any) bool {

	transitions := map[PropertyName]Animation{}

	setObject := func(obj DataObject) bool {
		if obj != nil {
			tag := strings.ToLower(obj.Tag())
			switch tag {
			case "", "_":
				ErrorLog("Invalid transition property name")

			default:
				transitions[PropertyName(tag)] = parseAnimation(obj)
				return true
			}
		}
		return false
	}

	switch value := value.(type) {
	case Params:
		for tag, val := range value {
			tag = defaultNormalize(tag)
			if tag == "" {
				ErrorLog("Invalid transition property name")
				return false
			}

			if val != nil {
				if animation, ok := val.(Animation); ok {
					transitions[PropertyName(tag)] = animation
				} else {
					notCompatibleType(Transition, val)
					return false
				}
			}
		}
		if len(transitions) == 0 {
			transitions = nil
		}
		properties.setRaw(Transition, transitions)
		return true

	case DataObject:
		if setObject(value) {
			properties.setRaw(Transition, transitions)
			return true
		}
		return false

	case DataNode:
		switch value.Type() {
		case ObjectNode:
			if setObject(value.Object()) {
				properties.setRaw(Transition, transitions)
				return true
			}
			return false

		case ArrayNode:
			for i := 0; i < value.ArraySize(); i++ {
				if obj := value.ArrayElement(i).Object(); obj != nil {
					if !setObject(obj) {
						return false
					}
				} else {
					notCompatibleType(Transition, value.ArrayElement(i))
					return false
				}
			}
			if len(transitions) == 0 {
				transitions = nil
			}
			properties.setRaw(Transition, transitions)
			return true
		}
	}

	notCompatibleType(Transition, value)
	return false

}

/*
	func (style *viewStyle) setTransition(tag PropertyName, value any) bool {
		setObject := func(obj DataObject) bool {
			if obj != nil {
				tag := defaultNormalize(tag)
				switch tag {
				case "", "_":
					ErrorLog("Invalid transition property name")

				default:
					style.transitions[tag] = parseAnimation(obj)
					return true
				}
			}
			return false
		}

		switch value := value.(type) {
		case Params:
			result := true
			for tag, val := range value {
				tag = defaultNormalize(tag)
				if tag == "" {
					ErrorLog("Invalid transition property name")
					result = false
				} else if val == nil {
					delete(style.transitions, tag)
				} else if animation, ok := val.(Animation); ok {
					style.transitions[tag] = animation
				} else {
					notCompatibleType(Transition, val)
					result = false
				}
			}
			return result

		case DataObject:
			return setObject(value)

		case DataNode:
			switch value.Type() {
			case ObjectNode:
				return setObject(value.Object())

			case ArrayNode:
				result := true
				for i := 0; i < value.ArraySize(); i++ {
					if obj := value.ArrayElement(i).Object(); obj != nil {
						result = setObject(obj) && result
					} else {
						notCompatibleType(tag, value.ArrayElement(i))
						result = false
					}
				}
				return result
			}
		}

		notCompatibleType(tag, value)
		return false
	}
*/

func viewStyleRemove(properties Properties, tag PropertyName) []PropertyName {
	switch tag {
	case BorderStyle, BorderColor, BorderWidth,
		BorderLeft, BorderLeftStyle, BorderLeftColor, BorderLeftWidth,
		BorderRight, BorderRightStyle, BorderRightColor, BorderRightWidth,
		BorderTop, BorderTopStyle, BorderTopColor, BorderTopWidth,
		BorderBottom, BorderBottomStyle, BorderBottomColor, BorderBottomWidth:
		if border := getBorderProperty(properties, Border); border != nil && border.deleteTag(tag) {
			return []PropertyName{Border}
		}

	case CellBorderStyle, CellBorderColor, CellBorderWidth,
		CellBorderLeft, CellBorderLeftStyle, CellBorderLeftColor, CellBorderLeftWidth,
		CellBorderRight, CellBorderRightStyle, CellBorderRightColor, CellBorderRightWidth,
		CellBorderTop, CellBorderTopStyle, CellBorderTopColor, CellBorderTopWidth,
		CellBorderBottom, CellBorderBottomStyle, CellBorderBottomColor, CellBorderBottomWidth:
		if border := getBorderProperty(properties, CellBorder); border != nil && border.deleteTag(tag) {
			return []PropertyName{CellBorder}
		}

	case MarginTop, MarginRight, MarginBottom, MarginLeft:
		return removeBoundsPropertySide(properties, Margin, tag)

	case PaddingTop, PaddingRight, PaddingBottom, PaddingLeft:
		return removeBoundsPropertySide(properties, Padding, tag)

	case CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft:
		return removeBoundsPropertySide(properties, CellPadding, tag)

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		if removeRadiusPropertyElement(properties, tag) {
			return []PropertyName{Radius, tag}
		}

	case OutlineStyle, OutlineWidth, OutlineColor:
		if outline := getOutlineProperty(properties); outline != nil {
			outline.Remove(tag)
			if outline.empty() {
				properties.setRaw(Outline, nil)
			}
			return []PropertyName{Outline, tag}
		}

	default:
		return propertiesRemove(properties, tag)
	}

	return []PropertyName{}
}

func viewStyleSet(style Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case Shadow, TextShadow:
		if setShadowProperty(style, tag, value) {
			return []PropertyName{tag}
		}

	case Background:
		return setBackgroundProperty(style, value)

	case Border, CellBorder:
		if border := newBorderProperty(value); border != nil {
			style.setRaw(tag, border)
			return []PropertyName{tag}
		} else {
			return nil
		}

	case BorderStyle, BorderColor, BorderWidth,
		BorderLeft, BorderLeftStyle, BorderLeftColor, BorderLeftWidth,
		BorderRight, BorderRightStyle, BorderRightColor, BorderRightWidth,
		BorderTop, BorderTopStyle, BorderTopColor, BorderTopWidth,
		BorderBottom, BorderBottomStyle, BorderBottomColor, BorderBottomWidth:

		return setBorderPropertyElement(style, Border, tag, value)

	case CellBorderStyle, CellBorderColor, CellBorderWidth,
		CellBorderLeft, CellBorderLeftStyle, CellBorderLeftColor, CellBorderLeftWidth,
		CellBorderRight, CellBorderRightStyle, CellBorderRightColor, CellBorderRightWidth,
		CellBorderTop, CellBorderTopStyle, CellBorderTopColor, CellBorderTopWidth,
		CellBorderBottom, CellBorderBottomStyle, CellBorderBottomColor, CellBorderBottomWidth:

		return setBorderPropertyElement(style, CellBorder, tag, value)

	case Radius:
		return setRadiusProperty(style, value)

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		if setRadiusPropertyElement(style, tag, value) {
			return []PropertyName{Radius, tag}
		}

	case Margin, Padding, CellPadding:
		return setBoundsProperty(style, tag, value)

	case MarginTop, MarginRight, MarginBottom, MarginLeft:
		return setBoundsPropertySide(style, Margin, tag, value)

	case PaddingTop, PaddingRight, PaddingBottom, PaddingLeft:
		return setBoundsPropertySide(style, Padding, tag, value)

	case CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft:
		return setBoundsPropertySide(style, CellPadding, tag, value)

	case HeadStyle, FootStyle:
		switch value := value.(type) {
		case string:
			style.setRaw(tag, value)

		case Params:
			style.setRaw(tag, value)

		case DataObject:
			if params := value.ToParams(); len(params) > 0 {
				style.setRaw(tag, params)
			} else {
				style.setRaw(tag, nil)
			}

		default:
			notCompatibleType(tag, value)
			return nil
		}
		return []PropertyName{tag}

	case CellStyle, ColumnStyle, RowStyle:
		switch value := value.(type) {
		case string:
			style.setRaw(tag, value)

		case Params:
			style.setRaw(tag, value)

		case DataObject:
			if params := value.ToParams(); len(params) > 0 {
				style.setRaw(tag, params)
			} else {
				style.setRaw(tag, nil)
			}

		case DataNode:
			switch value.Type() {
			case TextNode:
				if text := value.Text(); text != "" {
					style.setRaw(tag, text)
				} else {
					style.setRaw(tag, nil)
				}

			case ObjectNode:
				if obj := value.Object(); obj != nil {
					if params := obj.ToParams(); len(params) > 0 {
						style.setRaw(tag, params)
					} else {
						style.setRaw(tag, nil)
					}
				} else {
					notCompatibleType(tag, value)
					return nil
				}

			default:
				notCompatibleType(tag, value)
				return nil
			}

		default:
			notCompatibleType(tag, value)
			return nil
		}
		return []PropertyName{tag}

	case Outline:
		return setOutlineProperty(style, value)

	case OutlineStyle, OutlineWidth, OutlineColor:
		if outline := getOutlineProperty(style); outline != nil {
			if outline.Set(tag, value) {
				return []PropertyName{Outline, tag}
			}
		} else {
			outline := NewOutlineProperty(nil)
			if outline.Set(tag, value) {
				style.setRaw(Outline, outline)
				return []PropertyName{Outline, tag}
			}
		}
		return nil

	case Transform:
		if setTransformProperty(style, value) {
			return []PropertyName{Transform}
		} else {
			return nil
		}

	case Perspective, RotateX, RotateY, RotateZ, Rotate, SkewX, SkewY, ScaleX, ScaleY, ScaleZ,
		TranslateX, TranslateY, TranslateZ:
		return setTransformPropertyElement(style, tag, value)

	case Orientation:
		if text, ok := value.(string); ok {
			switch strings.ToLower(text) {
			case "vertical":
				style.setRaw(Orientation, TopDownOrientation)
				return []PropertyName{Orientation}

			case "horizontal":
				style.setRaw(Orientation, StartToEndOrientation)
				return []PropertyName{Orientation}
			}
		}

	case TextWeight:
		if n, ok := value.(int); ok {
			if n >= 100 && n%100 == 0 {
				n /= 100
				if n > 0 && n <= 9 {
					style.setRaw(TextWeight, n)
					return []PropertyName{TextWeight}
				}
			}
		}

	case Row, Column:
		return setRangeProperty(style, tag, value)

	case CellWidth, CellHeight:
		return setGridCellSize(style, tag, value)

	case ColumnSeparator:
		if separator := newColumnSeparatorProperty(value); separator != nil {
			style.setRaw(ColumnSeparator, separator)
			return []PropertyName{tag}
		}
		return nil

	case ColumnSeparatorStyle, ColumnSeparatorWidth, ColumnSeparatorColor:
		if separator := getColumnSeparatorProperty(style); separator != nil {
			if separator.Set(tag, value) {
				return []PropertyName{ColumnSeparator, tag}
			}
		} else {
			separator := newColumnSeparatorProperty(nil)
			if separator.Set(tag, value) {
				style.setRaw(ColumnSeparator, separator)
				return []PropertyName{ColumnSeparator, tag}
			}
		}
		return nil

	case Clip, ShapeOutside:
		return setClipShapeProperty(style, tag, value)

	case Filter, BackdropFilter:
		return setFilterProperty(style, tag, value)

	case Transition:
		if setTransitionProperty(style, value) {
			return []PropertyName{tag}
		} else {
			return nil
		}

	case AnimationTag:
		if setAnimationProperty(style, tag, value) {
			return []PropertyName{tag}
		} else {
			return nil
		}
	}

	return propertiesSet(style, tag, value)
}

func (style *viewStyle) Set(tag PropertyName, value any) bool {
	if value == nil {
		style.Remove(tag)
		return true
	}

	return viewStyleSet(style, normalizeViewStyleTag(tag), value) != nil
}

func (style *viewStyle) Remove(tag PropertyName) {
	viewStyleRemove(style, normalizeViewStyleTag(tag))
}
