package rui

/*
import (
	"bytes"
	"strconv"
	"testing"
)

func TestBoundsSet(t *testing.T) {

	session := createTestSession(t)

	obj := NewDataObject("Test")
	obj.SetPropertyValue("x", "10")
	obj.SetPropertyValue("padding", "8px")
	obj.SetPropertyValue("margins", "16mm,10pt,12in,auto")
	obj.SetPropertyValue("fail1", "x16mm")
	obj.SetPropertyValue("fail2", "16mm,10pt,12in")
	obj.SetPropertyValue("fail3", "x16mm,10pt,12in,auto")
	obj.SetPropertyValue("fail4", "16mm,x10pt,12in,auto")
	obj.SetPropertyValue("fail5", "16mm,10pt,x12in,auto")
	obj.SetPropertyValue("fail6", "16mm,10pt,12in,autoo")

	const failAttrsCount = 6

	var bounds Bounds
	if bounds.setProperty(obj, "padding", session) {
		if bounds.Left.Type != SizeInPixel || bounds.Left.Value != 8 ||
			bounds.Left != bounds.Right ||
			bounds.Left != bounds.Top ||
			bounds.Left != bounds.Bottom {
			t.Errorf("set padding error, result %v", bounds)
		}
	}

	if bounds.setProperty(obj, "margins", session) {
		if bounds.Top.Type != SizeInMM || bounds.Top.Value != 16 ||
			bounds.Right.Type != SizeInPt || bounds.Right.Value != 10 ||
			bounds.Bottom.Type != SizeInInch || bounds.Bottom.Value != 12 ||
			bounds.Left.Type != Auto {
			t.Errorf("set margins error, result %v", bounds)
		}
	}

	ignoreTestLog = true
	for i := 1; i <= failAttrsCount; i++ {
		if bounds.setProperty(obj, "fail"+strconv.Itoa(i), session) {
			t.Errorf("set 'fail' error, result %v", bounds)
		}
	}
	ignoreTestLog = false

	obj.SetPropertyValue("padding-left", "10mm")
	obj.SetPropertyValue("padding-top", "4pt")
	obj.SetPropertyValue("padding-right", "12in")
	obj.SetPropertyValue("padding-bottom", "8px")

	if bounds.setProperty(obj, "padding", session) {
		if bounds.Left.Type != SizeInMM || bounds.Left.Value != 10 ||
			bounds.Top.Type != SizeInPt || bounds.Top.Value != 4 ||
			bounds.Right.Type != SizeInInch || bounds.Right.Value != 12 ||
			bounds.Bottom.Type != SizeInPixel || bounds.Bottom.Value != 8 {
			t.Errorf("set margins error, result %v", bounds)
		}
	}

	for _, tag := range []string{"padding-left", "padding-top", "padding-right", "padding-bottom"} {
		if old, ok := obj.PropertyValue(tag); ok {
			ignoreTestLog = true
			obj.SetPropertyValue(tag, "x")
			if bounds.setProperty(obj, "padding", session) {
				t.Errorf("set \"%s\" value \"x\": result %v ", tag, bounds)
			}
			ignoreTestLog = false
			obj.SetPropertyValue(tag, old)
		}
	}
}

func TestBoundsWriteData(t *testing.T) {

	_ = createTestSession(t)

	bounds := Bounds{
		SizeUnit{SizeInPixel, 8},
		SizeUnit{SizeInInch, 10},
		SizeUnit{SizeInPt, 12},
		SizeUnit{Auto, 0},
	}

	buffer := new(bytes.Buffer)
	bounds.writeData(buffer)
	str := buffer.String()
	if str != `"8px,10in,12pt,auto"` {
		t.Errorf("result `%s`, expected `\"8px,10dip,12pt,auto\"`", str)
	}
}
*/
