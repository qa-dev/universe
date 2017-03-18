package event

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventService(t *testing.T) {
	queue := NewEventQueue()
	es := NewEventService(queue)
	assert.Equal(t, fmt.Sprintf("%p", queue), fmt.Sprintf("%p", es.queue))
}

func TestEventService_PushEvent(t *testing.T) {
	queue := NewEventQueue()
	es := NewEventService(queue)

	go func() {
		e := <-*queue
		assert.Equal(t, "test.event", e.Name, "Wrong event name generated")
	}()

	err := es.Publish(Event{"test.event", []byte("test")})
	assert.NoError(t, err)
}

func TestEventService_PushEvent_Blank(t *testing.T) {
	queue := NewEventQueue()
	es := NewEventService(queue)

	err := es.Publish(Event{"", []byte("test")})

	assert.Error(t, err)
	assert.Equal(t, "BLANK EVENT NAME", err.Error())
}
