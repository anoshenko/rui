package rui

import (
	"fmt"
	"iter"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// DataObject interface of a data object
type DataObject interface {
	DataValue

	// Tag returns data object tag
	Tag() string

	// Properties() returns an iterator to access the properties
	Properties() iter.Seq2[string, DataValue]

	// TextProperties() returns an iterator to access the text properties. Objects and Arrays are skipped.
	TextProperties() iter.Seq2[string, string]

	// PropertyCount returns properties count
	PropertyCount() int

	// PropertyTags() returns an iterator to access the property tags
	PropertyTags() iter.Seq[string]

	// SetSortedTags specifies whether properties will be sorted by tag
	SetSortedTags(sorted bool)

	// Property returns a data node corresponding to the specified chain of tags.
	// If such a property does not exist, then nil is returned.
	//
	// A chain of tags is specified by the "tag" and "tags" arguments.
	// For example
	//
	//  p := object.Property("tag1", "tag2", "tag2")
	//
	// equivalent to the following code
	//
	//  var p DataValue
	//  p = nil
	//  if p1 := object.Property("tag1"); p1 != nil {
	//  	if o1 := p1.Object(); o1 != nil {
	//  		if p2 := object.Property("tag2"); p2 != nil {
	//  			if o2 := p2.Object(); o2 != nil {
	//  				p = object.Property("tag3")
	//  			}
	//  		}
	//  	}
	//  }
	Property(tag string, tags ...string) DataValue

	// PropertyValue returns a string value of a property with the specified chain of tags (see Property method).
	PropertyText(tag string, tags ...string) (string, bool)

	// PropertyObject returns an object value of a property with the specified chain of tags (see Property method).
	// If such a property does not exist or the value is not a DataObject, then nil is returned.
	PropertyObject(tag string, tags ...string) DataObject

	// PropertyObject returns an array value of a property with the specified chain of tags (see Property method).
	// If such a property does not exist or the value is not a DataArray, then nil is returned.
	PropertyArray(tag string, tags ...string) DataArray

	// SetProperty sets a value of a property with a specific tag
	SetProperty(tag string, value DataValue)

	// SetPropertyText sets a string value of a property with a specific tag
	SetPropertyText(tag, value string)

	// PropertyByTag removes a data node corresponding to a property tag and returns it
	RemoveProperty(tag string) DataValue

	// ToParams create a params(map) representation of a data object
	ToParams() Params

	// ToStruct fills the struct fields with property values of the DataObject.
	// A pointer to the struct to be filled is passed as an argument.
	ToStruct(data any)

	// WriteToFile converts the data object to text and write it to the named file, creating it if necessary.
	// If the file does not exist, WriteToFile creates it with permissions perm (before umask);
	// otherwise WriteFile truncates it before writing, without changing permissions.
	WriteToFile(path string, perm os.FileMode) error
}

type dataObjectProperty struct {
	tag   string
	value DataValue
}

type dataObject struct {
	tag      string
	sorted   bool
	property []dataObjectProperty
}

func (value *dataObject) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	value.writeToBuffer(buffer, "")

	return buffer.String()
}

func (object *dataObject) writeToBuffer(buffer *strings.Builder, indent string) {

	if len(object.tag) > 0 {
		writeRuiText(buffer, object.tag, true)
	}

	if len(object.property) == 0 {
		buffer.WriteString("{}")
		return
	}

	indent2 := indent + "\t"
	buffer.WriteString("{\n")

	for _, property := range object.property {
		buffer.WriteString(indent2)
		writeRuiText(buffer, property.tag, true)
		buffer.WriteString(" = ")
		property.value.writeToBuffer(buffer, indent2)
		buffer.WriteString(",\n")
	}

	buffer.WriteString(indent)
	buffer.WriteString("}")

}

func (object *dataObject) Type() DataValueType {
	return ObjectDataValue
}

func (object *dataObject) Object() DataObject {
	return object
}

func (value *dataObject) Array() DataArray {
	return nil
}

func (object *dataObject) Tag() string {
	return object.tag
}

func (object *dataObject) Properties() iter.Seq2[string, DataValue] {
	return func(yield func(string, DataValue) bool) {
		for _, property := range object.property {
			if !yield(property.tag, property.value) {
				return
			}
		}
	}
}

func (object *dataObject) TextProperties() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for _, property := range object.property {
			if property.value.Type() == TextDataValue {
				if !yield(property.tag, property.value.String()) {
					return
				}
			}
		}
	}
}

func (object *dataObject) PropertyCount() int {
	return len(object.property)
}

func (object *dataObject) PropertyTags() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, property := range object.property {
			if !yield(property.tag) {
				return
			}
		}
	}
}

func (object *dataObject) SetSortedTags(sorted bool) {
	object.sorted = sorted
	if sorted {
		slices.SortFunc(object.property, func(e1, e2 dataObjectProperty) int {
			return strings.Compare(e1.tag, e2.tag)
		})
	}
}

func (object *dataObject) indexByTag(tag string) int {
	if object.property != nil {
		return slices.IndexFunc(object.property, func(property dataObjectProperty) bool {
			return property.tag == tag
		})
	}
	return -1
}

func (object *dataObject) Property(tag string, tags ...string) DataValue {
	if index := object.indexByTag(tag); index >= 0 {
		property := object.property[index].value
		if len(tags) > 0 {
			for _, t := range tags {
				obj := property.Object()
				if obj == nil {
					return nil
				}

				property = obj.Property(t)
				if property == nil {
					return nil
				}
			}
		}
		return property

	}
	return nil
}

func (object *dataObject) RemoveProperty(tag string) DataValue {
	if object.property != nil {
		if index := object.indexByTag(tag); index >= 0 {
			node := object.property[index]

			switch index {
			case 0:
				object.property = object.property[1:]

			case len(object.property) - 1:
				object.property = object.property[:len(object.property)-1]

			default:
				object.property = append(object.property[:index], object.property[index+1:]...)
			}

			return node.value
		}
	}
	return nil
}

func (object *dataObject) PropertyText(tag string, tags ...string) (string, bool) {
	if node := object.Property(tag, tag); node != nil && node.Type() == TextDataValue {
		return node.String(), true
	}
	return "", false
}

func (object *dataObject) PropertyObject(tag string, tags ...string) DataObject {
	if node := object.Property(tag); node != nil {
		return node.Object()
	}
	return nil
}

func (object *dataObject) PropertyArray(tag string, tags ...string) DataArray {
	if node := object.Property(tag); node != nil {
		return node.Array()
	}
	return nil
}

func (object *dataObject) SetProperty(tag string, value DataValue) {
	if value == nil {
		object.RemoveProperty(tag)
		return
	}

	if index := object.indexByTag(tag); index >= 0 {
		object.property[index].value = value
	} else {
		object.property = append(object.property, dataObjectProperty{
			tag:   tag,
			value: value,
		})

		if object.sorted {
			slices.SortFunc(object.property, func(e1, e2 dataObjectProperty) int {
				return strings.Compare(e1.tag, e2.tag)
			})
		}
	}
}

// SetPropertyValue - set a string property with tag by value
func (object *dataObject) SetPropertyText(tag, text string) {
	object.SetProperty(tag, NewTextData(text))
}

func (object *dataObject) ToParams() Params {
	params := Params{}
	for _, node := range object.property {
		switch node.value.Type() {
		case TextDataValue:
			if text := node.value.String(); text != "" {
				params[PropertyName(node.tag)] = text
			}

		case ObjectDataValue:
			if obj := node.value.Object(); obj != nil {
				params[PropertyName(node.tag)] = obj
			}

		case ArrayDataValue:
			if nodeArray := node.value.Array(); nodeArray != nil {
				array := make([]any, 0, nodeArray.Size())
				for data := range nodeArray.Elements() {
					if data != nil {
						switch data.Type() {
						case TextDataValue:
							if text := data.String(); text != "" {
								array = append(array, text)
							}

						case ArrayDataValue:
							if dataArray := data.Array(); dataArray != nil && dataArray.Size() > 0 {
								array = append(array, dataArray)
							}

						default:
							array = append(array, data)
						}
					}
				}

				if len(array) > 0 {
					params[PropertyName(node.tag)] = array
				}
			}
		}
	}

	return params
}

func copyToReflectSlice(array DataArray, value reflect.Value) {
	val := reflect.New(value.Type())
	val = val.Elem()

	switch val.Type().Elem().Kind() {
	case reflect.Pointer:
		for data := range array.Elements() {
			s := reflect.New(val.Type().Elem()).Elem()
			copyToReflectValue(data, s)
			val = reflect.Append(val, s)
		}

	case reflect.String:
		for text := range array.TextElements() {
			val = reflect.Append(val, reflect.ValueOf(text))
		}

	case reflect.Int:
		for text := range array.TextElements() {
			if n, err := strconv.ParseInt(text, 10, 64); err == nil {
				val = reflect.Append(val, reflect.ValueOf(int(n)))
			}
		}

	case reflect.Int8:
		for text := range array.TextElements() {
			if n, err := strconv.ParseInt(text, 10, 8); err == nil {
				val = reflect.Append(val, reflect.ValueOf(int8(n)))
			}
		}

	case reflect.Int16:
		for text := range array.TextElements() {
			if n, err := strconv.ParseInt(text, 10, 16); err == nil {
				val = reflect.Append(val, reflect.ValueOf(int16(n)))
			}
		}

	case reflect.Int32:
		for text := range array.TextElements() {
			if n, err := strconv.ParseInt(text, 10, 32); err == nil {
				val = reflect.Append(val, reflect.ValueOf(int32(n)))
			}
		}

	case reflect.Int64:
		for text := range array.TextElements() {
			if n, err := strconv.ParseInt(text, 10, 64); err == nil {
				val = reflect.Append(val, reflect.ValueOf(n))
			}
		}

	case reflect.Uint:
		for text := range array.TextElements() {
			if n, err := strconv.ParseUint(text, 10, 64); err == nil {
				val = reflect.Append(val, reflect.ValueOf(uint(n)))
			}
		}

	case reflect.Uint8:
		for text := range array.TextElements() {
			if n, err := strconv.ParseUint(text, 10, 8); err == nil {
				val = reflect.Append(val, reflect.ValueOf(uint8(n)))
			}
		}

	case reflect.Uint16:
		for text := range array.TextElements() {
			if n, err := strconv.ParseUint(text, 10, 16); err == nil {
				val = reflect.Append(val, reflect.ValueOf(uint16(n)))
			}
		}

	case reflect.Uint32:
		for text := range array.TextElements() {
			if len(text) > 0 {
				if text[0] == '#' {
					if n, err := strconv.ParseUint(text[1:], 16, 32); err == nil {
						val = reflect.Append(val, reflect.ValueOf(uint32(n)))
					}
				} else if n, err := strconv.ParseUint(text, 10, 32); err == nil {
					val = reflect.Append(val, reflect.ValueOf(uint32(n)))
				}
			}
		}

	case reflect.Uint64:
		for text := range array.TextElements() {
			if n, err := strconv.ParseUint(text, 10, 64); err == nil {
				val = reflect.Append(val, reflect.ValueOf(n))
			}
		}

	case reflect.Float32:
		for text := range array.TextElements() {
			if n, err := strconv.ParseFloat(text, 32); err == nil {
				val = reflect.Append(val, reflect.ValueOf(float32(n)))
			}
		}

	case reflect.Float64:
		for text := range array.TextElements() {
			if n, err := strconv.ParseFloat(text, 64); err == nil {
				val = reflect.Append(val, reflect.ValueOf(n))
			}
		}

	case reflect.Bool:
		for text := range array.TextElements() {
			switch text {
			case "true", "1":
				val = reflect.Append(val, reflect.ValueOf(true))

			case "false", "0":
				val = reflect.Append(val, reflect.ValueOf(false))
			}
		}

	case reflect.Struct:
		for data := range array.Elements() {
			if obj := data.Object(); obj != nil {
				s := reflect.New(val.Type().Elem())
				if s.Kind() == reflect.Pointer {
					s = s.Elem()
				}
				copyToReflectStruct(obj, s)
				val = reflect.Append(val, s)
			}
		}

	case reflect.Slice:
		for data := range array.Elements() {
			if array := data.Array(); array != nil {
				s := reflect.New(val.Type().Elem())
				if s.Kind() == reflect.Pointer {
					s = s.Elem()
				}
				copyToReflectSlice(array, s)
				val = reflect.Append(val, s)
			}
		}

	case reflect.Array:
		for data := range array.Elements() {
			if array := data.Array(); array != nil {
				s := reflect.New(val.Type().Elem())
				if s.Kind() == reflect.Pointer {
					s = s.Elem()
				}
				copyToReflectArray(array, s)
				val = reflect.Append(val, s)
			}
		}

	case reflect.Map:
		for data := range array.Elements() {
			if obj := data.Object(); obj != nil {
				s := reflect.MakeMap(val.Type().Elem())
				copyToReflectMap(obj, s)
				val = reflect.Append(val, s)
			}
		}

	default:
		ErrorLog("unsupported slice element type")
	}

	value.Set(val)
}

func copyToReflectArray(array DataArray, val reflect.Value) {
	i := 0
	maxSize := val.Len()

	switch val.Type().Elem().Kind() {
	case reflect.Pointer:
		for data := range array.Elements() {
			if i >= maxSize {
				break
			}
			value := val.Index(i)
			copyToReflectValue(data, value)
			i++
		}

	case reflect.String:
		for text := range array.TextElements() {
			if i >= maxSize {
				break
			}
			value := val.Index(i)
			value.SetString(text)
			i++
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for text := range array.TextElements() {
			if i >= maxSize {
				break
			}
			if n, err := strconv.ParseInt(text, 10, 64); err == nil {
				value := val.Index(i)
				value.SetInt(n)
				i++
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		for text := range array.TextElements() {
			if i >= maxSize {
				break
			}
			if len(text) > 0 {
				if text[0] == '#' {
					if n, err := strconv.ParseUint(text[1:], 16, 64); err == nil {
						value := val.Index(i)
						value.SetUint(n)
						i++
					}
				} else if n, err := strconv.ParseUint(text, 10, 64); err == nil {
					value := val.Index(i)
					value.SetUint(n)
					i++
				}
			}
		}

	case reflect.Float32, reflect.Float64:
		for text := range array.TextElements() {
			if n, err := strconv.ParseFloat(text, 64); err == nil {
				if i >= maxSize {
					break
				}
				value := val.Index(i)
				value.SetFloat(n)
				i++
			}
		}

	case reflect.Bool:
		for text := range array.TextElements() {
			if i >= maxSize {
				break
			}

			switch text {
			case "true", "1":
				value := val.Index(i)
				value.SetBool(true)
				i++

			case "false", "0":
				value := val.Index(i)
				value.SetBool(false)
				i++
			}
		}

	case reflect.Struct:
		for data := range array.Elements() {
			if obj := data.Object(); obj != nil {
				if i >= maxSize {
					break
				}
				value := val.Index(i)
				copyToReflectStruct(obj, value)
				i++
			}
		}

	case reflect.Map:
		for data := range array.Elements() {
			if obj := data.Object(); obj != nil {
				if i >= maxSize {
					break
				}
				value := val.Index(i)
				copyToReflectMap(obj, value)
				i++
			}
		}

	case reflect.Slice:
		for data := range array.Elements() {
			if array := data.Array(); array != nil {
				if i >= maxSize {
					break
				}
				value := val.Index(i)
				copyToReflectSlice(array, value)
				i++
			}
		}

	case reflect.Array:
		for data := range array.Elements() {
			if array := data.Array(); array != nil {
				if i >= maxSize {
					break
				}
				value := val.Index(i)
				copyToReflectArray(array, value)
				i++
			}
		}

	default:
		ErrorLog("unsupported array element type")
	}
}

func copyToReflectMap(object DataObject, val reflect.Value) {

	mapValue := reflect.MakeMap(val.Type())

	for tag, data := range object.Properties() {
		key := reflect.New(mapValue.Type().Key())
		if key.Kind() == reflect.Pointer {
			key = key.Elem()
		}

		switch key.Kind() {
		case reflect.String:
			key.SetString(tag)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if n, err := strconv.ParseInt(tag, 10, 64); err == nil {
				key.SetInt(n)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if n, err := strconv.ParseUint(tag, 10, 64); err == nil {
				key.SetUint(n)
			}

		case reflect.Float32, reflect.Float64:
			if n, err := strconv.ParseFloat(tag, 64); err == nil {
				key.SetFloat(n)
			}

		case reflect.Bool:
			switch tag {
			case "true", "1":
				key.SetBool(true)

			case "false", "0":
				key.SetBool(false)
			}

		default:
			ErrorLog("unsupported map key type")
			return
		}

		value := reflect.New(mapValue.Type().Elem()).Elem()
		copyToReflectValue(data, value)
		mapValue.SetMapIndex(key, value)
	}

	val.Set(mapValue)
}

func copyToReflectValue(data DataValue, value reflect.Value) {
	if value.Kind() == reflect.Pointer {
		v := reflect.New(value.Type().Elem())
		value.Set(v)
		value = v.Elem()
	}

	switch value.Kind() {
	case reflect.Struct:
		if obj := data.Object(); obj != nil {
			copyToReflectStruct(obj, value)
		}

	case reflect.String:
		if data.Type() == TextDataValue {
			value.SetString(data.String())
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if data.Type() == TextDataValue {
			if n, err := strconv.ParseInt(data.String(), 10, 64); err == nil {
				value.SetInt(n)
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if data.Type() == TextDataValue {
			if text := data.String(); len(text) > 0 {
				if text[0] == '#' {
					if n, err := strconv.ParseUint(text[1:], 16, 64); err == nil {
						value.SetUint(n)
					}
				} else if n, err := strconv.ParseUint(text, 10, 64); err == nil {
					value.SetUint(n)
				}
			}
		}

	case reflect.Float32, reflect.Float64:
		if data.Type() == TextDataValue {
			if n, err := strconv.ParseFloat(data.String(), 64); err == nil {
				value.SetFloat(n)
			}
		}

	case reflect.Bool:
		if data.Type() == TextDataValue {
			switch data.String() {
			case "true", "1":
				value.SetBool(true)

			case "false", "0":
				value.SetBool(false)
			}
		}

	case reflect.Slice:
		if array := data.Array(); array != nil {
			copyToReflectSlice(array, value)
		}

	case reflect.Array:
		if array := data.Array(); array != nil {
			copyToReflectArray(array, value)
		}

	case reflect.Map:
		if obj := data.Object(); obj != nil {
			copyToReflectMap(obj, value)
		}

	default:
		ErrorLog("unsupported element type")
	}
}

func copyToReflectStruct(object DataObject, val reflect.Value) {
	for i := range val.NumField() {
		value := val.Field(i)
		field := val.Type().Field(i)
		fieldName := field.Name
		if tag := field.Tag.Get("rui"); tag != "" {
			params := parseFieldTag(tag)
			if name, ok := params["id"]; ok && name != "" {
				fieldName = name
			}
		}

		if value.CanSet() {
			if data := object.Property(fieldName); data != nil {
				copyToReflectValue(data, value)
			}
		} else {
			DebugLogF("field %s can not be set", fieldName)
		}
	}
}

func (object *dataObject) ToStruct(data any) {
	if data != nil {

		value := reflect.ValueOf(data)
		if value.Kind() == reflect.Pointer {
			value = value.Elem()
		}

		if value.Kind() == reflect.Struct {
			copyToReflectStruct(object, value)
		}
	}
}

func (object *dataObject) WriteToFile(path string, perm os.FileMode) error {
	text := object.String()
	return os.WriteFile(path, []byte(text), perm)
}

func reflectArrayToDataValue(value reflect.Value) DataValue {
	count := value.Len()
	if count == 0 {
		ErrorLog("empty array")
		return nil
	}

	kind := value.Type().Elem().Kind()
	if kind == reflect.Pointer {
		kind = value.Type().Elem().Elem().Kind()
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.String, reflect.Float32, reflect.Float64, reflect.Bool:

		array := make([]DataValue, 0, count)
		for i := range count {
			if e := reflectValueToDataValue(value.Index(i)); e != nil {
				array = append(array, e)
			}
		}

		result := NewDataArray(array)
		if count < 10 {
			result.SetSingleLineString(true)
		}

		return result

	case reflect.Struct, reflect.Array, reflect.Slice:

		array := make([]DataValue, 0, count)
		for i := range count {
			if e := reflectValueToDataValue(value.Index(i)); e != nil {
				array = append(array, e)
			}
		}
		return NewDataArray(array)
	}

	ErrorLog("unsupported array element type")
	return nil
}

func reflectMapToDataValue(value reflect.Value) DataValue {
	keys := value.MapKeys()

	count := len(keys)
	if count == 0 {
		return nil
	}

	obj := new(dataObject)
	obj.tag = "_"
	obj.property = make([]dataObjectProperty, 0, count)

	for _, key := range keys {
		tag := ""
		switch key.Kind() {
		case reflect.String:
			tag = key.String()

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			tag = strconv.FormatInt(key.Int(), 10)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			tag = strconv.FormatUint(key.Uint(), 10)

		case reflect.Float32, reflect.Float64:
			tag = fmt.Sprintf("%g", key.Float())

		case reflect.Bool:
			tag = fmt.Sprintf("%t", key.Bool())

		default:
			ErrorLog("unsupported map key type")
			return nil
		}

		if tag != "" {
			if data := reflectValueToDataValue(value.MapIndex(key)); data != nil {
				obj.property = append(obj.property, dataObjectProperty{tag: tag, value: data})
			}
		} else {
			DebugLog("empty map key")
		}
	}

	if len(obj.property) == 0 {
		DebugLog("empty map")
		return nil
	}

	return obj
}

func reflectValueToDataValue(value reflect.Value) DataValue {
	for value.Kind() == reflect.Pointer && !value.IsNil() {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Struct:
		return reflectStructToDataObject(value)

	case reflect.String:
		return NewTextData(value.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewTextData(value.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewTextData(value.Uint())

	case reflect.Float32, reflect.Float64:
		return NewTextData(value.Float())

	case reflect.Bool:
		return NewTextData(value.Bool())

	case reflect.Slice, reflect.Array:
		return reflectArrayToDataValue(value)

	case reflect.Map:
		return reflectMapToDataValue(value)
	}

	ErrorLog("unsupported element type")
	return nil
}

func parseFieldTag(tag string) map[string]string {
	result := map[string]string{}
	for line := range strings.SplitSeq(tag, ";") {
		pair := strings.SplitN(line, "=", 2)
		if len(pair) == 2 {
			key := strings.Trim(pair[0], " \t")
			value := strings.Trim(pair[1], " \t")
			if key != "" && value != "" {
				result[key] = value
			}
		}
	}
	return result
}

func reflectStructToDataObject(val reflect.Value) DataObject {
	obj := new(dataObject)
	obj.tag = val.Type().Name()
	obj.property = []dataObjectProperty{}

	for field, value := range val.Fields() {
		if data := reflectValueToDataValue(value); data != nil {
			fieldName := field.Name
			if tag := field.Tag.Get("rui"); tag != "" {
				params := parseFieldTag(tag)
				if name, ok := params["id"]; ok && name != "" {
					fieldName = name
				}
			}
			obj.property = append(obj.property, dataObjectProperty{tag: fieldName, value: data})
		}
	}

	return obj
}

// NewDataObject create new DataObject with the tag and empty property list
func NewDataObject(tag string) DataObject {
	obj := new(dataObject)
	obj.tag = tag
	obj.property = []dataObjectProperty{}
	return obj
}

// NewDataObjectFromStruct converts a struct to DataObject.
//
// A struct or a pointer to a struct is passed as a function argument.
// If other types of data are passed, the function returns nil.
//
// Depending on the type, the following transformation of structure fields is performed:
//   - string, int, int8...int64, uint, uint8...uint64, float32, float64, and bool are converted into DataText;
//   - struct and map are converted into DataObject;
//   - array and slice are converted into DataArray.
//
// The names of the corresponding structure fields are used as names.
// A name can be redefined using a tag `rui:"id=your_name"`. Example:
//
//	type A struct {
//		X int `rui:"id=y"`
//	}
func NewDataObjectFromStruct(data any) DataObject {
	if data == nil {
		return nil
	}

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if value.Kind() == reflect.Struct {
		return reflectStructToDataObject(value)
	}

	return nil
}
