package rui

import (
	"io"
	"net/http"
	"os"
	"strings"
)

// SvgImageView represents an SvgImageView view
type SvgImageView interface {
	View
}

type svgImageViewData struct {
	viewData
}

// NewSvgImageView create new SvgImageView object and return it
func NewSvgImageView(session Session, params Params) SvgImageView {
	view := new(svgImageViewData)
	view.init(session)
	setInitParams(view, params)
	return view
}

func newSvgImageView(session Session) View {
	return NewSvgImageView(session, nil)
}

// Init initialize fields of imageView by default values
func (imageView *svgImageViewData) init(session Session) {
	imageView.viewData.init(session)
	imageView.tag = "SvgImageView"
	imageView.systemClass = "ruiSvgImageView"
}

func (imageView *svgImageViewData) String() string {
	return getViewString(imageView, nil)
}

func (imageView *svgImageViewData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case Source, "source":
		tag = Content

	case VerticalAlign:
		tag = CellVerticalAlign

	case HorizontalAlign:
		tag = CellHorizontalAlign
	}
	return tag
}

func (imageView *svgImageViewData) Remove(tag string) {
	imageView.remove(imageView.normalizeTag(tag))
}

func (imageView *svgImageViewData) remove(tag string) {
	imageView.viewData.remove(tag)

	if imageView.created {
		switch tag {
		case Content:
			updateInnerHTML(imageView.htmlID(), imageView.session)
		}
	}
}

func (imageView *svgImageViewData) Set(tag string, value any) bool {
	return imageView.set(imageView.normalizeTag(tag), value)
}

func (imageView *svgImageViewData) set(tag string, value any) bool {
	if value == nil {
		imageView.remove(tag)
		return true
	}

	switch tag {
	case Content:
		if text, ok := value.(string); ok {
			imageView.properties[Content] = text
			if imageView.created {
				updateInnerHTML(imageView.htmlID(), imageView.session)
			}
			imageView.propertyChangedEvent(Content)
			return true
		}
		notCompatibleType(Source, value)
		return false

	default:
		return imageView.viewData.set(tag, value)
	}
}

func (imageView *svgImageViewData) Get(tag string) any {
	return imageView.viewData.get(imageView.normalizeTag(tag))
}

func (imageView *svgImageViewData) htmlTag() string {
	return "div"
}

func (imageView *svgImageViewData) htmlSubviews(self View, buffer *strings.Builder) {
	if value := imageView.getRaw(Content); value != nil {
		if content, ok := value.(string); ok && content != "" {
			if strings.HasPrefix(content, "@") {
				if name, ok := imageView.session.ImageConstant(content[1:]); ok {
					content = name
				}
			}

			if image, ok := resources.images[content]; ok {
				if image.fs != nil {
					if data, err := image.fs.ReadFile(image.path); err == nil {
						buffer.WriteString(string(data))
						return
					} else {
						DebugLog(err.Error())
					}
				} else if data, err := os.ReadFile(image.path); err == nil {
					buffer.WriteString(string(data))
					return
				} else {
					DebugLog(err.Error())
				}
			}

			if strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") {
				resp, err := http.Get(content)
				if err == nil {
					defer resp.Body.Close()
					if body, err := io.ReadAll(resp.Body); err == nil {
						buffer.WriteString(string(body))
						return
					}
				}

				DebugLog(err.Error())
			}

			buffer.WriteString(content)
		}
	}
}

// GetSvgImageViewVerticalAlign return the vertical align of an SvgImageView subview: TopAlign (0), BottomAlign (1), CenterAlign (2)
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetSvgImageViewVerticalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CellVerticalAlign, LeftAlign, false)
}

// GetSvgImageViewHorizontalAlign return the vertical align of an SvgImageView subview: LeftAlign (0), RightAlign (1), CenterAlign (2)
// If the second argument (subviewID) is not specified or it is "" then a left position of the first argument (view) is returned
func GetSvgImageViewHorizontalAlign(view View, subviewID ...string) int {
	return enumStyledProperty(view, subviewID, CellHorizontalAlign, LeftAlign, false)
}
