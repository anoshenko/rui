package main

import (
	"github.com/anoshenko/rui"
)

const filePickerDemoText = `
GridLayout {
	width = 100%, height = 100%, cell-height = "auto, 1fr", cell-width = "1fr, auto",
	content = [
		FilePicker {
			id = filePicker, accept = "txt, html"
		},
		Button {
			id = fileDownload, row = 0, column = 1, content = "Download file", disabled = true,
		}
		EditView {
			id = selectedFileData, row = 1, column = 0:1, type = multiline, read-only = true, wrap = true,
		}
	]
}
`

var downloadedFile []byte = nil
var downloadedFilename = ""

func createFilePickerDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, filePickerDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "filePicker", rui.FileSelectedEvent, func(picker rui.FilePicker, files []rui.FileInfo) {
		if len(files) > 0 {
			picker.LoadFile(files[0], func(file rui.FileInfo, data []byte) {
				if data != nil {
					downloadedFile = data
					downloadedFilename = files[0].Name
					rui.Set(view, "selectedFileData", rui.Text, string(data))
					rui.Set(view, "fileDownload", rui.Disabled, false)
				} else {
					rui.Set(view, "selectedFileData", rui.Text, rui.LastError())
				}
			})
		}
	})

	rui.Set(view, "fileDownload", rui.ClickEvent, func() {
		view.Session().DownloadFileData(downloadedFilename, downloadedFile)
	})

	return view
}
