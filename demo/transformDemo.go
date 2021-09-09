package main

import (
	"fmt"

	"github.com/anoshenko/rui"
)

const transformDemoText = `
GridLayout {
	style = demoPage,
	content = [
		GridLayout {
			id = listLayout, width = 100%, height = 100%, cell-horizontal-align = center, cell-vertical-align = center,
			content = [
				TextView { id = transformView, width = 200px, height = 100px, 
					text = "View", text-align = center, text-size = 32pt, 
					background-color = #FFBBBBBB, radius = 8px, padding = 8px, margin = 4px,
					border = _{ style = solid, width = 1px, color = #FF000080 },
					shadow = _{ spread-radius = 1px, blur = 16px, color = #80000000},
				},	
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				"Perspective",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = PerspectiveEditor, type = slider, width = 120px, 
							min = -500, max = 500, step = 10, value = 0 
						},
						TextView {
							id = PerspectiveValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Perspective origin X (pixels)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = xPerspectiveOriginEditor, type = slider, width = 120px, 
							min = -1000, max = 1000, step = 10, value = 0 
						},
						TextView {
							id = xPerspectiveOriginValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Perspective origin Y (pixels)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = yPerspectiveOriginEditor, type = slider, width = 120px, 
							min = -1000, max = 1000, step = 10, value = 0 
						},
						TextView {
							id = yPerspectiveOriginValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Origin X (pixels)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = xOriginEditor, type = slider, width = 120px, 
							min = -1000, max = 1000, step = 10, value = 0 
						},
						TextView {
							id = xOriginValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Origin Y (pixels)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = yOriginEditor, type = slider, width = 120px, 
							min = -1000, max = 1000, step = 10, value = 0 
						},
						TextView {
							id = yOriginValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Origin Z (pixels)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = zOriginEditor, type = slider, width = 120px, 
							min = -1000, max = 1000, step = 10, value = 0 
						},
						TextView {
							id = zOriginValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Translate X (pixels)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = xTranslateEditor, type = slider, width = 120px, 
							min = -100, max = 100, step = 1, value = 0 
						},
						TextView {
							id = xTranslateValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Translate Y (pixels)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = yTranslateEditor, type = slider, width = 120px, 
							min = -100, max = 100, step = 1, value = 0 
						},
						TextView {
							id = yTranslateValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Translate Z (pixels)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = zTranslateEditor, type = slider, width = 120px, 
							min = -100, max = 100, step = 1, value = 0 
						},
						TextView {
							id = zTranslateValue, text = "0px", margin-left = 12px, width = 32px
						}
					]
				},
				"Scale X",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = xScaleEditor, type = slider, width = 120px, 
							min = -5, max = 5, step = 0.1, value = 1 
						},
						TextView {
							id = xScaleValue, text = "1", margin-left = 12px, width = 32px
						}
					]
				},
				"Scale Y",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = yScaleEditor, type = slider, width = 120px, 
							min = -5, max = 5, step = 0.1, value = 1 
						},
						TextView {
							id = yScaleValue, text = "1", margin-left = 12px, width = 32px
						}
					]
				},
				"Scale Z",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = zScaleEditor, type = slider, width = 120px, 
							min = -5, max = 5, step = 0.1, value = 1 
						},
						TextView {
							id = zScaleValue, text = "1", margin-left = 12px, width = 32px
						}
					]
				},
				"Skew X (degree)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = xSkewEditor, type = slider, width = 120px, 
							min = -90, max = 90, step = 1, value = 0 
						},
						TextView {
							id = xSkewValue, text = "0°", margin-left = 12px, width = 32px
						}
					]
				},
				"Skew Y (degree)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = ySkewEditor, type = slider, width = 120px, 
							min = -90, max = 90, step = 1, value = 0 
						},
						TextView {
							id = ySkewValue, text = "0°", margin-left = 12px, width = 32px
						}
					]
				},
				"Rotate (degree)",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = RotateEditor, type = slider, width = 120px, 
							min = -180, max = 180, step = 1, value = 0 
						},
						TextView {
							id = RotateValue, text = "0°", margin-left = 12px, width = 32px
						}
					]
				},
				"Rotate X",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = xRotateEditor, type = slider, width = 120px, 
							min = 0, max = 1, step = 0.01, value = 1 
						},
						TextView {
							id = xRotateValue, text = "1", margin-left = 12px, width = 32px
						}
					]
				},
				"Rotate Y",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = yRotateEditor, type = slider, width = 120px, 
							min = 0, max = 1, step = 0.01, value = 1 
						},
						TextView {
							id = yRotateValue, text = "1", margin-left = 12px, width = 32px
						}
					]
				},
				"Rotate Z",
				ListLayout { 
					margin-bottom = 12px, orientation = horizontal, vertical-align = center,
					content = [
						NumberPicker { id = zRotateEditor, type = slider, width = 120px, 
							min = 0, max = 1, step = 0.01, value = 1 
						},
						TextView {
							id = zRotateValue, text = "1", margin-left = 12px, width = 32px
						}
					]
				},
				Checkbox { id = backfaceVisibility, content = "backface-visibility", checked = true }				
			]
		}
	]		
}`

func transformDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, transformDemoText)
	if view == nil {
		return nil
	}

	// 		transform := rui.NewTransform(view.Session())

	transformView := rui.ViewByID(view, "transformView")
	if transformView == nil {
		return view
	}

	updateSliderText := func(tag string, value float64, unit string) {
		rui.Set(view, tag, rui.Text, fmt.Sprintf("%g%s", value, unit))
	}

	rui.Set(view, "PerspectiveEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.Perspective, rui.Px(newValue))
		updateSliderText("PerspectiveValue", newValue, "px")
	})

	rui.Set(view, "xPerspectiveOriginEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.PerspectiveOriginX, rui.Px(newValue))
		updateSliderText("xPerspectiveOriginValue", newValue, "px")
	})

	rui.Set(view, "yPerspectiveOriginEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.PerspectiveOriginY, rui.Px(newValue))
		updateSliderText("yPerspectiveOriginValue", newValue, "px")
	})

	rui.Set(view, "xOriginEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.OriginX, rui.Px(newValue))
		updateSliderText("xOriginValue", newValue, "px")
	})

	rui.Set(view, "yOriginEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.OriginY, rui.Px(newValue))
		updateSliderText("yOriginValue", newValue, "px")
	})

	rui.Set(view, "zOriginEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.OriginZ, rui.Px(newValue))
		updateSliderText("zOriginValue", newValue, "px")
	})

	rui.Set(view, "xTranslateEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.TranslateX, rui.Px(newValue))
		updateSliderText("xTranslateValue", newValue, "px")
	})

	rui.Set(view, "yTranslateEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.TranslateY, rui.Px(newValue))
		updateSliderText("yTranslateValue", newValue, "px")
	})

	rui.Set(view, "zTranslateEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.TranslateZ, rui.Px(newValue))
		updateSliderText("zTranslateValue", newValue, "px")
	})

	rui.Set(view, "xScaleEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.ScaleX, newValue)
		updateSliderText("xScaleValue", newValue, "")
	})

	rui.Set(view, "yScaleEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.ScaleY, newValue)
		updateSliderText("yScaleValue", newValue, "")
	})

	rui.Set(view, "zScaleEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.ScaleZ, newValue)
		updateSliderText("zScaleValue", newValue, "")
	})

	rui.Set(view, "RotateEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.Rotate, rui.Deg(newValue))
		updateSliderText("RotateValue", newValue, "°")
	})

	rui.Set(view, "xRotateEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.RotateX, newValue)
		updateSliderText("xRotateValue", newValue, "")
	})

	rui.Set(view, "yRotateEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.RotateY, newValue)
		updateSliderText("yRotateValue", newValue, "")
	})

	rui.Set(view, "zRotateEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.RotateZ, newValue)
		updateSliderText("zRotateValue", newValue, "")
	})

	rui.Set(view, "xSkewEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.SkewX, rui.Deg(newValue))
		updateSliderText("xSkewValue", newValue, "°")
	})

	rui.Set(view, "ySkewEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		transformView.Set(rui.SkewY, rui.Deg(newValue))
		updateSliderText("ySkewValue", newValue, "°")
	})

	rui.Set(view, "backfaceVisibility", rui.CheckboxChangedEvent, func(checkbox rui.Checkbox, checked bool) {
		transformView.Set(rui.BackfaceVisible, checked)
	})

	return view
}
