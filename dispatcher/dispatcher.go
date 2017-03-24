package dispatcher

import (
	"encoding/json"
	"net/http"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
	"github.com/qa-dev/universe/rabbitmq"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatcher struct {
	queue         *rabbitmq.RabbitMQ
	pluginStorage *plugins.PluginStorage
}

func NewDispatcher(queue *rabbitmq.RabbitMQ, storage *plugins.PluginStorage) *Dispatcher {
	return &Dispatcher{queue, storage}
}

func (d *Dispatcher) Run() {
	for cnt := 1; cnt <= runtime.NumCPU(); cnt++ {
		log.Infof("Run worker %d", cnt)
		go d.worker(cnt)
	}
}

func (d *Dispatcher) worker(num int) {
	msgs, err := d.queue.GetConsumer("consumer" + string(num))
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
