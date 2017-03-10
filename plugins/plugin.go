package plugins

import (
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/subscribe"
)

type Plugin interface {
	GetPluginInfo() *PluginInfo
	Subscribe(eventName string, subscribeData subscribe.SubscribeData)
	Unsubscribe(eventName string, unsubscribeData subscribe.UnsubscribeData)
	ProcessEvent(eventName string, eventData event.Event)
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
