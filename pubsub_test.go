package events

import (
	"encoding/json"
	"testing"
)

type localEvent struct {
	Name string
}

func (e *localEvent) EventName() string {
	return "some.topic.here"
}

func TestMain(t *testing.T) {

	b, _ := json.Marshal(&localEvent{"testas"})
	a := []byte("")
	Before(func(name string, b []byte, an any) {
		a = b
	})

	Subscribe(func(e *localEvent) {})

	Publish(&localEvent{"testas"})

	if string(a) != string(b) {
		t.Fatal("event not received")
	}
}

func TestSub(t *testing.T) {
	count := 0
	Subscribe(func(e *localEvent) {
		count++
	})

	Publish(&localEvent{"testas"})
	if count == 0 {
		t.Fatal("event not received")
	}
}

func TestPub(t *testing.T) {
	count := 0
	Subscribe(func(e *localEvent) {
		count++
	})

	Publish(&ExampleEvent{})
	if count == 1 {
		t.Fatal("event not received")
	}
}
