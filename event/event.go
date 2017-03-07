package event

type Event struct {
	Name    string `json:"name"`
	Payload []byte `json:"payload"`
}
