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
	consumerName := "consumer" + string(num)
	msgs, err := d.queue.GetConsumer(consumerName)
	if err != nil {
		log.Errorf("Error get consumer in event dispatcher worker %d", num)
	}
	for {
		if d.queue.IsOnline() == false {
			log.Infof("Worker %d lost connection to queue. Waiting...", num)
			backOnlineChan := make(chan bool)
			d.queue.NotifyReconnect(backOnlineChan)
			_ = <-backOnlineChan
			newMsgs, err := d.queue.GetConsumer(consumerName)
			if err != nil {
				log.Errorf("Error get consumer in event dispatcher worker %d", num)
			}
			msgs = newMsgs
			log.Infof("Worker %d established connection to queue", num)
		}
		data := <-msgs
		var ev event.Event
		err = json.Unmarshal(data.Body(), &ev)
		if err != nil {
			data.Reject()
			log.Errorf("Error unmarshal event in event dispatcher worker %d", num)
			continue
		}
		d.pluginStorage.ProcessEvent(&ev)
		data.Ack()
	}
}
