package dispatcher

import (
	"net/http"

	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatcher struct {
	queue         *event.Queue
	pluginStorage *plugins.PluginStorage
}

func NewDispatcher(queue *event.Queue, storage *plugins.PluginStorage) *Dispatcher {
	return &Dispatcher{queue, storage}
}

func (d *Dispatcher) Run() {
	go d.worker()
}

func (d *Dispatcher) worker() {
	for {
		e := <-*d.queue
		d.pluginStorage.ProcessEvent(e)
	}
}
