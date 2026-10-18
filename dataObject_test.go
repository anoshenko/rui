package rui

import "testing"

type TestDataXY struct {
	X float32
	Y float64
}
type TestDataObjectStruct struct {
	Items      []string `rui:"id=items"`
	Items2     [2]string
	Int_a      []int `rui:"id=intArray"`
	Int_a2     [3]int
	Int8_a     []int8
	Int16_a    []int16
	Int32_a    []int32
	Int64_a    []int64
	Uint_a     []uint
	Uint8_a    []uint8
	Uint16_a   []uint16
	Uint32_a   []uint32
	Uint64_a   []uint64
	Float32_a  []float32
	Float32_a2 [3]float32
	Float64_a  []float64
	Bool_a     []bool
	Bool_a2    [2]bool
	Text       string `rui:"id=text"`
	Int        int    `rui:"id=int"`
	Int8       int    `rui:"id=int8"`
	Int16      int    `rui:"id=int16"`
	Int32      int    `rui:"id=int32"`
	Int64      int    `rui:"id=int64"`
	Uint       uint   `rui:"id=uint"`
	Uint8      uint   `rui:"id=uint8"`
	Uint16     uint   `rui:"id=uint16"`
	Uint32     uint   `rui:"id=uint32"`
	Uint64     uint   `rui:"id=uint64"`
	Float32    float32
	Float64    float64
	Flag       bool
	B          TestDataXY
	Struct_a   []TestDataXY
	Struct_a2  [2]TestDataXY
	Matrix     [][]int
	Matrix2    [2][]int
	Matrix3    [][3]int
	Matrix4    [2][3]int
	Map        map[string]int
}

func TestDataObject(t *testing.T) {

	data := TestDataObjectStruct{
		Items:      []string{"Item 1", "Item2"},
		Items2:     [2]string{"Item 3", "Item4"},
		Int_a:      []int{-1, 2, 3},
		Int_a2:     [3]int{0, 1, 2},
		Int8_a:     []int8{4, 5, 6},
		Int16_a:    []int16{7, 8, 9},
		Int32_a:    []int32{10, 11, 12},
		Int64_a:    []int64{13, 14, 15},
		Uint_a:     []uint{1, 2, 3},
		Uint8_a:    []uint8{4, 5, 6},
		Uint16_a:   []uint16{7, 8, 9},
		Uint32_a:   []uint32{10, 11, 12},
		Uint64_a:   []uint64{13, 14, 15},
		Float32_a:  []float32{1.5, 2, 3},
		Float32_a2: [3]float32{2.5, 3.5, 4},
		Float64_a:  []float64{1, 2.5, 3},
		Bool_a:     []bool{false, true},
		Bool_a2:    [2]bool{true, false},
		Text:       "text",
		Int:        -1,
		Int8:       0,
		Int16:      1,
		Int32:      2,
		Int64:      3,
		Uint:       5,
		Uint8:      6,
		Uint16:     7,
		Uint32:     8,
		Uint64:     9,
		Float32:    10.5,
		Float64:    11.5,
		Flag:       true,
		B:          TestDataXY{X: 4, Y: 5.5},
		Struct_a:   []TestDataXY{{X: 1, Y: 2}, {X: 3, Y: 4}},
		Struct_a2:  [2]TestDataXY{{X: 5, Y: 6}, {X: 7, Y: 8}},
		Matrix: [][]int{
			{1, 2, 3},
			{4, 5},
		},
		Matrix2: [2][]int{
			{21, 22, 23},
			{24, 25},
		},
		Matrix3: [][3]int{
			{31, 32, 33},
			{34, 35, 36},
		},
		Matrix4: [2][3]int{
			{41, 42, 43},
			{44, 45, 46},
		},
		Map: map[string]int{
			"key 1": 12,
			"key2":  25,
		},
	}

	obj := NewDataObjectFromStruct(&data)
	if obj == nil {
		t.Error(`NewDataObjectFromStruct fail`)
	}

	expected := `TestDataObjectStruct{
	items = ["Item 1", Item2],
	Items2 = ["Item 3", Item4],
	intArray = [-1, 2, 3],
	Int_a2 = [0, 1, 2],
	Int8_a = [4, 5, 6],
	Int16_a = [7, 8, 9],
	Int32_a = [10, 11, 12],
	Int64_a = [13, 14, 15],
	Uint_a = [1, 2, 3],
	Uint8_a = [4, 5, 6],
	Uint16_a = [7, 8, 9],
	Uint32_a = [10, 11, 12],
	Uint64_a = [13, 14, 15],
	Float32_a = [1.5, 2, 3],
	Float32_a2 = [2.5, 3.5, 4],
	Float64_a = [1, 2.5, 3],
	Bool_a = [false, true],
	Bool_a2 = [true, false],
	text = text,
	int = -1,
	int8 = 0,
	int16 = 1,
	int32 = 2,
	int64 = 3,
	uint = 5,
	uint8 = 6,
	uint16 = 7,
	uint32 = 8,
	uint64 = 9,
	Float32 = 10.5,
	Float64 = 11.5,
	Flag = true,
	B = TestDataXY{
		X = 4,
		Y = 5.5,
	},
	Struct_a = [
		TestDataXY{
			X = 1,
			Y = 2,
		},
		TestDataXY{
			X = 3,
			Y = 4,
		},
	],
	Struct_a2 = [
		TestDataXY{
			X = 5,
			Y = 6,
		},
		TestDataXY{
			X = 7,
			Y = 8,
		},
	],
	Matrix = [
		[1, 2, 3],
		[4, 5],
	],
	Matrix2 = [
		[21, 22, 23],
		[24, 25],
	],
	Matrix3 = [
		[31, 32, 33],
		[34, 35, 36],
	],
	Matrix4 = [
		[41, 42, 43],
		[44, 45, 46],
	],
	Map = _{
		"key 1" = 12,
		key2 = 25,
	},
}`

	if mapObject := obj.PropertyObject("Map"); mapObject != nil {
		mapObject.SetSortedTags(true)
	}
	if str := obj.String(); str != expected {
		t.Error("NewDataObjectFromStruct result is fail:\n" + str)
		return
	}

	var s TestDataObjectStruct
	obj.ToStruct(&s)

	obj2 := NewDataObjectFromStruct(s)

	if mapObject := obj2.PropertyObject("Map"); mapObject != nil {
		mapObject.SetSortedTags(true)
	}

	if str := obj2.String(); str != expected {
		t.Error("CopyToStruct result is fail:\n" + str)
	}

}
