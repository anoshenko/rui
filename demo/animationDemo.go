package main

import "github.com/anoshenko/rui"

const animationDemoText = `
GridLayout {
	style = demoPage,
	content = [
		AbsoluteLayout {
			id = animationContainer, width = 100%, height = 100%,
			content = [ 
				View { 
					id = animatedView1, width = 32px, height = 32px, left = 16px, top = 16px, background-color = #FF0000FF 
				} 
			]
		},
		ListLayout {
			style = optionsPanel,
			content = [
				GridLayout {
					style = optionsTable,
					content = [
						TextView { row = 0, text = "Duration" },
						DropDownList { row = 0, column = 1, id = animationDuration, current = 0, items = ["4s", "8s", "12s"]},
						TextView { row = 1, text = "Delay" },
						DropDownList { row = 1, column = 1, id = animationDelay, current = 0, items = ["0s", "1s", "2s"]},
						TextView { row = 2, text = "Timing function" },
						DropDownList { row = 2, column = 1, id = animationTimingFunction, current = 0, items = ["ease", "linear", "steps(40)"]},
						TextView { row = 3, text = "Iteration Count" },
						DropDownList { row = 3, column = 1, id = animationIterationCount, current = 0, items = ["1", "3", "infinite"]},
						TextView { row = 4, text = "Direction" },
						DropDownList { row = 4, column = 1, id = animationDirection, current = 0, items = ["normal", "reverse", "alternate", "alternate-reverse"]},						
						Button { row = 5, column = 0:1, id = animationStart, content = Start },					
						Button { row = 6, column = 0:1, id = animationPause, content = Pause },					
					]
				}
			]
		}
	]
}`

func createAnimationDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, animationDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "animationStart", rui.ClickEvent, func() {
		frame := rui.GetViewFrame(view, "animationContainer")
		prop1 := rui.AnimatedProperty{
			Tag:  rui.Left,
			From: rui.Px(16),
			To:   rui.Px(16),
			KeyFrames: map[int]interface{}{
				25: rui.Px(frame.Width - 48),
				50: rui.Px(frame.Width - 48),
				75: rui.Px(16),
			},
		}
		prop2 := rui.AnimatedProperty{
			Tag:  rui.Top,
			From: rui.Px(16),
			To:   rui.Px(16),
			KeyFrames: map[int]interface{}{
				25: rui.Px(16),
				50: rui.Px(frame.Height - 48),
				75: rui.Px(frame.Height - 48),
			},
		}
		prop3 := rui.AnimatedProperty{
			Tag:  rui.Rotate,
			From: rui.Deg(0),
			To:   rui.Deg(360),
			KeyFrames: map[int]interface{}{
				25: rui.Deg(90),
				50: rui.Deg(180),
				75: rui.Deg(270),
			},
		}

		params := rui.Params{
			rui.PropertyTag:        []rui.AnimatedProperty{prop1, prop2, prop3},
			rui.Duration:           rui.GetDropDownCurrent(view, "animationDuration") * 4,
			rui.Delay:              rui.GetDropDownCurrent(view, "animationDelay"),
			rui.AnimationDirection: rui.GetDropDownCurrent(view, "animationDirection"),
		}

		switch rui.GetDropDownCurrent(view, "animationTimingFunction") {
		case 0:
			params[rui.TimingFunction] = rui.EaseTiming

		case 1:
			params[rui.TimingFunction] = rui.LinearTiming

		case 2:
			params[rui.TimingFunction] = rui.StepsTiming(40)
		}

		switch rui.GetDropDownCurrent(view, "animationIterationCount") {
		case 0:
			params[rui.IterationCount] = 1

		case 1:
			params[rui.IterationCount] = 3

		case 2:
			params[rui.IterationCount] = -1
		}

		rui.Set(view, "animatedView1", rui.AnimationTag, rui.NewAnimation(params))
	})

	rui.Set(view, "animationPause", rui.ClickEvent, func() {
		if rui.IsAnimationPaused(view, "animatedView1") {
			rui.Set(view, "animatedView1", rui.AnimationPaused, false)
			rui.Set(view, "animationPause", rui.Content, "Pause")
		} else {
			rui.Set(view, "animatedView1", rui.AnimationPaused, true)
			rui.Set(view, "animationPause", rui.Content, "Resume")
		}
	})

	rui.Set(view, "animatedView1", rui.AnimationStartEvent, func() {
		rui.Set(view, "animatedView1", rui.AnimationPaused, false)
		rui.Set(view, "animationPause", rui.Content, "Pause")
	})

	rui.Set(view, "animatedView1", rui.AnimationEndEvent, func() {
		rui.Set(view, "animatedView1", rui.AnimationPaused, false)
		rui.Set(view, "animationPause", rui.Content, "Pause")
	})

	rui.Set(view, "animatedView1", rui.AnimationCancelEvent, func() {
		rui.Set(view, "animatedView1", rui.AnimationPaused, false)
		rui.Set(view, "animationPause", rui.Content, "Pause")
	})

	rui.Set(view, "animatedView1", rui.AnimationIterationEvent, func() {
	})

	return view
}
