package events

import (
	"encoding/json"
	"errors"
	"testing"
)

type localEvent struct {
	Name string
}

func (e *localEvent) EventName() string {
	return "some.topic.here"
}

type localEventv1 struct {
	Name string
}

func (e *localEventv1) EventName() string {
	return "some.topic.here.v1"
}

type localEventv2 struct {
	Name string
}

func (e *localEventv2) EventName() string {
	return "some.topic.here.v2"
}
func TestMain(t *testing.T) {

	b, _ := json.Marshal(&localEventv1{"testas"})
	a := []byte("")
	Before(func(name string, b []byte, an any) error {
		a = b
		return nil
	})

	Subscribe(func(e *localEventv1) error {
		return errors.New("error 1")
	})

	err := Publish(&localEventv1{"testas"})

	if string(a) != string(b) {
		t.Fatal("event not equals")
	}

	if err == nil {
		t.Fatal("error not received")
	}
}

func TestSub(t *testing.T) {
	count := 0
	Subscribe(func(e *localEventv2) error {
		count++
		return nil
	})

	Publish(&localEventv2{"testas"})
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

func TestRequestReply(t *testing.T) {
	Subscribe(func(e *ExampleEvent) (any, error) {
		return &ExampleEvent{}, nil
	})

	r, err := RequestG[*ExampleEvent](&ExampleEvent{})

	if err != nil {
		t.Fatal("error not expected")
	}

	if r.EventName() != "ExampleEvent" {
		t.Fatal("wrong response")
	}
}

func TestRequestReplyWithName(t *testing.T) {
	Subscribe(func(e *ExampleWithNameEvent) (any, error) {
		return &ExampleWithNameEvent{Name: "test"}, nil
	})

	r, err := RequestG[*ExampleWithNameEvent](&ExampleWithNameEvent{})

	if err != nil {
		t.Fatal("error not expected")
	}

	if r.Name != "test" {
		t.Fatal("wrong response")
	}
}

func TestRequestReplyError(t *testing.T) {
	Subscribe(func(e *ExampleWithNameEvent) (any, error) {
		return &ExampleWithNameEvent{Name: "test"}, errors.New("error")
	})

	r, err := RequestG[*ExampleWithNameEvent](&ExampleWithNameEvent{})

	if err.Error() != "error" {
		t.Fatalf("error not expected %v", err)
	}

	if r.Name != "test" {
		t.Fatal("wrong response")
	}
}
