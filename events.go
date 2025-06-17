package rui

import (
	"reflect"
	"strings"
)

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
	DragEndEvent:            {jsEvent: "ondragend", jsFunc: "dragEndEvent"},
	DragEnterEvent:          {jsEvent: "ondragenter", jsFunc: "dragEnterEvent"},
	DragLeaveEvent:          {jsEvent: "ondragleave", jsFunc: "dragLeaveEvent"},
}

/*
type oneArgListener[V View, E any] interface {
	call(V, E)
	listener() func(V, E)
	rawListener() any
}

type oneArgListener0[V View, E any] struct {
	fn func()
}

type oneArgListenerV[V View, E any] struct {
	fn func(V)
}

type oneArgListenerE[V View, E any] struct {
	fn func(E)
}

type oneArgListenerVE[V View, E any] struct {
	fn func(V, E)
}

type oneArgListenerBinding[V View, E any] struct {
	name string
}

func newOneArgListener0[V View, E any](fn func()) oneArgListener[V, E] {
	obj := new(oneArgListener0[V, E])
	obj.fn = fn
	return obj
}

func (data *oneArgListener0[V, E]) call(_ V, _ E) {
	data.fn()
}

func (data *oneArgListener0[V, E]) listener() func(V, E) {
	return data.call
}

func (data *oneArgListener0[V, E]) rawListener() any {
	return data.fn
}

func newOneArgListenerV[V View, E any](fn func(V)) oneArgListener[V, E] {
	obj := new(oneArgListenerV[V, E])
	obj.fn = fn
	return obj
}

func (data *oneArgListenerV[V, E]) call(view V, _ E) {
	data.fn(view)
}

func (data *oneArgListenerV[V, E]) listener() func(V, E) {
	return data.call
}

func (data *oneArgListenerV[V, E]) rawListener() any {
	return data.fn
}

func newOneArgListenerE[V View, E any](fn func(E)) oneArgListener[V, E] {
	obj := new(oneArgListenerE[V, E])
	obj.fn = fn
	return obj
}

func (data *oneArgListenerE[V, E]) call(_ V, event E) {
	data.fn(event)
}

func (data *oneArgListenerE[V, E]) listener() func(V, E) {
	return data.call
}

func (data *oneArgListenerE[V, E]) rawListener() any {
	return data.fn
}

func newOneArgListenerVE[V View, E any](fn func(V, E)) oneArgListener[V, E] {
	obj := new(oneArgListenerVE[V, E])
	obj.fn = fn
	return obj
}

func (data *oneArgListenerVE[V, E]) call(view V, arg E) {
	data.fn(view, arg)
}

func (data *oneArgListenerVE[V, E]) listener() func(V, E) {
	return data.fn
}

func (data *oneArgListenerVE[V, E]) rawListener() any {
	return data.fn
}

func newOneArgListenerBinding[V View, E any](name string) oneArgListener[V, E] {
	obj := new(oneArgListenerBinding[V, E])
	obj.name = name
	return obj
}

func (data *oneArgListenerBinding[V, E]) call(view V, event E) {
	bind := view.binding()
	if bind == nil {
		ErrorLogF(`There is no a binding object for call "%s"`, data.name)
		return
	}

	val := reflect.ValueOf(bind)
	method := val.MethodByName(data.name)
	if !method.IsValid() {
		ErrorLogF(`The "%s" method is not valid`, data.name)
		return
	}

	methodType := method.Type()
	var args []reflect.Value = nil
	switch methodType.NumIn() {
	case 0:
		args = []reflect.Value{}

	case 1:
		inType := methodType.In(0)
		if inType == reflect.TypeOf(view) {
			args = []reflect.Value{reflect.ValueOf(view)}
		} else if inType == reflect.TypeOf(event) {
			args = []reflect.Value{reflect.ValueOf(event)}
		}

	case 2:
		if methodType.In(0) == reflect.TypeOf(view) && methodType.In(1) == reflect.TypeOf(event) {
			args = []reflect.Value{reflect.ValueOf(view), reflect.ValueOf(event)}
		}
	}

	if args != nil {
		method.Call(args)
	} else {
		ErrorLogF(`Unsupported prototype of "%s" method`, data.name)
	}
}

func (data *oneArgListenerBinding[V, E]) listener() func(V, E) {
	return data.call
}

func (data *oneArgListenerBinding[V, E]) rawListener() any {
	return data.name
}
*/

func valueToNoArgEventListeners[V any](view View, value any) ([]func(V), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case string:
		fn := func(arg V) {
			bind := view.binding()
			if bind == nil {
				ErrorLogF(`There is no a binding object for call "%s"`, value)
				return
			}

			val := reflect.ValueOf(bind)
			method := val.MethodByName(value)
			if !method.IsValid() {
				ErrorLogF(`The "%s" method is not valid`, value)
				return
			}

			methodType := method.Type()
			var args []reflect.Value = nil
			switch methodType.NumIn() {
			case 0:
				args = []reflect.Value{}

			case 1:
				inType := methodType.In(0)
				if inType == reflect.TypeOf(arg) {
					args = []reflect.Value{reflect.ValueOf(arg)}
				}
			}

			if args != nil {
				method.Call(args)
			} else {
				ErrorLogF(`Unsupported prototype of "%s" method`, value)
			}
		}
		return []func(V){fn}, true

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

/*
	func valueToOneArgEventListeners[V View, E any](view View, value any) ([]oneArgListener[V, E], bool) {
		if value == nil {
			return nil, true
		}

		switch value := value.(type) {
		case string:
			return []oneArgListener[V, E]{newOneArgListenerBinding[V, E](value)}, true

		case func(V, E):
			return []oneArgListener[V, E]{newOneArgListenerVE[V, E](value)}, true

		case func(V):
			return []oneArgListener[V, E]{newOneArgListenerV[V, E](value)}, true

		case func(E):
			return []oneArgListener[V, E]{newOneArgListenerE[V, E](value)}, true

		case func():
			return []oneArgListener[V, E]{newOneArgListener0[V, E](value)}, true

		case []func(V, E):
			result := make([]oneArgListener[V, E], 0, len(value))
			for _, fn := range value {
				if fn != nil {
					result = append(result, newOneArgListenerVE[V, E](fn))
				}
			}
			return result, len(result) > 0

		case []func(E):
			result := make([]oneArgListener[V, E], 0, len(value))
			for _, fn := range value {
				if fn != nil {
					result = append(result, newOneArgListenerE[V, E](fn))
				}
			}
			return result, len(result) > 0

		case []func(V):
			result := make([]oneArgListener[V, E], 0, len(value))
			for _, fn := range value {
				if fn != nil {
					result = append(result, newOneArgListenerV[V, E](fn))
				}
			}
			return result, len(result) > 0

		case []func():
			result := make([]oneArgListener[V, E], 0, len(value))
			for _, fn := range value {
				if fn != nil {
					result = append(result, newOneArgListener0[V, E](fn))
				}
			}
			return result, len(result) > 0

		case []any:
			result := make([]oneArgListener[V, E], 0, len(value))
			for _, v := range value {
				if v != nil {
					switch v := v.(type) {
					case func(V, E):
						result = append(result, newOneArgListenerVE[V, E](v))

					case func(E):
						result = append(result, newOneArgListenerE[V, E](v))

					case func(V):
						result = append(result, newOneArgListenerV[V, E](v))

					case func():
						result = append(result, newOneArgListener0[V, E](v))

					case string:
						result = append(result, newOneArgListenerBinding[V, E](v))

					default:
						return nil, false
					}
				}
			}
			return result, len(result) > 0
		}

		return nil, false
	}
*/
func valueToTwoArgEventListeners[V View, E any](view View, value any) ([]func(V, E, E), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case string:
		fn := func(view V, val1 E, val2 E) {
			bind := view.binding()
			if bind == nil {
				ErrorLogF(`There is no a binding object for call "%s"`, value)
				return
			}

			val := reflect.ValueOf(bind)
			method := val.MethodByName(value)
			if !method.IsValid() {
				ErrorLogF(`The "%s" method is not valid`, value)
				return
			}

			methodType := method.Type()
			var args []reflect.Value = nil
			switch methodType.NumIn() {
			case 0:
				args = []reflect.Value{}

			case 1:
				inType := methodType.In(0)
				if inType == reflect.TypeOf(view) {
					args = []reflect.Value{reflect.ValueOf(view)}
				} else if inType == reflect.TypeOf(val1) {
					args = []reflect.Value{reflect.ValueOf(val1)}
				}

			case 2:
				valType := reflect.TypeOf(val1)
				if methodType.In(0) == reflect.TypeOf(view) && methodType.In(1) == valType {
					args = []reflect.Value{reflect.ValueOf(view), reflect.ValueOf(val1)}
				} else if methodType.In(0) == valType && methodType.In(1) == valType {
					args = []reflect.Value{reflect.ValueOf(val1), reflect.ValueOf(val2)}
				}

			case 3:
				valType := reflect.TypeOf(val1)
				if methodType.In(0) == reflect.TypeOf(view) && methodType.In(1) == valType && methodType.In(2) == valType {
					args = []reflect.Value{reflect.ValueOf(view), reflect.ValueOf(val1), reflect.ValueOf(val2)}
				}
			}

			if args != nil {
				method.Call(args)
			} else {
				ErrorLogF(`Unsupported prototype of "%s" method`, value)
			}
		}
		return []func(V, E, E){fn}, true

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

/*
	func getOneArgEventListeners[V View, E any](view View, subviewID []string, tag PropertyName) []oneArgListener[V, E] {
		if view = getSubview(view, subviewID); view != nil {
			if value := view.Get(tag); value != nil {
				if result, ok := value.([]oneArgListener[V, E]); ok {
					return result
				}
			}
		}
		return []oneArgListener[V, E]{}
	}

	func getOneArgEventRawListeners[V View, E any](view View, subviewID []string, tag PropertyName) []any {
		listeners := getOneArgEventListeners[V, E](view, subviewID, tag)
		result := make([]any, len(listeners))
		for i, l := range listeners {
			result[i] = l.rawListener()
		}
		return result
	}
*/
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

func setNoArgEventListener[V View](view View, tag PropertyName, value any) []PropertyName {
	if listeners, ok := valueToNoArgEventListeners[V](view, value); ok {
		if len(listeners) > 0 {
			view.setRaw(tag, listeners)
		} else if view.getRaw(tag) != nil {
			view.setRaw(tag, nil)
		} else {
			return []PropertyName{}
		}
		return []PropertyName{tag}
	}
	notCompatibleType(tag, value)
	return nil
}

/*
func setOneArgEventListener[V View, T any](view View, tag PropertyName, value any) []PropertyName {
	if listeners, ok := valueToOneArgEventListeners[V, T](view, value); ok {
		if len(listeners) > 0 {
			view.setRaw(tag, listeners)
		} else if view.getRaw(tag) != nil {
			view.setRaw(tag, nil)
		} else {
			return []PropertyName{}
		}
		return []PropertyName{tag}
	}
	notCompatibleType(tag, value)
	return nil
}
*/

func setTwoArgEventListener[V View, T any](view View, tag PropertyName, value any) []PropertyName {
	listeners, ok := valueToTwoArgEventListeners[V, T](view, value)
	if !ok {
		notCompatibleType(tag, value)
		return nil
	}
	if len(listeners) > 0 {
		view.setRaw(tag, listeners)
	} else if view.getRaw(tag) != nil {
		view.setRaw(tag, nil)
	} else {
		return []PropertyName{}
	}
	return []PropertyName{tag}
}

func viewEventsHtml[T any](view View, events []PropertyName, buffer *strings.Builder) {
	for _, tag := range events {
		if js, ok := eventJsFunc[tag]; ok {
			if value := getOneArgEventListeners[View, T](view, nil, tag); len(value) > 0 {
				buffer.WriteString(js.jsEvent)
				buffer.WriteString(`="`)
				buffer.WriteString(js.jsFunc)
				buffer.WriteString(`(this, event)" `)
			}
		}
		/*if value := view.getRaw(tag); value != nil {
			if js, ok := eventJsFunc[tag]; ok {
				if listeners, ok := value.([] func(View, T)); ok && len(listeners) > 0 {
					buffer.WriteString(js.jsEvent)
					buffer.WriteString(`="`)
					buffer.WriteString(js.jsFunc)
					buffer.WriteString(`(this, event)" `)
				}
			}
		}
		*/
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
