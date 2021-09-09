package main

import (
	"fmt"

	"github.com/anoshenko/rui"
)

const popupDemoText = `
GridLayout {
	width = 100%, height = 100%, cell-height = "auto, 1fr",
	content = GridLayout {
		width = 100%, cell-width = "auto, 1fr",
		cell-vertical-align = center, gap = 8px,
		content = [
			Button {
				id = popupShowMessage, margin = 4px, content = "Show message",
			},
			Button {
				id = popupShowQuestion, row = 1, margin = 4px, content = "Show question",
			},
			TextView {
				id = popupShowQuestionResult, row = 1, column = 1, 
			},
			Button {
				id = popupShowCancellableQuestion, row = 2, margin = 4px, content = "Show cancellable question",
			},
			TextView {
				id = popupShowCancellableQuestionResult, row = 2, column = 1, 
			},
			Button {
				id = popupShowMenu, row = 3, margin = 4px, content = "Show menu",
			},
			TextView {
				id = popupShowMenuResult, row = 3, column = 1, 
			},
		]
	}
}
`

func createPopupDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, popupDemoText)
	if view == nil {
		return nil
	}

	rui.Set(view, "popupShowMessage", rui.ClickEvent, func() {
		rui.ShowMessage("Hello", "Hello world!!!", session)
	})

	rui.Set(view, "popupShowQuestion", rui.ClickEvent, func() {
		rui.ShowQuestion("Hello", "Are you alright?", session,
			func() {
				rui.Set(view, "popupShowQuestionResult", rui.Text, "Answer: Yes")
			},
			func() {
				rui.Set(view, "popupShowQuestionResult", rui.Text, "Answer: No")
			})
	})

	rui.Set(view, "popupShowCancellableQuestion", rui.ClickEvent, func() {
		rui.ShowCancellableQuestion("Hello", "Are you alright?", session,
			func() {
				rui.Set(view, "popupShowCancellableQuestionResult", rui.Text, "Answer: Yes")
			},
			func() {
				rui.Set(view, "popupShowCancellableQuestionResult", rui.Text, "Answer: No")
			},
			func() {
				rui.Set(view, "popupShowCancellableQuestionResult", rui.Text, "Answer: Cancel")
			})
	})

	rui.Set(view, "popupShowMenu", rui.ClickEvent, func() {
		rui.ShowMenu(session, rui.Params{
			rui.Items: []string{"Item 1", "Item 2", "Item 3", "Item 4"},
			rui.Title: "Menu",
			rui.PopupMenuResult: func(n int) {
				rui.Set(view, "popupShowMenuResult", rui.Text, fmt.Sprintf("Item %d selected", n+1))
			},
		})
	})

	return view
}
