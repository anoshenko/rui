package rui

import (
	"strings"
)

// Constants related to view's background description
const (
	// NoRepeat is value of the Repeat property of an background image:
	//
	// The image is not repeated (and hence the background image painting area
	// will not necessarily be entirely covered). The position of the non-repeated
	// background image is defined by the background-position CSS property.
	NoRepeat = 0

	// RepeatXY is value of the Repeat property of an background image:
	//
	// The image is repeated as much as needed to cover the whole background
	// image painting area. The last image will be clipped if it doesn't fit.
	RepeatXY = 1

	// RepeatX is value of the Repeat property of an background image:
	//
	// The image is repeated horizontally as much as needed to cover
	// the whole width background image painting area. The image is not repeated vertically.
	// The last image will be clipped if it doesn't fit.
	RepeatX = 2

	// RepeatY is value of the Repeat property of an background image:
	//
	// The image is repeated vertically as much as needed to cover
	// the whole height background image painting area. The image is not repeated horizontally.
	// The last image will be clipped if it doesn't fit.
	RepeatY = 3

	// RepeatRound is value of the Repeat property of an background image:
	//
	// As the allowed space increases in size, the repeated images will stretch (leaving no gaps)
	// until there is room (space left >= half of the image width) for another one to be added.
	// When the next image is added, all of the current ones compress to allow room.
	RepeatRound = 4

	// RepeatSpace is value of the Repeat property of an background image:
	//
	// The image is repeated as much as possible without clipping. The first and last images
	// are pinned to either side of the element, and whitespace is distributed evenly between the images.
	RepeatSpace = 5

	// ScrollAttachment is value of the Attachment property of an background image:
	//
	// The background is fixed relative to the element itself and does not scroll with its contents.
	// (It is effectively attached to the element's border.)
	ScrollAttachment = 0

	// FixedAttachment is value of the Attachment property of an background image:
	//
	// The background is fixed relative to the viewport. Even if an element has
	// a scrolling mechanism, the background doesn't move with the element.
	FixedAttachment = 1

	// LocalAttachment is value of the Attachment property of an background image:
	//
	// The background is fixed relative to the element's contents. If the element has a scrolling mechanism,
	// the background scrolls with the element's contents, and the background painting area
	// and background positioning area are relative to the scrollable area of the element
	// rather than to the border framing them.
	LocalAttachment = 2
)

type backgroundImage struct {
	backgroundElement
}

// NewBackgroundImage creates the new background image
//
// The following properties can be used:
//   - "src" [Source] - the name of the image in the "images" folder of the resources, or the URL of the image or inline-image.
//   - "width" [Width] - the width of the image.
//   - "height" [Height] - the height of the image.
//   - "image-horizontal-align" [ImageHorizontalAlign] - the horizontal alignment of the image relative to view's bounds.
//   - "image-vertical-align" [ImageVerticalAlign] - the vertical alignment of the image relative to view's  bounds.
//   - "repeat" [Repeat] - the repetition of the image.
//   - "fit" [Fit] - the image scaling parameters.
//   - "attachment" [Attachment] - defines whether a background image's position is fixed within the viewport or scrolls with its containing block.
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
