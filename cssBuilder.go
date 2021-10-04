package rui

import (
	"strings"
)

var systemStyles = map[string]string{
	"ruiApp":              "body",
	"ruiDefault":          "div",
	"ruiArticle":          "article",
	"ruiSection":          "section",
	"ruiAside":            "aside",
	"ruiHeader":           "header",
	"ruiMain":             "main",
	"ruiFooter":           "footer",
	"ruiNavigation":       "nav",
	"ruiFigure":           "figure",
	"ruiFigureCaption":    "figcaption",
	"ruiButton":           "button",
	"ruiP":                "p",
	"ruiParagraph":        "p",
	"ruiH1":               "h1",
	"ruiH2":               "h2",
	"ruiH3":               "h3",
	"ruiH4":               "h4",
	"ruiH5":               "h5",
	"ruiH6":               "h6",
	"ruiBlockquote":       "blockquote",
	"ruiCode":             "code",
	"ruiTable":            "table",
	"ruiTableHead":        "thead",
	"ruiTableFoot":        "tfoot",
	"ruiTableRow":         "tr",
	"ruiTableColumn":      "col",
	"ruiTableCell":        "td",
	"ruiDropDownList":     "select",
	"ruiDropDownListItem": "option",
}

var disabledStyles = []string{
	"ruiRoot",
	"ruiPopupLayer",
	"ruiAbsoluteLayout",
	"ruiGridLayout",
	"ruiListLayout",
	"ruiStackLayout",
	"ruiStackPageLayout",
	"ruiTabsLayout",
	"ruiImageView",
	"ruiListView",
}

type cssBuilder interface {
	add(key, value string)
	addValues(key, separator string, values ...string)
}

type viewCSSBuilder struct {
	buffer *strings.Builder
}

type cssValueBuilder struct {
	buffer *strings.Builder
}

type cssStyleBuilder struct {
	buffer *strings.Builder
	media  bool
}

func (builder *viewCSSBuilder) finish() string {
	if builder.buffer == nil {
		return ""
	}

	result := builder.buffer.String()
	freeStringBuilder(builder.buffer)
	builder.buffer = nil
	return result
}

func (builder *viewCSSBuilder) add(key, value string) {
	if value != "" {
		if builder.buffer == nil {
			builder.buffer = allocStringBuilder()
		} else if builder.buffer.Len() > 0 {
			builder.buffer.WriteRune(' ')
		}

		builder.buffer.WriteString(key)
		builder.buffer.WriteString(": ")
		builder.buffer.WriteString(value)
		builder.buffer.WriteRune(';')
	}
}

func (builder *viewCSSBuilder) addValues(key, separator string, values ...string) {
	if len(values) == 0 {
		return
	}

	if builder.buffer == nil {
		builder.buffer = allocStringBuilder()
	} else if builder.buffer.Len() > 0 {
		builder.buffer.WriteRune(' ')
	}

	builder.buffer.WriteString(key)
	builder.buffer.WriteString(": ")
	for i, value := range values {
		if i > 0 {
			builder.buffer.WriteString(separator)
		}
		builder.buffer.WriteString(value)
	}
	builder.buffer.WriteRune(';')
}

func (builder *cssValueBuilder) finish() string {
	if builder.buffer == nil {
		return ""
	}

	result := builder.buffer.String()
	freeStringBuilder(builder.buffer)
	builder.buffer = nil
	return result
}

func (builder *cssValueBuilder) add(key, value string) {
	if value != "" {
		if builder.buffer == nil {
			builder.buffer = allocStringBuilder()
		}
		builder.buffer.WriteString(value)
	}
}

func (builder *cssValueBuilder) addValues(key, separator string, values ...string) {
	if len(values) > 0 {
		if builder.buffer == nil {
			builder.buffer = allocStringBuilder()
		}
		for i, value := range values {
			if i > 0 {
				builder.buffer.WriteString(separator)
			}
			builder.buffer.WriteString(value)
		}
	}
}

func (builder *cssStyleBuilder) init() {
	builder.buffer = allocStringBuilder()
	builder.buffer.Grow(16 * 1024)
}

func (builder *cssStyleBuilder) finish() string {
	if builder.buffer == nil {
		return ""
	}

	result := builder.buffer.String()
	freeStringBuilder(builder.buffer)
	builder.buffer = nil
	return result
}

func (builder *cssStyleBuilder) startMedia(rule string) {
	if builder.buffer == nil {
		builder.init()
	}
	builder.buffer.WriteString(`@media screen`)
	builder.buffer.WriteString(rule)
	builder.buffer.WriteString(` {\n`)
	builder.media = true
}

func (builder *cssStyleBuilder) endMedia() {
	if builder.buffer == nil {
		builder.init()
	}
	builder.buffer.WriteString(`}\n`)
	builder.media = false
}

func (builder *cssStyleBuilder) startStyle(name string) {
	for _, disabledName := range disabledStyles {
		if name == disabledName {
			return
		}
	}

	if builder.buffer == nil {
		builder.init()
	}
	if builder.media {
		builder.buffer.WriteString(`\t`)
	}

	if sysName, ok := systemStyles[name]; ok {
		builder.buffer.WriteString(sysName)
	} else {
		builder.buffer.WriteRune('.')
		builder.buffer.WriteString(name)
	}

	builder.buffer.WriteString(` {\n`)
}

func (builder *cssStyleBuilder) endStyle() {
	if builder.buffer == nil {
		builder.init()
	}
	if builder.media {
		builder.buffer.WriteString(`\t`)
	}
	builder.buffer.WriteString(`}\n`)
}

func (builder *cssStyleBuilder) startAnimation(name string) {
	if builder.buffer == nil {
		builder.init()
	}

	builder.media = true
	builder.buffer.WriteString(`\n@keyframes `)
	builder.buffer.WriteString(name)
	builder.buffer.WriteString(` {\n`)
}

func (builder *cssStyleBuilder) endAnimation() {
	if builder.buffer == nil {
		builder.init()
	}
	builder.buffer.WriteString(`}\n`)
	builder.media = false
}

func (builder *cssStyleBuilder) startAnimationFrame(name string) {
	if builder.buffer == nil {
		builder.init()
	}

	builder.buffer.WriteString(`\t`)
	builder.buffer.WriteString(name)
	builder.buffer.WriteString(` {\n`)
}

func (builder *cssStyleBuilder) endAnimationFrame() {
	if builder.buffer == nil {
		builder.init()
	}
	builder.buffer.WriteString(`\t}\n`)
}

func (builder *cssStyleBuilder) add(key, value string) {
	if value != "" {
		if builder.buffer == nil {
			builder.init()
		}
		if builder.media {
			builder.buffer.WriteString(`\t`)
		}
		builder.buffer.WriteString(`\t`)
		builder.buffer.WriteString(key)
		builder.buffer.WriteString(`: `)
		builder.buffer.WriteString(value)
		builder.buffer.WriteString(`;\n`)
	}
}

func (builder *cssStyleBuilder) addValues(key, separator string, values ...string) {
	if len(values) == 0 {
		return
	}

	if builder.buffer == nil {
		builder.init()
	}
	if builder.media {
		builder.buffer.WriteString(`\t`)
	}
	builder.buffer.WriteString(`\t`)
	builder.buffer.WriteString(key)
	builder.buffer.WriteString(`: `)
	for i, value := range values {
		if i > 0 {
			builder.buffer.WriteString(separator)
		}
		builder.buffer.WriteString(value)
	}
	builder.buffer.WriteString(`;\n`)
}
