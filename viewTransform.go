package rui

import (
	"fmt"
	"math"
	"strings"
)

// Constants for [Transform] specific properties
const (
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
	TransformTag PropertyName = "transform"

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
	Perspective PropertyName = "perspective"

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
	PerspectiveOriginX PropertyName = "perspective-origin-x"

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
	PerspectiveOriginY PropertyName = "perspective-origin-y"

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
	BackfaceVisible PropertyName = "backface-visibility"

	// OriginX is the constant for "origin-x" property tag.
	//
	// Used by `View`.
	// x-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	TransformOriginX PropertyName = "transform-origin-x"

	// OriginY is the constant for "origin-y" property tag.
	//
	// Used by `View`.
	// y-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	TransformOriginY PropertyName = "transform-origin-y"

	// OriginZ is the constant for "origin-z" property tag.
	//
	// Used by `View`.
	// z-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: `SizeUnit`, `SizeFunc`, `string`, `float`, `int`.
	//
	// Internal type is `SizeUnit`, other types converted to it during assignment.
	// See `SizeUnit` description for more details.
	TransformOriginZ PropertyName = "transform-origin-z"

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
	TranslateX PropertyName = "translate-x"

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
	TranslateY PropertyName = "translate-y"

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
	TranslateZ PropertyName = "translate-z"

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
	ScaleX PropertyName = "scale-x"

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
	ScaleY PropertyName = "scale-y"

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
	ScaleZ PropertyName = "scale-z"

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
	Rotate PropertyName = "rotate"

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
	RotateX PropertyName = "rotate-x"

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
	RotateY PropertyName = "rotate-y"

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
	RotateZ PropertyName = "rotate-z"

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
	SkewX PropertyName = "skew-x"

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
	SkewY PropertyName = "skew-y"
)

// Transform interface specifies view transformation parameters: the x-, y-, and z-axis translation values,
// the x-, y-, and z-axis scaling values, the angle to use to distort the element along the abscissa and ordinate,
// the angle of the view rotation.
// Valid property tags: Perspective ("perspective"),  TranslateX ("translate-x"), TranslateY ("translate-y"), TranslateZ ("translate-z"),
// ScaleX ("scale-x"), ScaleY ("scale-y"), ScaleZ ("scale-z"), Rotate ("rotate"), RotateX ("rotate-x"),
// RotateY ("rotate-y"), RotateZ ("rotate-z"), SkewX ("skew-x"), and SkewY ("skew-y")
type Transform interface {
	Properties
	fmt.Stringer
	stringWriter
	transformCSS(session Session) string
}

type transformData struct {
	dataProperty
}

// NewTransform creates a new transform property data and return its interface
func NewTransform(params Params) Transform {
	transform := new(transformData)
	transform.init()

	for tag, value := range params {
		transform.Set(tag, value)
	}
	return transform
}

func (transform *transformData) init() {
	transform.dataProperty.init()
	transform.set = transformSet
	transform.supportedProperties = []PropertyName{
		RotateX, RotateY, RotateZ, Rotate, SkewX, SkewY, ScaleX, ScaleY, ScaleZ,
		Perspective, TranslateX, TranslateY, TranslateZ,
	}
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
	for _, tag := range transform.supportedProperties {
		if value, ok := transform.properties[tag]; ok {
			if comma {
				buffer.WriteString(", ")
			}
			buffer.WriteString(string(tag))
			buffer.WriteString(" = ")
			writePropertyValue(buffer, tag, value, indent)
			comma = true
		}
	}
	buffer.WriteString(" }")
}

func transformSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {

	case RotateX, RotateY, RotateZ:
		return setFloatProperty(properties, tag, value, 0, 1)

	case Rotate, SkewX, SkewY:
		return setAngleProperty(properties, tag, value)

	case ScaleX, ScaleY, ScaleZ:
		return setFloatProperty(properties, tag, value, -math.MaxFloat64, math.MaxFloat64)

	case Perspective, TranslateX, TranslateY, TranslateZ:
		return setSizeProperty(properties, tag, value)
	}

	return nil
}

func setTransformProperty(properties Properties, value any) bool {

	setObject := func(obj DataObject) bool {
		transform := NewTransform(nil)
		ok := true
		for i := 0; i < obj.PropertyCount(); i++ {
			if prop := obj.Property(i); prop.Type() == TextNode {
				if !transform.Set(PropertyName(prop.Tag()), prop.Text()) {
					ok = false
				}
			} else {
				ok = false
			}
		}

		if !ok && transform.empty() {
			return false
		}

		properties.setRaw(TransformTag, transform)
		return true
	}

	switch value := value.(type) {
	case Transform:
		properties.setRaw(TransformTag, value)
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

func getTransformProperty(properties Properties) Transform {
	if val := properties.getRaw(TransformTag); val != nil {
		if transform, ok := val.(Transform); ok {
			return transform
		}
	}
	return nil
}

func setTransformPropertyElement(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case Perspective, RotateX, RotateY, RotateZ, Rotate, SkewX, SkewY, ScaleX, ScaleY, ScaleZ, TranslateX, TranslateY, TranslateZ:
		var result []PropertyName = nil
		if transform := getTransformProperty(properties); transform != nil {
			if result = transformSet(transform, tag, value); result != nil {
				result = append(result, TransformTag)
			}
		} else {
			transform := NewTransform(nil)
			if result = transformSet(transform, tag, value); result != nil {
				properties.setRaw(TransformTag, transform)
				result = append(result, TransformTag)
			}
		}
		return result
	}

	ErrorLogF(`"Transform" interface does not support the "%s" property`, tag)
	return nil
}

/*
func getTransform3D(style Properties, session Session) bool {
	perspective, ok := sizeProperty(style, Perspective, session)
	return ok && perspective.Type != Auto && perspective.Value != 0
}
*/

func getPerspectiveOrigin(style Properties, session Session) (SizeUnit, SizeUnit) {
	x, _ := sizeProperty(style, PerspectiveOriginX, session)
	y, _ := sizeProperty(style, PerspectiveOriginY, session)
	return x, y
}

func getTransformOrigin(style Properties, session Session) (SizeUnit, SizeUnit, SizeUnit) {
	x, _ := sizeProperty(style, TransformOriginX, session)
	y, _ := sizeProperty(style, TransformOriginY, session)
	z, _ := sizeProperty(style, TransformOriginZ, session)
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

func (transform *transformData) transformCSS(session Session) string {

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if perspective, ok := sizeProperty(transform, Perspective, session); ok && perspective.Type != Auto && perspective.Value != 0 {
		buffer.WriteString(`perspective(`)
		buffer.WriteString(perspective.cssString("0", session))
		buffer.WriteString(") ")
	}

	skewX, skewY, skewOK := transform.getSkew(session)
	if skewOK {
		buffer.WriteString(`skew(`)
		buffer.WriteString(skewX.cssString())
		buffer.WriteRune(',')
		buffer.WriteString(skewY.cssString())
		buffer.WriteString(") ")
	}

	x, y, z := transform.getTranslate(session)
	if z.Type != Auto && z.Value != 0 {

		buffer.WriteString(`translate3d(`)
		buffer.WriteString(x.cssString("0", session))
		buffer.WriteRune(',')
		buffer.WriteString(y.cssString("0", session))
		buffer.WriteRune(',')
		buffer.WriteString(z.cssString("0", session))
		buffer.WriteString(") ")

	} else if (x.Type != Auto && x.Value != 0) || (y.Type != Auto && y.Value != 0) {

		buffer.WriteString(`translate(`)
		buffer.WriteString(x.cssString("0", session))
		buffer.WriteRune(',')
		buffer.WriteString(y.cssString("0", session))
		buffer.WriteString(") ")
	}

	scaleX, okScaleX := floatTextProperty(transform, ScaleX, session, 1)
	scaleY, okScaleY := floatTextProperty(transform, ScaleY, session, 1)
	scaleZ, okScaleZ := floatTextProperty(transform, ScaleZ, session, 1)
	if okScaleZ {

		buffer.WriteString(`scale3d(`)
		buffer.WriteString(scaleX)
		buffer.WriteRune(',')
		buffer.WriteString(scaleY)
		buffer.WriteRune(',')
		buffer.WriteString(scaleZ)
		buffer.WriteString(") ")

	} else if okScaleX || okScaleY {

		buffer.WriteString(`scale(`)
		buffer.WriteString(scaleX)
		buffer.WriteRune(',')
		buffer.WriteString(scaleY)
		buffer.WriteString(") ")
	}

	if angle, ok := angleProperty(transform, Rotate, session); ok {
		rotateX, xOK := floatTextProperty(transform, RotateX, session, 1)
		rotateY, yOK := floatTextProperty(transform, RotateY, session, 1)
		rotateZ, zOK := floatTextProperty(transform, RotateZ, session, 1)

		if xOK || yOK || zOK {

			buffer.WriteString(`rotate3d(`)
			buffer.WriteString(rotateX)
			buffer.WriteRune(',')
			buffer.WriteString(rotateY)
			buffer.WriteRune(',')
			buffer.WriteString(rotateZ)
			buffer.WriteRune(',')
			buffer.WriteString(angle.cssString())
			buffer.WriteString(") ")

		} else {

			buffer.WriteString(`rotate(`)
			buffer.WriteString(angle.cssString())
			buffer.WriteString(") ")
		}
	}

	length := buffer.Len()
	if length == 0 {
		return ""
	}
	result := buffer.String()
	return result[:length-1]
}

func (style *viewStyle) writeViewTransformCSS(builder cssBuilder, session Session) {
	x, y := getPerspectiveOrigin(style, session)
	z := AutoSize()
	if css := transformOriginCSS(x, y, z, session); css != "" {
		builder.add(`perspective-origin`, css)
	}

	if backfaceVisible, ok := boolProperty(style, BackfaceVisible, session); ok {
		if backfaceVisible {
			builder.add(`backface-visibility`, `visible`)
		} else {
			builder.add(`backface-visibility`, `hidden`)
		}
	}

	x, y, z = getTransformOrigin(style, session)
	if css := transformOriginCSS(x, y, z, session); css != "" {
		builder.add(`transform-origin`, css)
	}

	if transform := getTransformProperty(style); transform != nil {
		builder.add(`transform`, transform.transformCSS(session))
	}
}

func transformOriginCSS(x, y, z SizeUnit, session Session) string {
	if z.Type == Auto && x.Type == Auto && y.Type == Auto {
		return ""
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	if x.Type == SizeInPercent {
		switch x.Value {
		case 0:
			buffer.WriteString("left")
		case 50:
			buffer.WriteString("center")
		case 100:
			buffer.WriteString("right")

		default:
			buffer.WriteString(x.cssString("center", session))
		}
	} else {
		buffer.WriteString(x.cssString("center", session))
	}

	buffer.WriteRune(' ')

	if y.Type == SizeInPercent {
		switch y.Value {
		case 0:
			buffer.WriteString("top")
		case 50:
			buffer.WriteString("center")
		case 100:
			buffer.WriteString("bottom")

		default:
			buffer.WriteString(y.cssString("center", session))
		}
	} else {
		buffer.WriteString(y.cssString("center", session))
	}

	if z.Type != Auto && z.Value != 0 {
		buffer.WriteRune(' ')
		buffer.WriteString(z.cssString("0", session))
	}

	return buffer.String()
}

/*
func (view *viewData) updateTransformProperty(tag PropertyName) bool {
	htmlID := view.htmlID()
	session := view.session

	switch tag {
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
				session.updateCSSProperty(htmlID, string(BackfaceVisible), "visible")
			} else {
				session.updateCSSProperty(htmlID, string(BackfaceVisible), "hidden")
			}
		}

	case OriginX, OriginY, OriginZ:
		x, y, z := getOrigin(view, session)
		value := ""

		if z.Type != Auto {
			value = x.cssString("50%", session) + " " + y.cssString("50%", session) + " " + z.cssString("50%", session)
		} else if x.Type != Auto || y.Type != Auto {
			value = x.cssString("50%", session) + " " + y.cssString("50%", session)
		}
		session.updateCSSProperty(htmlID, "transform-origin", value)

	case TransformTag, SkewX, SkewY, TranslateX, TranslateY, TranslateZ,
		ScaleX, ScaleY, ScaleZ, Rotate, RotateX, RotateY, RotateZ:
		if transform := getTransformProperty(view); transform != nil {
			session.updateCSSProperty(htmlID, "transform", transform.transformCSS(session))
		} else {
			session.updateCSSProperty(htmlID, "transform", "")
		}

	default:
		return false
	}

	return true
}
*/
