package rui

import "strings"

var eventJsFunc = map[PropertyName]struct{ jsEvent, jsFunc string }{
	FocusEvent:              {jsEvent: "onfocus", jsFunc: "focusEvent"},
	LostFocusEvent:          {jsEvent: "onblur", jsFunc: "blurEvent"},
	KeyDownEvent:            {jsEvent: "onkeydown", jsFunc: "keyDownEvent"},
	KeyUpEvent:              {jsEvent: "onkeyup", jsFunc: "keyUpEvent"},
	ClickEvent:              {jsEvent: "onclick", jsFunc: "clickEvent"},
	DoubleClickEvent:        {jsEvent: "ondblclick", jsFunc: "doubleClickEvent"},
	MouseDown:               {jsEvent: "onmousedown", jsFunc: "mouseDownEvent"},
	MouseUp:                 {jsEvent: "onmouseup", jsFunc: "mouseUpEvent"},
	MouseMove:               {jsEvent: "onmousemove", jsFunc: "mouseMoveEvent"},
	MouseOut:                {jsEvent: "onmouseout", jsFunc: "mouseOutEvent"},
	MouseOver:               {jsEvent: "onmouseover", jsFunc: "mouseOverEvent"},
	ContextMenuEvent:        {jsEvent: "oncontextmenu", jsFunc: "contextMenuEvent"},
	PointerDown:             {jsEvent: "onpointerdown", jsFunc: "pointerDownEvent"},
	PointerUp:               {jsEvent: "onpointerup", jsFunc: "pointerUpEvent"},
	PointerMove:             {jsEvent: "onpointermove", jsFunc: "pointerMoveEvent"},
	PointerCancel:           {jsEvent: "onpointercancel", jsFunc: "pointerCancelEvent"},
	PointerOut:              {jsEvent: "onpointerout", jsFunc: "pointerOutEvent"},
	PointerOver:             {jsEvent: "onpointerover", jsFunc: "pointerOverEvent"},
	TouchStart:              {jsEvent: "ontouchstart", jsFunc: "touchStartEvent"},
	TouchEnd:                {jsEvent: "ontouchend", jsFunc: "touchEndEvent"},
	TouchMove:               {jsEvent: "ontouchmove", jsFunc: "touchMoveEvent"},
	TouchCancel:             {jsEvent: "ontouchcancel", jsFunc: "touchCancelEvent"},
	TransitionRunEvent:      {jsEvent: "ontransitionrun", jsFunc: "transitionRunEvent"},
	TransitionStartEvent:    {jsEvent: "ontransitionstart", jsFunc: "transitionStartEvent"},
	TransitionEndEvent:      {jsEvent: "ontransitionend", jsFunc: "transitionEndEvent"},
	TransitionCancelEvent:   {jsEvent: "ontransitioncancel", jsFunc: "transitionCancelEvent"},
	AnimationStartEvent:     {jsEvent: "onanimationstart", jsFunc: "animationStartEvent"},
	AnimationEndEvent:       {jsEvent: "onanimationend", jsFunc: "animationEndEvent"},
	AnimationIterationEvent: {jsEvent: "onanimationiteration", jsFunc: "animationIterationEvent"},
	AnimationCancelEvent:    {jsEvent: "onanimationcancel", jsFunc: "animationCancelEvent"},
}

func valueToNoArgEventListeners[V any](value any) ([]func(V), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(V):
		return []func(V){value}, true

	case func():
		fn := func(V) {
			value()
		}
		return []func(V){fn}, true

	case []func(V):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(V) {
				v()
			}
		}
		return listeners, true

	case []any:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(V):
				listeners[i] = v

			case func():
				listeners[i] = func(V) {
					v()
				}

			default:
				return nil, false
			}
		}
		return listeners, true
	}

	return nil, false
}

func valueToOneArgEventListeners[V View, E any](value any) ([]func(V, E), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(V, E):
		return []func(V, E){value}, true

	case func(E):
		fn := func(_ V, event E) {
			value(event)
		}
		return []func(V, E){fn}, true

	case func(V):
		fn := func(view V, _ E) {
			value(view)
		}
		return []func(V, E){fn}, true

	case func():
		fn := func(V, E) {
			value()
		}
		return []func(V, E){fn}, true

	case []func(V, E):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(_ V, event E) {
				v(event)
			}
		}
		return listeners, true

	case []func(V):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(view V, _ E) {
				v(view)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(V, E) {
				v()
			}
		}
		return listeners, true

	case []any:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(V, E):
				listeners[i] = v

			case func(E):
				listeners[i] = func(_ V, event E) {
					v(event)
				}

			case func(V):
				listeners[i] = func(view V, _ E) {
					v(view)
				}

			case func():
				listeners[i] = func(V, E) {
					v()
				}

			default:
				return nil, false
			}
		}
		return listeners, true
	}

	return nil, false
}

func valueToTwoArgEventListeners[V View, E any](value any) ([]func(V, E, E), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(V, E, E):
		return []func(V, E, E){value}, true

	case func(V, E):
		fn := func(v V, val, _ E) {
			value(v, val)
		}
		return []func(V, E, E){fn}, true

	case func(E, E):
		fn := func(_ V, val, old E) {
			value(val, old)
		}
		return []func(V, E, E){fn}, true

	case func(E):
		fn := func(_ V, val, _ E) {
			value(val)
		}
		return []func(V, E, E){fn}, true

	case func(V):
		fn := func(v V, _, _ E) {
			value(v)
		}
		return []func(V, E, E){fn}, true

	case func():
		fn := func(V, E, E) {
			value()
		}
		return []func(V, E, E){fn}, true

	case []func(V, E, E):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(V, E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(view V, val, _ E) {
				fn(view, val)
			}
		}
		return listeners, true

	case []func(E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(_ V, val, _ E) {
				fn(val)
			}
		}
		return listeners, true

	case []func(E, E):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(_ V, val, old E) {
				fn(val, old)
			}
		}
		return listeners, true

	case []func(V):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(view V, _, _ E) {
				fn(view)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, fn := range value {
			if fn == nil {
				return nil, false
			}
			listeners[i] = func(V, E, E) {
				fn()
			}
		}
		return listeners, true

	case []any:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(V, E, E), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch fn := v.(type) {
			case func(V, E, E):
				listeners[i] = fn

			case func(V, E):
				listeners[i] = func(view V, val, _ E) {
					fn(view, val)
				}

			case func(E, E):
				listeners[i] = func(_ V, val, old E) {
					fn(val, old)
				}

			case func(E):
				listeners[i] = func(_ V, val, _ E) {
					fn(val)
				}

			case func(V):
				listeners[i] = func(view V, _, _ E) {
					fn(view)
				}

			case func():
				listeners[i] = func(V, E, E) {
					fn()
				}

			default:
				return nil, false
			}
		}
		return listeners, true
	}

	return nil, false
}

func getNoArgEventListeners[V View](view View, subviewID []string, tag PropertyName) []func(V) {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(V)); ok {
				return result
			}
		}
	}
	return []func(V){}
}

func getOneArgEventListeners[V View, E any](view View, subviewID []string, tag PropertyName) []func(V, E) {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(V, E)); ok {
				return result
			}
		}
	}
	return []func(V, E){}
}

func getTwoArgEventListeners[V View, E any](view View, subviewID []string, tag PropertyName) []func(V, E, E) {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]func(V, E, E)); ok {
				return result
			}
		}
	}
	return []func(V, E, E){}
}

func setNoArgEventListener[V View](properties Properties, tag PropertyName, value any) []PropertyName {
	if listeners, ok := valueToNoArgEventListeners[V](value); ok {
		if len(listeners) > 0 {
			properties.setRaw(tag, listeners)
		} else if properties.getRaw(tag) != nil {
			properties.setRaw(tag, nil)
		} else {
			return []PropertyName{}
		}
		return []PropertyName{tag}
	}
	notCompatibleType(tag, value)
	return nil
}

func setOneArgEventListener[V View, T any](properties Properties, tag PropertyName, value any) []PropertyName {
	if listeners, ok := valueToOneArgEventListeners[V, T](value); ok {
		if len(listeners) > 0 {
			properties.setRaw(tag, listeners)
		} else if properties.getRaw(tag) != nil {
			properties.setRaw(tag, nil)
		} else {
			return []PropertyName{}
		}
		return []PropertyName{tag}
	}
	notCompatibleType(tag, value)
	return nil
}

func setTwoArgEventListener[V View, T any](properties Properties, tag PropertyName, value any) []PropertyName {
	listeners, ok := valueToTwoArgEventListeners[V, T](value)
	if !ok {
		notCompatibleType(tag, value)
		return nil
	} else if len(listeners) > 0 {
		properties.setRaw(tag, listeners)
	} else if properties.getRaw(tag) != nil {
		properties.setRaw(tag, nil)
	} else {
		return []PropertyName{}
	}
	return []PropertyName{tag}
}

func viewEventsHtml[T any](view View, events []PropertyName, buffer *strings.Builder) {
	for _, tag := range events {
		if value := view.getRaw(tag); value != nil {
			if js, ok := eventJsFunc[tag]; ok {
				if listeners, ok := value.([]func(View, T)); ok && len(listeners) > 0 {
					buffer.WriteString(js.jsEvent)
					buffer.WriteString(`="`)
					buffer.WriteString(js.jsFunc)
					buffer.WriteString(`(this, event)" `)
				}
			}
		}
	}
}

func updateEventListenerHtml(view View, tag PropertyName) {
	if js, ok := eventJsFunc[tag]; ok {
		value := view.getRaw(tag)
		session := view.Session()
		htmlID := view.htmlID()
		if value == nil {
			session.removeProperty(view.htmlID(), js.jsEvent)
		} else {
			session.updateProperty(htmlID, js.jsEvent, js.jsFunc+"(this, event)")
		}
	}
}
