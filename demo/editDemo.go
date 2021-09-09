package main

import (
	"github.com/anoshenko/rui"
)

const editDemoText = `
ListLayout {
	width = 100%, height = 100%, orientation = vertical, padding = 16px,
	content = [
		GridLayout {
			row-gap = 12px, column-gap = 8px, cell-width = "auto, 1fr, auto", cell-vertical-align = center,
			content = [
				TextView { row = 0, column = 0, text = "User name" },
				EditView { row = 0, column = 1, id = editUserName, min-width = 200px, hint = "Required", type = text },
				TextView { row = 1, column = 0, text = "Password" },
				EditView { row = 1, column = 1, id = editPassword, min-width = 200px, hint = "8 characters minimum", type = password },
				TextView { row = 2, column = 0, text = "Confirm password" },
				EditView { row = 2, column = 1, id = editConfirmPassword, min-width = 200px, hint = "Required", type = password },
				TextView { row = 2, column = 2, id = confirmLabel, text = "" },
				TextView { row = 3, column = 0, text = "Main e-mail" },
				EditView { row = 3, column = 1, id = editMainEmail, min-width = 200px, hint = "Required", type = email },
				TextView { row = 4, column = 0, text = "Additional e-mails" },
				EditView { row = 4, column = 1, id = editAdditionalEmails, min-width = 200px, hint = "Optional", type = emails },
				TextView { row = 5, column = 0, text = "Home page" },
				EditView { row = 5, column = 1, id = editHomePage, min-width = 200px, hint = "Optional", type = url },
				TextView { row = 6, column = 0, text = "Phone" },
				EditView { row = 6, column = 1, id = editPhone, min-width = 200px, hint = "Optional", type = phone },
				EditView { row = 7, column = 0:1, id = editMultiLine, height = 200px, type = multiline },
				Checkbox { row = 7, column = 2, id = editMultiLineWrap, content = "Wrap", margin-left = 12px }
			]
		},
	]
}
`

func createEditDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, editDemoText)
	if view == nil {
		return nil
	}

	setConfirmLabel := func(password, confirmPassword string) {
		if password == confirmPassword {
			rui.Set(view, "confirmLabel", rui.TextColor, rui.Green)
			if password != "" {
				rui.Set(view, "confirmLabel", rui.Text, "✓")
			} else {
				rui.Set(view, "confirmLabel", rui.Text, "")
			}
		} else {
			rui.Set(view, "confirmLabel", rui.TextColor, rui.Red)
			rui.Set(view, "confirmLabel", rui.Text, "✗")
		}
	}

	rui.Set(view, "editPassword", rui.EditTextChangedEvent, func(edit rui.EditView, newText string) {
		setConfirmLabel(newText, rui.GetText(view, "editConfirmPassword"))
	})

	rui.Set(view, "editConfirmPassword", rui.EditTextChangedEvent, func(edit rui.EditView, newText string) {
		setConfirmLabel(rui.GetText(view, "editPassword"), newText)
	})

	rui.Set(view, "editMultiLineWrap", rui.CheckboxChangedEvent, func(checkbox rui.Checkbox, checked bool) {
		rui.Set(view, "editMultiLine", rui.Wrap, checked)
	})

	return view
}
