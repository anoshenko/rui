package rui

import (
	"encoding/base64"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const stringBuilderCap = 4096

var stringBuilderPool = sync.Pool{
	New: func() any {
		result := new(strings.Builder)
		result.Grow(stringBuilderCap)
		return result
	},
}

func allocStringBuilder() *strings.Builder {
	if builder := stringBuilderPool.Get(); builder != nil {
		return builder.(*strings.Builder)
	}

	result := new(strings.Builder)
	result.Grow(stringBuilderCap)
	return result
}

func freeStringBuilder(builder *strings.Builder) {
	if builder != nil && builder.Cap() == stringBuilderCap {
		builder.Reset()
		stringBuilderPool.Put(builder)
	}
}

// GetLocalIP return IP address of the machine interface
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

func dataIntProperty(data DataObject, tag string) (int, bool) {
	if value, ok := data.PropertyText(tag); ok {
		if n, err := strconv.Atoi(value); err == nil {
			return n, true
		}
	}
	return 0, false
}

func dataBoolProperty(data DataObject, tag string) bool {
	if value, ok := data.PropertyText(tag); ok && value == "1" {
		return true
	}
	return false
}

func dataFloatProperty(data DataObject, tag string) float64 {
	if value, ok := data.PropertyText(tag); ok {
		if n, err := strconv.ParseFloat(value, 64); err == nil {
			return n
		}
	}
	return 0
}

// InlineImageFromResource reads image from resources and converts it to an inline image.
// Supported png, jpeg, gif, and svg files
func InlineImageFromResource(filename string) (string, bool) {
	if image, ok := resources.images[filename]; ok && image.fs != nil {
		dataType := map[string]string{
			".svg":  "data:image/svg+xml",
			".png":  "data:image/png",
			".jpg":  "data:image/jpg",
			".jpeg": "data:image/jpg",
			".gif":  "data:image/gif",
		}
		ext := strings.ToLower(filepath.Ext(filename))
		if prefix, ok := dataType[ext]; ok {
			if data, err := image.fs.ReadFile(image.path); err == nil {
				return prefix + ";base64," + base64.StdEncoding.EncodeToString(data), true
			} else {
				DebugLog(err.Error())
			}
		} else {
			DebugLogF(`InlineImageFromResource("%s") error: Unsupported file`, filename)
		}
	} else {
		DebugLogF(`The resource image "%s" not found`, filename)
	}

	return "", false
}

// InlineFileFromData creates inline image from the image data and mimetype.
func InlineFileFromData(data []byte, mimeType string) string {
	if data == nil {
		return ""
	}

	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data)
}

func textToHtml(text string) string {
	replace := []struct{ old, new string }{
		{old: "&", new: "&amp;"},
		{old: "\"", new: "&quot;"},
		{old: "'", new: "&apos;"},
		//{old: "<", new: "&lt;"},
		//{old: ">", new: "&gt;"},
	}

	for _, r := range replace {
		text = strings.ReplaceAll(text, r.old, r.new)
	}

	return text
}

func writeStrings(buffer *strings.Builder, strs ...string) {
	for _, str := range strs {
		buffer.WriteString(str)
	}
}

func isQuotesNeeded(text string) bool {
	switch len(text) {
	case 0:
		return true

	case 1:
		simple := (text[0] >= '0' && text[0] <= '9') || (text[0] >= 'A' && text[0] <= 'Z') ||
			(text[0] >= 'a' && text[0] <= 'z') || text[0] == '_'
		return !simple
	}

	for _, ch := range text {
		if (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
			ch == '+' || ch == '-' || ch == '@' || ch == '/' || ch == '_' || ch == '.' ||
			ch == ':' || ch == '#' || ch == '%' || ch == 'π' || ch == '°' {
		} else {
			return true
		}
	}

	return false
}

func replaceEscapeSymbols(text string) string {
	replace := []struct{ old, new string }{
		{old: "\\", new: `\\`},
		{old: "\t", new: `\t`},
		{old: "\r", new: `\r`},
		{old: "\n", new: `\n`},
		{old: "\"", new: `\"`},
	}
	for _, s := range replace {
		text = strings.ReplaceAll(text, s.old, s.new)
	}
	return text
}

func writeRuiText(buffer *strings.Builder, text string, objectName bool) {
	if isQuotesNeeded(text) || (objectName && strings.Contains(text, "/")) {
		if strings.ContainsAny(text, "\\\t\r\n\"") {
			text = replaceEscapeSymbols(text)
		}

		buffer.WriteRune('"')
		buffer.WriteString(text)
		buffer.WriteRune('"')
	} else {
		buffer.WriteString(text)
	}
}
