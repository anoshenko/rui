package rui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	DefaultMedia   = 0
	PortraitMedia  = 1
	LandscapeMedia = 2
)

type MediaStyle struct {
	Orientation int
	MaxWidth    int
	MaxHeight   int
	Styles      map[string]ViewStyle
}

type theme struct {
	name           string
	constants      map[string]string
	touchConstants map[string]string
	colors         map[string]string
	darkColors     map[string]string
	images         map[string]string
	darkImages     map[string]string
	styles         map[string]ViewStyle
	mediaStyles    []MediaStyle
}

type Theme interface {
	fmt.Stringer
	Name() string
	Constant(tag string) (string, string)
	SetConstant(tag string, value, touchUIValue string)
	// ConstantTags returns the list of all available constants
	ConstantTags() []string
	Color(tag string) (string, string)
	SetColor(tag, color, darkUIColor string)
	// ColorTags returns the list of all available color constants
	ColorTags() []string
	Image(tag string) (string, string)
	SetImage(tag, image, darkUIImage string)
	// ImageConstantTags returns the list of all available image constants
	ImageConstantTags() []string
	Append(anotherTheme Theme)

	constant(tag string, touchUI bool) string
	color(tag string, darkUI bool) string
	image(tag string, darkUI bool) string
	style(tag string) ViewStyle
	cssText(session Session) string
	data() *theme
}

func (rule MediaStyle) cssText() string {
	builder := allocStringBuilder()
	defer freeStringBuilder(builder)

	switch rule.Orientation {
	case PortraitMedia:
		builder.WriteString(" and (orientation: portrait)")

	case LandscapeMedia:
		builder.WriteString(" and (orientation: landscape)")
	}

	if rule.MaxWidth > 0 {
		builder.WriteString(" and (max-width: ")
		builder.WriteString(strconv.Itoa(rule.MaxWidth))
		builder.WriteString("px)")
	}

	if rule.MaxHeight > 0 {
		builder.WriteString(" and (max-height: ")
		builder.WriteString(strconv.Itoa(rule.MaxHeight))
		builder.WriteString("px)")
	}

	return builder.String()
}

func parseMediaRule(text string) (MediaStyle, bool) {
	rule := MediaStyle{
		Orientation: DefaultMedia,
		MaxWidth:    0,
		MaxHeight:   0,
		Styles:      map[string]ViewStyle{},
	}

	elements := strings.Split(text, ":")
	for i := 1; i < len(elements); i++ {
		switch element := elements[i]; element {
		case "portrait":
			if rule.Orientation != DefaultMedia {
				ErrorLog(`Duplicate orientation tag in the style section "` + text + `"`)
				return rule, false
			}
			rule.Orientation = PortraitMedia

		case "landscape":
			if rule.Orientation != DefaultMedia {
				ErrorLog(`Duplicate orientation tag in the style section "` + text + `"`)
				return rule, false
			}
			rule.Orientation = LandscapeMedia

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
				if rule.MaxWidth != 0 {
					ErrorLog(`Duplicate "width" tag in the style section "` + text + `"`)
					return rule, false
				}
				rule.MaxWidth = size
			} else if size, ok := elementSize("height"); !ok || size > 0 {
				if !ok {
					return rule, false
				}
				if rule.MaxHeight != 0 {
					ErrorLog(`Duplicate "height" tag in the style section "` + text + `"`)
					return rule, false
				}
				rule.MaxHeight = size
			} else {
				ErrorLogF(`Unknown elemnet "%s" in the style section name "%s"`, element, text)
				return rule, false
			}
		}
	}
	return rule, true
}

var defaultTheme = NewTheme("")

func NewTheme(name string) Theme {
	result := new(theme)
	result.init()
	result.name = name
	return result
}

func CreateThemeFromText(text string) (Theme, bool) {
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
	theme.images = map[string]string{}
	theme.darkImages = map[string]string{}
	theme.styles = map[string]ViewStyle{}
	theme.mediaStyles = []MediaStyle{}
}

func (theme *theme) Name() string {
	return theme.name
}

func (theme *theme) Constant(tag string) (string, string) {
	return theme.constants[tag], theme.touchConstants[tag]
}

func (theme *theme) SetConstant(tag, value, touchUIValue string) {
	value = strings.Trim(value, " \t")
	if value == "" {
		delete(theme.constants, tag)
		delete(theme.touchConstants, tag)
	} else {
		theme.constants[tag] = value
		touchUIValue = strings.Trim(touchUIValue, " \t")
		if touchUIValue == "" {
			delete(theme.touchConstants, tag)
		} else {
			theme.touchConstants[tag] = touchUIValue
		}
	}
}

func (theme *theme) Color(tag string) (string, string) {
	return theme.colors[tag], theme.darkColors[tag]
}

func (theme *theme) SetColor(tag, color, darkUIColor string) {
	color = strings.Trim(color, " \t")
	if color == "" {
		delete(theme.colors, tag)
		delete(theme.darkColors, tag)
	} else {
		theme.colors[tag] = color
		darkUIColor = strings.Trim(darkUIColor, " \t")
		if darkUIColor == "" {
			delete(theme.darkColors, tag)
		} else {
			theme.darkColors[tag] = darkUIColor
		}
	}
}

func (theme *theme) Image(tag string) (string, string) {
	return theme.images[tag], theme.darkImages[tag]
}

func (theme *theme) SetImage(tag, image, darkUIImage string) {
	image = strings.Trim(image, " \t")
	if image == "" {
		delete(theme.images, tag)
		delete(theme.darkImages, tag)
	} else {
		theme.images[tag] = image
		darkUIImage = strings.Trim(darkUIImage, " \t")
		if darkUIImage == "" {
			delete(theme.darkImages, tag)
		} else {
			theme.darkImages[tag] = darkUIImage
		}
	}
}

func (theme *theme) ConstantTags() []string {
	keys := make([]string, 0, len(theme.constants))
	for k := range theme.constants {
		keys = append(keys, k)
	}

	for tag := range theme.touchConstants {
		if _, ok := theme.constants[tag]; !ok {
			keys = append(keys, tag)
		}
	}

	sort.Strings(keys)
	return keys
}

func (theme *theme) ColorTags() []string {
	keys := make([]string, 0, len(theme.colors))
	for k := range theme.colors {
		keys = append(keys, k)
	}

	for tag := range theme.darkColors {
		if _, ok := theme.colors[tag]; !ok {
			keys = append(keys, tag)
		}
	}

	sort.Strings(keys)
	return keys
}

func (theme *theme) ImageConstantTags() []string {
	keys := make([]string, 0, len(theme.colors))
	for k := range theme.images {
		keys = append(keys, k)
	}

	for tag := range theme.darkImages {
		if _, ok := theme.images[tag]; !ok {
			keys = append(keys, tag)
		}
	}

	sort.Strings(keys)
	return keys
}

func (theme *theme) data() *theme {
	return theme
}

func (theme *theme) Append(anotherTheme Theme) {
	if theme.constants == nil {
		theme.init()
	}

	another := anotherTheme.data()
	for tag, constant := range another.constants {
		theme.constants[tag] = constant
	}

	for tag, constant := range another.touchConstants {
		theme.touchConstants[tag] = constant
	}

	for tag, color := range another.colors {
		theme.colors[tag] = color
	}

	for tag, color := range another.darkColors {
		theme.darkColors[tag] = color
	}

	for tag, image := range another.images {
		theme.images[tag] = image
	}

	for tag, image := range another.darkImages {
		theme.darkImages[tag] = image
	}

	for tag, style := range another.styles {
		theme.styles[tag] = style
	}

	for _, anotherMedia := range another.mediaStyles {
		exists := false
		for _, media := range theme.mediaStyles {
			if anotherMedia.MaxHeight == media.MaxHeight &&
				anotherMedia.MaxWidth == media.MaxWidth &&
				anotherMedia.Orientation == media.Orientation {
				for tag, style := range anotherMedia.Styles {
					media.Styles[tag] = style
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

	for tag, style := range theme.styles {
		/*var style viewStyle
		style.init()
		for tag, value := range obj {
			style.Set(tag, value)
		}*/
		builder.startStyle(tag)
		style.cssViewStyle(&builder, session)
		builder.endStyle()
	}

	for _, media := range theme.mediaStyles {
		builder.startMedia(media.cssText())
		for tag, style := range media.Styles {
			/*var style viewStyle
			style.init()
			for tag, value := range obj {
				style.Set(tag, value)
			}*/
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

	objToStyle := func(obj DataObject) ViewStyle {
		params := Params{}
		for i := 0; i < obj.PropertyCount(); i++ {
			if node := obj.Property(i); node != nil {
				switch node.Type() {
				case ArrayNode:
					params[node.Tag()] = node.ArrayElements()

				case ObjectNode:
					params[node.Tag()] = node.Object()

				default:
					params[node.Tag()] = node.Text()
				}
			}
		}
		return NewViewStyle(params)
	}

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

			case "images":
				if d.Type() == ObjectNode {
					if obj := d.Object(); obj != nil {
						objCount := obj.PropertyCount()
						for k := 0; k < objCount; k++ {
							if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
								theme.images[prop.Tag()] = prop.Text()
							}
						}
					}
				}

			case "images:dark":
				if d.Type() == ObjectNode {
					if obj := d.Object(); obj != nil {
						objCount := obj.PropertyCount()
						for k := 0; k < objCount; k++ {
							if prop := obj.Property(k); prop != nil && prop.Type() == TextNode {
								theme.darkImages[prop.Tag()] = prop.Text()
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
								theme.styles[obj.Tag()] = objToStyle(obj)
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
									rule.Styles[obj.Tag()] = objToStyle(obj)
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
			if theme.mediaStyles[i].Orientation != theme.mediaStyles[j].Orientation {
				return theme.mediaStyles[i].Orientation < theme.mediaStyles[j].Orientation
			}
			if theme.mediaStyles[i].MaxWidth != theme.mediaStyles[j].MaxWidth {
				return theme.mediaStyles[i].MaxWidth < theme.mediaStyles[j].MaxWidth
			}
			return theme.mediaStyles[i].MaxHeight < theme.mediaStyles[j].MaxHeight
		})
	}
}

func (theme *theme) constant(tag string, touchUI bool) string {
	result := ""
	if touchUI {
		if value, ok := theme.touchConstants[tag]; ok {
			result = value
		}
	}
	if result == "" {
		if value, ok := theme.constants[tag]; ok {
			result = value
		}
	}
	return result
}

func (theme *theme) color(tag string, darkUI bool) string {
	result := ""
	if darkUI {
		if value, ok := theme.darkColors[tag]; ok {
			result = value
		}
	}
	if result == "" {
		if value, ok := theme.colors[tag]; ok {
			result = value
		}
	}
	return result
}

func (theme *theme) image(tag string, darkUI bool) string {
	result := ""
	if darkUI {
		if value, ok := theme.darkImages[tag]; ok {
			result = value
		}
	}
	if result == "" {
		if value, ok := theme.images[tag]; ok {
			result = value
		}
	}
	return result
}

func (theme *theme) style(tag string) ViewStyle {
	if style, ok := theme.styles[tag]; ok {
		return style
	}

	return nil
}

func (theme *theme) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	writeString := func(text string) {
		if strings.ContainsAny(text, " \t\n\r\\\"'`,;{}[]()") {
			replace := []struct{ old, new string }{
				{old: "\\", new: `\\`},
				{old: "\t", new: `\t`},
				{old: "\r", new: `\r`},
				{old: "\n", new: `\n`},
				{old: "\"", new: `\"`},
			}
			for _, s := range replace {
				text = strings.Replace(text, s.old, s.new, -1)
			}
			buffer.WriteRune('"')
			buffer.WriteString(text)
			buffer.WriteRune('"')
		} else {
			buffer.WriteString(text)
		}
	}

	writeConstants := func(tag string, constants map[string]string) {
		count := len(constants)
		if count == 0 {
			return
		}

		buffer.WriteString("\t")
		buffer.WriteString(tag)
		buffer.WriteString(" = _{\n")

		tags := make([]string, 0, count)
		for name := range constants {
			tags = append(tags, name)
		}
		sort.Strings(tags)
		for _, name := range tags {
			if value, ok := constants[name]; ok && value != "" {
				buffer.WriteString("\t\t")
				writeString(name)
				buffer.WriteString(" = ")
				writeString(value)
				buffer.WriteString(",\n")
			}
		}

		buffer.WriteString("\t},\n")
	}

	buffer.WriteString("theme {\n")
	writeConstants("colors", theme.colors)
	writeConstants("colors:dark", theme.darkColors)
	writeConstants("images", theme.images)
	writeConstants("images:dark", theme.darkImages)
	writeConstants("constants", theme.constants)
	writeConstants("constants:touch", theme.touchConstants)

	buffer.WriteString("}\n")
	return buffer.String()
}
