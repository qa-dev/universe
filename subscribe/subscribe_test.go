package subscribe

import (
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/storage"
	"github.com/stretchr/testify/assert"
	"os"
)

var amqpUri string

func init() {
	amqpUri = os.Getenv("AMQP_URI")
	if amqpUri == "" {
		log.Fatal("AMQP_URI is required to run rabbitmq tests")
	}
}

func TestNewSubscribeService(t *testing.T) {
	storageUnit := storage.NewStorage()
	subscribeService := NewSubscribeService(storageUnit)
	assert.Equal(t, fmt.Sprintf("%p", storageUnit), fmt.Sprintf("%p", subscribeService.storage))
}

func TestSubscribeService_ProcessSubscribe(t *testing.T) {
	storageUnit := storage.NewStorage()
	subscribeService := NewSubscribeService(storageUnit)
	subscribe := Subscribe{EventName: "test.event", WebHookPath: "testpath"}
	assert.Equal(t, 0, len(subscribeService.storage.Data))
	err := subscribeService.ProcessSubscribe(subscribe)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(subscribeService.storage.Data))
}

func TestSubscribeService_ProcessSubscribe_BlankEventName(t *testing.T) {
	storageUnit := storage.NewStorage()
	subscribeService := NewSubscribeService(storageUnit)
	subscribe := Subscribe{EventName: "", WebHookPath: "testpath"}
	err := subscribeService.ProcessSubscribe(subscribe)
	assert.Error(t, err)
	assert.Equal(t, "BLANK EVENT NAME", err.Error())
}

func TestSubscribeService_ProcessSubscribe_BlankWebHook(t *testing.T) {
	storageUnit := storage.NewStorage()
	subscribeService := NewSubscribeService(storageUnit)
	subscribe := Subscribe{EventName: "test.event", WebHookPath: ""}
	err := subscribeService.ProcessSubscribe(subscribe)
	assert.Error(t, err)
	assert.Equal(t, "BLANK WEBHOOK PATH", err.Error())
}
