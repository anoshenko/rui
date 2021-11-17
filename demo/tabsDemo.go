package main

import (
	"fmt"

	"github.com/anoshenko/rui"
)

const tabsDemoText = `
GridLayout {
	style = demoPage,
	content = [
		TabsLayout { id = tabsLayout, width = 100%, height = 100%, tabs = top, tab-close-button = true,
			content = [
				View { width = 300px, height = 200px, background-color = #FFFF0000, title = "Red tab", icon = red_icon.svg },
				View { width = 400px, height = 250px, background-color = #FF00FF00, title = "Green tab", icon = green_icon.svg },
				View { width = 100px, height = 400px, background-color = #FF0000FF, title = "Blue tab", icon = blue_icon.svg },
				View { width = 300px, height = 200px, background-color = #FF000000, title = "Black tab", icon = black_icon.svg },
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Tabs location" },
						DropDownList { row = 0, column = 1, id = tabsTypeList, 
							items = ["top", "bottom", "left", "right", "left list", "right list", "hidden"]
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

	rui.Set(view, "tabsLayout", rui.TabCloseEvent, func(index int) {
		rui.ShowMessage("", fmt.Sprintf(`The close button of the tab "%d" was clicked`, index), session)
	})
	return view
}
