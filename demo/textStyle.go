package main

import "github.com/anoshenko/rui"

const textStyleDemoText = `
GridLayout {
	style = demoPage,
	content = [
		GridLayout {
			width = 100%, height = 100%, cell-vertical-align = center, cell-horizontal-align = center,
			content = [
				TextView {
					id = textStyleText, padding = 16px, max-width = 80%,
					border = _{ style = solid, width = 1px, color = darkgray },
					text = "Twenty years from now you will be more disappointed by the things that you didn't do than by the ones you did do. So throw off the bowlines. Sail away from the safe harbor. Catch the trade winds in your sails. Explore. Dream. Discover."
				}
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Font name" },
						DropDownList { row = 0, column = 1, id = textStyleFont, current = 0, items = ["default", "serif", "sans-serif", "\"Courier new\",  monospace", "cursive", "fantasy"]},
						TextView { row = 1, text = "Text size" },
						DropDownList { row = 1, column = 1, id = textStyleSize, current = 0, items = ["1em", "14pt", "12px", "1.5em"]},
						TextView { row = 2, text = "Text color" },
						ColorPicker { row = 2, column = 1, id = textStyleColor },
						TextView { row = 3, text = "Text weight" },
						DropDownList { row = 3, column = 1, id = textStyleWeight, current = 0, items = ["default", "thin", "extra-light", "light", "normal", "medium", "semi-bold", "bold", "extra-bold", "black"]},
						Checkbox { row = 4, column = 0:1, id = textStyleItalic, content = "Italic" },
						Checkbox { row = 5, column = 0:1, id = textStyleSmallCaps, content = "Small-caps" },
						Checkbox { row = 6, column = 0:1, id = textStyleStrikethrough, content = "Strikethrough" },
						Checkbox { row = 7, column = 0:1, id = textStyleOverline, content = "Overline" },
						Checkbox { row = 8, column = 0:1, id = textStyleUnderline, content = "Underline" },
						TextView { row = 9, text = "Line style" },
						DropDownList { row = 9, column = 1, id = textStyleLineStyle, current = 0, items = ["default", "solid", "dashed", "dotted", "double", "wavy"]},
						TextView { row = 10, text = "Line thickness" },
						DropDownList { row = 10, column = 1, id = textStyleLineThickness, current = 0, items = ["default", "1px", "1.5px", "2px", "3px", "4px"]},
						TextView { row = 11, text = "Line color" },
						ColorPicker { row = 11, column = 1, id = textStyleLineColor },
						TextView { row = 12, text = "Shadow" },
						DropDownList { row = 12, column = 1, id = textStyleShadow, current = 0, items = ["none", "gray, (x, y)=(1px, 1px), blur=0", "blue, (x, y)=(-2px, -2px), blur=1", "green, (x, y)=(0px, 0px), blur=3px"]},
					]
				}
			]
		}
	]
}
`

func createTextStyleDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, textStyleDemoText)
	if view == nil {
		return nil
	}

	rui.SetChangeListener(view, "textStyleFont", rui.Current, func(v rui.View, tag string) {
		fonts := []string{"", "serif", "sans-serif", "\"Courier new\", monospace", "cursive", "fantasy"}
		if number := rui.GetCurrent(v, ""); number > 0 && number < len(fonts) {
			rui.Set(view, "textStyleText", rui.FontName, fonts[number])
		} else {
			rui.Set(view, "textStyleText", rui.FontName, nil)
		}
	})
	/*
		rui.Set(view, "textStyleFont", rui.DropDownEvent, func(number int) {
			fonts := []string{"", "serif", "sans-serif", "\"Courier new\", monospace", "cursive", "fantasy"}
			if number > 0 && number < len(fonts) {
				rui.Set(view, "textStyleText", rui.FontName, fonts[number])
			} else {
				rui.Set(view, "textStyleText", rui.FontName, nil)
			}
		})
	*/

	rui.Set(view, "textStyleSize", rui.DropDownEvent, func(number int) {
		sizes := []string{"1em", "14pt", "12px", "1.5em"}
		if number >= 0 && number < len(sizes) {
			rui.Set(view, "textStyleText", rui.TextSize, sizes[number])
		}
	})

	rui.Set(view, "textStyleColor", rui.ColorChangedEvent, func(color rui.Color) {
		rui.Set(view, "textStyleText", rui.TextColor, color)
	})

	rui.Set(view, "textStyleWeight", rui.DropDownEvent, func(number int) {
		rui.Set(view, "textStyleText", rui.TextWeight, number)
	})

	rui.Set(view, "textStyleItalic", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "textStyleText", rui.Italic, state)
	})

	rui.Set(view, "textStyleSmallCaps", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "textStyleText", rui.SmallCaps, state)
	})

	rui.Set(view, "textStyleStrikethrough", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "textStyleText", rui.Strikethrough, state)
	})

	rui.Set(view, "textStyleOverline", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "textStyleText", rui.Overline, state)
	})

	rui.Set(view, "textStyleUnderline", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "textStyleText", rui.Underline, state)
	})

	rui.Set(view, "textStyleLineStyle", rui.DropDownEvent, func(number int) {
		styles := []string{"inherit", "solid", "dashed", "dotted", "double", "wavy"}
		if number > 0 && number < len(styles) {
			rui.Set(view, "textStyleText", rui.TextLineStyle, styles[number])
		} else {
			rui.Set(view, "textStyleText", rui.TextLineStyle, nil)
		}
	})

	rui.Set(view, "textStyleLineThickness", rui.DropDownEvent, func(number int) {
		sizes := []string{"", "1px", "1.5px", "2px", "3px", "4px"}
		if number > 0 && number < len(sizes) {
			rui.Set(view, "textStyleText", rui.TextLineThickness, sizes[number])
		} else {
			rui.Set(view, "textStyleText", rui.TextLineThickness, nil)
		}
	})

	rui.Set(view, "textStyleLineColor", rui.ColorChangedEvent, func(color rui.Color) {
		rui.Set(view, "textStyleText", rui.TextLineColor, color)
	})

	rui.Set(view, "textStyleShadow", rui.DropDownEvent, func(number int) {
		switch number {
		case 0:
			rui.Set(view, "textStyleText", rui.TextShadow, nil)

		case 1:
			rui.Set(view, "textStyleText", rui.TextShadow, rui.NewTextShadow(rui.Px(1), rui.Px(1), rui.Px(0), rui.Gray))

		case 2:
			rui.Set(view, "textStyleText", rui.TextShadow, rui.NewTextShadow(rui.Px(-2), rui.Px(-2), rui.Px(1), rui.Blue))

		case 3:
			rui.Set(view, "textStyleText", rui.TextShadow, rui.NewTextShadow(rui.Px(0), rui.Px(0), rui.Px(3), rui.Green))
		}
	})

	return view
}
