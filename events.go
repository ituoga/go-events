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
func Register[T any]() {
	evt := *new(T)
	if _, ok := reg[GetEventName(evt)]; ok {
		panic("event already registered" + " + " + GetEventName(evt))
	}
	reg[GetEventName(evt)] = reflect.TypeOf(evt).Elem()
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

func GetUnmarshal[T any](body []byte) T {
	evt := Get(GetEventName(*new(T)))
	err := json.Unmarshal(body, &evt)
	if err != nil {
		log.Fatalf("error unmarshalling event: %s", err)
	}
	return evt.(T)
}

func MustMarshal(evt any) []byte {
	body, err := json.Marshal(evt)
	if err != nil {
		log.Fatalf("error marshalling event: %s", err)
	}
	return body
}

func GetEventName(event any) string {
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
