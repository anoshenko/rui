package rui

import (
	"encoding/base64"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

var stringBuilders []*strings.Builder = make([]*strings.Builder, 4096)
var stringBuilderCount = 0

func allocStringBuilder() *strings.Builder {
	for stringBuilderCount > 0 {
		stringBuilderCount--
		result := stringBuilders[stringBuilderCount]
		if result != nil {
			stringBuilders[stringBuilderCount] = nil
			result.Reset()
			return result
		}
	}

	result := new(strings.Builder)
	result.Grow(4096)
	return result
}

func freeStringBuilder(builder *strings.Builder) {
	if builder != nil {
		if stringBuilderCount == len(stringBuilders) {
			stringBuilders = append(stringBuilders, builder)
		} else {
			stringBuilders[stringBuilderCount] = builder
		}
		stringBuilderCount++
	}
}

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
	if value, ok := data.PropertyValue(tag); ok {
		if n, err := strconv.Atoi(value); err == nil {
			return n, true
		}
	}
	return 0, false
}

func dataBoolProperty(data DataObject, tag string) bool {
	if value, ok := data.PropertyValue(tag); ok && value == "1" {
		return true
	}
	return false
}

func dataFloatProperty(data DataObject, tag string) float64 {
	if value, ok := data.PropertyValue(tag); ok {
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
