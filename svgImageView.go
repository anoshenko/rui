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
	return new(svgImageViewData) // NewSvgImageView(session, nil)
}

// Init initialize fields of imageView by default values
func (imageView *svgImageViewData) init(session Session) {
	imageView.viewData.init(session)
	imageView.tag = "SvgImageView"
	imageView.systemClass = "ruiSvgImageView"
	imageView.normalize = normalizeSvgImageViewTag
	imageView.set = imageView.setFunc
	imageView.changed = imageView.propertyChanged

}

func normalizeSvgImageViewTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
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

func (imageView *svgImageViewData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {
	case Content:
		if text, ok := value.(string); ok {
			imageView.setRaw(Content, text)
			return []PropertyName{tag}
		}
		notCompatibleType(Source, value)
		return nil

	default:
		return imageView.viewData.setFunc(tag, value)
	}
}

func (imageView *svgImageViewData) propertyChanged(tag PropertyName) {
	switch tag {
	case Content:
		updateInnerHTML(imageView.htmlID(), imageView.Session())

	default:
		imageView.viewData.propertyChanged(tag)
	}
}

func (imageView *svgImageViewData) htmlTag() string {
	return "div"
}

func (imageView *svgImageViewData) writeSvg(data []byte, buffer *strings.Builder) {
	text := string(data)
	index := strings.Index(text, "<svg")
	if index > 0 {
		text = text[index:]
	}

	index = strings.Index(text, "\n")
	for index >= 0 {
		if index > 0 && text[index-1] == '\r' {
			buffer.WriteString(text[:index-1])
		} else {
			buffer.WriteString(text[:index])
		}

		end := len(text)
		index++
		for index < end && (text[index] == ' ' || text[index] == '\t' || text[index] == '\r' || text[index] == '\n') {
			index++
		}

		text = text[index:]
		index = strings.Index(text, "\n")
	}

	buffer.WriteString(text)
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
						imageView.writeSvg(data, buffer)
						return
					} else {
						DebugLog(err.Error())
					}
				} else if data, err := os.ReadFile(image.path); err == nil {
					imageView.writeSvg(data, buffer)
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
						imageView.writeSvg(body, buffer)
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
