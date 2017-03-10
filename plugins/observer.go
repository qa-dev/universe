package plugins

import (
	"github.com/qa-dev/universe/event"
)

func (o *Observable) Register(v Plugin) {
	o.plugins = append(o.plugins, v)
}

func (o *Observable) ProcessEvent(name string, eventData event.Event) {
	for _, ob := range o.plugins {
		go func(o Plugin) {
			o.ProcessEvent(name, eventData)
		}(ob)
	}
}
