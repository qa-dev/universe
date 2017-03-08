package event

import (
	"fmt"
	"testing"

	"github.com/qa-dev/universe/config"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.SetTestDitectory()
}

func TestNewEventService(t *testing.T) {
	channel := make(chan Event)
	es := NewEventService(channel)
	assert.Equal(t, fmt.Sprintf("%p", channel), fmt.Sprintf("%p", es.ch))
}

func TestEventService_PushEvent(t *testing.T) {
	channel := make(chan Event)
	es := NewEventService(channel)

	go func() {
		e := <-channel
		assert.Equal(t, "test.event", e.Name, "Wrong event name generated")
	}()

	err := es.Publish(Event{"test.event", []byte("test")})
	assert.NoError(t, err)
}

func TestEventService_PushEvent_Blank(t *testing.T) {
	channel := make(chan Event)
	es := NewEventService(channel)

	err := es.Publish(Event{"", []byte("test")})

	assert.Error(t, err)
	assert.Equal(t, "BLANK EVENT NAME", err.Error())
}
