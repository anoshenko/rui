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

	// BorderBoxClip is value of the BackgroundClip property:
	// The background extends to the outside edge of the border (but underneath the border in z-ordering).
	BorderBoxClip = 0
	// PaddingBoxClip is value of the BackgroundClip property:
	// The background extends to the outside edge of the padding. No background is drawn beneath the border.
	PaddingBoxClip = 1
	// ContentBoxClip is value of the BackgroundClip property:
	// The background is painted within (clipped to) the content box.
	ContentBoxClip = 2
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
	propertyList
}

type backgroundImage struct {
	backgroundElement
}

// NewBackgroundImage creates the new background image
func createBackground(obj DataObject) BackgroundElement {
	var result BackgroundElement = nil

	switch obj.Tag() {
	case "image":
		image := new(backgroundImage)
		image.properties = map[string]any{}
		result = image

	case "linear-gradient":
		gradient := new(backgroundLinearGradient)
		gradient.properties = map[string]any{}
		result = gradient

	case "radial-gradient":
		gradient := new(backgroundRadialGradient)
		gradient.properties = map[string]any{}
		result = gradient

	case "conic-gradient":
		gradient := new(backgroundConicGradient)
		gradient.properties = map[string]any{}
		result = gradient

	default:
		return nil
	}

	count := obj.PropertyCount()
	for i := 0; i < count; i++ {
		if node := obj.Property(i); node.Type() == TextNode {
			if value := node.Text(); value != "" {
				result.Set(node.Tag(), value)
			}
		}
	}

	return result
}

// NewBackgroundImage creates the new background image
func NewBackgroundImage(params Params) BackgroundElement {
	result := new(backgroundImage)
	result.properties = map[string]any{}
	for tag, value := range params {
		result.Set(tag, value)
	}
	return result
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

func (image *backgroundImage) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
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

func (image *backgroundImage) Set(tag string, value any) bool {
	tag = image.normalizeTag(tag)
	switch tag {
	case Attachment, Width, Height, Repeat, ImageHorizontalAlign, ImageVerticalAlign,
		backgroundFit, Source:
		return image.backgroundElement.Set(tag, value)
	}

	return false
}

func (image *backgroundImage) Get(tag string) any {
	return image.backgroundElement.Get(image.normalizeTag(tag))
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
	image.writeToBuffer(buffer, indent, image.Tag(), []string{
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
