package event

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/rabbitmq"
)

type EventService struct {
	rmq *rabbitmq.RabbitMQ
}

func NewEventService(rmq *rabbitmq.RabbitMQ) *EventService {
	return &EventService{rmq}
}

func (e *EventService) Publish(ev Event) error {
	if ev.Name == "" {
		log.Println("Got blank event name")
		return errors.New("BLANK EVENT NAME")
	}
	log.Println("Got event name", ev.Name)
	e.rmq.PublishWithPriority(ev, 1)
	return nil
}
