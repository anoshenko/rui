package rui

import (
	"fmt"
	"strings"
)

// Constants for [FilterProperty] specific properties and events
const (
	// Blur is the constant for "blur" property tag.
	//
	// Used by FilterProperty.
	// Applies a Gaussian blur. The value of radius defines the value of the standard deviation to the Gaussian function, or
	// how many pixels on the screen blend into each other, so a larger value will create more blur. The lacuna value for
	// interpolation is 0. The parameter is specified as a length in pixels.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Blur PropertyName = "blur"

	// Brightness is the constant for "brightness" property tag.
	//
	// Used by FilterProperty.
	// Applies a linear multiplier to input image, making it appear more or less bright. A value of 0% will create an image
	// that is completely black. A value of 100% leaves the input unchanged. Other values are linear multipliers on the
	// effect. Values of an amount over 100% are allowed, providing brighter results.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Brightness PropertyName = "brightness"

	// Contrast is the constant for "contrast" property tag.
	//
	// Used by FilterProperty.
	// Adjusts the contrast of the input. A value of 0% will create an image that is completely black. A value of 100% leaves
	// the input unchanged. Values of amount over 100% are allowed, providing results with less contrast.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Contrast PropertyName = "contrast"

	// DropShadow is the constant for "drop-shadow" property tag.
	//
	// Used by FilterProperty.
	// Applies a drop shadow effect to the input image. A drop shadow is effectively a blurred, offset version of the input
	// image's alpha mask drawn in a particular color, composited below the image. Shadow parameters are set using the
	// ShadowProperty interface.
	//
	// Supported types: []ShadowProperty, ShadowProperty, string.
	//
	// Internal type is []ShadowProperty, other types converted to it during assignment.
	// See ShadowProperty description for more details.
	//
	// Conversion rules:
	//   - []ShadowProperty - stored as is, no conversion performed.
	//   - ShadowProperty - converted to []ShadowProperty.
	//   - string - string representation of ShadowProperty. Example: "_{blur = 1em, color = black, spread-radius = 0.5em}".
	DropShadow PropertyName = "drop-shadow"

	// Grayscale is the constant for "grayscale" property tag.
	//
	// Used by FilterProperty.
	// Converts the input image to grayscale. The value of ‘amount’ defines the proportion of the conversion. A value of 100%
	// is completely grayscale. A value of 0% leaves the input unchanged. Values between 0% and 100% are linear multipliers on
	// the effect.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Grayscale PropertyName = "grayscale"

	// HueRotate is the constant for "hue-rotate" property tag.
	//
	// Used by FilterProperty.
	// Applies a hue rotation on the input image. The value of ‘angle’ defines the number of degrees around the color circle
	// the input samples will be adjusted. A value of 0deg leaves the input unchanged. If the ‘angle’ parameter is missing, a
	// value of 0deg is used. Though there is no maximum value, the effect of values above 360deg wraps around.
	//
	// Supported types: AngleUnit, string, float, int.
	//
	// Internal type is AngleUnit, other types will be converted to it during assignment.
	// See AngleUnit description for more details.
	//
	// Conversion rules:
	//   - AngleUnit - stored as is, no conversion performed.
	//   - string - must contain string representation of AngleUnit. If numeric value will be provided without any suffix then AngleUnit with value and Radian value type will be created.
	//   - float - a new AngleUnit value will be created with Radian as a type.
	//   - int - a new AngleUnit value will be created with Radian as a type.
	HueRotate PropertyName = "hue-rotate"

	// Invert is the constant for "invert" property tag.
	//
	// Used by FilterProperty.
	// Inverts the samples in the input image. The value of ‘amount’ defines the proportion of the conversion. A value of 100%
	// is completely inverted. A value of 0% leaves the input unchanged. Values between 0% and 100% are linear multipliers on
	// the effect.
	//
	// Supported types: float64, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Invert PropertyName = "invert"

	// Saturate is the constant for "saturate" property tag.
	//
	// Used by FilterProperty.
	// Saturates the input image. The value of ‘amount’ defines the proportion of the conversion. A value of 0% is completely
	// un-saturated. A value of 100% leaves the input unchanged. Other values are linear multipliers on the effect. Values of
	// amount over 100% are allowed, providing super-saturated results.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Saturate PropertyName = "saturate"

	// Sepia is the constant for "sepia" property tag.
	//
	// Used by FilterProperty.
	// Converts the input image to sepia. The value of ‘amount’ defines the proportion of the conversion. A value of 100% is
	// completely sepia. A value of 0% leaves the input unchanged. Values between 0% and 100% are linear multipliers on the
	// effect.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Sepia PropertyName = "sepia"
)

// FilterProperty defines an applied to a View a graphical effects like blur or color shift.
// Allowable properties are Blur, Brightness, Contrast, DropShadow, Grayscale, HueRotate, Invert, Opacity, Saturate, and Sepia
type FilterProperty interface {
	Properties
	fmt.Stringer
	stringWriter
	cssStyle(session Session) string
}

type filterData struct {
	dataProperty
}

// NewFilterProperty creates the new FilterProperty
func NewFilterProperty(params Params) FilterProperty {
	if len(params) > 0 {
		filter := new(filterData)
		filter.init()
		for tag, value := range params {
			if !filter.Set(tag, value) {
				return nil
			}
		}
		return filter
	}
	return nil
}

func newFilterProperty(obj DataObject) FilterProperty {
	filter := new(filterData)
	filter.init()
	for i := range obj.PropertyCount() {
		if node := obj.Property(i); node != nil {
			tag := node.Tag()
			switch node.Type() {
			case TextNode:
				filter.Set(PropertyName(tag), node.Text())

			case ObjectNode:
				if tag == string(HueRotate) {
					// TODO
				} else {
					ErrorLog(`Invalid value of "` + tag + `"`)
				}

			default:
				ErrorLog(`Invalid value of "` + tag + `"`)
			}
		}
	}

	if len(filter.properties) > 0 {
		return filter
	}
	ErrorLog("Empty view filter")
	return nil
}

func (filter *filterData) init() {
	filter.dataProperty.init()
	filter.set = filterDataSet
	filter.supportedProperties = []PropertyName{Blur, Brightness, Contrast, Saturate, Grayscale, Invert, Opacity, Sepia, HueRotate, DropShadow}
}

func filterDataSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case Blur, Brightness, Contrast, Saturate:
		return setFloatProperty(properties, tag, value, 0, 10000)

	case Grayscale, Invert, Opacity, Sepia:
		return setFloatProperty(properties, tag, value, 0, 100)

	case HueRotate:
		return setAngleProperty(properties, tag, value)

	case DropShadow:
		if setShadowProperty(properties, tag, value) {
			return []PropertyName{tag}
		}
	}

	ErrorLogF(`"%s" property is not supported by the view filter`, tag)
	return nil
}

func (filter *filterData) String() string {
	return runStringWriter(filter)
}

func (filter *filterData) writeString(buffer *strings.Builder, indent string) {
	filter.writeToBuffer(buffer, indent, "filter", filter.AllTags())
}

func (filter *filterData) cssStyle(session Session) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if value, ok := floatTextProperty(filter, Blur, session, 0); ok {
		buffer.WriteString(string(Blur))
		buffer.WriteRune('(')
		buffer.WriteString(value)
		buffer.WriteString("px)")
	}

	for _, tag := range []PropertyName{Brightness, Contrast, Saturate, Grayscale, Invert, Opacity, Sepia} {
		if value, ok := floatTextProperty(filter, tag, session, 0); ok {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(string(tag))
			buffer.WriteRune('(')
			buffer.WriteString(value)
			buffer.WriteString("%)")
		}
	}

	if value, ok := angleProperty(filter, HueRotate, session); ok {
		if buffer.Len() > 0 {
			buffer.WriteRune(' ')
		}
		buffer.WriteString(string(HueRotate))
		buffer.WriteRune('(')
		buffer.WriteString(value.cssString())
		buffer.WriteRune(')')
	}

	var lead string
	if buffer.Len() > 0 {
		lead = " drop-shadow("
	} else {
		lead = "drop-shadow("
	}

	for _, shadow := range getShadows(filter, DropShadow) {
		if shadow.cssTextStyle(buffer, session, lead) {
			buffer.WriteRune(')')
			lead = " drop-shadow("
		}
	}

	return buffer.String()
}

func setFilterProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	switch value := value.(type) {
	case FilterProperty:
		properties.setRaw(tag, value)
		return []PropertyName{tag}

	case string:
		if obj := NewDataObject(value); obj == nil {
			if filter := newFilterProperty(obj); filter != nil {
				properties.setRaw(tag, filter)
				return []PropertyName{tag}
			}
		}

	case DataObject:
		if filter := newFilterProperty(value); filter != nil {
			properties.setRaw(tag, filter)
			return []PropertyName{tag}
		}

	case DataValue:
		if value.IsObject() {
			if filter := newFilterProperty(value.Object()); filter != nil {
				properties.setRaw(tag, filter)
				return []PropertyName{tag}
			}
		}
	}

	notCompatibleType(tag, value)
	return nil
}

// GetFilter returns a View graphical effects like blur or color shift.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetFilter(view View, subviewID ...string) FilterProperty {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.getRaw(Filter); value != nil {
			if filter, ok := value.(FilterProperty); ok {
				return filter
			}
		}
		if value := valueFromStyle(view, Filter); value != nil {
			if filter, ok := value.(FilterProperty); ok {
				return filter
			}
		}
	}

	return nil
}

// GetBackdropFilter returns the area behind a View graphical effects like blur or color shift.
// If the second argument (subviewID) is not specified or it is "" then a top position of the first argument (view) is returned
func GetBackdropFilter(view View, subviewID ...string) FilterProperty {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.getRaw(BackdropFilter); value != nil {
			if filter, ok := value.(FilterProperty); ok {
				return filter
			}
		}
		if value := valueFromStyle(view, BackdropFilter); value != nil {
			if filter, ok := value.(FilterProperty); ok {
				return filter
			}
		}
	}

	return nil
}
