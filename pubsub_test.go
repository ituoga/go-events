package events

import "testing"

type localEvent struct {
	Name string
}

func (e *localEvent) EventName() string {
	return "some.topic.here"
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
