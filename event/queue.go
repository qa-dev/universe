package event

type Queue chan *Event

func NewEventQueue() *Queue {
	q := make(Queue)
	return &q
}
