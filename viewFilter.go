package rui

import (
	"fmt"
	"strings"
)

const (
	// Blur is the constant for the "blur" property tag of the ViewFilter interface.
	// The "blur" float64 property applies a Gaussian blur. The value of radius defines the value
	// of the standard deviation to the Gaussian function, or how many pixels on the screen blend
	// into each other, so a larger value will create more blur. The lacuna value for interpolation is 0.
	// The parameter is specified as a length in pixels.
	Blur = "blur"

	// Brightness is the constant for the "brightness" property tag of the ViewFilter interface.
	// The "brightness" float64 property applies a linear multiplier to input image, making it appear more
	// or less bright. A value of 0% will create an image that is completely black.
	// A value of 100% leaves the input unchanged. Other values are linear multipliers on the effect.
	// Values of an amount over 100% are allowed, providing brighter results.
	Brightness = "brightness"

	// Contrast is the constant for the "contrast" property tag of the ViewFilter interface.
	// The "contrast" float64 property adjusts the contrast of the input.
	// A value of 0% will create an image that is completely black. A value of 100% leaves the input unchanged.
	// Values of amount over 100% are allowed, providing results with less contrast.
	Contrast = "contrast"

	// DropShadow is the constant for the "drop-shadow" property tag of the ViewFilter interface.
	// The "drop-shadow" property applies a drop shadow effect to the input image.
	// A drop shadow is effectively a blurred, offset version of the input image's alpha mask
	// drawn in a particular color, composited below the image.
	// Shadow parameters are set using the ViewShadow interface
	DropShadow = "drop-shadow"

	// Grayscale is the constant for the "grayscale" property tag of the ViewFilter interface.
	// The "grayscale" float64 property converts the input image to grayscale.
	// The value of ‘amount’ defines the proportion of the conversion.
	// A value of 100% is completely grayscale. A value of 0% leaves the input unchanged.
	// Values between 0% and 100% are linear multipliers on the effect.
	Grayscale = "grayscale"

	// HueRotate is the constant for the "hue-rotate" property tag of the ViewFilter interface.
	// The "hue-rotate" AngleUnit property applies a hue rotation on the input image.
	// The value of ‘angle’ defines the number of degrees around the color circle the input samples will be adjusted.
	// A value of 0deg leaves the input unchanged. If the ‘angle’ parameter is missing, a value of 0deg is used.
	// Though there is no maximum value, the effect of values above 360deg wraps around.
	HueRotate = "hue-rotate"

	// Invert is the constant for the "invert" property tag of the ViewFilter interface.
	// The "invert" float64 property inverts the samples in the input image.
	// The value of ‘amount’ defines the proportion of the conversion.
	// A value of 100% is completely inverted. A value of 0% leaves the input unchanged.
	// Values between 0% and 100% are linear multipliers on the effect.
	Invert = "invert"

	// Saturate is the constant for the "saturate" property tag of the ViewFilter interface.
	// The "saturate" float64 property saturates the input image.
	// The value of ‘amount’ defines the proportion of the conversion.
	// A value of 0% is completely un-saturated. A value of 100% leaves the input unchanged.
	// Other values are linear multipliers on the effect.
	// Values of amount over 100% are allowed, providing super-saturated results.
	Saturate = "saturate"

	// Sepia is the constant for the "sepia" property tag of the ViewFilter interface.
	// The "sepia" float64 property converts the input image to sepia.
	// The value of ‘amount’ defines the proportion of the conversion.
	// A value of 100% is completely sepia. A value of 0% leaves the input unchanged.
	// Values between 0% and 100% are linear multipliers on the effect.
	Sepia = "sepia"

	//Opacity = "opacity"
)

// ViewFilter defines an applied to a View a graphical effects like blur or color shift.
// Allowable properties are Blur, Brightness, Contrast, DropShadow, Grayscale, HueRotate, Invert, Opacity, Saturate, and Sepia
type ViewFilter interface {
	Properties
	fmt.Stringer
	stringWriter
	cssStyle(session Session) string
}

type viewFilter struct {
	propertyList
}

// NewViewFilter creates the new ViewFilter
func NewViewFilter(params Params) ViewFilter {
	if params != nil {
		filter := new(viewFilter)
		filter.init()
		for tag, value := range params {
			filter.Set(tag, value)
		}
		if len(filter.properties) > 0 {
			return filter
		}
	}
	return nil
}

func newViewFilter(obj DataObject) ViewFilter {
	filter := new(viewFilter)
	filter.init()
	for i := 0; i < obj.PropertyCount(); i++ {
		if node := obj.Property(i); node != nil {
			tag := node.Tag()
			switch node.Type() {
			case TextNode:
				filter.Set(tag, node.Text())

			case ObjectNode:
				if tag == HueRotate {
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

func (filter *viewFilter) Set(tag string, value any) bool {
	if value == nil {
		filter.Remove(tag)
		return true
	}

	switch strings.ToLower(tag) {
	case Blur, Brightness, Contrast, Saturate:
		return filter.setFloatProperty(tag, value, 0, 10000)

	case Grayscale, Invert, Opacity, Sepia:
		return filter.setFloatProperty(tag, value, 0, 100)

	case HueRotate:
		return filter.setAngleProperty(tag, value)

	case DropShadow:
		return filter.setShadow(tag, value)
	}

	ErrorLogF(`"%s" property is not supported by the view filter`, tag)
	return false
}

func (filter *viewFilter) String() string {
	return runStringWriter(filter)
}

func (filter *viewFilter) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("filter { ")
	comma := false
	tags := filter.AllTags()
	for _, tag := range tags {
		if value, ok := filter.properties[tag]; ok {
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

func (filter *viewFilter) cssStyle(session Session) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if value, ok := floatTextProperty(filter, Blur, session, 0); ok {
		buffer.WriteString(Blur)
		buffer.WriteRune('(')
		buffer.WriteString(value)
		buffer.WriteString("px)")
	}

	for _, tag := range []string{Brightness, Contrast, Saturate, Grayscale, Invert, Opacity, Sepia} {
		if value, ok := floatTextProperty(filter, tag, session, 0); ok {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(tag)
			buffer.WriteRune('(')
			buffer.WriteString(value)
			buffer.WriteString("%)")
		}
	}

	if value, ok := angleProperty(filter, HueRotate, session); ok {
		if buffer.Len() > 0 {
			buffer.WriteRune(' ')
		}
		buffer.WriteString(HueRotate)
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

func (style *viewStyle) setFilter(tag string, value any) bool {
	switch value := value.(type) {
	case ViewFilter:
		style.properties[tag] = value
		return true

	case string:
		if obj := NewDataObject(value); obj == nil {
			if filter := newViewFilter(obj); filter != nil {
				style.properties[tag] = filter
				return true
			}
		}
	case DataObject:
		if filter := newViewFilter(value); filter != nil {
			style.properties[tag] = filter
			return true
		}

	case DataValue:
		if value.IsObject() {
			if filter := newViewFilter(value.Object()); filter != nil {
				style.properties[tag] = filter
				return true
			}
		}
	}

	notCompatibleType(tag, value)
	return false
}

// GetFilter returns a View graphical effects like blur or color shift.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetFilter(view View, subviewID string) ViewFilter {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.getRaw(Filter); value != nil {
			if filter, ok := value.(ViewFilter); ok {
				return filter
			}
		}
	}

	return nil
}

// GetBackdropFilter returns the area behind a View graphical effects like blur or color shift.
// If the second argument (subviewID) is "" then a top position of the first argument (view) is returned
func GetBackdropFilter(view View, subviewID string) ViewFilter {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}
	if view != nil {
		if value := view.getRaw(BackdropFilter); value != nil {
			if filter, ok := value.(ViewFilter); ok {
				return filter
			}
		}
	}

	return nil
}
