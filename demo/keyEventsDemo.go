package main

import (
	"strconv"
	"strings"

	"github.com/anoshenko/rui"
)

func keyEventHandle(view rui.View, event rui.KeyEvent, tag string) {
	var buffer strings.Builder

	buffer.WriteString(tag)
	buffer.WriteString(`: TimeStamp = `)
	buffer.WriteString(strconv.FormatUint(event.TimeStamp, 10))
	buffer.WriteString(`, Key = "`)
	buffer.WriteString(event.Key)
	buffer.WriteString(`", Code = "`)
	buffer.WriteString(event.Code)
	buffer.WriteString(`"`)

	appendBool := func(name string, value bool) {
		buffer.WriteString(`, `)
		buffer.WriteString(name)
		if value {
			buffer.WriteString(` = true`)
		} else {
			buffer.WriteString(` = false`)
		}
	}
	appendBool("Repeat", event.Repeat)
	appendBool("CtrlKey", event.CtrlKey)
	appendBool("ShiftKey", event.ShiftKey)
	appendBool("AltKey", event.AltKey)
	appendBool("MetaKey", event.MetaKey)
	buffer.WriteString(";\n\n")

	rui.AppendEditText(view, "", buffer.String())
	rui.ScrollViewToEnd(view, "")
}

func createKeyEventsDemo(session rui.Session) rui.View {
	return rui.NewEditView(session, rui.Params{
		rui.Width:        rui.Percent(100),
		rui.Height:       rui.Percent(100),
		rui.ReadOnly:     true,
		rui.EditWrap:     true,
		rui.Hint:         "Set the focus and press a key",
		rui.EditViewType: rui.MultiLineText,
		rui.KeyDownEvent: func(view rui.View, event rui.KeyEvent) {
			keyEventHandle(view, event, rui.KeyDownEvent)
		},
		rui.KeyUpEvent: func(view rui.View, event rui.KeyEvent) {
			keyEventHandle(view, event, rui.KeyUpEvent)
		},
	})
}
