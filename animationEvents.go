package rui

import "strings"

const (
	// TransitionRunEvent is the constant for "transition-run-event" property tag.
	// The "transition-run-event" is fired when a transition is first created,
	// i.e. before any transition delay has begun.
	TransitionRunEvent = "transition-run-event"

	// TransitionStartEvent is the constant for "transition-end-event" property tag.
	// The "transition-start-event" is fired when a transition has actually started,
	// i.e., after "delay" has ended.
	TransitionStartEvent = "transition-start-event"

	// TransitionEndEvent is the constant for "transition-end-event" property tag.
	// The "transition-end-event" is fired when a transition has completed.
	TransitionEndEvent = "transition-end-event"

	// TransitionCancelEvent is the constant for "transition-cancel-event" property tag.
	// The "transition-cancel-event" is fired when a transition is cancelled. The transition is cancelled when:
	// * A new property transition has begun.
	// * The "visibility" property is set to "gone".
	// * The transition is stopped before it has run to completion, e.g. by moving the mouse off a hover-transitioning view.
	TransitionCancelEvent = "transition-cancel-event"

	// AnimationStartEvent is the constant for "animation-start-event" property tag.
	// The "animation-start-event" is fired when an animation has started.
	// If there is an animation-delay, this event will fire once the delay period has expired.
	AnimationStartEvent = "animation-start-event"

	// AnimationEndEvent is the constant for "animation-end-event" property tag.
	// The "animation-end-event" is fired when aт фnimation has completed.
	// If the animation aborts before reaching completion, such as if the element is removed
	// or the animation is removed from the element, the "animation-end-event" is not fired.
	AnimationEndEvent = "animation-end-event"

	// AnimationCancelEvent is the constant for "animation-cancel-event" property tag.
	// The "animation-cancel-event" is fired when an animation unexpectedly aborts.
	// In other words, any time it stops running without sending the "animation-end-event".
	// This might happen when the animation-name is changed such that the animation is removed,
	// or when the animating view is hidden. Therefore, either directly or because any of its
	// containing views are hidden.
	// The event is not supported by all browsers.
	AnimationCancelEvent = "animation-cancel-event"

	// AnimationIterationEvent is the constant for "animation-iteration-event" property tag.
	// The "animation-iteration-event" is fired when an iteration of an animation ends,
	// and another one begins. This event does not occur at the same time as the animationend event,
	// and therefore does not occur for animations with an "iteration-count" of one.
	AnimationIterationEvent = "animation-iteration-event"
)

var transitionEvents = map[string]struct{ jsEvent, jsFunc string }{
	TransitionRunEvent:    {jsEvent: "ontransitionrun", jsFunc: "transitionRunEvent"},
	TransitionStartEvent:  {jsEvent: "ontransitionstart", jsFunc: "transitionStartEvent"},
	TransitionEndEvent:    {jsEvent: "ontransitionend", jsFunc: "transitionEndEvent"},
	TransitionCancelEvent: {jsEvent: "ontransitioncancel", jsFunc: "transitionCancelEvent"},
}

func (view *viewData) setTransitionListener(tag string, value any) bool {
	listeners, ok := valueToEventListeners[View, string](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeTransitionListener(tag)
	} else if js, ok := transitionEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			view.session.updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)")
		}
	} else {
		return false
	}
	return true
}

func (view *viewData) removeTransitionListener(tag string) {
	delete(view.properties, tag)
	if view.created {
		if js, ok := transitionEvents[tag]; ok {
			view.session.removeProperty(view.htmlID(), js.jsEvent)
		}
	}
}

func transitionEventsHtml(view View, buffer *strings.Builder) {
	for tag, js := range transitionEvents {
		if value := view.getRaw(tag); value != nil {
			if listeners, ok := value.([]func(View, string)); ok && len(listeners) > 0 {
				buffer.WriteString(js.jsEvent + `="` + js.jsFunc + `(this, event)" `)
			}
		}
	}
}

func (view *viewData) handleTransitionEvents(tag string, data DataObject) {
	if property, ok := data.PropertyValue("property"); ok {
		if tag == TransitionEndEvent || tag == TransitionCancelEvent {
			if animation, ok := view.singleTransition[property]; ok {
				delete(view.singleTransition, property)
				if animation != nil {
					view.transitions[property] = animation
				} else {
					delete(view.transitions, property)
				}
				view.updateTransitionCSS()
			}
		}

		for _, listener := range getEventListeners[View, string](view, nil, tag) {
			listener(view, property)
		}
	}
}

var animationEvents = map[string]struct{ jsEvent, jsFunc string }{
	AnimationStartEvent:     {jsEvent: "onanimationstart", jsFunc: "animationStartEvent"},
	AnimationEndEvent:       {jsEvent: "onanimationend", jsFunc: "animationEndEvent"},
	AnimationIterationEvent: {jsEvent: "onanimationiteration", jsFunc: "animationIterationEvent"},
	AnimationCancelEvent:    {jsEvent: "onanimationcancel", jsFunc: "animationCancelEvent"},
}

func (view *viewData) setAnimationListener(tag string, value any) bool {
	listeners, ok := valueToEventListeners[View, string](value)
	if !ok {
		notCompatibleType(tag, value)
		return false
	}

	if listeners == nil {
		view.removeAnimationListener(tag)
	} else if js, ok := animationEvents[tag]; ok {
		view.properties[tag] = listeners
		if view.created {
			view.session.updateProperty(view.htmlID(), js.jsEvent, js.jsFunc+"(this, event)")
		}
	} else {
		return false
	}
	return true
}

func (view *viewData) removeAnimationListener(tag string) {
	delete(view.properties, tag)
	if view.created {
		if js, ok := animationEvents[tag]; ok {
			view.session.removeProperty(view.htmlID(), js.jsEvent)
		}
	}
}

func animationEventsHtml(view View, buffer *strings.Builder) {
	for tag, js := range animationEvents {
		if value := view.getRaw(tag); value != nil {
			if listeners, ok := value.([]func(View)); ok && len(listeners) > 0 {
				buffer.WriteString(js.jsEvent + `="` + js.jsFunc + `(this, event)" `)
			}
		}
	}
}

func (view *viewData) handleAnimationEvents(tag string, data DataObject) {
	if listeners := getEventListeners[View, string](view, nil, tag); len(listeners) > 0 {
		id := ""
		if name, ok := data.PropertyValue("name"); ok {
			for _, animation := range GetAnimation(view) {
				if name == animation.animationName() {
					id, _ = stringProperty(animation, ID, view.Session())
				}
			}
		}
		for _, listener := range listeners {
			listener(view, id)
		}
	}
}

// GetTransitionRunListeners returns the "transition-run-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTransitionRunListeners(view View, subviewID ...string) []func(View, string) {
	return getEventListeners[View, string](view, subviewID, TransitionRunEvent)
}

// GetTransitionStartListeners returns the "transition-start-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTransitionStartListeners(view View, subviewID ...string) []func(View, string) {
	return getEventListeners[View, string](view, subviewID, TransitionStartEvent)
}

// GetTransitionEndListeners returns the "transition-end-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTransitionEndListeners(view View, subviewID ...string) []func(View, string) {
	return getEventListeners[View, string](view, subviewID, TransitionEndEvent)
}

// GetTransitionCancelListeners returns the "transition-cancel-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTransitionCancelListeners(view View, subviewID ...string) []func(View, string) {
	return getEventListeners[View, string](view, subviewID, TransitionCancelEvent)
}

// GetAnimationStartListeners returns the "animation-start-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAnimationStartListeners(view View, subviewID ...string) []func(View, string) {
	return getEventListeners[View, string](view, subviewID, AnimationStartEvent)
}

// GetAnimationEndListeners returns the "animation-end-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAnimationEndListeners(view View, subviewID ...string) []func(View, string) {
	return getEventListeners[View, string](view, subviewID, AnimationEndEvent)
}

// GetAnimationCancelListeners returns the "animation-cancel-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAnimationCancelListeners(view View, subviewID ...string) []func(View, string) {
	return getEventListeners[View, string](view, subviewID, AnimationCancelEvent)
}

// GetAnimationIterationListeners returns the "animation-iteration-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAnimationIterationListeners(view View, subviewID ...string) []func(View, string) {
	return getEventListeners[View, string](view, subviewID, AnimationIterationEvent)
}
