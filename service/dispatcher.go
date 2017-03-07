package service

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"net/http"

	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/storage"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatcher struct {
	ch         chan event.Event
	storage    *storage.Storage
	httpClient ClientInterface
}

func NewDispatcher(ch chan event.Event, storage *storage.Storage, client ClientInterface) *Dispatcher {
	return &Dispatcher{ch, storage, client}
}

func (d *Dispatcher) Run() {
	go d.worker()
}

func (d *Dispatcher) worker() {
	for {
		e := <-d.ch

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
	}
}
