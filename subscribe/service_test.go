package subscribe

import (
	"testing"

	"github.com/qa-dev/universe/plugins"
	"github.com/qa-dev/universe/plugins/log"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeService_ProcessSubscribe_BlankPluginName(t *testing.T) {
	storage := plugins.NewPluginStorage()
	subSer := NewSubscribeService(storage)
	err := subSer.ProcessSubscribe("", []byte(""))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BLANK PLUGIN NAME")
}

func TestSubscribeService_ProcessSubscribe_WrongPluginName(t *testing.T) {
	storage := plugins.NewPluginStorage()
	subSer := NewSubscribeService(storage)
	err := subSer.ProcessSubscribe("pew", []byte(""))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No plugin found")
}

func TestSubscribeService_ProcessSubscribe(t *testing.T) {
	storage := plugins.NewPluginStorage()
	storage.Register(log.NewLog())
	subSer := NewSubscribeService(storage)
	err := subSer.ProcessSubscribe("log", []byte(""))
	assert.NoError(t, err)
}
