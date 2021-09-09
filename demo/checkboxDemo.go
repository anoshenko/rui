package main

import "github.com/anoshenko/rui"

const checkboxDemoText = `
GridLayout {
	style = demoPage,
	content = [
		GridLayout {
			width = 100%, height = 100%, cell-vertical-align = center, cell-horizontal-align = center,
			content = [
				GridLayout {
					width = 250px, height = 80px,
					border = _{ style = solid, width = 1px, color = gray },
					content = [
						Checkbox {
							id = checkbox, width = 100%, height = 100%,
							content = "Checkbox content"
						}
					]
				}
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Vertical align" },
						DropDownList { row = 0, column = 1, id = checkboxVAlign, current = 0, items = ["top", "bottom", "center", "stretch"]},
						TextView { row = 1, text = "Horizontal align" },
						DropDownList { row = 1, column = 1, id = checkboxHAlign, current = 0, items = ["left", "right", "center", "stretch"]},
						TextView { row = 2, text = "Checkbox vertical align" },
						DropDownList { row = 2, column = 1, id = checkboxBoxVAlign, current = 0, items = ["top", "bottom", "center"]},
						TextView { row = 3, text = "Checkbox horizontal align" },
						DropDownList { row = 3, column = 1, id = checkboxBoxHAlign, current = 0, items = ["left", "right", "center"]},
					]
				}
			]
		}
	]
}
`

func createCheckboxDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, checkboxDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "checkboxVAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "checkbox", rui.VerticalAlign, number)
	})

	rui.Set(view, "checkboxHAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "checkbox", rui.HorizontalAlign, number)
	})

	rui.Set(view, "checkboxBoxVAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "checkbox", rui.CheckboxVerticalAlign, number)
	})

	rui.Set(view, "checkboxBoxHAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "checkbox", rui.CheckboxHorizontalAlign, number)
	})

	return view
}
