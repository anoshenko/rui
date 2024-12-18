package rui

// Constants which describe values for view's animation events properties
const (
	// TransitionRunEvent is the constant for "transition-run-event" property tag.
	//
	// Used by View:
	// Is fired when a transition is first created, i.e. before any transition delay has begun.
	//
	// General listener format:
	//  func(view rui.View, propertyName rui.PropertyName).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - propertyName - Name of the property.
	//
	// Allowed listener formats:
	//  func(view rui.View),
	//  func(propertyName rui.PropertyName)
	//  func().
	TransitionRunEvent PropertyName = "transition-run-event"

	// TransitionStartEvent is the constant for "transition-start-event" property tag.
	//
	// Used by View:
	// Is fired when a transition has actually started, i.e., after "delay" has ended.
	//
	// General listener format:
	//  func(view rui.View, propertyName rui.PropertyName).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - propertyName - Name of the property.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(propertyName rui.PropertyName)
	//  func()
	TransitionStartEvent PropertyName = "transition-start-event"

	// TransitionEndEvent is the constant for "transition-end-event" property tag.
	//
	// Used by View:
	// Is fired when a transition has completed.
	//
	// General listener format:
	//  func(view rui.View, propertyName rui.PropertyName).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - propertyName - Name of the property.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(propertyName rui.PropertyName)
	//  func()
	TransitionEndEvent PropertyName = "transition-end-event"

	// TransitionCancelEvent is the constant for "transition-cancel-event" property tag.
	//
	// Used by View:
	// Is fired when a transition is cancelled. The transition is cancelled when:
	//   - A new property transition has begun.
	//   - The "visibility" property is set to "gone".
	//   - The transition is stopped before it has run to completion, e.g. by moving the mouse off a hover-transitioning view.
	//
	// General listener format:
	//  func(view rui.View, propertyName rui.PropertyName).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - propertyName - Name of the property.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(propertyName rui.PropertyName)
	//  func()
	TransitionCancelEvent PropertyName = "transition-cancel-event"

	// AnimationStartEvent is the constant for "animation-start-event" property tag.
	//
	// Used by View:
	// Fired when an animation has started. If there is an "animation-delay", this event will fire once the delay period has
	// expired.
	//
	// General listener format:
	//  func(view rui.View, animationId string).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - animationId - Id of the animation.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(animationId string)
	//  func()
	AnimationStartEvent PropertyName = "animation-start-event"

	// AnimationEndEvent is the constant for "animation-end-event" property tag.
	//
	// Used by View:
	// Fired when an animation has completed. If the animation aborts before reaching completion, such as if the element is
	// removed or the animation is removed from the element, the "animation-end-event" is not fired.
	//
	// General listener format:
	//  func(view rui.View, animationId string).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - animationId - Id of the animation.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(animationId string)
	//  func()
	AnimationEndEvent PropertyName = "animation-end-event"

	// AnimationCancelEvent is the constant for "animation-cancel-event" property tag.
	//
	// Used by View:
	// Fired when an animation unexpectedly aborts. In other words, any time it stops running without sending the
	// "animation-end-event". This might happen when the animation-name is changed such that the animation is removed, or when
	// the animating view is hidden. Therefore, either directly or because any of its containing views are hidden. The event
	// is not supported by all browsers.
	//
	// General listener format:
	//  func(view rui.View, animationId string).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - animationId - Id of the animation.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(animationId string)
	//  func()
	AnimationCancelEvent PropertyName = "animation-cancel-event"

	// AnimationIterationEvent is the constant for "animation-iteration-event" property tag.
	//
	// Used by View:
	// Fired when an iteration of an animation ends, and another one begins. This event does not occur at the same time as the
	// animation end event, and therefore does not occur for animations with an "iteration-count" of one.
	//
	// General listener format:
	//  func(view rui.View, animationId string).
	//
	// where:
	//   - view - Interface of a view which generated this event,
	//   - animationId - Id of the animation.
	//
	// Allowed listener formats:
	//  func(view rui.View)
	//  func(animationId string)
	//  func()
	AnimationIterationEvent PropertyName = "animation-iteration-event"
)

/*
func setTransitionListener(properties Properties, tag PropertyName, value any) bool {
	if listeners, ok := valueToOneArgEventListeners[View, string](value); ok {
		if len(listeners) == 0 {
			properties.setRaw(tag, nil)
		} else {
			properties.setRaw(tag, listeners)
		}
		return true
	}
	notCompatibleType(tag, value)
	return false
}

func (view *viewData) removeTransitionListener(tag PropertyName) {
	delete(view.properties, tag)
	if view.created {
		if js, ok := eventJsFunc[tag]; ok {
			view.session.removeProperty(view.htmlID(), js.jsEvent)
		}
	}
}

func transitionEventsHtml(view View, buffer *strings.Builder) {
	for _, tag := range []PropertyName{TransitionRunEvent, TransitionStartEvent, TransitionEndEvent, TransitionCancelEvent} {
		if value := view.getRaw(tag); value != nil {
			if js, ok := eventJsFunc[tag]; ok {
				if listeners, ok := value.([]func(View, string)); ok && len(listeners) > 0 {
					buffer.WriteString(js.jsEvent)
					buffer.WriteString(`="`)
					buffer.WriteString(js.jsFunc)
					buffer.WriteString(`(this, event)" `)
				}
			}
		}
	}
}
*/

func (view *viewData) handleTransitionEvents(tag PropertyName, data DataObject) {
	if propertyName, ok := data.PropertyValue("property"); ok {
		property := PropertyName(propertyName)
		if tag == TransitionEndEvent || tag == TransitionCancelEvent {
			if animation, ok := view.singleTransition[property]; ok {
				delete(view.singleTransition, property)
				setTransition(view, property, animation)
				session := view.session
				session.updateCSSProperty(view.htmlID(), "transition", transitionCSS(view, session))
			}
		}

		for _, listener := range getOneArgEventListeners[View, PropertyName](view, nil, tag) {
			listener(view, property)
		}
	}
}

/*
	func setAnimationListener(properties Properties, tag PropertyName, value any) bool {
		if listeners, ok := valueToOneArgEventListeners[View, string](value); ok {
			if len(listeners) == 0 {
				properties.setRaw(tag, nil)
			} else {
				properties.setRaw(tag, listeners)
			}
			return true
		}
		notCompatibleType(tag, value)
		return false
	}

func (view *viewData) removeAnimationListener(tag PropertyName) {
	delete(view.properties, tag)
	if view.created {
		if js, ok := eventJsFunc[tag]; ok {
			view.session.removeProperty(view.htmlID(), js.jsEvent)
		}
	}
}

func animationEventsHtml(view View, buffer *strings.Builder) {
	for _, tag := range []PropertyName{AnimationStartEvent, AnimationEndEvent, AnimationIterationEvent, AnimationCancelEvent} {
		if value := view.getRaw(tag); value != nil {
			if js, ok := eventJsFunc[tag]; ok {
				if listeners, ok := value.([]func(View, string)); ok && len(listeners) > 0 {
					buffer.WriteString(js.jsEvent)
					buffer.WriteString(`="`)
					buffer.WriteString(js.jsFunc)
					buffer.WriteString(`(this, event)" `)
				}
			}
		}
	}
}
*/

func (view *viewData) handleAnimationEvents(tag PropertyName, data DataObject) {
	if listeners := getOneArgEventListeners[View, string](view, nil, tag); len(listeners) > 0 {
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
	return getOneArgEventListeners[View, string](view, subviewID, TransitionRunEvent)
}

// GetTransitionStartListeners returns the "transition-start-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTransitionStartListeners(view View, subviewID ...string) []func(View, string) {
	return getOneArgEventListeners[View, string](view, subviewID, TransitionStartEvent)
}

// GetTransitionEndListeners returns the "transition-end-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTransitionEndListeners(view View, subviewID ...string) []func(View, string) {
	return getOneArgEventListeners[View, string](view, subviewID, TransitionEndEvent)
}

// GetTransitionCancelListeners returns the "transition-cancel-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetTransitionCancelListeners(view View, subviewID ...string) []func(View, string) {
	return getOneArgEventListeners[View, string](view, subviewID, TransitionCancelEvent)
}

// GetAnimationStartListeners returns the "animation-start-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAnimationStartListeners(view View, subviewID ...string) []func(View, string) {
	return getOneArgEventListeners[View, string](view, subviewID, AnimationStartEvent)
}

// GetAnimationEndListeners returns the "animation-end-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAnimationEndListeners(view View, subviewID ...string) []func(View, string) {
	return getOneArgEventListeners[View, string](view, subviewID, AnimationEndEvent)
}

// GetAnimationCancelListeners returns the "animation-cancel-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAnimationCancelListeners(view View, subviewID ...string) []func(View, string) {
	return getOneArgEventListeners[View, string](view, subviewID, AnimationCancelEvent)
}

// GetAnimationIterationListeners returns the "animation-iteration-event" listener list.
// If there are no listeners then the empty list is returned.
// If the second argument (subviewID) is not specified or it is "" then a value from the first argument (view) is returned.
func GetAnimationIterationListeners(view View, subviewID ...string) []func(View, string) {
	return getOneArgEventListeners[View, string](view, subviewID, AnimationIterationEvent)
}
