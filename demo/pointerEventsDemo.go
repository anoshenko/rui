package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anoshenko/rui"
)

const pointerEventsDemoText = `
GridLayout {
	width = 100%, height = 100%, cell-height = "1fr, auto",
	content = [
		GridLayout {
			padding = 12px,
			content = [
				GridLayout {
					id = pointerEventsTest, cell-horizontal-align = center, cell-vertical-align = center,
					border = _{ style = solid, width = 1px, color = gray},
					content = [
						TextView {
							id = pointerEventsText, text = "Test box",
						}
					]
				}
			],
		},
		Resizable {
			row = 1, side = top, background-color = lightgrey, height = 200px,
			content = EditView {
				id = pointerEventsLog, type = multiline, read-only = true, wrap = true,
			}
		},
	]
}
`

func createPointerEventsDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, pointerEventsDemoText)
	if view == nil {
		return nil
	}

	addToLog := func(tag string, event rui.PointerEvent) {
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

		appendFloat := func(name string, value float64) {
			buffer.WriteString(fmt.Sprintf(`, %s = %g`, name, value))
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
		appendFloat("Width", event.Width)
		appendFloat("Height", event.Height)
		appendFloat("Pressure", event.Pressure)
		appendFloat("TangentialPressure", event.TangentialPressure)
		appendFloat("TiltX", event.TiltX)
		appendFloat("TiltY", event.TiltY)
		appendFloat("Twist", event.Twist)

		buffer.WriteString(`, PointerType = `)
		buffer.WriteString(event.PointerType)

		appendBool("IsPrimary", event.IsPrimary)
		appendBool("CtrlKey", event.CtrlKey)
		appendBool("ShiftKey", event.ShiftKey)
		appendBool("AltKey", event.AltKey)
		appendBool("MetaKey", event.MetaKey)
		buffer.WriteString(";\n\n")

		rui.AppendEditText(view, "pointerEventsLog", buffer.String())
		rui.ScrollViewToEnd(view, "pointerEventsLog")
	}

	rui.SetParams(view, "pointerEventsTest", rui.Params{
		rui.PointerDown: func(v rui.View, event rui.PointerEvent) {
			addToLog("pointer-down", event)
		},
		rui.PointerUp: func(v rui.View, event rui.PointerEvent) {
			addToLog("pointer-up", event)
		},
		rui.PointerOut: func(v rui.View, event rui.PointerEvent) {
			addToLog("pointer-out", event)
			rui.Set(view, "pointerEventsText", rui.Text, "Pointer out")
		},
		rui.PointerOver: func(v rui.View, event rui.PointerEvent) {
			addToLog("pointer-over", event)
		},
		rui.PointerCancel: func(v rui.View, event rui.PointerEvent) {
			addToLog("pointer-cancel", event)
		},
		rui.PointerMove: func(v rui.View, event rui.PointerEvent) {
			rui.Set(view, "pointerEventsText", rui.Text,
				fmt.Sprintf("(X:Y): (%g : %g)<br>Client (X:Y): (%g : %g)<br>Screen (X:Y): (%g : %g)",
					event.X, event.Y, event.ClientX, event.ClientY, event.ScreenX, event.ScreenY))
		},
	})

	return view
}
