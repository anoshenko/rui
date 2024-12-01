package rui

import (
	"fmt"
)

const (
	// BorderBox is the value of the following properties:
	//
	// * BackgroundClip - The background extends to the outside edge of the border (but underneath the border in z-ordering).
	//
	// * BackgroundOrigin - The background is positioned relative to the border box.
	//
	// * MaskClip - The painted content is clipped to the border box.
	//
	// * MaskOrigin - The mask is positioned relative to the border box.
	BorderBox = 0

	// PaddingBox is value of the BackgroundClip and MaskClip property:
	//
	// * BackgroundClip - The background extends to the outside edge of the padding. No background is drawn beneath the border.
	//
	// * BackgroundOrigin - The background is positioned relative to the padding box.
	//
	// * MaskClip - The painted content is clipped to the padding box.
	//
	// * MaskOrigin - The mask is positioned relative to the padding box.
	PaddingBox = 1

	// ContentBox is value of the BackgroundClip and MaskClip property:
	//
	// * BackgroundClip - The background is painted within (clipped to) the content box.
	//
	// * BackgroundOrigin - The background is positioned relative to the content box.
	//
	// * MaskClip - The painted content is clipped to the content box.
	//
	// * MaskOrigin - The mask is positioned relative to the content box.
	ContentBox = 2
)

// BackgroundElement describes the background element
type BackgroundElement interface {
	Properties
	fmt.Stringer
	stringWriter
	cssStyle(session Session) string

	// Tag returns type of the background element.
	// Possible values are: "image", "conic-gradient", "linear-gradient" and "radial-gradient"
	Tag() string

	// Clone creates a new copy of BackgroundElement
	Clone() BackgroundElement
}

type backgroundElement struct {
	dataProperty
}

// NewBackgroundImage creates the new background image
func createBackground(obj DataObject) BackgroundElement {
	var result BackgroundElement = nil

	switch obj.Tag() {
	case "image":
		result = NewBackgroundImage(nil)

	case "linear-gradient":
		result = NewBackgroundLinearGradient(nil)

	case "radial-gradient":
		result = NewBackgroundRadialGradient(nil)

	case "conic-gradient":
		result = NewBackgroundConicGradient(nil)

	default:
		return nil
	}

	count := obj.PropertyCount()
	for i := 0; i < count; i++ {
		if node := obj.Property(i); node.Type() == TextNode {
			if value := node.Text(); value != "" {
				result.Set(PropertyName(node.Tag()), value)
			}
		}
	}

	return result
}

func parseBackgroundValue(value any) []BackgroundElement {

	switch value := value.(type) {
	case BackgroundElement:
		return []BackgroundElement{value}

	case []BackgroundElement:
		return value

	case []DataValue:
		background := []BackgroundElement{}
		for _, el := range value {
			if el.IsObject() {
				if element := createBackground(el.Object()); element != nil {
					background = append(background, element)
				} else {
					return nil
				}
			} else if obj := ParseDataText(el.Value()); obj != nil {
				if element := createBackground(obj); element != nil {
					background = append(background, element)
				} else {
					return nil
				}
			} else {
				return nil
			}
		}
		return background

	case DataObject:
		if element := createBackground(value); element != nil {
			return []BackgroundElement{element}
		}

	case []DataObject:
		background := []BackgroundElement{}
		for _, obj := range value {
			if element := createBackground(obj); element != nil {
				background = append(background, element)
			} else {
				return nil
			}
		}
		return background

	case string:
		if obj := ParseDataText(value); obj != nil {
			if element := createBackground(obj); element != nil {
				return []BackgroundElement{element}
			}
		}
	}

	return nil
}

func setBackgroundProperty(properties Properties, tag PropertyName, value any) []PropertyName {

	background := parseBackgroundValue(value)
	if background == nil {
		notCompatibleType(tag, value)
		return nil
	}

	if len(background) > 0 {
		properties.setRaw(tag, background)
	} else if properties.getRaw(tag) != nil {
		properties.setRaw(tag, nil)
	} else {
		return []PropertyName{}
	}

	return []PropertyName{tag}
}

func backgroundCSS(properties Properties, session Session) string {

	if value := properties.getRaw(Background); value != nil {
		if backgrounds, ok := value.([]BackgroundElement); ok && len(backgrounds) > 0 {
			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)

			for _, background := range backgrounds {
				if value := background.cssStyle(session); value != "" {
					if buffer.Len() > 0 {
						buffer.WriteString(", ")
					}
					buffer.WriteString(value)
				}
			}

			if buffer.Len() > 0 {
				backgroundColor, _ := colorProperty(properties, BackgroundColor, session)
				if backgroundColor != 0 {
					buffer.WriteRune(' ')
					buffer.WriteString(backgroundColor.cssString())
				}
				return buffer.String()
			}
		}
	}
	return ""
}

func maskCSS(properties Properties, session Session) string {

	if value := properties.getRaw(Mask); value != nil {
		if backgrounds, ok := value.([]BackgroundElement); ok && len(backgrounds) > 0 {
			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)

			for _, background := range backgrounds {
				if value := background.cssStyle(session); value != "" {
					if buffer.Len() > 0 {
						buffer.WriteString(", ")
					}
					buffer.WriteString(value)
				}
			}
			return buffer.String()
		}
	}
	return ""
}

func backgroundStyledPropery(view View, subviewID []string, tag PropertyName) []BackgroundElement {
	var background []BackgroundElement = nil

	if view = getSubview(view, subviewID); view != nil {
		if value := view.getRaw(tag); value != nil {
			if backgrounds, ok := value.([]BackgroundElement); ok {
				background = backgrounds
			}
		} else if value := valueFromStyle(view, tag); value != nil {
			background = parseBackgroundValue(value)
		}
	}

	if count := len(background); count > 0 {
		result := make([]BackgroundElement, count)
		copy(result, background)
		return result
	}

	return []BackgroundElement{}
}

// GetBackground returns the view background.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetBackground(view View, subviewID ...string) []BackgroundElement {
	return backgroundStyledPropery(view, subviewID, Background)
}

// GetMask returns the view mask.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetMask(view View, subviewID ...string) []BackgroundElement {
	return backgroundStyledPropery(view, subviewID, Mask)
}

// GetBackgroundClip returns a "background-clip" of the subview. Returns one of next values:
//
// BorderBox (0), PaddingBox (1), ContentBox (2)
//
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetBackgroundClip(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, BackgroundClip, 0, false)
}

// GetBackgroundOrigin returns a "background-origin" of the subview. Returns one of next values:
//
// BorderBox (0), PaddingBox (1), ContentBox (2)
//
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetBackgroundOrigin(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, BackgroundOrigin, 0, false)
}

// GetMaskClip returns a "mask-clip" of the subview. Returns one of next values:
//
// BorderBox (0), PaddingBox (1), ContentBox (2)
//
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetMaskClip(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, MaskClip, 0, false)
}

// GetMaskOrigin returns a "mask-origin" of the subview. Returns one of next values:
//
// BorderBox (0), PaddingBox (1), ContentBox (2)
//
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetMaskOrigin(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, MaskOrigin, 0, false)
}
