package rui

import (
	"fmt"
	"strings"
)

// Constants which represent [ImageView] specific properties and events
const (
	// LoadedEvent is the constant for "loaded-event" property tag.
	//
	// Used by `ImageView`.
	// Occur when the image has been loaded.
	//
	// General listener format:
	// `func(image rui.ImageView)`.
	//
	// where:
	// image - Interface of an image view which generated this event.
	//
	// Allowed listener formats:
	// `func()`.
	LoadedEvent = "loaded-event"

	// ErrorEvent is the constant for "error-event" property tag.
	//
	// Used by `ImageView`.
	// Occur when the image loading has been failed.
	//
	// General listener format:
	// `func(image rui.ImageView)`.
	//
	// where:
	// image - Interface of an image view which generated this event.
	//
	// Allowed listener formats:
	// `func()`.
	ErrorEvent = "error-event"

	// NoneFit - value of the "object-fit" property of an ImageView. The replaced content is not resized
	NoneFit = 0

	// ContainFit - value of the "object-fit" property of an ImageView. The replaced content
	// is scaled to maintain its aspect ratio while fitting within the element’s content box.
	// The entire object is made to fill the box, while preserving its aspect ratio, so the object
	// will be "letterboxed" if its aspect ratio does not match the aspect ratio of the box.
	ContainFit = 1

	// CoverFit - value of the "object-fit" property of an ImageView. The replaced content
	// is sized to maintain its aspect ratio while filling the element’s entire content box.
	// If the object's aspect ratio does not match the aspect ratio of its box, then the object will be clipped to fit.
	CoverFit = 2

	// FillFit - value of the "object-fit" property of an ImageView. The replaced content is sized
	// to fill the element’s content box. The entire object will completely fill the box.
	// If the object's aspect ratio does not match the aspect ratio of its box, then the object will be stretched to fit.
	FillFit = 3

	// ScaleDownFit - value of the "object-fit" property of an ImageView. The content is sized as
	// if NoneFit or ContainFit were specified, whichever would result in a smaller concrete object size.
	ScaleDownFit = 4
)

// ImageView represents an ImageView view
type ImageView interface {
	View
	// NaturalSize returns the intrinsic, density-corrected size (width, height) of the image in pixels.
	// If the image hasn't been loaded yet or an load error has occurred, then (0, 0) is returned.
	NaturalSize() (float64, float64)
	// CurrentSource() return the full URL of the image currently visible in the ImageView.
	// If the image hasn't been loaded yet or an load error has occurred, then "" is returned.
	CurrentSource() string
}

type imageViewData struct {
	viewData
	naturalWidth  float64
	naturalHeight float64
	currentSrc    string
}

// NewImageView create new ImageView object and return it
func NewImageView(session Session, params Params) ImageView {
	view := new(imageViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newImageView(session Session) View {
	return NewImageView(session, nil)
}

// Init initialize fields of imageView by default values
func (imageView *imageViewData) init(session Session) {
	imageView.viewData.init(session)
	imageView.tag = "ImageView"
	imageView.systemClass = "ruiImageView"
}

func (imageView *imageViewData) String() string {
	return getViewString(imageView, nil)
}

func (imageView *imageViewData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case "source":
		tag = Source

	case "src-set", "source-set":
		tag = SrcSet

	case VerticalAlign:
		tag = ImageVerticalAlign

	case HorizontalAlign:
		tag = ImageHorizontalAlign

	case altTag:
		tag = AltText
	}
	return tag
}

func (imageView *imageViewData) Remove(tag string) {
	imageView.remove(imageView.normalizeTag(tag))
}

func (imageView *imageViewData) remove(tag string) {
	imageView.viewData.remove(tag)
	if imageView.created {
		switch tag {
		case Source:
			imageView.session.updateProperty(imageView.htmlID(), "src", "")
			imageView.session.removeProperty(imageView.htmlID(), "srcset")

		case AltText:
			updateInnerHTML(imageView.htmlID(), imageView.session)

		case ImageVerticalAlign, ImageHorizontalAlign:
			updateCSSStyle(imageView.htmlID(), imageView.session)
		}
	}
}

func (imageView *imageViewData) Set(tag string, value any) bool {
	return imageView.set(imageView.normalizeTag(tag), value)
}

func (imageView *imageViewData) set(tag string, value any) bool {
	if value == nil {
		imageView.remove(tag)
		return true
	}

	switch tag {
	case Source:
		if text, ok := value.(string); ok {
			imageView.properties[tag] = text
			if imageView.created {
				src, srcset := imageView.src(text)
				imageView.session.updateProperty(imageView.htmlID(), "src", src)

				if srcset != "" {
					imageView.session.updateProperty(imageView.htmlID(), "srcset", srcset)
				} else {
					imageView.session.removeProperty(imageView.htmlID(), "srcset")
				}
			}
			imageView.propertyChangedEvent(Source)
			return true
		}
		notCompatibleType(Source, value)

	case SrcSet:
		if text, ok := value.(string); ok {
			if text == "" {
				delete(imageView.properties, tag)
			} else {
				imageView.properties[tag] = text
			}
			if imageView.created {
				_, srcset := imageView.src(text)
				if srcset != "" {
					imageView.session.updateProperty(imageView.htmlID(), "srcset", srcset)
				} else {
					imageView.session.removeProperty(imageView.htmlID(), "srcset")
				}
			}
			imageView.propertyChangedEvent(Source)
			return true
		}
		notCompatibleType(Source, value)

	case AltText:
		if text, ok := value.(string); ok {
			imageView.properties[AltText] = text
			if imageView.created {
				updateInnerHTML(imageView.htmlID(), imageView.session)
			}
			imageView.propertyChangedEvent(Source)
			return true
		}
		notCompatibleType(tag, value)

	case LoadedEvent, ErrorEvent:
		if listeners, ok := valueToNoParamListeners[ImageView](value); ok {
			if listeners == nil {
				delete(imageView.properties, tag)
			} else {
				imageView.properties[tag] = listeners
			}
			return true
		}

	default:
		if imageView.viewData.set(tag, value) {
			if imageView.created {
				switch tag {
				case ImageVerticalAlign, ImageHorizontalAlign:
					updateCSSStyle(imageView.htmlID(), imageView.session)
				}
			}
			return true
		}
	}

	return false
}

func (imageView *imageViewData) Get(tag string) any {
	return imageView.viewData.get(imageView.normalizeTag(tag))
}

func (imageView *imageViewData) imageListeners(tag string) []func(ImageView) {
	if value := imageView.getRaw(tag); value != nil {
		if listeners, ok := value.([]func(ImageView)); ok {
			return listeners
		}
	}
	return []func(ImageView){}
}

func (imageView *imageViewData) srcSet(path string) string {
	if value := imageView.getRaw(SrcSet); value != nil {
		if text, ok := value.(string); ok {
			srcset := strings.Split(text, ",")
			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)
			for i, src := range srcset {
				if i > 0 {
					buffer.WriteString(", ")
				}
				src = strings.Trim(src, " \t\n")
				buffer.WriteString(src)
				if index := strings.LastIndex(src, "@"); index > 0 {
					if ext := strings.LastIndex(src, "."); ext > index {
						buffer.WriteRune(' ')
						buffer.WriteString(src[index+1 : ext])
					}
				} else {
					buffer.WriteString(" 1x")
				}
			}
			return buffer.String()
		}
	}

	if srcset, ok := resources.imageSrcSets[path]; ok {
		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)
		for i, src := range srcset {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(src.path)
			buffer.WriteString(fmt.Sprintf(" %gx", src.scale))
		}
		return buffer.String()
	}
	return ""
}

func (imageView *imageViewData) htmlTag() string {
	return "img"
}

func (imageView *imageViewData) src(src string) (string, string) {
	if src != "" && src[0] == '@' {
		if image, ok := imageView.Session().ImageConstant(src[1:]); ok {
			src = image
		} else {
			src = ""
		}
	}

	if src != "" {
		return src, imageView.srcSet(src)
	}
	return "", ""
}

func (imageView *imageViewData) htmlProperties(self View, buffer *strings.Builder) {

	imageView.viewData.htmlProperties(self, buffer)

	if imageResource, ok := imageProperty(imageView, Source, imageView.Session()); ok && imageResource != "" {
		if src, srcset := imageView.src(imageResource); src != "" {
			buffer.WriteString(` src="`)
			buffer.WriteString(src)
			buffer.WriteString(`"`)
			if srcset != "" {
				buffer.WriteString(` srcset="`)
				buffer.WriteString(srcset)
				buffer.WriteString(`"`)
			}
		}
	}

	if text := GetImageViewAltText(imageView); text != "" {
		buffer.WriteString(` alt="`)
		buffer.WriteString(text)
		buffer.WriteString(`"`)
	}

	buffer.WriteString(` onload="imageLoaded(this, event)"`)

	if len(imageView.imageListeners(ErrorEvent)) > 0 {
		buffer.WriteString(` onerror="imageError(this, event)"`)
	}
}

func (imageView *imageViewData) cssStyle(self View, builder cssBuilder) {
	imageView.viewData.cssStyle(self, builder)

	if value, ok := enumProperty(imageView, Fit, imageView.session, 0); ok {
		builder.add("object-fit", enumProperties[Fit].cssValues[value])
	} else {
		builder.add("object-fit", "none")
	}

	vAlign := GetImageViewVerticalAlign(imageView)
	hAlign := GetImageViewHorizontalAlign(imageView)
	if vAlign != CenterAlign || hAlign != CenterAlign {
		var position string
		switch hAlign {
		case LeftAlign:
			position = "left"
		case RightAlign:
			position = "right"
		default:
			position = "center"
		}

		switch vAlign {
		case TopAlign:
			position += " top"
		case BottomAlign:
			position += " bottom"
		default:
			position += " center"
		}

		builder.add("object-position", position)
	}
}

func (imageView *imageViewData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case "imageViewError":
		for _, listener := range imageView.imageListeners(ErrorEvent) {
			listener(imageView)
		}

	case "imageViewLoaded":
		imageView.naturalWidth = dataFloatProperty(data, "natural-width")
		imageView.naturalHeight = dataFloatProperty(data, "natural-height")
		imageView.currentSrc, _ = data.PropertyValue("current-src")

		for _, listener := range imageView.imageListeners(LoadedEvent) {
			listener(imageView)
		}

	default:
		return imageView.viewData.handleCommand(self, command, data)
	}
	return true
}

func (imageView *imageViewData) NaturalSize() (float64, float64) {
	return imageView.naturalWidth, imageView.naturalHeight
}

func (imageView *imageViewData) CurrentSource() string {
	return imageView.currentSrc
}

// GetImageViewSource returns the image URL of an ImageView subview.
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetImageViewSource(view View, subviewID ...string) string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if image, ok := imageProperty(view, Source, view.Session()); ok {
			return image
		}
	}

	return ""
}

// GetImageViewAltText returns an alternative text description of an ImageView subview.
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetImageViewAltText(view View, subviewID ...string) string {
	if len(subviewID) > 0 && subviewID[0] != "" {
		view = ViewByID(view, subviewID[0])
	}

	if view != nil {
		if value := view.getRaw(AltText); value != nil {
			if text, ok := value.(string); ok {
				text, _ = view.Session().GetString(text)
				return text
			}
		}
	}
	return ""
}

// GetImageViewFit returns how the content of a replaced ImageView subview:
// NoneFit (0), ContainFit (1), CoverFit (2), FillFit (3), or ScaleDownFit (4).
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetImageViewFit(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, Fit, NoneFit, false)
}

// GetImageViewVerticalAlign return the vertical align of an ImageView subview: TopAlign (0), BottomAlign (1), CenterAlign (2)
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetImageViewVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, ImageVerticalAlign, LeftAlign, false)
}

// GetImageViewHorizontalAlign return the vertical align of an ImageView subview: LeftAlign (0), RightAlign (1), CenterAlign (2)
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetImageViewHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, ImageHorizontalAlign, LeftAlign, false)
}
