package rui

/*
import (
	"fmt"
	"strconv"
)

type AnimationTags struct {
	Tag        string
	Start, End interface{}
}

type AnimationKeyFrame struct {
	KeyFrame       int
	TimingFunction string
	Params         Params
}

type AnimationScenario interface {
	fmt.Stringer
	ruiStringer
	Name() string
	cssString(session Session) string
}

type animationScenario struct {
	name      string
	tags      []AnimationTags
	keyFrames []AnimationKeyFrame
	cssText   string
}

var animationScenarios = []string{}

func addAnimationScenario(name string) string {
	animationScenarios = append(animationScenarios, name)
	return name
}

func registerAnimationScenario() string {
	find := func(text string) bool {
		for _, scenario := range animationScenarios {
			if scenario == text {
				return true
			}
		}
		return false
	}

	n := 1
	name := fmt.Sprintf("scenario%08d", n)
	for find(name) {
		n++
		name = fmt.Sprintf("scenario%08d", n)
	}

	animationScenarios = append(animationScenarios, name)
	return name
}

func NewAnimationScenario(tags []AnimationTags, keyFrames []AnimationKeyFrame) AnimationScenario {
	if tags == nil {
		ErrorLog(`Nil "tags" argument is not allowed.`)
		return nil
	}

	if len(tags) == 0 {
		ErrorLog(`An empty "tags" argument is not allowed.`)
		return nil
	}

	animation := new(animationScenario)
	animation.tags = tags
	if keyFrames == nil && len(keyFrames) > 0 {
		animation.keyFrames = keyFrames
	}
	animation.name = registerAnimationScenario()

	return animation
}

func (animation *animationScenario) Name() string {
	return animation.name
}

func (animation *animationScenario) String() string {
	writer := newRUIWriter()
	animation.ruiString(writer)
	return writer.finish()
}

func (animation *animationScenario) ruiString(writer ruiWriter) {
	// TODO
}

func valueToCSS(tag string, value interface{}, session Session) string {
	if value == nil {
		return ""
	}

	convertFloat := func(val float64) string {
		if _, ok := sizeProperties[tag]; ok {
			return fmt.Sprintf("%gpx", val)
		}
		return fmt.Sprintf("%g", val)
	}

	switch value := value.(type) {
	case string:
		value, ok := session.resolveConstants(value)
		if !ok {
			return ""
		}
		if _, ok := sizeProperties[tag]; ok {
			var size SizeUnit
			if size.SetValue(value) {
				return size.cssString("auto")
			}
			return ""
		}
		if isPropertyInList(tag, colorProperties) {
			var color Color
			if color.SetValue(value) {
				return color.cssString()
			}
			return ""
		}
		if isPropertyInList(tag, angleProperties) {
			var angle AngleUnit
			if angle.SetValue(value) {
				return angle.cssString()
			}
			return ""
		}
		if _, ok := enumProperties[tag]; ok {
			var size SizeUnit
			if size.SetValue(value) {
				return size.cssString("auto")
			}
			return ""
		}
		return value

	case SizeUnit:
		return value.cssString("auto")

	case AngleUnit:
		return value.cssString()

	case Color:
		return value.cssString()

	case float32:
		return convertFloat(float64(value))

	case float64:
		return convertFloat(value)

	default:
		if n, ok := isInt(value); ok {
			if prop, ok := enumProperties[tag]; ok {
				values := prop.cssValues
				if n >= 0 && n < len(values) {
					return values[n]
				}
				return ""
			}

			return convertFloat(float64(n))
		}
	}
	return ""
}

func (animation *animationScenario) cssString(session Session) string {
	if animation.cssText != "" {

		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)

		writeValue := func(tag string, value interface{}) {
			if cssValue := valueToCSS(tag, value); cssValue != "" {
				buffer.WriteString("    ")
				buffer.WriteString(tag)
				buffer.WriteString(": ")
				buffer.WriteString(cssValue)
				buffer.WriteString(";\n")
			}
		}

		buffer.WriteString(`@keyframes `)
		buffer.WriteString(animation.name)

		buffer.WriteString(" {\n  from {\n")
		for _, property := range animation.tags {
			writeValue(property.Tag, property.Start)
		}

		buffer.WriteString(" }\n  to {\n")
		for _, property := range animation.tags {
			writeValue(property.Tag, property.End)
		}
		buffer.WriteString(" }\n")

		if animation.keyFrames != nil {
			for _, keyFrame := range animation.keyFrames {
				if keyFrame.KeyFrame > 0 && keyFrame.KeyFrame < 100 &&
					keyFrame.Params != nil && len(keyFrame.Params) > 0 {

					buffer.WriteString("  ")
					buffer.WriteString(strconv.Itoa(keyFrame.KeyFrame))
					buffer.WriteString("% {\n")
					for tag, value := range keyFrame.Params {
						writeValue(tag, value)
					}
					buffer.WriteString(" }\n")

				}
			}
		}
		buffer.WriteString("}\n")

		animation.cssText = buffer.String()
	}

	return animation.cssText
}
*/
