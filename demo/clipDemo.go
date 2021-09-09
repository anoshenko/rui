package main

import "github.com/anoshenko/rui"

const clipDemoText = `
GridLayout {
	width = 100%, height = 100%, cell-height = "auto, 1fr",
	cell-horizontal-align = center, cell-vertical-align = center,
	content = [
		DropDownList {
			id = clipType, current = 0, margin = 8px, max-width = 100%,
			items = ["none", 
				"inset(20%, 10%, 20%, 10%, 16px / 32px)", 
				"circle(50%, 45%, 45%)", 
				"ellipse(50%, 50%, 35%, 50%)",
				"polygon(50%, 2.4%, 34.5%, 33.8%, 0%, 38.8%, 25%, 63.1%, 19.1%, 97.6%, 50%, 81.3%, 80.9%, 97.6%, 75%, 63.1%, 100%, 38.8%, 65.5%, 33.8%)"],
		},
		ImageView {
			id = clipImage, row = 1, src = "cat.jpg",
		}
	]
}
`

func createClipDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, clipDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "clipType", rui.DropDownEvent, func(number int) {
		switch number {
		case 0:
			rui.Set(view, "clipImage", rui.Clip, nil)

		case 1:
			rui.Set(view, "clipImage", rui.Clip, rui.InsetClip(rui.Percent(20), rui.Percent(10),
				rui.Percent(20), rui.Percent(10), rui.NewRadiusProperty(rui.Params{
					rui.X: rui.Px(16),
					rui.Y: rui.Px(32),
				})))
		case 2:
			rui.Set(view, "clipImage", rui.Clip, rui.CircleClip(rui.Percent(50), rui.Percent(45), rui.Percent(45)))

		case 3:
			rui.Set(view, "clipImage", rui.Clip, rui.EllipseClip(rui.Percent(50), rui.Percent(50), rui.Percent(35), rui.Percent(50)))

		case 4:
			rui.Set(view, "clipImage", rui.Clip, rui.PolygonClip([]interface{}{"50%", "2.4%", "34.5%", "33.8%", "0%", "38.8%", "25%", "63.1%", "19.1%", "97.6%", "50%", "81.3%", "80.9%", "97.6%", "75%", "63.1%", "100%", "38.8%", "65.5%", "33.8%"}))
		}
	})

	return view
}
