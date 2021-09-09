package main

import (
	"github.com/anoshenko/rui"
)

const viewDemoText = `
GridLayout {
	style = demoPage,
	content = [
		GridLayout {
			width = 100%, height = 100%, cell-vertical-align = center, cell-horizontal-align = center,
			content = [
				View { 
					id = demoView, width = 250px, height = 150px,
					background-color = #FFFF0000, border-width = 1px
				}
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Border style" },
						DropDownList { row = 0, column = 1, id = viewBorderStyle, current = 0,
							items = ["none", "solid", "dashed", "dotted", "double", "4 styles"]
						},
						TextView { row = 1, text = "Border width" },
						DropDownList { row = 1, column = 1, id = viewBorderWidth, current = 0,
							items = ["1px", "2px", "3.5px", "5px", "1px,2px,3px,4px"]
						},
						TextView { row = 2, text = "Border color" },
						DropDownList { row = 2, column = 1, id = viewBorderColor, current = 0,
							items = ["black", "blue", "4 colors"]
						},
						TextView { row = 3, text = "Radius" },
						DropDownList { row = 3, column = 1, id = viewRadius, current = 0,
							items = ["0", "8px", "12px/24px", "0 48px/24px 0 48px/24px"]
						},
                
					]
				}
			]
		}
	]
}
`

func viewDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, viewDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "viewBorderStyle", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		if number < 5 {
			rui.Set(view, "demoView", rui.BorderStyle, number)
		} else {
			rui.Set(view, "demoView", rui.BorderTopStyle, 1)
			rui.Set(view, "demoView", rui.BorderRightStyle, 2)
			rui.Set(view, "demoView", rui.BorderBottomStyle, 3)
			rui.Set(view, "demoView", rui.BorderLeftStyle, 4)
		}
	})

	rui.Set(view, "viewBorderWidth", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		widths := []rui.SizeUnit{rui.Px(1), rui.Px(2), rui.Px(3.5), rui.Px(5)}
		if number < len(widths) {
			rui.Set(view, "demoView", rui.BorderWidth, widths[number])
		} else {
			rui.SetParams(view, "demoView", rui.Params{
				rui.BorderTopWidth:    rui.Px(1),
				rui.BorderRightWidth:  rui.Px(2),
				rui.BorderBottomWidth: rui.Px(3),
				rui.BorderLeftWidth:   rui.Px(4)})
		}
	})

	rui.Set(view, "viewBorderColor", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		colors := []rui.Color{rui.Black, rui.Blue}
		if number < len(colors) {
			rui.Set(view, "demoView", rui.BorderColor, colors[number])
		} else {
			rui.SetParams(view, "demoView", rui.Params{
				rui.BorderTopColor:    rui.Blue,
				rui.BorderRightColor:  rui.Green,
				rui.BorderBottomColor: rui.Magenta,
				rui.BorderLeftColor:   rui.Aqua})
		}
	})

	rui.Set(view, "viewRadius", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		switch number {
		case 0:
			rui.Set(view, "demoView", rui.Radius, nil)

		case 1:
			rui.Set(view, "demoView", rui.Radius, rui.Px(8))

		case 2:
			rui.Set(view, "demoView", rui.RadiusX, rui.Px(12))
			rui.Set(view, "demoView", rui.RadiusY, rui.Px(24))

		case 3:
			rui.Set(view, "demoView", rui.Radius, nil)
			rui.Set(view, "demoView", rui.RadiusTopRightX, rui.Px(48))
			rui.Set(view, "demoView", rui.RadiusTopRightY, rui.Px(24))
			rui.Set(view, "demoView", rui.RadiusBottomLeftX, rui.Px(48))
			rui.Set(view, "demoView", rui.RadiusBottomLeftY, rui.Px(24))
		}
	})

	return view
}
