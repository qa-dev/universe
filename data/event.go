package data

type Event struct {
	Name    string `json:"name"`
	Payload []byte `json:"payload"`
}
