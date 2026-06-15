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

	switch strings.Count(value, ":") {
	case 0:
		if r.First, err = strconv.Atoi(value); err != nil {
			ErrorLog(`Invalid range value "` + value + `" (` + err.Error() + ")")
			return false
		}
		r.Last = r.First
		return true

	case 1:
		if first, last, ok := strings.Cut(value, ":"); ok {
			if r.First, err = strconv.Atoi(strings.Trim(first, " \t\n\r")); err != nil {
				ErrorLog(`Invalid first range value "` + value + `" (` + err.Error() + ")")
				return false
			}
			if r.Last, err = strconv.Atoi(strings.Trim(last, " \t\n\r")); err != nil {
				ErrorLog(`Invalid last range value "` + value + `" (` + err.Error() + ")")
				return false
			}
			return true
		}
	}

	ErrorLog("Invalid range value: " + value)
	return false
}

func setRangeProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	if result := setSimpleProperty(properties, tag, value); result != nil {
		return result
	}

	var r Range
	switch value := value.(type) {
	case string:
		if !r.setValue(value) {
			invalidPropertyValue(tag, value)
			return nil
		}

	case Range:
		r = value

	default:
		if n, ok := isInt(value); ok {
			r.First = n
			r.Last = n
		} else {
			notCompatibleType(tag, value)
			return nil
		}
	}

	return setPropertyValue(properties, tag, r)
}
