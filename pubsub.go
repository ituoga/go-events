package events

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	ErrNoSubscribers = errors.New("no subscribers for this event")
)

var (
	// PubSub is the pubsub instance
	subscribers map[string][]reflect.Value

	beforePub func(string, []byte, any) error
)

func init() {
	subscribers = make(map[string][]reflect.Value)
}

func Before(fn func(string, []byte, any) error) {
	beforePub = fn
}

// Subscribe subscribes to an event
func Subscribe(fn any) {

	handlerVal := reflect.ValueOf(fn)

	funcType := handlerVal.Type()
	if funcType.NumIn() != 1 {
		panic("Handler must take exactly one argument")
	}

	var topic string

	argType := funcType.In(0)
	if argType.Implements(reflect.TypeOf((*Eventer)(nil)).Elem()) {
		dummyArg := reflect.New(argType.Elem()).Interface().(Eventer)
		topic = dummyArg.EventName()
		// panic("Handler argument must implement MyType interface")
	} else {
		topic = getStructName(argType)
		if topic == "" {
			panic("Handler argument must implement Eventer interface or be a struct")
		}
	}

	if _, ok := subscribers[topic]; !ok {
		subscribers[topic] = make([]reflect.Value, 1)
	}

	// subscribers[topic] = append(subscribers[topic], handlerVal)
	subscribers[topic][0] = handlerVal
}

// Publish publishes an event
func Publish(event any) error {
	if _, ok := subscribers[getEventName(event)]; !ok {
		return ErrNoSubscribers
	}
	if beforePub != nil {
		b, err := json.Marshal(event)
		if err != nil {
			return err
		}
		err = beforePub(getEventName(event), b, event)
		if err != nil {
			return err
		}
	}
	for _, fn := range subscribers[getEventName(event)] {
		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(event)
		results := fn.Call(in)
		if len(results) > 0 {
			if err, ok := results[len(results)-1].Interface().(error); ok && err != nil {
				return err
			}
		}
	}
	return nil
}

func RequestEvent(name string, b []byte) (any, error) {
	return RequestGE[any](name, b)
}

func RequestEventBytes(name string, b []byte) ([]byte, error) {
	response, err := RequestGE[any](name, b)
	if err != nil {
		return nil, err
	}
	return json.Marshal(response)
}

func Request(event any) (any, error) {
	return RequestG[any](event)
}

func RequestG[T any](event any) (T, error) {
	if _, ok := subscribers[getEventName(event)]; !ok {
		return *new(T), ErrNoSubscribers
	}
	if beforePub != nil {
		b, err := json.Marshal(event)
		if err != nil {
			return *new(T), err
		}
		err = beforePub(getEventName(event), b, event)
		if err != nil {
			return *new(T), err
		}
	}
	var results []reflect.Value
	for _, fn := range subscribers[getEventName(event)] {
		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(event)
		results = append(results, fn.Call(in)...)
	}
	if len(results) == 0 {
		return *new(T), ErrNoSubscribers
	}
	if len(results) == 2 {
		if results[1].Interface() != nil {
			if results[0].Interface() != nil {
				return results[0].Interface().(T), results[1].Interface().(error)
			}
			return *new(T), results[1].Interface().(error)
		}
		return results[0].Interface().(T), nil
	}
	return *new(T), errors.New("not implemented")
}

func RequestGE[T any](name string, b []byte) (T, error) {
	if _, ok := subscribers[name]; !ok {
		return *new(T), ErrNoSubscribers
	}
	// if beforePub != nil {
	// 	b, err := json.Marshal(event)
	// 	if err != nil {
	// 		return *new(T), err
	// 	}
	// 	err = beforePub(getEventName(event), b, event)
	// 	if err != nil {
	// 		return *new(T), err
	// 	}
	// }

	event := Get(name)
	var results []reflect.Value
	for _, fn := range subscribers[name] {
		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(event)
		results = append(results, fn.Call(in)...)
	}
	if len(results) == 0 {
		return *new(T), ErrNoSubscribers
	}
	if len(results) == 2 {
		if results[1].Interface() != nil {
			if results[0].Interface() != nil {
				return results[0].Interface().(T), results[1].Interface().(error)
			}
			return *new(T), results[1].Interface().(error)
		}
		return results[0].Interface().(T), nil
	}
	return *new(T), errors.New("not implemented")
}

func getEventName(event any) string {
	if e, ok := event.(Eventer); ok {
		return e.EventName()
	}
	return getStructName(reflect.TypeOf(event))
}

func getStructName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		return t.Name()
	}

	return ""
}
