package main

import (
	"fmt"
	"strings"

	"github.com/anoshenko/rui"
)

const filterDemoText = `
GridLayout {
	style = demoPage,
	content = [
		GridLayout {
			width = 100%, height = 100%, cell-vertical-align = center, cell-horizontal-align = center,
			content = [
				ImageView { id = filterImage, src = "mountain.svg" },
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						Checkbox { id = blurCheckbox, row = 0, content = "Blur" },
						NumberPicker { id = blurSlider, row = 1, type = slider, min = 0, max = 25, step = 0.1, disabled = true },
						TextView { id = blurValue, row = 1, column = 1, text = "0px", width = 40px },
						Checkbox { id = brightnessCheckbox, row = 2, content = "Brightness" },
						NumberPicker { id = brightnessSlider, row = 3, type = slider, min = 0, max = 200, step = 1, disabled = true },
						TextView { id = brightnessValue, row = 3, column = 1, text = "0%", width = 40px },
						Checkbox { id = contrastCheckbox, row = 4, content = "Contrast" },
						NumberPicker { id = contrastSlider, row = 5, type = slider, min = 0, max = 200, step = 1, disabled = true },
						TextView { id = contrastValue, row = 5, column = 1, text = "0%", width = 40px },
						Checkbox { id = grayscaleCheckbox, row = 6, content = "Grayscale" },
						NumberPicker { id = grayscaleSlider, row = 7, type = slider, min = 0, max = 100, step = 1, disabled = true },
						TextView { id = grayscaleValue, row = 7, column = 1, text = "0%", width = 40px },
						Checkbox { id = invertCheckbox, row = 8, content = "Invert" },
						NumberPicker { id = invertSlider, row = 9, type = slider, min = 0, max = 100, step = 1, disabled = true },
						TextView { id = invertValue, row = 9, column = 1, text = "0%", width = 40px },
						Checkbox { id = saturateCheckbox, row = 10, content = "Saturate" },
						NumberPicker { id = saturateSlider, row = 11, type = slider, min = 0, max = 200, step = 1, disabled = true },
						TextView { id = saturateValue, row = 11, column = 1, text = "0%", width = 40px },
						Checkbox { id = sepiaCheckbox, row = 12, content = "Sepia" },
						NumberPicker { id = sepiaSlider, row = 13, type = slider, min = 0, max = 100, step = 1, disabled = true },
						TextView { id = sepiaValue, row = 13, column = 1, text = "0%", width = 40px },
						Checkbox { id = opacityCheckbox, row = 14, content = "Opacity" },
						NumberPicker { id = opacitySlider, row = 15, type = slider, min = 0, max = 100, step = 1, disabled = true },
						TextView { id = opacityValue, row = 15, column = 1, text = "0%", width = 40px },
						Checkbox { id = huerotateCheckbox, row = 16, content = "hue-rotate" },
						NumberPicker { id = huerotateSlider, row = 17, type = slider, min = 0, max = 720, step = 1, disabled = true },
						TextView { id = huerotateValue, row = 17, column = 1, text = "0°", width = 40px },
						Checkbox { id = shadowCheckbox, row = 18, content = "drop-shadow" },
						ColorPicker { id = dropShadowColor, row = 18, column = 1, color = black, disabled = true },
						NumberPicker { id = shadowXSlider, row = 19, type = slider, min = -20, max = 20, step = 1, value = 0, disabled = true },
						TextView { id = shadowXValue, row = 19, column = 1, text = "x:0px", width = 40px },
						NumberPicker { id = shadowYSlider, row = 20, type = slider, min = -20, max = 20, step = 1, value = 0, disabled = true },
						TextView { id = shadowYValue, row = 20, column = 1, text = "y:0px", width = 40px },
						NumberPicker { id = shadowBlurSlider, row = 21, type = slider, min = 0, max = 40, step = 1, disabled = true },
						TextView { id = shadowBlurValue, row = 21, column = 1, text = "b:0px", width = 40px },
					]
				}
			]
		}
	]
}
`

var filterParams = rui.Params{}

func createFilterDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, filterDemoText)
	if view == nil {
		return nil
	}

	setEvents := func(tag string) {
		rui.Set(view, tag+"Checkbox", rui.CheckboxChangedEvent, func(state bool) {
			slider := tag + "Slider"
			rui.Set(view, slider, rui.Disabled, !state)
			if state {
				filterParams[tag] = rui.GetNumberPickerValue(view, slider)
			} else {
				delete(filterParams, tag)
			}
			rui.Set(view, "filterImage", rui.Filter, rui.NewViewFilter(filterParams))
		})

		rui.Set(view, tag+"Slider", rui.NumberChangedEvent, func(value float64) {
			var text string
			if tag == rui.Blur {
				text = fmt.Sprintf("%.2gpx", value)
			} else {
				text = fmt.Sprintf("%g%%", value)
			}
			rui.Set(view, tag+"Value", rui.Text, text)
			filterParams[tag] = value
			rui.Set(view, "filterImage", rui.Filter, rui.NewViewFilter(filterParams))
		})
	}

	for _, tag := range []string{rui.Blur, rui.Brightness, rui.Contrast, rui.Grayscale, rui.Invert, rui.Saturate, rui.Sepia, rui.Opacity} {
		setEvents(tag)
	}

	rui.Set(view, "huerotateCheckbox", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "huerotateSlider", rui.Disabled, !state)
		if state {
			filterParams[rui.HueRotate] = rui.AngleUnit{
				Type:  rui.Degree,
				Value: rui.GetNumberPickerValue(view, "huerotateSlider"),
			}
		} else {
			delete(filterParams, rui.HueRotate)
		}
		rui.Set(view, "filterImage", rui.Filter, rui.NewViewFilter(filterParams))
	})

	rui.Set(view, "huerotateSlider", rui.NumberChangedEvent, func(value float64) {
		rui.Set(view, "huerotateValue", rui.Text, fmt.Sprintf("%g°", value))
		filterParams[rui.HueRotate] = rui.AngleUnit{Type: rui.Degree, Value: value}
		rui.Set(view, "filterImage", rui.Filter, rui.NewViewFilter(filterParams))
	})

	updateShadow := func() {
		xOff := rui.SizeUnit{
			Type:  rui.SizeInPixel,
			Value: rui.GetNumberPickerValue(view, "shadowXSlider"),
		}
		yOff := rui.SizeUnit{
			Type:  rui.SizeInPixel,
			Value: rui.GetNumberPickerValue(view, "shadowYSlider"),
		}
		blur := rui.SizeUnit{
			Type:  rui.SizeInPixel,
			Value: rui.GetNumberPickerValue(view, "shadowBlurSlider"),
		}
		color := rui.GetColorPickerValue(view, "dropShadowColor")

		filterParams[rui.DropShadow] = rui.NewTextShadow(xOff, yOff, blur, color)
		rui.Set(view, "filterImage", rui.Filter, rui.NewViewFilter(filterParams))
	}

	rui.Set(view, "shadowCheckbox", rui.CheckboxChangedEvent, func(state bool) {
		for _, tag := range []string{"shadowXSlider", "shadowYSlider", "shadowBlurSlider", "dropShadowColor"} {
			rui.Set(view, tag, rui.Disabled, !state)
		}
		if state {
			updateShadow()
		} else {
			delete(filterParams, rui.DropShadow)
			rui.Set(view, "filterImage", rui.Filter, rui.NewViewFilter(filterParams))
		}
	})

	for _, tag := range []string{"shadowXSlider", "shadowYSlider", "shadowBlurSlider"} {
		rui.Set(view, tag, rui.NumberChangedEvent, func(picker rui.NumberPicker, value float64) {
			tag := strings.Replace(picker.ID(), "Slider", "Value", -1)
			text := rui.GetText(view, tag)
			rui.Set(view, tag, rui.Text, fmt.Sprintf("%s%gpx", text[:2], value))
			updateShadow()
		})
	}

	rui.Set(view, "dropShadowColor", rui.ColorChangedEvent, func(value rui.Color) {
		updateShadow()
	})

	return view
}
