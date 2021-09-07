package rui

import (
	"fmt"
	"strconv"
	"strings"
)

const (
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

// StepsTiming return a timing function along stepCount stops along the transition, diplaying each stop for equal lengths of time
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

// AnimationFinishedListener describes the end of an animation event handler
type AnimationFinishedListener interface {
	// OnAnimationFinished is called when a property animation is finished
	OnAnimationFinished(view View, property string)
}

type Animation struct {
	// Duration defines the time in seconds an animation should take to complete
	Duration float64
	// TimingFunction defines how intermediate values are calculated for a property being affected
	// by an animation effect. If the value is "" then the "ease" function is used
	TimingFunction string
	// Delay defines the duration in seconds to wait before starting a property's animation.
	Delay float64
	// FinishListener defines the end of an animation event handler
	FinishListener AnimationFinishedListener
}

type animationFinishedFunc struct {
	finishFunc func(View, string)
}

func (listener animationFinishedFunc) OnAnimationFinished(view View, property string) {
	if listener.finishFunc != nil {
		listener.finishFunc(view, property)
	}
}

func AnimationFinishedFunc(finishFunc func(View, string)) AnimationFinishedListener {
	listener := new(animationFinishedFunc)
	listener.finishFunc = finishFunc
	return listener
}

func validateTimingFunction(timingFunction string) bool {
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

func (view *viewData) SetAnimated(tag string, value interface{}, animation Animation) bool {
	timingFunction, ok := view.session.resolveConstants(animation.TimingFunction)
	if !ok || animation.Duration <= 0 || !validateTimingFunction(timingFunction) {
		if view.Set(tag, value) {
			if animation.FinishListener != nil {
				animation.FinishListener.OnAnimationFinished(view, tag)
			}
			return true
		}
		return false
	}

	updateProperty(view.htmlID(), "ontransitionend", "transitionEndEvent(this, event)", view.session)
	updateProperty(view.htmlID(), "ontransitioncancel", "transitionCancelEvent(this, event)", view.session)
	animation.TimingFunction = timingFunction
	view.animation[tag] = animation
	view.updateTransitionCSS()

	result := view.Set(tag, value)
	if !result {
		delete(view.animation, tag)
		view.updateTransitionCSS()
	}

	return result
}

func (view *viewData) transitionCSS() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)
	for tag, animation := range view.animation {
		if buffer.Len() > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(tag)
		buffer.WriteString(fmt.Sprintf(" %gs", animation.Duration))
		if animation.TimingFunction != "" {
			buffer.WriteRune(' ')
			buffer.WriteString(animation.TimingFunction)
		}
		if animation.Delay > 0 {
			if animation.TimingFunction == "" {
				buffer.WriteString(" ease")
			}
			buffer.WriteString(fmt.Sprintf(" %gs", animation.Delay))
		}
	}
	return buffer.String()
}

func (view *viewData) updateTransitionCSS() {
	updateCSSProperty(view.htmlID(), "transition", view.transitionCSS(), view.Session())
}

// SetAnimated sets the property with name "tag" of the "rootView" subview with "viewID" id by value. Result:
//  true - success,
//  false - error (incompatible type or invalid format of a string value, see AppLog).
func SetAnimated(rootView View, viewID, tag string, value interface{}, animation Animation) bool {
	if view := ViewByID(rootView, viewID); view != nil {
		return view.SetAnimated(tag, value, animation)
	}
	return false
}
