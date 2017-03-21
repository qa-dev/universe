package plugins

import (
	"github.com/qa-dev/universe/event"
)

type Plugin interface {
	GetPluginInfo() *PluginInfo
	Subscribe(input []byte) error
	Unsubscribe(input []byte) error
	ProcessEvent(eventData *event.Event)
	Loaded()
}

type PluginInfo struct {
	Name string
	Tag  string
}
