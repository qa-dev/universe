package service

import (
	"errors"
	"log"

	"github.com/qa-dev/universe/data"
)

type EventService struct {
	ch chan data.Event
}

func NewEventService(ch chan data.Event) *EventService {
	return &EventService{ch}
}

func (e *EventService) PushEvent(eventName string, payload []byte) error {
	if eventName == "" {
		log.Println("Got blank event name")
		return errors.New("BLANK EVENT NAME")
	}
	log.Println("Got event name", eventName)
	var event data.Event
	event.Name = eventName
	event.Payload = payload
	e.ch <- event
	return nil
}
