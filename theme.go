package rui

import (
	"sort"
	"strconv"
	"strings"
)

const (
	defaultMedia   = 0
	portraitMedia  = 1
	landscapeMedia = 2
)

type mediaStyle struct {
	orientation int
	width       int
	height      int
	styles      map[string]DataObject
}

func (rule mediaStyle) cssText() string {
	builder := allocStringBuilder()
	defer freeStringBuilder(builder)

	switch rule.orientation {
	case portraitMedia:
		builder.WriteString(" and (orientation: portrait)")

	case landscapeMedia:
		builder.WriteString(" and (orientation: landscape)")
	}

	if rule.width > 0 {
		builder.WriteString(" and (max-width: ")
		builder.WriteString(strconv.Itoa(rule.width))
		builder.WriteString("px)")
	}

	if rule.height > 0 {
		builder.WriteString(" and (max-height: ")
		builder.WriteString(strconv.Itoa(rule.height))
		builder.WriteString("px)")
	}

	return builder.String()
}

func parseMediaRule(text string) (mediaStyle, bool) {
	rule := mediaStyle{orientation: defaultMedia, width: 0, height: 0, styles: map[string]DataObject{}}
	elements := strings.Split(text, ":")
	for i := 1; i < len(elements); i++ {
		switch element := elements[i]; element {
		case "portrait":
			if rule.orientation != defaultMedia {
				ErrorLog(`Duplicate orientation tag in the style section "` + text + `"`)
				return rule, false
			}
			rule.orientation = portraitMedia

		case "landscape":
			if rule.orientation != defaultMedia {
				ErrorLog(`Duplicate orientation tag in the style section "` + text + `"`)
				return rule, false
			}
			rule.orientation = landscapeMedia

		default:
			elementSize := func(name string) (int, bool) {
				if strings.HasPrefix(element, name) {
					size, err := strconv.Atoi(element[len(name):])
					if err == nil && size > 0 {
						return size, true
					}
					ErrorLogF(`Invalid style section name "%s": %s`, text, err.Error())
					return 0, false
				}
				return 0, true
			}

			if size, ok := elementSize("width"); !ok || size > 0 {
				if !ok {
					return rule, false
				}
				if rule.width != 0 {
					ErrorLog(`Duplicate "width" tag in the style section "` + text + `"`)
					return rule, false
				}
				rule.width = size
			} else if size, ok := elementSize("height"); !ok || size > 0 {
				if !ok {
					return rule, false
				}
				if rule.height != 0 {
					ErrorLog(`Duplicate "height" tag in the style section "` + text + `"`)
					return rule, false
				}
				rule.height = size
			} else {
				ErrorLogF(`Unknown elemnet "%s" in the style section name "%s"`, element, text)
				return rule, false
			}
		}
	}
	return rule, true
}

type theme struct {
	name           string
	constants      map[string]string
	touchConstants map[string]string
	colors         map[string]string
	darkColors     map[string]string
	styles         map[string]DataObject
	mediaStyles    []mediaStyle
}

var defaultTheme = new(theme)

func newTheme(text string) (*theme, bool) {
	result := new(theme)
	result.init()
	ok := result.addText(text)
	return result, ok
}

func (theme *theme) init() {
	theme.constants = map[string]string{}
	theme.touchConstants = map[string]string{}
	theme.colors = map[string]string{}
	theme.darkColors = map[string]string{}
	theme.styles = map[string]DataObject{}
	theme.mediaStyles = []mediaStyle{}
}

func (theme *theme) concat(anotherTheme *theme) {
	if theme.constants == nil {
		theme.init()
	}

	for tag, constant := range anotherTheme.constants {
		theme.constants[tag] = constant
	}

	for tag, constant := range anotherTheme.touchConstants {
		theme.touchConstants[tag] = constant
	}

	for tag, color := range anotherTheme.colors {
		theme.colors[tag] = color
	}

	for tag, color := range anotherTheme.darkColors {
		theme.darkColors[tag] = color
	}

	for tag, style := range anotherTheme.styles {
		theme.styles[tag] = style
	}

	for _, anotherMedia := range anotherTheme.mediaStyles {
		exists := false
		for _, media := range theme.mediaStyles {
			if anotherMedia.height == media.height &&
				anotherMedia.width == media.width &&
				anotherMedia.orientation == media.orientation {
				for tag, style := range anotherMedia.styles {
					media.styles[tag] = style
				}
				exists = true
				break
			}
		}
		if !exists {
			theme.mediaStyles = append(theme.mediaStyles, anotherMedia)
		}
	}
}

func (theme *theme) cssText(session Session) string {
	if theme.styles == nil {
		theme.init()
		return ""
	}

	var builder cssStyleBuilder
	builder.init()

	for tag, obj := range theme.styles {
		var style viewStyle
		style.init()
		parseProperties(&style, obj)
		builder.startStyle(tag)
		style.cssViewStyle(&builder, session)
		builder.endStyle()
	}

	for _, media := range theme.mediaStyles {
		builder.startMedia(media.cssText())
		for tag, obj := range media.styles {
			var style viewStyle
			style.init()
			parseProperties(&style, obj)
			builder.startStyle(tag)
			style.cssViewStyle(&builder, session)
			builder.endStyle()
		}
		builder.endMedia()
	}

	return builder.finish()
}

func (theme *theme) addText(themeText string) bool {
	data := ParseDataText(themeText)
	if data == nil {
		return false
	}

	theme.addData(data)
	return true
}

func (theme *theme) addData(data DataObject) {
	if theme.constants == nil {
		theme.init()
	}

	if data.IsObject() && data.Tag() == "theme" {
		theme.parseThemeData(data)
	}
}

func (theme *theme) parseThemeData(data DataObject) {
	count := data.PropertyCount()

	for i := 0; i < count; i++ {
		if d := data.Property(i); d != nil {
			switch tag := d.Tag(); tag {
			case "constants":
				if d.Type() == ObjectNode {
					if obj := d.Object(); obj != nil {
						objCount := obj.PropertyCount()
						for k := 0; k < objCount; k++ {
							if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
								theme.constants[prop.Tag()] = prop.Text()
							}
						}
					}
				}

			case "constants:touch":
				if d.Type() == ObjectNode {
					if obj := d.Object(); obj != nil {
						objCount := obj.PropertyCount()
						for k := 0; k < objCount; k++ {
							if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
								theme.touchConstants[prop.Tag()] = prop.Text()
							}
						}
					}
				}

			case "colors":
				if d.Type() == ObjectNode {
					if obj := d.Object(); obj != nil {
						objCount := obj.PropertyCount()
						for k := 0; k < objCount; k++ {
							if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
								theme.colors[prop.Tag()] = prop.Text()
							}
						}
					}
				}

			case "colors:dark":
				if d.Type() == ObjectNode {
					if obj := d.Object(); obj != nil {
						objCount := obj.PropertyCount()
						for k := 0; k < objCount; k++ {
							if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
								theme.darkColors[prop.Tag()] = prop.Text()
							}
						}
					}
				}

			case "styles":
				if d.Type() == ArrayNode {
					arraySize := d.ArraySize()
					for k := 0; k < arraySize; k++ {
						if element := d.ArrayElement(k); element != nil && element.IsObject() {
							if obj := element.Object(); obj != nil {
								theme.styles[obj.Tag()] = obj
							}
						}
					}
				}

			default:
				if d.Type() == ArrayNode && strings.HasPrefix(tag, "styles:") {
					if rule, ok := parseMediaRule(tag); ok {
						arraySize := d.ArraySize()
						for k := 0; k < arraySize; k++ {
							if element := d.ArrayElement(k); element != nil && element.IsObject() {
								if obj := element.Object(); obj != nil {
									rule.styles[obj.Tag()] = obj
								}
							}
						}
						theme.mediaStyles = append(theme.mediaStyles, rule)
					}
				}
			}
		}
	}

	if len(theme.mediaStyles) > 0 {
		sort.SliceStable(theme.mediaStyles, func(i, j int) bool {
			if theme.mediaStyles[i].orientation != theme.mediaStyles[j].orientation {
				return theme.mediaStyles[i].orientation < theme.mediaStyles[j].orientation
			}
			if theme.mediaStyles[i].width != theme.mediaStyles[j].width {
				return theme.mediaStyles[i].width < theme.mediaStyles[j].width
			}
			return theme.mediaStyles[i].height < theme.mediaStyles[j].height
		})
	}
}
