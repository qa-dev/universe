package service

import (
	"bytes"
	"net/http"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/rabbitmq"
	"github.com/qa-dev/universe/storage"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatcher struct {
	rmq        *rabbitmq.RabbitMQ
	storage    *storage.Storage
	httpClient ClientInterface
}

func NewDispatcher(rmq *rabbitmq.RabbitMQ, storage *storage.Storage, client ClientInterface) *Dispatcher {
	return &Dispatcher{rmq, storage, client}
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
			panic(err)
		}

		d.storage.Mutex.Lock()

		if val, ok := d.storage.Data[e.Name]; ok {
			for _, hookPath := range val {
				log.Println("Sending event", e.Name, "to", hookPath)
				req, err := http.NewRequest("POST", hookPath, bytes.NewBuffer(e.Payload))
				req.Header.Set("Content-Type", "application/json")

				resp, err := d.httpClient.Do(req)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				log.Println("Status of sending event", e.Name, "is", resp.Status)
				resp.Body.Close()
			}
		} else {
			log.Println("No subscribers for event", e.Name)
		}

		d.storage.Mutex.Unlock()

		rawData.Ack(false)
	}
}
