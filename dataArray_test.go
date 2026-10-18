package rui

import "testing"

func TestDataArray(t *testing.T) {
	array := NewTextDataArrayOf(0, 1, 2, 3, 4, 5)
	array.SetSingleLineString(true)

	if array.String() != "[0, 1, 2, 3, 4, 5]" {
		t.Error(`NewTextDataArrayOf(0, 1, 2, 3, 4, 5) fail. Result: ` + array.String())
	}

	array.AppendInt(6, 7)
	if array.String() != "[0, 1, 2, 3, 4, 5, 6, 7]" {
		t.Error(`array.AppendInt(6, 7) fail. Result: ` + array.String())
	}

	array.InsertInt(0, 8, 9)
	if array.String() != "[8, 9, 0, 1, 2, 3, 4, 5, 6, 7]" {
		t.Error(`array.InsertInt(0, 8, 9) fail. Result: ` + array.String())
	}

	array.InsertInt(2, 10, 11, 12)
	if array.String() != "[8, 9, 10, 11, 12, 0, 1, 2, 3, 4, 5, 6, 7]" {
		t.Error(`array.InsertInt(2, 10, 11, 12) fail. Result: ` + array.String())
	}

	array.InsertInt(100, 13, 14)
	if array.String() != "[8, 9, 10, 11, 12, 0, 1, 2, 3, 4, 5, 6, 7, 13, 14]" {
		t.Error(`array.InsertInt(100, 13, 14) fail. Result: ` + array.String())
	}

	array.Remove(0)
	if array.String() != "[9, 10, 11, 12, 0, 1, 2, 3, 4, 5, 6, 7, 13, 14]" {
		t.Error(`array.Remove(0) fail. Result: ` + array.String())
	}

	array.Remove(1)
	if array.String() != "[9, 11, 12, 0, 1, 2, 3, 4, 5, 6, 7, 13, 14]" {
		t.Error(`array.Remove(1) fail. Result: ` + array.String())
	}

	array.Remove(12)
	if array.String() != "[9, 11, 12, 0, 1, 2, 3, 4, 5, 6, 7, 13]" {
		t.Error(`array.Remove(12) fail. Result: ` + array.String())
	}

	array.RemoveN(0, 2)
	if array.String() != "[12, 0, 1, 2, 3, 4, 5, 6, 7, 13]" {
		t.Error(`array.RemoveN(0, 2) fail. Result: ` + array.String())
	}

	array.RemoveN(1, 3)
	if array.String() != "[12, 3, 4, 5, 6, 7, 13]" {
		t.Error(`array.RemoveN(1, 3) fail. Result: ` + array.String())
	}

	array.RemoveN(5, 12)
	if array.String() != "[12, 3, 4, 5, 6]" {
		t.Error(`array.RemoveN(5, 12) fail. Result: ` + array.String())
	}
}
