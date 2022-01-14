package main

import "github.com/anoshenko/rui"

const tableViewDemoText = `
GridLayout {
	style = demoPage,
	content = [
		ColumnLayout {
			width = 100%, height = 100%,
			content = TableView {
				id = demoTableView1, margin = 24px, 
			}
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Cell gap" },
						DropDownList { row = 0, column = 1, id = tableCellGap, current = 0, items = ["0", "2px"]},
						TextView { row = 1, text = "Table border" },
						DropDownList { row = 1, column = 1, id = tableBorder, current = 0, items = ["none", "solid black 1px", "4 colors"]},
						TextView { row = 2, text = "Cell border" },
						DropDownList { row = 2, column = 1, id = tableCellBorder, current = 0, items = ["none", "solid black 1px", "4 colors"]},
						TextView { row = 3, text = "Cell padding" },
						DropDownList { row = 3, column = 1, id = tableCellPadding, current = 0, items = ["default", "4px", "8px, 16px, 8px, 16px"]},
						TextView { row = 4, text = "Head style" },
						DropDownList { row = 4, column = 1, id = tableHeadStyle, current = 0, items = ["none", "tableHead1", "rui.Params"]},
						TextView { row = 5, text = "Foot style" },
						DropDownList { row = 5, column = 1, id = tableFootStyle, current = 0, items = ["none", "tableFoot1", "rui.Params"]},
						Checkbox { row = 6, column = 0:1, id = tableRowStyle, content = "Row style" },
						Checkbox { row = 7, column = 0:1, id = tableColumnStyle, content = "Column style" },
						TextView { row = 8, text = "Selection mode" },
						DropDownList { row = 8, column = 1, id = tableSelectionMode, current = 0, items = ["none", "cell", "row"]},
						Checkbox { row = 9, column = 0:1, id = tableDisableHead, content = "Disable head selection" },
						Checkbox { row = 10, column = 0:1, id = tableDisableFoot, content = "Disable foot selection" },
					]
				}
			]
		}
	]
}
`

type demoTableAllowSelection struct {
	index []int
}

func (allow *demoTableAllowSelection) AllowCellSelection(row, column int) bool {
	return allow.AllowRowSelection(row)
}

func (allow *demoTableAllowSelection) AllowRowSelection(row int) bool {
	if allow.index != nil {
		for _, index := range allow.index {
			if index == row {
				return false
			}
		}
	}
	return true
}

func newDemoTableAllowSelection(index []int) *demoTableAllowSelection {
	result := new(demoTableAllowSelection)
	result.index = index
	return result
}

func createTableViewDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, tableViewDemoText)
	if view == nil {
		return nil
	}

	content := [][]interface{}{
		{"Cell content", "Cell value", rui.HorizontalTableJoin{}},
		{rui.VerticalTableJoin{}, "Type", "Value"},
		{"Text", "string", "Text"},
		{"Number", "int", 10},
		{rui.VerticalTableJoin{}, "float", 10.95},
		{"Boolean", "true", true},
		{rui.VerticalTableJoin{}, "false", false},
		{"Color", "red", rui.Red},
		{rui.VerticalTableJoin{}, "green", rui.Green},
		{rui.VerticalTableJoin{}, "yellow", rui.Yellow},
		{"View", "Button", rui.NewButton(session, rui.Params{
			rui.Content: "OK",
		})},
		{"Foot line", rui.HorizontalTableJoin{}, rui.HorizontalTableJoin{}},
	}

	rui.SetParams(view, "demoTableView1", rui.Params{
		rui.Content:    content,
		rui.HeadHeight: 2,
		rui.FootHeight: 1,
	})

	setBorder := func(borderTag string, number int) {
		switch number {
		case 1:
			rui.Set(view, "demoTableView1", borderTag, rui.NewBorder(rui.Params{
				rui.Style:    rui.SolidLine,
				rui.ColorTag: rui.Black,
				rui.Width:    rui.Px(1),
			}))

		case 2:
			rui.Set(view, "demoTableView1", borderTag, rui.NewBorder(rui.Params{
				rui.Style:       rui.SolidLine,
				rui.LeftColor:   rui.Blue,
				rui.RightColor:  rui.Magenta,
				rui.TopColor:    rui.Red,
				rui.BottomColor: rui.Green,
				rui.Width:       rui.Px(2),
			}))

		default:
			rui.Set(view, "demoTableView1", borderTag, nil)
		}
	}

	rui.Set(view, "tableSelectionMode", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "demoTableView1", rui.SelectionMode, number)
		switch rui.GetCurrent(view, "tableSelectionMode") {
		case rui.CellSelection:
			// TODO

		case rui.RowSelection:
			// TODO
		}
	})

	rui.Set(view, "tableDisableHead", rui.CheckboxChangedEvent, func(checked bool) {
		if checked {
			if rui.IsCheckboxChecked(view, "tableDisableFoot") {
				rui.Set(view, "demoTableView1", rui.AllowSelection, newDemoTableAllowSelection([]int{0, 1, 11}))
			} else {
				rui.Set(view, "demoTableView1", rui.AllowSelection, newDemoTableAllowSelection([]int{0, 1}))
			}
		} else {
			if rui.IsCheckboxChecked(view, "tableDisableFoot") {
				rui.Set(view, "demoTableView1", rui.AllowSelection, newDemoTableAllowSelection([]int{11}))
			} else {
				rui.Set(view, "demoTableView1", rui.AllowSelection, nil)
			}
		}
	})

	rui.Set(view, "tableDisableFoot", rui.CheckboxChangedEvent, func(checked bool) {
		if checked {
			if rui.IsCheckboxChecked(view, "tableDisableHead") {
				rui.Set(view, "demoTableView1", rui.AllowSelection, newDemoTableAllowSelection([]int{0, 1, 11}))
			} else {
				rui.Set(view, "demoTableView1", rui.AllowSelection, newDemoTableAllowSelection([]int{11}))
			}
		} else {
			if rui.IsCheckboxChecked(view, "tableDisableHead") {
				rui.Set(view, "demoTableView1", rui.AllowSelection, newDemoTableAllowSelection([]int{0, 1}))
			} else {
				rui.Set(view, "demoTableView1", rui.AllowSelection, nil)
			}
		}
	})

	rui.Set(view, "tableCellGap", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		if number == 0 {
			rui.Set(view, "demoTableView1", rui.Gap, rui.Px(0))
		} else {
			rui.Set(view, "demoTableView1", rui.Gap, rui.Px(2))
		}
	})

	rui.Set(view, "tableBorder", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		setBorder(rui.Border, number)
	})

	rui.Set(view, "tableCellBorder", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		setBorder(rui.CellBorder, number)
	})

	rui.Set(view, "tableHeadStyle", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		switch number {
		case 1:
			rui.Set(view, "demoTableView1", rui.HeadStyle, "tableHead1")

		case 2:
			rui.Set(view, "demoTableView1", rui.HeadStyle, rui.Params{
				rui.CellBorder: rui.NewBorder(rui.Params{
					rui.Style:    rui.SolidLine,
					rui.ColorTag: rui.Green,
					rui.Width:    "2px",
				}),
				rui.CellPadding:     "8px",
				rui.BackgroundColor: rui.LightGrey,
			})

		default:
			rui.Set(view, "demoTableView1", rui.HeadStyle, nil)
		}
	})

	rui.Set(view, "tableFootStyle", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		switch number {
		case 1:
			rui.Set(view, "demoTableView1", rui.FootStyle, "tableFoot1")

		case 2:
			rui.Set(view, "demoTableView1", rui.FootStyle, rui.Params{
				rui.Border: rui.NewBorder(rui.Params{
					rui.Style:    rui.SolidLine,
					rui.ColorTag: rui.Black,
					rui.Width:    "2px",
				}),
				rui.CellPadding:     "4px",
				rui.BackgroundColor: rui.LightYellow,
			})

		default:
			rui.Set(view, "demoTableView1", rui.FootStyle, nil)
		}
	})

	rui.Set(view, "tableCellPadding", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		switch number {
		case 1:
			rui.Set(view, "demoTableView1", rui.CellPadding, rui.Px(4))

		case 2:
			rui.Set(view, "demoTableView1", rui.CellPadding, rui.Bounds{
				Left:   rui.Px(16),
				Right:  rui.Px(16),
				Top:    rui.Px(8),
				Bottom: rui.Px(8),
			})

		default:
			rui.Set(view, "demoTableView1", rui.CellPadding, nil)
		}
	})

	rui.Set(view, "tableRowStyle", rui.CheckboxChangedEvent, func(checked bool) {
		if checked {
			rui.Set(view, "demoTableView1", rui.RowStyle, []rui.Params{
				{rui.BackgroundColor: 0xffeaece5},
				{rui.BackgroundColor: 0xfff0efef},
				{rui.BackgroundColor: 0xffe0e2e4},
				{rui.BackgroundColor: 0xffbccad6},
				{rui.BackgroundColor: 0xffcfe0e8},
				{rui.BackgroundColor: 0xffb7d7e8},
				{rui.BackgroundColor: 0xffdaebe8},
				{rui.BackgroundColor: 0xfff1e3dd},
				{rui.BackgroundColor: 0xfffbefcc},
				{rui.BackgroundColor: 0xfffff2df},
				{rui.BackgroundColor: 0xffffeead},
				{rui.BackgroundColor: 0xfff2e394},
			})
		} else {
			rui.Set(view, "demoTableView1", rui.RowStyle, nil)
		}
	})

	rui.Set(view, "tableColumnStyle", rui.CheckboxChangedEvent, func(checked bool) {
		if checked {
			rui.Set(view, "demoTableView1", rui.ColumnStyle, []rui.Params{
				{rui.BackgroundColor: 0xffeaece5},
				{rui.BackgroundColor: 0xffdaebe8},
				{rui.BackgroundColor: 0xfff2e394},
			})
		} else {
			rui.Set(view, "demoTableView1", rui.RowStyle, nil)
		}
	})

	return view
}
