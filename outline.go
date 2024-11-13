package rui

import (
	"fmt"
	"strings"
)

// OutlineProperty defines a view's outside border
type OutlineProperty interface {
	Properties
	stringWriter
	fmt.Stringer

	// ViewOutline returns style color and line width of an outline
	ViewOutline(session Session) ViewOutline
}

type outlinePropertyData struct {
	dataProperty
}

// NewOutlineProperty creates the new OutlineProperty.
// The following properties can be used:
//
// "color" (ColorTag). Determines the line color (Color);
//
// "width" (Width). Determines the line thickness (SizeUnit).
func NewOutlineProperty(params Params) OutlineProperty {
	outline := new(outlinePropertyData)
	outline.init()
	for tag, value := range params {
		outline.Set(tag, value)
	}
	return outline
}

func (outline *outlinePropertyData) init() {
	outline.propertyList.init()
	outline.normalize = normalizeOutlineTag
	outline.set = outlineSet
	outline.supportedProperties = []PropertyName{Style, Width, ColorTag}
}

func (outline *outlinePropertyData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("_{ ")
	comma := false
	for _, tag := range []PropertyName{Style, Width, ColorTag} {
		if value, ok := outline.properties[tag]; ok {
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

func (outline *outlinePropertyData) String() string {
	return runStringWriter(outline)
}

func normalizeOutlineTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	return PropertyName(strings.TrimPrefix(string(tag), "outline-"))
}

func outlineSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case Style:
		return setEnumProperty(properties, Style, value, enumProperties[BorderStyle].values)

	case Width:
		if width, ok := value.(SizeUnit); ok {
			switch width.Type {
			case SizeInFraction, SizeInPercent:
				notCompatibleType(tag, value)
				return nil
			}
		}
		return setSizeProperty(properties, Width, value)

	case ColorTag:
		return setColorProperty(properties, ColorTag, value)

	default:
		ErrorLogF(`"%s" property is not compatible with the OutlineProperty`, tag)
	}
	return nil
}

func (outline *outlinePropertyData) ViewOutline(session Session) ViewOutline {
	style, _ := valueToEnum(outline.getRaw(Style), BorderStyle, session, NoneLine)
	width, _ := sizeProperty(outline, Width, session)
	color, _ := colorProperty(outline, ColorTag, session)
	return ViewOutline{Style: style, Width: width, Color: color}
}

// ViewOutline describes parameters of a view border
type ViewOutline struct {
	// Style of the outline line
	Style int

	// Color of the outline line
	Color Color

	// Width of the outline line
	Width SizeUnit
}

func (outline ViewOutline) cssValue(builder cssBuilder, session Session) {
	values := enumProperties[BorderStyle].cssValues
	if outline.Style > 0 && outline.Style < len(values) && outline.Color.Alpha() > 0 &&
		outline.Width.Type != Auto && outline.Width.Type != SizeInFraction &&
		outline.Width.Type != SizeInPercent && outline.Width.Value > 0 {
		builder.addValues("outline", " ", outline.Width.cssString("0", session), values[outline.Style], outline.Color.cssString())
	}
}

func (outline ViewOutline) cssString(session Session) string {
	var builder cssValueBuilder
	outline.cssValue(&builder, session)
	return builder.finish()
}

func getOutlineProperty(properties Properties) OutlineProperty {
	if value := properties.Get(Outline); value != nil {
		if outline, ok := value.(OutlineProperty); ok {
			return outline
		}
	}

	return nil
}

func setOutlineProperty(properties Properties, value any) []PropertyName {
	switch value := value.(type) {
	case OutlineProperty:
		properties.setRaw(Outline, value)

	case ViewOutline:
		properties.setRaw(Outline, NewOutlineProperty(Params{Style: value.Style, Width: value.Width, ColorTag: value.Color}))

	case ViewBorder:
		properties.setRaw(Outline, NewOutlineProperty(Params{Style: value.Style, Width: value.Width, ColorTag: value.Color}))

	case DataObject:
		outline := NewOutlineProperty(nil)
		for _, tag := range []PropertyName{Style, Width, ColorTag} {
			if text, ok := value.PropertyValue(string(tag)); ok && text != "" {
				outline.Set(tag, text)
			}
		}
		properties.setRaw(Outline, outline)

	default:
		notCompatibleType(Outline, value)
		return nil
	}

	return []PropertyName{Outline}
}
