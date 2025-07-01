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

type noArgListener[V View] interface {
	Run(V)
	rawListener() any
}

type noArgListener0[V View] struct {
	fn func()
}

type noArgListenerV[V View] struct {
	fn func(V)
}

type noArgListenerBinding[V View] struct {
	name string
}

func newNoArgListener0[V View](fn func()) noArgListener[V] {
	obj := new(noArgListener0[V])
	obj.fn = fn
	return obj
}

func (data *noArgListener0[V]) Run(_ V) {
	data.fn()
}

func (data *noArgListener0[V]) rawListener() any {
	return data.fn
}

func newNoArgListenerV[V View](fn func(V)) noArgListener[V] {
	obj := new(noArgListenerV[V])
	obj.fn = fn
	return obj
}

func (data *noArgListenerV[V]) Run(view V) {
	data.fn(view)
}

func (data *noArgListenerV[V]) rawListener() any {
	return data.fn
}

func newNoArgListenerBinding[V View](name string) noArgListener[V] {
	obj := new(noArgListenerBinding[V])
	obj.name = name
	return obj
}

func (data *noArgListenerBinding[V]) Run(view V) {
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
		if equalType(methodType.In(0), reflect.TypeOf(view)) {
			args = []reflect.Value{reflect.ValueOf(view)}
		}
	}

	if args != nil {
		method.Call(args)
	} else {
		ErrorLogF(`Unsupported prototype of "%s" method`, data.name)
	}
}

func equalType(inType reflect.Type, argType reflect.Type) bool {
	return inType == argType || (inType.Kind() == reflect.Interface && argType.Implements(inType))
}

func (data *noArgListenerBinding[V]) rawListener() any {
	return data.name
}

func valueToNoArgEventListeners[V View](value any) ([]noArgListener[V], bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case []noArgListener[V]:
		return value, true

	case noArgListener[V]:
		return []noArgListener[V]{value}, true

	case string:
		return []noArgListener[V]{newNoArgListenerBinding[V](value)}, true

	case func(V):
		return []noArgListener[V]{newNoArgListenerV(value)}, true

	case func():
		return []noArgListener[V]{newNoArgListener0[V](value)}, true

	case []func(V):
		result := make([]noArgListener[V], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newNoArgListenerV(fn))
			}
		}
		return result, len(result) > 0

	case []func():
		result := make([]noArgListener[V], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newNoArgListener0[V](fn))
			}
		}
		return result, len(result) > 0

	case []any:
		result := make([]noArgListener[V], 0, len(value))
		for _, v := range value {
			if v != nil {
				switch v := v.(type) {
				case func(V):
					result = append(result, newNoArgListenerV(v))

				case func():
					result = append(result, newNoArgListener0[V](v))

				case string:
					result = append(result, newNoArgListenerBinding[V](v))

				default:
					return nil, false
				}
			}
		}
		return result, len(result) > 0
	}

	return nil, false
}

func setNoArgEventListener[V View](view View, tag PropertyName, value any) []PropertyName {
	if listeners, ok := valueToNoArgEventListeners[V](value); ok {
		return setArrayPropertyValue(view, tag, listeners)
	}
	notCompatibleType(tag, value)
	return nil
}

func getNoArgEventListeners[V View](view View, subviewID []string, tag PropertyName) []noArgListener[V] {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]noArgListener[V]); ok {
				return result
			}
		}
	}
	return []noArgListener[V]{}
}

func getNoArgEventRawListeners[V View](view View, subviewID []string, tag PropertyName) []any {
	listeners := getNoArgEventListeners[V](view, subviewID, tag)
	result := make([]any, len(listeners))
	for i, l := range listeners {
		result[i] = l.rawListener()
	}
	return result
}

func getNoArgBinding[V View](listeners []noArgListener[V]) string {
	for _, listener := range listeners {
		raw := listener.rawListener()
		if text, ok := raw.(string); ok && text != "" {
			return text
		}
	}
	return ""
}
