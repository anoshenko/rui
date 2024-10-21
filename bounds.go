package rui

import (
	"fmt"
	"strings"
)

// BorderProperty is an interface of a bounds property data
type BoundsProperty interface {
	Properties
	fmt.Stringer
	stringWriter

	// Bounds returns top, right, bottom and left size of the bounds
	Bounds(session Session) Bounds
}

type boundsPropertyData struct {
	propertyList
}

// NewBoundsProperty creates the new BoundsProperty object.
// The following SizeUnit properties can be used: "left" (Left), "right" (Right), "top" (Top), and "bottom" (Bottom).
func NewBoundsProperty(params Params) BoundsProperty {
	bounds := new(boundsPropertyData)
	bounds.properties = map[string]any{}
	if params != nil {
		for _, tag := range []string{Top, Right, Bottom, Left} {
			if value, ok := params[tag]; ok {
				bounds.Set(tag, value)
			}
		}
	}
	return bounds
}

func (bounds *boundsPropertyData) normalizeTag(tag string) string {
	tag = strings.ToLower(tag)
	switch tag {
	case MarginTop, PaddingTop, CellPaddingTop,
		"top-margin", "top-padding", "top-cell-padding":
		tag = Top

	case MarginRight, PaddingRight, CellPaddingRight,
		"right-margin", "right-padding", "right-cell-padding":
		tag = Right

	case MarginBottom, PaddingBottom, CellPaddingBottom,
		"bottom-margin", "bottom-padding", "bottom-cell-padding":
		tag = Bottom

	case MarginLeft, PaddingLeft, CellPaddingLeft,
		"left-margin", "left-padding", "left-cell-padding":
		tag = Left
	}

	return tag
}

func (bounds *boundsPropertyData) String() string {
	return runStringWriter(bounds)
}

func (bounds *boundsPropertyData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("_{ ")
	comma := false
	for _, tag := range []string{Top, Right, Bottom, Left} {
		if value, ok := bounds.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(tag)
			buffer.WriteString(" = ")
			writePropertyValue(buffer, tag, value, indent)
			comma = true
		}
	}
	buffer.WriteString(" }")
}

func (bounds *boundsPropertyData) Remove(tag string) {
	bounds.propertyList.Remove(bounds.normalizeTag(tag))
}

func (bounds *boundsPropertyData) Set(tag string, value any) bool {
	if value == nil {
		bounds.Remove(tag)
		return true
	}

	tag = bounds.normalizeTag(tag)

	switch tag {
	case Top, Right, Bottom, Left:
		return bounds.setSizeProperty(tag, value)

	default:
		ErrorLogF(`"%s" property is not compatible with the BoundsProperty`, tag)
	}

	return false
}

func (bounds *boundsPropertyData) Get(tag string) any {
	tag = bounds.normalizeTag(tag)
	if value, ok := bounds.properties[tag]; ok {
		return value
	}

	return nil
}

func (bounds *boundsPropertyData) Bounds(session Session) Bounds {
	top, _ := sizeProperty(bounds, Top, session)
	right, _ := sizeProperty(bounds, Right, session)
	bottom, _ := sizeProperty(bounds, Bottom, session)
	left, _ := sizeProperty(bounds, Left, session)
	return Bounds{Top: top, Right: right, Bottom: bottom, Left: left}
}

// Bounds describe bounds of rectangle.
type Bounds struct {
	Top, Right, Bottom, Left SizeUnit
}

// DefaultBounds return bounds with Top, Right, Bottom and Left fields set to Auto
func DefaultBounds() Bounds {
	return Bounds{
		Top:    SizeUnit{Type: Auto, Value: 0},
		Right:  SizeUnit{Type: Auto, Value: 0},
		Bottom: SizeUnit{Type: Auto, Value: 0},
		Left:   SizeUnit{Type: Auto, Value: 0},
	}
}

// SetAll set the Top, Right, Bottom and Left field to the equal value
func (bounds *Bounds) SetAll(value SizeUnit) {
	bounds.Top = value
	bounds.Right = value
	bounds.Bottom = value
	bounds.Left = value
}

func (bounds *Bounds) setFromProperties(tag, topTag, rightTag, bottomTag, leftTag string, properties Properties, session Session) {
	bounds.Top = AutoSize()
	if size, ok := sizeProperty(properties, tag, session); ok {
		bounds.Top = size
	}
	bounds.Right = bounds.Top
	bounds.Bottom = bounds.Top
	bounds.Left = bounds.Top

	if size, ok := sizeProperty(properties, topTag, session); ok {
		bounds.Top = size
	}
	if size, ok := sizeProperty(properties, rightTag, session); ok {
		bounds.Right = size
	}
	if size, ok := sizeProperty(properties, bottomTag, session); ok {
		bounds.Bottom = size
	}
	if size, ok := sizeProperty(properties, leftTag, session); ok {
		bounds.Left = size
	}
}

/*
func (bounds *Bounds) allFieldsAuto() bool {
	return bounds.Left.Type == Auto &&
		bounds.Top.Type == Auto &&
		bounds.Right.Type == Auto &&
		bounds.Bottom.Type == Auto
}

func (bounds *Bounds) allFieldsZero() bool {
	return (bounds.Left.Type == Auto || bounds.Left.Value == 0) &&
		(bounds.Top.Type == Auto || bounds.Top.Value == 0) &&
		(bounds.Right.Type == Auto || bounds.Right.Value == 0) &&
		(bounds.Bottom.Type == Auto || bounds.Bottom.Value == 0)
}
*/

func (bounds *Bounds) allFieldsEqual() bool {
	if bounds.Left.Type == bounds.Top.Type &&
		bounds.Left.Type == bounds.Right.Type &&
		bounds.Left.Type == bounds.Bottom.Type {
		return bounds.Left.Type == Auto ||
			(bounds.Left.Value == bounds.Top.Value &&
				bounds.Left.Value == bounds.Right.Value &&
				bounds.Left.Value == bounds.Bottom.Value)
	}

	return false
}

/*
func (bounds Bounds) writeCSSString(buffer *strings.Builder, textForAuto string) {
	buffer.WriteString(bounds.Top.cssString(textForAuto))
	if !bounds.allFieldsEqual() {
		buffer.WriteRune(' ')
		buffer.WriteString(bounds.Right.cssString(textForAuto))
		buffer.WriteRune(' ')
		buffer.WriteString(bounds.Bottom.cssString(textForAuto))
		buffer.WriteRune(' ')
		buffer.WriteString(bounds.Left.cssString(textForAuto))
	}
}
*/

// String convert Bounds to string
func (bounds *Bounds) String() string {
	if bounds.allFieldsEqual() {
		return bounds.Top.String()
	}
	return bounds.Top.String() + "," + bounds.Right.String() + "," +
		bounds.Bottom.String() + "," + bounds.Left.String()
}

func (bounds *Bounds) cssValue(tag string, builder cssBuilder, session Session) {
	if bounds.allFieldsEqual() {
		builder.add(tag, bounds.Top.cssString("0", session))
	} else {
		builder.addValues(tag, " ",
			bounds.Top.cssString("0", session),
			bounds.Right.cssString("0", session),
			bounds.Bottom.cssString("0", session),
			bounds.Left.cssString("0", session))
	}
}

func (bounds *Bounds) cssString(session Session) string {
	var builder cssValueBuilder
	bounds.cssValue("", &builder, session)
	return builder.finish()
}

func (properties *propertyList) setBounds(tag string, value any) bool {
	if !properties.setSimpleProperty(tag, value) {
		switch value := value.(type) {
		case string:
			if strings.Contains(value, ",") {
				values := split4Values(value)
				count := len(values)
				switch count {
				case 1:
					value = values[0]

				case 4:
					bounds := NewBoundsProperty(nil)
					for i, tag := range []string{Top, Right, Bottom, Left} {
						if !bounds.Set(tag, values[i]) {
							notCompatibleType(tag, value)
							return false
						}
					}
					properties.properties[tag] = bounds
					return true

				default:
					notCompatibleType(tag, value)
					return false
				}
			}
			return properties.setSizeProperty(tag, value)

		case SizeUnit:
			properties.properties[tag] = value

		case float32:
			properties.properties[tag] = Px(float64(value))

		case float64:
			properties.properties[tag] = Px(value)

		case Bounds:
			bounds := NewBoundsProperty(nil)
			if value.Top.Type != Auto {
				bounds.Set(Top, value.Top)
			}
			if value.Right.Type != Auto {
				bounds.Set(Right, value.Right)
			}
			if value.Bottom.Type != Auto {
				bounds.Set(Bottom, value.Bottom)
			}
			if value.Left.Type != Auto {
				bounds.Set(Left, value.Left)
			}
			properties.properties[tag] = bounds

		case BoundsProperty:
			properties.properties[tag] = value

		case DataObject:
			bounds := NewBoundsProperty(nil)
			for _, tag := range []string{Top, Right, Bottom, Left} {
				if text, ok := value.PropertyValue(tag); ok {
					if !bounds.Set(tag, text) {
						notCompatibleType(tag, value)
						return false
					}
				}
			}
			properties.properties[tag] = bounds

		default:
			if n, ok := isInt(value); ok {
				properties.properties[tag] = Px(float64(n))
			} else {
				notCompatibleType(tag, value)
				return false
			}
		}
	}

	return true
}

func (properties *propertyList) boundsProperty(tag string) BoundsProperty {
	if value, ok := properties.properties[tag]; ok {
		switch value := value.(type) {
		case string:
			bounds := NewBoundsProperty(nil)
			for _, t := range []string{Top, Right, Bottom, Left} {
				bounds.Set(t, value)
			}
			return bounds

		case SizeUnit:
			bounds := NewBoundsProperty(nil)
			for _, t := range []string{Top, Right, Bottom, Left} {
				bounds.Set(t, value)
			}
			return bounds

		case BoundsProperty:
			return value

		case Bounds:
			return NewBoundsProperty(Params{
				Top:    value.Top,
				Right:  value.Right,
				Bottom: value.Bottom,
				Left:   value.Left})
		}
	}

	return NewBoundsProperty(nil)
}

func (properties *propertyList) removeBoundsSide(mainTag, sideTag string) {
	bounds := properties.boundsProperty(mainTag)
	if bounds.Get(sideTag) != nil {
		bounds.Remove(sideTag)
		properties.properties[mainTag] = bounds
	}
}

func (properties *propertyList) setBoundsSide(mainTag, sideTag string, value any) bool {
	bounds := properties.boundsProperty(mainTag)
	if bounds.Set(sideTag, value) {
		properties.properties[mainTag] = bounds
		return true
	}

	notCompatibleType(sideTag, value)
	return false
}

func boundsProperty(properties Properties, tag string, session Session) (Bounds, bool) {
	if value := properties.Get(tag); value != nil {
		switch value := value.(type) {
		case string:
			if text, ok := session.resolveConstants(value); ok {
				if size, ok := StringToSizeUnit(text); ok {
					return Bounds{Left: size, Top: size, Right: size, Bottom: size}, true
				}
			}

		case SizeUnit:
			return Bounds{Left: value, Top: value, Right: value, Bottom: value}, true

		case Bounds:
			return value, true

		case BoundsProperty:
			return value.Bounds(session), true

		default:
			notCompatibleType(tag, value)
		}
	}

	return DefaultBounds(), false
}
