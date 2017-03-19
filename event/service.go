package event

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/queue"
)

type EventService struct {
	queue *queue.Queue
}

func NewEventService(queue *queue.Queue) *EventService {
	return &EventService{queue}
}

func (e *EventService) Publish(ev *Event) error {
	if ev.Name == "" {
		log.Println("Got blank event name")
		return errors.New("BLANK EVENT NAME")
	}
	log.Println("Got event name", ev.Name)
	e.queue.SendEvent(ev)
	return nil
}
