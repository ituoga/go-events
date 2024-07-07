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
func Subscribe(fn any) {

	handlerVal := reflect.ValueOf(fn)

	funcType := handlerVal.Type()
	if funcType.NumIn() != 1 {
		panic("Handler must take exactly one argument")
	}

	argType := funcType.In(0)
	if !argType.Implements(reflect.TypeOf((*Eventer)(nil)).Elem()) {
		panic("Handler argument must implement MyType interface")
	}

	dummyArg := reflect.New(argType.Elem()).Interface().(Eventer)
	topic := dummyArg.EventName()

	if _, ok := subscribers[topic]; !ok {
		subscribers[topic] = []reflect.Value{}
	}

	subscribers[topic] = append(subscribers[topic], handlerVal)
}

// Publish publishes an event
func Publish(event Eventer) {
	if _, ok := subscribers[event.EventName()]; !ok {
		return
	}
	for _, fn := range subscribers[event.EventName()] {
		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(event)
		fn.Call(in)
	}
}
