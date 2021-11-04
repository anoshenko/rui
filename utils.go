package rui

import (
	"net"
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
