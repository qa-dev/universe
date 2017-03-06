package service

import (
	"bytes"
	"log"
	"net/http"

	"github.com/qa-dev/Universe/data"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatcher struct {
	ch         chan data.Event
	storage    *data.Storage
	httpClient ClientInterface
}

func NewDispatcher(ch chan data.Event, storage *data.Storage, client ClientInterface) *Dispatcher {
	return &Dispatcher{ch, storage, client}
}

func (d *Dispatcher) Run() {
	go d.worker()
}

func (d *Dispatcher) worker() {
	for {
		event := <-d.ch

		d.storage.Mutex.Lock()

		if val, ok := d.storage.Data[event.Name]; ok {
			for _, hookPath := range val {
				log.Println("Sending event", event.Name, "to", hookPath)
				req, err := http.NewRequest("POST", hookPath, bytes.NewBuffer(event.Payload))
				req.Header.Set("Content-Type", "application/json")

				resp, err := d.httpClient.Do(req)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				log.Println("Status of sending event", event.Name, "is", resp.Status)
				resp.Body.Close()
			}
		} else {
			log.Println("No subscribers for event", event.Name)
		}

		d.storage.Mutex.Unlock()
	}
}
