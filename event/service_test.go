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
	q := rabbitmq.NewRabbitMQ(amqpUri, "test_new_dispatcher")
	time.Sleep(2 * time.Second)
	es := NewEventService(q)
	assert.Equal(t, fmt.Sprintf("%p", q), fmt.Sprintf("%p", es.queue))
}

func TestEventService_PushEvent(t *testing.T) {
	q := rabbitmq.NewRabbitMQ(amqpUri, "test_new_dispatcher")
	time.Sleep(2 * time.Second)
	es := NewEventService(q)

	go func() {
		msgs, err := q.GetConsumer("test_consumer")
		assert.NoError(t, err)
		data := <-msgs
		var e Event
		err = json.Unmarshal(data.Body(), &e)
		data.Ack()
		assert.NoError(t, err)
		assert.Equal(t, "test.event", e.Name, "Wrong event name generated")
	}()

	err := es.Publish(&Event{"test.event", []byte("test")})
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
}

func TestEventService_PushEvent_Blank(t *testing.T) {
	q := rabbitmq.NewRabbitMQ(amqpUri, "test_new_dispatcher")
	time.Sleep(2 * time.Second)
	es := NewEventService(q)

	err := es.Publish(&Event{"", []byte("test")})

	assert.Error(t, err)
	assert.Equal(t, "BLANK EVENT NAME", err.Error())
}
