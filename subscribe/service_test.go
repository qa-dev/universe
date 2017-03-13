package subscribe

import (
	"testing"

	"github.com/qa-dev/universe/plugins"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeService_ProcessSubscribe_BlankPluginName(t *testing.T) {
	storage := plugins.NewPluginStorage()
	subSer := NewSubscribeService(storage)
	err := subSer.ProcessSubscribe("", []byte(""))
	assert.Error(t, err)
}
