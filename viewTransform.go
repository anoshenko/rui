package rui

import (
	"strconv"
)

const (
	// Perspective is the name of the SizeUnit property that determines the distance between the z = 0 plane
	// and the user in order to give a 3D-positioned element some perspective. Each 3D element
	// with z > 0 becomes larger; each 3D-element with z < 0 becomes smaller.
	// The default value is 0 (no 3D effects).
	Perspective = "perspective"
	// PerspectiveOriginX is the name of the SizeUnit property that determines the x-coordinate of the position
	// at which the viewer is looking. It is used as the vanishing point by the Perspective property.
	// The default value is 50%.
	PerspectiveOriginX = "perspective-origin-x"
	// PerspectiveOriginY is the name of the SizeUnit property that determines the y-coordinate of the position
	// at which the viewer is looking. It is used as the vanishing point by the Perspective property.
	// The default value is 50%.
	PerspectiveOriginY = "perspective-origin-y"
	// BackfaceVisible is the name of the bool property that sets whether the back face of an element is
	// visible when turned towards the user. Values:
	//   true - the back face is visible when turned towards the user (default value).
	//   false - the back face is hidden, effectively making the element invisible when turned away from the user.
	BackfaceVisible = "backface-visibility"
	// OriginX is the name of the SizeUnit property that determines the x-coordinate of the point around which
	// a view transformation is applied.
	// The default value is 50%.
	OriginX = "origin-x"
	// OriginY is the name of the SizeUnit property that determines the y-coordinate of the point around which
	// a view transformation is applied.
	// The default value is 50%.
	OriginY = "origin-y"
	// OriginZ is the name of the SizeUnit property that determines the z-coordinate of the point around which
	// a view transformation is applied.
	// The default value is 50%.
	OriginZ = "origin-z"
	// TranslateX is the name of the SizeUnit property that specify the x-axis translation value
	// of a 2D/3D translation
	TranslateX = "translate-x"
	// TranslateY is the name of the SizeUnit property that specify the y-axis translation value
	// of a 2D/3D translation
	TranslateY = "translate-y"
	// TranslateZ is the name of the SizeUnit property that specify the z-axis translation value
	// of a 3D translation
	TranslateZ = "translate-z"
	// ScaleX is the name of the float property that specify the x-axis scaling value of a 2D/3D scale
	// The default value is 1.
	ScaleX = "scale-x"
	// ScaleY is the name of the float property that specify the y-axis scaling value of a 2D/3D scale
	// The default value is 1.
	ScaleY = "scale-y"
	// ScaleZ is the name of the float property that specify the z-axis scaling value of a 3D scale
	// The default value is 1.
	ScaleZ = "scale-z"
	// Rotate is the name of the AngleUnit property that determines the angle of the view rotation.
	// A positive angle denotes a clockwise rotation, a negative angle a counter-clockwise one.
	Rotate = "rotate"
	// RotateX is the name of the float property that determines the x-coordinate of the vector denoting
	// the axis of rotation which could between 0 and 1.
	RotateX = "rotate-x"
	// RotateY is the name of the float property that determines the y-coordinate of the vector denoting
	// the axis of rotation which could between 0 and 1.
	RotateY = "rotate-y"
	// RotateZ is the name of the float property that determines the z-coordinate of the vector denoting
	// the axis of rotation which could between 0 and 1.
	RotateZ = "rotate-z"
	// SkewX is the name of the AngleUnit property that representing the angle to use to distort
	// the element along the abscissa. The default value is 0.
	SkewX = "skew-x"
	// SkewY is the name of the AngleUnit property that representing the angle to use to distort
	// the element along the ordinate. The default value is 0.
	SkewY = "skew-y"
)

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

func getSkew(style Properties, session Session) (AngleUnit, AngleUnit) {
	skewX, _ := angleProperty(style, SkewX, session)
	skewY, _ := angleProperty(style, SkewY, session)
	return skewX, skewY
}

func getTranslate(style Properties, session Session) (SizeUnit, SizeUnit, SizeUnit) {
	x, _ := sizeProperty(style, TranslateX, session)
	y, _ := sizeProperty(style, TranslateY, session)
	z, _ := sizeProperty(style, TranslateZ, session)
	return x, y, z
}

func getScale(style Properties, session Session) (float64, float64, float64) {
	scaleX, _ := floatProperty(style, ScaleX, session, 1)
	scaleY, _ := floatProperty(style, ScaleY, session, 1)
	scaleZ, _ := floatProperty(style, ScaleZ, session, 1)
	return scaleX, scaleY, scaleZ
}

func getRotate(style Properties, session Session) (float64, float64, float64, AngleUnit) {
	rotateX, _ := floatProperty(style, RotateX, session, 1)
	rotateY, _ := floatProperty(style, RotateY, session, 1)
	rotateZ, _ := floatProperty(style, RotateZ, session, 1)
	angle, _ := angleProperty(style, Rotate, session)
	return rotateX, rotateY, rotateZ, angle
}

func (style *viewStyle) transform(session Session) string {

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	skewX, skewY := getSkew(style, session)
	if skewX.Value != 0 || skewY.Value != 0 {
		buffer.WriteString(`skew(`)
		buffer.WriteString(skewX.cssString())
		buffer.WriteRune(',')
		buffer.WriteString(skewY.cssString())
		buffer.WriteRune(')')
	}

	x, y, z := getTranslate(style, session)
	scaleX, scaleY, scaleZ := getScale(style, session)
	if getTransform3D(style, session) {
		if (x.Type != Auto && x.Value != 0) || (y.Type != Auto && y.Value != 0) || (z.Type != Auto && z.Value != 0) {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`translate3d(`)
			buffer.WriteString(x.cssString("0"))
			buffer.WriteRune(',')
			buffer.WriteString(y.cssString("0"))
			buffer.WriteRune(',')
			buffer.WriteString(z.cssString("0"))
			buffer.WriteRune(')')
		}

		if scaleX != 1 || scaleY != 1 || scaleZ != 1 {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`scale3d(`)
			buffer.WriteString(strconv.FormatFloat(scaleX, 'g', -1, 64))
			buffer.WriteRune(',')
			buffer.WriteString(strconv.FormatFloat(scaleY, 'g', -1, 64))
			buffer.WriteRune(',')
			buffer.WriteString(strconv.FormatFloat(scaleZ, 'g', -1, 64))
			buffer.WriteRune(')')
		}

		rotateX, rotateY, rotateZ, angle := getRotate(style, session)
		if angle.Value != 0 && (rotateX != 0 || rotateY != 0 || rotateZ != 0) {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`rotate3d(`)
			buffer.WriteString(strconv.FormatFloat(rotateX, 'g', -1, 64))
			buffer.WriteRune(',')
			buffer.WriteString(strconv.FormatFloat(rotateY, 'g', -1, 64))
			buffer.WriteRune(',')
			buffer.WriteString(strconv.FormatFloat(rotateZ, 'g', -1, 64))
			buffer.WriteRune(',')
			buffer.WriteString(angle.cssString())
			buffer.WriteRune(')')
		}
	} else {
		if (x.Type != Auto && x.Value != 0) || (y.Type != Auto && y.Value != 0) {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`translate(`)
			buffer.WriteString(x.cssString("0"))
			buffer.WriteRune(',')
			buffer.WriteString(y.cssString("0"))
			buffer.WriteRune(')')
		}

		if scaleX != 1 || scaleY != 1 {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(`scale(`)
			buffer.WriteString(strconv.FormatFloat(scaleX, 'g', -1, 64))
			buffer.WriteRune(',')
			buffer.WriteString(strconv.FormatFloat(scaleY, 'g', -1, 64))
			buffer.WriteRune(')')
		}

		angle, _ := angleProperty(style, Rotate, session)
		if angle.Value != 0 {
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
	if getTransform3D(style, session) {
		if perspective, ok := sizeProperty(style, Perspective, session); ok && perspective.Type != Auto && perspective.Value != 0 {
			builder.add(`perspective`, perspective.cssString("0"))
		}

		x, y := getPerspectiveOrigin(style, session)
		if x.Type != Auto || y.Type != Auto {
			builder.addValues(`perspective-origin`, ` `, x.cssString("50%"), y.cssString("50%"))
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
			builder.addValues(`transform-origin`, ` `, x.cssString("50%"), y.cssString("50%"), z.cssString("0"))
		}
	} else {
		x, y, _ := getOrigin(style, session)
		if x.Type != Auto || y.Type != Auto {
			builder.addValues(`transform-origin`, ` `, x.cssString("50%"), y.cssString("50%"))
		}
	}

	builder.add(`transform`, style.transform(session))
}

func (view *viewData) updateTransformProperty(tag string) bool {
	htmlID := view.htmlID()
	session := view.session

	switch tag {
	case Perspective:
		updateCSSStyle(htmlID, session)

	case PerspectiveOriginX, PerspectiveOriginY:
		if getTransform3D(view, session) {
			x, y := GetPerspectiveOrigin(view, "")
			value := ""
			if x.Type != Auto || y.Type != Auto {
				value = x.cssString("50%") + " " + y.cssString("50%")
			}
			updateCSSProperty(htmlID, "perspective-origin", value, session)
		}

	case BackfaceVisible:
		if getTransform3D(view, session) {
			if GetBackfaceVisible(view, "") {
				updateCSSProperty(htmlID, BackfaceVisible, "visible", session)
			} else {
				updateCSSProperty(htmlID, BackfaceVisible, "hidden", session)
			}
		}

	case OriginX, OriginY, OriginZ:
		x, y, z := getOrigin(view, session)
		value := ""
		if getTransform3D(view, session) {
			if x.Type != Auto || y.Type != Auto || z.Type != Auto {
				value = x.cssString("50%") + " " + y.cssString("50%") + " " + z.cssString("50%")
			}
		} else {
			if x.Type != Auto || y.Type != Auto {
				value = x.cssString("50%") + " " + y.cssString("50%")
			}
		}
		updateCSSProperty(htmlID, "transform-origin", value, session)

	case SkewX, SkewY, TranslateX, TranslateY, TranslateZ, ScaleX, ScaleY, ScaleZ, Rotate, RotateX, RotateY, RotateZ:
		updateCSSProperty(htmlID, "transform", view.transform(session), session)

	default:
		return false
	}

	return true
}
