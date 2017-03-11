package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPluginWeb(t *testing.T) {
	p := NewPluginWeb()
	assert.NotNil(t, p.storage)
}

func TestPluginWeb_GetPluginInfo(t *testing.T) {
	p := NewPluginWeb()
	assert.Equal(t, "web", p.GetPluginInfo().Tag)
	assert.Equal(t, "Web", p.GetPluginInfo().Name)
}

func TestPluginWeb_Subscribe(t *testing.T) {
	p := NewPluginWeb()
	inJson := []byte("{\"event_name\": \"test\", \"url\": \"hello\"}")
	err := p.Subscribe(inJson)
	assert.NoError(t, err)
	assert.Len(t, p.storage.Data, 1)
	assert.Len(t, p.storage.Data["test"], 1)
	assert.Equal(t, "hello", p.storage.Data["test"][0])
}

func TestPluginWeb_Unsubscribe(t *testing.T) {
	p := NewPluginWeb()
	inJson := []byte("{\"event_name\": \"test\", \"url\": \"hello\"}")
	err := p.Subscribe(inJson)
	assert.NoError(t, err)
	assert.Len(t, p.storage.Data, 1)
	err = p.Unsubscribe(inJson)
	assert.NoError(t, err)
	assert.Len(t, p.storage.Data, 0)
}
