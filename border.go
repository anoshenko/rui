package rui

import (
	"fmt"
	"strings"
)

// Constants related to view's border description
const (
	// NoneLine constant specifies that there is no border
	NoneLine = 0
	// SolidLine constant specifies the border/line as a solid line
	SolidLine = 1
	// DashedLine constant specifies the border/line as a dashed line
	DashedLine = 2
	// DottedLine constant specifies the border/line as a dotted line
	DottedLine = 3
	// DoubleLine constant specifies the border/line as a double solid line
	DoubleLine = 4
	// DoubleLine constant specifies the border/line as a double solid line
	WavyLine = 5

	// LeftStyle is the constant for "left-style" property tag.
	LeftStyle = "left-style"
	// RightStyle is the constant for "-right-style" property tag.
	RightStyle = "right-style"
	// TopStyle is the constant for "top-style" property tag.
	TopStyle = "top-style"
	// BottomStyle is the constant for "bottom-style" property tag.
	BottomStyle = "bottom-style"
	// LeftWidth is the constant for "left-width" property tag.
	LeftWidth = "left-width"
	// RightWidth is the constant for "-right-width" property tag.
	RightWidth = "right-width"
	// TopWidth is the constant for "top-width" property tag.
	TopWidth = "top-width"
	// BottomWidth is the constant for "bottom-width" property tag.
	BottomWidth = "bottom-width"
	// LeftColor is the constant for "left-color" property tag.
	LeftColor = "left-color"
	// RightColor is the constant for "-right-color" property tag.
	RightColor = "right-color"
	// TopColor is the constant for "top-color" property tag.
	TopColor = "top-color"
	// BottomColor is the constant for "bottom-color" property tag.
	BottomColor = "bottom-color"
)

// BorderProperty is the interface of a view border data
type BorderProperty interface {
	Properties
	fmt.Stringer
	stringWriter

	// ViewBorders returns top, right, bottom and left borders information all together
	ViewBorders(session Session) ViewBorders

	delete(tag string)
	cssStyle(builder cssBuilder, session Session)
	cssWidth(builder cssBuilder, session Session)
	cssColor(builder cssBuilder, session Session)
	cssStyleValue(session Session) string
	cssWidthValue(session Session) string
	cssColorValue(session Session) string
}

type borderProperty struct {
	propertyList
}

func newBorderProperty(value any) BorderProperty {
	border := new(borderProperty)
	border.properties = map[string]any{}

	if value != nil {
		switch value := value.(type) {
		case BorderProperty:
			return value

		case DataNode:
			if value.Type() == ObjectNode {
				_ = border.setBorderObject(value.Object())
			} else {
				return nil
			}

		case DataObject:
			_ = border.setBorderObject(value)

		case ViewBorder:
			border.properties[Style] = value.Style
			border.properties[Width] = value.Width
			border.properties[ColorTag] = value.Color

		case ViewBorders:
			if value.Left.Style == value.Right.Style &&
				value.Left.Style == value.Top.Style &&
				value.Left.Style == value.Bottom.Style {
				border.properties[Style] = value.Left.Style
			} else {
				border.properties[LeftStyle] = value.Left.Style
				border.properties[RightStyle] = value.Right.Style
				border.properties[TopStyle] = value.Top.Style
				border.properties[BottomStyle] = value.Bottom.Style
			}
			if value.Left.Width.Equal(value.Right.Width) &&
				value.Left.Width.Equal(value.Top.Width) &&
				value.Left.Width.Equal(value.Bottom.Width) {
				border.properties[Width] = value.Left.Width
			} else {
				border.properties[LeftWidth] = value.Left.Width
				border.properties[RightWidth] = value.Right.Width
				border.properties[TopWidth] = value.Top.Width
				border.properties[BottomWidth] = value.Bottom.Width
			}
			if value.Left.Color == value.Right.Color &&
				value.Left.Color == value.Top.Color &&
				value.Left.Color == value.Bottom.Color {
				border.properties[ColorTag] = value.Left.Color
			} else {
				border.properties[LeftColor] = value.Left.Color
				border.properties[RightColor] = value.Right.Color
				border.properties[TopColor] = value.Top.Color
				border.properties[BottomColor] = value.Bottom.Color
			}

		default:
			invalidPropertyValue(Border, value)
			return nil
		}
	}
	return border
}

// NewBorder creates the new BorderProperty
func NewBorder(params Params) BorderProperty {
	border := new(borderProperty)
	border.properties = map[string]any{}
	if params != nil {
		for _, tag := range []string{Style, Width, ColorTag, Left, Right, Top, Bottom,
			LeftStyle, RightStyle, TopStyle, BottomStyle,
			LeftWidth, RightWidth, TopWidth, BottomWidth,
			LeftColor, RightColor, TopColor, BottomColor} {
			if value, ok := params[tag]; ok && value != nil {
				border.Set(tag, value)
			}
		}
	}
	return border
}

func (border *borderProperty) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case BorderLeft, CellBorderLeft:
		return Left

	case BorderRight, CellBorderRight:
		return Right

	case BorderTop, CellBorderTop:
		return Top

	case BorderBottom, CellBorderBottom:
		return Bottom

	case BorderStyle, CellBorderStyle:
		return Style

	case BorderLeftStyle, CellBorderLeftStyle, "style-left":
		return LeftStyle

	case BorderRightStyle, CellBorderRightStyle, "style-right":
		return RightStyle

	case BorderTopStyle, CellBorderTopStyle, "style-top":
		return TopStyle

	case BorderBottomStyle, CellBorderBottomStyle, "style-bottom":
		return BottomStyle

	case BorderWidth, CellBorderWidth:
		return Width

	case BorderLeftWidth, CellBorderLeftWidth, "width-left":
		return LeftWidth

	case BorderRightWidth, CellBorderRightWidth, "width-right":
		return RightWidth

	case BorderTopWidth, CellBorderTopWidth, "width-top":
		return TopWidth

	case BorderBottomWidth, CellBorderBottomWidth, "width-bottom":
		return BottomWidth

	case BorderColor, CellBorderColor:
		return ColorTag

	case BorderLeftColor, CellBorderLeftColor, "color-left":
		return LeftColor

	case BorderRightColor, CellBorderRightColor, "color-right":
		return RightColor

	case BorderTopColor, CellBorderTopColor, "color-top":
		return TopColor

	case BorderBottomColor, CellBorderBottomColor, "color-bottom":
		return BottomColor
	}

	return tag
}

func (border *borderProperty) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("_{ ")
	comma := false

	write := func(tag string, value any) {
		if comma {
			buffer.WriteString(", ")
		}
		buffer.WriteString(tag)
		buffer.WriteString(" = ")
		writePropertyValue(buffer, BorderStyle, value, indent)
		comma = true
	}

	for _, tag := range []string{Style, Width, ColorTag} {
		if value, ok := border.properties[tag]; ok {
			write(tag, value)
		}
	}

	for _, side := range []string{Top, Right, Bottom, Left} {
		style, okStyle := border.properties[side+"-"+Style]
		width, okWidth := border.properties[side+"-"+Width]
		color, okColor := border.properties[side+"-"+ColorTag]
		if okStyle || okWidth || okColor {
			if comma {
				buffer.WriteString(", ")
				comma = false
			}

			buffer.WriteString(side)
			buffer.WriteString(" = _{ ")
			if okStyle {
				write(Style, style)
			}
			if okWidth {
				write(Width, width)
			}
			if okColor {
				write(ColorTag, color)
			}
			buffer.WriteString(" }")
			comma = true
		}
	}

	buffer.WriteString(" }")
}

func (border *borderProperty) String() string {
	return runStringWriter(border)
}

func (border *borderProperty) setSingleBorderObject(prefix string, obj DataObject) bool {
	result := true
	if text, ok := obj.PropertyValue(Style); ok {
		if !border.setEnumProperty(prefix+"-style", text, enumProperties[BorderStyle].values) {
			result = false
		}
	}
	if text, ok := obj.PropertyValue(ColorTag); ok {
		if !border.setColorProperty(prefix+"-color", text) {
			result = false
		}
	}
	if text, ok := obj.PropertyValue("width"); ok {
		if !border.setSizeProperty(prefix+"-width", text) {
			result = false
		}
	}
	return result
}

func (border *borderProperty) setBorderObject(obj DataObject) bool {
	result := true

	for _, side := range []string{Top, Right, Bottom, Left} {
		if node := obj.PropertyByTag(side); node != nil {
			if node.Type() == ObjectNode {
				if !border.setSingleBorderObject(side, node.Object()) {
					result = false
				}
			} else {
				notCompatibleType(side, node)
				result = false
			}
		}
	}

	if text, ok := obj.PropertyValue(Style); ok {
		values := split4Values(text)
		styles := enumProperties[BorderStyle].values
		switch len(values) {
		case 1:
			if !border.setEnumProperty(Style, values[0], styles) {
				result = false
			}

		case 4:
			for n, tag := range [4]string{TopStyle, RightStyle, BottomStyle, LeftStyle} {
				if !border.setEnumProperty(tag, values[n], styles) {
					result = false
				}
			}

		default:
			notCompatibleType(Style, text)
			result = false
		}
	}

	if text, ok := obj.PropertyValue(ColorTag); ok {
		values := split4Values(text)
		switch len(values) {
		case 1:
			if !border.setColorProperty(ColorTag, values[0]) {
				return false
			}

		case 4:
			for n, tag := range [4]string{TopColor, RightColor, BottomColor, LeftColor} {
				if !border.setColorProperty(tag, values[n]) {
					return false
				}
			}

		default:
			notCompatibleType(ColorTag, text)
			result = false
		}
	}

	if text, ok := obj.PropertyValue(Width); ok {
		values := split4Values(text)
		switch len(values) {
		case 1:
			if !border.setSizeProperty(Width, values[0]) {
				result = false
			}

		case 4:
			for n, tag := range [4]string{TopWidth, RightWidth, BottomWidth, LeftWidth} {
				if !border.setSizeProperty(tag, values[n]) {
					result = false
				}
			}

		default:
			notCompatibleType(Width, text)
			result = false
		}
	}

	return result
}

func (border *borderProperty) Remove(tag string) {
	tag = border.normalizeTag(tag)

	switch tag {
	case Style:
		for _, t := range []string{tag, TopStyle, RightStyle, BottomStyle, LeftStyle} {
			delete(border.properties, t)
		}

	case Width:
		for _, t := range []string{tag, TopWidth, RightWidth, BottomWidth, LeftWidth} {
			delete(border.properties, t)
		}

	case ColorTag:
		for _, t := range []string{tag, TopColor, RightColor, BottomColor, LeftColor} {
			delete(border.properties, t)
		}

	case Left, Right, Top, Bottom:
		border.Remove(tag + "-style")
		border.Remove(tag + "-width")
		border.Remove(tag + "-color")

	case LeftStyle, RightStyle, TopStyle, BottomStyle:
		delete(border.properties, tag)
		if style, ok := border.properties[Style]; ok && style != nil {
			for _, t := range []string{TopStyle, RightStyle, BottomStyle, LeftStyle} {
				if t != tag {
					if _, ok := border.properties[t]; !ok {
						border.properties[t] = style
					}
				}
			}
		}

	case LeftWidth, RightWidth, TopWidth, BottomWidth:
		delete(border.properties, tag)
		if width, ok := border.properties[Width]; ok && width != nil {
			for _, t := range []string{TopWidth, RightWidth, BottomWidth, LeftWidth} {
				if t != tag {
					if _, ok := border.properties[t]; !ok {
						border.properties[t] = width
					}
				}
			}
		}

	case LeftColor, RightColor, TopColor, BottomColor:
		delete(border.properties, tag)
		if color, ok := border.properties[ColorTag]; ok && color != nil {
			for _, t := range []string{TopColor, RightColor, BottomColor, LeftColor} {
				if t != tag {
					if _, ok := border.properties[t]; !ok {
						border.properties[t] = color
					}
				}
			}
		}

	default:
		ErrorLogF(`"%s" property is not compatible with the BorderProperty`, tag)
	}
}

func (border *borderProperty) Set(tag string, value any) bool {
	if value == nil {
		border.Remove(tag)
		return true
	}

	tag = border.normalizeTag(tag)

	switch tag {
	case Style:
		if border.setEnumProperty(Style, value, enumProperties[BorderStyle].values) {
			for _, side := range []string{TopStyle, RightStyle, BottomStyle, LeftStyle} {
				delete(border.properties, side)
			}
			return true
		}

	case Width:
		if border.setSizeProperty(Width, value) {
			for _, side := range []string{TopWidth, RightWidth, BottomWidth, LeftWidth} {
				delete(border.properties, side)
			}
			return true
		}

	case ColorTag:
		if border.setColorProperty(ColorTag, value) {
			for _, side := range []string{TopColor, RightColor, BottomColor, LeftColor} {
				delete(border.properties, side)
			}
			return true
		}

	case LeftStyle, RightStyle, TopStyle, BottomStyle:
		return border.setEnumProperty(tag, value, enumProperties[BorderStyle].values)

	case LeftWidth, RightWidth, TopWidth, BottomWidth:
		return border.setSizeProperty(tag, value)

	case LeftColor, RightColor, TopColor, BottomColor:
		return border.setColorProperty(tag, value)

	case Left, Right, Top, Bottom:
		switch value := value.(type) {
		case string:
			if obj := ParseDataText(value); obj != nil {
				return border.setSingleBorderObject(tag, obj)
			}

		case DataObject:
			return border.setSingleBorderObject(tag, value)

		case BorderProperty:
			styleTag := tag + "-" + Style
			if style := value.Get(styleTag); value != nil {
				border.properties[styleTag] = style
			}
			colorTag := tag + "-" + ColorTag
			if color := value.Get(colorTag); value != nil {
				border.properties[colorTag] = color
			}
			widthTag := tag + "-" + Width
			if width := value.Get(widthTag); value != nil {
				border.properties[widthTag] = width
			}
			return true

		case ViewBorder:
			border.properties[tag+"-"+Style] = value.Style
			border.properties[tag+"-"+Width] = value.Width
			border.properties[tag+"-"+ColorTag] = value.Color
			return true
		}
		fallthrough

	default:
		ErrorLogF(`"%s" property is not compatible with the BorderProperty`, tag)
	}

	return false
}

func (border *borderProperty) Get(tag string) any {
	tag = border.normalizeTag(tag)

	if result, ok := border.properties[tag]; ok {
		return result
	}

	switch tag {
	case Left, Right, Top, Bottom:
		result := newBorderProperty(nil)
		if style, ok := border.properties[tag+"-"+Style]; ok {
			result.Set(Style, style)
		} else if style, ok := border.properties[Style]; ok {
			result.Set(Style, style)
		}
		if width, ok := border.properties[tag+"-"+Width]; ok {
			result.Set(Width, width)
		} else if width, ok := border.properties[Width]; ok {
			result.Set(Width, width)
		}
		if color, ok := border.properties[tag+"-"+ColorTag]; ok {
			result.Set(ColorTag, color)
		} else if color, ok := border.properties[ColorTag]; ok {
			result.Set(ColorTag, color)
		}
		return result

	case LeftStyle, RightStyle, TopStyle, BottomStyle:
		if style, ok := border.properties[tag]; ok {
			return style
		}
		return border.properties[Style]

	case LeftWidth, RightWidth, TopWidth, BottomWidth:
		if width, ok := border.properties[tag]; ok {
			return width
		}
		return border.properties[Width]

	case LeftColor, RightColor, TopColor, BottomColor:
		if color, ok := border.properties[tag]; ok {
			return color
		}
		return border.properties[ColorTag]
	}

	return nil
}

func (border *borderProperty) delete(tag string) {
	tag = border.normalizeTag(tag)
	remove := []string{}

	switch tag {
	case Style:
		remove = []string{Style, LeftStyle, RightStyle, TopStyle, BottomStyle}

	case Width:
		remove = []string{Width, LeftWidth, RightWidth, TopWidth, BottomWidth}

	case ColorTag:
		remove = []string{ColorTag, LeftColor, RightColor, TopColor, BottomColor}

	case Left, Right, Top, Bottom:
		if border.Get(Style) != nil {
			border.properties[tag+"-"+Style] = 0
			remove = []string{tag + "-" + ColorTag, tag + "-" + Width}
		} else {
			remove = []string{tag + "-" + Style, tag + "-" + ColorTag, tag + "-" + Width}
		}

	case LeftStyle, RightStyle, TopStyle, BottomStyle:
		if border.Get(Style) != nil {
			border.properties[tag] = 0
		} else {
			remove = []string{tag}
		}

	case LeftWidth, RightWidth, TopWidth, BottomWidth:
		if border.Get(Width) != nil {
			border.properties[tag] = AutoSize()
		} else {
			remove = []string{tag}
		}

	case LeftColor, RightColor, TopColor, BottomColor:
		if border.Get(ColorTag) != nil {
			border.properties[tag] = 0
		} else {
			remove = []string{tag}
		}
	}

	for _, tag := range remove {
		delete(border.properties, tag)
	}
}

func (border *borderProperty) ViewBorders(session Session) ViewBorders {

	defStyle, _ := valueToEnum(border.getRaw(Style), BorderStyle, session, NoneLine)
	defWidth, _ := sizeProperty(border, Width, session)
	defColor, _ := colorProperty(border, ColorTag, session)

	getBorder := func(prefix string) ViewBorder {
		var result ViewBorder
		var ok bool
		if result.Style, ok = valueToEnum(border.getRaw(prefix+Style), BorderStyle, session, NoneLine); !ok {
			result.Style = defStyle
		}
		if result.Width, ok = sizeProperty(border, prefix+Width, session); !ok {
			result.Width = defWidth
		}
		if result.Color, ok = colorProperty(border, prefix+ColorTag, session); !ok {
			result.Color = defColor
		}
		return result
	}

	return ViewBorders{
		Top:    getBorder("top-"),
		Left:   getBorder("left-"),
		Right:  getBorder("right-"),
		Bottom: getBorder("bottom-"),
	}
}

func (border *borderProperty) cssStyle(builder cssBuilder, session Session) {
	borders := border.ViewBorders(session)
	values := enumProperties[BorderStyle].cssValues
	if borders.Top.Style == borders.Right.Style &&
		borders.Top.Style == borders.Left.Style &&
		borders.Top.Style == borders.Bottom.Style {
		builder.add(BorderStyle, values[borders.Top.Style])
	} else {
		builder.addValues(BorderStyle, " ", values[borders.Top.Style],
			values[borders.Right.Style], values[borders.Bottom.Style], values[borders.Left.Style])
	}
}

func (border *borderProperty) cssWidth(builder cssBuilder, session Session) {
	borders := border.ViewBorders(session)
	if borders.Top.Width == borders.Right.Width &&
		borders.Top.Width == borders.Left.Width &&
		borders.Top.Width == borders.Bottom.Width {
		if borders.Top.Width.Type != Auto {
			builder.add("border-width", borders.Top.Width.cssString("0", session))
		}
	} else {
		builder.addValues("border-width", " ",
			borders.Top.Width.cssString("0", session),
			borders.Right.Width.cssString("0", session),
			borders.Bottom.Width.cssString("0", session),
			borders.Left.Width.cssString("0", session))
	}
}

func (border *borderProperty) cssColor(builder cssBuilder, session Session) {
	borders := border.ViewBorders(session)
	if borders.Top.Color == borders.Right.Color &&
		borders.Top.Color == borders.Left.Color &&
		borders.Top.Color == borders.Bottom.Color {
		if borders.Top.Color != 0 {
			builder.add("border-color", borders.Top.Color.cssString())
		}
	} else {
		builder.addValues("border-color", " ", borders.Top.Color.cssString(),
			borders.Right.Color.cssString(), borders.Bottom.Color.cssString(), borders.Left.Color.cssString())
	}
}

func (border *borderProperty) cssStyleValue(session Session) string {
	var builder cssValueBuilder
	border.cssStyle(&builder, session)
	return builder.finish()
}

func (border *borderProperty) cssWidthValue(session Session) string {
	var builder cssValueBuilder
	border.cssWidth(&builder, session)
	return builder.finish()
}

func (border *borderProperty) cssColorValue(session Session) string {
	var builder cssValueBuilder
	border.cssColor(&builder, session)
	return builder.finish()
}

// ViewBorder describes parameters of a view border
type ViewBorder struct {
	// Style of the border line
	Style int

	// Color of the border line
	Color Color

	// Width of the border line
	Width SizeUnit
}

// ViewBorders describes the top, right, bottom, and left border of a view
type ViewBorders struct {
	Top, Right, Bottom, Left ViewBorder
}

// AllTheSame returns true if all borders are the same
func (border *ViewBorders) AllTheSame() bool {
	return border.Top.Style == border.Right.Style &&
		border.Top.Style == border.Left.Style &&
		border.Top.Style == border.Bottom.Style &&
		border.Top.Color == border.Right.Color &&
		border.Top.Color == border.Left.Color &&
		border.Top.Color == border.Bottom.Color &&
		border.Top.Width.Equal(border.Right.Width) &&
		border.Top.Width.Equal(border.Left.Width) &&
		border.Top.Width.Equal(border.Bottom.Width)
}

func getBorder(style Properties, tag string) BorderProperty {
	if value := style.Get(tag); value != nil {
		if border, ok := value.(BorderProperty); ok {
			return border
		}
	}
	return nil
}
