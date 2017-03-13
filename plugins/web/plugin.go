package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/universe/event"
	"github.com/qa-dev/universe/plugins"
)

type PluginWeb struct {
	storage *Storage
	client  HttpRequester
}

type HttpRequester interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewPluginWeb() PluginWeb {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	storage := NewStorage()
	return PluginWeb{storage, httpClient}
}

func (p PluginWeb) GetPluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{Name: "Web", Tag: "web"}
}

func (p PluginWeb) Subscribe(input []byte) error {
	var subscribeData SubscribeData
	err := json.Unmarshal(input, &subscribeData)
	if err != nil || len(subscribeData.EventName) == 0 || len(subscribeData.Url) == 0 {
		log.Errorf("%+v", input)
		errorText := fmt.Sprintf("Invalid input data: %s", string(input))
		return errors.New(errorText)
	}
	p.storage.AddSubscriber(subscribeData)
	return nil
}

func (p PluginWeb) Unsubscribe(input []byte) error {
	var unsubscribeData UnsubscribeData
	err := json.Unmarshal(input, &unsubscribeData)
	if err != nil || len(unsubscribeData.EventName) == 0 || len(unsubscribeData.Url) == 0 {
		log.Errorf("%+v", input)
		errorText := fmt.Sprintf("Invalid input data: %s", string(input))
		return errors.New(errorText)
	}
	err = p.storage.RemoveSubscriber(unsubscribeData)

	return err
}

func (p PluginWeb) ProcessEvent(eventData event.Event) {
	for _, subscribeUrl := range p.storage.Data[eventData.Name] {
		go p.sendRequest(subscribeUrl, eventData.Payload)
	}
}

func (p PluginWeb) sendRequest(url string, payload []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO: log statistics
	p.client.Do(req)
}
