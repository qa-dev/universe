package event

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/rabbitmq"
	"github.com/stretchr/testify/assert"
)

var amqpUri string

func init() {
	amqpUri = os.Getenv("AMQP_URI")
	if amqpUri == "" {
		log.Fatal("AMQP_URI is required to run rabbitmq tests")
	}
}

func TestNewEventService(t *testing.T) {
	rmq := &rabbitmq.RabbitMQ{}
	es := NewEventService(rmq)
	assert.Equal(t, fmt.Sprintf("%p", rmq), fmt.Sprintf("%p", es.rmq))
}

func TestEventService_PushEvent(t *testing.T) {
	rmq := rabbitmq.NewRabbitMQ(amqpUri, "test_event_service_push_event_queue")
	defer rmq.Close()
	// Даем время на подключение
	time.Sleep(5 * time.Second)
	es := NewEventService(rmq)

	consumeObj, err := rmq.Consume("test_consumer")
	assert.NoError(t, err)

	go func() {
		raw := <-consumeObj
		var e Event
		err = json.Unmarshal(raw.Body(), &e)
		assert.NoError(t, err)
		assert.Equal(t, "test.event", e.Name, "Wrong event name generated")
	}()

	err = es.Publish(Event{"test.event", []byte("test")})
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
}

func TestEventService_PushEvent_Blank(t *testing.T) {
	rmq := rabbitmq.NewRabbitMQ(amqpUri, "test_event_service_push_event_queue")
	defer rmq.Close()
	// Даем время на подключение
	time.Sleep(5 * time.Second)
	es := NewEventService(rmq)

	err := es.Publish(Event{"", []byte("test")})

	assert.Error(t, err)
	assert.Equal(t, "BLANK EVENT NAME", err.Error())
}
