package plugins

import (
	"errors"

	"github.com/qa-dev/universe/event"
)


type PluginStorage struct {
	plugins []Plugin
}

func NewPluginStorage() *PluginStorage {
	return &PluginStorage{}
}

func (o *PluginStorage) Register(v Plugin) {
	o.plugins = append(o.plugins, v)
}

func (o *PluginStorage) ProcessEvent(eventData *event.Event) {
	for _, ob := range o.plugins {
		ob.ProcessEvent(eventData)
	}
}

func (o *PluginStorage) ProcessSubscribe(pluginName string, input []byte) error {
	for _, ob := range o.plugins {
		if ob.GetPluginInfo().Tag == pluginName {
			return ob.Subscribe(input)
		}
	}

	return errors.New("No plugin found")
}

func (o *PluginStorage) ProcessUnsubscribe(pluginName string, input []byte) error {
	for _, ob := range o.plugins {
		if ob.GetPluginInfo().Tag == pluginName {
			return ob.Unsubscribe(input)
		}
	}

	return errors.New("No plugin found")
}

func (o *PluginStorage) GetPlugins() []Plugin {
	return o.plugins
}
