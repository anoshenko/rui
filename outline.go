package rui

import (
	"fmt"
	"strings"
)

type OutlineProperty interface {
	Properties
	stringWriter
	fmt.Stringer
	ViewOutline(session Session) ViewOutline
}

type outlinePropertyData struct {
	propertyList
}

func NewOutlineProperty(params Params) OutlineProperty {
	outline := new(outlinePropertyData)
	outline.properties = map[string]any{}
	for tag, value := range params {
		outline.Set(tag, value)
	}
	return outline
}

func (outline *outlinePropertyData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("_{ ")
	comma := false
	for _, tag := range []string{Style, Width, ColorTag} {
		if value, ok := outline.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(tag)
			buffer.WriteString(" = ")
			writePropertyValue(buffer, BorderStyle, value, indent)
			comma = true
		}
	}

	buffer.WriteString(" }")
}

func (outline *outlinePropertyData) String() string {
	return runStringWriter(outline)
}

func (outline *outlinePropertyData) normalizeTag(tag string) string {
	return strings.TrimPrefix(strings.ToLower(tag), "outline-")
}

func (outline *outlinePropertyData) Remove(tag string) {
	delete(outline.properties, outline.normalizeTag(tag))
}

func (outline *outlinePropertyData) Set(tag string, value any) bool {
	if value == nil {
		outline.Remove(tag)
		return true
	}

	tag = outline.normalizeTag(tag)
	switch tag {
	case Style:
		return outline.setEnumProperty(Style, value, enumProperties[BorderStyle].values)

	case Width:
		if width, ok := value.(SizeUnit); ok {
			switch width.Type {
			case SizeInFraction, SizeInPercent:
				notCompatibleType(tag, value)
				return false
			}
		}
		return outline.setSizeProperty(Width, value)

	case ColorTag:
		return outline.setColorProperty(ColorTag, value)

	default:
		ErrorLogF(`"%s" property is not compatible with the OutlineProperty`, tag)
	}
	return false
}

func (outline *outlinePropertyData) Get(tag string) any {
	return outline.propertyList.Get(outline.normalizeTag(tag))
}

func (outline *outlinePropertyData) ViewOutline(session Session) ViewOutline {
	style, _ := valueToEnum(outline.getRaw(Style), BorderStyle, session, NoneLine)
	width, _ := sizeProperty(outline, Width, session)
	color, _ := colorProperty(outline, ColorTag, session)
	return ViewOutline{Style: style, Width: width, Color: color}
}

// ViewOutline describes parameters of a view border
type ViewOutline struct {
	Style int
	Color Color
	Width SizeUnit
}

func (outline ViewOutline) cssValue(builder cssBuilder) {
	values := enumProperties[BorderStyle].cssValues
	if outline.Style > 0 && outline.Style < len(values) && outline.Color.Alpha() > 0 &&
		outline.Width.Type != Auto && outline.Width.Type != SizeInFraction &&
		outline.Width.Type != SizeInPercent && outline.Width.Value > 0 {
		builder.addValues("outline", " ", outline.Width.cssString("0"), values[outline.Style], outline.Color.cssString())
	}
}

func (outline ViewOutline) cssString() string {
	var builder cssValueBuilder
	outline.cssValue(&builder)
	return builder.finish()
}

func getOutline(properties Properties) OutlineProperty {
	if value := properties.Get(Outline); value != nil {
		if outline, ok := value.(OutlineProperty); ok {
			return outline
		}
	}

	return nil
}

func (style *viewStyle) setOutline(value any) bool {
	switch value := value.(type) {
	case OutlineProperty:
		style.properties[Outline] = value

	case ViewOutline:
		style.properties[Outline] = NewOutlineProperty(Params{Style: value.Style, Width: value.Width, ColorTag: value.Color})

	case ViewBorder:
		style.properties[Outline] = NewOutlineProperty(Params{Style: value.Style, Width: value.Width, ColorTag: value.Color})

	case DataObject:
		outline := NewOutlineProperty(nil)
		for _, tag := range []string{Style, Width, ColorTag} {
			if text, ok := value.PropertyValue(tag); ok && text != "" {
				outline.Set(tag, text)
			}
		}
		style.properties[Outline] = outline

	default:
		notCompatibleType(Outline, value)
		return false
	}

	return true
}
