package rui

/*
import (
	"testing"
)

func TestSizeUnitNew(t *testing.T) {
	_ = createTestSession(t)
	size := SizeUnit{SizeInPixel, 10}
	if Px(10) != size {
		t.Error("Px(10) error")
	}

	size = SizeUnit{SizeInPercent, 10}
	if Percent(10) != size {
		t.Error("Percent(10) error")
	}

	size = SizeUnit{SizeInPt, 10}
	if Pt(10) != size {
		t.Error("Pt(10) error")
	}

	size = SizeUnit{SizeInCM, 10}
	if Cm(10) != size {
		t.Error("Dip(10) error")
	}

	size = SizeUnit{SizeInMM, 10}
	if Mm(10) != size {
		t.Error("Mm(10) error")
	}

	size = SizeUnit{SizeInInch, 10}
	if Inch(10) != size {
		t.Error("Inch(10) error")
	}
}

func TestSizeUnitSet(t *testing.T) {
	_ = createTestSession(t)

	obj := new(dataObject)
	obj.SetPropertyValue("x", "20")
	obj.SetPropertyValue("size", "10mm")

	size := SizeUnit{Auto, 0}
	if size.setProperty(obj, "size", new(sessionData), nil) && (size.Type != SizeInMM || size.Value != 10) {
		t.Errorf("result: Type = %d, Value = %g", size.Type, size.Value)
	}
}

func TestSizeUnitSetValue(t *testing.T) {
	_ = createTestSession(t)

	type testData struct {
		text string
		size SizeUnit
	}

	testValues := []testData{
		testData{"auto", SizeUnit{Auto, 0}},
		testData{"1.5em", SizeUnit{SizeInEM, 1.5}},
		testData{"2ex", SizeUnit{SizeInEX, 2}},
		testData{"20px", SizeUnit{SizeInPixel, 20}},
		testData{"100%", SizeUnit{SizeInPercent, 100}},
		testData{"14pt", SizeUnit{SizeInPt, 14}},
		testData{"10pc", SizeUnit{SizeInPc, 10}},
		testData{"0.1in", SizeUnit{SizeInInch, 0.1}},
		testData{"10mm", SizeUnit{SizeInMM, 10}},
		testData{"90.5cm", SizeUnit{SizeInCM, 90.5}},
	}

	var size SizeUnit
	for _, data := range testValues {
		if size.SetValue(data.text) && size != data.size {
			t.Errorf("set \"%s\" result: Type = %d, Value = %g", data.text, size.Type, size.Value)
		}
	}

	failValues := []string{
		"xxx",
		"10.10.10px",
		"1000",
		"5km",
	}

	for _, text := range failValues {
		size.SetValue(text)
	}
}

func TestSizeUnitWriteData(t *testing.T) {
	_ = createTestSession(t)
	type testData struct {
		text string
		size SizeUnit
	}

	testValues := []testData{
		testData{"auto", SizeUnit{Auto, 0}},
		testData{"1.5em", SizeUnit{SizeInEM, 1.5}},
		testData{"2ex", SizeUnit{SizeInEX, 2}},
		testData{"20px", SizeUnit{SizeInPixel, 20}},
		testData{"100%", SizeUnit{SizeInPercent, 100}},
		testData{"14pt", SizeUnit{SizeInPt, 14}},
		testData{"10pc", SizeUnit{SizeInPc, 10}},
		testData{"0.1in", SizeUnit{SizeInInch, 0.1}},
		testData{"10mm", SizeUnit{SizeInMM, 10}},
		testData{"90.5cm", SizeUnit{SizeInCM, 90.5}},
	}

	buffer := new(bytes.Buffer)
	for _, data := range testValues {
		buffer.Reset()
		buffer.WriteString(data.size.String())
		str := buffer.String()
		if str != data.text {
			t.Errorf("result: \"%s\", expected: \"%s\"", str, data.text)
		}
	}
}
*/
