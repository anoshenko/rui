package rui

func (animation *animationData) Start(view View, listener func(view View, animation Animation, event string)) bool {
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
	if value := view.Get(AnimationTag); value != nil {
		if oldAnimation, ok := value.([]Animation); ok && len(oldAnimation) > 0 {
			animation.oldAnimation = oldAnimation
		}
	}

	animation.oldListeners = map[string][]func(View, string){}

	setListeners := func(event string, listener func(View, string)) {
		var listeners []func(View, string) = nil
		if value := view.Get(event); value != nil {
			if oldListeners, ok := value.([]func(View, string)); ok && len(oldListeners) > 0 {
				listeners = oldListeners
			}
		}

		if listeners == nil {
			view.Set(event, listener)
		} else {
			animation.oldListeners[event] = listeners
			view.Set(event, append(listeners, listener))
		}
	}

	setListeners(AnimationStartEvent, animation.onAnimationStart)
	setListeners(AnimationEndEvent, animation.onAnimationEnd)
	setListeners(AnimationCancelEvent, animation.onAnimationCancel)
	setListeners(AnimationIterationEvent, animation.onAnimationIteration)

	view.Set(AnimationTag, animation)
	return true
}

func (animation *animationData) finish() {
	if animation.view != nil {
		for _, event := range []string{AnimationStartEvent, AnimationEndEvent, AnimationCancelEvent, AnimationIterationEvent} {
			if listeners, ok := animation.oldListeners[event]; ok {
				animation.view.Set(event, listeners)
			} else {
				animation.view.Remove(event)
			}
		}

		if animation.oldAnimation != nil {
			animation.view.Set(AnimationTag, animation.oldAnimation)
			animation.oldAnimation = nil
		} else {
			animation.view.Set(AnimationTag, "")
		}

		animation.oldListeners = map[string][]func(View, string){}

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

func (animation *animationData) onAnimationStart(view View, _ string) {
	if animation.view != nil && animation.listener != nil {
		animation.listener(animation.view, animation, AnimationStartEvent)
	}
}

func (animation *animationData) onAnimationEnd(view View, _ string) {
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

func (animation *animationData) onAnimationIteration(view View, _ string) {
	if animation.view != nil && animation.listener != nil {
		animation.listener(animation.view, animation, AnimationIterationEvent)
	}
}

func (animation *animationData) onAnimationCancel(view View, _ string) {
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
