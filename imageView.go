package rui

import (
	"fmt"
	"strings"
)

const (
	// LoadedEvent is the constant for the "loaded-event" property tag.
	// The "loaded-event" event occurs event occurs when the image has been loaded.
	LoadedEvent = "loaded-event"
	// ErrorEvent is the constant for the "error-event" property tag.
	// The "error-event" event occurs event occurs when the image loading failed.
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

// ImageView - image View
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
	view.Init(session)
	setInitParams(view, params)
	return view
}

func newImageView(session Session) View {
	return NewImageView(session, nil)
}

// Init initialize fields of imageView by default values
func (imageView *imageViewData) Init(session Session) {
	imageView.viewData.Init(session)
	imageView.tag = "ImageView"
	//imageView.systemClass = "ruiImageView"

}

func (imageView *imageViewData) String() string {
	return getViewString(imageView)
}

func (imageView *imageViewData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case "source":
		tag = Source

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
			updateProperty(imageView.htmlID(), "src", "", imageView.session)
			removeProperty(imageView.htmlID(), "srcset", imageView.session)

		case AltText:
			updateInnerHTML(imageView.htmlID(), imageView.session)

		case ImageVerticalAlign, ImageHorizontalAlign:
			updateCSSStyle(imageView.htmlID(), imageView.session)
		}
	}
}

func (imageView *imageViewData) Set(tag string, value interface{}) bool {
	return imageView.set(imageView.normalizeTag(tag), value)
}

func valueToImageListeners(value interface{}) ([]func(ImageView), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(ImageView):
		return []func(ImageView){value}, true

	case func():
		fn := func(ImageView) {
			value()
		}
		return []func(ImageView){fn}, true

	case []func(ImageView):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(ImageView), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(ImageView) {
				v()
			}
		}
		return listeners, true

	case []interface{}:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(ImageView), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(ImageView):
				listeners[i] = v

			case func():
				listeners[i] = func(ImageView) {
					v()
				}

			default:
				return nil, false
			}
		}
		return listeners, true
	}

	return nil, false
}

func (imageView *imageViewData) set(tag string, value interface{}) bool {
	if value == nil {
		imageView.remove(tag)
		return true
	}

	switch tag {
	case Source:
		if text, ok := value.(string); ok {
			imageView.properties[Source] = text
			if imageView.created {
				src := text
				if src != "" && src[0] == '@' {
					src, _ = imageProperty(imageView, Source, imageView.session)
				}
				updateProperty(imageView.htmlID(), "src", src, imageView.session)
				if srcset := imageView.srcSet(src); srcset != "" {
					updateProperty(imageView.htmlID(), "srcset", srcset, imageView.session)
				} else {
					removeProperty(imageView.htmlID(), "srcset", imageView.session)
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
		if listeners, ok := valueToImageListeners(value); ok {
			imageView.properties[tag] = listeners
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

func (imageView *imageViewData) Get(tag string) interface{} {
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

/*
func (imageView *imageViewData) closeHTMLTag() bool {
	return false
}
*/

func (imageView *imageViewData) htmlProperties(self View, buffer *strings.Builder) {
	imageView.viewData.htmlProperties(self, buffer)

	if imageResource, ok := imageProperty(imageView, Source, imageView.Session()); ok && imageResource != "" {
		if imageResource[0] == '@' {
			if image, ok := imageView.Session().ImageConstant(imageResource[1:]); ok {
				imageResource = image
			} else {
				imageResource = ""
			}
		}

		if imageResource != "" {
			buffer.WriteString(` src="`)
			buffer.WriteString(imageResource)
			buffer.WriteString(`"`)
			if srcset := imageView.srcSet(imageResource); srcset != "" {
				buffer.WriteString(` srcset="`)
				buffer.WriteString(srcset)
				buffer.WriteString(`"`)
			}
		}
	}

	if text := GetImageViewAltText(imageView, ""); text != "" {
		buffer.WriteString(` alt="`)
		buffer.WriteString(textToJS(text))
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

	vAlign := GetImageViewVerticalAlign(imageView, "")
	hAlign := GetImageViewHorizontalAlign(imageView, "")
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
// If the second argument (subviewID) is "" then a left position of the first argument (view) is returned
func GetImageViewSource(view View, subviewID string) string {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		if image, ok := imageProperty(view, Source, view.Session()); ok {
			return image
		}
	}

	return ""
}

// GetImageViewAltText returns an alternative text description of an ImageView subview.
// If the second argument (subviewID) is "" then a left position of the first argument (view) is returned
func GetImageViewAltText(view View, subviewID string) string {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
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
// If the second argument (subviewID) is "" then a left position of the first argument (view) is returned
func GetImageViewFit(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if value, ok := enumStyledProperty(view, Fit, NoneFit); ok {
		return value
	}
	return 0
}

// GetImageViewVerticalAlign return the vertical align of an ImageView subview: TopAlign (0), BottomAlign (1), CenterAlign (2)
// If the second argument (subviewID) is "" then a left position of the first argument (view) is returned
func GetImageViewVerticalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if align, ok := enumStyledProperty(view, ImageVerticalAlign, LeftAlign); ok {
		return align
	}
	return CenterAlign
}

// GetImageViewHorizontalAlign return the vertical align of an ImageView subview: LeftAlign (0), RightAlign (1), CenterAlign (2)
// If the second argument (subviewID) is "" then a left position of the first argument (view) is returned
func GetImageViewHorizontalAlign(view View, subviewID string) int {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if align, ok := enumStyledProperty(view, ImageHorizontalAlign, LeftAlign); ok {
		return align
	}
	return CenterAlign
}
