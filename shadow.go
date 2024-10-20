package rui

import (
	"fmt"
	"strings"
)

// Constants for [ViewShadow] specific properties
const (
	// ColorTag is the constant for "color" property tag.
	//
	// Used by `ColumnSeparatorProperty`, `BorderProperty`, `OutlineProperty`, `ViewShadow`.
	//
	// Usage in `ColumnSeparatorProperty`:
	// Line color.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	//
	// Usage in `BorderProperty`:
	// Border line color.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	//
	// Usage in `OutlineProperty`:
	// Outline line color.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	//
	// Usage in `ViewShadow`:
	// Color property of the shadow.
	//
	// Supported types: `Color`, `string`.
	//
	// Internal type is `Color`, other types converted to it during assignment.
	// See `Color` description for more details.
	ColorTag = "color"

	// Inset is the constant for "inset" property tag.
	//
	// Used by `ViewShadow`.
	// Controls whether to draw shadow inside the frame or outside. Inset shadows are drawn inside the border(even transparent
	// ones), above the background, but below content.
	//
	// Supported types: `bool`, `int`, `string`.
	//
	// Values:
	// `true` or `1` or "true", "yes", "on", "1" - Drop shadow inside the frame(as if the content was depressed inside the box).
	// `false` or `0` or "false", "no", "off", "0" - Shadow is assumed to be a drop shadow(as if the box were raised above the content).
	Inset = "inset"

	// XOffset is the constant for "x-offset" property tag.
	//
	// Used by `ViewShadow`.
	// Determines the shadow horizontal offset. Negative values place the shadow to the left of the element.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	XOffset = "x-offset"

	// YOffset is the constant for "y-offset" property tag.
	//
	// Used by `ViewShadow`.
	// Determines the shadow vertical offset. Negative values place the shadow above the element.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	YOffset = "y-offset"

	// BlurRadius is the constant for "blur" property tag.
	//
	// Used by `ViewShadow`.
	// Determines the radius of the blur effect. The larger this value, the bigger the blur, so the shadow becomes bigger and
	// lighter. Negative values are not allowed.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	BlurRadius = "blur"

	// SpreadRadius is the constant for "spread-radius" property tag.
	//
	// Used by `ViewShadow`.
	// Positive values will cause the shadow to expand and grow bigger, negative values will cause the shadow to shrink.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	SpreadRadius = "spread-radius"
)

// ViewShadow contains attributes of the view shadow
type ViewShadow interface {
	Properties
	fmt.Stringer
	stringWriter
	cssStyle(buffer *strings.Builder, session Session, lead string) bool
	cssTextStyle(buffer *strings.Builder, session Session, lead string) bool
	visible(session Session) bool
}

type viewShadowData struct {
	propertyList
}

// NewViewShadow create the new shadow for a view. Arguments:
//
// offsetX, offsetY is x and y offset of the shadow;
//
// blurRadius is the blur radius of the shadow;
//
// spreadRadius is the spread radius of the shadow;
//
// color is the color of the shadow.
func NewViewShadow(offsetX, offsetY, blurRadius, spreadRadius SizeUnit, color Color) ViewShadow {
	return NewShadowWithParams(Params{
		XOffset:      offsetX,
		YOffset:      offsetY,
		BlurRadius:   blurRadius,
		SpreadRadius: spreadRadius,
		ColorTag:     color,
	})
}

// NewInsetViewShadow create the new inset shadow for a view. Arguments:
//
// offsetX, offsetY is x and y offset of the shadow;
//
// blurRadius is the blur radius of the shadow;
//
// spreadRadius is the spread radius of the shadow;
//
// color is the color of the shadow.
func NewInsetViewShadow(offsetX, offsetY, blurRadius, spreadRadius SizeUnit, color Color) ViewShadow {
	return NewShadowWithParams(Params{
		XOffset:      offsetX,
		YOffset:      offsetY,
		BlurRadius:   blurRadius,
		SpreadRadius: spreadRadius,
		ColorTag:     color,
		Inset:        true,
	})
}

// NewTextShadow create the new text shadow. Arguments:
//
// offsetX, offsetY is the x- and y-offset of the shadow;
//
// blurRadius is the blur radius of the shadow;
//
// color is the color of the shadow.
func NewTextShadow(offsetX, offsetY, blurRadius SizeUnit, color Color) ViewShadow {
	return NewShadowWithParams(Params{
		XOffset:    offsetX,
		YOffset:    offsetY,
		BlurRadius: blurRadius,
		ColorTag:   color,
	})
}

// NewShadowWithParams create the new shadow for a view.
// The following properties can be used:
//
// "color" (ColorTag). Determines the color of the shadow (Color);
//
// "x-offset" (XOffset). Determines the shadow horizontal offset (SizeUnit);
//
// "y-offset" (YOffset). Determines the shadow vertical offset (SizeUnit);
//
// "blur" (BlurRadius). Determines the radius of the blur effect (SizeUnit);
//
// "spread-radius" (SpreadRadius). Positive values (SizeUnit) will cause the shadow to expand and grow bigger, negative values will cause the shadow to shrink;
//
// "inset" (Inset). Controls (bool) whether to draw shadow inside the frame or outside.
func NewShadowWithParams(params Params) ViewShadow {
	shadow := new(viewShadowData)
	shadow.propertyList.init()
	if params != nil {
		for _, tag := range []string{ColorTag, Inset, XOffset, YOffset, BlurRadius, SpreadRadius} {
			if value, ok := params[tag]; ok && value != nil {
				shadow.Set(tag, value)
			}
		}
	}
	return shadow
}

// parseViewShadow parse DataObject and create ViewShadow object
func parseViewShadow(object DataObject) ViewShadow {
	shadow := new(viewShadowData)
	shadow.propertyList.init()
	parseProperties(shadow, object)
	return shadow
}

func (shadow *viewShadowData) Remove(tag string) {
	delete(shadow.properties, strings.ToLower(tag))
}

func (shadow *viewShadowData) Set(tag string, value any) bool {
	if value == nil {
		shadow.Remove(tag)
		return true
	}

	tag = strings.ToLower(tag)
	switch tag {
	case ColorTag, Inset, XOffset, YOffset, BlurRadius, SpreadRadius:
		return shadow.propertyList.Set(tag, value)
	}

	ErrorLogF(`"%s" property is not supported by Shadow`, tag)
	return false
}

func (shadow *viewShadowData) Get(tag string) any {
	return shadow.propertyList.Get(strings.ToLower(tag))
}

func (shadow *viewShadowData) cssStyle(buffer *strings.Builder, session Session, lead string) bool {
	color, _ := colorProperty(shadow, ColorTag, session)
	offsetX, _ := sizeProperty(shadow, XOffset, session)
	offsetY, _ := sizeProperty(shadow, YOffset, session)
	blurRadius, _ := sizeProperty(shadow, BlurRadius, session)
	spreadRadius, _ := sizeProperty(shadow, SpreadRadius, session)

	if color.Alpha() == 0 ||
		((offsetX.Type == Auto || offsetX.Value == 0) &&
			(offsetY.Type == Auto || offsetY.Value == 0) &&
			(blurRadius.Type == Auto || blurRadius.Value == 0) &&
			(spreadRadius.Type == Auto || spreadRadius.Value == 0)) {
		return false
	}

	buffer.WriteString(lead)
	if inset, _ := boolProperty(shadow, Inset, session); inset {
		buffer.WriteString("inset ")
	}

	buffer.WriteString(offsetX.cssString("0", session))
	buffer.WriteByte(' ')
	buffer.WriteString(offsetY.cssString("0", session))
	buffer.WriteByte(' ')
	buffer.WriteString(blurRadius.cssString("0", session))
	buffer.WriteByte(' ')
	buffer.WriteString(spreadRadius.cssString("0", session))
	buffer.WriteByte(' ')
	buffer.WriteString(color.cssString())
	return true
}

func (shadow *viewShadowData) cssTextStyle(buffer *strings.Builder, session Session, lead string) bool {
	color, _ := colorProperty(shadow, ColorTag, session)
	offsetX, _ := sizeProperty(shadow, XOffset, session)
	offsetY, _ := sizeProperty(shadow, YOffset, session)
	blurRadius, _ := sizeProperty(shadow, BlurRadius, session)

	if color.Alpha() == 0 ||
		((offsetX.Type == Auto || offsetX.Value == 0) &&
			(offsetY.Type == Auto || offsetY.Value == 0) &&
			(blurRadius.Type == Auto || blurRadius.Value == 0)) {
		return false
	}

	buffer.WriteString(lead)
	buffer.WriteString(offsetX.cssString("0", session))
	buffer.WriteByte(' ')
	buffer.WriteString(offsetY.cssString("0", session))
	buffer.WriteByte(' ')
	buffer.WriteString(blurRadius.cssString("0", session))
	buffer.WriteByte(' ')
	buffer.WriteString(color.cssString())
	return true
}

func (shadow *viewShadowData) visible(session Session) bool {
	color, _ := colorProperty(shadow, ColorTag, session)
	offsetX, _ := sizeProperty(shadow, XOffset, session)
	offsetY, _ := sizeProperty(shadow, YOffset, session)
	blurRadius, _ := sizeProperty(shadow, BlurRadius, session)
	spreadRadius, _ := sizeProperty(shadow, SpreadRadius, session)

	if color.Alpha() == 0 ||
		((offsetX.Type == Auto || offsetX.Value == 0) &&
			(offsetY.Type == Auto || offsetY.Value == 0) &&
			(blurRadius.Type == Auto || blurRadius.Value == 0) &&
			(spreadRadius.Type == Auto || spreadRadius.Value == 0)) {
		return false
	}
	return true
}

func (shadow *viewShadowData) String() string {
	return runStringWriter(shadow)
}

func (shadow *viewShadowData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("_{ ")
	comma := false
	for _, tag := range shadow.AllTags() {
		if value, ok := shadow.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(tag)
			buffer.WriteString(" = ")
			writePropertyValue(buffer, tag, value, indent)
			comma = true
		}
	}
	buffer.WriteString(" }")
}

func (properties *propertyList) setShadow(tag string, value any) bool {

	if value == nil {
		delete(properties.properties, tag)
		return true
	}

	switch value := value.(type) {
	case ViewShadow:
		properties.properties[tag] = []ViewShadow{value}

	case []ViewShadow:
		if len(value) == 0 {
			delete(properties.properties, tag)
		} else {
			properties.properties[tag] = value
		}

	case DataValue:
		if !value.IsObject() {
			return false
		}
		properties.properties[tag] = []ViewShadow{parseViewShadow(value.Object())}

	case []DataValue:
		shadows := []ViewShadow{}
		for _, data := range value {
			if data.IsObject() {
				shadows = append(shadows, parseViewShadow(data.Object()))
			}
		}
		if len(shadows) == 0 {
			return false
		}
		properties.properties[tag] = shadows

	case string:
		obj := NewDataObject(value)
		if obj == nil {
			notCompatibleType(tag, value)
			return false
		}
		properties.properties[tag] = []ViewShadow{parseViewShadow(obj)}

	default:
		notCompatibleType(tag, value)
		return false
	}

	return true
}

func getShadows(properties Properties, tag string) []ViewShadow {
	if value := properties.Get(tag); value != nil {
		switch value := value.(type) {
		case []ViewShadow:
			return value

		case ViewShadow:
			return []ViewShadow{value}
		}
	}
	return []ViewShadow{}
}

func shadowCSS(properties Properties, tag string, session Session) string {
	shadows := getShadows(properties, tag)
	if len(shadows) == 0 {
		return ""
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	lead := ""
	if tag == Shadow {
		for _, shadow := range shadows {
			if shadow.cssStyle(buffer, session, lead) {
				lead = ", "
			}
		}
	} else {
		for _, shadow := range shadows {
			if shadow.cssTextStyle(buffer, session, lead) {
				lead = ", "
			}
		}
	}
	return buffer.String()
}
