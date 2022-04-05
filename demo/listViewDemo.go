package main

import (
	"github.com/anoshenko/rui"
)

const listViewDemoText = `
GridLayout {
	style = demoPage,
	content = [
		ListView {
			id = listView, width = 100%, height = 100%, orientation = vertical,
			items = ["Item 1", "Item 2", "Item 3", "Item 4", "Item 5", "Item 6", "Item 7", "Item 8", "Item 9", "Item 10", "Item 11", "Item 12", "Item 13", "Item 14", "Item 15", "Item 16", "Item 17", "Item 18"]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Orientation" },
						DropDownList { row = 0, column = 1, id = listViewOrientation, current = 0, items = ["vertical", "horizontal", "bottom up", "end to start"]},
						TextView { row = 1, text = "Wrap" },
						DropDownList { row = 1, column = 1, id = listWrap, current = 0, items = ["off", "on", "reverse"]},
						TextView { row = 2, text = "Item height" },
						DropDownList { row = 2, column = 1, id = listItemHeight, current = 0, items = ["auto", "25%", "50px"]},
						TextView { row = 3, text = "Item width" },
						DropDownList { row = 3, column = 1, id = listItemWidth, current = 0, items = ["auto", "25%", "200px"]},
						TextView { row = 4, text = "Item vertical align" },
						DropDownList { row = 4, column = 1, id = listItemVAlign, current = 0, items = ["top", "bottom", "center"]},
						TextView { row = 5, text = "Item horizontal align" },
						DropDownList { row = 5, column = 1, id = listItemHAlign, current = 0, items = ["left", "right", "center"]},
						TextView { row = 6, text = "Checkbox" },
						DropDownList { row = 6, column = 1, id = listCheckbox, current = 0, items = ["none", "single", "multiple"]},
						TextView { row = 7, text = "Checkbox vertical align" },
						DropDownList { row = 7, column = 1, id = listCheckboxVAlign, current = 0, items = ["top", "bottom", "center"]},
						TextView { row = 8, text = "Checkbox horizontal align" },
						DropDownList { row = 8, column = 1, id = listCheckboxHAlign, current = 0, items = ["left", "right", "center"]},
						Button { row = 9, column = 0:1, id = listSetChecked, content = "set checked 1,4,8" }
					]
				}
			]
		}
	]
}`

func createListViewDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, listViewDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "listViewOrientation", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listView", rui.Orientation, number)
	})

	rui.Set(view, "listWrap", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listView", rui.Wrap, number)
	})

	setItemSize := func(tag string, number int, values []rui.SizeUnit) {
		if number >= 0 && number < len(values) {
			rui.Set(view, "listView", tag, values[number])
		}
	}

	rui.Set(view, "listItemWidth", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		setItemSize(rui.ItemWidth, number, []rui.SizeUnit{rui.AutoSize(), rui.Percent(25), rui.Px(200)})
	})

	rui.Set(view, "listItemHeight", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		setItemSize(rui.ItemHeight, number, []rui.SizeUnit{rui.AutoSize(), rui.Percent(25), rui.Px(50)})
	})

	rui.Set(view, "listItemVAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listView", rui.VerticalAlign, number)
	})

	rui.Set(view, "listItemHAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listView", rui.HorizontalAlign, number)
	})

	rui.Set(view, "listCheckbox", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listView", rui.ItemCheckbox, number)
	})

	rui.Set(view, "listCheckboxVAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listView", rui.CheckboxVerticalAlign, number)
	})

	rui.Set(view, "listCheckboxHAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "listView", rui.CheckboxHorizontalAlign, number)
	})

	rui.Set(view, "listSetChecked", rui.ClickEvent, func() {
		rui.Set(view, "listView", rui.Checked, "1, 4, 8")
	})

	return view
}
