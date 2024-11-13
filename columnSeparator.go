package rui

import (
	"fmt"
	"strings"
)

// ColumnSeparatorProperty is the interface of a view separator data
type ColumnSeparatorProperty interface {
	Properties
	fmt.Stringer
	stringWriter

	// ViewBorder returns column separator description in a form of ViewBorder
	ViewBorder(session Session) ViewBorder

	cssValue(session Session) string
}

type columnSeparatorProperty struct {
	dataProperty
}

func newColumnSeparatorProperty(value any) ColumnSeparatorProperty {

	if value == nil {
		separator := new(columnSeparatorProperty)
		separator.init()
		return separator
	}

	switch value := value.(type) {
	case ColumnSeparatorProperty:
		return value

	case DataObject:
		separator := new(columnSeparatorProperty)
		separator.init()
		for _, tag := range []PropertyName{Style, Width, ColorTag} {
			if val, ok := value.PropertyValue(string(tag)); ok && val != "" {
				propertiesSet(separator, tag, value)
			}
		}
		return separator

	case ViewBorder:
		separator := new(columnSeparatorProperty)
		separator.init()
		separator.properties = map[PropertyName]any{
			Style:    value.Style,
			Width:    value.Width,
			ColorTag: value.Color,
		}
		return separator
	}

	invalidPropertyValue(Border, value)
	return nil
}

// NewColumnSeparator creates the new ColumnSeparatorProperty.
// The following properties can be used:
//
// "style" (Style). Determines the line style (int). Valid values: 0 (NoneLine), 1 (SolidLine), 2 (DashedLine), 3 (DottedLine), or 4 (DoubleLine);
//
// "color" (ColorTag). Determines the line color (Color);
//
// "width" (Width). Determines the line thickness (SizeUnit).
func NewColumnSeparator(params Params) ColumnSeparatorProperty {
	separator := new(columnSeparatorProperty)
	separator.init()
	if params != nil {
		for _, tag := range []PropertyName{Style, Width, ColorTag} {
			if value, ok := params[tag]; ok && value != nil {
				separator.Set(tag, value)
			}
		}
	}
	return separator
}

func (separator *columnSeparatorProperty) init() {
	separator.dataProperty.init()
	separator.normalize = normalizeVolumnSeparatorTag
	separator.supportedProperties = []PropertyName{Style, Width, ColorTag}
}

func normalizeVolumnSeparatorTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case ColumnSeparatorStyle, "separator-style":
		return Style

	case ColumnSeparatorWidth, "separator-width":
		return Width

	case ColumnSeparatorColor, "separator-color":
		return ColorTag
	}

	return tag
}

func (separator *columnSeparatorProperty) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("_{ ")
	comma := false
	for _, tag := range []PropertyName{Style, Width, ColorTag} {
		if value, ok := separator.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(string(tag))
			buffer.WriteString(" = ")
			writePropertyValue(buffer, BorderStyle, value, indent)
			comma = true
		}
	}

	buffer.WriteString(" }")
}

func (separator *columnSeparatorProperty) String() string {
	return runStringWriter(separator)
}

func getColumnSeparatorProperty(properties Properties) ColumnSeparatorProperty {
	if val := properties.getRaw(ColumnSeparator); val != nil {
		if separator, ok := val.(ColumnSeparatorProperty); ok {
			return separator
		}
	}
	return nil
}

func (separator *columnSeparatorProperty) ViewBorder(session Session) ViewBorder {
	style, _ := valueToEnum(separator.getRaw(Style), BorderStyle, session, NoneLine)
	width, _ := sizeProperty(separator, Width, session)
	color, _ := colorProperty(separator, ColorTag, session)

	return ViewBorder{
		Style: style,
		Width: width,
		Color: color,
	}
}

func (separator *columnSeparatorProperty) cssValue(session Session) string {
	value := separator.ViewBorder(session)
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if value.Width.Type != Auto && value.Width.Type != SizeInFraction &&
		(value.Width.Value > 0 || value.Width.Type == SizeFunction) {
		buffer.WriteString(value.Width.cssString("", session))
	}

	styles := enumProperties[BorderStyle].cssValues
	if value.Style > 0 && value.Style < len(styles) {
		if buffer.Len() > 0 {
			buffer.WriteRune(' ')
		}
		buffer.WriteString(styles[value.Style])
	}

	if value.Color != 0 {
		if buffer.Len() > 0 {
			buffer.WriteRune(' ')
		}
		buffer.WriteString(value.Color.cssString())
	}

	return buffer.String()
}
