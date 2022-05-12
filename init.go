package rui

func init() {
	if theme, ok := CreateThemeFromText(defaultThemeText); ok {
		defaultTheme = theme
	}
}
