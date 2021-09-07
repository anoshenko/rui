package rui

import (
	"testing"
)

var stopTestLogFlag = false
var testLogDone chan int
var ignoreTestLog = false

func createTestLog(t *testing.T, ignore bool) {
	ignoreTestLog = ignore
	SetErrorLog(func(text string) {
		if ignoreTestLog {
			t.Log(text)
		} else {
			t.Error(text)
		}
	})
	SetDebugLog(func(text string) {
		t.Log(text)
	})
}

/*
func createTestSession(t *testing.T) *sessionData {
	session := new(sessionData)
	createTestLog(t, false)
	return session
}

func TestSessionConstants(t *testing.T) {
	session := createTestSession(t)

	customTheme := `
	theme {
		colors = _{
			textColor = #FF080808,
			myColor = #81234567
		},
		colors:dark = _{
			textColor = #FFF0F0F0,
			myColor = #87654321
		},
		constants = _{
			defaultPadding = 10px,
			myConstant = 100%
			const1 = "@const2, 10px; @const3"
			const2 = "20mm / @const4"
			const3 = "@const5 : 30pt"
			const4 = "40%"
			const5 = "50px"
		},
		constants:touch = _{
			defaultPadding = 20px,
			myConstant = 80%,
		},
	}
	`

	SetErrorLog(func(text string) {
		t.Error(text)
	})

	theme, ok := newTheme(customTheme)
	if !ok {
		return
	}

	session.SetCustomTheme(theme)

	type constPair struct {
		tag, value string
	}

	testConstants := func(constants []constPair) {
		for _, constant := range constants {
			if value, ok := session.Constant(constant.tag); ok {
				if value != constant.value {
					t.Error(constant.tag + " = " + value + ". Need: " + constant.value)
				}
			}
		}
	}

	testConstants([]constPair{
		{tag: "defaultPadding", value: "10px"},
		{tag: "myConstant", value: "100%"},
		{tag: "buttonMargin", value: "4px"},
	})

	session.SetConstant("myConstant", "25px")

	testConstants([]constPair{
		{tag: "defaultPadding", value: "10px"},
		{tag: "myConstant", value: "25px"},
		{tag: "buttonMargin", value: "4px"},
	})

	session.touchScreen = true

	testConstants([]constPair{
		{tag: "defaultPadding", value: "20px"},
		{tag: "myConstant", value: "80%"},
		{tag: "buttonMargin", value: "4px"},
	})

	session.SetTouchConstant("myConstant", "30pt")

	testConstants([]constPair{
		{tag: "defaultPadding", value: "20px"},
		{tag: "myConstant", value: "30pt"},
		{tag: "buttonMargin", value: "4px"},
	})

	if value, ok := session.Constant("const1"); ok {
		if value != "20mm/40%,10px;50px:30pt" {
			t.Error("const1 = " + value + ". Need: 20mm/40%,10px;50px:30pt")
		}
	}
}
*/
