package rui

import (
	"fmt"
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
	theme := session.getCurrentTheme()
	for {
		result := theme.constant(tag, session.touchScreen)
		if result == "" {
			ErrorLogF(`"%v" constant not found`, tag)
			return "", false
		}

		if len(result) < 2 || !strings.ContainsRune(result, '@') {
			return result, true
		}

		for _, separator := range []string{",", " ", ":", ";", "|", "/"} {
			if strings.Contains(result, separator) {
				return session.resolveConstantsNext(result, tags)
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
	if !strings.ContainsRune(value, '@') {
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

func (session *sessionData) getCurrentTheme() Theme {
	if session.currentTheme != nil {
		return session.currentTheme
	}

	if session.customTheme != nil {
		session.currentTheme = NewTheme("")
		session.currentTheme.Append(defaultTheme)
		session.currentTheme.Append(session.customTheme)
		return session.currentTheme
	}

	return defaultTheme
}

// Color return the color with "tag" name or 0 if it is not exists
func (session *sessionData) Color(tag string) (Color, bool) {
	tags := []string{tag}
	theme := session.getCurrentTheme()
	for {
		result := theme.color(tag, session.darkTheme)
		if result == "" {
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
	theme := session.getCurrentTheme()
	for {
		result := theme.image(tag, session.darkTheme)
		if result == "" {
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

func (session *sessionData) checkboxImage(checked bool, accentColor Color) string {

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
		if accentColor != 0 {
			backgroundColor = accentColor
		} else if backgroundColor, ok = session.Color("ruiHighlightColor"); !ok {
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

func (session *sessionData) checkboxOffImage(accentColor Color) string {
	if accentColor != 0 {
		return session.checkboxImage(false, accentColor)
	}
	if session.checkboxOff == "" {
		session.checkboxOff = session.checkboxImage(false, accentColor)
	}
	return session.checkboxOff
}

func (session *sessionData) checkboxOnImage(accentColor Color) string {
	if accentColor != 0 {
		return session.checkboxImage(true, accentColor)
	}
	if session.checkboxOn == "" {
		session.checkboxOn = session.checkboxImage(true, accentColor)
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

func (session *sessionData) radiobuttonOnImage(accentColor Color) string {
	if session.radiobuttonOn == "" {
		var borderColor, backgroundColor Color
		var ok bool

		if accentColor != 0 {
			borderColor = accentColor
		} else if borderColor, ok = session.Color("ruiHighlightColor"); !ok {
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

	if len(session.languages) > 0 {
		return session.languages[0]
	}

	return "en"
}

func (session *sessionData) SetLanguage(lang string) {
	lang = strings.Trim(lang, " \t\n\r")
	if lang != session.language {
		session.language = lang

		if session.rootView != nil && session.bridge != nil {
			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)

			viewHTML(session.rootView, buffer, "")
			session.bridge.updateInnerHTML("ruiRootView", buffer.String())
		}
	}
}
