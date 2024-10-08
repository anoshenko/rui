package rui

import (
	"strings"
)

func (style *viewStyle) setRange(tag string, value any) bool {
	switch value := value.(type) {
	case string:
		if strings.Contains(value, "@") {
			style.properties[tag] = value
			return true
		}
		var r Range
		if !r.setValue(value) {
			invalidPropertyValue(tag, value)
			return false
		}
		style.properties[tag] = r

	case int:
		style.properties[tag] = Range{First: value, Last: value}

	case Range:
		style.properties[tag] = value

	default:
		notCompatibleType(tag, value)
		return false
	}
	return true
}

func (style *viewStyle) setBackground(value any) bool {
	background := []BackgroundElement{}

	switch value := value.(type) {
	case BackgroundElement:
		background = []BackgroundElement{value}

	case []BackgroundElement:
		background = value

	case []DataValue:
		for _, el := range value {
			if el.IsObject() {
				if element := createBackground(el.Object()); element != nil {
					background = append(background, element)
				}
			} else if obj := ParseDataText(el.Value()); obj != nil {
				if element := createBackground(obj); element != nil {
					background = append(background, element)
				}
			}
		}

	case DataObject:
		if element := createBackground(value); element != nil {
			background = []BackgroundElement{element}
		}

	case []DataObject:
		for _, obj := range value {
			if element := createBackground(obj); element != nil {
				background = append(background, element)
			}
		}

	case string:
		if obj := ParseDataText(value); obj != nil {
			if element := createBackground(obj); element != nil {
				background = []BackgroundElement{element}
			}
		}
	}

	if len(background) > 0 {
		style.properties[Background] = background
		return true
	}
	return false
}

func (style *viewStyle) setTransition(tag string, value any) bool {
	setObject := func(obj DataObject) bool {
		if obj != nil {
			tag := strings.ToLower(tag)
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
			tag = strings.ToLower(strings.Trim(tag, " \t"))
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

func (style *viewStyle) setAnimation(tag string, value any) bool {

	set := func(animations []Animation) {
		style.properties[tag] = animations
		for _, animation := range animations {
			animation.used()
		}
	}

	switch value := value.(type) {
	case Animation:
		set([]Animation{value})
		return true

	case []Animation:
		set(value)
		return true

	case DataObject:
		if animation := parseAnimation(value); animation.hasAnimatedProperty() {
			set([]Animation{animation})
			return true
		}

	case DataNode:
		animations := []Animation{}
		result := true
		for i := 0; i < value.ArraySize(); i++ {
			if obj := value.ArrayElement(i).Object(); obj != nil {
				if anim := parseAnimation(obj); anim.hasAnimatedProperty() {
					animations = append(animations, anim)
				} else {
					result = false
				}
			} else {
				notCompatibleType(tag, value.ArrayElement(i))
				result = false
			}
		}
		if result && len(animations) > 0 {
			set(animations)
		}
		return result
	}

	notCompatibleType(tag, value)
	return false
}

func (style *viewStyle) Remove(tag string) {
	style.remove(strings.ToLower(tag))
}

func (style *viewStyle) remove(tag string) {
	switch tag {
	case BorderStyle, BorderColor, BorderWidth,
		BorderLeft, BorderLeftStyle, BorderLeftColor, BorderLeftWidth,
		BorderRight, BorderRightStyle, BorderRightColor, BorderRightWidth,
		BorderTop, BorderTopStyle, BorderTopColor, BorderTopWidth,
		BorderBottom, BorderBottomStyle, BorderBottomColor, BorderBottomWidth:
		if border := getBorder(style, Border); border != nil {
			border.delete(tag)
		}

	case CellBorderStyle, CellBorderColor, CellBorderWidth,
		CellBorderLeft, CellBorderLeftStyle, CellBorderLeftColor, CellBorderLeftWidth,
		CellBorderRight, CellBorderRightStyle, CellBorderRightColor, CellBorderRightWidth,
		CellBorderTop, CellBorderTopStyle, CellBorderTopColor, CellBorderTopWidth,
		CellBorderBottom, CellBorderBottomStyle, CellBorderBottomColor, CellBorderBottomWidth:
		if border := getBorder(style, CellBorder); border != nil {
			border.delete(tag)
		}

	case MarginTop, MarginRight, MarginBottom, MarginLeft,
		"top-margin", "right-margin", "bottom-margin", "left-margin":
		style.removeBoundsSide(Margin, tag)

	case PaddingTop, PaddingRight, PaddingBottom, PaddingLeft,
		"top-padding", "right-padding", "bottom-padding", "left-padding":
		style.removeBoundsSide(Padding, tag)

	case CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft:
		style.removeBoundsSide(CellPadding, tag)

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		style.removeRadiusElement(tag)

	case OutlineStyle, OutlineWidth, OutlineColor:
		if outline := getOutline(style); outline != nil {
			outline.Remove(tag)
		}

	default:
		style.propertyList.remove(tag)
	}
}

func (style *viewStyle) Set(tag string, value any) bool {
	return style.set(strings.ToLower(tag), value)
}

func (style *viewStyle) set(tag string, value any) bool {
	if value == nil {
		style.remove(tag)
		return true
	}

	switch tag {
	case Shadow, TextShadow:
		return style.setShadow(tag, value)

	case Background:
		return style.setBackground(value)

	case Border, CellBorder:
		if border := newBorderProperty(value); border != nil {
			style.properties[tag] = border
			return true
		}

	case BorderStyle, BorderColor, BorderWidth,
		BorderLeft, BorderLeftStyle, BorderLeftColor, BorderLeftWidth,
		BorderRight, BorderRightStyle, BorderRightColor, BorderRightWidth,
		BorderTop, BorderTopStyle, BorderTopColor, BorderTopWidth,
		BorderBottom, BorderBottomStyle, BorderBottomColor, BorderBottomWidth:

		border := getBorder(style, Border)
		if border == nil {
			border = NewBorder(nil)
			if border.Set(tag, value) {
				style.properties[Border] = border
				return true
			}
			return false
		}
		return border.Set(tag, value)

	case CellBorderStyle, CellBorderColor, CellBorderWidth,
		CellBorderLeft, CellBorderLeftStyle, CellBorderLeftColor, CellBorderLeftWidth,
		CellBorderRight, CellBorderRightStyle, CellBorderRightColor, CellBorderRightWidth,
		CellBorderTop, CellBorderTopStyle, CellBorderTopColor, CellBorderTopWidth,
		CellBorderBottom, CellBorderBottomStyle, CellBorderBottomColor, CellBorderBottomWidth:

		border := getBorder(style, CellBorder)
		if border == nil {
			border = NewBorder(nil)
			if border.Set(tag, value) {
				style.properties[CellBorder] = border
				return true
			}
			return false
		}
		return border.Set(tag, value)

	case Radius:
		return style.setRadius(value)

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		return style.setRadiusElement(tag, value)

	case Margin, Padding, CellPadding:
		return style.setBounds(tag, value)

	case MarginTop, MarginRight, MarginBottom, MarginLeft,
		"top-margin", "right-margin", "bottom-margin", "left-margin":
		return style.setBoundsSide(Margin, tag, value)

	case PaddingTop, PaddingRight, PaddingBottom, PaddingLeft,
		"top-padding", "right-padding", "bottom-padding", "left-padding":
		return style.setBoundsSide(Padding, tag, value)

	case CellPaddingTop, CellPaddingRight, CellPaddingBottom, CellPaddingLeft:
		return style.setBoundsSide(CellPadding, tag, value)

	case HeadStyle, FootStyle:
		switch value := value.(type) {
		case string:
			style.properties[tag] = value
			return true

		case Params:
			style.properties[tag] = value
			return true

		case DataObject:
			if params := value.ToParams(); len(params) > 0 {
				style.properties[tag] = params
			}
			return true
		}

	case CellStyle, ColumnStyle, RowStyle:
		switch value := value.(type) {
		case string:
			style.properties[tag] = value
			return true

		case Params:
			style.properties[tag] = value
			return true

		case DataObject:
			if params := value.ToParams(); len(params) > 0 {
				style.properties[tag] = params
			}
			return true

		case DataNode:
			switch value.Type() {
			case TextNode:
				if text := value.Text(); text != "" {
					style.properties[tag] = text
				}
				return true

			case ObjectNode:
				if obj := value.Object(); obj != nil {
					if params := obj.ToParams(); len(params) > 0 {
						style.properties[tag] = params
					}
					return true
				}

			case ArrayNode:
				// TODO
			}
		}

	case Outline:
		return style.setOutline(value)

	case OutlineStyle, OutlineWidth, OutlineColor:
		if outline := getOutline(style); outline != nil {
			return outline.Set(tag, value)
		}
		style.properties[Outline] = NewOutlineProperty(Params{tag: value})
		return true

	case TransformTag:
		return style.setTransform(value)

	case RotateX, RotateY, RotateZ, Rotate, SkewX, SkewY, ScaleX, ScaleY, ScaleZ,
		TranslateX, TranslateY, TranslateZ:
		return style.setTransformProperty(tag, value)

	case Orientation:
		if text, ok := value.(string); ok {
			switch strings.ToLower(text) {
			case "vertical":
				style.properties[Orientation] = TopDownOrientation
				return true

			case "horizontal":
				style.properties[Orientation] = StartToEndOrientation
				return true
			}
		}

	case TextWeight:
		if n, ok := value.(int); ok && n >= 100 && n%100 == 0 {
			n /= 100
			if n > 0 && n <= 9 {
				style.properties[TextWeight] = n
				return true
			}
		}

	case Row, Column:
		return style.setRange(tag, value)

	case CellWidth, CellHeight:
		return style.setGridCellSize(tag, value)

	case ColumnSeparator:
		if separator := newColumnSeparatorProperty(value); separator != nil {
			style.properties[ColumnSeparator] = separator
			return true
		}
		return false

	case ColumnSeparatorStyle, ColumnSeparatorWidth, ColumnSeparatorColor:
		var separator ColumnSeparatorProperty = nil
		if val, ok := style.properties[ColumnSeparator]; ok {
			separator = val.(ColumnSeparatorProperty)
		}
		if separator == nil {
			separator = newColumnSeparatorProperty(nil)
		}

		if separator.Set(tag, value) {
			style.properties[ColumnSeparator] = separator
			return true
		}
		return false

	case Clip, ShapeOutside:
		return style.setClipShape(tag, value)

	case Filter, BackdropFilter:
		return style.setFilter(tag, value)

	case Transition:
		return style.setTransition(tag, value)

	case AnimationTag:
		return style.setAnimation(tag, value)
	}

	return style.propertyList.set(tag, value)
}
