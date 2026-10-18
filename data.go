package rui

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// DataValueType defines the type of DataNode
type DataValueType int

// Constants which are used to describe a value type, see [DataValue]
const (
	// TextDataValue - the DataValue value is a string.
	// The String() method of DataValue interface is used to access the value.
	TextDataValue DataValueType = 0

	// ObjectDataValue - the DataValue value is a object.
	// The Object() method of DataValue interface is used to access the object.
	ObjectDataValue DataValueType = 1

	// ArrayDataNode - the DataValue value is an array of DataValue.
	// The Array() method of DataValue interface is used to access the array.
	ArrayDataValue DataValueType = 2
)

// DataValue interface of a data node value
type DataValue interface {
	fmt.Stringer

	writeToBuffer(buffer *strings.Builder, indent string)

	// IsObject returns "true" if data value is an object
	Type() DataValueType

	// Object returns data value as a data object if Type() == ObjectDataValue, and nil otherwise.
	Object() DataObject

	// Array returns data value as a data array if Type() == ArrayDataValue, and nil otherwise.
	Array() DataArray
}

/******************************************************************************/
type dataStringValue struct {
	text string
}

// NewTextData creates an DataValue object of TextDataValue type.
//
// The argument is converted to a string using the following functions:
//   - int: strconv.FormatInt(value, 10)
//   - uint: strconv.FormatUint(value, 10)
//   - bool: fmt.Sprintf("%t", value)
//   - float: fmt.Sprintf("%g", value)
func NewTextData[T string | int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 | bool | float32 | float64](value T) DataValue {

	data := new(dataStringValue)
	switch value := any(value).(type) {
	case string:
		data.text = value

	case int:
		data.text = strconv.FormatInt(int64(value), 10)

	case int8:
		data.text = strconv.FormatInt(int64(value), 10)

	case int16:
		data.text = strconv.FormatInt(int64(value), 10)

	case int32:
		data.text = strconv.FormatInt(int64(value), 10)

	case int64:
		data.text = strconv.FormatInt(value, 10)

	case uint:
		data.text = strconv.FormatUint(uint64(value), 10)

	case uint8:
		data.text = strconv.FormatUint(uint64(value), 10)

	case uint16:
		data.text = strconv.FormatUint(uint64(value), 10)

	case uint32:
		data.text = strconv.FormatUint(uint64(value), 10)

	case uint64:
		data.text = strconv.FormatUint(uint64(value), 10)

	case bool:
		data.text = fmt.Sprintf("%t", value)

	case float32:
		data.text = fmt.Sprintf("%g", float64(value))

	case float64:
		data.text = fmt.Sprintf("%g", value)
	}

	return data
}

func (value *dataStringValue) String() string {
	return value.text
}

func (value *dataStringValue) writeToBuffer(buffer *strings.Builder, indent string) {
	writeRuiText(buffer, value.text, false)
}

func (value *dataStringValue) Type() DataValueType {
	return TextDataValue
}

func (value *dataStringValue) Object() DataObject {
	return nil
}

func (value *dataStringValue) Array() DataArray {
	return nil
}

/******************************************************************************/
type dataParser struct {
	data      []rune
	size      int
	pos       int
	line      int
	lineStart int
}

func (parser *dataParser) skipSpaces(skipNewLine bool) {
	for parser.pos < parser.size {
		switch parser.data[parser.pos] {
		case '\n':
			if !skipNewLine {
				return
			}
			parser.line++
			parser.lineStart = parser.pos + 1

		case '/':
			if parser.pos+1 < parser.size {
				switch parser.data[parser.pos+1] {
				case '/':
					parser.pos += 2
					for parser.pos < parser.size && parser.data[parser.pos] != '\n' {
						parser.pos++
					}
					parser.pos--

				case '*':
					parser.pos += 3
					for {
						if parser.pos >= parser.size {
							ErrorLog("Unexpected end of file")
							return
						}
						if parser.data[parser.pos-1] == '*' && parser.data[parser.pos] == '/' {
							break
						}
						if parser.data[parser.pos-1] == '\n' {
							parser.line++
							parser.lineStart = parser.pos
						}
						parser.pos++
					}

				default:
					return
				}
			}

		case ' ', '\t':
			// do nothing

		default:
			if !unicode.IsSpace(parser.data[parser.pos]) {
				return
			}
		}
		parser.pos++
	}
}

func (parser *dataParser) skipSpacesAndCommas() {
	parser.skipSpaces(true)
	for parser.data[parser.pos] == ',' && parser.pos < parser.size {
		parser.pos++
		parser.skipSpaces(true)
	}
}

func (parser *dataParser) parseTag() (string, error) {
	parser.skipSpaces(true)
	startPos := parser.pos
	switch parser.data[parser.pos] {
	case '`':
		parser.pos++
		startPos++
		for parser.data[parser.pos] != '`' {
			parser.pos++
			if parser.pos >= parser.size {
				return string(parser.data[startPos:parser.size]), errors.New("unexpected end of text")
			}
		}
		str := string(parser.data[startPos:parser.pos])
		parser.pos++
		return str, nil

	case '\'', '"':
		stopSymbol := parser.data[parser.pos]
		parser.pos++
		startPos++
		slash := false
		for stopSymbol != parser.data[parser.pos] {
			if parser.data[parser.pos] == '\\' {
				parser.pos += 2
				slash = true
			} else {
				parser.pos++
			}
			if parser.pos >= parser.size {
				return string(parser.data[startPos:parser.size]), errors.New("unexpected end of text")
			}
		}

		if !slash {
			str := string(parser.data[startPos:parser.pos])
			parser.pos++
			parser.skipSpaces(false)
			return str, nil
		}

		buffer := make([]rune, parser.pos-startPos+1)
		n1 := 0
		n2 := startPos

		invalidEscape := func() (string, error) {
			str := string(parser.data[startPos:parser.pos])
			parser.pos++
			return str, fmt.Errorf(`invalid escape sequence in "%s" (position %d)`, str, n2-2-startPos)
		}

		for n2 < parser.pos {
			if parser.data[n2] != '\\' {
				buffer[n1] = parser.data[n2]
				n2++
			} else {
				n2 += 2
				switch parser.data[n2-1] {
				case 'n':
					buffer[n1] = '\n'

				case 'r':
					buffer[n1] = '\r'

				case 't':
					buffer[n1] = '\t'

				case '"':
					buffer[n1] = '"'

				case '\'':
					buffer[n1] = '\''

				case '\\':
					buffer[n1] = '\\'

				case 'x', 'X':
					if n2+2 > parser.pos {
						return invalidEscape()
					}
					x := 0
					for range 2 {
						ch := parser.data[n2]
						if ch >= '0' && ch <= '9' {
							x = x*16 + int(ch-'0')
						} else if ch >= 'a' && ch <= 'f' {
							x = x*16 + int(ch-'a'+10)
						} else if ch >= 'A' && ch <= 'F' {
							x = x*16 + int(ch-'A'+10)
						} else {
							return invalidEscape()
						}
						n2++
					}
					buffer[n1] = rune(x)

				case 'u', 'U':
					if n2+4 > parser.pos {
						return invalidEscape()
					}
					x := 0
					for range 4 {
						ch := parser.data[n2]
						if ch >= '0' && ch <= '9' {
							x = x*16 + int(ch-'0')
						} else if ch >= 'a' && ch <= 'f' {
							x = x*16 + int(ch-'a'+10)
						} else if ch >= 'A' && ch <= 'F' {
							x = x*16 + int(ch-'A'+10)
						} else {
							return invalidEscape()
						}
						n2++
					}
					buffer[n1] = rune(x)

				default:
					str := string(parser.data[startPos:parser.pos])
					return str, fmt.Errorf(`invalid escape sequence in "%s" (position %d)`, str, n2-2-startPos)
				}
			}
			n1++
		}

		parser.pos++
		parser.skipSpaces(false)
		return string(buffer[0:n1]), nil
	}

	for parser.pos < parser.size && !parser.stopSymbol(parser.data[parser.pos]) {
		parser.pos++
	}

	endPos := parser.pos
	parser.skipSpaces(false)
	if startPos == endPos {
		//ErrorLog("empty tag")
		return "", nil
	}
	return string(parser.data[startPos:endPos]), nil
}

func (parser *dataParser) stopSymbol(symbol rune) bool {
	return unicode.IsSpace(symbol) ||
		slices.Contains([]rune{'=', '{', '}', '[', ']', ',', ' ', '\t', '\n', '\'', '"', '`', '/'}, symbol)
}

func (parser *dataParser) parseNode() (string, DataValue, error) {

	tag, err := parser.parseTag()
	if err != nil {
		return "", nil, err
	}

	parser.skipSpaces(true)
	if parser.data[parser.pos] != '=' {
		return "", nil, fmt.Errorf("expected '=' after a tag name (line: %d, position: %d)", parser.line, parser.pos-parser.lineStart)
	}

	parser.pos++
	parser.skipSpaces(true)
	switch parser.data[parser.pos] {
	case '[':
		array, err := parser.parseArray()
		if err != nil {
			return "", nil, err
		}
		return tag, NewDataArray(array), nil

	case '{':
		obj, err := parser.parseObject("_")
		if err != nil {
			return "", nil, err
		}
		return tag, obj, nil

	case '}', ']', '=':
		return "", nil, fmt.Errorf(`expected '[', '{' or a tag name after '=' (line: %d, position: %d)`, parser.line, parser.pos-parser.lineStart)

	default:
		text, err := parser.parseTag()
		if err != nil {
			return "", nil, err
		}

		if parser.data[parser.pos] == '{' {
			obj, err := parser.parseObject(text)
			if err != nil {
				return "", nil, err
			}
			return tag, obj, nil
		}

		return tag, NewTextData(text), nil
	}
}

func (parser *dataParser) parseObject(tag string) (DataObject, error) {
	if parser.data[parser.pos] != '{' {
		return nil, fmt.Errorf(`expected '{' (line: %d, position: %d)`, parser.line, parser.pos-parser.lineStart)
	}
	parser.pos++

	obj := NewDataObject(tag)

	for parser.pos < parser.size {
		parser.skipSpaces(true)
		if parser.data[parser.pos] == '}' {
			parser.pos++
			parser.skipSpaces(false)
			return obj, nil
		}

		property, value, err := parser.parseNode()
		if err != nil {
			return nil, err
		}
		obj.SetProperty(property, value)

		switch parser.data[parser.pos] {
		case '}':
			parser.pos++
			parser.skipSpaces(true)
			return obj, nil

		case ',':
			parser.pos++

		case '\n':
			// do nothing

		default:
			return nil, fmt.Errorf(`expected '}', '\n' or ',' (line: %d, position: %d)`, parser.line, parser.pos-parser.lineStart)
		}

		parser.skipSpacesAndCommas()
	}

	return nil, errors.New("unexpected end of text")
}

func (parser *dataParser) parseArray() ([]DataValue, error) {
	parser.pos++
	parser.skipSpaces(true)

	array := []DataValue{}

	parser.skipSpacesAndCommas()
	for parser.pos < parser.size {
		if parser.data[parser.pos] != ']' {
			if parser.data[parser.pos] == '[' {

				subarray, err := parser.parseArray()
				if err != nil {
					return nil, err
				}
				array = append(array, NewDataArray(subarray))

			} else {
				tag, err := parser.parseTag()
				if err != nil {
					return nil, err
				}

				if parser.data[parser.pos] == '{' {
					obj, err := parser.parseObject(tag)
					if err != nil {
						return nil, err
					}
					array = append(array, obj)
				} else {
					array = append(array, NewTextData(tag))
				}
			}
		}

		switch parser.data[parser.pos] {
		case ']':
			parser.pos++
			parser.skipSpaces(true)
			return array, nil

		case '\n':
			// do nothing

		case ',':
			parser.pos++

		default:
			return nil, fmt.Errorf(`expected ']' or ',' (line: %d, position: %d)`, parser.line, parser.pos-parser.lineStart)
		}

		parser.skipSpacesAndCommas()
	}

	return nil, errors.New("unexpected end of text")
}

// ParseDataText parses a text.
// If successful, it returns a DataObject and nil as the second element;
// otherwise, it returns nil and an error.
func ParseDataText(text string) (DataObject, error) {

	if strings.ContainsAny(text, "\r") {
		text = strings.ReplaceAll(text, "\r\n", "\n")
		text = strings.ReplaceAll(text, "\r", "\n")
	}

	parser := dataParser{
		data:      append([]rune(text), rune(0)),
		pos:       0,
		line:      1,
		lineStart: 0,
	}
	parser.size = len(parser.data) - 1

	tag, err := parser.parseTag()
	if err != nil {
		return nil, err
	}
	return parser.parseObject(tag)
}

// ReadDataObject reads a text file from the specified path, parses it.
// If successful, it returns a DataObject and nil as the second element;
// otherwise, it returns nil and an error.
func ReadDataObject(path string) (DataObject, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		debugLog(err.Error())
		return nil, err
	}

	return ParseDataText(string(data))
}
