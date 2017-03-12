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

var Obs *Observable

type Observable struct {
	plugins []Plugin
}

type PluginInfo struct {
	Name string
	Tag  string
}

func init() {
	Obs = NewObservable()
}

func NewObservable() *Observable {
	return &Observable{}
}
