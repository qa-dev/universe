package rabbitmq

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Реализуем интерфейс amqp.Acknowledger
type MockAcknowledger struct {
	mock.Mock
}

func (m *MockAcknowledger) Ack(tag uint64, multiple bool) error {
	args := m.Called(tag, multiple)
	return args.Error(0)
}
func (m *MockAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	args := m.Called(tag, multiple, requeue)
	return args.Error(0)
}
func (m *MockAcknowledger) Reject(tag uint64, requeue bool) error {
	args := m.Called(tag, requeue)
	return args.Error(0)
}

func TestDelivery(t *testing.T) {
	a := assert.New(t)

	tag := uint64(123)
	ma := new(MockAcknowledger)
	ma.On("Ack", tag, false).Return(nil)

	d := amqp.Delivery{
		Acknowledger: ma,
		Body:         []byte("Test body"),
		DeliveryTag:  tag,
		Priority:     10,
	}

	amqpCh := make(chan amqp.Delivery)

	deliveryCh := newDelivery(amqpCh)

	amqpCh <- d

	msg := <-deliveryCh
	a.Equal(d.Body, msg.Body(), "Unexpected body of recieved msg")
	a.Equal(d.Priority, msg.Priority(), "Unexpected priority of recieved msg")
	a.NoError(msg.Ack())
}
