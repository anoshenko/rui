package main

import "github.com/anoshenko/rui"

const backgroundDemoText = `
GridLayout {
	style = demoPage,
	content = [
		ListLayout {
			width = 100%, height = 100%, padding = 32px,
			content = [
				TextView {
					id = backgroundView, width = 100%, height = 150%, padding = 16px,
					text = "Sample text", text-size = 4em, 
					border = _{ style = dotted, width = 8px, color = #FF008800 },
					background = image { src = cat.jpg }
				}
			]
		},		
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Image" },
						DropDownList { row = 0, column = 1, id = backgroundImage1, current = 0, items = ["cat.jpg", "winds.png", "gifsInEmail.gif", "mountain.svg"]},
						TextView { row = 1, text = "Fit" },
						DropDownList { row = 1, column = 1, id = backgroundFit1, current = 0, items = ["none", "contain", "cover"]},
						TextView { row = 2, text = "Horizontal align" },
						DropDownList { row = 2, column = 1, id = backgroundHAlign1, current = 0, items = ["left", "right", "center"]},
						TextView { row = 3, text = "Vertical align" },
						DropDownList { row = 3, column = 1, id = backgroundVAlign1, current = 0, items = ["top", "bottom", "center"]},
						TextView { row = 4, text = "Repeat" },
						DropDownList { row = 4, column = 1, id = backgroundRepeat1, current = 0, items = ["no-repeat", "repeat", "repeat-x", "repeat-y", "round", "space"]},
						TextView { row = 5, text = "Clip" },
						DropDownList { row = 5, column = 1, id = backgroundClip1, current = 0, items = ["padding-box", "border-box", "content-box", "text"]},
						TextView { row = 6, text = "Attachment" },
						DropDownList { row = 6, column = 1, id = backgroundAttachment1, current = 0, items = ["scroll", "fixed", "local"]},
					]
				}
			]
		}
	]
}
`

func createBackgroundDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, backgroundDemoText)
	if view == nil {
		return nil
	}

	updateBackground1 := func(list rui.DropDownList, number int) {
		images := []string{"cat.jpg", "winds.png", "gifsInEmail.gif", "mountain.svg"}
		image := rui.NewBackgroundImage(rui.Params{
			rui.Source:          images[rui.GetCurrent(view, "backgroundImage1")],
			rui.Fit:             rui.GetCurrent(view, "backgroundFit1"),
			rui.HorizontalAlign: rui.GetCurrent(view, "backgroundHAlign1"),
			rui.VerticalAlign:   rui.GetCurrent(view, "backgroundVAlign1"),
			rui.Repeat:          rui.GetCurrent(view, "backgroundRepeat1"),
			rui.BackgroundClip:  rui.GetCurrent(view, "backgroundClip1"),
			rui.Attachment:      rui.GetCurrent(view, "backgroundAttachment1"),
		})
		rui.Set(view, "backgroundView", rui.Background, image)
	}

	for _, id := range []string{
		"backgroundImage1",
		"backgroundFit1",
		"backgroundHAlign1",
		"backgroundVAlign1",
		"backgroundRepeat1",
		"backgroundClip1",
		"backgroundAttachment1",
	} {
		rui.Set(view, id, rui.DropDownEvent, updateBackground1)
	}

	return view
}
