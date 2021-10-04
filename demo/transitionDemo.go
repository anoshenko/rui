package main

import (
	"github.com/anoshenko/rui"
)

const transitionDemoText = `
ListLayout {
	width = 100%, height = 100%, orientation = vertical, padding = 12px,
	content = [
		"ease",
		View { id = bar1, width = 20%, style = transitionBar },
		"ease-in",
		View { id = bar2, width = 20%, style = transitionBar },
		"ease-out",
		View { id = bar3, width = 20%, style = transitionBar },
		"ease-in-out",
		View { id = bar4, width = 20%, style = transitionBar },
		"linear",
		View { id = bar5, width = 20%, style = transitionBar },
		"steps(5)",
		View { id = bar6, width = 20%, style = transitionBar },
		"cubic-bezier(0.1, -0.6, 0.2, 0)",
		View { id = bar7, width = 20%, style = transitionBar },
		Button { id = startTransition, content = "Start" }
	]
}`

func createTransitionDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, transitionDemoText)
	if view == nil {
		return nil
	}

	bars := map[string]string{
		"bar1": rui.EaseTiming,
		"bar2": rui.EaseInTiming,
		"bar3": rui.EaseOutTiming,
		"bar4": rui.EaseInOutTiming,
		"bar5": rui.LinearTiming,
		"bar6": rui.StepsTiming(5),
		"bar7": rui.CubicBezierTiming(0.1, -0.6, 0.2, 0),
	}

	rui.Set(view, "startTransition", rui.ClickEvent, func(button rui.View) {

		for id, timing := range bars {
			animation := rui.NewAnimation(rui.Params{
				rui.Duration:       2,
				rui.TimingFunction: timing,
			})

			if bar := rui.ViewByID(view, id); bar != nil {
				if rui.GetWidth(bar, "").Value == 100 {
					bar.Remove(rui.TransitionEndEvent)
					bar.SetAnimated(rui.Width, rui.Percent(20), animation)
				} else {
					bar.Set(rui.TransitionEndEvent, func(v rui.View, tag string) {
						bar.Remove(rui.TransitionEndEvent)
						bar.SetAnimated(rui.Width, rui.Percent(20), animation)
					})
					bar.SetAnimated(rui.Width, rui.Percent(100), animation)
				}
			}
		}
	})

	return view
}
