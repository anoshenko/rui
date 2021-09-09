package main

import (
	"github.com/anoshenko/rui"
)

const tabsDemoText = `
GridLayout {
	style = demoPage,
	content = [
		TabsLayout { id = tabsLayout, width = 100%, height = 100%, tabs = top, 
			content = [
				View { width = 300px, height = 200px, background-color = #FFFF0000, title = "Red tab"},
				View { width = 400px, height = 250px, background-color = #FF00FF00, title = "Green tab"},
				View { width = 100px, height = 400px, background-color = #FF0000FF, title = "Blue tab"},
				View { width = 300px, height = 200px, background-color = #FF000000, title = "Black tab"},
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Tabs location" },
						DropDownList { row = 0, column = 1, id = tabsTypeList, current = 1,
							items = ["hidden", "top", "bottom", "left", "right", "left list", "right list"]
						}
					]
				}
			]
		}
	]
}
`

func createTabsDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, tabsDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "tabsTypeList", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "tabsLayout", rui.Tabs, number)
	})

	return view
}
