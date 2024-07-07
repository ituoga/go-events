package events

import "github.com/google/uuid"

func UUID() *string {
	id, _ := uuid.NewV7()
	rid := id.String()
	return &rid
}
