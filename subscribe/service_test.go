package subscribe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscribeService_ProcessSubscribe_BlankPluginName(t *testing.T) {
	subSer := NewSubscribeService()
	err := subSer.ProcessSubscribe("", []byte(""))
	assert.Error(t, err)
}
