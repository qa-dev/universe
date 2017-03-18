package event

import (
	"errors"

	log "github.com/Sirupsen/logrus"
)

type EventService struct {
	queue *Queue
}

func NewEventService(queue *Queue) *EventService {
	return &EventService{queue}
}

func (e *EventService) Publish(ev Event) error {
	if ev.Name == "" {
		log.Println("Got blank event name")
		return errors.New("BLANK EVENT NAME")
	}
	log.Println("Got event name", ev.Name)
	*e.queue <- &ev
	return nil
}
