package rui

import (
	"fmt"
	"strings"
)

// Constants related to view's background description
const (
	// NoRepeat is value of the Repeat property of an background image:
	// The image is not repeated (and hence the background image painting area
	// will not necessarily be entirely covered). The position of the non-repeated
	// background image is defined by the background-position CSS property.
	NoRepeat = 0

	// RepeatXY is value of the Repeat property of an background image:
	// The image is repeated as much as needed to cover the whole background
	// image painting area. The last image will be clipped if it doesn't fit.
	RepeatXY = 1

	// RepeatX is value of the Repeat property of an background image:
	// The image is repeated horizontally as much as needed to cover
	// the whole width background image painting area. The image is not repeated vertically.
	// The last image will be clipped if it doesn't fit.
	RepeatX = 2

	// RepeatY is value of the Repeat property of an background image:
	// The image is repeated vertically as much as needed to cover
	// the whole height background image painting area. The image is not repeated horizontally.
	// The last image will be clipped if it doesn't fit.
	RepeatY = 3

	// RepeatRound is value of the Repeat property of an background image:
	// As the allowed space increases in size, the repeated images will stretch (leaving no gaps)
	// until there is room (space left >= half of the image width) for another one to be added.
	// When the next image is added, all of the current ones compress to allow room.
	RepeatRound = 4

	// RepeatSpace is value of the Repeat property of an background image:
	// The image is repeated as much as possible without clipping. The first and last images
	// are pinned to either side of the element, and whitespace is distributed evenly between the images.
	RepeatSpace = 5

	// ScrollAttachment is value of the Attachment property of an background image:
	// The background is fixed relative to the element itself and does not scroll with its contents.
	// (It is effectively attached to the element's border.)
	ScrollAttachment = 0

	// FixedAttachment is value of the Attachment property of an background image:
	// The background is fixed relative to the viewport. Even if an element has
	// a scrolling mechanism, the background doesn't move with the element.
	FixedAttachment = 1

	// LocalAttachment is value of the Attachment property of an background image:
	// The background is fixed relative to the element's contents. If the element has a scrolling mechanism,
	// the background scrolls with the element's contents, and the background painting area
	// and background positioning area are relative to the scrollable area of the element
	// rather than to the border framing them.
	LocalAttachment = 2

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

type backgroundImage struct {
	backgroundElement
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

// NewBackgroundImage creates the new background image
func NewBackgroundImage(params Params) BackgroundElement {
	result := new(backgroundImage)
	result.init()
	for tag, value := range params {
		result.Set(tag, value)
	}
	return result
}

func (image *backgroundImage) init() {
	image.backgroundElement.init()
	image.normalize = normalizeBackgroundImageTag
	image.supportedProperties = []PropertyName{
		Attachment, Width, Height, Repeat, ImageHorizontalAlign, ImageVerticalAlign, backgroundFit, Source,
	}
}

func (image *backgroundImage) Tag() string {
	return "image"
}

func (image *backgroundImage) Clone() BackgroundElement {
	result := NewBackgroundImage(nil)
	for tag, value := range image.properties {
		result.setRaw(tag, value)
	}
	return result
}

func normalizeBackgroundImageTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case "source":
		tag = Source

	case Fit:
		tag = backgroundFit

	case HorizontalAlign:
		tag = ImageHorizontalAlign

	case VerticalAlign:
		tag = ImageVerticalAlign
	}

	return tag
}

func (image *backgroundImage) cssStyle(session Session) string {
	if src, ok := imageProperty(image, Source, session); ok && src != "" {
		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)

		buffer.WriteString(`url(`)
		buffer.WriteString(src)
		buffer.WriteRune(')')

		attachment, _ := enumProperty(image, Attachment, session, NoRepeat)
		values := enumProperties[Attachment].values
		if attachment > 0 && attachment < len(values) {
			buffer.WriteRune(' ')
			buffer.WriteString(values[attachment])
		}

		align, _ := enumProperty(image, ImageHorizontalAlign, session, LeftAlign)
		values = enumProperties[ImageHorizontalAlign].values
		if align >= 0 && align < len(values) {
			buffer.WriteRune(' ')
			buffer.WriteString(values[align])
		} else {
			buffer.WriteString(` left`)
		}

		align, _ = enumProperty(image, ImageVerticalAlign, session, TopAlign)
		values = enumProperties[ImageVerticalAlign].values
		if align >= 0 && align < len(values) {
			buffer.WriteRune(' ')
			buffer.WriteString(values[align])
		} else {
			buffer.WriteString(` top`)
		}

		fit, _ := enumProperty(image, backgroundFit, session, NoneFit)
		values = enumProperties[backgroundFit].values
		if fit > 0 && fit < len(values) {

			buffer.WriteString(` / `)
			buffer.WriteString(values[fit])

		} else {

			width, _ := sizeProperty(image, Width, session)
			height, _ := sizeProperty(image, Height, session)

			if width.Type != Auto || height.Type != Auto {
				buffer.WriteString(` / `)
				buffer.WriteString(width.cssString("auto", session))
				buffer.WriteRune(' ')
				buffer.WriteString(height.cssString("auto", session))
			}
		}

		repeat, _ := enumProperty(image, Repeat, session, NoRepeat)
		values = enumProperties[Repeat].values
		if repeat >= 0 && repeat < len(values) {
			buffer.WriteRune(' ')
			buffer.WriteString(values[repeat])
		} else {
			buffer.WriteString(` no-repeat`)
		}

		return buffer.String()
	}

	return ""
}

func (image *backgroundImage) writeString(buffer *strings.Builder, indent string) {
	image.writeToBuffer(buffer, indent, image.Tag(), []PropertyName{
		Source,
		Width,
		Height,
		ImageHorizontalAlign,
		ImageVerticalAlign,
		backgroundFit,
		Repeat,
		Attachment,
	})
}

func (image *backgroundImage) String() string {
	return runStringWriter(image)
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
