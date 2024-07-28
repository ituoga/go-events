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
	if _, ok := reg[getEventName(evt)]; ok {
		panic("event already registered" + " + " + getEventName(evt))
	}
	reg[getEventName(evt)] = reflect.TypeOf(evt).Elem()
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
	evt := Get(getEventName(*new(T)))
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
