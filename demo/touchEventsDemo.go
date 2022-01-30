package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anoshenko/rui"
)

const touchEventsDemoText = `
GridLayout {
	width = 100%, height = 100%, cell-height = "1fr, auto",
	content = [
		GridLayout {
			padding = 12px,
			content = [
				GridLayout {
					id = touchEventsTest, cell-horizontal-align = center, cell-vertical-align = center,
					height = 100%,
					border = _{ style = solid, width = 1px, color = gray},
					content = [
						TextView {
							id = touchEventsText, text = "Test box",
						}
					]
				}
			],
		},
		Resizable {
			row = 1, side = top, background-color = lightgrey, height = 300px,
			content = EditView {
				id = touchEventsLog, type = multiline, read-only = true, wrap = true,
			}
		},
	]
}
`

func createTouchEventsDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, touchEventsDemoText)
	if view == nil {
		return nil
	}

	addToLog := func(tag string, event rui.TouchEvent) {
		var buffer strings.Builder

		appendBool := func(name string, value bool) {
			buffer.WriteString(`, `)
			buffer.WriteString(name)
			if value {
				buffer.WriteString(` = true`)
			} else {
				buffer.WriteString(` = false`)
			}
		}

		/*
			appendInt := func(name string, value int) {
				buffer.WriteString(`, `)
				buffer.WriteString(name)
				buffer.WriteString(` = `)
				buffer.WriteString(strconv.Itoa(value))
			}*/

		appendFloat := func(name string, value float64) {
			buffer.WriteString(fmt.Sprintf(`, %s = %g`, name, value))
		}

		appendPoint := func(name string, x, y float64) {
			buffer.WriteString(fmt.Sprintf(`, %s = (%g:%g)`, name, x, y))
		}

		buffer.WriteString(tag)
		buffer.WriteString(`: TimeStamp = `)
		buffer.WriteString(strconv.FormatUint(event.TimeStamp, 10))

		buffer.WriteString(`, touches = [`)
		for i, touch := range event.Touches {
			if i > 0 {
				buffer.WriteString(`, `)
			}
			buffer.WriteString(`{ Identifier = `)
			buffer.WriteString(strconv.Itoa(touch.Identifier))
			appendPoint("(X:Y)", touch.X, touch.Y)
			appendPoint("Client (X:Y)", touch.ClientX, touch.ClientY)
			appendPoint("Screen (X:Y)", touch.ScreenX, touch.ScreenY)
			appendPoint("Radius (X:Y)", touch.RadiusX, touch.RadiusY)
			appendFloat("RotationAngle", touch.RotationAngle)
			appendFloat("Force", touch.Force)
			buffer.WriteString(`}`)
		}
		buffer.WriteString(`]`)

		appendBool("CtrlKey", event.CtrlKey)
		appendBool("ShiftKey", event.ShiftKey)
		appendBool("AltKey", event.AltKey)
		appendBool("MetaKey", event.MetaKey)
		buffer.WriteString(";\n\n")

		rui.AppendEditText(view, "touchEventsLog", buffer.String())
		rui.ScrollViewToEnd(view, "touchEventsLog")
	}

	rui.SetParams(view, "touchEventsTest", rui.Params{
		rui.TouchStart: func(v rui.View, event rui.TouchEvent) {
			addToLog("touch-start", event)
		},
		rui.TouchEnd: func(v rui.View, event rui.TouchEvent) {
			addToLog("touch-end", event)
		},
		rui.TouchCancel: func(v rui.View, event rui.TouchEvent) {
			addToLog("touch-cancel", event)
		},
		rui.TouchMove: func(v rui.View, event rui.TouchEvent) {
			addToLog("touch-move", event)
		},
	})

	return view
}
