package main

import (
	"github.com/anoshenko/rui"
)

const absoluteLayoutDemoText = `
AbsoluteLayout {
	width = 100%, height = 100%,
	content = [ View { id = view1, width = 32px, height = 32px, left = 100px, top = 200px, background-color = #FF0000FF } ]
}
`

func createAbsoluteLayoutDemo(session rui.Session) rui.View {
	return rui.CreateViewFromText(session, absoluteLayoutDemoText)
}
