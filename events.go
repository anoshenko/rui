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
