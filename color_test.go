package rui

import (
	"bytes"
	"testing"
)

func TestColorARGB(t *testing.T) {
	color := Color(0x7FFE8743)
	a, r, g, b := color.ARGB()
	if a != 0x7F {
		t.Error("a != 0x7F")
	}
	if r != 0xFE {
		t.Error("r != 0xFE")
	}
	if g != 0x87 {
		t.Error("g != 0x87")
	}
	if b != 0x43 {
		t.Error("b != 0x43")
	}

	if color.Alpha() != 0x7F {
		t.Error("color.Alpha() != 0x7F")
	}

	if color.Red() != 0xFE {
		t.Error("color.Red() != 0xFE")
	}

	if color.Green() != 0x87 {
		t.Error("color.Green() != 0x87")
	}

	if color.Blue() != 0x43 {
		t.Error("color.Blue() != 0x43")
	}
}

func TestColorSetValue(t *testing.T) {
	createTestLog(t, true)

	testData := []struct{ src, result string }{
		{"#7F102040", "rgba(16,32,64,.50)"},
		{"#102040", "rgb(16,32,64)"},
		{"#8124", "rgba(17,34,68,.53)"},
		{"rgba(17,34,67,.5)", "rgba(17,34,67,.50)"},
		{"rgb(.25,50%,96)", "rgb(63,127,96)"},
		{"rgba(.25,50%,96,100%)", "rgb(63,127,96)"},
	}

	for _, data := range testData {
		color, ok := StringToColor(data.src)
		if !ok {
			t.Errorf(`color.SetValue("%s") fail`, data.src)
		}
		result := color.cssString()
		if result != data.result {
			t.Errorf(`color.cssString() = "%s", expected: "%s"`, result, data.result)
		}
	}
}

func TestColorWriteData(t *testing.T) {
	testCSS := func(t *testing.T, color Color, result string) {
		buffer := new(bytes.Buffer)
		buffer.WriteString(color.cssString())
		str := buffer.String()
		if str != result {
			t.Errorf("color = %#X, expected = \"%s\", result = \"%s\"", color, result, str)
		}
	}

	buffer := new(bytes.Buffer)
	color := Color(0x7FFE8743)
	color.writeData(buffer)
	str := buffer.String()
	if str != "#7FFE8743" {
		t.Errorf(`color = %#X, expected = "#7FFE8743", result = "%s"`, color, str)
	}

	testCSS(t, Color(0x7FFE8743), "rgba(254,135,67,.50)")
	testCSS(t, Color(0xFFFE8743), "rgb(254,135,67)")
	testCSS(t, Color(0x05FE8743), "rgba(254,135,67,.02)")
}

func TestColorSetData(t *testing.T) {
	test := func(t *testing.T, data string, result Color) {
		color, ok := StringToColor(data)
		if !ok {
			t.Errorf("data = \"%s\", fail result", data)
		} else if color != result {
			t.Errorf("data = \"%s\", expected = %#X, result = %#X", data, result, color)
		}
	}

	test(t, "#7Ffe8743", 0x7FFE8743)
	test(t, "#fE8743", 0xFFFE8743)
	test(t, "#AE43", 0xAAEE4433)
	test(t, "#E43", 0xFFEE4433)

	failData := []string{
		"",
		"7FfeG743",
		"#7Ffe87439",
		"#7FfeG743",
		"#7Ffe874",
		"#feG743",
		"#7Ffe8",
		"#fG73",
		"#GF3",
	}

	for _, data := range failData {
		if color, ok := StringToColor(data); ok {
			t.Errorf("data = \"%s\", success, result = %#X", data, color)
		}
	}
}
