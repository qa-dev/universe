package web

import (
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
	"github.com/qa-dev/universe/subscribe"
)

type PluginWeb struct {
}

func (p *PluginWeb) GetPluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{Name: "Web", Tag: "WEB"}
}

func (p *PluginWeb) Subscribe(eventName string, subscribeData subscribe.SubscribeData) {

}

func (p *PluginWeb) Unsubscribe(eventName string, unsubscribeData subscribe.UnsubscribeData) {

}

func (o *PluginWeb) ProcessEvent(eventName string, eventData event.Event) {

}
