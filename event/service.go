package event

import (
	"errors"
	"log"
)

type EventService struct {
	ch chan Event
}

func NewEventService(ch chan Event) *EventService {
	return &EventService{ch}
}

func (e *EventService) Publish(ev Event) error {
	if ev.Name == "" {
		log.Println("Got blank event name")
		return errors.New("BLANK EVENT NAME")
	}
	log.Println("Got event name", ev.Name)
	e.ch <- ev
	return nil
}
