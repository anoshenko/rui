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
	dataProperty
}

// NewBoundsProperty creates the new BoundsProperty object.
//
// The following SizeUnit properties can be used: "left" (Left), "right" (Right), "top" (Top), and "bottom" (Bottom).
func NewBoundsProperty(params Params) BoundsProperty {
	bounds := new(boundsPropertyData)
	bounds.init()

	if params != nil {
		for _, tag := range bounds.supportedProperties {
			if value, ok := params[tag]; ok && value != nil {
				bounds.set(bounds, tag, value)
			}
		}
	}
	return bounds
}

// NewBounds creates the new BoundsProperty object.
//
// The arguments specify the boundaries in a clockwise direction: "top", "right", "bottom", and "left".
//
// If the argument is specified as int or float64, the value is considered to be in pixels.
func NewBounds[topType SizeUnit | int | float64, rightType SizeUnit | int | float64, bottomType SizeUnit | int | float64, leftType SizeUnit | int | float64](
	top topType, right rightType, bottom bottomType, left leftType) BoundsProperty {
	return NewBoundsProperty(Params{
		Top:    top,
		Right:  right,
		Bottom: bottom,
		Left:   left,
	})
}

func (bounds *boundsPropertyData) init() {
	bounds.dataProperty.init()
	bounds.normalize = normalizeBoundsTag
	bounds.supportedProperties = []PropertyName{Top, Right, Bottom, Left}
}

func normalizeBoundsTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
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

func (bounds *Bounds) setFromProperties(tag, topTag, rightTag, bottomTag, leftTag PropertyName, properties Properties, session Session) {
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

// String convert Bounds to string
func (bounds *Bounds) String() string {
	if bounds.allFieldsEqual() {
		return bounds.Top.String()
	}
	return bounds.Top.String() + "," + bounds.Right.String() + "," +
		bounds.Bottom.String() + "," + bounds.Left.String()
}

func (bounds *Bounds) cssValue(tag PropertyName, builder cssBuilder, session Session) {
	if bounds.allFieldsEqual() {
		builder.add(string(tag), bounds.Top.cssString("0", session))
	} else {
		builder.addValues(string(tag), " ",
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

func setBoundsProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	if !setSimpleProperty(properties, tag, value) {
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
					for i, tag := range []PropertyName{Top, Right, Bottom, Left} {
						if !bounds.Set(tag, values[i]) {
							return nil
						}
					}
					properties.setRaw(tag, bounds)
					return []PropertyName{tag}

				default:
					notCompatibleType(tag, value)
					return nil
				}
			}
			return setSizeProperty(properties, tag, value)

		case SizeUnit:
			properties.setRaw(tag, value)

		case float32:
			properties.setRaw(tag, Px(float64(value)))

		case float64:
			properties.setRaw(tag, Px(value))

		case Bounds:
			bounds := NewBoundsProperty(nil)
			if value.Top.Type != Auto {
				bounds.setRaw(Top, value.Top)
			}
			if value.Right.Type != Auto {
				bounds.setRaw(Right, value.Right)
			}
			if value.Bottom.Type != Auto {
				bounds.setRaw(Bottom, value.Bottom)
			}
			if value.Left.Type != Auto {
				bounds.setRaw(Left, value.Left)
			}
			properties.setRaw(tag, bounds)

		case BoundsProperty:
			properties.setRaw(tag, value)

		case DataObject:
			bounds := NewBoundsProperty(nil)
			for _, tag := range []PropertyName{Top, Right, Bottom, Left} {
				if text, ok := value.PropertyValue(string(tag)); ok {
					if !bounds.Set(tag, text) {
						notCompatibleType(tag, value)
						return nil
					}
				}
			}
			properties.setRaw(tag, bounds)

		default:
			if n, ok := isInt(value); ok {
				properties.setRaw(tag, Px(float64(n)))
			} else {
				notCompatibleType(tag, value)
				return nil
			}
		}
	}

	return []PropertyName{tag}
}

func removeBoundsPropertySide(properties Properties, mainTag, sideTag PropertyName) []PropertyName {
	if bounds := getBoundsProperty(properties, mainTag); bounds != nil {
		if bounds.getRaw(sideTag) != nil {
			bounds.Remove(sideTag)
			if bounds.empty() {
				bounds = nil
			}
			properties.setRaw(mainTag, bounds)
			return []PropertyName{mainTag, sideTag}
		}
	}
	return []PropertyName{}
}

func setBoundsPropertySide(properties Properties, mainTag, sideTag PropertyName, value any) []PropertyName {
	if value == nil {
		return removeBoundsPropertySide(properties, mainTag, sideTag)
	}

	bounds := getBoundsProperty(properties, mainTag)
	if bounds == nil {
		bounds = NewBoundsProperty(nil)
	}
	if bounds.Set(sideTag, value) {
		properties.setRaw(mainTag, bounds)
		return []PropertyName{mainTag, sideTag}
	}

	notCompatibleType(sideTag, value)
	return nil
}

func getBoundsProperty(properties Properties, tag PropertyName) BoundsProperty {
	if value := properties.getRaw(tag); value != nil {
		switch value := value.(type) {
		case string:
			bounds := NewBoundsProperty(nil)
			for _, t := range []PropertyName{Top, Right, Bottom, Left} {
				bounds.Set(t, value)
			}
			return bounds

		case SizeUnit:
			bounds := NewBoundsProperty(nil)
			for _, t := range []PropertyName{Top, Right, Bottom, Left} {
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

	return nil
}

func getBounds(properties Properties, tag PropertyName, session Session) (Bounds, bool) {
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
