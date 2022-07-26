package rui

import "sort"

// Params defines a type of a parameters list
type Params map[string]any

func (params Params) Get(tag string) any {
	return params.getRaw(tag)
}

func (params Params) getRaw(tag string) any {
	if value, ok := params[tag]; ok {
		return value
	}
	return nil
}

func (params Params) Set(tag string, value any) bool {
	params.setRaw(tag, value)
	return true
}

func (params Params) setRaw(tag string, value any) {
	if value != nil {
		params[tag] = value
	} else {
		delete(params, tag)
	}
}

func (params Params) Remove(tag string) {
	delete(params, tag)
}

func (params Params) Clear() {
	for tag := range params {
		delete(params, tag)
	}
}

func (params Params) AllTags() []string {
	tags := make([]string, 0, len(params))
	for t := range params {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	return tags
}
