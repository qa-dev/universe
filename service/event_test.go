package service

import (
	"fmt"
	"testing"

	"github.com/qa-dev/universe/data"
	"github.com/stretchr/testify/assert"
)

func TestNewEventService(t *testing.T) {
	channel := make(chan data.Event)
	es := NewEventService(channel)
	assert.Equal(t, fmt.Sprintf("%p", channel), fmt.Sprintf("%p", es.ch))
}

func TestEventService_PushEvent(t *testing.T) {
	channel := make(chan data.Event)
	es := NewEventService(channel)

	go func() {
		event := <-channel
		assert.Equal(t, "test.event", event.Name, "Wrong event name generated")
	}()

	err := es.PushEvent("test.event", []byte("test"))
	assert.NoError(t, err)
}

func TestEventService_PushEvent_Blank(t *testing.T) {
	channel := make(chan data.Event)
	es := NewEventService(channel)

	err := es.PushEvent("", []byte("test"))

	assert.Error(t, err)
	assert.Equal(t, "BLANK EVENT NAME", err.Error())
}
