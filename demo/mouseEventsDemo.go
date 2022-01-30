package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anoshenko/rui"
)

const mouseEventsDemoText = `
GridLayout {
	width = 100%, height = 100%, cell-height = "1fr, auto",
	content = [
		GridLayout {
			padding = 12px,
			content = [
				GridLayout {
					id = mouseEventsTest, cell-horizontal-align = center, cell-vertical-align = center,
					height = 100%,
					border = _{ style = solid, width = 1px, color = gray},
					content = [
						TextView {
							id = mouseEventsText, text = "Test box",
						}
					]
				}
			],
		},
		Resizable {
			row = 1, side = top, background-color = lightgrey, height = 200px,
			content = EditView {
				id = mouseEventsLog, type = multiline, read-only = true, wrap = true,
			}
		},
	]
}
`

func createMouseEventsDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, mouseEventsDemoText)
	if view == nil {
		return nil
	}

	addToLog := func(tag string, event rui.MouseEvent) {
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

		appendInt := func(name string, value int) {
			buffer.WriteString(`, `)
			buffer.WriteString(name)
			buffer.WriteString(` = `)
			buffer.WriteString(strconv.Itoa(value))
		}

		appendPoint := func(name string, x, y float64) {
			buffer.WriteString(fmt.Sprintf(`, %s = (%g:%g)`, name, x, y))
		}

		buffer.WriteString(tag)
		buffer.WriteString(`: TimeStamp = `)
		buffer.WriteString(strconv.FormatUint(event.TimeStamp, 10))

		appendInt("Button", event.Button)
		appendInt("Buttons", event.Buttons)
		appendPoint("(X:Y)", event.X, event.Y)
		appendPoint("Client (X:Y)", event.ClientX, event.ClientY)
		appendPoint("Screen (X:Y)", event.ScreenX, event.ScreenY)
		appendBool("CtrlKey", event.CtrlKey)
		appendBool("ShiftKey", event.ShiftKey)
		appendBool("AltKey", event.AltKey)
		appendBool("MetaKey", event.MetaKey)
		buffer.WriteString(";\n\n")

		rui.AppendEditText(view, "mouseEventsLog", buffer.String())
		rui.ScrollViewToEnd(view, "mouseEventsLog")
	}

	rui.SetParams(view, "mouseEventsTest", rui.Params{
		rui.MouseDown: func(v rui.View, event rui.MouseEvent) {
			addToLog("mouse-down", event)
		},
		rui.MouseUp: func(v rui.View, event rui.MouseEvent) {
			addToLog("mouse-up", event)
		},
		rui.MouseOut: func(v rui.View, event rui.MouseEvent) {
			addToLog("mouse-out", event)
			rui.Set(view, "mouseEventsText", rui.Text, "Mouse out")
		},
		rui.MouseOver: func(v rui.View, event rui.MouseEvent) {
			addToLog("mouse-over", event)
		},
		rui.MouseMove: func(v rui.View, event rui.MouseEvent) {
			rui.Set(view, "mouseEventsText", rui.Text,
				fmt.Sprintf("(X:Y): (%g : %g)<br>Client (X:Y): (%g : %g)<br>Screen (X:Y): (%g : %g)",
					event.X, event.Y, event.ClientX, event.ClientY, event.ScreenX, event.ScreenY))
		},
	})

	return view
}
