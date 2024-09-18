package rui

import (
	"fmt"
	"math"
	"strings"
)

// Constants for [Transform] specific properties
const (
	// Perspective is the constant for "perspective" property tag.
	//
	// Used by `View`.
	// Distance between the z-plane and the user in order to give a 3D-positioned element some perspective. Each 3D element 
	// with z > 0 becomes larger, each 3D-element with z < 0 becomes smaller. The default value is 0 (no 3D effects).
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	Perspective = "perspective"

	// PerspectiveOriginX is the constant for "perspective-origin-x" property tag.
	//
	// Used by `View`.
	// x-coordinate of the position at which the viewer is looking. It is used as the vanishing point by the "perspective" 
	// property. The default value is 50%.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	PerspectiveOriginX = "perspective-origin-x"

	// PerspectiveOriginY is the constant for "perspective-origin-y" property tag.
	//
	// Used by `View`.
	// y-coordinate of the position at which the viewer is looking. It is used as the vanishing point by the "perspective" 
	// property. The default value is 50%.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	PerspectiveOriginY = "perspective-origin-y"

	// BackfaceVisible is the constant for "backface-visibility" property tag.
	//
	// Used by `View`.
	// Controls whether the back face of a view is visible when turned towards the user. Default value is `true`.
	//
	// Supported types: `bool`, `int`, `string`.
	//
	// Values:
	// `true` or `1` or "true", "yes", "on", "1" - Back face is visible when turned towards the user.
	// `false` or `0` or "false", "no", "off", "0" - Back face is hidden, effectively making the view invisible when turned away from the user.
	BackfaceVisible = "backface-visibility"

	// OriginX is the constant for "origin-x" property tag.
	//
	// Used by `View`.
	// x-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	OriginX = "origin-x"

	// OriginY is the constant for "origin-y" property tag.
	//
	// Used by `View`.
	// y-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	OriginY = "origin-y"

	// OriginZ is the constant for "origin-z" property tag.
	//
	// Used by `View`.
	// z-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	OriginZ = "origin-z"

	// TransformTag is the constant for "transform" property tag.
	//
	// Used by `View`.
	// Specify translation, scale and rotation over x, y and z axes as well as a distorsion of a view along x and y axes.
	//
	// Supported types: `Transform`, `string`.
	//
	// See `Transform` description for more details.
	//
	// Conversion rules:
	// `Transform` - stored as is, no conversion performed.
	// `string` - string representation of `Transform` interface. Example: "_{translate-x = 10px, scale-y = 1.1}".
	TransformTag = "transform"

	// TranslateX is the constant for "translate-x" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// x-axis translation value of a 2D/3D translation.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	//
	// Usage in `Transform`:
	// x-axis translation value of a 2D/3D translation.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	TranslateX = "translate-x"

	// TranslateY is the constant for "translate-y" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// y-axis translation value of a 2D/3D translation.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	//
	// Usage in `Transform`:
	// x-axis translation value of a 2D/3D translation.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	TranslateY = "translate-y"

	// TranslateZ is the constant for "translate-z" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// z-axis translation value of a 3D translation.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	//
	// Usage in `Transform`:
	// z-axis translation value of a 3D translation.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	TranslateZ = "translate-z"

	// ScaleX is the constant for "scale-x" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// x-axis scaling value of a 2D/3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original 
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	//
	// Usage in `Transform`:
	// x-axis scaling value of a 2D/3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original 
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	ScaleX = "scale-x"

	// ScaleY is the constant for "scale-y" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// y-axis scaling value of a 2D/3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original 
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	//
	// Usage in `Transform`:
	// y-axis scaling value of a 2D/3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original 
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	ScaleY = "scale-y"

	// ScaleZ is the constant for "scale-z" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// z-axis scaling value of a 3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original 
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	//
	// Usage in `Transform`:
	// z-axis scaling value of a 3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original 
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	ScaleZ = "scale-z"

	// Rotate is the constant for "rotate" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// Angle of the view rotation. A positive angle denotes a clockwise rotation, a negative angle a counter-clockwise.
	//
	// Supported types: `AngleUnit`, `string`, `float`, `int`.
	//
	// Internal type is `AngleUnit`, other types will be converted to it during assignment.
	// See `AngleUnit` description for more details.
	//
	// Conversion rules:
	// `AngleUnit` - stored as is, no conversion performed.
	// `string` - must contain string representation of `AngleUnit`. If numeric value will be provided without any suffix then `AngleUnit` with value and `Radian` value type will be created.
	// `float` - a new `AngleUnit` value will be created with `Radian` as a type.
	// `int` - a new `AngleUnit` value will be created with `Radian` as a type.
	//
	// Usage in `Transform`:
	// Angle of the view rotation. A positive angle denotes a clockwise rotation, a negative angle a counter-clockwise.
	//
	// Supported types: `AngleUnit`, `string`, `float`, `int`.
	//
	// Internal type is `AngleUnit`, other types will be converted to it during assignment.
	// See `AngleUnit` description for more details.
	//
	// Conversion rules:
	// `AngleUnit` - stored as is, no conversion performed.
	// `string` - must contain string representation of `AngleUnit`. If numeric value will be provided without any suffix then `AngleUnit` with value and `Radian` value type will be created.
	// `float` - a new `AngleUnit` value will be created with `Radian` as a type.
	// `int` - a new `AngleUnit` value will be created with `Radian` as a type.
	Rotate = "rotate"

	// RotateX is the constant for "rotate-x" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// x-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	//
	// Usage in `Transform`:
	// x-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	RotateX = "rotate-x"

	// RotateY is the constant for "rotate-y" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// y-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	//
	// Usage in `Transform`:
	// y-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	RotateY = "rotate-y"

	// RotateZ is the constant for "rotate-z" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// z-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	//
	// Usage in `Transform`:
	// z-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: `float`, `int`, `string`.
	//
	// Internal type is `float`, other types converted to it during assignment.
	RotateZ = "rotate-z"

	// SkewX is the constant for "skew-x" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// Angle to use to distort the element along the abscissa. The default value is 0.
	//
	// Supported types: `AngleUnit`, `string`, `float`, `int`.
	//
	// Internal type is `AngleUnit`, other types will be converted to it during assignment.
	// See `AngleUnit` description for more details.
	//
	// Conversion rules:
	// `AngleUnit` - stored as is, no conversion performed.
	// `string` - must contain string representation of `AngleUnit`. If numeric value will be provided without any suffix then `AngleUnit` with value and `Radian` value type will be created.
	// `float` - a new `AngleUnit` value will be created with `Radian` as a type.
	// `int` - a new `AngleUnit` value will be created with `Radian` as a type.
	//
	// Usage in `Transform`:
	// Angle to use to distort the element along the abscissa. The default value is 0.
	//
	// Supported types: `AngleUnit`, `string`, `float`, `int`.
	//
	// Internal type is `AngleUnit`, other types will be converted to it during assignment.
	// See `AngleUnit` description for more details.
	//
	// Conversion rules:
	// `AngleUnit` - stored as is, no conversion performed.
	// `string` - must contain string representation of `AngleUnit`. If numeric value will be provided without any suffix then `AngleUnit` with value and `Radian` value type will be created.
	// `float` - a new `AngleUnit` value will be created with `Radian` as a type.
	// `int` - a new `AngleUnit` value will be created with `Radian` as a type.
	SkewX = "skew-x"

	// SkewY is the constant for "skew-y" property tag.
	//
	// Used by `View`, `Transform`.
	//
	// Usage in `View`:
	// Angle to use to distort the element along the ordinate. The default value is 0.
	//
	// Supported types: `AngleUnit`, `string`, `float`, `int`.
	//
	// Internal type is `AngleUnit`, other types will be converted to it during assignment.
	// See `AngleUnit` description for more details.
	//
	// Conversion rules:
	// `AngleUnit` - stored as is, no conversion performed.
	// `string` - must contain string representation of `AngleUnit`. If numeric value will be provided without any suffix then `AngleUnit` with value and `Radian` value type will be created.
	// `float` - a new `AngleUnit` value will be created with `Radian` as a type.
	// `int` - a new `AngleUnit` value will be created with `Radian` as a type.
	//
	// Usage in `Transform`:
	// Angle to use to distort the element along the ordinate. The default value is 0.
	//
	// Supported types: `AngleUnit`, `string`, `float`, `int`.
	//
	// Internal type is `AngleUnit`, other types will be converted to it during assignment.
	// See `AngleUnit` description for more details.
	//
	// Conversion rules:
	// `AngleUnit` - stored as is, no conversion performed.
	// `string` - must contain string representation of `AngleUnit`. If numeric value will be provided without any suffix then `AngleUnit` with value and `Radian` value type will be created.
	// `float` - a new `AngleUnit` value will be created with `Radian` as a type.
	// `int` - a new `AngleUnit` value will be created with `Radian` as a type.
	SkewY = "skew-y"
)

// Transform interface specifies view transformation parameters: the x-, y-, and z-axis translation values,
// the x-, y-, and z-axis scaling values, the angle to use to distort the element along the abscissa and ordinate,
// the angle of the view rotation.
// Valid property tags: TranslateX ("translate-x"), TranslateY ("translate-y"), TranslateZ ("translate-z"),
// ScaleX ("scale-x"), ScaleY ("scale-y"), ScaleZ ("scale-z"), Rotate ("rotate"), RotateX ("rotate-x"),
// RotateY ("rotate-y"), RotateZ ("rotate-z"), SkewX ("skew-x"), and SkewY ("skew-y")
type Transform interface {
	Properties
	fmt.Stringer
	stringWriter
	transformCSS(session Session, transform3D bool) string
}

type transformData struct {
	propertyList
}

// NewTransform creates a new transform property data and return its interface
func NewTransform(params Params) Transform {
	transform := new(transformData)
	transform.properties = map[string]any{}
	for tag, value := range params {
		transform.Set(tag, value)
	}
	return transform
}

func (style *viewStyle) setTransform(value any) bool {

	setObject := func(obj DataObject) bool {
		transform := NewTransform(nil)
		ok := true
		for i := 0; i < obj.PropertyCount(); i++ {
			if prop := obj.Property(i); prop.Type() == TextNode {
				if !transform.Set(prop.Tag(), prop.Text()) {
					ok = false
				}
			} else {
				ok = false
			}
		}

		if !ok && len(transform.AllTags()) == 0 {
			return false
		}

		style.properties[TransformTag] = transform
		return true
	}

	switch value := value.(type) {
	case Transform:
		style.properties[TransformTag] = value
		return true

	case DataObject:
		return setObject(value)

	case DataNode:
		if obj := value.Object(); obj != nil {
			return setObject(obj)
		}
		notCompatibleType(TransformTag, value)
		return false

	case string:
		if obj := ParseDataText(value); obj != nil {
			return setObject(obj)
		}
		notCompatibleType(TransformTag, value)
		return false
	}

	return false
}

func (style *viewStyle) transformProperty() Transform {
	if val, ok := style.properties[TransformTag]; ok {
		if transform, ok := val.(Transform); ok {
			return transform
		}
	}
	return nil
}

func (style *viewStyle) setTransformProperty(tag string, value any) bool {
	switch tag {
	case RotateX, RotateY, RotateZ, Rotate, SkewX, SkewY, ScaleX, ScaleY, ScaleZ, TranslateX, TranslateY, TranslateZ:
		if transform := style.transformProperty(); transform != nil {
			return transform.Set(tag, value)
		}

		transform := NewTransform(nil)
		if !transform.Set(tag, value) {
			return false
		}

		style.properties[TransformTag] = transform
		return true
	}

	ErrorLogF(`"Transform" interface does not support the "%s" property`, tag)
	return false
}

func (transform *transformData) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)
	transform.writeString(buffer, "")
	return buffer.String()
}

func (transform *transformData) writeString(buffer *strings.Builder, indent string) {
	buffer.WriteString("_{ ")
	comma := false
	for _, tag := range []string{SkewX, SkewY, TranslateX, TranslateY, TranslateZ,
		ScaleX, ScaleY, ScaleZ, Rotate, RotateX, RotateY, RotateZ} {
		if value, ok := transform.properties[tag]; ok {
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

func (transform *transformData) Set(tag string, value any) bool {
	return transform.set(strings.ToLower(tag), value)
}

func (transform *transformData) set(tag string, value any) bool {
	if value == nil {
		_, exist := transform.properties[tag]
		if exist {
			delete(transform.properties, tag)
		}
		return exist
	}

	switch tag {

	case RotateX, RotateY, RotateZ:
		return transform.setFloatProperty(tag, value, 0, 1)

	case Rotate, SkewX, SkewY:
		return transform.setAngleProperty(tag, value)

	case ScaleX, ScaleY, ScaleZ:
		return transform.setFloatProperty(tag, value, -math.MaxFloat64, math.MaxFloat64)

	case TranslateX, TranslateY, TranslateZ:
		return transform.setSizeProperty(tag, value)
	}

	return false
}

func getTransform3D(style Properties, session Session) bool {
	perspective, ok := sizeProperty(style, Perspective, session)
	return ok && perspective.Type != Auto && perspective.Value != 0
}

func getPerspectiveOrigin(style Properties, session Session) (SizeUnit, SizeUnit) {
	x, _ := sizeProperty(style, PerspectiveOriginX, session)
	y, _ := sizeProperty(style, PerspectiveOriginY, session)
	return x, y
}

func getOrigin(style Properties, session Session) (SizeUnit, SizeUnit, SizeUnit) {
	x, _ := sizeProperty(style, OriginX, session)
	y, _ := sizeProperty(style, OriginY, session)
	z, _ := sizeProperty(style, OriginZ, session)
	return x, y, z
}

func (transform *transformData) getSkew(session Session) (AngleUnit, AngleUnit, bool) {
	skewX, okX := angleProperty(transform, SkewX, session)
	skewY, okY := angleProperty(transform, SkewY, session)
	return skewX, skewY, okX || okY
}

func (transform *transformData) getTranslate(session Session) (SizeUnit, SizeUnit, SizeUnit) {
	x, _ := sizeProperty(transform, TranslateX, session)
	y, _ := sizeProperty(transform, TranslateY, session)
	z, _ := sizeProperty(transform, TranslateZ, session)
	return x, y, z
}

func (transform *transformData) transformCSS(session Session, transform3D bool) string {

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	skewX, skewY, skewOK := transform.getSkew(session)
	if skewOK {
		buffer.WriteString(`skew(`)
		buffer.WriteString(skewX.cssString())
		buffer.WriteRune(',')
		buffer.WriteString(skewY.cssString())
		buffer.WriteRune(')')
	}

	x, y, z := transform.getTranslate(session)

	scaleX, okScaleX := floatTextProperty(transform, ScaleX, session, 1)
	scaleY, okScaleY := floatTextProperty(transform, ScaleY, session, 1)

	if transform3D {
		if x.Type != Auto || y.Type != Auto || z.Type != Auto {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`translate3d(`)
			buffer.WriteString(x.cssString("0", session))
			buffer.WriteRune(',')
			buffer.WriteString(y.cssString("0", session))
			buffer.WriteRune(',')
			buffer.WriteString(z.cssString("0", session))
			buffer.WriteRune(')')
		}

		scaleZ, okScaleZ := floatTextProperty(transform, ScaleZ, session, 1)
		if okScaleX || okScaleY || okScaleZ {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`scale3d(`)
			buffer.WriteString(scaleX)
			buffer.WriteRune(',')
			buffer.WriteString(scaleY)
			buffer.WriteRune(',')
			buffer.WriteString(scaleZ)
			buffer.WriteRune(')')
		}

		if angle, ok := angleProperty(transform, Rotate, session); ok {
			rotateX, _ := floatTextProperty(transform, RotateX, session, 1)
			rotateY, _ := floatTextProperty(transform, RotateY, session, 1)
			rotateZ, _ := floatTextProperty(transform, RotateZ, session, 1)

			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`rotate3d(`)
			buffer.WriteString(rotateX)
			buffer.WriteRune(',')
			buffer.WriteString(rotateY)
			buffer.WriteRune(',')
			buffer.WriteString(rotateZ)
			buffer.WriteRune(',')
			buffer.WriteString(angle.cssString())
			buffer.WriteRune(')')
		}

	} else {

		if x.Type != Auto || y.Type != Auto {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`translate(`)
			buffer.WriteString(x.cssString("0", session))
			buffer.WriteRune(',')
			buffer.WriteString(y.cssString("0", session))
			buffer.WriteRune(')')
		}

		if okScaleX || okScaleY {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`scale(`)
			buffer.WriteString(scaleX)
			buffer.WriteRune(',')
			buffer.WriteString(scaleY)
			buffer.WriteRune(')')
		}

		if angle, ok := angleProperty(transform, Rotate, session); ok {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`rotate(`)
			buffer.WriteString(angle.cssString())
			buffer.WriteRune(')')
		}
	}

	return buffer.String()
}

func (style *viewStyle) writeViewTransformCSS(builder cssBuilder, session Session) {
	transform3D := getTransform3D(style, session)
	if transform3D {
		if perspective, ok := sizeProperty(style, Perspective, session); ok && perspective.Type != Auto && perspective.Value != 0 {
			builder.add(`perspective`, perspective.cssString("0", session))
		}

		x, y := getPerspectiveOrigin(style, session)
		if x.Type != Auto || y.Type != Auto {
			builder.addValues(`perspective-origin`, ` `, x.cssString("50%", session), y.cssString("50%", session))
		}

		if backfaceVisible, ok := boolProperty(style, BackfaceVisible, session); ok {
			if backfaceVisible {
				builder.add(`backface-visibility`, `visible`)
			} else {
				builder.add(`backface-visibility`, `hidden`)
			}
		}

		x, y, z := getOrigin(style, session)
		if x.Type != Auto || y.Type != Auto || z.Type != Auto {
			builder.addValues(`transform-origin`, ` `, x.cssString("50%", session), y.cssString("50%", session), z.cssString("0", session))
		}
	} else {
		x, y, _ := getOrigin(style, session)
		if x.Type != Auto || y.Type != Auto {
			builder.addValues(`transform-origin`, ` `, x.cssString("50%", session), y.cssString("50%", session))
		}
	}

	if transform := style.transformProperty(); transform != nil {
		builder.add(`transform`, transform.transformCSS(session, transform3D))
	}
}

func (view *viewData) updateTransformProperty(tag string) bool {
	htmlID := view.htmlID()
	session := view.session

	switch tag {
	case Perspective:
		updateCSSStyle(htmlID, session)

	case PerspectiveOriginX, PerspectiveOriginY:
		if getTransform3D(view, session) {
			x, y := GetPerspectiveOrigin(view)
			value := ""
			if x.Type != Auto || y.Type != Auto {
				value = x.cssString("50%", session) + " " + y.cssString("50%", session)
			}
			session.updateCSSProperty(htmlID, "perspective-origin", value)
		}

	case BackfaceVisible:
		if getTransform3D(view, session) {
			if GetBackfaceVisible(view) {
				session.updateCSSProperty(htmlID, BackfaceVisible, "visible")
			} else {
				session.updateCSSProperty(htmlID, BackfaceVisible, "hidden")
			}
		}

	case OriginX, OriginY, OriginZ:
		x, y, z := getOrigin(view, session)
		value := ""
		if getTransform3D(view, session) {
			if x.Type != Auto || y.Type != Auto || z.Type != Auto {
				value = x.cssString("50%", session) + " " + y.cssString("50%", session) + " " + z.cssString("50%", session)
			}
		} else {
			if x.Type != Auto || y.Type != Auto {
				value = x.cssString("50%", session) + " " + y.cssString("50%", session)
			}
		}
		session.updateCSSProperty(htmlID, "transform-origin", value)

	case TransformTag, SkewX, SkewY, TranslateX, TranslateY, TranslateZ,
		ScaleX, ScaleY, ScaleZ, Rotate, RotateX, RotateY, RotateZ:
		if transform := view.transformProperty(); transform != nil {
			transform3D := getTransform3D(view, session)
			session.updateCSSProperty(htmlID, "transform", transform.transformCSS(session, transform3D))
		} else {
			session.updateCSSProperty(htmlID, "transform", "")
		}

	default:
		return false
	}

	return true
}
