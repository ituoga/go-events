package events

/*
ExampleEvent is an example event

type ExampleEvent struct{}

// EventName returns the name of the event

	func (e *ExampleEvent) EventName() string {
		return "ExampleEvent"
	}
*/
type ExampleEvent struct{}

// EventName returns the name of the event
func (e *ExampleEvent) EventName() string {
	return "ExampleEvent"
}

type ExampleWithNameEvent struct {
	Name string
}
