package rui

import (
	"fmt"
	"math"
	"strings"
)

// Constants for [TransformProperty] specific properties
const (
	// Transform is the constant for "transform" property tag.
	//
	// Used by View.
	// Specify translation, scale and rotation over x, y and z axes as well as a distortion of a view along x and y axes.
	//
	// Supported types: TransformProperty, string.
	//
	// See TransformProperty description for more details.
	//
	// Conversion rules:
	//   - TransformProperty - stored as is, no conversion performed.
	//   - string - string representation of TransformProperty interface. Example: "_{translate-x = 10px, scale-y = 1.1}".
	Transform PropertyName = "transform"

	// Perspective is the constant for "perspective" property tag.
	//
	// Used by View, TransformProperty.
	// Distance between the z-plane and the user in order to give a 3D-positioned element some perspective. Each 3D element
	// with z > 0 becomes larger, each 3D-element with z < 0 becomes smaller. The default value is 0 (no 3D effects).
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	Perspective PropertyName = "perspective"

	// PerspectiveOriginX is the constant for "perspective-origin-x" property tag.
	//
	// Used by View.
	// x-coordinate of the position at which the viewer is looking. It is used as the vanishing point by the "perspective"
	// property. The default value is 50%.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	PerspectiveOriginX PropertyName = "perspective-origin-x"

	// PerspectiveOriginY is the constant for "perspective-origin-y" property tag.
	//
	// Used by View.
	// y-coordinate of the position at which the viewer is looking. It is used as the vanishing point by the "perspective"
	// property. The default value is 50%.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	PerspectiveOriginY PropertyName = "perspective-origin-y"

	// BackfaceVisible is the constant for "backface-visibility" property tag.
	//
	// Used by View.
	// Controls whether the back face of a view is visible when turned towards the user. Default value is true.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", "1" - Back face is visible when turned towards the user.
	//   - false, 0, "false", "no", "off", "0" - Back face is hidden, effectively making the view invisible when turned away from the user.
	BackfaceVisible PropertyName = "backface-visibility"

	// TransformOriginX is the constant for "transform-origin-x" property tag.
	//
	// Used by View.
	// x-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TransformOriginX PropertyName = "transform-origin-x"

	// TransformOriginY is the constant for "transform-origin-y" property tag.
	//
	// Used by View.
	// y-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TransformOriginY PropertyName = "transform-origin-y"

	// TransformOriginZ is the constant for "transform-origin-z" property tag.
	//
	// Used by View.
	// z-coordinate of the point around which a view transformation is applied. The default value is 50%.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TransformOriginZ PropertyName = "transform-origin-z"

	// TranslateX is the constant for "translate-x" property tag.
	//
	// Used by View, TransformProperty.
	//
	// x-axis translation value of a 2D/3D translation.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TranslateX PropertyName = "translate-x"

	// TranslateY is the constant for "translate-y" property tag.
	//
	// Used by View, TransformProperty.
	//
	// x-axis translation value of a 2D/3D translation.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TranslateY PropertyName = "translate-y"

	// TranslateZ is the constant for "translate-z" property tag.
	//
	// Used by View, TransformProperty.
	//
	// z-axis translation value of a 3D translation.
	//
	// Supported types: SizeUnit, SizeFunc, string, float, int.
	//
	// Internal type is SizeUnit, other types converted to it during assignment.
	// See SizeUnit description for more details.
	TranslateZ PropertyName = "translate-z"

	// ScaleX is the constant for "scale-x" property tag.
	//
	// Used by View, TransformProperty.
	//
	// x-axis scaling value of a 2D/3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	ScaleX PropertyName = "scale-x"

	// ScaleY is the constant for "scale-y" property tag.
	//
	// Used by View, TransformProperty.
	//
	// y-axis scaling value of a 2D/3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	ScaleY PropertyName = "scale-y"

	// ScaleZ is the constant for "scale-z" property tag.
	//
	// Used by View, TransformProperty.
	//
	// z-axis scaling value of a 3D scale. The original scale is 1. Values between 0 and 1 are used to decrease original
	// scale, more than 1 - to increase. The default value is 1.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	ScaleZ PropertyName = "scale-z"

	// Rotate is the constant for "rotate" property tag.
	//
	// Used by View, TransformProperty.
	//
	// Angle of the view rotation. A positive angle denotes a clockwise rotation, a negative angle a counter-clockwise.
	//
	// Supported types: AngleUnit, string, float, int.
	//
	// Internal type is AngleUnit, other types will be converted to it during assignment.
	// See AngleUnit description for more details.
	//
	// Conversion rules:
	//   - AngleUnit - stored as is, no conversion performed.
	//   - string - must contain string representation of AngleUnit. If numeric value will be provided without any suffix then AngleUnit with value and Radian value type will be created.
	//   - float - a new AngleUnit value will be created with Radian as a type.
	//   - int - a new AngleUnit value will be created with Radian as a type.
	Rotate PropertyName = "rotate"

	// RotateX is the constant for "rotate-x" property tag.
	//
	// Used by View, TransformProperty.
	//
	// x-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	RotateX PropertyName = "rotate-x"

	// RotateY is the constant for "rotate-y" property tag.
	//
	// Used by View, TransformProperty.
	//
	// y-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	RotateY PropertyName = "rotate-y"

	// RotateZ is the constant for "rotate-z" property tag.
	//
	// Used by View, TransformProperty.
	//
	// z-coordinate of the vector denoting the axis of rotation in range 0 to 1.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	RotateZ PropertyName = "rotate-z"

	// SkewX is the constant for "skew-x" property tag.
	//
	// Used by View, TransformProperty.
	//
	// Angle to use to distort the element along the abscissa. The default value is 0.
	//
	// Supported types: AngleUnit, string, float, int.
	//
	// Internal type is AngleUnit, other types will be converted to it during assignment.
	// See AngleUnit description for more details.
	//
	// Conversion rules:
	//   - AngleUnit - stored as is, no conversion performed.
	//   - string - must contain string representation of AngleUnit. If numeric value will be provided without any suffix then AngleUnit with value and Radian value type will be created.
	//   - float - a new AngleUnit value will be created with Radian as a type.
	//   - int - a new AngleUnit value will be created with Radian as a type.
	SkewX PropertyName = "skew-x"

	// SkewY is the constant for "skew-y" property tag.
	//
	// Used by View, TransformProperty.
	//
	// Angle to use to distort the element along the ordinate. The default value is 0.
	//
	// Supported types: AngleUnit, string, float, int.
	//
	// Internal type is AngleUnit, other types will be converted to it during assignment.
	// See AngleUnit description for more details.
	//
	// Conversion rules:
	//   - AngleUnit - stored as is, no conversion performed.
	//   - string - must contain string representation of AngleUnit. If numeric value will be provided without any suffix then AngleUnit with value and Radian value type will be created.
	//   - float - a new AngleUnit value will be created with Radian as a type.
	//   - int - a new AngleUnit value will be created with Radian as a type.
	SkewY PropertyName = "skew-y"
)

// TransformProperty interface specifies view transformation parameters: the x-, y-, and z-axis translation values,
// the x-, y-, and z-axis scaling values, the angle to use to distort the element along the abscissa and ordinate,
// the angle of the view rotation.
//
// Valid property tags: Perspective ("perspective"),  TranslateX ("translate-x"), TranslateY ("translate-y"), TranslateZ ("translate-z"),
// ScaleX ("scale-x"), ScaleY ("scale-y"), ScaleZ ("scale-z"), Rotate ("rotate"), RotateX ("rotate-x"),
// RotateY ("rotate-y"), RotateZ ("rotate-z"), SkewX ("skew-x"), and SkewY ("skew-y")
type TransformProperty interface {
	Properties
	fmt.Stringer
	stringWriter
	transformCSS(session Session) string
	getSkew(session Session) (AngleUnit, AngleUnit, bool)
	getTranslate(session Session) (SizeUnit, SizeUnit, SizeUnit)
}

type transformPropertyData struct {
	dataProperty
}

// NewTransform creates a new transform property data and return its interface
//
// The following properties can be used:
//
// Perspective ("perspective"),  TranslateX ("translate-x"), TranslateY ("translate-y"), TranslateZ ("translate-z"),
// ScaleX ("scale-x"), ScaleY ("scale-y"), ScaleZ ("scale-z"), Rotate ("rotate"), RotateX ("rotate-x"),
// RotateY ("rotate-y"), RotateZ ("rotate-z"), SkewX ("skew-x"), and SkewY ("skew-y")
func NewTransformProperty(params Params) TransformProperty {
	transform := new(transformPropertyData)
	transform.init()

	for tag, value := range params {
		transform.Set(tag, value)
	}
	return transform
}

func (transform *transformPropertyData) init() {
	transform.dataProperty.init()
	transform.normalize = normalizeTransformTag
	transform.set = transformSet
	transform.supportedProperties = []PropertyName{
		RotateX, RotateY, RotateZ, Rotate, SkewX, SkewY, ScaleX, ScaleY, ScaleZ,
		Perspective, TranslateX, TranslateY, TranslateZ,
	}
}

func normalizeTransformTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)

	name := string(tag)
	if strings.HasPrefix(name, "push-") {
		tag = PropertyName(name[5:])
	}

	return tag
}

func (transform *transformPropertyData) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)
	transform.writeString(buffer, "")
	return buffer.String()
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

func valueToTransformProperty(value any) TransformProperty {

	parseObject := func(obj DataObject) TransformProperty {
		transform := NewTransformProperty(nil)
		ok := true
		for prop := range obj.Properties() {
			if prop.Type() == TextNode {
				if !transform.Set(PropertyName(prop.Tag()), prop.Text()) {
					ok = false
				}
			} else {
				ok = false
			}
		}

		if !ok && transform.empty() {
			return nil
		}
		return transform
	}

	switch value := value.(type) {
	case TransformProperty:
		return value

	case DataObject:
		return parseObject(value)

	case DataNode:
		if obj := value.Object(); obj != nil {
			return parseObject(obj)
		}

	case string:
		obj, err := ParseDataText(value)
		if err != nil {
			ErrorLog(err.Error())
		} else {
			return parseObject(obj)
		}
	}

	return nil
}

func setTransformProperty(properties Properties, tag PropertyName, value any) bool {
	if transform := valueToTransformProperty(value); transform != nil {
		properties.setRaw(tag, transform)
		return true
	}

	notCompatibleType(tag, value)
	return false
}

func getTransformProperty(properties Properties, tag PropertyName) TransformProperty {
	if val := properties.getRaw(tag); val != nil {
		if transform, ok := val.(TransformProperty); ok {
			return transform
		}
	}
	return nil
}

func setTransformPropertyElement(properties Properties, tag, transformTag PropertyName, value any) []PropertyName {
	srcTag := tag
	tag = normalizeTransformTag(tag)
	switch tag {
	case Perspective, RotateX, RotateY, RotateZ, Rotate, SkewX, SkewY, ScaleX, ScaleY, ScaleZ, TranslateX, TranslateY, TranslateZ:
		var result []PropertyName = nil
		if transform := getTransformProperty(properties, transformTag); transform != nil {
			if result = transformSet(transform, tag, value); result != nil {
				result = []PropertyName{srcTag, transformTag}
			}
		} else {
			transform := NewTransformProperty(nil)
			if result = transformSet(transform, tag, value); result != nil {
				properties.setRaw(transformTag, transform)
				result = []PropertyName{srcTag, transformTag}
			}
		}
		return result
	}

	ErrorLogF(`"TransformProperty" interface does not support the "%s" property`, tag)
	return nil
}

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

func (transform *transformPropertyData) getSkew(session Session) (AngleUnit, AngleUnit, bool) {
	skewX, okX := angleProperty(transform, SkewX, session)
	skewY, okY := angleProperty(transform, SkewY, session)
	return skewX, skewY, okX || okY
}

func (transform *transformPropertyData) getTranslate(session Session) (SizeUnit, SizeUnit, SizeUnit) {
	x, _ := sizeProperty(transform, TranslateX, session)
	y, _ := sizeProperty(transform, TranslateY, session)
	z, _ := sizeProperty(transform, TranslateZ, session)
	return x, y, z
}

func (transform *transformPropertyData) transformCSS(session Session) string {

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
		buffer.WriteString(x.cssString("0px", session))
		buffer.WriteRune(',')
		buffer.WriteString(y.cssString("0px", session))
		buffer.WriteRune(',')
		buffer.WriteString(z.cssString("0px", session))
		buffer.WriteString(") ")

	} else if (x.Type != Auto && x.Value != 0) || (y.Type != Auto && y.Value != 0) {

		buffer.WriteString(`translate(`)
		buffer.WriteString(x.cssString("0px", session))
		buffer.WriteRune(',')
		buffer.WriteString(y.cssString("0px", session))
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

	if transform := getTransformProperty(style, Transform); transform != nil {
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

// GetTransform returns a view transform:  translation, scale and rotation over x, y and z axes as well as a distortion of a view along x and y axes.
// The default value is nil (no transform)
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetTransform(view View, subviewID ...string) TransformProperty {
	return transformStyledProperty(view, subviewID, Transform)
}

// GetPerspective returns a distance between the z = 0 plane and the user in order to give a 3D-positioned
// element some perspective. Each 3D element with z > 0 becomes larger; each 3D-element with z < 0 becomes smaller.
// The default value is 0 (no 3D effects).
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetPerspective(view View, subviewID ...string) SizeUnit {
	return sizeStyledProperty(view, subviewID, Perspective, false)
}

// GetPerspectiveOrigin returns a x- and y-coordinate of the position at which the viewer is looking.
// It is used as the vanishing point by the Perspective property. The default value is (50%, 50%).
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetPerspectiveOrigin(view View, subviewID ...string) (SizeUnit, SizeUnit) {
	view = getSubview(view, subviewID)
	if view == nil {
		return AutoSize(), AutoSize()
	}
	return getPerspectiveOrigin(view, view.Session())
}

// GetBackfaceVisible returns a bool property that sets whether the back face of an element is
// visible when turned towards the user. Values:
// true - the back face is visible when turned towards the user (default value).
// false - the back face is hidden, effectively making the element invisible when turned away from the user.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetBackfaceVisible(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, BackfaceVisible, false)
}

// GetTransformOrigin returns a x-, y-, and z-coordinate of the point around which a view transformation is applied.
// The default value is (50%, 50%, 50%).
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetTransformOrigin(view View, subviewID ...string) (SizeUnit, SizeUnit, SizeUnit) {
	view = getSubview(view, subviewID)
	if view == nil {
		return AutoSize(), AutoSize(), AutoSize()
	}
	return getTransformOrigin(view, view.Session())
}

// GetTranslate returns a x-, y-, and z-axis translation value of a 2D/3D translation
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetTranslate(view View, subviewID ...string) (SizeUnit, SizeUnit, SizeUnit) {
	if transform := GetTransform(view, subviewID...); transform != nil {
		return transform.getTranslate(view.Session())
	}
	return AutoSize(), AutoSize(), AutoSize()
}

// GetSkew returns a angles to use to distort the element along the abscissa (x-axis)
// and the ordinate (y-axis). The default value is 0.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetSkew(view View, subviewID ...string) (AngleUnit, AngleUnit) {
	if transform := GetTransform(view, subviewID...); transform != nil {
		x, y, _ := transform.getSkew(view.Session())
		return x, y
	}
	return AngleUnit{Value: 0, Type: Radian}, AngleUnit{Value: 0, Type: Radian}
}

// GetScale returns a x-, y-, and z-axis scaling value of a 2D/3D scale. The default value is 1.
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetScale(view View, subviewID ...string) (float64, float64, float64) {
	if transform := GetTransform(view, subviewID...); transform != nil {
		session := view.Session()
		x, _ := floatProperty(transform, ScaleX, session, 1)
		y, _ := floatProperty(transform, ScaleY, session, 1)
		z, _ := floatProperty(transform, ScaleZ, session, 1)
		return x, y, z
	}
	return 1, 1, 1
}

// GetRotate returns a x-, y, z-coordinate of the vector denoting the axis of rotation, and the angle of the view rotation
//
// The second argument (subviewID) specifies the path to the child element whose value needs to be returned.
// If it is not specified then a value from the first argument (view) is returned.
func GetRotate(view View, subviewID ...string) (float64, float64, float64, AngleUnit) {
	if transform := GetTransform(view, subviewID...); transform != nil {
		session := view.Session()
		angle, _ := angleProperty(transform, Rotate, view.Session())
		rotateX, _ := floatProperty(transform, RotateX, session, 1)
		rotateY, _ := floatProperty(transform, RotateY, session, 1)
		rotateZ, _ := floatProperty(transform, RotateZ, session, 1)
		return rotateX, rotateY, rotateZ, angle
	}
	return 0, 0, 0, AngleUnit{Value: 0, Type: Radian}
}
