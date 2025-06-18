package rui

import "reflect"

type twoArgListener[V View, E any] interface {
	Run(V, E, E)
	rawListener() any
}

type twoArgListener0[V View, E any] struct {
	fn func()
}

type twoArgListenerV[V View, E any] struct {
	fn func(V)
}

type twoArgListenerE[V View, E any] struct {
	fn func(E)
}

type twoArgListenerVE[V View, E any] struct {
	fn func(V, E)
}

type twoArgListenerEE[V View, E any] struct {
	fn func(E, E)
}

type twoArgListenerVEE[V View, E any] struct {
	fn func(V, E, E)
}

type twoArgListenerBinding[V View, E any] struct {
	name string
}

func newTwoArgListener0[V View, E any](fn func()) twoArgListener[V, E] {
	obj := new(twoArgListener0[V, E])
	obj.fn = fn
	return obj
}

func (data *twoArgListener0[V, E]) Run(_ V, _ E, _ E) {
	data.fn()
}

func (data *twoArgListener0[V, E]) rawListener() any {
	return data.fn
}

func newTwoArgListenerV[V View, E any](fn func(V)) twoArgListener[V, E] {
	obj := new(twoArgListenerV[V, E])
	obj.fn = fn
	return obj
}

func (data *twoArgListenerV[V, E]) Run(view V, _ E, _ E) {
	data.fn(view)
}

func (data *twoArgListenerV[V, E]) rawListener() any {
	return data.fn
}

func newTwoArgListenerE[V View, E any](fn func(E)) twoArgListener[V, E] {
	obj := new(twoArgListenerE[V, E])
	obj.fn = fn
	return obj
}

func (data *twoArgListenerE[V, E]) Run(_ V, arg E, _ E) {
	data.fn(arg)
}

func (data *twoArgListenerE[V, E]) rawListener() any {
	return data.fn
}

func newTwoArgListenerVE[V View, E any](fn func(V, E)) twoArgListener[V, E] {
	obj := new(twoArgListenerVE[V, E])
	obj.fn = fn
	return obj
}

func (data *twoArgListenerVE[V, E]) Run(view V, arg E, _ E) {
	data.fn(view, arg)
}

func (data *twoArgListenerVE[V, E]) rawListener() any {
	return data.fn
}

func newTwoArgListenerEE[V View, E any](fn func(E, E)) twoArgListener[V, E] {
	obj := new(twoArgListenerEE[V, E])
	obj.fn = fn
	return obj
}

func (data *twoArgListenerEE[V, E]) Run(_ V, arg1 E, arg2 E) {
	data.fn(arg1, arg2)
}

func (data *twoArgListenerEE[V, E]) rawListener() any {
	return data.fn
}

func newTwoArgListenerVEE[V View, E any](fn func(V, E, E)) twoArgListener[V, E] {
	obj := new(twoArgListenerVEE[V, E])
	obj.fn = fn
	return obj
}

func (data *twoArgListenerVEE[V, E]) Run(view V, arg1 E, arg2 E) {
	data.fn(view, arg1, arg2)
}

func (data *twoArgListenerVEE[V, E]) rawListener() any {
	return data.fn
}

func newTwoArgListenerBinding[V View, E any](name string) twoArgListener[V, E] {
	obj := new(twoArgListenerBinding[V, E])
	obj.name = name
	return obj
}

func (data *twoArgListenerBinding[V, E]) Run(view V, arg1 E, arg2 E) {
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
		} else if inType == reflect.TypeOf(arg1) {
			args = []reflect.Value{reflect.ValueOf(arg1)}
		}

	case 2:
		valType := reflect.TypeOf(arg1)
		if methodType.In(0) == reflect.TypeOf(view) && methodType.In(1) == valType {
			args = []reflect.Value{reflect.ValueOf(view), reflect.ValueOf(arg1)}
		} else if methodType.In(0) == valType && methodType.In(1) == valType {
			args = []reflect.Value{reflect.ValueOf(arg1), reflect.ValueOf(arg2)}
		}

	case 3:
		valType := reflect.TypeOf(arg1)
		if methodType.In(0) == reflect.TypeOf(view) && methodType.In(1) == valType && methodType.In(2) == valType {
			args = []reflect.Value{reflect.ValueOf(view), reflect.ValueOf(arg1), reflect.ValueOf(arg2)}
		}
	}

	if args != nil {
		method.Call(args)
	} else {
		ErrorLogF(`Unsupported prototype of "%s" method`, data.name)
	}
}

func (data *twoArgListenerBinding[V, E]) rawListener() any {
	return data.name
}

func valueToTwoArgEventListeners[V View, E any](value any) ([]twoArgListener[V, E], bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case []twoArgListener[V, E]:
		return value, true

	case twoArgListener[V, E]:
		return []twoArgListener[V, E]{value}, true

	case string:
		return []twoArgListener[V, E]{newTwoArgListenerBinding[V, E](value)}, true

	case func(V, E):
		return []twoArgListener[V, E]{newTwoArgListenerVE(value)}, true

	case func(V):
		return []twoArgListener[V, E]{newTwoArgListenerV[V, E](value)}, true

	case func(E):
		return []twoArgListener[V, E]{newTwoArgListenerE[V](value)}, true

	case func():
		return []twoArgListener[V, E]{newTwoArgListener0[V, E](value)}, true

	case func(E, E):
		return []twoArgListener[V, E]{newTwoArgListenerEE[V](value)}, true

	case func(V, E, E):
		return []twoArgListener[V, E]{newTwoArgListenerVEE(value)}, true

	case []func(V, E):
		result := make([]twoArgListener[V, E], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newTwoArgListenerVE(fn))
			}
		}
		return result, len(result) > 0

	case []func(E):
		result := make([]twoArgListener[V, E], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newTwoArgListenerE[V](fn))
			}
		}
		return result, len(result) > 0

	case []func(V):
		result := make([]twoArgListener[V, E], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newTwoArgListenerV[V, E](fn))
			}
		}
		return result, len(result) > 0

	case []func():
		result := make([]twoArgListener[V, E], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newTwoArgListener0[V, E](fn))
			}
		}
		return result, len(result) > 0

	case []func(E, E):
		result := make([]twoArgListener[V, E], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newTwoArgListenerEE[V](fn))
			}
		}
		return result, len(result) > 0

	case []func(V, E, E):
		result := make([]twoArgListener[V, E], 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newTwoArgListenerVEE(fn))
			}
		}
		return result, len(result) > 0

	case []any:
		result := make([]twoArgListener[V, E], 0, len(value))
		for _, v := range value {
			if v != nil {
				switch v := v.(type) {
				case func(V, E):
					result = append(result, newTwoArgListenerVE(v))

				case func(E):
					result = append(result, newTwoArgListenerE[V](v))

				case func(V):
					result = append(result, newTwoArgListenerV[V, E](v))

				case func():
					result = append(result, newTwoArgListener0[V, E](v))

				case func(E, E):
					result = append(result, newTwoArgListenerEE[V](v))

				case func(V, E, E):
					result = append(result, newTwoArgListenerVEE(v))

				case string:
					result = append(result, newTwoArgListenerBinding[V, E](v))

				default:
					return nil, false
				}
			}
		}
		return result, len(result) > 0
	}

	return nil, false
}

func setTwoArgEventListener[V View, T any](view View, tag PropertyName, value any) []PropertyName {
	if listeners, ok := valueToTwoArgEventListeners[V, T](value); ok {
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

func getTwoArgEventListeners[V View, E any](view View, subviewID []string, tag PropertyName) []twoArgListener[V, E] {
	if view = getSubview(view, subviewID); view != nil {
		if value := view.Get(tag); value != nil {
			if result, ok := value.([]twoArgListener[V, E]); ok {
				return result
			}
		}
	}
	return []twoArgListener[V, E]{}
}

func getTwoArgEventRawListeners[V View, E any](view View, subviewID []string, tag PropertyName) []any {
	listeners := getTwoArgEventListeners[V, E](view, subviewID, tag)
	result := make([]any, len(listeners))
	for i, l := range listeners {
		result[i] = l.rawListener()
	}
	return result
}
