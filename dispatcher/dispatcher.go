package dispatcher

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
	"github.com/qa-dev/universe/queue"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatcher struct {
	queue         *queue.Queue
	pluginStorage *plugins.PluginStorage
}

func NewDispatcher(queue *queue.Queue, storage *plugins.PluginStorage) *Dispatcher {
	return &Dispatcher{queue, storage}
}

func (d *Dispatcher) Run() {
	go d.worker()
}

func (d *Dispatcher) worker() {
	msgs, err := d.queue.GetConsumer("consumer")
	if err != nil {
		log.Error("Error get consumer in event dispatcher worker")
	}
	for {
		if d.queue.IsOnline() == false {
			log.Info("Worker lost connection to queue. Waiting...")
			backOnlineChan := make(chan bool)
			d.queue.NotifyReconnect(backOnlineChan)
			_ = <-backOnlineChan
			newMsgs, err := d.queue.GetConsumer("consumer")
			if err != nil {
				log.Error("Error get consumer in event dispatcher worker")
			}
			msgs = newMsgs
			log.Info("Worker established connection to queue")
		}
		data := <-msgs
		var ev event.Event
		err = json.Unmarshal(data.Body(), &ev)
		if err != nil {
			data.Reject()
			log.Error("Error unmarchal event in event dispatcher worker")
			continue
		}
		d.pluginStorage.ProcessEvent(&ev)
		data.Ack()
	}
}
