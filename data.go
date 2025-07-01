package rui

import (
	"errors"
	"fmt"
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

	// PropertyByTag removes a data node corresponding to a property tag and returns it
	RemovePropertyByTag(tag string) DataNode
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

func (object *dataObject) RemovePropertyByTag(tag string) DataNode {
	if object.property != nil {
		for i, node := range object.property {
			if node.Tag() == tag {
				switch i {
				case 0:
					object.property = object.property[1:]

				case len(object.property) - 1:
					object.property = object.property[:len(object.property)-1]

				default:
					object.property = append(object.property[:i], object.property[i+1:]...)
				}

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

func (parser *dataParser) parseNode() (DataNode, error) {
	var tag string
	var err error

	if tag, err = parser.parseTag(); err != nil {
		return nil, err
	}

	parser.skipSpaces(true)
	if parser.data[parser.pos] != '=' {
		return nil, fmt.Errorf("expected '=' after a tag name (line: %d, position: %d)", parser.line, parser.pos-parser.lineStart)
	}

	parser.pos++
	parser.skipSpaces(true)
	switch parser.data[parser.pos] {
	case '[':
		node := new(dataNode)
		node.tag = tag

		if node.array, err = parser.parseArray(); err != nil {
			return nil, err
		}
		return node, nil

	case '{':
		node := new(dataNode)
		node.tag = tag
		if node.value, err = parser.parseObject("_"); err != nil {
			return nil, err
		}
		return node, nil

	case '}', ']', '=':
		return nil, fmt.Errorf(`expected '[', '{' or a tag name after '=' (line: %d, position: %d)`, parser.line, parser.pos-parser.lineStart)

	default:
		var str string
		if str, err = parser.parseTag(); err != nil {
			return nil, err
		}

		node := new(dataNode)
		node.tag = tag

		if parser.data[parser.pos] == '{' {
			if node.value, err = parser.parseObject(str); err != nil {
				return nil, err
			}
		} else {
			val := new(dataStringValue)
			val.value = str
			node.value = val
		}

		return node, nil
	}
}

func (parser *dataParser) parseObject(tag string) (DataObject, error) {
	if parser.data[parser.pos] != '{' {
		return nil, fmt.Errorf(`expected '{' (line: %d, position: %d)`, parser.line, parser.pos-parser.lineStart)
	}
	parser.pos++

	obj := new(dataObject)
	obj.tag = tag
	obj.property = []DataNode{}

	for parser.pos < parser.size {
		parser.skipSpaces(true)
		if parser.data[parser.pos] == '}' {
			parser.pos++
			parser.skipSpaces(false)
			return obj, nil
		}

		node, err := parser.parseNode()
		if err != nil {
			return nil, err
		}

		obj.property = append(obj.property, node)
		if parser.data[parser.pos] == '}' {
			parser.pos++
			parser.skipSpaces(true)
			return obj, nil
		} else if parser.data[parser.pos] != ',' && parser.data[parser.pos] != '\n' {
			return nil, fmt.Errorf(`expected '}', '\n' or ',' (line: %d, position: %d)`, parser.line, parser.pos-parser.lineStart)
		}

		if parser.data[parser.pos] != '\n' {
			parser.pos++
		}

		parser.skipSpaces(true)
		for parser.data[parser.pos] == ',' {
			parser.pos++
			parser.skipSpaces(true)
		}
	}

	return nil, errors.New("unexpected end of text")
}

func (parser *dataParser) parseArray() ([]DataValue, error) {
	parser.pos++
	parser.skipSpaces(true)

	array := []DataValue{}

	for parser.pos < parser.size {
		parser.skipSpaces(true)
		for parser.data[parser.pos] == ',' && parser.pos < parser.size {
			parser.pos++
			parser.skipSpaces(true)
		}

		if parser.pos >= parser.size {
			break
		}

		if parser.data[parser.pos] == ']' {
			parser.pos++
			parser.skipSpaces(true)
			return array, nil
		}

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
			val := new(dataStringValue)
			val.value = tag
			array = append(array, val)
		}

		switch parser.data[parser.pos] {
		case ']', ',', '\n':

		default:
			return nil, fmt.Errorf(`expected ']' or ',' (line: %d, position: %d)`, parser.line, parser.pos-parser.lineStart)
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

	return nil, errors.New("unexpected end of text")
}

// ParseDataText - parse text and return DataNode
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
