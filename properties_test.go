package rui

/*
import (
	"testing"
)

func TestProperties(t *testing.T) {

	createTestLog(t, true)

	list := new(propertyList)
	list.init()

	if !list.Set("name", "abc") {
		t.Error(`list.Set("name", "abc") fail`)
	}

	if !list.Has("name") {
		t.Error(`list.Has("name") fail`)
	}

	v := list.Get("name")
	if v == nil {
		t.Error(`list.Get("name") fail`)
	}
	if text, ok := v.(string); ok {
		if text != "abc" {
			t.Error(`list.Get("name") != "abc"`)
		}
	} else {
		t.Error(`list.Get("name") is not string`)
	}

	sizeValues := []any{"@small", "auto", "10px", Pt(20), AutoSize()}
	for _, value := range sizeValues {
		if !list.setSizeProperty("size", value) {
			t.Errorf(`setSizeProperty("size", %v) fail`, value)
		}
	}

	failSizeValues := []any{"@small,big", "abc", "10", Color(20), 100}
	for _, value := range failSizeValues {
		if list.setSizeProperty("size", value) {
			t.Errorf(`setSizeProperty("size", %v) success`, value)
		}
	}

	angleValues := []any{"@angle", "2pi", "π", "3deg", "60°", Rad(1.5), Deg(45), 1, 1.5}
	for _, value := range angleValues {
		if !list.setAngleProperty("angle", value) {
			t.Errorf(`setAngleProperty("angle", %v) fail`, value)
		}
	}

	failAngleValues := []any{"@angle,2", "pi32", "deg", "60°x", Color(0xFFFFFFFF)}
	for _, value := range failAngleValues {
		if list.setAngleProperty("angle", value) {
			t.Errorf(`setAngleProperty("angle", %v) success`, value)
		}
	}

	colorValues := []any{"@color", "#FF234567", "#234567", "rgba(30%, 128, 0.5, .25)", "rgb(30%, 128, 0.5)", Color(0xFFFFFFFF), 0xFFFFFFFF, White}
	for _, color := range colorValues {
		if !list.setColorProperty("color", color) {
			t.Errorf(`list.setColorProperty("color", %v) fail`, color)
		}
	}

	failColorValues := []any{"@color|2", "#FF234567FF", "#TT234567", "rgba(500%, 128, 10.5, .25)", 10.6}
	for _, color := range failColorValues {
		if list.setColorProperty("color", color) {
			t.Errorf(`list.setColorProperty("color", %v) success`, color)
		}
	}

	enumValues := []any{"@enum", "inherit", "on", Inherit, 2}
	inheritOffOn := inheritOffOnValues()
	for _, value := range enumValues {
		if !list.setEnumProperty("enum", value, inheritOffOn) {
			t.Errorf(`list.setEnumProperty("enum", %v, %v) fail`, value, inheritOffOn)
		}
	}

	failEnumValues := []any{"@enum 13", "inherit2", "onn", -1, 10}
	for _, value := range failEnumValues {
		if list.setEnumProperty("enum", value, inheritOffOn) {
			t.Errorf(`list.setEnumProperty("enum", %v, %v) success`, value, inheritOffOn)
		}
	}

	boolValues := []any{"@bool", "true", "yes ", "on", " 1", "false", "no", "off", "0", 0, 1, false, true}
	for _, value := range boolValues {
		if !list.setBoolProperty("bool", value) {
			t.Errorf(`list.setBoolProperty("bool", %v) fail`, value)
		}
	}

	failBoolValues := []any{"@bool,2", "tr", "ys", "10", -1, 10, 0.8}
	for _, value := range failBoolValues {
		if list.setBoolProperty("bool", value) {
			t.Errorf(`list.setBoolProperty("bool", %v) success`, value)
		}
	}

	intValues := []any{"@int", " 100", "-10 ", 0, 250}
	for _, value := range intValues {
		if !list.setIntProperty("int", value) {
			t.Errorf(`list.setIntProperty("int", %v) fail`, value)
		}
	}

	failIntValues := []any{"@int|10", "100i", "-1.0 ", 0.0}
	for _, value := range failIntValues {
		if list.setIntProperty("int", value) {
			t.Errorf(`list.setIntProperty("int", %v) success`, value)
		}
	}

	floatValues := []any{"@float", " 100.25", "-1.5e12 ", uint(0), 250, float32(10.2), float64(0)}
	for _, value := range floatValues {
		if !list.setFloatProperty("float", value) {
			t.Errorf(`list.setFloatProperty("float", %v) fail`, value)
		}
	}

	failFloatValues := []any{"@float|2", " 100.25i", "-1.5ee12 ", "abc"}
	for _, value := range failFloatValues {
		if list.setFloatProperty("float", value) {
			t.Errorf(`list.setFloatProperty("float", %v) success`, value)
		}
	}

	boundsValues := []any{"@bounds", "10px,20pt,@bottom,0", Em(2), []any{"@top", Px(10), AutoSize(), "14pt"}}
	for _, value := range boundsValues {
		if !list.setBoundsProperty("margin", value) {
			t.Errorf(`list.setBoundsProperty("margin", %v) fail`, value)
		}
	}
}
*/
