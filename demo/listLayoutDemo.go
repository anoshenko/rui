package main

import (
	"github.com/anoshenko/rui"
)

const listLayoutDemoText = `
GridLayout {
	style = demoPage,
	content = [
		ListLayout {
			id = listLayout, width = 100%, height = 100%, orientation = up-down,
			content = [
				GridLayout { width = 200px, height = 100px, content = ["View 1"], horizontal-align = center, vertical-align = center, 
					background-color = #FFAAAAAA, radius = 8px, padding = 8px, margin = 4px,
					border = _{ style = solid, width = 1px, color = black } 
				},
				GridLayout { width = 100px, height = 200px, content = ["View 2"], horizontal-align = center, vertical-align = center, 
					background-color = #FFB0B0B0, radius = 8px, padding = 8px, margin = 4px,
					border = _{ style = solid, width = 1px, color = black } 
				},
				GridLayout { width = 150px, height = 150px, content = ["View 3"], horizontal-align = center, vertical-align = center, 
					background-color = #FFBBBBBB, radius = 8px, padding = 8px, margin = 4px,
					border = _{ style = solid, width = 1px, color = black } 
				},
				GridLayout { width = 150px, height = 100px, content = ["View 4"], horizontal-align = center, vertical-align = center, 
					background-color = #FFC0C0C0, radius = 8px, padding = 8px, margin = 4px,
					border = _{ style = solid, width = 1px, color = black } 
				},
				GridLayout { width = 200px, height = 150px, content = ["View 5"], horizontal-align = center, vertical-align = center, 
					background-color = #FFCCCCCC, radius = 8px, padding = 8px, margin = 4px,
					border = _{ style = solid, width = 1px, color = black } 
				},
				GridLayout { width = 100px, height = 100px, content = ["View 6"], horizontal-align = center, vertical-align = center, 
					background-color = #FFDDDDDD, radius = 8px, padding = 8px, margin = 4px,
					border = _{ style = solid, width = 1px, color = black } 
				},
				GridLayout { width = 150px, height = 200px, content = ["View 7"], horizontal-align = center, vertical-align = center, 
					background-color = #FFEEEEEE, radius = 8px, padding = 8px, margin = 4px,
					border = _{ style = solid, width = 1px, color = black } 
				},
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Orientation" },
						DropDownList { row = 0, column = 1, id = listOrientation, current = 0,
							items = ["up-down", "start-to-end", "bottom-up", "end-to-start"]
						},
						TextView { row = 1, text = "Wrap" },
						DropDownList { row = 1, column = 1, id = listWrap, current = 0, items = ["off", "on", "reverse"]},
						TextView { row = 2, text = "Vertical align" },
						DropDownList { row = 2, column = 1, id = listVAlign, current = 0, items = ["top", "bottom", "center", "stretch"]},
						TextView { row = 3, text = "Horizontal align" },
						DropDownList { row = 3, column = 1, id = listHAlign, current = 0, items = ["left", "right", "center", "stretch"]},
					]
				},
			]
		}
	]		
}`

func createListLayoutDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, listLayoutDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "listOrientation", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listLayout", rui.Orientation, number)
	})

	rui.Set(view, "listWrap", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listLayout", rui.ListWrap, number)
	})

	rui.Set(view, "listHAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listLayout", rui.HorizontalAlign, number)
	})

	rui.Set(view, "listVAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listLayout", rui.VerticalAlign, number)
	})

	return view
}
