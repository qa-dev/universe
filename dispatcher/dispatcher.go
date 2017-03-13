package dispatcher

import (
	"encoding/json"
	"net/http"

	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
	"github.com/qa-dev/universe/rabbitmq"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatcher struct {
	rmq           *rabbitmq.RabbitMQ
	pluginStorage *plugins.PluginStorage
}

func NewDispatcher(rmq *rabbitmq.RabbitMQ, storage *plugins.PluginStorage) *Dispatcher {
	return &Dispatcher{rmq, storage}
}

func (d *Dispatcher) Run() {
	go d.worker()
}

func (d *Dispatcher) worker() {
	consumeObj, err := d.rmq.Consume("consumer")
	if err != nil {
		panic(err)
	}

	for {
		rawData := <-consumeObj
		var e event.Event
		err = json.Unmarshal(rawData.Body(), &e)
		if err != nil {
			// TODO log?
			panic(err)
		}

		d.pluginStorage.ProcessEvent(e)

		rawData.Ack(false)
	}
}
