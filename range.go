package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// Range defines range limits. The First and Last value are included in the range
type Range struct {
	First, Last int
}

// String returns a string representation of the Range struct
func (r Range) String() string {
	if r.First == r.Last {
		return fmt.Sprintf("%d", r.First)
	}
	return fmt.Sprintf("%d:%d", r.First, r.Last)
}

func (r *Range) setValue(value string) bool {
	var err error
	if strings.ContainsRune(value, ':') {
		values := strings.Split(value, ":")
		if len(values) != 2 {
			ErrorLog("Invalid range value: " + value)
			return false
		}
		if r.First, err = strconv.Atoi(strings.Trim(values[0], " \t\n\r")); err != nil {
			ErrorLog(`Invalid first range value "` + value + `" (` + err.Error() + ")")
			return false
		}
		if r.Last, err = strconv.Atoi(strings.Trim(values[1], " \t\n\r")); err != nil {
			ErrorLog(`Invalid last range value "` + value + `" (` + err.Error() + ")")
			return false
		}
		return true
	}

	if r.First, err = strconv.Atoi(value); err != nil {
		ErrorLog(`Invalid range value "` + value + `" (` + err.Error() + ")")
		return false
	}
	r.Last = r.First
	return true
}

func setRangeProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	switch value := value.(type) {
	case string:
		if setSimpleProperty(properties, tag, value) {
			return []PropertyName{tag}
		}

		var r Range
		if !r.setValue(value) {
			invalidPropertyValue(tag, value)
			return nil
		}
		properties.setRaw(tag, r)

	case Range:
		properties.setRaw(tag, value)

	default:
		if n, ok := isInt(value); ok {
			properties.setRaw(tag, Range{First: n, Last: n})
		} else {
			notCompatibleType(tag, value)
			return nil
		}
	}
	return []PropertyName{tag}
}
