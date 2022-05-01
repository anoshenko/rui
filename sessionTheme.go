package rui

import (
	"fmt"
	"sort"
	"strings"
)

func (session *sessionData) DarkTheme() bool {
	return session.darkTheme
}

func (session *sessionData) TouchScreen() bool {
	return session.touchScreen
}

func (session *sessionData) PixelRatio() float64 {
	return session.pixelRatio
}

func (session *sessionData) TextDirection() int {
	return session.textDirection
}

func (session *sessionData) constant(tag string, prevTags []string) (string, bool) {
	tags := append(prevTags, tag)
	result := ""
	theme := session.getCurrentTheme()
	for {
		ok := false
		if session.touchScreen {
			if theme.touchConstants != nil {
				result, ok = theme.touchConstants[tag]
			}
		}

		if !ok {
			result, ok = theme.constants[tag]
		}

		if !ok {
			ErrorLogF(`"%v" constant not found`, tag)
			return "", false
		}

		if len(result) < 2 || !strings.ContainsRune(result, '@') {
			return result, true
		}

		for _, separator := range []string{",", " ", ":", ";", "|", "/"} {
			if strings.Contains(result, separator) {
				result, ok = session.resolveConstantsNext(result, tags)
				return result, ok
			}
		}

		if result[0] != '@' {
			return result, true
		}

		tag = result[1:]
		for _, t := range tags {
			if t == tag {
				ErrorLogF(`"%v" constant is cyclic`, tag)
				return "", false
			}
		}
		tags = append(tags, tag)
	}
}

func (session *sessionData) resolveConstants(value string) (string, bool) {
	return session.resolveConstantsNext(value, []string{})
}

func (session *sessionData) resolveConstantsNext(value string, prevTags []string) (string, bool) {
	if !strings.Contains(value, "@") {
		return value, true
	}

	separators := []rune{',', ' ', ':', ';', '|', '/'}
	sep := rune(0)
	index := -1
	for _, s := range separators {
		if i := strings.IndexRune(value, s); i >= 0 {
			if i < index || index < 0 {
				sep = s
				index = i
			}
		}
	}

	ok := true
	if index >= 0 {
		v1 := strings.Trim(value[:index], " \t\n\r")
		v2 := strings.Trim(value[index+1:], " \t\n\r")
		if len(v1) > 1 && v1[0] == '@' {
			if v1, ok = session.constant(v1[1:], prevTags); !ok {
				return value, false
			}
			if v, ok := session.resolveConstantsNext(v1, prevTags); ok {
				v1 = v
			} else {
				return v1 + string(sep) + v2, false
			}
		}

		if v, ok := session.resolveConstantsNext(v2, prevTags); ok {
			v2 = v
		}

		return v1 + string(sep) + v2, ok

	} else if value[0] == '@' {

		if value, ok = session.constant(value[1:], prevTags); ok {
			return session.resolveConstantsNext(value, prevTags)
		}
	}

	return value, false
}

func (session *sessionData) Constant(tag string) (string, bool) {
	return session.constant(tag, []string{})
}

func (session *sessionData) getCurrentTheme() *theme {
	if session.currentTheme != nil {
		return session.currentTheme
	}

	if session.customTheme != nil {
		session.currentTheme = new(theme)
		session.currentTheme.init()
		session.currentTheme.concat(defaultTheme)
		session.currentTheme.concat(session.customTheme)
		return session.currentTheme
	}

	return defaultTheme
}

// Color return the color with "tag" name or 0 if it is not exists
func (session *sessionData) Color(tag string) (Color, bool) {
	tags := []string{tag}
	result := ""
	theme := session.getCurrentTheme()
	for {
		ok := false
		if session.darkTheme {
			if theme.darkColors != nil {
				result, ok = theme.darkColors[tag]
			}
		}

		if !ok {
			if theme.colors != nil {
				result, ok = theme.colors[tag]
			}
		}

		if !ok {
			ErrorLogF(`"%v" color not found`, tag)
			return 0, false
		}

		if len(result) == 0 || result[0] != '@' {
			color, err := stringToColor(result)
			if err != nil {
				ErrorLogF(`invalid value "%v" of "%v" color constant (%s)`, result, tag, err.Error())
				return 0, false
			}
			return color, true
		}

		tag = result[1:]
		for _, t := range tags {
			if t == tag {
				ErrorLogF(`"%v" color is cyclic`, tag)
				return 0, false
			}
		}

		tags = append(tags, tag)
	}
}

func (session *sessionData) ImageConstant(tag string) (string, bool) {
	tags := []string{tag}
	result := ""
	theme := session.getCurrentTheme()
	for {
		ok := false
		if session.darkTheme {
			if theme.darkImages != nil {
				result, ok = theme.darkImages[tag]
			}
		}

		if !ok {
			if theme.images != nil {
				result, ok = theme.images[tag]
			}
		}

		if !ok {
			ErrorLogF(`"%v" image not found`, tag)
			return "", false
		}

		if result == "" || result[0] != '@' {
			return result, true
		}

		tag = result[1:]
		for _, t := range tags {
			if t == tag {
				ErrorLogF(`"%v" image is cyclic`, tag)
				return "", false
			}
		}

		tags = append(tags, tag)
	}
}

func (session *sessionData) SetCustomTheme(name string) bool {
	if name == "" {
		if session.customTheme == nil {
			return true
		}
	} else if theme, ok := resources.themes[name]; ok {
		session.customTheme = theme
		session.currentTheme = nil
	} else {
		return false
	}

	session.reload()
	return true
}

const checkImage = `<svg width="16" height="16" version="1.1" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg"><path d="m4 8 3 4 5-8" fill="none" stroke="#fff" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5"/></svg>`

func (session *sessionData) checkboxImage(checked bool) string {

	var borderColor, backgroundColor Color
	var ok bool

	if borderColor, ok = session.Color("ruiDisabledTextColor"); !ok {
		if session.darkTheme {
			borderColor = 0xFFA0A0A0
		} else {
			borderColor = 0xFF202020
		}
	}

	if checked {
		if backgroundColor, ok = session.Color("ruiHighlightColor"); !ok {
			backgroundColor = 0xFF1A74E8
		}
	} else if backgroundColor, ok = session.Color("ruiBackgroundColor"); !ok {
		if session.darkTheme {
			backgroundColor = 0xFF040404
		} else {
			backgroundColor = 0xFFF0F0F0
		}
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(`<div style="width: 18px; height: 18px; background-color: `)
	buffer.WriteString(backgroundColor.cssString())
	buffer.WriteString(`; border: 1px solid `)
	buffer.WriteString(borderColor.cssString())
	buffer.WriteString(`; border-radius: 4px;">`)
	if checked {
		buffer.WriteString(checkImage)
	}
	buffer.WriteString(`</div>`)

	return buffer.String()
}

func (session *sessionData) checkboxOffImage() string {
	if session.checkboxOff == "" {
		session.checkboxOff = session.checkboxImage(false)
	}
	return session.checkboxOff
}

func (session *sessionData) checkboxOnImage() string {
	if session.checkboxOn == "" {
		session.checkboxOn = session.checkboxImage(true)
	}
	return session.checkboxOn
}

func (session *sessionData) radiobuttonOffImage() string {
	if session.radiobuttonOff == "" {
		var borderColor, backgroundColor Color
		var ok bool

		if borderColor, ok = session.Color("ruiDisabledTextColor"); !ok {
			if session.darkTheme {
				borderColor = 0xFFA0A0A0
			} else {
				borderColor = 0xFF202020
			}
		}

		if backgroundColor, ok = session.Color("ruiBackgroundColor"); !ok {
			if session.darkTheme {
				backgroundColor = 0xFFA0A0A0
			} else {
				backgroundColor = 0xFF202020
			}
		}

		session.radiobuttonOff = fmt.Sprintf(`<div style="width: 16px; height: 16px; background-color: %s; border: 1px solid %s; border-radius: 8px;"></div>`,
			backgroundColor.cssString(), borderColor.cssString())
	}
	return session.radiobuttonOff
}

func (session *sessionData) radiobuttonOnImage() string {
	if session.radiobuttonOn == "" {
		var borderColor, backgroundColor Color
		var ok bool

		if borderColor, ok = session.Color("ruiHighlightColor"); !ok {
			borderColor = 0xFF1A74E8
		}

		if backgroundColor, ok = session.Color("ruiHighlightTextColor"); !ok {
			backgroundColor = 0xFFFFFFFF
		}

		session.radiobuttonOn = fmt.Sprintf(`<div style="width: 16px; height: 16px; display: grid; justify-items: center; align-items: center; background-color: %s; border: 2px solid %s; border-radius: 8px;"><div style="width: 8px; height: 8px; background-color: %s; border-radius: 4px;"></div></div>`,
			backgroundColor.cssString(), borderColor.cssString(), borderColor.cssString())
	}
	return session.radiobuttonOn
}

func (session *sessionData) UserAgent() string {
	return session.userAgent
}

func (session *sessionData) Language() string {
	if session.language != "" {
		return session.language
	}

	if session.languages != nil && len(session.languages) > 0 {
		return session.languages[0]
	}

	return "en"
}

func (session *sessionData) SetLanguage(lang string) {
	lang = strings.Trim(lang, " \t\n\r")
	if lang != session.language {
		session.language = lang

		if session.rootView != nil {
			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)

			buffer.WriteString(`document.getElementById('ruiRootView').innerHTML = '`)
			viewHTML(session.rootView, buffer)
			buffer.WriteString("';\nscanElementsSize();")

			session.runScript(buffer.String())
		}
	}
}

func (session *sessionData) ConstantTags() []string {
	theme := session.getCurrentTheme()

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

func (session *sessionData) ColorTags() []string {
	theme := session.getCurrentTheme()

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

func (session *sessionData) ImageConstantTags() []string {
	theme := session.getCurrentTheme()

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
