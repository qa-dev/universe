package plugins

import (
	"errors"

	"github.com/qa-dev/universe/event"
)

func (o *Observable) Register(v Plugin) {
	o.plugins = append(o.plugins, v)
}

func (o *Observable) ProcessEvent(eventData event.Event) {
	for _, ob := range o.plugins {
		go func(o Plugin) {
			o.ProcessEvent(eventData)
		}(ob)
	}
}

func (o *Observable) ProcessSubscribe(pluginName string, input []byte) error {
	for _, ob := range o.plugins {
		if ob.GetPluginInfo().Tag == pluginName {
			return ob.Subscribe(input)
		}
	}

	return errors.New("No plugin found")
}

func (o *Observable) GetPlugins() []Plugin {
	return o.plugins
}
