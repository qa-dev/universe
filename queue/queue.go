package queue

import (
	"github.com/qa-dev/universe/rabbitmq"
)

type Queue struct {
	rabbit *rabbitmq.RabbitMQ
}

func NewQueue(rabbit *rabbitmq.RabbitMQ) *Queue {
	q := Queue{rabbit: rabbit}
	return &q
}

func (q *Queue) SendEvent(event interface{}) {
	q.rabbit.SendEvent(event)
}

func (q *Queue) GetConsumer(name string) (<-chan rabbitmq.Delivery, error) {
	return q.rabbit.GetConsumer(name)
}

func (q *Queue) IsOnline() bool {
	return q.rabbit.IsOnline()
}

func (q *Queue) NotifyReconnect(receiver chan bool) {
	q.rabbit.NotifyReconnect(receiver)
}
