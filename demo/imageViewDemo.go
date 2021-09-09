package main

import (
	"github.com/anoshenko/rui"
)

const imageViewDemoText = `
GridLayout {
	style = demoPage,
	content = [
		ImageView {
			id = imageView1, width = 100%, height = 100%, src = "cat.jpg",
			border = _{ style = solid, width = 1px, color = #FF008800 } 
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Image" },
						DropDownList { row = 0, column = 1, id = imageViewImage, current = 0, items = ["cat.jpg", "winds.png", "gifsInEmail.gif", "mountain.svg"]},
						TextView { row = 1, text = "Fit" },
						DropDownList { row = 1, column = 1, id = imageViewFit, current = 0, items = ["none", "fill", "contain", "cover", "scale-down"]},
						TextView { row = 2, text = "Horizontal align" },
						DropDownList { row = 2, column = 1, id = imageViewHAlign, current = 2, items = ["left", "right", "center"]},
						TextView { row = 3, text = "Vertical align" },
						DropDownList { row = 3, column = 1, id = imageViewVAlign, current = 2, items = ["top", "bottom", "center"]},
					]
				}
			]
		}
	]
}
`

func createImageViewDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, imageViewDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "imageViewImage", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		images := []string{"cat.jpg", "winds.png", "gifsInEmail.gif", "mountain.svg"}
		if number < len(images) {
			rui.Set(view, "imageView1", rui.Source, images[number])
		}
	})

	rui.Set(view, "imageViewFit", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "imageView1", rui.Fit, number)
	})

	rui.Set(view, "imageViewHAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "imageView1", rui.ImageHorizontalAlign, number)
	})

	rui.Set(view, "imageViewVAlign", rui.DropDownEvent, func(list rui.DropDownList, number int) {
		rui.Set(view, "imageView1", rui.ImageVerticalAlign, number)
	})

	return view
}
