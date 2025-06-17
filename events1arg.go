package rui

import (
	"reflect"
)

type oneArgListener[V View, E any] interface {
	Run(V, E)
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

func (data *oneArgListener0[V, E]) Run(_ V, _ E) {
	data.fn()
}

func (data *oneArgListener0[V, E]) rawListener() any {
	return data.fn
}

func newOneArgListenerV[V View, E any](fn func(V)) oneArgListener[V, E] {
	obj := new(oneArgListenerV[V, E])
	obj.fn = fn
	return obj
}

func (data *oneArgListenerV[V, E]) Run(view V, _ E) {
	data.fn(view)
}

func (data *oneArgListenerV[V, E]) rawListener() any {
	return data.fn
}

func newOneArgListenerE[V View, E any](fn func(E)) oneArgListener[V, E] {
	obj := new(oneArgListenerE[V, E])
	obj.fn = fn
	return obj
}

func (data *oneArgListenerE[V, E]) Run(_ V, event E) {
	data.fn(event)
}

func (data *oneArgListenerE[V, E]) rawListener() any {
	return data.fn
}

func newOneArgListenerVE[V View, E any](fn func(V, E)) oneArgListener[V, E] {
	obj := new(oneArgListenerVE[V, E])
	obj.fn = fn
	return obj
}

func (data *oneArgListenerVE[V, E]) Run(view V, arg E) {
	data.fn(view, arg)
}

func (data *oneArgListenerVE[V, E]) rawListener() any {
	return data.fn
}

func newOneArgListenerBinding[V View, E any](name string) oneArgListener[V, E] {
	obj := new(oneArgListenerBinding[V, E])
	obj.name = name
	return obj
}

func (data *oneArgListenerBinding[V, E]) Run(view V, event E) {
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

func (data *oneArgListenerBinding[V, E]) rawListener() any {
	return data.name
}

func valueToOneArgEventListeners[V View, E any](value any) ([]oneArgListener[V, E], bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case []oneArgListener[V, E]:
		return value, true

	case oneArgListener[V, E]:
		return []oneArgListener[V, E]{value}, true

	case string:
		return []oneArgListener[V, E]{newOneArgListenerBinding[V, E](value)}, true

	case func(V, E):
		return []oneArgListener[V, E]{newOneArgListenerVE(value)}, true

	case func(V):
		return []oneArgListener[V, E]{newOneArgListenerV[V, E](value)}, true

	case func(E):
		return []oneArgListener[V, E]{newOneArgListenerE[V](value)}, true

	case func():
		return []oneArgListener[V, E]{newOneArgListener0[V, E](value)}, true

	case []func(V, E):
		result := make([]oneArgListener[V, E], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newOneArgListenerVE(fn))
			}
		}
		return result, len(result) > 0

	case []func(E):
		result := make([]oneArgListener[V, E], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newOneArgListenerE[V](fn))
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
					result = append(result, newOneArgListenerVE(v))

				case func(E):
					result = append(result, newOneArgListenerE[V](v))

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

func setOneArgEventListener[V View, T any](view View, tag PropertyName, value any) []PropertyName {
	if listeners, ok := valueToOneArgEventListeners[V, T](value); ok {
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
