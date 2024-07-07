package events

import "reflect"

var (
	// PubSub is the pubsub instance
	subscribers map[string][]reflect.Value
)

func init() {
	subscribers = make(map[string][]reflect.Value)
}

// Subscribe subscribes to an event
func Subscribe(event string, fn any) {
	if _, ok := subscribers[event]; !ok {
		subscribers[event] = []reflect.Value{}
	}
	subscribers[event] = append(subscribers[event], reflect.ValueOf(fn))
}

// Publish publishes an event
func Publish(event string, args ...any) {
	if _, ok := subscribers[event]; !ok {
		return
	}
	for _, fn := range subscribers[event] {
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			in[i] = reflect.ValueOf(arg)
		}
		fn.Call(in)
	}
}
