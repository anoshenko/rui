package rui

import (
	"testing"
)

func TestSizeFunc(t *testing.T) {
	session := new(sessionData)
	session.getCurrentTheme().SetConstant("a1", "120px", "120px")

	SetErrorLog(func(text string) {
		t.Error(text)
	})
	SetDebugLog(func(text string) {
		t.Log(text)
	})

	testFunc := func(fn SizeFunc, str, css string) {
		if fn != nil {
			if text := fn.String(); str != text {
				t.Error("String() error.\nResult:   \"" + text + "\"\nExpected: \"" + str + `"`)
			}
			if text := fn.cssString(session); css != text {
				t.Error("cssString() error.\nResult:   \"" + text + "\"\nExpected: \"" + css + `"`)
			}
		}
	}

	testFunc(MinSize("100%", Px(10)), `min(100%, 10px)`, `min(100%, 10px)`)
	testFunc(MaxSize(Percent(100), "@a1"), `max(100%, @a1)`, `max(100%, 120px)`)
	testFunc(SumSize(Percent(100), "@a1"), `sum(100%, @a1)`, `calc(100% + 120px)`)
	testFunc(SubSize(Percent(100), "@a1"), `sub(100%, @a1)`, `calc(100% - 120px)`)
	testFunc(MulSize(Percent(100), "@a1"), `mul(100%, @a1)`, `calc(100% * 120px)`)
	testFunc(DivSize(Percent(100), "@a1"), `div(100%, @a1)`, `calc(100% / 120px)`)
	testFunc(ClampSize(Percent(20), "@a1", Percent(40)), `clamp(20%, @a1, 40%)`, `clamp(20%, 120px, 40%)`)

	testFunc(MaxSize(SubSize(Percent(100), "@a1"), "@a1"), `max(sub(100%, @a1), @a1)`, `max(100% - 120px, 120px)`)

	testParse := func(str, css string) {
		if fn := parseSizeFunc(str); fn != nil {
			testFunc(fn, str, css)
		}
	}

	testParse(`min(100%, 10px)`, `min(100%, 10px)`)
	testParse(`max(100%, @a1)`, `max(100%, 120px)`)
	testParse(`max(sub(100%, @a1), @a1)`, `max(100% - 120px, 120px)`)
	testParse(`mul(sub(100%, @a1), @a1)`, `calc((100% - 120px) * 120px)`)
	testParse(`mul(sub(100%, @a1), div(mul(@a1, 3), 2))`, `calc((100% - 120px) * ((120px * 3) / 2))`)
}
