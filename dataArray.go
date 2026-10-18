package rui

import (
	"iter"
	"slices"
	"strings"
)

// DataValue interface of a data node value
type DataArray interface {
	DataValue

	// Size returns returns the size of the array.
	Size() int

	// Element returns a value of an array element.
	Element(index int) DataValue

	// Slice returns an all element as a slice. The returned slice is a copy of the data.
	// Modifying the slice does not affect the internal structure of the DataArray.
	Slice() []DataValue

	// Elements returns an iterator to access the array elements
	Elements() iter.Seq[DataValue]

	// ElementsWithIndices returns an iterator to access the array elements with itheir indices
	ElementsWithIndices() iter.Seq2[int, DataValue]

	// TextElements returns an iterator for accessing the text elements of the array.
	// Elements of type DataArray and DataObject are skipped.
	TextElements() iter.Seq[string]

	// Append adds the specified elements to the end of the array.
	//
	// All nil elements are skipped (not added to the array).
	Append(values ...DataValue)

	// Append adds the specified elements to the end of the array.
	AppendText(values ...string)

	// Append adds the specified elements to the end of the array.
	AppendInt(values ...int)

	// Append adds the specified elements to the end of the array.
	AppendFloat(values ...float64)

	// Append adds the specified elements to the end of the array.
	AppendArray(array DataArray)

	// Insert inserts the specified elements into the array starting at the element with index.
	//
	// All nil elements are skipped (not added to the array).
	//
	// If index < 0 then the elements are inserted at the beginning of the array.
	//
	// If the index is greater than the array size, the elements are inserted at the end of the array.
	Insert(index int, values ...DataValue)

	// Insert inserts the specified elements into the array starting at the element with index.
	//
	// If index < 0 then the elements are inserted at the beginning of the array.
	//
	// If the index is greater than the array size, the elements are inserted at the end of the array.
	InsertText(index int, values ...string)

	// Insert inserts the specified elements into the array starting at the element with index.
	//
	// If index < 0 then the elements are inserted at the beginning of the array.
	//
	// If the index is greater than the array size, the elements are inserted at the end of the array.
	InsertInt(index int, values ...int)

	// Insert inserts the specified elements into the array starting at the element with index.
	//
	// If index < 0 then the elements are inserted at the beginning of the array.
	//
	// If the index is greater than the array size, the elements are inserted at the end of the array.
	InsertFloat(index int, values ...float64)

	// Remove removes the element at the specified index from the array and returns it. It returns nil if the index is out of bounds.
	Remove(index int) DataValue

	// RemoveN removes count elements from the array starting at the specified index and returns them.
	// It returns empty array if the index is out of bounds.
	RemoveN(index int, count int) []DataValue

	// SingleLineString returns how the result of the String() method will be formatted.
	//
	// If result is true, then all array elements are output on a single line (separated by ", ").
	//
	// If result is false, then each array element is output on a new line.
	SingleLineString() bool

	// SetSingleLineString Determines how the result of the String() method is formatted.
	//
	// If argument is true, then all array elements are output on a single line (separated by ", ").
	//
	// If argument is false, then each array element is output on a new line.
	SetSingleLineString(flag bool)
}

type dataArrayDataValue struct {
	array            []DataValue
	singleLineString bool
}

// NewDataArray creates an DataValue object of ArrayDataValue type.
//
// When creating the object, all nil elements are skipped (not added to the array).
func NewDataArray(array []DataValue) DataArray {
	data := new(dataArrayDataValue)

	if array == nil {
		data.array = []DataValue{}
	} else {
		data.array = make([]DataValue, 0, len(array))
		for _, element := range array {
			if element != nil {
				data.array = append(data.array, element)
			}
		}
	}

	return data
}

// NewDataArrayOf creates an DataValue object of ArrayDataValue type.
//
// When creating the object, all nil elements are skipped (not added to the array).
func NewDataArrayOf(array ...DataValue) DataArray {
	data := new(dataArrayDataValue)

	data.array = make([]DataValue, 0, len(array))
	for _, element := range array {
		if element != nil {
			data.array = append(data.array, element)
		}
	}

	return data
}

// NewTexDataArray creates a DataArray object (an DataValue object of ArrayDataValue type).
//
// All elements of the argument are converted to DataValue using the NewDataText function.
func NewTextDataArray[T string | int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 | bool | float32 | float64](array []T) DataArray {
	data := new(dataArrayDataValue)

	if array == nil {
		data.array = []DataValue{}
	} else {
		data.array = make([]DataValue, 0, len(array))
		for _, element := range array {
			data.array = append(data.array, NewTextData(element))
		}
	}

	if len(data.array) < 10 {
		data.singleLineString = true
	}

	return data
}

// NewTextDataArrayOf creates a DataArray object (an DataValue object of ArrayDataValue type).
//
// All arguments are converted to DataValue using the NewDataText function.
func NewTextDataArrayOf[T string | int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 | bool | float32 | float64](values ...T) DataArray {
	data := new(dataArrayDataValue)

	data.array = make([]DataValue, 0, len(values))
	for _, element := range values {
		data.array = append(data.array, NewTextData(element))
	}

	if len(data.array) < 10 {
		data.singleLineString = true
	}

	return data
}

func (value *dataArrayDataValue) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	value.writeToBuffer(buffer, "")

	return buffer.String()
}

func (value *dataArrayDataValue) writeToBuffer(buffer *strings.Builder, indent string) {

	if len(value.array) == 0 {
		buffer.WriteString("[]")
		return
	}

	indent2 := indent + "\t"

	if value.singleLineString {

		lead := "["
		for _, element := range value.array {
			buffer.WriteString(lead)
			element.writeToBuffer(buffer, indent2)
			lead = ", "
		}

	} else {

		buffer.WriteString("[\n")
		for _, element := range value.array {
			buffer.WriteString(indent2)
			element.writeToBuffer(buffer, indent2)
			buffer.WriteString(",\n")
		}

		buffer.WriteString(indent)
	}

	buffer.WriteString("]")
}

func (value *dataArrayDataValue) Type() DataValueType {
	return ArrayDataValue
}

func (value *dataArrayDataValue) Object() DataObject {
	return nil
}

func (value *dataArrayDataValue) Array() DataArray {
	return value
}

func (value *dataArrayDataValue) Size() int {
	return len(value.array)
}

func (value *dataArrayDataValue) Element(index int) DataValue {
	if index >= 0 && index < len(value.array) {
		return value.array[index]
	}
	return nil
}

func (value *dataArrayDataValue) Slice() []DataValue {
	return slices.Clone(value.array)
}

func (value *dataArrayDataValue) Elements() iter.Seq[DataValue] {
	return func(yield func(DataValue) bool) {
		for _, element := range value.array {
			if !yield(element) {
				return
			}
		}
	}
}

func (value *dataArrayDataValue) TextElements() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, element := range value.array {
			if element.Type() == TextDataValue {
				if !yield(element.String()) {
					return
				}
			}
		}
	}
}

func (value *dataArrayDataValue) ElementsWithIndices() iter.Seq2[int, DataValue] {
	return func(yield func(int, DataValue) bool) {
		for i, element := range value.array {
			if !yield(i, element) {
				return
			}
		}
	}
}

func (value *dataArrayDataValue) Append(data ...DataValue) {
	for _, d := range data {
		if d != nil {
			value.array = append(value.array, d)
		}
	}
}

func (value *dataArrayDataValue) AppendText(values ...string) {
	for _, v := range values {
		value.array = append(value.array, NewTextData(v))
	}
}

func (value *dataArrayDataValue) AppendInt(values ...int) {
	for _, v := range values {
		value.array = append(value.array, NewTextData(v))
	}
}

func (value *dataArrayDataValue) AppendFloat(values ...float64) {
	for _, v := range values {
		value.array = append(value.array, NewTextData(v))
	}
}

func (value *dataArrayDataValue) AppendArray(array DataArray) {
	if array != nil {
		for d := range array.Elements() {
			value.array = append(value.array, d)
		}
	}
}

func (value *dataArrayDataValue) Insert(index int, data ...DataValue) {
	if index >= len(value.array) {
		value.Append(data...)
		return
	}

	if index < 0 {
		index = 0
	}
	value.array = slices.Insert(value.array, index, data...)
}

func (value *dataArrayDataValue) InsertText(index int, values ...string) {
	data := make([]DataValue, len(values))
	for i, v := range values {
		data[i] = NewTextData(v)
	}
	value.Insert(index, data...)
}

func (value *dataArrayDataValue) InsertInt(index int, values ...int) {
	data := make([]DataValue, len(values))
	for i, v := range values {
		data[i] = NewTextData(v)
	}
	value.Insert(index, data...)
}

func (value *dataArrayDataValue) InsertFloat(index int, values ...float64) {
	data := make([]DataValue, len(values))
	for i, v := range values {
		data[i] = NewTextData(v)
	}
	value.Insert(index, data...)
}

func (value *dataArrayDataValue) Remove(index int) DataValue {
	size := len(value.array)
	if index < 0 || index >= size {
		return nil
	}

	result := value.array[index]
	switch index {
	case 0:
		value.array = value.array[1:]

	case size - 1:
		value.array = value.array[:index]

	default:
		value.array = append(value.array[:index], value.array[index+1:]...)
	}
	return result
}

func (value *dataArrayDataValue) RemoveN(index int, count int) []DataValue {
	size := len(value.array)
	if index < 0 || index >= size {
		return nil
	}

	if index+count >= size {
		result := value.array[index:]
		value.array = value.array[:index]
		return result
	}

	result := value.array[index : index+count]
	if index == 0 {
		value.array = value.array[count:]
	} else {
		value.array = append(value.array[:index], value.array[index+count:]...)
	}

	return result
}

func (value *dataArrayDataValue) SetSingleLineString(flag bool) {
	value.singleLineString = true
}

func (value *dataArrayDataValue) SingleLineString() bool {
	return value.singleLineString
}
