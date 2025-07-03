package rui

import (
	"slices"
	"strings"
)

// Properties interface of properties map
type Properties interface {
	// Get returns a value of the property with name defined by the argument.
	// The type of return value depends on the property. If the property is not set then nil is returned.
	Get(tag PropertyName) any
	getRaw(tag PropertyName) any

	// Set sets the value (second argument) of the property with name defined by the first argument.
	// Return "true" if the value has been set, in the opposite case "false" are returned and
	// a description of the error is written to the log
	Set(tag PropertyName, value any) bool
	setRaw(tag PropertyName, value any)

	// Remove removes the property with name defined by the argument
	Remove(tag PropertyName)

	// Clear removes all properties
	Clear()

	// AllTags returns an array of the set properties
	AllTags() []PropertyName

	empty() bool
}

type propertyList struct {
	properties map[PropertyName]any
	normalize  func(PropertyName) PropertyName
}

type dataProperty struct {
	propertyList
	supportedProperties []PropertyName
	get                 func(Properties, PropertyName) any
	set                 func(Properties, PropertyName, any) []PropertyName
	remove              func(Properties, PropertyName) []PropertyName
}

func defaultNormalize(tag PropertyName) PropertyName {
	return PropertyName(strings.ToLower(strings.Trim(string(tag), " \t")))
}

func (properties *propertyList) init() {
	properties.properties = map[PropertyName]any{}
	properties.normalize = defaultNormalize
	//properties.getFunc = properties.getRaw
	//properties.set = propertiesSet
	//properties.remove = propertiesRemove
}

func (properties *propertyList) empty() bool {
	return len(properties.properties) == 0
}

func (properties *propertyList) getRaw(tag PropertyName) any {
	if value, ok := properties.properties[tag]; ok {
		return value
	}
	return nil
}

func (properties *propertyList) setRaw(tag PropertyName, value any) {
	if value == nil {
		delete(properties.properties, tag)
	} else {
		properties.properties[tag] = value
	}
}

/*
	func (properties *propertyList) Remove(tag PropertyName) {
		properties.remove(properties, properties.normalize(tag))
	}
*/
func (properties *propertyList) Clear() {
	properties.properties = map[PropertyName]any{}
}

func (properties *propertyList) AllTags() []PropertyName {
	tags := make([]PropertyName, 0, len(properties.properties))
	for tag := range properties.properties {
		tags = append(tags, tag)
	}
	slices.Sort(tags)
	return tags
}

func (properties *propertyList) writeToBuffer(buffer *strings.Builder,
	indent string, objectTag string, tags []PropertyName) {

	buffer.WriteString(objectTag)
	buffer.WriteString(" {\n")

	indent2 := indent + "\t"

	for _, tag := range tags {
		if value, ok := properties.properties[tag]; ok {
			text := propertyValueToString(tag, value, indent2)
			if text != "" {
				buffer.WriteString(indent2)
				buffer.WriteString(string(tag))
				buffer.WriteString(" = ")
				buffer.WriteString(text)
				buffer.WriteString(",\n")
			}
		}
	}

	buffer.WriteString(indent)
	buffer.WriteString("}")
}

func parseProperties(properties Properties, object DataObject) {
	for node := range object.Properties() {
		switch node.Type() {
		case TextNode:
			properties.Set(PropertyName(node.Tag()), node.Text())

		case ObjectNode:
			properties.Set(PropertyName(node.Tag()), node.Object())

		case ArrayNode:
			properties.Set(PropertyName(node.Tag()), node.ArrayElements())
		}
	}
}

func propertiesGet(properties Properties, tag PropertyName) any {
	return properties.getRaw(tag)
}

func propertiesRemove(properties Properties, tag PropertyName) []PropertyName {
	if properties.getRaw(tag) == nil {
		return []PropertyName{}
	}
	properties.setRaw(tag, nil)
	return []PropertyName{tag}
}

func (data *dataProperty) init() {
	data.propertyList.init()
	data.get = propertiesGet
	data.set = propertiesSet
	data.remove = propertiesRemove
}

func (data *dataProperty) Get(tag PropertyName) any {
	return data.get(data, data.normalize(tag))
}

func (data *dataProperty) Remove(tag PropertyName) {
	data.remove(data, data.normalize(tag))
}

func (data *dataProperty) writeToBuffer(buffer *strings.Builder, indent string, objectName string, tags []PropertyName) {
	buffer.WriteString(objectName)
	buffer.WriteString("{ ")
	comma := false
	for _, tag := range tags {
		if value, ok := data.properties[tag]; ok {
			text := propertyValueToString(tag, value, indent)
			if text != "" {
				if comma {
					buffer.WriteString(", ")
				}
				buffer.WriteString(string(tag))
				buffer.WriteString(" = ")
				buffer.WriteString(text)
				comma = true
			}
		}
	}
	buffer.WriteString(" }")
}

func (data *dataProperty) writeString(buffer *strings.Builder, indent string) {
	data.writeToBuffer(buffer, indent, "_", data.AllTags())
}
