package main

import (
	"github.com/anoshenko/rui"
)

const filePickerDemoText = `
GridLayout {
	width = 100%, height = 100%, cell-height = "auto, 1fr",
	content = [
		FilePicker {
			id = filePicker, accept = "txt, html"
		},
		EditView {
			id = selectedFileData, row = 1, type = multiline, read-only = true, wrap = true,
		}
	]
}
`

func createFilePickerDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, filePickerDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "filePicker", rui.FileSelectedEvent, func(picker rui.FilePicker, files []rui.FileInfo) {
		if len(files) > 0 {
			picker.LoadFile(files[0], func(file rui.FileInfo, data []byte) {
				if data != nil {
					rui.Set(view, "selectedFileData", rui.Text, string(data))
				} else {
					rui.Set(view, "selectedFileData", rui.Text, rui.LastError())
				}
			})
		}
	})
	return view
}
