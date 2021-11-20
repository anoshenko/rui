package main

import "github.com/anoshenko/rui"

const stackLayoutDemoText = `
GridLayout {
	style = demoPage,
	content = [
		StackLayout {
			id = stackLayout, width = 100%, height = 100%
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						Button { row = 0, column = 0:1, id = pushRed, content = "Push red view" },
						Button { row = 1, column = 0:1, id = pushGreen, content = "Push green view" },
						Button { row = 2, column = 0:1, id = pushBlue, content = "Push blue view" },
						Button { row = 3, column = 0:1, id = popView, content = "Pop view" },
						TextView { row = 4, text = "Animation" },
						DropDownList { row = 4, column = 1, id = pushAnimation, current = 0, items = ["default", "start-to-end", "end-to-start", "top-down", "bottom-up"]},
						TextView { row = 5, text = "Timing" },
						DropDownList { row = 5, column = 1, id = pushTiming, current = 0, items = ["ease", "linear"]},
						TextView { row = 6, text = "Duration" },
						DropDownList { row = 6, column = 1, id = pushDuration, current = 0, items = ["0.5s", "1s", "2s"]},
					]
				}
			]
		}
	]
}
`

func createStackLayoutDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, stackLayoutDemoText)
	if view == nil {
		return nil
	}

	animation := func() int {
		return rui.GetCurrent(view, "pushAnimation")
	}

	/*
		transition := func() rui.ViewTransition {
			timing := rui.EaseTiming
			timings := []string{rui.EaseTiming, rui.LinearTiming}
			if n := rui.GetCurrent(view, "pushTiming"); n >= 0 && n < len(timings) {
				timing = timings[n]
			}

			duration := float64(0.5)
			durations := []float64{0.5, 1, 2}
			if n := rui.GetCurrent(view, "pushDuration"); n >= 0 && n < len(durations) {
				duration = durations[n]
			}

			return rui.NewTransition(duration, timing, session)
		}
	*/

	rui.Set(view, "pushRed", rui.ClickEvent, func(_ rui.View) {
		if stackLayout := rui.StackLayoutByID(view, "stackLayout"); stackLayout != nil {
			if v := rui.CreateViewFromText(session, `View { background-color = red }`); v != nil {
				stackLayout.Push(v, animation(), nil)
			}
		}
	})

	rui.Set(view, "pushGreen", rui.ClickEvent, func(_ rui.View) {
		if stackLayout := rui.StackLayoutByID(view, "stackLayout"); stackLayout != nil {
			if v := rui.CreateViewFromText(session, `View { background-color = green }`); v != nil {
				stackLayout.Push(v, animation(), nil)
			}
		}
	})

	rui.Set(view, "pushBlue", rui.ClickEvent, func(_ rui.View) {
		if stackLayout := rui.StackLayoutByID(view, "stackLayout"); stackLayout != nil {
			if v := rui.CreateViewFromText(session, `View { background-color = blue }`); v != nil {
				stackLayout.Push(v, animation(), nil)
			}
		}
	})

	rui.Set(view, "popView", rui.ClickEvent, func(_ rui.View) {
		if stackLayout := rui.StackLayoutByID(view, "stackLayout"); stackLayout != nil {
			stackLayout.Pop(animation(), nil)
		}
	})

	return view
}
