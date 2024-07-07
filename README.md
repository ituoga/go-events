# go-events

```go
package main

type Hello struct {
	Var string
}

func (*Hello) EventName() string {
	return "hello" // should be unique globally
}
```


```go
package main

import (
	"fmt"
	"log"

	"github.com/ituoga/go-events"
)

func init() {
    events.Register[*Hello]()
}

func main() {

    // Before runs before each Publish
	events.Before(func(name string, event []byte, e any) error {
		// Store to database
		log.Printf("%s: %s %T\n", name, event, e)

        // return error if database failes so no futher actions 
        // will be made
        // or return nil if success
		return nil //errors.New("error")
	})
   

    // Subscribe to Hello.EventName() 
	events.Subscribe(func(event *Hello) {
		fmt.Printf("%v %T\n", event, event)
	})

    // Publish Hello struct and print error if any
	log.Printf("%v", events.Publish(&Hello{"Var ;)"}))
}
```