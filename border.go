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
	//
	// Used by BorderProperty.
	// Left border line style.
	//
	// Supported types: int, string.
	//
	// Values:
	//  - 0 (NoneLine) or "none" - The border will not be drawn.
	//  - 1 (SolidLine) or "solid" - Solid line as a border.
	//  - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//  - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//  - 4 (DoubleLine) or "double" - Double line as a border.
	LeftStyle PropertyName = "left-style"

	// RightStyle is the constant for "right-style" property tag.
	//
	// Used by BorderProperty.
	// Right border line style.
	//
	// Supported types: int, string.
	//
	// Values:
	//  - 0 (NoneLine) or "none" - The border will not be drawn.
	//  - 1 (SolidLine) or "solid" - Solid line as a border.
	//  - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//  - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//  - 4 (DoubleLine) or "double" - Double line as a border.
	RightStyle PropertyName = "right-style"

	// TopStyle is the constant for "top-style" property tag.
	//
	// Used by BorderProperty.
	// Top border line style.
	//
	// Supported types: int, string.
	//
	// Values:
	//  - 0 (NoneLine) or "none" - The border will not be drawn.
	//  - 1 (SolidLine) or "solid" - Solid line as a border.
	//  - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//  - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//  - 4 (DoubleLine) or "double" - Double line as a border.
	TopStyle PropertyName = "top-style"

	// BottomStyle is the constant for "bottom-style" property tag.
	//
	// Used by BorderProperty.
	// Bottom border line style.
	//
	// Supported types: int, string.
	//
	// Values:
	//  - 0 (NoneLine) or "none" - The border will not be drawn.
	//  - 1 (SolidLine) or "solid" - Solid line as a border.
	//  - 2 (DashedLine) or "dashed" - Dashed line as a border.
	//  - 3 (DottedLine) or "dotted" - Dotted line as a border.
	//  - 4 (DoubleLine) or "double" - Double line as a border.
	BottomStyle PropertyName = "bottom-style"

	// LeftWidth is the constant for "left-width" property tag.
	//
	// Used by BorderProperty.
	// Left border line width.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	LeftWidth PropertyName = "left-width"

	// RightWidth is the constant for "right-width" property tag.
	//
	// Used by BorderProperty.
	// Right border line width.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	RightWidth PropertyName = "right-width"

	// TopWidth is the constant for "top-width" property tag.
	//
	// Used by BorderProperty.
	// Top border line width.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	TopWidth PropertyName = "top-width"

	// BottomWidth is the constant for "bottom-width" property tag.
	//
	// Used by BorderProperty.
	// Bottom border line width.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See [SizeUnit] description for more details.
	BottomWidth PropertyName = "bottom-width"

	// LeftColor is the constant for "left-color" property tag.
	//
	// Used by BorderProperty.
	// Left border line color.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	LeftColor PropertyName = "left-color"

	// RightColor is the constant for "right-color" property tag.
	//
	// Used by BorderProperty.
	// Right border line color.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	RightColor PropertyName = "right-color"

	// TopColor is the constant for "top-color" property tag.
	//
	// Used by BorderProperty.
	// Top border line color.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	TopColor PropertyName = "top-color"

	// BottomColor is the constant for "bottom-color" property tag.
	//
	// Used by BorderProperty.
	// Bottom border line color.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See [Color] description for more details.
	BottomColor PropertyName = "bottom-color"
)

// BorderProperty is the interface of a view border data
type BorderProperty interface {
	Properties
	fmt.Stringer
	stringWriter

	// ViewBorders returns top, right, bottom and left borders information all together
	ViewBorders(session Session) ViewBorders

	deleteTag(tag PropertyName) bool
	cssStyle(builder cssBuilder, session Session)
	cssWidth(builder cssBuilder, session Session)
	cssColor(builder cssBuilder, session Session)
	cssStyleValue(session Session) string
	cssWidthValue(session Session) string
	cssColorValue(session Session) string
}

type borderProperty struct {
	dataProperty
}

func newBorderProperty(value any) BorderProperty {
	border := new(borderProperty)
	border.init()

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

// NewBorder creates the new BorderProperty.
// The following properties can be used:
//
// "style" (Style). Determines the line style (int). Valid values: 0 (NoneLine), 1 (SolidLine), 2 (DashedLine), 3 (DottedLine), or 4 (DoubleLine);
//
// "color" (ColorTag). Determines the line color (Color);
//
// "width" (Width). Determines the line thickness (SizeUnit).
func NewBorder(params Params) BorderProperty {
	border := new(borderProperty)
	border.init()

	if params != nil {
		for _, tag := range []PropertyName{Style, Width, ColorTag, Left, Right, Top, Bottom,
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

func (border *borderProperty) init() {
	border.dataProperty.init()
	border.normalize = normalizeBorderTag
	border.get = borderGet
	border.set = borderSet
	border.remove = borderRemove
	border.supportedProperties = []PropertyName{
		Left,
		Right,
		Top,
		Bottom,
		Style,
		LeftStyle,
		RightStyle,
		TopStyle,
		BottomStyle,
		Width,
		LeftWidth,
		RightWidth,
		TopWidth,
		BottomWidth,
		ColorTag,
		LeftColor,
		RightColor,
		TopColor,
		BottomColor,
	}
}

func normalizeBorderTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
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

	write := func(tag PropertyName, value any) {
		text := propertyValueToString(tag, value, indent)
		if text != "" {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(string(tag))
			buffer.WriteString(" = ")
			buffer.WriteString(text)
			comma = true
		}
	}

	for _, tag := range []PropertyName{Style, Width, ColorTag} {
		if value, ok := border.properties[tag]; ok {
			write(tag, value)
		}
	}

	for _, side := range []PropertyName{Top, Right, Bottom, Left} {
		style, okStyle := border.properties[side+"-"+Style]
		width, okWidth := border.properties[side+"-"+Width]
		color, okColor := border.properties[side+"-"+ColorTag]
		if okStyle || okWidth || okColor {
			if comma {
				buffer.WriteString(", ")
				comma = false
			}

			buffer.WriteString(string(side))
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

func (border *borderProperty) setBorderObject(obj DataObject) bool {
	result := true
	for node := range obj.Properties() {
		tag := PropertyName(node.Tag())
		switch node.Type() {
		case TextNode:
			if borderSet(border, tag, node.Text()) == nil {
				result = false
			}

		case ObjectNode:
			if borderSet(border, tag, node.Object()) == nil {
				result = false
			}

		default:
			result = false
		}
	}
	return result
}

func borderRemove(properties Properties, tag PropertyName) []PropertyName {
	result := []PropertyName{}
	removeTag := func(t PropertyName) {
		if properties.getRaw(t) != nil {
			properties.setRaw(t, nil)
			result = append(result, t)
		}
	}

	switch tag {
	case Style:
		for _, t := range []PropertyName{tag, TopStyle, RightStyle, BottomStyle, LeftStyle} {
			removeTag(t)
		}

	case Width:
		for _, t := range []PropertyName{tag, TopWidth, RightWidth, BottomWidth, LeftWidth} {
			removeTag(t)
		}

	case ColorTag:
		for _, t := range []PropertyName{tag, TopColor, RightColor, BottomColor, LeftColor} {
			removeTag(t)
		}

	case Left, Right, Top, Bottom:
		removeTag(tag + "-style")
		removeTag(tag + "-width")
		removeTag(tag + "-color")

	case LeftStyle, RightStyle, TopStyle, BottomStyle:
		removeTag(tag)
		if style := properties.getRaw(Style); style != nil {
			for _, t := range []PropertyName{TopStyle, RightStyle, BottomStyle, LeftStyle} {
				if t != tag {
					if properties.getRaw(t) == nil {
						properties.setRaw(t, style)
						result = append(result, t)
					}
				}
			}
		}

	case LeftWidth, RightWidth, TopWidth, BottomWidth:
		removeTag(tag)
		if width := properties.getRaw(Width); width != nil {
			for _, t := range []PropertyName{TopWidth, RightWidth, BottomWidth, LeftWidth} {
				if t != tag {
					if properties.getRaw(t) == nil {
						properties.setRaw(t, width)
						result = append(result, t)
					}
				}
			}
		}

	case LeftColor, RightColor, TopColor, BottomColor:
		removeTag(tag)
		if color := properties.getRaw(ColorTag); color != nil {
			for _, t := range []PropertyName{TopColor, RightColor, BottomColor, LeftColor} {
				if t != tag {
					if properties.getRaw(t) == nil {
						properties.setRaw(t, color)
						result = append(result, t)
					}
				}
			}
		}

	default:
		ErrorLogF(`"%s" property is not compatible with the BorderProperty`, tag)
	}

	return result
}

func borderSet(properties Properties, tag PropertyName, value any) []PropertyName {

	setSingleBorderObject := func(prefix PropertyName, obj DataObject) []PropertyName {
		result := []PropertyName{}
		if text, ok := obj.PropertyValue(string(Style)); ok {
			props := setEnumProperty(properties, prefix+"-style", text, enumProperties[BorderStyle].values)
			if props == nil {
				return nil
			}
			result = append(result, props...)
		}
		if text, ok := obj.PropertyValue(string(ColorTag)); ok {
			props := setColorProperty(properties, prefix+"-color", text)
			if props == nil && len(result) == 0 {
				return nil
			}
			result = append(result, props...)
		}
		if text, ok := obj.PropertyValue("width"); ok {
			props := setSizeProperty(properties, prefix+"-width", text)
			if props == nil && len(result) == 0 {
				return nil
			}
			result = append(result, props...)
		}
		if len(result) > 0 {
			result = append(result, prefix)
		}
		return result
	}

	switch tag {
	case Style:
		if result := setEnumProperty(properties, Style, value, enumProperties[BorderStyle].values); result != nil {
			for _, side := range []PropertyName{TopStyle, RightStyle, BottomStyle, LeftStyle} {
				if value := properties.getRaw(side); value != nil {
					properties.setRaw(side, nil)
					result = append(result, side)
				}
			}
			return result
		}

	case Width:
		if result := setSizeProperty(properties, Width, value); result != nil {
			for _, side := range []PropertyName{TopWidth, RightWidth, BottomWidth, LeftWidth} {
				if value := properties.getRaw(side); value != nil {
					properties.setRaw(side, nil)
					result = append(result, side)
				}
			}
			return result
		}

	case ColorTag:
		if result := setColorProperty(properties, ColorTag, value); result != nil {
			for _, side := range []PropertyName{TopColor, RightColor, BottomColor, LeftColor} {
				if value := properties.getRaw(side); value != nil {
					properties.setRaw(side, nil)
					result = append(result, side)
				}
			}
			return result
		}

	case LeftStyle, RightStyle, TopStyle, BottomStyle:
		return setEnumProperty(properties, tag, value, enumProperties[BorderStyle].values)

	case LeftWidth, RightWidth, TopWidth, BottomWidth:
		return setSizeProperty(properties, tag, value)

	case LeftColor, RightColor, TopColor, BottomColor:
		return setColorProperty(properties, tag, value)

	case Left, Right, Top, Bottom:
		switch value := value.(type) {
		case string:
			obj, err := ParseDataText(value)
			if err != nil {
				ErrorLog(err.Error())
			} else {
				return setSingleBorderObject(tag, obj)
			}

		case DataObject:
			return setSingleBorderObject(tag, value)

		case BorderProperty:
			result := []PropertyName{}
			styleTag := tag + "-" + Style
			if style := value.Get(styleTag); value != nil {
				properties.setRaw(styleTag, style)
				result = append(result, styleTag)
			}
			colorTag := tag + "-" + ColorTag
			if color := value.Get(colorTag); value != nil {
				properties.setRaw(colorTag, color)
				result = append(result, colorTag)
			}
			widthTag := tag + "-" + Width
			if width := value.Get(widthTag); value != nil {
				properties.setRaw(widthTag, width)
				result = append(result, widthTag)
			}
			return result

		case ViewBorder:
			properties.setRaw(tag+"-"+Style, value.Style)
			properties.setRaw(tag+"-"+Width, value.Width)
			properties.setRaw(tag+"-"+ColorTag, value.Color)
			return []PropertyName{tag + "-" + Style, tag + "-" + Width, tag + "-" + ColorTag}
		}
		fallthrough

	default:
		ErrorLogF(`"%s" property is not compatible with the BorderProperty`, tag)
	}

	return nil
}

func borderGet(properties Properties, tag PropertyName) any {
	if result := properties.getRaw(tag); result != nil {
		return result
	}

	switch tag {
	case Left, Right, Top, Bottom:
		result := newBorderProperty(nil)
		if style := properties.getRaw(tag + "-" + Style); style != nil {
			result.Set(Style, style)
		} else if style := properties.getRaw(Style); style != nil {
			result.Set(Style, style)
		}
		if width := properties.getRaw(tag + "-" + Width); width != nil {
			result.Set(Width, width)
		} else if width := properties.getRaw(Width); width != nil {
			result.Set(Width, width)
		}
		if color := properties.getRaw(tag + "-" + ColorTag); color != nil {
			result.Set(ColorTag, color)
		} else if color := properties.getRaw(ColorTag); color != nil {
			result.Set(ColorTag, color)
		}
		return result

	case LeftStyle, RightStyle, TopStyle, BottomStyle:
		if style := properties.getRaw(tag); style != nil {
			return style
		}
		return properties.getRaw(Style)

	case LeftWidth, RightWidth, TopWidth, BottomWidth:
		if width := properties.getRaw(tag); width != nil {
			return width
		}
		return properties.getRaw(Width)

	case LeftColor, RightColor, TopColor, BottomColor:
		if color := properties.getRaw(tag); color != nil {
			return color
		}
		return properties.getRaw(ColorTag)
	}

	return nil
}

func (border *borderProperty) deleteTag(tag PropertyName) bool {

	result := false
	removeTags := func(tags []PropertyName) {
		for _, tag := range tags {
			if border.getRaw(tag) != nil {
				border.setRaw(tag, nil)
				result = true
			}
		}
	}

	switch tag {
	case Style:
		removeTags([]PropertyName{Style, LeftStyle, RightStyle, TopStyle, BottomStyle})

	case Width:
		removeTags([]PropertyName{Width, LeftWidth, RightWidth, TopWidth, BottomWidth})

	case ColorTag:
		removeTags([]PropertyName{ColorTag, LeftColor, RightColor, TopColor, BottomColor})

	case Left, Right, Top, Bottom:
		if border.Get(Style) != nil {
			border.properties[tag+"-"+Style] = 0
			result = true
			removeTags([]PropertyName{tag + "-" + ColorTag, tag + "-" + Width})
		} else {
			removeTags([]PropertyName{tag + "-" + Style, tag + "-" + ColorTag, tag + "-" + Width})
		}

	case LeftStyle, RightStyle, TopStyle, BottomStyle:
		if border.getRaw(tag) != nil {
			if border.Get(Style) != nil {
				border.properties[tag] = 0
				result = true
			} else {
				removeTags([]PropertyName{tag})
			}
		}

	case LeftWidth, RightWidth, TopWidth, BottomWidth:
		if border.getRaw(tag) != nil {
			if border.Get(Width) != nil {
				border.properties[tag] = AutoSize()
				result = true
			} else {
				removeTags([]PropertyName{tag})
			}
		}

	case LeftColor, RightColor, TopColor, BottomColor:
		if border.getRaw(tag) != nil {
			if border.Get(ColorTag) != nil {
				border.properties[tag] = 0
				result = true
			} else {
				removeTags([]PropertyName{tag})
			}
		}
	}

	return result
}

func (border *borderProperty) ViewBorders(session Session) ViewBorders {

	defStyle, _ := valueToEnum(border.getRaw(Style), BorderStyle, session, NoneLine)
	defWidth, _ := sizeProperty(border, Width, session)
	defColor, _ := colorProperty(border, ColorTag, session)

	getBorder := func(prefix PropertyName) ViewBorder {
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
		builder.add(string(BorderStyle), values[borders.Top.Style])
	} else {
		builder.addValues(string(BorderStyle), " ", values[borders.Top.Style],
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

func getBorderProperty(properties Properties, tag PropertyName) BorderProperty {
	if value := properties.getRaw(tag); value != nil {
		if border, ok := value.(BorderProperty); ok {
			return border
		}
	}
	return nil
}

func setBorderPropertyElement(properties Properties, mainTag, tag PropertyName, value any) []PropertyName {
	border := getBorderProperty(properties, mainTag)
	if border == nil {
		border = NewBorder(nil)
		if border.Set(tag, value) {
			properties.setRaw(mainTag, border)
			return []PropertyName{mainTag, tag}
		}
	} else if border.Set(tag, value) {
		return []PropertyName{mainTag, tag}
	}
	return nil
}
