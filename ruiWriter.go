package rui

/*
import (
	"fmt"
	"strconv"
	"strings"
)

type ruiWriter interface {
	startObject(tag string)
	startObjectProperty(tag, objectTag string)
	endObject()
	startArrayProperty(tag string)
	endObArray()
	writeProperty(tag string, value any)
	finish() string
}

type ruiStringer interface {
	ruiString(writer ruiWriter)
}

type ruiWriterData struct {
	buffer *strings.Builder
	indent string
}

func newRUIWriter() ruiWriter {
	writer := new(ruiWriterData)
	return writer
}

func (writer *ruiWriterData) writeIndent() {
	if writer.buffer == nil {
		writer.buffer = allocStringBuilder()
		writer.indent = ""
		return
	}

	if writer.indent != "" {
		writer.buffer.WriteString(writer.indent)
	}
}

func (writer *ruiWriterData) writeString(str string) {
	hasEsc := strings.ContainsAny(str, "\t\"\r\n")
	if hasEsc || strings.ContainsAny(str, " ,;'`[]{}()") {
		if !strings.Contains(str, "`") && (hasEsc || strings.Contains(str, `\`)) {
			writer.buffer.WriteRune('`')
			writer.buffer.WriteString(str)
			writer.buffer.WriteRune('`')
		} else {
			replace := []struct{ old, new string }{
				{old: `\`, new: `\\`},
				{old: "\t", new: `\t`},
				{old: "\r", new: `\r`},
				{old: "\n", new: `\n`},
				{old: "\"", new: `\"`},
			}
			for _, s := range replace {
				str = strings.Replace(str, s.old, s.new, -1)
			}
			writer.buffer.WriteRune('"')
			writer.buffer.WriteString(str)
			writer.buffer.WriteRune('"')
		}
	} else {
		writer.buffer.WriteString(str)
	}
}

func (writer *ruiWriterData) startObject(tag string) {
	writer.writeIndent()
	writer.indent += "\t"
	writer.writeString(tag)
	writer.buffer.WriteString(" {\n")
}

func (writer *ruiWriterData) startObjectProperty(tag, objectTag string) {
	writer.writeIndent()
	writer.indent += "\t"
	writer.writeString(tag)
	writer.writeString(" = ")
	if objectTag == "" {
		objectTag = "_"
	}
	writer.writeString(objectTag)
	writer.buffer.WriteString(" {\n")
}

func (writer *ruiWriterData) endObject() {
	if len(writer.indent) > 0 {
		writer.indent = writer.indent[1:]
	}
	writer.writeIndent()
	writer.buffer.WriteRune('}')
}

func (writer *ruiWriterData) startArrayProperty(tag string) {
	writer.writeIndent()
	writer.writeString(tag)
	writer.buffer.WriteString(" = [\n")
	writer.indent += "\t"
}

func (writer *ruiWriterData) endObArray() {
	if len(writer.indent) > 0 {
		writer.indent = writer.indent[1:]
	}
	writer.writeIndent()
	writer.buffer.WriteString("],\n")
}

func (writer *ruiWriterData) writeValue(value any) {

	switch value := value.(type) {
	case string:
		writer.writeString(value)

	case ruiStringer:
		value.ruiString(writer)
		// TODO

	case fmt.Stringer:
		writer.writeString(value.String())

	case float32:
		writer.writeString(fmt.Sprintf("%g", float64(value)))

	case float64:
		writer.writeString(fmt.Sprintf("%g", value))

	case []string:
		switch len(value) {
		case 0:
			writer.buffer.WriteString("[]\n")

		case 1:
			writer.writeString(value[0])

		default:
			writer.buffer.WriteString("[\n")
			writer.indent += "\t"
			for _, v := range value {
				writer.buffer.WriteString(writer.indent)
				writer.writeString(v)
				writer.buffer.WriteString(",\n")
			}

			writer.indent = writer.indent[1:]
			writer.buffer.WriteString(writer.indent)
			writer.buffer.WriteRune(']')
		}

	case []View:
		switch len(value) {
		case 0:
			writer.buffer.WriteString("[]\n")

		case 1:
			writer.writeValue(value[0])

		default:
			writer.buffer.WriteString("[\n")
			writer.indent += "\t"
			for _, v := range value {
				writer.buffer.WriteString(writer.indent)
				v.ruiString(writer)
				writer.buffer.WriteString(",\n")
			}

			writer.indent = writer.indent[1:]
			writer.buffer.WriteString(writer.indent)
			writer.buffer.WriteRune(']')
		}

	case []any:
		switch len(value) {
		case 0:
			writer.buffer.WriteString("[]\n")

		case 1:
			writer.writeValue(value[0])

		default:
			writer.buffer.WriteString("[\n")
			writer.indent += "\t"
			for _, v := range value {
				writer.buffer.WriteString(writer.indent)
				writer.writeValue(v)
				writer.buffer.WriteString(",\n")
			}

			writer.indent = writer.indent[1:]
			writer.buffer.WriteString(writer.indent)
			writer.buffer.WriteRune(']')
		}

	default:
		if n, ok := isInt(value); ok {
			writer.buffer.WriteString(strconv.Itoa(n))
		}
	}
	writer.buffer.WriteString(",\n")
}

func (writer *ruiWriterData) writeProperty(tag string, value any) {
	writer.writeIndent()
	writer.writeString(tag)
	writer.buffer.WriteString(" = ")
	writer.writeValue(value)
}

func (writer *ruiWriterData) finish() string {
	result := ""
	if writer.buffer != nil {
		result = writer.buffer.String()
		freeStringBuilder(writer.buffer)
		writer.buffer = nil
	}
	return result
}
*/
