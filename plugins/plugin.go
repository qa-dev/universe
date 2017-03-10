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

type PluginInfo struct {
	Name string
	Tag  string
}
