package plugins

import (
	"github.com/qa-dev/universe/event"
)

type Plugin interface {
	GetPluginInfo() *PluginInfo
	Subscribe(input []byte) error
	Unsubscribe(input []byte) error
	ProcessEvent(eventData event.Event)
}

type PluginStorage struct {
	plugins []Plugin
}

func NewPluginStorage() *PluginStorage {
	return &PluginStorage{}
}

type PluginInfo struct {
	Name string
	Tag  string
}
