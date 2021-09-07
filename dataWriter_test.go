package rui

/*
import (
	"testing"
)

func TestDataWriter(t *testing.T) {
	w := NewDataWriter()
	w.StartObject("root")
	w.WriteStringKey("key1", "text")
	w.WriteStringKey("key2", "text 2")
	w.WriteStringKey("key 3", "text4")
	w.WriteStringsKey("key4", []string{"text4.1", "text4.2", "text4.3"}, '|')
	w.WriteStringsKey("key5", []string{"text5.1", "text5.2", "text5.3"}, ',')
	w.WriteColorKey("color", Color(0x7FD18243))
	w.WriteColorsKey("colors", []Color{Color(0x7FD18243), Color(0xFF817263)}, ',')
	w.WriteIntKey("int", 43)
	w.WriteIntsKey("ints", []int{111, 222, 333}, '|')

	w.StartObjectKey("obj", "xxx")
	w.WriteSizeUnitKey("size", Px(16))
	w.WriteSizeUnitsKey("sizes", []SizeUnit{Px(8), Percent(100)}, ',')
	w.StartArray("array")
	w.WriteStringToArray("text")
	w.WriteColorToArray(Color(0x23456789))
	w.WriteIntToArray(1)
	w.WriteSizeUnitToArray(Inch(2))
	w.FinishArray()
	w.WriteBoundsKey("bounds1", Bounds{Px(8), Px(8), Px(8), Px(8)})
	w.WriteBoundsKey("bounds2", Bounds{Px(8), Pt(12), Mm(4.5), Inch(1.2)})
	w.FinishObject() // xxx

	w.FinishObject() // root

	text := w.String()
	expected := `root {
	key1 = text,
	key2 = "text 2",
	"key 3" = text4,
	key4 = text4.1|text4.2|text4.3,
	key5 = "text5.1,text5.2,text5.3",
	color = #7FD18243,
	colors = "#7FD18243,#FF817263",
	int = 43,
	ints = 111|222|333,
	obj = xxx {
		size = 16px,
		sizes = "8px,100%",
		array = [
			text,
			#23456789,
			1,
			2in
		],
		bounds1 = 8px,
		bounds2 = "8px,12pt,4.5mm,1.2in"
	}
}`

	if text != expected {
		t.Error("DataWriter test fail. Result:\n`" + text + "`\nExpected:\n`" + expected + "`")
	}

}
*/
