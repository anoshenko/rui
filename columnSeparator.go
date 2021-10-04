package rui

import (
	"fmt"
	"strings"
)

// ColumnSeparatorProperty is the interface of a view separator data
type ColumnSeparatorProperty interface {
	Properties
	ruiStringer
	fmt.Stringer
	ViewBorder(session Session) ViewBorder
	cssValue(session Session) string
}

type columnSeparatorProperty struct {
	propertyList
}

func newColumnSeparatorProperty(value interface{}) ColumnSeparatorProperty {

	if value == nil {
		separator := new(columnSeparatorProperty)
		separator.properties = map[string]interface{}{}
		return separator
	}

	switch value := value.(type) {
	case ColumnSeparatorProperty:
		return value

	case DataObject:
		separator := new(columnSeparatorProperty)
		separator.properties = map[string]interface{}{}
		for _, tag := range []string{Style, Width, ColorTag} {
			if val, ok := value.PropertyValue(tag); ok && val != "" {
				separator.set(tag, value)
			}
		}
		return separator

	case ViewBorder:
		separator := new(columnSeparatorProperty)
		separator.properties = map[string]interface{}{
			Style:    value.Style,
			Width:    value.Width,
			ColorTag: value.Color,
		}
		return separator
	}

	invalidPropertyValue(Border, value)
	return nil
}

// NewColumnSeparator creates the new ColumnSeparatorProperty
func NewColumnSeparator(params Params) ColumnSeparatorProperty {
	separator := new(columnSeparatorProperty)
	separator.properties = map[string]interface{}{}
	if params != nil {
		for _, tag := range []string{Style, Width, ColorTag} {
			if value, ok := params[tag]; ok && value != nil {
				separator.Set(tag, value)
			}
		}
	}
	return separator
}

func (separator *columnSeparatorProperty) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
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

func (separator *columnSeparatorProperty) ruiString(writer ruiWriter) {
	writer.startObject("_")
	for _, tag := range []string{Style, Width, ColorTag} {
		if value, ok := separator.properties[tag]; ok {
			writer.writeProperty(Style, value)
		}
	}
	writer.endObject()
}

func (separator *columnSeparatorProperty) String() string {
	writer := newRUIWriter()
	separator.ruiString(writer)
	return writer.finish()
}

func (separator *columnSeparatorProperty) Remove(tag string) {

	switch tag = separator.normalizeTag(tag); tag {
	case Style, Width, ColorTag:
		delete(separator.properties, tag)

	default:
		ErrorLogF(`"%s" property is not compatible with the ColumnSeparatorProperty`, tag)
	}
}

func (separator *columnSeparatorProperty) Set(tag string, value interface{}) bool {
	tag = separator.normalizeTag(tag)

	if value == nil {
		separator.remove(tag)
		return true
	}

	switch tag {
	case Style:
		return separator.setEnumProperty(Style, value, enumProperties[BorderStyle].values)

	case Width:
		return separator.setSizeProperty(Width, value)

	case ColorTag:
		return separator.setColorProperty(ColorTag, value)
	}

	ErrorLogF(`"%s" property is not compatible with the ColumnSeparatorProperty`, tag)
	return false
}

func (separator *columnSeparatorProperty) Get(tag string) interface{} {
	tag = separator.normalizeTag(tag)

	if result, ok := separator.properties[tag]; ok {
		return result
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

	if value.Width.Type != Auto && value.Width.Type != SizeInFraction && value.Width.Value > 0 {
		buffer.WriteString(value.Width.cssString(""))
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
