package rui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Constants which related to view's animation
const (
	// Animation is the constant for "animation" property tag.
	//
	// Used by View.
	// Sets and starts animations.
	//
	// Supported types: AnimationProperty, []AnimationProperty.
	//
	// Internal type is []AnimationProperty, other types converted to it during assignment.
	// See AnimationProperty description for more details.
	Animation PropertyName = "animation"

	// AnimationPaused is the constant for "animation-paused" property tag.
	//
	// Used by AnimationProperty.
	// Controls whether the animation is running or paused.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", or "1" - Animation is paused.
	//   - false, 0, "false", "no", "off", or "0" - Animation is playing.
	AnimationPaused PropertyName = "animation-paused"

	// Transition is the constant for "transition" property tag.
	//
	// Used by View.
	//
	// Sets transition animation of view properties. Each provided property must contain AnimationProperty which describe how
	// particular property will be animated on property value change. Transition animation can be applied to properties of the
	// type SizeUnit, Color, AngleUnit, float64 and composite properties that contain elements of the listed types(for
	// example, "shadow", "border", etc.). If we'll try to animate other properties with internal type like bool or
	// string no error will occur, simply there will be no animation.
	//
	// Supported types: Params.
	//
	// See Params description for more details.
	Transition PropertyName = "transition"

	// PropertyTag is the constant for "property" property tag.
	//
	// Used by AnimationProperty.
	//
	// Describes a scenario for changing a View's property. Used only for animation script.
	//
	// Supported types: []AnimatedProperty, AnimatedProperty.
	//
	// Internal type is []AnimatedProperty, other types converted to it during assignment.
	// See AnimatedProperty description for more details.
	PropertyTag PropertyName = "property"

	// Duration is the constant for "duration" property tag.
	//
	// Used by AnimationProperty.
	//
	// Sets the length of time in seconds that an animation takes to complete one cycle.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Duration PropertyName = "duration"

	// Delay is the constant for "delay" property tag.
	//
	// Used by AnimationProperty.
	//
	// Specifies the amount of time in seconds to wait from applying the animation to an element before beginning to perform
	// the animation. The animation can start later, immediately from its beginning or immediately and partway through the
	// animation.
	//
	// Supported types: float, int, string.
	//
	// Internal type is float, other types converted to it during assignment.
	Delay PropertyName = "delay"

	// TimingFunction is the constant for "timing-function" property tag.
	//
	// Used by AnimationProperty.
	//
	// Set how an animation progresses through the duration of each cycle.
	//
	// Supported types: string.
	//
	// Values:
	//   - "ease" (EaseTiming) - Speed increases towards the middle and slows down at the end.
	//   - "ease-in" (EaseInTiming) - Speed is slow at first, but increases in the end.
	//   - "ease-out" (EaseOutTiming) - Speed is fast at first, but decreases in the end.
	//   - "ease-in-out" (EaseInOutTiming) - Speed is slow at first, but quickly increases and at the end it decreases again.
	//   - "linear" (LinearTiming) - Constant speed.
	//   - "step(n)" (StepTiming(n int) function) - Timing function along stepCount stops along the transition, displaying each stop for equal lengths of time.
	//   - "cubic-bezier(x1, y1, x2, y2)" (CubicBezierTiming(x1, y1, x2, y2 float64) function) - Cubic-Bezier curve timing function. x1 and x2 must be in the range [0, 1].
	TimingFunction PropertyName = "timing-function"

	// IterationCount is the constant for "iteration-count" property tag.
	//
	// Used by AnimationProperty.
	//
	// Sets the number of times an animation sequence should be played before stopping. Used only for animation script.
	//
	// Supported types: int, string.
	//
	// Internal type is int, other types converted to it during assignment.
	IterationCount PropertyName = "iteration-count"

	// AnimationDirection is the constant for "animation-direction" property tag.
	//
	// Used by AnimationProperty.
	//
	// Whether an animation should play forward, backward, or alternate back and forth between playing the sequence forward
	// and backward. Used only for animation script.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (NormalAnimation) or "normal" - The animation plays forward every iteration, that is, when the animation ends, it is immediately reset to its starting position and played again.
	//   - 1 (ReverseAnimation) or "reverse" - The animation plays backwards, from the last position to the first, and then resets to the final position and plays again.
	//   - 2 (AlternateAnimation) or "alternate" - The animation changes direction in each cycle, that is, in the first cycle, it starts from the start position, reaches the end position, and in the second cycle, it continues from the end position and reaches the start position, and so on.
	//   - 3 (AlternateReverseAnimation) or "alternate-reverse" - The animation starts playing from the end position and reaches the start position, and in the next cycle, continuing from the start position, it goes to the end position.
	AnimationDirection PropertyName = "animation-direction"

	// NormalAnimation is value of the "animation-direction" property.
	//
	// The animation plays forwards each cycle. In other words, each time the animation cycles,
	// the animation will reset to the beginning state and start over again. This is the default value.
	NormalAnimation = 0

	// ReverseAnimation is value of the "animation-direction" property.
	//
	// The animation plays backwards each cycle. In other words, each time the animation cycles,
	// the animation will reset to the end state and start over again. Animation steps are performed
	// backwards, and timing functions are also reversed.
	//
	// For example, an "ease-in" timing function becomes "ease-out".
	ReverseAnimation = 1

	// AlternateAnimation is value of the "animation-direction" property.
	//
	// The animation reverses direction each cycle, with the first iteration being played forwards.
	// The count to determine if a cycle is even or odd starts at one.
	AlternateAnimation = 2

	// AlternateReverseAnimation is value of the "animation-direction" property.
	//
	// The animation reverses direction each cycle, with the first iteration being played backwards.
	// The count to determine if a cycle is even or odd starts at one.
	AlternateReverseAnimation = 3

	// EaseTiming - a timing function which increases in velocity towards the middle of the transition, slowing back down at the end
	EaseTiming = "ease"

	// EaseInTiming - a timing function which starts off slowly, with the transition speed increasing until complete
	EaseInTiming = "ease-in"

	// EaseOutTiming - a timing function which starts transitioning quickly, slowing down the transition continues.
	EaseOutTiming = "ease-out"

	// EaseInOutTiming - a timing function which starts transitioning slowly, speeds up, and then slows down again.
	EaseInOutTiming = "ease-in-out"

	// LinearTiming - a timing function at an even speed
	LinearTiming = "linear"
)

// StepsTiming return a timing function along stepCount stops along the transition, displaying each stop for equal lengths of time
func StepsTiming(stepCount int) string {
	return "steps(" + strconv.Itoa(stepCount) + ")"
}

// CubicBezierTiming return a cubic-Bezier curve timing function. x1 and x2 must be in the range [0, 1].
func CubicBezierTiming(x1, y1, x2, y2 float64) string {
	if x1 < 0 {
		x1 = 0
	} else if x1 > 1 {
		x1 = 1
	}
	if x2 < 0 {
		x2 = 0
	} else if x2 > 1 {
		x2 = 1
	}
	return fmt.Sprintf("cubic-bezier(%g, %g, %g, %g)", x1, y1, x2, y2)
}

// AnimatedProperty describes the change script of one property
type AnimatedProperty struct {
	// Tag is the name of the property
	Tag PropertyName
	// From is the initial value of the property
	From any
	// To is the final value of the property
	To any
	// KeyFrames is intermediate property values
	KeyFrames map[int]any
}

type animationData struct {
	dataProperty
	keyFramesName string
	usageCounter  int
	view          View
	listener      func(view View, animation AnimationProperty, event PropertyName)
	oldListeners  map[PropertyName][]func(View, PropertyName)
	oldAnimation  []AnimationProperty
}

// AnimationProperty interface is used to set animation parameters. Used properties:
//
// "property", "id", "duration", "delay", "timing-function", "iteration-count", and "animation-direction"
type AnimationProperty interface {
	Properties
	fmt.Stringer

	// Start starts the animation for the view specified by the first argument.
	// The second argument specifies the animation event listener (can be nil)
	Start(view View, listener func(view View, animation AnimationProperty, event PropertyName)) bool
	// Stop stops the animation
	Stop()
	// Pause pauses the animation
	Pause()
	// Resume resumes an animation that was stopped using the Pause method
	Resume()

	writeTransitionString(tag PropertyName, buffer *strings.Builder)
	animationCSS(session Session) string
	transitionCSS(buffer *strings.Builder, session Session)
	hasAnimatedProperty() bool
	animationName() string
	used()
	unused(session Session)
}

func parseAnimation(obj DataObject) AnimationProperty {
	animation := new(animationData)
	animation.init()

	for i := range obj.PropertyCount() {
		if node := obj.Property(i); node != nil {
			tag := PropertyName(node.Tag())
			if node.Type() == TextNode {
				animation.Set(tag, node.Text())
			} else {
				animation.Set(tag, node)
			}
		}
	}
	return animation
}

// NewAnimationProperty creates a new animation object and return its interface
//
// The following properties can be used:
//   - "id" (ID) - specifies the animation identifier. Used only for animation script.
//   - "duration" (Duration) - specifies the length of time in seconds that an animation takes to complete one cycle;
//   - "delay" (Delay) - specifies the amount of time in seconds to wait from applying the animation to an element before beginning to perform
//     the animation. The animation can start later, immediately from its beginning or immediately and partway through the animation;
//   - "timing-function" (TimingFunction) - specifies how an animation progresses through the duration of each cycle;
//   - "iteration-count" (IterationCount) - specifies the number of times an animation sequence should be played before stopping. Used only for animation script;
//   - "animation-direction" (AnimationDirection) - specifies whether an animation should play forward, backward,
//     or alternate back and forth between playing the sequence forward and backward. Used only for animation script;
//   - "property" (PropertyTag) - describes a scenario for changing a View's property. Used only for animation script.
func NewAnimationProperty(params Params) AnimationProperty {
	animation := new(animationData)
	animation.init()

	for tag, value := range params {
		animation.Set(tag, value)
	}
	return animation
}

// NewTransitionAnimation creates animation data for the transition.
//   - timingFunc - specifies how an animation progresses through the duration of each cycle. If it is "" then "easy" function is used;
//   - duration - specifies the length of time in seconds that an animation takes to complete one cycle. Must be > 0;
//   - delay - specifies the amount of time in seconds to wait from applying the animation to an element before beginning to perform
//     the animation. The animation can start later, immediately from its beginning or immediately and partway through the animation.
func NewTransitionAnimation(timingFunc string, duration float64, delay float64) AnimationProperty {
	animation := new(animationData)
	animation.init()

	if duration <= 0 {
		ErrorLog("Animation duration must be greater than 0")
		return nil
	}

	if !animation.Set(Duration, duration) {
		return nil
	}

	if timingFunc != "" {
		if !animation.Set(TimingFunction, timingFunc) {
			return nil
		}
	}

	if delay != 0 {
		animation.Set(Delay, delay)
	}

	return animation
}

// NewTransitionAnimation creates the animation scenario.
//   - id - specifies the animation identifier.
//   - timingFunc - specifies how an animation progresses through the duration of each cycle. If it is "" then "easy" function is used;
//   - duration - specifies the length of time in seconds that an animation takes to complete one cycle. Must be > 0;
//   - delay - specifies the amount of time in seconds to wait from applying the animation to an element before beginning to perform
//     the animation. The animation can start later, immediately from its beginning or immediately and partway through the animation.
//   - direction - specifies whether an animation should play forward, backward,
//     or alternate back and forth between playing the sequence forward and backward. Only the following values ​​can be used:
//     0 (NormalAnimation), 1 (ReverseAnimation), 2 (AlternateAnimation), and 3 (AlternateReverseAnimation);
//   - iterationCount - specifies the number of times an animation sequence should be played before stopping. A negative value specifies infinite repetition;
//   - property, properties - describes a scenario for changing a View's property.
func NewAnimation(id string, timingFunc string, duration float64, delay float64, direction int, iterationCount int, property AnimatedProperty, properties ...AnimatedProperty) AnimationProperty {
	animation := new(animationData)
	animation.init()

	if duration <= 0 {
		ErrorLog("Animation duration must be greater than 0")
		return nil
	}

	if !animation.Set(Duration, duration) {
		return nil
	}

	if id != "" {
		animation.Set(ID, id)
	}

	if timingFunc != "" {
		animation.Set(TimingFunction, timingFunc)
	}

	if delay != 0 {
		animation.Set(Delay, delay)
	}

	if direction > 0 {
		animation.Set(AnimationDirection, direction)
	}

	if iterationCount != 0 {
		animation.Set(IterationCount, iterationCount)
	}

	if len(properties) > 0 {
		animation.Set(PropertyTag, append([]AnimatedProperty{property}, properties...))
	} else {
		animation.Set(PropertyTag, property)
	}

	return animation
}

func (animation *animationData) init() {
	animation.dataProperty.init()
	animation.normalize = normalizeAnimation
	animation.set = animationSet
	animation.supportedProperties = []PropertyName{ID, PropertyTag, Duration, Delay, TimingFunction, IterationCount, AnimationDirection}
}

func (animation *animationData) animatedProperties() []AnimatedProperty {
	value := animation.getRaw(PropertyTag)
	if value == nil {
		ErrorLog("There are no animated properties.")
		return nil
	}

	props, ok := value.([]AnimatedProperty)
	if !ok {
		ErrorLog("Invalid animated properties.")
		return nil
	}

	if len(props) == 0 {
		ErrorLog("There are no animated properties.")
		return nil
	}

	return props
}

func (animation *animationData) hasAnimatedProperty() bool {
	return animation.animatedProperties() != nil
}

func (animation *animationData) animationName() string {
	return animation.keyFramesName
}

func (animation *animationData) used() {
	animation.usageCounter++
}

func (animation *animationData) unused(session Session) {
	animation.usageCounter--
	if animation.usageCounter <= 0 && animation.keyFramesName != "" {
		session.removeAnimation(animation.keyFramesName)
	}
}

func normalizeAnimation(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	if tag == Direction {
		return AnimationDirection
	}
	return tag
}

func animationSet(properties Properties, tag PropertyName, value any) []PropertyName {
	switch tag {
	case ID:
		if text, ok := value.(string); ok {
			text = strings.Trim(text, " \t\n\r")
			if text == "" {
				properties.setRaw(tag, nil)
			} else {
				properties.setRaw(tag, text)
			}
			return []PropertyName{tag}
		}
		notCompatibleType(tag, value)
		return nil

	case PropertyTag:
		switch value := value.(type) {
		case AnimatedProperty:
			if value.From == nil && value.KeyFrames != nil {
				if val, ok := value.KeyFrames[0]; ok {
					value.From = val
					delete(value.KeyFrames, 0)
				}
			}
			if value.To == nil && value.KeyFrames != nil {
				if val, ok := value.KeyFrames[100]; ok {
					value.To = val
					delete(value.KeyFrames, 100)
				}
			}

			if value.From == nil {
				ErrorLog("AnimatedProperty.From is nil")
			} else if value.To == nil {
				ErrorLog("AnimatedProperty.To is nil")
			} else {
				properties.setRaw(tag, []AnimatedProperty{value})
				return []PropertyName{tag}
			}

		case []AnimatedProperty:
			props := []AnimatedProperty{}
			for _, val := range value {
				if val.From == nil && val.KeyFrames != nil {
					if v, ok := val.KeyFrames[0]; ok {
						val.From = v
						delete(val.KeyFrames, 0)
					}
				}
				if val.To == nil && val.KeyFrames != nil {
					if v, ok := val.KeyFrames[100]; ok {
						val.To = v
						delete(val.KeyFrames, 100)
					}
				}

				if val.From == nil {
					ErrorLog("AnimatedProperty.From is nil")
				} else if val.To == nil {
					ErrorLog("AnimatedProperty.To is nil")
				} else {
					props = append(props, val)
				}
			}
			if len(props) > 0 {
				properties.setRaw(tag, props)
				return []PropertyName{tag}
			} else {
				ErrorLog("[]AnimatedProperty is empty")
			}

		case DataNode:
			parseObject := func(obj DataObject) (AnimatedProperty, bool) {
				result := AnimatedProperty{}
				for i := range obj.PropertyCount() {
					if node := obj.Property(i); node.Type() == TextNode {
						propTag := strings.ToLower(node.Tag())
						switch propTag {
						case "from", "0", "0%":
							result.From = node.Text()

						case "to", "100", "100%":
							result.To = node.Text()

						default:
							tagLen := len(propTag)
							if tagLen > 0 && propTag[tagLen-1] == '%' {
								propTag = propTag[:tagLen-1]
							}
							n, err := strconv.Atoi(propTag)
							if err != nil {
								ErrorLog(err.Error())
							} else if n < 0 || n > 100 {
								ErrorLogF(`key-frame "%d" is out of range`, n)
							} else {
								if result.KeyFrames == nil {
									result.KeyFrames = map[int]any{n: node.Text()}
								} else {
									result.KeyFrames[n] = node.Text()
								}
							}
						}
					}
				}
				if result.From != nil && result.To != nil {
					return result, true
				}
				return result, false
			}

			switch value.Type() {
			case ObjectNode:
				if prop, ok := parseObject(value.Object()); ok {
					properties.setRaw(tag, []AnimatedProperty{prop})
					return []PropertyName{tag}
				}

			case ArrayNode:
				props := []AnimatedProperty{}
				for _, val := range value.ArrayElements() {
					if val.IsObject() {
						if prop, ok := parseObject(val.Object()); ok {
							props = append(props, prop)
						}
					} else {
						notCompatibleType(tag, val)
					}
				}
				if len(props) > 0 {
					properties.setRaw(tag, props)
					return []PropertyName{tag}
				}

			default:
				notCompatibleType(tag, value)
			}

		default:
			notCompatibleType(tag, value)
		}

	case Duration:
		return setFloatProperty(properties, tag, value, 0, math.MaxFloat64)

	case Delay:
		return setFloatProperty(properties, tag, value, -math.MaxFloat64, math.MaxFloat64)

	case TimingFunction:
		if text, ok := value.(string); ok {
			properties.setRaw(tag, text)
			return []PropertyName{tag}
		}

	case IterationCount:
		return setIntProperty(properties, tag, value)

	case AnimationDirection:
		return setEnumProperty(properties, AnimationDirection, value, enumProperties[AnimationDirection].values)

	default:
		ErrorLogF(`The "%s" property is not supported by Animation`, tag)
	}

	return nil
}

func (animation *animationData) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString("animation {")

	for _, tag := range animation.AllTags() {
		if tag != PropertyTag {
			if value, ok := animation.properties[tag]; ok && value != nil {
				buffer.WriteString("\n\t")
				buffer.WriteString(string(tag))
				buffer.WriteString(" = ")
				writePropertyValue(buffer, tag, value, "\t")
				buffer.WriteRune(',')
			}
		}
	}

	writeProperty := func(prop AnimatedProperty, indent string) {
		buffer.WriteString(string(prop.Tag))
		buffer.WriteString("{\n")
		buffer.WriteString(indent)
		buffer.WriteString("from = ")
		writePropertyValue(buffer, "from", prop.From, indent)
		buffer.WriteString(",\n")
		buffer.WriteString(indent)
		buffer.WriteString("to = ")
		writePropertyValue(buffer, "to", prop.To, indent)
		for key, value := range prop.KeyFrames {
			buffer.WriteString(",\n")
			buffer.WriteString(indent)
			tag := strconv.Itoa(key) + "%"
			buffer.WriteString(tag)
			buffer.WriteString(" = ")
			writePropertyValue(buffer, PropertyName(tag), value, indent)
		}
		buffer.WriteString("\n")
		buffer.WriteString(indent[1:])
		buffer.WriteString("}")
	}

	if props := animation.animatedProperties(); len(props) > 0 {

		buffer.WriteString("\n\t")
		buffer.WriteString(string(PropertyTag))
		buffer.WriteString(" = ")
		if len(props) > 1 {
			buffer.WriteString("[\n")
			for _, prop := range props {
				buffer.WriteString("\t\t")
				writeProperty(prop, "\t\t\t")
				buffer.WriteString(",\n")
			}
			buffer.WriteString("\t]")
		} else {
			writeProperty(props[0], "\t\t")
		}
	}

	buffer.WriteString("\n}")
	return buffer.String()
}

func (animation *animationData) animationCSS(session Session) string {
	if animation.keyFramesName == "" {
		if props := animation.animatedProperties(); props != nil {
			animation.keyFramesName = session.registerAnimation(props)
		} else {
			return ""
		}
	}

	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	buffer.WriteString(animation.keyFramesName)

	if duration, ok := floatProperty(animation, Duration, session, 1); ok && duration > 0 {
		buffer.WriteString(fmt.Sprintf(" %gs ", duration))
	} else {
		buffer.WriteString(" 1s ")
	}

	buffer.WriteString(timingFunctionCSS(animation, TimingFunction, session))

	if delay, ok := floatProperty(animation, Delay, session, 0); ok && delay > 0 {
		buffer.WriteString(fmt.Sprintf(" %gs", delay))
	} else {
		buffer.WriteString(" 0s")
	}

	if iterationCount, _ := intProperty(animation, IterationCount, session, 0); iterationCount >= 0 {
		if iterationCount == 0 {
			iterationCount = 1
		}
		buffer.WriteString(fmt.Sprintf(" %d ", iterationCount))
	} else {
		buffer.WriteString(" infinite ")
	}

	direction, _ := enumProperty(animation, AnimationDirection, session, 0)
	values := enumProperties[AnimationDirection].cssValues
	if direction < 0 || direction >= len(values) {
		direction = 0
	}
	buffer.WriteString(values[direction])

	// TODO "animation-fill-mode"
	buffer.WriteString(" forwards")

	return buffer.String()
}

func (animation *animationData) transitionCSS(buffer *strings.Builder, session Session) {

	if duration, ok := floatProperty(animation, Duration, session, 1); ok && duration > 0 {
		buffer.WriteString(fmt.Sprintf(" %gs ", duration))
	} else {
		buffer.WriteString(" 1s ")
	}

	buffer.WriteString(timingFunctionCSS(animation, TimingFunction, session))

	if delay, ok := floatProperty(animation, Delay, session, 0); ok && delay > 0 {
		buffer.WriteString(fmt.Sprintf(" %gs", delay))
	}
}

func (animation *animationData) writeTransitionString(tag PropertyName, buffer *strings.Builder) {
	buffer.WriteString(string(tag))
	buffer.WriteString("{")
	lead := " "

	writeFloatProperty := func(name PropertyName) bool {
		if value := animation.getRaw(name); value != nil {
			buffer.WriteString(lead)
			buffer.WriteString(string(name))
			buffer.WriteString(" = ")
			writePropertyValue(buffer, name, value, "")
			lead = ", "
			return true
		}
		return false
	}

	if !writeFloatProperty(Duration) {
		buffer.WriteString(" duration = 1")
		lead = ", "
	}

	writeFloatProperty(Delay)

	if value := animation.getRaw(TimingFunction); value != nil {
		if timingFunction, ok := value.(string); ok && timingFunction != "" {
			buffer.WriteString(lead)
			buffer.WriteString(string(TimingFunction))
			buffer.WriteString(" = ")
			if strings.ContainsAny(timingFunction, " ,()") {
				buffer.WriteRune('"')
				buffer.WriteString(timingFunction)
				buffer.WriteRune('"')
			} else {
				buffer.WriteString(timingFunction)
			}
		}
	}

	buffer.WriteString(" }")
}

func timingFunctionCSS(properties Properties, tag PropertyName, session Session) string {
	if timingFunction, ok := stringProperty(properties, tag, session); ok {
		if timingFunction, ok = session.resolveConstants(timingFunction); ok && isTimingFunctionValid(timingFunction) {
			return timingFunction
		}
	}
	return ("ease")
}

func isTimingFunctionValid(timingFunction string) bool {
	switch timingFunction {
	case "", EaseTiming, EaseInTiming, EaseOutTiming, EaseInOutTiming, LinearTiming:
		return true
	}

	size := len(timingFunction)
	if size > 0 && timingFunction[size-1] == ')' {
		if index := strings.IndexRune(timingFunction, '('); index > 0 {
			args := timingFunction[index+1 : size-1]
			switch timingFunction[:index] {
			case "steps":
				if _, err := strconv.Atoi(strings.Trim(args, " \t\n")); err == nil {
					return true
				}

			case "cubic-bezier":
				if params := strings.Split(args, ","); len(params) == 4 {
					for _, param := range params {
						if _, err := strconv.ParseFloat(strings.Trim(param, " \t\n"), 64); err != nil {
							return false
						}
					}
					return true
				}
			}
		}
	}

	return false
}

// IsTimingFunctionValid returns "true" if the "timingFunction" argument is the valid timing function.
func IsTimingFunctionValid(timingFunction string, session Session) bool {
	if timingFunc, ok := session.resolveConstants(strings.Trim(timingFunction, " \t\n")); ok {
		return isTimingFunctionValid(timingFunc)
	}
	return false
}

func (session *sessionData) registerAnimation(props []AnimatedProperty) string {

	session.animationCounter++
	name := fmt.Sprintf("kf%06d", session.animationCounter)

	var cssBuilder cssStyleBuilder

	cssBuilder.init(0)
	cssBuilder.startAnimation(name)

	fromParams := Params{}
	toParams := Params{}
	frames := []int{}

	for _, prop := range props {
		fromParams[prop.Tag] = prop.From
		toParams[prop.Tag] = prop.To
		if len(prop.KeyFrames) > 0 {
			for frame := range prop.KeyFrames {
				needAppend := true
				for i, n := range frames {
					if n == frame {
						needAppend = false
						break
					} else if frame < n {
						needAppend = false
						frames = append(append(frames[:i], frame), frames[i+1:]...)
						break
					}
				}
				if needAppend {
					frames = append(frames, frame)
				}
			}
		}
	}

	cssBuilder.startAnimationFrame("from")
	NewViewStyle(fromParams).cssViewStyle(&cssBuilder, session)
	cssBuilder.endAnimationFrame()

	if len(frames) > 0 {
		for _, frame := range frames {
			params := Params{}
			for _, prop := range props {
				if prop.KeyFrames != nil {
					if value, ok := prop.KeyFrames[frame]; ok {
						params[prop.Tag] = value
					}
				}
			}

			if len(params) > 0 {
				cssBuilder.startAnimationFrame(strconv.Itoa(frame) + "%")
				NewViewStyle(params).cssViewStyle(&cssBuilder, session)
				cssBuilder.endAnimationFrame()
			}
		}
	}

	cssBuilder.startAnimationFrame("to")
	NewViewStyle(toParams).cssViewStyle(&cssBuilder, session)
	cssBuilder.endAnimationFrame()

	cssBuilder.endAnimation()
	session.addAnimationCSS(cssBuilder.finish())

	return name
}

func (view *viewData) SetAnimated(tag PropertyName, value any, animation AnimationProperty) bool {
	if animation == nil {
		return view.Set(tag, value)
	}

	session := view.Session()
	htmlID := view.htmlID()
	session.startUpdateScript(htmlID)

	session.updateProperty(htmlID, "ontransitionend", "transitionEndEvent(this, event)")
	session.updateProperty(htmlID, "ontransitioncancel", "transitionCancelEvent(this, event)")

	transitions := getTransitionProperty(view)
	var prevAnimation AnimationProperty = nil
	if transitions != nil {
		if prev, ok := transitions[tag]; ok {
			prevAnimation = prev
		}
	}
	view.singleTransition[tag] = prevAnimation
	setTransition(view, tag, animation)
	view.session.updateCSSProperty(view.htmlID(), "transition", transitionCSS(view, view.session))

	session.finishUpdateScript(htmlID)

	result := view.Set(tag, value)
	if !result {
		delete(view.singleTransition, tag)
		view.session.updateCSSProperty(view.htmlID(), "transition", transitionCSS(view, view.session))
	}

	return result
}

func animationCSS(properties Properties, session Session) string {
	if value := properties.getRaw(Animation); value != nil {
		if animations, ok := value.([]AnimationProperty); ok {
			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)

			for _, animation := range animations {
				if css := animation.animationCSS(session); css != "" {
					if buffer.Len() > 0 {
						buffer.WriteString(", ")
					}
					buffer.WriteString(css)
				}
			}

			return buffer.String()
		}
	}

	return ""
}

func transitionCSS(properties Properties, session Session) string {
	if transitions := getTransitionProperty(properties); len(transitions) > 0 {
		buffer := allocStringBuilder()
		defer freeStringBuilder(buffer)

		convert := map[PropertyName]string{
			CellHeight:        "grid-template-rows",
			CellWidth:         "grid-template-columns",
			Row:               "grid-row",
			Column:            "grid-column",
			Clip:              "clip-path",
			Shadow:            "box-shadow",
			ColumnSeparator:   "column-rule",
			FontName:          "font",
			TextSize:          "font-size",
			TextLineThickness: "text-decoration-thickness",
		}

		for tag, animation := range transitions {
			if buffer.Len() > 0 {
				buffer.WriteString(", ")
			}

			if cssTag, ok := convert[tag]; ok {
				buffer.WriteString(cssTag)
			} else {
				buffer.WriteString(string(tag))
			}
			animation.transitionCSS(buffer, session)
		}
		return buffer.String()
	}
	return ""
}

/*
func (view *viewData) updateTransitionCSS() {
	view.session.updateCSSProperty(view.htmlID(), "transition", transitionCSS(view, view.session))
}
*/

func (style *viewStyle) Transition(tag PropertyName) AnimationProperty {
	if transitions := getTransitionProperty(style); transitions != nil {
		if anim, ok := transitions[tag]; ok {
			return anim
		}
	}
	return nil
}

func (style *viewStyle) Transitions() map[PropertyName]AnimationProperty {
	result := map[PropertyName]AnimationProperty{}
	for tag, animation := range getTransitionProperty(style) {
		result[tag] = animation
	}
	return result
}

func (style *viewStyle) SetTransition(tag PropertyName, animation AnimationProperty) {
	setTransition(style, style.normalize(tag), animation)
}

func (view *viewData) SetTransition(tag PropertyName, animation AnimationProperty) {
	setTransition(view, view.normalize(tag), animation)
	if view.created {
		view.session.updateCSSProperty(view.htmlID(), "transition", transitionCSS(view, view.session))
	}
}

func setTransition(properties Properties, tag PropertyName, animation AnimationProperty) {
	transitions := getTransitionProperty(properties)

	if animation == nil {
		if transitions != nil {
			delete(transitions, tag)
			if len(transitions) == 0 {
				properties.setRaw(Transition, nil)
			}
		}
	} else if transitions != nil {
		transitions[tag] = animation
	} else {
		properties.setRaw(Transition, map[PropertyName]AnimationProperty{tag: animation})
	}
}

func getTransitionProperty(properties Properties) map[PropertyName]AnimationProperty {
	if value := properties.getRaw(Transition); value != nil {
		if transitions, ok := value.(map[PropertyName]AnimationProperty); ok {
			return transitions
		}
	}
	return nil
}

func setAnimationProperty(properties Properties, tag PropertyName, value any) bool {

	set := func(animations []AnimationProperty) {
		properties.setRaw(tag, animations)
		for _, animation := range animations {
			animation.used()
		}
	}

	switch value := value.(type) {
	case AnimationProperty:
		set([]AnimationProperty{value})
		return true

	case []AnimationProperty:
		set(value)
		return true

	case DataObject:
		if animation := parseAnimation(value); animation.hasAnimatedProperty() {
			set([]AnimationProperty{animation})
			return true
		}

	case DataNode:
		animations := []AnimationProperty{}
		result := true
		for i := range value.ArraySize() {
			if obj := value.ArrayElement(i).Object(); obj != nil {
				if anim := parseAnimation(obj); anim.hasAnimatedProperty() {
					animations = append(animations, anim)
				} else {
					result = false
				}
			} else {
				notCompatibleType(tag, value.ArrayElement(i))
				result = false
			}
		}
		if result && len(animations) > 0 {
			set(animations)
		}
		return result
	}

	notCompatibleType(tag, value)
	return false
}

// SetAnimated sets the property with name "tag" of the "rootView" subview with "viewID" id by value. Result:
// true - success,
// false - error (incompatible type or invalid format of a string value, see AppLog).
func SetAnimated(rootView View, viewID string, tag PropertyName, value any, animation AnimationProperty) bool {
	if view := ViewByID(rootView, viewID); view != nil {
		return view.SetAnimated(tag, value, animation)
	}
	return false
}

// IsAnimationPaused returns "true" if an animation of the subview is paused, "false" otherwise.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func IsAnimationPaused(view View, subviewID ...string) bool {
	return boolStyledProperty(view, subviewID, AnimationPaused, false)
}

// GetTransitions returns the subview transitions. The result is always non-nil.
// If the second argument (subviewID) is not specified or it is "" then transitions of the first argument (view) is returned
func GetTransitions(view View, subviewID ...string) map[PropertyName]AnimationProperty {
	if view = getSubview(view, subviewID); view != nil {
		return view.Transitions()
	}

	return map[PropertyName]AnimationProperty{}
}

// GetTransition returns the subview property transition. If there is no transition for the given property then nil is returned.
// If the second argument (subviewID) is not specified or it is "" then transitions of the first argument (view) is returned
func GetTransition(view View, subviewID string, tag PropertyName) AnimationProperty {
	if subviewID != "" {
		view = ViewByID(view, subviewID)
	}

	if view != nil {
		return view.Transition(tag)
	}

	return nil
}

// AddTransition adds the transition for the subview property.
// If the second argument (subviewID) is not specified or it is "" then the transition is added to the first argument (view)
func AddTransition(view View, subviewID string, tag PropertyName, animation AnimationProperty) bool {
	if tag != "" {
		if subviewID != "" {
			view = ViewByID(view, subviewID)
		}

		if view != nil {
			view.SetTransition(tag, animation)
			return true
		}
	}
	return false
}

// GetAnimation returns the subview animations. The result is always non-nil.
// If the second argument (subviewID) is not specified or it is "" then transitions of the first argument (view) is returned
func GetAnimation(view View, subviewID ...string) []AnimationProperty {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.getRaw(Animation); value != nil {
			if animations, ok := value.([]AnimationProperty); ok && animations != nil {
				return animations
			}
		}
	}

	return []AnimationProperty{}
}
