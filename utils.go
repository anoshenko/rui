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
