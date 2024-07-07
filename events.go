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

var reg = map[string]reflect.Type{}

// Register registers an pointer to event in the registry
// var _ events.Eventer = (*YourStruct)(nil)
// events.Register[*CustomEvent]()
func Register[T Eventer]() {
	evt := *new(T)
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

func GetUnmarshal(event string, body []byte) any {
	evt := Get(event)
	err := json.Unmarshal(body, &evt)
	if err != nil {
		log.Fatalf("error unmarshalling event: %s", err)
	}
	return evt
}
