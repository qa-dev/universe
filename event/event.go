package event

import "fmt"

type Event struct {
	Name    string `json:"name"`
	Payload []byte `json:"payload"`
}

func (e Event) String() string {
	return fmt.Sprintf("(\n\tName: %s\n\tPayload: %s\n)", e.Name, string(e.Payload))
}
