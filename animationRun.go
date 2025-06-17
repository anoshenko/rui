package rui

import "slices"

func (animation *animationData) Start(view View, listener func(view View, animation AnimationProperty, event PropertyName)) bool {
	if view == nil {
		ErrorLog("nil View in animation.Start() function")
		return false
	}
	if !animation.hasAnimatedProperty() {
		return false
	}

	animation.view = view
	animation.listener = listener

	animation.oldAnimation = nil

	//if  getOneArgEventListeners[View, PropertyName](view, nil, Animation)
	if value := view.Get(Animation); value != nil {
		if oldAnimation, ok := value.([]AnimationProperty); ok && len(oldAnimation) > 0 {
			animation.oldAnimation = oldAnimation
		}
	}

	animation.oldListeners = map[PropertyName][]oneArgListener[View, PropertyName]{}

	setListeners := func(event PropertyName, listener func(View, PropertyName)) {
		listeners := getOneArgEventListeners[View, PropertyName](view, nil, event)
		if len(listeners) > 0 {
			animation.oldListeners[event] = slices.Clone(listeners)
		}
		view.Set(event, append(listeners, newOneArgListenerVE(listener)))
	}

	setListeners(AnimationStartEvent, animation.onAnimationStart)
	setListeners(AnimationEndEvent, animation.onAnimationEnd)
	setListeners(AnimationCancelEvent, animation.onAnimationCancel)
	setListeners(AnimationIterationEvent, animation.onAnimationIteration)

	view.Set(Animation, animation)
	return true
}

func (animation *animationData) finish() {
	if animation.view != nil {
		for _, event := range []PropertyName{AnimationStartEvent, AnimationEndEvent, AnimationCancelEvent, AnimationIterationEvent} {
			if listeners, ok := animation.oldListeners[event]; ok && len(listeners) > 0 {
				animation.view.Set(event, listeners)
			} else {
				animation.view.Remove(event)
			}
		}

		if animation.oldAnimation != nil {
			animation.view.Set(Animation, animation.oldAnimation)
			animation.oldAnimation = nil
		} else {
			animation.view.Set(Animation, "")
		}

		animation.oldListeners = map[PropertyName][]oneArgListener[View, PropertyName]{}

		animation.view = nil
		animation.listener = nil
	}
}

func (animation *animationData) Stop() {
	animation.onAnimationCancel(animation.view, "")
}

func (animation *animationData) Pause() {
	if animation.view != nil {
		animation.view.Set(AnimationPaused, true)
	}
}

func (animation *animationData) Resume() {
	if animation.view != nil {
		animation.view.Remove(AnimationPaused)
	}
}

func (animation *animationData) onAnimationStart(view View, _ PropertyName) {
	if animation.view != nil && animation.listener != nil {
		animation.listener(animation.view, animation, AnimationStartEvent)
	}
}

func (animation *animationData) onAnimationEnd(view View, _ PropertyName) {
	if animation.view != nil {
		animationView := animation.view
		listener := animation.listener

		if value, ok := animation.properties[PropertyTag]; ok {
			if props, ok := value.([]AnimatedProperty); ok {
				for _, prop := range props {
					animationView.setRaw(prop.Tag, prop.To)
				}
			}
		}

		animation.finish()
		if listener != nil {
			listener(animationView, animation, AnimationEndEvent)
		}
	}
}

func (animation *animationData) onAnimationIteration(view View, _ PropertyName) {
	if animation.view != nil && animation.listener != nil {
		animation.listener(animation.view, animation, AnimationIterationEvent)
	}
}

func (animation *animationData) onAnimationCancel(view View, _ PropertyName) {
	if animation.view != nil {
		animationView := animation.view
		listener := animation.listener

		if value, ok := animation.properties[PropertyTag]; ok {
			if props, ok := value.([]AnimatedProperty); ok {
				for _, prop := range props {
					animationView.Set(prop.Tag, prop.To)
				}
			}
		}

		animation.finish()
		if listener != nil {
			listener(animationView, animation, AnimationCancelEvent)
		}
	}
}
