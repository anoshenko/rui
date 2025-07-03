package rui

import (
	"fmt"
	"strings"
)

// Constants for [ShadowProperty] specific properties
const (
	// ColorTag is the constant for "color" property tag.
	//
	// Used by ColumnSeparatorProperty, BorderProperty, OutlineProperty, ShadowProperty.
	//
	// # Usage in ColumnSeparatorProperty
	//
	// Line color.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See Color description for more details.
	//
	// # Usage in BorderProperty
	//
	// Border line color.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See Color description for more details.
	//
	// # Usage in OutlineProperty
	//
	// Outline line color.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See Color description for more details.
	//
	// # Usage in ShadowProperty
	//
	// Color property of the shadow.
	//
	// Supported types: Color, string.
	//
	// Internal type is Color, other types converted to it during assignment.
	// See Color description for more details.
	ColorTag PropertyName = "color"

	// Inset is the constant for "inset" property tag.
	//
	// Used by ShadowProperty.
	// Controls whether to draw shadow inside the frame or outside. Inset shadows are drawn inside the border(even transparent
	// ones), above the background, but below content.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", "1" - Drop shadow inside the frame(as if the content was depressed inside the box).
	//   - false, 0, "false", "no", "off", "0" - Shadow is assumed to be a drop shadow(as if the box were raised above the content).
	Inset PropertyName = "inset"

	// XOffset is the constant for "x-offset" property tag.
	//
	// Used by ShadowProperty.
	// Determines the shadow horizontal offset. Negative values place the shadow to the left of the element.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	XOffset PropertyName = "x-offset"

	// YOffset is the constant for "y-offset" property tag.
	//
	// Used by ShadowProperty.
	// Determines the shadow vertical offset. Negative values place the shadow above the element.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	YOffset PropertyName = "y-offset"

	// BlurRadius is the constant for "blur" property tag.
	//
	// Used by ShadowProperty.
	// Determines the radius of the blur effect. The larger this value, the bigger the blur, so the shadow becomes bigger and
	// lighter. Negative values are not allowed.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	BlurRadius PropertyName = "blur"

	// SpreadRadius is the constant for "spread-radius" property tag.
	//
	// Used by ShadowProperty.
	// Positive values will cause the shadow to expand and grow bigger, negative values will cause the shadow to shrink.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	SpreadRadius PropertyName = "spread-radius"
)

// ShadowProperty contains attributes of the view shadow
type ShadowProperty interface {
	Properties
	fmt.Stringer
	stringWriter
	cssStyle(buffer *strings.Builder, session Session, lead string) bool
	cssTextStyle(buffer *strings.Builder, session Session, lead string) bool
	visible(session Session) bool
}

type shadowPropertyData struct {
	dataProperty
}

// NewShadow create the new shadow property for a view. Arguments:
//   - offsetX, offsetY is x and y offset of the shadow (if the argument is specified as int or float64, the value is considered to be in pixels);
//   - blurRadius is the blur radius of the shadow (if the argument is specified as int or float64, the value is considered to be in pixels);
//   - spreadRadius is the spread radius of the shadow (if the argument is specified as int or float64, the value is considered to be in pixels);
//   - color is the color of the shadow.
func NewShadow[xOffsetType SizeUnit | int | float64, yOffsetType SizeUnit | int | float64, blurType SizeUnit | int | float64, spreadType SizeUnit | int | float64](
	xOffset xOffsetType, yOffset yOffsetType, blurRadius blurType, spreadRadius spreadType, color Color) ShadowProperty {
	return NewShadowProperty(Params{
		XOffset:      xOffset,
		YOffset:      yOffset,
		BlurRadius:   blurRadius,
		SpreadRadius: spreadRadius,
		ColorTag:     color,
	})
}

// NewInsetShadow create the new inset shadow property for a view. Arguments:
//   - offsetX, offsetY is x and y offset of the shadow (if the argument is specified as int or float64, the value is considered to be in pixels);
//   - blurRadius is the blur radius of the shadow (if the argument is specified as int or float64, the value is considered to be in pixels);
//   - spreadRadius is the spread radius of the shadow (if the argument is specified as int or float64, the value is considered to be in pixels);
//   - color is the color of the shadow.
func NewInsetShadow[xOffsetType SizeUnit | int | float64, yOffsetType SizeUnit | int | float64, blurType SizeUnit | int | float64, spreadType SizeUnit | int | float64](
	xOffset xOffsetType, yOffset yOffsetType, blurRadius blurType, spreadRadius spreadType, color Color) ShadowProperty {
	return NewShadowProperty(Params{
		XOffset:      xOffset,
		YOffset:      yOffset,
		BlurRadius:   blurRadius,
		SpreadRadius: spreadRadius,
		ColorTag:     color,
		Inset:        true,
	})
}

// NewTextShadow create the new text shadow property. Arguments:
//   - offsetX, offsetY is the x- and y-offset of the shadow (if the argument is specified as int or float64, the value is considered to be in pixels);
//   - blurRadius is the blur radius of the shadow (if the argument is specified as int or float64, the value is considered to be in pixels);
//   - color is the color of the shadow.
func NewTextShadow[xOffsetType SizeUnit | int | float64, yOffsetType SizeUnit | int | float64, blurType SizeUnit | int | float64](
	xOffset xOffsetType, yOffset yOffsetType, blurRadius blurType, color Color) ShadowProperty {
	return NewShadowProperty(Params{
		XOffset:    xOffset,
		YOffset:    yOffset,
		BlurRadius: blurRadius,
		ColorTag:   color,
	})
}

// NewShadowProperty create the new shadow property for a view.
//
// The following properties can be used:
//   - "color" (ColorTag). Determines the color of the shadow (Color);
//   - "x-offset" (XOffset). Determines the shadow horizontal offset (SizeUnit);
//   - "y-offset" (YOffset). Determines the shadow vertical offset (SizeUnit);
//   - "blur" (BlurRadius). Determines the radius of the blur effect (SizeUnit);
//   - "spread-radius" (SpreadRadius). Positive values (SizeUnit) will cause the shadow to expand and grow bigger, negative values will cause the shadow to shrink;
//   - "inset" (Inset). Controls (bool) whether to draw shadow inside the frame or outside.
func NewShadowProperty(params Params) ShadowProperty {
	shadow := new(shadowPropertyData)
	shadow.init()

	if params != nil {
		for _, tag := range []PropertyName{ColorTag, Inset, XOffset, YOffset, BlurRadius, SpreadRadius} {
			if value, ok := params[tag]; ok && value != nil {
				shadow.set(shadow, tag, value)
			}
		}
	}
	return shadow
}

// parseShadowProperty parse DataObject and create ShadowProperty object
func parseShadowProperty(object DataObject) ShadowProperty {
	shadow := new(shadowPropertyData)
	shadow.init()
	parseProperties(shadow, object)
	return shadow
}

func (shadow *shadowPropertyData) init() {
	shadow.dataProperty.init()
	shadow.supportedProperties = []PropertyName{ColorTag, Inset, XOffset, YOffset, BlurRadius, SpreadRadius}
}

func (shadow *shadowPropertyData) cssStyle(buffer *strings.Builder, session Session, lead string) bool {
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

func (shadow *shadowPropertyData) cssTextStyle(buffer *strings.Builder, session Session, lead string) bool {
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

func (shadow *shadowPropertyData) visible(session Session) bool {
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

func (shadow *shadowPropertyData) String() string {
	return runStringWriter(shadow)
}

func setShadowProperty(properties Properties, tag PropertyName, value any) bool {

	if value == nil {
		properties.setRaw(tag, nil)
		return true
	}

	switch value := value.(type) {
	case ShadowProperty:
		properties.setRaw(tag, []ShadowProperty{value})

	case []ShadowProperty:
		if len(value) == 0 {
			properties.setRaw(tag, nil)
		} else {
			properties.setRaw(tag, value)
		}

	case DataValue:
		if !value.IsObject() {
			return false
		}
		properties.setRaw(tag, []ShadowProperty{parseShadowProperty(value.Object())})

	case []DataValue:
		shadows := []ShadowProperty{}
		for _, data := range value {
			if data.IsObject() {
				shadows = append(shadows, parseShadowProperty(data.Object()))
			}
		}
		if len(shadows) == 0 {
			return false
		}
		properties.setRaw(tag, shadows)

	case string:
		obj := NewDataObject(value)
		if obj == nil {
			notCompatibleType(tag, value)
			return false
		}
		properties.setRaw(tag, []ShadowProperty{parseShadowProperty(obj)})

	default:
		notCompatibleType(tag, value)
		return false
	}

	return true
}

func getShadows(properties Properties, tag PropertyName) []ShadowProperty {
	if value := properties.Get(tag); value != nil {
		switch value := value.(type) {
		case []ShadowProperty:
			return value

		case ShadowProperty:
			return []ShadowProperty{value}
		}
	}
	return []ShadowProperty{}
}

func shadowCSS(properties Properties, tag PropertyName, session Session) string {
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
