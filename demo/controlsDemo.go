package main

import (
	"fmt"
	"math"
	"time"

	"github.com/anoshenko/rui"
)

const controlsDemoText = `
ListLayout {
	width = 100%, height = 100%, orientation = vertical, padding = 16px,
	content = [
		DetailsView {
			margin = 8px,
			summary = "Details title",
			content = "Details content"
		}
		ListLayout { orientation = horizontal, vertical-align = center, padding = 8px,
			border = _{ width = 1px, style = solid, color = #FF000000 }, radius = 4px,
			content = [
				Checkbox { id = controlsCheckbox, content = "Checkbox" },
				Button { id = controlsCheckboxButton, margin-left = 32px, content = "Check checkbox" },
			]
		},
		ListLayout { orientation = horizontal, margin-top = 16px, padding = 8px, vertical-align = center,
			border = _{ width = 1px, style = solid, color = #FF000000 }, radius = 4px,
			content = [
				Button { id = controlsProgressDec, content = "<<" },
				Button { id = controlsProgressInc, content = ">>", margin-left = 12px },
				ProgressBar { id = controlsProgress, max = 100, value = 50, margin-left = 12px  },
				TextView { id = controlsProgressLabel, text = "50 / 100", margin-left = 12px },
			]
		},
		ListLayout { orientation = horizontal, margin-top = 16px, padding = 8px, vertical-align = center,
			border = _{ width = 1px, style = solid, color = #FF000000 }, radius = 4px,
			content = [
				"Enter number (-5...10)",
				NumberPicker { id = controlsNumberEditor, type = editor, width = 80px, 
					margin-left = 12px, min = -5, max = 10, step = 0.1, value = 0 
				},
				NumberPicker { id = controlsNumberSlider, type = slider, width = 150px, 
					margin-left = 12px, min = -5, max = 10, step = 0.1, value = 0 
				}
			]
		},
		ListLayout { orientation = horizontal, margin-top = 16px, padding = 8px, vertical-align = center,
			border = _{ width = 1px, style = solid, color = #FF000000 }, radius = 4px,
			content = [
				"Select color",
				ColorPicker { id = controlsColorPicker, value = #0000FF,
					margin = _{ left = 12px, right = 24px} 
				},
				"Result",
				View { id = controlsColorResult, width = 24px, height = 24px, margin-left = 12px, background-color = #0000FF }
			]
		},
		ListLayout { orientation = horizontal, margin-top = 16px, padding = 8px, vertical-align = center,
			border = _{ width = 1px, style = solid, color = #FF000000 }, radius = 4px,
			content = [
				"Select a time and date:",
				TimePicker { id = controlsTimePicker, min = "00:00", margin-left = 12px },
				DatePicker { id = controlsDatePicker, min = "2001-01-01", margin-right = 24px },
				"Result:",
				TextView { id = controlsDateResult, margin-left = 12px }
			]
		},
		Button {
			id = controlsMessage, margin-top = 16px, content = "Show message"
		}
	]
}
`

func createControlsDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, controlsDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "controlsCheckbox", rui.CheckboxChangedEvent, func(checkbox rui.Checkbox, checked bool) {
		if checked {
			rui.Set(view, "controlsCheckboxButton", rui.Content, "Uncheck checkbox")
		} else {
			rui.Set(view, "controlsCheckboxButton", rui.Content, "Check checkbox")
		}
	})

	rui.Set(view, "controlsCheckboxButton", rui.ClickEvent, func(rui.View) {
		checked := rui.IsCheckboxChecked(view, "controlsCheckbox")
		rui.Set(view, "controlsCheckbox", rui.Checked, !checked)
	})

	setProgressBar := func(dx float64) {
		if value := rui.GetProgressBarValue(view, "controlsProgress"); value >= 0 {
			max := rui.GetProgressBarMax(view, "controlsProgress")
			newValue := math.Min(math.Max(0, value+dx), max)
			if newValue != value {
				rui.Set(view, "controlsProgress", rui.Value, newValue)
				rui.Set(view, "controlsProgressLabel", rui.Text, fmt.Sprintf("%g / %g", newValue, max))
			}
		}
	}

	rui.Set(view, "controlsProgressDec", rui.ClickEvent, func(rui.View) {
		setProgressBar(-1)
	})

	rui.Set(view, "controlsProgressInc", rui.ClickEvent, func(rui.View) {
		setProgressBar(+1)
	})

	rui.Set(view, "controlsNumberEditor", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		rui.Set(view, "controlsNumberSlider", rui.Value, newValue)
	})

	rui.Set(view, "controlsNumberSlider", rui.NumberChangedEvent, func(v rui.NumberPicker, newValue float64) {
		rui.Set(view, "controlsNumberEditor", rui.Value, newValue)
	})

	rui.Set(view, "controlsColorPicker", rui.ColorChangedEvent, func(v rui.ColorPicker, newColor rui.Color) {
		rui.Set(view, "controlsColorResult", rui.BackgroundColor, newColor)
	})

	rui.Set(view, "controlsTimePicker", rui.Value, demoTime)
	rui.Set(view, "controlsDatePicker", rui.Value, demoTime)

	rui.Set(view, "controlsTimePicker", rui.TimeChangedEvent, func(v rui.TimePicker, newDate time.Time) {
		demoTime = time.Date(demoTime.Year(), demoTime.Month(), demoTime.Day(), newDate.Hour(), newDate.Minute(),
			newDate.Second(), newDate.Nanosecond(), demoTime.Location())
		rui.Set(view, "controlsDateResult", rui.Text, demoTime.Format(time.RFC1123))
	})

	rui.Set(view, "controlsDatePicker", rui.DateChangedEvent, func(v rui.DatePicker, newDate time.Time) {
		demoTime = time.Date(newDate.Year(), newDate.Month(), newDate.Day(), demoTime.Hour(), demoTime.Minute(),
			demoTime.Second(), demoTime.Nanosecond(), demoTime.Location())
		rui.Set(view, "controlsDateResult", rui.Text, demoTime.Format(time.RFC1123))
	})

	rui.Set(view, "controlsMessage", rui.ClickEvent, func(rui.View) {
		rui.ShowMessage("Hello", "Hello world!!!", session)
	})

	return view
}

var demoTime = time.Now()
