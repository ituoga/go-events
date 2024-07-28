package events

import (
	"encoding/json"
	"log"
	"reflect"
)

/*
func init() {
	events.Register[*CustomerCreatedEvent]()
	events.Register[*CustomerDeletedEvent]()
}
*/

type Eventer interface {
	EventName() string
}

var reg = map[string]reflect.Type{}

// Register registers an pointer to event in the registry
// var _ events.Eventer = (*YourStruct)(nil)
// events.Register[*CustomEvent]()
func Register[T Eventer]() {
	evt := *new(T)
	if _, ok := reg[evt.EventName()]; ok {
		panic("event already registered" + " + " + evt.EventName())
	}
	reg[evt.EventName()] = reflect.TypeOf(evt).Elem()
}

// Get returns a new instance of the event from string
func Get(event string) any {
	evt, ok := reg[event]
	if !ok {
		log.Fatalf("event not found: %s", event)
	}
	elm := reflect.New(evt)
	return elm.Interface()
}

func GetUnmarshal[T any](event Eventer, body []byte) T {
	evt := Get(event.EventName())
	err := json.Unmarshal(body, &evt)
	if err != nil {
		log.Fatalf("error unmarshalling event: %s", err)
	}
	return evt.(T)
}

// func GetUnmarshal[T any](event string, body []byte) T {
// 	evt := Get(event)
// 	err := json.Unmarshal(body, &evt)
// 	if err != nil {
// 		log.Fatalf("error unmarshalling event: %s", err)
// 	}
// 	return evt.(T)
// }
