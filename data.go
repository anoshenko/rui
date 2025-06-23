package rui

import (
	"slices"
	"strings"
	"unicode"
)

// DataValue interface of a data node value
type DataValue interface {
	// IsObject returns "true" if data value is an object
	IsObject() bool

	// Object returns data value as a data object
	Object() DataObject

	// Value returns value as a string
	Value() string
}

// DataObject interface of a data object
type DataObject interface {
	DataValue

	// Tag returns data object tag
	Tag() string

	// PropertyCount returns properties count
	PropertyCount() int

	// Property returns a data node corresponding to a property with specific index
	Property(index int) DataNode

	// PropertyByTag returns a data node corresponding to a property tag
	PropertyByTag(tag string) DataNode

	// PropertyValue returns a string value of a property with a specific tag
	PropertyValue(tag string) (string, bool)

	// PropertyObject returns an object value of a property with a specific tag
	PropertyObject(tag string) DataObject

	// SetPropertyValue sets a string value of a property with a specific tag
	SetPropertyValue(tag, value string)

	// SetPropertyObject sets an object value of a property with a specific tag
	SetPropertyObject(tag string, object DataObject)

	// ToParams create a params(map) representation of a data object
	ToParams() Params
}

// DataNodeType defines the type of DataNode
type DataNodeType int

// Constants which are used to describe a node type, see [DataNode]
const (
	// TextNode - node is the pair "tag - text value". Syntax: <tag> = <text>
	TextNode DataNodeType = 0
	// ObjectNode - node is the pair "tag - object". Syntax: <tag> = <object name>{...}
	ObjectNode DataNodeType = 1
	// ArrayNode - node is the pair "tag - object". Syntax: <tag> = [...]
	ArrayNode DataNodeType = 2
)

// DataNode interface of a data node
type DataNode interface {
	// Tag returns a tag name
	Tag() string

	// Type returns a node type. Possible values are TextNode, ObjectNode and ArrayNode
	Type() DataNodeType

	// Text returns node text
	Text() string

	// Object returns node as object if that node type is an object
	Object() DataObject

	// ArraySize returns array size if that node type is an array
	ArraySize() int

	// ArrayElement returns a value of an array if that node type is an array
	ArrayElement(index int) DataValue

	// ArrayElements returns an array of objects if that node is an array
	ArrayElements() []DataValue

	// ArrayAsParams returns an array of a params(map) if that node is an array
	ArrayAsParams() []Params
}

/******************************************************************************/
type dataStringValue struct {
	value string
}

func (value *dataStringValue) Value() string {
	return value.value
}

func (value *dataStringValue) IsObject() bool {
	return false
}

func (value *dataStringValue) Object() DataObject {
	return nil
}

/******************************************************************************/
type dataObject struct {
	tag      string
	property []DataNode
}

// NewDataObject create new DataObject with the tag and empty property list
func NewDataObject(tag string) DataObject {
	obj := new(dataObject)
	obj.tag = tag
	obj.property = []DataNode{}
	return obj
}

func (object *dataObject) Value() string {
	return ""
}

func (object *dataObject) IsObject() bool {
	return true
}

func (object *dataObject) Object() DataObject {
	return object
}

func (object *dataObject) Tag() string {
	return object.tag
}

func (object *dataObject) PropertyCount() int {
	if object.property != nil {
		return len(object.property)
	}
	return 0
}

func (object *dataObject) Property(index int) DataNode {
	if object.property == nil || index < 0 || index >= len(object.property) {
		return nil
	}
	return object.property[index]
}

func (object *dataObject) PropertyByTag(tag string) DataNode {
	if object.property != nil {
		for _, node := range object.property {
			if node.Tag() == tag {
				return node
			}
		}
	}
	return nil
}

func (object *dataObject) PropertyValue(tag string) (string, bool) {
	if node := object.PropertyByTag(tag); node != nil && node.Type() == TextNode {
		return node.Text(), true
	}
	return "", false
}

func (object *dataObject) PropertyObject(tag string) DataObject {
	if node := object.PropertyByTag(tag); node != nil && node.Type() == ObjectNode {
		return node.Object()
	}
	return nil
}

func (object *dataObject) setNode(node DataNode) {
	if len(object.property) == 0 {
		object.property = []DataNode{node}
	} else {
		tag := node.Tag()
		for i, p := range object.property {
			if p.Tag() == tag {
				object.property[i] = node
				return
			}
		}

		object.property = append(object.property, node)
	}
}

// SetPropertyValue - set a string property with tag by value
func (object *dataObject) SetPropertyValue(tag, value string) {
	val := new(dataStringValue)
	val.value = value
	node := new(dataNode)
	node.tag = tag
	node.value = val
	object.setNode(node)
}

// SetPropertyObject - set a property with tag by object
func (object *dataObject) SetPropertyObject(tag string, obj DataObject) {
	node := new(dataNode)
	node.tag = tag
	node.value = obj
	object.setNode(node)
}

func (object *dataObject) ToParams() Params {
	params := Params{}
	for _, node := range object.property {
		switch node.Type() {
		case TextNode:
			if text := node.Text(); text != "" {
				params[PropertyName(node.Tag())] = text
			}

		case ObjectNode:
			if obj := node.Object(); obj != nil {
				params[PropertyName(node.Tag())] = node.Object()
			}

		case ArrayNode:
			array := []any{}
			for i := range node.ArraySize() {
				if data := node.ArrayElement(i); data != nil {
					if data.IsObject() {
						if obj := data.Object(); obj != nil {
							array = append(array, obj)
						}
					} else if text := data.Value(); text != "" {
						array = append(array, text)
					}
				}
			}
			if len(array) > 0 {
				params[PropertyName(node.Tag())] = array
			}
		}
	}

	return params
}

/******************************************************************************/
type dataNode struct {
	tag   string
	value DataValue
	array []DataValue
}

func (node *dataNode) Tag() string {
	return node.tag
}

func (node *dataNode) Type() DataNodeType {
	if node.array != nil {
		return ArrayNode
	}
	if node.value.IsObject() {
		return ObjectNode
	}
	return TextNode
}

func (node *dataNode) Text() string {
	if node.value != nil {
		return node.value.Value()
	}
	return ""
}

func (node *dataNode) Object() DataObject {
	if node.value != nil {
		return node.value.Object()
	}
	return nil
}

func (node *dataNode) ArraySize() int {
	if node.array != nil {
		return len(node.array)
	}
	return 0
}

func (node *dataNode) ArrayElement(index int) DataValue {
	if node.array != nil && index >= 0 && index < len(node.array) {
		return node.array[index]
	}
	return nil
}

func (node *dataNode) ArrayElements() []DataValue {
	if node.array != nil {
		return node.array
	}
	return []DataValue{}
}

func (node *dataNode) ArrayAsParams() []Params {
	result := []Params{}
	if node.array != nil {
		for _, data := range node.array {
			if data.IsObject() {
				if obj := data.Object(); obj != nil {
					if params := obj.ToParams(); len(params) > 0 {
						result = append(result, params)
					}
				}
			}
		}
	}
	return result
}

// ParseDataText - parse text and return DataNode
func ParseDataText(text string) DataObject {

	if strings.ContainsAny(text, "\r") {
		text = strings.ReplaceAll(text, "\r\n", "\n")
		text = strings.ReplaceAll(text, "\r", "\n")
	}
	data := append([]rune(text), rune(0))
	pos := 0
	size := len(data) - 1
	line := 1
	lineStart := 0

	skipSpaces := func(skipNewLine bool) {
		for pos < size {
			switch data[pos] {
			case '\n':
				if !skipNewLine {
					return
				}
				line++
				lineStart = pos + 1

			case '/':
				if pos+1 < size {
					switch data[pos+1] {
					case '/':
						pos += 2
						for pos < size && data[pos] != '\n' {
							pos++
						}
						pos--

					case '*':
						pos += 3
						for {
							if pos >= size {
								ErrorLog("Unexpected end of file")
								return
							}
							if data[pos-1] == '*' && data[pos] == '/' {
								break
							}
							if data[pos-1] == '\n' {
								line++
								lineStart = pos
							}
							pos++
						}

					default:
						return
					}
				}

			case ' ', '\t':
				// do nothing

			default:
				if !unicode.IsSpace(data[pos]) {
					return
				}
			}
			pos++
		}
	}

	parseTag := func() (string, bool) {
		skipSpaces(true)
		startPos := pos
		switch data[pos] {
		case '`':
			pos++
			startPos++
			for data[pos] != '`' {
				pos++
				if pos >= size {
					ErrorLog("Unexpected end of text")
					return string(data[startPos:size]), false
				}
			}
			str := string(data[startPos:pos])
			pos++
			return str, true

		case '\'', '"':
			stopSymbol := data[pos]
			pos++
			startPos++
			slash := false
			for stopSymbol != data[pos] {
				if data[pos] == '\\' {
					pos += 2
					slash = true
				} else {
					pos++
				}
				if pos >= size {
					ErrorLog("Unexpected end of text")
					return string(data[startPos:size]), false
				}
			}

			if !slash {
				str := string(data[startPos:pos])
				pos++
				skipSpaces(false)
				return str, true
			}

			buffer := make([]rune, pos-startPos+1)
			n1 := 0
			n2 := startPos

			invalidEscape := func() (string, bool) {
				str := string(data[startPos:pos])
				pos++
				ErrorLogF("Invalid escape sequence in \"%s\" (position %d)", str, n2-2-startPos)
				return str, false
			}

			for n2 < pos {
				if data[n2] != '\\' {
					buffer[n1] = data[n2]
					n2++
				} else {
					n2 += 2
					switch data[n2-1] {
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
						if n2+2 > pos {
							return invalidEscape()
						}
						x := 0
						for range 2 {
							ch := data[n2]
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
						if n2+4 > pos {
							return invalidEscape()
						}
						x := 0
						for range 4 {
							ch := data[n2]
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
						str := string(data[startPos:pos])
						ErrorLogF("Invalid escape sequence in \"%s\" (position %d)", str, n2-2-startPos)
						return str, false
					}
				}
				n1++
			}

			pos++
			skipSpaces(false)
			return string(buffer[0:n1]), true
		}

		stopSymbol := func(symbol rune) bool {
			return unicode.IsSpace(symbol) ||
				slices.Contains([]rune{'=', '{', '}', '[', ']', ',', ' ', '\t', '\n', '\'', '"', '`', '/'}, symbol)
		}

		for pos < size && !stopSymbol(data[pos]) {
			pos++
		}

		endPos := pos
		skipSpaces(false)
		if startPos == endPos {
			//ErrorLog("empty tag")
			return "", true
		}
		return string(data[startPos:endPos]), true
	}

	var parseObject func(tag string) DataObject
	var parseArray func() []DataValue

	parseNode := func() DataNode {
		var tag string
		var ok bool

		if tag, ok = parseTag(); !ok {
			return nil
		}

		skipSpaces(true)
		if data[pos] != '=' {
			ErrorLogF("expected '=' after a tag name (line: %d, position: %d)", line, pos-lineStart)
			return nil
		}

		pos++
		skipSpaces(true)
		switch data[pos] {
		case '[':
			node := new(dataNode)
			node.tag = tag

			if node.array = parseArray(); node.array == nil {
				return nil
			}
			return node

		case '{':
			node := new(dataNode)
			node.tag = tag
			if node.value = parseObject("_"); node.value == nil {
				return nil
			}
			return node

		case '}', ']', '=':
			ErrorLogF("Expected '[', '{' or a tag name after '=' (line: %d, position: %d)", line, pos-lineStart)
			return nil

		default:
			var str string
			if str, ok = parseTag(); !ok {
				return nil
			}

			node := new(dataNode)
			node.tag = tag

			if data[pos] == '{' {
				if node.value = parseObject(str); node.value == nil {
					return nil
				}
			} else {
				val := new(dataStringValue)
				val.value = str
				node.value = val
			}

			return node
		}
	}

	parseObject = func(tag string) DataObject {
		if data[pos] != '{' {
			ErrorLogF("Expected '{' (line: %d, position: %d)", line, pos-lineStart)
			return nil
		}
		pos++

		obj := new(dataObject)
		obj.tag = tag
		obj.property = []DataNode{}

		for pos < size {
			var node DataNode

			skipSpaces(true)
			if data[pos] == '}' {
				pos++
				skipSpaces(false)
				return obj
			}

			if node = parseNode(); node == nil {
				return nil
			}
			obj.property = append(obj.property, node)
			if data[pos] == '}' {
				pos++
				skipSpaces(true)
				return obj
			} else if data[pos] != ',' && data[pos] != '\n' {
				ErrorLogF(`Expected '}', '\n' or ',' (line: %d, position: %d)`, line, pos-lineStart)
				return nil
			}
			if data[pos] != '\n' {
				pos++
			}
			skipSpaces(true)
			for data[pos] == ',' {
				pos++
				skipSpaces(true)
			}
		}

		ErrorLog("Unexpected end of text")
		return nil
	}

	parseArray = func() []DataValue {
		pos++
		skipSpaces(true)

		array := []DataValue{}

		for pos < size {
			var tag string
			var ok bool

			skipSpaces(true)
			for data[pos] == ',' && pos < size {
				pos++
				skipSpaces(true)
			}

			if pos >= size {
				break
			}

			if data[pos] == ']' {
				pos++
				skipSpaces(true)
				return array
			}

			if tag, ok = parseTag(); !ok {
				return nil
			}

			if data[pos] == '{' {
				obj := parseObject(tag)
				if obj == nil {
					return nil
				}
				array = append(array, obj)
			} else {
				val := new(dataStringValue)
				val.value = tag
				array = append(array, val)
			}

			switch data[pos] {
			case ']', ',', '\n':

			default:
				ErrorLogF("Expected ']' or ',' (line: %d, position: %d)", line, pos-lineStart)
				return nil
			}

			/*
				if data[pos] == ']' {
					pos++
					skipSpaces()
					return array, nil
				} else if data[pos] != ',' {
					return nil, fmt.Errorf("Expected ']' or ',' (line: %d, position: %d)", line, pos-lineStart)
				}
				pos++
				skipSpaces()
			*/
		}

		ErrorLog("Unexpected end of text")
		return nil
	}

	if tag, ok := parseTag(); ok {
		return parseObject(tag)
	}
	return nil
}
