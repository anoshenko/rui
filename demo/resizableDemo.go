package main

import "github.com/anoshenko/rui"

const resizableDemoText = `
GridLayout {
	cell-width = "auto, 1fr, auto", cell-height = "auto, 1fr, auto",
	content = [
		Resizable {
			id = resizableTop, column = 0:2, row = 0, side = bottom,
			background-color = lightgrey,
			content = GridLayout {
				cell-vertical-align = center, cell-horizontal-align = center,
				background-color = yellow, padding = 8px, content = "Top",
			}
		},
		Resizable {
			id = resizableBottom, column = 0:2, row = 2, side = top,
			background-color = lightgrey,
			content = GridLayout {
				cell-vertical-align = center, cell-horizontal-align = center,
				background-color = lightcoral, padding = 8px, content = "Bottom",
			}
		},
		Resizable {
			id = resizableLeft, column = 0, row = 0:2, side = right,
			background-color = lightgrey,
			content = GridLayout {
				cell-vertical-align = center, cell-horizontal-align = center,
				background-color = lightskyblue, padding = 8px, content = "Left",
			}
		},
		Resizable {
			id = resizableRight, column = 2, row = 0:2, side = left,
			background-color = lightgrey,
			content = GridLayout {
				cell-vertical-align = center, cell-horizontal-align = center,
				background-color = lightpink, padding = 8px, content = "Right",
			}
		}
		GridLayout {
			column = 1, row = 1, cell-vertical-align = center, cell-horizontal-align = center,
			content = Resizable {
				id = resizableRight, side = all,
				background-color = lightgrey,
				content = GridLayout {
					cell-vertical-align = center, cell-horizontal-align = center,
					background-color = lightseagreen, padding = 8px, content = "Center",
				}
			}
		}
	]
}
`

func createResizableDemo(session rui.Session) rui.View {
	return rui.CreateViewFromText(session, resizableDemoText)
	/*
		return rui.NewGridLayout(session, rui.Params{
			rui.CellWidth:  []rui.SizeUnit{rui.AutoSize(), rui.Fr(1), rui.AutoSize()},
			rui.CellHeight: []rui.SizeUnit{rui.AutoSize(), rui.Fr(1), rui.AutoSize()},
			rui.Content: []rui.View{
				rui.NewResizable(session, rui.Params{
					rui.ID:              "resizableTop",
					rui.Column:          rui.Range{First: 0, Last: 2},
					rui.Row:             0,
					rui.Side:            rui.BottomSide,
					rui.BackgroundColor: rui.LightGray,
					rui.Content: rui.NewGridLayout(session, rui.Params{
						rui.BackgroundColor:     rui.Yellow,
						rui.CellHorizontalAlign: rui.CenterAlign,
						rui.CellVerticalAlign:   rui.CenterAlign,
						rui.Content:             "Top",
					}),
				}),
			},
		})
	*/
}
