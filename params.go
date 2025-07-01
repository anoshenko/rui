package rui

import (
	"slices"
)

// Params defines a type of a parameters list
type Params map[PropertyName]any

// Get returns a value of the property with name defined by the argument. The type of return value depends
// on the property. If the property is not set then nil is returned.
func (params Params) Get(tag PropertyName) any {
	return params.getRaw(tag)
}

func (params Params) getRaw(tag PropertyName) any {
	if value, ok := params[tag]; ok {
		return value
	}
	return nil
}

// Set sets the value (second argument) of the property with name defined by the first argument.
// Return "true" if the value has been set, in the opposite case "false" is returned and a description of an error is written to the log
func (params Params) Set(tag PropertyName, value any) bool {
	params.setRaw(tag, value)
	return true
}

func (params Params) setRaw(tag PropertyName, value any) {
	if value != nil {
		params[tag] = value
	} else {
		delete(params, tag)
	}
}

// Remove removes the property with name defined by the argument from a map.
func (params Params) Remove(tag PropertyName) {
	delete(params, tag)
}

// Clear removes all properties from a map.
func (params Params) Clear() {
	for tag := range params {
		delete(params, tag)
	}
}

// AllTags returns a sorted slice of all properties.
func (params Params) AllTags() []PropertyName {
	tags := make([]PropertyName, 0, len(params))
	for t := range params {
		tags = append(tags, t)
	}
	slices.Sort(tags)
	return tags
}

func (params Params) empty() bool {
	return len(params) == 0
}
