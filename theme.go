package rui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Constants used as a values for [MediaStyleParams] member Orientation
const (
	// DefaultMedia means that style appliance will not be related to client's window orientation
	DefaultMedia = 0

	// PortraitMedia means that style apply on clients with portrait window orientation
	PortraitMedia = 1

	// PortraitMedia means that style apply on clients with landscape window orientation
	LandscapeMedia = 2
)

// MediaStyleParams define rules when particular style will be applied
type MediaStyleParams struct {
	// Orientation for which particular style will be applied
	Orientation int

	// MinWidth for which particular style will be applied
	MinWidth int

	// MaxWidth for which particular style will be applied
	MaxWidth int

	// MinHeight for which particular style will be applied
	MinHeight int

	// MaxHeight for which particular style will be applied
	MaxHeight int
}

type mediaStyle struct {
	MediaStyleParams
	styles map[string]ViewStyle
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
	mediaStyles    []mediaStyle
}

// Theme interface to describe application's theme
type Theme interface {
	fmt.Stringer

	// Name returns a name of the theme
	Name() string

	// Constant returns normal and touch theme constant value with specific tag
	Constant(tag string) (string, string)

	// SetConstant sets a value for a constant
	SetConstant(tag string, value, touchUIValue string)

	// ConstantTags returns the list of all available constants
	ConstantTags() []string

	// Color returns normal and dark theme color constant value with specific tag
	Color(tag string) (string, string)

	// SetColor sets normal and dark theme color constant value with specific tag
	SetColor(tag, color, darkUIColor string)

	// ColorTags returns the list of all available color constants
	ColorTags() []string

	// Image returns normal and dark theme image constant value with specific tag
	Image(tag string) (string, string)

	// SetImage sets normal and dark theme image constant value with specific tag
	SetImage(tag, image, darkUIImage string)

	// ImageConstantTags returns the list of all available image constants
	ImageConstantTags() []string

	// Style returns view style by its tag
	Style(tag string) ViewStyle

	// SetStyle sets style for a tag
	SetStyle(tag string, style ViewStyle)

	// RemoveStyle removes style with provided tag
	RemoveStyle(tag string)

	// MediaStyle returns media style which correspond to provided media style parameters
	MediaStyle(tag string, params MediaStyleParams) ViewStyle

	// SetMediaStyle sets media style with provided media style parameters and a tag
	SetMediaStyle(tag string, params MediaStyleParams, style ViewStyle)

	// StyleTags returns all tags which describe a style
	StyleTags() []string

	// MediaStyles returns all media style settings which correspond to a style tag
	MediaStyles(tag string) []struct {
		Selectors string
		Params    MediaStyleParams
	}

	// Append theme to a list of themes
	Append(anotherTheme Theme)

	constant(tag string, touchUI bool) string
	color(tag string, darkUI bool) string
	image(tag string, darkUI bool) string
	style(tag string) ViewStyle
	cssText(session Session) string
	data() *theme
}

func (rule mediaStyle) cssText() string {
	builder := allocStringBuilder()
	defer freeStringBuilder(builder)

	switch rule.Orientation {
	case PortraitMedia:
		builder.WriteString(" and (orientation: portrait)")

	case LandscapeMedia:
		builder.WriteString(" and (orientation: landscape)")
	}

	writeSize := func(tag string, minSize, maxSize int) {
		if minSize != maxSize {
			if minSize > 0 {
				builder.WriteString(fmt.Sprintf(" and (min-%s: %d.001px)", tag, minSize))
			}
			if maxSize > 0 {
				builder.WriteString(fmt.Sprintf(" and (max-%s: %dpx)", tag, maxSize))
			}
		} else if minSize > 0 {
			builder.WriteString(fmt.Sprintf(" and (%s: %dpx)", tag, minSize))
		}
	}

	writeSize("width", rule.MinWidth, rule.MaxWidth)
	writeSize("height", rule.MinHeight, rule.MaxHeight)

	return builder.String()
}

func parseMediaRule(text string) (mediaStyle, bool) {
	rule := mediaStyle{
		styles: map[string]ViewStyle{},
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
			elementSize := func(name string) (int, int, bool, error) {
				if strings.HasPrefix(element, name) {
					var err error = nil
					min := 0
					max := 0
					data := element[len(name):]
					if pos := strings.Index(data, "-"); pos >= 0 {
						if pos > 0 {
							min, err = strconv.Atoi(data[:pos])
						}
						if err == nil && pos+1 < len(data) {
							max, err = strconv.Atoi(data[pos+1:])
						}
					} else {
						max, err = strconv.Atoi(data)
					}
					return min, max, true, err
				}
				return 0, 0, false, nil
			}

			if min, max, ok, err := elementSize("width"); ok {

				if err != nil {
					ErrorLogF(`Invalid style section name "%s": %s`, text, err.Error())
					return rule, false
				}
				if rule.MinWidth != 0 || rule.MaxWidth != 0 {
					ErrorLog(`Duplicate "width" tag in the style section "` + text + `"`)
					return rule, false
				}
				if min == 0 && max == 0 {
					ErrorLog(`Invalid arguments of "width" tag in the style section "` + text + `"`)
					return rule, false
				}

				rule.MinWidth = min
				rule.MaxWidth = max

			} else if min, max, ok, err := elementSize("height"); ok {

				if err != nil {
					ErrorLogF(`Invalid style section name "%s": %s`, text, err.Error())
					return rule, false
				}
				if rule.MinHeight != 0 || rule.MaxHeight != 0 {
					ErrorLog(`Duplicate "height" tag in the style section "` + text + `"`)
					return rule, false
				}
				if min == 0 && max == 0 {
					ErrorLog(`Invalid arguments of "height" tag in the style section "` + text + `"`)
					return rule, false
				}

				rule.MinHeight = min
				rule.MaxHeight = max

			} else {

				ErrorLogF(`Unknown element "%s" in the style section name "%s"`, element, text)
				return rule, false
			}
		}
	}
	return rule, true
}

var defaultTheme = NewTheme("")

// NewTheme creates a new theme with specific name and return its interface.
func NewTheme(name string) Theme {
	result := new(theme)
	result.init()
	result.name = name
	return result
}

// CreateThemeFromText creates a new theme from text and return its interface on success.
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
	theme.mediaStyles = []mediaStyle{}
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

func (theme *theme) Style(tag string) ViewStyle {
	if style, ok := theme.styles[tag]; ok {
		return style
	}
	return nil
}

func (theme *theme) SetStyle(tag string, style ViewStyle) {
	if style != nil {
		theme.styles[tag] = style
	} else {
		delete(theme.styles, tag)
	}
}

func (theme *theme) RemoveStyle(tag string) {
	tag2 := tag + ":"
	remove := func(styles map[string]ViewStyle) {
		tags := []string{tag}
		for t := range styles {
			if strings.HasPrefix(t, tag2) {
				tags = append(tags, t)
			}
		}
		for _, t := range tags {
			delete(styles, t)
		}
	}
	remove(theme.styles)
	for _, mediaStyle := range theme.mediaStyles {
		remove(mediaStyle.styles)
	}
}

func (theme *theme) MediaStyle(tag string, params MediaStyleParams) ViewStyle {
	for _, styles := range theme.mediaStyles {
		if styles.Orientation == params.Orientation &&
			styles.MaxWidth == params.MaxWidth &&
			styles.MinWidth == params.MinWidth &&
			styles.MaxHeight == params.MaxHeight &&
			styles.MinHeight == params.MinHeight {
			if style, ok := styles.styles[tag]; ok {
				return style
			}
		}
	}

	if params.Orientation == 0 && params.MaxWidth == 0 && params.MinWidth == 0 &&
		params.MaxHeight == 0 && params.MinHeight == 0 {
		return theme.style(tag)
	}

	return nil
}

func (theme *theme) SetMediaStyle(tag string, params MediaStyleParams, style ViewStyle) {
	if params.MaxWidth < 0 {
		params.MaxWidth = 0
	}
	if params.MinWidth < 0 {
		params.MinWidth = 0
	}
	if params.MaxHeight < 0 {
		params.MaxHeight = 0
	}
	if params.MinHeight < 0 {
		params.MinHeight = 0
	}

	if params.Orientation == 0 && params.MaxWidth == 0 && params.MinWidth == 0 &&
		params.MaxHeight == 0 && params.MinHeight == 0 {
		theme.SetStyle(tag, style)
		return
	}

	for i, styles := range theme.mediaStyles {
		if styles.Orientation == params.Orientation &&
			styles.MaxWidth == params.MaxWidth &&
			styles.MinWidth == params.MinWidth &&
			styles.MaxHeight == params.MaxHeight &&
			styles.MinHeight == params.MinHeight {
			if style != nil {
				theme.mediaStyles[i].styles[tag] = style
			} else {
				delete(theme.mediaStyles[i].styles, tag)
			}
			break
		}
	}

	if style != nil {
		theme.mediaStyles = append(theme.mediaStyles, mediaStyle{
			MediaStyleParams: params,
			styles:           map[string]ViewStyle{tag: style},
		})
		theme.sortMediaStyles()
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

func (theme *theme) StyleTags() []string {
	keys := make([]string, 0, len(theme.styles)*2)

	appendTag := func(k string) {
		n := sort.SearchStrings(keys, k)
		if n >= len(keys) {
			keys = append(keys, k)
		} else if keys[n] != k {
			if n == 0 {
				keys = append([]string{k}, keys...)
			} else {
				keys = append(keys[:n+1], keys[n:]...)
				keys[n] = k
			}
		}
	}

	for k := range theme.styles {
		if index := strings.IndexRune(k, ':'); index < 0 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for k := range theme.styles {
		if index := strings.IndexRune(k, ':'); index > 0 {
			appendTag(k[:index])
		}
	}

	for _, media := range theme.mediaStyles {
		for k := range media.styles {
			index := strings.IndexRune(k, ':')
			if index > 0 {
				appendTag(k[:index])
			} else if index < 0 {
				appendTag(k)
			}
		}
	}
	return keys
}

func (theme *theme) MediaStyles(tag string) []struct {
	Selectors string
	Params    MediaStyleParams
} {
	result := []struct {
		Selectors string
		Params    MediaStyleParams
	}{}

	prefix := tag + ":"
	prefixLen := len(prefix)
	for themeTag := range theme.styles {
		if strings.HasPrefix(themeTag, prefix) {
			result = append(result, struct {
				Selectors string
				Params    MediaStyleParams
			}{
				Selectors: themeTag[prefixLen:],
				Params:    MediaStyleParams{},
			})
		}
	}

	for _, media := range theme.mediaStyles {
		if _, ok := media.styles[tag]; ok {
			result = append(result, struct {
				Selectors string
				Params    MediaStyleParams
			}{
				Selectors: "",
				Params:    media.MediaStyleParams,
			})
		}
		for themeTag := range media.styles {
			if strings.HasPrefix(themeTag, prefix) {
				result = append(result, struct {
					Selectors string
					Params    MediaStyleParams
				}{
					Selectors: themeTag[prefixLen:],
					Params:    media.MediaStyleParams,
				})
			}
		}
	}
	return result
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
			if anotherMedia.MinHeight == media.MinHeight &&
				anotherMedia.MaxHeight == media.MaxHeight &&
				anotherMedia.MinWidth == media.MinWidth &&
				anotherMedia.MaxWidth == media.MaxWidth &&
				anotherMedia.Orientation == media.Orientation {
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
	builder.init(16)

	styleList := func(styles map[string]ViewStyle) []string {
		ruiStyles := []string{}
		customStyles := []string{}
		for tag := range styles {
			if strings.HasPrefix(tag, "rui") {
				ruiStyles = append(ruiStyles, tag)
			} else {
				customStyles = append(customStyles, tag)
			}
		}
		sort.Strings(ruiStyles)
		sort.Strings(customStyles)
		return append(ruiStyles, customStyles...)
	}

	for _, tag := range styleList(theme.styles) {
		if style := theme.styles[tag]; style != nil {
			builder.startStyle(tag)
			style.cssViewStyle(&builder, session)
			builder.endStyle()
		}
	}

	for _, media := range theme.mediaStyles {
		builder.startMedia(media.cssText())
		for _, tag := range styleList(media.styles) {
			if style := media.styles[tag]; style != nil {
				builder.startStyle(tag)
				style.cssViewStyle(&builder, session)
				builder.endStyle()
			}
		}
		builder.endMedia()
	}

	return builder.finish()
}

func (theme *theme) addText(themeText string) bool {
	if theme.constants == nil {
		theme.init()
	}

	data, err := ParseDataText(themeText)
	if err != nil {
		ErrorLog(err.Error())
		return false
	} else if !data.IsObject() || data.Tag() != "theme" {
		return false
	}

	objToStyle := func(obj DataObject) ViewStyle {
		params := Params{}
		for node := range obj.Properties() {
			switch node.Type() {
			case ArrayNode:
				params[PropertyName(node.Tag())] = node.ArrayElements()

			case ObjectNode:
				params[PropertyName(node.Tag())] = node.Object()

			default:
				params[PropertyName(node.Tag())] = node.Text()
			}
		}
		return NewViewStyle(params)
	}

	for d := range data.Properties() {
		switch tag := d.Tag(); tag {
		case "constants":
			if d.Type() == ObjectNode {
				if obj := d.Object(); obj != nil {
					for prop := range obj.Properties() {
						if prop.Type() == TextNode {
							theme.constants[prop.Tag()] = prop.Text()
						}
					}
				}
			}

		case "constants:touch":
			if d.Type() == ObjectNode {
				if obj := d.Object(); obj != nil {
					for prop := range obj.Properties() {
						if prop.Type() == TextNode {
							theme.touchConstants[prop.Tag()] = prop.Text()
						}
					}
				}
			}

		case "colors":
			if d.Type() == ObjectNode {
				if obj := d.Object(); obj != nil {
					for prop := range obj.Properties() {
						if prop.Type() == TextNode {
							theme.colors[prop.Tag()] = prop.Text()
						}
					}
				}
			}

		case "colors:dark":
			if d.Type() == ObjectNode {
				if obj := d.Object(); obj != nil {
					for prop := range obj.Properties() {
						if prop.Type() == TextNode {
							theme.darkColors[prop.Tag()] = prop.Text()
						}
					}
				}
			}

		case "images":
			if d.Type() == ObjectNode {
				if obj := d.Object(); obj != nil {
					for prop := range obj.Properties() {
						if prop.Type() == TextNode {
							theme.images[prop.Tag()] = prop.Text()
						}
					}
				}
			}

		case "images:dark":
			if d.Type() == ObjectNode {
				if obj := d.Object(); obj != nil {
					for prop := range obj.Properties() {
						if prop.Type() == TextNode {
							theme.darkImages[prop.Tag()] = prop.Text()
						}
					}
				}
			}

		case "styles":
			if d.Type() == ArrayNode {
				arraySize := d.ArraySize()
				for k := range arraySize {
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
					for k := range arraySize {
						if element := d.ArrayElement(k); element != nil && element.IsObject() {
							if obj := element.Object(); obj != nil {
								rule.styles[obj.Tag()] = objToStyle(obj)
							}
						}
					}
					theme.mediaStyles = append(theme.mediaStyles, rule)
				}
			}
		}
	}

	theme.sortMediaStyles()
	return true
}

func (theme *theme) sortMediaStyles() {
	if len(theme.mediaStyles) > 1 {
		sort.SliceStable(theme.mediaStyles, func(i, j int) bool {
			if theme.mediaStyles[i].Orientation != theme.mediaStyles[j].Orientation {
				return theme.mediaStyles[i].Orientation < theme.mediaStyles[j].Orientation
			}
			if theme.mediaStyles[i].MinWidth != theme.mediaStyles[j].MinWidth {
				return theme.mediaStyles[i].MinWidth < theme.mediaStyles[j].MinWidth
			}
			if theme.mediaStyles[i].MinHeight != theme.mediaStyles[j].MinHeight {
				return theme.mediaStyles[i].MinHeight < theme.mediaStyles[j].MinHeight
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
		if isQuotesNeeded(text) {
			buffer.WriteRune('"')
			buffer.WriteString(replaceEscapeSymbols(text))
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

	writeStyles := func(orientation, maxWidth, maxHeight int, styles map[string]ViewStyle) bool {
		count := len(styles)
		if count == 0 {
			return false
		}

		tags := make([]string, 0, count)
		for name := range styles {
			tags = append(tags, name)
		}
		sort.Strings(tags)

		buffer.WriteString("\tstyles")
		switch orientation {
		case PortraitMedia:
			buffer.WriteString(":portrait")

		case LandscapeMedia:
			buffer.WriteString(":landscape")
		}
		if maxWidth > 0 {
			fmt.Fprintf(buffer, ":width%d", maxWidth)
		}
		if maxHeight > 0 {
			fmt.Fprintf(buffer, ":height%d", maxHeight)
		}
		buffer.WriteString(" = [\n")

		for _, tag := range tags {
			if style, ok := styles[tag]; ok && len(style.AllTags()) > 0 {
				buffer.WriteString("\t\t")
				writeViewStyle(tag, style, buffer, "\t\t", nil)
				buffer.WriteString(",\n")
			}
		}
		buffer.WriteString("\t],\n")
		return true
	}

	writeStyles(0, 0, 0, theme.styles)
	for _, media := range theme.mediaStyles {
		//writeStyles(media.orientation, media.maxWidth, media.maxHeight, media.styles)
		if count := len(media.styles); count > 0 {

			tags := make([]string, 0, count)
			for name := range media.styles {
				tags = append(tags, name)
			}
			sort.Strings(tags)

			buffer.WriteString("\tstyles")
			switch media.Orientation {
			case PortraitMedia:
				buffer.WriteString(":portrait")

			case LandscapeMedia:
				buffer.WriteString(":landscape")
			}

			if media.MinWidth > 0 {
				fmt.Fprintf(buffer, ":width%d-", media.MinWidth)
				if media.MaxWidth > 0 {
					buffer.WriteString(strconv.Itoa(media.MaxWidth))
				}
			} else if media.MaxWidth > 0 {
				fmt.Fprintf(buffer, ":width%d", media.MaxWidth)
			}

			if media.MinHeight > 0 {
				fmt.Fprintf(buffer, ":height%d-", media.MinHeight)
				if media.MaxHeight > 0 {
					buffer.WriteString(strconv.Itoa(media.MaxHeight))
				}
			} else if media.MaxHeight > 0 {
				fmt.Fprintf(buffer, ":height%d", media.MaxHeight)
			}

			buffer.WriteString(" = [\n")

			for _, tag := range tags {
				if style, ok := media.styles[tag]; ok && len(style.AllTags()) > 0 {
					buffer.WriteString("\t\t")
					writeViewStyle(tag, style, buffer, "\t\t", nil)
					buffer.WriteString(",\n")
				}
			}
			buffer.WriteString("\t],\n")
		}
	}

	buffer.WriteString("}\n")
	return buffer.String()
}
