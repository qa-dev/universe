package subscribe

import (
	"fmt"
	"testing"

	"github.com/qa-dev/universe/data"
	"github.com/stretchr/testify/assert"
)

func TestNewSubscribeService(t *testing.T) {
	storage := data.NewStorage()
	subscribeService := NewSubscribeService(storage)
	assert.Equal(t, fmt.Sprintf("%p", storage), fmt.Sprintf("%p", subscribeService.storage))
}

func TestSubscribeService_ProcessSubscribe(t *testing.T) {
	storage := data.NewStorage()
	subscribeService := NewSubscribeService(storage)
	subscribe := Subscribe{EventName: "test.event", WebHookPath: "testpath"}
	assert.Equal(t, 0, len(subscribeService.storage.Data))
	err := subscribeService.ProcessSubscribe(subscribe)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(subscribeService.storage.Data))
}

func TestSubscribeService_ProcessSubscribe_BlankEventName(t *testing.T) {
	storage := data.NewStorage()
	subscribeService := NewSubscribeService(storage)
	subscribe := Subscribe{EventName: "", WebHookPath: "testpath"}
	err := subscribeService.ProcessSubscribe(subscribe)
	assert.Error(t, err)
	assert.Equal(t, "BLANK EVENT NAME", err.Error())
}

func TestSubscribeService_ProcessSubscribe_BlankWebHook(t *testing.T) {
	storage := data.NewStorage()
	subscribeService := NewSubscribeService(storage)
	subscribe := Subscribe{EventName: "test.event", WebHookPath: ""}
	err := subscribeService.ProcessSubscribe(subscribe)
	assert.Error(t, err)
	assert.Equal(t, "BLANK WEBHOOK PATH", err.Error())
}
